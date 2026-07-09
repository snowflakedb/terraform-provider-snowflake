package resources

import (
	"context"
	"errors"
	"fmt"

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

// TODO (next PRs): the following CreateIcebergTableOptions fields are not yet supported by this resource:
//   - structure/layout: PartitionBy, ClusterBy
//   - attached policies: RowAccessPolicy, AggregationPolicy
//   - CopyGrants and CopyTags
//   - ICEBERG_MERGE_ON_READ_BEHAVIOR (needs to be added to SDK)
//   - column-level extras (part of ColumnsAndConstraints): out-of-line constraints, and per-column
//     DefaultValue, NotNull, InlineConstraint, MaskingPolicy, ProjectionPolicy, Comment...
//   - https://docs.snowflake.com/en/sql-reference/parameters#label-iceberg-default-ddl-collation
var icebergTableSchema = collections.MergeMaps(
	icebergTableCommonSchema(),
	map[string]*schema.Schema{
		"column": basicColumnSchema(),
		"base_location": {
			Type:             schema.TypeString,
			Optional:         true,
			ForceNew:         true,
			DiffSuppressFunc: suppressIcebergTableBaseLocationSuffix,
			Description:      "The path to a directory where Snowflake can write data and metadata files for the Iceberg table. Specify a relative path from the table's `EXTERNAL_VOLUME` location.",
		},
		"path_layout": icebergTablePathLayoutSchema(),
		"error_logging": {
			Type:             schema.TypeString,
			Optional:         true,
			Default:          BooleanDefault,
			ValidateDiagFunc: validateBooleanString,
			Description:      externalChangesNotDetectedFieldDescription(booleanStringFieldDescription("Specifies whether error logging is enabled for the Iceberg table.")),
		},
		"change_tracking": {
			Type:             schema.TypeString,
			Optional:         true,
			ForceNew:         true,
			Default:          BooleanDefault,
			ValidateDiagFunc: validateBooleanString,
			Description:      externalChangesNotDetectedFieldDescription(booleanStringFieldDescription("Specifies whether to enable change tracking on the Iceberg table. Cannot be changed after creation.")),
		},
		"iceberg_version": {
			Type:             schema.TypeInt,
			Optional:         true,
			ForceNew:         true,
			Description:      "Specifies the Iceberg table format version.",
			ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
		},
		ParametersAttributeName: {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Outputs the result of `SHOW PARAMETERS IN ICEBERG TABLE` for the given Iceberg table.",
			Elem:        &schema.Resource{Schema: schemas.ShowIcebergTableSnowflakeManagedParametersSchema},
		},
	},
	icebergTableSnowflakeManagedParametersSchema(),
)

func IcebergTable() *schema.Resource {
	return &schema.Resource{
		// TODO (next PRs): Add PreviewFeature*ContextWrapper when this resource is moved to the production provider.
		CreateContext: TrackingCreateWrapper(resources.IcebergTable, CreateIcebergTable),
		ReadContext:   TrackingReadWrapper(resources.IcebergTable, ReadIcebergTableFunc(true)),
		UpdateContext: TrackingUpdateWrapper(resources.IcebergTable, UpdateIcebergTable),
		DeleteContext: TrackingDeleteWrapper(resources.IcebergTable, icebergTableDeleteFunc()),

		Description: "Resource used to manage a Snowflake-managed Iceberg table. For more information, check [the official documentation](https://docs.snowflake.com/en/sql-reference/sql/create-iceberg-table-snowflake).",

		Schema: icebergTableSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.IcebergTable, importIcebergTable),
		},
		Timeouts: defaultTimeouts,
		CustomizeDiff: customdiff.All(
			ComputedIfAnyAttributeChanged(icebergTableSchema, ShowOutputAttributeName, "comment"),
			ComputedIfAnyAttributeChanged(icebergTableSchema, DescribeOutputAttributeName, "column"),
			ComputedIfAnyAttributeChanged(icebergTableSchema, ParametersAttributeName,
				"external_volume", "catalog", "target_file_size", "storage_serialization_policy",
				"catalog_sync", "data_retention_time_in_days", "max_data_extension_time_in_days", "enable_data_compaction",
				"enable_iceberg_merge_on_read",
			),
			icebergTableSnowflakeManagedParametersCustomDiff,
		),
	}
}

func CreateIcebergTable(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)
	name := d.Get("name").(string)
	id := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, name)

	columns, err := parseBasicColumns(d.Get("column").([]any))
	if err != nil {
		return diag.FromErr(err)
	}
	columnsAndConstraints := *sdk.NewIcebergTableColumnsAndConstraintsRequest().WithColumns(toIcebergTableColumnRequests(columns))
	req := sdk.NewCreateIcebergTableRequest(id, columnsAndConstraints)

	if err := errors.Join(
		stringAttributeCreate(d, "comment", &req.Comment),
		stringAttributeCreate(d, "base_location", &req.BaseLocation),
		intAttributeCreate(d, "iceberg_version", &req.IcebergVersion),
		booleanStringAttributeCreate(d, "change_tracking", &req.ChangeTracking),
		booleanStringAttributeCreate(d, "error_logging", &req.ErrorLogging),
		attributeMappedValueCreate(d, "path_layout", &req.PathLayout, func(value any) (*sdk.IcebergTablePathLayout, error) {
			pathLayout, err := sdk.ToIcebergTablePathLayout(value.(string))
			if err != nil {
				return nil, err
			}
			return &pathLayout, nil
		}),
	); err != nil {
		return diag.FromErr(err)
	}

	diags := handleIcebergTableSnowflakeManagedParametersCreate(d, req)
	if diags.HasError() {
		return diags
	}

	if err := client.IcebergTables.Create(ctx, req); err != nil {
		return diag.FromErr(fmt.Errorf("error creating Iceberg table (%s), err = %w", id.FullyQualifiedName(), err))
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))
	return ReadIcebergTableFunc(false)(ctx, d, meta)
}

func ReadIcebergTableFunc(withExternalChangesMarking bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		return readIcebergTableWithParameterHandler(ctx, d, meta, handleIcebergTableSnowflakeManagedParameterRead, schemas.IcebergTableSnowflakeManagedParametersToSchema, func(d *schema.ResourceData, table *sdk.IcebergTable, details []sdk.IcebergTableDetails) error {
			if withExternalChangesMarking {
				var baseLocation string
				if table.BaseLocation != nil {
					baseLocation = NormalizeIcebergTableBaseLocation(*table.BaseLocation)
				}

				if err := handleExternalChangesToObjectInShow(
					d,
					outputMapping{"iceberg_table_format_version", "iceberg_version", table.IcebergTableFormatVersion, table.IcebergTableFormatVersion, nil},
					outputMapping{"base_location", "base_location", baseLocation, baseLocation, func(value any) any {
						return NormalizeIcebergTableBaseLocation(value.(string))
					}},
				); err != nil {
					return err
				}
			}

			// path_layout is not exposed by SHOW or DESCRIBE, so it is not read back (external changes are not detected).
			return errors.Join(
				d.Set("column", icebergTableColumnsToSchema(details)),
			)
		})
	}
}

func UpdateIcebergTable(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// TODO (next PRs): columns are ForceNew for now; handle the update properly
	// TODO (next PRs): comment needs to be altered separately - report this

	set := sdk.NewIcebergTableSetPropertiesRequest()
	unset := sdk.NewIcebergTableUnsetPropertiesRequest()
	if errs := errors.Join(
		stringAttributeUpdate(d, "comment", &set.Comment, &unset.Comment),
	); errs != nil {
		return diag.FromErr(errs)
	}
	if err := applyIcebergTableAlter(ctx, client, id, set, unset); err != nil {
		return diag.FromErr(err)
	}

	set = sdk.NewIcebergTableSetPropertiesRequest()
	unset = sdk.NewIcebergTableUnsetPropertiesRequest()
	if errs := errors.Join(
		booleanStringAttributeUpdate(d, "error_logging", &set.ErrorLogging, &unset.ErrorLogging),
	); errs != nil {
		return diag.FromErr(errs)
	}
	if diags := handleIcebergTableSnowflakeManagedParametersUpdate(d, set, unset); diags.HasError() {
		return diags
	}
	if err := applyIcebergTableAlter(ctx, client, id, set, unset); err != nil {
		return diag.FromErr(err)
	}

	return ReadIcebergTableFunc(false)(ctx, d, meta)
}

func handleIcebergTableSnowflakeManagedParametersCreate(d *schema.ResourceData, req *sdk.CreateIcebergTableRequest) diag.Diagnostics {
	if diags := JoinDiags(
		handleParameterCreateWithMapping(d, sdk.IcebergTableParameterExternalVolume, &req.ExternalVolume, sdk.ParseAccountObjectIdentifier),
		handleParameterCreateWithMapping(d, sdk.IcebergTableParameterCatalog, &req.Catalog, stringToStringEnumProvider(sdk.ToIcebergTableCatalog)),
		handleParameterCreateWithMapping(d, sdk.IcebergTableParameterTargetFileSize, &req.TargetFileSize, stringToStringEnumProvider(sdk.ToIcebergTableTargetFileSize)),
		handleParameterCreateWithMapping(d, sdk.IcebergTableParameterStorageSerializationPolicy, &req.StorageSerializationPolicy, stringToStringEnumProvider(sdk.ToStorageSerializationPolicy)),
		handleParameterCreate(d, sdk.IcebergTableParameterCatalogSync, &req.CatalogSync),
		handleParameterCreate(d, sdk.IcebergTableParameterDataRetentionTimeInDays, &req.DataRetentionTimeInDays),
		handleParameterCreate(d, sdk.IcebergTableParameterMaxDataExtensionTimeInDays, &req.MaxDataExtensionTimeInDays),
		handleParameterCreate(d, sdk.IcebergTableParameterEnableDataCompaction, &req.EnableDataCompaction),
		handleParameterCreate(d, sdk.IcebergTableParameterEnableIcebergMergeOnRead, &req.EnableIcebergMergeOnRead),
	); diags.HasError() {
		return diags
	}

	return nil
}

// handleIcebergTableParametersUpdate populates the set/unset requests for all alterable Iceberg table parameters.
// storage_serialization_policy is intentionally omitted: it is create-only (ForceNew) and cannot be altered.
func handleIcebergTableSnowflakeManagedParametersUpdate(d *schema.ResourceData, set *sdk.IcebergTableSetPropertiesRequest, unset *sdk.IcebergTableUnsetPropertiesRequest) diag.Diagnostics {
	return JoinDiags(
		handleParameterUpdate(d, sdk.IcebergTableParameterCatalogSync, &set.CatalogSync, &unset.CatalogSync),
		handleParameterUpdate(d, sdk.IcebergTableParameterDataRetentionTimeInDays, &set.DataRetentionTimeInDays, &unset.DataRetentionTimeInDays),
		handleParameterUpdate(d, sdk.IcebergTableParameterMaxDataExtensionTimeInDays, &set.MaxDataExtensionTimeInDays, &unset.MaxDataExtensionTimeInDays),
		handleParameterUpdate(d, sdk.IcebergTableParameterEnableDataCompaction, &set.EnableDataCompaction, &unset.EnableDataCompaction),
		handleParameterUpdate(d, sdk.IcebergTableParameterEnableIcebergMergeOnRead, &set.EnableIcebergMergeOnRead, &unset.EnableIcebergMergeOnRead),
		handleParameterUpdateWithMapping(d, sdk.IcebergTableParameterTargetFileSize, &set.TargetFileSize, &unset.TargetFileSize, stringToStringEnumProvider(sdk.ToIcebergTableTargetFileSize)),
	)
}

func toIcebergTableColumnRequests(columns []basicColumn) []sdk.IcebergTableColumnRequest {
	return collections.Map(columns, func(c basicColumn) sdk.IcebergTableColumnRequest {
		return *sdk.NewIcebergTableColumnRequest(c.Name, c.DataType)
	})
}

func icebergTableColumnsToSchema(details []sdk.IcebergTableDetails) []map[string]any {
	return collections.Map(details, func(d sdk.IcebergTableDetails) map[string]any {
		return map[string]any{
			"name": d.Name,
			"type": d.Type.ToSql(),
		}
	})
}
