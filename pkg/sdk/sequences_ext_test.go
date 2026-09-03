package sdk

func init() {
	id := sequencesTestIdSchemaObjectIdentifier
	renameTarget := randomSchemaObjectIdentifierInSchema(id.SchemaId())

	sequencesTests.Create.
		withExpectedSqlf(
			case_Sequences_sql_Create_basic,
			`CREATE SEQUENCE %s`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Sequences_sql_Create_all,
			func(opts *CreateSequenceOptions) {
				opts.IfNotExists = new(true)
				opts.Start = new(1)
				opts.Increment = new(1)
				opts.ValuesBehavior = new(ValuesBehaviorOrder)
				opts.Comment = new("comment")
			},
			`CREATE SEQUENCE IF NOT EXISTS %s START = 1 INCREMENT = 1 ORDER COMMENT = 'comment'`,
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Create_orReplace",
			func(opts *CreateSequenceOptions) { opts.OrReplace = new(true) },
			`CREATE OR REPLACE SEQUENCE %s`,
			id.FullyQualifiedName(),
		)

	sequencesTests.Alter.
		withDefaultOpts(func() *AlterSequenceOptions {
			return &AlterSequenceOptions{
				name:     id,
				IfExists: new(true),
			}
		}).
		withModifyAndExpectedSqlf(
			case_Sequences_sql_Alter_RenameTo,
			func(opts *AlterSequenceOptions) { opts.RenameTo = &renameTarget },
			`ALTER SEQUENCE IF EXISTS %s RENAME TO %s`,
			id.FullyQualifiedName(), renameTarget.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Sequences_sql_Alter_SetIncrement,
			func(opts *AlterSequenceOptions) { opts.SetIncrement = new(1) },
			`ALTER SEQUENCE IF EXISTS %s SET INCREMENT = 1`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Sequences_sql_Alter_Set,
			func(opts *AlterSequenceOptions) {
				opts.Set = &SequenceSet{
					Comment:        new("comment"),
					ValuesBehavior: new(ValuesBehaviorOrder),
				}
			},
			`ALTER SEQUENCE IF EXISTS %s SET ORDER COMMENT = 'comment'`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Sequences_sql_Alter_UnsetComment,
			func(opts *AlterSequenceOptions) { opts.UnsetComment = new(true) },
			`ALTER SEQUENCE IF EXISTS %s UNSET COMMENT`, id.FullyQualifiedName(),
		)

	sequencesTests.Show.
		withExpectedSql(case_Sequences_sql_Show_basic, `SHOW SEQUENCES`).
		withModifyAndExpectedSqlf(
			case_Sequences_sql_Show_all,
			func(opts *ShowSequenceOptions) {
				opts.Like = &Like{Pattern: new("pattern")}
				opts.In = &In{Account: new(true)}
			},
			`SHOW SEQUENCES LIKE 'pattern' IN ACCOUNT`,
		).
		withModifyAndExpectedSqlf(
			case_Sequences_sql_Show_Like,
			func(opts *ShowSequenceOptions) { opts.Like = &Like{Pattern: new("pattern")} },
			`SHOW SEQUENCES LIKE 'pattern'`,
		).
		withModifyAndExpectedSqlf(
			case_Sequences_sql_Show_In,
			func(opts *ShowSequenceOptions) { opts.In = &In{Account: new(true)} },
			`SHOW SEQUENCES IN ACCOUNT`,
		)

	sequencesTests.Describe.
		withExpectedSqlf(
			case_Sequences_sql_Describe_basic,
			`DESCRIBE SEQUENCE %s`, id.FullyQualifiedName(),
		)

	sequencesTests.Drop.
		withExpectedSqlf(
			case_Sequences_sql_Drop_basic,
			`DROP SEQUENCE %s`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Sequences_sql_Drop_all,
			func(opts *DropSequenceOptions) {
				opts.IfExists = new(true)
				opts.Constraint = &SequenceConstraint{Cascade: new(true)}
			},
			`DROP SEQUENCE IF EXISTS %s CASCADE`, id.FullyQualifiedName(),
		)
}
