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

type StorageIntegrationGcsDetailsAssert struct {
	*assert.SnowflakeObjectAssert[sdk.StorageIntegrationGcsDetails, sdk.AccountObjectIdentifier]
}

func StorageIntegrationGcsDetails(t *testing.T, id sdk.AccountObjectIdentifier) *StorageIntegrationGcsDetailsAssert {
	t.Helper()
	return &StorageIntegrationGcsDetailsAssert{
		assert.NewSnowflakeObjectAssertWithTestClientObjectProvider(sdk.ObjectType("StorageIntegrationGcsDetails"), id, func(testClient *helpers.TestClient) assert.ObjectProvider[sdk.StorageIntegrationGcsDetails, sdk.AccountObjectIdentifier] {
			return testClient.StorageIntegration.DescribeGcs
		}),
	}
}

func StorageIntegrationGcsDetailsFromObject(t *testing.T, storageIntegrationGcsDetails *sdk.StorageIntegrationGcsDetails) *StorageIntegrationGcsDetailsAssert {
	t.Helper()
	return &StorageIntegrationGcsDetailsAssert{
		assert.NewSnowflakeObjectAssertWithObject(sdk.ObjectType("StorageIntegrationGcsDetails"), storageIntegrationGcsDetails.ID(), storageIntegrationGcsDetails),
	}
}

func (s *StorageIntegrationGcsDetailsAssert) HasEnabled(expected bool) *StorageIntegrationGcsDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.StorageIntegrationGcsDetails) error {
		t.Helper()
		if o.Enabled != expected {
			return fmt.Errorf("expected enabled: %v; got: %v", expected, o.Enabled)
		}
		return nil
	})
	return s
}

func (s *StorageIntegrationGcsDetailsAssert) HasProvider(expected string) *StorageIntegrationGcsDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.StorageIntegrationGcsDetails) error {
		t.Helper()
		if o.Provider != expected {
			return fmt.Errorf("expected provider: %v; got: %v", expected, o.Provider)
		}
		return nil
	})
	return s
}

func (s *StorageIntegrationGcsDetailsAssert) HasAllowedLocations(expected ...string) *StorageIntegrationGcsDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.StorageIntegrationGcsDetails) error {
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

func (s *StorageIntegrationGcsDetailsAssert) HasBlockedLocations(expected ...string) *StorageIntegrationGcsDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.StorageIntegrationGcsDetails) error {
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

func (s *StorageIntegrationGcsDetailsAssert) HasComment(expected string) *StorageIntegrationGcsDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.StorageIntegrationGcsDetails) error {
		t.Helper()
		if o.Comment != expected {
			return fmt.Errorf("expected comment: %v; got: %v", expected, o.Comment)
		}
		return nil
	})
	return s
}

func (s *StorageIntegrationGcsDetailsAssert) HasUsePrivatelinkEndpoint(expected bool) *StorageIntegrationGcsDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.StorageIntegrationGcsDetails) error {
		t.Helper()
		if o.UsePrivatelinkEndpoint != expected {
			return fmt.Errorf("expected use privatelink endpoint: %v; got: %v", expected, o.UsePrivatelinkEndpoint)
		}
		return nil
	})
	return s
}

func (s *StorageIntegrationGcsDetailsAssert) HasServiceAccount(expected string) *StorageIntegrationGcsDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.StorageIntegrationGcsDetails) error {
		t.Helper()
		if o.ServiceAccount != expected {
			return fmt.Errorf("expected service account: %v; got: %v", expected, o.ServiceAccount)
		}
		return nil
	})
	return s
}
