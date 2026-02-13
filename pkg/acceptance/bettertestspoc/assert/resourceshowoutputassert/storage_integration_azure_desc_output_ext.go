package resourceshowoutputassert

import (
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (s *StorageIntegrationAzureDescribeOutputAssert) HasConsentUrlSet() *StorageIntegrationAzureDescribeOutputAssert {
	s.AddAssertion(assert.ResourceDescribeOutputValuePresent("consent_url"))
	return s
}

func (s *StorageIntegrationAzureDescribeOutputAssert) HasMultiTenantAppNameSet() *StorageIntegrationAzureDescribeOutputAssert {
	s.AddAssertion(assert.ResourceDescribeOutputValuePresent("multi_tenant_app_name"))
	return s
}

func (s *StorageIntegrationAzureDescribeOutputAssert) HasAllowedLocations(expected ...sdk.StorageLocation) *StorageIntegrationAzureDescribeOutputAssert {
	s.AddAssertion(assert.ResourceDescribeOutputValueSet("allowed_locations.#", strconv.FormatInt(int64(len(expected)), 10)))
	for _, v := range expected {
		s.AddAssertion(assert.ResourceDescribeOutputSetElem("allowed_locations.*", v.Path))
	}
	return s
}

func (s *StorageIntegrationAzureDescribeOutputAssert) HasBlockedLocations(expected ...sdk.StorageLocation) *StorageIntegrationAzureDescribeOutputAssert {
	s.AddAssertion(assert.ResourceDescribeOutputValueSet("blocked_locations.#", strconv.FormatInt(int64(len(expected)), 10)))
	for _, v := range expected {
		s.AddAssertion(assert.ResourceDescribeOutputSetElem("blocked_locations.*", v.Path))
	}
	return s
}
