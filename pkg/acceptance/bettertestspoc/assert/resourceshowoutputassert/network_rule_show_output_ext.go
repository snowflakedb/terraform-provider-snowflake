package resourceshowoutputassert

func (s *NetworkRuleShowOutputAssert) HasCreatedOnNotEmpty() *NetworkRuleShowOutputAssert {
	s.ValuePresent("created_on")
	return s
}

func (s *NetworkRuleShowOutputAssert) HasCommentEmpty() *NetworkRuleShowOutputAssert {
	s.StringValueSet("comment", "")
	return s
}
