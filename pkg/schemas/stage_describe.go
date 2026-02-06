package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var StageDescribeSchema = map[string]*schema.Schema{
	"directory_table": {
		Type: schema.TypeList,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"enable": {
					Type:     schema.TypeBool,
					Computed: true,
				},
				"auto_refresh": {
					Type:     schema.TypeBool,
					Computed: true,
				},
			},
		},
		Computed: true,
	},
	"file_format": {
		Type: schema.TypeList,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"format_name": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"csv":  csvFileFormatSchema,
				"json": jsonFileFormatSchema,
			},
		},
		Computed: true,
	},
}

var csvFileFormatSchema = &schema.Schema{
	Type:     schema.TypeList,
	Computed: true,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"record_delimiter": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"field_delimiter": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"file_extension": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"skip_header": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"parse_header": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"date_format": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"time_format": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"timestamp_format": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"binary_format": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"escape": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"escape_unenclosed_field": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"trim_space": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"field_optionally_enclosed_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"null_if": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"compression": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"error_on_column_count_mismatch": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"validate_utf8": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"skip_blank_lines": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"replace_invalid_characters": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"empty_field_as_null": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"skip_byte_order_mark": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"encoding": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"multi_line": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	},
}

var jsonFileFormatSchema = &schema.Schema{
	Type:     schema.TypeList,
	Computed: true,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"compression": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"date_format": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"time_format": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"timestamp_format": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"binary_format": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"trim_space": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"multi_line": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"null_if": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"file_extension": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enable_octal": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"allow_duplicate": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"strip_outer_array": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"strip_null_values": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"replace_invalid_characters": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"ignore_utf8_errors": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"skip_byte_order_mark": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	},
}

func StageDescribeToSchema(properties sdk.StageDetails) (map[string]any, error) {
	schema := make(map[string]any)

	if properties.DirectoryTable != nil {
		schema["directory_table"] = []map[string]any{
			{
				"enable":       properties.DirectoryTable.Enable,
				"auto_refresh": properties.DirectoryTable.AutoRefresh,
			},
		}
	}

	fileFormat := map[string]any{
		"format_name": "",
		"csv":         []any{},
		"json":        []any{},
	}
	switch {
	case properties.FileFormatName != nil:
		fileFormat["format_name"] = properties.FileFormatName.FullyQualifiedName()
	case properties.FileFormatCsv != nil:
		fileFormat["csv"] = []any{StageFileFormatCsvToSchema(properties.FileFormatCsv)}
	case properties.FileFormatJson != nil:
		fileFormat["json"] = []any{StageFileFormatJsonToSchema(properties.FileFormatJson)}
	}
	schema["file_format"] = []map[string]any{fileFormat}

	return schema, nil
}

func StageFileFormatJsonToSchema(json *sdk.FileFormatJson) map[string]any {
	return map[string]any{
		"type":                       json.Type,
		"compression":                json.Compression,
		"date_format":                json.DateFormat,
		"time_format":                json.TimeFormat,
		"timestamp_format":           json.TimestampFormat,
		"binary_format":              json.BinaryFormat,
		"trim_space":                 json.TrimSpace,
		"multi_line":                 json.MultiLine,
		"null_if":                    collections.Map(json.NullIf, func(v string) any { return v }),
		"file_extension":             json.FileExtension,
		"enable_octal":               json.EnableOctal,
		"allow_duplicate":            json.AllowDuplicate,
		"strip_outer_array":          json.StripOuterArray,
		"strip_null_values":          json.StripNullValues,
		"replace_invalid_characters": json.ReplaceInvalidCharacters,
		"ignore_utf8_errors":         json.IgnoreUtf8Errors,
		"skip_byte_order_mark":       json.SkipByteOrderMark,
	}
}

func StageFileFormatCsvToSchema(csv *sdk.FileFormatCsv) map[string]any {
	return map[string]any{
		"type":                           csv.Type,
		"record_delimiter":               csv.RecordDelimiter,
		"field_delimiter":                csv.FieldDelimiter,
		"file_extension":                 csv.FileExtension,
		"skip_header":                    csv.SkipHeader,
		"parse_header":                   csv.ParseHeader,
		"date_format":                    csv.DateFormat,
		"time_format":                    csv.TimeFormat,
		"timestamp_format":               csv.TimestampFormat,
		"binary_format":                  csv.BinaryFormat,
		"escape":                         csv.Escape,
		"escape_unenclosed_field":        csv.EscapeUnenclosedField,
		"trim_space":                     csv.TrimSpace,
		"field_optionally_enclosed_by":   csv.FieldOptionallyEnclosedBy,
		"null_if":                        collections.Map(csv.NullIf, func(v string) any { return v }),
		"compression":                    csv.Compression,
		"error_on_column_count_mismatch": csv.ErrorOnColumnCountMismatch,
		"validate_utf8":                  csv.ValidateUtf8,
		"skip_blank_lines":               csv.SkipBlankLines,
		"replace_invalid_characters":     csv.ReplaceInvalidCharacters,
		"empty_field_as_null":            csv.EmptyFieldAsNull,
		"skip_byte_order_mark":           csv.SkipByteOrderMark,
		"encoding":                       csv.Encoding,
		"multi_line":                     csv.MultiLine,
	}
}
