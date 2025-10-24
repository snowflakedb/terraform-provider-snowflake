package resourceassert

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func (a *ApiAuthenticationIntegrationWithClientCredentialsResourceAssert) HasOauthAllowedScopesLen(len int) *ApiAuthenticationIntegrationWithClientCredentialsResourceAssert {
	a.AddAssertion(assert.ValueSet("oauth_allowed_scopes.#", fmt.Sprintf("%d", len)))
	return a
}

func (a *ApiAuthenticationIntegrationWithClientCredentialsResourceAssert) HasOauthAllowedScopesElem(index int, value string) *ApiAuthenticationIntegrationWithClientCredentialsResourceAssert {
	a.AddAssertion(assert.ValueSet(fmt.Sprintf("oauth_allowed_scopes.%d", index), value))
	return a
}
