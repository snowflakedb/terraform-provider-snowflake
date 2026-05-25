package resourceshowoutputassert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func CortexAgentsDatasourceShowOutput(t *testing.T, name string) *CortexAgentShowOutputAssert {
	t.Helper()

	s := CortexAgentShowOutputAssert{
		ResourceAssert: assert.NewDatasourceAssert("data."+name, "show_output", "cortex_agents.0."),
	}
	s.AddAssertion(assert.ValueSet("show_output.#", "1"))
	return &s
}
