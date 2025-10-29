package resourceshowoutputassert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

// DatabasesDatasourceShowOutput is a temporary workaround to have better show output assertions in data source acceptance tests.
func DatabasesDatasourceShowOutput(t *testing.T, name string) *DatabaseShowOutputAssert {
	t.Helper()

	s := DatabaseShowOutputAssert{
		ResourceAssert: assert.NewDatasourceAssert(name, "show_output", "databases.0."),
	}
	s.AddAssertion(assert.ValueSet("show_output.#", "1"))
	return &s
}

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
