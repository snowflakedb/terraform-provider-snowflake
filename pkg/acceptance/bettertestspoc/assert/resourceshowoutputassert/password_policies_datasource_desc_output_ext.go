package resourceshowoutputassert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func PasswordPoliciesDatasourceDescribeOutput(t *testing.T, name string) *PasswordPolicyDescribeOutputAssert {
	t.Helper()

	s := PasswordPolicyDescribeOutputAssert{
		ResourceAssert: assert.NewDatasourceAssert("data."+name, "describe_output", "password_policies.0."),
	}
	s.AddAssertion(assert.ValueSet("describe_output.#", "1"))
	return &s
}
