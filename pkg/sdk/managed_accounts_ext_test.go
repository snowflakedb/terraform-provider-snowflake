package sdk

func init() {
	id := managedAccountsTestIdAccountObjectIdentifier

	managedAccountsTests.Create.
		withDefaultOpts(func() *CreateManagedAccountOptions {
			return &CreateManagedAccountOptions{
				name: id,
				CreateManagedAccountParams: CreateManagedAccountParams{
					AdminName:     "admin",
					AdminPassword: "password",
				},
			}
		}).
		withExpectedSqlf(
			case_ManagedAccounts_sql_Create_basic,
			"CREATE MANAGED ACCOUNT %s ADMIN_NAME = 'admin', ADMIN_PASSWORD = 'password', TYPE = READER",
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_ManagedAccounts_sql_Create_all,
			func(opts *CreateManagedAccountOptions) {
				opts.CreateManagedAccountParams.Comment = String("comment")
			},
			"CREATE MANAGED ACCOUNT %s ADMIN_NAME = 'admin', ADMIN_PASSWORD = 'password', TYPE = READER, COMMENT = 'comment'",
			id.FullyQualifiedName(),
		)

	managedAccountsTests.Drop.
		withExpectedSqlf(
			case_ManagedAccounts_sql_Drop_basic,
			"DROP MANAGED ACCOUNT %s",
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_ManagedAccounts_sql_Drop_all,
			func(opts *DropManagedAccountOptions) { opts.IfExists = Bool(true) },
			"DROP MANAGED ACCOUNT IF EXISTS %s",
			id.FullyQualifiedName(),
		)

	managedAccountsTests.Show.
		withExpectedSql(case_ManagedAccounts_sql_Show_basic, "SHOW MANAGED ACCOUNTS").
		withModifyAndExpectedSqlf(
			case_ManagedAccounts_sql_Show_all,
			func(opts *ShowManagedAccountOptions) { opts.Like = &Like{Pattern: String("myaccount")} },
			"SHOW MANAGED ACCOUNTS LIKE 'myaccount'",
		).
		withModifyAndExpectedSqlf(
			case_ManagedAccounts_sql_Show_Like,
			func(opts *ShowManagedAccountOptions) { opts.Like = &Like{Pattern: String("myaccount")} },
			"SHOW MANAGED ACCOUNTS LIKE 'myaccount'",
		)
}
