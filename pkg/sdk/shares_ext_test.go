package sdk

func init() {
	id := sharesTestIdAccountObjectIdentifier

	sharesTests.Create.
		withExpectedSqlf(
			case_Shares_sql_Create_basic,
			`CREATE SHARE %s`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Shares_sql_Create_all,
			func(opts *CreateShareOptions) {
				opts.OrReplace = new(true)
				opts.Comment = new("comment")
			},
			`CREATE OR REPLACE SHARE %s COMMENT = 'comment'`, id.FullyQualifiedName(),
		)

	accountId1 := NewAccountIdentifier("my-org", "myaccount")
	accountId2 := NewAccountIdentifier("my-org", "myaccount2")
	tagId := randomSchemaObjectIdentifier()

	sharesTests.Alter.
		withDefaultOpts(func() *AlterShareOptions {
			return &AlterShareOptions{
				IfExists: new(true),
				name:     id,
			}
		}).
		withModifyAndExpectedSqlf(
			case_Shares_sql_Alter_Add,
			func(opts *AlterShareOptions) {
				opts.Add = &ShareAdd{Accounts: []AccountIdentifier{accountId1}, ShareRestrictions: new(true)}
			},
			`ALTER SHARE IF EXISTS %s ADD ACCOUNTS = %s SHARE_RESTRICTIONS = true`,
			id.FullyQualifiedName(), accountId1.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Shares_sql_Alter_Remove,
			func(opts *AlterShareOptions) {
				opts.Remove = &ShareRemove{Accounts: []AccountIdentifier{accountId1, accountId2}}
			},
			`ALTER SHARE IF EXISTS %s REMOVE ACCOUNTS = %s, %s`,
			id.FullyQualifiedName(), accountId1.FullyQualifiedName(), accountId2.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Shares_sql_Alter_Set,
			func(opts *AlterShareOptions) {
				opts.Set = &ShareSet{Accounts: []AccountIdentifier{accountId1}, Comment: new("comment")}
			},
			`ALTER SHARE IF EXISTS %s SET ACCOUNTS = %s COMMENT = 'comment'`,
			id.FullyQualifiedName(), accountId1.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Shares_sql_Alter_Unset,
			func(opts *AlterShareOptions) {
				opts.Unset = &ShareUnset{Comment: new(true)}
			},
			`ALTER SHARE IF EXISTS %s UNSET COMMENT`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Shares_sql_Alter_SetTags,
			func(opts *AlterShareOptions) {
				opts.SetTags = []TagAssociation{{Name: tagId, Value: "v1"}}
			},
			`ALTER SHARE IF EXISTS %s SET TAG %s = 'v1'`,
			id.FullyQualifiedName(), tagId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Shares_sql_Alter_UnsetTags,
			func(opts *AlterShareOptions) {
				opts.UnsetTags = []ObjectIdentifier{tagId}
			},
			`ALTER SHARE IF EXISTS %s UNSET TAG %s`,
			id.FullyQualifiedName(), tagId.FullyQualifiedName(),
		)

	sharesTests.Drop.
		withExpectedSqlf(
			case_Shares_sql_Drop_basic,
			`DROP SHARE %s`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Shares_sql_Drop_all,
			func(opts *DropShareOptions) { opts.IfExists = new(true) },
			`DROP SHARE IF EXISTS %s`, id.FullyQualifiedName(),
		)

	sharesTests.Show.
		withExpectedSql(case_Shares_sql_Show_basic, `SHOW SHARES`).
		withModifyAndExpectedSqlf(
			case_Shares_sql_Show_all,
			func(opts *ShowShareOptions) {
				opts.Like = &Like{Pattern: new("pattern")}
				opts.StartsWith = new("prefix")
				opts.Limit = &LimitFrom{Rows: new(10), From: new("from")}
			},
			`SHOW SHARES LIKE 'pattern' STARTS WITH 'prefix' LIMIT 10 FROM 'from'`,
		).
		withModifyAndExpectedSqlf(
			case_Shares_sql_Show_Like,
			func(opts *ShowShareOptions) { opts.Like = &Like{Pattern: new("pattern")} },
			`SHOW SHARES LIKE 'pattern'`,
		).
		withModifyAndExpectedSqlf(
			case_Shares_sql_Show_StartsWith,
			func(opts *ShowShareOptions) { opts.StartsWith = new("prefix") },
			`SHOW SHARES STARTS WITH 'prefix'`,
		).
		withModifyAndExpectedSqlf(
			case_Shares_sql_Show_Limit,
			func(opts *ShowShareOptions) { opts.Limit = &LimitFrom{Rows: new(10), From: new("from")} },
			`SHOW SHARES LIMIT 10 FROM 'from'`,
		)

	externalShareId := randomExternalObjectIdentifier()

	sharesTests.Describe.
		withExpectedSqlf(
			case_Shares_sql_Describe_basic,
			`DESCRIBE SHARE %s`, sharesTestIdObjectIdentifier.FullyQualifiedName(),
		).
		// covers the DescribeConsumer path, which passes an ExternalObjectIdentifier through the same options struct
		withAdditionalSqlCasef(
			"sql_Describe_external",
			func(opts *DescribeShareOptions) { opts.name = externalShareId },
			`DESCRIBE SHARE %s`, externalShareId.FullyQualifiedName(),
		)
}
