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

func fileFormatDef() *g.QueryStruct {
	return g.NewQueryStruct("FileFormatOptions").
		OptionalQueryStructField(
			"CsvOptions",
			g.NewQueryStruct("FileFormatCsvOptions").
				SQLWithCustomFieldName("formatType", "TYPE = CSV").
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
				WithValidation(g.ConflictingFields, "SkipHeader", "ParseHeader"),
			g.KeywordOptions(),
		).
		OptionalQueryStructField(
			"JsonOptions",
			g.NewQueryStruct("FileFormatJsonOptions").
				SQLWithCustomFieldName("formatType", "TYPE = JSON").
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
				WithValidation(g.ConflictingFields, "IgnoreUtf8Errors", "ReplaceInvalidCharacters"),
			g.KeywordOptions(),
		).
		OptionalQueryStructField(
			"AvroOptions",
			g.NewQueryStruct("FileFormatAvroOptions").
				SQLWithCustomFieldName("formatType", "TYPE = AVRO").
				OptionalEnumAssignment("COMPRESSION", AvroCompressionEnumDef, g.ParameterOptions().NoQuotes()).
				OptionalBooleanAssignment("TRIM_SPACE", g.ParameterOptions()).
				OptionalBooleanAssignment("REPLACE_INVALID_CHARACTERS", g.ParameterOptions()).
				ListAssignment("NULL_IF", "NullString", g.ParameterOptions().Parentheses()),
			g.KeywordOptions(),
		).
		OptionalQueryStructField(
			"OrcOptions",
			g.NewQueryStruct("FileFormatOrcOptions").
				SQLWithCustomFieldName("formatType", "TYPE = ORC").
				OptionalBooleanAssignment("TRIM_SPACE", g.ParameterOptions()).
				OptionalBooleanAssignment("REPLACE_INVALID_CHARACTERS", g.ParameterOptions()).
				ListAssignment("NULL_IF", "NullString", g.ParameterOptions().Parentheses()),
			g.KeywordOptions(),
		).
		OptionalQueryStructField(
			"ParquetOptions",
			g.NewQueryStruct("FileFormatParquetOptions").
				SQLWithCustomFieldName("formatType", "TYPE = PARQUET").
				OptionalEnumAssignment("COMPRESSION", ParquetCompressionEnumDef, g.ParameterOptions().NoQuotes()).
				OptionalBooleanAssignment("SNAPPY_COMPRESSION", g.ParameterOptions()).
				OptionalBooleanAssignment("BINARY_AS_TEXT", g.ParameterOptions()).
				OptionalBooleanAssignment("USE_LOGICAL_TYPE", g.ParameterOptions()).
				OptionalBooleanAssignment("TRIM_SPACE", g.ParameterOptions()).
				OptionalBooleanAssignment("USE_VECTORIZED_SCANNER", g.ParameterOptions()).
				OptionalBooleanAssignment("REPLACE_INVALID_CHARACTERS", g.ParameterOptions()).
				ListAssignment("NULL_IF", "NullString", g.ParameterOptions().Parentheses()).
				WithValidation(g.ConflictingFields, "Compression", "SnappyCompression"),
			g.KeywordOptions(),
		).
		OptionalQueryStructField(
			"XmlOptions",
			g.NewQueryStruct("FileFormatXmlOptions").
				SQLWithCustomFieldName("formatType", "TYPE = XML").
				OptionalEnumAssignment("COMPRESSION", XmlCompressionEnumDef, g.ParameterOptions().NoQuotes()).
				OptionalBooleanAssignment("IGNORE_UTF8_ERRORS", g.ParameterOptions()).
				OptionalBooleanAssignment("PRESERVE_SPACE", g.ParameterOptions()).
				OptionalBooleanAssignment("STRIP_OUTER_ELEMENT", g.ParameterOptions()).
				OptionalBooleanAssignment("DISABLE_AUTO_CONVERT", g.ParameterOptions()).
				OptionalBooleanAssignment("REPLACE_INVALID_CHARACTERS", g.ParameterOptions()).
				OptionalBooleanAssignment("SKIP_BYTE_ORDER_MARK", g.ParameterOptions()).
				WithValidation(g.ConflictingFields, "IgnoreUtf8Errors", "ReplaceInvalidCharacters"),
			g.KeywordOptions(),
		).
		WithValidation(g.ConflictingFields, "CsvOptions", "JsonOptions", "AvroOptions", "OrcOptions", "ParquetOptions", "XmlOptions")
}

// fileFormatFlatOptionsDef models the flat, type-specific options shared by the standalone
// CREATE FILE FORMAT / ALTER FILE FORMAT ... SET statements (unlike fileFormatDef(), whose
// groups are nested for embedding as FILE_FORMAT = (TYPE = CSV, ...) inside other objects).
func fileFormatFlatOptionsDef() *g.QueryStruct {
	return g.NewQueryStruct("FileFormatObjectOptions").
		OptionalComment().
		// CSV
		OptionalAssignmentWithFieldName("COMPRESSION", CsvCompressionEnumDef.KindPtr(), g.ParameterOptions().NoQuotes(), "CsvCompression").
		OptionalQueryStructField("CsvRecordDelimiter", stageFileFormatStringOrNone(), g.ListOptions().NoParentheses().SQL("RECORD_DELIMITER =")).
		OptionalQueryStructField("CsvFieldDelimiter", stageFileFormatStringOrNone(), g.ListOptions().NoParentheses().SQL("FIELD_DELIMITER =")).
		OptionalAssignmentWithFieldName("MULTI_LINE", "*bool", g.ParameterOptions(), "CsvMultiLine").
		OptionalAssignmentWithFieldName("FILE_EXTENSION", "*string", g.ParameterOptions().SingleQuotes(), "CsvFileExtension").
		OptionalAssignmentWithFieldName("PARSE_HEADER", "*bool", g.ParameterOptions(), "CsvParseHeader").
		OptionalAssignmentWithFieldName("SKIP_HEADER", "*int", g.ParameterOptions(), "CsvSkipHeader").
		OptionalAssignmentWithFieldName("SKIP_BLANK_LINES", "*bool", g.ParameterOptions(), "CsvSkipBlankLines").
		OptionalQueryStructField("CsvDateFormat", stageFileFormatStringOrAuto(), g.ListOptions().NoParentheses().SQL("DATE_FORMAT =")).
		OptionalQueryStructField("CsvTimeFormat", stageFileFormatStringOrAuto(), g.ListOptions().NoParentheses().SQL("TIME_FORMAT =")).
		OptionalQueryStructField("CsvTimestampFormat", stageFileFormatStringOrAuto(), g.ListOptions().NoParentheses().SQL("TIMESTAMP_FORMAT =")).
		OptionalAssignmentWithFieldName("BINARY_FORMAT", BinaryFormatEnumDef.KindPtr(), g.ParameterOptions().NoQuotes(), "CsvBinaryFormat").
		OptionalQueryStructField("CsvEscape", stageFileFormatStringOrNone(), g.ListOptions().NoParentheses().SQL("ESCAPE =")).
		OptionalQueryStructField("CsvEscapeUnenclosedField", stageFileFormatStringOrNone(), g.ListOptions().NoParentheses().SQL("ESCAPE_UNENCLOSED_FIELD =")).
		OptionalAssignmentWithFieldName("TRIM_SPACE", "*bool", g.ParameterOptions(), "CsvTrimSpace").
		OptionalQueryStructField("CsvFieldOptionallyEnclosedBy", stageFileFormatStringOrNone(), g.ListOptions().NoParentheses().SQL("FIELD_OPTIONALLY_ENCLOSED_BY =")).
		ListAssignmentWithFieldName("NULL_IF", "NullString", g.ParameterOptions().Parentheses(), "CsvNullIf").
		OptionalAssignmentWithFieldName("ERROR_ON_COLUMN_COUNT_MISMATCH", "*bool", g.ParameterOptions(), "CsvErrorOnColumnCountMismatch").
		OptionalAssignmentWithFieldName("REPLACE_INVALID_CHARACTERS", "*bool", g.ParameterOptions(), "CsvReplaceInvalidCharacters").
		OptionalAssignmentWithFieldName("EMPTY_FIELD_AS_NULL", "*bool", g.ParameterOptions(), "CsvEmptyFieldAsNull").
		OptionalAssignmentWithFieldName("SKIP_BYTE_ORDER_MARK", "*bool", g.ParameterOptions(), "CsvSkipByteOrderMark").
		OptionalAssignmentWithFieldName("ENCODING", CsvEncodingEnumDef.KindPtr(), g.ParameterOptions().NoQuotes(), "CsvEncoding").
		// JSON
		OptionalAssignmentWithFieldName("COMPRESSION", JsonCompressionEnumDef.KindPtr(), g.ParameterOptions().NoQuotes(), "JsonCompression").
		OptionalQueryStructField("JsonDateFormat", stageFileFormatStringOrAuto(), g.ListOptions().NoParentheses().SQL("DATE_FORMAT =")).
		OptionalQueryStructField("JsonTimeFormat", stageFileFormatStringOrAuto(), g.ListOptions().NoParentheses().SQL("TIME_FORMAT =")).
		OptionalQueryStructField("JsonTimestampFormat", stageFileFormatStringOrAuto(), g.ListOptions().NoParentheses().SQL("TIMESTAMP_FORMAT =")).
		OptionalAssignmentWithFieldName("BINARY_FORMAT", BinaryFormatEnumDef.KindPtr(), g.ParameterOptions().NoQuotes(), "JsonBinaryFormat").
		OptionalAssignmentWithFieldName("TRIM_SPACE", "*bool", g.ParameterOptions(), "JsonTrimSpace").
		OptionalAssignmentWithFieldName("MULTI_LINE", "*bool", g.ParameterOptions(), "JsonMultiLine").
		ListAssignmentWithFieldName("NULL_IF", "NullString", g.ParameterOptions().Parentheses(), "JsonNullIf").
		OptionalAssignmentWithFieldName("FILE_EXTENSION", "*string", g.ParameterOptions().SingleQuotes(), "JsonFileExtension").
		OptionalAssignmentWithFieldName("ENABLE_OCTAL", "*bool", g.ParameterOptions(), "JsonEnableOctal").
		OptionalAssignmentWithFieldName("ALLOW_DUPLICATE", "*bool", g.ParameterOptions(), "JsonAllowDuplicate").
		OptionalAssignmentWithFieldName("STRIP_OUTER_ARRAY", "*bool", g.ParameterOptions(), "JsonStripOuterArray").
		OptionalAssignmentWithFieldName("STRIP_NULL_VALUES", "*bool", g.ParameterOptions(), "JsonStripNullValues").
		OptionalAssignmentWithFieldName("REPLACE_INVALID_CHARACTERS", "*bool", g.ParameterOptions(), "JsonReplaceInvalidCharacters").
		OptionalAssignmentWithFieldName("IGNORE_UTF8_ERRORS", "*bool", g.ParameterOptions(), "JsonIgnoreUtf8Errors").
		OptionalAssignmentWithFieldName("SKIP_BYTE_ORDER_MARK", "*bool", g.ParameterOptions(), "JsonSkipByteOrderMark").
		// AVRO
		OptionalAssignmentWithFieldName("COMPRESSION", AvroCompressionEnumDef.KindPtr(), g.ParameterOptions().NoQuotes(), "AvroCompression").
		OptionalAssignmentWithFieldName("TRIM_SPACE", "*bool", g.ParameterOptions(), "AvroTrimSpace").
		OptionalAssignmentWithFieldName("REPLACE_INVALID_CHARACTERS", "*bool", g.ParameterOptions(), "AvroReplaceInvalidCharacters").
		ListAssignmentWithFieldName("NULL_IF", "NullString", g.ParameterOptions().Parentheses(), "AvroNullIf").
		// ORC
		OptionalAssignmentWithFieldName("TRIM_SPACE", "*bool", g.ParameterOptions(), "OrcTrimSpace").
		OptionalAssignmentWithFieldName("REPLACE_INVALID_CHARACTERS", "*bool", g.ParameterOptions(), "OrcReplaceInvalidCharacters").
		ListAssignmentWithFieldName("NULL_IF", "NullString", g.ParameterOptions().Parentheses(), "OrcNullIf").
		// PARQUET
		OptionalAssignmentWithFieldName("COMPRESSION", ParquetCompressionEnumDef.KindPtr(), g.ParameterOptions().NoQuotes(), "ParquetCompression").
		OptionalAssignmentWithFieldName("SNAPPY_COMPRESSION", "*bool", g.ParameterOptions(), "ParquetSnappyCompression").
		OptionalAssignmentWithFieldName("BINARY_AS_TEXT", "*bool", g.ParameterOptions(), "ParquetBinaryAsText").
		OptionalAssignmentWithFieldName("USE_LOGICAL_TYPE", "*bool", g.ParameterOptions(), "ParquetUseLogicalType").
		OptionalAssignmentWithFieldName("TRIM_SPACE", "*bool", g.ParameterOptions(), "ParquetTrimSpace").
		OptionalAssignmentWithFieldName("USE_VECTORIZED_SCANNER", "*bool", g.ParameterOptions(), "ParquetUseVectorizedScanner").
		OptionalAssignmentWithFieldName("REPLACE_INVALID_CHARACTERS", "*bool", g.ParameterOptions(), "ParquetReplaceInvalidCharacters").
		ListAssignmentWithFieldName("NULL_IF", "NullString", g.ParameterOptions().Parentheses(), "ParquetNullIf").
		// XML
		OptionalAssignmentWithFieldName("COMPRESSION", XmlCompressionEnumDef.KindPtr(), g.ParameterOptions().NoQuotes(), "XmlCompression").
		OptionalAssignmentWithFieldName("IGNORE_UTF8_ERRORS", "*bool", g.ParameterOptions(), "XmlIgnoreUtf8Errors").
		OptionalAssignmentWithFieldName("PRESERVE_SPACE", "*bool", g.ParameterOptions(), "XmlPreserveSpace").
		OptionalAssignmentWithFieldName("STRIP_OUTER_ELEMENT", "*bool", g.ParameterOptions(), "XmlStripOuterElement").
		OptionalAssignmentWithFieldName("DISABLE_AUTO_CONVERT", "*bool", g.ParameterOptions(), "XmlDisableAutoConvert").
		OptionalAssignmentWithFieldName("REPLACE_INVALID_CHARACTERS", "*bool", g.ParameterOptions(), "XmlReplaceInvalidCharacters").
		OptionalAssignmentWithFieldName("SKIP_BYTE_ORDER_MARK", "*bool", g.ParameterOptions(), "XmlSkipByteOrderMark").
		WithValidation(g.ConflictingFields, "CsvSkipHeader", "CsvParseHeader").
		WithValidation(g.ConflictingFields, "JsonIgnoreUtf8Errors", "JsonReplaceInvalidCharacters").
		WithValidation(g.ConflictingFields, "ParquetCompression", "ParquetSnappyCompression").
		WithValidation(g.ConflictingFields, "XmlIgnoreUtf8Errors", "XmlReplaceInvalidCharacters").
		WithAdditionalValidations()
}

var fileFormatsDef = g.NewInterface(
	"FileFormats",
	"FileFormat",
	g.KindOfT[sdkcommons.SchemaObjectIdentifier](),
).
	CreateOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/create-file-format",
		g.NewQueryStruct("CreateFileFormat").
			Create().
			OrReplace().
			SQL("FILE FORMAT").
			IfNotExists().
			Name().
			EnumAssignmentWithFieldName("TYPE", FileFormatTypeEnumDef, g.ParameterOptions().Required().NoQuotes(), "FileFormatType").
			QueryStructField("FileFormatObjectOptions", fileFormatFlatOptionsDef(), g.ListOptions().NoParentheses().NoComma()).
			OptionalComment().
			WithValidation(g.ValidIdentifier, "name").
			WithAdditionalValidations(),
	).
	AlterOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/alter-file-format",
		g.NewQueryStruct("AlterFileFormat").
			Alter().
			SQL("FILE FORMAT").
			IfExists().
			Name().
			RenameTo().
			OptionalQueryStructField("Set", fileFormatFlatOptionsDef(), g.ListOptions().NoParentheses().NoComma().SQL("SET")).
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ExactlyOneValueSet, "RenameTo", "Set"),
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
		g.PlainStruct("FileFormatDetails").
			Enum("Type", FileFormatTypeEnumDef).
			Field("Options", "FileFormatObjectOptions"),
	).
	WithCustomInterfaceMethod(
		"DescribeDetails",
		"DescribeDetails returns the parsed, type-specific file format options.",
		[]*g.MethodParameter{g.NewMethodParameter("id", g.KindOfT[sdkcommons.SchemaObjectIdentifier]())},
		"*FileFormatDetails", "error",
	).
	CustomOperation(
		"DummyOperation",
		"not available",
		g.NewQueryStruct("FileFormatsDummyOperation").
			OptionalQueryStructField("FileFormat", g.RemoveValidations(fileFormatDef()), g.ListOptions().Parentheses().SQL("FILE_FORMAT =")).
			WithAdditionalValidations(),
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
