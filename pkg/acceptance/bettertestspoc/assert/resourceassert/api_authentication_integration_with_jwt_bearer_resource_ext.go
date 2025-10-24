package resourceassert

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func (a *ApiAuthenticationIntegrationWithJwtBearerResourceAssert) HasOauthAllowedScopesLen(len int) *ApiAuthenticationIntegrationWithJwtBearerResourceAssert {
	a.AddAssertion(assert.ValueSet("oauth_allowed_scopes.#", fmt.Sprintf("%d", len)))
	return a
}

func (a *ApiAuthenticationIntegrationWithJwtBearerResourceAssert) HasOauthAllowedScopesElem(index int, value string) *ApiAuthenticationIntegrationWithJwtBearerResourceAssert {
	a.AddAssertion(assert.ValueSet(fmt.Sprintf("oauth_allowed_scopes.%d", index), value))
	return a
}
