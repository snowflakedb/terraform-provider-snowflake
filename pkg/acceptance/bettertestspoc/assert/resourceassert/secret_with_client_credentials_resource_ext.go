package resourceassert

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func (s *SecretWithClientCredentialsResourceAssert) HasOauthScopes(expected ...string) *SecretWithClientCredentialsResourceAssert {
	s.AddAssertion(assert.ValueSet("oauth_scopes.#", fmt.Sprintf("%d", len(expected))))
	for _, val := range expected {
		s.AddAssertion(assert.SetElem("oauth_scopes", val))
	}
	return s
}

func (s *SecretWithClientCredentialsResourceAssert) HasOauthScopesLength(len int) *SecretWithClientCredentialsResourceAssert {
	s.AddAssertion(assert.ValueSet("oauth_scopes.#", fmt.Sprintf("%d", len)))
	return s
}
