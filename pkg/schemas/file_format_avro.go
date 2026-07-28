package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DescribeFileFormatAvroSchema represents output of DESCRIBE query for the single AVRO FileFormat.
var DescribeFileFormatAvroSchema = map[string]*schema.Schema{
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

var _ = DescribeFileFormatAvroSchema

// FileFormatAvroToSchema converts the SDK details for an AVRO file format into the DescribeOutputAttributeName schema,
// reusing the field mapping already defined for stages and adding the file format's own id.
func FileFormatAvroToSchema(fileFormatAvro *sdk.FileFormatAvro) map[string]any {
	fileFormatAvroSchema := StageFileFormatAvroToSchema(fileFormatAvro)
	fileFormatAvroSchema["id"] = fileFormatAvro.Id.FullyQualifiedName()
	return fileFormatAvroSchema
}

var _ = FileFormatAvroToSchema
