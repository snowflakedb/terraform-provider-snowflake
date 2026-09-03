package sdk

func init() {
	id := cortexSearchServicesTestIdSchemaObjectIdentifier
	warehouseId := randomAccountObjectIdentifier()
	tagId := NewAccountObjectIdentifier("tag1")
	tagId2 := NewAccountObjectIdentifier("tag2")
	databaseId := NewAccountObjectIdentifier("database")
	queryDefinition := "SELECT product_id, product_name, searchable_text FROM staging_table"

	cortexSearchServicesTests.Create.
		withDefaultOpts(func() *CreateCortexSearchServiceOptions {
			return &CreateCortexSearchServiceOptions{
				name:            id,
				On:              "searchable_text",
				TargetLag:       "1 minutes",
				Warehouse:       warehouseId,
				QueryDefinition: queryDefinition,
			}
		}).
		withExpectedSqlf(
			case_CortexSearchServices_sql_Create_basic,
			`CREATE CORTEX SEARCH SERVICE %s ON searchable_text WAREHOUSE = %s TARGET_LAG = '1 minutes' AS %s`,
			id.FullyQualifiedName(), warehouseId.FullyQualifiedName(), queryDefinition,
		).
		withModifyAndExpectedSqlf(
			case_CortexSearchServices_sql_Create_all,
			func(opts *CreateCortexSearchServiceOptions) {
				opts.IfNotExists = new(true)
				opts.Attributes = &Attributes{
					Columns: []string{"product_id", "product_name"},
				}
				opts.EmbeddingModel = new("snowflake-arctic-embed-l-v2.0")
				opts.PrimaryKey = []string{"product_id"}
				opts.RefreshMode = new(CortexSearchServiceRefreshModeFull)
				opts.Initialize = new(CortexSearchServiceInitializeOnCreate)
				opts.FullIndexBuildIntervalDays = new(30)
				opts.RequestLogging = new(true)
				opts.AutoSuspend = new(3600)
				opts.Comment = new("comment")
			},
			`CREATE CORTEX SEARCH SERVICE IF NOT EXISTS %s ON searchable_text PRIMARY KEY (product_id) ATTRIBUTES product_id, product_name WAREHOUSE = %s TARGET_LAG = '1 minutes' EMBEDDING_MODEL = 'snowflake-arctic-embed-l-v2.0' REFRESH_MODE = FULL INITIALIZE = ON_CREATE FULL_INDEX_BUILD_INTERVAL_DAYS = 30 REQUEST_LOGGING = true AUTO_SUSPEND = 3600 COMMENT = 'comment' AS %s`,
			id.FullyQualifiedName(), warehouseId.FullyQualifiedName(), queryDefinition,
		).
		withAdditionalSqlCasef(
			"sql_Create_embeddingModel",
			func(opts *CreateCortexSearchServiceOptions) {
				opts.EmbeddingModel = new("snowflake-arctic-embed-m-v1.5")
			},
			`CREATE CORTEX SEARCH SERVICE %s ON searchable_text WAREHOUSE = %s TARGET_LAG = '1 minutes' EMBEDDING_MODEL = 'snowflake-arctic-embed-m-v1.5' AS %s`,
			id.FullyQualifiedName(), warehouseId.FullyQualifiedName(), queryDefinition,
		).
		withAdditionalSqlCasef(
			"sql_Create_primaryKey",
			func(opts *CreateCortexSearchServiceOptions) {
				opts.PrimaryKey = []string{"id", "ts"}
			},
			`CREATE CORTEX SEARCH SERVICE %s ON searchable_text PRIMARY KEY (id, ts) WAREHOUSE = %s TARGET_LAG = '1 minutes' AS %s`,
			id.FullyQualifiedName(), warehouseId.FullyQualifiedName(), queryDefinition,
		).
		withAdditionalSqlCasef(
			"sql_Create_refreshModeAndInitialize",
			func(opts *CreateCortexSearchServiceOptions) {
				opts.RefreshMode = new(CortexSearchServiceRefreshModeIncremental)
				opts.Initialize = new(CortexSearchServiceInitializeOnSchedule)
			},
			`CREATE CORTEX SEARCH SERVICE %s ON searchable_text WAREHOUSE = %s TARGET_LAG = '1 minutes' REFRESH_MODE = INCREMENTAL INITIALIZE = ON_SCHEDULE AS %s`,
			id.FullyQualifiedName(), warehouseId.FullyQualifiedName(), queryDefinition,
		).
		withAdditionalSqlCasef(
			"sql_Create_orReplace",
			func(opts *CreateCortexSearchServiceOptions) { opts.OrReplace = new(true) },
			`CREATE OR REPLACE CORTEX SEARCH SERVICE %s ON searchable_text WAREHOUSE = %s TARGET_LAG = '1 minutes' AS %s`,
			id.FullyQualifiedName(), warehouseId.FullyQualifiedName(), queryDefinition,
		)

	cortexSearchServicesTests.Alter.
		withModifyAndExpectedSqlf(
			case_CortexSearchServices_sql_Alter_Suspend,
			func(opts *AlterCortexSearchServiceOptions) { opts.Suspend = new(true) },
			`ALTER CORTEX SEARCH SERVICE %s SUSPEND`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_CortexSearchServices_sql_Alter_Resume,
			func(opts *AlterCortexSearchServiceOptions) { opts.Resume = new(true) },
			`ALTER CORTEX SEARCH SERVICE %s RESUME`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_CortexSearchServices_sql_Alter_Refresh,
			func(opts *AlterCortexSearchServiceOptions) { opts.Refresh = new(true) },
			`ALTER CORTEX SEARCH SERVICE %s REFRESH`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_CortexSearchServices_sql_Alter_Set,
			func(opts *AlterCortexSearchServiceOptions) {
				opts.Set = &CortexSearchServiceSet{TargetLag: new("1 minutes")}
			},
			`ALTER CORTEX SEARCH SERVICE %s SET TARGET_LAG = '1 minutes'`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_CortexSearchServices_sql_Alter_SetDefaults,
			func(opts *AlterCortexSearchServiceOptions) {
				opts.SetDefaults = &CortexSearchServiceSetDefaults{AutoSuspend: new(true)}
			},
			`ALTER CORTEX SEARCH SERVICE %s SET AUTO_SUSPEND = NULL`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_CortexSearchServices_sql_Alter_SetPrimaryKey,
			func(opts *AlterCortexSearchServiceOptions) {
				opts.SetPrimaryKey = &CortexSearchServiceSetPrimaryKey{
					PrimaryKey: []string{"id", "ts"},
				}
			},
			`ALTER CORTEX SEARCH SERVICE %s SET PRIMARY KEY = (id, ts)`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_CortexSearchServices_sql_Alter_SetAttributes,
			func(opts *AlterCortexSearchServiceOptions) {
				opts.SetAttributes = &CortexSearchServiceSetAttributes{
					Columns: []string{"col1", "col2"},
				}
			},
			`ALTER CORTEX SEARCH SERVICE %s SET ATTRIBUTES (col1, col2)`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_CortexSearchServices_sql_Alter_UnsetPrimaryKey,
			func(opts *AlterCortexSearchServiceOptions) { opts.UnsetPrimaryKey = new(true) },
			`ALTER CORTEX SEARCH SERVICE %s UNSET PRIMARY KEY`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_CortexSearchServices_sql_Alter_UnsetAttributes,
			func(opts *AlterCortexSearchServiceOptions) { opts.UnsetAttributes = new(true) },
			`ALTER CORTEX SEARCH SERVICE %s UNSET ATTRIBUTES`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_CortexSearchServices_sql_Alter_SetTags,
			func(opts *AlterCortexSearchServiceOptions) {
				opts.SetTags = []TagAssociation{
					{Name: tagId, Value: "value1"},
					{Name: tagId2, Value: "value2"},
				}
			},
			`ALTER CORTEX SEARCH SERVICE %s SET TAG %s = 'value1', %s = 'value2'`,
			id.FullyQualifiedName(), tagId.FullyQualifiedName(), tagId2.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_CortexSearchServices_sql_Alter_UnsetTags,
			func(opts *AlterCortexSearchServiceOptions) {
				opts.UnsetTags = []ObjectIdentifier{tagId, tagId2}
			},
			`ALTER CORTEX SEARCH SERVICE %s UNSET TAG %s, %s`,
			id.FullyQualifiedName(), tagId.FullyQualifiedName(), tagId2.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_all",
			func(opts *AlterCortexSearchServiceOptions) {
				opts.Set = &CortexSearchServiceSet{
					TargetLag:                  new("1 minutes"),
					Warehouse:                  &warehouseId,
					FullIndexBuildIntervalDays: new(30),
					RequestLogging:             new(true),
					AutoSuspend:                new(3600),
					Comment:                    new("comment"),
				}
			},
			`ALTER CORTEX SEARCH SERVICE %s SET TARGET_LAG = '1 minutes' WAREHOUSE = %s FULL_INDEX_BUILD_INTERVAL_DAYS = 30 REQUEST_LOGGING = true AUTO_SUSPEND = 3600 COMMENT = 'comment'`,
			id.FullyQualifiedName(), warehouseId.FullyQualifiedName(),
		)

	cortexSearchServicesTests.Show.
		withExpectedSql(case_CortexSearchServices_sql_Show_basic, "SHOW CORTEX SEARCH SERVICES").
		withModifyAndExpectedSqlf(
			case_CortexSearchServices_sql_Show_all,
			func(opts *ShowCortexSearchServiceOptions) {
				opts.Like = &Like{Pattern: new("pattern")}
				opts.In = &In{Account: new(true)}
				opts.StartsWith = new("foo")
				opts.Limit = &LimitFrom{Rows: new(1), From: new("bar")}
			},
			`SHOW CORTEX SEARCH SERVICES LIKE 'pattern' IN ACCOUNT STARTS WITH 'foo' LIMIT 1 FROM 'bar'`,
		).
		withModifyAndExpectedSqlf(
			case_CortexSearchServices_sql_Show_Like,
			func(opts *ShowCortexSearchServiceOptions) {
				opts.Like = &Like{Pattern: new(id.Name())}
			},
			`SHOW CORTEX SEARCH SERVICES LIKE '%s'`, id.Name(),
		).
		withModifyAndExpectedSqlf(
			case_CortexSearchServices_sql_Show_In,
			func(opts *ShowCortexSearchServiceOptions) {
				opts.In = &In{Database: databaseId}
			},
			`SHOW CORTEX SEARCH SERVICES IN DATABASE %s`, databaseId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_CortexSearchServices_sql_Show_StartsWith,
			func(opts *ShowCortexSearchServiceOptions) { opts.StartsWith = new("foo") },
			`SHOW CORTEX SEARCH SERVICES STARTS WITH 'foo'`,
		).
		withModifyAndExpectedSqlf(
			case_CortexSearchServices_sql_Show_Limit,
			func(opts *ShowCortexSearchServiceOptions) {
				opts.Limit = &LimitFrom{Rows: new(1)}
			},
			`SHOW CORTEX SEARCH SERVICES LIMIT 1`,
		).
		withAdditionalSqlCasef(
			"sql_Show_LimitFrom",
			func(opts *ShowCortexSearchServiceOptions) {
				opts.Limit = &LimitFrom{Rows: new(1), From: new("foo")}
			},
			`SHOW CORTEX SEARCH SERVICES LIMIT 1 FROM 'foo'`,
		)

	cortexSearchServicesTests.Describe.
		withExpectedSqlf(
			case_CortexSearchServices_sql_Describe_basic,
			`DESCRIBE CORTEX SEARCH SERVICE %s`, id.FullyQualifiedName(),
		)

	cortexSearchServicesTests.Drop.
		withExpectedSqlf(
			case_CortexSearchServices_sql_Drop_basic,
			`DROP CORTEX SEARCH SERVICE %s`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_CortexSearchServices_sql_Drop_all,
			func(opts *DropCortexSearchServiceOptions) { opts.IfExists = new(true) },
			`DROP CORTEX SEARCH SERVICE IF EXISTS %s`, id.FullyQualifiedName(),
		)
}
