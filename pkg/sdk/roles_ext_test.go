package sdk

func init() {
	id := rolesTestIdAccountObjectIdentifier
	tagId := randomSchemaObjectIdentifier()
	tagId2 := randomSchemaObjectIdentifier()
	renameTarget := randomAccountObjectIdentifier()
	classId := randomSchemaObjectIdentifier()
	grantedRoleId := randomAccountObjectIdentifier()
	grantedUserId := randomAccountObjectIdentifier()

	rolesTests.Create.
		withExpectedSqlf(
			case_Roles_sql_Create_basic,
			`CREATE ROLE %s`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Roles_sql_Create_all,
			func(opts *CreateRoleOptions) {
				opts.IfNotExists = new(true)
				opts.Comment = new("comment")
				opts.Tag = []TagAssociation{{Name: tagId, Value: "v1"}}
			},
			`CREATE ROLE IF NOT EXISTS %s COMMENT = 'comment' TAG (%s = 'v1')`,
			id.FullyQualifiedName(), tagId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Create_orReplace",
			func(opts *CreateRoleOptions) { opts.OrReplace = new(true) },
			`CREATE OR REPLACE ROLE %s`, id.FullyQualifiedName(),
		)

	rolesTests.Alter.
		withModifyAndExpectedSqlf(
			case_Roles_sql_Alter_RenameTo,
			func(opts *AlterRoleOptions) { opts.RenameTo = &renameTarget },
			`ALTER ROLE %s RENAME TO %s`, id.FullyQualifiedName(), renameTarget.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Roles_sql_Alter_SetComment,
			func(opts *AlterRoleOptions) { opts.SetComment = new("some comment") },
			`ALTER ROLE %s SET COMMENT = 'some comment'`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Roles_sql_Alter_SetTags,
			func(opts *AlterRoleOptions) {
				opts.SetTags = []TagAssociation{
					{Name: tagId, Value: "tag-value"},
					{Name: tagId2, Value: "tag-value2"},
				}
			},
			`ALTER ROLE %s SET TAG %s = 'tag-value', %s = 'tag-value2'`,
			id.FullyQualifiedName(), tagId.FullyQualifiedName(), tagId2.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Roles_sql_Alter_UnsetComment,
			func(opts *AlterRoleOptions) { opts.UnsetComment = new(true) },
			`ALTER ROLE %s UNSET COMMENT`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Roles_sql_Alter_UnsetTags,
			func(opts *AlterRoleOptions) {
				opts.UnsetTags = []ObjectIdentifier{tagId, tagId2}
			},
			`ALTER ROLE %s UNSET TAG %s, %s`,
			id.FullyQualifiedName(), tagId.FullyQualifiedName(), tagId2.FullyQualifiedName(),
		)

	rolesTests.Drop.
		withExpectedSqlf(
			case_Roles_sql_Drop_basic,
			`DROP ROLE %s`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Roles_sql_Drop_all,
			func(opts *DropRoleOptions) { opts.IfExists = new(true) },
			`DROP ROLE IF EXISTS %s`, id.FullyQualifiedName(),
		)

	rolesTests.Show.
		withAdditionalValidationCase(
			"validation_Show_Like_patternRequired",
			func(opts *ShowRoleOptions) { opts.Like = &Like{} },
			ErrPatternRequiredForLikeKeyword,
		).
		withExpectedSql(case_Roles_sql_Show_basic, `SHOW ROLES`).
		withModifyAndExpectedSqlf(
			case_Roles_sql_Show_all,
			func(opts *ShowRoleOptions) {
				opts.Like = &Like{Pattern: new("new_role")}
				opts.InClass = &RolesInClass{Class: classId}
			},
			`SHOW ROLES LIKE 'new_role' IN CLASS %s`, classId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Roles_sql_Show_Like,
			func(opts *ShowRoleOptions) { opts.Like = &Like{Pattern: new("new_role")} },
			`SHOW ROLES LIKE 'new_role'`,
		).
		withModifyAndExpectedSqlf(
			case_Roles_sql_Show_InClass,
			func(opts *ShowRoleOptions) { opts.InClass = &RolesInClass{Class: classId} },
			`SHOW ROLES IN CLASS %s`, classId.FullyQualifiedName(),
		)

	rolesTests.Grant.
		withDefaultOpts(func() *GrantRoleOptions {
			return &GrantRoleOptions{
				name:  id,
				Grant: GrantRoleTo{User: &grantedUserId},
			}
		}).
		withAdditionalValidationCase(
			"validation_Grant_Grant_Role_ValidIdentifier",
			func(opts *GrantRoleOptions) {
				opts.Grant.User = nil
				opts.Grant.Role = new(emptyAccountObjectIdentifier)
			},
			errInvalidIdentifier("GrantRoleOptions.Grant", "Role"),
		).
		withAdditionalValidationCase(
			"validation_Grant_Grant_User_ValidIdentifier",
			func(opts *GrantRoleOptions) {
				opts.Grant.User = new(emptyAccountObjectIdentifier)
			},
			errInvalidIdentifier("GrantRoleOptions.Grant", "User"),
		).
		withExpectedSqlf(
			case_Roles_sql_Grant_basic,
			`GRANT ROLE %s TO USER %s`, id.FullyQualifiedName(), grantedUserId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Grant_toRole",
			func(opts *GrantRoleOptions) {
				opts.Grant.User = nil
				opts.Grant.Role = &grantedRoleId
			},
			`GRANT ROLE %s TO ROLE %s`, id.FullyQualifiedName(), grantedRoleId.FullyQualifiedName(),
		)

	rolesTests.Revoke.
		withDefaultOpts(func() *RevokeRoleOptions {
			return &RevokeRoleOptions{
				name:   id,
				Revoke: RevokeRoleFrom{User: &grantedUserId},
			}
		}).
		withExpectedSqlf(
			case_Roles_sql_Revoke_basic,
			`REVOKE ROLE %s FROM USER %s`, id.FullyQualifiedName(), grantedUserId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Revoke_fromRole",
			func(opts *RevokeRoleOptions) {
				opts.Revoke.User = nil
				opts.Revoke.Role = &grantedRoleId
			},
			`REVOKE ROLE %s FROM ROLE %s`, id.FullyQualifiedName(), grantedRoleId.FullyQualifiedName(),
		)
}
