package resourceshowoutputassert

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"

func (s *SecurityIntegrationShowOutputAssert) HasCreatedOnNotEmpty() *SecurityIntegrationShowOutputAssert {
	s.AddAssertion(assert.ValuePresent("show_output.0.created_on"))
	return s
}
