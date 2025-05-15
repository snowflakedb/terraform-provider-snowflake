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

func (c *ComputePoolAssert) HasNoComment() *ComputePoolAssert {
	c.AddAssertion(func(t *testing.T, o *sdk.ComputePool) error {
		t.Helper()
		if o.Comment != nil {
			return fmt.Errorf("expected comment to have nil; got: %s", *o.Comment)
		}
		return nil
	})
	return c
}

func (c *ComputePoolAssert) HasNoApplication() *ComputePoolAssert {
	c.AddAssertion(func(t *testing.T, o *sdk.ComputePool) error {
		t.Helper()
		if o.Application != nil {
			return fmt.Errorf("expected application to have nil; got: %s", *o.Application)
		}
		return nil
	})
	return c
}

func (c *ComputePoolAssert) HasNoBudget() *ComputePoolAssert {
	c.AddAssertion(func(t *testing.T, o *sdk.ComputePool) error {
		t.Helper()
		if o.Budget != nil {
			return fmt.Errorf("expected budget to have nil; got: %s", *o.Budget)
		}
		return nil
	})
	return c
}
