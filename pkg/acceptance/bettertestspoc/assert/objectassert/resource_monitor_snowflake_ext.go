package objectassert

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

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
		if o.SuspendImmediatelyAt != nil {
			return fmt.Errorf("expected suspend immediate at to be nil, was %v", *o.SuspendImmediatelyAt)
		}
		return nil
	})
	return r
}
