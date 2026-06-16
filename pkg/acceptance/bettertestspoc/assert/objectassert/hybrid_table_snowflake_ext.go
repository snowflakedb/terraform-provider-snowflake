package objectassert

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

// HasRowsNil asserts that SHOW HYBRID TABLES returned NULL for the rows column.
// SHOW reports NULL for empty hybrid tables until rows have been written, so the
// generator-emitted HasRows(0) helper (which rejects nil) cannot be used in
// create/import-time assertions.
func (h *HybridTableAssert) HasRowsNil() *HybridTableAssert {
	h.AddAssertion(func(t *testing.T, o *sdk.HybridTable) error {
		t.Helper()
		if o.Rows != nil {
			return fmt.Errorf("expected rows: nil; got: %d", *o.Rows)
		}
		return nil
	})
	return h
}

// HasBytesNil asserts that SHOW HYBRID TABLES returned NULL for the bytes column.
// SHOW reports NULL for empty hybrid tables until rows have been written, so the
// generator-emitted HasBytes(0) helper (which rejects nil) cannot be used in
// create/import-time assertions.
func (h *HybridTableAssert) HasBytesNil() *HybridTableAssert {
	h.AddAssertion(func(t *testing.T, o *sdk.HybridTable) error {
		t.Helper()
		if o.Bytes != nil {
			return fmt.Errorf("expected bytes: nil; got: %d", *o.Bytes)
		}
		return nil
	})
	return h
}
