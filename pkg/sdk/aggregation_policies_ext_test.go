package sdk

func init() {
	id := aggregationPoliciesTestIdSchemaObjectIdentifier
	renameTarget := randomSchemaObjectIdentifier()
	body := "AGGREGATION_CONSTRAINT(MIN_GROUP_SIZE => 5)"
	comment := "some comment"
	tag1 := NewAccountObjectIdentifier("tag1")
	tag2 := NewAccountObjectIdentifier("tag2")

	aggregationPoliciesTests.Create.
		withDefaultOpts(func() *CreateAggregationPolicyOptions {
			return &CreateAggregationPolicyOptions{
				name: id,
				Body: body,
			}
		}).
		withExpectedSqlf(
			case_AggregationPolicies_sql_Create_basic,
			"CREATE AGGREGATION POLICY %s AS () RETURNS AGGREGATION_CONSTRAINT -> %s",
			id.FullyQualifiedName(), body,
		).
		withModifyAndExpectedSqlf(
			case_AggregationPolicies_sql_Create_all,
			func(opts *CreateAggregationPolicyOptions) {
				opts.IfNotExists = new(true)
				opts.Comment = new(comment)
				opts.Tag = []TagAssociation{
					{Name: tag1, Value: "value1"},
					{Name: tag2, Value: "value2"},
				}
			},
			`CREATE AGGREGATION POLICY IF NOT EXISTS %s AS () RETURNS AGGREGATION_CONSTRAINT -> %s COMMENT = '%s' TAG ("tag1" = 'value1', "tag2" = 'value2')`,
			id.FullyQualifiedName(), body, comment,
		).
		withAdditionalSqlCasef(
			"sql_Create_orReplace",
			func(opts *CreateAggregationPolicyOptions) { opts.OrReplace = new(true) },
			"CREATE OR REPLACE AGGREGATION POLICY %s AS () RETURNS AGGREGATION_CONSTRAINT -> %s",
			id.FullyQualifiedName(), body,
		)

	aggregationPoliciesTests.Alter.
		withModify(
			case_AggregationPolicies_validation_Alter_opts_ExactlyOneValueSet_MoreThanOneSet,
			func(opts *AlterAggregationPolicyOptions) {
				opts.RenameTo = &renameTarget
				opts.Set = &AggregationPolicySet{SetBody: new(body)}
			},
		).
		withModifyAndExpectedSqlf(
			case_AggregationPolicies_sql_Alter_RenameTo,
			func(opts *AlterAggregationPolicyOptions) { opts.RenameTo = &renameTarget },
			"ALTER AGGREGATION POLICY %s RENAME TO %s",
			id.FullyQualifiedName(), renameTarget.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_AggregationPolicies_sql_Alter_SetTags,
			func(opts *AlterAggregationPolicyOptions) {
				opts.SetTags = []TagAssociation{
					{Name: tag1, Value: "value1"},
					{Name: tag2, Value: "value2"},
				}
			},
			`ALTER AGGREGATION POLICY %s SET TAG "tag1" = 'value1', "tag2" = 'value2'`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_AggregationPolicies_sql_Alter_UnsetTags,
			func(opts *AlterAggregationPolicyOptions) {
				opts.UnsetTags = []ObjectIdentifier{tag1, tag2}
			},
			`ALTER AGGREGATION POLICY %s UNSET TAG "tag1", "tag2"`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_AggregationPolicies_sql_Alter_Set,
			func(opts *AlterAggregationPolicyOptions) {
				opts.Set = &AggregationPolicySet{Comment: new(comment)}
			},
			"ALTER AGGREGATION POLICY %s SET COMMENT = '%s'",
			id.FullyQualifiedName(), comment,
		).
		withModifyAndExpectedSqlf(
			case_AggregationPolicies_sql_Alter_Unset,
			func(opts *AlterAggregationPolicyOptions) {
				opts.Unset = &AggregationPolicyUnset{Comment: new(true)}
			},
			"ALTER AGGREGATION POLICY %s UNSET COMMENT",
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_setBody",
			func(opts *AlterAggregationPolicyOptions) {
				opts.Set = &AggregationPolicySet{SetBody: new(body)}
			},
			"ALTER AGGREGATION POLICY %s SET BODY -> %s",
			id.FullyQualifiedName(), body,
		).
		withAdditionalSqlCasef(
			"sql_Alter_ifExists",
			func(opts *AlterAggregationPolicyOptions) {
				opts.IfExists = new(true)
				opts.Set = &AggregationPolicySet{SetBody: new(body)}
			},
			"ALTER AGGREGATION POLICY IF EXISTS %s SET BODY -> %s",
			id.FullyQualifiedName(), body,
		)

	aggregationPoliciesTests.Drop.
		withExpectedSqlf(
			case_AggregationPolicies_sql_Drop_basic,
			"DROP AGGREGATION POLICY %s",
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_AggregationPolicies_sql_Drop_all,
			func(opts *DropAggregationPolicyOptions) { opts.IfExists = new(true) },
			"DROP AGGREGATION POLICY IF EXISTS %s",
			id.FullyQualifiedName(),
		)

	aggregationPoliciesTests.Show.
		withExpectedSql(
			case_AggregationPolicies_sql_Show_basic,
			"SHOW AGGREGATION POLICIES",
		).
		withModifyAndExpectedSqlf(
			case_AggregationPolicies_sql_Show_all,
			func(opts *ShowAggregationPolicyOptions) {
				opts.Like = &Like{Pattern: String("like-pattern")}
				opts.In = &In{Account: new(true)}
				opts.StartsWith = String("abc")
				opts.Limit = &LimitFrom{Rows: Int(10)}
			},
			"SHOW AGGREGATION POLICIES LIKE 'like-pattern' IN ACCOUNT STARTS WITH 'abc' LIMIT 10",
		).
		withModifyAndExpectedSqlf(
			case_AggregationPolicies_sql_Show_Like,
			func(opts *ShowAggregationPolicyOptions) { opts.Like = &Like{Pattern: String("like-pattern")} },
			"SHOW AGGREGATION POLICIES LIKE 'like-pattern'",
		).
		withModifyAndExpectedSqlf(
			case_AggregationPolicies_sql_Show_In,
			func(opts *ShowAggregationPolicyOptions) { opts.In = &In{Account: new(true)} },
			"SHOW AGGREGATION POLICIES IN ACCOUNT",
		).
		withModifyAndExpectedSqlf(
			case_AggregationPolicies_sql_Show_StartsWith,
			func(opts *ShowAggregationPolicyOptions) { opts.StartsWith = String("abc") },
			"SHOW AGGREGATION POLICIES STARTS WITH 'abc'",
		).
		withModifyAndExpectedSqlf(
			case_AggregationPolicies_sql_Show_Limit,
			func(opts *ShowAggregationPolicyOptions) { opts.Limit = &LimitFrom{Rows: Int(10)} },
			"SHOW AGGREGATION POLICIES LIMIT 10",
		).
		withAdditionalSqlCasef(
			"sql_Show_In_Database",
			func(opts *ShowAggregationPolicyOptions) {
				opts.In = &In{Database: id.DatabaseId()}
			},
			"SHOW AGGREGATION POLICIES IN DATABASE %s",
			id.DatabaseId().FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Show_In_Schema",
			func(opts *ShowAggregationPolicyOptions) {
				opts.In = &In{Schema: id.SchemaId()}
			},
			"SHOW AGGREGATION POLICIES IN SCHEMA %s",
			id.SchemaId().FullyQualifiedName(),
		)

	aggregationPoliciesTests.Describe.
		withExpectedSqlf(
			case_AggregationPolicies_sql_Describe_basic,
			"DESCRIBE AGGREGATION POLICY %s",
			id.FullyQualifiedName(),
		)
}
