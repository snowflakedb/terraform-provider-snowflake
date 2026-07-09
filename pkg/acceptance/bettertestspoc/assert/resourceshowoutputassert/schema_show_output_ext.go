package resourceshowoutputassert

func (s *SchemaShowOutputAssert) HasCreatedOnNotEmpty() *SchemaShowOutputAssert {
	s.ValuePresent("created_on")
	return s
}

func (s *SchemaShowOutputAssert) HasOwnerNotEmpty() *SchemaShowOutputAssert {
	s.ValuePresent("owner")
	return s
}

func (s *SchemaShowOutputAssert) HasRetentionTimeNotEmpty() *SchemaShowOutputAssert {
	s.ValuePresent("retention_time")
	return s
}

func (s *SchemaShowOutputAssert) HasOwnerRoleTypeNotEmpty() *SchemaShowOutputAssert {
	s.ValuePresent("owner_role_type")
	return s
}
