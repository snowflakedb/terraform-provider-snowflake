package resourceshowoutputassert

import (
	"fmt"
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (s *StorageIntegrationAwsDescribeOutputAssert) HasIamUserArnSet() *StorageIntegrationAwsDescribeOutputAssert {
	s.AddAssertion(assert.ResourceDescribeOutputValuePresent("iam_user_arn"))
	return s
}

func (s *StorageIntegrationAwsDescribeOutputAssert) HasExternalIdSet() *StorageIntegrationAwsDescribeOutputAssert {
	s.AddAssertion(assert.ResourceDescribeOutputValuePresent("external_id"))
	return s
}

func (s *StorageIntegrationAwsDescribeOutputAssert) HasAllowedLocations(expected ...sdk.StorageLocation) *StorageIntegrationAwsDescribeOutputAssert {
	s.AddAssertion(assert.ResourceDescribeOutputValueSet("allowed_locations.#", strconv.FormatInt(int64(len(expected)), 10)))
	for i, v := range expected {
		s.AddAssertion(assert.ResourceDescribeOutputValueSet(fmt.Sprintf("allowed_locations.%d", i), v.Path))
	}
	return s
}

func (s *StorageIntegrationAwsDescribeOutputAssert) HasBlockedLocations(expected ...sdk.StorageLocation) *StorageIntegrationAwsDescribeOutputAssert {
	s.AddAssertion(assert.ResourceDescribeOutputValueSet("blocked_locations.#", strconv.FormatInt(int64(len(expected)), 10)))
	for i, v := range expected {
		s.AddAssertion(assert.ResourceDescribeOutputValueSet(fmt.Sprintf("blocked_locations.%d", i), v.Path))
	}
	return s
}
