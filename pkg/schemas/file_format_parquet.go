package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DescribeFileFormatParquetSchema represents output of DESCRIBE query for the single Parquet FileFormat.
var DescribeFileFormatParquetSchema = map[string]*schema.Schema{
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
		Elem:     &schema.Schema{Type: schema.TypeString},
		Computed: true,
	},
}

var _ = DescribeFileFormatParquetSchema

// FileFormatParquetToSchema converts the SDK details for a Parquet file format into the DescribeOutputAttributeName schema,
// reusing the field mapping already defined for stages and adding the file format's own id.
func FileFormatParquetToSchema(fileFormatParquet *sdk.FileFormatParquet) map[string]any {
	fileFormatParquetSchema := StageFileFormatParquetToSchema(fileFormatParquet)
	fileFormatParquetSchema["id"] = fileFormatParquet.Id.FullyQualifiedName()
	return fileFormatParquetSchema
}

var _ = FileFormatParquetToSchema
