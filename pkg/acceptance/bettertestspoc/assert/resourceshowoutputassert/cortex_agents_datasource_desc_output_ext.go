package resourceshowoutputassert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func CortexAgentsDatasourceDescribeOutput(t *testing.T, name string) *CortexAgentDescribeOutputAssert {
	t.Helper()

	s := CortexAgentDescribeOutputAssert{
		ResourceAssert: assert.NewDatasourceAssert("data."+name, "describe_output", "cortex_agents.0."),
	}
	s.AddAssertion(assert.ValueSet("describe_output.#", "1"))
	return &s
}
