package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func AwsStageDescribeSchema() map[string]*schema.Schema {
	return collections.MergeMaps(stageDescribeSchema, map[string]*schema.Schema{
		"privatelink": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"use_privatelink_endpoint": {
						Type:     schema.TypeBool,
						Computed: true,
					},
				},
			},
		},
		"location": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"url": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"aws_access_point_arn": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
	})
}

func CommonStageDescribeSchema() map[string]*schema.Schema {
	return stageDescribeSchema
}

var stageDescribeSchema = map[string]*schema.Schema{
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
				"csv": csvFileFormatSchema,
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

	if properties.FileFormatName != nil || properties.FileFormatCsv != nil {
		fileFormat := map[string]any{
			"format_name": "",
			"csv":         []any{},
		}
		if properties.FileFormatName != nil {
			fileFormat["format_name"] = properties.FileFormatName.FullyQualifiedName()
		} else if properties.FileFormatCsv != nil {
			fileFormat["csv"] = []any{StageFileFormatCsvToSchema(properties.FileFormatCsv)}
		}
		schema["file_format"] = []map[string]any{fileFormat}
	}

	return schema, nil
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

func AwsStageDescribeToSchema(properties sdk.StageDetails) (map[string]any, error) {
	schema, err := StageDescribeToSchema(properties)
	if err != nil {
		return nil, err
	}

	if properties.Location != nil {
		schema["location"] = []map[string]any{
			{
				"url":                  properties.Location.Url,
				"aws_access_point_arn": properties.Location.AwsAccessPointArn,
			},
		}
	}
	if properties.PrivateLink != nil {
		schema["privatelink"] = []map[string]any{
			{
				"use_privatelink_endpoint": properties.PrivateLink.UsePrivatelinkEndpoint,
			},
		}
	}
	return schema, nil
}
