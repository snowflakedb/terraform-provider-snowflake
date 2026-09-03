package sdk

func init() {
	id := passwordPoliciesTestIdSchemaObjectIdentifier
	renameTarget := randomSchemaObjectIdentifierInSchema(id.SchemaId())
	databaseId := NewAccountObjectIdentifier(id.DatabaseName())

	passwordPoliciesTests.Create.
		withExpectedSqlf(
			case_PasswordPolicies_sql_Create_basic,
			"CREATE PASSWORD POLICY %s", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_PasswordPolicies_sql_Create_all,
			func(opts *CreatePasswordPolicyOptions) {
				opts.IfNotExists = new(true)
				opts.PasswordMinLength = new(10)
				opts.PasswordMaxLength = new(20)
				opts.PasswordMinUpperCaseChars = new(1)
				opts.PasswordMinLowerCaseChars = new(1)
				opts.PasswordMinNumericChars = new(1)
				opts.PasswordMinSpecialChars = new(1)
				opts.PasswordMinAgeDays = new(30)
				opts.PasswordMaxAgeDays = new(30)
				opts.PasswordMaxRetries = new(5)
				opts.PasswordLockoutTimeMins = new(30)
				opts.PasswordHistory = new(15)
				opts.Comment = new("test comment")
			},
			`CREATE PASSWORD POLICY IF NOT EXISTS %s PASSWORD_MIN_LENGTH = 10 PASSWORD_MAX_LENGTH = 20 PASSWORD_MIN_UPPER_CASE_CHARS = 1 PASSWORD_MIN_LOWER_CASE_CHARS = 1 PASSWORD_MIN_NUMERIC_CHARS = 1 PASSWORD_MIN_SPECIAL_CHARS = 1 PASSWORD_MIN_AGE_DAYS = 30 PASSWORD_MAX_AGE_DAYS = 30 PASSWORD_MAX_RETRIES = 5 PASSWORD_LOCKOUT_TIME_MINS = 30 PASSWORD_HISTORY = 15 COMMENT = 'test comment'`,
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Create_orReplace",
			func(opts *CreatePasswordPolicyOptions) { opts.OrReplace = new(true) },
			"CREATE OR REPLACE PASSWORD POLICY %s", id.FullyQualifiedName(),
		)

	passwordPoliciesTests.Alter.
		withModify(
			case_PasswordPolicies_validation_Alter_opts_ExactlyOneValueSet_MoreThanOneSet,
			func(opts *AlterPasswordPolicyOptions) {
				opts.Set = &PasswordPolicySet{Comment: new("test comment")}
				opts.Unset = &PasswordPolicyUnset{Comment: new(true)}
			},
		).
		withModifyAndExpectedSqlf(
			case_PasswordPolicies_sql_Alter_Set,
			func(opts *AlterPasswordPolicyOptions) {
				opts.Set = &PasswordPolicySet{
					PasswordMinLength:         new(10),
					PasswordMaxLength:         new(20),
					PasswordMinUpperCaseChars: new(1),
					PasswordMinLowerCaseChars: new(1),
					PasswordMinNumericChars:   new(1),
					PasswordMinSpecialChars:   new(1),
					PasswordMinAgeDays:        new(30),
					PasswordMaxAgeDays:        new(30),
					PasswordMaxRetries:        new(5),
					PasswordLockoutTimeMins:   new(30),
					PasswordHistory:           new(15),
					Comment:                   new("test comment"),
				}
			},
			"ALTER PASSWORD POLICY %s SET PASSWORD_MIN_LENGTH = 10, PASSWORD_MAX_LENGTH = 20, PASSWORD_MIN_UPPER_CASE_CHARS = 1, PASSWORD_MIN_LOWER_CASE_CHARS = 1, PASSWORD_MIN_NUMERIC_CHARS = 1, PASSWORD_MIN_SPECIAL_CHARS = 1, PASSWORD_MIN_AGE_DAYS = 30, PASSWORD_MAX_AGE_DAYS = 30, PASSWORD_MAX_RETRIES = 5, PASSWORD_LOCKOUT_TIME_MINS = 30, PASSWORD_HISTORY = 15, COMMENT = 'test comment'",
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_PasswordPolicies_sql_Alter_Unset,
			func(opts *AlterPasswordPolicyOptions) {
				opts.Unset = &PasswordPolicyUnset{
					PasswordMinLength:         new(true),
					PasswordMaxLength:         new(true),
					PasswordMinUpperCaseChars: new(true),
					PasswordMinLowerCaseChars: new(true),
					PasswordMinNumericChars:   new(true),
					PasswordMinSpecialChars:   new(true),
					PasswordMinAgeDays:        new(true),
					PasswordMaxAgeDays:        new(true),
					PasswordMaxRetries:        new(true),
					PasswordLockoutTimeMins:   new(true),
					PasswordHistory:           new(true),
					Comment:                   new(true),
				}
			},
			"ALTER PASSWORD POLICY %s UNSET PASSWORD_MIN_LENGTH, PASSWORD_MAX_LENGTH, PASSWORD_MIN_UPPER_CASE_CHARS, PASSWORD_MIN_LOWER_CASE_CHARS, PASSWORD_MIN_NUMERIC_CHARS, PASSWORD_MIN_SPECIAL_CHARS, PASSWORD_MIN_AGE_DAYS, PASSWORD_MAX_AGE_DAYS, PASSWORD_MAX_RETRIES, PASSWORD_LOCKOUT_TIME_MINS, PASSWORD_HISTORY, COMMENT",
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_PasswordPolicies_sql_Alter_RenameTo,
			func(opts *AlterPasswordPolicyOptions) { opts.RenameTo = &renameTarget },
			"ALTER PASSWORD POLICY %s RENAME TO %s",
			id.FullyQualifiedName(), renameTarget.FullyQualifiedName(),
		)

	passwordPoliciesTests.Drop.
		withExpectedSqlf(
			case_PasswordPolicies_sql_Drop_basic,
			"DROP PASSWORD POLICY %s", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_PasswordPolicies_sql_Drop_all,
			func(opts *DropPasswordPolicyOptions) { opts.IfExists = new(true) },
			"DROP PASSWORD POLICY IF EXISTS %s", id.FullyQualifiedName(),
		)

	passwordPoliciesTests.Show.
		withExpectedSql(case_PasswordPolicies_sql_Show_basic, "SHOW PASSWORD POLICIES").
		withModifyAndExpectedSqlf(
			case_PasswordPolicies_sql_Show_all,
			func(opts *ShowPasswordPolicyOptions) {
				opts.Like = &Like{Pattern: new(id.Name())}
				opts.In = &ExtendedIn{In: In{Schema: id.SchemaId()}}
				opts.StartsWith = new("starts-with-pattern")
				opts.Limit = &LimitFrom{Rows: new(10), From: new("limit-from")}
			},
			"SHOW PASSWORD POLICIES LIKE '%s' IN SCHEMA %s STARTS WITH 'starts-with-pattern' LIMIT 10 FROM 'limit-from'",
			id.Name(), id.SchemaId().FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_PasswordPolicies_sql_Show_Like,
			func(opts *ShowPasswordPolicyOptions) { opts.Like = &Like{Pattern: new(id.Name())} },
			"SHOW PASSWORD POLICIES LIKE '%s'", id.Name(),
		).
		withModifyAndExpectedSqlf(
			case_PasswordPolicies_sql_Show_In,
			func(opts *ShowPasswordPolicyOptions) {
				opts.In = &ExtendedIn{In: In{Schema: id.SchemaId()}}
			},
			"SHOW PASSWORD POLICIES IN SCHEMA %s", id.SchemaId().FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_PasswordPolicies_sql_Show_On,
			func(opts *ShowPasswordPolicyOptions) { opts.On = &On{Account: new(true)} },
			"SHOW PASSWORD POLICIES ON ACCOUNT",
		).
		withModifyAndExpectedSqlf(
			case_PasswordPolicies_sql_Show_StartsWith,
			func(opts *ShowPasswordPolicyOptions) { opts.StartsWith = new("starts-with-pattern") },
			"SHOW PASSWORD POLICIES STARTS WITH 'starts-with-pattern'",
		).
		withModifyAndExpectedSqlf(
			case_PasswordPolicies_sql_Show_Limit,
			func(opts *ShowPasswordPolicyOptions) { opts.Limit = &LimitFrom{Rows: new(10)} },
			"SHOW PASSWORD POLICIES LIMIT 10",
		).
		withAdditionalSqlCasef(
			"sql_Show_Like_In_Account",
			func(opts *ShowPasswordPolicyOptions) {
				opts.Like = &Like{Pattern: new(id.Name())}
				opts.In = &ExtendedIn{In: In{Account: new(true)}}
			},
			"SHOW PASSWORD POLICIES LIKE '%s' IN ACCOUNT", id.Name(),
		).
		withAdditionalSqlCasef(
			"sql_Show_Like_In_Database",
			func(opts *ShowPasswordPolicyOptions) {
				opts.Like = &Like{Pattern: new(id.Name())}
				opts.In = &ExtendedIn{In: In{Database: databaseId}}
			},
			"SHOW PASSWORD POLICIES LIKE '%s' IN DATABASE %s",
			id.Name(), databaseId.FullyQualifiedName(),
		)

	passwordPoliciesTests.Describe.
		withExpectedSqlf(
			case_PasswordPolicies_sql_Describe_basic,
			"DESCRIBE PASSWORD POLICY %s", id.FullyQualifiedName(),
		)
}
