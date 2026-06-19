package resourceshowoutputassert

import (
	"fmt"
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func (a *ApiIntegrationGitHttpsApiDescribeOutputAssert) HasOauthAllowedScopes(expected ...string) *ApiIntegrationGitHttpsApiDescribeOutputAssert {
	a.AddAssertion(assert.ResourceDescribeOutputValueSet("oauth_allowed_scopes.#", strconv.Itoa(len(expected))))
	for i, scope := range expected {
		a.AddAssertion(assert.ResourceDescribeOutputValueSet(fmt.Sprintf("oauth_allowed_scopes.%d", i), scope))
	}
	return a
}

func (a *ApiIntegrationGitHttpsApiDescribeOutputAssert) HasAllowedPrefixes(expected ...string) *ApiIntegrationGitHttpsApiDescribeOutputAssert {
	a.AddAssertion(assert.ResourceDescribeOutputValueSet("allowed_prefixes.#", strconv.Itoa(len(expected))))
	for i, prefix := range expected {
		a.AddAssertion(assert.ResourceDescribeOutputValueSet(fmt.Sprintf("allowed_prefixes.%d", i), prefix))
	}
	return a
}

func (a *ApiIntegrationGitHttpsApiDescribeOutputAssert) HasBlockedPrefixes(expected ...string) *ApiIntegrationGitHttpsApiDescribeOutputAssert {
	a.AddAssertion(assert.ResourceDescribeOutputValueSet("blocked_prefixes.#", strconv.Itoa(len(expected))))
	for i, prefix := range expected {
		a.AddAssertion(assert.ResourceDescribeOutputValueSet(fmt.Sprintf("blocked_prefixes.%d", i), prefix))
	}
	return a
}
