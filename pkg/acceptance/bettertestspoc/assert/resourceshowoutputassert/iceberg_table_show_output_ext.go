package resourceshowoutputassert

import (
	"fmt"
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (i *IcebergTableShowOutputAssert) HasAutoRefreshStatusEmpty() *IcebergTableShowOutputAssert {
	i.StringValueSet("auto_refresh_status.#", "0")
	return i
}

func (i *IcebergTableShowOutputAssert) HasAutoRefreshStatusNotEmpty() *IcebergTableShowOutputAssert {
	i.StringValueSet("auto_refresh_status.#", "1")
	return i
}

func (i *IcebergTableShowOutputAssert) HasAutoRefreshStatus(expected sdk.IcebergTableAutoRefreshStatus) *IcebergTableShowOutputAssert {
	i.StringValueSet("auto_refresh_status.#", "1")
	i.StringValueSet("auto_refresh_status.0.current_snapshot_id", strconv.Itoa(expected.CurrentSnapshotId))
	i.StringValueSet("auto_refresh_status.0.pending_snapshot_count", strconv.Itoa(expected.PendingSnapshotCount))
	i.StringValueSet("auto_refresh_status.0.execution_state", expected.ExecutionState)
	return i
}

func (i *IcebergTableShowOutputAssert) HasPartitionSpecs(expected ...sdk.IcebergTablePartitionSpec) *IcebergTableShowOutputAssert {
	i.StringValueSet("partition_specs.#", strconv.Itoa(len(expected)))
	for specIdx, spec := range expected {
		i.StringValueSet(fmt.Sprintf("partition_specs.%d.spec_id", specIdx), strconv.Itoa(spec.SpecId))
		i.StringValueSet(fmt.Sprintf("partition_specs.%d.fields.#", specIdx), strconv.Itoa(len(spec.Fields)))
		for fieldIdx, field := range spec.Fields {
			i.StringValueSet(fmt.Sprintf("partition_specs.%d.fields.%d.name", specIdx, fieldIdx), field.Name)
			i.StringValueSet(fmt.Sprintf("partition_specs.%d.fields.%d.transform", specIdx, fieldIdx), field.Transform)
			i.StringValueSet(fmt.Sprintf("partition_specs.%d.fields.%d.source_id", specIdx, fieldIdx), strconv.Itoa(field.SourceId))
			i.StringValueSet(fmt.Sprintf("partition_specs.%d.fields.%d.field_id", specIdx, fieldIdx), strconv.Itoa(field.FieldId))
		}
	}
	return i
}
