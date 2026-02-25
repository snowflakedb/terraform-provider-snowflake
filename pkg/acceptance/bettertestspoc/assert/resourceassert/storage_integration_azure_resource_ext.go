package resourceassert

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (s *StorageIntegrationAzureResourceAssert) HasStorageAllowedLocationsStorageLocation(expected ...sdk.StorageLocation) *StorageIntegrationAzureResourceAssert {
	return s.HasStorageAllowedLocations(collections.Map(expected, func(v sdk.StorageLocation) string {
		return v.Path
	})...)
}

func (s *StorageIntegrationAzureResourceAssert) HasStorageBlockedLocationsStorageLocation(expected ...sdk.StorageLocation) *StorageIntegrationAzureResourceAssert {
	return s.HasStorageBlockedLocations(collections.Map(expected, func(v sdk.StorageLocation) string {
		return v.Path
	})...)
}
