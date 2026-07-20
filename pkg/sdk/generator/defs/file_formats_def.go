package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var (
	FileFormatTypeEnumDef = g.NewEnum(
		"FileFormatType", "FileFormatTypes",
		"CSV", "JSON", "AVRO", "ORC", "PARQUET", "XML",
	)
	BinaryFormatEnumDef = g.NewEnum(
		"BinaryFormat", "BinaryFormats",
		"HEX", "BASE64", "UTF8",
	)
	CsvCompressionEnumDef = g.NewEnum(
		"CsvCompression", "CsvCompressions",
		"AUTO", "GZIP", "BZ2", "BROTLI", "ZSTD", "DEFLATE", "RAW_DEFLATE", "NONE",
	)
	CsvEncodingEnumDef = g.NewEnum(
		"CsvEncoding", "CsvEncodings",
		"BIG5", "EUCJP", "EUCKR", "GB18030", "IBM420", "IBM424",
		"ISO2022CN", "ISO2022JP", "ISO2022KR",
		"ISO88591", "ISO88592", "ISO88595", "ISO88596", "ISO88597", "ISO88598", "ISO88599", "ISO885915",
		"KOI8R", "SHIFTJIS",
		"UTF8", "UTF16", "UTF16BE", "UTF16LE", "UTF32", "UTF32BE", "UTF32LE",
		"WINDOWS1250", "WINDOWS1251", "WINDOWS1252", "WINDOWS1253", "WINDOWS1254", "WINDOWS1255", "WINDOWS1256",
	)
	JsonCompressionEnumDef = g.NewEnum(
		"JsonCompression", "JsonCompressions",
		"AUTO", "GZIP", "BZ2", "BROTLI", "ZSTD", "DEFLATE", "RAW_DEFLATE", "NONE",
	)
	AvroCompressionEnumDef = g.NewEnum(
		"AvroCompression", "AvroCompressions",
		"AUTO", "GZIP", "BROTLI", "ZSTD", "DEFLATE", "RAW_DEFLATE", "NONE",
	)
	ParquetCompressionEnumDef = g.NewEnum(
		"ParquetCompression", "ParquetCompressions",
		"AUTO", "LZO", "SNAPPY", "NONE",
	)
	XmlCompressionEnumDef = g.NewEnum(
		"XmlCompression", "XmlCompressions",
		"AUTO", "GZIP", "BZ2", "BROTLI", "ZSTD", "DEFLATE", "RAW_DEFLATE", "NONE",
	)
)

func stageFileFormatStringOrAuto() *g.QueryStruct {
	return g.NewQueryStruct("StageFileFormatStringOrAuto").
		OptionalText("Value", g.KeywordOptions().SingleQuotes()).
		OptionalSQL("AUTO").
		WithValidation(g.ExactlyOneValueSet, "Value", "Auto")
}

func stageFileFormatStringOrNone() *g.QueryStruct {
	return g.NewQueryStruct("StageFileFormatStringOrNone").
		OptionalText("Value", g.KeywordOptions().SingleQuotes()).
		OptionalSQL("NONE").
		WithValidation(g.ExactlyOneValueSet, "Value", "None")
}

// csvFileFormatOptionFields appends the CSV-specific file format fields (unprefixed) onto qs.
// Reused both by fileFormatDef()'s nested embed struct and by the CreateCsv/AlterCsv operations.
func csvFileFormatOptionFields(qs *g.QueryStruct) *g.QueryStruct {
	return qs.
		OptionalEnumAssignment("COMPRESSION", CsvCompressionEnumDef, g.ParameterOptions().NoQuotes()).
		OptionalQueryStructField("RecordDelimiter", stageFileFormatStringOrNone(), g.ListOptions().NoParentheses().SQL("RECORD_DELIMITER =")).
		OptionalQueryStructField("FieldDelimiter", stageFileFormatStringOrNone(), g.ListOptions().NoParentheses().SQL("FIELD_DELIMITER =")).
		OptionalBooleanAssignment("MULTI_LINE", g.ParameterOptions()).
		OptionalTextAssignment("FILE_EXTENSION", g.ParameterOptions().SingleQuotes()).
		OptionalBooleanAssignment("PARSE_HEADER", g.ParameterOptions()).
		OptionalNumberAssignment("SKIP_HEADER", g.ParameterOptions()).
		OptionalBooleanAssignment("SKIP_BLANK_LINES", g.ParameterOptions()).
		OptionalQueryStructField("DateFormat", stageFileFormatStringOrAuto(), g.ListOptions().NoParentheses().SQL("DATE_FORMAT =")).
		OptionalQueryStructField("TimeFormat", stageFileFormatStringOrAuto(), g.ListOptions().NoParentheses().SQL("TIME_FORMAT =")).
		OptionalQueryStructField("TimestampFormat", stageFileFormatStringOrAuto(), g.ListOptions().NoParentheses().SQL("TIMESTAMP_FORMAT =")).
		OptionalEnumAssignment("BINARY_FORMAT", BinaryFormatEnumDef, g.ParameterOptions().NoQuotes()).
		OptionalQueryStructField("Escape", stageFileFormatStringOrNone(), g.ListOptions().NoParentheses().SQL("ESCAPE =")).
		OptionalQueryStructField("EscapeUnenclosedField", stageFileFormatStringOrNone(), g.ListOptions().NoParentheses().SQL("ESCAPE_UNENCLOSED_FIELD =")).
		OptionalBooleanAssignment("TRIM_SPACE", g.ParameterOptions()).
		OptionalQueryStructField("FieldOptionallyEnclosedBy", stageFileFormatStringOrNone(), g.ListOptions().NoParentheses().SQL("FIELD_OPTIONALLY_ENCLOSED_BY =")).
		ListAssignment("NULL_IF", "NullString", g.ParameterOptions().Parentheses()).
		OptionalBooleanAssignment("ERROR_ON_COLUMN_COUNT_MISMATCH", g.ParameterOptions()).
		OptionalBooleanAssignment("REPLACE_INVALID_CHARACTERS", g.ParameterOptions()).
		OptionalBooleanAssignment("EMPTY_FIELD_AS_NULL", g.ParameterOptions()).
		OptionalBooleanAssignment("SKIP_BYTE_ORDER_MARK", g.ParameterOptions()).
		OptionalEnumAssignment("ENCODING", CsvEncodingEnumDef, g.ParameterOptions().NoQuotes()).
		WithValidation(g.ConflictingFields, "SkipHeader", "ParseHeader")
}

func jsonFileFormatOptionFields(qs *g.QueryStruct) *g.QueryStruct {
	return qs.
		OptionalEnumAssignment("COMPRESSION", JsonCompressionEnumDef, g.ParameterOptions().NoQuotes()).
		OptionalQueryStructField("DateFormat", stageFileFormatStringOrAuto(), g.ListOptions().NoParentheses().SQL("DATE_FORMAT =")).
		OptionalQueryStructField("TimeFormat", stageFileFormatStringOrAuto(), g.ListOptions().NoParentheses().SQL("TIME_FORMAT =")).
		OptionalQueryStructField("TimestampFormat", stageFileFormatStringOrAuto(), g.ListOptions().NoParentheses().SQL("TIMESTAMP_FORMAT =")).
		OptionalEnumAssignment("BINARY_FORMAT", BinaryFormatEnumDef, g.ParameterOptions().NoQuotes()).
		OptionalBooleanAssignment("TRIM_SPACE", g.ParameterOptions()).
		OptionalBooleanAssignment("MULTI_LINE", g.ParameterOptions()).
		ListAssignment("NULL_IF", "NullString", g.ParameterOptions().Parentheses()).
		OptionalTextAssignment("FILE_EXTENSION", g.ParameterOptions().SingleQuotes()).
		OptionalBooleanAssignment("ENABLE_OCTAL", g.ParameterOptions()).
		OptionalBooleanAssignment("ALLOW_DUPLICATE", g.ParameterOptions()).
		OptionalBooleanAssignment("STRIP_OUTER_ARRAY", g.ParameterOptions()).
		OptionalBooleanAssignment("STRIP_NULL_VALUES", g.ParameterOptions()).
		OptionalBooleanAssignment("REPLACE_INVALID_CHARACTERS", g.ParameterOptions()).
		OptionalBooleanAssignment("IGNORE_UTF8_ERRORS", g.ParameterOptions()).
		OptionalBooleanAssignment("SKIP_BYTE_ORDER_MARK", g.ParameterOptions()).
		WithValidation(g.ConflictingFields, "IgnoreUtf8Errors", "ReplaceInvalidCharacters")
}

func avroFileFormatOptionFields(qs *g.QueryStruct) *g.QueryStruct {
	return qs.
		OptionalEnumAssignment("COMPRESSION", AvroCompressionEnumDef, g.ParameterOptions().NoQuotes()).
		OptionalBooleanAssignment("TRIM_SPACE", g.ParameterOptions()).
		OptionalBooleanAssignment("REPLACE_INVALID_CHARACTERS", g.ParameterOptions()).
		ListAssignment("NULL_IF", "NullString", g.ParameterOptions().Parentheses())
}

func orcFileFormatOptionFields(qs *g.QueryStruct) *g.QueryStruct {
	return qs.
		OptionalBooleanAssignment("TRIM_SPACE", g.ParameterOptions()).
		OptionalBooleanAssignment("REPLACE_INVALID_CHARACTERS", g.ParameterOptions()).
		ListAssignment("NULL_IF", "NullString", g.ParameterOptions().Parentheses())
}

func parquetFileFormatOptionFields(qs *g.QueryStruct) *g.QueryStruct {
	return qs.
		OptionalEnumAssignment("COMPRESSION", ParquetCompressionEnumDef, g.ParameterOptions().NoQuotes()).
		OptionalBooleanAssignment("SNAPPY_COMPRESSION", g.ParameterOptions()).
		OptionalBooleanAssignment("BINARY_AS_TEXT", g.ParameterOptions()).
		OptionalBooleanAssignment("USE_LOGICAL_TYPE", g.ParameterOptions()).
		OptionalBooleanAssignment("TRIM_SPACE", g.ParameterOptions()).
		OptionalBooleanAssignment("USE_VECTORIZED_SCANNER", g.ParameterOptions()).
		OptionalBooleanAssignment("REPLACE_INVALID_CHARACTERS", g.ParameterOptions()).
		ListAssignment("NULL_IF", "NullString", g.ParameterOptions().Parentheses()).
		WithValidation(g.ConflictingFields, "Compression", "SnappyCompression")
}

func xmlFileFormatOptionFields(qs *g.QueryStruct) *g.QueryStruct {
	return qs.
		OptionalEnumAssignment("COMPRESSION", XmlCompressionEnumDef, g.ParameterOptions().NoQuotes()).
		OptionalBooleanAssignment("IGNORE_UTF8_ERRORS", g.ParameterOptions()).
		OptionalBooleanAssignment("PRESERVE_SPACE", g.ParameterOptions()).
		OptionalBooleanAssignment("STRIP_OUTER_ELEMENT", g.ParameterOptions()).
		OptionalBooleanAssignment("DISABLE_AUTO_CONVERT", g.ParameterOptions()).
		OptionalBooleanAssignment("REPLACE_INVALID_CHARACTERS", g.ParameterOptions()).
		OptionalBooleanAssignment("SKIP_BYTE_ORDER_MARK", g.ParameterOptions()).
		WithValidation(g.ConflictingFields, "IgnoreUtf8Errors", "ReplaceInvalidCharacters")
}

// fileFormatCsvDetailsDef, fileFormatJsonDetailsDef, ... mirror csv/json/.../xmlFileFormatOptionFields
// for the DESCRIBE FILE FORMAT output, one struct per file format type, plus fileFormatAllDetailsDef
// which combines every type's fields (prefixed by type) for DescribeAllDetails.
var fileFormatCsvDetailsDef = g.PlainStruct("FileFormatCsvDetails").
	SchemaObjectIdentifier().
	OptionalField("Compression", CsvCompressionEnumDef.Kind()).
	OptionalField("RecordDelimiter", "StageFileFormatStringOrNone").
	OptionalField("FieldDelimiter", "StageFileFormatStringOrNone").
	OptionalText("FileExtension").
	OptionalNumber("SkipHeader").
	OptionalBool("ParseHeader").
	OptionalBool("SkipBlankLines").
	OptionalField("DateFormat", "StageFileFormatStringOrAuto").
	OptionalField("TimeFormat", "StageFileFormatStringOrAuto").
	OptionalField("TimestampFormat", "StageFileFormatStringOrAuto").
	OptionalField("BinaryFormat", BinaryFormatEnumDef.Kind()).
	OptionalField("Escape", "StageFileFormatStringOrNone").
	OptionalField("EscapeUnenclosedField", "StageFileFormatStringOrNone").
	OptionalBool("TrimSpace").
	OptionalField("FieldOptionallyEnclosedBy", "StageFileFormatStringOrNone").
	Field("NullIf", "[]NullString").
	OptionalBool("ErrorOnColumnCountMismatch").
	OptionalBool("ReplaceInvalidCharacters").
	OptionalBool("EmptyFieldAsNull").
	OptionalBool("SkipByteOrderMark").
	OptionalField("Encoding", CsvEncodingEnumDef.Kind())

var fileFormatJsonDetailsDef = g.PlainStruct("FileFormatJsonDetails").
	SchemaObjectIdentifier().
	OptionalField("Compression", JsonCompressionEnumDef.Kind()).
	OptionalField("DateFormat", "StageFileFormatStringOrAuto").
	OptionalField("TimeFormat", "StageFileFormatStringOrAuto").
	OptionalField("TimestampFormat", "StageFileFormatStringOrAuto").
	OptionalField("BinaryFormat", BinaryFormatEnumDef.Kind()).
	OptionalBool("TrimSpace").
	OptionalBool("MultiLine").
	Field("NullIf", "[]NullString").
	OptionalText("FileExtension").
	OptionalBool("EnableOctal").
	OptionalBool("AllowDuplicate").
	OptionalBool("StripOuterArray").
	OptionalBool("StripNullValues").
	OptionalBool("ReplaceInvalidCharacters").
	OptionalBool("IgnoreUtf8Errors").
	OptionalBool("SkipByteOrderMark")

var fileFormatAvroDetailsDef = g.PlainStruct("FileFormatAvroDetails").
	SchemaObjectIdentifier().
	OptionalField("Compression", AvroCompressionEnumDef.Kind()).
	OptionalBool("TrimSpace").
	OptionalBool("ReplaceInvalidCharacters").
	Field("NullIf", "[]NullString")

var fileFormatOrcDetailsDef = g.PlainStruct("FileFormatOrcDetails").
	SchemaObjectIdentifier().
	OptionalBool("TrimSpace").
	OptionalBool("ReplaceInvalidCharacters").
	Field("NullIf", "[]NullString")

var fileFormatParquetDetailsDef = g.PlainStruct("FileFormatParquetDetails").
	SchemaObjectIdentifier().
	OptionalField("Compression", ParquetCompressionEnumDef.Kind()).
	OptionalBool("TrimSpace").
	OptionalBool("BinaryAsText").
	OptionalBool("UseLogicalType").
	OptionalBool("UseVectorizedScanner").
	OptionalBool("ReplaceInvalidCharacters").
	Field("NullIf", "[]NullString")

var fileFormatXmlDetailsDef = g.PlainStruct("FileFormatXmlDetails").
	SchemaObjectIdentifier().
	OptionalField("Compression", XmlCompressionEnumDef.Kind()).
	OptionalBool("IgnoreUtf8Errors").
	OptionalBool("PreserveSpace").
	OptionalBool("StripOuterElement").
	OptionalBool("DisableAutoConvert").
	OptionalBool("ReplaceInvalidCharacters").
	OptionalBool("SkipByteOrderMark")

var fileFormatAllDetailsDef = g.PlainStruct("FileFormatAllDetails").
	SchemaObjectIdentifier().
	Field("Type", FileFormatTypeEnumDef.Kind()).
	OptionalField("CsvCompression", CsvCompressionEnumDef.Kind()).
	OptionalField("CsvRecordDelimiter", "StageFileFormatStringOrNone").
	OptionalField("CsvFieldDelimiter", "StageFileFormatStringOrNone").
	OptionalText("CsvFileExtension").
	OptionalNumber("CsvSkipHeader").
	OptionalBool("CsvParseHeader").
	OptionalBool("CsvSkipBlankLines").
	OptionalField("CsvDateFormat", "StageFileFormatStringOrAuto").
	OptionalField("CsvTimeFormat", "StageFileFormatStringOrAuto").
	OptionalField("CsvTimestampFormat", "StageFileFormatStringOrAuto").
	OptionalField("CsvBinaryFormat", BinaryFormatEnumDef.Kind()).
	OptionalField("CsvEscape", "StageFileFormatStringOrNone").
	OptionalField("CsvEscapeUnenclosedField", "StageFileFormatStringOrNone").
	OptionalBool("CsvTrimSpace").
	OptionalField("CsvFieldOptionallyEnclosedBy", "StageFileFormatStringOrNone").
	Field("CsvNullIf", "[]NullString").
	OptionalBool("CsvErrorOnColumnCountMismatch").
	OptionalBool("CsvReplaceInvalidCharacters").
	OptionalBool("CsvEmptyFieldAsNull").
	OptionalBool("CsvSkipByteOrderMark").
	OptionalField("CsvEncoding", CsvEncodingEnumDef.Kind()).
	OptionalField("JsonCompression", JsonCompressionEnumDef.Kind()).
	OptionalField("JsonDateFormat", "StageFileFormatStringOrAuto").
	OptionalField("JsonTimeFormat", "StageFileFormatStringOrAuto").
	OptionalField("JsonTimestampFormat", "StageFileFormatStringOrAuto").
	OptionalField("JsonBinaryFormat", BinaryFormatEnumDef.Kind()).
	OptionalBool("JsonTrimSpace").
	OptionalBool("JsonMultiLine").
	Field("JsonNullIf", "[]NullString").
	OptionalText("JsonFileExtension").
	OptionalBool("JsonEnableOctal").
	OptionalBool("JsonAllowDuplicate").
	OptionalBool("JsonStripOuterArray").
	OptionalBool("JsonStripNullValues").
	OptionalBool("JsonReplaceInvalidCharacters").
	OptionalBool("JsonIgnoreUtf8Errors").
	OptionalBool("JsonSkipByteOrderMark").
	OptionalField("AvroCompression", AvroCompressionEnumDef.Kind()).
	OptionalBool("AvroTrimSpace").
	OptionalBool("AvroReplaceInvalidCharacters").
	Field("AvroNullIf", "[]NullString").
	OptionalBool("OrcTrimSpace").
	OptionalBool("OrcReplaceInvalidCharacters").
	Field("OrcNullIf", "[]NullString").
	OptionalField("ParquetCompression", ParquetCompressionEnumDef.Kind()).
	OptionalBool("ParquetTrimSpace").
	OptionalBool("ParquetBinaryAsText").
	OptionalBool("ParquetUseLogicalType").
	OptionalBool("ParquetUseVectorizedScanner").
	OptionalBool("ParquetReplaceInvalidCharacters").
	Field("ParquetNullIf", "[]NullString").
	OptionalField("XmlCompression", XmlCompressionEnumDef.Kind()).
	OptionalBool("XmlIgnoreUtf8Errors").
	OptionalBool("XmlPreserveSpace").
	OptionalBool("XmlStripOuterElement").
	OptionalBool("XmlDisableAutoConvert").
	OptionalBool("XmlReplaceInvalidCharacters").
	OptionalBool("XmlSkipByteOrderMark")

// fileFormatDef models the nested, per-type FILE_FORMAT = (TYPE = CSV, ...) options used for
// embedding a file format directly into another object (e.g. CREATE/ALTER STAGE). Its structs
// are anchored to generation via a helper struct on the FileFormats Describe operation, since
// FileFormats itself only ever embeds them, but stages_def.go references "FileFormatOptions" by name.
func fileFormatDef() *g.QueryStruct {
	return g.NewQueryStruct("FileFormatOptions").
		OptionalQueryStructField(
			"CsvOptions",
			csvFileFormatOptionFields(g.NewQueryStruct("FileFormatCsvOptions").SQLWithCustomFieldName("formatType", "TYPE = CSV")),
			g.KeywordOptions(),
		).
		OptionalQueryStructField(
			"JsonOptions",
			jsonFileFormatOptionFields(g.NewQueryStruct("FileFormatJsonOptions").SQLWithCustomFieldName("formatType", "TYPE = JSON")),
			g.KeywordOptions(),
		).
		OptionalQueryStructField(
			"AvroOptions",
			avroFileFormatOptionFields(g.NewQueryStruct("FileFormatAvroOptions").SQLWithCustomFieldName("formatType", "TYPE = AVRO")),
			g.KeywordOptions(),
		).
		OptionalQueryStructField(
			"OrcOptions",
			orcFileFormatOptionFields(g.NewQueryStruct("FileFormatOrcOptions").SQLWithCustomFieldName("formatType", "TYPE = ORC")),
			g.KeywordOptions(),
		).
		OptionalQueryStructField(
			"ParquetOptions",
			parquetFileFormatOptionFields(g.NewQueryStruct("FileFormatParquetOptions").SQLWithCustomFieldName("formatType", "TYPE = PARQUET")),
			g.KeywordOptions(),
		).
		OptionalQueryStructField(
			"XmlOptions",
			xmlFileFormatOptionFields(g.NewQueryStruct("FileFormatXmlOptions").SQLWithCustomFieldName("formatType", "TYPE = XML")),
			g.KeywordOptions(),
		).
		WithValidation(g.ConflictingFields, "CsvOptions", "JsonOptions", "AvroOptions", "OrcOptions", "ParquetOptions", "XmlOptions")
}

// createFileFormat builds the CREATE FILE FORMAT ... TYPE = <sqlType> struct shared shape for a single type.
func createFileFormat(structName, sqlType string) *g.QueryStruct {
	return g.NewQueryStruct(structName).
		Create().
		OrReplace().
		SQL("FILE FORMAT").
		IfNotExists().
		Name().
		SQLWithCustomFieldName("formatType", "TYPE = "+sqlType)
}

// createFileFormatWithFields applies the type-specific option fields to createFileFormat and appends
// the comment field and name validation common to every Create<Type> operation.
func createFileFormatWithFields(structName, sqlType string, typeFields func(*g.QueryStruct) *g.QueryStruct) *g.QueryStruct {
	return typeFields(createFileFormat(structName, sqlType)).
		OptionalComment().
		WithValidation(g.ValidIdentifier, "name")
}

// alterFileFormatSetFields builds the SET (...) struct for a single type by applying its
// type-specific option fields and appending the comment field common to every Alter<Type> operation.
func alterFileFormatSetFields(structName string, typeFields func(*g.QueryStruct) *g.QueryStruct) *g.QueryStruct {
	return typeFields(g.NewQueryStruct(structName)).OptionalComment()
}

// alterFileFormat builds the shared ALTER FILE FORMAT ... [ RENAME TO | SET (...) ] shape for a single type.
func alterFileFormat(structName string, setFields *g.QueryStruct) *g.QueryStruct {
	return g.NewQueryStruct(structName).
		Alter().
		SQL("FILE FORMAT").
		IfExists().
		Name().
		RenameTo().
		OptionalQueryStructField("Set", setFields, g.ListOptions().NoParentheses().NoComma().SQL("SET")).
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ExactlyOneValueSet, "RenameTo", "Set")
}

var fileFormatsDef = g.NewInterface(
	"FileFormats",
	"FileFormat",
	g.KindOfT[sdkcommons.SchemaObjectIdentifier](),
).
	CustomOperation(
		"CreateCsv",
		"https://docs.snowflake.com/en/sql-reference/sql/create-file-format",
		createFileFormatWithFields("CreateCsvFileFormat", "CSV", csvFileFormatOptionFields),
	).
	CustomOperation(
		"CreateJson",
		"https://docs.snowflake.com/en/sql-reference/sql/create-file-format",
		createFileFormatWithFields("CreateJsonFileFormat", "JSON", jsonFileFormatOptionFields),
	).
	CustomOperation(
		"CreateAvro",
		"https://docs.snowflake.com/en/sql-reference/sql/create-file-format",
		createFileFormatWithFields("CreateAvroFileFormat", "AVRO", avroFileFormatOptionFields),
	).
	CustomOperation(
		"CreateOrc",
		"https://docs.snowflake.com/en/sql-reference/sql/create-file-format",
		createFileFormatWithFields("CreateOrcFileFormat", "ORC", orcFileFormatOptionFields),
	).
	CustomOperation(
		"CreateParquet",
		"https://docs.snowflake.com/en/sql-reference/sql/create-file-format",
		createFileFormatWithFields("CreateParquetFileFormat", "PARQUET", parquetFileFormatOptionFields),
	).
	CustomOperation(
		"CreateXml",
		"https://docs.snowflake.com/en/sql-reference/sql/create-file-format",
		createFileFormatWithFields("CreateXmlFileFormat", "XML", xmlFileFormatOptionFields),
	).
	CustomOperation(
		"AlterCsv",
		"https://docs.snowflake.com/en/sql-reference/sql/alter-file-format",
		alterFileFormat("AlterCsvFileFormat", alterFileFormatSetFields("AlterCsvFileFormatSet", csvFileFormatOptionFields)),
	).
	CustomOperation(
		"AlterJson",
		"https://docs.snowflake.com/en/sql-reference/sql/alter-file-format",
		alterFileFormat("AlterJsonFileFormat", alterFileFormatSetFields("AlterJsonFileFormatSet", jsonFileFormatOptionFields)),
	).
	CustomOperation(
		"AlterAvro",
		"https://docs.snowflake.com/en/sql-reference/sql/alter-file-format",
		alterFileFormat("AlterAvroFileFormat", alterFileFormatSetFields("AlterAvroFileFormatSet", avroFileFormatOptionFields)),
	).
	CustomOperation(
		"AlterOrc",
		"https://docs.snowflake.com/en/sql-reference/sql/alter-file-format",
		alterFileFormat("AlterOrcFileFormat", alterFileFormatSetFields("AlterOrcFileFormatSet", orcFileFormatOptionFields)),
	).
	CustomOperation(
		"AlterParquet",
		"https://docs.snowflake.com/en/sql-reference/sql/alter-file-format",
		alterFileFormat("AlterParquetFileFormat", alterFileFormatSetFields("AlterParquetFileFormatSet", parquetFileFormatOptionFields)),
	).
	CustomOperation(
		"AlterXml",
		"https://docs.snowflake.com/en/sql-reference/sql/alter-file-format",
		alterFileFormat("AlterXmlFileFormat", alterFileFormatSetFields("AlterXmlFileFormatSet", xmlFileFormatOptionFields)),
	).
	DropOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/drop-file-format",
		g.NewQueryStruct("DropFileFormat").
			Drop().
			SQL("FILE FORMAT").
			IfExists().
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	).
	ShowOperationWithPairedStructs(
		"https://docs.snowflake.com/en/sql-reference/sql/show-file-formats",
		g.StructPair("ShowFileFormatsRow", "FileFormat").
			Time("created_on").
			Text("name").
			Text("database_name").
			Text("schema_name").
			Enum("type", FileFormatTypeEnumDef, g.WithPlainFieldName("Type")).
			Text("owner").
			Text("comment").
			Text("owner_role_type").
			Text("format_options", g.WithPlainFieldName("FormatOptions")),
		g.NewQueryStruct("ShowFileFormats").
			Show().
			SQL("FILE FORMATS").
			OptionalLike().
			OptionalIn(),
		g.ShowByIDInFiltering,
		g.ShowByIDLikeFiltering,
	).
	DescribeOperationWithPairedStructs(
		g.DescriptionMappingKindSlice,
		"https://docs.snowflake.com/en/sql-reference/sql/desc-file-format",
		g.StructPair("descFileFormatsDbRow", "FileFormatProperty").
			Text("property", g.WithPlainFieldName("Name")).
			Text("property_type", g.WithPlainFieldName("Type")).
			Text("property_value", g.WithPlainFieldName("Value")).
			Text("property_default", g.WithPlainFieldName("Default")),
		g.NewQueryStruct("DescribeFileFormat").
			Describe().
			SQL("FILE FORMAT").
			Name().
			WithValidation(g.ValidIdentifier, "name"),
		fileFormatDef(),
		fileFormatCsvDetailsDef,
		fileFormatJsonDetailsDef,
		fileFormatAvroDetailsDef,
		fileFormatOrcDetailsDef,
		fileFormatParquetDetailsDef,
		fileFormatXmlDetailsDef,
		fileFormatAllDetailsDef,
	).
	WithCustomInterfaceMethod(
		"DescribeCsvDetails",
		"DescribeCsvDetails returns converted describe output for CSV file formats.",
		[]*g.MethodParameter{g.NewMethodParameter("id", g.KindOfT[sdkcommons.SchemaObjectIdentifier]())},
		"*FileFormatCsvDetails", "error",
	).
	WithCustomInterfaceMethod(
		"DescribeJsonDetails",
		"DescribeJsonDetails returns converted describe output for JSON file formats.",
		[]*g.MethodParameter{g.NewMethodParameter("id", g.KindOfT[sdkcommons.SchemaObjectIdentifier]())},
		"*FileFormatJsonDetails", "error",
	).
	WithCustomInterfaceMethod(
		"DescribeAvroDetails",
		"DescribeAvroDetails returns converted describe output for Avro file formats.",
		[]*g.MethodParameter{g.NewMethodParameter("id", g.KindOfT[sdkcommons.SchemaObjectIdentifier]())},
		"*FileFormatAvroDetails", "error",
	).
	WithCustomInterfaceMethod(
		"DescribeOrcDetails",
		"DescribeOrcDetails returns converted describe output for ORC file formats.",
		[]*g.MethodParameter{g.NewMethodParameter("id", g.KindOfT[sdkcommons.SchemaObjectIdentifier]())},
		"*FileFormatOrcDetails", "error",
	).
	WithCustomInterfaceMethod(
		"DescribeParquetDetails",
		"DescribeParquetDetails returns converted describe output for Parquet file formats.",
		[]*g.MethodParameter{g.NewMethodParameter("id", g.KindOfT[sdkcommons.SchemaObjectIdentifier]())},
		"*FileFormatParquetDetails", "error",
	).
	WithCustomInterfaceMethod(
		"DescribeXmlDetails",
		"DescribeXmlDetails returns converted describe output for XML file formats.",
		[]*g.MethodParameter{g.NewMethodParameter("id", g.KindOfT[sdkcommons.SchemaObjectIdentifier]())},
		"*FileFormatXmlDetails", "error",
	).
	WithCustomInterfaceMethod(
		"DescribeAllDetails",
		"DescribeAllDetails returns parsed describe output for any file format type.",
		[]*g.MethodParameter{g.NewMethodParameter("id", g.KindOfT[sdkcommons.SchemaObjectIdentifier]())},
		"*FileFormatAllDetails", "error",
	).
	WithEnums(
		FileFormatTypeEnumDef,
		BinaryFormatEnumDef,
		CsvCompressionEnumDef,
		CsvEncodingEnumDef,
		JsonCompressionEnumDef,
		AvroCompressionEnumDef,
		ParquetCompressionEnumDef,
		XmlCompressionEnumDef,
	)
