package objectassert

import (
	"fmt"
	"slices"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (w *WarehouseAssert) HasStateOneOf(expected ...sdk.WarehouseState) *WarehouseAssert {
	w.AddAssertion(func(t *testing.T, o *sdk.Warehouse) error {
		t.Helper()
		if !slices.Contains(expected, o.State) {
			return fmt.Errorf("expected state one of: %v; got: %v", expected, string(o.State))
		}
		return nil
	})
	return w
}

func (w *WarehouseAssert) HasNoGeneration() *WarehouseAssert {
	w.AddAssertion(func(t *testing.T, o *sdk.Warehouse) error {
		t.Helper()
		if o.Generation != nil {
			return fmt.Errorf("expected generation to be empty; got: %s", *o.Generation)
		}
		return nil
	})
	return w
}

func (w *WarehouseAssert) HasNoResourceConstraint() *WarehouseAssert {
	w.AddAssertion(func(t *testing.T, o *sdk.Warehouse) error {
		t.Helper()
		if o.ResourceConstraint != nil {
			return fmt.Errorf("expected resource constraint to be empty; got: %s", *o.ResourceConstraint)
		}
		return nil
	})
	return w
}

func (w *WarehouseAssert) HasNoSize() *WarehouseAssert {
	w.AddAssertion(func(t *testing.T, o *sdk.Warehouse) error {
		t.Helper()
		if o.Size != nil {
			return fmt.Errorf("expected size to be empty; got: %s", *o.Size)
		}
		return nil
	})
	return w
}

func (w *WarehouseAssert) HasNoMaxClusterCount() *WarehouseAssert {
	w.AddAssertion(func(t *testing.T, o *sdk.Warehouse) error {
		t.Helper()
		if o.MaxClusterCount != nil {
			return fmt.Errorf("expected max cluster count to be empty; got: %d", *o.MaxClusterCount)
		}
		return nil
	})
	return w
}

func (w *WarehouseAssert) HasNoMinClusterCount() *WarehouseAssert {
	w.AddAssertion(func(t *testing.T, o *sdk.Warehouse) error {
		t.Helper()
		if o.MinClusterCount != nil {
			return fmt.Errorf("expected min cluster count to be empty; got: %d", *o.MinClusterCount)
		}
		return nil
	})
	return w
}

func (w *WarehouseAssert) HasNoScalingPolicy() *WarehouseAssert {
	w.AddAssertion(func(t *testing.T, o *sdk.Warehouse) error {
		t.Helper()
		if o.ScalingPolicy != nil {
			return fmt.Errorf("expected scaling policy to be empty; got: %s", *o.ScalingPolicy)
		}
		return nil
	})
	return w
}

func (w *WarehouseAssert) HasNoAutoSuspend() *WarehouseAssert {
	w.AddAssertion(func(t *testing.T, o *sdk.Warehouse) error {
		t.Helper()
		if o.AutoSuspend != nil {
			return fmt.Errorf("expected auto suspend to be empty; got: %d", *o.AutoSuspend)
		}
		return nil
	})
	return w
}

func (w *WarehouseAssert) HasNoEnableQueryAcceleration() *WarehouseAssert {
	w.AddAssertion(func(t *testing.T, o *sdk.Warehouse) error {
		t.Helper()
		if o.EnableQueryAcceleration != nil {
			return fmt.Errorf("expected enable query acceleration to be empty; got: %t", *o.EnableQueryAcceleration)
		}
		return nil
	})
	return w
}

func (w *WarehouseAssert) HasNoQueryAccelerationMaxScaleFactor() *WarehouseAssert {
	w.AddAssertion(func(t *testing.T, o *sdk.Warehouse) error {
		t.Helper()
		if o.QueryAccelerationMaxScaleFactor != nil {
			return fmt.Errorf("expected query acceleration max scale factor to be empty; got: %d", *o.QueryAccelerationMaxScaleFactor)
		}
		return nil
	})
	return w
}

func (w *WarehouseAssert) HasNoMaxQueryPerformanceLevel() *WarehouseAssert {
	w.AddAssertion(func(t *testing.T, o *sdk.Warehouse) error {
		t.Helper()
		if o.MaxQueryPerformanceLevel != nil {
			return fmt.Errorf("expected max query performance level to be empty; got: %s", *o.MaxQueryPerformanceLevel)
		}
		return nil
	})
	return w
}

func (w *WarehouseAssert) HasNoQueryThroughputMultiplier() *WarehouseAssert {
	w.AddAssertion(func(t *testing.T, o *sdk.Warehouse) error {
		t.Helper()
		if o.QueryThroughputMultiplier != nil {
			return fmt.Errorf("expected query throughput multiplier to be empty; got: %d", *o.QueryThroughputMultiplier)
		}
		return nil
	})
	return w
}
