package sdk

func init() {
	id := applicationRolesTestIdDatabaseObjectIdentifier
	roleId := randomAccountObjectIdentifier()
	appRoleId := randomDatabaseObjectIdentifier()
	appId := randomAccountObjectIdentifier()

	applicationRolesTests.Grant.
		withDefaultOpts(func() *GrantApplicationRoleOptions {
			return &GrantApplicationRoleOptions{
				name: id,
				To:   KindOfRole{RoleName: &roleId},
			}
		}).
		withExpectedSqlf(
			case_ApplicationRoles_sql_Grant_basic,
			`GRANT APPLICATION ROLE %s TO ROLE %s`,
			id.FullyQualifiedName(), roleId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Grant_toApplicationRole",
			func(opts *GrantApplicationRoleOptions) {
				opts.To = KindOfRole{ApplicationRoleName: &appRoleId}
			},
			`GRANT APPLICATION ROLE %s TO APPLICATION ROLE %s`,
			id.FullyQualifiedName(), appRoleId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Grant_toApplication",
			func(opts *GrantApplicationRoleOptions) {
				opts.To = KindOfRole{ApplicationName: &appId}
			},
			`GRANT APPLICATION ROLE %s TO APPLICATION %s`,
			id.FullyQualifiedName(), appId.FullyQualifiedName(),
		)

	applicationRolesTests.Revoke.
		withDefaultOpts(func() *RevokeApplicationRoleOptions {
			return &RevokeApplicationRoleOptions{
				name: id,
				From: KindOfRole{RoleName: &roleId},
			}
		}).
		withExpectedSqlf(
			case_ApplicationRoles_sql_Revoke_basic,
			`REVOKE APPLICATION ROLE %s FROM ROLE %s`,
			id.FullyQualifiedName(), roleId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Revoke_fromApplicationRole",
			func(opts *RevokeApplicationRoleOptions) {
				opts.From = KindOfRole{ApplicationRoleName: &appRoleId}
			},
			`REVOKE APPLICATION ROLE %s FROM APPLICATION ROLE %s`,
			id.FullyQualifiedName(), appRoleId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Revoke_fromApplication",
			func(opts *RevokeApplicationRoleOptions) {
				opts.From = KindOfRole{ApplicationName: &appId}
			},
			`REVOKE APPLICATION ROLE %s FROM APPLICATION %s`,
			id.FullyQualifiedName(), appId.FullyQualifiedName(),
		)

	applicationRolesTests.Show.
		withDefaultOpts(func() *ShowApplicationRoleOptions {
			return &ShowApplicationRoleOptions{
				ApplicationName: appId,
			}
		}).
		withExpectedSqlf(
			case_ApplicationRoles_sql_Show_basic,
			`SHOW APPLICATION ROLES IN APPLICATION %s`,
			appId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_ApplicationRoles_sql_Show_all,
			func(opts *ShowApplicationRoleOptions) {
				opts.Limit = &LimitFrom{Rows: new(123), From: new("some limit")}
			},
			`SHOW APPLICATION ROLES IN APPLICATION %s LIMIT 123 FROM 'some limit'`,
			appId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_ApplicationRoles_sql_Show_Limit,
			func(opts *ShowApplicationRoleOptions) {
				opts.Limit = &LimitFrom{Rows: new(123), From: new("some limit")}
			},
			`SHOW APPLICATION ROLES IN APPLICATION %s LIMIT 123 FROM 'some limit'`,
			appId.FullyQualifiedName(),
		)
}
