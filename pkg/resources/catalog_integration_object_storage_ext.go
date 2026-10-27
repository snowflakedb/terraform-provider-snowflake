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
	catalogIntegrationObjectStorageCreateExt      = TrackingCreateWrapper(resources.CatalogIntegrationObjectStorage, CreateCatalogIntegrationObjectStorage)
	catalogIntegrationObjectStorageReadExt        = TrackingReadWrapper(resources.CatalogIntegrationObjectStorage, ReadCatalogIntegrationObjectStorageFunc(true))
	catalogIntegrationObjectStorageUpdateExt      = TrackingUpdateWrapper(resources.CatalogIntegrationObjectStorage, UpdateCatalogIntegrationObjectStorage)
	catalogIntegrationObjectStorageDescriptionExt = "Resource used to manage catalog integration objects for Apache Iceberg™ table files or Delta table files in object storage. For more information, check [catalog integration documentation](https://docs.snowflake.com/en/sql-reference/sql/create-catalog-integration-object-storage)."
	// No TrackingCustomDiffWrapper: the resource did not use one.
	catalogIntegrationObjectStorageCustomizeDiffExt = customdiff.All(
		ComputedIfAnyAttributeChanged(catalogIntegrationObjectStorageSchema, ShowOutputAttributeName, "enabled", "comment"),
		ComputedIfAnyAttributeChanged(catalogIntegrationObjectStorageSchema, DescribeOutputAttributeName, "enabled", "refresh_interval_seconds", "comment"),
		RecreateWhenCatalogSourceChangedExternally(sdk.CatalogIntegrationCatalogSourceTypeObjectStore),
	)
	catalogIntegrationObjectStorageImporterExt = &schema.ResourceImporter{
		StateContext: TrackingImportWrapper(resources.CatalogIntegrationObjectStorage, ImportCatalogIntegrationObjectStorage),
	}
)

// =============================================================================
// Request / mapping / Read-Update-Import helpers
// =============================================================================

// SET only. Catalog integration ALTER has no UNSET request (SNOW-3243983).
type catalogIntegrationObjectStorageAlterRequestsExt struct {
	set *sdk.CatalogIntegrationSetRequest
}

func catalogIntegrationObjectStorageParseIdExt(id string) (sdk.AccountObjectIdentifier, error) {
	return sdk.ParseAccountObjectIdentifier(id)
}

func catalogIntegrationObjectStorageNewCreateRequestExt(d *schema.ResourceData, id sdk.AccountObjectIdentifier) (*sdk.CreateObjectStorageCatalogIntegrationRequest, error) {
	enabled := d.Get("enabled").(bool)
	tableFormat, err := sdk.ToCatalogIntegrationTableFormat(d.Get("table_format").(string))
	if err != nil {
		return nil, err
	}

	return sdk.NewCreateObjectStorageCatalogIntegrationRequest(id, tableFormat, enabled), nil
}

func catalogIntegrationObjectStorageApplyCreateOptionalsExt(d *schema.ResourceData, request *sdk.CreateObjectStorageCatalogIntegrationRequest) error {
	return errors.Join(
		intAttributeCreateBuilder(d, "refresh_interval_seconds", request.WithRefreshIntervalSeconds),
		stringAttributeCreateBuilder(d, "comment", request.WithComment),
	)
}

func catalogIntegrationObjectStorageHandleExternalChangesInDescribeExt(d *schema.ResourceData, details *sdk.CatalogIntegrationObjectStorageDetails, withExternalChangesMarking bool) diag.Diagnostics {
	if withExternalChangesMarking {
		if err := handleExternalChangesToObjectInFlatDescribe(d, catalogIntegrationObjectStorageOutputMappingsExt(details)...); err != nil {
			return diag.FromErr(err)
		}
	}
	// refresh_interval_seconds is not copied into config; drift is marked on describe output.
	if err := errors.Join(
		d.Set("catalog_source", string(details.CatalogSource)),
		d.Set("table_format", string(details.TableFormat)),
	); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func catalogIntegrationObjectStorageOutputMappingsExt(details *sdk.CatalogIntegrationObjectStorageDetails) []outputMapping {
	return []outputMapping{
		{"refresh_interval_seconds", "refresh_interval_seconds", details.RefreshIntervalSeconds, details.RefreshIntervalSeconds, nil},
	}
}

func catalogIntegrationObjectStorageSetConfigFieldsExt(d *schema.ResourceData, catalogIntegration *sdk.CatalogIntegration) error {
	return errors.Join(
		d.Set("name", catalogIntegration.Name),
		d.Set("enabled", catalogIntegration.Enabled),
		d.Set("comment", catalogIntegration.Comment),
	)
}

func catalogIntegrationObjectStorageApplySetUnsetExt(d *schema.ResourceData) (*catalogIntegrationObjectStorageAlterRequestsExt, error) {
	set := sdk.NewCatalogIntegrationSetRequest()
	if errs := handleCatalogIntegrationCommonPropsUpdate(d, set); errs != nil {
		return nil, errs
	}
	return &catalogIntegrationObjectStorageAlterRequestsExt{set: set}, nil
}

func catalogIntegrationObjectStorageSetFieldsNotSetByReadExt(d *schema.ResourceData, _ sdk.AccountObjectIdentifier, details *sdk.CatalogIntegrationObjectStorageDetails, importPermissiveMode bool) error {
	return importIntAttribute(d, "refresh_interval_seconds", importPermissiveMode, details.RefreshIntervalSeconds)
}

// =============================================================================
// SDK incision points
// =============================================================================

func catalogIntegrationObjectStorageCreateInSdkExt(ctx context.Context, meta any, _ sdk.AccountObjectIdentifier, req *sdk.CreateObjectStorageCatalogIntegrationRequest) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	if err := client.CatalogIntegrations.CreateObjectStorage(ctx, req); err != nil {
		return diag.FromErr(fmt.Errorf("error creating object storage catalog integration, err = %w", err))
	}
	return nil
}

func catalogIntegrationObjectStorageShowByIdSafelyInSdkExt(ctx context.Context, d *schema.ResourceData, meta any, id sdk.AccountObjectIdentifier) (*sdk.CatalogIntegration, diag.Diagnostics) {
	client := meta.(*provider.Context).Client
	s, err := client.CatalogIntegrations.ShowByIDSafely(ctx, id)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			d.SetId("")
			return nil, diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to query object storage catalog integration. Marking the resource as removed.",
					Detail:   fmt.Sprintf("object storage catalog integration id: %s, Err: %s", id.FullyQualifiedName(), err),
				},
			}
		}
		return nil, diag.FromErr(err)
	}
	return s, nil
}

// Describe plus the object-store catalog-source guard. Import does not call ShowByID.
func catalogIntegrationObjectStorageShowByIdInSdkExt(ctx context.Context, meta any, id sdk.AccountObjectIdentifier) (*sdk.CatalogIntegrationObjectStorageDetails, error) {
	client := meta.(*provider.Context).Client
	details, err := client.CatalogIntegrations.DescribeObjectStorageDetails(ctx, id)
	if err != nil {
		return nil, err
	}
	if details.CatalogSource != sdk.CatalogIntegrationCatalogSourceTypeObjectStore {
		return nil, fmt.Errorf("invalid catalog source type, expected %s, got %s", sdk.CatalogIntegrationCatalogSourceTypeObjectStore, details.CatalogSource)
	}
	return details, nil
}

func catalogIntegrationObjectStorageDescribeInSdkExt(ctx context.Context, meta any, id sdk.AccountObjectIdentifier) (*sdk.CatalogIntegrationObjectStorageDetails, error) {
	client := meta.(*provider.Context).Client
	details, err := client.CatalogIntegrations.DescribeObjectStorageDetails(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("could not describe object storage catalog integration (%s), err = %w", helpers.EncodeResourceIdentifier(id), err)
	}
	return details, nil
}

func catalogIntegrationObjectStorageAlterSetInSdkExt(ctx context.Context, _ *schema.ResourceData, meta any, id sdk.AccountObjectIdentifier, reqs *catalogIntegrationObjectStorageAlterRequestsExt) error {
	if !reflect.DeepEqual(*reqs.set, *sdk.NewCatalogIntegrationSetRequest()) {
		client := meta.(*provider.Context).Client
		req := sdk.NewAlterCatalogIntegrationRequest(id).WithSet(*reqs.set)
		if err := client.CatalogIntegrations.Alter(ctx, req); err != nil {
			return fmt.Errorf("error updating catalog integration (%s), err = %w", helpers.EncodeResourceIdentifier(id), err)
		}
	}
	return nil
}

func catalogIntegrationObjectStorageAlterUnsetInSdkExt(context.Context, any, sdk.AccountObjectIdentifier, *catalogIntegrationObjectStorageAlterRequestsExt) error {
	return nil
}
