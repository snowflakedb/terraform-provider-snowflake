package objectassert

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

// HasQueryThroughputMultiplierSet asserts the multiplier is present.
// Snowflake does not guarantee a specific default after a warehouse type change (SNOW-3687301).
func (w *WarehouseAdaptiveAssert) HasQueryThroughputMultiplierSet() *WarehouseAdaptiveAssert {
	w.AddAssertion(func(t *testing.T, o *sdk.WarehouseAdaptive) error {
		t.Helper()
		if o.QueryThroughputMultiplier == nil {
			return fmt.Errorf("expected query throughput multiplier to have value; got: nil")
		}
		return nil
	})
	return w
}
