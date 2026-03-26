package resourceshowoutputassert

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func (c *OAuthRestAuthenticationDescribeOutputAssert) HasOauthAllowedScopes(expected ...string) *OAuthRestAuthenticationDescribeOutputAssert {
	c.AddAssertion(assert.ResourceDescribeOutputValueSet("rest_authentication.0.oauth_allowed_scopes.#", fmt.Sprintf("%d", len(expected))))
	for i, v := range expected {
		c.AddAssertion(assert.ResourceDescribeOutputValueSet(fmt.Sprintf("rest_authentication.0.oauth_allowed_scopes.%d", i), v))
	}
	return c
}
