package sdk

func init() {
	id := snowflakeIntelligencesTestIdAccountObjectIdentifier
	agentId := randomSchemaObjectIdentifier()

	snowflakeIntelligencesTests.Create.
		withExpectedSqlf(
			case_SnowflakeIntelligences_sql_Create_basic,
			`CREATE SNOWFLAKE INTELLIGENCE %s`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_SnowflakeIntelligences_sql_Create_all,
			func(opts *CreateSnowflakeIntelligenceOptions) { opts.IfNotExists = new(true) },
			`CREATE SNOWFLAKE INTELLIGENCE IF NOT EXISTS %s`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Create_orReplace",
			func(opts *CreateSnowflakeIntelligenceOptions) { opts.OrReplace = new(true) },
			`CREATE OR REPLACE SNOWFLAKE INTELLIGENCE %s`, id.FullyQualifiedName(),
		)

	snowflakeIntelligencesTests.Alter.
		withModifyAndExpectedSqlf(
			case_SnowflakeIntelligences_sql_Alter_AddAgent,
			func(opts *AlterSnowflakeIntelligenceOptions) {
				opts.IfExists = new(true)
				opts.AddAgent = &agentId
			},
			`ALTER SNOWFLAKE INTELLIGENCE IF EXISTS %s ADD AGENT %s`,
			id.FullyQualifiedName(), agentId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_SnowflakeIntelligences_sql_Alter_DropAgent,
			func(opts *AlterSnowflakeIntelligenceOptions) {
				opts.DropAgent = &agentId
			},
			`ALTER SNOWFLAKE INTELLIGENCE %s DROP AGENT %s`,
			id.FullyQualifiedName(), agentId.FullyQualifiedName(),
		)

	snowflakeIntelligencesTests.Drop.
		withExpectedSqlf(
			case_SnowflakeIntelligences_sql_Drop_basic,
			`DROP SNOWFLAKE INTELLIGENCE %s`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_SnowflakeIntelligences_sql_Drop_all,
			func(opts *DropSnowflakeIntelligenceOptions) { opts.IfExists = new(true) },
			`DROP SNOWFLAKE INTELLIGENCE IF EXISTS %s`, id.FullyQualifiedName(),
		)

	snowflakeIntelligencesTests.Show.
		withExpectedSql(case_SnowflakeIntelligences_sql_Show_basic, `SHOW SNOWFLAKE INTELLIGENCES`).
		withModifyAndExpectedSqlf(
			case_SnowflakeIntelligences_sql_Show_all,
			func(*ShowSnowflakeIntelligenceOptions) {},
			`SHOW SNOWFLAKE INTELLIGENCES`,
		)

	snowflakeIntelligencesTests.Describe.
		withExpectedSqlf(
			case_SnowflakeIntelligences_sql_Describe_basic,
			`DESCRIBE SNOWFLAKE INTELLIGENCE %s`, id.FullyQualifiedName(),
		)

	snowflakeIntelligencesTests.ShowAgents.
		withExpectedSqlf(
			case_SnowflakeIntelligences_sql_ShowAgents_basic,
			`SHOW AGENTS IN SNOWFLAKE INTELLIGENCE %s`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_SnowflakeIntelligences_sql_ShowAgents_all,
			func(*ShowAgentsSnowflakeIntelligenceOptions) {},
			`SHOW AGENTS IN SNOWFLAKE INTELLIGENCE %s`, id.FullyQualifiedName(),
		)
}
