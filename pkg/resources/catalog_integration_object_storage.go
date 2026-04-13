package resources

import (
	"context"
	"errors"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var catalogIntegrationObjectStorageSchema = func() map[string]*schema.Schema {
	objectStorageSchema := map[string]*schema.Schema{
		"table_format": {
			Type:             schema.TypeString,
			Required:         true,
			ForceNew:         true,
			Description:      "Specifies the table format. " + enumValuesDescription(sdk.AllCatalogIntegrationTableFormats),
			DiffSuppressFunc: NormalizeAndCompare(sdk.ToCatalogIntegrationTableFormat),
			ValidateDiagFunc: sdkValidation(sdk.ToCatalogIntegrationTableFormat),
		},
	}
	return collections.MergeMaps(
		catalogIntegrationCommonSchema(schemas.DescribeCatalogIntegrationObjectStorageDetailsSchema), objectStorageSchema)
}()

func CatalogIntegrationObjectStorage() *schema.Resource {
	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.CatalogIntegrationObjectStorageResource), TrackingCreateWrapper(resources.CatalogIntegrationObjectStorage, CreateCatalogIntegrationObjectStorage)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.CatalogIntegrationObjectStorageResource), TrackingReadWrapper(resources.CatalogIntegrationObjectStorage, ReadCatalogIntegrationObjectStorageFunc(true))),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.CatalogIntegrationObjectStorageResource), TrackingUpdateWrapper(resources.CatalogIntegrationObjectStorage, UpdateCatalogIntegrationObjectStorage)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.CatalogIntegrationObjectStorageResource), TrackingDeleteWrapper(resources.CatalogIntegrationObjectStorage, deleteCatalogIntegrationFunc())),
		Description:   "Resource used to manage catalog integration objects for Apache Iceberg™ table files or Delta table files in object storage. For more information, check [catalog integration documentation](https://docs.snowflake.com/en/sql-reference/sql/create-catalog-integration-object-storage).",

		Schema: catalogIntegrationObjectStorageSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.CatalogIntegrationObjectStorage, ImportCatalogIntegrationObjectStorage),
		},
		Timeouts: defaultTimeouts,
		CustomizeDiff: customdiff.All(
			ComputedIfAnyAttributeChanged(catalogIntegrationObjectStorageSchema, ShowOutputAttributeName, "enabled", "comment"),
			ComputedIfAnyAttributeChanged(catalogIntegrationObjectStorageSchema, DescribeOutputAttributeName, "enabled", "refresh_interval_seconds", "comment"),
		),
	}
}

func ImportCatalogIntegrationObjectStorage(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return nil, err
	}

	details, err := client.CatalogIntegrations.DescribeObjectStorageDetails(ctx, id)
	if err != nil {
		return nil, err
	}
	if details.CatalogSource != sdk.CatalogIntegrationCatalogSourceTypeObjectStorage {
		return nil, fmt.Errorf("invalid catalog source type, expected %s, got %s", sdk.CatalogIntegrationCatalogSourceTypeObjectStorage, details.CatalogSource)
	}

	return []*schema.ResourceData{d}, nil
}

func CreateCatalogIntegrationObjectStorage(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	name := d.Get("name").(string)
	id, err := sdk.ParseAccountObjectIdentifier(name)
	if err != nil {
		return diag.FromErr(err)
	}
	enabled := d.Get("enabled").(bool)
	tableFormat, err := sdk.ToCatalogIntegrationTableFormat(d.Get("table_format").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	request := sdk.NewCreateCatalogIntegrationRequest(id, enabled)
	objectStorageRequest := sdk.NewObjectStorageParamsRequest(tableFormat)
	errs := errors.Join(
		intAttributeCreateBuilder(d, "refresh_interval_seconds", request.WithRefreshIntervalSeconds),
		stringAttributeCreateBuilder(d, "comment", request.WithComment),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}

	if err := client.CatalogIntegrations.Create(ctx, request.WithObjectStorageCatalogSourceParams(*objectStorageRequest)); err != nil {
		return diag.FromErr(fmt.Errorf("error creating object storage catalog integration, err = %w", err))
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))
	return ReadCatalogIntegrationObjectStorageFunc(false)(ctx, d, meta)
}

func ReadCatalogIntegrationObjectStorageFunc(withExternalChangesMarking bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		client := meta.(*provider.Context).Client
		id, err := sdk.ParseAccountObjectIdentifier(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}

		s, err := client.CatalogIntegrations.ShowByIDSafely(ctx, id)
		if err != nil {
			if errors.Is(err, sdk.ErrObjectNotFound) {
				d.SetId("")
				return diag.Diagnostics{
					diag.Diagnostic{
						Severity: diag.Warning,
						Summary:  "Failed to query object storage catalog integration. Marking the resource as removed.",
						Detail:   fmt.Sprintf("object storage catalog integration id: %s, Err: %s", id.FullyQualifiedName(), err),
					},
				}
			}
			return diag.FromErr(err)
		}

		details, err := client.CatalogIntegrations.DescribeObjectStorageDetails(ctx, id)
		if err != nil {
			return diag.FromErr(fmt.Errorf("could not describe object storage catalog integration (%s), err = %w", d.Id(), err))
		}

		if withExternalChangesMarking {
			if err = handleExternalChangesToObjectInFlatDescribe(d,
				outputMapping{"refresh_interval_seconds", "refresh_interval_seconds", details.RefreshIntervalSeconds, details.RefreshIntervalSeconds, nil},
			); err != nil {
				return diag.FromErr(err)
			}
		}

		errs := errors.Join(
			d.Set("name", details.Id.Name()),
			d.Set("enabled", details.Enabled),
			// not reading refresh_interval_seconds on purpose (handled as external change to describe output)
			d.Set("comment", details.Comment),
			d.Set("table_format", string(details.TableFormat)),
			d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
			d.Set(ShowOutputAttributeName, []map[string]any{schemas.CatalogIntegrationToSchema(s)}),
			d.Set(DescribeOutputAttributeName, []map[string]any{schemas.CatalogIntegrationObjectStorageDetailsToSchema(details)}),
		)

		return diag.FromErr(errs)
	}
}

func UpdateCatalogIntegrationObjectStorage(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	if err := handleCatalogIntegrationUpdate(ctx, d, meta); err != nil {
		return diag.FromErr(err)
	}
	return ReadCatalogIntegrationObjectStorageFunc(false)(ctx, d, meta)
}
