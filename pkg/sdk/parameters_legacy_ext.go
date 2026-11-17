package sdk

import (
	"context"
	"errors"
	"fmt"
	"strconv"
)

// The legacy account parameter write path reproduces the SQL that the provider sent before the parameter catalog,
// up to v2.22.x. It is used by snowflake_account_parameter and snowflake_object_parameter (with on_account) when the
// ACCOUNT_PARAMETER_CATALOG_WRITE_PATH or OBJECT_PARAMETER_PARSE_ON_SET experiment is disabled. Do not add new
// parameters here: new parameters are supported only by the catalog write path (SetAccountParameter,
// UnsetAccountParameter, and the *ViaCatalog functions).
// The tables are frozen and reproduce the SQL sent by v2.22.1.

type legacyParameterFormat int

const (
	// legacyRaw sends the value as written, so string values must be quoted manually.
	legacyRaw legacyParameterFormat = iota
	// legacySingleQuoted sends the value as a single-quoted string literal.
	legacySingleQuoted
	legacyBool
	legacyInt
	// legacyNonNegativeFloat sends the value as written after checking that it is a non-negative float.
	legacyNonNegativeFloat
)

func (f legacyParameterFormat) render(parameter string, value string) (string, error) {
	switch f {
	case legacyRaw:
		return value, nil
	case legacySingleQuoted:
		return SingleQuotes.Modify(value), nil
	case legacyBool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return "", fmt.Errorf("boolean value (\"true\"/\"false\") expected for %v parameter, but got %v instead", parameter, value)
		}
		return strconv.FormatBool(b), nil
	case legacyInt:
		i, err := strconv.Atoi(value)
		if err != nil {
			return "", fmt.Errorf("failed to parse parameter %s, expected an integer, but got %v", parameter, value)
		}
		return strconv.Itoa(i), nil
	case legacyNonNegativeFloat:
		if value == "" {
			return value, nil
		}
		if f, err := strconv.ParseFloat(value, 64); err != nil || f < 0 {
			return "", fmt.Errorf("%s must be a non-negative float, got %v", parameter, value)
		}
		return value, nil
	default:
		return "", fmt.Errorf("unknown legacy parameter format %d", f)
	}
}

type legacyAccountParameterSet struct {
	alterAccount   bool   `ddl:"static" sql:"ALTER ACCOUNT"`
	set            bool   `ddl:"static" sql:"SET"`
	parameterKey   string `ddl:"keyword"`
	equals         bool   `ddl:"static" sql:"="`
	parameterValue string `ddl:"keyword"`
}

type legacyAccountParameterUnset struct {
	alterAccount bool   `ddl:"static" sql:"ALTER ACCOUNT"`
	unset        bool   `ddl:"static" sql:"UNSET"`
	parameterKey string `ddl:"keyword"`
}

// legacyParameterNotSupportedError returns ErrLegacyParameterNotSupported only when the catalog write path supports
// the parameter, so that enabling the experiment helps. Otherwise, it returns ErrParameterNotSupported.
func legacyParameterNotSupportedError[P ~string](parameter P) error {
	if errors.Is(NewAccountParametersUnsetRequest().UnsetParameterFromRaw(string(parameter)), ErrParameterNotSupported) {
		return fmt.Errorf("%w: %s", ErrParameterNotSupported, parameter)
	}
	return fmt.Errorf("%w: %s", ErrLegacyParameterNotSupported, parameter)
}

func legacySetParameterOnAccountSql[P ~string](formats map[P]legacyParameterFormat, parameter P, value string) (string, error) {
	format, ok := formats[parameter]
	if !ok {
		return "", legacyParameterNotSupportedError(parameter)
	}
	rendered, err := format.render(string(parameter), value)
	if err != nil {
		return "", err
	}
	return structToSQL(&legacyAccountParameterSet{parameterKey: string(parameter), parameterValue: rendered})
}

func legacyUnsetParameterOnAccountSql[P ~string](formats map[P]legacyParameterFormat, parameter P) (string, error) {
	if _, ok := formats[parameter]; !ok {
		return "", legacyParameterNotSupportedError(parameter)
	}
	return structToSQL(&legacyAccountParameterUnset{parameterKey: string(parameter)})
}

func (v *parameters) execLegacy(ctx context.Context, sql string, err error) error {
	if err != nil {
		return err
	}
	_, err = v.client.exec(ctx, sql)
	return err
}

// SetAccountParameterLegacy uses the legacy write path. SetAccountParameter is the catalog write path.
func (v *parameters) SetAccountParameterLegacy(ctx context.Context, parameter AccountParameter, value string) error {
	sql, err := legacySetParameterOnAccountSql(legacyAccountParameterFormats, parameter, value)
	return v.execLegacy(ctx, sql, err)
}

// UnsetAccountParameterLegacy uses the legacy write path. UnsetAccountParameter is the catalog write path.
func (v *parameters) UnsetAccountParameterLegacy(ctx context.Context, parameter AccountParameter) error {
	sql, err := legacyUnsetParameterOnAccountSql(legacyAccountParameterFormats, parameter)
	return v.execLegacy(ctx, sql, err)
}

// SetObjectParameterOnAccountLegacy uses the legacy write path. SetObjectParameterOnAccountViaCatalog is the catalog write path.
func (v *parameters) SetObjectParameterOnAccountLegacy(ctx context.Context, parameter ObjectParameter, value string) error {
	sql, err := legacySetParameterOnAccountSql(legacyObjectParameterOnAccountFormats, parameter, value)
	return v.execLegacy(ctx, sql, err)
}

// UnsetObjectParameterOnAccountLegacy uses the legacy write path. UnsetObjectParameterOnAccountViaCatalog is the catalog write path.
func (v *parameters) UnsetObjectParameterOnAccountLegacy(ctx context.Context, parameter ObjectParameter) error {
	sql, err := legacyUnsetParameterOnAccountSql(legacyObjectParameterOnAccountFormats, parameter)
	return v.execLegacy(ctx, sql, err)
}

// legacyAccountParameterFormats is the effective v2.22.1 account -> session -> object parameter fallback chain of SetAccountParameter (now SetAccountParameterLegacy).
var legacyAccountParameterFormats = map[AccountParameter]legacyParameterFormat{
	"ABORT_DETACHED_QUERY":                                            legacyBool,
	"ACTIVE_PYTHON_PROFILER":                                          legacySingleQuoted,
	"ALLOWED_SPCS_WORKLOAD_TYPES":                                     legacySingleQuoted,
	"ALLOW_BIND_VALUES_ACCESS":                                        legacyBool,
	"ALLOW_CLIENT_MFA_CACHING":                                        legacyBool,
	"ALLOW_ID_TOKEN":                                                  legacyBool,
	"AUTOCOMMIT":                                                      legacyBool,
	"BINARY_INPUT_FORMAT":                                             legacySingleQuoted,
	"BINARY_OUTPUT_FORMAT":                                            legacySingleQuoted,
	"CATALOG":                                                         legacyRaw,
	"CLIENT_ENABLE_LOG_INFO_STATEMENT_PARAMETERS":                     legacyBool,
	"CLIENT_ENCRYPTION_KEY_SIZE":                                      legacyInt,
	"CLIENT_MEMORY_LIMIT":                                             legacyInt,
	"CLIENT_METADATA_REQUEST_USE_CONNECTION_CTX":                      legacyBool,
	"CLIENT_METADATA_USE_SESSION_DATABASE":                            legacyBool,
	"CLIENT_PREFETCH_THREADS":                                         legacyInt,
	"CLIENT_RESULT_CHUNK_SIZE":                                        legacyInt,
	"CLIENT_RESULT_COLUMN_CASE_INSENSITIVE":                           legacyBool,
	"CLIENT_SESSION_KEEP_ALIVE":                                       legacyBool,
	"CLIENT_SESSION_KEEP_ALIVE_HEARTBEAT_FREQUENCY":                   legacyInt,
	"CLIENT_TIMESTAMP_TYPE_MAPPING":                                   legacySingleQuoted,
	"CORTEX_CODE_CLI_DAILY_EST_CREDIT_LIMIT_PER_USER":                 legacyInt,
	"CORTEX_CODE_DESKTOP_DAILY_EST_CREDIT_LIMIT_PER_USER":             legacyInt,
	"CORTEX_CODE_SNOWSIGHT_DAILY_EST_CREDIT_LIMIT_PER_USER":           legacyInt,
	"CORTEX_ENABLED_CROSS_REGION":                                     legacySingleQuoted,
	"CORTEX_MODELS_ALLOWLIST":                                         legacySingleQuoted,
	"CSV_TIMESTAMP_FORMAT":                                            legacySingleQuoted,
	"DATA_METRIC_SCHEDULE":                                            legacySingleQuoted,
	"DATA_RETENTION_TIME_IN_DAYS":                                     legacyInt,
	"DATE_INPUT_FORMAT":                                               legacySingleQuoted,
	"DATE_OUTPUT_FORMAT":                                              legacySingleQuoted,
	"DEFAULT_DBT_VERSION":                                             legacySingleQuoted,
	"DEFAULT_DDL_COLLATION":                                           legacySingleQuoted,
	"DEFAULT_NOTEBOOK_COMPUTE_POOL_CPU":                               legacySingleQuoted,
	"DEFAULT_NOTEBOOK_COMPUTE_POOL_GPU":                               legacySingleQuoted,
	"DEFAULT_STREAMLIT_COMPUTE_POOL":                                  legacySingleQuoted,
	"DISABLE_USER_PRIVILEGE_GRANTS":                                   legacyBool,
	"DISALLOWED_SPCS_WORKLOAD_TYPES":                                  legacySingleQuoted,
	"ENABLE_BUDGET_EVENT_LOGGING":                                     legacyBool,
	"ENABLE_CORTEX_ANALYST":                                           legacyBool,
	"ENABLE_DATA_COMPACTION":                                          legacyBool,
	"ENABLE_GET_DDL_USE_DATA_TYPE_ALIAS":                              legacyBool,
	"ENABLE_ICEBERG_MERGE_ON_READ":                                    legacyBool,
	"ENABLE_IDENTIFIER_FIRST_LOGIN":                                   legacyBool,
	"ENABLE_INTERNAL_STAGES_PRIVATELINK":                              legacyBool,
	"ENABLE_NOTEBOOK_CREATION_IN_PERSONAL_DB":                         legacyBool,
	"ENABLE_PERSONAL_DATABASE":                                        legacyBool,
	"ENABLE_PER_ACCOUNT_APP_SERVICE_PRIVATELINK_URL":                  legacyBool,
	"ENABLE_SPCS_BLOCK_STORAGE_SNOWFLAKE_FULL_ENCRYPTION_ENFORCEMENT": legacyBool,
	"ENABLE_TAG_PROPAGATION_EVENT_LOGGING":                            legacyBool,
	"ENABLE_TRI_SECRET_AND_REKEY_OPT_OUT_FOR_IMAGE_REPOSITORY":        legacyBool,
	"ENABLE_TRI_SECRET_AND_REKEY_OPT_OUT_FOR_SPCS_BLOCK_STORAGE":      legacyBool,
	"ENABLE_UNHANDLED_EXCEPTIONS_REPORTING":                           legacyBool,
	"ENABLE_UNLOAD_PHYSICAL_TYPE_OPTIMIZATION":                        legacyBool,
	"ENABLE_UNREDACTED_QUERY_SYNTAX_ERROR":                            legacyBool,
	"ENFORCE_NETWORK_RULES_FOR_INTERNAL_STAGES":                       legacyBool,
	"ERROR_ON_NONDETERMINISTIC_MERGE":                                 legacyBool,
	"ERROR_ON_NONDETERMINISTIC_UPDATE":                                legacyBool,
	"EVENT_TABLE":                                                     legacySingleQuoted,
	"EXTERNAL_OAUTH_ADD_PRIVILEGED_ROLES_TO_BLOCKED_LIST":             legacyBool,
	"GEOGRAPHY_OUTPUT_FORMAT":                                         legacySingleQuoted,
	"GEOMETRY_OUTPUT_FORMAT":                                          legacySingleQuoted,
	"HYBRID_TABLE_LOCK_TIMEOUT":                                       legacyInt,
	"ICEBERG_VERSION_DEFAULT":                                         legacyInt,
	"INITIAL_REPLICATION_SIZE_LIMIT_IN_TB":                            legacyNonNegativeFloat,
	"JDBC_TREAT_DECIMAL_AS_INT":                                       legacyBool,
	"JDBC_TREAT_TIMESTAMP_NTZ_AS_UTC":                                 legacyBool,
	"JDBC_USE_SESSION_TIMEZONE":                                       legacyBool,
	"JSON_INDENT":                                                     legacyInt,
	"JS_TREAT_INTEGER_AS_BIGINT":                                      legacyBool,
	"LOCK_TIMEOUT":                                                    legacyInt,
	"LOG_EVENT_LEVEL":                                                 legacyRaw,
	"LOG_LEVEL":                                                       legacyRaw,
	"MAX_CONCURRENCY_LEVEL":                                           legacyInt,
	"MAX_DATA_EXTENSION_TIME_IN_DAYS":                                 legacyInt,
	"METRIC_LEVEL":                                                    legacyRaw,
	"MIN_DATA_RETENTION_TIME_IN_DAYS":                                 legacyInt,
	"MULTI_STATEMENT_COUNT":                                           legacyInt,
	"NETWORK_POLICY":                                                  legacySingleQuoted,
	"NOORDER_SEQUENCE_AS_DEFAULT":                                     legacyBool,
	"OAUTH_ADD_PRIVILEGED_ROLES_TO_BLOCKED_LIST":                      legacyBool,
	"ODBC_TREAT_DECIMAL_AS_INT":                                       legacyBool,
	"PERIODIC_DATA_REKEYING":                                          legacyBool,
	"PIPE_EXECUTION_PAUSED":                                           legacyBool,
	"PREVENT_LOAD_FROM_INLINE_URL":                                    legacyBool,
	"PREVENT_UNLOAD_TO_INLINE_URL":                                    legacyBool,
	"PREVENT_UNLOAD_TO_INTERNAL_STAGES":                               legacyBool,
	"PYTHON_PROFILER_MODULES":                                         legacyRaw,
	"PYTHON_PROFILER_TARGET_STAGE":                                    legacyRaw,
	"QUERY_TAG":                                                       legacySingleQuoted,
	"QUOTED_IDENTIFIERS_IGNORE_CASE":                                  legacyBool,
	"READ_CONSISTENCY_MODE":                                           legacySingleQuoted,
	"REQUIRE_STORAGE_INTEGRATION_FOR_STAGE_CREATION":                  legacyBool,
	"REQUIRE_STORAGE_INTEGRATION_FOR_STAGE_OPERATION":                 legacyBool,
	"ROWS_PER_RESULTSET":                                              legacyInt,
	"ROW_TIMESTAMP_DEFAULT":                                           legacyBool,
	"S3_STAGE_VPCE_DNS_NAME":                                          legacySingleQuoted,
	"SEARCH_PATH":                                                     legacySingleQuoted,
	"SHARE_RESTRICTIONS":                                              legacyBool,
	"SIMULATED_DATA_SHARING_CONSUMER":                                 legacySingleQuoted,
	"SQL_TRACE_QUERY_TEXT":                                            legacySingleQuoted,
	"SSO_LOGIN_PAGE":                                                  legacyBool,
	"STATEMENT_QUEUED_TIMEOUT_IN_SECONDS":                             legacyInt,
	"STATEMENT_TIMEOUT_IN_SECONDS":                                    legacyInt,
	"STORAGE_SERIALIZATION_POLICY":                                    legacyRaw,
	"STRICT_JSON_OUTPUT":                                              legacyBool,
	"SUSPEND_TASK_AFTER_NUM_FAILURES":                                 legacyInt,
	"TASK_AUTO_RETRY_ATTEMPTS":                                        legacyInt,
	"TIMESTAMP_DAY_IS_ALWAYS_24H":                                     legacyBool,
	"TIMESTAMP_INPUT_FORMAT":                                          legacySingleQuoted,
	"TIMESTAMP_LTZ_OUTPUT_FORMAT":                                     legacySingleQuoted,
	"TIMESTAMP_NTZ_OUTPUT_FORMAT":                                     legacySingleQuoted,
	"TIMESTAMP_OUTPUT_FORMAT":                                         legacySingleQuoted,
	"TIMESTAMP_TYPE_MAPPING":                                          legacySingleQuoted,
	"TIMESTAMP_TZ_OUTPUT_FORMAT":                                      legacySingleQuoted,
	"TIMEZONE":                                                        legacySingleQuoted,
	"TIME_INPUT_FORMAT":                                               legacySingleQuoted,
	"TIME_OUTPUT_FORMAT":                                              legacySingleQuoted,
	"TRACE_LEVEL":                                                     legacySingleQuoted,
	"TRANSACTION_ABORT_ON_ERROR":                                      legacyBool,
	"TRANSACTION_DEFAULT_ISOLATION_LEVEL":                             legacySingleQuoted,
	"TWO_DIGIT_CENTURY_START":                                         legacyInt,
	"UNSUPPORTED_DDL_ACTION":                                          legacySingleQuoted,
	"USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE":                        legacyRaw,
	"USER_TASK_TIMEOUT_MS":                                            legacyInt,
	"USE_CACHED_RESULT":                                               legacyBool,
	"USE_WORKSPACES_FOR_SQL":                                          legacySingleQuoted,
	"WEEK_OF_YEAR_POLICY":                                             legacyInt,
	"WEEK_START":                                                      legacyInt,
}

// legacyObjectParameterOnAccountFormats is the v2.22.1 SetObjectParameterOnAccount (now SetObjectParameterOnAccountLegacy) mapping.
var legacyObjectParameterOnAccountFormats = map[ObjectParameter]legacyParameterFormat{
	"CATALOG":                                  legacyRaw,
	"DATA_METRIC_SCHEDULE":                     legacySingleQuoted,
	"DATA_RETENTION_TIME_IN_DAYS":              legacyInt,
	"DEFAULT_DDL_COLLATION":                    legacySingleQuoted,
	"DEFAULT_NOTEBOOK_COMPUTE_POOL_CPU":        legacySingleQuoted,
	"DEFAULT_NOTEBOOK_COMPUTE_POOL_GPU":        legacySingleQuoted,
	"ENABLE_DATA_COMPACTION":                   legacyBool,
	"ENABLE_ICEBERG_MERGE_ON_READ":             legacyBool,
	"ENABLE_NOTEBOOK_CREATION_IN_PERSONAL_DB":  legacyBool,
	"ENABLE_UNREDACTED_QUERY_SYNTAX_ERROR":     legacyBool,
	"ICEBERG_VERSION_DEFAULT":                  legacyInt,
	"LOG_EVENT_LEVEL":                          legacyRaw,
	"LOG_LEVEL":                                legacyRaw,
	"MAX_CONCURRENCY_LEVEL":                    legacyInt,
	"MAX_DATA_EXTENSION_TIME_IN_DAYS":          legacyInt,
	"NETWORK_POLICY":                           legacySingleQuoted,
	"PIPE_EXECUTION_PAUSED":                    legacyBool,
	"PREVENT_UNLOAD_TO_INTERNAL_STAGES":        legacyBool,
	"ROW_TIMESTAMP_DEFAULT":                    legacyBool,
	"SHARE_RESTRICTIONS":                       legacyBool,
	"STATEMENT_QUEUED_TIMEOUT_IN_SECONDS":      legacyInt,
	"STATEMENT_TIMEOUT_IN_SECONDS":             legacyInt,
	"STORAGE_SERIALIZATION_POLICY":             legacyRaw,
	"SUSPEND_TASK_AFTER_NUM_FAILURES":          legacyInt,
	"TASK_AUTO_RETRY_ATTEMPTS":                 legacyInt,
	"TRACE_LEVEL":                              legacyRaw,
	"USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE": legacyRaw,
	"USER_TASK_TIMEOUT_MS":                     legacyInt,
}
