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

// StageDatasourceDescribeSchema is a helper function used to get the schema for the describe output of the stage data source.
// It supports all stage types and file format types.
func StageDatasourceDescribeSchema() map[string]*schema.Schema {
	// For now, only the aws stage has any additional fields in describe.
	return AwsStageDescribeSchema()
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
				"csv":     csvFileFormatSchema,
				"json":    jsonFileFormatSchema,
				"avro":    avroFileFormatSchema,
				"orc":     orcFileFormatSchema,
				"parquet": parquetFileFormatSchema,
				"xml":     xmlFileFormatSchema,
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

var avroFileFormatSchema = &schema.Schema{
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
			"trim_space": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"replace_invalid_characters": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"null_if": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	},
}

var orcFileFormatSchema = &schema.Schema{
	Type:     schema.TypeList,
	Computed: true,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"trim_space": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"replace_invalid_characters": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"null_if": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	},
}

var parquetFileFormatSchema = &schema.Schema{
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
			"binary_as_text": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"use_logical_type": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"trim_space": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"use_vectorized_scanner": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"replace_invalid_characters": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"null_if": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	},
}

var xmlFileFormatSchema = &schema.Schema{
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
			"ignore_utf8_errors": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"preserve_space": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"strip_outer_element": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"disable_auto_convert": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"replace_invalid_characters": {
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
		"avro":        []any{},
		"orc":         []any{},
		"parquet":     []any{},
		"xml":         []any{},
	}
	switch {
	case properties.FileFormatName != nil:
		fileFormat["format_name"] = properties.FileFormatName.FullyQualifiedName()
	case properties.FileFormatCsv != nil:
		fileFormat["csv"] = []any{StageFileFormatCsvToSchema(properties.FileFormatCsv)}
	case properties.FileFormatJson != nil:
		fileFormat["json"] = []any{StageFileFormatJsonToSchema(properties.FileFormatJson)}
	case properties.FileFormatAvro != nil:
		fileFormat["avro"] = []any{StageFileFormatAvroToSchema(properties.FileFormatAvro)}
	case properties.FileFormatOrc != nil:
		fileFormat["orc"] = []any{StageFileFormatOrcToSchema(properties.FileFormatOrc)}
	case properties.FileFormatParquet != nil:
		fileFormat["parquet"] = []any{StageFileFormatParquetToSchema(properties.FileFormatParquet)}
	case properties.FileFormatXml != nil:
		fileFormat["xml"] = []any{StageFileFormatXmlToSchema(properties.FileFormatXml)}
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

// StageDatasourceToDatasourceSchema is a helper function used to fill the object fields in the data source.
// It supports all stage types and file format types.
func StageDatasourceToDatasourceSchema(properties sdk.StageDetails) (map[string]any, error) {
	// For now, only the aws stage has any additional fields in describe.
	return AwsStageDescribeToSchema(properties)
}

func StageFileFormatAvroToSchema(avro *sdk.FileFormatAvro) map[string]any {
	return map[string]any{
		"type":                       avro.Type,
		"compression":                avro.Compression,
		"trim_space":                 avro.TrimSpace,
		"replace_invalid_characters": avro.ReplaceInvalidCharacters,
		"null_if":                    collections.Map(avro.NullIf, func(v string) any { return v }),
	}
}

func StageFileFormatOrcToSchema(orc *sdk.FileFormatOrc) map[string]any {
	return map[string]any{
		"type":                       orc.Type,
		"trim_space":                 orc.TrimSpace,
		"replace_invalid_characters": orc.ReplaceInvalidCharacters,
		"null_if":                    collections.Map(orc.NullIf, func(v string) any { return v }),
	}
}

func StageFileFormatParquetToSchema(parquet *sdk.FileFormatParquet) map[string]any {
	return map[string]any{
		"type":                       parquet.Type,
		"compression":                parquet.Compression,
		"binary_as_text":             parquet.BinaryAsText,
		"use_logical_type":           parquet.UseLogicalType,
		"trim_space":                 parquet.TrimSpace,
		"use_vectorized_scanner":     parquet.UseVectorizedScanner,
		"replace_invalid_characters": parquet.ReplaceInvalidCharacters,
		"null_if":                    collections.Map(parquet.NullIf, func(v string) any { return v }),
	}
}

func StageFileFormatXmlToSchema(xml *sdk.FileFormatXml) map[string]any {
	return map[string]any{
		"type":                       xml.Type,
		"compression":                xml.Compression,
		"ignore_utf8_errors":         xml.IgnoreUtf8Errors,
		"preserve_space":             xml.PreserveSpace,
		"strip_outer_element":        xml.StripOuterElement,
		"disable_auto_convert":       xml.DisableAutoConvert,
		"replace_invalid_characters": xml.ReplaceInvalidCharacters,
		"skip_byte_order_mark":       xml.SkipByteOrderMark,
	}
}
