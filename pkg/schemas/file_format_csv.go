package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DescribeFileFormatCsvSchema represents output of DESCRIBE query for the single CSV FileFormat.
var DescribeFileFormatCsvSchema = map[string]*schema.Schema{
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
	"skip_blank_lines": {
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
		Elem:     &schema.Schema{Type: schema.TypeString},
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
}

var _ = DescribeFileFormatCsvSchema

// FileFormatCsvToSchema converts the SDK details for a CSV file format into the DescribeOutputAttributeName schema,
// reusing the field mapping already defined for stages and adding the file format's own id.
func FileFormatCsvToSchema(fileFormatCsv *sdk.FileFormatCsv) map[string]any {
	fileFormatCsvSchema := StageFileFormatCsvToSchema(fileFormatCsv)
	fileFormatCsvSchema["id"] = fileFormatCsv.Id.FullyQualifiedName()
	return fileFormatCsvSchema
}

var _ = FileFormatCsvToSchema
