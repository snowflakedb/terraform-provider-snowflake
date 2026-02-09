package resourceshowoutputassert

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"

func (s *StorageIntegrationAwsDescribeOutputAssert) HasIamUserArnSet() *StorageIntegrationAwsDescribeOutputAssert {
	s.AddAssertion(assert.ResourceDescribeOutputValuePresent("iam_user_arn"))
	return s
}

func (s *StorageIntegrationAwsDescribeOutputAssert) HasExternalIdSet() *StorageIntegrationAwsDescribeOutputAssert {
	s.AddAssertion(assert.ResourceDescribeOutputValuePresent("external_id"))
	return s
}
