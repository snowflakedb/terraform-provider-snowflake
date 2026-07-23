package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DescribeFileFormatOrcSchema represents output of DESCRIBE query for the single ORC FileFormat.
var DescribeFileFormatOrcSchema = map[string]*schema.Schema{
	"id": {
		Type:     schema.TypeString,
		Computed: true,
	},
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
		Elem:     &schema.Schema{Type: schema.TypeString},
		Computed: true,
	},
}

var _ = DescribeFileFormatOrcSchema

// FileFormatOrcToSchema converts the SDK details for an ORC file format into the DescribeOutputAttributeName schema,
// reusing the field mapping already defined for stages and adding the file format's own id.
func FileFormatOrcToSchema(fileFormatOrc *sdk.FileFormatOrc) map[string]any {
	fileFormatOrcSchema := StageFileFormatOrcToSchema(fileFormatOrc)
	fileFormatOrcSchema["id"] = fileFormatOrc.Id.FullyQualifiedName()
	return fileFormatOrcSchema
}

var _ = FileFormatOrcToSchema
