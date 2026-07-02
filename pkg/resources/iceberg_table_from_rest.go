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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var icebergTableFromRestSchema = collections.MergeMaps(
	icebergTableCommonSchema(),
	map[string]*schema.Schema{
		"catalog_table_name": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  "Specifies the name of the table as it appears in the external catalog.",
		},
		"catalog_namespace": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "Specifies the namespace (or database) in the external catalog that the table belongs to. If not specified, the catalog integration's default namespace is used.",
		},
		"path_layout": {
			Type:             schema.TypeString,
			Optional:         true,
			ForceNew:         true,
			ValidateDiagFunc: StringInSlice(sdk.AsStringList(sdk.AllIcebergTablePathLayouts), true),
			Description:      externalChangesNotDetectedFieldDescription(fmt.Sprintf("Specifies the storage layout for the Iceberg table's Parquet files. Valid values are: %v. Cannot be changed after creation.", sdk.AllIcebergTablePathLayouts)),
		},
		"auto_refresh": {
			Type:             schema.TypeString,
			Optional:         true,
			ValidateDiagFunc: validateBooleanString,
			Default:          BooleanDefault,
			Description:      booleanStringFieldDescription("Specifies whether Snowflake should periodically refresh the Iceberg table metadata from the external catalog."),
			DiffSuppressFunc: SuppressIfAny(
				IgnoreChangeToCurrentSnowflakePlainValueInOutputWithMapping("show_output", "auto_refresh_status", func(x any) any {
					return len(x.([]any)) != 0
				}),
			),
		},
		// Override the shared parameters output to expose the additional REST-specific parameters.
		ParametersAttributeName: {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Outputs the result of `SHOW PARAMETERS IN ICEBERG TABLE` for the given Iceberg table.",
			Elem:        &schema.Resource{Schema: schemas.ShowIcebergTableFromRestParametersSchema},
		},
	},
	icebergTableFromRestParametersSchema(),
)

func IcebergTableFromRest() *schema.Resource {
	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.IcebergTableFromRestResource), TrackingCreateWrapper(resources.IcebergTableFromRest, CreateIcebergTableFromRest)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.IcebergTableFromRestResource), TrackingReadWrapper(resources.IcebergTableFromRest, ReadIcebergTableFromRestFunc(true))),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.IcebergTableFromRestResource), TrackingUpdateWrapper(resources.IcebergTableFromRest, UpdateIcebergTableFromRest)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.IcebergTableFromRestResource), TrackingDeleteWrapper(resources.IcebergTableFromRest, icebergTableDeleteFunc())),

		Description: "Resource used to manage an Iceberg table whose metadata is managed by an external Iceberg REST catalog. For more information, check [the official documentation](https://docs.snowflake.com/en/sql-reference/sql/create-iceberg-table-rest).",

		Schema: icebergTableFromRestSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.IcebergTableFromRest, importIcebergTableFromRest),
		},
		Timeouts: defaultTimeouts,
		CustomizeDiff: customdiff.All(
			ComputedIfAnyAttributeChanged(icebergTableFromRestSchema, ShowOutputAttributeName, "comment", "auto_refresh", "catalog_table_name", "catalog_namespace"),
			ComputedIfAnyAttributeChanged(icebergTableFromRestSchema, ParametersAttributeName, "external_volume", "catalog", "replace_invalid_characters", "target_file_size", "storage_serialization_policy", "enable_iceberg_merge_on_read", "iceberg_merge_on_read_behavior"),
			icebergTableFromRestParametersCustomDiff,
		),
	}
}

func importIcebergTableFromRest(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return nil, err
	}

	table, err := client.IcebergTables.ShowByIDSafely(ctx, id)
	if err != nil {
		return nil, err
	}

	if _, err := ImportName[sdk.SchemaObjectIdentifier](ctx, d, nil); err != nil {
		return nil, err
	}

	if err := d.Set("auto_refresh", booleanStringFromBool(table.AutoRefreshStatus != nil)); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

func CreateIcebergTableFromRest(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)
	name := d.Get("name").(string)
	catalogTableName := d.Get("catalog_table_name").(string)

	id := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, name)
	req := sdk.NewCreateFromIcebergRestIcebergTableRequest(id, catalogTableName)

	if err := stringAttributeCreate(d, "comment", &req.Comment); err != nil {
		return diag.FromErr(err)
	}
	if err := stringAttributeCreate(d, "catalog_namespace", &req.CatalogNamespace); err != nil {
		return diag.FromErr(err)
	}
	if err := attributeMappedValueCreate(d, "path_layout", &req.PathLayout, func(value any) (*sdk.IcebergTablePathLayout, error) {
		pathLayout, err := sdk.ToIcebergTablePathLayout(value.(string))
		if err != nil {
			return nil, err
		}
		return &pathLayout, nil
	}); err != nil {
		return diag.FromErr(err)
	}
	if err := booleanStringAttributeCreate(d, "auto_refresh", &req.AutoRefresh); err != nil {
		return diag.FromErr(err)
	}
	if diags := handleIcebergTableParametersCreate(d, &req.ExternalVolume, &req.Catalog, &req.ReplaceInvalidCharacters); diags.HasError() {
		return diags
	}
	if diags := JoinDiags(
		handleParameterCreateWithMapping(d, sdk.IcebergTableParameterTargetFileSize, &req.TargetFileSize, stringToStringEnumProvider(sdk.ToIcebergTableTargetFileSize)),
		handleParameterCreateWithMapping(d, sdk.IcebergTableParameterStorageSerializationPolicy, &req.StorageSerializationPolicy, stringToStringEnumProvider(sdk.ToStorageSerializationPolicy)),
		handleParameterCreateWithMapping(d, sdk.IcebergTableParameterIcebergMergeOnReadBehavior, &req.IcebergMergeOnReadBehavior, stringToStringEnumProvider(sdk.ToIcebergTableIcebergMergeOnReadBehavior)),
		handleParameterCreate(d, sdk.IcebergTableParameterEnableIcebergMergeOnRead, &req.EnableIcebergMergeOnRead),
	); diags.HasError() {
		return diags
	}

	if err := client.IcebergTables.CreateFromIcebergRest(ctx, req); err != nil {
		return diag.FromErr(fmt.Errorf("error creating Iceberg table from REST catalog (%s), err = %w", id.FullyQualifiedName(), err))
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))
	return ReadIcebergTableFromRestFunc(false)(ctx, d, meta)
}

func ReadIcebergTableFromRestFunc(withExternalChangesMarking bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		// path_layout is not exposed by SHOW or DESCRIBE, so it is not read back (external changes are not detected).
		return readIcebergTableWithParameterHandler(ctx, d, meta, handleIcebergTableFromRestParameterRead, schemas.IcebergTableFromRestParametersToSchema, func(d *schema.ResourceData, table *sdk.IcebergTable) error {
			var catalogTableName string
			if table.CatalogTableName != nil {
				catalogTableName = *table.CatalogTableName
			}
			if withExternalChangesMarking {
				var autoRefreshSet string
				if table.AutoRefreshStatus != nil {
					autoRefreshSet = table.AutoRefreshStatus.ExecutionState
				}
				if err := handleExternalChangesToObject(
					d,
					"show_output.0.auto_refresh_status",
					outputMapping{"execution_state", "auto_refresh", autoRefreshSet, booleanStringFromBool(table.AutoRefreshStatus != nil), nil},
				); err != nil {
					return err
				}
				var catalogNamespace string
				if table.CatalogNamespace != nil {
					catalogNamespace = *table.CatalogNamespace
				}
				if err := handleExternalChangesToObjectInShow(
					d,
					outputMapping{"catalog_namespace", "catalog_namespace", catalogNamespace, catalogNamespace, nil},
				); err != nil {
					return err
				}
			}
			return errors.Join(
				d.Set("catalog_table_name", catalogTableName),
			)
		})
	}
}

func UpdateIcebergTableFromRest(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// TODO(SNOW-3735539): Altering IcebergTableParameterReplaceInvalidCharacters and IcebergTableParameterTargetFileSize with comment at the same time
	// does not cause any changes. Therefore, changes in these parameters are extracted as a separate call.
	// After this is fixed in Snowflake, we should have one SET call for all changes.
	set := sdk.NewIcebergTableSetPropertiesRequest()
	unset := sdk.NewIcebergTableUnsetPropertiesRequest()
	if errs := errors.Join(
		stringAttributeUpdate(d, "comment", &set.Comment, &unset.Comment),
	); errs != nil {
		return diag.FromErr(errs)
	}
	if err := booleanStringAttributeUnsetFallbackUpdate(d, "auto_refresh", &set.AutoRefresh, false); err != nil {
		return diag.FromErr(err)
	}

	if err := applyIcebergTableAlter(ctx, client, id, set, unset); err != nil {
		return diag.FromErr(err)
	}

	set = sdk.NewIcebergTableSetPropertiesRequest()
	unset = sdk.NewIcebergTableUnsetPropertiesRequest()
	if diags := JoinDiags(
		handleParameterUpdateWithMapping(d, sdk.IcebergTableParameterTargetFileSize, &set.TargetFileSize, &unset.TargetFileSize, stringToStringEnumProvider(sdk.ToIcebergTableTargetFileSize)),
		handleParameterUpdate(d, sdk.IcebergTableParameterReplaceInvalidCharacters, &set.ReplaceInvalidCharacters, &unset.ReplaceInvalidCharacters),
		handleParameterUpdate(d, sdk.IcebergTableParameterEnableIcebergMergeOnRead, &set.EnableIcebergMergeOnRead, &unset.EnableIcebergMergeOnRead),
	); diags.HasError() {
		return diags
	}

	if err := applyIcebergTableAlter(ctx, client, id, set, unset); err != nil {
		return diag.FromErr(err)
	}

	return ReadIcebergTableFromRestFunc(false)(ctx, d, meta)
}
