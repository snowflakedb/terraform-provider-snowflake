package resourceshowoutputassert

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func (s *DatabaseShowOutputAssert) HasCreatedOnNotEmpty() *DatabaseShowOutputAssert {
	s.AddAssertion(assert.ResourceShowOutputValuePresent("created_on"))
	return s
}

func (s *DatabaseShowOutputAssert) HasIsCurrentNotEmpty() *DatabaseShowOutputAssert {
	s.AddAssertion(assert.ResourceShowOutputValuePresent("is_current"))
	return s
}

func (s *DatabaseShowOutputAssert) HasOwnerNotEmpty() *DatabaseShowOutputAssert {
	s.AddAssertion(assert.ResourceShowOutputValuePresent("owner"))
	return s
}

func (s *DatabaseShowOutputAssert) HasRetentionTimeNotEmpty() *DatabaseShowOutputAssert {
	s.AddAssertion(assert.ResourceShowOutputValuePresent("retention_time"))
	return s
}

func (s *DatabaseShowOutputAssert) HasOwnerRoleTypeNotEmpty() *DatabaseShowOutputAssert {
	s.AddAssertion(assert.ResourceShowOutputValuePresent("owner_role_type"))
	return s
}

func (s *DatabaseShowOutputAssert) HasOriginEmpty() *DatabaseShowOutputAssert {
	s.AddAssertion(assert.ResourceShowOutputValueSet("origin", ""))
	return s
}
