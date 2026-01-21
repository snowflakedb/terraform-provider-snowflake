package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

// Dual-value field helpers: String OR Keyword
// Wrapped in functions to avoid generator's nested field parent conflict
func stageFileFormatStringOrAuto() *g.QueryStruct {
	return g.NewQueryStruct("StageFileFormatStringOrAuto").
		OptionalTextAssignment("Value", g.ParameterOptions().SingleQuotes()).
		OptionalSQL("AUTO").
		WithValidation(g.ExactlyOneValueSet, "Value", "Auto")
}

func stageFileFormatStringOrNone() *g.QueryStruct {
	return g.NewQueryStruct("StageFileFormatStringOrNone").
		OptionalTextAssignment("Value", g.ParameterOptions().SingleQuotes()).
		OptionalSQL("NONE").
		WithValidation(g.ExactlyOneValueSet, "Value", "None")
}

func fileFormatDef() *g.QueryStruct {
	return g.NewQueryStruct("FileFormatOptions").
		OptionalQueryStructField(
			"CsvOptions",
			g.NewQueryStruct("FileFormatCsvOptions").
				PredefinedQueryStructField("formatType", "string", g.StaticOptions().SQL("TYPE = CSV")).
				OptionalAssignment("COMPRESSION", g.KindOfTPointer[sdkcommons.CsvCompression](), g.ParameterOptions().NoQuotes()).
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
				OptionalAssignment("BINARY_FORMAT", g.KindOfTPointer[sdkcommons.BinaryFormat](), g.ParameterOptions().NoQuotes()).
				OptionalQueryStructField("Escape", stageFileFormatStringOrNone(), g.ListOptions().NoParentheses().SQL("ESCAPE =")).
				OptionalQueryStructField("EscapeUnenclosedField", stageFileFormatStringOrNone(), g.ListOptions().NoParentheses().SQL("ESCAPE_UNENCLOSED_FIELD =")).
				OptionalBooleanAssignment("TRIM_SPACE", g.ParameterOptions()).
				OptionalQueryStructField("FieldOptionallyEnclosedBy", stageFileFormatStringOrNone(), g.ListOptions().NoParentheses().SQL("FIELD_OPTIONALLY_ENCLOSED_BY =")).
				ListAssignment("NULL_IF", "NullString", g.ParameterOptions().Parentheses()).
				OptionalBooleanAssignment("ERROR_ON_COLUMN_COUNT_MISMATCH", g.ParameterOptions()).
				OptionalBooleanAssignment("REPLACE_INVALID_CHARACTERS", g.ParameterOptions()).
				OptionalBooleanAssignment("EMPTY_FIELD_AS_NULL", g.ParameterOptions()).
				OptionalBooleanAssignment("SKIP_BYTE_ORDER_MARK", g.ParameterOptions()).
				OptionalTextAssignment("ENCODING", g.ParameterOptions().SingleQuotes()),
			// TODO: SKIP_HEADER and PARSE_HEADER are not compatible to be used together
			g.KeywordOptions(),
		).
		OptionalQueryStructField(
			"JsonOptions",
			g.NewQueryStruct("FileFormatJsonOptions").
				PredefinedQueryStructField("formatType", "string", g.StaticOptions().SQL("TYPE = JSON")).
				OptionalAssignment("COMPRESSION", g.KindOfTPointer[sdkcommons.JsonCompression](), g.ParameterOptions().NoQuotes()).
				OptionalQueryStructField("DateFormat", stageFileFormatStringOrAuto(), g.ListOptions().NoParentheses().SQL("DATE_FORMAT =")).
				OptionalQueryStructField("TimeFormat", stageFileFormatStringOrAuto(), g.ListOptions().NoParentheses().SQL("TIME_FORMAT =")).
				OptionalQueryStructField("TimestampFormat", stageFileFormatStringOrAuto(), g.ListOptions().NoParentheses().SQL("TIMESTAMP_FORMAT =")).
				OptionalAssignment("BINARY_FORMAT", g.KindOfTPointer[sdkcommons.BinaryFormat](), g.ParameterOptions().NoQuotes()).
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
				OptionalBooleanAssignment("SKIP_BYTE_ORDER_MARK", g.ParameterOptions()),
			// IGNORE_UTF8_ERRORS and REPLACE_INVALID_CHARACTERS are not compatible to be used together
			g.KeywordOptions(),
		).
		OptionalQueryStructField(
			"AvroOptions",
			g.NewQueryStruct("FileFormatAvroOptions").
				PredefinedQueryStructField("formatType", "string", g.StaticOptions().SQL("TYPE = AVRO")).
				OptionalAssignment("COMPRESSION", g.KindOfTPointer[sdkcommons.AvroCompression](), g.ParameterOptions().NoQuotes()).
				OptionalBooleanAssignment("TRIM_SPACE", g.ParameterOptions()).
				OptionalBooleanAssignment("REPLACE_INVALID_CHARACTERS", g.ParameterOptions()).
				ListAssignment("NULL_IF", "NullString", g.ParameterOptions().Parentheses()),
			g.KeywordOptions(),
		).
		OptionalQueryStructField(
			"OrcOptions",
			g.NewQueryStruct("FileFormatOrcOptions").
				PredefinedQueryStructField("formatType", "string", g.StaticOptions().SQL("TYPE = ORC")).
				OptionalBooleanAssignment("TRIM_SPACE", g.ParameterOptions()).
				OptionalBooleanAssignment("REPLACE_INVALID_CHARACTERS", g.ParameterOptions()).
				ListAssignment("NULL_IF", "NullString", g.ParameterOptions().Parentheses()),
			g.KeywordOptions(),
		).
		OptionalQueryStructField(
			"ParquetOptions",
			g.NewQueryStruct("FileFormatParquetOptions").
				PredefinedQueryStructField("formatType", "string", g.StaticOptions().SQL("TYPE = PARQUET")).
				OptionalAssignment("COMPRESSION", g.KindOfTPointer[sdkcommons.ParquetCompression](), g.ParameterOptions().NoQuotes()).
				OptionalBooleanAssignment("SNAPPY_COMPRESSION", g.ParameterOptions()).
				OptionalBooleanAssignment("BINARY_AS_TEXT", g.ParameterOptions()).
				OptionalBooleanAssignment("USE_LOGICAL_TYPE", g.ParameterOptions()).
				OptionalBooleanAssignment("TRIM_SPACE", g.ParameterOptions()).
				OptionalBooleanAssignment("USE_VECTORIZED_SCANNER", g.ParameterOptions()).
				OptionalBooleanAssignment("REPLACE_INVALID_CHARACTERS", g.ParameterOptions()).
				ListAssignment("NULL_IF", "NullString", g.ParameterOptions().Parentheses()),
			// COMPRESSION and SNAPPY_COMPRESSION options for parquet format is not allowed
			g.KeywordOptions(),
		).
		OptionalQueryStructField(
			"XmlOptions",
			g.NewQueryStruct("FileFormatXmlOptions").
				PredefinedQueryStructField("formatType", "string", g.StaticOptions().SQL("TYPE = XML")).
				OptionalAssignment("COMPRESSION", g.KindOfTPointer[sdkcommons.XmlCompression](), g.ParameterOptions().NoQuotes()).
				OptionalBooleanAssignment("IGNORE_UTF8_ERRORS", g.ParameterOptions()).
				OptionalBooleanAssignment("PRESERVE_SPACE", g.ParameterOptions()).
				OptionalBooleanAssignment("STRIP_OUTER_ELEMENT", g.ParameterOptions()).
				OptionalBooleanAssignment("DISABLE_AUTO_CONVERT", g.ParameterOptions()).
				OptionalBooleanAssignment("REPLACE_INVALID_CHARACTERS", g.ParameterOptions()).
				OptionalBooleanAssignment("SKIP_BYTE_ORDER_MARK", g.ParameterOptions()),
			// IGNORE_UTF8_ERRORS and REPLACE_INVALID_CHARACTERS are not compatible to be used together.

			g.KeywordOptions(),
		).
		WithValidation(g.ExactlyOneValueSet, "CsvOptions", "JsonOptions", "AvroOptions", "OrcOptions", "ParquetOptions", "XmlOptions")
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
		// PredefinedQueryStructField("Options", "*FileFormat", g.ListOptions().Parentheses().SQL("FILE_FORMAT =")),
		// fileFormatDef(),
		// fileFormatDef(),
	)
