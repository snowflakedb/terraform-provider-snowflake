package sdk

func init() {
	id := dataMovementRulesTestIdSchemaObjectIdentifier
	renameTarget := randomSchemaObjectIdentifier()
	maxRows := "(1000)"
	newMaxRows := "(500)"
	comment := "some comment"

	dataMovementRulesTests.Create.
		withDefaultOpts(func() *CreateDataMovementRuleOptions {
			return &CreateDataMovementRuleOptions{
				name:    id,
				Type:    DataMovementRuleTypeCopyIntoExternalStage,
				MaxRows: maxRows,
			}
		}).
		withExpectedSqlf(
			case_DataMovementRules_sql_Create_basic,
			"CREATE DATA MOVEMENT RULE %s TYPE = 'COPY_INTO_EXTERNAL_STAGE' MAX_ROWS AS () RETURNS INTEGER -> %s",
			id.FullyQualifiedName(), maxRows,
		).
		withModifyAndExpectedSqlf(
			case_DataMovementRules_sql_Create_all,
			func(opts *CreateDataMovementRuleOptions) {
				opts.IfNotExists = new(true)
				opts.Type = DataMovementRuleTypeSnowsightUi
				opts.Comment = new(comment)
			},
			"CREATE DATA MOVEMENT RULE IF NOT EXISTS %s TYPE = 'SNOWSIGHT_UI' MAX_ROWS AS () RETURNS INTEGER -> %s COMMENT = '%s'",
			id.FullyQualifiedName(), maxRows, comment,
		).
		withAdditionalSqlCasef(
			"sql_Create_orReplace",
			func(opts *CreateDataMovementRuleOptions) { opts.OrReplace = new(true) },
			"CREATE OR REPLACE DATA MOVEMENT RULE %s TYPE = 'COPY_INTO_EXTERNAL_STAGE' MAX_ROWS AS () RETURNS INTEGER -> %s",
			id.FullyQualifiedName(), maxRows,
		).
		withAdditionalSqlCasef(
			"sql_Create_expression",
			func(opts *CreateDataMovementRuleOptions) {
				opts.MaxRows = "CASE WHEN SYS_CONTEXT('SNOWFLAKE$SESSION', 'ROLE') = 'COMP_TEAM_ROLE' THEN NULL ELSE 0 END"
			},
			"CREATE DATA MOVEMENT RULE %s TYPE = 'COPY_INTO_EXTERNAL_STAGE' MAX_ROWS AS () RETURNS INTEGER -> CASE WHEN SYS_CONTEXT('SNOWFLAKE$SESSION', 'ROLE') = 'COMP_TEAM_ROLE' THEN NULL ELSE 0 END",
			id.FullyQualifiedName(),
		)

	dataMovementRulesTests.Alter.
		withModify(
			case_DataMovementRules_validation_Alter_opts_ExactlyOneValueSet_MoreThanOneSet,
			func(opts *AlterDataMovementRuleOptions) {
				opts.RenameTo = &renameTarget
				opts.Set = &DataMovementRuleSet{SetMaxRows: new(newMaxRows)}
			},
		).
		withModifyAndExpectedSqlf(
			case_DataMovementRules_sql_Alter_RenameTo,
			func(opts *AlterDataMovementRuleOptions) { opts.RenameTo = &renameTarget },
			"ALTER DATA MOVEMENT RULE %s RENAME TO %s",
			id.FullyQualifiedName(), renameTarget.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_DataMovementRules_sql_Alter_Set,
			func(opts *AlterDataMovementRuleOptions) {
				opts.Set = &DataMovementRuleSet{Comment: new(comment)}
			},
			"ALTER DATA MOVEMENT RULE %s SET COMMENT = '%s'",
			id.FullyQualifiedName(), comment,
		).
		withModifyAndExpectedSqlf(
			case_DataMovementRules_sql_Alter_Unset,
			func(opts *AlterDataMovementRuleOptions) {
				opts.Unset = &DataMovementRuleUnset{Comment: new(true)}
			},
			"ALTER DATA MOVEMENT RULE %s UNSET COMMENT",
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_setMaxRows",
			func(opts *AlterDataMovementRuleOptions) {
				opts.Set = &DataMovementRuleSet{SetMaxRows: new(newMaxRows)}
			},
			"ALTER DATA MOVEMENT RULE %s SET MAX_ROWS AS () RETURNS INTEGER -> %s",
			id.FullyQualifiedName(), newMaxRows,
		).
		withAdditionalSqlCasef(
			"sql_Alter_ifExists",
			func(opts *AlterDataMovementRuleOptions) {
				opts.IfExists = new(true)
				opts.Set = &DataMovementRuleSet{SetMaxRows: new(newMaxRows)}
			},
			"ALTER DATA MOVEMENT RULE IF EXISTS %s SET MAX_ROWS AS () RETURNS INTEGER -> %s",
			id.FullyQualifiedName(), newMaxRows,
		)

	dataMovementRulesTests.Drop.
		withExpectedSqlf(
			case_DataMovementRules_sql_Drop_basic,
			"DROP DATA MOVEMENT RULE %s",
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_DataMovementRules_sql_Drop_all,
			func(opts *DropDataMovementRuleOptions) { opts.IfExists = new(true) },
			"DROP DATA MOVEMENT RULE IF EXISTS %s",
			id.FullyQualifiedName(),
		)

	dataMovementRulesTests.Show.
		withExpectedSql(
			case_DataMovementRules_sql_Show_basic,
			"SHOW DATA MOVEMENT RULES",
		).
		withModifyAndExpectedSqlf(
			case_DataMovementRules_sql_Show_all,
			func(opts *ShowDataMovementRuleOptions) {
				opts.Like = &Like{Pattern: String("like-pattern")}
				opts.In = &In{Account: new(true)}
				opts.Limit = &LimitFrom{Rows: Int(10)}
			},
			"SHOW DATA MOVEMENT RULES LIKE 'like-pattern' IN ACCOUNT LIMIT 10",
		).
		withModifyAndExpectedSqlf(
			case_DataMovementRules_sql_Show_Like,
			func(opts *ShowDataMovementRuleOptions) { opts.Like = &Like{Pattern: String("like-pattern")} },
			"SHOW DATA MOVEMENT RULES LIKE 'like-pattern'",
		).
		withModifyAndExpectedSqlf(
			case_DataMovementRules_sql_Show_In,
			func(opts *ShowDataMovementRuleOptions) { opts.In = &In{Account: new(true)} },
			"SHOW DATA MOVEMENT RULES IN ACCOUNT",
		).
		withModifyAndExpectedSqlf(
			case_DataMovementRules_sql_Show_Limit,
			func(opts *ShowDataMovementRuleOptions) { opts.Limit = &LimitFrom{Rows: Int(10)} },
			"SHOW DATA MOVEMENT RULES LIMIT 10",
		).
		withAdditionalSqlCasef(
			"sql_Show_In_Database",
			func(opts *ShowDataMovementRuleOptions) {
				opts.In = &In{Database: id.DatabaseId()}
			},
			"SHOW DATA MOVEMENT RULES IN DATABASE %s",
			id.DatabaseId().FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Show_In_Schema",
			func(opts *ShowDataMovementRuleOptions) {
				opts.In = &In{Schema: id.SchemaId()}
			},
			"SHOW DATA MOVEMENT RULES IN SCHEMA %s",
			id.SchemaId().FullyQualifiedName(),
		)

	dataMovementRulesTests.Describe.
		withExpectedSqlf(
			case_DataMovementRules_sql_Describe_basic,
			"DESCRIBE DATA MOVEMENT RULE %s",
			id.FullyQualifiedName(),
		)
}
