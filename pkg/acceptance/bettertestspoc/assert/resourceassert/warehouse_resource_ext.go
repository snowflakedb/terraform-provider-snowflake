package resourceassert

import (
	r "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/resources"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/assert"
)

func (w *WarehouseResourceAssert) HasDefaultMaxConcurrencyLevel() *WarehouseResourceAssert {
	w.AddAssertion(assert.ValueSet("max_concurrency_level", "8"))
	return w
}

func (w *WarehouseResourceAssert) HasDefaultStatementQueuedTimeoutInSeconds() *WarehouseResourceAssert {
	w.AddAssertion(assert.ValueSet("statement_queued_timeout_in_seconds", "0"))
	return w
}

func (w *WarehouseResourceAssert) HasDefaultStatementTimeoutInSeconds() *WarehouseResourceAssert {
	w.AddAssertion(assert.ValueSet("statement_timeout_in_seconds", "172800"))
	return w
}

func (w *WarehouseResourceAssert) HasAllDefault() *WarehouseResourceAssert {
	return w.HasDefaultMaxConcurrencyLevel().
		HasNoWarehouseType().
		HasNoWarehouseSize().
		HasNoMaxClusterCount().
		HasNoMinClusterCount().
		HasNoScalingPolicy().
		HasAutoSuspendString(r.IntDefaultString).
		HasAutoResumeString(r.BooleanDefault).
		HasNoInitiallySuspended().
		HasNoResourceMonitor().
		HasEnableQueryAccelerationString(r.BooleanDefault).
		HasQueryAccelerationMaxScaleFactorString(r.IntDefaultString).
		HasDefaultMaxConcurrencyLevel().
		HasDefaultStatementQueuedTimeoutInSeconds().
		HasDefaultStatementTimeoutInSeconds()
}
