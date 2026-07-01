package resourceshowoutputassert

func (p *PostgresInstanceShowOutputAssert) HasCreatedOnNotEmpty() *PostgresInstanceShowOutputAssert {
	p.ValuePresent("created_on")
	return p
}
