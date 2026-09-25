package defs

import (
	"slices"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/parameterdefs"
)

var (
	onSchema               = []parameterdefs.ParameterLevel{parameterdefs.ParameterLevelAccount, parameterdefs.ParameterLevelAccountExt, parameterdefs.ParameterLevelDatabase, parameterdefs.ParameterLevelSchema}
	onTable                = []parameterdefs.ParameterLevel{parameterdefs.ParameterLevelAccount, parameterdefs.ParameterLevelAccountExt, parameterdefs.ParameterLevelDatabase, parameterdefs.ParameterLevelSchema, parameterdefs.ParameterLevelTable}
	onTask                 = []parameterdefs.ParameterLevel{parameterdefs.ParameterLevelAccount, parameterdefs.ParameterLevelAccountExt, parameterdefs.ParameterLevelDatabase, parameterdefs.ParameterLevelSchema, parameterdefs.ParameterLevelTask}
	onFunctionAndProcedure = []parameterdefs.ParameterLevel{parameterdefs.ParameterLevelAccount, parameterdefs.ParameterLevelAccountExt, parameterdefs.ParameterLevelDatabase, parameterdefs.ParameterLevelSchema, parameterdefs.ParameterLevelFunction, parameterdefs.ParameterLevelProcedure}
	onLog                  = []parameterdefs.ParameterLevel{parameterdefs.ParameterLevelAccount, parameterdefs.ParameterLevelAccountExt, parameterdefs.ParameterLevelDatabase, parameterdefs.ParameterLevelSchema, parameterdefs.ParameterLevelProject, parameterdefs.ParameterLevelProcedure, parameterdefs.ParameterLevelFunction, parameterdefs.ParameterLevelTable, parameterdefs.ParameterLevelTask, parameterdefs.ParameterLevelService}
)

var (
	Catalog = parameterdefs.ParameterDef{
		SqlName:     "CATALOG",
		Kind:        g.KindOfT[sdkcommons.AccountObjectIdentifier](),
		Levels:      onTable,
		Description: "The database parameter that specifies the default catalog to use for Iceberg tables.",
	}
	DataRetentionTimeInDays = parameterdefs.ParameterDef{
		SqlName:     "DATA_RETENTION_TIME_IN_DAYS",
		Kind:        g.KindInt,
		Levels:      onTable,
		Description: "Specifies the number of days for which Time Travel actions (CLONE and UNDROP) can be performed on the database, as well as specifying the default Time Travel retention time for all schemas created in the database. For more details, see [Understanding & Using Time Travel](https://docs.snowflake.com/en/user-guide/data-time-travel).",
	}
	DefaultDdlCollation = parameterdefs.ParameterDef{
		SqlName:     "DEFAULT_DDL_COLLATION",
		Kind:        g.KindOfT[sdkcommons.StringAllowEmpty](),
		Levels:      onTable,
		Description: "Specifies a default collation specification for all schemas and tables added to the database. It can be overridden on schema or table level. For more information, see [collation specification](https://docs.snowflake.com/en/sql-reference/collation#label-collation-specification).",
	}
	DefaultNotebookComputePoolCpu = parameterdefs.ParameterDef{
		SqlName:     "DEFAULT_NOTEBOOK_COMPUTE_POOL_CPU",
		Kind:        g.KindString,
		Levels:      onSchema,
		Description: "Sets the preferred CPU compute pool used for Notebooks on CPU Container Runtime.",
	}
	DefaultNotebookComputePoolGpu = parameterdefs.ParameterDef{
		SqlName:     "DEFAULT_NOTEBOOK_COMPUTE_POOL_GPU",
		Kind:        g.KindString,
		Levels:      onSchema,
		Description: "Sets the preferred GPU compute pool used for Notebooks on GPU Container Runtime.",
	}
	EnableConsoleOutput = parameterdefs.ParameterDef{
		SqlName:     "ENABLE_CONSOLE_OUTPUT",
		Kind:        g.KindBool,
		Levels:      []parameterdefs.ParameterLevel{parameterdefs.ParameterLevelAccountExt, parameterdefs.ParameterLevelDatabase, parameterdefs.ParameterLevelSchema, parameterdefs.ParameterLevelTable},
		Description: "If true, enables stdout/stderr fast path logging for anonymous stored procedures.",
	}
	ExternalVolume = parameterdefs.ParameterDef{
		SqlName:     "EXTERNAL_VOLUME",
		Kind:        g.KindOfT[sdkcommons.AccountObjectIdentifier](),
		Levels:      onTable,
		Description: "The database parameter that specifies the default external volume to use for Iceberg tables.",
	}
	LogEventLevel = parameterdefs.ParameterDef{
		SqlName:     "LOG_EVENT_LEVEL",
		Kind:        g.KindOfT[sdkcommons.LogLevel](),
		Levels:      append(slices.Clone(onLog), parameterdefs.ParameterLevelSession, parameterdefs.ParameterLevelUser),
		Description: "Specifies the severity level of log events (rows with record type EVENT) that should be ingested and made available in the active event table. Log events at the specified level (and at more severe levels) are ingested.",
	}
	LogLevel = parameterdefs.ParameterDef{
		SqlName:     "LOG_LEVEL",
		Kind:        g.KindOfT[sdkcommons.LogLevel](),
		Levels:      append(slices.Clone(onLog), parameterdefs.ParameterLevelSession, parameterdefs.ParameterLevelUser),
		Description: "Specifies the severity level of messages that should be ingested and made available in the active event table. Messages at the specified level (and at more severe levels) are ingested.",
	}
	MaxDataExtensionTimeInDays = parameterdefs.ParameterDef{
		SqlName:     "MAX_DATA_EXTENSION_TIME_IN_DAYS",
		Kind:        g.KindInt,
		Levels:      onTable,
		Description: "Object parameter that specifies the maximum number of days for which Snowflake can extend the data retention period for tables in the database to prevent streams on the tables from becoming stale.",
	}
	QuotedIdentifiersIgnoreCase = parameterdefs.ParameterDef{
		SqlName:     "QUOTED_IDENTIFIERS_IGNORE_CASE",
		Kind:        g.KindBool,
		Levels:      append(slices.Clone(onTable), parameterdefs.ParameterLevelSession, parameterdefs.ParameterLevelUser),
		Description: "If true, the case of quoted identifiers is ignored.",
	}
	ReplaceInvalidCharacters = parameterdefs.ParameterDef{
		SqlName:     "REPLACE_INVALID_CHARACTERS",
		Kind:        g.KindBool,
		Levels:      onTable,
		Description: "Specifies whether to replace invalid UTF-8 characters with the Unicode replacement character in query results for an Iceberg table. You can only set this parameter for tables that use an external Iceberg catalog.",
	}
	StorageSerializationPolicy = parameterdefs.ParameterDef{
		SqlName:     "STORAGE_SERIALIZATION_POLICY",
		Kind:        g.KindOfT[sdkcommons.StorageSerializationPolicy](),
		Levels:      onTable,
		Description: "The storage serialization policy for Iceberg tables that use Snowflake as the catalog. COMPATIBLE: Snowflake performs encoding and compression of data files that ensures interoperability with third-party compute engines. OPTIMIZED: Snowflake performs encoding and compression of data files that ensures the best table performance within Snowflake.",
	}
	SuspendTaskAfterNumFailures = parameterdefs.ParameterDef{
		SqlName:     "SUSPEND_TASK_AFTER_NUM_FAILURES",
		Kind:        g.KindInt,
		Levels:      onTask,
		Description: "How many times a task must fail in a row before it is automatically suspended. 0 disables auto-suspending.",
	}
	TaskAutoRetryAttempts = parameterdefs.ParameterDef{
		SqlName:     "TASK_AUTO_RETRY_ATTEMPTS",
		Kind:        g.KindInt,
		Levels:      onTask,
		Description: "Maximum automatic retries allowed for a user task.",
	}
	TraceLevel = parameterdefs.ParameterDef{
		SqlName:     "TRACE_LEVEL",
		Kind:        g.KindOfT[sdkcommons.TraceLevel](),
		Levels:      append(slices.Clone(onFunctionAndProcedure), parameterdefs.ParameterLevelSession, parameterdefs.ParameterLevelUser),
		Description: "Controls how trace events are ingested into the event table.",
	}
	UserTaskManagedInitialWarehouseSize = parameterdefs.ParameterDef{
		SqlName:     "USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE",
		Kind:        g.KindOfT[sdkcommons.WarehouseSize](),
		Levels:      onTask,
		Description: "The initial size of warehouse to use for managed warehouses in the absence of history.",
	}
	UserTaskMinimumTriggerIntervalInSeconds = parameterdefs.ParameterDef{
		SqlName:     "USER_TASK_MINIMUM_TRIGGER_INTERVAL_IN_SECONDS",
		Kind:        g.KindInt,
		Levels:      onTask,
		Description: "Minimum amount of time between Triggered Task executions in seconds.",
	}
	UserTaskTimeoutMs = parameterdefs.ParameterDef{
		SqlName:     "USER_TASK_TIMEOUT_MS",
		Kind:        g.KindInt,
		Levels:      onTask,
		Description: "User task execution timeout in milliseconds.",
	}
)

var AllParameters = []parameterdefs.ParameterDef{
	Catalog,
	DataRetentionTimeInDays,
	DefaultDdlCollation,
	DefaultNotebookComputePoolCpu,
	DefaultNotebookComputePoolGpu,
	EnableConsoleOutput,
	ExternalVolume,
	LogEventLevel,
	LogLevel,
	MaxDataExtensionTimeInDays,
	QuotedIdentifiersIgnoreCase,
	ReplaceInvalidCharacters,
	StorageSerializationPolicy,
	SuspendTaskAfterNumFailures,
	TaskAutoRetryAttempts,
	TraceLevel,
	UserTaskManagedInitialWarehouseSize,
	UserTaskMinimumTriggerIntervalInSeconds,
	UserTaskTimeoutMs,
}

// ParameterDefsForLevel returns the catalog entries assigned to a generator consumer level.
//
// TODO [SNOW-4160077]: Remove or relocate this temporary public API when the
// parameter generator becomes the single source of truth for SDK, resource,
// and data source parameter handling.
func ParameterDefsForLevel(level parameterdefs.ParameterLevel) []parameterdefs.ParameterDef {
	return collections.Filter(AllParameters, func(p parameterdefs.ParameterDef) bool {
		return slices.Contains(p.Levels, level)
	})
}
