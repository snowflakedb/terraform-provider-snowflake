package resourceshowoutputassert

import (
	"strconv"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

// StreamsDatasourceShowOutput is a temporary workaround to have better show output assertions in data source acceptance tests.
func StreamsDatasourceShowOutput(t *testing.T, name string) *StreamShowOutputAssert {
	t.Helper()

	s := StreamShowOutputAssert{
		ResourceAssert: assert.NewDatasourceAssert("data."+name, "show_output", "streams.0."),
	}
	s.AddAssertion(assert.ValueSet("show_output.#", "1"))
	return &s
}

func (s *StreamShowOutputAssert) HasCreatedOnNotEmpty() *StreamShowOutputAssert {
	s.AddAssertion(assert.ResourceShowOutputValuePresent("created_on"))
	return s
}

func (s *StreamShowOutputAssert) HasStaleAfterNotEmpty() *StreamShowOutputAssert {
	s.AddAssertion(assert.ResourceShowOutputValuePresent("stale_after"))
	return s
}

func (s *StreamShowOutputAssert) HasBaseTables(ids ...sdk.SchemaObjectIdentifier) *StreamShowOutputAssert {
	s.AddAssertion(assert.ResourceShowOutputValueSet("base_tables.#", strconv.FormatInt(int64(len(ids)), 10)))
	for i := range ids {
		s.AddAssertion(assert.ResourceShowOutputSetElem("base_tables", ids[i].FullyQualifiedName()))
	}
	return s
}
