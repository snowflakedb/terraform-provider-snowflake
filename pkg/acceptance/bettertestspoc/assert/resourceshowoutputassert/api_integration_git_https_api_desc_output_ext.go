package resourceshowoutputassert

import (
	"fmt"
	"strconv"
)

func (a *ApiIntegrationGitHttpsApiDescribeOutputAssert) HasOauthAllowedScopes(expected ...string) *ApiIntegrationGitHttpsApiDescribeOutputAssert {
	a.StringValueSet("oauth_allowed_scopes.#", strconv.Itoa(len(expected)))
	for i, scope := range expected {
		a.StringValueSet(fmt.Sprintf("oauth_allowed_scopes.%d", i), scope)
	}
	return a
}

func (a *ApiIntegrationGitHttpsApiDescribeOutputAssert) HasAllowedPrefixes(expected ...string) *ApiIntegrationGitHttpsApiDescribeOutputAssert {
	a.StringValueSet("allowed_prefixes.#", strconv.Itoa(len(expected)))
	for i, prefix := range expected {
		a.StringValueSet(fmt.Sprintf("allowed_prefixes.%d", i), prefix)
	}
	return a
}

func (a *ApiIntegrationGitHttpsApiDescribeOutputAssert) HasBlockedPrefixes(expected ...string) *ApiIntegrationGitHttpsApiDescribeOutputAssert {
	a.StringValueSet("blocked_prefixes.#", strconv.Itoa(len(expected)))
	for i, prefix := range expected {
		a.StringValueSet(fmt.Sprintf("blocked_prefixes.%d", i), prefix)
	}
	return a
}
