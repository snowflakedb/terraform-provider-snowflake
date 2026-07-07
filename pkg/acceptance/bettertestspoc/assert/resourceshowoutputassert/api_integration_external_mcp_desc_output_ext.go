package resourceshowoutputassert

import (
	"fmt"
	"strconv"
)

func (a *ApiIntegrationExternalMcpDescribeOutputAssert) HasAllowedPrefixes(expected ...string) *ApiIntegrationExternalMcpDescribeOutputAssert {
	a.StringValueSet("allowed_prefixes.#", strconv.Itoa(len(expected)))
	for i, prefix := range expected {
		a.StringValueSet(fmt.Sprintf("allowed_prefixes.%d", i), prefix)
	}
	return a
}

func (a *ApiIntegrationExternalMcpDescribeOutputAssert) HasBlockedPrefixes(expected ...string) *ApiIntegrationExternalMcpDescribeOutputAssert {
	a.StringValueSet("blocked_prefixes.#", strconv.Itoa(len(expected)))
	for i, prefix := range expected {
		a.StringValueSet(fmt.Sprintf("blocked_prefixes.%d", i), prefix)
	}
	return a
}
