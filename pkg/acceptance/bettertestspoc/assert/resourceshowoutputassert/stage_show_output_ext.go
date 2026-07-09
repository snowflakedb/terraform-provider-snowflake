package resourceshowoutputassert

func (s *StageShowOutputAssert) HasCreatedOnNotEmpty() *StageShowOutputAssert {
	s.ValuePresent("created_on")
	return s
}

func (s *StageShowOutputAssert) HasCommentEmpty() *StageShowOutputAssert {
	s.StringValueSet("comment", "")
	return s
}

func (s *StageShowOutputAssert) HasStorageIntegrationEmpty() *StageShowOutputAssert {
	s.StringValueSet("storage_integration", "")
	return s
}

func (s *StageShowOutputAssert) HasRegionEmpty() *StageShowOutputAssert {
	s.StringValueSet("region", "")
	return s
}

func (s *StageShowOutputAssert) HasCloudEmpty() *StageShowOutputAssert {
	s.StringValueSet("cloud", "")
	return s
}

func (s *StageShowOutputAssert) HasEndpointEmpty() *StageShowOutputAssert {
	s.StringValueSet("endpoint", "")
	return s
}

func (s *StageShowOutputAssert) HasUrlEmpty() *StageShowOutputAssert {
	s.StringValueSet("url", "")
	return s
}
