package resources

import (
	"context"
	"errors"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var openflowDeploymentSchema = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("Identifier for the openflow deployment; must be unique for the account."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"deployment_type": {
		Type:             schema.TypeString,
		Optional:         true,
		ForceNew:         true,
		Default:          string(sdk.OpenflowDeploymentTypeSnowflake),
		ValidateDiagFunc: sdkValidation(sdk.ToOpenflowDeploymentType),
		DiffSuppressFunc: SuppressIfAny(NormalizeAndCompare(sdk.ToOpenflowDeploymentType)),
		Description:      fmt.Sprintf("Type of the deployment. Valid values: %s.", possibleValuesListed(sdk.AllOpenflowDeploymentTypes)),
	},
	"vpc_type": {
		Type:             schema.TypeString,
		Optional:         true,
		ForceNew:         true,
		ValidateDiagFunc: sdkValidation(sdk.ToOpenflowVpcType),
		DiffSuppressFunc: SuppressIfAny(NormalizeAndCompare(sdk.ToOpenflowVpcType)),
		Description:      fmt.Sprintf("VPC type for BYOC deployments. Valid values: %s.", possibleValuesListed(sdk.AllOpenflowVpcTypes)),
	},
	"custom_ingress_hostname": {
		Type:        schema.TypeString,
		Optional:    true,
		ForceNew:    true,
		Description: "Custom ingress hostname for BYOC deployments.",
	},
	"use_private_link": {
		Type:        schema.TypeBool,
		Optional:    true,
		ForceNew:    true,
		Default:     false,
		Description: "Whether to use private link for the deployment.",
	},
	"use_user_auth_over_privatelink": {
		Type:        schema.TypeBool,
		Optional:    true,
		ForceNew:    true,
		Default:     false,
		Description: "Whether to use user auth over private link.",
	},
	"event_table": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Fully-qualified name of an event table for the deployment (db.schema.table), or empty to unset.",
	},
	"display_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Human-readable display name for the deployment.",
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Comment for the deployment.",
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Output of SHOW OPENFLOW DEPLOYMENTS for this deployment.",
		Elem:        &schema.Resource{Schema: schemas.ShowOpenflowDeploymentSchema},
	},
	DescribeOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Output of DESCRIBE OPENFLOW DEPLOYMENT for this deployment.",
		Elem:        &schema.Resource{Schema: schemas.DescribeOpenflowDeploymentSchema},
	},
}

func OpenflowDeployment() *schema.Resource {
	return &schema.Resource{
		CreateContext: TrackingCreateWrapper(resources.OpenflowDeployment, CreateOpenflowDeployment),
		ReadContext:   TrackingReadWrapper(resources.OpenflowDeployment, ReadOpenflowDeploymentFunc(true)),
		UpdateContext: TrackingUpdateWrapper(resources.OpenflowDeployment, UpdateOpenflowDeployment),
		DeleteContext: TrackingDeleteWrapper(resources.OpenflowDeployment, DeleteOpenflowDeployment),
		Description:   "Manages an Openflow Deployment. Deployments are the account-scoped top-level objects for Openflow SOM (Service Object Model).",

		CustomizeDiff: TrackingCustomDiffWrapper(resources.OpenflowDeployment, customdiff.All(
			ComputedIfAnyAttributeChanged(openflowDeploymentSchema, ShowOutputAttributeName, "comment", "display_name", "event_table"),
			ComputedIfAnyAttributeChanged(openflowDeploymentSchema, DescribeOutputAttributeName, "comment", "display_name", "event_table"),
		)),

		Schema: openflowDeploymentSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.OpenflowDeployment, ImportOpenflowDeployment),
		},
		Timeouts: defaultTimeouts,
	}
}

func ImportOpenflowDeployment(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return nil, err
	}
	deployment, err := client.OpenflowDeployments.ShowByID(ctx, id)
	if err != nil {
		return nil, err
	}
	errs := errors.Join(
		d.Set("name", deployment.Name),
		d.Set("deployment_type", string(deployment.DeploymentType)),
		d.Set("use_private_link", deployment.UsePrivateLink),
		d.Set("use_user_auth_over_privatelink", deployment.UseUserAuthOverPrivatelink),
	)
	if deployment.VpcType != nil {
		errs = errors.Join(errs, d.Set("vpc_type", string(*deployment.VpcType)))
	}
	if deployment.CustomIngressHostname != nil {
		errs = errors.Join(errs, d.Set("custom_ingress_hostname", *deployment.CustomIngressHostname))
	}
	if errs != nil {
		return nil, errs
	}
	return []*schema.ResourceData{d}, nil
}

func CreateOpenflowDeployment(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id := sdk.NewAccountObjectIdentifier(d.Get("name").(string))

	request := sdk.NewCreateOpenflowDeploymentRequest(id).WithIfNotExists(true)

	if v, ok := d.GetOk("deployment_type"); ok {
		dt, err := sdk.ToOpenflowDeploymentType(v.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		request.WithDeploymentType(dt)
	}
	if v, ok := d.GetOk("vpc_type"); ok {
		vt, err := sdk.ToOpenflowVpcType(v.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		request.WithVpcType(vt)
	}
	if v, ok := d.GetOk("custom_ingress_hostname"); ok {
		request.WithCustomIngressHostname(v.(string))
	}
	if v := d.Get("use_private_link").(bool); v {
		request.WithUsePrivateLink(v)
	}
	if v := d.Get("use_user_auth_over_privatelink").(bool); v {
		request.WithUseUserAuthOverPrivatelink(v)
	}
	if v, ok := d.GetOk("event_table"); ok {
		request.WithEventTable(v.(string))
	}
	if v, ok := d.GetOk("display_name"); ok {
		request.WithDisplayName(v.(string))
	}
	if v, ok := d.GetOk("comment"); ok {
		request.WithComment(v.(string))
	}

	if err := client.OpenflowDeployments.Create(ctx, request); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(helpers.EncodeResourceIdentifier(id))

	// BYOC deployments start INACTIVE (user must provision CloudFormation), so
	// use the longer timeout. SNOWFLAKE type goes CREATING → ACTIVE directly.
	deploymentType := sdk.OpenflowDeploymentType(d.Get("deployment_type").(string))
	if deploymentType == sdk.OpenflowDeploymentTypeByoc {
		if err := waitForOpenflowDeploymentActiveByoc(ctx, client, id); err != nil {
			return diag.Diagnostics{{
				Severity: diag.Warning,
				Summary:  "BYOC deployment is waiting for CloudFormation provisioning",
				Detail: fmt.Sprintf("Deployment %s is INACTIVE. Provision the CloudFormation stack in your VPC, then run `terraform apply` again. Error: %s", id.Name(), err),
			}}
		}
	} else {
		if err := waitForOpenflowDeploymentActive(ctx, client, id); err != nil {
			return diag.FromErr(err)
		}
	}

	return ReadOpenflowDeploymentFunc(false)(ctx, d, meta)
}

func ReadOpenflowDeploymentFunc(withExternalChangesMarking bool) schema.ReadContextFunc {
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
				return diag.Diagnostics{{
					Severity: diag.Warning,
					Summary:  "Failed to query openflow deployment. Marking as removed.",
					Detail:   fmt.Sprintf("Deployment id: %s, Err: %s", id.FullyQualifiedName(), err),
				}}
			}
			return diag.FromErr(err)
		}

		details, err := client.OpenflowDeployments.Describe(ctx, id)
		if err != nil {
			return diag.FromErr(err)
		}

		if withExternalChangesMarking {
			var eventTable, displayName, comment string
			if deployment.EventTable != nil {
				eventTable = *deployment.EventTable
			}
			if deployment.DisplayName != nil {
				displayName = *deployment.DisplayName
			}
			if deployment.Comment != nil {
				comment = *deployment.Comment
			}
			if err = handleExternalChangesToObjectInShow(d,
				outputMapping{"event_table", "event_table", eventTable, eventTable, nil},
				outputMapping{"display_name", "display_name", displayName, displayName, nil},
				outputMapping{"comment", "comment", comment, comment, nil},
			); err != nil {
				return diag.FromErr(err)
			}
		}

		if err = setStateToValuesFromConfig(d, openflowDeploymentSchema, []string{
			"event_table", "display_name", "comment",
		}); err != nil {
			return diag.FromErr(err)
		}

		errs := errors.Join(
			d.Set(ShowOutputAttributeName, []map[string]any{schemas.OpenflowDeploymentToSchema(deployment)}),
			d.Set(DescribeOutputAttributeName, []map[string]any{schemas.OpenflowDeploymentDetailsToSchema(*details)}),
			d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
			d.Set("deployment_type", string(deployment.DeploymentType)),
			d.Set("use_private_link", deployment.UsePrivateLink),
			d.Set("use_user_auth_over_privatelink", deployment.UseUserAuthOverPrivatelink),
		)
		if deployment.VpcType != nil {
			errs = errors.Join(errs, d.Set("vpc_type", string(*deployment.VpcType)))
		}
		if deployment.CustomIngressHostname != nil {
			errs = errors.Join(errs, d.Set("custom_ingress_hostname", *deployment.CustomIngressHostname))
		}
		if deployment.EventTable != nil {
			errs = errors.Join(errs, d.Set("event_table", *deployment.EventTable))
		}
		if deployment.DisplayName != nil {
			errs = errors.Join(errs, d.Set("display_name", *deployment.DisplayName))
		}
		if deployment.Comment != nil {
			errs = errors.Join(errs, d.Set("comment", *deployment.Comment))
		}
		if errs != nil {
			return diag.FromErr(errs)
		}
		return nil
	}
}

func UpdateOpenflowDeployment(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	set := sdk.NewOpenflowDeploymentSetRequest()
	unset := sdk.NewOpenflowDeploymentUnsetRequest()

	if d.HasChange("comment") {
		if v, ok := d.GetOk("comment"); ok {
			set.WithComment(v.(string))
		} else {
			unset.WithComment(true)
		}
	}
	if d.HasChange("display_name") {
		if v, ok := d.GetOk("display_name"); ok {
			set.WithDisplayName(v.(string))
		} else {
			unset.WithDisplayName(true)
		}
	}
	if d.HasChange("event_table") {
		if v, ok := d.GetOk("event_table"); ok {
			set.WithEventTable(v.(string))
		} else {
			unset.WithEventTable(true)
		}
	}

	if set.Comment != nil || set.DisplayName != nil || set.EventTable != nil {
		if err := client.OpenflowDeployments.Alter(ctx, sdk.NewAlterOpenflowDeploymentRequest(id).WithSet(*set)); err != nil {
			return diag.FromErr(err)
		}
	}
	if unset.Comment != nil || unset.DisplayName != nil || unset.EventTable != nil {
		if err := client.OpenflowDeployments.Alter(ctx, sdk.NewAlterOpenflowDeploymentRequest(id).WithUnset(*unset)); err != nil {
			return diag.FromErr(err)
		}
	}

	return ReadOpenflowDeploymentFunc(false)(ctx, d, meta)
}

func DeleteOpenflowDeployment(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	if err := client.OpenflowDeployments.DropSafely(ctx, id); err != nil {
		return diag.FromErr(err)
	}
	if err := waitForOpenflowDeploymentDeleted(ctx, client, id); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}
