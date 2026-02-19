package resourceassert

import (
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (s *StorageIntegrationAwsResourceAssert) HasStorageAllowedLocations(expected ...sdk.StorageLocation) *StorageIntegrationAwsResourceAssert {
	s.AddAssertion(assert.ValueSet("storage_allowed_locations.#", strconv.FormatInt(int64(len(expected)), 10)))
	for _, v := range expected {
		s.AddAssertion(assert.SetElem("storage_allowed_locations", v.Path))
	}
	return s
}

func (s *StorageIntegrationAwsResourceAssert) HasStorageBlockedLocations(expected ...sdk.StorageLocation) *StorageIntegrationAwsResourceAssert {
	s.AddAssertion(assert.ValueSet("storage_blocked_locations.#", strconv.FormatInt(int64(len(expected)), 10)))
	for _, v := range expected {
		s.AddAssertion(assert.SetElem("storage_blocked_locations", v.Path))
	}
	return s
}
