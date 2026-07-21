package objectassert

import (
	"fmt"
	"strings"
	"testing"
	"time"

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

func (s *ServiceDetailsAssert) HasCreatedOnNotEmpty() *ServiceDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.ServiceDetails) error {
		t.Helper()
		if o.CreatedOn == (time.Time{}) {
			return fmt.Errorf("expected created_on to be not empty")
		}
		return nil
	})
	return s
}

func (s *ServiceDetailsAssert) HasUpdatedOnNotEmpty() *ServiceDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.ServiceDetails) error {
		t.Helper()
		if o.UpdatedOn == (time.Time{}) {
			return fmt.Errorf("expected updated_on to be not empty")
		}
		return nil
	})
	return s
}

func (s *ServiceDetailsAssert) HasSpecNotEmpty() *ServiceDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.ServiceDetails) error {
		t.Helper()
		if o.Spec == "" {
			return fmt.Errorf("expected spec to be not empty")
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

func (s *ServiceDetailsAssert) HasDnsNameNotEmpty() *ServiceDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.ServiceDetails) error {
		t.Helper()
		if o.DnsName == "" {
			return fmt.Errorf("expected dns name to be not empty")
		}
		return nil
	})
	return s
}

func (s *ServiceDetailsAssert) HasSpecDigestNotEmpty() *ServiceDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.ServiceDetails) error {
		t.Helper()
		if o.SpecDigest == "" {
			return fmt.Errorf("expected spec digest to be not empty")
		}
		return nil
	})
	return s
}
