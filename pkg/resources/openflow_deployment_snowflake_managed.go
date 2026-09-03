package resources

import (
	"context"
	"errors"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// openflowDeploymentSnowflakeManagedSchema is the common deployment schema with no additions: a
// Snowflake-managed deployment takes none of the BYOC networking options, which is the main reason the two
// types are separate resources.
var openflowDeploymentSnowflakeManagedSchema = openflowDeploymentCommonSchema()

func OpenflowDeploymentSnowflakeManaged() *schema.Resource {
	return &schema.Resource{
		// TODO(SNOW-4039167): Add the PreviewFeature*ContextWrapper calls when this resource is moved to the
		// production provider. It is registered only in the acceptance test provider for now, so there is no
		// preview feature to gate on yet.
		CreateContext: TrackingCreateWrapper(resources.OpenflowDeploymentSnowflakeManaged, CreateOpenflowDeploymentSnowflakeManaged),
		ReadContext:   TrackingReadWrapper(resources.OpenflowDeploymentSnowflakeManaged, ReadOpenflowDeploymentSnowflakeManaged(true)),
		UpdateContext: TrackingUpdateWrapper(resources.OpenflowDeploymentSnowflakeManaged, UpdateOpenflowDeploymentSnowflakeManaged),
		DeleteContext: TrackingDeleteWrapper(resources.OpenflowDeploymentSnowflakeManaged, deleteOpenflowDeployment),
		Description: joinWithSpace(
			"Resource used to manage Snowflake-managed Openflow deployments.",
			"A deployment is the account-level container that Openflow runtimes are created in.",
			"Snowflake runs it, so creation is complete once the deployment reaches the ACTIVE state, which takes several minutes.",
			"For more information, check [Openflow deployment documentation](https://docs.snowflake.com/en/sql-reference/sql/create-openflow-deployment).",
		),

		CustomizeDiff: TrackingCustomDiffWrapper(resources.OpenflowDeploymentSnowflakeManaged, customdiff.All(
			RecreateWhenResourceTypeChangedExternally("type", sdk.OpenflowDeploymentTypeSnowflake, sdk.ToOpenflowDeploymentType),
			ComputedIfAnyAttributeChanged(openflowDeploymentSnowflakeManagedSchema, ShowOutputAttributeName, "name", "display_name", "comment"),
			ComputedIfAnyAttributeChanged(openflowDeploymentSnowflakeManagedSchema, DescribeOutputAttributeName, "name", "display_name", "comment"),
			ComputedIfAnyAttributeChanged(openflowDeploymentSnowflakeManagedSchema, FullyQualifiedNameAttributeName, "name"),
			ComputedIfAnyAttributeChanged(openflowDeploymentSnowflakeManagedSchema, ParametersAttributeName, openflowDeploymentEventTableAttribute),
			openflowDeploymentParametersCustomDiff,
		)),

		Schema: openflowDeploymentSnowflakeManagedSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.OpenflowDeploymentSnowflakeManaged, ImportOpenflowDeploymentSnowflakeManaged),
		},

		Timeouts: openflowDeploymentTimeouts,
	}
}

func ImportOpenflowDeploymentSnowflakeManaged(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return nil, err
	}

	deployment, err := client.OpenflowDeployments.ShowByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if deployment.Type != sdk.OpenflowDeploymentTypeSnowflake {
		return nil, errors.New("openflow deployment " + id.Name() + " is of type " + string(deployment.Type) +
			"; import it as snowflake_openflow_deployment_byoc instead")
	}

	parameters, err := client.OpenflowDeployments.ShowParameters(ctx, id)
	if err != nil {
		return nil, err
	}

	if errs := errors.Join(
		d.Set("name", deployment.Name),
		handleOpenflowDeploymentParameterRead(d, parameters),
		importOpenflowDeploymentOptionals(d, deployment),
	); errs != nil {
		return nil, errs
	}
	return []*schema.ResourceData{d}, nil
}

func CreateOpenflowDeploymentSnowflakeManaged(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id := sdk.NewAccountObjectIdentifier(d.Get("name").(string))

	request := sdk.NewCreateOpenflowDeploymentRequest(id, sdk.OpenflowDeploymentTypeSnowflake)
	if err := createOpenflowDeploymentCommonFields(d, request); err != nil {
		return diag.FromErr(err)
	}

	if err := client.OpenflowDeployments.Create(ctx, request); err != nil {
		return diag.FromErr(err)
	}
	// Set the id before polling so a timeout leaves a resource that can be recovered with terraform
	// import rather than an orphaned deployment.
	d.SetId(helpers.EncodeResourceIdentifier(id))

	if err := waitForOpenflowDeploymentStatus(
		ctx, client, id, d.Timeout(schema.TimeoutCreate),
		[]sdk.OpenflowDeploymentStatus{sdk.OpenflowDeploymentStatusActive},
		[]sdk.OpenflowDeploymentStatus{
			sdk.OpenflowDeploymentStatusCreateFailed,
			sdk.OpenflowDeploymentStatusDeleted,
		},
	); err != nil {
		return diag.FromErr(err)
	}

	return ReadOpenflowDeploymentSnowflakeManaged(false)(ctx, d, meta)
}

func ReadOpenflowDeploymentSnowflakeManaged(withExternalChangesMarking bool) schema.ReadContextFunc {
	return readOpenflowDeploymentFunc(openflowDeploymentSnowflakeManagedSchema, withExternalChangesMarking,
		func(d *schema.ResourceData, deployment *sdk.OpenflowDeployment) error { return nil })
}

func UpdateOpenflowDeploymentSnowflakeManaged(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if err := updateOpenflowDeploymentCommonFields(ctx, client, d, id); err != nil {
		return diag.FromErr(err)
	}

	return ReadOpenflowDeploymentSnowflakeManaged(false)(ctx, d, meta)
}
