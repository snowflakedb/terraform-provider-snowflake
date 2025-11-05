package resourceshowoutputassert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func AuthenticationPoliciesDatasourceShowOutput(t *testing.T, name string) *AuthenticationPolicyShowOutputAssert {
	t.Helper()

	a := AuthenticationPolicyShowOutputAssert{
		ResourceAssert: assert.NewDatasourceAssert("data."+name, "show_output", "authentication_policies.0."),
	}
	a.AddAssertion(assert.ValueSet("show_output.#", "1"))
	return &a
}

func (a *AuthenticationPolicyShowOutputAssert) HasCreatedOnNotEmpty() *AuthenticationPolicyShowOutputAssert {
	a.AddAssertion(assert.ResourceShowOutputValuePresent("created_on"))
	return a
}
