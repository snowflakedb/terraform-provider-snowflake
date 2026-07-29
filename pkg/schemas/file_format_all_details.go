package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

// DescribeFileFormatAllDetailsSchema represents output of DESCRIBE query for a file format of any type.
// Because DESCRIBE FILE FORMAT returns a different set of properties for every file format type, this schema is
// a union of all the per-type schemas; only the fields applicable to the given file format type are filled.
var DescribeFileFormatAllDetailsSchema = collections.MergeMaps(DescribeFileFormatCsvSchema, DescribeFileFormatJsonSchema, DescribeFileFormatAvroSchema, DescribeFileFormatOrcSchema, DescribeFileFormatParquetSchema, DescribeFileFormatXmlSchema)

var _ = DescribeFileFormatAllDetailsSchema

// FileFormatAllDetailsToSchema converts the SDK details for any file format type into the DescribeOutputAttributeName
// schema of the file formats data source, reusing the per-type mappers. Fields that do not apply to the given file
// format type are not set (Terraform fills them with zero values).
func FileFormatAllDetailsToSchema(details sdk.FileFormatAllDetails) map[string]any {
	switch {
	case details.Csv != nil:
		return FileFormatCsvToSchema(details.Csv)
	case details.Json != nil:
		return FileFormatJsonToSchema(details.Json)
	case details.Avro != nil:
		return FileFormatAvroToSchema(details.Avro)
	case details.Orc != nil:
		return FileFormatOrcToSchema(details.Orc)
	case details.Parquet != nil:
		return FileFormatParquetToSchema(details.Parquet)
	case details.Xml != nil:
		return FileFormatXmlToSchema(details.Xml)
	default:
		return map[string]any{
			"id":   details.Id.FullyQualifiedName(),
			"type": string(details.Type),
		}
	}
}

var _ = FileFormatAllDetailsToSchema
