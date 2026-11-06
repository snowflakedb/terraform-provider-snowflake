package resourceshowoutputassert

// HasPostgresTables expects the order DESCRIBE_MIRROR returns: sorted, in Postgres quoting.
func (p *PostgresMirrorDescribeOutputAssert) HasPostgresTables(expected ...string) *PostgresMirrorDescribeOutputAssert {
	p.ListContainsExactlyStringValuesInOrder("postgres_tables", expected...)
	return p
}

// HasPostgresSchemas expects the order DESCRIBE_MIRROR returns: sorted, in Postgres quoting.
func (p *PostgresMirrorDescribeOutputAssert) HasPostgresSchemas(expected ...string) *PostgresMirrorDescribeOutputAssert {
	p.ListContainsExactlyStringValuesInOrder("postgres_schemas", expected...)
	return p
}

func (p *PostgresMirrorDescribeOutputAssert) HasStatusNotEmpty() *PostgresMirrorDescribeOutputAssert {
	p.ValuePresent("status")
	return p
}
