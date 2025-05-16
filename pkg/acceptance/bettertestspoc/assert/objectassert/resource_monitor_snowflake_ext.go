package objectassert

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
)

func (r *ResourceMonitorAssert) HasNonEmptyStartTime() *ResourceMonitorAssert {
	r.AddAssertion(func(t *testing.T, o *sdk.ResourceMonitor) error {
		t.Helper()
		if o.StartTime == "" {
			return fmt.Errorf("expected start time to be non empty")
		}
		return nil
	})
	return r
}

func (r *ResourceMonitorAssert) HasNonEmptyEndTime() *ResourceMonitorAssert {
	r.AddAssertion(func(t *testing.T, o *sdk.ResourceMonitor) error {
		t.Helper()
		if o.StartTime == "" {
			return fmt.Errorf("expected end time to be non empty")
		}
		return nil
	})
	return r
}

func (r *ResourceMonitorAssert) HasSuspendAtNil() *ResourceMonitorAssert {
	r.AddAssertion(func(t *testing.T, o *sdk.ResourceMonitor) error {
		t.Helper()
		if o.SuspendAt != nil {
			return fmt.Errorf("expected suspend at to be nil, was %v", *o.SuspendAt)
		}
		return nil
	})
	return r
}

func (r *ResourceMonitorAssert) HasSuspendImmediateAtNil() *ResourceMonitorAssert {
	r.AddAssertion(func(t *testing.T, o *sdk.ResourceMonitor) error {
		t.Helper()
		if o.SuspendImmediateAt != nil {
			return fmt.Errorf("expected suspend immediate at to be nil, was %v", *o.SuspendImmediateAt)
		}
		return nil
	})
	return r
}
