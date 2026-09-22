package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (icebergTableDetailsToSchemaMapper) additionalSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func (icebergTableDetailsToSchemaMapper) additionalToSchema(src *sdk.IcebergTableDetails, dst map[string]any) {
	dst["type"] = src.TypeString()
}

func icebergTablePartitionSpecsToSchema(partitionSpecs []sdk.IcebergTablePartitionSpec) []map[string]any {
	result := make([]map[string]any, len(partitionSpecs))
	for i, spec := range partitionSpecs {
		fields := make([]map[string]any, len(spec.Fields))
		for j, field := range spec.Fields {
			fields[j] = map[string]any{
				"name":      field.Name,
				"transform": field.Transform,
				"source_id": field.SourceId,
				"field_id":  field.FieldId,
			}
		}
		result[i] = map[string]any{
			"spec_id": spec.SpecId,
			"fields":  fields,
		}
	}
	return result
}

func IcebergTableDetailsListToSchema(details []sdk.IcebergTableDetails) []map[string]any {
	result := make([]map[string]any, len(details))
	for i := range details {
		result[i] = IcebergTableDetailsToSchema(&details[i])
	}
	return result
}
