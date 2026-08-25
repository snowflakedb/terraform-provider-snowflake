package sdk

func init() {
	id := pipesTestIdSchemaObjectIdentifier
	tagId1 := NewAccountObjectIdentifier("tag_name1")
	tagId2 := NewAccountObjectIdentifier("tag_name2")

	pipesTests.Create.
		withDefaultOpts(func() *CreatePipeOptions {
			return &CreatePipeOptions{
				name:          id,
				copyStatement: "<copy_statement>",
			}
		}).
		withAdditionalValidationCase(
			"validation_Create_copyStatement_required",
			func(opts *CreatePipeOptions) { opts.copyStatement = "" },
			errNotSet("CreatePipeOptions", "copyStatement"),
		).
		withExpectedSqlf(
			case_Pipes_sql_Create_basic,
			`CREATE PIPE %s AS <copy_statement>`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Pipes_sql_Create_all,
			func(opts *CreatePipeOptions) {
				opts.IfNotExists = Bool(true)
				opts.AutoIngest = Bool(true)
				opts.ErrorIntegration = String("some_error_integration")
				opts.AwsSnsTopic = String("some aws sns topic")
				opts.Integration = String("some integration")
				opts.Comment = String("some comment")
			},
			`CREATE PIPE IF NOT EXISTS %s AUTO_INGEST = true ERROR_INTEGRATION = some_error_integration AWS_SNS_TOPIC = 'some aws sns topic' INTEGRATION = 'some integration' COMMENT = 'some comment' AS <copy_statement>`,
			id.FullyQualifiedName(),
		)

	pipesTests.Alter.
		withModifyAndExpectedSqlf(
			case_Pipes_sql_Alter_Set,
			func(opts *AlterPipeOptions) {
				opts.IfExists = Bool(true)
				opts.Set = &PipeSet{
					ErrorIntegration:    String("new_error_integration"),
					PipeExecutionPaused: Bool(true),
					Comment:             String("new comment"),
				}
			},
			`ALTER PIPE IF EXISTS %s SET ERROR_INTEGRATION = new_error_integration, PIPE_EXECUTION_PAUSED = true, COMMENT = 'new comment'`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Pipes_sql_Alter_Unset,
			func(opts *AlterPipeOptions) {
				opts.IfExists = Bool(true)
				opts.Unset = &PipeUnset{
					ErrorIntegration:    Bool(true),
					PipeExecutionPaused: Bool(true),
					Comment:             Bool(true),
				}
			},
			`ALTER PIPE IF EXISTS %s UNSET ERROR_INTEGRATION, PIPE_EXECUTION_PAUSED, COMMENT`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Pipes_sql_Alter_SetTags,
			func(opts *AlterPipeOptions) {
				opts.SetTags = []TagAssociation{
					{Name: tagId1, Value: "v1"},
				}
			},
			`ALTER PIPE %s SET TAG %s = 'v1'`, id.FullyQualifiedName(), tagId1.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Pipes_sql_Alter_UnsetTags,
			func(opts *AlterPipeOptions) {
				opts.UnsetTags = []ObjectIdentifier{tagId1}
			},
			`ALTER PIPE %s UNSET TAG %s`, id.FullyQualifiedName(), tagId1.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Pipes_sql_Alter_Refresh,
			func(opts *AlterPipeOptions) { opts.Refresh = &PipeRefresh{} },
			`ALTER PIPE %s REFRESH`, id.FullyQualifiedName(),
		)

	pipesTests.Alter.
		withAdditionalSqlCasef(
			"sql_Alter_SetTags_multiple",
			func(opts *AlterPipeOptions) {
				opts.SetTags = []TagAssociation{
					{Name: tagId1, Value: "v1"},
					{Name: tagId2, Value: "v2"},
				}
			},
			`ALTER PIPE %s SET TAG %s = 'v1', %s = 'v2'`,
			id.FullyQualifiedName(), tagId1.FullyQualifiedName(), tagId2.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_UnsetTags_multiple",
			func(opts *AlterPipeOptions) {
				opts.UnsetTags = []ObjectIdentifier{tagId1, tagId2}
			},
			`ALTER PIPE %s UNSET TAG %s, %s`,
			id.FullyQualifiedName(), tagId1.FullyQualifiedName(), tagId2.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Refresh_all",
			func(opts *AlterPipeOptions) {
				opts.IfExists = Bool(true)
				opts.Refresh = &PipeRefresh{
					Prefix:        String("/d1"),
					ModifiedAfter: String("2018-07-30T13:56:46-07:00"),
				}
			},
			`ALTER PIPE IF EXISTS %s REFRESH PREFIX = '/d1' MODIFIED_AFTER = '2018-07-30T13:56:46-07:00'`,
			id.FullyQualifiedName(),
		)

	pipesTests.Drop.
		withExpectedSqlf(
			case_Pipes_sql_Drop_basic,
			`DROP PIPE %s`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Pipes_sql_Drop_all,
			func(opts *DropPipeOptions) { opts.IfExists = Bool(true) },
			`DROP PIPE IF EXISTS %s`, id.FullyQualifiedName(),
		)

	pipesTests.Show.
		withAdditionalValidationCase(
			"validation_Show_Like_patternRequired",
			func(opts *ShowPipeOptions) { opts.Like = &Like{} },
			ErrPatternRequiredForLikeKeyword,
		).
		withAdditionalValidationCase(
			"validation_Show_In_noScope",
			func(opts *ShowPipeOptions) { opts.In = &In{} },
			errExactlyOneOf("ShowPipeOptions.In", "Account", "Database", "Schema"),
		).
		withAdditionalValidationCase(
			"validation_Show_In_moreThanOneScope",
			func(opts *ShowPipeOptions) {
				opts.In = &In{Account: Bool(true), Database: id.DatabaseId()}
			},
			errExactlyOneOf("ShowPipeOptions.In", "Account", "Database", "Schema"),
		).
		withExpectedSql(case_Pipes_sql_Show_basic, `SHOW PIPES`).
		withModifyAndExpectedSqlf(
			case_Pipes_sql_Show_Like,
			func(opts *ShowPipeOptions) { opts.Like = &Like{Pattern: String(id.Name())} },
			`SHOW PIPES LIKE '%s'`, id.Name(),
		).
		withModifyAndExpectedSqlf(
			case_Pipes_sql_Show_In,
			func(opts *ShowPipeOptions) { opts.In = &In{Account: Bool(true)} },
			`SHOW PIPES IN ACCOUNT`,
		).
		withModifyAndExpectedSqlf(
			case_Pipes_sql_Show_all,
			func(opts *ShowPipeOptions) {
				opts.Like = &Like{Pattern: String(id.Name())}
				opts.In = &In{Schema: id.SchemaId()}
			},
			`SHOW PIPES LIKE '%s' IN SCHEMA %s`, id.Name(), id.SchemaId().FullyQualifiedName(),
		)

	pipesTests.Show.
		withAdditionalSqlCasef(
			"sql_Show_In_Database",
			func(opts *ShowPipeOptions) { opts.In = &In{Database: id.DatabaseId()} },
			`SHOW PIPES IN DATABASE %s`, id.DatabaseId().FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Show_In_Schema",
			func(opts *ShowPipeOptions) { opts.In = &In{Schema: id.SchemaId()} },
			`SHOW PIPES IN SCHEMA %s`, id.SchemaId().FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Show_Like_In_Account",
			func(opts *ShowPipeOptions) {
				opts.Like = &Like{Pattern: String(id.Name())}
				opts.In = &In{Account: Bool(true)}
			},
			`SHOW PIPES LIKE '%s' IN ACCOUNT`, id.Name(),
		).
		withAdditionalSqlCasef(
			"sql_Show_Like_In_Database",
			func(opts *ShowPipeOptions) {
				opts.Like = &Like{Pattern: String(id.Name())}
				opts.In = &In{Database: id.DatabaseId()}
			},
			`SHOW PIPES LIKE '%s' IN DATABASE %s`, id.Name(), id.DatabaseId().FullyQualifiedName(),
		)

	pipesTests.Describe.
		withExpectedSqlf(
			case_Pipes_sql_Describe_basic,
			`DESCRIBE PIPE %s`, id.FullyQualifiedName(),
		)
}
