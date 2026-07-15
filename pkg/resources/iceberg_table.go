package resources

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// TODO (next PRs): the following CreateIcebergTableOptions fields are not yet supported by this resource:
//   - CopyGrants and CopyTags
//   - Use ALTER TABLE for handling column changes
//   - ICEBERG_MERGE_ON_READ_BEHAVIOR (needs to be added to SDK)
//   - https://docs.snowflake.com/en/sql-reference/parameters#label-iceberg-default-ddl-collation
var icebergTableSchema = collections.MergeMaps(
	icebergTableCommonSchema(),
	map[string]*schema.Schema{
		"column":                 columnSchema(),
		"primary_key_constraint": primaryKeyConstraintSchema(),
		"unique_constraint":      uniqueConstraintSchema(),
		"foreign_key_constraint": foreignKeyConstraintSchema(),
		"check_constraint":       checkConstraintSchema(),
		"row_access_policy":      rowAccessPolicyFieldSchema("Iceberg table"),
		"aggregation_policy":     aggregationPolicySchema("Iceberg table"),
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
		"partition_by": icebergTablePartitionBySchema(),
		"cluster_by": {
			Type:          schema.TypeList,
			Elem:          &schema.Schema{Type: schema.TypeString},
			Optional:      true,
			ConflictsWith: []string{"partition_by"},
			Description:   externalChangesNotDetectedFieldDescription("A list of one or more table columns/expressions to be used as clustering key(s) for the table."),
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
			// ComputedIf is missing on purpose - diff suppression is not enough to avoid the output field being marked as computed.
			ComputedIfAnyAttributeChanged(
				icebergTableSchema, ParametersAttributeName,
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

	columns, err := parseIcebergTableColumns(d)
	if err != nil {
		return diag.FromErr(err)
	}
	outOfLineConstraints, err := parseOutOfLineConstraints(d)
	if err != nil {
		return diag.FromErr(err)
	}
	columnsAndConstraints := *sdk.NewIcebergTableColumnsAndConstraintsRequest(columns)
	if len(outOfLineConstraints) > 0 {
		columnsAndConstraints.WithOutOfLineConstraint(outOfLineConstraints)
	}
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
		attributeMappedValueCreateBuilderNested(d, "partition_by", req.WithPartitionBy, parseIcebergTablePartitionBy),
		attributeMappedValueCreateBuilder(d, "cluster_by", req.WithClusterBy, func(value []any) ([]string, error) {
			return expandStringList(value), nil
		}),
	); err != nil {
		return diag.FromErr(err)
	}

	if v := d.Get("row_access_policy"); len(v.([]any)) > 0 {
		policyId, columns, err := extractPolicyWithColumnsSet(v, "on")
		if err != nil {
			return diag.FromErr(err)
		}
		req.WithRowAccessPolicy(*sdk.NewIcebergTableRowAccessPolicyRequest(policyId, columns))
	}

	if v := d.Get("aggregation_policy"); len(v.([]any)) > 0 {
		id, columns, err := extractPolicyWithColumnsSet(v, "entity_key")
		if err != nil {
			return diag.FromErr(err)
		}
		aggregationPolicyReq := sdk.NewIcebergTableAggregationPolicyRequest(id)
		if len(columns) > 0 {
			aggregationPolicyReq.WithEntityKey(columns)
		}
		req.WithAggregationPolicy(*aggregationPolicyReq)
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

// parseIcebergTableColumns parses the "column" list from the resource data into IcebergTableColumnRequests.
func parseIcebergTableColumns(d *schema.ResourceData) ([]sdk.IcebergTableColumnRequest, error) {
	raw := d.Get("column").([]any)
	indices := make([]int, len(raw))
	for i := range indices {
		indices[i] = i
	}
	return collections.MapErr(indices, func(i int) (sdk.IcebergTableColumnRequest, error) {
		return parseIcebergTableColumn(d, i)
	})
}

// parseIcebergTableColumn parses a single column at the given index (e.g. "column.0.") into an IcebergTableColumnRequest.
func parseIcebergTableColumn(d *schema.ResourceData, index int) (sdk.IcebergTableColumnRequest, error) {
	prefix := fmt.Sprintf("column.%d.", index)
	name := d.Get(prefix + "name").(string)
	dataType, err := datatypes.ParseDataType(d.Get(prefix + "type").(string))
	if err != nil {
		return sdk.IcebergTableColumnRequest{}, fmt.Errorf("parsing data type of column %q: %w", name, err)
	}
	req := sdk.NewIcebergTableColumnRequest(name, dataType)
	if err := booleanStringAttributeCreate(d, prefix+"not_null", &req.NotNull); err != nil {
		return sdk.IcebergTableColumnRequest{}, fmt.Errorf("parsing not_null for column %q: %w", name, err)
	}

	if err := errors.Join(
		stringAttributeCreateBuilder(d, prefix+"comment", func(v string) *sdk.IcebergTableColumnRequest { return req.WithComment(v) }),
		attributeMappedValueCreateBuilderNested(d, prefix+"default", func(v sdk.ColumnDefaultValue) *sdk.IcebergTableColumnRequest {
			return req.WithDefaultValue(v)
		}, func(d *schema.ResourceData) (sdk.ColumnDefaultValue, error) {
			return parseIcebergColumnDefaultValue(d, index)
		}),
		attributeMappedValueCreateBuilderNested(d, prefix+"masking_policy", func(v sdk.TableColumnMaskingPolicyRequest) *sdk.IcebergTableColumnRequest {
			return req.WithMaskingPolicy(v)
		}, func(d *schema.ResourceData) (sdk.TableColumnMaskingPolicyRequest, error) {
			return parseColumnMaskingPolicy(d, prefix+"masking_policy.0.")
		}),
		attributeMappedValueCreateBuilderNested(d, prefix+"projection_policy", func(v sdk.TableColumnProjectionPolicyRequest) *sdk.IcebergTableColumnRequest {
			return req.WithProjectionPolicy(v)
		}, func(d *schema.ResourceData) (sdk.TableColumnProjectionPolicyRequest, error) {
			return parseColumnProjectionPolicy(d, prefix+"projection_policy.0.")
		}),
	); err != nil {
		return sdk.IcebergTableColumnRequest{}, fmt.Errorf("parsing column %q: %w", name, err)
	}

	return *req, nil
}

// parseIcebergColumnDefaultValue reads the default expression directly from the raw config (rather than
// d.GetOk) so that an explicitly configured empty-string expression is not mistaken for an unset one.
func parseIcebergColumnDefaultValue(d *schema.ResourceData, index int) (sdk.ColumnDefaultValue, error) {
	defaultValue := sdk.ColumnDefaultValue{}
	path := cty.GetAttrPath("column").IndexInt(index).GetAttr("default").IndexInt(0).GetAttr("expression")
	configValue, diags := d.GetRawConfigAt(path)
	if diags.HasError() {
		return defaultValue, fmt.Errorf("reading raw config for column %d default expression: %v", index, diags)
	}
	if !configValue.IsNull() {
		expression := configValue.AsString()
		defaultValue.Expression = &expression
	}
	return defaultValue, nil
}

func ReadIcebergTableFunc(withExternalChangesMarking bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		client := meta.(*provider.Context).Client
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

			id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
			if err != nil {
				return err
			}
			policyRefs, err := client.PolicyReferences.GetForEntity(ctx, sdk.NewGetForEntityPolicyReferenceRequest(id, sdk.PolicyEntityDomainTable))
			if err != nil {
				return err
			}
			if err := handlePolicyReferences(policyRefs, d); err != nil {
				return err
			}

			var partitionBy []map[string]any
			if currentPartitionSpec, err := collections.FindFirst(table.PartitionSpecs, func(spec sdk.IcebergTablePartitionSpec) bool {
				return spec.SpecId == table.CurrentPartitionSpecId
			}); err == nil {
				partitionBy, err = icebergTablePartitionSpecFieldsToSchema(currentPartitionSpec.Fields)
				if err != nil {
					return fmt.Errorf("could not parse partition spec fields for Iceberg table (%s): %w", id.FullyQualifiedName(), err)
				}
			}

			// TODO (next PRs):
			// path_layout, error_logging and change_tracking are not exposed by SHOW or DESCRIBE, so they are not read back (external changes are not detected).
			// cluster_by is not read back either. SHOW/DESCRIBE ICEBERG TABLE do not expose the
			// clustering key, and even for regular tables Snowflake returns a transformed clustering expression
			// rather than the original DDL text, so external changes cannot be reliably detected. See
			// https://docs.snowflake.com/en/user-guide/tables-clustering-keys#defining-a-clustering-key-for-a-table
			// add these limitations to the documentation and report this to Snowflake
			return errors.Join(
				d.Set("column", handleIcebergTableColumns(details, policyRefs)),
				d.Set("partition_by", partitionBy),
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

	if d.HasChange("cluster_by") {
		clusterBy := expandStringList(d.Get("cluster_by").([]any))
		alterReq := sdk.NewAlterIcebergTableRequest(id)
		if len(clusterBy) > 0 {
			alterReq.WithClusteringAction(*sdk.NewIcebergTableClusteringActionRequest().WithClusterBy(clusterBy))
		} else {
			alterReq.WithClusteringAction(*sdk.NewIcebergTableClusteringActionRequest().WithDropClusteringKey(true))
		}
		if err := client.IcebergTables.Alter(ctx, alterReq); err != nil {
			return diag.FromErr(fmt.Errorf("error updating cluster_by on %v: %w", d.Id(), err))
		}
	}

	if d.HasChange("row_access_policy") {
		var addReq *sdk.ViewAddRowAccessPolicyRequest
		var dropReq *sdk.ViewDropRowAccessPolicyRequest
		err := rowAccessPolicyAlterRequests(d, func(id sdk.SchemaObjectIdentifier, columns []sdk.Column) {
			addReq = sdk.NewViewAddRowAccessPolicyRequest(id, columns)
		}, func(id sdk.SchemaObjectIdentifier) {
			dropReq = sdk.NewViewDropRowAccessPolicyRequest(id)
		})
		if err != nil {
			return diag.FromErr(err)
		}
		alterReq := sdk.NewAlterIcebergTableRequest(id)
		switch {
		case addReq != nil && dropReq != nil:
			alterReq.WithDropAndAddRowAccessPolicy(*sdk.NewIcebergTableDropAndAddRowAccessPolicyRequest(
				*sdk.NewIcebergTableDropRowAccessPolicyRequest(dropReq.RowAccessPolicy),
				*sdk.NewIcebergTableAddRowAccessPolicyRequest(addReq.RowAccessPolicy, addReq.On),
			))
		case addReq != nil:
			alterReq.WithAddRowAccessPolicy(*sdk.NewViewAddRowAccessPolicyRequest(addReq.RowAccessPolicy, addReq.On))
		case dropReq != nil:
			alterReq.WithDropRowAccessPolicy(*sdk.NewViewDropRowAccessPolicyRequest(dropReq.RowAccessPolicy))
		}
		if err := client.IcebergTables.Alter(ctx, alterReq); err != nil {
			return diag.FromErr(fmt.Errorf("error updating row access policy on %v err = %w", d.Id(), err))
		}
	}

	if d.HasChange("aggregation_policy") {
		newId, newColumns, isSet, err := aggregationPolicyAlterState(d)
		if err != nil {
			return diag.FromErr(err)
		}
		if isSet {
			aggregationPolicyReq := sdk.NewViewSetAggregationPolicyRequest(newId)
			if len(newColumns) > 0 {
				aggregationPolicyReq.WithEntityKey(newColumns)
			}
			if err := client.IcebergTables.Alter(ctx, sdk.NewAlterIcebergTableRequest(id).WithSetAggregationPolicy(*aggregationPolicyReq.WithForce(true))); err != nil {
				return diag.FromErr(fmt.Errorf("error setting aggregation policy for iceberg table %v: %w", d.Id(), err))
			}
		} else {
			if err := client.IcebergTables.Alter(ctx, sdk.NewAlterIcebergTableRequest(id).WithUnsetAggregationPolicy(*sdk.NewViewUnsetAggregationPolicyRequest())); err != nil {
				return diag.FromErr(fmt.Errorf("error unsetting aggregation policy for iceberg table %v: %w", d.Id(), err))
			}
		}
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

func handleIcebergTableColumns(columns []sdk.IcebergTableDetails, policyRefs []sdk.PolicyReference) []map[string]any {
	if len(columns) == 0 {
		return nil
	}

	return collections.Map(columns, func(column sdk.IcebergTableDetails) map[string]any {
		columnState := map[string]any{
			"name":     column.Name,
			"type":     column.Type.ToSql(),
			"not_null": booleanStringFromBool(!column.IsNullable),
		}
		if column.Comment != nil {
			columnState["comment"] = *column.Comment
		}
		if column.Default != nil {
			columnState["default"] = []map[string]any{{"expression": *column.Default}}
		}
		columnPoliciesState, err := columnPoliciesToState(column.Name, policyRefs)
		if err != nil {
			log.Printf("[DEBUG] could not convert column policies to state for column %q: %v", column.Name, err)
		}
		return collections.MergeMaps(columnState, columnPoliciesState)
	})
}

// icebergTablePartitionKinds lists the top-level fields of a single partition_by block; exactly one must be set.
var icebergTablePartitionKinds = []string{"identity", "bucket", "truncate", "year", "month", "day", "hour"}

func icebergTablePartitionBySchema() *schema.Schema {
	return &schema.Schema{
		Type:          schema.TypeList,
		Optional:      true,
		ForceNew:      true,
		ConflictsWith: []string{"cluster_by"},
		Description:   "Defines the partitioning for the Iceberg table. Cannot be changed after creation. Exactly one of identity, bucket, truncate, year, month, day, or hour must be set for each entry. Cannot be used together with `cluster_by`.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"identity": {
					Type:        schema.TypeString,
					Optional:    true,
					ForceNew:    true,
					Description: "Name of the column to use as-is for partitioning.",
				},
				"bucket": {
					Type:        schema.TypeList,
					Optional:    true,
					ForceNew:    true,
					MaxItems:    1,
					Description: "Partitions the table by hashing the column into a fixed number of buckets.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"num_buckets": {Type: schema.TypeInt, Required: true, ForceNew: true, Description: "Number of buckets to hash the column values into."},
							"column":      {Type: schema.TypeString, Required: true, ForceNew: true, Description: "Name of the column to bucket."},
						},
					},
				},
				"truncate": {
					Type:        schema.TypeList,
					Optional:    true,
					ForceNew:    true,
					MaxItems:    1,
					Description: "Partitions the table by truncating the column value to a fixed width.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"width":  {Type: schema.TypeInt, Required: true, ForceNew: true, Description: "Width to truncate the column value to."},
							"column": {Type: schema.TypeString, Required: true, ForceNew: true, Description: "Name of the column to truncate."},
						},
					},
				},
				"year":  icebergTablePartitionTimeSchema("Partitions the table by the year component of the column."),
				"month": icebergTablePartitionTimeSchema("Partitions the table by the month component of the column."),
				"day":   icebergTablePartitionTimeSchema("Partitions the table by the day component of the column."),
				"hour":  icebergTablePartitionTimeSchema("Partitions the table by the hour component of the column."),
			},
		},
	}
}

func icebergTablePartitionTimeSchema(description string) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		ForceNew:    true,
		Description: description,
	}
}

// parseIcebergTablePartitionBy parses the partition_by blocks from the resource data to SDK objects.
func parseIcebergTablePartitionBy(d *schema.ResourceData) ([]sdk.IcebergTablePartitionExpressionRequest, error) {
	entries := d.Get("partition_by").([]any)
	requests := make([]sdk.IcebergTablePartitionExpressionRequest, len(entries))

	for i := range entries {
		prefix := fmt.Sprintf("partition_by.%d.", i)
		req := sdk.NewIcebergTablePartitionExpressionRequest()

		err := errors.Join(
			attributeMappedValueCreateBuilderNested(d, prefix+"identity", req.WithIdentity, func(d *schema.ResourceData) (string, error) {
				return d.Get(prefix + "identity").(string), nil
			}),
			attributeMappedValueCreateBuilderNested(d, prefix+"bucket", req.WithBucket, parseIcebergTablePartitionBucket(prefix+"bucket.0.")),
			attributeMappedValueCreateBuilderNested(d, prefix+"truncate", req.WithTruncate, parseIcebergTablePartitionTruncate(prefix+"truncate.0.")),
			attributeMappedValueCreateBuilderNested(d, prefix+"year", req.WithYear, parseIcebergTablePartitionTime(prefix+"year.0.", func(args sdk.IcebergTablePartitionTimeArgsRequest) sdk.IcebergTablePartitionYearRequest {
				return *sdk.NewIcebergTablePartitionYearRequest().WithArgs(args)
			})),
			attributeMappedValueCreateBuilderNested(d, prefix+"month", req.WithMonth, parseIcebergTablePartitionTime(prefix+"month.0.", func(args sdk.IcebergTablePartitionTimeArgsRequest) sdk.IcebergTablePartitionMonthRequest {
				return *sdk.NewIcebergTablePartitionMonthRequest().WithArgs(args)
			})),
			attributeMappedValueCreateBuilderNested(d, prefix+"day", req.WithDay, parseIcebergTablePartitionTime(prefix+"day.0.", func(args sdk.IcebergTablePartitionTimeArgsRequest) sdk.IcebergTablePartitionDayRequest {
				return *sdk.NewIcebergTablePartitionDayRequest().WithArgs(args)
			})),
			attributeMappedValueCreateBuilderNested(d, prefix+"hour", req.WithHour, parseIcebergTablePartitionTime(prefix+"hour.0.", func(args sdk.IcebergTablePartitionTimeArgsRequest) sdk.IcebergTablePartitionHourRequest {
				return *sdk.NewIcebergTablePartitionHourRequest().WithArgs(args)
			})),
		)
		if err != nil {
			return nil, err
		}

		requests[i] = *req
	}
	return requests, nil
}

// parseIcebergTablePartitionBucket returns a mapper parsing a bucket partition block at the given prefix.
func parseIcebergTablePartitionBucket(prefix string) func(d *schema.ResourceData) (sdk.IcebergTablePartitionBucketRequest, error) {
	return func(d *schema.ResourceData) (sdk.IcebergTablePartitionBucketRequest, error) {
		return *sdk.NewIcebergTablePartitionBucketRequest().WithArgs(
			*sdk.NewIcebergTablePartitionBucketArgsRequest(d.Get(prefix+"num_buckets").(int), d.Get(prefix+"column").(string)),
		), nil
	}
}

// parseIcebergTablePartitionTruncate returns a mapper parsing a truncate partition block at the given prefix.
func parseIcebergTablePartitionTruncate(prefix string) func(d *schema.ResourceData) (sdk.IcebergTablePartitionTruncateRequest, error) {
	return func(d *schema.ResourceData) (sdk.IcebergTablePartitionTruncateRequest, error) {
		return *sdk.NewIcebergTablePartitionTruncateRequest().WithArgs(
			*sdk.NewIcebergTablePartitionTruncateArgsRequest(d.Get(prefix+"width").(int), d.Get(prefix+"column").(string)),
		), nil
	}
}

// parseIcebergTablePartitionTime returns a mapper parsing a time-based (year/month/day/hour) partition block at the given prefix.
func parseIcebergTablePartitionTime[T any](key string, build func(sdk.IcebergTablePartitionTimeArgsRequest) T) func(d *schema.ResourceData) (T, error) {
	return func(d *schema.ResourceData) (T, error) {
		return build(*sdk.NewIcebergTablePartitionTimeArgsRequest(d.Get(key).(string))), nil
	}
}

// icebergTablePartitionNumericTransformPattern matches the bucket[N] and truncate[N] transforms, capturing
// the transform kind and its numeric argument.
var icebergTablePartitionNumericTransformPattern = regexp.MustCompile(`^(bucket|truncate)\[(\d+)]$`)

// icebergTablePartitionSpecFieldsToSchema converts the fields of a SHOW ICEBERG TABLES partition spec back
// into the partition_by schema shape. Snowflake derives each field's name from its source column and
// transform (e.g. identity -> "<column>", bucket[4] -> "<column>_bucket_4", year -> "<column>_year"), so the
// source column name can be recovered by stripping the transform-specific suffix.
func icebergTablePartitionSpecFieldsToSchema(fields []sdk.IcebergTablePartitionSpecField) ([]map[string]any, error) {
	entries := make([]map[string]any, len(fields))
	for i, field := range fields {
		switch field.Transform {
		case "identity":
			entries[i] = map[string]any{"identity": field.Name}
		case "year", "month", "day", "hour":
			entries[i] = map[string]any{field.Transform: strings.TrimSuffix(field.Name, "_"+field.Transform)}
		default:
			match := icebergTablePartitionNumericTransformPattern.FindStringSubmatch(field.Transform)
			if match == nil {
				return nil, fmt.Errorf("unsupported Iceberg table partition transform: %s", field.Transform)
			}
			kind, arg := match[1], match[2]
			n, err := strconv.Atoi(arg)
			if err != nil {
				return nil, err
			}
			switch kind {
			case "bucket":
				column := strings.TrimSuffix(field.Name, "_bucket_"+arg)
				entries[i] = map[string]any{"bucket": []map[string]any{{"num_buckets": n, "column": column}}}
			case "truncate":
				column := strings.TrimSuffix(field.Name, "_trunc_"+arg)
				entries[i] = map[string]any{"truncate": []map[string]any{{"width": n, "column": column}}}
			}
		}
	}
	return entries, nil
}
