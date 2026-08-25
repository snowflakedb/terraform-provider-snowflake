package sdk

func init() {
	id := connectionsTestIdAccountObjectIdentifier
	externalId := randomExternalObjectIdentifier()
	enableAccountId := randomAccountIdentifier()
	enableSecondAccountId := randomAccountIdentifier()
	disableAccountId := randomAccountIdentifier()

	connectionsTests.Create.
		withExpectedSqlf(
			case_Connections_sql_Create_basic,
			"CREATE CONNECTION %s", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Connections_sql_Create_all,
			func(opts *CreateConnectionOptions) {
				opts.IfNotExists = Bool(true)
				opts.Comment = String("comment")
			},
			"CREATE CONNECTION IF NOT EXISTS %s COMMENT = 'comment'", id.FullyQualifiedName(),
		)

	connectionsTests.Create.
		withAdditionalSqlCasef(
			"sql_Create_asReplicaOf",
			func(opts *CreateConnectionOptions) { opts.AsReplicaOf = &externalId },
			"CREATE CONNECTION %s AS REPLICA OF %s", id.FullyQualifiedName(), externalId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Create_asReplicaOf_all",
			func(opts *CreateConnectionOptions) {
				opts.IfNotExists = Bool(true)
				opts.AsReplicaOf = &externalId
				opts.Comment = String("comment")
			},
			"CREATE CONNECTION IF NOT EXISTS %s AS REPLICA OF %s COMMENT = 'comment'", id.FullyQualifiedName(), externalId.FullyQualifiedName(),
		)

	connectionsTests.Alter.
		withModifyAndExpectedSqlf(
			case_Connections_sql_Alter_EnableConnectionFailover,
			func(opts *AlterConnectionOptions) {
				opts.EnableConnectionFailover = &EnableConnectionFailover{
					ToAccounts: []AccountIdentifier{enableAccountId, enableSecondAccountId},
				}
			},
			"ALTER CONNECTION %s ENABLE FAILOVER TO ACCOUNTS %s, %s", id.FullyQualifiedName(), enableAccountId.FullyQualifiedName(), enableSecondAccountId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Connections_sql_Alter_DisableConnectionFailover,
			func(opts *AlterConnectionOptions) { opts.DisableConnectionFailover = &DisableConnectionFailover{} },
			"ALTER CONNECTION %s DISABLE FAILOVER", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Connections_sql_Alter_Primary,
			func(opts *AlterConnectionOptions) { opts.Primary = Bool(true) },
			"ALTER CONNECTION %s PRIMARY", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Connections_sql_Alter_Set,
			func(opts *AlterConnectionOptions) { opts.Set = &ConnectionSet{Comment: String("test comment")} },
			"ALTER CONNECTION %s SET COMMENT = 'test comment'", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Connections_sql_Alter_Unset,
			func(opts *AlterConnectionOptions) { opts.Unset = &ConnectionUnset{Comment: Bool(true)} },
			"ALTER CONNECTION %s UNSET COMMENT", id.FullyQualifiedName(),
		)

	connectionsTests.Alter.
		withAdditionalSqlCasef(
			"sql_Alter_DisableConnectionFailover_toAccounts",
			func(opts *AlterConnectionOptions) {
				opts.DisableConnectionFailover = &DisableConnectionFailover{
					ToAccounts: &ToAccounts{[]AccountIdentifier{disableAccountId}},
				}
			},
			"ALTER CONNECTION %s DISABLE FAILOVER TO ACCOUNTS %s", id.FullyQualifiedName(), disableAccountId.FullyQualifiedName(),
		)

	connectionsTests.Drop.
		withExpectedSqlf(
			case_Connections_sql_Drop_basic,
			"DROP CONNECTION %s", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Connections_sql_Drop_all,
			func(opts *DropConnectionOptions) { opts.IfExists = Bool(true) },
			"DROP CONNECTION IF EXISTS %s", id.FullyQualifiedName(),
		)

	connectionsTests.Show.
		withExpectedSql(case_Connections_sql_Show_basic, "SHOW CONNECTIONS").
		withModifyAndExpectedSqlf(
			case_Connections_sql_Show_all,
			func(opts *ShowConnectionOptions) { opts.Like = &Like{String("test_connection_name")} },
			"SHOW CONNECTIONS LIKE 'test_connection_name'",
		).
		withModifyAndExpectedSqlf(
			case_Connections_sql_Show_Like,
			func(opts *ShowConnectionOptions) { opts.Like = &Like{String("test_connection_name")} },
			"SHOW CONNECTIONS LIKE 'test_connection_name'",
		)
}
