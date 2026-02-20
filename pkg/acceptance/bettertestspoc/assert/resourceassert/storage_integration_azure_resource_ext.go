package resourceassert

import (
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (s *StorageIntegrationAzureResourceAssert) HasStorageAllowedLocations(expected ...sdk.StorageLocation) *StorageIntegrationAzureResourceAssert {
	s.AddAssertion(assert.ValueSet("storage_allowed_locations.#", strconv.FormatInt(int64(len(expected)), 10)))
	for _, v := range expected {
		s.AddAssertion(assert.SetElem("storage_allowed_locations", v.Path))
	}
	return s
}

func (s *StorageIntegrationAzureResourceAssert) HasStorageBlockedLocations(expected ...sdk.StorageLocation) *StorageIntegrationAzureResourceAssert {
	s.AddAssertion(assert.ValueSet("storage_blocked_locations.#", strconv.FormatInt(int64(len(expected)), 10)))
	for _, v := range expected {
		s.AddAssertion(assert.SetElem("storage_blocked_locations", v.Path))
	}
	return s
}
