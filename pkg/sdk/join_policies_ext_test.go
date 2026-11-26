package sdk

func init() {
	id := joinPoliciesTestIdSchemaObjectIdentifier
	renameTarget := randomSchemaObjectIdentifier()
	body := "JOIN_CONSTRAINT(JOIN_REQUIRED => TRUE)"
	comment := "some comment"
	tag1 := NewAccountObjectIdentifier("tag1")
	tag2 := NewAccountObjectIdentifier("tag2")

	joinPoliciesTests.Create.
		withDefaultOpts(func() *CreateJoinPolicyOptions {
			return &CreateJoinPolicyOptions{
				name: id,
				Body: body,
			}
		}).
		withExpectedSqlf(
			case_JoinPolicies_sql_Create_basic,
			"CREATE JOIN POLICY %s AS () RETURNS JOIN_CONSTRAINT -> %s",
			id.FullyQualifiedName(), body,
		).
		withModifyAndExpectedSqlf(
			case_JoinPolicies_sql_Create_all,
			func(opts *CreateJoinPolicyOptions) {
				opts.IfNotExists = new(true)
				opts.Comment = new(comment)
				opts.Tag = []TagAssociation{
					{Name: tag1, Value: "value1"},
					{Name: tag2, Value: "value2"},
				}
			},
			`CREATE JOIN POLICY IF NOT EXISTS %s AS () RETURNS JOIN_CONSTRAINT -> %s COMMENT = '%s' TAG ("tag1" = 'value1', "tag2" = 'value2')`,
			id.FullyQualifiedName(), body, comment,
		).
		withAdditionalSqlCasef(
			"sql_Create_orReplace",
			func(opts *CreateJoinPolicyOptions) { opts.OrReplace = new(true) },
			"CREATE OR REPLACE JOIN POLICY %s AS () RETURNS JOIN_CONSTRAINT -> %s",
			id.FullyQualifiedName(), body,
		)

	joinPoliciesTests.Alter.
		withModify(
			case_JoinPolicies_validation_Alter_opts_ExactlyOneValueSet_MoreThanOneSet,
			func(opts *AlterJoinPolicyOptions) {
				opts.RenameTo = &renameTarget
				opts.Set = &JoinPolicySet{SetBody: new(body)}
			},
		).
		withModifyAndExpectedSqlf(
			case_JoinPolicies_sql_Alter_RenameTo,
			func(opts *AlterJoinPolicyOptions) { opts.RenameTo = &renameTarget },
			"ALTER JOIN POLICY %s RENAME TO %s",
			id.FullyQualifiedName(), renameTarget.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_JoinPolicies_sql_Alter_SetTags,
			func(opts *AlterJoinPolicyOptions) {
				opts.SetTags = []TagAssociation{
					{Name: tag1, Value: "value1"},
					{Name: tag2, Value: "value2"},
				}
			},
			`ALTER JOIN POLICY %s SET TAG "tag1" = 'value1', "tag2" = 'value2'`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_JoinPolicies_sql_Alter_UnsetTags,
			func(opts *AlterJoinPolicyOptions) {
				opts.UnsetTags = []ObjectIdentifier{tag1, tag2}
			},
			`ALTER JOIN POLICY %s UNSET TAG "tag1", "tag2"`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_JoinPolicies_sql_Alter_Set,
			func(opts *AlterJoinPolicyOptions) {
				opts.Set = &JoinPolicySet{Comment: new(comment)}
			},
			"ALTER JOIN POLICY %s SET COMMENT = '%s'",
			id.FullyQualifiedName(), comment,
		).
		withModifyAndExpectedSqlf(
			case_JoinPolicies_sql_Alter_Unset,
			func(opts *AlterJoinPolicyOptions) {
				opts.Unset = &JoinPolicyUnset{Comment: new(true)}
			},
			"ALTER JOIN POLICY %s UNSET COMMENT",
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_setBody",
			func(opts *AlterJoinPolicyOptions) {
				opts.Set = &JoinPolicySet{SetBody: new(body)}
			},
			"ALTER JOIN POLICY %s SET BODY -> %s",
			id.FullyQualifiedName(), body,
		).
		withAdditionalSqlCasef(
			"sql_Alter_ifExists",
			func(opts *AlterJoinPolicyOptions) {
				opts.IfExists = new(true)
				opts.Set = &JoinPolicySet{SetBody: new(body)}
			},
			"ALTER JOIN POLICY IF EXISTS %s SET BODY -> %s",
			id.FullyQualifiedName(), body,
		)

	joinPoliciesTests.Drop.
		withExpectedSqlf(
			case_JoinPolicies_sql_Drop_basic,
			"DROP JOIN POLICY %s",
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_JoinPolicies_sql_Drop_all,
			func(opts *DropJoinPolicyOptions) { opts.IfExists = new(true) },
			"DROP JOIN POLICY IF EXISTS %s",
			id.FullyQualifiedName(),
		)

	joinPoliciesTests.Show.
		withExpectedSql(
			case_JoinPolicies_sql_Show_basic,
			"SHOW JOIN POLICIES",
		).
		withModifyAndExpectedSqlf(
			case_JoinPolicies_sql_Show_all,
			func(opts *ShowJoinPolicyOptions) {
				opts.Like = &Like{Pattern: String("like-pattern")}
				opts.In = &In{Account: new(true)}
				opts.StartsWith = String("abc")
				opts.Limit = &LimitFrom{Rows: Int(10)}
			},
			"SHOW JOIN POLICIES LIKE 'like-pattern' IN ACCOUNT STARTS WITH 'abc' LIMIT 10",
		).
		withModifyAndExpectedSqlf(
			case_JoinPolicies_sql_Show_Like,
			func(opts *ShowJoinPolicyOptions) { opts.Like = &Like{Pattern: String("like-pattern")} },
			"SHOW JOIN POLICIES LIKE 'like-pattern'",
		).
		withModifyAndExpectedSqlf(
			case_JoinPolicies_sql_Show_In,
			func(opts *ShowJoinPolicyOptions) { opts.In = &In{Account: new(true)} },
			"SHOW JOIN POLICIES IN ACCOUNT",
		).
		withModifyAndExpectedSqlf(
			case_JoinPolicies_sql_Show_StartsWith,
			func(opts *ShowJoinPolicyOptions) { opts.StartsWith = String("abc") },
			"SHOW JOIN POLICIES STARTS WITH 'abc'",
		).
		withModifyAndExpectedSqlf(
			case_JoinPolicies_sql_Show_Limit,
			func(opts *ShowJoinPolicyOptions) { opts.Limit = &LimitFrom{Rows: Int(10)} },
			"SHOW JOIN POLICIES LIMIT 10",
		).
		withAdditionalSqlCasef(
			"sql_Show_In_Database",
			func(opts *ShowJoinPolicyOptions) {
				opts.In = &In{Database: id.DatabaseId()}
			},
			"SHOW JOIN POLICIES IN DATABASE %s",
			id.DatabaseId().FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Show_In_Schema",
			func(opts *ShowJoinPolicyOptions) {
				opts.In = &In{Schema: id.SchemaId()}
			},
			"SHOW JOIN POLICIES IN SCHEMA %s",
			id.SchemaId().FullyQualifiedName(),
		)

	joinPoliciesTests.Describe.
		withExpectedSqlf(
			case_JoinPolicies_sql_Describe_basic,
			"DESCRIBE JOIN POLICY %s",
			id.FullyQualifiedName(),
		)
}
