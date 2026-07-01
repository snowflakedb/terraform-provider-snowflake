package resourceshowoutputassert

func (s *SessionPolicyShowOutputAssert) HasCreatedOnNotEmpty() *SessionPolicyShowOutputAssert {
	s.ValuePresent("created_on")
	return s
}
