package objectassert

import (
	"fmt"
	"slices"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (w *WarehouseInteractiveAssert) HasTables(expected ...sdk.SchemaObjectIdentifier) *WarehouseInteractiveAssert {
	w.AddAssertion(func(t *testing.T, o *sdk.WarehouseInteractive) error {
		t.Helper()
		// SHOW WAREHOUSES does not guarantee a stable order for the tables list, so compare order-insensitively.
		mapped := collections.Map(o.Tables, func(item sdk.SchemaObjectIdentifier) string { return item.FullyQualifiedName() })
		mappedExpected := collections.Map(expected, func(item sdk.SchemaObjectIdentifier) string { return item.FullyQualifiedName() })
		slices.Sort(mapped)
		slices.Sort(mappedExpected)
		if !slices.Equal(mapped, mappedExpected) {
			return fmt.Errorf("expected tables: %v; got: %v", mappedExpected, mapped)
		}
		return nil
	})
	return w
}
