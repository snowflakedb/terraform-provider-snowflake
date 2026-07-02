package resourceshowoutputassert

import (
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (s *StorageIntegrationGcsDescribeOutputAssert) HasServiceAccountSet() *StorageIntegrationGcsDescribeOutputAssert {
	s.ValuePresent("service_account")
	return s
}

func (s *StorageIntegrationGcsDescribeOutputAssert) HasAllowedLocations(expected ...sdk.StorageLocation) *StorageIntegrationGcsDescribeOutputAssert {
	s.StringValueSet("allowed_locations.#", strconv.FormatInt(int64(len(expected)), 10))
	for _, v := range expected {
		s.SetContainsElem("allowed_locations", v.Path)
	}
	return s
}

func (s *StorageIntegrationGcsDescribeOutputAssert) HasBlockedLocations(expected ...sdk.StorageLocation) *StorageIntegrationGcsDescribeOutputAssert {
	s.StringValueSet("blocked_locations.#", strconv.FormatInt(int64(len(expected)), 10))
	for _, v := range expected {
		s.SetContainsElem("blocked_locations", v.Path)
	}
	return s
}
