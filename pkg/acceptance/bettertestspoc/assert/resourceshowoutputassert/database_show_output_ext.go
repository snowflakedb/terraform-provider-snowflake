package resourceshowoutputassert

func (s *DatabaseShowOutputAssert) HasCreatedOnNotEmpty() *DatabaseShowOutputAssert {
	s.ValuePresent("created_on")
	return s
}

func (s *DatabaseShowOutputAssert) HasIsCurrentNotEmpty() *DatabaseShowOutputAssert {
	s.ValuePresent("is_current")
	return s
}

func (s *DatabaseShowOutputAssert) HasOwnerNotEmpty() *DatabaseShowOutputAssert {
	s.ValuePresent("owner")
	return s
}

func (s *DatabaseShowOutputAssert) HasRetentionTimeNotEmpty() *DatabaseShowOutputAssert {
	s.ValuePresent("retention_time")
	return s
}

func (s *DatabaseShowOutputAssert) HasOwnerRoleTypeNotEmpty() *DatabaseShowOutputAssert {
	s.ValuePresent("owner_role_type")
	return s
}

func (s *DatabaseShowOutputAssert) HasOriginEmpty() *DatabaseShowOutputAssert {
	s.ValueSet("origin", "")
	return s
}
