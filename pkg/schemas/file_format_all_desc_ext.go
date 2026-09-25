package schemas

import (
	"maps"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (fileFormatAllDetailsToSchemaMapper) additionalSchema() map[string]*schema.Schema {
	out := collections.MergeMaps(
		DescribeFileFormatCsvSchema,
		DescribeFileFormatJsonSchema,
		DescribeFileFormatAvroSchema,
		DescribeFileFormatOrcSchema,
		DescribeFileFormatParquetSchema,
		DescribeFileFormatXmlSchema,
	)
	delete(out, "id")
	delete(out, "type")
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
