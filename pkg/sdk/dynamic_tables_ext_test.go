package sdk

func init() {
	id := dynamicTablesTestIdSchemaObjectIdentifier
	warehouseId := NewAccountObjectIdentifier("warehouse_name")
	databaseId := NewAccountObjectIdentifier("database")
	storageLifecyclePolicyId := randomSchemaObjectIdentifier()
	query := "SELECT product_id, product_name FROM staging_table"

	dynamicTablesTests.Create.
		withDefaultOpts(func() *CreateDynamicTableOptions {
			return &CreateDynamicTableOptions{
				name: id,
				TargetLag: TargetLag{
					MaximumDuration: new("1 minutes"),
				},
				Warehouse: warehouseId,
				Query:     query,
			}
		}).
		withModifyAndExpectedSqlf(
			case_DynamicTables_sql_Create_basic,
			func(opts *CreateDynamicTableOptions) { opts.OrReplace = new(true) },
			`CREATE OR REPLACE DYNAMIC TABLE %s TARGET_LAG = '1 minutes' WAREHOUSE = %s AS %s`,
			id.FullyQualifiedName(), warehouseId.FullyQualifiedName(), query,
		).
		withModifyAndExpectedSqlf(
			case_DynamicTables_sql_Create_all,
			func(opts *CreateDynamicTableOptions) {
				opts.OrReplace = new(true)
				opts.Initialize = new(DynamicTableInitializeOnSchedule)
				opts.RefreshMode = new(DynamicTableRefreshModeFull)
				opts.Comment = new("comment")
			},
			`CREATE OR REPLACE DYNAMIC TABLE %s TARGET_LAG = '1 minutes' INITIALIZE = ON_SCHEDULE REFRESH_MODE = FULL WAREHOUSE = %s COMMENT = 'comment' AS %s`,
			id.FullyQualifiedName(), warehouseId.FullyQualifiedName(), query,
		)

	dynamicTablesTests.Alter.
		withAdditionalValidationCase(
			"validation_Alter_AddStorageLifecyclePolicy_invalidIdentifier",
			func(opts *AlterDynamicTableOptions) {
				opts.AddStorageLifecyclePolicy = &DynamicTableAddStorageLifecyclePolicy{
					StorageLifecyclePolicy: emptySchemaObjectIdentifier,
					On:                     []Column{{Value: "FIRST_COLUMN"}},
				}
			},
			errInvalidIdentifier("DynamicTableAddStorageLifecyclePolicy", "StorageLifecyclePolicy"),
		).
		withAdditionalValidationCase(
			"validation_Alter_AddStorageLifecyclePolicy_On_notSet",
			func(opts *AlterDynamicTableOptions) {
				opts.AddStorageLifecyclePolicy = &DynamicTableAddStorageLifecyclePolicy{
					StorageLifecyclePolicy: storageLifecyclePolicyId,
				}
			},
			errNotSet("DynamicTableAddStorageLifecyclePolicy", "On"),
		).
		withAdditionalValidationCase(
			"validation_Alter_Set_Warehouse_invalidIdentifier",
			func(opts *AlterDynamicTableOptions) {
				opts.Set = &DynamicTableSet{Warehouse: new(emptyAccountObjectIdentifier)}
			},
			errInvalidIdentifier("DynamicTableSet", "Warehouse"),
		).
		withModifyAndExpectedSqlf(
			case_DynamicTables_sql_Alter_Suspend,
			func(opts *AlterDynamicTableOptions) { opts.Suspend = new(true) },
			`ALTER DYNAMIC TABLE %s SUSPEND`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_DynamicTables_sql_Alter_Resume,
			func(opts *AlterDynamicTableOptions) { opts.Resume = new(true) },
			`ALTER DYNAMIC TABLE %s RESUME`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_DynamicTables_sql_Alter_Refresh,
			func(opts *AlterDynamicTableOptions) { opts.Refresh = new(true) },
			`ALTER DYNAMIC TABLE %s REFRESH`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_DynamicTables_sql_Alter_Set,
			func(opts *AlterDynamicTableOptions) {
				opts.Set = &DynamicTableSet{
					TargetLag: &TargetLag{MaximumDuration: new("1 minutes")},
					Warehouse: &warehouseId,
				}
			},
			`ALTER DYNAMIC TABLE %s SET TARGET_LAG = '1 minutes' WAREHOUSE = %s`,
			id.FullyQualifiedName(), warehouseId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_DynamicTables_sql_Alter_SetComment,
			func(opts *AlterDynamicTableOptions) { opts.SetComment = new("some comment") },
			`ALTER DYNAMIC TABLE %s SET COMMENT = 'some comment'`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_DynamicTables_sql_Alter_AddStorageLifecyclePolicy,
			func(opts *AlterDynamicTableOptions) {
				opts.AddStorageLifecyclePolicy = &DynamicTableAddStorageLifecyclePolicy{
					StorageLifecyclePolicy: storageLifecyclePolicyId,
					On:                     []Column{{Value: "FIRST_COLUMN"}},
				}
			},
			`ALTER DYNAMIC TABLE %s ADD STORAGE LIFECYCLE POLICY %s ON ("FIRST_COLUMN")`,
			id.FullyQualifiedName(), storageLifecyclePolicyId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_DynamicTables_sql_Alter_DropStorageLifecyclePolicy,
			func(opts *AlterDynamicTableOptions) { opts.DropStorageLifecyclePolicy = new(true) },
			`ALTER DYNAMIC TABLE %s DROP STORAGE LIFECYCLE POLICY`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_AddStorageLifecyclePolicy_multipleColumns",
			func(opts *AlterDynamicTableOptions) {
				opts.AddStorageLifecyclePolicy = &DynamicTableAddStorageLifecyclePolicy{
					StorageLifecyclePolicy: storageLifecyclePolicyId,
					On:                     []Column{{Value: "FIRST_COLUMN"}, {Value: "SECOND_COLUMN"}},
				}
			},
			`ALTER DYNAMIC TABLE %s ADD STORAGE LIFECYCLE POLICY %s ON ("FIRST_COLUMN", "SECOND_COLUMN")`,
			id.FullyQualifiedName(), storageLifecyclePolicyId.FullyQualifiedName(),
		)

	dynamicTablesTests.Drop.
		withExpectedSqlf(
			case_DynamicTables_sql_Drop_basic,
			`DROP DYNAMIC TABLE %s`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_DynamicTables_sql_Drop_all,
			func(opts *DropDynamicTableOptions) { opts.IfExists = new(true) },
			`DROP DYNAMIC TABLE IF EXISTS %s`, id.FullyQualifiedName(),
		)

	dynamicTablesTests.Show.
		withAdditionalValidationCase(
			"validation_Show_Like_patternRequired",
			func(opts *ShowDynamicTableOptions) { opts.Like = &Like{} },
			ErrPatternRequiredForLikeKeyword,
		).
		withAdditionalValidationCase(
			"validation_Show_In_noScope",
			func(opts *ShowDynamicTableOptions) { opts.In = &In{} },
			errExactlyOneOf("ShowDynamicTableOptions.In", "Account", "Database", "Schema"),
		).
		withAdditionalValidationCase(
			"validation_Show_In_moreThanOneScope",
			func(opts *ShowDynamicTableOptions) {
				opts.In = &In{Account: new(true), Database: databaseId}
			},
			errExactlyOneOf("ShowDynamicTableOptions.In", "Account", "Database", "Schema"),
		).
		withExpectedSql(case_DynamicTables_sql_Show_basic, `SHOW DYNAMIC TABLES`).
		withModifyAndExpectedSqlf(
			case_DynamicTables_sql_Show_all,
			func(opts *ShowDynamicTableOptions) {
				opts.Like = &Like{Pattern: new(id.Name())}
				opts.In = &In{Database: databaseId}
				opts.StartsWith = new("prefix")
				opts.Limit = &LimitFrom{Rows: new(10), From: new("foo")}
			},
			`SHOW DYNAMIC TABLES LIKE '%s' IN DATABASE %s STARTS WITH 'prefix' LIMIT 10 FROM 'foo'`,
			id.Name(), databaseId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_DynamicTables_sql_Show_Like,
			func(opts *ShowDynamicTableOptions) {
				opts.Like = &Like{Pattern: new(id.Name())}
			},
			`SHOW DYNAMIC TABLES LIKE '%s'`, id.Name(),
		).
		withModifyAndExpectedSqlf(
			case_DynamicTables_sql_Show_In,
			func(opts *ShowDynamicTableOptions) {
				opts.In = &In{Database: databaseId}
			},
			`SHOW DYNAMIC TABLES IN DATABASE %s`, databaseId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_DynamicTables_sql_Show_StartsWith,
			func(opts *ShowDynamicTableOptions) { opts.StartsWith = new("prefix") },
			`SHOW DYNAMIC TABLES STARTS WITH 'prefix'`,
		).
		withModifyAndExpectedSqlf(
			case_DynamicTables_sql_Show_Limit,
			func(opts *ShowDynamicTableOptions) {
				opts.Limit = &LimitFrom{Rows: new(10), From: new("foo")}
			},
			`SHOW DYNAMIC TABLES LIMIT 10 FROM 'foo'`,
		).
		withAdditionalSqlCasef(
			"sql_Show_Like_In",
			func(opts *ShowDynamicTableOptions) {
				opts.Like = &Like{Pattern: new(id.Name())}
				opts.In = &In{Database: databaseId}
			},
			`SHOW DYNAMIC TABLES LIKE '%s' IN DATABASE %s`,
			id.Name(), databaseId.FullyQualifiedName(),
		)

	dynamicTablesTests.Describe.
		withExpectedSqlf(
			case_DynamicTables_sql_Describe_basic,
			`DESCRIBE DYNAMIC TABLE %s`, id.FullyQualifiedName(),
		)
}
