package objectassert

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (s *ServiceDetailsAssert) HasCurrentInstancesBetween(min, max int) *ServiceDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.ServiceDetails) error {
		t.Helper()
		if o.CurrentInstances < min || max < o.CurrentInstances {
			return fmt.Errorf("expected current instances to be between %d and %d; got: %d", min, max, o.CurrentInstances)
		}
		return nil
	})
	return s
}

func (s *ServiceDetailsAssert) HasTargetInstancesBetween(min, max int) *ServiceDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.ServiceDetails) error {
		t.Helper()
		if o.TargetInstances < min || max < o.TargetInstances {
			return fmt.Errorf("expected target instances to be between %d and %d; got: %d", min, max, o.TargetInstances)
		}
		return nil
	})
	return s
}

func (s *ServiceDetailsAssert) HasSpecThatContains(content string) *ServiceDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.ServiceDetails) error {
		t.Helper()
		if !strings.Contains(o.Spec, content) {
			return fmt.Errorf("expected spec to contain: %v; got: %v", content, o.Spec)
		}
		return nil
	})
	return s
}
