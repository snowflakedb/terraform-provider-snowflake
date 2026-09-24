package schemas

import (
	"maps"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (fileFormatAllDetailsToSchemaMapper) additionalSchema() map[string]*schema.Schema {
	out := make(map[string]*schema.Schema)
	for _, s := range []map[string]*schema.Schema{
		DescribeFileFormatCsvSchema,
		DescribeFileFormatJsonSchema,
		DescribeFileFormatAvroSchema,
		DescribeFileFormatOrcSchema,
		DescribeFileFormatParquetSchema,
		DescribeFileFormatXmlSchema,
	} {
		for k, v := range s {
			if k == "id" || k == "type" {
				continue
			}
			out[k] = v
		}
	}
	return out
}

func (fileFormatAllDetailsToSchemaMapper) additionalToSchema(src *sdk.FileFormatAllDetails, dst map[string]any) {
	var nested map[string]any
	switch {
	case src.Csv != nil:
		nested = FileFormatCsvToSchema(src.Csv)
	case src.Json != nil:
		nested = FileFormatJsonToSchema(src.Json)
	case src.Avro != nil:
		nested = FileFormatAvroToSchema(src.Avro)
	case src.Orc != nil:
		nested = FileFormatOrcToSchema(src.Orc)
	case src.Parquet != nil:
		nested = FileFormatParquetToSchema(src.Parquet)
	case src.Xml != nil:
		nested = FileFormatXmlToSchema(src.Xml)
	default:
		return
	}
	maps.Copy(dst, nested)
}
