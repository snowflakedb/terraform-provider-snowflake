package resourceshowoutputassert

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"

func (s *NetworkRuleShowOutputAssert) HasCreatedOnNotEmpty() *NetworkRuleShowOutputAssert {
	s.AddAssertion(assert.ResourceShowOutputValuePresent("created_on"))
	return s
}

func (s *NetworkRuleShowOutputAssert) HasCommentEmpty() *NetworkRuleShowOutputAssert {
	s.AddAssertion(assert.ResourceShowOutputValueSet("comment", ""))
	return s
}
