// Code generated by assertions generator; DO NOT EDIT.

package resourceassert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

type TaskResourceAssert struct {
	*assert.ResourceAssert
}

func TaskResource(t *testing.T, name string) *TaskResourceAssert {
	t.Helper()

	return &TaskResourceAssert{
		ResourceAssert: assert.NewResourceAssert(name, "resource"),
	}
}

func ImportedTaskResource(t *testing.T, id string) *TaskResourceAssert {
	t.Helper()

	return &TaskResourceAssert{
		ResourceAssert: assert.NewImportedResourceAssert(id, "imported resource"),
	}
}

///////////////////////////////////
// Attribute value string checks //
///////////////////////////////////

func (t *TaskResourceAssert) HasAbortDetachedQueryString(expected string) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("abort_detached_query", expected))
	return t
}

func (t *TaskResourceAssert) HasAfterString(expected string) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("after", expected))
	return t
}

func (t *TaskResourceAssert) HasAllowOverlappingExecutionString(expected string) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("allow_overlapping_execution", expected))
	return t
}

func (t *TaskResourceAssert) HasAutocommitString(expected string) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("autocommit", expected))
	return t
}

func (t *TaskResourceAssert) HasBinaryInputFormatString(expected string) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("binary_input_format", expected))
	return t
}

func (t *TaskResourceAssert) HasBinaryOutputFormatString(expected string) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("binary_output_format", expected))
	return t
}

func (t *TaskResourceAssert) HasClientMemoryLimitString(expected string) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("client_memory_limit", expected))
	return t
}

func (t *TaskResourceAssert) HasClientMetadataRequestUseConnectionCtxString(expected string) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("client_metadata_request_use_connection_ctx", expected))
	return t
}

func (t *TaskResourceAssert) HasClientPrefetchThreadsString(expected string) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("client_prefetch_threads", expected))
	return t
}

func (t *TaskResourceAssert) HasClientResultChunkSizeString(expected string) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("client_result_chunk_size", expected))
	return t
}

func (t *TaskResourceAssert) HasClientResultColumnCaseInsensitiveString(expected string) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("client_result_column_case_insensitive", expected))
	return t
}

func (t *TaskResourceAssert) HasClientSessionKeepAliveString(expected string) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("client_session_keep_alive", expected))
	return t
}

func (t *TaskResourceAssert) HasClientSessionKeepAliveHeartbeatFrequencyString(expected string) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("client_session_keep_alive_heartbeat_frequency", expected))
	return t
}

func (t *TaskResourceAssert) HasClientTimestampTypeMappingString(expected string) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("client_timestamp_type_mapping", expected))
	return t
}

func (t *TaskResourceAssert) HasCommentString(expected string) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("comment", expected))
	return t
}

func (t *TaskResourceAssert) HasConfigString(expected string) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("config", expected))
	return t
}

func (t *TaskResourceAssert) HasDatabaseString(expected string) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("database", expected))
	return t
}

func (t *TaskResourceAssert) HasDateInputFormatString(expected string) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("date_input_format", expected))
	return t
}

func (t *TaskResourceAssert) HasDateOutputFormatString(expected string) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("date_output_format", expected))
	return t
}

func (t *TaskResourceAssert) HasEnableUnloadPhysicalTypeOptimizationString(expected string) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("enable_unload_physical_type_optimization", expected))
	return t
}

func (t *TaskResourceAssert) HasErrorIntegrationString(expected string) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("error_integration", expected))
	return t
}

func (t *TaskResourceAssert) HasErrorOnNondeterministicMergeString(expected string) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("error_on_nondeterministic_merge", expected))
	return t
}

func (t *TaskResourceAssert) HasErrorOnNondeterministicUpdateString(expected string) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("error_on_nondeterministic_update", expected))
	return t
}

func (t *TaskResourceAssert) HasFinalizeString(expected string) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("finalize", expected))
	return t
}

func (t *TaskResourceAssert) HasFullyQualifiedNameString(expected string) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("fully_qualified_name", expected))
	return t
}

func (t *TaskResourceAssert) HasGeographyOutputFormatString(expected string) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("geography_output_format", expected))
	return t
}

func (t *TaskResourceAssert) HasGeometryOutputFormatString(expected string) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("geometry_output_format", expected))
	return t
}

func (t *TaskResourceAssert) HasJdbcTreatTimestampNtzAsUtcString(expected string) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("jdbc_treat_timestamp_ntz_as_utc", expected))
	return t
}

func (t *TaskResourceAssert) HasJdbcUseSessionTimezoneString(expected string) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("jdbc_use_session_timezone", expected))
	return t
}

func (t *TaskResourceAssert) HasJsonIndentString(expected string) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("json_indent", expected))
	return t
}

func (t *TaskResourceAssert) HasLockTimeoutString(expected string) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("lock_timeout", expected))
	return t
}

func (t *TaskResourceAssert) HasLogLevelString(expected string) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("log_level", expected))
	return t
}

func (t *TaskResourceAssert) HasMultiStatementCountString(expected string) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("multi_statement_count", expected))
	return t
}

func (t *TaskResourceAssert) HasNameString(expected string) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("name", expected))
	return t
}

func (t *TaskResourceAssert) HasNoorderSequenceAsDefaultString(expected string) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("noorder_sequence_as_default", expected))
	return t
}

func (t *TaskResourceAssert) HasOdbcTreatDecimalAsIntString(expected string) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("odbc_treat_decimal_as_int", expected))
	return t
}

func (t *TaskResourceAssert) HasQueryTagString(expected string) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("query_tag", expected))
	return t
}

func (t *TaskResourceAssert) HasQuotedIdentifiersIgnoreCaseString(expected string) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("quoted_identifiers_ignore_case", expected))
	return t
}

func (t *TaskResourceAssert) HasRowsPerResultsetString(expected string) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("rows_per_resultset", expected))
	return t
}

func (t *TaskResourceAssert) HasS3StageVpceDnsNameString(expected string) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("s3_stage_vpce_dns_name", expected))
	return t
}

// TODO: Bring back
//func (t *TaskResourceAssert) HasScheduleString(expected string) *TaskResourceAssert {
//	t.AddAssertion(assert.ValueSet("schedule", expected))
//	return t
//}

func (t *TaskResourceAssert) HasSchemaString(expected string) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("schema", expected))
	return t
}

func (t *TaskResourceAssert) HasSearchPathString(expected string) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("search_path", expected))
	return t
}

func (t *TaskResourceAssert) HasSqlStatementString(expected string) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("sql_statement", expected))
	return t
}

func (t *TaskResourceAssert) HasStartedString(expected string) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("started", expected))
	return t
}

func (t *TaskResourceAssert) HasStatementQueuedTimeoutInSecondsString(expected string) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("statement_queued_timeout_in_seconds", expected))
	return t
}

func (t *TaskResourceAssert) HasStatementTimeoutInSecondsString(expected string) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("statement_timeout_in_seconds", expected))
	return t
}

func (t *TaskResourceAssert) HasStrictJsonOutputString(expected string) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("strict_json_output", expected))
	return t
}

func (t *TaskResourceAssert) HasSuspendTaskAfterNumFailuresString(expected string) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("suspend_task_after_num_failures", expected))
	return t
}

func (t *TaskResourceAssert) HasTaskAutoRetryAttemptsString(expected string) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("task_auto_retry_attempts", expected))
	return t
}

func (t *TaskResourceAssert) HasTimeInputFormatString(expected string) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("time_input_format", expected))
	return t
}

func (t *TaskResourceAssert) HasTimeOutputFormatString(expected string) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("time_output_format", expected))
	return t
}

func (t *TaskResourceAssert) HasTimestampDayIsAlways24hString(expected string) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("timestamp_day_is_always_24h", expected))
	return t
}

func (t *TaskResourceAssert) HasTimestampInputFormatString(expected string) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("timestamp_input_format", expected))
	return t
}

func (t *TaskResourceAssert) HasTimestampLtzOutputFormatString(expected string) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("timestamp_ltz_output_format", expected))
	return t
}

func (t *TaskResourceAssert) HasTimestampNtzOutputFormatString(expected string) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("timestamp_ntz_output_format", expected))
	return t
}

func (t *TaskResourceAssert) HasTimestampOutputFormatString(expected string) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("timestamp_output_format", expected))
	return t
}

func (t *TaskResourceAssert) HasTimestampTypeMappingString(expected string) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("timestamp_type_mapping", expected))
	return t
}

func (t *TaskResourceAssert) HasTimestampTzOutputFormatString(expected string) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("timestamp_tz_output_format", expected))
	return t
}

func (t *TaskResourceAssert) HasTimezoneString(expected string) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("timezone", expected))
	return t
}

func (t *TaskResourceAssert) HasTraceLevelString(expected string) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("trace_level", expected))
	return t
}

func (t *TaskResourceAssert) HasTransactionAbortOnErrorString(expected string) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("transaction_abort_on_error", expected))
	return t
}

func (t *TaskResourceAssert) HasTransactionDefaultIsolationLevelString(expected string) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("transaction_default_isolation_level", expected))
	return t
}

func (t *TaskResourceAssert) HasTwoDigitCenturyStartString(expected string) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("two_digit_century_start", expected))
	return t
}

func (t *TaskResourceAssert) HasUnsupportedDdlActionString(expected string) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("unsupported_ddl_action", expected))
	return t
}

func (t *TaskResourceAssert) HasUseCachedResultString(expected string) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("use_cached_result", expected))
	return t
}

func (t *TaskResourceAssert) HasUserTaskManagedInitialWarehouseSizeString(expected string) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("user_task_managed_initial_warehouse_size", expected))
	return t
}

func (t *TaskResourceAssert) HasUserTaskMinimumTriggerIntervalInSecondsString(expected string) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("user_task_minimum_trigger_interval_in_seconds", expected))
	return t
}

func (t *TaskResourceAssert) HasUserTaskTimeoutMsString(expected string) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("user_task_timeout_ms", expected))
	return t
}

func (t *TaskResourceAssert) HasWarehouseString(expected string) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("warehouse", expected))
	return t
}

func (t *TaskResourceAssert) HasWeekOfYearPolicyString(expected string) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("week_of_year_policy", expected))
	return t
}

func (t *TaskResourceAssert) HasWeekStartString(expected string) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("week_start", expected))
	return t
}

func (t *TaskResourceAssert) HasWhenString(expected string) *TaskResourceAssert {
	t.AddAssertion(assert.ValueSet("when", expected))
	return t
}

////////////////////////////
// Attribute empty checks //
////////////////////////////

func (t *TaskResourceAssert) HasNoAbortDetachedQuery() *TaskResourceAssert {
	t.AddAssertion(assert.ValueNotSet("abort_detached_query"))
	return t
}

func (t *TaskResourceAssert) HasNoAfter() *TaskResourceAssert {
	t.AddAssertion(assert.ValueNotSet("after"))
	return t
}

func (t *TaskResourceAssert) HasNoAllowOverlappingExecution() *TaskResourceAssert {
	t.AddAssertion(assert.ValueNotSet("allow_overlapping_execution"))
	return t
}

func (t *TaskResourceAssert) HasNoAutocommit() *TaskResourceAssert {
	t.AddAssertion(assert.ValueNotSet("autocommit"))
	return t
}

func (t *TaskResourceAssert) HasNoBinaryInputFormat() *TaskResourceAssert {
	t.AddAssertion(assert.ValueNotSet("binary_input_format"))
	return t
}

func (t *TaskResourceAssert) HasNoBinaryOutputFormat() *TaskResourceAssert {
	t.AddAssertion(assert.ValueNotSet("binary_output_format"))
	return t
}

func (t *TaskResourceAssert) HasNoClientMemoryLimit() *TaskResourceAssert {
	t.AddAssertion(assert.ValueNotSet("client_memory_limit"))
	return t
}

func (t *TaskResourceAssert) HasNoClientMetadataRequestUseConnectionCtx() *TaskResourceAssert {
	t.AddAssertion(assert.ValueNotSet("client_metadata_request_use_connection_ctx"))
	return t
}

func (t *TaskResourceAssert) HasNoClientPrefetchThreads() *TaskResourceAssert {
	t.AddAssertion(assert.ValueNotSet("client_prefetch_threads"))
	return t
}

func (t *TaskResourceAssert) HasNoClientResultChunkSize() *TaskResourceAssert {
	t.AddAssertion(assert.ValueNotSet("client_result_chunk_size"))
	return t
}

func (t *TaskResourceAssert) HasNoClientResultColumnCaseInsensitive() *TaskResourceAssert {
	t.AddAssertion(assert.ValueNotSet("client_result_column_case_insensitive"))
	return t
}

func (t *TaskResourceAssert) HasNoClientSessionKeepAlive() *TaskResourceAssert {
	t.AddAssertion(assert.ValueNotSet("client_session_keep_alive"))
	return t
}

func (t *TaskResourceAssert) HasNoClientSessionKeepAliveHeartbeatFrequency() *TaskResourceAssert {
	t.AddAssertion(assert.ValueNotSet("client_session_keep_alive_heartbeat_frequency"))
	return t
}

func (t *TaskResourceAssert) HasNoClientTimestampTypeMapping() *TaskResourceAssert {
	t.AddAssertion(assert.ValueNotSet("client_timestamp_type_mapping"))
	return t
}

func (t *TaskResourceAssert) HasNoComment() *TaskResourceAssert {
	t.AddAssertion(assert.ValueNotSet("comment"))
	return t
}

func (t *TaskResourceAssert) HasNoConfig() *TaskResourceAssert {
	t.AddAssertion(assert.ValueNotSet("config"))
	return t
}

func (t *TaskResourceAssert) HasNoDatabase() *TaskResourceAssert {
	t.AddAssertion(assert.ValueNotSet("database"))
	return t
}

func (t *TaskResourceAssert) HasNoDateInputFormat() *TaskResourceAssert {
	t.AddAssertion(assert.ValueNotSet("date_input_format"))
	return t
}

func (t *TaskResourceAssert) HasNoDateOutputFormat() *TaskResourceAssert {
	t.AddAssertion(assert.ValueNotSet("date_output_format"))
	return t
}

func (t *TaskResourceAssert) HasNoEnableUnloadPhysicalTypeOptimization() *TaskResourceAssert {
	t.AddAssertion(assert.ValueNotSet("enable_unload_physical_type_optimization"))
	return t
}

func (t *TaskResourceAssert) HasNoErrorIntegration() *TaskResourceAssert {
	t.AddAssertion(assert.ValueNotSet("error_integration"))
	return t
}

func (t *TaskResourceAssert) HasNoErrorOnNondeterministicMerge() *TaskResourceAssert {
	t.AddAssertion(assert.ValueNotSet("error_on_nondeterministic_merge"))
	return t
}

func (t *TaskResourceAssert) HasNoErrorOnNondeterministicUpdate() *TaskResourceAssert {
	t.AddAssertion(assert.ValueNotSet("error_on_nondeterministic_update"))
	return t
}

func (t *TaskResourceAssert) HasNoFinalize() *TaskResourceAssert {
	t.AddAssertion(assert.ValueNotSet("finalize"))
	return t
}

func (t *TaskResourceAssert) HasNoFullyQualifiedName() *TaskResourceAssert {
	t.AddAssertion(assert.ValueNotSet("fully_qualified_name"))
	return t
}

func (t *TaskResourceAssert) HasNoGeographyOutputFormat() *TaskResourceAssert {
	t.AddAssertion(assert.ValueNotSet("geography_output_format"))
	return t
}

func (t *TaskResourceAssert) HasNoGeometryOutputFormat() *TaskResourceAssert {
	t.AddAssertion(assert.ValueNotSet("geometry_output_format"))
	return t
}

func (t *TaskResourceAssert) HasNoJdbcTreatTimestampNtzAsUtc() *TaskResourceAssert {
	t.AddAssertion(assert.ValueNotSet("jdbc_treat_timestamp_ntz_as_utc"))
	return t
}

func (t *TaskResourceAssert) HasNoJdbcUseSessionTimezone() *TaskResourceAssert {
	t.AddAssertion(assert.ValueNotSet("jdbc_use_session_timezone"))
	return t
}

func (t *TaskResourceAssert) HasNoJsonIndent() *TaskResourceAssert {
	t.AddAssertion(assert.ValueNotSet("json_indent"))
	return t
}

func (t *TaskResourceAssert) HasNoLockTimeout() *TaskResourceAssert {
	t.AddAssertion(assert.ValueNotSet("lock_timeout"))
	return t
}

func (t *TaskResourceAssert) HasNoLogLevel() *TaskResourceAssert {
	t.AddAssertion(assert.ValueNotSet("log_level"))
	return t
}

func (t *TaskResourceAssert) HasNoMultiStatementCount() *TaskResourceAssert {
	t.AddAssertion(assert.ValueNotSet("multi_statement_count"))
	return t
}

func (t *TaskResourceAssert) HasNoName() *TaskResourceAssert {
	t.AddAssertion(assert.ValueNotSet("name"))
	return t
}

func (t *TaskResourceAssert) HasNoNoorderSequenceAsDefault() *TaskResourceAssert {
	t.AddAssertion(assert.ValueNotSet("noorder_sequence_as_default"))
	return t
}

func (t *TaskResourceAssert) HasNoOdbcTreatDecimalAsInt() *TaskResourceAssert {
	t.AddAssertion(assert.ValueNotSet("odbc_treat_decimal_as_int"))
	return t
}

func (t *TaskResourceAssert) HasNoQueryTag() *TaskResourceAssert {
	t.AddAssertion(assert.ValueNotSet("query_tag"))
	return t
}

func (t *TaskResourceAssert) HasNoQuotedIdentifiersIgnoreCase() *TaskResourceAssert {
	t.AddAssertion(assert.ValueNotSet("quoted_identifiers_ignore_case"))
	return t
}

func (t *TaskResourceAssert) HasNoRowsPerResultset() *TaskResourceAssert {
	t.AddAssertion(assert.ValueNotSet("rows_per_resultset"))
	return t
}

func (t *TaskResourceAssert) HasNoS3StageVpceDnsName() *TaskResourceAssert {
	t.AddAssertion(assert.ValueNotSet("s3_stage_vpce_dns_name"))
	return t
}

func (t *TaskResourceAssert) HasNoSchedule() *TaskResourceAssert {
	t.AddAssertion(assert.ValueNotSet("schedule"))
	return t
}

func (t *TaskResourceAssert) HasNoSchema() *TaskResourceAssert {
	t.AddAssertion(assert.ValueNotSet("schema"))
	return t
}

func (t *TaskResourceAssert) HasNoSearchPath() *TaskResourceAssert {
	t.AddAssertion(assert.ValueNotSet("search_path"))
	return t
}

func (t *TaskResourceAssert) HasNoSqlStatement() *TaskResourceAssert {
	t.AddAssertion(assert.ValueNotSet("sql_statement"))
	return t
}

func (t *TaskResourceAssert) HasNoStarted() *TaskResourceAssert {
	t.AddAssertion(assert.ValueNotSet("started"))
	return t
}

func (t *TaskResourceAssert) HasNoStatementQueuedTimeoutInSeconds() *TaskResourceAssert {
	t.AddAssertion(assert.ValueNotSet("statement_queued_timeout_in_seconds"))
	return t
}

func (t *TaskResourceAssert) HasNoStatementTimeoutInSeconds() *TaskResourceAssert {
	t.AddAssertion(assert.ValueNotSet("statement_timeout_in_seconds"))
	return t
}

func (t *TaskResourceAssert) HasNoStrictJsonOutput() *TaskResourceAssert {
	t.AddAssertion(assert.ValueNotSet("strict_json_output"))
	return t
}

func (t *TaskResourceAssert) HasNoSuspendTaskAfterNumFailures() *TaskResourceAssert {
	t.AddAssertion(assert.ValueNotSet("suspend_task_after_num_failures"))
	return t
}

func (t *TaskResourceAssert) HasNoTaskAutoRetryAttempts() *TaskResourceAssert {
	t.AddAssertion(assert.ValueNotSet("task_auto_retry_attempts"))
	return t
}

func (t *TaskResourceAssert) HasNoTimeInputFormat() *TaskResourceAssert {
	t.AddAssertion(assert.ValueNotSet("time_input_format"))
	return t
}

func (t *TaskResourceAssert) HasNoTimeOutputFormat() *TaskResourceAssert {
	t.AddAssertion(assert.ValueNotSet("time_output_format"))
	return t
}

func (t *TaskResourceAssert) HasNoTimestampDayIsAlways24h() *TaskResourceAssert {
	t.AddAssertion(assert.ValueNotSet("timestamp_day_is_always_24h"))
	return t
}

func (t *TaskResourceAssert) HasNoTimestampInputFormat() *TaskResourceAssert {
	t.AddAssertion(assert.ValueNotSet("timestamp_input_format"))
	return t
}

func (t *TaskResourceAssert) HasNoTimestampLtzOutputFormat() *TaskResourceAssert {
	t.AddAssertion(assert.ValueNotSet("timestamp_ltz_output_format"))
	return t
}

func (t *TaskResourceAssert) HasNoTimestampNtzOutputFormat() *TaskResourceAssert {
	t.AddAssertion(assert.ValueNotSet("timestamp_ntz_output_format"))
	return t
}

func (t *TaskResourceAssert) HasNoTimestampOutputFormat() *TaskResourceAssert {
	t.AddAssertion(assert.ValueNotSet("timestamp_output_format"))
	return t
}

func (t *TaskResourceAssert) HasNoTimestampTypeMapping() *TaskResourceAssert {
	t.AddAssertion(assert.ValueNotSet("timestamp_type_mapping"))
	return t
}

func (t *TaskResourceAssert) HasNoTimestampTzOutputFormat() *TaskResourceAssert {
	t.AddAssertion(assert.ValueNotSet("timestamp_tz_output_format"))
	return t
}

func (t *TaskResourceAssert) HasNoTimezone() *TaskResourceAssert {
	t.AddAssertion(assert.ValueNotSet("timezone"))
	return t
}

func (t *TaskResourceAssert) HasNoTraceLevel() *TaskResourceAssert {
	t.AddAssertion(assert.ValueNotSet("trace_level"))
	return t
}

func (t *TaskResourceAssert) HasNoTransactionAbortOnError() *TaskResourceAssert {
	t.AddAssertion(assert.ValueNotSet("transaction_abort_on_error"))
	return t
}

func (t *TaskResourceAssert) HasNoTransactionDefaultIsolationLevel() *TaskResourceAssert {
	t.AddAssertion(assert.ValueNotSet("transaction_default_isolation_level"))
	return t
}

func (t *TaskResourceAssert) HasNoTwoDigitCenturyStart() *TaskResourceAssert {
	t.AddAssertion(assert.ValueNotSet("two_digit_century_start"))
	return t
}

func (t *TaskResourceAssert) HasNoUnsupportedDdlAction() *TaskResourceAssert {
	t.AddAssertion(assert.ValueNotSet("unsupported_ddl_action"))
	return t
}

func (t *TaskResourceAssert) HasNoUseCachedResult() *TaskResourceAssert {
	t.AddAssertion(assert.ValueNotSet("use_cached_result"))
	return t
}

func (t *TaskResourceAssert) HasNoUserTaskManagedInitialWarehouseSize() *TaskResourceAssert {
	t.AddAssertion(assert.ValueNotSet("user_task_managed_initial_warehouse_size"))
	return t
}

func (t *TaskResourceAssert) HasNoUserTaskMinimumTriggerIntervalInSeconds() *TaskResourceAssert {
	t.AddAssertion(assert.ValueNotSet("user_task_minimum_trigger_interval_in_seconds"))
	return t
}

func (t *TaskResourceAssert) HasNoUserTaskTimeoutMs() *TaskResourceAssert {
	t.AddAssertion(assert.ValueNotSet("user_task_timeout_ms"))
	return t
}

func (t *TaskResourceAssert) HasNoWarehouse() *TaskResourceAssert {
	t.AddAssertion(assert.ValueNotSet("warehouse"))
	return t
}

func (t *TaskResourceAssert) HasNoWeekOfYearPolicy() *TaskResourceAssert {
	t.AddAssertion(assert.ValueNotSet("week_of_year_policy"))
	return t
}

func (t *TaskResourceAssert) HasNoWeekStart() *TaskResourceAssert {
	t.AddAssertion(assert.ValueNotSet("week_start"))
	return t
}

func (t *TaskResourceAssert) HasNoWhen() *TaskResourceAssert {
	t.AddAssertion(assert.ValueNotSet("when"))
	return t
}
