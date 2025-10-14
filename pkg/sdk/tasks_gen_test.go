package sdk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTasks_Create(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	sql := "SELECT CURRENT_TIMESTAMP"

	// Minimal valid CreateTaskOptions
	defaultOpts := func() *CreateTaskOptions {
		return &CreateTaskOptions{
			name: id,
			sql:  sql,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateTaskOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: conflicting fields for [opts.OrReplace opts.IfNotExists]", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.IfNotExists = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateTaskOptions", "OrReplace", "IfNotExists"))
	})

	t.Run("validation: exactly one field from [opts.Warehouse.Warehouse opts.Warehouse.UserTaskManagedInitialWarehouseSize] should be present", func(t *testing.T) {
		opts := defaultOpts()
		opts.Warehouse = &CreateTaskWarehouse{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateTaskOptions.Warehouse", "Warehouse", "UserTaskManagedInitialWarehouseSize"))
	})

	t.Run("validation: opts.SessionParameters.SessionParameters should be valid", func(t *testing.T) {
		opts := defaultOpts()
		opts.SessionParameters = &SessionParameters{
			JsonIndent: Int(-1),
		}
		assertOptsInvalidJoinedErrors(t, opts, errIntValue("SessionParameters", "JsonIndent", IntErrGreaterOrEqual, 0))
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "CREATE TASK %s AS %s", id.FullyQualifiedName(), sql)
	})

	t.Run("with initial warehouse size", func(t *testing.T) {
		opts := defaultOpts()
		opts.Warehouse = &CreateTaskWarehouse{
			UserTaskManagedInitialWarehouseSize: Pointer(WarehouseSizeXSmall),
		}
		assertOptsValidAndSQLEquals(t, opts, "CREATE TASK %s USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE = 'XSMALL' AS %s", id.FullyQualifiedName(), sql)
	})

	t.Run("with target completion interval", func(t *testing.T) {
		opts := defaultOpts()
		opts.TargetCompletionInterval = String("10 MINUTES")
		assertOptsValidAndSQLEquals(t, opts, "CREATE TASK %s TARGET_COMPLETION_INTERVAL = '10 MINUTES' AS %s", id.FullyQualifiedName(), sql)
	})

	t.Run("with serverless task min statement size", func(t *testing.T) {
		opts := defaultOpts()
		opts.ServerlessTaskMinStatementSize = Pointer(WarehouseSizeSmall)
		assertOptsValidAndSQLEquals(t, opts, "CREATE TASK %s SERVERLESS_TASK_MIN_STATEMENT_SIZE = 'SMALL' AS %s", id.FullyQualifiedName(), sql)
	})

	t.Run("with serverless task max statement size", func(t *testing.T) {
		opts := defaultOpts()
		opts.ServerlessTaskMaxStatementSize = Pointer(WarehouseSizeLarge)
		assertOptsValidAndSQLEquals(t, opts, "CREATE TASK %s SERVERLESS_TASK_MAX_STATEMENT_SIZE = 'LARGE' AS %s", id.FullyQualifiedName(), sql)
	})

	t.Run("with all serverless task parameters", func(t *testing.T) {
		opts := defaultOpts()
		opts.TargetCompletionInterval = String("5 MINUTES")
		opts.ServerlessTaskMinStatementSize = Pointer(WarehouseSizeXSmall)
		opts.ServerlessTaskMaxStatementSize = Pointer(WarehouseSizeMedium)
		assertOptsValidAndSQLEquals(t, opts, "CREATE TASK %s TARGET_COMPLETION_INTERVAL = '5 MINUTES' SERVERLESS_TASK_MIN_STATEMENT_SIZE = 'XSMALL' SERVERLESS_TASK_MAX_STATEMENT_SIZE = 'MEDIUM' AS %s", id.FullyQualifiedName(), sql)
	})

	t.Run("all options", func(t *testing.T) {
		warehouseId := randomAccountObjectIdentifier()
		otherTaskId := randomSchemaObjectIdentifier()
		tagId := randomSchemaObjectIdentifier()
		finalizerId := randomSchemaObjectIdentifier()
		opts := defaultOpts()

		opts.OrReplace = Bool(true)
		opts.Warehouse = &CreateTaskWarehouse{
			Warehouse: &warehouseId,
		}
		opts.Schedule = String("10 MINUTE")
		opts.Config = String(`$${"output_dir": "/temp/test_directory/", "learning_rate": 0.1}$$`)
		opts.AllowOverlappingExecution = Bool(true)
		opts.SessionParameters = &SessionParameters{
			JsonIndent:  Int(10),
			LockTimeout: Int(5),
		}
		opts.UserTaskTimeoutMs = Int(5)
		opts.SuspendTaskAfterNumFailures = Int(6)
		opts.ErrorIntegration = Pointer(NewAccountObjectIdentifier("some_error_integration"))
		opts.Comment = String("some comment")
		opts.Finalize = &finalizerId
		opts.TaskAutoRetryAttempts = Int(10)
		opts.Tag = []TagAssociation{{
			Name:  tagId,
			Value: "v1",
		}}
		opts.UserTaskMinimumTriggerIntervalInSeconds = Int(10)
		opts.After = []SchemaObjectIdentifier{otherTaskId}
		opts.When = String(`SYSTEM$STREAM_HAS_DATA('MYSTREAM')`)

		assertOptsValidAndSQLEquals(t, opts, "CREATE OR REPLACE TASK %s WAREHOUSE = %s SCHEDULE = '10 MINUTE' CONFIG = $${\"output_dir\": \"/temp/test_directory/\", \"learning_rate\": 0.1}$$ ALLOW_OVERLAPPING_EXECUTION = true JSON_INDENT = 10, LOCK_TIMEOUT = 5 USER_TASK_TIMEOUT_MS = 5 SUSPEND_TASK_AFTER_NUM_FAILURES = 6 ERROR_INTEGRATION = \"some_error_integration\" COMMENT = 'some comment' FINALIZE = %s TASK_AUTO_RETRY_ATTEMPTS = 10 TAG (%s = 'v1') USER_TASK_MINIMUM_TRIGGER_INTERVAL_IN_SECONDS = 10 AFTER %s WHEN SYSTEM$STREAM_HAS_DATA('MYSTREAM') AS SELECT CURRENT_TIMESTAMP", id.FullyQualifiedName(), warehouseId.FullyQualifiedName(), finalizerId.FullyQualifiedName(), tagId.FullyQualifiedName(), otherTaskId.FullyQualifiedName())
	})
}

func TestTasks_CreateOrAlter(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	sql := "SELECT CURRENT_TIMESTAMP"

	// Minimal valid CreateTaskOptions
	defaultOpts := func() *CreateOrAlterTaskOptions {
		return &CreateOrAlterTaskOptions{
			name: id,
			sql:  sql,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateOrAlterTaskOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: exactly one field from [opts.Warehouse.Warehouse opts.Warehouse.UserTaskManagedInitialWarehouseSize] should be present", func(t *testing.T) {
		opts := defaultOpts()
		opts.Warehouse = &CreateTaskWarehouse{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateOrAlterTaskOptions.CreateTaskWarehouse", "Warehouse", "UserTaskManagedInitialWarehouseSize"))
	})

	t.Run("validation: opts.SessionParameters.SessionParameters should be valid", func(t *testing.T) {
		opts := defaultOpts()
		opts.SessionParameters = &SessionParameters{
			JsonIndent: Int(-1),
		}
		assertOptsInvalidJoinedErrors(t, opts, errIntValue("SessionParameters", "JsonIndent", IntErrGreaterOrEqual, 0))
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "CREATE OR ALTER TASK %s AS %s", id.FullyQualifiedName(), sql)
	})

	t.Run("all options", func(t *testing.T) {
		warehouseId := randomAccountObjectIdentifier()
		otherTaskId := randomSchemaObjectIdentifier()
		finalizerId := randomSchemaObjectIdentifier()
		opts := defaultOpts()

		opts.Warehouse = &CreateTaskWarehouse{
			Warehouse: &warehouseId,
		}
		opts.Schedule = String("10 MINUTE")
		opts.Config = String(`$${"output_dir": "/temp/test_directory/", "learning_rate": 0.1}$$`)
		opts.AllowOverlappingExecution = Bool(true)
		opts.UserTaskTimeoutMs = Int(5)
		opts.SessionParameters = &SessionParameters{
			JsonIndent:  Int(10),
			LockTimeout: Int(5),
		}
		opts.SuspendTaskAfterNumFailures = Int(6)
		opts.ErrorIntegration = Pointer(NewAccountObjectIdentifier("some_error_integration"))
		opts.Comment = String("some comment")
		opts.Finalize = &finalizerId
		opts.TaskAutoRetryAttempts = Int(10)
		opts.After = []SchemaObjectIdentifier{otherTaskId}
		opts.When = String(`SYSTEM$STREAM_HAS_DATA('MYSTREAM')`)

		assertOptsValidAndSQLEquals(t, opts, "CREATE OR ALTER TASK %s WAREHOUSE = %s SCHEDULE = '10 MINUTE' CONFIG = $${\"output_dir\": \"/temp/test_directory/\", \"learning_rate\": 0.1}$$ ALLOW_OVERLAPPING_EXECUTION = true USER_TASK_TIMEOUT_MS = 5 JSON_INDENT = 10, LOCK_TIMEOUT = 5 SUSPEND_TASK_AFTER_NUM_FAILURES = 6 ERROR_INTEGRATION = \"some_error_integration\" COMMENT = 'some comment' FINALIZE = %s TASK_AUTO_RETRY_ATTEMPTS = 10 AFTER %s WHEN SYSTEM$STREAM_HAS_DATA('MYSTREAM') AS SELECT CURRENT_TIMESTAMP", id.FullyQualifiedName(), warehouseId.FullyQualifiedName(), finalizerId.FullyQualifiedName(), otherTaskId.FullyQualifiedName())
	})
}

func TestTasks_Clone(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	sourceId := randomSchemaObjectIdentifier()

	// Minimal valid CloneTaskOptions
	defaultOpts := func() *CloneTaskOptions {
		return &CloneTaskOptions{
			name:       id,
			sourceTask: sourceId,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CloneTaskOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: valid identifier for [opts.sourceTask]", func(t *testing.T) {
		opts := defaultOpts()
		opts.sourceTask = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "CREATE TASK %s CLONE %s", id.FullyQualifiedName(), sourceId.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.CopyGrants = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "CREATE OR REPLACE TASK %s CLONE %s COPY GRANTS", id.FullyQualifiedName(), sourceId.FullyQualifiedName())
	})
}

func TestTasks_Alter(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	otherTaskId := randomSchemaObjectIdentifier()

	// Minimal valid AlterTaskOptions
	defaultOpts := func() *AlterTaskOptions {
		return &AlterTaskOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *AlterTaskOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: exactly one field from [opts.Resume opts.Suspend opts.RemoveAfter opts.AddAfter opts.Set opts.Unset opts.SetTags opts.UnsetTags opts.ModifyAs opts.ModifyWhen] should be present", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterTaskOptions", "Resume", "Suspend", "RemoveAfter", "AddAfter", "Set", "Unset", "SetTags", "UnsetTags", "SetFinalize", "UnsetFinalize", "ModifyAs", "ModifyWhen", "RemoveWhen"))
	})

	t.Run("validation: exactly one field from [opts.Resume opts.Suspend opts.RemoveAfter opts.AddAfter opts.Set opts.Unset opts.SetTags opts.UnsetTags opts.ModifyAs opts.ModifyWhen] should be present - more present", func(t *testing.T) {
		opts := defaultOpts()
		opts.Resume = Bool(true)
		opts.Suspend = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterTaskOptions", "Resume", "Suspend", "RemoveAfter", "AddAfter", "Set", "Unset", "SetTags", "UnsetTags", "SetFinalize", "UnsetFinalize", "ModifyAs", "ModifyWhen", "RemoveWhen"))
	})

	t.Run("validation: at least one of the fields [opts.Set.Warehouse opts.Set.UserTaskManagedInitialWarehouseSize opts.Set.Schedule opts.Set.Config opts.Set.AllowOverlappingExecution opts.Set.UserTaskTimeoutMs opts.Set.SuspendTaskAfterNumFailures opts.Set.ErrorIntegration opts.Set.Comment opts.Set.SessionParameters] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &TaskSet{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterTaskOptions.Set", "Warehouse", "UserTaskManagedInitialWarehouseSize", "Schedule", "Config", "AllowOverlappingExecution", "UserTaskTimeoutMs", "SuspendTaskAfterNumFailures", "ErrorIntegration", "Comment", "SessionParameters", "TaskAutoRetryAttempts", "UserTaskMinimumTriggerIntervalInSeconds", "TargetCompletionInterval", "ServerlessTaskMinStatementSize", "ServerlessTaskMaxStatementSize"))
	})

	t.Run("validation: conflicting fields for [opts.Set.Warehouse opts.Set.UserTaskManagedInitialWarehouseSize]", func(t *testing.T) {
		warehouseId := randomAccountObjectIdentifier()
		opts := defaultOpts()
		opts.Set = &TaskSet{}
		opts.Set.Warehouse = &warehouseId
		opts.Set.UserTaskManagedInitialWarehouseSize = Pointer(WarehouseSizeXSmall)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("AlterTaskOptions.Set", "Warehouse", "UserTaskManagedInitialWarehouseSize"))
	})

	t.Run("validation: opts.Set.SessionParameters.SessionParameters should be valid", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &TaskSet{
			SessionParameters: &SessionParameters{
				JsonIndent: Int(-1),
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errIntValue("SessionParameters", "JsonIndent", IntErrGreaterOrEqual, 0))
	})

	t.Run("validation: at least one of the fields [opts.Unset.Warehouse opts.Unset.Schedule opts.Unset.Config opts.Unset.AllowOverlappingExecution opts.Unset.UserTaskTimeoutMs opts.Unset.SuspendTaskAfterNumFailures opts.Unset.ErrorIntegration opts.Unset.Comment opts.Unset.SessionParametersUnset] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &TaskUnset{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterTaskOptions.Unset", "Warehouse", "UserTaskManagedInitialWarehouseSize", "Schedule", "Config", "AllowOverlappingExecution", "UserTaskTimeoutMs", "SuspendTaskAfterNumFailures", "ErrorIntegration", "Comment", "SessionParametersUnset", "TaskAutoRetryAttempts", "UserTaskMinimumTriggerIntervalInSeconds", "TargetCompletionInterval", "ServerlessTaskMinStatementSize", "ServerlessTaskMaxStatementSize"))
	})

	t.Run("validation: opts.Unset.SessionParametersUnset.SessionParametersUnset should be valid", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &TaskUnset{}
		opts.Unset.SessionParametersUnset = &SessionParametersUnset{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("SessionParametersUnset", "AbortDetachedQuery", "ActivePythonProfiler", "Autocommit", "BinaryInputFormat", "BinaryOutputFormat", "ClientEnableLogInfoStatementParameters", "ClientMemoryLimit", "ClientMetadataRequestUseConnectionCtx", "ClientPrefetchThreads", "ClientResultChunkSize", "ClientResultColumnCaseInsensitive", "ClientMetadataUseSessionDatabase", "ClientSessionKeepAlive", "ClientSessionKeepAliveHeartbeatFrequency", "ClientTimestampTypeMapping", "CsvTimestampFormat", "DateInputFormat", "DateOutputFormat", "EnableUnloadPhysicalTypeOptimization", "ErrorOnNondeterministicMerge", "ErrorOnNondeterministicUpdate", "GeographyOutputFormat", "GeometryOutputFormat", "HybridTableLockTimeout", "JdbcTreatDecimalAsInt", "JdbcTreatTimestampNtzAsUtc", "JdbcUseSessionTimezone", "JsonIndent", "JsTreatIntegerAsBigInt", "LockTimeout", "LogLevel", "MultiStatementCount", "NoorderSequenceAsDefault", "OdbcTreatDecimalAsInt", "PythonProfilerModules", "PythonProfilerTargetStage", "QueryTag", "QuotedIdentifiersIgnoreCase", "RowsPerResultset", "S3StageVpceDnsName", "SearchPath", "SimulatedDataSharingConsumer", "StatementQueuedTimeoutInSeconds", "StatementTimeoutInSeconds", "StrictJsonOutput", "TimestampDayIsAlways24h", "TimestampInputFormat", "TimestampLTZOutputFormat", "TimestampNTZOutputFormat", "TimestampOutputFormat", "TimestampTypeMapping", "TimestampTZOutputFormat", "Timezone", "TimeInputFormat", "TimeOutputFormat", "TraceLevel", "TransactionAbortOnError", "TransactionDefaultIsolationLevel", "TwoDigitCenturyStart", "UnsupportedDDLAction", "UseCachedResult", "WeekOfYearPolicy", "WeekStart"))
	})

	t.Run("alter resume", func(t *testing.T) {
		opts := defaultOpts()
		opts.Resume = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "ALTER TASK %s RESUME", id.FullyQualifiedName())
	})

	t.Run("alter suspend", func(t *testing.T) {
		opts := defaultOpts()
		opts.Suspend = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "ALTER TASK %s SUSPEND", id.FullyQualifiedName())
	})

	t.Run("alter remove after", func(t *testing.T) {
		opts := defaultOpts()
		opts.RemoveAfter = []SchemaObjectIdentifier{otherTaskId}
		assertOptsValidAndSQLEquals(t, opts, "ALTER TASK %s REMOVE AFTER %s", id.FullyQualifiedName(), otherTaskId.FullyQualifiedName())
	})

	t.Run("alter add after", func(t *testing.T) {
		opts := defaultOpts()
		opts.AddAfter = []SchemaObjectIdentifier{otherTaskId}
		assertOptsValidAndSQLEquals(t, opts, "ALTER TASK %s ADD AFTER %s", id.FullyQualifiedName(), otherTaskId.FullyQualifiedName())
	})

	t.Run("alter set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &TaskSet{
			Comment: String("some comment"),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER TASK %s SET COMMENT = 'some comment'", id.FullyQualifiedName())
	})

	t.Run("alter set: multiple", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &TaskSet{
			UserTaskTimeoutMs: Int(2000),
			Comment:           String("some comment"),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER TASK %s SET USER_TASK_TIMEOUT_MS = 2000, COMMENT = 'some comment'", id.FullyQualifiedName())
	})

	t.Run("alter set warehouse", func(t *testing.T) {
		warehouseId := randomAccountObjectIdentifier()
		opts := defaultOpts()
		opts.Set = &TaskSet{
			Warehouse: &warehouseId,
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER TASK %s SET WAREHOUSE = %s", id.FullyQualifiedName(), warehouseId.FullyQualifiedName())
	})

	t.Run("alter set session parameter", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &TaskSet{
			SessionParameters: &SessionParameters{
				JsonIndent: Int(15),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER TASK %s SET JSON_INDENT = 15", id.FullyQualifiedName())
	})

	t.Run("alter set target completion interval", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &TaskSet{
			TargetCompletionInterval: String("15 MINUTES"),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER TASK %s SET TARGET_COMPLETION_INTERVAL = '15 MINUTES'", id.FullyQualifiedName())
	})

	t.Run("alter set serverless task min statement size", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &TaskSet{
			ServerlessTaskMinStatementSize: Pointer(WarehouseSizeXSmall),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER TASK %s SET SERVERLESS_TASK_MIN_STATEMENT_SIZE = 'XSMALL'", id.FullyQualifiedName())
	})

	t.Run("alter set serverless task max statement size", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &TaskSet{
			ServerlessTaskMaxStatementSize: Pointer(WarehouseSizeXLarge),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER TASK %s SET SERVERLESS_TASK_MAX_STATEMENT_SIZE = 'XLARGE'", id.FullyQualifiedName())
	})

	t.Run("alter set multiple serverless task parameters", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &TaskSet{
			TargetCompletionInterval:        String("20 MINUTES"),
			ServerlessTaskMinStatementSize:  Pointer(WarehouseSizeSmall),
			ServerlessTaskMaxStatementSize:  Pointer(WarehouseSizeLarge),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER TASK %s SET TARGET_COMPLETION_INTERVAL = '20 MINUTES', SERVERLESS_TASK_MIN_STATEMENT_SIZE = 'SMALL', SERVERLESS_TASK_MAX_STATEMENT_SIZE = 'LARGE'", id.FullyQualifiedName())
	})

	t.Run("alter unset", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &TaskUnset{
			Comment: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER TASK %s UNSET COMMENT", id.FullyQualifiedName())
	})

	t.Run("alter unset: multiple", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &TaskUnset{
			UserTaskTimeoutMs: Bool(true),
			Comment:           Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER TASK %s UNSET USER_TASK_TIMEOUT_MS, COMMENT", id.FullyQualifiedName())
	})

	t.Run("alter unset target completion interval", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &TaskUnset{
			TargetCompletionInterval: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER TASK %s UNSET TARGET_COMPLETION_INTERVAL", id.FullyQualifiedName())
	})

	t.Run("alter unset serverless task min statement size", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &TaskUnset{
			ServerlessTaskMinStatementSize: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER TASK %s UNSET SERVERLESS_TASK_MIN_STATEMENT_SIZE", id.FullyQualifiedName())
	})

	t.Run("alter unset serverless task max statement size", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &TaskUnset{
			ServerlessTaskMaxStatementSize: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER TASK %s UNSET SERVERLESS_TASK_MAX_STATEMENT_SIZE", id.FullyQualifiedName())
	})

	t.Run("alter unset multiple serverless task parameters", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &TaskUnset{
			TargetCompletionInterval:       Bool(true),
			ServerlessTaskMinStatementSize: Bool(true),
			ServerlessTaskMaxStatementSize: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER TASK %s UNSET TARGET_COMPLETION_INTERVAL, SERVERLESS_TASK_MIN_STATEMENT_SIZE, SERVERLESS_TASK_MAX_STATEMENT_SIZE", id.FullyQualifiedName())
	})

	t.Run("alter set tags", func(t *testing.T) {
		opts := defaultOpts()
		opts.SetTags = []TagAssociation{
			{
				Name:  NewAccountObjectIdentifier("tag1"),
				Value: "value1",
			},
			{
				Name:  NewAccountObjectIdentifier("tag2"),
				Value: "value2",
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER TASK %s SET TAG "tag1" = 'value1', "tag2" = 'value2'`, id.FullyQualifiedName())
	})

	t.Run("alter unset tags", func(t *testing.T) {
		opts := defaultOpts()
		opts.UnsetTags = []ObjectIdentifier{
			NewAccountObjectIdentifier("tag1"),
			NewAccountObjectIdentifier("tag2"),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER TASK %s UNSET TAG "tag1", "tag2"`, id.FullyQualifiedName())
	})

	t.Run("alter modify as", func(t *testing.T) {
		opts := defaultOpts()
		opts.ModifyAs = String("new as")
		assertOptsValidAndSQLEquals(t, opts, "ALTER TASK %s MODIFY AS new as", id.FullyQualifiedName())
	})

	t.Run("alter modify when", func(t *testing.T) {
		opts := defaultOpts()
		opts.ModifyWhen = String("new when")
		assertOptsValidAndSQLEquals(t, opts, "ALTER TASK %s MODIFY WHEN new when", id.FullyQualifiedName())
	})

	t.Run("alter set finalize", func(t *testing.T) {
		opts := defaultOpts()
		finalizeId := randomSchemaObjectIdentifier()
		opts.SetFinalize = &finalizeId
		assertOptsValidAndSQLEquals(t, opts, "ALTER TASK %s SET FINALIZE = %s", id.FullyQualifiedName(), finalizeId.FullyQualifiedName())
	})

	t.Run("alter unset finalize", func(t *testing.T) {
		opts := defaultOpts()
		opts.UnsetFinalize = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "ALTER TASK %s UNSET FINALIZE", id.FullyQualifiedName())
	})

	t.Run("alter remove when", func(t *testing.T) {
		opts := defaultOpts()
		opts.RemoveWhen = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "ALTER TASK %s REMOVE WHEN", id.FullyQualifiedName())
	})
}

func TestTasks_Drop(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	// Minimal valid DropTaskOptions
	defaultOpts := func() *DropTaskOptions {
		return &DropTaskOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DropTaskOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "DROP TASK %s", id.FullyQualifiedName())
	})
}

func TestTasks_Show(t *testing.T) {
	// Minimal valid ShowTaskOptions
	defaultOpts := func() *ShowTaskOptions {
		return &ShowTaskOptions{}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *ShowTaskOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "SHOW TASKS")
	})

	t.Run("in application", func(t *testing.T) {
		opts := defaultOpts()
		id := randomAccountObjectIdentifier()
		opts.In = &ExtendedIn{
			Application: id,
		}
		assertOptsValidAndSQLEquals(t, opts, "SHOW TASKS IN APPLICATION %s", id.FullyQualifiedName())
	})

	t.Run("in application package", func(t *testing.T) {
		opts := defaultOpts()
		id := randomAccountObjectIdentifier()
		opts.In = &ExtendedIn{
			ApplicationPackage: id,
		}
		assertOptsValidAndSQLEquals(t, opts, "SHOW TASKS IN APPLICATION PACKAGE %s", id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.Terse = Bool(true)
		opts.Like = &Like{
			Pattern: String("myaccount"),
		}
		opts.In = &ExtendedIn{
			In: In{
				Account: Bool(true),
			},
		}
		opts.StartsWith = String("abc")
		opts.RootOnly = Bool(true)
		opts.Limit = &LimitFrom{Rows: Int(10)}
		assertOptsValidAndSQLEquals(t, opts, "SHOW TERSE TASKS LIKE 'myaccount' IN ACCOUNT STARTS WITH 'abc' ROOT ONLY LIMIT 10")
	})
}

func TestTasks_Describe(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	// Minimal valid DescribeTaskOptions
	defaultOpts := func() *DescribeTaskOptions {
		return &DescribeTaskOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DescribeTaskOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "DESCRIBE TASK %s", id.FullyQualifiedName())
	})
}

func TestTasks_Execute(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	// Minimal valid ExecuteTaskOptions
	defaultOpts := func() *ExecuteTaskOptions {
		return &ExecuteTaskOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *ExecuteTaskOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "EXECUTE TASK %s", id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.RetryLast = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "EXECUTE TASK %s RETRY LAST", id.FullyQualifiedName())
	})
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
			Error:                "invalid schedule format",
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
				assert.EqualValues(t, tc.ExpectedTaskSchedule, taskSchedule)
				assert.NoError(t, err)
			}
		})
	}
}
