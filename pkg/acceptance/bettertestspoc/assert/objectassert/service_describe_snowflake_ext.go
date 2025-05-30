package objectassert

import (
	"fmt"
	"slices"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type ServiceDetailsAssert struct {
	*assert.SnowflakeObjectAssert[sdk.ServiceDetails, sdk.SchemaObjectIdentifier]
}

func ServiceDetails(t *testing.T, id sdk.SchemaObjectIdentifier) *ServiceDetailsAssert {
	t.Helper()
	return &ServiceDetailsAssert{
		assert.NewSnowflakeObjectAssertWithTestClientObjectProvider(sdk.ObjectTypeService, id, func(testClient *helpers.TestClient) assert.ObjectProvider[sdk.ServiceDetails, sdk.SchemaObjectIdentifier] {
			return testClient.Service.Describe
		}),
	}
}

func (s *ServiceDetailsAssert) HasName(expected string) *ServiceDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.ServiceDetails) error {
		t.Helper()
		if o.Name != expected {
			return fmt.Errorf("expected name: %v; got: %v", expected, o.Name)
		}
		return nil
	})
	return s
}

func (s *ServiceDetailsAssert) HasStatus(expected sdk.ServiceStatus) *ServiceDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.ServiceDetails) error {
		t.Helper()
		if o.Status != expected {
			return fmt.Errorf("expected status: %v; got: %v", expected, o.Status)
		}
		return nil
	})
	return s
}

func (s *ServiceDetailsAssert) HasDatabaseName(expected string) *ServiceDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.ServiceDetails) error {
		t.Helper()
		if o.DatabaseName != expected {
			return fmt.Errorf("expected database name: %v; got: %v", expected, o.DatabaseName)
		}
		return nil
	})
	return s
}

func (s *ServiceDetailsAssert) HasSchemaName(expected string) *ServiceDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.ServiceDetails) error {
		t.Helper()
		if o.SchemaName != expected {
			return fmt.Errorf("expected schema name: %v; got: %v", expected, o.SchemaName)
		}
		return nil
	})
	return s
}

func (s *ServiceDetailsAssert) HasOwner(expected string) *ServiceDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.ServiceDetails) error {
		t.Helper()
		if o.Owner != expected {
			return fmt.Errorf("expected owner: %v; got: %v", expected, o.Owner)
		}
		return nil
	})
	return s
}

func (s *ServiceDetailsAssert) HasComputePool(expected sdk.AccountObjectIdentifier) *ServiceDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.ServiceDetails) error {
		t.Helper()
		if o.ComputePool.Name() != expected.Name() {
			return fmt.Errorf("expected compute pool: %v; got: %v", expected.Name(), o.ComputePool.Name())
		}
		return nil
	})
	return s
}

func (s *ServiceDetailsAssert) HasSpec(expected string) *ServiceDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.ServiceDetails) error {
		t.Helper()
		if o.Spec != expected {
			return fmt.Errorf("expected spec: %v; got: %v", expected, o.Spec)
		}
		return nil
	})
	return s
}

func (s *ServiceDetailsAssert) HasDnsName(expected string) *ServiceDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.ServiceDetails) error {
		t.Helper()
		if o.DnsName != expected {
			return fmt.Errorf("expected dns name: %v; got: %v", expected, o.DnsName)
		}
		return nil
	})
	return s
}

func (s *ServiceDetailsAssert) HasCurrentInstances(expected int) *ServiceDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.ServiceDetails) error {
		t.Helper()
		if o.CurrentInstances != expected {
			return fmt.Errorf("expected current instances: %v; got: %v", expected, o.CurrentInstances)
		}
		return nil
	})
	return s
}

func (s *ServiceDetailsAssert) HasTargetInstances(expected int) *ServiceDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.ServiceDetails) error {
		t.Helper()
		if o.TargetInstances != expected {
			return fmt.Errorf("expected target instances: %v; got: %v", expected, o.TargetInstances)
		}
		return nil
	})
	return s
}

func (s *ServiceDetailsAssert) HasMinReadyInstances(expected int) *ServiceDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.ServiceDetails) error {
		t.Helper()
		if o.MinReadyInstances != expected {
			return fmt.Errorf("expected min ready instances: %v; got: %v", expected, o.MinReadyInstances)
		}
		return nil
	})
	return s
}

func (s *ServiceDetailsAssert) HasMinInstances(expected int) *ServiceDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.ServiceDetails) error {
		t.Helper()
		if o.MinInstances != expected {
			return fmt.Errorf("expected min instances: %v; got: %v", expected, o.MinInstances)
		}
		return nil
	})
	return s
}

func (s *ServiceDetailsAssert) HasMaxInstances(expected int) *ServiceDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.ServiceDetails) error {
		t.Helper()
		if o.MaxInstances != expected {
			return fmt.Errorf("expected max instances: %v; got: %v", expected, o.MaxInstances)
		}
		return nil
	})
	return s
}

func (s *ServiceDetailsAssert) HasAutoResume(expected bool) *ServiceDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.ServiceDetails) error {
		t.Helper()
		if o.AutoResume != expected {
			return fmt.Errorf("expected auto resume: %v; got: %v", expected, o.AutoResume)
		}
		return nil
	})
	return s
}

func (s *ServiceDetailsAssert) HasExternalAccessIntegrations(expected ...sdk.AccountObjectIdentifier) *ServiceDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.ServiceDetails) error {
		t.Helper()
		for _, expected := range expected {
			if !slices.Contains(o.ExternalAccessIntegrations, expected) {
				return fmt.Errorf("expected external access integrations: %v; got: %v", expected, o.ExternalAccessIntegrations)
			}
		}
		return nil
	})
	return s
}

func (s *ServiceDetailsAssert) HasCreatedOn(expected time.Time) *ServiceDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.ServiceDetails) error {
		t.Helper()
		if o.CreatedOn != expected {
			return fmt.Errorf("expected created on: %v; got: %v", expected, o.CreatedOn)
		}
		return nil
	})
	return s
}

func (s *ServiceDetailsAssert) HasUpdatedOn(expected time.Time) *ServiceDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.ServiceDetails) error {
		t.Helper()
		if o.UpdatedOn != expected {
			return fmt.Errorf("expected updated on: %v; got: %v", expected, o.UpdatedOn)
		}
		return nil
	})
	return s
}

func (s *ServiceDetailsAssert) HasResumedOn(expected time.Time) *ServiceDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.ServiceDetails) error {
		t.Helper()
		if o.ResumedOn == nil {
			return fmt.Errorf("expected resumed on to have value; got: nil")
		}
		if *o.ResumedOn != expected {
			return fmt.Errorf("expected resumed on: %v; got: %v", expected, *o.ResumedOn)
		}
		return nil
	})
	return s
}

func (s *ServiceDetailsAssert) HasSuspendedOn(expected time.Time) *ServiceDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.ServiceDetails) error {
		t.Helper()
		if o.SuspendedOn == nil {
			return fmt.Errorf("expected suspended on to have value; got: nil")
		}
		if *o.SuspendedOn != expected {
			return fmt.Errorf("expected suspended on: %v; got: %v", expected, *o.SuspendedOn)
		}
		return nil
	})
	return s
}

func (s *ServiceDetailsAssert) HasAutoSuspendSecs(expected int) *ServiceDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.ServiceDetails) error {
		t.Helper()
		if o.AutoSuspendSecs != expected {
			return fmt.Errorf("expected auto suspend secs: %v; got: %v", expected, o.AutoSuspendSecs)
		}
		return nil
	})
	return s
}

func (s *ServiceDetailsAssert) HasComment(expected string) *ServiceDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.ServiceDetails) error {
		t.Helper()
		if o.Comment == nil {
			return fmt.Errorf("expected comment to have value; got: nil")
		}
		if *o.Comment != expected {
			return fmt.Errorf("expected comment: %v; got: %v", expected, *o.Comment)
		}
		return nil
	})
	return s
}

func (s *ServiceDetailsAssert) HasOwnerRoleType(expected string) *ServiceDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.ServiceDetails) error {
		t.Helper()
		if o.OwnerRoleType != expected {
			return fmt.Errorf("expected owner role type: %v; got: %v", expected, o.OwnerRoleType)
		}
		return nil
	})
	return s
}

func (s *ServiceDetailsAssert) HasQueryWarehouse(expected sdk.AccountObjectIdentifier) *ServiceDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.ServiceDetails) error {
		t.Helper()
		if o.QueryWarehouse == nil {
			return fmt.Errorf("expected query warehouse to have value; got: nil")
		}
		if o.QueryWarehouse.Name() != expected.Name() {
			return fmt.Errorf("expected query warehouse: %v; got: %v", expected.Name(), o.QueryWarehouse.Name())
		}
		return nil
	})
	return s
}

func (s *ServiceDetailsAssert) HasIsJob(expected bool) *ServiceDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.ServiceDetails) error {
		t.Helper()
		if o.IsJob != expected {
			return fmt.Errorf("expected is job: %v; got: %v", expected, o.IsJob)
		}
		return nil
	})
	return s
}

func (s *ServiceDetailsAssert) HasIsAsyncJob(expected bool) *ServiceDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.ServiceDetails) error {
		t.Helper()
		if o.IsAsyncJob != expected {
			return fmt.Errorf("expected is async job: %v; got: %v", expected, o.IsAsyncJob)
		}
		return nil
	})
	return s
}

func (s *ServiceDetailsAssert) HasSpecDigest(expected string) *ServiceDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.ServiceDetails) error {
		t.Helper()
		if o.SpecDigest != expected {
			return fmt.Errorf("expected spec digest: %v; got: %v", expected, o.SpecDigest)
		}
		return nil
	})
	return s
}

func (s *ServiceDetailsAssert) HasIsUpgrading(expected bool) *ServiceDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.ServiceDetails) error {
		t.Helper()
		if o.IsUpgrading != expected {
			return fmt.Errorf("expected is upgrading: %v; got: %v", expected, o.IsUpgrading)
		}
		return nil
	})
	return s
}

func (s *ServiceDetailsAssert) HasManagingObjectDomain(expected string) *ServiceDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.ServiceDetails) error {
		t.Helper()
		if o.ManagingObjectDomain == nil {
			return fmt.Errorf("expected managing object domain to have value; got: nil")
		}
		if *o.ManagingObjectDomain != expected {
			return fmt.Errorf("expected managing object domain: %v; got: %v", expected, *o.ManagingObjectDomain)
		}
		return nil
	})
	return s
}

func (s *ServiceDetailsAssert) HasManagingObjectName(expected string) *ServiceDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.ServiceDetails) error {
		t.Helper()
		if o.ManagingObjectName == nil {
			return fmt.Errorf("expected managing object name to have value; got: nil")
		}
		if *o.ManagingObjectName != expected {
			return fmt.Errorf("expected managing object name: %v; got: %v", expected, *o.ManagingObjectName)
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

func (s *ServiceDetailsAssert) HasNoResumedOn() *ServiceDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.ServiceDetails) error {
		t.Helper()
		if o.ResumedOn != nil {
			return fmt.Errorf("expected resumed_on to be empty")
		}
		return nil
	})
	return s
}

func (s *ServiceDetailsAssert) HasNoSuspendedOn() *ServiceDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.ServiceDetails) error {
		t.Helper()
		if o.SuspendedOn != nil {
			return fmt.Errorf("expected suspended_on to be empty")
		}
		return nil
	})
	return s
}

func (s *ServiceDetailsAssert) HasNoQueryWarehouse() *ServiceDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.ServiceDetails) error {
		t.Helper()
		if o.QueryWarehouse != nil {
			return fmt.Errorf("expected query_warehouse to be empty")
		}
		return nil
	})
	return s
}

func (s *ServiceDetailsAssert) HasNoExternalAccessIntegrations() *ServiceDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.ServiceDetails) error {
		t.Helper()
		if len(o.ExternalAccessIntegrations) > 0 {
			return fmt.Errorf("expected external access integrations to be empty")
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

func (s *ServiceDetailsAssert) HasNoComment() *ServiceDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.ServiceDetails) error {
		t.Helper()
		if o.Comment != nil {
			return fmt.Errorf("expected comment to be empty")
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

func (s *ServiceDetailsAssert) HasNoManagingObjectDomain() *ServiceDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.ServiceDetails) error {
		t.Helper()
		if o.ManagingObjectDomain != nil {
			return fmt.Errorf("expected managing object domain to be empty")
		}
		return nil
	})
	return s
}

func (s *ServiceDetailsAssert) HasNoManagingObjectName() *ServiceDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.ServiceDetails) error {
		t.Helper()
		if o.ManagingObjectName != nil {
			return fmt.Errorf("expected managing object name to be empty")
		}
		return nil
	})
	return s
}
