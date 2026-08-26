package sdk

func init() {
	id := databaseRolesTestIdDatabaseObjectIdentifier
	renameTarget := randomDatabaseObjectIdentifierInDatabase(id.DatabaseId())
	databaseId := id.DatabaseId()
	tagId := NewAccountObjectIdentifier("123")
	tagId2 := NewAccountObjectIdentifier("456")
	accountRoleId := randomAccountObjectIdentifier()
	databaseRoleId := randomDatabaseObjectIdentifier()
	shareId := randomAccountObjectIdentifier()

	databaseRolesTests.Create.
		withExpectedSqlf(
			case_DatabaseRoles_sql_Create_basic,
			`CREATE DATABASE ROLE %s`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_DatabaseRoles_sql_Create_all,
			func(opts *CreateDatabaseRoleOptions) {
				opts.IfNotExists = new(true)
				opts.Comment = new("some comment")
			},
			`CREATE DATABASE ROLE IF NOT EXISTS %s COMMENT = 'some comment'`,
			id.FullyQualifiedName(),
		)

	databaseRolesTests.Alter.
		withAdditionalValidationCase(
			"validation_Alter_RenameTo_differentDatabase",
			func(opts *AlterDatabaseRoleOptions) {
				opts.RenameTo = new(randomDatabaseObjectIdentifier())
			},
			ErrDifferentDatabase,
		).
		withModifyAndExpectedSqlf(
			case_DatabaseRoles_sql_Alter_RenameTo,
			func(opts *AlterDatabaseRoleOptions) { opts.RenameTo = &renameTarget },
			`ALTER DATABASE ROLE %s RENAME TO %s`,
			id.FullyQualifiedName(), renameTarget.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_DatabaseRoles_sql_Alter_Set,
			func(opts *AlterDatabaseRoleOptions) {
				opts.IfExists = new(true)
				opts.Set = &DatabaseRoleSet{Comment: new("new comment")}
			},
			`ALTER DATABASE ROLE IF EXISTS %s SET COMMENT = 'new comment'`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_DatabaseRoles_sql_Alter_Unset,
			func(opts *AlterDatabaseRoleOptions) {
				opts.IfExists = new(true)
				opts.Unset = &DatabaseRoleUnset{Comment: new(true)}
			},
			`ALTER DATABASE ROLE IF EXISTS %s UNSET COMMENT`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_DatabaseRoles_sql_Alter_SetTags,
			func(opts *AlterDatabaseRoleOptions) {
				opts.IfExists = new(true)
				opts.SetTags = []TagAssociation{
					{Name: tagId, Value: "value-123"},
					{Name: tagId2, Value: "value-123"},
				}
			},
			`ALTER DATABASE ROLE IF EXISTS %s SET TAG %s = 'value-123', %s = 'value-123'`,
			id.FullyQualifiedName(), tagId.FullyQualifiedName(), tagId2.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_DatabaseRoles_sql_Alter_UnsetTags,
			func(opts *AlterDatabaseRoleOptions) {
				opts.IfExists = new(true)
				opts.UnsetTags = []ObjectIdentifier{tagId, tagId2}
			},
			`ALTER DATABASE ROLE IF EXISTS %s UNSET TAG %s, %s`,
			id.FullyQualifiedName(), tagId.FullyQualifiedName(), tagId2.FullyQualifiedName(),
		)

	databaseRolesTests.Drop.
		withExpectedSqlf(
			case_DatabaseRoles_sql_Drop_basic,
			`DROP DATABASE ROLE %s`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_DatabaseRoles_sql_Drop_all,
			func(opts *DropDatabaseRoleOptions) { opts.IfExists = new(true) },
			`DROP DATABASE ROLE IF EXISTS %s`, id.FullyQualifiedName(),
		)

	databaseRolesTests.Show.
		withDefaultOpts(func() *ShowDatabaseRoleOptions {
			return &ShowDatabaseRoleOptions{
				Database: databaseId,
			}
		}).
		withAdditionalValidationCase(
			"validation_Show_Like_patternRequired",
			func(opts *ShowDatabaseRoleOptions) { opts.Like = &Like{} },
			ErrPatternRequiredForLikeKeyword,
		).
		withExpectedSqlf(
			case_DatabaseRoles_sql_Show_basic,
			`SHOW DATABASE ROLES IN DATABASE %s`, databaseId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_DatabaseRoles_sql_Show_all,
			func(opts *ShowDatabaseRoleOptions) {
				opts.Like = &Like{Pattern: new(id.Name())}
				opts.Limit = &LimitFrom{Rows: new(10), From: new("name")}
			},
			`SHOW DATABASE ROLES LIKE '%s' IN DATABASE %s LIMIT 10 FROM 'name'`,
			id.Name(), databaseId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_DatabaseRoles_sql_Show_Like,
			func(opts *ShowDatabaseRoleOptions) {
				opts.Like = &Like{Pattern: new(id.Name())}
			},
			`SHOW DATABASE ROLES LIKE '%s' IN DATABASE %s`,
			id.Name(), databaseId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_DatabaseRoles_sql_Show_Limit,
			func(opts *ShowDatabaseRoleOptions) {
				opts.Limit = &LimitFrom{Rows: new(10), From: new("name")}
			},
			`SHOW DATABASE ROLES IN DATABASE %s LIMIT 10 FROM 'name'`,
			databaseId.FullyQualifiedName(),
		)

	databaseRolesTests.Grant.
		withDefaultOpts(func() *GrantDatabaseRoleOptions {
			return &GrantDatabaseRoleOptions{
				name: id,
				To:   DatabaseRoleKindOfRole{DatabaseRoleName: &databaseRoleId},
			}
		}).
		withExpectedSqlf(
			case_DatabaseRoles_sql_Grant_basic,
			`GRANT DATABASE ROLE %s TO DATABASE ROLE %s`,
			id.FullyQualifiedName(), databaseRoleId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Grant_toAccountRole",
			func(opts *GrantDatabaseRoleOptions) {
				opts.To = DatabaseRoleKindOfRole{AccountRoleName: &accountRoleId}
			},
			`GRANT DATABASE ROLE %s TO ROLE %s`,
			id.FullyQualifiedName(), accountRoleId.FullyQualifiedName(),
		)

	databaseRolesTests.Revoke.
		withDefaultOpts(func() *RevokeDatabaseRoleOptions {
			return &RevokeDatabaseRoleOptions{
				name: id,
				From: DatabaseRoleKindOfRole{DatabaseRoleName: &databaseRoleId},
			}
		}).
		withExpectedSqlf(
			case_DatabaseRoles_sql_Revoke_basic,
			`REVOKE DATABASE ROLE %s FROM DATABASE ROLE %s`,
			id.FullyQualifiedName(), databaseRoleId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Revoke_fromAccountRole",
			func(opts *RevokeDatabaseRoleOptions) {
				opts.From = DatabaseRoleKindOfRole{AccountRoleName: &accountRoleId}
			},
			`REVOKE DATABASE ROLE %s FROM ROLE %s`,
			id.FullyQualifiedName(), accountRoleId.FullyQualifiedName(),
		)

	databaseRolesTests.GrantToShare.
		withDefaultOpts(func() *GrantToShareDatabaseRoleOptions {
			return &GrantToShareDatabaseRoleOptions{
				name:  id,
				Share: shareId,
			}
		}).
		withExpectedSqlf(
			case_DatabaseRoles_sql_GrantToShare_basic,
			`GRANT DATABASE ROLE %s TO SHARE %s`,
			id.FullyQualifiedName(), shareId.FullyQualifiedName(),
		)

	databaseRolesTests.RevokeFromShare.
		withDefaultOpts(func() *RevokeFromShareDatabaseRoleOptions {
			return &RevokeFromShareDatabaseRoleOptions{
				name:  id,
				Share: shareId,
			}
		}).
		withExpectedSqlf(
			case_DatabaseRoles_sql_RevokeFromShare_basic,
			`REVOKE DATABASE ROLE %s FROM SHARE %s`,
			id.FullyQualifiedName(), shareId.FullyQualifiedName(),
		)
}
