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
	// TODO [next PRs]: these are sets so the order varies
	// for i, v := range expected {
	// 	s.AddAssertion(assert.ResourceDescribeOutputValueSet(fmt.Sprintf("allowed_locations.%d", i), v.Path))
	// }
	return s
}

func (s *StorageIntegrationAzureDescribeOutputAssert) HasBlockedLocations(expected ...sdk.StorageLocation) *StorageIntegrationAzureDescribeOutputAssert {
	s.AddAssertion(assert.ResourceDescribeOutputValueSet("blocked_locations.#", strconv.FormatInt(int64(len(expected)), 10)))
	// TODO [next PRs]: these are sets so the order varies
	// for i, v := range expected {
	// 	s.AddAssertion(assert.ResourceDescribeOutputValueSet(fmt.Sprintf("blocked_locations.%d", i), v.Path))
	// }
	return s
}
