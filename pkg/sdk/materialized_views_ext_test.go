package sdk

func init() {
	id := materializedViewsTestIdSchemaObjectIdentifier
	sql := "SELECT id FROM t"

	rowAccessPolicyId := randomSchemaObjectIdentifier()
	tag1Id := randomSchemaObjectIdentifier()
	tag2Id := randomSchemaObjectIdentifier()
	maskingPolicy1Id := randomSchemaObjectIdentifier()
	maskingPolicy2Id := randomSchemaObjectIdentifier()
	renameTarget := randomSchemaObjectIdentifier()

	materializedViewsTests.Create.
		withDefaultOpts(func() *CreateMaterializedViewOptions {
			return &CreateMaterializedViewOptions{
				name: id,
				sql:  sql,
			}
		}).
		withExpectedSqlf(
			case_MaterializedViews_sql_Create_basic,
			"CREATE MATERIALIZED VIEW %s AS %s", id.FullyQualifiedName(), sql,
		).
		withModifyAndExpectedSqlf(
			case_MaterializedViews_sql_Create_all,
			func(opts *CreateMaterializedViewOptions) {
				opts.IfNotExists = new(true)
				opts.Secure = new(true)
				opts.Columns = []MaterializedViewColumn{
					{Name: "column_without_comment"},
					{Name: "column_with_comment", Comment: new("column 2 comment")},
				}
				opts.ColumnsMaskingPolicies = []MaterializedViewColumnMaskingPolicy{
					{
						Name:          "column",
						MaskingPolicy: maskingPolicy1Id,
						Using:         []string{"a", "b"},
						Tag:           []TagAssociation{{Name: tag1Id, Value: "v1"}},
					},
					{Name: "column 2", MaskingPolicy: maskingPolicy2Id},
				}
				opts.Comment = new("comment")
				opts.RowAccessPolicy = &MaterializedViewRowAccessPolicy{RowAccessPolicy: rowAccessPolicyId, On: []string{"c", "d"}}
				opts.Tag = []TagAssociation{{Name: tag2Id, Value: "v2"}}
				opts.ClusterBy = &MaterializedViewClusterBy{Expressions: []MaterializedViewClusterByExpression{{"column_without_comment"}, {"column_with_comment"}}}
			},
			`CREATE SECURE MATERIALIZED VIEW IF NOT EXISTS %s ("column_without_comment", "column_with_comment" COMMENT 'column 2 comment') column MASKING POLICY %s USING (a, b) TAG (%s = 'v1'), column 2 MASKING POLICY %s COMMENT = 'comment' ROW ACCESS POLICY %s ON (c, d) TAG (%s = 'v2') CLUSTER BY ("column_without_comment", "column_with_comment") AS %s`,
			id.FullyQualifiedName(), maskingPolicy1Id.FullyQualifiedName(), tag1Id.FullyQualifiedName(), maskingPolicy2Id.FullyQualifiedName(), rowAccessPolicyId.FullyQualifiedName(), tag2Id.FullyQualifiedName(), sql,
		).
		withAdditionalSqlCasef(
			"sql_Create_orReplace",
			func(opts *CreateMaterializedViewOptions) {
				opts.OrReplace = new(true)
				opts.CopyGrants = new(true)
			},
			"CREATE OR REPLACE MATERIALIZED VIEW %s COPY GRANTS AS %s", id.FullyQualifiedName(), sql,
		)

	materializedViewsTests.Alter.
		withModifyAndExpectedSqlf(
			case_MaterializedViews_sql_Alter_RenameTo,
			func(opts *AlterMaterializedViewOptions) { opts.RenameTo = &renameTarget },
			"ALTER MATERIALIZED VIEW %s RENAME TO %s", id.FullyQualifiedName(), renameTarget.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_MaterializedViews_sql_Alter_ClusterBy,
			func(opts *AlterMaterializedViewOptions) {
				opts.ClusterBy = &MaterializedViewClusterBy{Expressions: []MaterializedViewClusterByExpression{{"column"}}}
			},
			`ALTER MATERIALIZED VIEW %s CLUSTER BY ("column")`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_MaterializedViews_sql_Alter_DropClusteringKey,
			func(opts *AlterMaterializedViewOptions) { opts.DropClusteringKey = new(true) },
			"ALTER MATERIALIZED VIEW %s DROP CLUSTERING KEY", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_MaterializedViews_sql_Alter_SuspendRecluster,
			func(opts *AlterMaterializedViewOptions) { opts.SuspendRecluster = new(true) },
			"ALTER MATERIALIZED VIEW %s SUSPEND RECLUSTER", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_MaterializedViews_sql_Alter_ResumeRecluster,
			func(opts *AlterMaterializedViewOptions) { opts.ResumeRecluster = new(true) },
			"ALTER MATERIALIZED VIEW %s RESUME RECLUSTER", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_MaterializedViews_sql_Alter_Suspend,
			func(opts *AlterMaterializedViewOptions) { opts.Suspend = new(true) },
			"ALTER MATERIALIZED VIEW %s SUSPEND", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_MaterializedViews_sql_Alter_Resume,
			func(opts *AlterMaterializedViewOptions) { opts.Resume = new(true) },
			"ALTER MATERIALIZED VIEW %s RESUME", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_MaterializedViews_sql_Alter_Set,
			func(opts *AlterMaterializedViewOptions) { opts.Set = &MaterializedViewSet{Secure: new(true)} },
			"ALTER MATERIALIZED VIEW %s SET SECURE", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_MaterializedViews_sql_Alter_Unset,
			func(opts *AlterMaterializedViewOptions) { opts.Unset = &MaterializedViewUnset{Secure: new(true)} },
			"ALTER MATERIALIZED VIEW %s UNSET SECURE", id.FullyQualifiedName(),
		)

	materializedViewsTests.Drop.
		withExpectedSqlf(
			case_MaterializedViews_sql_Drop_basic,
			"DROP MATERIALIZED VIEW %s", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_MaterializedViews_sql_Drop_all,
			func(opts *DropMaterializedViewOptions) { opts.IfExists = new(true) },
			"DROP MATERIALIZED VIEW IF EXISTS %s", id.FullyQualifiedName(),
		)

	materializedViewsTests.Show.
		withExpectedSql(case_MaterializedViews_sql_Show_basic, "SHOW MATERIALIZED VIEWS").
		withModifyAndExpectedSqlf(
			case_MaterializedViews_sql_Show_all,
			func(opts *ShowMaterializedViewOptions) {
				opts.Like = &Like{Pattern: new("myaccount")}
				opts.In = &In{Account: new(true)}
			},
			"SHOW MATERIALIZED VIEWS LIKE 'myaccount' IN ACCOUNT",
		).
		withModifyAndExpectedSqlf(
			case_MaterializedViews_sql_Show_Like,
			func(opts *ShowMaterializedViewOptions) { opts.Like = &Like{Pattern: new("pattern")} },
			"SHOW MATERIALIZED VIEWS LIKE 'pattern'",
		).
		withModifyAndExpectedSqlf(
			case_MaterializedViews_sql_Show_In,
			func(opts *ShowMaterializedViewOptions) { opts.In = &In{Account: new(true)} },
			"SHOW MATERIALIZED VIEWS IN ACCOUNT",
		)

	materializedViewsTests.Describe.
		withExpectedSqlf(
			case_MaterializedViews_sql_Describe_basic,
			"DESCRIBE MATERIALIZED VIEW %s", id.FullyQualifiedName(),
		)
}
