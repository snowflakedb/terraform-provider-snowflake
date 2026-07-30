package schemas

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

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
