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

var fileFormatsDef = g.NewInterface(
	"FileFormats",
	"FileFormat",
	g.KindOfT[sdkcommons.SchemaObjectIdentifier](),
).
	CustomOperation(
		"DummyOperation",
		"not available",
		g.NewQueryStruct("FileFormatsDummyOperation").
			OptionalQueryStructField("FileFormat", fileFormatDef(), g.ListOptions().Parentheses().SQL("FILE_FORMAT =")),
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
