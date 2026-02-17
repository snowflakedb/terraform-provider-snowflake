package objectassert

import (
	"fmt"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (a *HybridTableAssert) HasCreatedOnNotEmpty() *HybridTableAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.HybridTable) error {
		t.Helper()
		if o.CreatedOn == (time.Time{}) {
			return fmt.Errorf("expected created_on to be not empty")
		}
		return nil
	})
	return a
}

func (a *HybridTableAssert) HasNonZeroRows() *HybridTableAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.HybridTable) error {
		t.Helper()
		if o.Rows == 0 {
			return fmt.Errorf("expected rows to be greater than 0, got: %d", o.Rows)
		}
		return nil
	})
	return a
}

func (a *HybridTableAssert) HasCommentNotEmpty() *HybridTableAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.HybridTable) error {
		t.Helper()
		if o.Comment == "" {
			return fmt.Errorf("expected comment to be not empty")
		}
		return nil
	})
	return a
}

func (a *HybridTableAssert) HasOwnerNotEmpty() *HybridTableAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.HybridTable) error {
		t.Helper()
		if o.Owner == "" {
			return fmt.Errorf("expected owner to be not empty")
		}
		return nil
	})
	return a
}

func (a *HybridTableAssert) HasCommentEmpty() *HybridTableAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.HybridTable) error {
		t.Helper()
		if o.Comment != "" {
			return fmt.Errorf("expected comment to be empty, got: %s", o.Comment)
		}
		return nil
	})
	return a
}

func (a *HybridTableAssert) HasRowsGreaterThanOrEqual(expected int) *HybridTableAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.HybridTable) error {
		t.Helper()
		if o.Rows < expected {
			return fmt.Errorf("expected rows to be >= %d, got: %d", expected, o.Rows)
		}
		return nil
	})
	return a
}

func (a *HybridTableAssert) HasBytesGreaterThanOrEqual(expected int) *HybridTableAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.HybridTable) error {
		t.Helper()
		if o.Bytes < expected {
			return fmt.Errorf("expected bytes to be >= %d, got: %d", expected, o.Bytes)
		}
		return nil
	})
	return a
}
