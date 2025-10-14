package resourceshowoutputassert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func (c *SemanticViewShowOutputAssert) HasCreatedOnNotEmpty() *SemanticViewShowOutputAssert {
	c.AddAssertion(assert.ResourceShowOutputValuePresent("created_on"))
	return c
}

// SemanticViewsDatasourceShowOutput is a temporary workaround to have better show output assertions in data source acceptance tests.
func SemanticViewsDatasourceShowOutput(t *testing.T, name string) *SemanticViewShowOutputAssert {
	t.Helper()

	s := SemanticViewShowOutputAssert{
		ResourceAssert: assert.NewDatasourceAssert(name, "show_output", "semantic_views.0."),
	}
	s.AddAssertion(assert.ValueSet("show_output.#", "1"))
	return &s
}
