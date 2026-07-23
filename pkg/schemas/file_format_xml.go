package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DescribeFileFormatXmlSchema represents output of DESCRIBE query for the single XML FileFormat.
var DescribeFileFormatXmlSchema = map[string]*schema.Schema{
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
	"disable_snowflake_data": {
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
}

var _ = DescribeFileFormatXmlSchema

// FileFormatXmlToSchema converts the SDK details for an XML file format into the DescribeOutputAttributeName schema,
// reusing the field mapping already defined for stages and adding the file format's own id.
func FileFormatXmlToSchema(fileFormatXml *sdk.FileFormatXml) map[string]any {
	fileFormatXmlSchema := StageFileFormatXmlToSchema(fileFormatXml)
	fileFormatXmlSchema["id"] = fileFormatXml.Id.FullyQualifiedName()
	return fileFormatXmlSchema
}

var _ = FileFormatXmlToSchema
