package resourceshowoutputassert

import (
	"strconv"
)

func (s *SessionPolicyDescribeOutputAssert) HasAllowedSecondaryRoles(expected ...string) *SessionPolicyDescribeOutputAssert {
	s.StringValueSet("allowed_secondary_roles.#", strconv.FormatInt(int64(len(expected)), 10))
	for _, v := range expected {
		s.SetContainsElem("allowed_secondary_roles", v)
	}
	return s
}

func (s *SessionPolicyDescribeOutputAssert) HasBlockedSecondaryRoles(expected ...string) *SessionPolicyDescribeOutputAssert {
	s.StringValueSet("blocked_secondary_roles.#", strconv.FormatInt(int64(len(expected)), 10))
	for _, v := range expected {
		s.SetContainsElem("blocked_secondary_roles", v)
	}
	return s
}
