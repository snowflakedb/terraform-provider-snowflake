package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var (
	ExternalTableFileFormatTypeEnumDef = g.NewEnum(
		"ExternalTableFileFormatType", "ExternalTableFileFormatTypes",
		"CSV", "JSON", "AVRO", "ORC", "PARQUET",
	)
	ExternalTableCsvCompressionEnumDef = g.NewEnum(
		"ExternalTableCsvCompression", "ExternalTableCsvCompressions",
		"AUTO", "GZIP", "BZ2", "BROTLI", "ZSTD", "DEFLATE", "RAW_DEFLATE", "NONE",
	)
	ExternalTableJsonCompressionEnumDef = g.NewEnum(
		"ExternalTableJsonCompression", "ExternalTableJsonCompressions",
		"AUTO", "GZIP", "BZ2", "BROTLI", "ZSTD", "DEFLATE", "RAW_DEFLATE", "NONE",
	)
	ExternalTableAvroCompressionEnumDef = g.NewEnum(
		"ExternalTableAvroCompression", "ExternalTableAvroCompressions",
		"AUTO", "GZIP", "BZ2", "BROTLI", "ZSTD", "DEFLATE", "RAW_DEFLATE", "NONE",
	)
	ExternalTableParquetCompressionEnumDef = g.NewEnum(
		"ExternalTableParquetCompression", "ExternalTableParquetCompressions",
		"AUTO", "SNAPPY", "NONE",
	)
)

func cloudProviderParams() *g.QueryStruct {
	return g.NewQueryStruct("CloudProviderParams").
		OptionalAssignmentWithFieldName("INTEGRATION", "*string", g.ParameterOptions().SingleQuotes(), "GoogleCloudStorageIntegration").
		OptionalAssignmentWithFieldName("INTEGRATION", "*string", g.ParameterOptions().SingleQuotes(), "MicrosoftAzureIntegration").
		WithValidation(g.ConflictingFields, "GoogleCloudStorageIntegration", "MicrosoftAzureIntegration")
}

func rawFileFormat() *g.QueryStruct {
	return g.NewQueryStruct("RawFileFormat").
		Text("Format", g.KeywordOptions().Required())
}

func asExpression() *g.QueryStruct {
	return g.NewQueryStruct("AsExpression").
		Text("Expression", g.KeywordOptions().Required())
}

func externalTableFileFormatTypeOptions() *g.QueryStruct {
	return g.NewQueryStruct("ExternalTableFileFormatTypeOptions").
		OptionalAssignmentWithFieldName("COMPRESSION", ExternalTableCsvCompressionEnumDef.KindPtr(), g.ParameterOptions().NoQuotes(), "CSVCompression").
		OptionalAssignmentWithFieldName("RECORD_DELIMITER", "*string", g.ParameterOptions().SingleQuotes(), "CSVRecordDelimiter").
		OptionalAssignmentWithFieldName("FIELD_DELIMITER", "*string", g.ParameterOptions().SingleQuotes(), "CSVFieldDelimiter").
		OptionalAssignmentWithFieldName("SKIP_HEADER", "*int", g.ParameterOptions(), "CSVSkipHeader").
		OptionalAssignmentWithFieldName("SKIP_BLANK_LINES", "*bool", g.ParameterOptions(), "CSVSkipBlankLines").
		OptionalAssignmentWithFieldName("ESCAPE_UNENCLOSED_FIELD", "*string", g.ParameterOptions().SingleQuotes(), "CSVEscapeUnenclosedField").
		OptionalAssignmentWithFieldName("TRIM_SPACE", "*bool", g.ParameterOptions(), "CSVTrimSpace").
		OptionalAssignmentWithFieldName("FIELD_OPTIONALLY_ENCLOSED_BY", "*string", g.ParameterOptions().SingleQuotes(), "CSVFieldOptionallyEnclosedBy").
		ListAssignmentWithFieldName("NULL_IF", "NullString", g.ParameterOptions().Parentheses(), "CSVNullIf").
		OptionalAssignmentWithFieldName("EMPTY_FIELD_AS_NULL", "*bool", g.ParameterOptions(), "CSVEmptyFieldAsNull").
		PredefinedQueryStructField("CSVEncoding", "*CsvEncoding", g.ParameterOptions().SingleQuotes().SQL("ENCODING")).
		OptionalAssignmentWithFieldName("COMPRESSION", ExternalTableJsonCompressionEnumDef.KindPtr(), g.ParameterOptions().NoQuotes(), "JSONCompression").
		OptionalAssignmentWithFieldName("ALLOW_DUPLICATE", "*bool", g.ParameterOptions(), "JSONAllowDuplicate").
		OptionalAssignmentWithFieldName("STRIP_OUTER_ARRAY", "*bool", g.ParameterOptions(), "JSONStripOuterArray").
		OptionalAssignmentWithFieldName("STRIP_NULL_VALUES", "*bool", g.ParameterOptions(), "JSONStripNullValues").
		OptionalAssignmentWithFieldName("REPLACE_INVALID_CHARACTERS", "*bool", g.ParameterOptions(), "JSONReplaceInvalidCharacters").
		OptionalAssignmentWithFieldName("COMPRESSION", ExternalTableAvroCompressionEnumDef.KindPtr(), g.ParameterOptions().NoQuotes(), "AvroCompression").
		OptionalAssignmentWithFieldName("REPLACE_INVALID_CHARACTERS", "*bool", g.ParameterOptions(), "AvroReplaceInvalidCharacters").
		OptionalAssignmentWithFieldName("TRIM_SPACE", "*bool", g.ParameterOptions(), "ORCTrimSpace").
		OptionalAssignmentWithFieldName("REPLACE_INVALID_CHARACTERS", "*bool", g.ParameterOptions(), "ORCReplaceInvalidCharacters").
		ListAssignmentWithFieldName("NULL_IF", "NullString", g.ParameterOptions().Parentheses(), "ORCNullIf").
		OptionalAssignmentWithFieldName("COMPRESSION", ExternalTableParquetCompressionEnumDef.KindPtr(), g.ParameterOptions().NoQuotes(), "ParquetCompression").
		OptionalAssignmentWithFieldName("BINARY_AS_TEXT", "*bool", g.ParameterOptions(), "ParquetBinaryAsText").
		OptionalAssignmentWithFieldName("REPLACE_INVALID_CHARACTERS", "*bool", g.ParameterOptions(), "ParquetReplaceInvalidCharacters")
}

func externalTableFileFormat() *g.QueryStruct {
	return g.NewQueryStruct("ExternalTableFileFormat").
		OptionalAssignmentWithFieldName("FORMAT_NAME", "*string", g.ParameterOptions().SingleQuotes(), "Name").
		OptionalAssignmentWithFieldName("TYPE", ExternalTableFileFormatTypeEnumDef.KindPtr(), g.ParameterOptions().NoQuotes(), "FileFormatType").
		OptionalQueryStructField("Options", externalTableFileFormatTypeOptions(), nil).
		WithValidation(g.ExactlyOneValueSet, "Name", "FileFormatType").
		WithAdditionalValidations()
}

func externalTableColumn() *g.QueryStruct {
	return g.NewQueryStruct("ExternalTableColumn").
		Text("Name", g.KeywordOptions().Required()).
		PredefinedQueryStructField("DataType", g.KindOfT[sdkcommons.DataType](), g.KeywordOptions().Required()).
		QueryStructField("AsExpression", asExpression(), g.ListOptions().Parentheses().SQL("AS").Required()).
		OptionalSQL("NOT NULL").
		PredefinedQueryStructField("InlineConstraint", "*ColumnInlineConstraint", g.KeywordOptions())
}

func withExternalTableFileFormats(qs *g.QueryStruct) *g.QueryStruct {
	return qs.
		OptionalQueryStructField("FileFormat", externalTableFileFormat(), g.ListOptions().Parentheses().NoComma().SQL("FILE_FORMAT =")).
		// RawFileFormat was introduced, because of the decision taken during https://github.com/Snowflake-Labs/terraform-provider-snowflake/pull/2228
		// that for now the snowflake_external_table resource should continue on using raw file format, which wasn't previously supported by the new SDK.
		// In the future it should most likely be replaced by a more structured version FileFormat
		OptionalQueryStructField("RawFileFormat", rawFileFormat(), g.ListOptions().Parentheses().SQL("FILE_FORMAT ="))
}

func withExternalTablePolicyAndTags(qs *g.QueryStruct) *g.QueryStruct {
	return qs.
		OptionalComment().
		PredefinedQueryStructField("RowAccessPolicy", "*TableRowAccessPolicyLegacy", g.KeywordOptions()).
		OptionalTags()
}

func withExternalTableCreateTail(qs *g.QueryStruct) *g.QueryStruct {
	return withExternalTablePolicyAndTags(qs.OptionalSQL("COPY GRANTS"))
}

func refreshExternalTable() *g.QueryStruct {
	return g.NewQueryStruct("RefreshExternalTable").
		PredefinedQueryStructField("Path", "string", g.ParameterOptions().NoEquals().SingleQuotes())
}

func externalTableFile() *g.QueryStruct {
	return g.NewQueryStruct("ExternalTableFile").
		Text("Name", g.KeywordOptions().SingleQuotes().Required())
}

func externalTablePartition() *g.QueryStruct {
	return g.NewQueryStruct("Partition").
		Text("ColumnName", g.KeywordOptions().Required()).
		PredefinedQueryStructField("Value", "string", g.ParameterOptions().SingleQuotes().Required())
}

func externalTableDropOption() *g.QueryStruct {
	return g.NewQueryStruct("ExternalTableDropOption").
		OptionalSQL("RESTRICT").
		OptionalSQL("CASCADE").
		WithValidation(g.ConflictingFields, "Restrict", "Cascade")
}

var externalTablePairs = g.StructPair("externalTableRow", "ExternalTable").
	Time("created_on").
	Text("name").
	Text("database_name").
	Text("schema_name").
	Bool("invalid").
	OptionalText("invalid_reason", g.WithRequiredInPlain()).
	Text("owner").
	Text("comment").
	Text("stage").
	Text("location").
	Text("file_format_name").
	OptionalText("file_format_type", g.WithRequiredInPlain()).
	Text("cloud").
	Text("region").
	OptionalText("notification_channel", g.WithRequiredInPlain()).
	OptionalTime("last_refreshed_on", g.WithRequiredInPlain()).
	Text("table_format").
	OptionalText("last_refresh_details", g.WithRequiredInPlain()).
	Text("owner_role_type")

var externalTableColumnDetailsPairs = g.StructPair("externalTableColumnDetailsRow", "ExternalTableColumnDetails").
	Text("name").
	Field("type", "DataType", "DataType").
	Text("kind").
	BoolFromText("null?", g.WithPlainFieldName("IsNullable"), g.WithDbFieldName("Null")).
	OptionalText("default").
	BoolFromText("primary key", g.WithPlainFieldName("IsPrimary")).
	BoolFromText("unique key", g.WithPlainFieldName("IsUnique")).
	OptionalBoolFromText("check").
	OptionalText("expression").
	OptionalText("comment").
	OptionalText("policy name", g.WithPlainFieldName("PolicyName"))

var externalTableStageDetailsPairs = g.StructPair("externalTableStageDetailsRow", "ExternalTableStageDetails").
	Text("parent_property").
	Text("property").
	Text("property_type").
	Text("property_value").
	Text("property_default")

// TODO: ALTER column actions from https://docs.snowflake.com/en/sql-reference/sql/alter-table#external-table-column-actions-exttablecolumnaction
// are omitted — they were not in the pre-migration SDK.
// TODO: CREATE options present in https://docs.snowflake.com/en/sql-reference/sql/create-external-table
// but absent from the current SDK (e.g. additional TABLE_FORMAT values, newer file-format options) are omitted.
var externalTablesDef = g.NewInterface(
	"ExternalTables",
	"ExternalTable",
	g.KindOfT[sdkcommons.SchemaObjectIdentifier](),
).CreateOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/create-external-table",
	withExternalTableCreateTail(
		withExternalTableFileFormats(
			g.NewQueryStruct("CreateExternalTable").
				Create().
				OrReplace().
				SQL("EXTERNAL TABLE").
				IfNotExists().
				Name().
				ListQueryStructField("Columns", externalTableColumn(), g.ListOptions().Parentheses()).
				OptionalQueryStructField("CloudProviderParams", cloudProviderParams(), nil).
				PredefinedQueryStructField("PartitionBy", "[]string", g.KeywordOptions().Parentheses().SQL("PARTITION BY")).
				TextAssignment("LOCATION", g.ParameterOptions()).
				OptionalBooleanAssignment("REFRESH_ON_CREATE", g.ParameterOptions()).
				OptionalBooleanAssignment("AUTO_REFRESH", g.ParameterOptions()).
				OptionalTextAssignment("PATTERN", g.ParameterOptions().SingleQuotes()),
		).
			OptionalTextAssignment("AWS_SNS_TOPIC", g.ParameterOptions().SingleQuotes()),
	).
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists").
		WithValidation(g.ValidateValueSet, "Location").
		WithValidation(g.ExactlyOneValueSet, "RawFileFormat", "FileFormat"),
).CustomOperation(
	"CreateWithManualPartitioning",
	"https://docs.snowflake.com/en/sql-reference/sql/create-external-table",
	withExternalTableCreateTail(
		withExternalTableFileFormats(
			g.NewQueryStruct("CreateWithManualPartitioningExternalTable").
				Create().
				OrReplace().
				SQL("EXTERNAL TABLE").
				IfNotExists().
				Name().
				ListQueryStructField("Columns", externalTableColumn(), g.ListOptions().Parentheses()).
				OptionalQueryStructField("CloudProviderParams", cloudProviderParams(), nil).
				PredefinedQueryStructField("PartitionBy", "[]string", g.KeywordOptions().Parentheses().SQL("PARTITION BY")).
				TextAssignment("LOCATION", g.ParameterOptions()).
				SQL("PARTITION_TYPE = USER_SPECIFIED"),
		),
	).
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists").
		WithValidation(g.ValidateValueSet, "Location").
		WithValidation(g.ExactlyOneValueSet, "RawFileFormat", "FileFormat"),
).CustomOperation(
	"CreateDeltaLake",
	"https://docs.snowflake.com/en/sql-reference/sql/create-external-table",
	withExternalTableCreateTail(
		withExternalTableFileFormats(
			g.NewQueryStruct("CreateDeltaLakeExternalTable").
				Create().
				OrReplace().
				SQL("EXTERNAL TABLE").
				IfNotExists().
				Name().
				ListQueryStructField("Columns", externalTableColumn(), g.ListOptions().Parentheses()).
				OptionalQueryStructField("CloudProviderParams", cloudProviderParams(), nil).
				PredefinedQueryStructField("PartitionBy", "[]string", g.KeywordOptions().Parentheses().SQL("PARTITION BY")).
				TextAssignment("LOCATION", g.ParameterOptions()).
				// TODO: REFRESH_ON_CREATE and AUTO_REFRESH are kept to match the pre-migration SDK,
				// even though the Delta Lake syntax block in the docs omits them. The delta lake
				// integration test does set both (to false) and passes, so Snowflake seems to accept
				// them here. Confirm whether TRUE is accepted too, and drop them if it isn't.
				OptionalBooleanAssignment("REFRESH_ON_CREATE", g.ParameterOptions()).
				OptionalBooleanAssignment("AUTO_REFRESH", g.ParameterOptions()),
		).
			SQL("TABLE_FORMAT = DELTA"),
	).
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists").
		WithValidation(g.ValidateValueSet, "Location").
		WithValidation(g.ExactlyOneValueSet, "RawFileFormat", "FileFormat"),
).CustomOperation(
	"CreateUsingTemplate",
	"https://docs.snowflake.com/en/sql-reference/sql/create-external-table#variant-syntax",
	withExternalTablePolicyAndTags(
		withExternalTableFileFormats(
			g.NewQueryStruct("CreateExternalTableUsingTemplate").
				Create().
				OrReplace().
				SQL("EXTERNAL TABLE").
				Name().
				OptionalSQL("COPY GRANTS").
				PredefinedQueryStructField("Query", "[]string", g.ParameterOptions().NoEquals().Parentheses().SQL("USING TEMPLATE")).
				OptionalQueryStructField("CloudProviderParams", cloudProviderParams(), nil).
				PredefinedQueryStructField("PartitionBy", "[]string", g.KeywordOptions().Parentheses().SQL("PARTITION BY")).
				TextAssignment("LOCATION", g.ParameterOptions()).
				OptionalBooleanAssignment("REFRESH_ON_CREATE", g.ParameterOptions()).
				OptionalBooleanAssignment("AUTO_REFRESH", g.ParameterOptions()).
				OptionalTextAssignment("PATTERN", g.ParameterOptions().SingleQuotes()),
		).
			OptionalTextAssignment("AWS_SNS_TOPIC", g.ParameterOptions().SingleQuotes()),
	).
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ValidateValueSet, "Query").
		WithValidation(g.ValidateValueSet, "Location").
		WithValidation(g.ExactlyOneValueSet, "RawFileFormat", "FileFormat"),
).AlterOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/alter-external-table",
	g.NewQueryStruct("AlterExternalTable").
		Alter().
		SQL("EXTERNAL TABLE").
		IfExists().
		Name().
		OptionalQueryStructField("Refresh", refreshExternalTable(), g.KeywordOptions().SQL("REFRESH")).
		ListQueryStructField("AddFiles", externalTableFile(), g.KeywordOptions().NoQuotes().Parentheses().SQL("ADD FILES")).
		ListQueryStructField("RemoveFiles", externalTableFile(), g.KeywordOptions().NoQuotes().Parentheses().SQL("REMOVE FILES")).
		OptionalAssignmentWithFieldName("SET AUTO_REFRESH", "*bool", g.ParameterOptions(), "AutoRefresh").
		// TODO: SET TAG / UNSET TAG are kept to match the pre-migration SDK, but the ALTER EXTERNAL
		// TABLE docs state that tags on an external table have to be set through ALTER TABLE. The
		// snowflake_external_table resource calls both on update, yet neither is covered by an
		// integration or acceptance test. Verify against an account before relying on them.
		OptionalSetTags().
		OptionalUnsetTags().
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ExactlyOneValueSet, "Refresh", "AddFiles", "RemoveFiles", "AutoRefresh", "SetTags", "UnsetTags"),
).CustomOperation(
	"AlterPartitions",
	"https://docs.snowflake.com/en/sql-reference/sql/alter-external-table",
	g.NewQueryStruct("AlterExternalTablePartition").
		Alter().
		SQL("EXTERNAL TABLE").
		IfExists().
		Name().
		ListQueryStructField("AddPartitions", externalTablePartition(), g.KeywordOptions().Parentheses().SQL("ADD PARTITION")).
		OptionalSQL("DROP PARTITION").
		Assignment("LOCATION", "string", g.ParameterOptions().NoEquals().SingleQuotes()).
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ConflictingFields, "AddPartitions", "DropPartition"),
).DropOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/drop-external-table",
	g.NewQueryStruct("DropExternalTable").
		Drop().
		SQL("EXTERNAL TABLE").
		IfExists().
		Name().
		OptionalQueryStructField("DropOption", externalTableDropOption(), nil).
		WithValidation(g.ValidIdentifier, "name"),
).ShowOperationWithPairedStructs(
	"https://docs.snowflake.com/en/sql-reference/sql/show-external-tables",
	externalTablePairs,
	g.NewQueryStruct("ShowExternalTable").
		Show().
		Terse().
		SQL("EXTERNAL TABLES").
		OptionalLike().
		OptionalIn().
		OptionalStartsWith().
		PredefinedQueryStructField("LimitFrom", "*LimitFrom", g.KeywordOptions().SQL("LIMIT")),
	g.ShowByIDInFiltering,
	g.ShowByIDLikeFiltering,
).CustomShowOperationWithPairedStructs(
	"DescribeColumns",
	g.ShowMappingKindSlice,
	"https://docs.snowflake.com/en/sql-reference/sql/desc-external-table",
	externalTableColumnDetailsPairs,
	g.NewQueryStruct("DescribeExternalTableColumns").
		Describe().
		SQL("EXTERNAL TABLE").
		Name().
		SQL("TYPE = COLUMNS").
		WithValidation(g.ValidIdentifier, "name"),
).CustomShowOperationWithPairedStructs(
	"DescribeStage",
	g.ShowMappingKindSlice,
	"https://docs.snowflake.com/en/sql-reference/sql/desc-external-table",
	externalTableStageDetailsPairs,
	g.NewQueryStruct("DescribeExternalTableStage").
		Describe().
		SQL("EXTERNAL TABLE").
		Name().
		SQL("TYPE = STAGE").
		WithValidation(g.ValidIdentifier, "name"),
).WithShowByIDFindPredicateKind(g.ShowByIDFindPredicateFullID).WithEnums(
	ExternalTableFileFormatTypeEnumDef,
	ExternalTableCsvCompressionEnumDef,
	ExternalTableJsonCompressionEnumDef,
	ExternalTableAvroCompressionEnumDef,
	ExternalTableParquetCompressionEnumDef,
	// TODO(next-pr): drop this allow-list once SDK unit tests move to generated + *_ext_test.go
).WithAllowedGenerationParts(
	g.PartDefault,
	g.PartDto,
	g.PartDtoBuilders,
	g.PartImpl,
	g.PartValidations,
	g.PartEnums,
)
