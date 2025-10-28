package resourceassert

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func (a *ApiAuthenticationIntegrationWithAuthorizationCodeGrantResourceAssert) HasOauthAllowedScopesLen(len int) *ApiAuthenticationIntegrationWithAuthorizationCodeGrantResourceAssert {
	a.AddAssertion(assert.ValueSet("oauth_allowed_scopes.#", fmt.Sprintf("%d", len)))
	return a
}

func (a *ApiAuthenticationIntegrationWithAuthorizationCodeGrantResourceAssert) HasOauthAllowedScopesElem(index int, value string) *ApiAuthenticationIntegrationWithAuthorizationCodeGrantResourceAssert {
	a.AddAssertion(assert.ValueSet(fmt.Sprintf("oauth_allowed_scopes.%d", index), value))
	return a
}
