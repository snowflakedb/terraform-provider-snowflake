package sdk

import "time"

func init() {
	id := schemasTestIdDatabaseObjectIdentifier
	externalVolumeId := randomAccountObjectIdentifier()
	catalogId := randomAccountObjectIdentifier()
	tagId := randomSchemaObjectIdentifier()
	cloneSourceId := NewAccountObjectIdentifier("sch1")
	cloneTimestamp := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	renameTarget := randomDatabaseObjectIdentifierInDatabase(id.DatabaseId())
	swapWithId := randomDatabaseObjectIdentifierInDatabase(id.DatabaseId())
	tagId1 := NewAccountObjectIdentifier("tag1")
	tagId2 := NewAccountObjectIdentifier("tag2")
	databaseId := NewAccountObjectIdentifier("database_name")

	schemasTests.Create.
		withExpectedSqlf(
			case_Schemas_sql_Create_basic,
			`CREATE SCHEMA %s`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Schemas_sql_Create_all,
			func(opts *CreateSchemaOptions) {
				opts.Transient = new(true)
				opts.IfNotExists = new(true)
				opts.WithManagedAccess = new(true)
				opts.DataRetentionTimeInDays = new(1)
				opts.MaxDataExtensionTimeInDays = new(1)
				opts.ExternalVolume = &externalVolumeId
				opts.Catalog = &catalogId
				opts.PipeExecutionPaused = new(true)
				opts.ReplaceInvalidCharacters = new(true)
				opts.DefaultDdlCollation = new(StringAllowEmpty{Value: "en_US-trim"})
				opts.DefaultNotebookComputePoolCpu = new("CPU_X64_S")
				opts.DefaultNotebookComputePoolGpu = new("GPU_NV_S")
				opts.StorageSerializationPolicy = new(StorageSerializationPolicyCompatible)
				opts.LogLevel = new(LogLevelInfo)
				opts.TraceLevel = new(TraceLevelPropagate)
				opts.SuspendTaskAfterNumFailures = new(10)
				opts.TaskAutoRetryAttempts = new(10)
				opts.UserTaskManagedInitialWarehouseSize = new(WarehouseSizeMedium)
				opts.UserTaskTimeoutMs = new(12000)
				opts.UserTaskMinimumTriggerIntervalInSeconds = new(30)
				opts.QuotedIdentifiersIgnoreCase = new(true)
				opts.EnableConsoleOutput = new(true)
				opts.Comment = new("comment")
				opts.Tag = []TagAssociation{{Name: tagId, Value: "v1"}}
			},
			`CREATE TRANSIENT SCHEMA IF NOT EXISTS %s WITH MANAGED ACCESS DATA_RETENTION_TIME_IN_DAYS = 1 MAX_DATA_EXTENSION_TIME_IN_DAYS = 1 `+
				`EXTERNAL_VOLUME = %s CATALOG = %s PIPE_EXECUTION_PAUSED = true REPLACE_INVALID_CHARACTERS = true DEFAULT_DDL_COLLATION = 'en_US-trim' DEFAULT_NOTEBOOK_COMPUTE_POOL_CPU = 'CPU_X64_S' DEFAULT_NOTEBOOK_COMPUTE_POOL_GPU = 'GPU_NV_S' STORAGE_SERIALIZATION_POLICY = COMPATIBLE `+
				`LOG_LEVEL = 'INFO' TRACE_LEVEL = 'PROPAGATE' SUSPEND_TASK_AFTER_NUM_FAILURES = 10 TASK_AUTO_RETRY_ATTEMPTS = 10 USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE = MEDIUM `+
				`USER_TASK_TIMEOUT_MS = 12000 USER_TASK_MINIMUM_TRIGGER_INTERVAL_IN_SECONDS = 30 QUOTED_IDENTIFIERS_IGNORE_CASE = true ENABLE_CONSOLE_OUTPUT = true `+
				`COMMENT = 'comment' TAG (%s = 'v1')`,
			id.FullyQualifiedName(), externalVolumeId.FullyQualifiedName(), catalogId.FullyQualifiedName(), tagId.FullyQualifiedName(),
		)

	schemasTests.Clone.
		withDefaultOpts(func() *CloneSchemaOptions {
			return &CloneSchemaOptions{
				name:  id,
				Clone: Clone{SourceObject: cloneSourceId},
			}
		}).
		withExpectedSqlf(
			case_Schemas_sql_Clone_basic,
			`CREATE SCHEMA %s CLONE %s`, id.FullyQualifiedName(), cloneSourceId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Clone_all",
			func(opts *CloneSchemaOptions) {
				opts.OrReplace = new(true)
				opts.Clone = Clone{
					SourceObject: cloneSourceId,
					At:           &TimeTravel{Timestamp: new(cloneTimestamp)},
				}
			},
			`CREATE OR REPLACE SCHEMA %s CLONE %s AT (TIMESTAMP => '2021-01-01 00:00:00 +0000 UTC')`,
			id.FullyQualifiedName(), cloneSourceId.FullyQualifiedName(),
		)

	schemasTests.Alter.
		withModifyAndExpectedSqlf(
			case_Schemas_sql_Alter_RenameTo,
			func(opts *AlterSchemaOptions) {
				opts.IfExists = new(true)
				opts.RenameTo = &renameTarget
			},
			`ALTER SCHEMA IF EXISTS %s RENAME TO %s`, id.FullyQualifiedName(), renameTarget.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Schemas_sql_Alter_SwapWith,
			func(opts *AlterSchemaOptions) { opts.SwapWith = &swapWithId },
			`ALTER SCHEMA %s SWAP WITH %s`, id.FullyQualifiedName(), swapWithId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Schemas_sql_Alter_Set,
			func(opts *AlterSchemaOptions) {
				opts.Set = &SchemaSet{
					DataRetentionTimeInDays:                 new(1),
					MaxDataExtensionTimeInDays:              new(1),
					ExternalVolume:                          &externalVolumeId,
					Catalog:                                 &catalogId,
					PipeExecutionPaused:                     new(true),
					ReplaceInvalidCharacters:                new(true),
					DefaultDdlCollation:                     new(StringAllowEmpty{Value: "en_US-trim"}),
					DefaultNotebookComputePoolCpu:           new("CPU_X64_S"),
					DefaultNotebookComputePoolGpu:           new("GPU_NV_S"),
					StorageSerializationPolicy:              new(StorageSerializationPolicyCompatible),
					LogLevel:                                new(LogLevelInfo),
					TraceLevel:                              new(TraceLevelPropagate),
					SuspendTaskAfterNumFailures:             new(10),
					TaskAutoRetryAttempts:                   new(10),
					UserTaskManagedInitialWarehouseSize:     new(WarehouseSizeMedium),
					UserTaskTimeoutMs:                       new(12000),
					UserTaskMinimumTriggerIntervalInSeconds: new(30),
					QuotedIdentifiersIgnoreCase:             new(true),
					EnableConsoleOutput:                     new(true),
					Comment:                                 new("comment"),
				}
			},
			`ALTER SCHEMA %s SET DATA_RETENTION_TIME_IN_DAYS = 1, MAX_DATA_EXTENSION_TIME_IN_DAYS = 1, `+
				`EXTERNAL_VOLUME = %s, CATALOG = %s, PIPE_EXECUTION_PAUSED = true, REPLACE_INVALID_CHARACTERS = true, DEFAULT_DDL_COLLATION = 'en_US-trim', DEFAULT_NOTEBOOK_COMPUTE_POOL_CPU = 'CPU_X64_S', DEFAULT_NOTEBOOK_COMPUTE_POOL_GPU = 'GPU_NV_S', STORAGE_SERIALIZATION_POLICY = COMPATIBLE, `+
				`LOG_LEVEL = 'INFO', TRACE_LEVEL = 'PROPAGATE', SUSPEND_TASK_AFTER_NUM_FAILURES = 10, TASK_AUTO_RETRY_ATTEMPTS = 10, USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE = MEDIUM, `+
				`USER_TASK_TIMEOUT_MS = 12000, USER_TASK_MINIMUM_TRIGGER_INTERVAL_IN_SECONDS = 30, QUOTED_IDENTIFIERS_IGNORE_CASE = true, ENABLE_CONSOLE_OUTPUT = true, `+
				`COMMENT = 'comment'`,
			id.FullyQualifiedName(), externalVolumeId.FullyQualifiedName(), catalogId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Schemas_sql_Alter_Unset,
			func(opts *AlterSchemaOptions) {
				opts.Unset = &SchemaUnset{
					DataRetentionTimeInDays:                 new(true),
					MaxDataExtensionTimeInDays:              new(true),
					ExternalVolume:                          new(true),
					Catalog:                                 new(true),
					PipeExecutionPaused:                     new(true),
					ReplaceInvalidCharacters:                new(true),
					DefaultDdlCollation:                     new(true),
					DefaultNotebookComputePoolCpu:           new(true),
					DefaultNotebookComputePoolGpu:           new(true),
					StorageSerializationPolicy:              new(true),
					LogLevel:                                new(true),
					TraceLevel:                              new(true),
					SuspendTaskAfterNumFailures:             new(true),
					TaskAutoRetryAttempts:                   new(true),
					UserTaskManagedInitialWarehouseSize:     new(true),
					UserTaskTimeoutMs:                       new(true),
					UserTaskMinimumTriggerIntervalInSeconds: new(true),
					QuotedIdentifiersIgnoreCase:             new(true),
					EnableConsoleOutput:                     new(true),
					Comment:                                 new(true),
				}
			},
			`ALTER SCHEMA %s UNSET DATA_RETENTION_TIME_IN_DAYS, MAX_DATA_EXTENSION_TIME_IN_DAYS, EXTERNAL_VOLUME, CATALOG, PIPE_EXECUTION_PAUSED, `+
				`REPLACE_INVALID_CHARACTERS, DEFAULT_DDL_COLLATION, DEFAULT_NOTEBOOK_COMPUTE_POOL_CPU, DEFAULT_NOTEBOOK_COMPUTE_POOL_GPU, STORAGE_SERIALIZATION_POLICY, LOG_LEVEL, TRACE_LEVEL, SUSPEND_TASK_AFTER_NUM_FAILURES, TASK_AUTO_RETRY_ATTEMPTS, `+
				`USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE, USER_TASK_TIMEOUT_MS, USER_TASK_MINIMUM_TRIGGER_INTERVAL_IN_SECONDS, QUOTED_IDENTIFIERS_IGNORE_CASE, ENABLE_CONSOLE_OUTPUT, COMMENT`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Schemas_sql_Alter_SetTags,
			func(opts *AlterSchemaOptions) {
				opts.SetTags = []TagAssociation{
					{Name: tagId1, Value: "value1"},
					{Name: tagId2, Value: "value2"},
				}
			},
			`ALTER SCHEMA %s SET TAG %s = 'value1', %s = 'value2'`,
			id.FullyQualifiedName(), tagId1.FullyQualifiedName(), tagId2.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Schemas_sql_Alter_UnsetTags,
			func(opts *AlterSchemaOptions) {
				opts.UnsetTags = []ObjectIdentifier{tagId1, tagId2}
			},
			`ALTER SCHEMA %s UNSET TAG %s, %s`,
			id.FullyQualifiedName(), tagId1.FullyQualifiedName(), tagId2.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Schemas_sql_Alter_EnableManagedAccess,
			func(opts *AlterSchemaOptions) { opts.EnableManagedAccess = new(true) },
			`ALTER SCHEMA %s ENABLE MANAGED ACCESS`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Schemas_sql_Alter_DisableManagedAccess,
			func(opts *AlterSchemaOptions) { opts.DisableManagedAccess = new(true) },
			`ALTER SCHEMA %s DISABLE MANAGED ACCESS`, id.FullyQualifiedName(),
		)

	schemasTests.Drop.
		withExpectedSqlf(
			case_Schemas_sql_Drop_basic,
			`DROP SCHEMA %s`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Schemas_sql_Drop_all,
			func(opts *DropSchemaOptions) {
				opts.IfExists = new(true)
				opts.Cascade = new(true)
			},
			`DROP SCHEMA IF EXISTS %s CASCADE`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Drop_restrict",
			func(opts *DropSchemaOptions) {
				opts.IfExists = new(true)
				opts.Restrict = new(true)
			},
			`DROP SCHEMA IF EXISTS %s RESTRICT`, id.FullyQualifiedName(),
		)

	schemasTests.Undrop.
		withExpectedSqlf(
			case_Schemas_sql_Undrop_basic,
			`UNDROP SCHEMA %s`, id.FullyQualifiedName(),
		)

	schemasTests.Show.
		withExpectedSql(case_Schemas_sql_Show_basic, `SHOW SCHEMAS`).
		withModifyAndExpectedSqlf(
			case_Schemas_sql_Show_all,
			func(opts *ShowSchemaOptions) {
				opts.Terse = new(true)
				opts.History = new(true)
				opts.Like = &Like{Pattern: new("schema_pattern")}
				opts.In = &ExtendedIn{In: In{Database: databaseId}}
				opts.StartsWith = new("schema_pattern")
				opts.Limit = &LimitFrom{Rows: new(3), From: new("name_string")}
			},
			`SHOW TERSE SCHEMAS HISTORY LIKE 'schema_pattern' IN DATABASE %s STARTS WITH 'schema_pattern' LIMIT 3 FROM 'name_string'`,
			databaseId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Schemas_sql_Show_Like,
			func(opts *ShowSchemaOptions) {
				opts.Terse = new(true)
				opts.History = new(true)
				opts.Like = &Like{Pattern: new("schema_pattern")}
			},
			`SHOW TERSE SCHEMAS HISTORY LIKE 'schema_pattern'`,
		).
		withModifyAndExpectedSqlf(
			case_Schemas_sql_Show_In,
			func(opts *ShowSchemaOptions) {
				opts.Terse = new(true)
				opts.History = new(true)
				opts.In = &ExtendedIn{In: In{Account: new(true)}}
			},
			`SHOW TERSE SCHEMAS HISTORY IN ACCOUNT`,
		).
		withAdditionalSqlCasef(
			"sql_Show_In_database",
			func(opts *ShowSchemaOptions) {
				opts.Terse = new(true)
				opts.History = new(true)
				opts.In = &ExtendedIn{In: In{Database: databaseId}}
			},
			`SHOW TERSE SCHEMAS HISTORY IN DATABASE %s`,
			databaseId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Schemas_sql_Show_StartsWith,
			func(opts *ShowSchemaOptions) {
				opts.Terse = new(true)
				opts.History = new(true)
				opts.StartsWith = new("schema_pattern")
			},
			`SHOW TERSE SCHEMAS HISTORY STARTS WITH 'schema_pattern'`,
		).
		withModifyAndExpectedSqlf(
			case_Schemas_sql_Show_Limit,
			func(opts *ShowSchemaOptions) {
				opts.Terse = new(true)
				opts.History = new(true)
				opts.Limit = &LimitFrom{Rows: new(3), From: new("name_string")}
			},
			`SHOW TERSE SCHEMAS HISTORY LIMIT 3 FROM 'name_string'`,
		)

	schemasTests.Describe.
		withExpectedSqlf(
			case_Schemas_sql_Describe_basic,
			`DESCRIBE SCHEMA %s`, id.FullyQualifiedName(),
		)
}
