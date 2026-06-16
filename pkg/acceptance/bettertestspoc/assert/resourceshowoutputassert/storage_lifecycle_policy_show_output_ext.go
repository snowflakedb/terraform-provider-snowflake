package resourceshowoutputassert

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func (s *StorageLifecyclePolicyShowOutputAssert) HasCreatedOnNotEmpty() *StorageLifecyclePolicyShowOutputAssert {
	s.AddAssertion(assert.ResourceShowOutputValuePresent("created_on"))
	return s
}
