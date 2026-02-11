package model

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

func InternalStageWithId(id sdk.SchemaObjectIdentifier) *InternalStageModel {
	return InternalStage("test", id.DatabaseName(), id.SchemaName(), id.Name())
}

func (i *InternalStageModel) WithDirectoryEnabled(enable string) *InternalStageModel {
	i.Directory = tfconfig.ListVariable(tfconfig.ObjectVariable(
		map[string]tfconfig.Variable{
			"enable": tfconfig.StringVariable(enable),
		},
	))
	return i
}

func (i *InternalStageModel) WithDirectoryEnabledAndAutoRefresh(enable bool, autoRefresh string) *InternalStageModel {
	i.Directory = tfconfig.ListVariable(tfconfig.ObjectVariable(
		map[string]tfconfig.Variable{
			"enable":       tfconfig.BoolVariable(enable),
			"auto_refresh": tfconfig.StringVariable(autoRefresh),
		},
	))
	return i
}

func (i *InternalStageModel) WithEncryptionSnowflakeFull() *InternalStageModel {
	return i.WithEncryptionValue(
		tfconfig.ListVariable(tfconfig.ObjectVariable(
			map[string]tfconfig.Variable{
				"snowflake_full": tfconfig.ListVariable(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
					"any": tfconfig.StringVariable(string(config.SnowflakeProviderConfigSingleAttributeWorkaround)),
				})),
			},
		)),
	)
}

func (i *InternalStageModel) WithEncryptionSnowflakeSse() *InternalStageModel {
	return i.WithEncryptionValue(
		tfconfig.ListVariable(tfconfig.ObjectVariable(
			map[string]tfconfig.Variable{
				"snowflake_sse": tfconfig.ListVariable(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
					"any": tfconfig.StringVariable(string(config.SnowflakeProviderConfigSingleAttributeWorkaround)),
				})),
			},
		)),
	)
}

func (i *InternalStageModel) WithEncryptionBothTypes() *InternalStageModel {
	return i.WithEncryptionValue(
		tfconfig.ListVariable(tfconfig.ObjectVariable(
			map[string]tfconfig.Variable{
				"snowflake_full": tfconfig.ListVariable(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
					"any": tfconfig.StringVariable(string(config.SnowflakeProviderConfigSingleAttributeWorkaround)),
				})),
				"snowflake_sse": tfconfig.ListVariable(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
					"any": tfconfig.StringVariable(string(config.SnowflakeProviderConfigSingleAttributeWorkaround)),
				})),
			},
		)),
	)
}

// WithFileFormatName sets a named file format reference.
func (e *InternalStageModel) WithFileFormatName(formatName string) *InternalStageModel {
	return e.WithFileFormatValue(stageFileFormatName(formatName))
}

// WithFileFormatCsv sets inline CSV file format with the provided options.
func (e *InternalStageModel) WithFileFormatCsv(opts sdk.FileFormatCsvOptions) *InternalStageModel {
	return e.WithFileFormatValue(stageFileFormatCsv(opts))
}

func (i *InternalStageModel) WithFileFormatCsvConflictingOptions() *InternalStageModel {
	return i.WithFileFormatCsv(sdk.FileFormatCsvOptions{
		SkipHeader:  sdk.Pointer(1),
		ParseHeader: sdk.Pointer(true),
	})
}

func (i *InternalStageModel) WithFileFormatCsvInvalidSkipHeader() *InternalStageModel {
	return i.WithFileFormatCsv(sdk.FileFormatCsvOptions{
		SkipHeader: sdk.Pointer(-1),
	})
}

func (i *InternalStageModel) WithFileFormatInvalidFormatName() *InternalStageModel {
	return i.WithFileFormatValue(
		tfconfig.ListVariable(tfconfig.ObjectVariable(
			map[string]tfconfig.Variable{
				"format_name": tfconfig.StringVariable("invalid"),
			},
		)),
	)
}

func (i *InternalStageModel) WithFileFormatCsvInvalidEncoding() *InternalStageModel {
	return i.WithFileFormatCsv(sdk.FileFormatCsvOptions{
		Encoding: sdk.Pointer(sdk.CsvEncoding("INVALID")),
	})
}

func (i *InternalStageModel) WithFileFormatCsvInvalidBooleanString() *InternalStageModel {
	return i.WithFileFormatValue(
		tfconfig.ListVariable(tfconfig.ObjectVariable(
			map[string]tfconfig.Variable{
				"csv": tfconfig.ListVariable(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
					"multi_line": tfconfig.StringVariable("invalid"),
				})),
			},
		)),
	)
}

func (i *InternalStageModel) WithFileFormatCsvInvalidBinaryFormat() *InternalStageModel {
	return i.WithFileFormatCsv(sdk.FileFormatCsvOptions{
		BinaryFormat: sdk.Pointer(sdk.BinaryFormat("INVALID")),
	})
}

func (i *InternalStageModel) WithFileFormatCsvInvalidCompression() *InternalStageModel {
	return i.WithFileFormatCsv(sdk.FileFormatCsvOptions{
		Compression: sdk.Pointer(sdk.CsvCompression("INVALID")),
	})
}

func (i *InternalStageModel) WithFileFormatMultipleFormats() *InternalStageModel {
	return i.WithFileFormatValue(
		tfconfig.ListVariable(tfconfig.ObjectVariable(
			map[string]tfconfig.Variable{
				"format_name": tfconfig.StringVariable("some_format"),
				"csv": tfconfig.ListVariable(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
					"field_delimiter": tfconfig.StringVariable(","),
				})),
			},
		)),
	)
}

// WithFileFormatJson sets inline JSON file format with the provided options.
func (i *InternalStageModel) WithFileFormatJson(opts sdk.FileFormatJsonOptions) *InternalStageModel {
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

	return i.WithFileFormatValue(
		tfconfig.ListVariable(tfconfig.ObjectVariable(
			map[string]tfconfig.Variable{
				"json": tfconfig.ListVariable(tfconfig.ObjectVariable(jsonMap)),
			},
		)),
	)
}

func (i *InternalStageModel) WithFileFormatJsonInvalidCompression() *InternalStageModel {
	return i.WithFileFormatJson(sdk.FileFormatJsonOptions{
		Compression: sdk.Pointer(sdk.JsonCompression("INVALID")),
	})
}

func (i *InternalStageModel) WithFileFormatJsonInvalidBinaryFormat() *InternalStageModel {
	return i.WithFileFormatJson(sdk.FileFormatJsonOptions{
		BinaryFormat: sdk.Pointer(sdk.BinaryFormat("INVALID")),
	})
}

func (i *InternalStageModel) WithFileFormatJsonInvalidBooleanString() *InternalStageModel {
	return i.WithFileFormatValue(
		tfconfig.ListVariable(tfconfig.ObjectVariable(
			map[string]tfconfig.Variable{
				"json": tfconfig.ListVariable(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
					"multi_line": tfconfig.StringVariable("invalid"),
				})),
			},
		)),
	)
}

func (i *InternalStageModel) WithFileFormatJsonConflictingOptions() *InternalStageModel {
	return i.WithFileFormatJson(sdk.FileFormatJsonOptions{
		ReplaceInvalidCharacters: sdk.Pointer(true),
		IgnoreUtf8Errors:         sdk.Pointer(true),
	})
}

// WithFileFormatAvro sets inline AVRO file format with the provided options.
func (i *InternalStageModel) WithFileFormatAvro(opts sdk.FileFormatAvroOptions) *InternalStageModel {
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

	return i.WithFileFormatValue(
		tfconfig.ListVariable(tfconfig.ObjectVariable(
			map[string]tfconfig.Variable{
				"avro": tfconfig.ListVariable(tfconfig.ObjectVariable(avroMap)),
			},
		)),
	)
}

// WithFileFormatOrc sets inline ORC file format with the provided options.
func (i *InternalStageModel) WithFileFormatOrc(opts sdk.FileFormatOrcOptions) *InternalStageModel {
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

	return i.WithFileFormatValue(
		tfconfig.ListVariable(tfconfig.ObjectVariable(
			map[string]tfconfig.Variable{
				"orc": tfconfig.ListVariable(tfconfig.ObjectVariable(orcMap)),
			},
		)),
	)
}

// WithFileFormatParquet sets inline Parquet file format with the provided options.
func (i *InternalStageModel) WithFileFormatParquet(opts sdk.FileFormatParquetOptions) *InternalStageModel {
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

	return i.WithFileFormatValue(
		tfconfig.ListVariable(tfconfig.ObjectVariable(
			map[string]tfconfig.Variable{
				"parquet": tfconfig.ListVariable(tfconfig.ObjectVariable(parquetMap)),
			},
		)),
	)
}

func (i *InternalStageModel) WithFileFormatParquetInvalidCompression() *InternalStageModel {
	return i.WithFileFormatValue(
		tfconfig.ListVariable(tfconfig.ObjectVariable(
			map[string]tfconfig.Variable{
				"parquet": tfconfig.ListVariable(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
					"compression": tfconfig.StringVariable("INVALID"),
				})),
			},
		)),
	)
}

func (i *InternalStageModel) WithFileFormatParquetInvalidBooleanString() *InternalStageModel {
	return i.WithFileFormatValue(
		tfconfig.ListVariable(tfconfig.ObjectVariable(
			map[string]tfconfig.Variable{
				"parquet": tfconfig.ListVariable(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
					"trim_space": tfconfig.StringVariable("invalid"),
				})),
			},
		)),
	)
}

// WithFileFormatXml sets inline XML file format with the provided options.
func (i *InternalStageModel) WithFileFormatXml(opts sdk.FileFormatXmlOptions) *InternalStageModel {
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

	return i.WithFileFormatValue(
		tfconfig.ListVariable(tfconfig.ObjectVariable(
			map[string]tfconfig.Variable{
				"xml": tfconfig.ListVariable(tfconfig.ObjectVariable(xmlMap)),
			},
		)),
	)
}

func (i *InternalStageModel) WithFileFormatXmlConflictingOptions() *InternalStageModel {
	return i.WithFileFormatXml(sdk.FileFormatXmlOptions{
		ReplaceInvalidCharacters: sdk.Pointer(true),
		IgnoreUtf8Errors:         sdk.Pointer(true),
	})
}

func (i *InternalStageModel) WithFileFormatXmlInvalidCompression() *InternalStageModel {
	return i.WithFileFormatValue(
		tfconfig.ListVariable(tfconfig.ObjectVariable(
			map[string]tfconfig.Variable{
				"xml": tfconfig.ListVariable(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
					"compression": tfconfig.StringVariable("INVALID"),
				})),
			},
		)),
	)
}

func (i *InternalStageModel) WithFileFormatXmlInvalidBooleanString() *InternalStageModel {
	return i.WithFileFormatValue(
		tfconfig.ListVariable(tfconfig.ObjectVariable(
			map[string]tfconfig.Variable{
				"xml": tfconfig.ListVariable(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
					"preserve_space": tfconfig.StringVariable("invalid"),
				})),
			},
		)),
	)
}
