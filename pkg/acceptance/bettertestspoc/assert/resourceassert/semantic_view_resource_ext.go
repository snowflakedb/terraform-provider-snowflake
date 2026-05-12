package resourceassert

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"

func (s *SemanticViewResourceAssert) HasNoTables() *SemanticViewResourceAssert {
	s.AddAssertion(assert.ValueNotSet("tables"))
	return s
}
