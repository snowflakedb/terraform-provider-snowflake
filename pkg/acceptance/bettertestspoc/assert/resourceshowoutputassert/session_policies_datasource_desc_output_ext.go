package resourceshowoutputassert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func SessionPoliciesDatasourceDescribeOutput(t *testing.T, name string) *SessionPolicyDescribeOutputAssert {
	t.Helper()

	s := SessionPolicyDescribeOutputAssert{
		ResourceAssert: assert.NewDatasourceAssert("data."+name, "describe_output", "session_policies.0."),
	}
	s.AddAssertion(assert.ValueSet("describe_output.#", "1"))
	return &s
}
