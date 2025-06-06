package resourceshowoutputassert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

// ComputePoolsDatasourceShowOutput is a temporary workaround to have better show output assertions in data source acceptance tests.
func ComputePoolsDatasourceShowOutput(t *testing.T, name string) *ComputePoolShowOutputAssert {
	t.Helper()

	s := ComputePoolShowOutputAssert{
		ResourceAssert: assert.NewDatasourceAssert("data."+name, "show_output", "compute_pools.0."),
	}
	s.AddAssertion(assert.ValueSet("show_output.#", "1"))
	return &s
}

func (c *ComputePoolShowOutputAssert) HasCreatedOnNotEmpty() *ComputePoolShowOutputAssert {
	c.AddAssertion(assert.ResourceShowOutputValuePresent("created_on"))
	return c
}

func (c *ComputePoolShowOutputAssert) HasResumedOnNotEmpty() *ComputePoolShowOutputAssert {
	c.AddAssertion(assert.ResourceShowOutputValuePresent("resumed_on"))
	return c
}

func (c *ComputePoolShowOutputAssert) HasUpdatedOnNotEmpty() *ComputePoolShowOutputAssert {
	c.AddAssertion(assert.ResourceShowOutputValuePresent("updated_on"))
	return c
}
