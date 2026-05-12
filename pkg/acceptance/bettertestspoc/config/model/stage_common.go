package model

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

// stageFileFormatName sets a named file format reference.
func stageFileFormatName(formatName string) tfconfig.Variable {
	return tfconfig.ListVariable(tfconfig.ObjectVariable(
		map[string]tfconfig.Variable{
			"format_name": tfconfig.StringVariable(formatName),
		},
	))
}

// stageFileFormatCsv sets inline CSV file format with the provided options.
func stageFileFormatCsv(opts sdk.FileFormatCsvOptions) tfconfig.Variable {
	csvMap := make(map[string]tfconfig.Variable)

	if opts.Compression != nil {
		csvMap["compression"] = tfconfig.StringVariable(string(*opts.Compression))
	}
	if opts.RecordDelimiter != nil {
		if opts.RecordDelimiter.None != nil && *opts.RecordDelimiter.None {
			csvMap["record_delimiter"] = tfconfig.StringVariable("NONE")
		} else if opts.RecordDelimiter.Value != nil {
			csvMap["record_delimiter"] = tfconfig.StringVariable(*opts.RecordDelimiter.Value)
		}
	}
	if opts.FieldDelimiter != nil {
		if opts.FieldDelimiter.None != nil && *opts.FieldDelimiter.None {
			csvMap["field_delimiter"] = tfconfig.StringVariable("NONE")
		} else if opts.FieldDelimiter.Value != nil {
			csvMap["field_delimiter"] = tfconfig.StringVariable(*opts.FieldDelimiter.Value)
		}
	}
	if opts.MultiLine != nil {
		csvMap["multi_line"] = tfconfig.BoolVariable(*opts.MultiLine)
	}
	if opts.FileExtension != nil {
		csvMap["file_extension"] = tfconfig.StringVariable(*opts.FileExtension)
	}
	if opts.ParseHeader != nil {
		csvMap["parse_header"] = tfconfig.BoolVariable(*opts.ParseHeader)
	}
	if opts.SkipHeader != nil {
		csvMap["skip_header"] = tfconfig.IntegerVariable(*opts.SkipHeader)
	}
	if opts.SkipBlankLines != nil {
		csvMap["skip_blank_lines"] = tfconfig.BoolVariable(*opts.SkipBlankLines)
	}
	if opts.DateFormat != nil {
		if opts.DateFormat.Auto != nil && *opts.DateFormat.Auto {
			csvMap["date_format"] = tfconfig.StringVariable("AUTO")
		} else if opts.DateFormat.Value != nil {
			csvMap["date_format"] = tfconfig.StringVariable(*opts.DateFormat.Value)
		}
	}
	if opts.TimeFormat != nil {
		if opts.TimeFormat.Auto != nil && *opts.TimeFormat.Auto {
			csvMap["time_format"] = tfconfig.StringVariable("AUTO")
		} else if opts.TimeFormat.Value != nil {
			csvMap["time_format"] = tfconfig.StringVariable(*opts.TimeFormat.Value)
		}
	}
	if opts.TimestampFormat != nil {
		if opts.TimestampFormat.Auto != nil && *opts.TimestampFormat.Auto {
			csvMap["timestamp_format"] = tfconfig.StringVariable("AUTO")
		} else if opts.TimestampFormat.Value != nil {
			csvMap["timestamp_format"] = tfconfig.StringVariable(*opts.TimestampFormat.Value)
		}
	}
	if opts.BinaryFormat != nil {
		csvMap["binary_format"] = tfconfig.StringVariable(string(*opts.BinaryFormat))
	}
	if opts.Escape != nil {
		if opts.Escape.None != nil && *opts.Escape.None {
			csvMap["escape"] = tfconfig.StringVariable("NONE")
		} else if opts.Escape.Value != nil {
			csvMap["escape"] = tfconfig.StringVariable(*opts.Escape.Value)
		}
	}
	if opts.EscapeUnenclosedField != nil {
		if opts.EscapeUnenclosedField.None != nil && *opts.EscapeUnenclosedField.None {
			csvMap["escape_unenclosed_field"] = tfconfig.StringVariable("NONE")
		} else if opts.EscapeUnenclosedField.Value != nil {
			csvMap["escape_unenclosed_field"] = tfconfig.StringVariable(*opts.EscapeUnenclosedField.Value)
		}
	}
	if opts.TrimSpace != nil {
		csvMap["trim_space"] = tfconfig.BoolVariable(*opts.TrimSpace)
	}
	if opts.FieldOptionallyEnclosedBy != nil {
		if opts.FieldOptionallyEnclosedBy.None != nil && *opts.FieldOptionallyEnclosedBy.None {
			csvMap["field_optionally_enclosed_by"] = tfconfig.StringVariable("NONE")
		} else if opts.FieldOptionallyEnclosedBy.Value != nil {
			csvMap["field_optionally_enclosed_by"] = tfconfig.StringVariable(*opts.FieldOptionallyEnclosedBy.Value)
		}
	}
	if len(opts.NullIf) > 0 {
		nullIfVars := make([]tfconfig.Variable, len(opts.NullIf))
		for idx, v := range opts.NullIf {
			nullIfVars[idx] = tfconfig.StringVariable(v.S)
		}
		csvMap["null_if"] = tfconfig.ListVariable(nullIfVars...)
	}
	if opts.ErrorOnColumnCountMismatch != nil {
		csvMap["error_on_column_count_mismatch"] = tfconfig.BoolVariable(*opts.ErrorOnColumnCountMismatch)
	}
	if opts.ReplaceInvalidCharacters != nil {
		csvMap["replace_invalid_characters"] = tfconfig.BoolVariable(*opts.ReplaceInvalidCharacters)
	}
	if opts.EmptyFieldAsNull != nil {
		csvMap["empty_field_as_null"] = tfconfig.BoolVariable(*opts.EmptyFieldAsNull)
	}
	if opts.SkipByteOrderMark != nil {
		csvMap["skip_byte_order_mark"] = tfconfig.BoolVariable(*opts.SkipByteOrderMark)
	}
	if opts.Encoding != nil {
		csvMap["encoding"] = tfconfig.StringVariable(string(*opts.Encoding))
	}

	// Workaround for empty objects - Terraform requires at least one attribute
	if len(csvMap) == 0 {
		csvMap["any"] = tfconfig.StringVariable(string(config.SnowflakeProviderConfigSingleAttributeWorkaround))
	}

	return tfconfig.ListVariable(tfconfig.ObjectVariable(
		map[string]tfconfig.Variable{
			"csv": tfconfig.ListVariable(tfconfig.ObjectVariable(csvMap)),
		},
	))
}

// stageFileFormatJson sets inline JSON file format with the provided options.
func stageFileFormatJson(opts sdk.FileFormatJsonOptions) tfconfig.Variable {
	jsonMap := make(map[string]tfconfig.Variable)

	if opts.Compression != nil {
		jsonMap["compression"] = tfconfig.StringVariable(string(*opts.Compression))
	}
	if opts.DateFormat != nil {
		if opts.DateFormat.Auto != nil && *opts.DateFormat.Auto {
			jsonMap["date_format"] = tfconfig.StringVariable("AUTO")
		} else if opts.DateFormat.Value != nil {
			jsonMap["date_format"] = tfconfig.StringVariable(*opts.DateFormat.Value)
		}
	}
	if opts.TimeFormat != nil {
		if opts.TimeFormat.Auto != nil && *opts.TimeFormat.Auto {
			jsonMap["time_format"] = tfconfig.StringVariable("AUTO")
		} else if opts.TimeFormat.Value != nil {
			jsonMap["time_format"] = tfconfig.StringVariable(*opts.TimeFormat.Value)
		}
	}
	if opts.TimestampFormat != nil {
		if opts.TimestampFormat.Auto != nil && *opts.TimestampFormat.Auto {
			jsonMap["timestamp_format"] = tfconfig.StringVariable("AUTO")
		} else if opts.TimestampFormat.Value != nil {
			jsonMap["timestamp_format"] = tfconfig.StringVariable(*opts.TimestampFormat.Value)
		}
	}
	if opts.BinaryFormat != nil {
		jsonMap["binary_format"] = tfconfig.StringVariable(string(*opts.BinaryFormat))
	}
	if opts.TrimSpace != nil {
		jsonMap["trim_space"] = tfconfig.BoolVariable(*opts.TrimSpace)
	}
	if opts.MultiLine != nil {
		jsonMap["multi_line"] = tfconfig.BoolVariable(*opts.MultiLine)
	}
	if len(opts.NullIf) > 0 {
		nullIfVars := make([]tfconfig.Variable, len(opts.NullIf))
		for idx, v := range opts.NullIf {
			nullIfVars[idx] = tfconfig.StringVariable(v.S)
		}
		jsonMap["null_if"] = tfconfig.ListVariable(nullIfVars...)
	}
	if opts.FileExtension != nil {
		jsonMap["file_extension"] = tfconfig.StringVariable(*opts.FileExtension)
	}
	if opts.EnableOctal != nil {
		jsonMap["enable_octal"] = tfconfig.BoolVariable(*opts.EnableOctal)
	}
	if opts.AllowDuplicate != nil {
		jsonMap["allow_duplicate"] = tfconfig.BoolVariable(*opts.AllowDuplicate)
	}
	if opts.StripOuterArray != nil {
		jsonMap["strip_outer_array"] = tfconfig.BoolVariable(*opts.StripOuterArray)
	}
	if opts.StripNullValues != nil {
		jsonMap["strip_null_values"] = tfconfig.BoolVariable(*opts.StripNullValues)
	}
	if opts.ReplaceInvalidCharacters != nil {
		jsonMap["replace_invalid_characters"] = tfconfig.BoolVariable(*opts.ReplaceInvalidCharacters)
	}
	if opts.IgnoreUtf8Errors != nil {
		jsonMap["ignore_utf8_errors"] = tfconfig.BoolVariable(*opts.IgnoreUtf8Errors)
	}
	if opts.SkipByteOrderMark != nil {
		jsonMap["skip_byte_order_mark"] = tfconfig.BoolVariable(*opts.SkipByteOrderMark)
	}

	// Workaround for empty objects - Terraform requires at least one attribute
	if len(jsonMap) == 0 {
		jsonMap["any"] = tfconfig.StringVariable(string(config.SnowflakeProviderConfigSingleAttributeWorkaround))
	}

	return tfconfig.ListVariable(tfconfig.ObjectVariable(
		map[string]tfconfig.Variable{
			"json": tfconfig.ListVariable(tfconfig.ObjectVariable(jsonMap)),
		},
	))
}

// stageFileFormatAvro sets inline AVRO file format with the provided options.
func stageFileFormatAvro(opts sdk.FileFormatAvroOptions) tfconfig.Variable {
	avroMap := make(map[string]tfconfig.Variable)

	if opts.Compression != nil {
		avroMap["compression"] = tfconfig.StringVariable(string(*opts.Compression))
	}
	if opts.TrimSpace != nil {
		avroMap["trim_space"] = tfconfig.BoolVariable(*opts.TrimSpace)
	}
	if opts.ReplaceInvalidCharacters != nil {
		avroMap["replace_invalid_characters"] = tfconfig.BoolVariable(*opts.ReplaceInvalidCharacters)
	}
	if len(opts.NullIf) > 0 {
		nullIfVars := make([]tfconfig.Variable, len(opts.NullIf))
		for idx, v := range opts.NullIf {
			nullIfVars[idx] = tfconfig.StringVariable(v.S)
		}
		avroMap["null_if"] = tfconfig.ListVariable(nullIfVars...)
	}

	// Workaround for empty objects - Terraform requires at least one attribute
	if len(avroMap) == 0 {
		avroMap["any"] = tfconfig.StringVariable(string(config.SnowflakeProviderConfigSingleAttributeWorkaround))
	}

	return tfconfig.ListVariable(tfconfig.ObjectVariable(
		map[string]tfconfig.Variable{
			"avro": tfconfig.ListVariable(tfconfig.ObjectVariable(avroMap)),
		},
	))
}

// stageFileFormatOrc sets inline ORC file format with the provided options.
func stageFileFormatOrc(opts sdk.FileFormatOrcOptions) tfconfig.Variable {
	orcMap := make(map[string]tfconfig.Variable)

	if opts.TrimSpace != nil {
		orcMap["trim_space"] = tfconfig.BoolVariable(*opts.TrimSpace)
	}
	if opts.ReplaceInvalidCharacters != nil {
		orcMap["replace_invalid_characters"] = tfconfig.BoolVariable(*opts.ReplaceInvalidCharacters)
	}
	if len(opts.NullIf) > 0 {
		nullIfVars := make([]tfconfig.Variable, len(opts.NullIf))
		for idx, v := range opts.NullIf {
			nullIfVars[idx] = tfconfig.StringVariable(v.S)
		}
		orcMap["null_if"] = tfconfig.ListVariable(nullIfVars...)
	}

	// Workaround for empty objects - Terraform requires at least one attribute
	if len(orcMap) == 0 {
		orcMap["any"] = tfconfig.StringVariable(string(config.SnowflakeProviderConfigSingleAttributeWorkaround))
	}

	return tfconfig.ListVariable(tfconfig.ObjectVariable(
		map[string]tfconfig.Variable{
			"orc": tfconfig.ListVariable(tfconfig.ObjectVariable(orcMap)),
		},
	))
}

// stageFileFormatParquet sets inline Parquet file format with the provided options.
func stageFileFormatParquet(opts sdk.FileFormatParquetOptions) tfconfig.Variable {
	parquetMap := make(map[string]tfconfig.Variable)

	if opts.Compression != nil {
		parquetMap["compression"] = tfconfig.StringVariable(string(*opts.Compression))
	}
	if opts.BinaryAsText != nil {
		parquetMap["binary_as_text"] = tfconfig.BoolVariable(*opts.BinaryAsText)
	}
	if opts.UseLogicalType != nil {
		parquetMap["use_logical_type"] = tfconfig.BoolVariable(*opts.UseLogicalType)
	}
	if opts.TrimSpace != nil {
		parquetMap["trim_space"] = tfconfig.BoolVariable(*opts.TrimSpace)
	}
	if opts.UseVectorizedScanner != nil {
		parquetMap["use_vectorized_scanner"] = tfconfig.BoolVariable(*opts.UseVectorizedScanner)
	}
	if opts.ReplaceInvalidCharacters != nil {
		parquetMap["replace_invalid_characters"] = tfconfig.BoolVariable(*opts.ReplaceInvalidCharacters)
	}
	if len(opts.NullIf) > 0 {
		nullIfVars := make([]tfconfig.Variable, len(opts.NullIf))
		for idx, v := range opts.NullIf {
			nullIfVars[idx] = tfconfig.StringVariable(v.S)
		}
		parquetMap["null_if"] = tfconfig.ListVariable(nullIfVars...)
	}

	// Workaround for empty objects - Terraform requires at least one attribute
	if len(parquetMap) == 0 {
		parquetMap["any"] = tfconfig.StringVariable(string(config.SnowflakeProviderConfigSingleAttributeWorkaround))
	}

	return tfconfig.ListVariable(tfconfig.ObjectVariable(
		map[string]tfconfig.Variable{
			"parquet": tfconfig.ListVariable(tfconfig.ObjectVariable(parquetMap)),
		},
	))
}

// stageFileFormatXml sets inline XML file format with the provided options.
func stageFileFormatXml(opts sdk.FileFormatXmlOptions) tfconfig.Variable {
	xmlMap := make(map[string]tfconfig.Variable)

	if opts.Compression != nil {
		xmlMap["compression"] = tfconfig.StringVariable(string(*opts.Compression))
	}
	if opts.IgnoreUtf8Errors != nil {
		xmlMap["ignore_utf8_errors"] = tfconfig.BoolVariable(*opts.IgnoreUtf8Errors)
	}
	if opts.PreserveSpace != nil {
		xmlMap["preserve_space"] = tfconfig.BoolVariable(*opts.PreserveSpace)
	}
	if opts.StripOuterElement != nil {
		xmlMap["strip_outer_element"] = tfconfig.BoolVariable(*opts.StripOuterElement)
	}
	if opts.DisableAutoConvert != nil {
		xmlMap["disable_auto_convert"] = tfconfig.BoolVariable(*opts.DisableAutoConvert)
	}
	if opts.ReplaceInvalidCharacters != nil {
		xmlMap["replace_invalid_characters"] = tfconfig.BoolVariable(*opts.ReplaceInvalidCharacters)
	}
	if opts.SkipByteOrderMark != nil {
		xmlMap["skip_byte_order_mark"] = tfconfig.BoolVariable(*opts.SkipByteOrderMark)
	}

	// Workaround for empty objects - Terraform requires at least one attribute
	if len(xmlMap) == 0 {
		xmlMap["any"] = tfconfig.StringVariable(string(config.SnowflakeProviderConfigSingleAttributeWorkaround))
	}

	return tfconfig.ListVariable(tfconfig.ObjectVariable(
		map[string]tfconfig.Variable{
			"xml": tfconfig.ListVariable(tfconfig.ObjectVariable(xmlMap)),
		},
	))
}
