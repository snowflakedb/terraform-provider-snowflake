package model

import (
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

func ExternalGcsStageWithId(id sdk.SchemaObjectIdentifier, storageIntegration, url string) *ExternalGcsStageModel {
	return ExternalGcsStage("test", id.DatabaseName(), id.SchemaName(), id.Name(), storageIntegration, url)
}

func (e *ExternalGcsStageModel) WithDirectoryEnabled(enable string) *ExternalGcsStageModel {
	e.Directory = tfconfig.ListVariable(tfconfig.ObjectVariable(
		map[string]tfconfig.Variable{
			"enable": tfconfig.StringVariable(enable),
		},
	))
	return e
}

func (e *ExternalGcsStageModel) WithDirectoryEnabledAndOptions(opts sdk.ExternalGCSDirectoryTableOptionsRequest) *ExternalGcsStageModel {
	directoryMap := map[string]tfconfig.Variable{
		"enable": tfconfig.BoolVariable(opts.Enable),
	}
	if opts.RefreshOnCreate != nil {
		directoryMap["refresh_on_create"] = tfconfig.StringVariable(strconv.FormatBool(*opts.RefreshOnCreate))
	}
	if opts.AutoRefresh != nil {
		directoryMap["auto_refresh"] = tfconfig.StringVariable(strconv.FormatBool(*opts.AutoRefresh))
	}
	if opts.NotificationIntegration != nil {
		directoryMap["notification_integration"] = tfconfig.StringVariable(*opts.NotificationIntegration)
	}
	e.Directory = tfconfig.ListVariable(tfconfig.ObjectVariable(directoryMap))
	return e
}

func (e *ExternalGcsStageModel) WithInvalidAutoRefresh() *ExternalGcsStageModel {
	e.Directory = tfconfig.ListVariable(tfconfig.ObjectVariable(
		map[string]tfconfig.Variable{
			"enable":       tfconfig.BoolVariable(true),
			"auto_refresh": tfconfig.StringVariable("invalid"),
		},
	))
	return e
}

func (e *ExternalGcsStageModel) WithInvalidRefreshOnCreate() *ExternalGcsStageModel {
	e.Directory = tfconfig.ListVariable(tfconfig.ObjectVariable(
		map[string]tfconfig.Variable{
			"enable":            tfconfig.BoolVariable(true),
			"refresh_on_create": tfconfig.StringVariable("invalid"),
		},
	))
	return e
}

// WithFileFormatName sets a named file format reference.
func (e *ExternalGcsStageModel) WithFileFormatName(formatName string) *ExternalGcsStageModel {
	return e.WithFileFormatValue(
		tfconfig.ListVariable(tfconfig.ObjectVariable(
			map[string]tfconfig.Variable{
				"format_name": tfconfig.StringVariable(formatName),
			},
		)),
	)
}

// WithFileFormatCsv sets inline CSV file format with the provided options.
func (e *ExternalGcsStageModel) WithFileFormatCsv(opts sdk.FileFormatCsvOptions) *ExternalGcsStageModel {
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

	return e.WithFileFormatValue(
		tfconfig.ListVariable(tfconfig.ObjectVariable(
			map[string]tfconfig.Variable{
				"csv": tfconfig.ListVariable(tfconfig.ObjectVariable(csvMap)),
			},
		)),
	)
}

func (e *ExternalGcsStageModel) WithEncryptionGcsSseKms(kmsKeyId string) *ExternalGcsStageModel {
	encryptionMap := map[string]tfconfig.Variable{
		"gcs_sse_kms": tfconfig.ListVariable(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"kms_key_id": tfconfig.StringVariable(kmsKeyId),
		})),
	}
	return e.WithEncryptionValue(tfconfig.ListVariable(tfconfig.ObjectVariable(encryptionMap)))
}

func (e *ExternalGcsStageModel) WithEncryptionGcsSseKmsNoKey() *ExternalGcsStageModel {
	encryptionMap := map[string]tfconfig.Variable{
		"gcs_sse_kms": tfconfig.ListVariable(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"any": tfconfig.StringVariable(string(config.SnowflakeProviderConfigSingleAttributeWorkaround)),
		})),
	}
	return e.WithEncryptionValue(tfconfig.ListVariable(tfconfig.ObjectVariable(encryptionMap)))
}

func (e *ExternalGcsStageModel) WithEncryptionNone() *ExternalGcsStageModel {
	return e.WithEncryptionValue(
		tfconfig.ListVariable(tfconfig.ObjectVariable(
			map[string]tfconfig.Variable{
				"none": tfconfig.ListVariable(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
					"any": tfconfig.StringVariable(string(config.SnowflakeProviderConfigSingleAttributeWorkaround)),
				})),
			},
		)),
	)
}

func (e *ExternalGcsStageModel) WithEncryptionBothTypes() *ExternalGcsStageModel {
	return e.WithEncryptionValue(
		tfconfig.ListVariable(tfconfig.ObjectVariable(
			map[string]tfconfig.Variable{
				"gcs_sse_kms": tfconfig.ListVariable(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
					"kms_key_id": tfconfig.StringVariable("foo"),
				})),
				"none": tfconfig.ListVariable(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
					"any": tfconfig.StringVariable(string(config.SnowflakeProviderConfigSingleAttributeWorkaround)),
				})),
			},
		)),
	)
}

func (e *ExternalGcsStageModel) WithEncryptionNoneTypeSpecified() *ExternalGcsStageModel {
	return e.WithEncryptionValue(
		tfconfig.ListVariable(tfconfig.ObjectVariable(
			map[string]tfconfig.Variable{
				"any": tfconfig.StringVariable(string(config.SnowflakeProviderConfigSingleAttributeWorkaround)),
			},
		)),
	)
}
