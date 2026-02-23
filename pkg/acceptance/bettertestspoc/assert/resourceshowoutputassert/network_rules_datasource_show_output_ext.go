package resourceshowoutputassert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func NetworkRulesDatasourceShowOutput(t *testing.T, datasourceReference string) *NetworkRuleShowOutputAssert {
	t.Helper()

	n := NetworkRuleShowOutputAssert{
		ResourceAssert: assert.NewDatasourceAssert(datasourceReference, "show_output", "network_rules.0."),
	}
	n.AddAssertion(assert.ValueSet("show_output.#", "1"))
	return &n
}
