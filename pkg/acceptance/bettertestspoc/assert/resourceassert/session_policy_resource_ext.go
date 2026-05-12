package resourceassert

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func (s *SessionPolicyResourceAssert) HasNoAllowedSecondaryRoles() *SessionPolicyResourceAssert {
	s.AddAssertion(assert.ValueSet("allowed_secondary_roles.#", "1"))
	s.AddAssertion(assert.ValueSet("allowed_secondary_roles.0.none", "true"))
	return s
}

func (s *SessionPolicyResourceAssert) HasAllowedSecondaryRoles(expected ...string) *SessionPolicyResourceAssert {
	if len(expected) == 0 {
		return s.HasNoAllowedSecondaryRoles()
	}
	s.AddAssertion(assert.ValueSet("allowed_secondary_roles.#", "1"))
	s.SetContainsExactlyStringValues("allowed_secondary_roles.0.roles", expected...)
	return s
}

func (s *SessionPolicyResourceAssert) HasAllAllowedSecondaryRoles() *SessionPolicyResourceAssert {
	s.AddAssertion(assert.ValueSet("allowed_secondary_roles.#", "1"))
	s.AddAssertion(assert.ValueSet("allowed_secondary_roles.0.all", "true"))
	return s
}

func (s *SessionPolicyResourceAssert) HasNoBlockedSecondaryRoles() *SessionPolicyResourceAssert {
	s.AddAssertion(assert.ValueSet("blocked_secondary_roles.#", "1"))
	s.AddAssertion(assert.ValueSet("blocked_secondary_roles.0.none", "true"))
	return s
}

func (s *SessionPolicyResourceAssert) HasBlockedSecondaryRoles(expected ...string) *SessionPolicyResourceAssert {
	if len(expected) == 0 {
		return s.HasNoBlockedSecondaryRoles()
	}
	s.AddAssertion(assert.ValueSet("blocked_secondary_roles.#", "1"))
	s.SetContainsExactlyStringValues("blocked_secondary_roles.0.roles", expected...)
	return s
}

func (s *SessionPolicyResourceAssert) HasAllBlockedSecondaryRoles() *SessionPolicyResourceAssert {
	s.AddAssertion(assert.ValueSet("blocked_secondary_roles.#", "1"))
	s.AddAssertion(assert.ValueSet("blocked_secondary_roles.0.all", "true"))
	return s
}
