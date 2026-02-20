package resourceshowoutputassert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func NetworkRulesDatasourceDescribeOutput(t *testing.T, datasourceReference string) *NetworkRuleDescOutputAssert {
	t.Helper()

	n := NetworkRuleDescOutputAssert{
		ResourceAssert: assert.NewDatasourceAssert(datasourceReference, "describe_output", "network_rules.0."),
	}
	n.AddAssertion(assert.ValueSet("describe_output.#", "1"))
	return &n
}
