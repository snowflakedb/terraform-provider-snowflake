package resourceshowoutputassert

import (
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (s *StorageIntegrationAllDescribeOutputAssert) HasIamUserArnSet() *StorageIntegrationAllDescribeOutputAssert {
	s.ValuePresent("iam_user_arn")
	return s
}

func (s *StorageIntegrationAllDescribeOutputAssert) HasExternalIdSet() *StorageIntegrationAllDescribeOutputAssert {
	s.ValuePresent("external_id")
	return s
}

func (s *StorageIntegrationAllDescribeOutputAssert) HasConsentUrlSet() *StorageIntegrationAllDescribeOutputAssert {
	s.ValuePresent("consent_url")
	return s
}

func (s *StorageIntegrationAllDescribeOutputAssert) HasMultiTenantAppNameSet() *StorageIntegrationAllDescribeOutputAssert {
	s.ValuePresent("multi_tenant_app_name")
	return s
}

func (s *StorageIntegrationAllDescribeOutputAssert) HasServiceAccountSet() *StorageIntegrationAllDescribeOutputAssert {
	s.ValuePresent("service_account")
	return s
}

func (s *StorageIntegrationAllDescribeOutputAssert) HasAllowedLocations(expected ...sdk.StorageLocation) *StorageIntegrationAllDescribeOutputAssert {
	s.StringValueSet("allowed_locations.#", strconv.FormatInt(int64(len(expected)), 10))
	for _, v := range expected {
		s.SetContainsElem("allowed_locations", v.Path)
	}
	return s
}

func (s *StorageIntegrationAllDescribeOutputAssert) HasBlockedLocations(expected ...sdk.StorageLocation) *StorageIntegrationAllDescribeOutputAssert {
	s.StringValueSet("blocked_locations.#", strconv.FormatInt(int64(len(expected)), 10))
	for _, v := range expected {
		s.SetContainsElem("blocked_locations", v.Path)
	}
	return s
}
