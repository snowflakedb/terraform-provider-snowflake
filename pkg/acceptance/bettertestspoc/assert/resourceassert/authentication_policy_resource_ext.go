package resourceassert

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (s *AuthenticationPolicyResourceAssert) HasAuthenticationMethods(expected ...sdk.AuthenticationMethodsOption) *AuthenticationPolicyResourceAssert {
	s.AddAssertion(assert.ValueSet("authentication_methods.#", fmt.Sprintf("%d", len(expected))))
	for _, v := range expected {
		s.AddAssertion(assert.SetElem("authentication_methods", string(v)))
	}
	return s
}

func (s *AuthenticationPolicyResourceAssert) HasMfaAuthenticationMethods(expected ...sdk.MfaAuthenticationMethodsOption) *AuthenticationPolicyResourceAssert {
	s.AddAssertion(assert.ValueSet("mfa_authentication_methods.#", fmt.Sprintf("%d", len(expected))))
	for _, v := range expected {
		s.AddAssertion(assert.SetElem("mfa_authentication_methods", string(v)))
	}
	return s
}

func (s *AuthenticationPolicyResourceAssert) HasClientTypes(expected ...sdk.ClientTypesOption) *AuthenticationPolicyResourceAssert {
	s.AddAssertion(assert.ValueSet("client_types.#", fmt.Sprintf("%d", len(expected))))
	for _, v := range expected {
		s.AddAssertion(assert.SetElem("client_types", string(v)))
	}
	return s
}

func (s *AuthenticationPolicyResourceAssert) HasSecurityIntegrations(expected ...string) *AuthenticationPolicyResourceAssert {
	s.AddAssertion(assert.ValueSet("security_integrations.#", fmt.Sprintf("%d", len(expected))))
	for _, v := range expected {
		s.AddAssertion(assert.SetElem("security_integrations", v))
	}
	return s
}
