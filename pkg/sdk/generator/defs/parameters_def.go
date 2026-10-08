package defs

import (
	"slices"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/parameterdefs"
)

var (
	onAccount                  = []parameterdefs.ParameterLevel{parameterdefs.ParameterLevelAccount, parameterdefs.ParameterLevelAccountExt}
	onAccountExt               = []parameterdefs.ParameterLevel{parameterdefs.ParameterLevelAccountExt}
	onSchema                   = []parameterdefs.ParameterLevel{parameterdefs.ParameterLevelAccount, parameterdefs.ParameterLevelAccountExt, parameterdefs.ParameterLevelDatabase, parameterdefs.ParameterLevelSchema}
	onTable                    = []parameterdefs.ParameterLevel{parameterdefs.ParameterLevelAccount, parameterdefs.ParameterLevelAccountExt, parameterdefs.ParameterLevelDatabase, parameterdefs.ParameterLevelSchema, parameterdefs.ParameterLevelTable}
	onTableAndIcebergTable     = []parameterdefs.ParameterLevel{parameterdefs.ParameterLevelAccount, parameterdefs.ParameterLevelAccountExt, parameterdefs.ParameterLevelDatabase, parameterdefs.ParameterLevelSchema, parameterdefs.ParameterLevelTable, parameterdefs.ParameterLevelIcebergTable}
	onAccountAndIcebergTable   = []parameterdefs.ParameterLevel{parameterdefs.ParameterLevelAccount, parameterdefs.ParameterLevelAccountExt, parameterdefs.ParameterLevelIcebergTable}
	onIcebergTableOnly         = []parameterdefs.ParameterLevel{parameterdefs.ParameterLevelIcebergTable}
	onTask                     = []parameterdefs.ParameterLevel{parameterdefs.ParameterLevelAccount, parameterdefs.ParameterLevelAccountExt, parameterdefs.ParameterLevelDatabase, parameterdefs.ParameterLevelSchema, parameterdefs.ParameterLevelTask}
	onFunctionAndProcedure     = []parameterdefs.ParameterLevel{parameterdefs.ParameterLevelAccount, parameterdefs.ParameterLevelAccountExt, parameterdefs.ParameterLevelDatabase, parameterdefs.ParameterLevelSchema, parameterdefs.ParameterLevelFunction, parameterdefs.ParameterLevelProcedure}
	onLog                      = []parameterdefs.ParameterLevel{parameterdefs.ParameterLevelAccount, parameterdefs.ParameterLevelAccountExt, parameterdefs.ParameterLevelDatabase, parameterdefs.ParameterLevelSchema, parameterdefs.ParameterLevelProject, parameterdefs.ParameterLevelProcedure, parameterdefs.ParameterLevelFunction, parameterdefs.ParameterLevelTable, parameterdefs.ParameterLevelTask}
	onUser                     = append(slices.Clone(onAccount), parameterdefs.ParameterLevelSession, parameterdefs.ParameterLevelUser)
	onUserOnly                 = append(slices.Clone(onAccount), parameterdefs.ParameterLevelUser)
	onLogUser                  = append(slices.Clone(onLog), parameterdefs.ParameterLevelSession, parameterdefs.ParameterLevelUser)
	onWarehouse                = []parameterdefs.ParameterLevel{parameterdefs.ParameterLevelAccount, parameterdefs.ParameterLevelAccountExt, parameterdefs.ParameterLevelWarehouse, parameterdefs.ParameterLevelWarehouseInteractive}
	onWarehouseAll             = []parameterdefs.ParameterLevel{parameterdefs.ParameterLevelAccount, parameterdefs.ParameterLevelAccountExt, parameterdefs.ParameterLevelWarehouse, parameterdefs.ParameterLevelWarehouseAdaptive, parameterdefs.ParameterLevelWarehouseInteractive}
	onWarehouseInteractiveOnly = []parameterdefs.ParameterLevel{parameterdefs.ParameterLevelWarehouseInteractive}
)

const parameterTypeSnowflakeDefault = "sdk.ParameterTypeSnowflakeDefault"

var (
	AbortDetachedQuery = parameterdefs.ParameterDef{
		SqlName:      "ABORT_DETACHED_QUERY",
		Kind:         g.KindBool,
		Levels:       append(slices.Clone(onUser), parameterdefs.ParameterLevelTask),
		Description:  "Specifies the action that Snowflake performs for in-progress queries if connectivity is lost due to abrupt termination of a session (e.g. network outage, browser termination, service interruption).",
		DefaultValue: "false",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	ActivePythonProfiler = parameterdefs.ParameterDef{
		SqlName:     "ACTIVE_PYTHON_PROFILER",
		Kind:        g.KindOfT[sdkcommons.ActivePythonProfiler](),
		Levels:      onAccount,
		Description: "Sets the profiler to use for the session when [profiling Python handler code](https://docs.snowflake.com/en/developer-guide/stored-procedure/python/procedure-python-profiler).",
	}
	AllowBindValuesAccess = parameterdefs.ParameterDef{
		SqlName:     "ALLOW_BIND_VALUES_ACCESS",
		Kind:        g.KindBool,
		Levels:      onAccount,
		Description: "Used to allow clients to access bind variable values.",
	}
	AllowClientMfaCaching = parameterdefs.ParameterDef{
		SqlName:     "ALLOW_CLIENT_MFA_CACHING",
		Kind:        g.KindBool,
		Levels:      onAccount,
		Description: "Specifies whether an MFA token can be saved in the client-side operating system keystore to promote continuous, secure connectivity without users needing to respond to an MFA prompt at the start of each connection attempt to Snowflake. For details and the list of supported Snowflake-provided clients, see [Using MFA token caching to minimize the number of prompts during authentication — optional.](https://docs.snowflake.com/en/user-guide/security-mfa.html#label-mfa-token-caching)",
	}
	AllowIdToken = parameterdefs.ParameterDef{
		SqlName:     "ALLOW_ID_TOKEN",
		Kind:        g.KindBool,
		Levels:      onAccount,
		Description: "Specifies whether a connection token can be saved in the client-side operating system keystore to promote continuous, secure connectivity without users needing to enter login credentials at the start of each connection attempt to Snowflake. For details and the list of supported Snowflake-provided clients, see [Using connection caching to minimize the number of prompts for authentication — optional.](https://docs.snowflake.com/en/user-guide/admin-security-fed-auth-use.html#label-browser-based-sso-connection-caching)",
	}
	AllowRowTimestamp = parameterdefs.ParameterDef{
		SqlName:      "ALLOW_ROW_TIMESTAMP",
		Kind:         g.KindBool,
		Levels:       onIcebergTableOnly,
		Description:  "Specifies whether a hidden row-level timestamp column is maintained on a Snowflake-managed Apache Iceberg™ table.",
		DefaultValue: "false",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	AllowedSpcsWorkloadTypes = parameterdefs.ParameterDef{
		SqlName:     "ALLOWED_SPCS_WORKLOAD_TYPES",
		Kind:        g.KindString,
		Levels:      onAccount,
		Description: "Used to specify the workload types that are allowed in your account to deploy to Snowpark Container Services.",
	}
	Autocommit = parameterdefs.ParameterDef{
		SqlName:      "AUTOCOMMIT",
		Kind:         g.KindBool,
		Levels:       append(slices.Clone(onUser), parameterdefs.ParameterLevelTask),
		Description:  "Specifies whether autocommit is enabled for the session. Autocommit determines whether a DML statement, when executed without an active transaction, is automatically committed after the statement successfully completes. For more information, see [Transactions](https://docs.snowflake.com/en/sql-reference/transactions).",
		DefaultValue: "true",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	BaseLocationPrefix = parameterdefs.ParameterDef{
		SqlName:     "BASE_LOCATION_PREFIX",
		Kind:        g.KindString,
		Levels:      onAccount,
		Description: "Specifies a prefix for Snowflake to use in the write path for Snowflake-managed Apache Iceberg™ tables. For more information, see [data and metadata directories for Iceberg tables](https://docs.snowflake.com/en/user-guide/tables-iceberg-storage.html#label-tables-iceberg-configure-external-volume-base-location).",
	}
	BinaryInputFormat = parameterdefs.ParameterDef{
		SqlName:      "BINARY_INPUT_FORMAT",
		Kind:         g.KindOfT[sdkcommons.BinaryInputFormat](),
		Levels:       append(slices.Clone(onUser), parameterdefs.ParameterLevelTask),
		Description:  "The format of VARCHAR values passed as input to VARCHAR-to-BINARY conversion functions. For more information, see [Binary input and output](https://docs.snowflake.com/en/sql-reference/binary-input-output).",
		DefaultValue: "sdk.BinaryInputFormatHex",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	BinaryOutputFormat = parameterdefs.ParameterDef{
		SqlName:      "BINARY_OUTPUT_FORMAT",
		Kind:         g.KindOfT[sdkcommons.BinaryOutputFormat](),
		Levels:       append(slices.Clone(onUser), parameterdefs.ParameterLevelTask),
		Description:  "The format for VARCHAR values returned as output by BINARY-to-VARCHAR conversion functions. For more information, see [Binary input and output](https://docs.snowflake.com/en/sql-reference/binary-input-output).",
		DefaultValue: "sdk.BinaryOutputFormatHex",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	Catalog = parameterdefs.ParameterDef{
		SqlName:      "CATALOG",
		Kind:         g.KindOfT[sdkcommons.AccountObjectIdentifier](),
		Levels:       onTableAndIcebergTable,
		Description:  "The parameter that specifies the default catalog to use for Iceberg tables.",
		DefaultValue: "",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	CatalogSync = parameterdefs.ParameterDef{
		SqlName:      "CATALOG_SYNC",
		Kind:         g.KindString,
		Levels:       onAccountAndIcebergTable,
		Description:  "Specifies the name of your catalog integration for [Snowflake Open Catalog](https://other-docs.snowflake.com/en/opencatalog/overview). Snowflake syncs tables that use the specified catalog integration with your Snowflake Open Catalog account. For more information, see [Sync a Snowflake-managed table with Snowflake Open Catalog](https://docs.snowflake.com/en/user-guide/tables-iceberg-open-catalog-sync).",
		DefaultValue: "",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	ClientEnableLogInfoStatementParameters = parameterdefs.ParameterDef{
		SqlName:     "CLIENT_ENABLE_LOG_INFO_STATEMENT_PARAMETERS",
		Kind:        g.KindBool,
		Levels:      onAccount,
		Description: "Enables users to log the data values bound to [PreparedStatements](https://docs.snowflake.com/en/developer-guide/jdbc/jdbc-api.html#label-jdbc-api-preparedstatement).",
	}
	ClientEncryptionKeySize = parameterdefs.ParameterDef{
		SqlName:     "CLIENT_ENCRYPTION_KEY_SIZE",
		Kind:        g.KindInt,
		Levels:      onAccount,
		Description: "Specifies the AES encryption key size, in bits, used by Snowflake to encrypt/decrypt files stored on internal stages (for loading/unloading data) when you use the SNOWFLAKE_FULL encryption type.",
	}
	ClientMemoryLimit = parameterdefs.ParameterDef{
		SqlName:      "CLIENT_MEMORY_LIMIT",
		Kind:         g.KindInt,
		Levels:       append(slices.Clone(onUser), parameterdefs.ParameterLevelTask),
		Description:  "Parameter that specifies the maximum amount of memory the JDBC driver or ODBC driver should use for the result set from queries (in MB).",
		DefaultValue: "1536",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	ClientMetadataRequestUseConnectionCtx = parameterdefs.ParameterDef{
		SqlName:      "CLIENT_METADATA_REQUEST_USE_CONNECTION_CTX",
		Kind:         g.KindBool,
		Levels:       append(slices.Clone(onUser), parameterdefs.ParameterLevelTask),
		Description:  "For specific ODBC functions and JDBC methods, this parameter can change the default search scope from all databases/schemas to the current database/schema. The narrower search typically returns fewer rows and executes more quickly.",
		DefaultValue: "false",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	ClientMetadataUseSessionDatabase = parameterdefs.ParameterDef{
		SqlName:     "CLIENT_METADATA_USE_SESSION_DATABASE",
		Kind:        g.KindBool,
		Levels:      onAccount,
		Description: "For specific ODBC functions and JDBC methods, this parameter can change the default search scope from all databases to the current database. The narrower search typically returns fewer rows and executes more quickly ([more details on the usage](https://docs.snowflake.com/en/sql-reference/parameters#client-metadata-use-session-database)).",
	}
	ClientPrefetchThreads = parameterdefs.ParameterDef{
		SqlName:      "CLIENT_PREFETCH_THREADS",
		Kind:         g.KindInt,
		Levels:       append(slices.Clone(onUser), parameterdefs.ParameterLevelTask),
		Description:  "Parameter that specifies the number of threads used by the client to pre-fetch large result sets. The driver will attempt to honor the parameter value, but defines the minimum and maximum values (depending on your system’s resources) to improve performance.",
		DefaultValue: "4",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	ClientResultChunkSize = parameterdefs.ParameterDef{
		SqlName:      "CLIENT_RESULT_CHUNK_SIZE",
		Kind:         g.KindInt,
		Levels:       append(slices.Clone(onUser), parameterdefs.ParameterLevelTask),
		Description:  "Parameter that specifies the maximum size of each set (or chunk) of query results to download (in MB). The JDBC driver downloads query results in chunks.",
		DefaultValue: "160",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	ClientResultColumnCaseInsensitive = parameterdefs.ParameterDef{
		SqlName:      "CLIENT_RESULT_COLUMN_CASE_INSENSITIVE",
		Kind:         g.KindBool,
		Levels:       append(slices.Clone(onUser), parameterdefs.ParameterLevelTask),
		Description:  "Parameter that indicates whether to match column name case-insensitively in ResultSet.get* methods in JDBC.",
		DefaultValue: "false",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	ClientSessionKeepAlive = parameterdefs.ParameterDef{
		SqlName:      "CLIENT_SESSION_KEEP_ALIVE",
		Kind:         g.KindBool,
		Levels:       append(slices.Clone(onUser), parameterdefs.ParameterLevelTask),
		Description:  "Parameter that indicates whether to force a user to log in again after a period of inactivity in the session.",
		DefaultValue: "false",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	ClientSessionKeepAliveHeartbeatFrequency = parameterdefs.ParameterDef{
		SqlName:      "CLIENT_SESSION_KEEP_ALIVE_HEARTBEAT_FREQUENCY",
		Kind:         g.KindInt,
		Levels:       append(slices.Clone(onUser), parameterdefs.ParameterLevelTask),
		Description:  "Number of seconds in-between client attempts to update the token for the session.",
		DefaultValue: "3600",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	ClientTimestampTypeMapping = parameterdefs.ParameterDef{
		SqlName:      "CLIENT_TIMESTAMP_TYPE_MAPPING",
		Kind:         g.KindOfT[sdkcommons.ClientTimestampTypeMapping](),
		Levels:       append(slices.Clone(onUser), parameterdefs.ParameterLevelTask),
		Description:  "Specifies the [TIMESTAMP_* variation](https://docs.snowflake.com/en/sql-reference/data-types-datetime.html#label-datatypes-timestamp-variations) to use when binding timestamp variables for JDBC or ODBC applications that use the bind API to load data.",
		DefaultValue: "sdk.ClientTimestampTypeMappingLtz",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	CortexCodeCliDailyEstCreditLimitPerUser = parameterdefs.ParameterDef{
		SqlName:     "CORTEX_CODE_CLI_DAILY_EST_CREDIT_LIMIT_PER_USER",
		Kind:        g.KindInt,
		Levels:      onAccount,
		Description: "Daily estimated credit limit per user for Cortex Code CLI usage. Set to `-1` for the default (unlimited), `0` to block usage, or a positive value to cap a user's estimated credit usage over a rolling 24-hour window. For more information, see [Cortex Code credit usage limits](https://docs.snowflake.com/en/user-guide/cortex-code/credit-usage-limit).",
	}
	CortexCodeDesktopDailyEstCreditLimitPerUser = parameterdefs.ParameterDef{
		SqlName:     "CORTEX_CODE_DESKTOP_DAILY_EST_CREDIT_LIMIT_PER_USER",
		Kind:        g.KindInt,
		Levels:      onAccount,
		Description: "Daily estimated credit limit per user for Cortex Code Desktop usage. Set to `-1` for the default (unlimited), `0` to block usage, or a positive value to cap a user's estimated credit usage over a rolling 24-hour window. For more information, see [Cortex Code credit usage limits](https://docs.snowflake.com/en/user-guide/cortex-code/credit-usage-limit).",
	}
	CortexCodeSnowsightDailyEstCreditLimitPerUser = parameterdefs.ParameterDef{
		SqlName:     "CORTEX_CODE_SNOWSIGHT_DAILY_EST_CREDIT_LIMIT_PER_USER",
		Kind:        g.KindInt,
		Levels:      onAccount,
		Description: "Daily estimated credit limit per user for Cortex Code Snowsight usage. Set to `-1` for the default (unlimited), `0` to block usage, or a positive value to cap a user's estimated credit usage over a rolling 24-hour window. For more information, see [Cortex Code credit usage limits](https://docs.snowflake.com/en/user-guide/cortex-code/credit-usage-limit).",
	}
	CortexEnabledCrossRegion = parameterdefs.ParameterDef{
		SqlName:     "CORTEX_ENABLED_CROSS_REGION",
		Kind:        g.KindString,
		Levels:      onAccount,
		Description: "Specifies the regions where an inference request may be processed in case the request cannot be processed in the region where request is originally placed. Specifying DISABLED disables cross-region inferencing. For examples and details, see [Cross-region inference](https://docs.snowflake.com/en/user-guide/snowflake-cortex/cross-region-inference).",
	}
	CortexModelsAllowlist = parameterdefs.ParameterDef{
		SqlName:     "CORTEX_MODELS_ALLOWLIST",
		Kind:        g.KindString,
		Levels:      onAccount,
		Description: "Specifies the models that users in the account can access. Use this parameter to allowlist models for all users in the account. If you need to provide specific users with access beyond what you’ve specified in the allowlist, use role-based access control instead. For more information, see [Model allowlist](https://docs.snowflake.com/en/user-guide/snowflake-cortex/aisql.html#label-cortex-llm-allowlist).",
	}
	CsvTimestampFormat = parameterdefs.ParameterDef{
		SqlName:     "CSV_TIMESTAMP_FORMAT",
		Kind:        g.KindString,
		Levels:      onAccount,
		Description: "Specifies the format for TIMESTAMP values in CSV files downloaded from Snowsight. If this parameter is not set, [TIMESTAMP_LTZ_OUTPUT_FORMAT](https://docs.snowflake.com/en/sql-reference/parameters#label-timestamp-ltz-output-format) will be used for TIMESTAMP_LTZ values, [TIMESTAMP_TZ_OUTPUT_FORMAT](https://docs.snowflake.com/en/sql-reference/parameters#label-timestamp-tz-output-format) will be used for TIMESTAMP_TZ and [TIMESTAMP_NTZ_OUTPUT_FORMAT](https://docs.snowflake.com/en/sql-reference/parameters#label-timestamp-ntz-output-format) for TIMESTAMP_NTZ values. For more information, see [Date and time input and output formats](https://docs.snowflake.com/en/sql-reference/date-time-input-output) or [Download your query results](https://docs.snowflake.com/en/user-guide/ui-snowsight-query.html#label-snowsight-download-query-results).",
	}
	DataMetricSchedule = parameterdefs.ParameterDef{
		SqlName:      "DATA_METRIC_SCHEDULE",
		Kind:         g.KindString,
		Levels:       onAccountAndIcebergTable,
		Description:  "Specifies the schedule to run the data metric functions associated to the table. All data metric functions on the table or view follow the same schedule.",
		DefaultValue: "60 MINUTES",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	DataRetentionTimeInDays = parameterdefs.ParameterDef{
		SqlName:      "DATA_RETENTION_TIME_IN_DAYS",
		Kind:         g.KindInt,
		Levels:       append(slices.Clone(onTableAndIcebergTable), parameterdefs.ParameterLevelHybridTable),
		Description:  "Specifies the number of days for which Time Travel actions (CLONE and UNDROP) can be performed on the database, as well as specifying the default Time Travel retention time for all schemas created in the database. For more details, see [Understanding & Using Time Travel](https://docs.snowflake.com/en/user-guide/data-time-travel).",
		DefaultValue: "1",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	DateInputFormat = parameterdefs.ParameterDef{
		SqlName:      "DATE_INPUT_FORMAT",
		Kind:         g.KindString,
		Levels:       append(slices.Clone(onUser), parameterdefs.ParameterLevelTask),
		Description:  "Specifies the input format for the DATE data type. For more information, see [Date and time input and output formats](https://docs.snowflake.com/en/sql-reference/date-time-input-output).",
		DefaultValue: "AUTO",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	DateOutputFormat = parameterdefs.ParameterDef{
		SqlName:      "DATE_OUTPUT_FORMAT",
		Kind:         g.KindString,
		Levels:       append(slices.Clone(onUser), parameterdefs.ParameterLevelTask),
		Description:  "Specifies the display format for the DATE data type. For more information, see [Date and time input and output formats](https://docs.snowflake.com/en/sql-reference/date-time-input-output).",
		DefaultValue: "YYYY-MM-DD",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	DefaultDbtVersion = parameterdefs.ParameterDef{
		SqlName:     "DEFAULT_DBT_VERSION",
		Kind:        g.KindString,
		Levels:      onAccount,
		Description: "Used to set the default version for all future dbt project objects created in an account.",
	}
	DefaultDdlCollation = parameterdefs.ParameterDef{
		SqlName:      "DEFAULT_DDL_COLLATION",
		Kind:         g.KindOfT[sdkcommons.StringAllowEmpty](),
		Levels:       onTableAndIcebergTable,
		Description:  "Specifies a default collation specification for all schemas and tables added to the database. It can be overridden on schema or table level. For more information, see [collation specification](https://docs.snowflake.com/en/sql-reference/collation#label-collation-specification).",
		DefaultValue: "",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	DefaultNotebookComputePoolCpu = parameterdefs.ParameterDef{
		SqlName:      "DEFAULT_NOTEBOOK_COMPUTE_POOL_CPU",
		Kind:         g.KindString,
		Levels:       onSchema,
		Description:  "Sets the preferred CPU compute pool used for Notebooks on CPU Container Runtime.",
		DefaultValue: "SYSTEM_COMPUTE_POOL_CPU",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	DefaultNotebookComputePoolGpu = parameterdefs.ParameterDef{
		SqlName:      "DEFAULT_NOTEBOOK_COMPUTE_POOL_GPU",
		Kind:         g.KindString,
		Levels:       onSchema,
		Description:  "Sets the preferred GPU compute pool used for Notebooks on GPU Container Runtime.",
		DefaultValue: "SYSTEM_COMPUTE_POOL_GPU",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	DefaultNullOrdering = parameterdefs.ParameterDef{
		SqlName:     "DEFAULT_NULL_ORDERING",
		Kind:        g.KindOfT[sdkcommons.DefaultNullOrdering](),
		Levels:      onAccount,
		Description: "Specifies the default ordering of NULL values in a result set ([more details](https://docs.snowflake.com/en/sql-reference/parameters#default-null-ordering)).",
	}
	DefaultStreamlitComputePool = parameterdefs.ParameterDef{
		SqlName:     "DEFAULT_STREAMLIT_COMPUTE_POOL",
		Kind:        g.KindString,
		Levels:      onAccount,
		Description: "Specifies the name of the default compute pool to use when creating container-runtime Streamlit apps.",
	}
	DefaultStreamlitNotebookWarehouse = parameterdefs.ParameterDef{
		SqlName:     "DEFAULT_STREAMLIT_NOTEBOOK_WAREHOUSE",
		Kind:        g.KindOfT[sdkcommons.AccountObjectIdentifier](),
		Levels:      onAccount,
		Description: "Specifies the name of the default warehouse to use when creating a notebook.",
	}
	DisableUiDownloadButton = parameterdefs.ParameterDef{
		SqlName:     "DISABLE_UI_DOWNLOAD_BUTTON",
		Kind:        g.KindBool,
		Levels:      onAccount,
		Description: "Controls whether users in an account see a button to download data in Snowsight or the Classic Console, such as a table returned from running a query in a worksheet. If the button to download is hidden in Snowsight or the Classic Console, users can still download or export data using [third-party software](https://docs.snowflake.com/en/user-guide/ecosystem).",
	}
	DisableUserPrivilegeGrants = parameterdefs.ParameterDef{
		SqlName:     "DISABLE_USER_PRIVILEGE_GRANTS",
		Kind:        g.KindBool,
		Levels:      onAccount,
		Description: "Controls whether users in an account can grant privileges directly to other users. Disabling user privilege grants (that is, setting DISABLE_USER_PRIVILEGE_GRANTS to TRUE) does not affect existing grants to users. Existing grants to users continue to confer privileges to those users. For more information, see [GRANT <privileges> … TO USER](https://docs.snowflake.com/en/sql-reference/sql/grant-privilege-user).",
	}
	DisallowedSpcsWorkloadTypes = parameterdefs.ParameterDef{
		SqlName:     "DISALLOWED_SPCS_WORKLOAD_TYPES",
		Kind:        g.KindString,
		Levels:      onAccount,
		Description: "Used to specify the workload types that are disallowed in your account to deploy to Snowpark Container Services.",
	}
	EnableAutomaticSensitiveDataClassificationLog = parameterdefs.ParameterDef{
		SqlName:     "ENABLE_AUTOMATIC_SENSITIVE_DATA_CLASSIFICATION_LOG",
		Kind:        g.KindBool,
		Levels:      onAccount,
		Description: "Controls whether events from [automatic sensitive data classification](https://docs.snowflake.com/en/user-guide/classify-auto) are logged in the user event table.",
	}
	EnableBudgetEventLogging = parameterdefs.ParameterDef{
		SqlName:     "ENABLE_BUDGET_EVENT_LOGGING",
		Kind:        g.KindBool,
		Levels:      onAccount,
		Description: "Controls whether events from budgets are logged to the event table.",
	}
	EnableConsoleOutput = parameterdefs.ParameterDef{
		SqlName:      "ENABLE_CONSOLE_OUTPUT",
		Kind:         g.KindBool,
		Levels:       []parameterdefs.ParameterLevel{parameterdefs.ParameterLevelAccountExt, parameterdefs.ParameterLevelDatabase, parameterdefs.ParameterLevelSchema, parameterdefs.ParameterLevelTable, parameterdefs.ParameterLevelFunction, parameterdefs.ParameterLevelProcedure},
		Description:  "If true, enables stdout/stderr fast path logging for anonymous stored procedures.",
		DefaultValue: "false",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	EnableCortexAnalyst = parameterdefs.ParameterDef{
		SqlName:     "ENABLE_CORTEX_ANALYST",
		Kind:        g.KindBool,
		Levels:      onAccount,
		Description: "Controls whether Cortex Analyst is enabled for the account.",
	}
	EnableDataCompaction = parameterdefs.ParameterDef{
		SqlName:      "ENABLE_DATA_COMPACTION",
		Kind:         g.KindBool,
		Levels:       onAccountAndIcebergTable,
		Description:  "Specifies whether Snowflake should enable data compaction on Snowflake-managed Apache Iceberg™ tables.",
		DefaultValue: "true",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	EnableEgressCostOptimizer = parameterdefs.ParameterDef{
		SqlName:     "ENABLE_EGRESS_COST_OPTIMIZER",
		Kind:        g.KindBool,
		Levels:      onAccount,
		Description: "Enables or disables the Listing Cross-cloud auto-fulfillment Egress cost optimizer.",
	}
	EnableGetDdlUseDataTypeAlias = parameterdefs.ParameterDef{
		SqlName:     "ENABLE_GET_DDL_USE_DATA_TYPE_ALIAS",
		Kind:        g.KindBool,
		Levels:      onAccount,
		Description: "Specifies whether the output returned by the GET_DDL function contains data type synonyms specified in the original DDL statement. Data type synonyms are also called data type aliases.",
	}
	EnableIcebergMergeOnRead = parameterdefs.ParameterDef{
		SqlName:      "ENABLE_ICEBERG_MERGE_ON_READ",
		Kind:         g.KindBool,
		Levels:       onAccountAndIcebergTable,
		Description:  "Specifies whether to enable merge-on-read behavior for Snowflake-managed Apache Iceberg™ tables.",
		DefaultValue: "true",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	EnableIdentifierFirstLogin = parameterdefs.ParameterDef{
		SqlName:     "ENABLE_IDENTIFIER_FIRST_LOGIN",
		Kind:        g.KindBool,
		Levels:      onAccount,
		Description: "Determines the login flow for users. When enabled, Snowflake prompts users for their username or email address before presenting authentication methods. For details, see [Identifier-first login](https://docs.snowflake.com/en/user-guide/identifier-first-login).",
	}
	EnableInternalStagesPrivatelink = parameterdefs.ParameterDef{
		SqlName:     "ENABLE_INTERNAL_STAGES_PRIVATELINK",
		Kind:        g.KindBool,
		Levels:      onAccount,
		Description: "Specifies whether the [SYSTEM$GET_PRIVATELINK_CONFIG](https://docs.snowflake.com/en/sql-reference/functions/system_get_privatelink_config) function returns the private-internal-stages key in the query result. The corresponding value in the query result is used during the configuration process for private connectivity to internal stages.",
	}
	EnableNotebookCreationInPersonalDb = parameterdefs.ParameterDef{
		SqlName:     "ENABLE_NOTEBOOK_CREATION_IN_PERSONAL_DB",
		Kind:        g.KindBool,
		Levels:      onAccount,
		Description: "Used to enable or disable private notebooks on a Snowflake account.",
	}
	EnablePerAccountAppServicePrivatelinkUrl = parameterdefs.ParameterDef{
		SqlName:     "ENABLE_PER_ACCOUNT_APP_SERVICE_PRIVATELINK_URL",
		Kind:        g.KindBool,
		Levels:      onAccount,
		Description: "Specifies whether Snowflake generates a per-account private connectivity URL for the Snowflake App Service. For more information, see [Snowflake documentation](https://docs.snowflake.com/en/sql-reference/parameters#enable-per-account-app-service-privatelink-url).",
	}
	EnablePersonalDatabase = parameterdefs.ParameterDef{
		SqlName: "ENABLE_PERSONAL_DATABASE",
		Kind:    g.KindBool,
		Levels:  onAccountExt,
	}
	EnableSpcsBlockStorageSnowflakeFullEncryptionEnforcement = parameterdefs.ParameterDef{
		SqlName:     "ENABLE_SPCS_BLOCK_STORAGE_SNOWFLAKE_FULL_ENCRYPTION_ENFORCEMENT",
		Kind:        g.KindBool,
		Levels:      onAccount,
		Description: "Used to enable enforcement of SNOWFLAKE_FULL encryption for Snowpark Container Services block-storage volumes and snapshots.",
	}
	EnableTagPropagationEventLogging = parameterdefs.ParameterDef{
		SqlName:     "ENABLE_TAG_PROPAGATION_EVENT_LOGGING",
		Kind:        g.KindBool,
		Levels:      onAccount,
		Description: "Controls whether Snowflake collects telemetry data for tag propagation.",
	}
	EnableTriSecretAndRekeyOptOutForImageRepository = parameterdefs.ParameterDef{
		SqlName:     "ENABLE_TRI_SECRET_AND_REKEY_OPT_OUT_FOR_IMAGE_REPOSITORY",
		Kind:        g.KindBool,
		Levels:      onAccount,
		Description: "Specifies choice for the [image repository](https://docs.snowflake.com/en/developer-guide/snowpark-container-services/working-with-registry-repository.html#label-registry-and-repository-image-repository) to opt out of Tri-Secret Secure and [Periodic rekeying](https://docs.snowflake.com/en/user-guide/security-encryption-manage.html#label-periodic-rekeying).",
	}
	EnableTriSecretAndRekeyOptOutForSpcsBlockStorage = parameterdefs.ParameterDef{
		SqlName:     "ENABLE_TRI_SECRET_AND_REKEY_OPT_OUT_FOR_SPCS_BLOCK_STORAGE",
		Kind:        g.KindBool,
		Levels:      onAccount,
		Description: "Specifies the choice for the [Snowpark Container Services block storage volume](https://docs.snowflake.com/en/developer-guide/snowpark-container-services/block-storage-volume) to opt out of Tri-Secret Secure and [Periodic rekeying](https://docs.snowflake.com/en/user-guide/security-encryption-manage.html#label-periodic-rekeying).",
	}
	EnableUnhandledExceptionsReporting = parameterdefs.ParameterDef{
		SqlName:     "ENABLE_UNHANDLED_EXCEPTIONS_REPORTING",
		Kind:        g.KindBool,
		Levels:      onAccount,
		Description: "Specifies whether Snowflake may capture – in an event table – log messages or trace event data for unhandled exceptions in procedure or UDF handler code. For more information, see [Capturing messages from unhandled exceptions](https://docs.snowflake.com/en/developer-guide/logging-tracing/unhandled-exception-messages).",
	}
	EnableUnloadPhysicalTypeOptimization = parameterdefs.ParameterDef{
		SqlName:      "ENABLE_UNLOAD_PHYSICAL_TYPE_OPTIMIZATION",
		Kind:         g.KindBool,
		Levels:       append(slices.Clone(onUser), parameterdefs.ParameterLevelTask),
		Description:  "Specifies whether to set the schema for unloaded Parquet files based on the logical column data types (i.e. the types in the unload SQL query or source table) or on the unloaded column values (i.e. the smallest data types and precision that support the values in the output columns of the unload SQL statement or source table).",
		DefaultValue: "true",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	EnableUnredactedQuerySyntaxError = parameterdefs.ParameterDef{
		SqlName:      "ENABLE_UNREDACTED_QUERY_SYNTAX_ERROR",
		Kind:         g.KindBool,
		Levels:       onUserOnly,
		Description:  "Controls whether query text is redacted if a SQL query fails due to a syntax or parsing error. If FALSE, the content of a failed query is redacted in the views, pages, and functions that provide a query history. Only users with a role that is granted or inherits the AUDIT privilege can set the ENABLE_UNREDACTED_QUERY_SYNTAX_ERROR parameter. When using the ALTER USER command to set the parameter to TRUE for a particular user, modify the user that you want to see the query text, not the user who executed the query (if those are different users).",
		DefaultValue: "false",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	EnableUnredactedSecureObjectError = parameterdefs.ParameterDef{
		SqlName:     "ENABLE_UNREDACTED_SECURE_OBJECT_ERROR",
		Kind:        g.KindBool,
		Levels:      onAccount,
		Description: "Controls whether error messages related to secure objects are redacted in metadata. For more information, see [Secure objects: Redaction of information in error messages](https://docs.snowflake.com/en/release-notes/bcr-bundles/un-bundled/bcr-1858). Only users with a role that is granted or inherits the AUDIT privilege can set the ENABLE_UNREDACTED_SECURE_OBJECT_ERROR parameter. When using the ALTER USER command to set the parameter to TRUE for a particular user, modify the user that you want to see the redacted error messages in metadata, not the user who caused the error.",
	}
	EnforceNetworkRulesForInternalStages = parameterdefs.ParameterDef{
		SqlName:     "ENFORCE_NETWORK_RULES_FOR_INTERNAL_STAGES",
		Kind:        g.KindBool,
		Levels:      onAccount,
		Description: "Specifies whether a network policy that uses network rules can restrict access to AWS internal stages. This parameter has no effect on network policies that do not use network rules. This account-level parameter affects both account-level and user-level network policies. For details about using network policies and network rules to restrict access to AWS internal stages, including the use of this parameter, see [Protecting internal stages on AWS](https://docs.snowflake.com/en/user-guide/network-policies.html#label-network-policies-rules-stages).",
	}
	ErrorOnNondeterministicMerge = parameterdefs.ParameterDef{
		SqlName:      "ERROR_ON_NONDETERMINISTIC_MERGE",
		Kind:         g.KindBool,
		Levels:       append(slices.Clone(onUser), parameterdefs.ParameterLevelTask),
		Description:  "Specifies whether to return an error when the [MERGE](https://docs.snowflake.com/en/sql-reference/sql/merge) command is used to update or delete a target row that joins multiple source rows and the system cannot determine the action to perform on the target row.",
		DefaultValue: "true",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	ErrorOnNondeterministicUpdate = parameterdefs.ParameterDef{
		SqlName:      "ERROR_ON_NONDETERMINISTIC_UPDATE",
		Kind:         g.KindBool,
		Levels:       append(slices.Clone(onUser), parameterdefs.ParameterLevelTask),
		Description:  "Specifies whether to return an error when the [UPDATE](https://docs.snowflake.com/en/sql-reference/sql/update) command is used to update a target row that joins multiple source rows and the system cannot determine the action to perform on the target row.",
		DefaultValue: "false",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	EventTable = parameterdefs.ParameterDef{
		SqlName:      "EVENT_TABLE",
		Kind:         g.KindOfT[sdkcommons.SchemaObjectIdentifier](),
		Levels:       append(slices.Clone(onAccount), parameterdefs.ParameterLevelOpenflowDeployment),
		Description:  "Specifies the name of the event table for logging messages from stored procedures and UDFs contained by the object with which the event table is associated. Associating an event table with a database is available in [Enterprise Edition or higher](https://docs.snowflake.com/en/user-guide/intro-editions).",
		DefaultValue: "",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	ExternalOauthAddPrivilegedRolesToBlockedList = parameterdefs.ParameterDef{
		SqlName:     "EXTERNAL_OAUTH_ADD_PRIVILEGED_ROLES_TO_BLOCKED_LIST",
		Kind:        g.KindBool,
		Levels:      onAccount,
		Description: "Determines whether the ACCOUNTADMIN, ORGADMIN, GLOBALORGADMIN, and SECURITYADMIN roles can be used as the primary role when creating a Snowflake session based on the access token from the External OAuth authorization server.",
	}
	ExternalVolume = parameterdefs.ParameterDef{
		SqlName:      "EXTERNAL_VOLUME",
		Kind:         g.KindOfT[sdkcommons.AccountObjectIdentifier](),
		Levels:       onTableAndIcebergTable,
		Description:  "The parameter that specifies the default external volume to use for Iceberg tables.",
		DefaultValue: "",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	FallbackWarehouse = parameterdefs.ParameterDef{
		SqlName:      "FALLBACK_WAREHOUSE",
		Kind:         g.KindOfT[sdkcommons.AccountObjectIdentifier](),
		Levels:       onWarehouseInteractiveOnly,
		Description:  "If not null, specifies the warehouse to use as a fallback for statements that timed out.",
		DefaultValue: "",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	GeographyOutputFormat = parameterdefs.ParameterDef{
		SqlName:      "GEOGRAPHY_OUTPUT_FORMAT",
		Kind:         g.KindOfT[sdkcommons.GeographyOutputFormat](),
		Levels:       append(slices.Clone(onUser), parameterdefs.ParameterLevelTask),
		Description:  "Display format for [GEOGRAPHY values](https://docs.snowflake.com/en/sql-reference/data-types-geospatial.html#label-data-types-geography).",
		DefaultValue: "sdk.GeographyOutputFormatGeoJSON",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	GeometryOutputFormat = parameterdefs.ParameterDef{
		SqlName:      "GEOMETRY_OUTPUT_FORMAT",
		Kind:         g.KindOfT[sdkcommons.GeometryOutputFormat](),
		Levels:       append(slices.Clone(onUser), parameterdefs.ParameterLevelTask),
		Description:  "Display format for [GEOMETRY values](https://docs.snowflake.com/en/sql-reference/data-types-geospatial.html#label-data-types-geometry).",
		DefaultValue: "sdk.GeometryOutputFormatGeoJSON",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	HybridTableLockTimeout = parameterdefs.ParameterDef{
		SqlName:     "HYBRID_TABLE_LOCK_TIMEOUT",
		Kind:        g.KindInt,
		Levels:      onAccount,
		Description: "Number of seconds to wait while trying to acquire row-level locks on a hybrid table, before timing out and aborting the statement.",
	}
	IcebergMergeOnReadBehavior = parameterdefs.ParameterDef{
		SqlName:      "ICEBERG_MERGE_ON_READ_BEHAVIOR",
		Kind:         IcebergTableIcebergMergeOnReadBehaviorEnumDef.Kind(),
		Levels:       onIcebergTableOnly,
		Description:  "Specifies the merge-on-read behavior for Snowflake-managed Apache Iceberg™ tables.",
		DefaultValue: `sdk.IcebergTableIcebergMergeOnReadBehavior("auto")`,
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	IcebergVersionDefault = parameterdefs.ParameterDef{
		SqlName:     "ICEBERG_VERSION_DEFAULT",
		Kind:        g.KindInt,
		Levels:      onAccount,
		Description: "Specifies the Iceberg specification version to conform to when creating new Snowflake-managed Iceberg tables.",
	}
	InitialReplicationSizeLimitInTb = parameterdefs.ParameterDef{
		SqlName:     "INITIAL_REPLICATION_SIZE_LIMIT_IN_TB",
		Kind:        g.KindFloat,
		Levels:      onAccount,
		Description: "Sets the maximum estimated size limit for the initial replication of a primary database to a secondary database (in TB). Set this parameter on any account that stores a secondary database. This size limit helps prevent accounts from accidentally incurring large database replication charges. To remove the size limit, set the value to 0.0. It is required to pass numbers with scale of at least 1 (e.g. 20.5, 32.25, 33.333, etc.).",
	}
	JdbcTreatDecimalAsInt = parameterdefs.ParameterDef{
		SqlName:      "JDBC_TREAT_DECIMAL_AS_INT",
		Kind:         g.KindBool,
		Levels:       onUser,
		Description:  "Specifies how JDBC processes columns that have a scale of zero (0).",
		DefaultValue: "true",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	JdbcTreatTimestampNtzAsUtc = parameterdefs.ParameterDef{
		SqlName:      "JDBC_TREAT_TIMESTAMP_NTZ_AS_UTC",
		Kind:         g.KindBool,
		Levels:       append(slices.Clone(onUser), parameterdefs.ParameterLevelTask),
		Description:  "Specifies how JDBC processes TIMESTAMP_NTZ values ([more details](https://docs.snowflake.com/en/sql-reference/parameters#jdbc-treat-timestamp-ntz-as-utc)).",
		DefaultValue: "false",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	JdbcUseSessionTimezone = parameterdefs.ParameterDef{
		SqlName:      "JDBC_USE_SESSION_TIMEZONE",
		Kind:         g.KindBool,
		Levels:       append(slices.Clone(onUser), parameterdefs.ParameterLevelTask),
		Description:  "Specifies whether the JDBC Driver uses the time zone of the JVM or the time zone of the session (specified by the [TIMEZONE](https://docs.snowflake.com/en/sql-reference/parameters#label-timezone) parameter) for the getDate(), getTime(), and getTimestamp() methods of the ResultSet class.",
		DefaultValue: "true",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	JsTreatIntegerAsBigint = parameterdefs.ParameterDef{
		SqlName:     "JS_TREAT_INTEGER_AS_BIGINT",
		Kind:        g.KindBool,
		Levels:      onAccount,
		Description: "Specifies how the Snowflake Node.js Driver processes numeric columns that have a scale of zero (0), for example INTEGER or NUMBER(p, 0).",
	}
	JsonIndent = parameterdefs.ParameterDef{
		SqlName:      "JSON_INDENT",
		Kind:         g.KindInt,
		Levels:       append(slices.Clone(onUser), parameterdefs.ParameterLevelTask),
		Description:  "Specifies the number of blank spaces to indent each new element in JSON output in the session. Also specifies whether to insert newline characters after each element.",
		DefaultValue: "2",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	ListingAutoFulfillmentReplicationRefreshSchedule = parameterdefs.ParameterDef{
		SqlName:     "LISTING_AUTO_FULFILLMENT_REPLICATION_REFRESH_SCHEDULE",
		Kind:        g.KindString,
		Levels:      onAccount,
		Description: "Sets the time interval used to refresh the application package based data products to other regions.",
	}
	LockTimeout = parameterdefs.ParameterDef{
		SqlName:      "LOCK_TIMEOUT",
		Kind:         g.KindInt,
		Levels:       append(slices.Clone(onUser), parameterdefs.ParameterLevelTask),
		Description:  "Number of seconds to wait while trying to lock a resource, before timing out and aborting the statement.",
		DefaultValue: "43200",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	LogEventLevel = parameterdefs.ParameterDef{
		SqlName:      "LOG_EVENT_LEVEL",
		Kind:         g.KindOfT[sdkcommons.LogLevel](),
		Levels:       append(slices.Clone(onLog), parameterdefs.ParameterLevelIcebergTable, parameterdefs.ParameterLevelSession, parameterdefs.ParameterLevelUser),
		Description:  "Specifies the severity level of log events (rows with record type EVENT) that should be ingested and made available in the active event table. Log events at the specified level (and at more severe levels) are ingested.",
		DefaultValue: "sdk.LogLevelOff",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	LogLevel = parameterdefs.ParameterDef{
		SqlName:      "LOG_LEVEL",
		Kind:         g.KindOfT[sdkcommons.LogLevel](),
		Levels:       onLogUser,
		Description:  "Specifies the severity level of messages that should be ingested and made available in the active event table. Messages at the specified level (and at more severe levels) are ingested.",
		DefaultValue: "sdk.LogLevelOff",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	MaxConcurrencyLevel = parameterdefs.ParameterDef{
		SqlName:      "MAX_CONCURRENCY_LEVEL",
		Kind:         g.KindInt,
		Levels:       onWarehouse,
		Description:  "Specifies the concurrency level for SQL statements (that is, queries and DML) executed by a warehouse.",
		DefaultValue: "8",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	MaxDataExtensionTimeInDays = parameterdefs.ParameterDef{
		SqlName:      "MAX_DATA_EXTENSION_TIME_IN_DAYS",
		Kind:         g.KindInt,
		Levels:       append(slices.Clone(onTableAndIcebergTable), parameterdefs.ParameterLevelHybridTable),
		Description:  "Object parameter that specifies the maximum number of days for which Snowflake can extend the data retention period for tables in the database to prevent streams on the tables from becoming stale.",
		DefaultValue: "14",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	MetricLevel = parameterdefs.ParameterDef{
		SqlName:      "METRIC_LEVEL",
		Kind:         g.KindOfT[sdkcommons.MetricLevel](),
		Levels:       []parameterdefs.ParameterLevel{parameterdefs.ParameterLevelAccount, parameterdefs.ParameterLevelAccountExt, parameterdefs.ParameterLevelFunction, parameterdefs.ParameterLevelProcedure},
		Description:  "Controls how metrics data is ingested into the event table. For more information about metric levels, see [Setting levels for logging, metrics, and tracing](https://docs.snowflake.com/en/developer-guide/logging-tracing/telemetry-levels).",
		DefaultValue: "sdk.MetricLevelNone",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	MinDataRetentionTimeInDays = parameterdefs.ParameterDef{
		SqlName:     "MIN_DATA_RETENTION_TIME_IN_DAYS",
		Kind:        g.KindInt,
		Levels:      onAccount,
		Description: "Minimum number of days for which Snowflake retains historical data for performing Time Travel actions (SELECT, CLONE, UNDROP) on an object. If a minimum number of days for data retention is set on an account, the data retention period for an object is determined by MAX([DATA_RETENTION_TIME_IN_DAYS](https://docs.snowflake.com/en/sql-reference/parameters#label-data-retention-time-in-days), MIN_DATA_RETENTION_TIME_IN_DAYS).",
	}
	MultiStatementCount = parameterdefs.ParameterDef{
		SqlName:      "MULTI_STATEMENT_COUNT",
		Kind:         g.KindInt,
		Levels:       append(slices.Clone(onUser), parameterdefs.ParameterLevelTask),
		Description:  "Number of statements to execute when using the multi-statement capability.",
		DefaultValue: "1",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	NetworkPolicy = parameterdefs.ParameterDef{
		SqlName:      "NETWORK_POLICY",
		Kind:         g.KindOfT[sdkcommons.AccountObjectIdentifier](),
		Levels:       onUserOnly,
		Description:  "Specifies the network policy to enforce for your account. Network policies enable restricting access to your account based on users’ IP address. For more details, see [Controlling network traffic with network policies](https://docs.snowflake.com/en/user-guide/network-policies).",
		DefaultValue: "",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	NoorderSequenceAsDefault = parameterdefs.ParameterDef{
		SqlName:      "NOORDER_SEQUENCE_AS_DEFAULT",
		Kind:         g.KindBool,
		Levels:       append(slices.Clone(onUser), parameterdefs.ParameterLevelTask),
		Description:  "Specifies whether the ORDER or NOORDER property is set by default when you create a new sequence or add a new table column. The ORDER and NOORDER properties determine whether or not the values are generated for the sequence or auto-incremented column in [increasing or decreasing order](https://docs.snowflake.com/en/user-guide/querying-sequences.html#label-querying-sequences-increasing-values).",
		DefaultValue: "true",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	OauthAddPrivilegedRolesToBlockedList = parameterdefs.ParameterDef{
		SqlName:     "OAUTH_ADD_PRIVILEGED_ROLES_TO_BLOCKED_LIST",
		Kind:        g.KindBool,
		Levels:      onAccount,
		Description: "Determines whether the ACCOUNTADMIN, ORGADMIN, GLOBALORGADMIN, and SECURITYADMIN roles can be used as the primary role when creating a Snowflake session based on the access token from Snowflake’s authorization server.",
	}
	OdbcTreatDecimalAsInt = parameterdefs.ParameterDef{
		SqlName:      "ODBC_TREAT_DECIMAL_AS_INT",
		Kind:         g.KindBool,
		Levels:       append(slices.Clone(onUser), parameterdefs.ParameterLevelTask),
		Description:  "Specifies how ODBC processes columns that have a scale of zero (0).",
		DefaultValue: "false",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	OptimizeDataLayout = parameterdefs.ParameterDef{
		SqlName:      "OPTIMIZE_DATA_LAYOUT",
		Kind:         g.KindBool,
		Levels:       onIcebergTableOnly,
		Description:  "Specifies whether Snowflake should optimize the data layout (e.g. file sizes) of a Snowflake-managed Apache Iceberg™ table.",
		DefaultValue: "true",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	PeriodicDataRekeying = parameterdefs.ParameterDef{
		SqlName:     "PERIODIC_DATA_REKEYING",
		Kind:        g.KindBool,
		Levels:      onAccount,
		Description: "It enables/disables re-encryption of table data with new keys on a yearly basis to provide additional levels of data protection ([more details](https://docs.snowflake.com/en/sql-reference/parameters#periodic-data-rekeying)).",
	}
	PipeExecutionPaused = parameterdefs.ParameterDef{
		SqlName:     "PIPE_EXECUTION_PAUSED",
		Kind:        g.KindBool,
		Levels:      []parameterdefs.ParameterLevel{parameterdefs.ParameterLevelAccount, parameterdefs.ParameterLevelAccountExt, parameterdefs.ParameterLevelSchema},
		Description: "Specifies whether to pause a running pipe, primarily in preparation for transferring ownership of the pipe to a different role.",
	}
	PreventLoadFromInlineUrl = parameterdefs.ParameterDef{
		SqlName: "PREVENT_LOAD_FROM_INLINE_URL",
		Kind:    g.KindBool,
		Levels:  onAccountExt,
	}
	PreventUnloadToInlineUrl = parameterdefs.ParameterDef{
		SqlName:     "PREVENT_UNLOAD_TO_INLINE_URL",
		Kind:        g.KindBool,
		Levels:      onAccount,
		Description: "Specifies whether to prevent ad hoc data unload operations to external cloud storage locations (that is, [COPY INTO location](https://docs.snowflake.com/en/sql-reference/sql/copy-into-location) statements that specify the cloud storage URL and access settings directly in the statement). For an example, see [Unloading data from a table directly to files in an external location](https://docs.snowflake.com/en/sql-reference/sql/copy-into-location.html#label-copy-into-location-ad-hoc).",
	}
	PreventUnloadToInternalStages = parameterdefs.ParameterDef{
		SqlName:      "PREVENT_UNLOAD_TO_INTERNAL_STAGES",
		Kind:         g.KindBool,
		Levels:       onUserOnly,
		Description:  "Specifies whether to prevent data unload operations to internal (Snowflake) stages using [COPY INTO location](https://docs.snowflake.com/en/sql-reference/sql/copy-into-location) statements.",
		DefaultValue: "false",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	PythonProfilerModules = parameterdefs.ParameterDef{
		SqlName:     "PYTHON_PROFILER_MODULES",
		Kind:        g.KindString,
		Levels:      onAccount,
		Description: "Specifies the list of Python modules to include in a report when [profiling Python handler code](https://docs.snowflake.com/en/developer-guide/stored-procedure/python/procedure-python-profiler).",
	}
	PythonProfilerTargetStage = parameterdefs.ParameterDef{
		SqlName:     "PYTHON_PROFILER_TARGET_STAGE",
		Kind:        g.KindOfT[sdkcommons.SchemaObjectIdentifier](),
		Levels:      onAccount,
		Description: "Specifies the fully-qualified name of the stage in which to save a report when [profiling Python handler code](https://docs.snowflake.com/en/developer-guide/stored-procedure/python/procedure-python-profiler).",
	}
	QueryTag = parameterdefs.ParameterDef{
		SqlName:      "QUERY_TAG",
		Kind:         g.KindString,
		Levels:       append(slices.Clone(onUser), parameterdefs.ParameterLevelTask),
		Description:  "Optional string that can be used to tag queries and other SQL statements executed within a session. The tags are displayed in the output of the [QUERY_HISTORY, QUERY_HISTORY_BY_*](https://docs.snowflake.com/en/sql-reference/functions/query_history) functions.",
		DefaultValue: "",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	QuotedIdentifiersIgnoreCase = parameterdefs.ParameterDef{
		SqlName:      "QUOTED_IDENTIFIERS_IGNORE_CASE",
		Kind:         g.KindBool,
		Levels:       append(slices.Clone(onTableAndIcebergTable), parameterdefs.ParameterLevelTask, parameterdefs.ParameterLevelSession, parameterdefs.ParameterLevelUser),
		Description:  "If true, the case of quoted identifiers is ignored.",
		DefaultValue: "false",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	ReadConsistencyMode = parameterdefs.ParameterDef{
		SqlName:     "READ_CONSISTENCY_MODE",
		Kind:        g.KindString,
		Levels:      onAccount,
		Description: "Defines the level of consistency guarantees that are required for sessions with near-concurrent changes.",
	}
	ReplaceInvalidCharacters = parameterdefs.ParameterDef{
		SqlName:      "REPLACE_INVALID_CHARACTERS",
		Kind:         g.KindBool,
		Levels:       onTableAndIcebergTable,
		Description:  "Specifies whether to replace invalid UTF-8 characters with the Unicode replacement character in query results for an Iceberg table. You can only set this parameter for tables that use an external Iceberg catalog.",
		DefaultValue: "false",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	RequireStorageIntegrationForStageCreation = parameterdefs.ParameterDef{
		SqlName:     "REQUIRE_STORAGE_INTEGRATION_FOR_STAGE_CREATION",
		Kind:        g.KindBool,
		Levels:      onAccount,
		Description: "Specifies whether to require a storage integration object as cloud credentials when creating a named external stage (using [CREATE STAGE](https://docs.snowflake.com/en/sql-reference/sql/create-stage)) to access a private cloud storage location.",
	}
	RequireStorageIntegrationForStageOperation = parameterdefs.ParameterDef{
		SqlName:     "REQUIRE_STORAGE_INTEGRATION_FOR_STAGE_OPERATION",
		Kind:        g.KindBool,
		Levels:      onAccount,
		Description: "Specifies whether to require using a named external stage that references a storage integration object as cloud credentials when loading data from or unloading data to a private cloud storage location.",
	}
	RowTimestampDefault = parameterdefs.ParameterDef{
		SqlName:     "ROW_TIMESTAMP_DEFAULT",
		Kind:        g.KindBool,
		Levels:      onAccount,
		Description: "Use this parameter to set row timestamps by default for new tables in a container.",
	}
	RowsPerResultset = parameterdefs.ParameterDef{
		SqlName:      "ROWS_PER_RESULTSET",
		Kind:         g.KindInt,
		Levels:       append(slices.Clone(onUser), parameterdefs.ParameterLevelTask),
		Description:  "Specifies the maximum number of rows returned in a result set. A value of 0 specifies no maximum.",
		DefaultValue: "0",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	S3StageVpceDnsName = parameterdefs.ParameterDef{
		SqlName:      "S3_STAGE_VPCE_DNS_NAME",
		Kind:         g.KindString,
		Levels:       append(slices.Clone(onUser), parameterdefs.ParameterLevelTask),
		Description:  "Specifies the DNS name of an Amazon S3 interface endpoint. Requests sent to the internal stage of an account via [AWS PrivateLink for Amazon S3](https://docs.aws.amazon.com/AmazonS3/latest/userguide/privatelink-interface-endpoints.html) use this endpoint to connect. For more information, see [Accessing Internal stages with dedicated interface endpoints](https://docs.snowflake.com/en/user-guide/private-internal-stages-aws.html#label-aws-privatelink-internal-stage-network-isolation).",
		DefaultValue: "",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	SearchPath = parameterdefs.ParameterDef{
		SqlName:      "SEARCH_PATH",
		Kind:         g.KindString,
		Levels:       onUser,
		Description:  "Specifies the path to search to resolve unqualified object names in queries. For more information, see [Name resolution in queries](https://docs.snowflake.com/en/sql-reference/name-resolution.html#label-object-name-resolution-search-path). Comma-separated list of identifiers. An identifier can be a fully or partially qualified schema name.",
		DefaultValue: "$current, $public",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	ServerlessTaskMaxStatementSize = parameterdefs.ParameterDef{
		SqlName:      "SERVERLESS_TASK_MAX_STATEMENT_SIZE",
		Kind:         g.KindOfT[sdkcommons.WarehouseSize](),
		Levels:       append(slices.Clone(onAccount), parameterdefs.ParameterLevelTask),
		Description:  "Specifies the maximum allowed warehouse size for [Serverless tasks](https://docs.snowflake.com/en/user-guide/tasks-intro.html#label-tasks-compute-resources-serverless).",
		DefaultValue: `sdk.WarehouseSize("X2Large")`,
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	ServerlessTaskMinStatementSize = parameterdefs.ParameterDef{
		SqlName:      "SERVERLESS_TASK_MIN_STATEMENT_SIZE",
		Kind:         g.KindOfT[sdkcommons.WarehouseSize](),
		Levels:       append(slices.Clone(onAccount), parameterdefs.ParameterLevelTask),
		Description:  "Specifies the minimum allowed warehouse size for [Serverless tasks](https://docs.snowflake.com/en/user-guide/tasks-intro.html#label-tasks-compute-resources-serverless).",
		DefaultValue: "sdk.WarehouseSizeXSmall",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	ServiceCallerTokenValiditySecs = parameterdefs.ParameterDef{
		SqlName:      "SERVICE_CALLER_TOKEN_VALIDITY_SECS",
		Kind:         g.KindInt,
		Levels:       []parameterdefs.ParameterLevel{parameterdefs.ParameterLevelService},
		Description:  "Controls how long a caller's rights login token is valid for Snowpark Container Services.",
		DefaultValue: "120",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	ShareRestrictions = parameterdefs.ParameterDef{
		SqlName: "SHARE_RESTRICTIONS",
		Kind:    g.KindBool,
		Levels:  onAccountExt,
	}
	SimulatedDataSharingConsumer = parameterdefs.ParameterDef{
		SqlName:      "SIMULATED_DATA_SHARING_CONSUMER",
		Kind:         g.KindString,
		Levels:       onUser,
		Description:  "Specifies the name of a consumer account to simulate for testing/validating shared data, particularly shared secure views. When this parameter is set in a session, shared views return rows as if executed in the specified consumer account rather than the provider account.",
		DefaultValue: "",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	SqlTraceQueryText = parameterdefs.ParameterDef{
		SqlName:     "SQL_TRACE_QUERY_TEXT",
		Kind:        g.KindString,
		Levels:      onAccount,
		Description: "Used to specify whether to capture the SQL text of a traced SQL statement.",
	}
	SsoLoginPage = parameterdefs.ParameterDef{
		SqlName:     "SSO_LOGIN_PAGE",
		Kind:        g.KindBool,
		Levels:      onAccount,
		Description: "This deprecated parameter disables preview mode for testing SSO (after enabling federated authentication) before rolling it out to users.",
	}
	StatementQueuedTimeoutInSeconds = parameterdefs.ParameterDef{
		SqlName:      "STATEMENT_QUEUED_TIMEOUT_IN_SECONDS",
		Kind:         g.KindInt,
		Levels:       append(slices.Clone(onWarehouseAll), parameterdefs.ParameterLevelTask, parameterdefs.ParameterLevelSession, parameterdefs.ParameterLevelUser),
		Description:  "Amount of time, in seconds, a SQL statement (query, DDL, DML, etc.) remains queued for a warehouse before it is canceled by the system. This parameter can be used in conjunction with the [MAX_CONCURRENCY_LEVEL](https://docs.snowflake.com/en/sql-reference/parameters#label-max-concurrency-level) parameter to ensure a warehouse is never backlogged.",
		DefaultValue: "0",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	StatementTimeoutInSeconds = parameterdefs.ParameterDef{
		SqlName:      "STATEMENT_TIMEOUT_IN_SECONDS",
		Kind:         g.KindInt,
		Levels:       append(slices.Clone(onWarehouseAll), parameterdefs.ParameterLevelTask, parameterdefs.ParameterLevelSession, parameterdefs.ParameterLevelUser),
		Description:  "Amount of time, in seconds, after which a running SQL statement (query, DDL, DML, etc.) is canceled by the system.",
		DefaultValue: "172800",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	StorageSerializationPolicy = parameterdefs.ParameterDef{
		SqlName:      "STORAGE_SERIALIZATION_POLICY",
		Kind:         g.KindOfT[sdkcommons.StorageSerializationPolicy](),
		Levels:       onTableAndIcebergTable,
		Description:  "The storage serialization policy for Iceberg tables that use Snowflake as the catalog. COMPATIBLE: Snowflake performs encoding and compression of data files that ensures interoperability with third-party compute engines. OPTIMIZED: Snowflake performs encoding and compression of data files that ensures the best table performance within Snowflake.",
		DefaultValue: "sdk.StorageSerializationPolicyOptimized",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	StrictJsonOutput = parameterdefs.ParameterDef{
		SqlName:      "STRICT_JSON_OUTPUT",
		Kind:         g.KindBool,
		Levels:       append(slices.Clone(onUser), parameterdefs.ParameterLevelTask),
		Description:  "This parameter specifies whether JSON output in a session is compatible with the general standard (as described by [http://json.org](http://json.org)). By design, Snowflake allows JSON input that contains non-standard values; however, these non-standard values might result in Snowflake outputting JSON that is incompatible with other platforms and languages. This parameter, when enabled, ensures that Snowflake outputs valid/compatible JSON.",
		DefaultValue: "false",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	SuspendTaskAfterNumFailures = parameterdefs.ParameterDef{
		SqlName:      "SUSPEND_TASK_AFTER_NUM_FAILURES",
		Kind:         g.KindInt,
		Levels:       onTask,
		Description:  "How many times a task must fail in a row before it is automatically suspended. 0 disables auto-suspending.",
		DefaultValue: "10",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	TargetFileSize = parameterdefs.ParameterDef{
		SqlName:      "TARGET_FILE_SIZE",
		Kind:         IcebergTableTargetFileSizeEnumDef.Kind(),
		Levels:       onIcebergTableOnly,
		Description:  "Specifies the target file size for the Parquet data files generated for a Snowflake-managed Apache Iceberg™ table.",
		DefaultValue: "sdk.IcebergTableTargetFileSizeAuto",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	TaskAutoRetryAttempts = parameterdefs.ParameterDef{
		SqlName:      "TASK_AUTO_RETRY_ATTEMPTS",
		Kind:         g.KindInt,
		Levels:       onTask,
		Description:  "Maximum automatic retries allowed for a user task.",
		DefaultValue: "0",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	TimeInputFormat = parameterdefs.ParameterDef{
		SqlName:      "TIME_INPUT_FORMAT",
		Kind:         g.KindString,
		Levels:       append(slices.Clone(onUser), parameterdefs.ParameterLevelTask),
		Description:  "Specifies the input format for the TIME data type. For more information, see [Date and time input and output formats](https://docs.snowflake.com/en/sql-reference/date-time-input-output). Any valid, supported time format or AUTO (AUTO specifies that Snowflake attempts to automatically detect the format of times stored in the system during the session).",
		DefaultValue: "AUTO",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	TimeOutputFormat = parameterdefs.ParameterDef{
		SqlName:      "TIME_OUTPUT_FORMAT",
		Kind:         g.KindString,
		Levels:       append(slices.Clone(onUser), parameterdefs.ParameterLevelTask),
		Description:  "Specifies the display format for the TIME data type. For more information, see [Date and time input and output formats](https://docs.snowflake.com/en/sql-reference/date-time-input-output).",
		DefaultValue: "HH24:MI:SS",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	TimestampDayIsAlways24h = parameterdefs.ParameterDef{
		SqlName:      "TIMESTAMP_DAY_IS_ALWAYS_24H",
		Kind:         g.KindBool,
		Levels:       append(slices.Clone(onUser), parameterdefs.ParameterLevelTask),
		Description:  "Specifies whether the [DATEADD](https://docs.snowflake.com/en/sql-reference/functions/dateadd) function (and its aliases) always consider a day to be exactly 24 hours for expressions that span multiple days.",
		DefaultValue: "false",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	TimestampInputFormat = parameterdefs.ParameterDef{
		SqlName:      "TIMESTAMP_INPUT_FORMAT",
		Kind:         g.KindString,
		Levels:       append(slices.Clone(onUser), parameterdefs.ParameterLevelTask),
		Description:  "Specifies the input format for the TIMESTAMP data type alias. For more information, see [Date and time input and output formats](https://docs.snowflake.com/en/sql-reference/date-time-input-output). Any valid, supported timestamp format or AUTO (AUTO specifies that Snowflake attempts to automatically detect the format of timestamps stored in the system during the session).",
		DefaultValue: "AUTO",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	TimestampLtzOutputFormat = parameterdefs.ParameterDef{
		SqlName:      "TIMESTAMP_LTZ_OUTPUT_FORMAT",
		Kind:         g.KindString,
		Levels:       append(slices.Clone(onUser), parameterdefs.ParameterLevelTask),
		Description:  "Specifies the display format for the TIMESTAMP_LTZ data type. If no format is specified, defaults to [TIMESTAMP_OUTPUT_FORMAT](https://docs.snowflake.com/en/sql-reference/parameters#label-timestamp-output-format). For more information, see [Date and time input and output formats](https://docs.snowflake.com/en/sql-reference/date-time-input-output).",
		DefaultValue: "",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	TimestampNtzOutputFormat = parameterdefs.ParameterDef{
		SqlName:      "TIMESTAMP_NTZ_OUTPUT_FORMAT",
		Kind:         g.KindString,
		Levels:       append(slices.Clone(onUser), parameterdefs.ParameterLevelTask),
		Description:  "Specifies the display format for the TIMESTAMP_NTZ data type.",
		DefaultValue: "YYYY-MM-DD HH24:MI:SS.FF3",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	TimestampOutputFormat = parameterdefs.ParameterDef{
		SqlName:      "TIMESTAMP_OUTPUT_FORMAT",
		Kind:         g.KindString,
		Levels:       append(slices.Clone(onUser), parameterdefs.ParameterLevelTask),
		Description:  "Specifies the display format for the TIMESTAMP data type alias. For more information, see [Date and time input and output formats](https://docs.snowflake.com/en/sql-reference/date-time-input-output).",
		DefaultValue: "YYYY-MM-DD HH24:MI:SS.FF3 TZHTZM",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	TimestampTypeMapping = parameterdefs.ParameterDef{
		SqlName:      "TIMESTAMP_TYPE_MAPPING",
		Kind:         g.KindOfT[sdkcommons.TimestampTypeMapping](),
		Levels:       append(slices.Clone(onUser), parameterdefs.ParameterLevelTask),
		Description:  "Specifies the TIMESTAMP_* variation that the TIMESTAMP data type alias maps to.",
		DefaultValue: "sdk.TimestampTypeMappingNtz",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	TimestampTzOutputFormat = parameterdefs.ParameterDef{
		SqlName:      "TIMESTAMP_TZ_OUTPUT_FORMAT",
		Kind:         g.KindString,
		Levels:       append(slices.Clone(onUser), parameterdefs.ParameterLevelTask),
		Description:  "Specifies the display format for the TIMESTAMP_TZ data type. If no format is specified, defaults to [TIMESTAMP_OUTPUT_FORMAT](https://docs.snowflake.com/en/sql-reference/parameters#label-timestamp-output-format). For more information, see [Date and time input and output formats](https://docs.snowflake.com/en/sql-reference/date-time-input-output).",
		DefaultValue: "",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	Timezone = parameterdefs.ParameterDef{
		SqlName:      "TIMEZONE",
		Kind:         g.KindString,
		Levels:       append(slices.Clone(onUser), parameterdefs.ParameterLevelTask),
		Description:  "Specifies the time zone for the session. You can specify a [time zone name](https://data.iana.org/time-zones/tzdb-2021a/zone1970.tab) or a [link name](https://data.iana.org/time-zones/tzdb-2021a/backward) from release 2021a of the [IANA Time Zone Database](https://www.iana.org/time-zones) (e.g. America/Los_Angeles, Europe/London, UTC, Etc/GMT, etc.).",
		DefaultValue: "America/Los_Angeles",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	TraceLevel = parameterdefs.ParameterDef{
		SqlName:      "TRACE_LEVEL",
		Kind:         g.KindOfT[sdkcommons.TraceLevel](),
		Levels:       append(slices.Clone(onFunctionAndProcedure), parameterdefs.ParameterLevelTask, parameterdefs.ParameterLevelSession, parameterdefs.ParameterLevelUser),
		Description:  "Controls how trace events are ingested into the event table.",
		DefaultValue: "sdk.TraceLevelOff",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	TransactionAbortOnError = parameterdefs.ParameterDef{
		SqlName:      "TRANSACTION_ABORT_ON_ERROR",
		Kind:         g.KindBool,
		Levels:       append(slices.Clone(onUser), parameterdefs.ParameterLevelTask),
		Description:  "Specifies the action to perform when a statement issued within a non-autocommit transaction returns with an error.",
		DefaultValue: "false",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	TransactionDefaultIsolationLevel = parameterdefs.ParameterDef{
		SqlName:      "TRANSACTION_DEFAULT_ISOLATION_LEVEL",
		Kind:         g.KindOfT[sdkcommons.TransactionDefaultIsolationLevel](),
		Levels:       append(slices.Clone(onUser), parameterdefs.ParameterLevelTask),
		Description:  "Specifies the isolation level for transactions in the user session.",
		DefaultValue: "sdk.TransactionDefaultIsolationLevelReadCommitted",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	TwoDigitCenturyStart = parameterdefs.ParameterDef{
		SqlName:      "TWO_DIGIT_CENTURY_START",
		Kind:         g.KindInt,
		Levels:       append(slices.Clone(onUser), parameterdefs.ParameterLevelTask),
		Description:  "Specifies the “century start” year for 2-digit years (i.e. the earliest year such dates can represent). This parameter prevents ambiguous dates when importing or converting data with the `YY` date format component (i.e. years represented as 2 digits).",
		DefaultValue: "1970",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	UnsupportedDdlAction = parameterdefs.ParameterDef{
		SqlName:     "UNSUPPORTED_DDL_ACTION",
		Kind:        g.KindOfT[sdkcommons.UnsupportedDDLAction](),
		Levels:      append(slices.Clone(onUser), parameterdefs.ParameterLevelTask),
		Description: "Determines if an unsupported (i.e. non-default) value specified for a constraint property returns an error.",
		// TODO [SNOW-1501905]: quick workaround for now: lowercase for ignore in snowflake by default but uppercase for FAIL
		DefaultValue: "sdk.UnsupportedDDLAction(strings.ToLower(string(sdk.UnsupportedDDLActionIgnore)))",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	UseCachedResult = parameterdefs.ParameterDef{
		SqlName:      "USE_CACHED_RESULT",
		Kind:         g.KindBool,
		Levels:       append(slices.Clone(onUser), parameterdefs.ParameterLevelTask),
		Description:  "Specifies whether to reuse persisted query results, if available, when a matching query is submitted.",
		DefaultValue: "true",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	UseWorkspacesForSql = parameterdefs.ParameterDef{
		SqlName:     "USE_WORKSPACES_FOR_SQL",
		Kind:        g.KindString,
		Levels:      onAccount,
		Description: "Controls whether the Workspaces editor is the default SQL editing experience for the account. Valid values are `always` and `never`.",
	}
	UserTaskManagedInitialWarehouseSize = parameterdefs.ParameterDef{
		SqlName:      "USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE",
		Kind:         g.KindOfT[sdkcommons.WarehouseSize](),
		Levels:       onTask,
		Description:  "The initial size of warehouse to use for managed warehouses in the absence of history.",
		DefaultValue: "sdk.WarehouseSizeMedium",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	UserTaskMinimumTriggerIntervalInSeconds = parameterdefs.ParameterDef{
		SqlName:      "USER_TASK_MINIMUM_TRIGGER_INTERVAL_IN_SECONDS",
		Kind:         g.KindInt,
		Levels:       onTask,
		Description:  "Minimum amount of time between Triggered Task executions in seconds.",
		DefaultValue: "30",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	UserTaskTimeoutMs = parameterdefs.ParameterDef{
		SqlName:      "USER_TASK_TIMEOUT_MS",
		Kind:         g.KindInt,
		Levels:       onTask,
		Description:  "User task execution timeout in milliseconds.",
		DefaultValue: "3600000",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	WeekOfYearPolicy = parameterdefs.ParameterDef{
		SqlName:      "WEEK_OF_YEAR_POLICY",
		Kind:         g.KindInt,
		Levels:       append(slices.Clone(onUser), parameterdefs.ParameterLevelTask),
		Description:  "Specifies how the weeks in a given year are computed. `0`: The semantics used are equivalent to the ISO semantics, in which a week belongs to a given year if at least 4 days of that week are in that year. `1`: January 1 is included in the first week of the year and December 31 is included in the last week of the year.",
		DefaultValue: "0",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
	WeekStart = parameterdefs.ParameterDef{
		SqlName:      "WEEK_START",
		Kind:         g.KindInt,
		Levels:       append(slices.Clone(onUser), parameterdefs.ParameterLevelTask),
		Description:  "Specifies the first day of the week (used by week-related date functions). `0`: Legacy Snowflake behavior is used (i.e. ISO-like semantics). `1` (Monday) to `7` (Sunday): All the week-related functions use weeks that start on the specified day of the week.",
		DefaultValue: "0",
		DefaultLevel: parameterTypeSnowflakeDefault,
	}
)

var AllParameters = []parameterdefs.ParameterDef{
	AbortDetachedQuery,
	ActivePythonProfiler,
	AllowBindValuesAccess,
	AllowClientMfaCaching,
	AllowIdToken,
	AllowRowTimestamp,
	AllowedSpcsWorkloadTypes,
	Autocommit,
	BaseLocationPrefix,
	BinaryInputFormat,
	BinaryOutputFormat,
	Catalog,
	CatalogSync,
	ClientEnableLogInfoStatementParameters,
	ClientEncryptionKeySize,
	ClientMemoryLimit,
	ClientMetadataRequestUseConnectionCtx,
	ClientMetadataUseSessionDatabase,
	ClientPrefetchThreads,
	ClientResultChunkSize,
	ClientResultColumnCaseInsensitive,
	ClientSessionKeepAlive,
	ClientSessionKeepAliveHeartbeatFrequency,
	ClientTimestampTypeMapping,
	CortexCodeCliDailyEstCreditLimitPerUser,
	CortexCodeDesktopDailyEstCreditLimitPerUser,
	CortexCodeSnowsightDailyEstCreditLimitPerUser,
	CortexEnabledCrossRegion,
	CortexModelsAllowlist,
	CsvTimestampFormat,
	DataMetricSchedule,
	DataRetentionTimeInDays,
	DateInputFormat,
	DateOutputFormat,
	DefaultDbtVersion,
	DefaultDdlCollation,
	DefaultNotebookComputePoolCpu,
	DefaultNotebookComputePoolGpu,
	DefaultNullOrdering,
	DefaultStreamlitComputePool,
	DefaultStreamlitNotebookWarehouse,
	DisableUiDownloadButton,
	DisableUserPrivilegeGrants,
	DisallowedSpcsWorkloadTypes,
	EnableAutomaticSensitiveDataClassificationLog,
	EnableBudgetEventLogging,
	EnableConsoleOutput,
	EnableCortexAnalyst,
	EnableDataCompaction,
	EnableEgressCostOptimizer,
	EnableGetDdlUseDataTypeAlias,
	EnableIcebergMergeOnRead,
	EnableIdentifierFirstLogin,
	EnableInternalStagesPrivatelink,
	EnableNotebookCreationInPersonalDb,
	EnablePerAccountAppServicePrivatelinkUrl,
	EnablePersonalDatabase,
	EnableSpcsBlockStorageSnowflakeFullEncryptionEnforcement,
	EnableTagPropagationEventLogging,
	EnableTriSecretAndRekeyOptOutForImageRepository,
	EnableTriSecretAndRekeyOptOutForSpcsBlockStorage,
	EnableUnhandledExceptionsReporting,
	EnableUnloadPhysicalTypeOptimization,
	EnableUnredactedQuerySyntaxError,
	EnableUnredactedSecureObjectError,
	EnforceNetworkRulesForInternalStages,
	ErrorOnNondeterministicMerge,
	ErrorOnNondeterministicUpdate,
	EventTable,
	ExternalOauthAddPrivilegedRolesToBlockedList,
	ExternalVolume,
	FallbackWarehouse,
	GeographyOutputFormat,
	GeometryOutputFormat,
	HybridTableLockTimeout,
	IcebergMergeOnReadBehavior,
	IcebergVersionDefault,
	InitialReplicationSizeLimitInTb,
	JdbcTreatDecimalAsInt,
	JdbcTreatTimestampNtzAsUtc,
	JdbcUseSessionTimezone,
	JsTreatIntegerAsBigint,
	JsonIndent,
	ListingAutoFulfillmentReplicationRefreshSchedule,
	LockTimeout,
	LogEventLevel,
	LogLevel,
	MaxConcurrencyLevel,
	MaxDataExtensionTimeInDays,
	MetricLevel,
	MinDataRetentionTimeInDays,
	MultiStatementCount,
	NetworkPolicy,
	NoorderSequenceAsDefault,
	OauthAddPrivilegedRolesToBlockedList,
	OdbcTreatDecimalAsInt,
	OptimizeDataLayout,
	PeriodicDataRekeying,
	PipeExecutionPaused,
	PreventLoadFromInlineUrl,
	PreventUnloadToInlineUrl,
	PreventUnloadToInternalStages,
	PythonProfilerModules,
	PythonProfilerTargetStage,
	QueryTag,
	QuotedIdentifiersIgnoreCase,
	ReadConsistencyMode,
	ReplaceInvalidCharacters,
	RequireStorageIntegrationForStageCreation,
	RequireStorageIntegrationForStageOperation,
	RowTimestampDefault,
	RowsPerResultset,
	S3StageVpceDnsName,
	SearchPath,
	ServerlessTaskMaxStatementSize,
	ServerlessTaskMinStatementSize,
	ServiceCallerTokenValiditySecs,
	ShareRestrictions,
	SimulatedDataSharingConsumer,
	SqlTraceQueryText,
	SsoLoginPage,
	StatementQueuedTimeoutInSeconds,
	StatementTimeoutInSeconds,
	StorageSerializationPolicy,
	StrictJsonOutput,
	SuspendTaskAfterNumFailures,
	TargetFileSize,
	TaskAutoRetryAttempts,
	TimeInputFormat,
	TimeOutputFormat,
	TimestampDayIsAlways24h,
	TimestampInputFormat,
	TimestampLtzOutputFormat,
	TimestampNtzOutputFormat,
	TimestampOutputFormat,
	TimestampTypeMapping,
	TimestampTzOutputFormat,
	Timezone,
	TraceLevel,
	TransactionAbortOnError,
	TransactionDefaultIsolationLevel,
	TwoDigitCenturyStart,
	UnsupportedDdlAction,
	UseCachedResult,
	UseWorkspacesForSql,
	UserTaskManagedInitialWarehouseSize,
	UserTaskMinimumTriggerIntervalInSeconds,
	UserTaskTimeoutMs,
	WeekOfYearPolicy,
	WeekStart,
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

func AssertionParameterType(kind string) string {
	switch kind {
	case g.KindBool, g.KindInt, g.KindString:
		return kind
	case g.KindOfT[sdkcommons.StringAllowEmpty]():
		return g.KindString
	}
	if _, err := g.ToObjectIdentifierKind(kind); err == nil {
		return g.KindString
	}
	return "sdk." + kind
}
