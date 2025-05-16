package resourceassert

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/assert"
)

func (s *SecretWithClientCredentialsResourceAssert) HasOauthScopesLength(len int) *SecretWithClientCredentialsResourceAssert {
	s.AddAssertion(assert.ValueSet("oauth_scopes.#", fmt.Sprintf("%d", len)))
	return s
}
