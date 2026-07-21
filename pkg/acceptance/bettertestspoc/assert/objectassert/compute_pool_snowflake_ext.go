package objectassert

import (
	"fmt"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (c *ComputePoolAssert) HasCreatedOnNotEmpty() *ComputePoolAssert {
	c.AddAssertion(func(t *testing.T, o *sdk.ComputePool) error {
		t.Helper()
		if o.CreatedOn == (time.Time{}) {
			return fmt.Errorf("expected created_on to be not empty")
		}
		return nil
	})
	return c
}

func (c *ComputePoolAssert) HasResumedOnNotEmpty() *ComputePoolAssert {
	c.AddAssertion(func(t *testing.T, o *sdk.ComputePool) error {
		t.Helper()
		if o.ResumedOn == (time.Time{}) {
			return fmt.Errorf("expected resumed_on to be not empty")
		}
		return nil
	})
	return c
}

func (c *ComputePoolAssert) HasUpdatedOnNotEmpty() *ComputePoolAssert {
	c.AddAssertion(func(t *testing.T, o *sdk.ComputePool) error {
		t.Helper()
		if o.UpdatedOn == (time.Time{}) {
			return fmt.Errorf("expected updated_on to be not empty")
		}
		return nil
	})
	return c
}
