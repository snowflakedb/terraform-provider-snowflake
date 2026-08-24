package sdk

func init() {
	id := notebooksTestIdSchemaObjectIdentifier
	stageId := randomSchemaObjectIdentifier()
	stageLocation := NewStageLocation(stageId, "dir/subdir")
	var from Location = stageLocation
	queryWarehouseId := NewAccountObjectIdentifier("sample_qwh")
	warehouseId := NewAccountObjectIdentifier("sample_wh")
	computePoolId := NewAccountObjectIdentifier("sample_cp")
	secretId := NewSchemaObjectIdentifier("db_name", "sc_name", "n_name")
	eaiId := NewAccountObjectIdentifier("test")
	tagId := NewAccountObjectIdentifier("tag1")
	renameTarget := randomSchemaObjectIdentifier()
	databaseId := NewAccountObjectIdentifier("database-name")

	notebooksTests.Create.
		withModifyAndExpectedSqlf(
			case_Notebooks_sql_Create_basic,
			func(opts *CreateNotebookOptions) { opts.OrReplace = new(true) },
			"CREATE OR REPLACE NOTEBOOK %s", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Notebooks_sql_Create_all,
			func(opts *CreateNotebookOptions) {
				opts.IfNotExists = new(true)
				opts.From = &from
				opts.MainFile = new("main_file")
				opts.Comment = new("comment")
				opts.QueryWarehouse = &queryWarehouseId
				opts.IdleAutoShutdownTimeSeconds = new(3600)
				opts.Warehouse = &warehouseId
				opts.RuntimeName = new("sample")
				opts.ComputePool = &computePoolId
				opts.RuntimeEnvironmentVersion = new("WH-RUNTIME-2.0")
				opts.DefaultVersion = new("FIRST")
			},
			`CREATE NOTEBOOK IF NOT EXISTS %s FROM '@\"%s\".\"%s\".\"%s\"/dir/subdir' MAIN_FILE = 'main_file' COMMENT = 'comment' QUERY_WAREHOUSE = %s IDLE_AUTO_SHUTDOWN_TIME_SECONDS = 3600 WAREHOUSE = %s RUNTIME_NAME = 'sample' COMPUTE_POOL = %s RUNTIME_ENVIRONMENT_VERSION = 'WH-RUNTIME-2.0' DEFAULT_VERSION = FIRST`,
			id.FullyQualifiedName(), stageId.DatabaseName(), stageId.SchemaName(), stageId.Name(), queryWarehouseId.FullyQualifiedName(), warehouseId.FullyQualifiedName(), computePoolId.FullyQualifiedName(),
		)

	notebooksTests.Alter.
		withModifyAndExpectedSqlf(
			case_Notebooks_sql_Alter_Set,
			func(opts *AlterNotebookOptions) {
				opts.Set = &NotebookSet{Comment: new("comment")}
			},
			"ALTER NOTEBOOK %s SET COMMENT = 'comment'", id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_all",
			func(opts *AlterNotebookOptions) {
				opts.Set = &NotebookSet{
					Comment:                     new("comment"),
					QueryWarehouse:              &queryWarehouseId,
					IdleAutoShutdownTimeSeconds: new(3600),
					Secrets:                     &SecretsList{[]SecretReference{{VariableName: "var_name", Name: secretId}}},
					MainFile:                    new("main_file"),
					Warehouse:                   &warehouseId,
					RuntimeName:                 new("runtime_name"),
					ComputePool:                 &computePoolId,
					ExternalAccessIntegrations:  []AccountObjectIdentifier{eaiId},
					RuntimeEnvironmentVersion:   new("WH-RUNTIME-2.0"),
				}
			},
			"ALTER NOTEBOOK %s SET COMMENT = 'comment' QUERY_WAREHOUSE = %s IDLE_AUTO_SHUTDOWN_TIME_SECONDS = 3600 SECRETS = ('var_name' = %s) MAIN_FILE = 'main_file' WAREHOUSE = %s RUNTIME_NAME = 'runtime_name' COMPUTE_POOL = %s EXTERNAL_ACCESS_INTEGRATIONS = (%s) RUNTIME_ENVIRONMENT_VERSION = 'WH-RUNTIME-2.0'",
			id.FullyQualifiedName(), queryWarehouseId.FullyQualifiedName(), secretId.FullyQualifiedName(), warehouseId.FullyQualifiedName(), computePoolId.FullyQualifiedName(), eaiId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Notebooks_sql_Alter_Unset,
			func(opts *AlterNotebookOptions) {
				opts.Unset = &NotebookUnset{
					Comment:                    new(true),
					QueryWarehouse:             new(true),
					Secrets:                    new(true),
					Warehouse:                  new(true),
					RuntimeName:                new(true),
					ComputePool:                new(true),
					ExternalAccessIntegrations: new(true),
					RuntimeEnvironmentVersion:  new(true),
				}
			},
			"ALTER NOTEBOOK %s UNSET COMMENT, QUERY_WAREHOUSE, SECRETS, WAREHOUSE, RUNTIME_NAME, COMPUTE_POOL, EXTERNAL_ACCESS_INTEGRATIONS, RUNTIME_ENVIRONMENT_VERSION",
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Notebooks_sql_Alter_SetTags,
			func(opts *AlterNotebookOptions) {
				opts.SetTags = []TagAssociation{{Name: tagId, Value: "value1"}}
			},
			`ALTER NOTEBOOK %s SET TAG %s = 'value1'`,
			id.FullyQualifiedName(), tagId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Notebooks_sql_Alter_UnsetTags,
			func(opts *AlterNotebookOptions) {
				opts.UnsetTags = []ObjectIdentifier{tagId}
			},
			`ALTER NOTEBOOK %s UNSET TAG %s`,
			id.FullyQualifiedName(), tagId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Notebooks_sql_Alter_RenameTo,
			func(opts *AlterNotebookOptions) { opts.RenameTo = &renameTarget },
			`ALTER NOTEBOOK %s RENAME TO %s`,
			id.FullyQualifiedName(), renameTarget.FullyQualifiedName(),
		)

	notebooksTests.Drop.
		withExpectedSqlf(
			case_Notebooks_sql_Drop_basic,
			"DROP NOTEBOOK %s", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Notebooks_sql_Drop_all,
			func(opts *DropNotebookOptions) { opts.IfExists = new(true) },
			"DROP NOTEBOOK IF EXISTS %s", id.FullyQualifiedName(),
		)

	notebooksTests.Describe.
		withExpectedSqlf(
			case_Notebooks_sql_Describe_basic,
			"DESCRIBE NOTEBOOK %s", id.FullyQualifiedName(),
		)

	notebooksTests.Show.
		withExpectedSql(case_Notebooks_sql_Show_basic, "SHOW NOTEBOOKS").
		withModifyAndExpectedSqlf(
			case_Notebooks_sql_Show_all,
			func(opts *ShowNotebookOptions) {
				opts.Like = &Like{Pattern: new("notebook-name")}
				opts.In = &In{Database: databaseId}
				opts.Limit = &LimitFrom{Rows: new(10)}
				opts.StartsWith = new("prefix")
			},
			`SHOW NOTEBOOKS LIKE 'notebook-name' IN DATABASE %s LIMIT 10 STARTS WITH 'prefix'`,
			databaseId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Notebooks_sql_Show_Like,
			func(opts *ShowNotebookOptions) { opts.Like = &Like{Pattern: new("notebook-name")} },
			`SHOW NOTEBOOKS LIKE 'notebook-name'`,
		).
		withModifyAndExpectedSqlf(
			case_Notebooks_sql_Show_In,
			func(opts *ShowNotebookOptions) { opts.In = &In{Database: databaseId} },
			`SHOW NOTEBOOKS IN DATABASE %s`,
			databaseId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Notebooks_sql_Show_Limit,
			func(opts *ShowNotebookOptions) { opts.Limit = &LimitFrom{Rows: new(10)} },
			`SHOW NOTEBOOKS LIMIT 10`,
		).
		withModifyAndExpectedSqlf(
			case_Notebooks_sql_Show_StartsWith,
			func(opts *ShowNotebookOptions) { opts.StartsWith = new("prefix") },
			`SHOW NOTEBOOKS STARTS WITH 'prefix'`,
		)
}
