package resourceshowoutputassert

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

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
