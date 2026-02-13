package resourceassert

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func (a *ApiAuthenticationIntegrationWithClientCredentialsResourceAssert) HasOauthAllowedScopes(values ...string) *ApiAuthenticationIntegrationWithClientCredentialsResourceAssert {
	a.AddAssertion(assert.ValueSet("oauth_allowed_scopes.#", fmt.Sprintf("%d", len(values))))
	for _, value := range values {
		a.AddAssertion(assert.SetElem("oauth_allowed_scopes.*", value))
	}
	return a
}
