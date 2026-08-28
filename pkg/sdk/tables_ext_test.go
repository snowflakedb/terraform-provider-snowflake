package sdk

func init() {
	id := tablesTestIdSchemaObjectIdentifier
	databaseId := NewAccountObjectIdentifier(id.DatabaseName())

	tablesTests.DescribeSearchOptimization.
		withExpectedSqlf(
			case_Tables_sql_DescribeSearchOptimization_basic,
			"DESCRIBE SEARCH OPTIMIZATION ON %s", id.FullyQualifiedName(),
		)

	tablesTests.SelectTableConstraints.
		withDefaultOpts(func() *SelectTableConstraintsTableOptions {
			return &SelectTableConstraintsTableOptions{
				Database:    databaseId,
				TableSchema: id.SchemaName(),
				TableName:   id.Name(),
			}
		}).
		withExpectedSqlf(
			case_Tables_sql_SelectTableConstraints_basic,
			"SELECT * FROM %s . INFORMATION_SCHEMA.TABLE_CONSTRAINTS WHERE TABLE_SCHEMA = '%s' AND TABLE_NAME = '%s'",
			databaseId.FullyQualifiedName(), id.SchemaName(), id.Name(),
		)

	tablesTests.SelectCheckConstraints.
		withDefaultOpts(func() *SelectCheckConstraintsTableOptions {
			return &SelectCheckConstraintsTableOptions{
				Database:         databaseId,
				ConstraintSchema: id.SchemaName(),
				ConstraintTable:  id.Name(),
			}
		}).
		withExpectedSqlf(
			case_Tables_sql_SelectCheckConstraints_basic,
			"SELECT * FROM %s . INFORMATION_SCHEMA.CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = '%s' AND CONSTRAINT_TABLE = '%s'",
			databaseId.FullyQualifiedName(), id.SchemaName(), id.Name(),
		)
}
