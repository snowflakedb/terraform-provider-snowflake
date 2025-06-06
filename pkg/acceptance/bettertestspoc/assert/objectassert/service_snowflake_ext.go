package objectassert

import (
	"fmt"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (s *ServiceAssert) HasCreatedOnNotEmpty() *ServiceAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.Service) error {
		t.Helper()
		if o.CreatedOn == (time.Time{}) {
			return fmt.Errorf("expected created_on to be not empty")
		}
		return nil
	})
	return s
}

func (s *ServiceAssert) HasUpdatedOnNotEmpty() *ServiceAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.Service) error {
		t.Helper()
		if o.UpdatedOn == (time.Time{}) {
			return fmt.Errorf("expected updated_on to be not empty")
		}
		return nil
	})
	return s
}

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

func (s *ServiceAssert) HasNoResumedOn() *ServiceAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.Service) error {
		t.Helper()
		if o.ResumedOn != nil {
			return fmt.Errorf("expected resumed_on to be empty")
		}
		return nil
	})
	return s
}

func (s *ServiceAssert) HasNoSuspendedOn() *ServiceAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.Service) error {
		t.Helper()
		if o.SuspendedOn != nil {
			return fmt.Errorf("expected suspended_on to be empty")
		}
		return nil
	})
	return s
}

func (s *ServiceAssert) HasNoQueryWarehouse() *ServiceAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.Service) error {
		t.Helper()
		if o.QueryWarehouse != nil {
			return fmt.Errorf("expected query_warehouse to be empty")
		}
		return nil
	})
	return s
}

func (s *ServiceAssert) HasNoExternalAccessIntegrations() *ServiceAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.Service) error {
		t.Helper()
		if len(o.ExternalAccessIntegrations) > 0 {
			return fmt.Errorf("expected external access integrations to be empty")
		}
		return nil
	})
	return s
}

func (s *ServiceAssert) HasDnsNameNotEmpty() *ServiceAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.Service) error {
		t.Helper()
		if o.DnsName == "" {
			return fmt.Errorf("expected dns name to be not empty")
		}
		return nil
	})
	return s
}

func (s *ServiceAssert) HasNoComment() *ServiceAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.Service) error {
		t.Helper()
		if o.Comment != nil {
			return fmt.Errorf("expected comment to be empty")
		}
		return nil
	})
	return s
}

func (s *ServiceAssert) HasSpecDigestNotEmpty() *ServiceAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.Service) error {
		t.Helper()
		if o.SpecDigest == "" {
			return fmt.Errorf("expected spec digest to be not empty")
		}
		return nil
	})
	return s
}

func (s *ServiceAssert) HasNoManagingObjectDomain() *ServiceAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.Service) error {
		t.Helper()
		if o.ManagingObjectDomain != nil {
			return fmt.Errorf("expected managing object domain to be empty")
		}
		return nil
	})
	return s
}

func (s *ServiceAssert) HasNoManagingObjectName() *ServiceAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.Service) error {
		t.Helper()
		if o.ManagingObjectName != nil {
			return fmt.Errorf("expected managing object name to be empty")
		}
		return nil
	})
	return s
}
