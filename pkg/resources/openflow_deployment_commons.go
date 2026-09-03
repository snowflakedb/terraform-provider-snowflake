package resources

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// The two deployment types have different lifecycles, so they are modeled as two resources over the single
// OPENFLOW DEPLOYMENT object, as the catalog_integration_* resources are. Snowflake-managed settles on
// ACTIVE; BYOC settles on INACTIVE and only becomes ACTIVE once the customer provisions their own cloud
// infrastructure, which cannot be driven from SQL.
//
// Shared here: the common attributes, the read path, the update path and the two-step teardown.

const (
	// openflowDeploymentCreateTimeout covers a Snowflake-managed create, observed at just over six minutes.
	openflowDeploymentCreateTimeout = 30 * time.Minute
	openflowDeploymentUpdateTimeout = 30 * time.Minute
	// openflowDeploymentDeleteTimeout covers TERMINATE tearing infrastructure down, then DROP.
	openflowDeploymentDeleteTimeout = 30 * time.Minute

	// openflowDeploymentEventTableAttribute is not part of SHOW or DESCRIBE output; it is a
	// deployment-level parameter read back through SHOW PARAMETERS.
	openflowDeploymentEventTableAttribute = "event_table"
	openflowDeploymentEventTableParameter = "EVENT_TABLE"
)

var openflowDeploymentTimeouts = &schema.ResourceTimeout{
	Create: schema.DefaultTimeout(openflowDeploymentCreateTimeout),
	Read:   schema.DefaultTimeout(defaultReadTimeout),
	Update: schema.DefaultTimeout(openflowDeploymentUpdateTimeout),
	Delete: schema.DefaultTimeout(openflowDeploymentDeleteTimeout),
}

// openflowDeploymentCommonSchema returns the attributes shared by both deployment resources. Callers add
// their own type-specific attributes on top.
func openflowDeploymentCommonSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:             schema.TypeString,
			Required:         true,
			Description:      blocklistedCharactersFieldDescription("Specifies the identifier for the Openflow deployment; must be unique for the account."),
			DiffSuppressFunc: suppressIdentifierQuoting,
		},
		"display_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "A free-text alias for the deployment. Shown in the Openflow UI in place of the deployment's identifier when set.",
		},
		openflowDeploymentEventTableAttribute: {
			Type:             schema.TypeString,
			Optional:         true,
			Computed:         true,
			Description:      "Fully qualified name of an event table the deployment logs to. For more information, check [EVENT_TABLE documentation](https://docs.snowflake.com/en/sql-reference/parameters#event-table).",
			DiffSuppressFunc: suppressIdentifierQuoting,
		},
		"comment": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Specifies a comment for the Openflow deployment.",
		},
		"type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Specifies the type of the Openflow deployment. This field is used to detect when the deployment type was changed outside of Terraform and to recreate the resource when that happens.",
		},
		FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
		ShowOutputAttributeName: {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Outputs the result of `SHOW OPENFLOW DEPLOYMENTS` for the given deployment.",
			Elem: &schema.Resource{
				Schema: schemas.ShowOpenflowDeploymentSchema,
			},
		},
		ParametersAttributeName: {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Outputs the result of `SHOW PARAMETERS IN OPENFLOW DEPLOYMENT` for the given deployment.",
			Elem:        &schema.Resource{Schema: schemas.ShowOpenflowDeploymentParametersSchema},
		},
		DescribeOutputAttributeName: {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Outputs the result of `DESCRIBE OPENFLOW DEPLOYMENT` for the given deployment.",
			Elem: &schema.Resource{
				Schema: schemas.DescribeOpenflowDeploymentSchema,
			},
		},
	}
}

// createOpenflowDeploymentCommonFields applies the shared attributes to a create request.
func createOpenflowDeploymentCommonFields(d *schema.ResourceData, request *sdk.CreateOpenflowDeploymentRequest) error {
	// Only the table form of EVENT_TABLE is exposed. NONE means "drop all events", a different intent from
	// omitting the attribute ("inherit the account's"), so it would need its own representation in config.
	if err := errors.Join(
		stringAttributeCreateBuilder(d, "display_name", request.WithDisplayName),
		stringAttributeCreateBuilder(d, "comment", request.WithComment),
	); err != nil {
		return err
	}
	if diags := handleOpenflowDeploymentParameterCreate(d, request); diags.HasError() {
		return fmt.Errorf("%s", diags[0].Summary)
	}
	return nil
}

// importOpenflowDeploymentOptionals sets the optional attributes only when Snowflake reports a value, so an
// unset one stays absent from state as it does after a create.
func importOpenflowDeploymentOptionals(d *schema.ResourceData, deployment *sdk.OpenflowDeployment) error {
	var errs []error
	if deployment.DisplayName != nil && *deployment.DisplayName != "" {
		errs = append(errs, d.Set("display_name", *deployment.DisplayName))
	}
	if deployment.Comment != nil && *deployment.Comment != "" {
		errs = append(errs, d.Set("comment", *deployment.Comment))
	}
	return errors.Join(errs...)
}

// readOpenflowDeploymentFunc builds the read handler shared by both deployment resources.
// resourceSchema is the concrete resource's schema, and typeSpecificFields lists the attributes that
// resource owns beyond the common ones, so drift handling covers them too.
func readOpenflowDeploymentFunc(
	resourceSchema map[string]*schema.Schema,
	withExternalChangesMarking bool,
	setTypeSpecificFields func(d *schema.ResourceData, deployment *sdk.OpenflowDeployment) error,
) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		client := meta.(*provider.Context).Client
		id, err := sdk.ParseAccountObjectIdentifier(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}

		deployment, err := client.OpenflowDeployments.ShowByIDSafely(ctx, id)
		if err != nil {
			if errors.Is(err, sdk.ErrObjectNotFound) {
				d.SetId("")
				return diag.Diagnostics{
					diag.Diagnostic{
						Severity: diag.Warning,
						Summary:  "Failed to query Openflow deployment. Marking the resource as removed.",
						Detail:   fmt.Sprintf("Openflow deployment id: %s, Err: %s", id.FullyQualifiedName(), err),
					},
				}
			}
			return diag.FromErr(err)
		}

		deploymentDetails, err := client.OpenflowDeployments.Describe(ctx, id)
		if err != nil {
			return diag.FromErr(err)
		}

		parameters, err := client.OpenflowDeployments.ShowParameters(ctx, id)
		if err != nil {
			return diag.FromErr(fmt.Errorf("reading parameters of openflow deployment %s: %w", id.Name(), err))
		}

		// Flattened up front: the drift comparison is an `any` equality against the string in show_output, so a
		// *string would never match and would mark drift on every read.
		var displayName, comment string
		if deployment.DisplayName != nil {
			displayName = *deployment.DisplayName
		}
		if deployment.Comment != nil {
			comment = *deployment.Comment
		}

		if withExternalChangesMarking {
			if err = handleExternalChangesToObjectInShow(
				d,
				outputMapping{"display_name", "display_name", displayName, displayName, nil},
				outputMapping{"comment", "comment", comment, comment, nil},
			); err != nil {
				return diag.FromErr(err)
			}
		}

		if err = handleOpenflowDeploymentParameterRead(d, parameters); err != nil {
			return diag.FromErr(err)
		}

		if err = setStateToValuesFromConfig(d, resourceSchema, []string{
			"display_name",
			"comment",
		}); err != nil {
			return diag.FromErr(err)
		}

		errs := errors.Join(
			d.Set(ShowOutputAttributeName, []map[string]any{schemas.OpenflowDeploymentToSchema(deployment)}),
			d.Set(DescribeOutputAttributeName, []map[string]any{schemas.OpenflowDeploymentDetailsToSchema(*deploymentDetails)}),
			d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
			d.Set("type", string(deployment.Type)),
			d.Set(ParametersAttributeName, []map[string]any{schemas.OpenflowDeploymentParametersToSchema(parameters, meta.(*provider.Context))}),
			setTypeSpecificFields(d, deployment),
		)
		if errs != nil {
			return diag.FromErr(errs)
		}
		return nil
	}
}

// updateOpenflowDeploymentCommonFields applies changes to the shared attributes. All are metadata only, so
// none needs polling. A rename replaces the resource id, so a caller needing the identifier afterwards must
// re-read it from d.
func updateOpenflowDeploymentCommonFields(ctx context.Context, client *sdk.Client, d *schema.ResourceData, id sdk.AccountObjectIdentifier) error {
	// Renaming first means the SET and UNSET below address the deployment by the name it now has.
	if d.HasChange("name") {
		newId := sdk.NewAccountObjectIdentifier(d.Get("name").(string))
		if err := client.OpenflowDeployments.Alter(ctx, sdk.NewAlterOpenflowDeploymentRequest(id).WithRenameTo(newId)); err != nil {
			d.Partial(true)
			return fmt.Errorf("error renaming openflow deployment from %v to %v: %w", id.FullyQualifiedName(), newId.FullyQualifiedName(), err)
		}
		d.SetId(helpers.EncodeResourceIdentifier(newId))
		id = newId
	}

	set, unset := sdk.NewOpenflowDeploymentSetRequest(), sdk.NewOpenflowDeploymentUnsetRequest()

	if err := errors.Join(
		stringAttributeUpdate(d, "display_name", &set.DisplayName, &unset.DisplayName),
		stringAttributeUpdate(d, "comment", &set.Comment, &unset.Comment),
	); err != nil {
		return err
	}
	if diags := handleOpenflowDeploymentParameterUpdate(d, set, unset); diags.HasError() {
		return fmt.Errorf("%s", diags[0].Summary)
	}

	if (*set != sdk.OpenflowDeploymentSetRequest{}) {
		if err := client.OpenflowDeployments.Alter(ctx, sdk.NewAlterOpenflowDeploymentRequest(id).WithSet(*set)); err != nil {
			return err
		}
	}
	if (*unset != sdk.OpenflowDeploymentUnsetRequest{}) {
		if err := client.OpenflowDeployments.Alter(ctx, sdk.NewAlterOpenflowDeploymentRequest(id).WithUnset(*unset)); err != nil {
			return err
		}
	}
	return nil
}

// deleteOpenflowDeployment tears the deployment down in two steps: TERMINATE, which is asynchronous, then
// DROP once the deployment reports DELETED.
//
// A BYOC customer is expected to have torn their own cloud resources down first. Terminating before that can
// leave infrastructure still billing in their account, which Terraform cannot detect, hence the warning in
// the resource documentation.
func deleteOpenflowDeployment(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	timeout := d.Timeout(schema.TimeoutDelete)

	deployment, err := client.OpenflowDeployments.ShowByIDSafely(ctx, id)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	switch deployment.Status {
	case sdk.OpenflowDeploymentStatusDeleting, sdk.OpenflowDeploymentStatusDeleted:
		// Already terminated or on its way; fall through to the wait and drop.
	default:
		// TERMINATE is refused while the deployment is still coming up, so a destroy that runs while
		// provisioning is in flight has to let it settle first.
		if err := waitForOpenflowDeploymentSettled(ctx, client, id, timeout); err != nil {
			return diag.FromErr(err)
		}
		if err := client.OpenflowDeployments.Alter(ctx, sdk.NewAlterOpenflowDeploymentRequest(id).WithIfExists(true).WithTerminate(true)); err != nil {
			return diag.FromErr(err)
		}
	}

	if err := waitForOpenflowDeploymentTerminated(ctx, client, id, timeout); err != nil {
		return diag.FromErr(err)
	}

	if err := client.OpenflowDeployments.DropSafely(ctx, id); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
