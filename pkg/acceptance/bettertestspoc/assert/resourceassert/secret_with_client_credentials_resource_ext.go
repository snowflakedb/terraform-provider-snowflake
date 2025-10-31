package resourceassert

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func (s *SecretWithClientCredentialsResourceAssert) HasOauthScopes(expected ...string) *SecretWithClientCredentialsResourceAssert {
	s.AddAssertion(assert.ValueSet("oauth_scopes.#", fmt.Sprintf("%d", len(expected))))
	for i, val := range expected {
		s.AddAssertion(assert.ValueSet(fmt.Sprintf("oauth_scopes.%d", i), val))
	}
	return s
}

func (s *SecretWithClientCredentialsResourceAssert) HasOauthScopesLength(len int) *SecretWithClientCredentialsResourceAssert {
	s.AddAssertion(assert.ValueSet("oauth_scopes.#", fmt.Sprintf("%d", len)))
	return s
}
