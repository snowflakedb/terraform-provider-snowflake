package resourceshowoutputassert

func (p *PostgresInstanceShowOutputAssert) HasCreatedOnNotEmpty() *PostgresInstanceShowOutputAssert {
	p.ValuePresent("created_on")
	return p
}

func (p *PostgresInstanceShowOutputAssert) HasIsHa(expected bool) *PostgresInstanceShowOutputAssert {
	p.BoolValueSet("is_ha", expected)
	return p
}
