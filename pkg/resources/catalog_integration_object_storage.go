package resources

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var catalogIntegrationObjectStorageSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "Specifies the identifier (i.e. name) of the catalog integration; must be unique in your account.",
	},
	"enabled": {
		Type:     schema.TypeBool,
		Required: true,
		Description: joinWithSpace("Specifies whether the catalog integration is available for use for Iceberg tables.",
			"`true` allows users to create new Iceberg tables that reference this integration. Existing Iceberg tables that reference this integration function normally.",
			"`false` prevents users from creating new Iceberg tables that reference this integration. Existing Iceberg tables that reference this integration cannot access the catalog in the table definition."),
	},
	"refresh_interval_seconds": {
		Type:     schema.TypeInt,
		Optional: true,
		Description: joinWithSpace("Specifies the number of seconds to wait between attempts to poll the external Iceberg catalog for metadata updates for automated refresh.",
			"For Delta-based tables, specifies the number of seconds to wait between attempts to poll your external cloud storage for new metadata."),
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInDescribe("refresh_interval_seconds"),
		ValidateFunc:     validation.IntAtLeast(1),
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Default:     "",
		Description: "Specifies a comment for the catalog integration.",
	},
	"table_format": {
		Type:         schema.TypeString,
		Required:     true,
		ForceNew:     true,
		Description:  "Specifies the table format. Supported values: ICEBERG, DELTA.",
		ValidateFunc: validation.StringInSlice([]string{string(sdk.CatalogIntegrationTableFormatIceberg), string(sdk.CatalogIntegrationTableFormatDelta)}, true),
	},
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW CATALOG INTEGRATIONS` for the given catalog integration.",
		Elem: &schema.Resource{
			Schema: schemas.ShowCatalogIntegrationSchema,
		},
	},
	DescribeOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `DESCRIBE CATALOG INTEGRATION` for the given catalog integration.",
		Elem: &schema.Resource{
			Schema: schemas.DescribeCatalogIntegrationObjectStorageDetailsSchema,
		},
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}

func CatalogIntegrationObjectStorage() *schema.Resource {
	deleteFunc := ResourceDeleteContextFunc(
		sdk.ParseAccountObjectIdentifier,
		func(client *sdk.Client) DropSafelyFunc[sdk.AccountObjectIdentifier] {
			return client.CatalogIntegrations.DropSafely
		},
	)

	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.CatalogIntegrationObjectStorageResource), TrackingCreateWrapper(resources.CatalogIntegrationObjectStorage, CreateCatalogIntegrationObjectStorage)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.CatalogIntegrationObjectStorageResource), TrackingReadWrapper(resources.CatalogIntegrationObjectStorage, ReadCatalogIntegrationObjectStorageFunc(true))),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.CatalogIntegrationObjectStorageResource), TrackingUpdateWrapper(resources.CatalogIntegrationObjectStorage, UpdateCatalogIntegrationObjectStorage)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.CatalogIntegrationObjectStorageResource), TrackingDeleteWrapper(resources.CatalogIntegrationObjectStorage, deleteFunc)),
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

func ImportCatalogIntegrationObjectStorage(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
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

func CreateCatalogIntegrationObjectStorage(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	name := d.Get("name").(string)
	id := sdk.NewAccountObjectIdentifierFromFullyQualifiedName(name)
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

func UpdateCatalogIntegrationObjectStorage(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	set := sdk.NewCatalogIntegrationSetRequest()

	errs := errors.Join(
		booleanAttributeUpdateSetOnly(d, "enabled", &set.Enabled),
		// TODO [SNOW-3243983]: UNSET not implemented
		intAttributeUnsetFallbackUpdateWithZeroDefault(d, "refresh_interval_seconds", &set.RefreshIntervalSeconds, 30),
		stringAttributeUpdateSetOnly(d, "comment", &set.Comment),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}

	if !reflect.DeepEqual(*set, *sdk.NewCatalogIntegrationSetRequest()) {
		req := sdk.NewAlterCatalogIntegrationRequest(id).WithSet(*set)
		if err := client.CatalogIntegrations.Alter(ctx, req); err != nil {
			return diag.FromErr(fmt.Errorf("error updating object storage catalog integration (%s), err = %w", d.Id(), err))
		}
	}
	return ReadCatalogIntegrationObjectStorageFunc(false)(ctx, d, meta)
}
