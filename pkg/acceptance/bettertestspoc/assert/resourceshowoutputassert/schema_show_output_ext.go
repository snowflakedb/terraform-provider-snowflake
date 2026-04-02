package resourceshowoutputassert

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func (s *SchemaShowOutputAssert) HasCreatedOnNotEmpty() *SchemaShowOutputAssert {
	s.AddAssertion(assert.ResourceShowOutputValuePresent("created_on"))
	return s
}

func (s *SchemaShowOutputAssert) HasOwnerNotEmpty() *SchemaShowOutputAssert {
	s.AddAssertion(assert.ResourceShowOutputValuePresent("owner"))
	return s
}

func (s *SchemaShowOutputAssert) HasRetentionTimeNotEmpty() *SchemaShowOutputAssert {
	s.AddAssertion(assert.ResourceShowOutputValuePresent("retention_time"))
	return s
}

func (s *SchemaShowOutputAssert) HasOwnerRoleTypeNotEmpty() *SchemaShowOutputAssert {
	s.AddAssertion(assert.ResourceShowOutputValuePresent("owner_role_type"))
	return s
}
