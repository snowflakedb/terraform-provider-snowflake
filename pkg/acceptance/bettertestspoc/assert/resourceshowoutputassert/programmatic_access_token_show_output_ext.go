package resourceshowoutputassert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

// ProgrammaticAccessTokensDatasourceShowOutput is a temporary workaround to have better show output assertions in data source acceptance tests.
func ProgrammaticAccessTokensDatasourceShowOutput(t *testing.T, name string) *ProgrammaticAccessTokenShowOutputAssert {
	t.Helper()

	s := ProgrammaticAccessTokenShowOutputAssert{
		ResourceAssert: assert.NewDatasourceAssert(name, "show_output", "user_programmatic_access_tokens.0."),
	}
	s.AddAssertion(assert.ValueSet("show_output.#", "1"))
	return &s
}

func (p *ProgrammaticAccessTokenShowOutputAssert) HasExpiresAtNotEmpty() *ProgrammaticAccessTokenShowOutputAssert {
	p.AddAssertion(assert.ResourceShowOutputValuePresent("expires_at"))
	return p
}

func (p *ProgrammaticAccessTokenShowOutputAssert) HasCreatedOnNotEmpty() *ProgrammaticAccessTokenShowOutputAssert {
	p.AddAssertion(assert.ResourceShowOutputValuePresent("created_on"))
	return p
}

func (p *ProgrammaticAccessTokenShowOutputAssert) HasMinsToBypassNetworkPolicyRequirementNotEmpty() *ProgrammaticAccessTokenShowOutputAssert {
	p.AddAssertion(assert.ResourceShowOutputValuePresent("mins_to_bypass_network_policy_requirement"))
	return p
}

func (p *ProgrammaticAccessTokenShowOutputAssert) HasRoleRestrictionEmpty() *ProgrammaticAccessTokenShowOutputAssert {
	p.AddAssertion(assert.ResourceShowOutputValueSet("role_restriction", ""))
	return p
}
