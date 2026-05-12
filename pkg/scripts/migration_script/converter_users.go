package main

import (
	"errors"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

var _ ConvertibleCsvRow[UserRepresentation] = new(UserCsvRow)

type UserCsvRow struct {
	Comment                                       string `csv:"comment"`
	CreatedOn                                     string `csv:"created_on"`
	DaysToExpiry                                  string `csv:"days_to_expiry"`
	DefaultNamespace                              string `csv:"default_namespace"`
	DefaultRole                                   string `csv:"default_role"`
	DefaultSecondaryRoles                         string `csv:"default_secondary_roles"`
	DefaultWarehouse                              string `csv:"default_warehouse"`
	Disabled                                      string `csv:"disabled"`
	DisplayName                                   string `csv:"display_name"`
	Email                                         string `csv:"email"`
	ExpiresAtTime                                 string `csv:"expires_at_time"`
	ExtAuthnDuo                                   string `csv:"ext_authn_duo"`
	ExtAuthnUid                                   string `csv:"ext_authn_uid"`
	CustomLandingPageUrl                          string `csv:"custom_landing_page_url"`
	CustomLandingPageUrlFlushNextUiLoad           string `csv:"custom_landing_page_url_flush_next_ui_load"`
	FirstName                                     string `csv:"first_name"`
	HasMfa                                        string `csv:"has_mfa"`
	HasPassword                                   string `csv:"has_password"`
	HasRsaPublicKey                               string `csv:"has_rsa_public_key"`
	HasWorkloadIdentity                           string `csv:"has_workload_identity"`
	LastName                                      string `csv:"last_name"`
	LastSuccessLogin                              string `csv:"last_success_login"`
	LockedUntilTime                               string `csv:"locked_until_time"`
	LoginName                                     string `csv:"login_name"`
	MiddleName                                    string `csv:"middle_name"`
	MinsToBypassMfa                               string `csv:"mins_to_bypass_mfa"`
	MinsToBypassNetworkPolicy                     string `csv:"mins_to_bypass_network_policy"`
	MinsToUnlock                                  string `csv:"mins_to_unlock"`
	MustChangePassword                            string `csv:"must_change_password"`
	Name                                          string `csv:"name"`
	Owner                                         string `csv:"owner"`
	Password                                      string `csv:"password"`
	PasswordLastSetTime                           string `csv:"password_last_set_time"`
	RsaPublicKey                                  string `csv:"rsa_public_key"`
	RsaPublicKey2                                 string `csv:"rsa_public_key2"`
	RsaPublicKey2Fp                               string `csv:"rsa_public_key2_fp"`
	RsaPublicKeyFp                                string `csv:"rsa_public_key_fp"`
	SnowflakeLock                                 string `csv:"snowflake_lock"`
	SnowflakeSupport                              string `csv:"snowflake_support"`
	Type                                          string `csv:"type"`
	AbortDetachedQueryLevel                       string `csv:"abort_detached_query_level"`
	AbortDetachedQueryValue                       string `csv:"abort_detached_query_value"`
	ActivePythonProfilerLevel                     string `csv:"active_python_profiler_level"`
	ActivePythonProfilerValue                     string `csv:"active_python_profiler_value"`
	AutocommitLevel                               string `csv:"autocommit_level"`
	AutocommitValue                               string `csv:"autocommit_value"`
	BinaryInputFormatLevel                        string `csv:"binary_input_format_level"`
	BinaryInputFormatValue                        string `csv:"binary_input_format_value"`
	BinaryOutputFormatLevel                       string `csv:"binary_output_format_level"`
	BinaryOutputFormatValue                       string `csv:"binary_output_format_value"`
	ClientEnableLogInfoStatementParametersLevel   string `csv:"client_enable_log_info_statement_parameters_level"`
	ClientEnableLogInfoStatementParametersValue   string `csv:"client_enable_log_info_statement_parameters_value"`
	ClientMemoryLimitLevel                        string `csv:"client_memory_limit_level"`
	ClientMemoryLimitValue                        string `csv:"client_memory_limit_value"`
	ClientMetadataRequestUseConnectionCtxLevel    string `csv:"client_metadata_request_use_connection_ctx_level"`
	ClientMetadataRequestUseConnectionCtxValue    string `csv:"client_metadata_request_use_connection_ctx_value"`
	ClientMetadataUseSessionDatabaseLevel         string `csv:"client_metadata_use_session_database_level"`
	ClientMetadataUseSessionDatabaseValue         string `csv:"client_metadata_use_session_database_value"`
	ClientPrefetchThreadsLevel                    string `csv:"client_prefetch_threads_level"`
	ClientPrefetchThreadsValue                    string `csv:"client_prefetch_threads_value"`
	ClientResultChunkSizeLevel                    string `csv:"client_result_chunk_size_level"`
	ClientResultChunkSizeValue                    string `csv:"client_result_chunk_size_value"`
	ClientResultColumnCaseInsensitiveLevel        string `csv:"client_result_column_case_insensitive_level"`
	ClientResultColumnCaseInsensitiveValue        string `csv:"client_result_column_case_insensitive_value"`
	ClientSessionKeepAliveLevel                   string `csv:"client_session_keep_alive_level"`
	ClientSessionKeepAliveValue                   string `csv:"client_session_keep_alive_value"`
	ClientSessionKeepAliveHeartbeatFrequencyLevel string `csv:"client_session_keep_alive_heartbeat_frequency_level"`
	ClientSessionKeepAliveHeartbeatFrequencyValue string `csv:"client_session_keep_alive_heartbeat_frequency_value"`
	ClientTimestampTypeMappingLevel               string `csv:"client_timestamp_type_mapping_level"`
	ClientTimestampTypeMappingValue               string `csv:"client_timestamp_type_mapping_value"`
	CsvTimestampFormatLevel                       string `csv:"csv_timestamp_format_level"`
	CsvTimestampFormatValue                       string `csv:"csv_timestamp_format_value"`
	DateInputFormatLevel                          string `csv:"date_input_format_level"`
	DateInputFormatValue                          string `csv:"date_input_format_value"`
	DateOutputFormatLevel                         string `csv:"date_output_format_level"`
	DateOutputFormatValue                         string `csv:"date_output_format_value"`
	EnableUnloadPhysicalTypeOptimizationLevel     string `csv:"enable_unload_physical_type_optimization_level"`
	EnableUnloadPhysicalTypeOptimizationValue     string `csv:"enable_unload_physical_type_optimization_value"`
	EnableUnredactedQuerySyntaxErrorLevel         string `csv:"enable_unredacted_query_syntax_error_level"`
	EnableUnredactedQuerySyntaxErrorValue         string `csv:"enable_unredacted_query_syntax_error_value"`
	ErrorOnNondeterministicMergeLevel             string `csv:"error_on_nondeterministic_merge_level"`
	ErrorOnNondeterministicMergeValue             string `csv:"error_on_nondeterministic_merge_value"`
	ErrorOnNondeterministicUpdateLevel            string `csv:"error_on_nondeterministic_update_level"`
	ErrorOnNondeterministicUpdateValue            string `csv:"error_on_nondeterministic_update_value"`
	GeographyOutputFormatLevel                    string `csv:"geography_output_format_level"`
	GeographyOutputFormatValue                    string `csv:"geography_output_format_value"`
	GeometryOutputFormatLevel                     string `csv:"geometry_output_format_level"`
	GeometryOutputFormatValue                     string `csv:"geometry_output_format_value"`
	HybridTableLockTimeoutLevel                   string `csv:"hybrid_table_lock_timeout_level"`
	HybridTableLockTimeoutValue                   string `csv:"hybrid_table_lock_timeout_value"`
	JdbcTreatDecimalAsIntLevel                    string `csv:"jdbc_treat_decimal_as_int_level"`
	JdbcTreatDecimalAsIntValue                    string `csv:"jdbc_treat_decimal_as_int_value"`
	JdbcTreatTimestampNtzAsUtcLevel               string `csv:"jdbc_treat_timestamp_ntz_as_utc_level"`
	JdbcTreatTimestampNtzAsUtcValue               string `csv:"jdbc_treat_timestamp_ntz_as_utc_value"`
	JdbcUseSessionTimezoneLevel                   string `csv:"jdbc_use_session_timezone_level"`
	JdbcUseSessionTimezoneValue                   string `csv:"jdbc_use_session_timezone_value"`
	JsTreatIntegerAsBigIntLevel                   string `csv:"js_treat_integer_as_big_int_level"`
	JsTreatIntegerAsBigIntValue                   string `csv:"js_treat_integer_as_big_int_value"`
	JsonIndentLevel                               string `csv:"json_indent_level"`
	JsonIndentValue                               string `csv:"json_indent_value"`
	LockTimeoutLevel                              string `csv:"lock_timeout_level"`
	LockTimeoutValue                              string `csv:"lock_timeout_value"`
	LogLevelLevel                                 string `csv:"log_level_level"`
	LogLevelValue                                 string `csv:"log_level_value"`
	MultiStatementCountLevel                      string `csv:"multi_statement_count_level"`
	MultiStatementCountValue                      string `csv:"multi_statement_count_value"`
	NetworkPolicyLevel                            string `csv:"network_policy_level"`
	NetworkPolicyValue                            string `csv:"network_policy_value"`
	NoorderSequenceAsDefaultLevel                 string `csv:"noorder_sequence_as_default_level"`
	NoorderSequenceAsDefaultValue                 string `csv:"noorder_sequence_as_default_value"`
	OdbcTreatDecimalAsIntLevel                    string `csv:"odbc_treat_decimal_as_int_level"`
	OdbcTreatDecimalAsIntValue                    string `csv:"odbc_treat_decimal_as_int_value"`
	PreventUnloadToInternalStagesLevel            string `csv:"prevent_unload_to_internal_stages_level"`
	PreventUnloadToInternalStagesValue            string `csv:"prevent_unload_to_internal_stages_value"`
	PythonProfilerModulesLevel                    string `csv:"python_profiler_modules_level"`
	PythonProfilerModulesValue                    string `csv:"python_profiler_modules_value"`
	PythonProfilerTargetStageLevel                string `csv:"python_profiler_target_stage_level"`
	PythonProfilerTargetStageValue                string `csv:"python_profiler_target_stage_value"`
	QueryTagLevel                                 string `csv:"query_tag_level"`
	QueryTagValue                                 string `csv:"query_tag_value"`
	QuotedIdentifiersIgnoreCaseLevel              string `csv:"quoted_identifiers_ignore_case_level"`
	QuotedIdentifiersIgnoreCaseValue              string `csv:"quoted_identifiers_ignore_case_value"`
	RowsPerResultsetLevel                         string `csv:"rows_per_resultset_level"`
	RowsPerResultsetValue                         string `csv:"rows_per_resultset_value"`
	S3StageVpceDnsNameLevel                       string `csv:"s3_stage_vpce_dns_name_level"`
	S3StageVpceDnsNameValue                       string `csv:"s3_stage_vpce_dns_name_value"`
	SearchPathLevel                               string `csv:"search_path_level"`
	SearchPathValue                               string `csv:"search_path_value"`
	SimulatedDataSharingConsumerLevel             string `csv:"simulated_data_sharing_consumer_level"`
	SimulatedDataSharingConsumerValue             string `csv:"simulated_data_sharing_consumer_value"`
	StatementQueuedTimeoutInSecondsLevel          string `csv:"statement_queued_timeout_in_seconds_level"`
	StatementQueuedTimeoutInSecondsValue          string `csv:"statement_queued_timeout_in_seconds_value"`
	StatementTimeoutInSecondsLevel                string `csv:"statement_timeout_in_seconds_level"`
	StatementTimeoutInSecondsValue                string `csv:"statement_timeout_in_seconds_value"`
	StrictJsonOutputLevel                         string `csv:"strict_json_output_level"`
	StrictJsonOutputValue                         string `csv:"strict_json_output_value"`
	TimeInputFormatLevel                          string `csv:"time_input_format_level"`
	TimeInputFormatValue                          string `csv:"time_input_format_value"`
	TimeOutputFormatLevel                         string `csv:"time_output_format_level"`
	TimeOutputFormatValue                         string `csv:"time_output_format_value"`
	TimestampDayIsAlways24hLevel                  string `csv:"timestamp_day_is_always24h_level"`
	TimestampDayIsAlways24hValue                  string `csv:"timestamp_day_is_always24h_value"`
	TimestampInputFormatLevel                     string `csv:"timestamp_input_format_level"`
	TimestampInputFormatValue                     string `csv:"timestamp_input_format_value"`
	TimestampLTZOutputFormatLevel                 string `csv:"timestamp_l_t_z_output_format_level"`
	TimestampLTZOutputFormatValue                 string `csv:"timestamp_l_t_z_output_format_value"`
	TimestampNTZOutputFormatLevel                 string `csv:"timestamp_n_t_z_output_format_level"`
	TimestampNTZOutputFormatValue                 string `csv:"timestamp_n_t_z_output_format_value"`
	TimestampOutputFormatLevel                    string `csv:"timestamp_output_format_level"`
	TimestampOutputFormatValue                    string `csv:"timestamp_output_format_value"`
	TimestampTZOutputFormatLevel                  string `csv:"timestamp_t_z_output_format_level"`
	TimestampTZOutputFormatValue                  string `csv:"timestamp_t_z_output_format_value"`
	TimestampTypeMappingLevel                     string `csv:"timestamp_type_mapping_level"`
	TimestampTypeMappingValue                     string `csv:"timestamp_type_mapping_value"`
	TimezoneLevel                                 string `csv:"timezone_level"`
	TimezoneValue                                 string `csv:"timezone_value"`
	TraceLevelLevel                               string `csv:"trace_level_level"`
	TraceLevelValue                               string `csv:"trace_level_value"`
	TransactionAbortOnErrorLevel                  string `csv:"transaction_abort_on_error_level"`
	TransactionAbortOnErrorValue                  string `csv:"transaction_abort_on_error_value"`
	TransactionDefaultIsolationLevelLevel         string `csv:"transaction_default_isolation_level_level"`
	TransactionDefaultIsolationLevelValue         string `csv:"transaction_default_isolation_level_value"`
	TwoDigitCenturyStartLevel                     string `csv:"two_digit_century_start_level"`
	TwoDigitCenturyStartValue                     string `csv:"two_digit_century_start_value"`
	UnsupportedDDLActionLevel                     string `csv:"unsupported_d_d_l_action_level"`
	UnsupportedDDLActionValue                     string `csv:"unsupported_d_d_l_action_value"`
	UseCachedResultLevel                          string `csv:"use_cached_result_level"`
	UseCachedResultValue                          string `csv:"use_cached_result_value"`
	WeekOfYearPolicyLevel                         string `csv:"week_of_year_policy_level"`
	WeekOfYearPolicyValue                         string `csv:"week_of_year_policy_value"`
	WeekStartLevel                                string `csv:"week_start_level"`
	WeekStartValue                                string `csv:"week_start_value"`
}

type UserRepresentation struct {
	sdk.User

	// describe output fields (not in sdk.User)
	MiddleName                          string
	Password                            string
	PasswordLastSetTime                 string
	SnowflakeSupport                    bool
	MinsToBypassNetworkPolicy           string
	RsaPublicKey                        string
	RsaPublicKeyFp                      string
	RsaPublicKey2                       string
	RsaPublicKey2Fp                     string
	CustomLandingPageUrl                string
	CustomLandingPageUrlFlushNextUiLoad bool

	// parameters
	AbortDetachedQuery                       *bool
	ActivePythonProfiler                     *string
	Autocommit                               *bool
	BinaryInputFormat                        *string
	BinaryOutputFormat                       *string
	ClientEnableLogInfoStatementParameters   *bool
	ClientMemoryLimit                        *int
	ClientMetadataRequestUseConnectionCtx    *bool
	ClientMetadataUseSessionDatabase         *bool
	ClientPrefetchThreads                    *int
	ClientResultChunkSize                    *int
	ClientResultColumnCaseInsensitive        *bool
	ClientSessionKeepAlive                   *bool
	ClientSessionKeepAliveHeartbeatFrequency *int
	ClientTimestampTypeMapping               *string
	CsvTimestampFormat                       *string
	DateInputFormat                          *string
	DateOutputFormat                         *string
	EnableUnloadPhysicalTypeOptimization     *bool
	EnableUnredactedQuerySyntaxError         *bool
	ErrorOnNondeterministicMerge             *bool
	ErrorOnNondeterministicUpdate            *bool
	GeographyOutputFormat                    *string
	GeometryOutputFormat                     *string
	HybridTableLockTimeout                   *int
	JdbcTreatDecimalAsInt                    *bool
	JdbcTreatTimestampNtzAsUtc               *bool
	JdbcUseSessionTimezone                   *bool
	JsTreatIntegerAsBigInt                   *bool
	JsonIndent                               *int
	LockTimeout                              *int
	LogLevel                                 *string
	MultiStatementCount                      *int
	NetworkPolicy                            *string
	NoorderSequenceAsDefault                 *bool
	OdbcTreatDecimalAsInt                    *bool
	PreventUnloadToInternalStages            *bool
	PythonProfilerModules                    *string
	PythonProfilerTargetStage                *string
	QueryTag                                 *string
	QuotedIdentifiersIgnoreCase              *bool
	RowsPerResultset                         *int
	S3StageVpceDnsName                       *string
	SearchPath                               *string
	SimulatedDataSharingConsumer             *string
	StatementQueuedTimeoutInSeconds          *int
	StatementTimeoutInSeconds                *int
	StrictJsonOutput                         *bool
	TimeInputFormat                          *string
	TimeOutputFormat                         *string
	TimestampDayIsAlways24h                  *bool
	TimestampInputFormat                     *string
	TimestampLTZOutputFormat                 *string
	TimestampNTZOutputFormat                 *string
	TimestampOutputFormat                    *string
	TimestampTZOutputFormat                  *string
	TimestampTypeMapping                     *string
	Timezone                                 *string
	TraceLevel                               *string
	TransactionAbortOnError                  *bool
	TransactionDefaultIsolationLevel         *string
	TwoDigitCenturyStart                     *int
	UnsupportedDDLAction                     *string
	UseCachedResult                          *bool
	WeekOfYearPolicy                         *int
	WeekStart                                *int
}

func (row UserCsvRow) convert() (*UserRepresentation, error) {
	userRepresentation := &UserRepresentation{
		User: sdk.User{
			Name:                  row.Name,
			LoginName:             row.LoginName,
			DisplayName:           row.DisplayName,
			FirstName:             row.FirstName,
			LastName:              row.LastName,
			Email:                 row.Email,
			MinsToUnlock:          row.MinsToUnlock,
			DaysToExpiry:          row.DaysToExpiry,
			Comment:               csvUnescape(row.Comment),
			Disabled:              row.Disabled == "true",
			MustChangePassword:    row.MustChangePassword == "true",
			SnowflakeLock:         row.SnowflakeLock == "true",
			DefaultWarehouse:      row.DefaultWarehouse,
			DefaultNamespace:      row.DefaultNamespace,
			DefaultRole:           row.DefaultRole,
			DefaultSecondaryRoles: row.DefaultSecondaryRoles,
			ExtAuthnDuo:           row.ExtAuthnDuo == "true",
			ExtAuthnUid:           row.ExtAuthnUid,
			MinsToBypassMfa:       row.MinsToBypassMfa,
			Owner:                 row.Owner,
			HasPassword:           row.HasPassword == "true",
			HasRsaPublicKey:       row.HasRsaPublicKey == "true",
			Type:                  row.Type,
			HasMfa:                row.HasMfa == "true",
			HasWorkloadIdentity:   row.HasWorkloadIdentity == "true",
		},
		// describe output fields (apply csvUnescape for fields that may contain escape sequences)
		MiddleName:                          row.MiddleName,
		Password:                            row.Password,
		PasswordLastSetTime:                 row.PasswordLastSetTime,
		SnowflakeSupport:                    row.SnowflakeSupport == "true",
		MinsToBypassNetworkPolicy:           row.MinsToBypassNetworkPolicy,
		RsaPublicKey:                        csvUnescape(row.RsaPublicKey),
		RsaPublicKeyFp:                      row.RsaPublicKeyFp,
		RsaPublicKey2:                       csvUnescape(row.RsaPublicKey2),
		RsaPublicKey2Fp:                     row.RsaPublicKey2Fp,
		CustomLandingPageUrl:                row.CustomLandingPageUrl,
		CustomLandingPageUrlFlushNextUiLoad: row.CustomLandingPageUrlFlushNextUiLoad == "true",
	}

	handler := newParameterHandler(sdk.ParameterTypeUser)
	errs := errors.Join(
		handler.handleBooleanParameter(row.AbortDetachedQueryLevel, row.AbortDetachedQueryValue, &userRepresentation.AbortDetachedQuery),
		handler.handleStringParameter(row.ActivePythonProfilerLevel, row.ActivePythonProfilerValue, &userRepresentation.ActivePythonProfiler),
		handler.handleBooleanParameter(row.AutocommitLevel, row.AutocommitValue, &userRepresentation.Autocommit),
		handler.handleStringParameter(row.BinaryInputFormatLevel, row.BinaryInputFormatValue, &userRepresentation.BinaryInputFormat),
		handler.handleStringParameter(row.BinaryOutputFormatLevel, row.BinaryOutputFormatValue, &userRepresentation.BinaryOutputFormat),
		handler.handleBooleanParameter(row.ClientEnableLogInfoStatementParametersLevel, row.ClientEnableLogInfoStatementParametersValue, &userRepresentation.ClientEnableLogInfoStatementParameters),
		handler.handleIntegerParameter(row.ClientMemoryLimitLevel, row.ClientMemoryLimitValue, &userRepresentation.ClientMemoryLimit),
		handler.handleBooleanParameter(row.ClientMetadataRequestUseConnectionCtxLevel, row.ClientMetadataRequestUseConnectionCtxValue, &userRepresentation.ClientMetadataRequestUseConnectionCtx),
		handler.handleBooleanParameter(row.ClientMetadataUseSessionDatabaseLevel, row.ClientMetadataUseSessionDatabaseValue, &userRepresentation.ClientMetadataUseSessionDatabase),
		handler.handleIntegerParameter(row.ClientPrefetchThreadsLevel, row.ClientPrefetchThreadsValue, &userRepresentation.ClientPrefetchThreads),
		handler.handleIntegerParameter(row.ClientResultChunkSizeLevel, row.ClientResultChunkSizeValue, &userRepresentation.ClientResultChunkSize),
		handler.handleBooleanParameter(row.ClientResultColumnCaseInsensitiveLevel, row.ClientResultColumnCaseInsensitiveValue, &userRepresentation.ClientResultColumnCaseInsensitive),
		handler.handleBooleanParameter(row.ClientSessionKeepAliveLevel, row.ClientSessionKeepAliveValue, &userRepresentation.ClientSessionKeepAlive),
		handler.handleIntegerParameter(row.ClientSessionKeepAliveHeartbeatFrequencyLevel, row.ClientSessionKeepAliveHeartbeatFrequencyValue, &userRepresentation.ClientSessionKeepAliveHeartbeatFrequency),
		handler.handleStringParameter(row.ClientTimestampTypeMappingLevel, row.ClientTimestampTypeMappingValue, &userRepresentation.ClientTimestampTypeMapping),
		handler.handleStringParameter(row.CsvTimestampFormatLevel, row.CsvTimestampFormatValue, &userRepresentation.CsvTimestampFormat),
		handler.handleStringParameter(row.DateInputFormatLevel, row.DateInputFormatValue, &userRepresentation.DateInputFormat),
		handler.handleStringParameter(row.DateOutputFormatLevel, row.DateOutputFormatValue, &userRepresentation.DateOutputFormat),
		handler.handleBooleanParameter(row.EnableUnloadPhysicalTypeOptimizationLevel, row.EnableUnloadPhysicalTypeOptimizationValue, &userRepresentation.EnableUnloadPhysicalTypeOptimization),
		handler.handleBooleanParameter(row.EnableUnredactedQuerySyntaxErrorLevel, row.EnableUnredactedQuerySyntaxErrorValue, &userRepresentation.EnableUnredactedQuerySyntaxError),
		handler.handleBooleanParameter(row.ErrorOnNondeterministicMergeLevel, row.ErrorOnNondeterministicMergeValue, &userRepresentation.ErrorOnNondeterministicMerge),
		handler.handleBooleanParameter(row.ErrorOnNondeterministicUpdateLevel, row.ErrorOnNondeterministicUpdateValue, &userRepresentation.ErrorOnNondeterministicUpdate),
		handler.handleStringParameter(row.GeographyOutputFormatLevel, row.GeographyOutputFormatValue, &userRepresentation.GeographyOutputFormat),
		handler.handleStringParameter(row.GeometryOutputFormatLevel, row.GeometryOutputFormatValue, &userRepresentation.GeometryOutputFormat),
		handler.handleIntegerParameter(row.HybridTableLockTimeoutLevel, row.HybridTableLockTimeoutValue, &userRepresentation.HybridTableLockTimeout),
		handler.handleBooleanParameter(row.JdbcTreatDecimalAsIntLevel, row.JdbcTreatDecimalAsIntValue, &userRepresentation.JdbcTreatDecimalAsInt),
		handler.handleBooleanParameter(row.JdbcTreatTimestampNtzAsUtcLevel, row.JdbcTreatTimestampNtzAsUtcValue, &userRepresentation.JdbcTreatTimestampNtzAsUtc),
		handler.handleBooleanParameter(row.JdbcUseSessionTimezoneLevel, row.JdbcUseSessionTimezoneValue, &userRepresentation.JdbcUseSessionTimezone),
		handler.handleBooleanParameter(row.JsTreatIntegerAsBigIntLevel, row.JsTreatIntegerAsBigIntValue, &userRepresentation.JsTreatIntegerAsBigInt),
		handler.handleIntegerParameter(row.JsonIndentLevel, row.JsonIndentValue, &userRepresentation.JsonIndent),
		handler.handleIntegerParameter(row.LockTimeoutLevel, row.LockTimeoutValue, &userRepresentation.LockTimeout),
		handler.handleStringParameter(row.LogLevelLevel, row.LogLevelValue, &userRepresentation.LogLevel),
		handler.handleIntegerParameter(row.MultiStatementCountLevel, row.MultiStatementCountValue, &userRepresentation.MultiStatementCount),
		handler.handleStringParameter(row.NetworkPolicyLevel, row.NetworkPolicyValue, &userRepresentation.NetworkPolicy),
		handler.handleBooleanParameter(row.NoorderSequenceAsDefaultLevel, row.NoorderSequenceAsDefaultValue, &userRepresentation.NoorderSequenceAsDefault),
		handler.handleBooleanParameter(row.OdbcTreatDecimalAsIntLevel, row.OdbcTreatDecimalAsIntValue, &userRepresentation.OdbcTreatDecimalAsInt),
		handler.handleBooleanParameter(row.PreventUnloadToInternalStagesLevel, row.PreventUnloadToInternalStagesValue, &userRepresentation.PreventUnloadToInternalStages),
		handler.handleStringParameter(row.PythonProfilerModulesLevel, row.PythonProfilerModulesValue, &userRepresentation.PythonProfilerModules),
		handler.handleStringParameter(row.PythonProfilerTargetStageLevel, row.PythonProfilerTargetStageValue, &userRepresentation.PythonProfilerTargetStage),
		handler.handleStringParameter(row.QueryTagLevel, row.QueryTagValue, &userRepresentation.QueryTag),
		handler.handleBooleanParameter(row.QuotedIdentifiersIgnoreCaseLevel, row.QuotedIdentifiersIgnoreCaseValue, &userRepresentation.QuotedIdentifiersIgnoreCase),
		handler.handleIntegerParameter(row.RowsPerResultsetLevel, row.RowsPerResultsetValue, &userRepresentation.RowsPerResultset),
		handler.handleStringParameter(row.S3StageVpceDnsNameLevel, row.S3StageVpceDnsNameValue, &userRepresentation.S3StageVpceDnsName),
		handler.handleStringParameter(row.SearchPathLevel, row.SearchPathValue, &userRepresentation.SearchPath),
		handler.handleStringParameter(row.SimulatedDataSharingConsumerLevel, row.SimulatedDataSharingConsumerValue, &userRepresentation.SimulatedDataSharingConsumer),
		handler.handleIntegerParameter(row.StatementQueuedTimeoutInSecondsLevel, row.StatementQueuedTimeoutInSecondsValue, &userRepresentation.StatementQueuedTimeoutInSeconds),
		handler.handleIntegerParameter(row.StatementTimeoutInSecondsLevel, row.StatementTimeoutInSecondsValue, &userRepresentation.StatementTimeoutInSeconds),
		handler.handleBooleanParameter(row.StrictJsonOutputLevel, row.StrictJsonOutputValue, &userRepresentation.StrictJsonOutput),
		handler.handleStringParameter(row.TimeInputFormatLevel, row.TimeInputFormatValue, &userRepresentation.TimeInputFormat),
		handler.handleStringParameter(row.TimeOutputFormatLevel, row.TimeOutputFormatValue, &userRepresentation.TimeOutputFormat),
		handler.handleBooleanParameter(row.TimestampDayIsAlways24hLevel, row.TimestampDayIsAlways24hValue, &userRepresentation.TimestampDayIsAlways24h),
		handler.handleStringParameter(row.TimestampInputFormatLevel, row.TimestampInputFormatValue, &userRepresentation.TimestampInputFormat),
		handler.handleStringParameter(row.TimestampLTZOutputFormatLevel, row.TimestampLTZOutputFormatValue, &userRepresentation.TimestampLTZOutputFormat),
		handler.handleStringParameter(row.TimestampNTZOutputFormatLevel, row.TimestampNTZOutputFormatValue, &userRepresentation.TimestampNTZOutputFormat),
		handler.handleStringParameter(row.TimestampOutputFormatLevel, row.TimestampOutputFormatValue, &userRepresentation.TimestampOutputFormat),
		handler.handleStringParameter(row.TimestampTZOutputFormatLevel, row.TimestampTZOutputFormatValue, &userRepresentation.TimestampTZOutputFormat),
		handler.handleStringParameter(row.TimestampTypeMappingLevel, row.TimestampTypeMappingValue, &userRepresentation.TimestampTypeMapping),
		handler.handleStringParameter(row.TimezoneLevel, row.TimezoneValue, &userRepresentation.Timezone),
		handler.handleStringParameter(row.TraceLevelLevel, row.TraceLevelValue, &userRepresentation.TraceLevel),
		handler.handleBooleanParameter(row.TransactionAbortOnErrorLevel, row.TransactionAbortOnErrorValue, &userRepresentation.TransactionAbortOnError),
		handler.handleStringParameter(row.TransactionDefaultIsolationLevelLevel, row.TransactionDefaultIsolationLevelValue, &userRepresentation.TransactionDefaultIsolationLevel),
		handler.handleIntegerParameter(row.TwoDigitCenturyStartLevel, row.TwoDigitCenturyStartValue, &userRepresentation.TwoDigitCenturyStart),
		handler.handleStringParameter(row.UnsupportedDDLActionLevel, row.UnsupportedDDLActionValue, &userRepresentation.UnsupportedDDLAction),
		handler.handleBooleanParameter(row.UseCachedResultLevel, row.UseCachedResultValue, &userRepresentation.UseCachedResult),
		handler.handleIntegerParameter(row.WeekOfYearPolicyLevel, row.WeekOfYearPolicyValue, &userRepresentation.WeekOfYearPolicy),
		handler.handleIntegerParameter(row.WeekStartLevel, row.WeekStartValue, &userRepresentation.WeekStart),
	)
	if errs != nil {
		return nil, errs
	}

	return userRepresentation, nil
}
