package sdk

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

// TODO: add more tests
func TestSetObjectParameterOnObject(t *testing.T) {
	id := randomAccountObjectIdentifier()

	defaultOpts := func() *setParameterOnObject {
		return &setParameterOnObject{
			objectType:       ObjectTypeUser,
			objectIdentifier: id,
			parameterKey:     "ENABLE_UNREDACTED_QUERY_SYNTAX_ERROR",
			parameterValue:   "TRUE",
		}
	}

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSqlEqualsf(t, opts, "ALTER USER %s SET ENABLE_UNREDACTED_QUERY_SYNTAX_ERROR = TRUE", id.FullyQualifiedName())
	})
}

func TestUnsetObjectParameterOnObject(t *testing.T) {
	t.Run("unset on user", func(t *testing.T) {
		id := randomAccountObjectIdentifier()
		opts := &unsetParameterOnObject{
			objectType:       ObjectTypeUser,
			objectIdentifier: id,
			parameterKey:     "ENABLE_UNREDACTED_QUERY_SYNTAX_ERROR",
		}
		assertOptsValidAndSqlEqualsf(t, opts, "ALTER USER %s UNSET ENABLE_UNREDACTED_QUERY_SYNTAX_ERROR", id.FullyQualifiedName())
	})

	t.Run("unset on database", func(t *testing.T) {
		id := randomAccountObjectIdentifier()
		opts := &unsetParameterOnObject{
			objectType:       ObjectTypeDatabase,
			objectIdentifier: id,
			parameterKey:     "USER_TASK_TIMEOUT_MS",
		}
		assertOptsValidAndSqlEqualsf(t, opts, "ALTER DATABASE %s UNSET USER_TASK_TIMEOUT_MS", id.FullyQualifiedName())
	})

	t.Run("unset on schema", func(t *testing.T) {
		id := randomDatabaseObjectIdentifier()
		opts := &unsetParameterOnObject{
			objectType:       ObjectTypeSchema,
			objectIdentifier: id,
			parameterKey:     "DATA_RETENTION_TIME_IN_DAYS",
		}
		assertOptsValidAndSqlEqualsf(t, opts, "ALTER SCHEMA %s UNSET DATA_RETENTION_TIME_IN_DAYS", id.FullyQualifiedName())
	})
}

func TestUnsetObjectParameterNetworkPolicyOnUser(t *testing.T) {
	opts := &AlterUserOptions{
		name: NewAccountObjectIdentifierFromFullyQualifiedName("TEST_USER"),
		Unset: &UserUnset{
			NetworkPolicy: Bool(true),
		},
	}
	t.Run("Unset User Network Policy", func(t *testing.T) {
		assertOptsValidAndSqlEqualsf(t, opts, `ALTER USER "TEST_USER" UNSET NETWORK_POLICY`)
	})
}

// Proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3344 is fixed.
func TestToAccountParameter(t *testing.T) {
	type test struct {
		input string
		want  AccountParameter
	}

	valid := []test{
		// Case insensitive.
		{input: "allow_client_mfa_caching", want: AccountParameterAllowClientMfaCaching},

		// Supported Values.
		{input: "ALLOW_CLIENT_MFA_CACHING", want: AccountParameterAllowClientMfaCaching},
		{input: "ALLOW_ID_TOKEN", want: AccountParameterAllowIdToken},
		{input: "CLIENT_ENCRYPTION_KEY_SIZE", want: AccountParameterClientEncryptionKeySize},
		{input: "CORTEX_CODE_CLI_DAILY_EST_CREDIT_LIMIT_PER_USER", want: AccountParameterCortexCodeCliDailyEstCreditLimitPerUser},
		{input: "CORTEX_CODE_DESKTOP_DAILY_EST_CREDIT_LIMIT_PER_USER", want: AccountParameterCortexCodeDesktopDailyEstCreditLimitPerUser},
		{input: "CORTEX_CODE_SNOWSIGHT_DAILY_EST_CREDIT_LIMIT_PER_USER", want: AccountParameterCortexCodeSnowsightDailyEstCreditLimitPerUser},
		{input: "CORTEX_ENABLED_CROSS_REGION", want: AccountParameterCortexEnabledCrossRegion},
		{input: "DEFAULT_STREAMLIT_COMPUTE_POOL", want: AccountParameterDefaultStreamlitComputePool},
		{input: "DISABLE_USER_PRIVILEGE_GRANTS", want: AccountParameterDisableUserPrivilegeGrants},
		{input: "ENABLE_IDENTIFIER_FIRST_LOGIN", want: AccountParameterEnableIdentifierFirstLogin},
		{input: "ENABLE_INTERNAL_STAGES_PRIVATELINK", want: AccountParameterEnableInternalStagesPrivatelink},
		{input: "ENABLE_TRI_SECRET_AND_REKEY_OPT_OUT_FOR_IMAGE_REPOSITORY", want: AccountParameterEnableTriSecretAndRekeyOptOutForImageRepository},
		{input: "ENABLE_TRI_SECRET_AND_REKEY_OPT_OUT_FOR_SPCS_BLOCK_STORAGE", want: AccountParameterEnableTriSecretAndRekeyOptOutForSpcsBlockStorage},
		{input: "ENABLE_UNHANDLED_EXCEPTIONS_REPORTING", want: AccountParameterEnableUnhandledExceptionsReporting},
		{input: "ENFORCE_NETWORK_RULES_FOR_INTERNAL_STAGES", want: AccountParameterEnforceNetworkRulesForInternalStages},
		{input: "EVENT_TABLE", want: AccountParameterEventTable},
		{input: "EXTERNAL_OAUTH_ADD_PRIVILEGED_ROLES_TO_BLOCKED_LIST", want: AccountParameterExternalOauthAddPrivilegedRolesToBlockedList},
		{input: "INITIAL_REPLICATION_SIZE_LIMIT_IN_TB", want: AccountParameterInitialReplicationSizeLimitInTb},
		{input: "MIN_DATA_RETENTION_TIME_IN_DAYS", want: AccountParameterMinDataRetentionTimeInDays},
		{input: "METRIC_LEVEL", want: AccountParameterMetricLevel},
		{input: "NETWORK_POLICY", want: AccountParameterNetworkPolicy},
		{input: "OAUTH_ADD_PRIVILEGED_ROLES_TO_BLOCKED_LIST", want: AccountParameterOauthAddPrivilegedRolesToBlockedList},
		{input: "PERIODIC_DATA_REKEYING", want: AccountParameterPeriodicDataRekeying},
		{input: "PREVENT_UNLOAD_TO_INLINE_URL", want: AccountParameterPreventUnloadToInlineUrl},
		{input: "REQUIRE_STORAGE_INTEGRATION_FOR_STAGE_CREATION", want: AccountParameterRequireStorageIntegrationForStageCreation},
		{input: "REQUIRE_STORAGE_INTEGRATION_FOR_STAGE_OPERATION", want: AccountParameterRequireStorageIntegrationForStageOperation},
		{input: "SSO_LOGIN_PAGE", want: AccountParameterSsoLoginPage},

		// Session Parameters (inherited)
		{input: "ABORT_DETACHED_QUERY", want: AccountParameterAbortDetachedQuery},
		{input: "ACTIVE_PYTHON_PROFILER", want: AccountParameterActivePythonProfiler},
		{input: "AUTOCOMMIT", want: AccountParameterAutocommit},
		{input: "BINARY_INPUT_FORMAT", want: AccountParameterBinaryInputFormat},
		{input: "BINARY_OUTPUT_FORMAT", want: AccountParameterBinaryOutputFormat},
		{input: "CLIENT_ENABLE_LOG_INFO_STATEMENT_PARAMETERS", want: AccountParameterClientEnableLogInfoStatementParameters},
		{input: "CLIENT_MEMORY_LIMIT", want: AccountParameterClientMemoryLimit},
		{input: "CLIENT_METADATA_REQUEST_USE_CONNECTION_CTX", want: AccountParameterClientMetadataRequestUseConnectionCtx},
		{input: "CLIENT_METADATA_USE_SESSION_DATABASE", want: AccountParameterClientMetadataUseSessionDatabase},
		{input: "CLIENT_PREFETCH_THREADS", want: AccountParameterClientPrefetchThreads},
		{input: "CLIENT_RESULT_CHUNK_SIZE", want: AccountParameterClientResultChunkSize},
		{input: "CLIENT_RESULT_COLUMN_CASE_INSENSITIVE", want: AccountParameterClientResultColumnCaseInsensitive},
		{input: "CLIENT_SESSION_KEEP_ALIVE", want: AccountParameterClientSessionKeepAlive},
		{input: "CLIENT_SESSION_KEEP_ALIVE_HEARTBEAT_FREQUENCY", want: AccountParameterClientSessionKeepAliveHeartbeatFrequency},
		{input: "CLIENT_TIMESTAMP_TYPE_MAPPING", want: AccountParameterClientTimestampTypeMapping},
		{input: "CSV_TIMESTAMP_FORMAT", want: AccountParameterCsvTimestampFormat},
		{input: "DATE_INPUT_FORMAT", want: AccountParameterDateInputFormat},
		{input: "DATE_OUTPUT_FORMAT", want: AccountParameterDateOutputFormat},
		{input: "ENABLE_UNLOAD_PHYSICAL_TYPE_OPTIMIZATION", want: AccountParameterEnableUnloadPhysicalTypeOptimization},
		{input: "ERROR_ON_NONDETERMINISTIC_MERGE", want: AccountParameterErrorOnNondeterministicMerge},
		{input: "ERROR_ON_NONDETERMINISTIC_UPDATE", want: AccountParameterErrorOnNondeterministicUpdate},
		{input: "GEOGRAPHY_OUTPUT_FORMAT", want: AccountParameterGeographyOutputFormat},
		{input: "GEOMETRY_OUTPUT_FORMAT", want: AccountParameterGeometryOutputFormat},
		{input: "HYBRID_TABLE_LOCK_TIMEOUT", want: AccountParameterHybridTableLockTimeout},
		{input: "JDBC_TREAT_DECIMAL_AS_INT", want: AccountParameterJdbcTreatDecimalAsInt},
		{input: "JDBC_TREAT_TIMESTAMP_NTZ_AS_UTC", want: AccountParameterJdbcTreatTimestampNtzAsUtc},
		{input: "JDBC_USE_SESSION_TIMEZONE", want: AccountParameterJdbcUseSessionTimezone},
		{input: "JSON_INDENT", want: AccountParameterJsonIndent},
		{input: "JS_TREAT_INTEGER_AS_BIGINT", want: AccountParameterJsTreatIntegerAsBigint},
		{input: "LOCK_TIMEOUT", want: AccountParameterLockTimeout},
		{input: "MULTI_STATEMENT_COUNT", want: AccountParameterMultiStatementCount},
		{input: "NOORDER_SEQUENCE_AS_DEFAULT", want: AccountParameterNoorderSequenceAsDefault},
		{input: "ODBC_TREAT_DECIMAL_AS_INT", want: AccountParameterOdbcTreatDecimalAsInt},
		{input: "PYTHON_PROFILER_MODULES", want: AccountParameterPythonProfilerModules},
		{input: "PYTHON_PROFILER_TARGET_STAGE", want: AccountParameterPythonProfilerTargetStage},
		{input: "QUERY_TAG", want: AccountParameterQueryTag},
		{input: "QUOTED_IDENTIFIERS_IGNORE_CASE", want: AccountParameterQuotedIdentifiersIgnoreCase},
		{input: "ROWS_PER_RESULTSET", want: AccountParameterRowsPerResultset},
		{input: "S3_STAGE_VPCE_DNS_NAME", want: AccountParameterS3StageVpceDnsName},
		{input: "SEARCH_PATH", want: AccountParameterSearchPath},
		{input: "SIMULATED_DATA_SHARING_CONSUMER", want: AccountParameterSimulatedDataSharingConsumer},
		{input: "STRICT_JSON_OUTPUT", want: AccountParameterStrictJsonOutput},
		{input: "TIME_INPUT_FORMAT", want: AccountParameterTimeInputFormat},
		{input: "TIME_OUTPUT_FORMAT", want: AccountParameterTimeOutputFormat},
		{input: "TIMESTAMP_DAY_IS_ALWAYS_24H", want: AccountParameterTimestampDayIsAlways24h},
		{input: "TIMESTAMP_INPUT_FORMAT", want: AccountParameterTimestampInputFormat},
		{input: "TIMESTAMP_LTZ_OUTPUT_FORMAT", want: AccountParameterTimestampLtzOutputFormat},
		{input: "TIMESTAMP_NTZ_OUTPUT_FORMAT", want: AccountParameterTimestampNtzOutputFormat},
		{input: "TIMESTAMP_OUTPUT_FORMAT", want: AccountParameterTimestampOutputFormat},
		{input: "TIMESTAMP_TYPE_MAPPING", want: AccountParameterTimestampTypeMapping},
		{input: "TIMESTAMP_TZ_OUTPUT_FORMAT", want: AccountParameterTimestampTzOutputFormat},
		{input: "TIMEZONE", want: AccountParameterTimezone},
		{input: "TRANSACTION_ABORT_ON_ERROR", want: AccountParameterTransactionAbortOnError},
		{input: "TRANSACTION_DEFAULT_ISOLATION_LEVEL", want: AccountParameterTransactionDefaultIsolationLevel},
		{input: "TWO_DIGIT_CENTURY_START", want: AccountParameterTwoDigitCenturyStart},
		{input: "UNSUPPORTED_DDL_ACTION", want: AccountParameterUnsupportedDdlAction},
		{input: "USE_CACHED_RESULT", want: AccountParameterUseCachedResult},
		{input: "WEEK_OF_YEAR_POLICY", want: AccountParameterWeekOfYearPolicy},
		{input: "WEEK_START", want: AccountParameterWeekStart},

		// Object Parameters (inherited)
		{input: "CATALOG", want: AccountParameterCatalog},
		{input: "DATA_RETENTION_TIME_IN_DAYS", want: AccountParameterDataRetentionTimeInDays},
		{input: "DEFAULT_DDL_COLLATION", want: AccountParameterDefaultDdlCollation},
		{input: "DEFAULT_NOTEBOOK_COMPUTE_POOL_CPU", want: AccountParameterDefaultNotebookComputePoolCpu},
		{input: "DEFAULT_NOTEBOOK_COMPUTE_POOL_GPU", want: AccountParameterDefaultNotebookComputePoolGpu},
		{input: "EXTERNAL_VOLUME", want: AccountParameterExternalVolume},
		{input: "LOG_LEVEL", want: AccountParameterLogLevel},
		{input: "MAX_CONCURRENCY_LEVEL", want: AccountParameterMaxConcurrencyLevel},
		{input: "MAX_DATA_EXTENSION_TIME_IN_DAYS", want: AccountParameterMaxDataExtensionTimeInDays},
		{input: "PIPE_EXECUTION_PAUSED", want: AccountParameterPipeExecutionPaused},
		{input: "REPLACE_INVALID_CHARACTERS", want: AccountParameterReplaceInvalidCharacters},
		{input: "STATEMENT_QUEUED_TIMEOUT_IN_SECONDS", want: AccountParameterStatementQueuedTimeoutInSeconds},
		{input: "STATEMENT_TIMEOUT_IN_SECONDS", want: AccountParameterStatementTimeoutInSeconds},
		{input: "STORAGE_SERIALIZATION_POLICY", want: AccountParameterStorageSerializationPolicy},
		{input: "SUSPEND_TASK_AFTER_NUM_FAILURES", want: AccountParameterSuspendTaskAfterNumFailures},
		{input: "TRACE_LEVEL", want: AccountParameterTraceLevel},
		{input: "USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE", want: AccountParameterUserTaskManagedInitialWarehouseSize},
		{input: "USER_TASK_TIMEOUT_MS", want: AccountParameterUserTaskTimeoutMs},
		{input: "TASK_AUTO_RETRY_ATTEMPTS", want: AccountParameterTaskAutoRetryAttempts},
		{input: "USER_TASK_MINIMUM_TRIGGER_INTERVAL_IN_SECONDS", want: AccountParameterUserTaskMinimumTriggerIntervalInSeconds},
		{input: "METRIC_LEVEL", want: AccountParameterMetricLevel},

		// User Parameters (inherited)
		{input: "ENABLE_UNREDACTED_QUERY_SYNTAX_ERROR", want: AccountParameterEnableUnredactedQuerySyntaxError},
		{input: "PREVENT_UNLOAD_TO_INTERNAL_STAGES", want: AccountParameterPreventUnloadToInternalStages},
	}

	invalid := []test{
		{input: ""},
		{input: "foo"},
	}

	for _, tc := range valid {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ToAccountParameter(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}

	for _, tc := range invalid {
		t.Run(tc.input, func(t *testing.T) {
			_, err := ToAccountParameter(tc.input)
			require.Error(t, err)
		})
	}
}

func TestSetObjectParameterOnObjectViaCatalog(t *testing.T) {
	cases := []struct {
		name        string
		object      Object
		parameter   ObjectParameter
		value       string
		wantHandled bool
		wantErr     bool
	}{
		{
			name:      "unknown key falls back for warehouse",
			object:    Object{ObjectType: ObjectTypeWarehouse, Name: NewAccountObjectIdentifier("some_object")},
			parameter: "DATA_RETENTION_TIME_IN_DAYS",
			value:     "1",
		},
		{
			name:        "warehouse parse error propagates",
			object:      Object{ObjectType: ObjectTypeWarehouse, Name: NewAccountObjectIdentifier("some_object")},
			parameter:   "MAX_CONCURRENCY_LEVEL",
			value:       "not-an-int",
			wantHandled: true,
			wantErr:     true,
		},
		{
			name:      "unknown key falls back for user",
			object:    Object{ObjectType: ObjectTypeUser, Name: NewAccountObjectIdentifier("some_object")},
			parameter: "DATA_RETENTION_TIME_IN_DAYS",
			value:     "1",
		},
		{
			name:        "user parse error propagates",
			object:      Object{ObjectType: ObjectTypeUser, Name: NewAccountObjectIdentifier("some_object")},
			parameter:   "ENABLE_UNREDACTED_QUERY_SYNTAX_ERROR",
			value:       "not-a-bool",
			wantHandled: true,
			wantErr:     true,
		},
		{
			name:        "task parse error propagates",
			object:      Object{ObjectType: ObjectTypeTask, Name: NewSchemaObjectIdentifier("some_db", "some_schema", "some_task")},
			parameter:   "ABORT_DETACHED_QUERY",
			value:       "not-a-bool",
			wantHandled: true,
			wantErr:     true,
		},
		{
			name:        "service parse error propagates",
			object:      Object{ObjectType: ObjectTypeService, Name: NewSchemaObjectIdentifier("some_db", "some_schema", "some_service")},
			parameter:   "SERVICE_CALLER_TOKEN_VALIDITY_SECS",
			value:       "not-an-int",
			wantHandled: true,
			wantErr:     true,
		},
		{
			name:        "hybrid table parse error propagates",
			object:      Object{ObjectType: ObjectTypeHybridTable, Name: NewSchemaObjectIdentifier("some_db", "some_schema", "some_table")},
			parameter:   "DATA_RETENTION_TIME_IN_DAYS",
			value:       "not-an-int",
			wantHandled: true,
			wantErr:     true,
		},
		{
			name:        "iceberg table parse error propagates",
			object:      Object{ObjectType: ObjectTypeIcebergTable, Name: NewSchemaObjectIdentifier("some_db", "some_schema", "some_table")},
			parameter:   "DATA_RETENTION_TIME_IN_DAYS",
			value:       "not-an-int",
			wantHandled: true,
			wantErr:     true,
		},
		{
			name:      "function without arguments falls back",
			object:    Object{ObjectType: ObjectTypeFunction, Name: NewSchemaObjectIdentifier("some_db", "some_schema", "some_fn")},
			parameter: "LOG_LEVEL",
			value:     "INFO",
		},
		{
			name:      "unknown function enum value falls back",
			object:    Object{ObjectType: ObjectTypeFunction, Name: NewSchemaObjectIdentifierWithArguments("some_db", "some_schema", "some_fn", dataTypeVarchar)},
			parameter: "LOG_LEVEL",
			value:     "not-a-level",
		},
		{
			name:      "procedure without arguments falls back",
			object:    Object{ObjectType: ObjectTypeProcedure, Name: NewSchemaObjectIdentifier("some_db", "some_schema", "some_proc")},
			parameter: "LOG_LEVEL",
			value:     "INFO",
		},
		{
			name:        "procedure parse error propagates",
			object:      Object{ObjectType: ObjectTypeProcedure, Name: NewSchemaObjectIdentifierWithArguments("some_db", "some_schema", "some_proc", dataTypeVarchar)},
			parameter:   "ENABLE_CONSOLE_OUTPUT",
			value:       "not-a-bool",
			wantHandled: true,
			wantErr:     true,
		},
		{
			name:      "no catalog writer for table",
			object:    Object{ObjectType: ObjectTypeTable, Name: NewSchemaObjectIdentifier("some_db", "some_schema", "some_table")},
			parameter: "DATA_RETENTION_TIME_IN_DAYS",
			value:     "1",
		},
		{
			name:      "no catalog writer for pipe",
			object:    Object{ObjectType: ObjectTypePipe, Name: NewSchemaObjectIdentifier("some_db", "some_schema", "some_pipe")},
			parameter: "LOG_EVENT_LEVEL",
			value:     "INFO",
		},
		{
			name:      "unknown key falls back for database",
			object:    Object{ObjectType: ObjectTypeDatabase, Name: NewAccountObjectIdentifier("some_db")},
			parameter: "NOT_A_REAL_PARAMETER",
			value:     "1",
		},
		{
			name:      "unknown key falls back for schema",
			object:    Object{ObjectType: ObjectTypeSchema, Name: NewDatabaseObjectIdentifier("some_db", "some_schema")},
			parameter: "NOT_A_REAL_PARAMETER",
			value:     "1",
		},
		{
			name:      "wrong identifier type falls back",
			object:    Object{ObjectType: ObjectTypeDatabase, Name: NewDatabaseObjectIdentifier("some_db", "some_schema")},
			parameter: "DATA_RETENTION_TIME_IN_DAYS",
			value:     "1",
		},
		{
			name:        "genuine parse error propagates",
			object:      Object{ObjectType: ObjectTypeDatabase, Name: NewAccountObjectIdentifier("some_db")},
			parameter:   "DATA_RETENTION_TIME_IN_DAYS",
			value:       "not-an-int",
			wantHandled: true,
			wantErr:     true,
		},
		{
			name:      "unknown database enum value falls back",
			object:    Object{ObjectType: ObjectTypeDatabase, Name: NewAccountObjectIdentifier("some_db")},
			parameter: "LOG_LEVEL",
			value:     "not-a-level",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := &parameters{client: &Client{}}
			handled, err := p.SetObjectParameterOnObjectViaCatalog(context.Background(), tc.object, tc.parameter, tc.value)
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tc.wantHandled, handled)
		})
	}
}

func TestUnsetObjectParameterOnObjectViaCatalog(t *testing.T) {
	cases := []struct {
		name        string
		object      Object
		parameter   ObjectParameter
		wantHandled bool
	}{
		{name: "unknown key falls back for warehouse", object: Object{ObjectType: ObjectTypeWarehouse, Name: NewAccountObjectIdentifier("some_object")}, parameter: "DATA_RETENTION_TIME_IN_DAYS"},
		{name: "unknown key falls back for user", object: Object{ObjectType: ObjectTypeUser, Name: NewAccountObjectIdentifier("some_object")}, parameter: "DATA_RETENTION_TIME_IN_DAYS"},
		{name: "function without arguments falls back", object: Object{ObjectType: ObjectTypeFunction, Name: NewSchemaObjectIdentifier("some_db", "some_schema", "some_fn")}, parameter: "LOG_LEVEL"},
		{name: "no catalog writer for table", object: Object{ObjectType: ObjectTypeTable, Name: NewSchemaObjectIdentifier("some_db", "some_schema", "some_table")}, parameter: "DATA_RETENTION_TIME_IN_DAYS"},
		{name: "unknown key falls back for schema", object: Object{ObjectType: ObjectTypeSchema, Name: NewDatabaseObjectIdentifier("some_db", "some_schema")}, parameter: "NOT_A_REAL_PARAMETER"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := &parameters{client: &Client{}}
			handled, err := p.UnsetObjectParameterOnObjectViaCatalog(context.Background(), tc.object, tc.parameter)
			require.NoError(t, err)
			require.Equal(t, tc.wantHandled, handled)
		})
	}
}

func TestSetObjectParameterOnAccountViaCatalog(t *testing.T) {
	cases := []struct {
		name        string
		parameter   ObjectParameter
		value       string
		wantHandled bool
		wantErr     bool
	}{
		{name: "unknown key falls back", parameter: "NOT_A_REAL_PARAMETER", value: "1"},
		{name: "genuine parse error propagates", parameter: "MIN_DATA_RETENTION_TIME_IN_DAYS", value: "not-an-int", wantHandled: true, wantErr: true},
		{name: "unknown enum value falls back", parameter: "LOG_LEVEL", value: "not-a-level"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := &parameters{client: &Client{}}
			handled, err := p.SetObjectParameterOnAccountViaCatalog(context.Background(), tc.parameter, tc.value)
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tc.wantHandled, handled)
		})
	}
}

func TestUnsetObjectParameterOnAccountViaCatalog(t *testing.T) {
	p := &parameters{client: &Client{}}
	handled, err := p.UnsetObjectParameterOnAccountViaCatalog(context.Background(), ObjectParameter("NOT_A_REAL_PARAMETER"))
	require.NoError(t, err)
	require.False(t, handled)
}

// The expected SQL of the legacy write path is the SQL sent by v2.22.x, and invalid values fail as in v2.22.x.
func TestSetParameterOnAccountLegacySql(t *testing.T) {
	cases := []struct {
		name      string
		formats   string
		parameter string
		value     string
		wantSql   string
		wantErr   string
	}{
		{name: "raw string keeps manual quotes", formats: "account", parameter: "PYTHON_PROFILER_MODULES", value: "'module_a,module_b'", wantSql: "ALTER ACCOUNT SET PYTHON_PROFILER_MODULES = 'module_a,module_b'"},
		{name: "raw identifier is not quoted", formats: "account", parameter: "CATALOG", value: "my_catalog", wantSql: "ALTER ACCOUNT SET CATALOG = my_catalog"},
		{name: "raw enum is not normalized", formats: "account", parameter: "LOG_LEVEL", value: "info", wantSql: "ALTER ACCOUNT SET LOG_LEVEL = info"},
		{name: "identifier as single-quoted string", formats: "account", parameter: "NETWORK_POLICY", value: "my_policy", wantSql: "ALTER ACCOUNT SET NETWORK_POLICY = 'my_policy'"},
		{name: "qualified identifier as single-quoted string", formats: "account", parameter: "EVENT_TABLE", value: "DB.SCHEMA.EVENTS", wantSql: "ALTER ACCOUNT SET EVENT_TABLE = 'DB.SCHEMA.EVENTS'"},
		{name: "single-quoted string escapes quotes", formats: "account", parameter: "QUERY_TAG", value: "it's", wantSql: `ALTER ACCOUNT SET QUERY_TAG = 'it\'s'`},
		{name: "enum single-quoted on the account parameter path", formats: "account", parameter: "TRACE_LEVEL", value: "ALWAYS", wantSql: "ALTER ACCOUNT SET TRACE_LEVEL = 'ALWAYS'"},
		{name: "bool is normalized", formats: "account", parameter: "ENABLE_PERSONAL_DATABASE", value: "TRUE", wantSql: "ALTER ACCOUNT SET ENABLE_PERSONAL_DATABASE = true"},
		{name: "bool accepts 1", formats: "account", parameter: "ENABLE_PERSONAL_DATABASE", value: "1", wantSql: "ALTER ACCOUNT SET ENABLE_PERSONAL_DATABASE = true"},
		{name: "int is normalized", formats: "account", parameter: "STATEMENT_TIMEOUT_IN_SECONDS", value: "042", wantSql: "ALTER ACCOUNT SET STATEMENT_TIMEOUT_IN_SECONDS = 42"},
		{name: "non-negative float", formats: "account", parameter: "INITIAL_REPLICATION_SIZE_LIMIT_IN_TB", value: "9.9", wantSql: "ALTER ACCOUNT SET INITIAL_REPLICATION_SIZE_LIMIT_IN_TB = 9.9"},
		// v2.22.x rejected this value in the provider; it is left to Snowflake now.
		{name: "int range is not checked", formats: "account", parameter: "CLIENT_ENCRYPTION_KEY_SIZE", value: "42", wantSql: "ALTER ACCOUNT SET CLIENT_ENCRYPTION_KEY_SIZE = 42"},
		{name: "invalid bool", formats: "account", parameter: "ENABLE_PERSONAL_DATABASE", value: "yes", wantErr: `boolean value ("true"/"false") expected for ENABLE_PERSONAL_DATABASE parameter, but got yes instead`},
		{name: "invalid int", formats: "account", parameter: "STATEMENT_TIMEOUT_IN_SECONDS", value: "9.9", wantErr: "failed to parse parameter STATEMENT_TIMEOUT_IN_SECONDS, expected an integer, but got 9.9"},
		{name: "negative float", formats: "account", parameter: "INITIAL_REPLICATION_SIZE_LIMIT_IN_TB", value: "-0.5", wantErr: "INITIAL_REPLICATION_SIZE_LIMIT_IN_TB must be a non-negative float, got -0.5"},
		{name: "object parameter: enum is raw", formats: "object", parameter: "TRACE_LEVEL", value: "ALWAYS", wantSql: "ALTER ACCOUNT SET TRACE_LEVEL = ALWAYS"},
		{name: "object parameter: identifier as single-quoted string", formats: "object", parameter: "NETWORK_POLICY", value: "my_policy", wantSql: "ALTER ACCOUNT SET NETWORK_POLICY = 'my_policy'"},
		{name: "object parameter: raw identifier", formats: "object", parameter: "CATALOG", value: "my_catalog", wantSql: "ALTER ACCOUNT SET CATALOG = my_catalog"},
		{name: "object parameter: int", formats: "object", parameter: "DATA_RETENTION_TIME_IN_DAYS", value: "1", wantSql: "ALTER ACCOUNT SET DATA_RETENTION_TIME_IN_DAYS = 1"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var sql string
			var err error
			if tc.formats == "object" {
				sql, err = legacySetParameterOnAccountSql(legacyObjectParameterOnAccountFormats, ObjectParameter(tc.parameter), tc.value)
			} else {
				sql, err = legacySetParameterOnAccountSql(legacyAccountParameterFormats, AccountParameter(tc.parameter), tc.value)
			}
			if tc.wantErr != "" {
				require.EqualError(t, err, tc.wantErr)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.wantSql, sql)
		})
	}
}

func TestUnsetParameterOnAccountLegacySql(t *testing.T) {
	sql, err := legacyUnsetParameterOnAccountSql(legacyAccountParameterFormats, AccountParameterTimezone)
	require.NoError(t, err)
	require.Equal(t, "ALTER ACCOUNT UNSET TIMEZONE", sql)

	sql, err = legacyUnsetParameterOnAccountSql(legacyObjectParameterOnAccountFormats, ObjectParameterNetworkPolicy)
	require.NoError(t, err)
	require.Equal(t, "ALTER ACCOUNT UNSET NETWORK_POLICY", sql)
}

func TestParameterOnAccountLegacy_NotSupported(t *testing.T) {
	// The errors are returned before any SQL is sent.
	p := &parameters{client: &Client{}}
	ctx := context.Background()

	t.Run("supported only by the catalog write path", func(t *testing.T) {
		require.ErrorIs(t, p.SetAccountParameterLegacy(ctx, AccountParameterDisableUiDownloadButton, "true"), ErrLegacyParameterNotSupported)
		require.ErrorIs(t, p.UnsetAccountParameterLegacy(ctx, AccountParameterDisableUiDownloadButton), ErrLegacyParameterNotSupported)
		require.ErrorIs(t, p.SetObjectParameterOnAccountLegacy(ctx, ObjectParameter("ENABLE_CONSOLE_OUTPUT"), "true"), ErrLegacyParameterNotSupported)
		require.ErrorIs(t, p.UnsetObjectParameterOnAccountLegacy(ctx, ObjectParameter("ENABLE_CONSOLE_OUTPUT")), ErrLegacyParameterNotSupported)
	})

	t.Run("not supported by any write path", func(t *testing.T) {
		err := p.SetAccountParameterLegacy(ctx, AccountParameterServiceCallerTokenValiditySecs, "60")
		require.ErrorIs(t, err, ErrParameterNotSupported)
		require.NotErrorIs(t, err, ErrLegacyParameterNotSupported)

		err = p.UnsetObjectParameterOnAccountLegacy(ctx, ObjectParameter("FALLBACK_WAREHOUSE"))
		require.ErrorIs(t, err, ErrParameterNotSupported)
		require.NotErrorIs(t, err, ErrLegacyParameterNotSupported)
	})
}
