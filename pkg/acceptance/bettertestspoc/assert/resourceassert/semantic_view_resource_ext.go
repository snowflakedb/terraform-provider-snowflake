package resourceassert

func (s *SemanticViewResourceAssert) HasNoTables() *SemanticViewResourceAssert {
	s.ValueNotSet("tables")
	return s
}
