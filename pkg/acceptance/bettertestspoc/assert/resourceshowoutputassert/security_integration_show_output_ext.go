package resourceshowoutputassert

func (s *SecurityIntegrationShowOutputAssert) HasCreatedOnNotEmpty() *SecurityIntegrationShowOutputAssert {
	s.ValuePresent("created_on")
	return s
}
