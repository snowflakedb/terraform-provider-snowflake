package resourceshowoutputassert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func SessionPoliciesDatasourceShowOutput(t *testing.T, name string) *SessionPolicyShowOutputAssert {
	t.Helper()

	s := SessionPolicyShowOutputAssert{
		ResourceAssert: assert.NewDatasourceAssert("data."+name, "show_output", "session_policies.0."),
	}
	s.AddAssertion(assert.ValueSet("show_output.#", "1"))
	return &s
}
