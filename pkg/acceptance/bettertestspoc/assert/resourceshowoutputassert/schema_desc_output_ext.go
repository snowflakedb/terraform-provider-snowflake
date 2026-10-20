package resourceshowoutputassert

func (s *SchemaDescribeOutputAssert) HasCreatedOnNotEmpty() *SchemaDescribeOutputAssert {
	s.ValuePresent("created_on")
	return s
}

func (s *SchemaDescribeOutputAssert) HasNameNotEmpty() *SchemaDescribeOutputAssert {
	s.ValuePresent("name")
	return s
}
