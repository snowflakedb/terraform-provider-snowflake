package resources

import (
	"context"
	"errors"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// openflowDeploymentByocSchema adds the networking options that only apply when the deployment runs in the
// customer's own cloud account. All of them are ForceNew: Snowflake has no ALTER for them.
var openflowDeploymentByocSchema = func() map[string]*schema.Schema {
	byocSchema := map[string]*schema.Schema{
		"vpc_type": {
			Type: schema.TypeString,
			// Snowflake requires VPC_TYPE on every BYOC deployment and has no default for it, so marking it
			// required turns a mid-apply SQL failure into a plan-time error naming the argument. Not defaulted
			// on purpose: the two values describe different network topologies.
			Required:         true,
			ForceNew:         true,
			ValidateDiagFunc: sdkValidation(sdk.ToOpenflowVpcType),
			DiffSuppressFunc: SuppressIfAny(NormalizeAndCompare(sdk.ToOpenflowVpcType)),
			Description:      fmt.Sprintf("Specifies whether the deployment's VPC is created by Snowflake or supplied by you. Valid values are (case-insensitive): %s.", possibleValuesListed(sdk.AllOpenflowVpcTypes)),
		},
		"custom_ingress_hostname": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "Specifies a custom hostname for ingress into the deployment.",
		},
		"use_private_link": {
			Type:             schema.TypeString,
			Optional:         true,
			ForceNew:         true,
			ValidateDiagFunc: validateBooleanString,
			DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShow("use_private_link"),
			Description:      booleanStringFieldDescription("Specifies whether the deployment is reached over private link."),
			Default:          BooleanDefault,
		},
		"use_user_auth_over_privatelink": {
			Type:             schema.TypeString,
			Optional:         true,
			ForceNew:         true,
			ValidateDiagFunc: validateBooleanString,
			DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShow("use_user_auth_over_private_link"),
			Description:      booleanStringFieldDescription("Specifies whether user authentication is performed over private link."),
			Default:          BooleanDefault,
		},
	}
	return collections.MergeMaps(openflowDeploymentCommonSchema(), byocSchema)
}()

func OpenflowDeploymentByoc() *schema.Resource {
	return &schema.Resource{
		// TODO(SNOW-4039167): Add the PreviewFeature*ContextWrapper calls when this resource is moved to the
		// production provider. It is registered only in the acceptance test provider for now, so there is no
		// preview feature to gate on yet.
		CreateContext: TrackingCreateWrapper(resources.OpenflowDeploymentByoc, CreateOpenflowDeploymentByoc),
		ReadContext:   TrackingReadWrapper(resources.OpenflowDeploymentByoc, ReadOpenflowDeploymentByoc(true)),
		UpdateContext: TrackingUpdateWrapper(resources.OpenflowDeploymentByoc, UpdateOpenflowDeploymentByoc),
		DeleteContext: TrackingDeleteWrapper(resources.OpenflowDeploymentByoc, deleteOpenflowDeployment),
		Description: joinWithSpace(
			"Resource used to manage BYOC (bring your own cloud) Openflow deployments, which run in your own cloud account.",
			"Creation finishes when the deployment reaches the INACTIVE state, which is as far as Terraform can take it:",
			"the deployment only becomes ACTIVE after you apply the deployment's CloudFormation template in your cloud account, at which point it registers itself with Snowflake.",
			"That template can only be downloaded from the Snowflake UI - it is not retrievable over SQL - so provisioning that infrastructure is a separate step outside this resource.",
			"For more information, check [Openflow deployment documentation](https://docs.snowflake.com/en/sql-reference/sql/create-openflow-deployment).",
		),

		CustomizeDiff: TrackingCustomDiffWrapper(resources.OpenflowDeploymentByoc, customdiff.All(
			RecreateWhenResourceTypeChangedExternally("type", sdk.OpenflowDeploymentTypeByoc, sdk.ToOpenflowDeploymentType),
			ComputedIfAnyAttributeChanged(openflowDeploymentByocSchema, ShowOutputAttributeName, "name", "display_name", "comment"),
			ComputedIfAnyAttributeChanged(openflowDeploymentByocSchema, DescribeOutputAttributeName, "name", "display_name", "comment"),
			ComputedIfAnyAttributeChanged(openflowDeploymentByocSchema, FullyQualifiedNameAttributeName, "name"),
			ComputedIfAnyAttributeChanged(openflowDeploymentByocSchema, ParametersAttributeName, openflowDeploymentEventTableAttribute),
			openflowDeploymentParametersCustomDiff,
		)),

		Schema: openflowDeploymentByocSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.OpenflowDeploymentByoc, ImportOpenflowDeploymentByoc),
		},

		Timeouts: openflowDeploymentTimeouts,
	}
}

func ImportOpenflowDeploymentByoc(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return nil, err
	}

	deployment, err := client.OpenflowDeployments.ShowByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if deployment.Type != sdk.OpenflowDeploymentTypeByoc {
		return nil, errors.New("openflow deployment " + id.Name() + " is of type " + string(deployment.Type) +
			"; import it as snowflake_openflow_deployment_snowflake_managed instead")
	}

	parameters, err := client.OpenflowDeployments.ShowParameters(ctx, id)
	if err != nil {
		return nil, err
	}

	errs := errors.Join(
		d.Set("name", deployment.Name),
		handleOpenflowDeploymentParameterRead(d, parameters),
		importOpenflowDeploymentOptionals(d, deployment),
		d.Set("custom_ingress_hostname", deployment.CustomIngressHostname),
		d.Set("use_private_link", BooleanDefault),
		d.Set("use_user_auth_over_privatelink", BooleanDefault),
	)
	if deployment.VpcType != nil {
		errs = errors.Join(errs, d.Set("vpc_type", string(*deployment.VpcType)))
	}
	if errs != nil {
		return nil, errs
	}
	return []*schema.ResourceData{d}, nil
}

func CreateOpenflowDeploymentByoc(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id := sdk.NewAccountObjectIdentifier(d.Get("name").(string))

	request := sdk.NewCreateOpenflowDeploymentRequest(id, sdk.OpenflowDeploymentTypeByoc)
	if err := errors.Join(
		createOpenflowDeploymentCommonFields(d, request),
		stringAttributeCreateBuilder(d, "custom_ingress_hostname", request.WithCustomIngressHostname),
		booleanStringAttributeCreateBuilder(d, "use_private_link", request.WithUsePrivateLink),
		booleanStringAttributeCreateBuilder(d, "use_user_auth_over_privatelink", request.WithUseUserAuthOverPrivatelink),
		attributeMappedValueCreateBuilder(d, "vpc_type", request.WithVpcType, sdk.ToOpenflowVpcType),
	); err != nil {
		return diag.FromErr(err)
	}

	if err := client.OpenflowDeployments.Create(ctx, request); err != nil {
		return diag.FromErr(err)
	}
	// Set the id before polling so a timeout leaves a resource that can be recovered with terraform
	// import rather than an orphaned deployment.
	d.SetId(helpers.EncodeResourceIdentifier(id))

	// INACTIVE is terminal for a BYOC create: reaching ACTIVE needs the customer's own cloud infrastructure.
	// ACTIVE is accepted too, in case it registers itself while we are still polling.
	if err := waitForOpenflowDeploymentStatus(
		ctx, client, id, d.Timeout(schema.TimeoutCreate),
		[]sdk.OpenflowDeploymentStatus{
			sdk.OpenflowDeploymentStatusInactive,
			sdk.OpenflowDeploymentStatusActive,
		},
		[]sdk.OpenflowDeploymentStatus{
			sdk.OpenflowDeploymentStatusCreateFailed,
			sdk.OpenflowDeploymentStatusDeleted,
		},
	); err != nil {
		return diag.FromErr(err)
	}

	return ReadOpenflowDeploymentByoc(false)(ctx, d, meta)
}

func ReadOpenflowDeploymentByoc(withExternalChangesMarking bool) schema.ReadContextFunc {
	return readOpenflowDeploymentFunc(openflowDeploymentByocSchema, withExternalChangesMarking,
		func(d *schema.ResourceData, deployment *sdk.OpenflowDeployment) error {
			errs := errors.Join(
				d.Set("custom_ingress_hostname", deployment.CustomIngressHostname),
			)
			if deployment.VpcType != nil {
				errs = errors.Join(errs, d.Set("vpc_type", string(*deployment.VpcType)))
			}
			return errs
		})
}

func UpdateOpenflowDeploymentByoc(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// Only the shared metadata attributes are updatable; the networking options are all ForceNew.
	if err := updateOpenflowDeploymentCommonFields(ctx, client, d, id); err != nil {
		return diag.FromErr(err)
	}

	return ReadOpenflowDeploymentByoc(false)(ctx, d, meta)
}
