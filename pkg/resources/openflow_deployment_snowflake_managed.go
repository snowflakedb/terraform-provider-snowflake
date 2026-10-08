package resources

import (
	"context"
	"errors"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// openflowDeploymentSnowflakeManagedSchema is the common deployment schema with no additions: a
// Snowflake-managed deployment takes none of the BYOC networking options, which is the main reason the two
// types are separate resources.
//
// This must stay the base schema only (not merged with openflowDeploymentParametersSchema): that map is
// populated by an init() in another file, which runs after package-level var initializers, so merging it
// here at var-init time would silently copy it while still empty. The merge happens in
// OpenflowDeploymentSnowflakeManaged's Schema field instead, which runs after all init()s.
var openflowDeploymentSnowflakeManagedSchema = openflowDeploymentCommonSchema()

func OpenflowDeploymentSnowflakeManaged() *schema.Resource {
	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.OpenflowDeploymentSnowflakeManagedResource), TrackingCreateWrapper(resources.OpenflowDeploymentSnowflakeManaged, CreateOpenflowDeploymentSnowflakeManaged)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.OpenflowDeploymentSnowflakeManagedResource), TrackingReadWrapper(resources.OpenflowDeploymentSnowflakeManaged, ReadOpenflowDeploymentSnowflakeManaged(true))),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.OpenflowDeploymentSnowflakeManagedResource), TrackingUpdateWrapper(resources.OpenflowDeploymentSnowflakeManaged, UpdateOpenflowDeploymentSnowflakeManaged)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.OpenflowDeploymentSnowflakeManagedResource), TrackingDeleteWrapper(resources.OpenflowDeploymentSnowflakeManaged, deleteOpenflowDeployment)),
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
			ComputedIfAnyAttributeChanged(openflowDeploymentParametersSchema, ParametersAttributeName, strings.ToLower(string(sdk.OpenflowDeploymentParameterEventTable))),
			openflowDeploymentParametersCustomDiff,
		)),

		Schema: collections.MergeMaps(openflowDeploymentSnowflakeManagedSchema, openflowDeploymentParametersSchema),
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

	parameters, err := client.OpenflowDeployments.ShowParametersDetails(ctx, id)
	if err != nil {
		return nil, err
	}
	if diags := handleOpenflowDeploymentParameterRead(d, parameters); diags.HasError() {
		return nil, errorFromDiagnostics(diags)
	}

	if errs := errors.Join(
		d.Set("name", deployment.Name),
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
