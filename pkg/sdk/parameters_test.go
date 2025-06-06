package sdk

import (
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
		assertOptsValidAndSQLEquals(t, opts, "ALTER USER %s SET ENABLE_UNREDACTED_QUERY_SYNTAX_ERROR = TRUE", id.FullyQualifiedName())
	})
}

func TestUnsetObjectParameterNetworkPolicyOnAccount(t *testing.T) {
	opts := &AlterAccountOptions{
		Unset: &AccountUnset{
			LegacyParameters: &AccountLevelParametersUnset{
				ObjectParameters: &ObjectParametersUnset{
					NetworkPolicy: Bool(true),
				},
			},
		},
	}
	t.Run("Unset Account Network Policy", func(t *testing.T) {
		assertOptsValidAndSQLEquals(t, opts, "ALTER ACCOUNT UNSET NETWORK_POLICY")
	})
}

func TestUnsetObjectParameterNetworkPolicyOnUser(t *testing.T) {
	opts := &AlterUserOptions{
		name: NewAccountObjectIdentifierFromFullyQualifiedName("TEST_USER"),
		Unset: &UserUnset{
			ObjectParameters: &UserObjectParametersUnset{
				NetworkPolicy: Bool(true),
			},
		},
	}
	t.Run("Unset User Network Policy", func(t *testing.T) {
		assertOptsValidAndSQLEquals(t, opts, `ALTER USER "TEST_USER" UNSET NETWORK_POLICY`)
	})
}

// Proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3344 is fixed.
func TestSetAccountParameterEnforceNetworkRulesForInternalStages(t *testing.T) {
	opts := &AlterAccountOptions{
		Set: &AccountSet{
			LegacyParameters: &AccountLevelParameters{
				AccountParameters: &AccountParameters{
					EnforceNetworkRulesForInternalStages: Bool(true),
				},
			},
		},
	}
	t.Run("Set Enforce Network Rules for Internal Stages", func(t *testing.T) {
		assertOptsValidAndSQLEquals(t, opts, "ALTER ACCOUNT SET ENFORCE_NETWORK_RULES_FOR_INTERNAL_STAGES = true")
	})
}

func TestToAccountParameter(t *testing.T) {
	type test struct {
		input string
		want  AccountParameter
	}

	valid := []test{
		// Case insensitive.
		{input: "allow_client_mfa_caching", want: AccountParameterAllowClientMFACaching},

		// Supported Values.
		{input: "ALLOW_CLIENT_MFA_CACHING", want: AccountParameterAllowClientMFACaching},
		{input: "ALLOW_ID_TOKEN", want: AccountParameterAllowIDToken},
		{input: "CLIENT_ENCRYPTION_KEY_SIZE", want: AccountParameterClientEncryptionKeySize},
		{input: "CORTEX_ENABLED_CROSS_REGION", want: AccountParameterCortexEnabledCrossRegion},
		{input: "DISABLE_USER_PRIVILEGE_GRANTS", want: AccountParameterDisableUserPrivilegeGrants},
		{input: "ENABLE_IDENTIFIER_FIRST_LOGIN", want: AccountParameterEnableIdentifierFirstLogin},
		{input: "ENABLE_INTERNAL_STAGES_PRIVATELINK", want: AccountParameterEnableInternalStagesPrivatelink},
		{input: "ENABLE_TRI_SECRET_AND_REKEY_OPT_OUT_FOR_IMAGE_REPOSITORY", want: AccountParameterEnableTriSecretAndRekeyOptOutForImageRepository},
		{input: "ENABLE_TRI_SECRET_AND_REKEY_OPT_OUT_FOR_SPCS_BLOCK_STORAGE", want: AccountParameterEnableTriSecretAndRekeyOptOutForSpcsBlockStorage},
		{input: "ENABLE_UNHANDLED_EXCEPTIONS_REPORTING", want: AccountParameterEnableUnhandledExceptionsReporting},
		{input: "ENFORCE_NETWORK_RULES_FOR_INTERNAL_STAGES", want: AccountParameterEnforceNetworkRulesForInternalStages},
		{input: "EVENT_TABLE", want: AccountParameterEventTable},
		{input: "EXTERNAL_OAUTH_ADD_PRIVILEGED_ROLES_TO_BLOCKED_LIST", want: AccountParameterExternalOAuthAddPrivilegedRolesToBlockedList},
		{input: "INITIAL_REPLICATION_SIZE_LIMIT_IN_TB", want: AccountParameterInitialReplicationSizeLimitInTB},
		{input: "MIN_DATA_RETENTION_TIME_IN_DAYS", want: AccountParameterMinDataRetentionTimeInDays},
		{input: "METRIC_LEVEL", want: AccountParameterMetricLevel},
		{input: "NETWORK_POLICY", want: AccountParameterNetworkPolicy},
		{input: "OAUTH_ADD_PRIVILEGED_ROLES_TO_BLOCKED_LIST", want: AccountParameterOAuthAddPrivilegedRolesToBlockedList},
		{input: "PERIODIC_DATA_REKEYING", want: AccountParameterPeriodicDataRekeying},
		{input: "PREVENT_LOAD_FROM_INLINE_URL", want: AccountParameterPreventLoadFromInlineURL},
		{input: "PREVENT_UNLOAD_TO_INLINE_URL", want: AccountParameterPreventUnloadToInlineURL},
		{input: "REQUIRE_STORAGE_INTEGRATION_FOR_STAGE_CREATION", want: AccountParameterRequireStorageIntegrationForStageCreation},
		{input: "REQUIRE_STORAGE_INTEGRATION_FOR_STAGE_OPERATION", want: AccountParameterRequireStorageIntegrationForStageOperation},
		{input: "SSO_LOGIN_PAGE", want: AccountParameterSSOLoginPage},

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
		{input: "JS_TREAT_INTEGER_AS_BIGINT", want: AccountParameterJsTreatIntegerAsBigInt},
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
		{input: "DEFAULT_DDL_COLLATION", want: AccountParameterDefaultDDLCollation},
		{input: "EXTERNAL_VOLUME", want: AccountParameterExternalVolume},
		{input: "LOG_LEVEL", want: AccountParameterLogLevel},
		{input: "MAX_CONCURRENCY_LEVEL", want: AccountParameterMaxConcurrencyLevel},
		{input: "MAX_DATA_EXTENSION_TIME_IN_DAYS", want: AccountParameterMaxDataExtensionTimeInDays},
		{input: "PIPE_EXECUTION_PAUSED", want: AccountParameterPipeExecutionPaused},
		{input: "REPLACE_INVALID_CHARACTERS", want: AccountParameterReplaceInvalidCharacters},
		{input: "STATEMENT_QUEUED_TIMEOUT_IN_SECONDS", want: AccountParameterStatementQueuedTimeoutInSeconds},
		{input: "STATEMENT_TIMEOUT_IN_SECONDS", want: AccountParameterStatementTimeoutInSeconds},
		{input: "STORAGE_SERIALIZATION_POLICY", want: AccountParameterStorageSerializationPolicy},
		{input: "SHARE_RESTRICTIONS", want: AccountParameterShareRestrictions},
		{input: "SUSPEND_TASK_AFTER_NUM_FAILURES", want: AccountParameterSuspendTaskAfterNumFailures},
		{input: "TRACE_LEVEL", want: AccountParameterTraceLevel},
		{input: "USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE", want: AccountParameterUserTaskManagedInitialWarehouseSize},
		{input: "USER_TASK_TIMEOUT_MS", want: AccountParameterUserTaskTimeoutMs},
		{input: "TASK_AUTO_RETRY_ATTEMPTS", want: AccountParameterTaskAutoRetryAttempts},
		{input: "USER_TASK_MINIMUM_TRIGGER_INTERVAL_IN_SECONDS", want: AccountParameterUserTaskMinimumTriggerIntervalInSeconds},
		{input: "METRIC_LEVEL", want: AccountParameterMetricLevel},
		{input: "ENABLE_CONSOLE_OUTPUT", want: AccountParameterEnableConsoleOutput},

		// User Parameters (inherited)
		{input: "ENABLE_PERSONAL_DATABASE", want: AccountParameterEnablePersonalDatabase},
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
