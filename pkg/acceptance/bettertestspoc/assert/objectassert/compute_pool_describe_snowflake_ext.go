package objectassert

import (
	"fmt"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (a *ComputePoolDetailsAssert) HasCreatedOnNotEmpty() *ComputePoolDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ComputePoolDetails) error {
		t.Helper()
		if o.CreatedOn == (time.Time{}) {
			return fmt.Errorf("expected created_on to be not empty")
		}
		return nil
	})
	return a
}

func (a *ComputePoolDetailsAssert) HasResumedOnNotEmpty() *ComputePoolDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ComputePoolDetails) error {
		t.Helper()
		if o.ResumedOn == (time.Time{}) {
			return fmt.Errorf("expected resumed_on to be not empty")
		}
		return nil
	})
	return a
}

func (a *ComputePoolDetailsAssert) HasUpdatedOnNotEmpty() *ComputePoolDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ComputePoolDetails) error {
		t.Helper()
		if o.UpdatedOn == (time.Time{}) {
			return fmt.Errorf("expected updated_on to be not empty")
		}
		return nil
	})
	return a
}
