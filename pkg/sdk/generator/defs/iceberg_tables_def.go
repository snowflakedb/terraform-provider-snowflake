package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var (
	IcebergTableTargetFileSizeEnumDef = g.NewEnum(
		"IcebergTableTargetFileSize", "IcebergTableTargetFileSizes",
		"AUTO", "16MB", "32MB", "64MB", "128MB",
	)
	IcebergTablePathLayoutEnumDef = g.NewEnum(
		"IcebergTablePathLayout", "IcebergTablePathLayouts",
		"FLAT", "HIERARCHICAL",
	)
	IcebergTableDescribeTypeEnumDef = g.NewEnum(
		"IcebergTableDescribeType", "IcebergTableDescribeTypes",
		"COLUMNS", "STAGE",
	)
	IcebergTableLogEventLevelEnumDef = g.NewEnum(
		"IcebergTableLogEventLevel", "IcebergTableLogEventLevels",
		"ERROR", "WARN", "DEBUG",
	)
	IcebergTableTypeEnumDef = g.NewEnum(
		"IcebergTableType", "IcebergTableTypes",
		"MANAGED", "UNMANAGED", "NOT ICEBERG",
	)
	IcebergTableCatalogEnumDef = g.NewEnum(
		"IcebergTableCatalog", "IcebergTableCatalogs",
		"SNOWFLAKE",
	)
)

var icebergTableColumn = g.NewQueryStruct("IcebergTableColumn").
	Text("Name", g.KeywordOptions().Required().DoubleQuotes()).
	PredefinedQueryStructField("ColumnType", "datatypes.DataType", g.ParameterOptions().NoQuotes().NoEquals().Required()).
	PredefinedQueryStructField("DefaultValue", g.KindOfTPointer[sdkcommons.ColumnDefaultValue](), g.KeywordOptions()).
	OptionalSQL("NOT NULL").
	OptionalQueryStructField("InlineConstraint", tableColumnInlineConstraint(), g.KeywordOptions()).
	OptionalQueryStructField("MaskingPolicy", tableColumnMaskingPolicy, g.KeywordOptions()).
	OptionalQueryStructField("ProjectionPolicy", tableColumnProjectionPolicy, g.KeywordOptions()).
	OptionalTags().
	OptionalTextAssignment("COMMENT", g.ParameterOptions().NoEquals().SingleQuotes())

var icebergTableColumnsAndConstraints = g.NewQueryStruct("IcebergTableColumnsAndConstraints").
	ListQueryStructField("Columns", icebergTableColumn, g.KeywordOptions()).
	ListQueryStructField("OutOfLineConstraint", tableOutOfLineConstraint(), g.KeywordOptions())

var icebergTablePartitionBucketArgs = g.NewQueryStruct("IcebergTablePartitionBucketArgs").
	Number("NumBuckets", g.KeywordOptions().Required()).
	Text("Column", g.KeywordOptions().Required().DoubleQuotes())

var icebergTablePartitionBucket = g.NewQueryStruct("IcebergTablePartitionBucket").
	SQL("BUCKET").
	QueryStructField("Args", icebergTablePartitionBucketArgs, g.ListOptions().Parentheses())

var icebergTablePartitionTruncateArgs = g.NewQueryStruct("IcebergTablePartitionTruncateArgs").
	Number("Width", g.KeywordOptions().Required()).
	Text("Column", g.KeywordOptions().Required().DoubleQuotes())

var icebergTablePartitionTruncate = g.NewQueryStruct("IcebergTablePartitionTruncate").
	SQL("TRUNCATE").
	QueryStructField("Args", icebergTablePartitionTruncateArgs, g.ListOptions().Parentheses())

var icebergTablePartitionTimeArgs = g.NewQueryStruct("IcebergTablePartitionTimeArgs").
	Text("Column", g.KeywordOptions().Required().DoubleQuotes())

var icebergTablePartitionYear = g.NewQueryStruct("IcebergTablePartitionYear").
	SQL("YEAR").
	QueryStructField("Args", icebergTablePartitionTimeArgs, g.ListOptions().Parentheses())

var icebergTablePartitionMonth = g.NewQueryStruct("IcebergTablePartitionMonth").
	SQL("MONTH").
	QueryStructField("Args", icebergTablePartitionTimeArgs, g.ListOptions().Parentheses())

var icebergTablePartitionDay = g.NewQueryStruct("IcebergTablePartitionDay").
	SQL("DAY").
	QueryStructField("Args", icebergTablePartitionTimeArgs, g.ListOptions().Parentheses())

var icebergTablePartitionHour = g.NewQueryStruct("IcebergTablePartitionHour").
	SQL("HOUR").
	QueryStructField("Args", icebergTablePartitionTimeArgs, g.ListOptions().Parentheses())

var icebergTablePartitionExpression = g.NewQueryStruct("IcebergTablePartitionExpression").
	OptionalText("Identity", g.KeywordOptions().DoubleQuotes()).
	OptionalQueryStructField("Bucket", icebergTablePartitionBucket, g.KeywordOptions()).
	OptionalQueryStructField("Truncate", icebergTablePartitionTruncate, g.KeywordOptions()).
	OptionalQueryStructField("Year", icebergTablePartitionYear, g.KeywordOptions()).
	OptionalQueryStructField("Month", icebergTablePartitionMonth, g.KeywordOptions()).
	OptionalQueryStructField("Day", icebergTablePartitionDay, g.KeywordOptions()).
	OptionalQueryStructField("Hour", icebergTablePartitionHour, g.KeywordOptions())

// TODO [next PR]: validation is not generated properly as this is used as an array; using the additionalValidations above for now
// .WithValidation(g.ExactlyOneValueSet, "Identity", "Bucket", "Truncate", "Year", "Month", "Day", "Hour")

var icebergTableSetProperties = g.NewQueryStruct("IcebergTableSetProperties").
	OptionalBooleanAssignment("REPLACE_INVALID_CHARACTERS", g.ParameterOptions()).
	OptionalTextAssignment("CATALOG_SYNC", g.ParameterOptions().SingleQuotes()).
	OptionalNumberAssignment("DATA_RETENTION_TIME_IN_DAYS", g.ParameterOptions()).
	OptionalNumberAssignment("MAX_DATA_EXTENSION_TIME_IN_DAYS", g.ParameterOptions()).
	OptionalBooleanAssignment("AUTO_REFRESH", g.ParameterOptions()).
	OptionalAssignment("TARGET_FILE_SIZE", IcebergTableTargetFileSizeEnumDef.KindPtr(), g.ParameterOptions().SingleQuotes()).
	PredefinedQueryStructField("Contact", "[]TableContact", g.KeywordOptions().Parentheses().SQL("CONTACT")).
	OptionalAssignment("LOG_EVENT_LEVEL", IcebergTableLogEventLevelEnumDef.KindPtr(), g.ParameterOptions().NoQuotes()).
	OptionalBooleanAssignment("ERROR_LOGGING", g.ParameterOptions()).
	OptionalBooleanAssignment("ENABLE_DATA_COMPACTION", g.ParameterOptions()).
	OptionalBooleanAssignment("ENABLE_ICEBERG_MERGE_ON_READ", g.ParameterOptions()).
	OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
	WithValidation(g.AtLeastOneValueSet, "ReplaceInvalidCharacters", "CatalogSync", "DataRetentionTimeInDays", "MaxDataExtensionTimeInDays", "AutoRefresh", "TargetFileSize", "Contact", "LogEventLevel", "ErrorLogging", "EnableDataCompaction", "EnableIcebergMergeOnRead", "Comment")

var icebergTableUnsetProperties = g.NewQueryStruct("IcebergTableUnsetProperties").
	OptionalSQL("REPLACE_INVALID_CHARACTERS").
	OptionalSQL("CATALOG_SYNC").
	OptionalSQL("DATA_RETENTION_TIME_IN_DAYS").
	OptionalSQL("MAX_DATA_EXTENSION_TIME_IN_DAYS").
	OptionalSQL("TARGET_FILE_SIZE").
	OptionalSQL("LOG_EVENT_LEVEL").
	OptionalSQL("ERROR_LOGGING").
	OptionalSQL("ENABLE_DATA_COMPACTION").
	OptionalSQL("ENABLE_ICEBERG_MERGE_ON_READ").
	OptionalSQL("COMMENT").
	WithValidation(g.AtLeastOneValueSet, "ReplaceInvalidCharacters", "CatalogSync", "DataRetentionTimeInDays", "MaxDataExtensionTimeInDays", "TargetFileSize", "LogEventLevel", "ErrorLogging", "EnableDataCompaction", "EnableIcebergMergeOnRead", "Comment")

var icebergTableAddColumnAction = g.NewQueryStruct("IcebergTableAddColumnAction").
	SQL("ADD COLUMN").
	OptionalSQL("IF NOT EXISTS").
	Text("Name", g.KeywordOptions().Required().DoubleQuotes()).
	PredefinedQueryStructField("ColumnType", "datatypes.DataType", g.ParameterOptions().NoQuotes().NoEquals().Required()).
	OptionalQueryStructField("InlineConstraint", tableColumnInlineConstraint(), g.KeywordOptions()).
	PredefinedQueryStructField("DefaultValue", g.KindOfTPointer[sdkcommons.ColumnDefaultValue](), g.KeywordOptions()).
	OptionalQueryStructField("MaskingPolicy", tableColumnMaskingPolicy, g.KeywordOptions()).
	OptionalQueryStructField("ProjectionPolicy", tableColumnProjectionPolicy, g.KeywordOptions()).
	OptionalTags()

var icebergTableAlterColumnAction = g.NewQueryStruct("IcebergTableAlterColumnAction").
	SQL("COLUMN").
	Text("ColumnName", g.KeywordOptions().Required().DoubleQuotes()).
	OptionalSQL("SET NOT NULL").
	OptionalSQL("DROP NOT NULL").
	PredefinedQueryStructField("DataType", "*datatypes.DataType", g.ParameterOptions().NoEquals().SQL("SET DATA TYPE")).
	OptionalTextAssignment("COMMENT", g.ParameterOptions().NoEquals().SingleQuotes()).
	OptionalSQL("UNSET COMMENT").
	OptionalTextAssignment("SET WRITE DEFAULT", g.ParameterOptions().NoEquals()).
	OptionalSQL("DROP WRITE DEFAULT")

// TODO [next PR]: validation is not generated properly as this is used as an array; using the additionalValidations above for now
// .WithValidation(g.ExactlyOneValueSet, "SetNotNull", "DropNotNull", "DataType", "Comment", "UnsetComment", "SetWriteDefault", "DropWriteDefault")

var icebergTableClusteringAction = g.NewQueryStruct("IcebergTableClusteringAction").
	ListAssignment("CLUSTER BY", "string", g.ParameterOptions().NoEquals().Parentheses()).
	OptionalQueryStructField(
		"ChangeReclusterState",
		g.NewQueryStruct("IcebergTableReclusterChangeState").
			PredefinedQueryStructField("State", g.KindOfTPointer[sdkcommons.ReclusterState](), g.KeywordOptions()).
			SQL("RECLUSTER"),
		g.KeywordOptions(),
	).
	OptionalSQL("DROP CLUSTERING KEY").
	WithValidation(g.ExactlyOneValueSet, "ClusterBy", "ChangeReclusterState", "DropClusteringKey")

var icebergTableAggregationPolicy = g.NewQueryStruct("IcebergTableAggregationPolicy").
	Identifier("AggregationPolicy", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SQL("AGGREGATION POLICY").Required()).
	WithValidation(g.ValidIdentifier, "AggregationPolicy")

var tableRowAccessPolicy = g.NewQueryStruct("IcebergTableRowAccessPolicy").
	Identifier("Name", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SQL("ROW ACCESS POLICY").Required()).
	ListAssignment("ON", "Column", g.ParameterOptions().Required().NoEquals().Parentheses()).
	WithValidation(g.ValidIdentifier, "Name")

var icebergTablesDef = g.NewInterface(
	"IcebergTables",
	"IcebergTable",
	g.KindOfT[sdkcommons.SchemaObjectIdentifier](),
).CreateOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/create-iceberg-table-snowflake",
	g.NewQueryStruct("CreateIcebergTable").
		Create().
		OrReplace().
		OptionalSQL("TRANSIENT").
		SQL("ICEBERG TABLE").
		IfNotExists().
		Name().
		QueryStructField("ColumnsAndConstraints", icebergTableColumnsAndConstraints, g.ListOptions().Parentheses()).
		ListQueryStructField("PartitionBy", icebergTablePartitionExpression, g.KeywordOptions().Parentheses().SQL("PARTITION BY")).
		OptionalAssignment("PATH_LAYOUT", IcebergTablePathLayoutEnumDef.KindPtr(), g.ParameterOptions().NoQuotes()).
		PredefinedQueryStructField("ClusterBy", "[]string", g.KeywordOptions().Parentheses().SQL("CLUSTER BY")).
		OptionalIdentifier("ExternalVolume", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("EXTERNAL_VOLUME")).
		OptionalAssignment("CATALOG", IcebergTableCatalogEnumDef.KindPtr(), g.ParameterOptions().SingleQuotes()).
		OptionalTextAssignment("BASE_LOCATION", g.ParameterOptions().SingleQuotes()).
		OptionalAssignment("TARGET_FILE_SIZE", IcebergTableTargetFileSizeEnumDef.KindPtr(), g.ParameterOptions().SingleQuotes()).
		OptionalTextAssignment("CATALOG_SYNC", g.ParameterOptions().SingleQuotes()).
		PredefinedQueryStructField("StorageSerializationPolicy", "*StorageSerializationPolicy", g.ParameterOptions().SQL("STORAGE_SERIALIZATION_POLICY")).
		OptionalNumberAssignment("DATA_RETENTION_TIME_IN_DAYS", g.ParameterOptions()).
		OptionalNumberAssignment("MAX_DATA_EXTENSION_TIME_IN_DAYS", g.ParameterOptions()).
		OptionalBooleanAssignment("CHANGE_TRACKING", g.ParameterOptions()).
		OptionalSQL("COPY GRANTS").
		OptionalBooleanAssignment("ERROR_LOGGING", g.ParameterOptions()).
		OptionalComment().
		OptionalNumberAssignment("ICEBERG_VERSION", g.ParameterOptions()).
		OptionalBooleanAssignment("ENABLE_ICEBERG_MERGE_ON_READ", g.ParameterOptions()).
		OptionalQueryStructField("RowAccessPolicy", tableRowAccessPolicy, g.KeywordOptions()).
		OptionalQueryStructField("AggregationPolicy", icebergTableAggregationPolicy, g.KeywordOptions()).
		OptionalTags().
		OptionalBooleanAssignment("ENABLE_DATA_COMPACTION", g.ParameterOptions()).
		PredefinedQueryStructField("Contact", "[]TableContact", g.KeywordOptions().Parentheses().SQL("WITH CONTACT")).
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists").
		WithAdditionalValidations(),
).AlterOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/alter-iceberg-table",
	g.NewQueryStruct("AlterIcebergTable").
		Alter().
		SQL("ICEBERG TABLE").
		IfExists().
		Name().
		OptionalQueryStructField(
			"AddColumnAction",
			icebergTableAddColumnAction,
			g.KeywordOptions(),
		).
		OptionalQueryStructField(
			"DropColumnAction",
			tableDropColumnAction,
			g.KeywordOptions(),
		).
		OptionalQueryStructField(
			"RenameColumnAction",
			tableRenameColumnAction,
			g.KeywordOptions(),
		).
		ListQueryStructField(
			"AlterColumnAction",
			icebergTableAlterColumnAction,
			g.KeywordOptions().SQL("ALTER"),
		).
		OptionalQueryStructField(
			"SetMaskingPolicyOnColumn",
			tableSetColumnMaskingPolicy,
			g.KeywordOptions(),
		).
		OptionalQueryStructField(
			"UnsetMaskingPolicyOnColumn",
			tableUnsetColumnMaskingPolicy,
			g.KeywordOptions(),
		).
		OptionalQueryStructField(
			"SetProjectionPolicyOnColumn",
			tableSetColumnProjectionPolicy,
			g.KeywordOptions(),
		).
		OptionalQueryStructField(
			"UnsetProjectionPolicyOnColumn",
			tableUnsetColumnProjectionPolicy,
			g.KeywordOptions(),
		).
		OptionalQueryStructField(
			"SetTagsOnColumn",
			tableSetColumnTags,
			g.KeywordOptions(),
		).
		OptionalQueryStructField(
			"UnsetTagsOnColumn",
			tableUnsetColumnTags,
			g.KeywordOptions(),
		).
		OptionalQueryStructField(
			"ClusteringAction",
			icebergTableClusteringAction,
			g.KeywordOptions(),
		).
		OptionalQueryStructField(
			"Set",
			icebergTableSetProperties,
			g.KeywordOptions().SQL("SET"),
		).
		OptionalQueryStructField(
			"Unset",
			icebergTableUnsetProperties,
			g.KeywordOptions().SQL("UNSET"),
		).
		OptionalSetTags().
		OptionalUnsetTags().
		PredefinedQueryStructField("AddRowAccessPolicy", "*ViewAddRowAccessPolicy", g.KeywordOptions()).
		PredefinedQueryStructField("DropRowAccessPolicy", "*ViewDropRowAccessPolicy", g.KeywordOptions()).
		PredefinedQueryStructField("DropAndAddRowAccessPolicy", "*ViewDropAndAddRowAccessPolicy", g.ListOptions().NoParentheses()).
		OptionalSQL("DROP ALL ROW ACCESS POLICIES").
		// TODO(next PR): add ALTER ICEBERG TABLE ... REFRESH (separate operation; see https://docs.snowflake.com/en/sql-reference/sql/alter-iceberg-table-refresh)
		OptionalQueryStructField("SetAggregationPolicy", tableSetAggregationPolicy, g.KeywordOptions()).
		OptionalQueryStructField("UnsetAggregationPolicy", tableUnsetAggregationPolicy, g.KeywordOptions()).
		OptionalQueryStructField("SetJoinPolicy", tableSetJoinPolicy, g.KeywordOptions()).
		OptionalQueryStructField("UnsetJoinPolicy", tableUnsetJoinPolicy, g.KeywordOptions()).
		OptionalQueryStructField(
			"SearchOptimizationAction",
			tableSearchOptimizationAction,
			g.KeywordOptions(),
		).
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ExactlyOneValueSet, "AddColumnAction", "DropColumnAction", "RenameColumnAction", "AlterColumnAction", "SetMaskingPolicyOnColumn", "UnsetMaskingPolicyOnColumn", "SetProjectionPolicyOnColumn", "UnsetProjectionPolicyOnColumn", "SetTagsOnColumn", "UnsetTagsOnColumn", "ClusteringAction", "Set", "Unset", "SetTags", "UnsetTags", "AddRowAccessPolicy", "DropRowAccessPolicy", "DropAndAddRowAccessPolicy", "DropAllRowAccessPolicies", "SetAggregationPolicy", "UnsetAggregationPolicy", "SetJoinPolicy", "UnsetJoinPolicy", "SearchOptimizationAction").
		WithAdditionalValidations(),
).DropOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/drop-iceberg-table",
	g.NewQueryStruct("DropIcebergTable").
		Drop().
		SQL("ICEBERG TABLE").
		IfExists().
		Name().
		OptionalSQL("CASCADE").
		OptionalSQL("RESTRICT").
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ConflictingFields, "Cascade", "Restrict"),
).ShowOperationWithPairedStructs(
	"https://docs.snowflake.com/en/sql-reference/sql/show-iceberg-tables",
	g.StructPair("icebergTableRow", "IcebergTable").
		Time("created_on").
		Text("name").
		Text("database_name").
		Text("schema_name").
		OptionalText("owner").
		OptionalAccountObjectIdentifier("external_volume_name", g.WithPlainFieldName("ExternalVolumeName")).
		OptionalAccountObjectIdentifier("catalog_name", g.WithPlainFieldName("CatalogName")).
		Enum("iceberg_table_type", IcebergTableTypeEnumDef).
		OptionalText("catalog_table_name").
		OptionalText("catalog_namespace").
		Text("base_location").
		Field("can_write_metadata", "string", "bool").
		OptionalText("comment").
		OptionalText("name_mapping").
		Text("owner_role_type").
		Text("catalog_sync_name").
		Text("auto_refresh_status").
		Text("partition_specs").
		Number("current_partition_spec_id").
		Number("iceberg_table_format_version"),
	g.NewQueryStruct("ShowIcebergTables").
		Show().
		Terse().
		SQL("ICEBERG TABLES").
		OptionalLike().
		OptionalIn().
		OptionalStartsWith().
		OptionalLimitFrom(),
	// NOTE: TYPE=COLUMNS returns the same data as omitting the TYPE parameter.
	//       TYPE=STAGE returns underlying stage details, but this is not needed for the resource.
	g.ShowByIDInFiltering,
	g.ShowByIDLikeFiltering,
).DescribeOperationWithPairedStructs(
	g.DescriptionMappingKindSlice,
	"https://docs.snowflake.com/en/sql-reference/sql/desc-iceberg-table",
	g.StructPair("icebergTableDetailsRow", "IcebergTableDetails").
		Text("name").
		Text("type").
		Text("source iceberg type").
		Text("kind").
		PlainField("null", "bool", g.WithPlainFieldName("IsNullable")).
		OptionalText("default").
		PlainField("primary key", "bool", g.WithPlainFieldName("PrimaryKey")).
		PlainField("unique key", "bool", g.WithPlainFieldName("UniqueKey")).
		OptionalText("check").
		OptionalText("expression").
		OptionalText("comment").
		OptionalText("policy name", g.WithPlainFieldName("PolicyName")).
		OptionalText("privacy domain", g.WithPlainFieldName("PrivacyDomain")).
		OptionalText("name mapping", g.WithPlainFieldName("NameMapping")).
		OptionalText("write default", g.WithPlainFieldName("WriteDefault")),
	g.NewQueryStruct("DescribeIcebergTable").
		Describe().
		SQL("ICEBERG TABLE").
		Name().
		OptionalAssignmentWithFieldName("TYPE", IcebergTableDescribeTypeEnumDef.Kind(), g.ParameterOptions().NoQuotes(), "DescribeType").
		WithValidation(g.ValidIdentifier, "name"),
).WithEnums(
	IcebergTableTargetFileSizeEnumDef,
	IcebergTablePathLayoutEnumDef,
	IcebergTableDescribeTypeEnumDef,
	TableSearchMethodEnumDef,
	IcebergTableLogEventLevelEnumDef,
	IcebergTableTypeEnumDef,
	IcebergTableCatalogEnumDef,
)
