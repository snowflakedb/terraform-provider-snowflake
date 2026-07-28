package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// FileFormatsDatasourceDescribeSchema is a helper function used to get the schema for the describe output of the file formats data source.
// Because DESCRIBE FILE FORMAT returns a different set of properties for every file format type, the type-specific
// properties are nested under the field named after the file format type (only one of them is filled at a time).
func FileFormatsDatasourceDescribeSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"type": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"csv":     nestedDescribeFileFormatSchema(DescribeFileFormatCsvSchema),
		"json":    nestedDescribeFileFormatSchema(DescribeFileFormatJsonSchema),
		"avro":    nestedDescribeFileFormatSchema(DescribeFileFormatAvroSchema),
		"orc":     nestedDescribeFileFormatSchema(DescribeFileFormatOrcSchema),
		"parquet": nestedDescribeFileFormatSchema(DescribeFileFormatParquetSchema),
		"xml":     nestedDescribeFileFormatSchema(DescribeFileFormatXmlSchema),
	}
}

func nestedDescribeFileFormatSchema(describeSchema map[string]*schema.Schema) *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: describeSchema,
		},
	}
}

// FileFormatAllDetailsToSchema converts the SDK details for any file format type into the DescribeOutputAttributeName
// schema of the file formats data source, reusing the per-type mappers.
func FileFormatAllDetailsToSchema(details sdk.FileFormatAllDetails) map[string]any {
	detailsSchema := map[string]any{
		"id":   details.Id.FullyQualifiedName(),
		"type": string(details.Type),
	}
	switch {
	case details.Csv != nil:
		detailsSchema["csv"] = []map[string]any{FileFormatCsvToSchema(details.Csv)}
	case details.Json != nil:
		detailsSchema["json"] = []map[string]any{FileFormatJsonToSchema(details.Json)}
	case details.Avro != nil:
		detailsSchema["avro"] = []map[string]any{FileFormatAvroToSchema(details.Avro)}
	case details.Orc != nil:
		detailsSchema["orc"] = []map[string]any{FileFormatOrcToSchema(details.Orc)}
	case details.Parquet != nil:
		detailsSchema["parquet"] = []map[string]any{FileFormatParquetToSchema(details.Parquet)}
	case details.Xml != nil:
		detailsSchema["xml"] = []map[string]any{FileFormatXmlToSchema(details.Xml)}
	}
	return detailsSchema
}
