package sdk

func init() {
	primary := randomExternalObjectIdentifier()

	replicationFunctionsTests.ShowReplicationAccounts.
		withExpectedSql(
			case_ReplicationFunctions_sql_ShowReplicationAccounts_basic,
			"SHOW REPLICATION ACCOUNTS",
		).
		withModifyAndExpectedSqlf(
			case_ReplicationFunctions_sql_ShowReplicationAccounts_all,
			func(opts *ShowReplicationAccountsOptions) {},
			"SHOW REPLICATION ACCOUNTS",
		)

	replicationFunctionsTests.ShowReplicationDatabases.
		withExpectedSql(
			case_ReplicationFunctions_sql_ShowReplicationDatabases_basic,
			"SHOW REPLICATION DATABASES",
		).
		withModifyAndExpectedSqlf(
			case_ReplicationFunctions_sql_ShowReplicationDatabases_all,
			func(opts *ShowReplicationDatabasesOptions) {
				opts.Like = &Like{Pattern: new("mydb")}
				opts.WithPrimary = &primary
			},
			"SHOW REPLICATION DATABASES LIKE 'mydb' WITH PRIMARY %s", primary.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_ReplicationFunctions_sql_ShowReplicationDatabases_Like,
			func(opts *ShowReplicationDatabasesOptions) {
				opts.Like = &Like{Pattern: new("mydb")}
			},
			"SHOW REPLICATION DATABASES LIKE 'mydb'",
		).
		withModifyAndExpectedSqlf(
			case_ReplicationFunctions_sql_ShowReplicationDatabases_WithPrimary,
			func(opts *ShowReplicationDatabasesOptions) {
				opts.WithPrimary = &primary
			},
			"SHOW REPLICATION DATABASES WITH PRIMARY %s", primary.FullyQualifiedName(),
		)

	replicationFunctionsTests.ShowRegions.
		withExpectedSql(
			case_ReplicationFunctions_sql_ShowRegions_basic,
			"SHOW REGIONS",
		).
		withModifyAndExpectedSqlf(
			case_ReplicationFunctions_sql_ShowRegions_all,
			func(opts *ShowRegionsOptions) {
				opts.Like = &Like{Pattern: new("mydb")}
			},
			"SHOW REGIONS LIKE 'mydb'",
		).
		withModifyAndExpectedSqlf(
			case_ReplicationFunctions_sql_ShowRegions_Like,
			func(opts *ShowRegionsOptions) {
				opts.Like = &Like{Pattern: new("mydb")}
			},
			"SHOW REGIONS LIKE 'mydb'",
		)
}
