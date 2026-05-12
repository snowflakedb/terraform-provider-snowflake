package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleUserValidations(t *testing.T) {
	testCases := []struct {
		name           string
		inputRows      [][]string
		outputContains []string
	}{
		{
			name: "multiple users of different types",
			inputRows: [][]string{
				{"comment", "name", "type", "login_name"},
				{"person comment", "TEST_USER_PERSON", "PERSON", "test_person_login"},
				{"service comment", "TEST_USER_SERVICE", "SERVICE", "test_service_login"},
				{"legacy comment", "TEST_USER_LEGACY", "LEGACY_SERVICE", "test_legacy_login"},
			},
			outputContains: []string{
				`resource "snowflake_user" "snowflake_generated_user_TEST_USER_PERSON"`,
				`resource "snowflake_service_user" "snowflake_generated_service_user_TEST_USER_SERVICE"`,
				`resource "snowflake_legacy_service_user" "snowflake_generated_legacy_service_user_TEST_USER_LEGACY"`,
			},
		},
		{
			name: "empty type defaults to PERSON",
			inputRows: [][]string{
				{"comment", "name", "type"},
				{"", "TEST_USER", ""},
			},
			outputContains: []string{
				`resource "snowflake_user" "snowflake_generated_user_TEST_USER"`,
			},
		},
		{
			name: "lowercase type is handled correctly",
			inputRows: [][]string{
				{"comment", "name", "type"},
				{"", "TEST_USER_1", "person"},
				{"", "TEST_USER_2", "service"},
				{"", "TEST_USER_3", "legacy_service"},
			},
			outputContains: []string{
				`resource "snowflake_user" "snowflake_generated_user_TEST_USER_1"`,
				`resource "snowflake_service_user" "snowflake_generated_service_user_TEST_USER_2"`,
				`resource "snowflake_legacy_service_user" "snowflake_generated_legacy_service_user_TEST_USER_3"`,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := HandleUsers(&Config{
				ObjectType: ObjectTypeUsers,
				ImportFlag: ImportStatementTypeBlock,
			}, tc.inputRows)

			assert.NoError(t, err)
			for _, expectedOutput := range tc.outputContains {
				assert.Contains(t, output, expectedOutput)
			}
		})
	}
}

func TestHandleUserMappings(t *testing.T) {
	testCases := []struct {
		name           string
		inputRows      [][]string
		expectedOutput string
	}{
		{
			name: "basic user (PERSON)",
			inputRows: [][]string{
				{"name", "type"},
				{"TEST_USER", "PERSON"},
			},
			expectedOutput: `
resource "snowflake_user" "snowflake_generated_user_TEST_USER" {
  name = "TEST_USER"
}
import {
  to = snowflake_user.snowflake_generated_user_TEST_USER
  id = "\"TEST_USER\""
}
`,
		},
		{
			name: "complete user (PERSON) with all basic attributes",
			inputRows: [][]string{
				{"comment", "default_namespace", "default_role", "default_secondary_roles", "default_warehouse", "disabled", "display_name", "email", "first_name", "last_name", "login_name", "must_change_password", "name", "type"},
				{"User comment", "DB.SCHEMA", "ANALYST", "ALL", "COMPUTE_WH", "true", "Test User Display", "test@example.com", "John", "Doe", "test_login", "true", "TEST_USER", "PERSON"},
			},
			expectedOutput: `
resource "snowflake_user" "snowflake_generated_user_TEST_USER" {
  name = "TEST_USER"
  comment = "User comment"
  default_namespace = "DB.SCHEMA"
  default_role = "ANALYST"
  default_secondary_roles_option = "ALL"
  default_warehouse = "COMPUTE_WH"
  disabled = "true"
  display_name = "Test User Display"
  email = "test@example.com"
  first_name = "John"
  last_name = "Doe"
  login_name = "test_login"
  must_change_password = "true"
}
import {
  to = snowflake_user.snowflake_generated_user_TEST_USER
  id = "\"TEST_USER\""
}
`,
		},
		{
			name: "user with parameters",
			inputRows: [][]string{
				{"name", "type", "abort_detached_query_level", "abort_detached_query_value", "timezone_level", "timezone_value"},
				{"TEST_USER", "PERSON", "USER", "true", "USER", "America/New_York"},
			},
			expectedOutput: `
resource "snowflake_user" "snowflake_generated_user_TEST_USER" {
  name = "TEST_USER"
  abort_detached_query = true
  timezone = "America/New_York"
}
import {
  to = snowflake_user.snowflake_generated_user_TEST_USER
  id = "\"TEST_USER\""
}
`,
		},
		{
			name: "user with describe output fields (middle_name, rsa_public_key, etc.)",
			inputRows: [][]string{
				{"name", "type", "login_name", "display_name", "first_name", "middle_name", "last_name", "email", "comment", "custom_landing_page_url", "custom_landing_page_url_flush_next_ui_load", "mins_to_bypass_network_policy", "snowflake_support", "password", "password_last_set_time", "rsa_public_key", "rsa_public_key_fp", "rsa_public_key2", "rsa_public_key2_fp"},
				{"TEST_USER_DESCRIBE", "PERSON", "test_login", "Test Display", "John", "Michael", "Doe", "john@example.com", "Describe user", "https://landing.page", "true", "30", "false", "********", "2025-12-03 10:00:00", "RSA_KEY_1", "KEY_FP_1", "RSA_KEY_2", "KEY_FP_2"},
			},
			expectedOutput: `
resource "snowflake_user" "snowflake_generated_user_TEST_USER_DESCRIBE" {
  name = "TEST_USER_DESCRIBE"
  comment = "Describe user"
  display_name = "Test Display"
  email = "john@example.com"
  first_name = "John"
  last_name = "Doe"
  login_name = "test_login"
  middle_name = "Michael"
  rsa_public_key = <<EOT
RSA_KEY_1
EOT
  rsa_public_key_2 = <<EOT
RSA_KEY_2
EOT
}
import {
  to = snowflake_user.snowflake_generated_user_TEST_USER_DESCRIBE
  id = "\"TEST_USER_DESCRIBE\""
}
`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := HandleUsers(&Config{
				ObjectType: ObjectTypeUsers,
				ImportFlag: ImportStatementTypeBlock,
			}, tc.inputRows)

			assert.NoError(t, err)
			assert.Equal(t, strings.TrimLeft(tc.expectedOutput, "\n"), strings.TrimLeft(output, "\n"))
		})
	}
}

func TestHandleServiceUserMappings(t *testing.T) {
	testCases := []struct {
		name           string
		inputRows      [][]string
		expectedOutput string
	}{
		{
			name: "basic service user",
			inputRows: [][]string{
				{"name", "type"},
				{"SVC_USER", "SERVICE"},
			},
			expectedOutput: `
resource "snowflake_service_user" "snowflake_generated_service_user_SVC_USER" {
  name = "SVC_USER"
}
import {
  to = snowflake_service_user.snowflake_generated_service_user_SVC_USER
  id = "\"SVC_USER\""
}
`,
		},
		{
			name: "complete service user with all basic attributes",
			inputRows: [][]string{
				{"comment", "default_namespace", "default_role", "default_secondary_roles", "default_warehouse", "disabled", "display_name", "email", "login_name", "name", "type"},
				{"Service user comment", "DB.SCHEMA", "SVC_ROLE", "ALL", "SVC_WH", "false", "Service User", "svc@example.com", "svc_login", "SVC_USER", "SERVICE"},
			},
			expectedOutput: `
resource "snowflake_service_user" "snowflake_generated_service_user_SVC_USER" {
  name = "SVC_USER"
  comment = "Service user comment"
  default_namespace = "DB.SCHEMA"
  default_role = "SVC_ROLE"
  default_secondary_roles_option = "ALL"
  default_warehouse = "SVC_WH"
  display_name = "Service User"
  email = "svc@example.com"
  login_name = "svc_login"
}
import {
  to = snowflake_service_user.snowflake_generated_service_user_SVC_USER
  id = "\"SVC_USER\""
}
`,
		},
		{
			name: "service user with parameters",
			inputRows: [][]string{
				{"name", "type", "autocommit_level", "autocommit_value", "log_level_level", "log_level_value"},
				{"SVC_USER", "SERVICE", "USER", "false", "USER", "INFO"},
			},
			expectedOutput: `
resource "snowflake_service_user" "snowflake_generated_service_user_SVC_USER" {
  name = "SVC_USER"
  autocommit = false
  log_level = "INFO"
}
import {
  to = snowflake_service_user.snowflake_generated_service_user_SVC_USER
  id = "\"SVC_USER\""
}
`,
		},
		{
			name: "service user with describe fields",
			inputRows: [][]string{
				{"name", "type", "login_name", "display_name", "email", "comment", "custom_landing_page_url", "rsa_public_key", "rsa_public_key_fp", "rsa_public_key2", "rsa_public_key2_fp"},
				{"SVC_DESCRIBE", "SERVICE", "svc_login", "Service Display", "svc@example.com", "Service with describe", "https://svc.landing.page", "RSA_SVC_KEY", "SVC_KEY_FP", "RSA_SVC_KEY2", "SVC_KEY_FP2"},
			},
			expectedOutput: `
resource "snowflake_service_user" "snowflake_generated_service_user_SVC_DESCRIBE" {
  name = "SVC_DESCRIBE"
  comment = "Service with describe"
  display_name = "Service Display"
  email = "svc@example.com"
  login_name = "svc_login"
  rsa_public_key = <<EOT
RSA_SVC_KEY
EOT
  rsa_public_key_2 = <<EOT
RSA_SVC_KEY2
EOT
}
import {
  to = snowflake_service_user.snowflake_generated_service_user_SVC_DESCRIBE
  id = "\"SVC_DESCRIBE\""
}
`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := HandleUsers(&Config{
				ObjectType: ObjectTypeUsers,
				ImportFlag: ImportStatementTypeBlock,
			}, tc.inputRows)

			assert.NoError(t, err)
			assert.Equal(t, strings.TrimLeft(tc.expectedOutput, "\n"), strings.TrimLeft(output, "\n"))
		})
	}
}

func TestHandleLegacyServiceUserMappings(t *testing.T) {
	testCases := []struct {
		name           string
		inputRows      [][]string
		expectedOutput string
	}{
		{
			name: "basic legacy service user",
			inputRows: [][]string{
				{"name", "type"},
				{"LEGACY_SVC_USER", "LEGACY_SERVICE"},
			},
			expectedOutput: `
resource "snowflake_legacy_service_user" "snowflake_generated_legacy_service_user_LEGACY_SVC_USER" {
  name = "LEGACY_SVC_USER"
}
import {
  to = snowflake_legacy_service_user.snowflake_generated_legacy_service_user_LEGACY_SVC_USER
  id = "\"LEGACY_SVC_USER\""
}
`,
		},
		{
			name: "complete legacy service user with all basic attributes",
			inputRows: [][]string{
				{"comment", "default_namespace", "default_role", "default_secondary_roles", "default_warehouse", "disabled", "display_name", "email", "login_name", "must_change_password", "name", "type"},
				{"Legacy service comment", "DB.SCHEMA", "LEGACY_ROLE", "ALL", "LEGACY_WH", "true", "Legacy Service User", "legacy@example.com", "legacy_login", "true", "LEGACY_SVC_USER", "LEGACY_SERVICE"},
			},
			expectedOutput: `
resource "snowflake_legacy_service_user" "snowflake_generated_legacy_service_user_LEGACY_SVC_USER" {
  name = "LEGACY_SVC_USER"
  comment = "Legacy service comment"
  default_namespace = "DB.SCHEMA"
  default_role = "LEGACY_ROLE"
  default_secondary_roles_option = "ALL"
  default_warehouse = "LEGACY_WH"
  disabled = "true"
  display_name = "Legacy Service User"
  email = "legacy@example.com"
  login_name = "legacy_login"
  must_change_password = "true"
}
import {
  to = snowflake_legacy_service_user.snowflake_generated_legacy_service_user_LEGACY_SVC_USER
  id = "\"LEGACY_SVC_USER\""
}
`,
		},
		{
			name: "legacy service user with parameters",
			inputRows: [][]string{
				{"name", "type", "statement_timeout_in_seconds_level", "statement_timeout_in_seconds_value", "trace_level_level", "trace_level_value"},
				{"LEGACY_SVC_USER", "LEGACY_SERVICE", "USER", "3600", "USER", "ON_EVENT"},
			},
			expectedOutput: `
resource "snowflake_legacy_service_user" "snowflake_generated_legacy_service_user_LEGACY_SVC_USER" {
  name = "LEGACY_SVC_USER"
  statement_timeout_in_seconds = 3600
  trace_level = "ON_EVENT"
}
import {
  to = snowflake_legacy_service_user.snowflake_generated_legacy_service_user_LEGACY_SVC_USER
  id = "\"LEGACY_SVC_USER\""
}
`,
		},
		{
			name: "legacy service user with describe fields",
			inputRows: [][]string{
				{"name", "type", "login_name", "display_name", "email", "comment", "rsa_public_key", "rsa_public_key_fp", "rsa_public_key2", "rsa_public_key2_fp"},
				{"LEGACY_SVC_DESCRIBE", "LEGACY_SERVICE", "legacy_login", "Legacy Display", "legacy@example.com", "Legacy with describe", "RSA_LEGACY_KEY", "LEGACY_KEY_FP", "RSA_LEGACY_KEY2", "LEGACY_KEY_FP2"},
			},
			expectedOutput: `
resource "snowflake_legacy_service_user" "snowflake_generated_legacy_service_user_LEGACY_SVC_DESCRIBE" {
  name = "LEGACY_SVC_DESCRIBE"
  comment = "Legacy with describe"
  display_name = "Legacy Display"
  email = "legacy@example.com"
  login_name = "legacy_login"
  rsa_public_key = <<EOT
RSA_LEGACY_KEY
EOT
  rsa_public_key_2 = <<EOT
RSA_LEGACY_KEY2
EOT
}
import {
  to = snowflake_legacy_service_user.snowflake_generated_legacy_service_user_LEGACY_SVC_DESCRIBE
  id = "\"LEGACY_SVC_DESCRIBE\""
}
`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := HandleUsers(&Config{
				ObjectType: ObjectTypeUsers,
				ImportFlag: ImportStatementTypeBlock,
			}, tc.inputRows)

			assert.NoError(t, err)
			assert.Equal(t, strings.TrimLeft(tc.expectedOutput, "\n"), strings.TrimLeft(output, "\n"))
		})
	}
}

func TestHandleUserImportStatementFormat(t *testing.T) {
	testCases := []struct {
		name           string
		inputRows      [][]string
		importFlag     ImportStatementType
		expectedOutput string
	}{
		{
			name: "user with statement import format",
			inputRows: [][]string{
				{"name", "type"},
				{"TEST_USER", "PERSON"},
			},
			importFlag: ImportStatementTypeStatement,
			expectedOutput: `
resource "snowflake_user" "snowflake_generated_user_TEST_USER" {
  name = "TEST_USER"
}
# terraform import snowflake_user.snowflake_generated_user_TEST_USER '"TEST_USER"'
`,
		},
		{
			name: "service user with statement import format",
			inputRows: [][]string{
				{"name", "type"},
				{"SVC_USER", "SERVICE"},
			},
			importFlag: ImportStatementTypeStatement,
			expectedOutput: `
resource "snowflake_service_user" "snowflake_generated_service_user_SVC_USER" {
  name = "SVC_USER"
}
# terraform import snowflake_service_user.snowflake_generated_service_user_SVC_USER '"SVC_USER"'
`,
		},
		{
			name: "legacy service user with statement import format",
			inputRows: [][]string{
				{"name", "type"},
				{"LEGACY_USER", "LEGACY_SERVICE"},
			},
			importFlag: ImportStatementTypeStatement,
			expectedOutput: `
resource "snowflake_legacy_service_user" "snowflake_generated_legacy_service_user_LEGACY_USER" {
  name = "LEGACY_USER"
}
# terraform import snowflake_legacy_service_user.snowflake_generated_legacy_service_user_LEGACY_USER '"LEGACY_USER"'
`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := HandleUsers(&Config{
				ObjectType: ObjectTypeUsers,
				ImportFlag: tc.importFlag,
			}, tc.inputRows)

			assert.NoError(t, err)
			assert.Equal(t, strings.TrimLeft(tc.expectedOutput, "\n"), strings.TrimLeft(output, "\n"))
		})
	}
}

func TestHandleUserWithAllParameters(t *testing.T) {
	// Test user with all supported parameters set at USER level
	inputRows := [][]string{
		{
			"name", "type",
			"abort_detached_query_level", "abort_detached_query_value",
			"autocommit_level", "autocommit_value",
			"binary_input_format_level", "binary_input_format_value",
			"binary_output_format_level", "binary_output_format_value",
			"client_memory_limit_level", "client_memory_limit_value",
			"client_metadata_request_use_connection_ctx_level", "client_metadata_request_use_connection_ctx_value",
			"client_prefetch_threads_level", "client_prefetch_threads_value",
			"client_result_chunk_size_level", "client_result_chunk_size_value",
			"client_result_column_case_insensitive_level", "client_result_column_case_insensitive_value",
			"client_session_keep_alive_level", "client_session_keep_alive_value",
			"client_session_keep_alive_heartbeat_frequency_level", "client_session_keep_alive_heartbeat_frequency_value",
			"client_timestamp_type_mapping_level", "client_timestamp_type_mapping_value",
			"date_input_format_level", "date_input_format_value",
			"date_output_format_level", "date_output_format_value",
			"enable_unload_physical_type_optimization_level", "enable_unload_physical_type_optimization_value",
			"enable_unredacted_query_syntax_error_level", "enable_unredacted_query_syntax_error_value",
			"error_on_nondeterministic_merge_level", "error_on_nondeterministic_merge_value",
			"error_on_nondeterministic_update_level", "error_on_nondeterministic_update_value",
			"geography_output_format_level", "geography_output_format_value",
			"geometry_output_format_level", "geometry_output_format_value",
			"jdbc_treat_decimal_as_int_level", "jdbc_treat_decimal_as_int_value",
			"jdbc_treat_timestamp_ntz_as_utc_level", "jdbc_treat_timestamp_ntz_as_utc_value",
			"jdbc_use_session_timezone_level", "jdbc_use_session_timezone_value",
			"json_indent_level", "json_indent_value",
			"lock_timeout_level", "lock_timeout_value",
			"log_level_level", "log_level_value",
			"multi_statement_count_level", "multi_statement_count_value",
			"network_policy_level", "network_policy_value",
			"noorder_sequence_as_default_level", "noorder_sequence_as_default_value",
			"odbc_treat_decimal_as_int_level", "odbc_treat_decimal_as_int_value",
			"prevent_unload_to_internal_stages_level", "prevent_unload_to_internal_stages_value",
			"query_tag_level", "query_tag_value",
			"quoted_identifiers_ignore_case_level", "quoted_identifiers_ignore_case_value",
			"rows_per_resultset_level", "rows_per_resultset_value",
			"s3_stage_vpce_dns_name_level", "s3_stage_vpce_dns_name_value",
			"search_path_level", "search_path_value",
			"simulated_data_sharing_consumer_level", "simulated_data_sharing_consumer_value",
			"statement_queued_timeout_in_seconds_level", "statement_queued_timeout_in_seconds_value",
			"statement_timeout_in_seconds_level", "statement_timeout_in_seconds_value",
			"strict_json_output_level", "strict_json_output_value",
			"time_input_format_level", "time_input_format_value",
			"time_output_format_level", "time_output_format_value",
			"timestamp_day_is_always24h_level", "timestamp_day_is_always24h_value",
			"timestamp_input_format_level", "timestamp_input_format_value",
			"timestamp_l_t_z_output_format_level", "timestamp_l_t_z_output_format_value",
			"timestamp_n_t_z_output_format_level", "timestamp_n_t_z_output_format_value",
			"timestamp_output_format_level", "timestamp_output_format_value",
			"timestamp_t_z_output_format_level", "timestamp_t_z_output_format_value",
			"timestamp_type_mapping_level", "timestamp_type_mapping_value",
			"timezone_level", "timezone_value",
			"trace_level_level", "trace_level_value",
			"transaction_abort_on_error_level", "transaction_abort_on_error_value",
			"transaction_default_isolation_level_level", "transaction_default_isolation_level_value",
			"two_digit_century_start_level", "two_digit_century_start_value",
			"unsupported_d_d_l_action_level", "unsupported_d_d_l_action_value",
			"use_cached_result_level", "use_cached_result_value",
			"week_of_year_policy_level", "week_of_year_policy_value",
			"week_start_level", "week_start_value",
		},
		{
			"ALL_PARAMS_USER", "PERSON",
			"USER", "true", // abort_detached_query
			"USER", "false", // autocommit
			"USER", "HEX", // binary_input_format
			"USER", "BASE64", // binary_output_format
			"USER", "2048", // client_memory_limit
			"USER", "true", // client_metadata_request_use_connection_ctx
			"USER", "8", // client_prefetch_threads
			"USER", "128", // client_result_chunk_size
			"USER", "true", // client_result_column_case_insensitive
			"USER", "true", // client_session_keep_alive
			"USER", "1800", // client_session_keep_alive_heartbeat_frequency
			"USER", "TIMESTAMP_NTZ", // client_timestamp_type_mapping
			"USER", "YYYY-MM-DD", // date_input_format
			"USER", "YYYY-MM-DD", // date_output_format
			"USER", "true", // enable_unload_physical_type_optimization
			"USER", "true", // enable_unredacted_query_syntax_error
			"USER", "false", // error_on_nondeterministic_merge
			"USER", "false", // error_on_nondeterministic_update
			"USER", "GeoJSON", // geography_output_format
			"USER", "WKT", // geometry_output_format
			"USER", "false", // jdbc_treat_decimal_as_int
			"USER", "true", // jdbc_treat_timestamp_ntz_as_utc
			"USER", "true", // jdbc_use_session_timezone
			"USER", "4", // json_indent
			"USER", "43200", // lock_timeout
			"USER", "INFO", // log_level
			"USER", "0", // multi_statement_count
			"USER", "NETWORK_POLICY_1", // network_policy
			"USER", "true", // noorder_sequence_as_default
			"USER", "false", // odbc_treat_decimal_as_int
			"USER", "true", // prevent_unload_to_internal_stages
			"USER", "my_query_tag", // query_tag
			"USER", "true", // quoted_identifiers_ignore_case
			"USER", "100000", // rows_per_resultset
			"USER", "vpce-dns-name", // s3_stage_vpce_dns_name
			"USER", "$current, $public", // search_path
			"USER", "CONSUMER_ACCOUNT", // simulated_data_sharing_consumer
			"USER", "300", // statement_queued_timeout_in_seconds
			"USER", "86400", // statement_timeout_in_seconds
			"USER", "true", // strict_json_output
			"USER", "HH24:MI:SS", // time_input_format
			"USER", "HH24:MI:SS", // time_output_format
			"USER", "true", // timestamp_day_is_always_24h
			"USER", "YYYY-MM-DD HH24:MI:SS", // timestamp_input_format
			"USER", "YYYY-MM-DD HH24:MI:SS TZHTZM", // timestamp_ltz_output_format
			"USER", "YYYY-MM-DD HH24:MI:SS", // timestamp_ntz_output_format
			"USER", "YYYY-MM-DD HH24:MI:SS.FF3", // timestamp_output_format
			"USER", "YYYY-MM-DD HH24:MI:SS TZHTZM", // timestamp_tz_output_format
			"USER", "TIMESTAMP_NTZ", // timestamp_type_mapping
			"USER", "America/New_York", // timezone
			"USER", "ON_EVENT", // trace_level
			"USER", "true", // transaction_abort_on_error
			"USER", "READ COMMITTED", // transaction_default_isolation_level
			"USER", "1970", // two_digit_century_start
			"USER", "IGNORE", // unsupported_ddl_action
			"USER", "false", // use_cached_result
			"USER", "1", // week_of_year_policy
			"USER", "1", // week_start
		},
	}

	expectedOutput := `resource "snowflake_user" "snowflake_generated_user_ALL_PARAMS_USER" {
  name = "ALL_PARAMS_USER"
  abort_detached_query = true
  autocommit = false
  binary_input_format = "HEX"
  binary_output_format = "BASE64"
  client_memory_limit = 2048
  client_metadata_request_use_connection_ctx = true
  client_prefetch_threads = 8
  client_result_chunk_size = 128
  client_result_column_case_insensitive = true
  client_session_keep_alive = true
  client_session_keep_alive_heartbeat_frequency = 1800
  client_timestamp_type_mapping = "TIMESTAMP_NTZ"
  date_input_format = "YYYY-MM-DD"
  date_output_format = "YYYY-MM-DD"
  enable_unload_physical_type_optimization = true
  enable_unredacted_query_syntax_error = true
  error_on_nondeterministic_merge = false
  error_on_nondeterministic_update = false
  geography_output_format = "GeoJSON"
  geometry_output_format = "WKT"
  jdbc_treat_decimal_as_int = false
  jdbc_treat_timestamp_ntz_as_utc = true
  jdbc_use_session_timezone = true
  json_indent = 4
  lock_timeout = 43200
  log_level = "INFO"
  multi_statement_count = 0
  network_policy = "NETWORK_POLICY_1"
  noorder_sequence_as_default = true
  odbc_treat_decimal_as_int = false
  prevent_unload_to_internal_stages = true
  query_tag = "my_query_tag"
  quoted_identifiers_ignore_case = true
  rows_per_resultset = 100000
  s3_stage_vpce_dns_name = "vpce-dns-name"
  search_path = "$current, $public"
  simulated_data_sharing_consumer = "CONSUMER_ACCOUNT"
  statement_queued_timeout_in_seconds = 300
  statement_timeout_in_seconds = 86400
  strict_json_output = true
  time_input_format = "HH24:MI:SS"
  time_output_format = "HH24:MI:SS"
  timestamp_day_is_always_24h = true
  timestamp_input_format = "YYYY-MM-DD HH24:MI:SS"
  timestamp_ltz_output_format = "YYYY-MM-DD HH24:MI:SS TZHTZM"
  timestamp_ntz_output_format = "YYYY-MM-DD HH24:MI:SS"
  timestamp_output_format = "YYYY-MM-DD HH24:MI:SS.FF3"
  timestamp_type_mapping = "TIMESTAMP_NTZ"
  timestamp_tz_output_format = "YYYY-MM-DD HH24:MI:SS TZHTZM"
  timezone = "America/New_York"
  trace_level = "ON_EVENT"
  transaction_abort_on_error = true
  transaction_default_isolation_level = "READ COMMITTED"
  two_digit_century_start = 1970
  unsupported_ddl_action = "IGNORE"
  use_cached_result = false
  week_of_year_policy = 1
  week_start = 1
}
import {
  to = snowflake_user.snowflake_generated_user_ALL_PARAMS_USER
  id = "\"ALL_PARAMS_USER\""
}
`

	output, err := HandleUsers(&Config{
		ObjectType: ObjectTypeUsers,
		ImportFlag: ImportStatementTypeBlock,
	}, inputRows)

	assert.NoError(t, err)
	assert.Equal(t, expectedOutput, output)
}

func TestHandleUserPasswordParameter(t *testing.T) {
	// Password from DESCRIBE is always masked as "********" and should NOT be included in output
	// This test verifies that password field is correctly ignored during mapping
	testCases := []struct {
		name             string
		inputRows        [][]string
		shouldContain    []string
		shouldNotContain []string
	}{
		{
			name: "person user with masked password - password should not be in output",
			inputRows: [][]string{
				{"name", "type", "password", "login_name", "comment"},
				{"TEST_USER", "PERSON", "********", "test_login", "Test user comment"},
			},
			shouldContain: []string{
				`name = "TEST_USER"`,
				`login_name = "test_login"`,
				`comment = "Test user comment"`,
			},
			shouldNotContain: []string{
				"password =",
				"********",
			},
		},
		{
			name: "legacy service user with masked password - password should not be in output",
			inputRows: [][]string{
				{"name", "type", "password", "login_name"},
				{"LEGACY_SVC", "LEGACY_SERVICE", "********", "legacy_login"},
			},
			shouldContain: []string{
				`name = "LEGACY_SVC"`,
				`login_name = "legacy_login"`,
			},
			shouldNotContain: []string{
				"password =",
				"********",
			},
		},
		{
			name: "user with empty password - no password field in output",
			inputRows: [][]string{
				{"name", "type", "password", "login_name"},
				{"TEST_USER", "PERSON", "", "test_login"},
			},
			shouldContain: []string{
				`name = "TEST_USER"`,
				`login_name = "test_login"`,
			},
			shouldNotContain: []string{
				"password =",
			},
		},
		{
			name: "service user - cannot have password field",
			inputRows: [][]string{
				{"name", "type", "password", "login_name"},
				{"SVC_USER", "SERVICE", "********", "svc_login"},
			},
			shouldContain: []string{
				`name = "SVC_USER"`,
				`login_name = "svc_login"`,
			},
			shouldNotContain: []string{
				"password =",
				"********",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := HandleUsers(&Config{
				ObjectType: ObjectTypeUsers,
				ImportFlag: ImportStatementTypeBlock,
			}, tc.inputRows)

			assert.NoError(t, err)
			for _, expected := range tc.shouldContain {
				assert.Contains(t, output, expected, "Expected output to contain: %s", expected)
			}
			for _, notExpected := range tc.shouldNotContain {
				assert.NotContains(t, output, notExpected, "Expected output NOT to contain: %s", notExpected)
			}
		})
	}
}
