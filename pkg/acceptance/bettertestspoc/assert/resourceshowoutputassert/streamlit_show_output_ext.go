package resourceshowoutputassert

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func (s *StreamlitShowOutputAssert) HasCreatedOnNotEmpty() *StreamlitShowOutputAssert {
	s.AddAssertion(assert.ResourceShowOutputValuePresent("created_on"))
	return s
}

func (s *StreamlitShowOutputAssert) HasUrlIdNotEmpty() *StreamlitShowOutputAssert {
	s.AddAssertion(assert.ResourceShowOutputValuePresent("url_id"))
	return s
}

func (s *StreamlitShowOutputAssert) HasOwnerNotEmpty() *StreamlitShowOutputAssert {
	s.AddAssertion(assert.ResourceShowOutputValuePresent("owner"))
	return s
}
