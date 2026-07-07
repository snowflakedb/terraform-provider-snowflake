package resourceshowoutputassert

import (
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
