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

func (i *InternalStageModel) WithFileFormatAvro(opts sdk.FileFormatAvroOptions) *InternalStageModel {
	return i.WithFileFormatValue(stageFileFormatAvro(opts))
}

func (i *InternalStageModel) WithFileFormatOrc(opts sdk.FileFormatOrcOptions) *InternalStageModel {
	return i.WithFileFormatValue(stageFileFormatOrc(opts))
}

func (i *InternalStageModel) WithFileFormatJson(opts sdk.FileFormatJsonOptions) *InternalStageModel {
	return i.WithFileFormatValue(stageFileFormatJson(opts))
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

func (i *InternalStageModel) WithFileFormatParquet(opts sdk.FileFormatParquetOptions) *InternalStageModel {
	return i.WithFileFormatValue(stageFileFormatParquet(opts))
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

func (i *InternalStageModel) WithFileFormatXml(opts sdk.FileFormatXmlOptions) *InternalStageModel {
	return i.WithFileFormatValue(stageFileFormatXml(opts))
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
