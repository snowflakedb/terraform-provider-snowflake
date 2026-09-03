package sdk

func init() {
	id := streamlitsTestIdSchemaObjectIdentifier
	warehouse := NewAccountObjectIdentifier("test_warehouse")
	integration := NewAccountObjectIdentifier("integration")
	renameTarget := randomSchemaObjectIdentifier()

	streamlitsTests.Create.
		withDefaultOpts(func() *CreateStreamlitOptions {
			return &CreateStreamlitOptions{
				name:         id,
				RootLocation: "@test",
				MainFile:     "manifest.yml",
			}
		}).
		withExpectedSqlf(
			case_Streamlits_sql_Create_basic,
			`CREATE STREAMLIT %s ROOT_LOCATION = '@test' MAIN_FILE = 'manifest.yml'`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Streamlits_sql_Create_all,
			func(opts *CreateStreamlitOptions) {
				opts.IfNotExists = new(true)
				opts.QueryWarehouse = &warehouse
				opts.ExternalAccessIntegrations = []AccountObjectIdentifier{integration}
				opts.Title = new("foo")
				opts.Comment = new("test")
			},
			`CREATE STREAMLIT IF NOT EXISTS %s ROOT_LOCATION = '@test' MAIN_FILE = 'manifest.yml' QUERY_WAREHOUSE = %s EXTERNAL_ACCESS_INTEGRATIONS = (%s) TITLE = 'foo' COMMENT = 'test'`,
			id.FullyQualifiedName(), warehouse.FullyQualifiedName(), integration.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Create_orReplace",
			func(opts *CreateStreamlitOptions) { opts.OrReplace = new(true) },
			`CREATE OR REPLACE STREAMLIT %s ROOT_LOCATION = '@test' MAIN_FILE = 'manifest.yml'`,
			id.FullyQualifiedName(),
		)

	streamlitsTests.Alter.
		withDefaultOpts(func() *AlterStreamlitOptions {
			return &AlterStreamlitOptions{
				IfExists: new(true),
				name:     id,
			}
		}).
		withModify(
			case_Streamlits_validation_Alter_opts_ExactlyOneValueSet_MoreThanOneSet,
			func(opts *AlterStreamlitOptions) {
				opts.RenameTo = &renameTarget
				opts.Set = &StreamlitSet{Comment: new("test")}
			},
		).
		withModifyAndExpectedSqlf(
			case_Streamlits_sql_Alter_RenameTo,
			func(opts *AlterStreamlitOptions) { opts.RenameTo = &renameTarget },
			`ALTER STREAMLIT IF EXISTS %s RENAME TO %s`,
			id.FullyQualifiedName(), renameTarget.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Streamlits_sql_Alter_Set,
			func(opts *AlterStreamlitOptions) {
				opts.Set = &StreamlitSet{
					RootLocation:               new("@test"),
					MainFile:                   new("manifest.yml"),
					QueryWarehouse:             &warehouse,
					ExternalAccessIntegrations: []AccountObjectIdentifier{integration},
					Comment:                    new("test"),
					Title:                      new("foo"),
				}
			},
			`ALTER STREAMLIT IF EXISTS %s SET ROOT_LOCATION = '@test' MAIN_FILE = 'manifest.yml' QUERY_WAREHOUSE = %s EXTERNAL_ACCESS_INTEGRATIONS = (%s) COMMENT = 'test' TITLE = 'foo'`,
			id.FullyQualifiedName(), warehouse.FullyQualifiedName(), integration.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Streamlits_sql_Alter_Unset,
			func(opts *AlterStreamlitOptions) {
				opts.Unset = &StreamlitUnset{
					QueryWarehouse:             new(true),
					Comment:                    new(true),
					Title:                      new(true),
					ExternalAccessIntegrations: new(true),
				}
			},
			`ALTER STREAMLIT IF EXISTS %s UNSET QUERY_WAREHOUSE, COMMENT, TITLE, EXTERNAL_ACCESS_INTEGRATIONS`,
			id.FullyQualifiedName(),
		)

	streamlitsTests.Drop.
		withExpectedSqlf(
			case_Streamlits_sql_Drop_basic,
			`DROP STREAMLIT %s`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Streamlits_sql_Drop_all,
			func(opts *DropStreamlitOptions) { opts.IfExists = new(true) },
			`DROP STREAMLIT IF EXISTS %s`, id.FullyQualifiedName(),
		)

	streamlitsTests.Show.
		withAdditionalValidationCase(
			"validation_Show_Like_patternRequired",
			func(opts *ShowStreamlitOptions) { opts.Like = &Like{} },
			ErrPatternRequiredForLikeKeyword,
		).
		withExpectedSql(case_Streamlits_sql_Show_basic, "SHOW STREAMLITS").
		withModifyAndExpectedSqlf(
			case_Streamlits_sql_Show_all,
			func(opts *ShowStreamlitOptions) {
				opts.Terse = new(true)
				opts.Like = &Like{Pattern: new("pattern")}
				opts.In = &In{Account: new(true)}
				opts.Limit = &LimitFrom{Rows: new(123), From: new("from pattern")}
			},
			`SHOW TERSE STREAMLITS LIKE 'pattern' IN ACCOUNT LIMIT 123 FROM 'from pattern'`,
		).
		withModifyAndExpectedSqlf(
			case_Streamlits_sql_Show_Like,
			func(opts *ShowStreamlitOptions) { opts.Like = &Like{Pattern: new("pattern")} },
			`SHOW STREAMLITS LIKE 'pattern'`,
		).
		withModifyAndExpectedSqlf(
			case_Streamlits_sql_Show_In,
			func(opts *ShowStreamlitOptions) { opts.In = &In{Account: new(true)} },
			`SHOW STREAMLITS IN ACCOUNT`,
		).
		withModifyAndExpectedSqlf(
			case_Streamlits_sql_Show_Limit,
			func(opts *ShowStreamlitOptions) {
				opts.Limit = &LimitFrom{Rows: new(123), From: new("from pattern")}
			},
			`SHOW STREAMLITS LIMIT 123 FROM 'from pattern'`,
		).
		withAdditionalSqlCasef(
			"sql_Show_Terse",
			func(opts *ShowStreamlitOptions) { opts.Terse = new(true) },
			`SHOW TERSE STREAMLITS`,
		)

	streamlitsTests.Describe.
		withExpectedSqlf(
			case_Streamlits_sql_Describe_basic,
			`DESCRIBE STREAMLIT %s`, id.FullyQualifiedName(),
		)
}
