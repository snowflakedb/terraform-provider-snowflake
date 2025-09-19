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
