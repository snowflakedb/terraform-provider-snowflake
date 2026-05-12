package resourceassert

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (s *StorageIntegrationGcsResourceAssert) HasStorageAllowedLocationsStorageLocation(expected ...sdk.StorageLocation) *StorageIntegrationGcsResourceAssert {
	return s.HasStorageAllowedLocations(collections.Map(expected, func(v sdk.StorageLocation) string {
		return v.Path
	})...)
}

func (s *StorageIntegrationGcsResourceAssert) HasStorageBlockedLocationsStorageLocation(expected ...sdk.StorageLocation) *StorageIntegrationGcsResourceAssert {
	return s.HasStorageBlockedLocations(collections.Map(expected, func(v sdk.StorageLocation) string {
		return v.Path
	})...)
}
