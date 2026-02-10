package model

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

// CsvFileFormatOptions holds CSV file format configuration options.
type CsvFileFormatOptions struct {
	Compression                sdk.CsvCompression
	RecordDelimiter            string
	FieldDelimiter             string
	MultiLine                  *bool
	FileExtension              string
	ParseHeader                *bool
	SkipHeader                 *int
	SkipBlankLines             *bool
	DateFormat                 string
	TimeFormat                 string
	TimestampFormat            string
	BinaryFormat               sdk.BinaryFormat
	Escape                     string
	EscapeUnenclosedField      string
	TrimSpace                  *bool
	FieldOptionallyEnclosedBy  string
	NullIf                     []string
	ErrorOnColumnCountMismatch *bool
	ReplaceInvalidCharacters   *bool
	EmptyFieldAsNull           *bool
	SkipByteOrderMark          *bool
	Encoding                   sdk.CsvEncoding
}

// JsonFileFormatOptions holds JSON file format configuration options.
type JsonFileFormatOptions struct {
	Compression              sdk.JsonCompression
	DateFormat               string
	TimeFormat               string
	TimestampFormat          string
	BinaryFormat             sdk.BinaryFormat
	TrimSpace                *bool
	MultiLine                *bool
	NullIf                   []string
	FileExtension            string
	EnableOctal              *bool
	AllowDuplicate           *bool
	StripOuterArray          *bool
	StripNullValues          *bool
	ReplaceInvalidCharacters *bool
	IgnoreUtf8Errors         *bool
	SkipByteOrderMark        *bool
}

// AvroFileFormatOptions holds AVRO file format configuration options.
type AvroFileFormatOptions struct {
	Compression              sdk.AvroCompression
	TrimSpace                *bool
	ReplaceInvalidCharacters *bool
	NullIf                   []string
}

// OrcFileFormatOptions holds ORC file format configuration options.
type OrcFileFormatOptions struct {
	TrimSpace                *bool
	ReplaceInvalidCharacters *bool
	NullIf                   []string
}

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
func (i *InternalStageModel) WithFileFormatName(formatName string) *InternalStageModel {
	return i.WithFileFormatValue(
		tfconfig.ListVariable(tfconfig.ObjectVariable(
			map[string]tfconfig.Variable{
				"format_name": tfconfig.StringVariable(formatName),
			},
		)),
	)
}

// WithFileFormatCsv sets inline CSV file format with the provided options.
func (i *InternalStageModel) WithFileFormatCsv(opts CsvFileFormatOptions) *InternalStageModel {
	csvMap := make(map[string]tfconfig.Variable)

	if opts.Compression != "" {
		csvMap["compression"] = tfconfig.StringVariable(string(opts.Compression))
	}
	if opts.RecordDelimiter != "" {
		csvMap["record_delimiter"] = tfconfig.StringVariable(opts.RecordDelimiter)
	}
	if opts.FieldDelimiter != "" {
		csvMap["field_delimiter"] = tfconfig.StringVariable(opts.FieldDelimiter)
	}
	if opts.MultiLine != nil {
		csvMap["multi_line"] = tfconfig.BoolVariable(*opts.MultiLine)
	}
	if opts.FileExtension != "" {
		csvMap["file_extension"] = tfconfig.StringVariable(opts.FileExtension)
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
	if opts.DateFormat != "" {
		csvMap["date_format"] = tfconfig.StringVariable(opts.DateFormat)
	}
	if opts.TimeFormat != "" {
		csvMap["time_format"] = tfconfig.StringVariable(opts.TimeFormat)
	}
	if opts.TimestampFormat != "" {
		csvMap["timestamp_format"] = tfconfig.StringVariable(opts.TimestampFormat)
	}
	if opts.BinaryFormat != "" {
		csvMap["binary_format"] = tfconfig.StringVariable(string(opts.BinaryFormat))
	}
	if opts.Escape != "" {
		csvMap["escape"] = tfconfig.StringVariable(opts.Escape)
	}
	if opts.EscapeUnenclosedField != "" {
		csvMap["escape_unenclosed_field"] = tfconfig.StringVariable(opts.EscapeUnenclosedField)
	}
	if opts.TrimSpace != nil {
		csvMap["trim_space"] = tfconfig.BoolVariable(*opts.TrimSpace)
	}
	if opts.FieldOptionallyEnclosedBy != "" {
		csvMap["field_optionally_enclosed_by"] = tfconfig.StringVariable(opts.FieldOptionallyEnclosedBy)
	}
	if len(opts.NullIf) > 0 {
		nullIfVars := make([]tfconfig.Variable, len(opts.NullIf))
		for idx, v := range opts.NullIf {
			nullIfVars[idx] = tfconfig.StringVariable(v)
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
	if opts.Encoding != "" {
		csvMap["encoding"] = tfconfig.StringVariable(string(opts.Encoding))
	}

	// Workaround for empty objects - Terraform requires at least one attribute
	if len(csvMap) == 0 {
		csvMap["any"] = tfconfig.StringVariable(string(config.SnowflakeProviderConfigSingleAttributeWorkaround))
	}

	return i.WithFileFormatValue(
		tfconfig.ListVariable(tfconfig.ObjectVariable(
			map[string]tfconfig.Variable{
				"csv": tfconfig.ListVariable(tfconfig.ObjectVariable(csvMap)),
			},
		)),
	)
}

func (i *InternalStageModel) WithFileFormatCsvConflictingOptions() *InternalStageModel {
	return i.WithFileFormatCsv(CsvFileFormatOptions{
		SkipHeader:  sdk.Pointer(1),
		ParseHeader: sdk.Pointer(true),
	})
}

func (i *InternalStageModel) WithFileFormatCsvInvalidSkipHeader() *InternalStageModel {
	return i.WithFileFormatCsv(CsvFileFormatOptions{
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
	return i.WithFileFormatCsv(CsvFileFormatOptions{
		Encoding: "INVALID",
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
	return i.WithFileFormatCsv(CsvFileFormatOptions{
		BinaryFormat: "INVALID",
	})
}

func (i *InternalStageModel) WithFileFormatCsvInvalidCompression() *InternalStageModel {
	return i.WithFileFormatCsv(CsvFileFormatOptions{
		Compression: "INVALID",
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
func (i *InternalStageModel) WithFileFormatJson(opts JsonFileFormatOptions) *InternalStageModel {
	jsonMap := make(map[string]tfconfig.Variable)

	if opts.Compression != "" {
		jsonMap["compression"] = tfconfig.StringVariable(string(opts.Compression))
	}
	if opts.DateFormat != "" {
		jsonMap["date_format"] = tfconfig.StringVariable(opts.DateFormat)
	}
	if opts.TimeFormat != "" {
		jsonMap["time_format"] = tfconfig.StringVariable(opts.TimeFormat)
	}
	if opts.TimestampFormat != "" {
		jsonMap["timestamp_format"] = tfconfig.StringVariable(opts.TimestampFormat)
	}
	if opts.BinaryFormat != "" {
		jsonMap["binary_format"] = tfconfig.StringVariable(string(opts.BinaryFormat))
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
			nullIfVars[idx] = tfconfig.StringVariable(v)
		}
		jsonMap["null_if"] = tfconfig.ListVariable(nullIfVars...)
	}
	if opts.FileExtension != "" {
		jsonMap["file_extension"] = tfconfig.StringVariable(opts.FileExtension)
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
	return i.WithFileFormatJson(JsonFileFormatOptions{
		Compression: "INVALID",
	})
}

func (i *InternalStageModel) WithFileFormatJsonInvalidBinaryFormat() *InternalStageModel {
	return i.WithFileFormatJson(JsonFileFormatOptions{
		BinaryFormat: "INVALID",
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
	return i.WithFileFormatJson(JsonFileFormatOptions{
		ReplaceInvalidCharacters: sdk.Pointer(true),
		IgnoreUtf8Errors:         sdk.Pointer(true),
	})
}

// WithFileFormatAvro sets inline AVRO file format with the provided options.
func (i *InternalStageModel) WithFileFormatAvro(opts AvroFileFormatOptions) *InternalStageModel {
	avroMap := make(map[string]tfconfig.Variable)

	if opts.Compression != "" {
		avroMap["compression"] = tfconfig.StringVariable(string(opts.Compression))
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
			nullIfVars[idx] = tfconfig.StringVariable(v)
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
func (i *InternalStageModel) WithFileFormatOrc(opts OrcFileFormatOptions) *InternalStageModel {
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
			nullIfVars[idx] = tfconfig.StringVariable(v)
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
