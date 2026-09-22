package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (icebergTableToSchemaMapper) additionalSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// auto_refresh_status is returned as a JSON blob and parsed into a struct in the SDK,
		// so it is represented here as a nested object instead of the generator's default string field.
		"auto_refresh_status": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"current_snapshot_id": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"last_snapshot_time": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"pending_snapshot_count": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"execution_state": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"last_updated_time": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		// partition_specs is returned as a JSON blob and parsed into a struct slice in the
		// SDK, so it is represented here as a nested object list instead of the generator's default string field.
		"partition_specs": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"spec_id": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"fields": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"name": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"transform": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"source_id": {
									Type:     schema.TypeInt,
									Computed: true,
								},
								"field_id": {
									Type:     schema.TypeInt,
									Computed: true,
								},
							},
						},
					},
				},
			},
		},
	}
}

func (icebergTableToSchemaMapper) additionalToSchema(src *sdk.IcebergTable, dst map[string]any) {
	// Map the parsed auto_refresh_status struct into a nested object (empty list when absent).
	if src.AutoRefreshStatus != nil {
		autoRefreshStatus := map[string]any{
			"current_snapshot_id":    src.AutoRefreshStatus.CurrentSnapshotId,
			"pending_snapshot_count": src.AutoRefreshStatus.PendingSnapshotCount,
			"execution_state":        src.AutoRefreshStatus.ExecutionState,
			"last_updated_time":      src.AutoRefreshStatus.LastUpdatedTime,
		}
		if src.AutoRefreshStatus.LastSnapshotTime != nil {
			autoRefreshStatus["last_snapshot_time"] = *src.AutoRefreshStatus.LastSnapshotTime
		}
		dst["auto_refresh_status"] = []map[string]any{autoRefreshStatus}
	}
	dst["partition_specs"] = icebergTablePartitionSpecsToSchema(src.PartitionSpecs)
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
