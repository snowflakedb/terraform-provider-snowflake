package sdk

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	sql := "SELECT CURRENT_TIMESTAMP"
	warehouseId := randomAccountObjectIdentifier()
	otherTaskId := randomSchemaObjectIdentifier()
	tagId := randomSchemaObjectIdentifier()
	finalizerId := randomSchemaObjectIdentifier()

	tasksTests.Create.
		withDefaultOpts(func() *CreateTaskOptions {
			return &CreateTaskOptions{
				name: tasksTestIdSchemaObjectIdentifier,
				sql:  sql,
			}
		}).
		withModify(
			case_Tasks_validation_Create_opts_Warehouse_ExactlyOneValueSet_MoreThanOneSet,
			func(opts *CreateTaskOptions) {
				opts.Warehouse = &CreateTaskWarehouse{
					Warehouse:                           &warehouseId,
					UserTaskManagedInitialWarehouseSize: new(WarehouseSizeXSmall),
				}
			},
		).
		withAdditionalValidationCase(
			"validation_Create_SessionParameters_shouldBeValid",
			func(opts *CreateTaskOptions) {
				opts.SessionParameters = &SessionParameters{JsonIndent: new(-1)}
			},
			errIntValue("SessionParameters", "JsonIndent", IntErrGreaterOrEqual, 0),
		).
		withExpectedSqlf(
			case_Tasks_sql_Create_basic,
			"CREATE TASK %s AS %s", tasksTestIdSchemaObjectIdentifier.FullyQualifiedName(), sql,
		).
		withModifyAndExpectedSqlf(
			case_Tasks_sql_Create_all,
			func(opts *CreateTaskOptions) {
				opts.IfNotExists = new(true)
				opts.Warehouse = &CreateTaskWarehouse{Warehouse: &warehouseId}
				opts.Schedule = new("10 MINUTE")
				opts.Config = new(`{"output_dir": "/temp/test_directory/", "learning_rate": 0.1}`)
				opts.AllowOverlappingExecution = new(true)
				opts.SessionParameters = &SessionParameters{
					JsonIndent:  new(10),
					LockTimeout: new(5),
				}
				opts.UserTaskTimeoutMs = new(5)
				opts.SuspendTaskAfterNumFailures = new(6)
				opts.ErrorIntegration = new(NewAccountObjectIdentifier("some_error_integration"))
				opts.Comment = new("some comment")
				opts.Finalize = &finalizerId
				opts.TaskAutoRetryAttempts = new(10)
				opts.Tag = []TagAssociation{{Name: tagId, Value: "v1"}}
				opts.UserTaskMinimumTriggerIntervalInSeconds = new(10)
				opts.TargetCompletionInterval = new("10 MINUTES")
				opts.ServerlessTaskMinStatementSize = new(WarehouseSizeSmall)
				opts.ServerlessTaskMaxStatementSize = new(WarehouseSizeLarge)
				opts.After = []SchemaObjectIdentifier{otherTaskId}
				opts.When = new(`SYSTEM$STREAM_HAS_DATA('MYSTREAM')`)
			},
			`CREATE TASK IF NOT EXISTS %s WAREHOUSE = %s SCHEDULE = '10 MINUTE' CONFIG = $${"output_dir": "/temp/test_directory/", "learning_rate": 0.1}$$ ALLOW_OVERLAPPING_EXECUTION = true JSON_INDENT = 10, LOCK_TIMEOUT = 5 USER_TASK_TIMEOUT_MS = 5 SUSPEND_TASK_AFTER_NUM_FAILURES = 6 ERROR_INTEGRATION = "some_error_integration" COMMENT = 'some comment' FINALIZE = %s TASK_AUTO_RETRY_ATTEMPTS = 10 TAG (%s = 'v1') USER_TASK_MINIMUM_TRIGGER_INTERVAL_IN_SECONDS = 10 TARGET_COMPLETION_INTERVAL = '10 MINUTES' SERVERLESS_TASK_MIN_STATEMENT_SIZE = 'SMALL' SERVERLESS_TASK_MAX_STATEMENT_SIZE = 'LARGE' AFTER %s WHEN SYSTEM$STREAM_HAS_DATA('MYSTREAM') AS SELECT CURRENT_TIMESTAMP`,
			tasksTestIdSchemaObjectIdentifier.FullyQualifiedName(), warehouseId.FullyQualifiedName(), finalizerId.FullyQualifiedName(), tagId.FullyQualifiedName(), otherTaskId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Create_withInitialWarehouseSize",
			func(opts *CreateTaskOptions) {
				opts.Warehouse = &CreateTaskWarehouse{
					UserTaskManagedInitialWarehouseSize: new(WarehouseSizeXSmall),
				}
			},
			"CREATE TASK %s USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE = 'XSMALL' AS %s",
			tasksTestIdSchemaObjectIdentifier.FullyQualifiedName(), sql,
		).
		withAdditionalSqlCasef(
			"sql_Create_orReplace",
			func(opts *CreateTaskOptions) { opts.OrReplace = new(true) },
			"CREATE OR REPLACE TASK %s AS %s", tasksTestIdSchemaObjectIdentifier.FullyQualifiedName(), sql,
		)

	tasksTests.CreateOrAlter.
		withDefaultOpts(func() *CreateOrAlterTaskOptions {
			return &CreateOrAlterTaskOptions{
				name: tasksTestIdSchemaObjectIdentifier,
				sql:  sql,
			}
		}).
		withModify(
			case_Tasks_validation_CreateOrAlter_opts_Warehouse_ExactlyOneValueSet_MoreThanOneSet,
			func(opts *CreateOrAlterTaskOptions) {
				opts.Warehouse = &CreateTaskWarehouse{
					Warehouse:                           &warehouseId,
					UserTaskManagedInitialWarehouseSize: new(WarehouseSizeXSmall),
				}
			},
		).
		withAdditionalValidationCase(
			"validation_CreateOrAlter_SessionParameters_shouldBeValid",
			func(opts *CreateOrAlterTaskOptions) {
				opts.SessionParameters = &SessionParameters{JsonIndent: new(-1)}
			},
			errIntValue("SessionParameters", "JsonIndent", IntErrGreaterOrEqual, 0),
		).
		withExpectedSqlf(
			case_Tasks_sql_CreateOrAlter_basic,
			"CREATE OR ALTER TASK %s AS %s", tasksTestIdSchemaObjectIdentifier.FullyQualifiedName(), sql,
		).
		withModifyAndExpectedSqlf(
			case_Tasks_sql_CreateOrAlter_all,
			func(opts *CreateOrAlterTaskOptions) {
				opts.Warehouse = &CreateTaskWarehouse{Warehouse: &warehouseId}
				opts.Schedule = new("10 MINUTE")
				opts.Config = new(`{"output_dir": "/temp/test_directory/", "learning_rate": 0.1}`)
				opts.AllowOverlappingExecution = new(true)
				opts.UserTaskTimeoutMs = new(5)
				opts.SessionParameters = &SessionParameters{
					JsonIndent:  new(10),
					LockTimeout: new(5),
				}
				opts.SuspendTaskAfterNumFailures = new(6)
				opts.ErrorIntegration = new(NewAccountObjectIdentifier("some_error_integration"))
				opts.Comment = new("some comment")
				opts.Finalize = &finalizerId
				opts.TaskAutoRetryAttempts = new(10)
				opts.After = []SchemaObjectIdentifier{otherTaskId}
				opts.When = new(`SYSTEM$STREAM_HAS_DATA('MYSTREAM')`)
			},
			`CREATE OR ALTER TASK %s WAREHOUSE = %s SCHEDULE = '10 MINUTE' CONFIG = $${"output_dir": "/temp/test_directory/", "learning_rate": 0.1}$$ ALLOW_OVERLAPPING_EXECUTION = true USER_TASK_TIMEOUT_MS = 5 JSON_INDENT = 10, LOCK_TIMEOUT = 5 SUSPEND_TASK_AFTER_NUM_FAILURES = 6 ERROR_INTEGRATION = "some_error_integration" COMMENT = 'some comment' FINALIZE = %s TASK_AUTO_RETRY_ATTEMPTS = 10 AFTER %s WHEN SYSTEM$STREAM_HAS_DATA('MYSTREAM') AS SELECT CURRENT_TIMESTAMP`,
			tasksTestIdSchemaObjectIdentifier.FullyQualifiedName(), warehouseId.FullyQualifiedName(), finalizerId.FullyQualifiedName(), otherTaskId.FullyQualifiedName(),
		)

	sourceTaskId := randomSchemaObjectIdentifier()

	tasksTests.Clone.
		withDefaultOpts(func() *CloneTaskOptions {
			return &CloneTaskOptions{
				name:       tasksTestIdSchemaObjectIdentifier,
				sourceTask: sourceTaskId,
			}
		}).
		withExpectedSqlf(
			case_Tasks_sql_Clone_basic,
			"CREATE TASK %s CLONE %s", tasksTestIdSchemaObjectIdentifier.FullyQualifiedName(), sourceTaskId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Clone_all",
			func(opts *CloneTaskOptions) {
				opts.OrReplace = new(true)
				opts.CopyGrants = new(true)
			},
			"CREATE OR REPLACE TASK %s CLONE %s COPY GRANTS", tasksTestIdSchemaObjectIdentifier.FullyQualifiedName(), sourceTaskId.FullyQualifiedName(),
		)

	finalizeId := randomSchemaObjectIdentifier()

	tasksTests.Alter.
		withModify(
			case_Tasks_validation_Alter_opts_Set_ConflictingFields,
			func(opts *AlterTaskOptions) {
				opts.Set = &TaskSet{
					Warehouse:                           &warehouseId,
					UserTaskManagedInitialWarehouseSize: new(WarehouseSizeXSmall),
				}
			},
		).
		withAdditionalValidationCase(
			"validation_Alter_Set_SessionParameters_shouldBeValid",
			func(opts *AlterTaskOptions) {
				opts.Set = &TaskSet{
					SessionParameters: &SessionParameters{JsonIndent: new(-1)},
				}
			},
			errIntValue("SessionParameters", "JsonIndent", IntErrGreaterOrEqual, 0),
		).
		withAdditionalValidationCase(
			"validation_Alter_Unset_SessionParametersUnset_shouldBeValid",
			func(opts *AlterTaskOptions) {
				opts.Unset = &TaskUnset{
					SessionParametersUnset: &SessionParametersUnset{},
				}
			},
			errAtLeastOneOf("SessionParametersUnset", "AbortDetachedQuery", "ActivePythonProfiler", "Autocommit", "BinaryInputFormat", "BinaryOutputFormat", "ClientEnableLogInfoStatementParameters", "ClientMemoryLimit", "ClientMetadataRequestUseConnectionCtx", "ClientPrefetchThreads", "ClientResultChunkSize", "ClientResultColumnCaseInsensitive", "ClientMetadataUseSessionDatabase", "ClientSessionKeepAlive", "ClientSessionKeepAliveHeartbeatFrequency", "ClientTimestampTypeMapping", "CsvTimestampFormat", "DateInputFormat", "DateOutputFormat", "EnableCortexAnalyst", "EnableGetDdlUseDataTypeAlias", "EnableUnloadPhysicalTypeOptimization", "ErrorOnNondeterministicMerge", "ErrorOnNondeterministicUpdate", "GeographyOutputFormat", "GeometryOutputFormat", "HybridTableLockTimeout", "JdbcTreatDecimalAsInt", "JdbcTreatTimestampNtzAsUtc", "JdbcUseSessionTimezone", "JsonIndent", "JsTreatIntegerAsBigInt", "LockTimeout", "LogLevel", "LogEventLevel", "MultiStatementCount", "NoorderSequenceAsDefault", "OdbcTreatDecimalAsInt", "PythonProfilerModules", "PythonProfilerTargetStage", "QueryTag", "QuotedIdentifiersIgnoreCase", "RowsPerResultset", "S3StageVpceDnsName", "SearchPath", "SimulatedDataSharingConsumer", "StatementQueuedTimeoutInSeconds", "StatementTimeoutInSeconds", "StrictJsonOutput", "TimestampDayIsAlways24h", "TimestampInputFormat", "TimestampLTZOutputFormat", "TimestampNTZOutputFormat", "TimestampOutputFormat", "TimestampTypeMapping", "TimestampTZOutputFormat", "Timezone", "TimeInputFormat", "TimeOutputFormat", "TraceLevel", "TransactionAbortOnError", "TransactionDefaultIsolationLevel", "TwoDigitCenturyStart", "UnsupportedDDLAction", "UseCachedResult", "WeekOfYearPolicy", "WeekStart"),
		).
		withModifyAndExpectedSqlf(
			case_Tasks_sql_Alter_Resume,
			func(opts *AlterTaskOptions) { opts.Resume = new(true) },
			"ALTER TASK %s RESUME", tasksTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Tasks_sql_Alter_Suspend,
			func(opts *AlterTaskOptions) { opts.Suspend = new(true) },
			"ALTER TASK %s SUSPEND", tasksTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Tasks_sql_Alter_RemoveAfter,
			func(opts *AlterTaskOptions) { opts.RemoveAfter = []SchemaObjectIdentifier{otherTaskId} },
			"ALTER TASK %s REMOVE AFTER %s", tasksTestIdSchemaObjectIdentifier.FullyQualifiedName(), otherTaskId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Tasks_sql_Alter_AddAfter,
			func(opts *AlterTaskOptions) { opts.AddAfter = []SchemaObjectIdentifier{otherTaskId} },
			"ALTER TASK %s ADD AFTER %s", tasksTestIdSchemaObjectIdentifier.FullyQualifiedName(), otherTaskId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Tasks_sql_Alter_Set,
			func(opts *AlterTaskOptions) {
				opts.Set = &TaskSet{
					UserTaskTimeoutMs:              new(2000),
					Comment:                        new("some comment"),
					TargetCompletionInterval:       new("15 MINUTES"),
					ServerlessTaskMinStatementSize: new(WarehouseSizeXSmall),
					ServerlessTaskMaxStatementSize: new(WarehouseSizeXLarge),
				}
			},
			"ALTER TASK %s SET USER_TASK_TIMEOUT_MS = 2000, COMMENT = 'some comment', TARGET_COMPLETION_INTERVAL = '15 MINUTES', SERVERLESS_TASK_MIN_STATEMENT_SIZE = 'XSMALL', SERVERLESS_TASK_MAX_STATEMENT_SIZE = 'XLARGE'",
			tasksTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_singleField",
			func(opts *AlterTaskOptions) {
				opts.Set = &TaskSet{Comment: new("some comment")}
			},
			"ALTER TASK %s SET COMMENT = 'some comment'", tasksTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_warehouse",
			func(opts *AlterTaskOptions) {
				opts.Set = &TaskSet{Warehouse: &warehouseId}
			},
			"ALTER TASK %s SET WAREHOUSE = %s", tasksTestIdSchemaObjectIdentifier.FullyQualifiedName(), warehouseId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_sessionParameter",
			func(opts *AlterTaskOptions) {
				opts.Set = &TaskSet{
					SessionParameters: &SessionParameters{JsonIndent: new(15)},
				}
			},
			"ALTER TASK %s SET JSON_INDENT = 15", tasksTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Tasks_sql_Alter_Unset,
			func(opts *AlterTaskOptions) {
				opts.Unset = &TaskUnset{
					UserTaskTimeoutMs:              new(true),
					Comment:                        new(true),
					TargetCompletionInterval:       new(true),
					ServerlessTaskMinStatementSize: new(true),
					ServerlessTaskMaxStatementSize: new(true),
				}
			},
			"ALTER TASK %s UNSET USER_TASK_TIMEOUT_MS, COMMENT, TARGET_COMPLETION_INTERVAL, SERVERLESS_TASK_MIN_STATEMENT_SIZE, SERVERLESS_TASK_MAX_STATEMENT_SIZE",
			tasksTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Unset_singleField",
			func(opts *AlterTaskOptions) {
				opts.Unset = &TaskUnset{Comment: new(true)}
			},
			"ALTER TASK %s UNSET COMMENT", tasksTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Tasks_sql_Alter_SetTags,
			func(opts *AlterTaskOptions) {
				opts.SetTags = []TagAssociation{
					{Name: NewAccountObjectIdentifier("tag1"), Value: "value1"},
					{Name: NewAccountObjectIdentifier("tag2"), Value: "value2"},
				}
			},
			`ALTER TASK %s SET TAG "tag1" = 'value1', "tag2" = 'value2'`, tasksTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Tasks_sql_Alter_UnsetTags,
			func(opts *AlterTaskOptions) {
				opts.UnsetTags = []ObjectIdentifier{
					NewAccountObjectIdentifier("tag1"),
					NewAccountObjectIdentifier("tag2"),
				}
			},
			`ALTER TASK %s UNSET TAG "tag1", "tag2"`, tasksTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Tasks_sql_Alter_ModifyAs,
			func(opts *AlterTaskOptions) { opts.ModifyAs = new("new as") },
			"ALTER TASK %s MODIFY AS new as", tasksTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Tasks_sql_Alter_ModifyWhen,
			func(opts *AlterTaskOptions) { opts.ModifyWhen = new("new when") },
			"ALTER TASK %s MODIFY WHEN new when", tasksTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Tasks_sql_Alter_SetFinalize,
			func(opts *AlterTaskOptions) { opts.SetFinalize = &finalizeId },
			"ALTER TASK %s SET FINALIZE = %s", tasksTestIdSchemaObjectIdentifier.FullyQualifiedName(), finalizeId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Tasks_sql_Alter_UnsetFinalize,
			func(opts *AlterTaskOptions) { opts.UnsetFinalize = new(true) },
			"ALTER TASK %s UNSET FINALIZE", tasksTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Tasks_sql_Alter_RemoveWhen,
			func(opts *AlterTaskOptions) { opts.RemoveWhen = new(true) },
			"ALTER TASK %s REMOVE WHEN", tasksTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		)

	tasksTests.Drop.
		withExpectedSqlf(
			case_Tasks_sql_Drop_basic,
			"DROP TASK %s", tasksTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Tasks_sql_Drop_all,
			func(opts *DropTaskOptions) { opts.IfExists = new(true) },
			"DROP TASK IF EXISTS %s", tasksTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		)

	applicationId := randomAccountObjectIdentifier()
	applicationPackageId := randomAccountObjectIdentifier()

	tasksTests.Show.
		withExpectedSql(case_Tasks_sql_Show_basic, "SHOW TASKS").
		withModifyAndExpectedSqlf(
			case_Tasks_sql_Show_all,
			func(opts *ShowTaskOptions) {
				opts.Terse = new(true)
				opts.Like = &Like{Pattern: new("myaccount")}
				opts.In = &ExtendedIn{In: In{Account: new(true)}}
				opts.StartsWith = new("abc")
				opts.RootOnly = new(true)
				opts.Limit = &LimitFrom{Rows: new(10)}
			},
			"SHOW TERSE TASKS LIKE 'myaccount' IN ACCOUNT STARTS WITH 'abc' ROOT ONLY LIMIT 10",
		).
		withModifyAndExpectedSqlf(
			case_Tasks_sql_Show_Like,
			func(opts *ShowTaskOptions) { opts.Like = &Like{Pattern: new("pattern")} },
			"SHOW TASKS LIKE 'pattern'",
		).
		withModifyAndExpectedSqlf(
			case_Tasks_sql_Show_In,
			func(opts *ShowTaskOptions) { opts.In = &ExtendedIn{In: In{Account: new(true)}} },
			"SHOW TASKS IN ACCOUNT",
		).
		withAdditionalSqlCasef(
			"sql_Show_inApplication",
			func(opts *ShowTaskOptions) { opts.In = &ExtendedIn{Application: applicationId} },
			"SHOW TASKS IN APPLICATION %s", applicationId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Show_inApplicationPackage",
			func(opts *ShowTaskOptions) { opts.In = &ExtendedIn{ApplicationPackage: applicationPackageId} },
			"SHOW TASKS IN APPLICATION PACKAGE %s", applicationPackageId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Tasks_sql_Show_StartsWith,
			func(opts *ShowTaskOptions) { opts.StartsWith = new("abc") },
			"SHOW TASKS STARTS WITH 'abc'",
		).
		withModifyAndExpectedSqlf(
			case_Tasks_sql_Show_Limit,
			func(opts *ShowTaskOptions) { opts.Limit = &LimitFrom{Rows: new(10)} },
			"SHOW TASKS LIMIT 10",
		)

	tasksTests.Describe.
		withExpectedSqlf(
			case_Tasks_sql_Describe_basic,
			"DESCRIBE TASK %s", tasksTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		)

	tasksTests.Execute.
		withExpectedSqlf(
			case_Tasks_sql_Execute_basic,
			"EXECUTE TASK %s", tasksTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Execute_all",
			func(opts *ExecuteTaskOptions) { opts.RetryLast = new(true) },
			"EXECUTE TASK %s RETRY LAST", tasksTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		)
}

func TestParseTaskSchedule(t *testing.T) {
	testCases := map[string]struct {
		Schedule             string
		ExpectedTaskSchedule *TaskSchedule
		Error                string
	}{
		"valid schedule: m minutes": {
			Schedule:             "5 m",
			ExpectedTaskSchedule: &TaskSchedule{Minutes: 5},
		},
		"valid schedule: M minutes": {
			Schedule:             "5 M",
			ExpectedTaskSchedule: &TaskSchedule{Minutes: 5},
		},
		"valid schedule: MINUTE minutes": {
			Schedule:             "5 MINUTE",
			ExpectedTaskSchedule: &TaskSchedule{Minutes: 5},
		},
		"valid schedule: MINUTES minutes": {
			Schedule:             "5 MINUTES",
			ExpectedTaskSchedule: &TaskSchedule{Minutes: 5},
		},
		"valid schedule: s seconds": {
			Schedule:             "30 s",
			ExpectedTaskSchedule: &TaskSchedule{Seconds: 30},
		},
		"valid schedule: S seconds": {
			Schedule:             "30 S",
			ExpectedTaskSchedule: &TaskSchedule{Seconds: 30},
		},
		"valid schedule: SECOND seconds": {
			Schedule:             "30 SECOND",
			ExpectedTaskSchedule: &TaskSchedule{Seconds: 30},
		},
		"valid schedule: SECONDS seconds": {
			Schedule:             "30 SECONDS",
			ExpectedTaskSchedule: &TaskSchedule{Seconds: 30},
		},
		"valid schedule: h hours": {
			Schedule:             "2 h",
			ExpectedTaskSchedule: &TaskSchedule{Hours: 2},
		},
		"valid schedule: H hours": {
			Schedule:             "2 H",
			ExpectedTaskSchedule: &TaskSchedule{Hours: 2},
		},
		"valid schedule: HOUR hours": {
			Schedule:             "2 HOUR",
			ExpectedTaskSchedule: &TaskSchedule{Hours: 2},
		},
		"valid schedule: HOURS hours": {
			Schedule:             "2 HOURS",
			ExpectedTaskSchedule: &TaskSchedule{Hours: 2},
		},
		"valid schedule: cron": {
			Schedule:             "USING CRON * * * * * UTC",
			ExpectedTaskSchedule: &TaskSchedule{Cron: "* * * * * UTC"},
		},
		"valid schedule: cron with case sensitive location": {
			Schedule:             "USING CRON * * * * * America/Loc_Angeles",
			ExpectedTaskSchedule: &TaskSchedule{Cron: "* * * * * America/Loc_Angeles"},
		},
		"invalid schedule: wrong schedule format": {
			Schedule:             "SOME SCHEDULE",
			ExpectedTaskSchedule: nil,
			Error:                `strconv.Atoi: parsing "SOME": invalid syntax`,
		},
		"invalid schedule: wrong minutes format": {
			Schedule:             "a5 MINUTE",
			ExpectedTaskSchedule: nil,
			Error:                `strconv.Atoi: parsing "A5": invalid syntax`,
		},
		// currently, cron expressions are not validated (they are on Snowflake level)
		"invalid schedule: wrong cron format": {
			Schedule:             "USING CRON some_cron",
			ExpectedTaskSchedule: &TaskSchedule{Cron: "some_cron"},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			taskSchedule, err := ParseTaskSchedule(tc.Schedule)
			if tc.Error != "" {
				assert.Nil(t, taskSchedule)
				assert.ErrorContains(t, err, tc.Error)
			} else {
				assert.Equal(t, tc.ExpectedTaskSchedule, taskSchedule)
				assert.NoError(t, err)
			}
		})
	}
}

func TestParseTargetCompletionInterval(t *testing.T) {
	valid := map[string]struct {
		Input    string
		Expected *TaskTargetCompletionInterval
	}{
		"valid hours singular": {
			Input:    "1 HOUR",
			Expected: &TaskTargetCompletionInterval{Hours: Pointer(1)},
		},
		"valid hours plural": {
			Input:    "2 HOURS",
			Expected: &TaskTargetCompletionInterval{Hours: Pointer(2)},
		},
		"valid hours - short form": {
			Input:    "2 h",
			Expected: &TaskTargetCompletionInterval{Hours: Pointer(2)},
		},
		"valid minutes singular": {
			Input:    "1 MINUTE",
			Expected: &TaskTargetCompletionInterval{Minutes: Pointer(1)},
		},
		"valid minutes plural": {
			Input:    "10 MINUTES",
			Expected: &TaskTargetCompletionInterval{Minutes: Pointer(10)},
		},
		"valid minutes - short form": {
			Input:    "5 m",
			Expected: &TaskTargetCompletionInterval{Minutes: Pointer(5)},
		},
		"valid seconds singular": {
			Input:    "1 SECOND",
			Expected: &TaskTargetCompletionInterval{Seconds: Pointer(1)},
		},
		"valid seconds plural": {
			Input:    "30 SECONDS",
			Expected: &TaskTargetCompletionInterval{Seconds: Pointer(30)},
		},
		"valid seconds - short form": {
			Input:    "30 s",
			Expected: &TaskTargetCompletionInterval{Seconds: Pointer(30)},
		},
		"valid lowercase": {
			Input:    "5 minutes",
			Expected: &TaskTargetCompletionInterval{Minutes: Pointer(5)},
		},
		"leading/trailing spaces": {
			Input:    " 7 HOURS ",
			Expected: &TaskTargetCompletionInterval{Hours: Pointer(7)},
		},
	}

	for name, tc := range valid {
		t.Run(name, func(t *testing.T) {
			got, err := parseTargetCompletionInterval(tc.Input)
			require.NoError(t, err)
			assert.Equal(t, tc.Expected, got)
		})
	}
	invalid := map[string]struct {
		Input string
		Error string
	}{
		"invalid format: missing value": {
			Input: "MINUTES",
			Error: "invalid task target completion interval format",
		},
		"invalid format: extra parts": {
			Input: "1 HOURS EXTRA",
			Error: "invalid task target completion interval format",
		},
		"invalid value: not a number": {
			Input: "foo HOURS",
			Error: "invalid task target completion interval value",
		},
		"invalid unit: nonsense": {
			Input: "5 CATS",
			Error: "invalid task target completion interval unit",
		},
		"empty input": {
			Input: "",
			Error: "invalid task target completion interval format",
		},
	}
	for name, tc := range invalid {
		t.Run(name, func(t *testing.T) {
			got, err := parseTargetCompletionInterval(tc.Input)
			assert.Empty(t, got)
			assert.ErrorContains(t, err, tc.Error)
		})
	}
}
