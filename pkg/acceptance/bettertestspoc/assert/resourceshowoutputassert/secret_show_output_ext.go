package resourceshowoutputassert

func (s *SecretShowOutputAssert) HasCreatedOnNotEmpty() *SecretShowOutputAssert {
	s.ValuePresent("created_on")
	return s
}
