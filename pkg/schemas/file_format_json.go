package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DescribeFileFormatJsonSchema represents output of DESCRIBE query for the single JSON FileFormat.
var DescribeFileFormatJsonSchema = map[string]*schema.Schema{
	"id": {
		Type:     schema.TypeString,
		Computed: true,
	},
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
		Elem:     &schema.Schema{Type: schema.TypeString},
		Computed: true,
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
}

var _ = DescribeFileFormatJsonSchema

func FileFormatJsonToSchema(fileFormatJson *sdk.FileFormatJson) map[string]any {
	fileFormatJsonSchema := make(map[string]any)
	fileFormatJsonSchema["id"] = fileFormatJson.Id.FullyQualifiedName()
	fileFormatJsonSchema["type"] = fileFormatJson.Type
	fileFormatJsonSchema["compression"] = fileFormatJson.Compression
	fileFormatJsonSchema["date_format"] = fileFormatJson.DateFormat
	fileFormatJsonSchema["time_format"] = fileFormatJson.TimeFormat
	fileFormatJsonSchema["timestamp_format"] = fileFormatJson.TimestampFormat
	fileFormatJsonSchema["binary_format"] = fileFormatJson.BinaryFormat
	fileFormatJsonSchema["trim_space"] = fileFormatJson.TrimSpace
	fileFormatJsonSchema["multi_line"] = fileFormatJson.MultiLine
	// Adjusted manually
	fileFormatJsonSchema["null_if"] = collections.Map(fileFormatJson.NullIf, func(v string) any { return v })
	fileFormatJsonSchema["file_extension"] = fileFormatJson.FileExtension
	fileFormatJsonSchema["enable_octal"] = fileFormatJson.EnableOctal
	fileFormatJsonSchema["allow_duplicate"] = fileFormatJson.AllowDuplicate
	fileFormatJsonSchema["strip_outer_array"] = fileFormatJson.StripOuterArray
	fileFormatJsonSchema["strip_null_values"] = fileFormatJson.StripNullValues
	fileFormatJsonSchema["replace_invalid_characters"] = fileFormatJson.ReplaceInvalidCharacters
	fileFormatJsonSchema["ignore_utf8_errors"] = fileFormatJson.IgnoreUtf8Errors
	fileFormatJsonSchema["skip_byte_order_mark"] = fileFormatJson.SkipByteOrderMark
	return fileFormatJsonSchema
}

var _ = FileFormatJsonToSchema
