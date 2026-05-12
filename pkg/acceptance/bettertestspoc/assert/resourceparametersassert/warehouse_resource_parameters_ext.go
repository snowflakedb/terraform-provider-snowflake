package resourceparametersassert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func WarehousesDatasourceParameters(t *testing.T, datasourceReference string) *WarehouseResourceParametersAssert {
	t.Helper()

	w := WarehouseResourceParametersAssert{
		ResourceAssert: assert.NewDatasourceAssert(datasourceReference, "parameters", "warehouses.0."),
	}
	return &w
}

func (w *WarehouseResourceParametersAssert) HasDefaultMaxConcurrencyLevel() *WarehouseResourceParametersAssert {
	return w.
		HasMaxConcurrencyLevel(8).
		HasMaxConcurrencyLevelLevel("")
}

func (w *WarehouseResourceParametersAssert) HasDefaultStatementQueuedTimeoutInSeconds() *WarehouseResourceParametersAssert {
	return w.
		HasStatementQueuedTimeoutInSeconds(0).
		HasStatementQueuedTimeoutInSecondsLevel("")
}

func (w *WarehouseResourceParametersAssert) HasDefaultStatementTimeoutInSeconds() *WarehouseResourceParametersAssert {
	return w.
		HasStatementTimeoutInSeconds(172800).
		HasStatementTimeoutInSecondsLevel("")
}
