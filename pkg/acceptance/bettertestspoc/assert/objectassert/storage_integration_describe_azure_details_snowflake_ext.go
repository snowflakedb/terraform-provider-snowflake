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

type StorageIntegrationAzureDetailsAssert struct {
	*assert.SnowflakeObjectAssert[sdk.StorageIntegrationAzureDetails, sdk.AccountObjectIdentifier]
}

func StorageIntegrationAzureDetails(t *testing.T, id sdk.AccountObjectIdentifier) *StorageIntegrationAzureDetailsAssert {
	t.Helper()
	return &StorageIntegrationAzureDetailsAssert{
		assert.NewSnowflakeObjectAssertWithTestClientObjectProvider(sdk.ObjectType("StorageIntegrationAzureDetails"), id, func(testClient *helpers.TestClient) assert.ObjectProvider[sdk.StorageIntegrationAzureDetails, sdk.AccountObjectIdentifier] {
			return testClient.StorageIntegration.DescribeAzure
		}),
	}
}

func StorageIntegrationAzureDetailsFromObject(t *testing.T, storageIntegrationAzureDetails *sdk.StorageIntegrationAzureDetails) *StorageIntegrationAzureDetailsAssert {
	t.Helper()
	return &StorageIntegrationAzureDetailsAssert{
		assert.NewSnowflakeObjectAssertWithObject(sdk.ObjectType("StorageIntegrationAzureDetails"), storageIntegrationAzureDetails.ID(), storageIntegrationAzureDetails),
	}
}

func (s *StorageIntegrationAzureDetailsAssert) HasEnabled(expected bool) *StorageIntegrationAzureDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.StorageIntegrationAzureDetails) error {
		t.Helper()
		if o.Enabled != expected {
			return fmt.Errorf("expected enabled: %v; got: %v", expected, o.Enabled)
		}
		return nil
	})
	return s
}

func (s *StorageIntegrationAzureDetailsAssert) HasProvider(expected string) *StorageIntegrationAzureDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.StorageIntegrationAzureDetails) error {
		t.Helper()
		if o.Provider != expected {
			return fmt.Errorf("expected provider: %v; got: %v", expected, o.Provider)
		}
		return nil
	})
	return s
}

func (s *StorageIntegrationAzureDetailsAssert) HasAllowedLocations(expected ...string) *StorageIntegrationAzureDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.StorageIntegrationAzureDetails) error {
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

func (s *StorageIntegrationAzureDetailsAssert) HasBlockedLocations(expected ...string) *StorageIntegrationAzureDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.StorageIntegrationAzureDetails) error {
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

func (s *StorageIntegrationAzureDetailsAssert) HasComment(expected string) *StorageIntegrationAzureDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.StorageIntegrationAzureDetails) error {
		t.Helper()
		if o.Comment != expected {
			return fmt.Errorf("expected comment: %v; got: %v", expected, o.Comment)
		}
		return nil
	})
	return s
}

func (s *StorageIntegrationAzureDetailsAssert) HasUsePrivatelinkEndpoint(expected bool) *StorageIntegrationAzureDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.StorageIntegrationAzureDetails) error {
		t.Helper()
		if o.UsePrivatelinkEndpoint != expected {
			return fmt.Errorf("expected use privatelink endpoint: %v; got: %v", expected, o.UsePrivatelinkEndpoint)
		}
		return nil
	})
	return s
}

func (s *StorageIntegrationAzureDetailsAssert) HasTenantId(expected string) *StorageIntegrationAzureDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.StorageIntegrationAzureDetails) error {
		t.Helper()
		if o.TenantId != expected {
			return fmt.Errorf("expected tenant id: %v; got: %v", expected, o.TenantId)
		}
		return nil
	})
	return s
}

func (s *StorageIntegrationAzureDetailsAssert) HasConsentUrl(expected string) *StorageIntegrationAzureDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.StorageIntegrationAzureDetails) error {
		t.Helper()
		if o.ConsentUrl != expected {
			return fmt.Errorf("expected consent url: %v; got: %v", expected, o.ConsentUrl)
		}
		return nil
	})
	return s
}

func (s *StorageIntegrationAzureDetailsAssert) HasMultiTenantAppName(expected string) *StorageIntegrationAzureDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.StorageIntegrationAzureDetails) error {
		t.Helper()
		if o.MultiTenantAppName != expected {
			return fmt.Errorf("expected multi tenant app name: %v; got: %v", expected, o.MultiTenantAppName)
		}
		return nil
	})
	return s
}
