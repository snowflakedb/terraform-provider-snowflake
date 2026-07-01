package resourceshowoutputassert

import (
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (s *StorageIntegrationAzureDescribeOutputAssert) HasConsentUrlSet() *StorageIntegrationAzureDescribeOutputAssert {
	s.ValuePresent("consent_url")
	return s
}

func (s *StorageIntegrationAzureDescribeOutputAssert) HasMultiTenantAppNameSet() *StorageIntegrationAzureDescribeOutputAssert {
	s.ValuePresent("multi_tenant_app_name")
	return s
}

func (s *StorageIntegrationAzureDescribeOutputAssert) HasAllowedLocations(expected ...sdk.StorageLocation) *StorageIntegrationAzureDescribeOutputAssert {
	s.StringValueSet("allowed_locations.#", strconv.FormatInt(int64(len(expected)), 10))
	for _, v := range expected {
		s.SetContainsElem("allowed_locations", v.Path)
	}
	return s
}

func (s *StorageIntegrationAzureDescribeOutputAssert) HasBlockedLocations(expected ...sdk.StorageLocation) *StorageIntegrationAzureDescribeOutputAssert {
	s.StringValueSet("blocked_locations.#", strconv.FormatInt(int64(len(expected)), 10))
	for _, v := range expected {
		s.SetContainsElem("blocked_locations", v.Path)
	}
	return s
}
