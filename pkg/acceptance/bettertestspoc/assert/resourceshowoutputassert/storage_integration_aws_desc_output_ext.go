package resourceshowoutputassert

import (
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (s *StorageIntegrationAwsDescribeOutputAssert) HasIamUserArnSet() *StorageIntegrationAwsDescribeOutputAssert {
	s.ValuePresent("iam_user_arn")
	return s
}

func (s *StorageIntegrationAwsDescribeOutputAssert) HasExternalIdSet() *StorageIntegrationAwsDescribeOutputAssert {
	s.ValuePresent("external_id")
	return s
}

func (s *StorageIntegrationAwsDescribeOutputAssert) HasAllowedLocations(expected ...sdk.StorageLocation) *StorageIntegrationAwsDescribeOutputAssert {
	s.StringValueSet("allowed_locations.#", strconv.FormatInt(int64(len(expected)), 10))
	for _, v := range expected {
		s.SetContainsElem("allowed_locations", v.Path)
	}
	return s
}

func (s *StorageIntegrationAwsDescribeOutputAssert) HasBlockedLocations(expected ...sdk.StorageLocation) *StorageIntegrationAwsDescribeOutputAssert {
	s.StringValueSet("blocked_locations.#", strconv.FormatInt(int64(len(expected)), 10))
	for _, v := range expected {
		s.SetContainsElem("blocked_locations", v.Path)
	}
	return s
}
