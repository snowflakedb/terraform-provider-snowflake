package resources

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	openflowRuntimeCreateTimeout = 30 * time.Minute
	openflowRuntimeUpdateTimeout = 30 * time.Minute
	openflowRuntimeDeleteTimeout = 30 * time.Minute
)

var openflowRuntimeTimeouts = &schema.ResourceTimeout{
	Create: schema.DefaultTimeout(openflowRuntimeCreateTimeout),
	Read:   schema.DefaultTimeout(defaultReadTimeout),
	Update: schema.DefaultTimeout(openflowRuntimeUpdateTimeout),
	Delete: schema.DefaultTimeout(openflowRuntimeDeleteTimeout),
}

var openflowRuntimeSchema = map[string]*schema.Schema{
	"database": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("The database in which to create the Openflow runtime."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"schema": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("The schema in which to create the Openflow runtime."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      blocklistedCharactersFieldDescription("Specifies the identifier for the Openflow runtime; must be unique for the schema in which the runtime is created."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"deployment": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("Specifies the Openflow deployment the runtime runs in. Snowflake has no ALTER for it, so changing it recreates the runtime."),
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"node_type": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      fmt.Sprintf("Specifies the size of the runtime's nodes. Snowflake has no ALTER for it, so changing it recreates the runtime. Valid values are (case-insensitive): %s.", possibleValuesListed(sdk.AllOpenflowRuntimeNodeTypes)),
		ValidateDiagFunc: sdkValidation(sdk.ToOpenflowRuntimeNodeType),
		DiffSuppressFunc: SuppressIfAny(NormalizeAndCompare(sdk.ToOpenflowRuntimeNodeType)),
	},
	"min_nodes": {
		Type:             schema.TypeInt,
		Required:         true,
		ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
		Description:      "Specifies the minimum number of nodes the runtime scales down to.",
	},
	"max_nodes": {
		Type:             schema.TypeInt,
		Required:         true,
		ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
		Description:      "Specifies the maximum number of nodes the runtime scales up to.",
	},
	"execute_as_role": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      blocklistedCharactersFieldDescription("Specifies the role the runtime executes as."),
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"external_access_integrations": {
		Type:             schema.TypeSet,
		Optional:         true,
		MinItems:         1,
		Description:      "Specifies the names of the external access integrations that allow the runtime to reach external sites.",
		DiffSuppressFunc: NormalizeAndCompareIdentifiersInSet("external_access_integrations"),
		Elem: &schema.Schema{
			Type:             schema.TypeString,
			ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		},
	},
	"display_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "A free-text alias for the runtime. Shown in the Openflow UI in place of the runtime's identifier when set.",
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the Openflow runtime.",
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW OPENFLOW RUNTIMES` for the given runtime.",
		Elem: &schema.Resource{
			Schema: schemas.ShowOpenflowRuntimeSchema,
		},
	},
	DescribeOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `DESCRIBE OPENFLOW RUNTIME` for the given runtime.",
		Elem: &schema.Resource{
			Schema: schemas.DescribeOpenflowRuntimeSchema,
		},
	},
}

func OpenflowRuntime() *schema.Resource {
	return &schema.Resource{
		// TODO(SNOW-4039167): Add the PreviewFeature*ContextWrapper calls when this resource is moved to the
		// production provider. It is registered only in the acceptance test provider for now, so there is no
		// preview feature to gate on yet.
		CreateContext: TrackingCreateWrapper(resources.OpenflowRuntime, CreateOpenflowRuntime),
		ReadContext:   TrackingReadWrapper(resources.OpenflowRuntime, ReadOpenflowRuntime(true)),
		UpdateContext: TrackingUpdateWrapper(resources.OpenflowRuntime, UpdateOpenflowRuntime),
		DeleteContext: TrackingDeleteWrapper(resources.OpenflowRuntime, DeleteOpenflowRuntime),
		Description: joinWithSpace(
			"Resource used to manage Openflow runtimes, the compute a deployment runs connectors on.",
			"Every mutating statement is asynchronous, so create and update return once the runtime reaches the ACTIVE state.",
			"For more information, check [Openflow runtime documentation](https://docs.snowflake.com/en/sql-reference/sql/create-openflow-runtime).",
		),

		CustomizeDiff: TrackingCustomDiffWrapper(resources.OpenflowRuntime, customdiff.All(
			ComputedIfAnyAttributeChanged(openflowRuntimeSchema, ShowOutputAttributeName, "name", "display_name", "comment", "min_nodes", "max_nodes", "execute_as_role", "external_access_integrations"),
			ComputedIfAnyAttributeChanged(openflowRuntimeSchema, DescribeOutputAttributeName, "name", "display_name", "comment", "min_nodes", "max_nodes", "execute_as_role", "external_access_integrations"),
			ComputedIfAnyAttributeChanged(openflowRuntimeSchema, FullyQualifiedNameAttributeName, "name"),
		)),

		Schema:   openflowRuntimeSchema,
		Timeouts: openflowRuntimeTimeouts,
		Importer: &schema.ResourceImporter{
			StateContext: ImportOpenflowRuntime,
		},
	}
}

func ImportOpenflowRuntime(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return nil, err
	}
	runtime, err := client.OpenflowRuntimes.ShowByID(ctx, id)
	if err != nil {
		return nil, err
	}

	errs := errors.Join(
		d.Set("database", id.DatabaseName()),
		d.Set("schema", id.SchemaName()),
		d.Set("name", id.Name()),
		d.Set("deployment", runtime.Deployment),
		d.Set("node_type", string(runtime.NodeType)),
		d.Set("min_nodes", runtime.MinNodes),
		d.Set("max_nodes", runtime.MaxNodes),
		d.Set("execute_as_role", runtime.ExecuteAsRole),
		d.Set("external_access_integrations", collections.Map(runtime.ExternalAccessIntegrations, sdk.AccountObjectIdentifier.Name)),
	)
	if runtime.DisplayName != nil && *runtime.DisplayName != "" {
		errs = errors.Join(errs, d.Set("display_name", *runtime.DisplayName))
	}
	if runtime.Comment != nil && *runtime.Comment != "" {
		errs = errors.Join(errs, d.Set("comment", *runtime.Comment))
	}
	if errs != nil {
		return nil, errs
	}
	return []*schema.ResourceData{d}, nil
}

func CreateOpenflowRuntime(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id := sdk.NewSchemaObjectIdentifier(d.Get("database").(string), d.Get("schema").(string), d.Get("name").(string))

	deployment, err := sdk.ParseAccountObjectIdentifier(d.Get("deployment").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	nodeType, err := sdk.ToOpenflowRuntimeNodeType(d.Get("node_type").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	executeAsRole, err := sdk.ParseAccountObjectIdentifier(d.Get("execute_as_role").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	request := sdk.NewCreateOpenflowRuntimeRequest(id, deployment, nodeType, d.Get("min_nodes").(int), d.Get("max_nodes").(int), executeAsRole)
	if err := errors.Join(
		stringAttributeCreateBuilder(d, "display_name", request.WithDisplayName),
		stringAttributeCreateBuilder(d, "comment", request.WithComment),
		attributeMappedValueCreateBuilder(d, "external_access_integrations", request.WithExternalAccessIntegrations, ToOpenflowRuntimeExternalAccessIntegrationsRequest),
	); err != nil {
		return diag.FromErr(err)
	}

	if err := client.OpenflowRuntimes.Create(ctx, request); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(helpers.EncodeResourceIdentifier(id))

	if err := waitForOpenflowRuntimeActive(ctx, client, id, d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.FromErr(err)
	}
	return ReadOpenflowRuntime(false)(ctx, d, meta)
}

func ReadOpenflowRuntime(withExternalChangesMarking bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		client := meta.(*provider.Context).Client
		id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}

		runtime, err := client.OpenflowRuntimes.ShowByIDSafely(ctx, id)
		if err != nil {
			if errors.Is(err, sdk.ErrObjectNotFound) {
				d.SetId("")
				return diag.Diagnostics{
					diag.Diagnostic{
						Severity: diag.Warning,
						Summary:  "Failed to query openflow runtime. Marking the resource as removed.",
						Detail:   fmt.Sprintf("Openflow runtime: %s, Err: %s", id.FullyQualifiedName(), err),
					},
				}
			}
			return diag.FromErr(err)
		}

		runtimeDetails, err := client.OpenflowRuntimes.Describe(ctx, id)
		if err != nil {
			return diag.FromErr(err)
		}

		// Flattened up front: the drift comparison is an `any` equality against the string in show_output, so a
		// *string would never match and would mark drift on every read.
		var displayName, comment string
		if runtime.DisplayName != nil {
			displayName = *runtime.DisplayName
		}
		if runtime.Comment != nil {
			comment = *runtime.Comment
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

		if err = setStateToValuesFromConfig(d, openflowRuntimeSchema, []string{
			"display_name",
			"comment",
		}); err != nil {
			return diag.FromErr(err)
		}

		if errs := errors.Join(
			// Read back so that an external change to a ForceNew attribute diverges from config and forces
			// replacement on the next plan.
			d.Set("deployment", runtime.Deployment),
			d.Set("node_type", string(runtime.NodeType)),
			d.Set("min_nodes", runtime.MinNodes),
			d.Set("max_nodes", runtime.MaxNodes),
			d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
			d.Set(ShowOutputAttributeName, []map[string]any{schemas.OpenflowRuntimeToSchema(runtime)}),
			d.Set("execute_as_role", runtime.ExecuteAsRole),
			d.Set(DescribeOutputAttributeName, []map[string]any{schemas.OpenflowRuntimeDetailsToSchema(*runtimeDetails)}),
		); errs != nil {
			return diag.FromErr(errs)
		}
		return nil
	}
}

func UpdateOpenflowRuntime(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	timeout := d.Timeout(schema.TimeoutUpdate)

	// ALTER is refused while the runtime is in a transient status, so an apply that lands while another
	// change is still in flight has to let it settle first.
	if err := waitForOpenflowRuntimeSettled(ctx, client, id, timeout); err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("name") {
		newId := sdk.NewSchemaObjectIdentifierInSchema(id.SchemaId(), d.Get("name").(string))
		if err := client.OpenflowRuntimes.Alter(ctx, sdk.NewAlterOpenflowRuntimeRequest(id).WithRenameTo(newId)); err != nil {
			return diag.FromErr(err)
		}
		if err := waitForOpenflowRuntimeUpdated(ctx, client, newId, timeout); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(helpers.EncodeResourceIdentifier(newId))
		id = newId
	}

	set, unset := sdk.NewOpenflowRuntimeSetRequest(), sdk.NewOpenflowRuntimeUnsetRequest()
	// execute_as_role is required, so it has no UNSET and the helper's unset target can never be reached.
	var executeAsRoleUnset *bool
	if err := errors.Join(
		stringAttributeUpdate(d, "display_name", &set.DisplayName, &unset.DisplayName),
		stringAttributeUpdate(d, "comment", &set.Comment, &unset.Comment),
		attributeMappedValueUpdate(d, "external_access_integrations", &set.ExternalAccessIntegrations, &unset.ExternalAccessIntegrations, ToOpenflowRuntimeExternalAccessIntegrationsRequest),
		intAttributeUpdateSetOnly(d, "min_nodes", &set.MinNodes),
		intAttributeUpdateSetOnly(d, "max_nodes", &set.MaxNodes),
		accountObjectIdentifierAttributeUpdate(d, "execute_as_role", &set.ExecuteAsRole, &executeAsRoleUnset),
	); err != nil {
		return diag.FromErr(err)
	}

	if (*set != sdk.OpenflowRuntimeSetRequest{}) {
		if err := client.OpenflowRuntimes.Alter(ctx, sdk.NewAlterOpenflowRuntimeRequest(id).WithSet(*set)); err != nil {
			return diag.FromErr(err)
		}
		if err := waitForOpenflowRuntimeUpdated(ctx, client, id, timeout); err != nil {
			return diag.FromErr(err)
		}
	}
	if (*unset != sdk.OpenflowRuntimeUnsetRequest{}) {
		if err := client.OpenflowRuntimes.Alter(ctx, sdk.NewAlterOpenflowRuntimeRequest(id).WithUnset(*unset)); err != nil {
			return diag.FromErr(err)
		}
		if err := waitForOpenflowRuntimeUpdated(ctx, client, id, timeout); err != nil {
			return diag.FromErr(err)
		}
	}

	return ReadOpenflowRuntime(false)(ctx, d, meta)
}

func DeleteOpenflowRuntime(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	timeout := d.Timeout(schema.TimeoutDelete)

	runtime, err := client.OpenflowRuntimes.ShowByIDSafely(ctx, id)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	switch runtime.Status {
	case sdk.OpenflowRuntimeStatusDeleting, sdk.OpenflowRuntimeStatusDeleted:
		// Already terminated or on its way; fall through to the wait and drop.
	default:
		// TERMINATE is refused while the runtime is still in a transient status, so a destroy that runs
		// while a create or an update is in flight has to let it settle first.
		if err := waitForOpenflowRuntimeSettled(ctx, client, id, timeout); err != nil {
			return diag.FromErr(err)
		}
		if err := client.OpenflowRuntimes.Alter(ctx, sdk.NewAlterOpenflowRuntimeRequest(id).WithIfExists(true).WithTerminate(true)); err != nil {
			return diag.FromErr(err)
		}
	}

	// DROP is only permitted from DELETED.
	if err := waitForOpenflowRuntimeTerminated(ctx, client, id, timeout); err != nil {
		return diag.FromErr(err)
	}
	if err := client.OpenflowRuntimes.DropSafely(ctx, id); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}

func ToOpenflowRuntimeExternalAccessIntegrationsRequest(value any) (sdk.OpenflowRuntimeExternalAccessIntegrationsRequest, error) {
	raw := expandStringList(value.(*schema.Set).List())
	integrations := make([]sdk.AccountObjectIdentifier, len(raw))
	for i, v := range raw {
		integrations[i] = sdk.NewAccountObjectIdentifier(v)
	}
	return sdk.OpenflowRuntimeExternalAccessIntegrationsRequest{
		ExternalAccessIntegrations: integrations,
	}, nil
}
