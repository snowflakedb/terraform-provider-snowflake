package objectassert

import (
	"fmt"
	"slices"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type StorageIntegrationAwsDetailsAssert struct {
	*assert.SnowflakeObjectAssert[sdk.StorageIntegrationAwsDetails, sdk.AccountObjectIdentifier]
}

func StorageIntegrationAwsDetails(t *testing.T, id sdk.AccountObjectIdentifier) *StorageIntegrationAwsDetailsAssert {
	t.Helper()
	return &StorageIntegrationAwsDetailsAssert{
		assert.NewSnowflakeObjectAssertWithTestClientObjectProvider(sdk.ObjectType("StorageIntegrationAwsDetails"), id, func(testClient *helpers.TestClient) assert.ObjectProvider[sdk.StorageIntegrationAwsDetails, sdk.AccountObjectIdentifier] {
			return testClient.StorageIntegration.DescribeAws
		}),
	}
}

func StorageIntegrationAwsDetailsFromObject(t *testing.T, storageIntegrationAwsDetails *sdk.StorageIntegrationAwsDetails) *StorageIntegrationAwsDetailsAssert {
	t.Helper()
	return &StorageIntegrationAwsDetailsAssert{
		assert.NewSnowflakeObjectAssertWithObject(sdk.ObjectType("StorageIntegrationAwsDetails"), storageIntegrationAwsDetails.ID(), storageIntegrationAwsDetails),
	}
}

func (s *StorageIntegrationAwsDetailsAssert) HasEnabled(expected bool) *StorageIntegrationAwsDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.StorageIntegrationAwsDetails) error {
		t.Helper()
		if o.Enabled != expected {
			return fmt.Errorf("expected enabled: %v; got: %v", expected, o.Enabled)
		}
		return nil
	})
	return s
}

func (s *StorageIntegrationAwsDetailsAssert) HasProvider(expected string) *StorageIntegrationAwsDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.StorageIntegrationAwsDetails) error {
		t.Helper()
		if o.Provider != expected {
			return fmt.Errorf("expected provider: %v; got: %v", expected, o.Provider)
		}
		return nil
	})
	return s
}

func (s *StorageIntegrationAwsDetailsAssert) HasAllowedLocations(expected ...string) *StorageIntegrationAwsDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.StorageIntegrationAwsDetails) error {
		t.Helper()
		mapped := collections.Map(o.AllowedLocations, func(item string) any { return item })
		mappedExpected := collections.Map(expected, func(item string) any { return item })
		if !slices.Equal(mapped, mappedExpected) {
			return fmt.Errorf("expected allowed locations: %v; got: %v", expected, o.AllowedLocations)
		}
		return nil
	})
	return s
}

func (s *StorageIntegrationAwsDetailsAssert) HasBlockedLocations(expected ...string) *StorageIntegrationAwsDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.StorageIntegrationAwsDetails) error {
		t.Helper()
		mapped := collections.Map(o.BlockedLocations, func(item string) any { return item })
		mappedExpected := collections.Map(expected, func(item string) any { return item })
		if !slices.Equal(mapped, mappedExpected) {
			return fmt.Errorf("expected blocked locations: %v; got: %v", expected, o.BlockedLocations)
		}
		return nil
	})
	return s
}

func (s *StorageIntegrationAwsDetailsAssert) HasComment(expected string) *StorageIntegrationAwsDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.StorageIntegrationAwsDetails) error {
		t.Helper()
		if o.Comment != expected {
			return fmt.Errorf("expected comment: %v; got: %v", expected, o.Comment)
		}
		return nil
	})
	return s
}

func (s *StorageIntegrationAwsDetailsAssert) HasUsePrivatelinkEndpoint(expected bool) *StorageIntegrationAwsDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.StorageIntegrationAwsDetails) error {
		t.Helper()
		if o.UsePrivatelinkEndpoint != expected {
			return fmt.Errorf("expected use privatelink endpoint: %v; got: %v", expected, o.UsePrivatelinkEndpoint)
		}
		return nil
	})
	return s
}

func (s *StorageIntegrationAwsDetailsAssert) HasIamUserArn(expected string) *StorageIntegrationAwsDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.StorageIntegrationAwsDetails) error {
		t.Helper()
		if o.IamUserArn != expected {
			return fmt.Errorf("expected iam user arn: %v; got: %v", expected, o.IamUserArn)
		}
		return nil
	})
	return s
}

func (s *StorageIntegrationAwsDetailsAssert) HasRoleArn(expected string) *StorageIntegrationAwsDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.StorageIntegrationAwsDetails) error {
		t.Helper()
		if o.RoleArn != expected {
			return fmt.Errorf("expected role arn: %v; got: %v", expected, o.RoleArn)
		}
		return nil
	})
	return s
}

func (s *StorageIntegrationAwsDetailsAssert) HasObjectAcl(expected string) *StorageIntegrationAwsDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.StorageIntegrationAwsDetails) error {
		t.Helper()
		if o.ObjectAcl != expected {
			return fmt.Errorf("expected object acl: %v; got: %v", expected, o.ObjectAcl)
		}
		return nil
	})
	return s
}

func (s *StorageIntegrationAwsDetailsAssert) HasExternalId(expected string) *StorageIntegrationAwsDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.StorageIntegrationAwsDetails) error {
		t.Helper()
		if o.ExternalId != expected {
			return fmt.Errorf("expected external id: %v; got: %v", expected, o.ExternalId)
		}
		return nil
	})
	return s
}
