package schemas

import (
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

// FileFormatJsonToSchema converts the SDK details for a JSON file format into the DescribeOutputAttributeName schema,
// reusing the field mapping already defined for stages and adding the file format's own id.
func FileFormatJsonToSchema(fileFormatJson *sdk.FileFormatJson) map[string]any {
	fileFormatJsonSchema := StageFileFormatJsonToSchema(fileFormatJson)
	fileFormatJsonSchema["id"] = fileFormatJson.Id.FullyQualifiedName()
	return fileFormatJsonSchema
}

var _ = FileFormatJsonToSchema
