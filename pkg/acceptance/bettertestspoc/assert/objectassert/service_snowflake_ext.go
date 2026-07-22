package objectassert

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (s *ServiceAssert) HasResumedOnNotEmpty() *ServiceAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.Service) error {
		t.Helper()
		if o.ResumedOn == nil {
			return fmt.Errorf("expected resumed_on to be not empty")
		}
		return nil
	})
	return s
}

func (s *ServiceAssert) HasSuspendedOnNotEmpty() *ServiceAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.Service) error {
		t.Helper()
		if o.SuspendedOn == nil {
			return fmt.Errorf("expected suspended_on to be not empty")
		}
		return nil
	})
	return s
}

func (s *ServiceAssert) HasCurrentInstancesBetween(min, max int) *ServiceAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.Service) error {
		t.Helper()
		if o.CurrentInstances < min || max < o.CurrentInstances {
			return fmt.Errorf("expected current instances to be between %d and %d; got: %d", min, max, o.CurrentInstances)
		}
		return nil
	})
	return s
}

func (s *ServiceAssert) HasTargetInstancesBetween(min, max int) *ServiceAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.Service) error {
		t.Helper()
		if o.TargetInstances < min || max < o.TargetInstances {
			return fmt.Errorf("expected target instances to be between %d and %d; got: %d", min, max, o.TargetInstances)
		}
		return nil
	})
	return s
}
