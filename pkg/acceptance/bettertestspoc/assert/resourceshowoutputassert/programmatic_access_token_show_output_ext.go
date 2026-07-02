package resourceshowoutputassert

func (p *ProgrammaticAccessTokenShowOutputAssert) HasExpiresAtNotEmpty() *ProgrammaticAccessTokenShowOutputAssert {
	p.ValuePresent("expires_at")
	return p
}

func (p *ProgrammaticAccessTokenShowOutputAssert) HasCreatedOnNotEmpty() *ProgrammaticAccessTokenShowOutputAssert {
	p.ValuePresent("created_on")
	return p
}

func (p *ProgrammaticAccessTokenShowOutputAssert) HasMinsToBypassNetworkPolicyRequirementNotEmpty() *ProgrammaticAccessTokenShowOutputAssert {
	p.ValuePresent("mins_to_bypass_network_policy_requirement")
	return p
}

func (p *ProgrammaticAccessTokenShowOutputAssert) HasRoleRestrictionEmpty() *ProgrammaticAccessTokenShowOutputAssert {
	p.StringValueSet("role_restriction", "")
	return p
}
