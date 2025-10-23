package resourceassert

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (s *AuthenticationPolicyResourceAssert) HasAuthenticationMethods(expected ...sdk.AuthenticationMethodsOption) *AuthenticationPolicyResourceAssert {
	s.AddAssertion(assert.ValueSet("authentication_methods.#", fmt.Sprintf("%d", len(expected))))
	for i, v := range expected {
		s.AddAssertion(assert.ValueSet(fmt.Sprintf("authentication_methods.%d", i), string(v)))
	}
	return s
}

func (s *AuthenticationPolicyResourceAssert) HasMfaAuthenticationMethods(expected ...sdk.MfaAuthenticationMethodsOption) *AuthenticationPolicyResourceAssert {
	s.AddAssertion(assert.ValueSet("mfa_authentication_methods.#", fmt.Sprintf("%d", len(expected))))
	for i, v := range expected {
		s.AddAssertion(assert.ValueSet(fmt.Sprintf("mfa_authentication_methods.%d", i), string(v)))
	}
	return s
}

func (s *AuthenticationPolicyResourceAssert) HasClientTypes(expected ...sdk.ClientTypesOption) *AuthenticationPolicyResourceAssert {
	s.AddAssertion(assert.ValueSet("client_types.#", fmt.Sprintf("%d", len(expected))))
	for i, v := range expected {
		s.AddAssertion(assert.ValueSet(fmt.Sprintf("client_types.%d", i), string(v)))
	}
	return s
}

func (s *AuthenticationPolicyResourceAssert) HasSecurityIntegrations(expected ...string) *AuthenticationPolicyResourceAssert {
	s.AddAssertion(assert.ValueSet("security_integrations.#", fmt.Sprintf("%d", len(expected))))
	for i, v := range expected {
		s.AddAssertion(assert.ValueSet(fmt.Sprintf("security_integrations.%d", i), v))
	}
	return s
}
