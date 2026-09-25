package schemas

import (
	"maps"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TODO [next PRs]: first move only. Generated Stage DESCRIBE schemas are empty because every
// public key is a nested struct/slice (unmapped today); this ext owns the whole describe_output.
// Native slices/nested structs should generate those keys and shrink this file.

func stageDirectoryTableDescribeSchema() *schema.Schema {
	return &schema.Schema{
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
				"last_refreshed_on": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
		Computed: true,
	}
}

func stageFileFormatDescribeSchema() *schema.Schema {
	return &schema.Schema{
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
	}
}

func stageCommonDescribeSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"directory_table": stageDirectoryTableDescribeSchema(),
		"file_format":     stageFileFormatDescribeSchema(),
	}
}

func stageAwsLocationDescribeSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"url": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"aws_access_point_arn": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func stageCompatibleLocationDescribeSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"url": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
		},
	}
}

func stagePrivateLinkDescribeSchema() *schema.Schema {
	return &schema.Schema{
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
	}
}

func mapStageCommonDescribe(src *sdk.StageCommon, dst map[string]any) {
	if src.DirectoryTable != nil {
		dst["directory_table"] = []map[string]any{
			{
				"enable":            src.DirectoryTable.Enable,
				"auto_refresh":      src.DirectoryTable.AutoRefresh,
				"last_refreshed_on": src.DirectoryTable.LastRefreshedOn,
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
	case src.FileFormatName != nil:
		fileFormat["format_name"] = src.FileFormatName.FullyQualifiedName()
	case src.FileFormatCsv != nil:
		fileFormat["csv"] = []any{stageFileFormatCsvToSchema(src.FileFormatCsv)}
	case src.FileFormatJson != nil:
		fileFormat["json"] = []any{stageFileFormatJsonToSchema(src.FileFormatJson)}
	case src.FileFormatAvro != nil:
		fileFormat["avro"] = []any{stageFileFormatAvroToSchema(src.FileFormatAvro)}
	case src.FileFormatOrc != nil:
		fileFormat["orc"] = []any{stageFileFormatOrcToSchema(src.FileFormatOrc)}
	case src.FileFormatParquet != nil:
		fileFormat["parquet"] = []any{stageFileFormatParquetToSchema(src.FileFormatParquet)}
	case src.FileFormatXml != nil:
		fileFormat["xml"] = []any{stageFileFormatXmlToSchema(src.FileFormatXml)}
	}
	dst["file_format"] = []map[string]any{fileFormat}
}

func stageCommonFromAws(src *sdk.StageAws) *sdk.StageCommon {
	return &sdk.StageCommon{
		FileFormatName:    src.FileFormatName,
		FileFormatCsv:     src.FileFormatCsv,
		FileFormatJson:    src.FileFormatJson,
		FileFormatAvro:    src.FileFormatAvro,
		FileFormatOrc:     src.FileFormatOrc,
		FileFormatParquet: src.FileFormatParquet,
		FileFormatXml:     src.FileFormatXml,
		DirectoryTable:    src.DirectoryTable,
	}
}

func stageCommonFromAwsCompatible(src *sdk.StageAwsCompatible) *sdk.StageCommon {
	return &sdk.StageCommon{
		FileFormatName:    src.FileFormatName,
		FileFormatCsv:     src.FileFormatCsv,
		FileFormatJson:    src.FileFormatJson,
		FileFormatAvro:    src.FileFormatAvro,
		FileFormatOrc:     src.FileFormatOrc,
		FileFormatParquet: src.FileFormatParquet,
		FileFormatXml:     src.FileFormatXml,
		DirectoryTable:    src.DirectoryTable,
	}
}

func mapStagePrivateLink(pl *sdk.StagePrivateLink, dst map[string]any) {
	if pl == nil {
		return
	}
	dst["privatelink"] = []map[string]any{
		{
			"use_privatelink_endpoint": pl.UsePrivatelinkEndpoint,
		},
	}
}

func mapStageAwsLocation(loc *sdk.StageLocationDetails, dst map[string]any) {
	if loc == nil {
		return
	}
	dst["location"] = []map[string]any{
		{
			"url":                  loc.Url,
			"aws_access_point_arn": loc.AwsAccessPointArn,
		},
	}
}

func mapStageCompatibleLocation(loc *sdk.StageLocationDetails, dst map[string]any) {
	if loc == nil {
		return
	}
	dst["location"] = []map[string]any{
		{
			"url": loc.Url,
		},
	}
}

func stageAwsDescribeSchema() map[string]*schema.Schema {
	return collections.MergeMaps(stageCommonDescribeSchema(), map[string]*schema.Schema{
		"privatelink": stagePrivateLinkDescribeSchema(),
		"location":    stageAwsLocationDescribeSchema(),
	})
}

func (stageDetailsToSchemaMapper) additionalSchema() map[string]*schema.Schema {
	// Union / stages DS matches AWS: it is the public-key superset (common + privatelink + location).
	return stageAwsDescribeSchema()
}

func (stageDetailsToSchemaMapper) additionalToSchema(src *sdk.StageDetails, dst map[string]any) {
	maps.Copy(dst, StageAwsToSchema(src.AsAws()))
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

func stageFileFormatJsonToSchema(json *sdk.FileFormatJson) map[string]any {
	return map[string]any{
		"type":                       string(json.Type),
		"compression":                string(json.Compression),
		"date_format":                json.DateFormat,
		"time_format":                json.TimeFormat,
		"timestamp_format":           json.TimestampFormat,
		"binary_format":              string(json.BinaryFormat),
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

func stageFileFormatCsvToSchema(csv *sdk.FileFormatCsv) map[string]any {
	return map[string]any{
		"type":                           string(csv.Type),
		"record_delimiter":               csv.RecordDelimiter,
		"field_delimiter":                csv.FieldDelimiter,
		"file_extension":                 csv.FileExtension,
		"skip_header":                    csv.SkipHeader,
		"parse_header":                   csv.ParseHeader,
		"date_format":                    csv.DateFormat,
		"time_format":                    csv.TimeFormat,
		"timestamp_format":               csv.TimestampFormat,
		"binary_format":                  string(csv.BinaryFormat),
		"escape":                         csv.Escape,
		"escape_unenclosed_field":        csv.EscapeUnenclosedField,
		"trim_space":                     csv.TrimSpace,
		"field_optionally_enclosed_by":   csv.FieldOptionallyEnclosedBy,
		"null_if":                        collections.Map(csv.NullIf, func(v string) any { return v }),
		"compression":                    string(csv.Compression),
		"error_on_column_count_mismatch": csv.ErrorOnColumnCountMismatch,
		"validate_utf8":                  csv.ValidateUtf8,
		"skip_blank_lines":               csv.SkipBlankLines,
		"replace_invalid_characters":     csv.ReplaceInvalidCharacters,
		"empty_field_as_null":            csv.EmptyFieldAsNull,
		"skip_byte_order_mark":           csv.SkipByteOrderMark,
		"encoding":                       string(csv.Encoding),
		"multi_line":                     csv.MultiLine,
	}
}

func stageFileFormatAvroToSchema(avro *sdk.FileFormatAvro) map[string]any {
	return map[string]any{
		"type":                       string(avro.Type),
		"compression":                string(avro.Compression),
		"trim_space":                 avro.TrimSpace,
		"replace_invalid_characters": avro.ReplaceInvalidCharacters,
		"null_if":                    collections.Map(avro.NullIf, func(v string) any { return v }),
	}
}

func stageFileFormatOrcToSchema(orc *sdk.FileFormatOrc) map[string]any {
	return map[string]any{
		"type":                       string(orc.Type),
		"trim_space":                 orc.TrimSpace,
		"replace_invalid_characters": orc.ReplaceInvalidCharacters,
		"null_if":                    collections.Map(orc.NullIf, func(v string) any { return v }),
	}
}

func stageFileFormatParquetToSchema(parquet *sdk.FileFormatParquet) map[string]any {
	return map[string]any{
		"type":                       string(parquet.Type),
		"compression":                string(parquet.Compression),
		"binary_as_text":             parquet.BinaryAsText,
		"use_logical_type":           parquet.UseLogicalType,
		"trim_space":                 parquet.TrimSpace,
		"use_vectorized_scanner":     parquet.UseVectorizedScanner,
		"replace_invalid_characters": parquet.ReplaceInvalidCharacters,
		"null_if":                    collections.Map(parquet.NullIf, func(v string) any { return v }),
	}
}

func stageFileFormatXmlToSchema(xml *sdk.FileFormatXml) map[string]any {
	return map[string]any{
		"type":                       string(xml.Type),
		"compression":                string(xml.Compression),
		"ignore_utf8_errors":         xml.IgnoreUtf8Errors,
		"preserve_space":             xml.PreserveSpace,
		"strip_outer_element":        xml.StripOuterElement,
		"disable_auto_convert":       xml.DisableAutoConvert,
		"replace_invalid_characters": xml.ReplaceInvalidCharacters,
		"skip_byte_order_mark":       xml.SkipByteOrderMark,
		// disable_snowflake_data is deprecated in Snowflake
	}
}
