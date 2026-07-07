package resourceshowoutputassert

func (s *StreamlitShowOutputAssert) HasCreatedOnNotEmpty() *StreamlitShowOutputAssert {
	s.ValuePresent("created_on")
	return s
}

func (s *StreamlitShowOutputAssert) HasUrlIdNotEmpty() *StreamlitShowOutputAssert {
	s.ValuePresent("url_id")
	return s
}

func (s *StreamlitShowOutputAssert) HasOwnerNotEmpty() *StreamlitShowOutputAssert {
	s.ValuePresent("owner")
	return s
}
