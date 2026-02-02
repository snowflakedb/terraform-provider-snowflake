package resourceshowoutputassert

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func (s *StageShowOutputAssert) HasCreatedOnNotEmpty() *StageShowOutputAssert {
	s.AddAssertion(assert.ResourceShowOutputValuePresent("created_on"))
	return s
}

func (s *StageShowOutputAssert) HasCommentEmpty() *StageShowOutputAssert {
	s.AddAssertion(assert.ResourceShowOutputValueSet("comment", ""))
	return s
}
