package resourceshowoutputassert

import (
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func (s *SessionPolicyDescribeOutputAssert) HasAllowedSecondaryRoles(expected ...string) *SessionPolicyDescribeOutputAssert {
	s.AddAssertion(assert.ResourceDescribeOutputValueSet("allowed_secondary_roles.#", strconv.FormatInt(int64(len(expected)), 10)))
	for _, v := range expected {
		s.AddAssertion(assert.ResourceDescribeOutputSetElem("allowed_secondary_roles", v))
	}
	return s
}

func (s *SessionPolicyDescribeOutputAssert) HasBlockedSecondaryRoles(expected ...string) *SessionPolicyDescribeOutputAssert {
	s.AddAssertion(assert.ResourceDescribeOutputValueSet("blocked_secondary_roles.#", strconv.FormatInt(int64(len(expected)), 10)))
	for _, v := range expected {
		s.AddAssertion(assert.ResourceDescribeOutputSetElem("blocked_secondary_roles", v))
	}
	return s
}
