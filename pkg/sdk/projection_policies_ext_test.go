package sdk

func init() {
	id := projectionPoliciesTestIdSchemaObjectIdentifier
	renameTarget := randomSchemaObjectIdentifier()
	body := "PROJECTION_CONSTRAINT(ALLOW => false)"
	comment := "some comment"
	tag1 := NewAccountObjectIdentifier("tag1")
	tag2 := NewAccountObjectIdentifier("tag2")

	projectionPoliciesTests.Create.
		withDefaultOpts(func() *CreateProjectionPolicyOptions {
			return &CreateProjectionPolicyOptions{
				name: id,
				Body: body,
			}
		}).
		withExpectedSqlf(
			case_ProjectionPolicies_sql_Create_basic,
			"CREATE PROJECTION POLICY %s AS () RETURNS PROJECTION_CONSTRAINT -> %s",
			id.FullyQualifiedName(), body,
		).
		withModifyAndExpectedSqlf(
			case_ProjectionPolicies_sql_Create_all,
			func(opts *CreateProjectionPolicyOptions) {
				opts.IfNotExists = new(true)
				opts.Comment = new(comment)
				opts.Tag = []TagAssociation{
					{Name: tag1, Value: "value1"},
					{Name: tag2, Value: "value2"},
				}
			},
			`CREATE PROJECTION POLICY IF NOT EXISTS %s AS () RETURNS PROJECTION_CONSTRAINT -> %s COMMENT = '%s' TAG ("tag1" = 'value1', "tag2" = 'value2')`,
			id.FullyQualifiedName(), body, comment,
		).
		withAdditionalSqlCasef(
			"sql_Create_orReplace",
			func(opts *CreateProjectionPolicyOptions) { opts.OrReplace = new(true) },
			"CREATE OR REPLACE PROJECTION POLICY %s AS () RETURNS PROJECTION_CONSTRAINT -> %s",
			id.FullyQualifiedName(), body,
		)

	projectionPoliciesTests.Alter.
		withModify(
			case_ProjectionPolicies_validation_Alter_opts_ExactlyOneValueSet_MoreThanOneSet,
			func(opts *AlterProjectionPolicyOptions) {
				opts.RenameTo = &renameTarget
				opts.Set = &ProjectionPolicySet{SetBody: new(body)}
			},
		).
		withModifyAndExpectedSqlf(
			case_ProjectionPolicies_sql_Alter_RenameTo,
			func(opts *AlterProjectionPolicyOptions) { opts.RenameTo = &renameTarget },
			"ALTER PROJECTION POLICY %s RENAME TO %s",
			id.FullyQualifiedName(), renameTarget.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_ProjectionPolicies_sql_Alter_SetTags,
			func(opts *AlterProjectionPolicyOptions) {
				opts.SetTags = []TagAssociation{
					{Name: tag1, Value: "value1"},
					{Name: tag2, Value: "value2"},
				}
			},
			`ALTER PROJECTION POLICY %s SET TAG "tag1" = 'value1', "tag2" = 'value2'`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_ProjectionPolicies_sql_Alter_UnsetTags,
			func(opts *AlterProjectionPolicyOptions) {
				opts.UnsetTags = []ObjectIdentifier{tag1, tag2}
			},
			`ALTER PROJECTION POLICY %s UNSET TAG "tag1", "tag2"`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_ProjectionPolicies_sql_Alter_Set,
			func(opts *AlterProjectionPolicyOptions) {
				opts.Set = &ProjectionPolicySet{Comment: new(comment)}
			},
			"ALTER PROJECTION POLICY %s SET COMMENT = '%s'",
			id.FullyQualifiedName(), comment,
		).
		withModifyAndExpectedSqlf(
			case_ProjectionPolicies_sql_Alter_Unset,
			func(opts *AlterProjectionPolicyOptions) {
				opts.Unset = &ProjectionPolicyUnset{Comment: new(true)}
			},
			"ALTER PROJECTION POLICY %s UNSET COMMENT",
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_setBody",
			func(opts *AlterProjectionPolicyOptions) {
				opts.Set = &ProjectionPolicySet{SetBody: new(body)}
			},
			"ALTER PROJECTION POLICY %s SET BODY -> %s",
			id.FullyQualifiedName(), body,
		).
		withAdditionalSqlCasef(
			"sql_Alter_ifExists",
			func(opts *AlterProjectionPolicyOptions) {
				opts.IfExists = new(true)
				opts.Set = &ProjectionPolicySet{SetBody: new(body)}
			},
			"ALTER PROJECTION POLICY IF EXISTS %s SET BODY -> %s",
			id.FullyQualifiedName(), body,
		)

	projectionPoliciesTests.Drop.
		withExpectedSqlf(
			case_ProjectionPolicies_sql_Drop_basic,
			"DROP PROJECTION POLICY %s",
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_ProjectionPolicies_sql_Drop_all,
			func(opts *DropProjectionPolicyOptions) { opts.IfExists = new(true) },
			"DROP PROJECTION POLICY IF EXISTS %s",
			id.FullyQualifiedName(),
		)

	projectionPoliciesTests.Show.
		withExpectedSql(
			case_ProjectionPolicies_sql_Show_basic,
			"SHOW PROJECTION POLICIES",
		).
		withModifyAndExpectedSqlf(
			case_ProjectionPolicies_sql_Show_all,
			func(opts *ShowProjectionPolicyOptions) {
				opts.Like = &Like{Pattern: String("like-pattern")}
				opts.In = &ExtendedIn{In: In{Account: new(true)}}
				opts.StartsWith = String("abc")
				opts.Limit = &LimitFrom{Rows: Int(10)}
			},
			"SHOW PROJECTION POLICIES LIKE 'like-pattern' IN ACCOUNT STARTS WITH 'abc' LIMIT 10",
		).
		withModifyAndExpectedSqlf(
			case_ProjectionPolicies_sql_Show_Like,
			func(opts *ShowProjectionPolicyOptions) { opts.Like = &Like{Pattern: String("like-pattern")} },
			"SHOW PROJECTION POLICIES LIKE 'like-pattern'",
		).
		withModifyAndExpectedSqlf(
			case_ProjectionPolicies_sql_Show_In,
			func(opts *ShowProjectionPolicyOptions) { opts.In = &ExtendedIn{In: In{Account: new(true)}} },
			"SHOW PROJECTION POLICIES IN ACCOUNT",
		).
		withModifyAndExpectedSqlf(
			case_ProjectionPolicies_sql_Show_StartsWith,
			func(opts *ShowProjectionPolicyOptions) { opts.StartsWith = String("abc") },
			"SHOW PROJECTION POLICIES STARTS WITH 'abc'",
		).
		withModifyAndExpectedSqlf(
			case_ProjectionPolicies_sql_Show_Limit,
			func(opts *ShowProjectionPolicyOptions) { opts.Limit = &LimitFrom{Rows: Int(10)} },
			"SHOW PROJECTION POLICIES LIMIT 10",
		).
		withAdditionalSqlCasef(
			"sql_Show_In_Database",
			func(opts *ShowProjectionPolicyOptions) {
				opts.In = &ExtendedIn{In: In{Database: id.DatabaseId()}}
			},
			"SHOW PROJECTION POLICIES IN DATABASE %s",
			id.DatabaseId().FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Show_In_Schema",
			func(opts *ShowProjectionPolicyOptions) {
				opts.In = &ExtendedIn{In: In{Schema: id.SchemaId()}}
			},
			"SHOW PROJECTION POLICIES IN SCHEMA %s",
			id.SchemaId().FullyQualifiedName(),
		)

	projectionPoliciesTests.Describe.
		withExpectedSqlf(
			case_ProjectionPolicies_sql_Describe_basic,
			"DESCRIBE PROJECTION POLICY %s",
			id.FullyQualifiedName(),
		)
}
