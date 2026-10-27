package resources

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// =============================================================================
// Constructor
// =============================================================================

var (
	catalogIntegrationAwsGlueCreateExt = TrackingCreateWrapper(resources.CatalogIntegrationAwsGlue, CreateCatalogIntegrationAwsGlue)
	catalogIntegrationAwsGlueReadExt   = TrackingReadWrapper(resources.CatalogIntegrationAwsGlue, ReadCatalogIntegrationAwsGlueFunc(true))
	catalogIntegrationAwsGlueUpdateExt = TrackingUpdateWrapper(resources.CatalogIntegrationAwsGlue, UpdateCatalogIntegrationAwsGlue)
	// Same delete wiring as deleteCatalogIntegrationFunc: ParseAccountObjectIdentifier and CatalogIntegrations.DropSafely.
	catalogIntegrationAwsGlueDeleteExt = TrackingDeleteWrapper(resources.CatalogIntegrationAwsGlue, ResourceDeleteContextFunc(
		catalogIntegrationAwsGlueParseIdExt,
		catalogIntegrationAwsGlueDropSafelyInSdkExt,
	))
	catalogIntegrationAwsGlueDescriptionExt = "Resource used to manage AWS Glue catalog integration objects. For more information, check [catalog integration documentation](https://docs.snowflake.com/en/sql-reference/sql/create-catalog-integration-glue)."
	// CustomizeDiff was not wrapped in TrackingCustomDiffWrapper.
	catalogIntegrationAwsGlueCustomizeDiffExt = customdiff.All(
		ComputedIfAnyAttributeChanged(catalogIntegrationAwsGlueSchema, ShowOutputAttributeName, "enabled", "comment"),
		ComputedIfAnyAttributeChanged(catalogIntegrationAwsGlueSchema, DescribeOutputAttributeName, "enabled", "refresh_interval_seconds", "comment"),
		RecreateWhenCatalogSourceChangedExternally(sdk.CatalogIntegrationCatalogSourceTypeGlue),
	)
	catalogIntegrationAwsGlueImporterExt = &schema.ResourceImporter{
		StateContext: TrackingImportWrapper(resources.CatalogIntegrationAwsGlue, ImportCatalogIntegrationAwsGlue),
	}
)

// =============================================================================
// Request / mapping / Read-Update-Import helpers
// =============================================================================

func catalogIntegrationAwsGlueParseIdExt(id string) (sdk.AccountObjectIdentifier, error) {
	return sdk.ParseAccountObjectIdentifier(id)
}

func catalogIntegrationAwsGlueNewCreateRequestExt(d *schema.ResourceData, id sdk.AccountObjectIdentifier) (*sdk.CreateCatalogIntegrationRequest, error) {
	enabled := d.Get("enabled").(bool)
	glueAwsRoleArn := d.Get("glue_aws_role_arn").(string)
	glueCatalogId := d.Get("glue_catalog_id").(string)

	request := sdk.NewCreateCatalogIntegrationRequest(id, enabled)
	awsGlueRequest := sdk.NewAwsGlueParamsRequest(glueAwsRoleArn, glueCatalogId)
	// Attached before optionals so glue_region and catalog_namespace can be written
	// onto this params pointer. The create request matches the previous post-optional copy.
	return request.WithAwsGlueCatalogSourceParams(*awsGlueRequest), nil
}

func catalogIntegrationAwsGlueApplyCreateOptionalsExt(d *schema.ResourceData, request *sdk.CreateCatalogIntegrationRequest) error {
	return errors.Join(
		intAttributeCreateBuilder(d, "refresh_interval_seconds", request.WithRefreshIntervalSeconds),
		stringAttributeCreateBuilder(d, "comment", request.WithComment),
		stringAttributeCreateBuilder(d, "glue_region", request.AwsGlueCatalogSourceParams.WithGlueRegion),
		stringAttributeCreateBuilder(d, "catalog_namespace", request.AwsGlueCatalogSourceParams.WithCatalogNamespace),
	)
}

// SET only. Catalog integration ALTER has no UNSET request (SNOW-3243983).
type catalogIntegrationAwsGlueAlterRequestsExt struct {
	set *sdk.CatalogIntegrationSetRequest
}

func catalogIntegrationAwsGlueHandleExternalChangesInDescribeExt(d *schema.ResourceData, details *sdk.CatalogIntegrationAwsGlueDetails, withExternalChangesMarking bool) diag.Diagnostics {
	if withExternalChangesMarking {
		if err := handleExternalChangesToObjectInFlatDescribe(d, catalogIntegrationAwsGlueOutputMappingsExt(details)...); err != nil {
			return diag.FromErr(err)
		}
	}
	// refresh_interval_seconds and glue_region are not copied into config; drift is marked on describe output.
	if err := errors.Join(
		d.Set("catalog_source", string(details.CatalogSource)),
		d.Set("glue_aws_role_arn", details.GlueAwsRoleArn),
		d.Set("glue_catalog_id", details.GlueCatalogId),
		d.Set("catalog_namespace", details.CatalogNamespace),
	); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func catalogIntegrationAwsGlueOutputMappingsExt(details *sdk.CatalogIntegrationAwsGlueDetails) []outputMapping {
	return []outputMapping{
		{"refresh_interval_seconds", "refresh_interval_seconds", details.RefreshIntervalSeconds, details.RefreshIntervalSeconds, nil},
		{"glue_region", "glue_region", details.GlueRegion, details.GlueRegion, nil},
	}
}

func catalogIntegrationAwsGlueSetConfigFieldsExt(d *schema.ResourceData, catalogIntegration *sdk.CatalogIntegration) error {
	return errors.Join(
		d.Set("name", catalogIntegration.Name),
		d.Set("enabled", catalogIntegration.Enabled),
		d.Set("comment", catalogIntegration.Comment),
	)
}

func catalogIntegrationAwsGlueApplySetUnsetExt(d *schema.ResourceData) (*catalogIntegrationAwsGlueAlterRequestsExt, error) {
	set := sdk.NewCatalogIntegrationSetRequest()
	if errs := handleCatalogIntegrationCommonPropsUpdate(d, set); errs != nil {
		return nil, errs
	}
	return &catalogIntegrationAwsGlueAlterRequestsExt{set: set}, nil
}

func catalogIntegrationAwsGlueSetFieldsNotSetByReadExt(_ *schema.ResourceData, _ sdk.AccountObjectIdentifier, _ *sdk.CatalogIntegrationAwsGlueDetails, _ bool) error {
	return nil
}

// =============================================================================
// SDK incision points
// =============================================================================

func catalogIntegrationAwsGlueCreateInSdkExt(ctx context.Context, meta any, _ sdk.AccountObjectIdentifier, req *sdk.CreateCatalogIntegrationRequest) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	if err := client.CatalogIntegrations.Create(ctx, req); err != nil {
		return diag.FromErr(fmt.Errorf("error creating AWS Glue catalog integration, err = %w", err))
	}
	return nil
}

func catalogIntegrationAwsGlueShowByIdSafelyInSdkExt(ctx context.Context, d *schema.ResourceData, meta any, id sdk.AccountObjectIdentifier) (*sdk.CatalogIntegration, diag.Diagnostics) {
	client := meta.(*provider.Context).Client
	s, err := client.CatalogIntegrations.ShowByIDSafely(ctx, id)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			d.SetId("")
			return nil, diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to query AWS Glue catalog integration. Marking the resource as removed.",
					Detail:   fmt.Sprintf("AWS Glue catalog integration id: %s, Err: %s", id.FullyQualifiedName(), err),
				},
			}
		}
		return nil, diag.FromErr(err)
	}
	return s, nil
}

func catalogIntegrationAwsGlueDescribeInSdkExt(ctx context.Context, meta any, id sdk.AccountObjectIdentifier) (*sdk.CatalogIntegrationAwsGlueDetails, error) {
	client := meta.(*provider.Context).Client
	details, err := client.CatalogIntegrations.DescribeAwsGlueDetails(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("could not describe AWS Glue catalog integration (%s), err = %w", helpers.EncodeResourceIdentifier(id), err)
	}
	return details, nil
}

// Describe plus the Glue catalog-source guard. Import does not call ShowByID.
func catalogIntegrationAwsGlueShowByIdInSdkExt(ctx context.Context, meta any, id sdk.AccountObjectIdentifier) (*sdk.CatalogIntegrationAwsGlueDetails, error) {
	client := meta.(*provider.Context).Client
	details, err := client.CatalogIntegrations.DescribeAwsGlueDetails(ctx, id)
	if err != nil {
		return nil, err
	}
	if details.CatalogSource != sdk.CatalogIntegrationCatalogSourceTypeGlue {
		return nil, fmt.Errorf("invalid catalog source type, expected %s, got %s", sdk.CatalogIntegrationCatalogSourceTypeGlue, details.CatalogSource)
	}
	return details, nil
}

func catalogIntegrationAwsGlueAlterSetInSdkExt(ctx context.Context, _ *schema.ResourceData, meta any, id sdk.AccountObjectIdentifier, reqs *catalogIntegrationAwsGlueAlterRequestsExt) error {
	if reflect.DeepEqual(*reqs.set, *sdk.NewCatalogIntegrationSetRequest()) {
		return nil
	}
	client := meta.(*provider.Context).Client
	req := sdk.NewAlterCatalogIntegrationRequest(id).WithSet(*reqs.set)
	if err := client.CatalogIntegrations.Alter(ctx, req); err != nil {
		return fmt.Errorf("error updating catalog integration (%s), err = %w", helpers.EncodeResourceIdentifier(id), err)
	}
	return nil
}

func catalogIntegrationAwsGlueAlterUnsetInSdkExt(context.Context, any, sdk.AccountObjectIdentifier, *catalogIntegrationAwsGlueAlterRequestsExt) error {
	return nil
}

func catalogIntegrationAwsGlueDropSafelyInSdkExt(client *sdk.Client) DropSafelyFunc[sdk.AccountObjectIdentifier] {
	return client.CatalogIntegrations.DropSafely
}
