package resourceshowoutputassert

import (
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (s *StorageIntegrationGcsDescribeOutputAssert) HasServiceAccountSet() *StorageIntegrationGcsDescribeOutputAssert {
	s.AddAssertion(assert.ResourceDescribeOutputValuePresent("service_account"))
	return s
}

func (s *StorageIntegrationGcsDescribeOutputAssert) HasAllowedLocations(expected ...sdk.StorageLocation) *StorageIntegrationGcsDescribeOutputAssert {
	s.AddAssertion(assert.ResourceDescribeOutputValueSet("allowed_locations.#", strconv.FormatInt(int64(len(expected)), 10)))
	for _, v := range expected {
		s.AddAssertion(assert.ResourceDescribeOutputSetElem("allowed_locations", v.Path))
	}
	return s
}

func (s *StorageIntegrationGcsDescribeOutputAssert) HasBlockedLocations(expected ...sdk.StorageLocation) *StorageIntegrationGcsDescribeOutputAssert {
	s.AddAssertion(assert.ResourceDescribeOutputValueSet("blocked_locations.#", strconv.FormatInt(int64(len(expected)), 10)))
	for _, v := range expected {
		s.AddAssertion(assert.ResourceDescribeOutputSetElem("blocked_locations", v.Path))
	}
	return s
}
