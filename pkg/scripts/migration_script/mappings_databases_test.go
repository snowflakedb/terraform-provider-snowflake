package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleDatabaseMappings(t *testing.T) {
	testCases := []struct {
		name           string
		inputRows      [][]string
		expectedOutput string
	}{
		{
			name: "database without parameters",
			inputRows: [][]string{
				{"comment", "created_on", "dropped_on", "is_current", "is_default", "kind", "name", "options", "origin", "owner", "owner_role_type", "resource_group", "retention_time"},
				{"", "2024-06-06 00:00:00.000 +0000 UTC", "", "false", "false", "STANDARD", "DB1", "", "", "ADMIN", "ROLE", "", "1"},
			},
			expectedOutput: `
resource "snowflake_database" "snowflake_generated_database_DB1" {
  name = "DB1"
}
import {
  to = snowflake_database.snowflake_generated_database_DB1
  id = "\"DB1\""
}
`,
		},
		{
			name: "minimal database",
			inputRows: [][]string{
				{"comment", "created_on", "dropped_on", "is_current", "is_default", "kind", "name", "options", "origin", "owner", "owner_role_type", "resource_group", "retention_time", "catalog_level", "catalog_value", "data_retention_time_in_days_level", "data_retention_time_in_days_value", "default_ddl_collation_level", "default_ddl_collation_value", "enable_console_output_level", "enable_console_output_value", "external_volume_level", "external_volume_value", "log_level_level", "log_level_value", "max_data_extension_time_in_days_level", "max_data_extension_time_in_days_value", "quoted_identifiers_ignore_case_level", "quoted_identifiers_ignore_case_value", "replace_invalid_characters_level", "replace_invalid_characters_value", "storage_serialization_policy_level", "storage_serialization_policy_value", "suspend_task_after_num_failures_level", "suspend_task_after_num_failures_value", "task_auto_retry_attempts_level", "task_auto_retry_attempts_value", "trace_level_level", "trace_level_value", "user_task_managed_initial_warehouse_size_level", "user_task_managed_initial_warehouse_size_value", "user_task_minimum_trigger_interval_in_seconds_level", "user_task_minimum_trigger_interval_in_seconds_value", "user_task_timeout_ms_level", "user_task_timeout_ms_value"},
				{"", "2024-06-06 00:00:00.000 +0000 UTC", "", "false", "false", "STANDARD", "DB1", "", "", "ADMIN", "ROLE", "", "1", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
			},
			expectedOutput: `
resource "snowflake_database" "snowflake_generated_database_DB1" {
  name = "DB1"
}
import {
  to = snowflake_database.snowflake_generated_database_DB1
  id = "\"DB1\""
}
`,
		},
		{
			name: "database with all fields",
			inputRows: [][]string{
				{"comment", "created_on", "dropped_on", "is_current", "is_default", "kind", "name", "options", "origin", "owner", "owner_role_type", "resource_group", "retention_time", "catalog_level", "catalog_value", "data_retention_time_in_days_level", "data_retention_time_in_days_value", "default_ddl_collation_level", "default_ddl_collation_value", "enable_console_output_level", "enable_console_output_value", "external_volume_level", "external_volume_value", "log_level_level", "log_level_value", "max_data_extension_time_in_days_level", "max_data_extension_time_in_days_value", "quoted_identifiers_ignore_case_level", "quoted_identifiers_ignore_case_value", "replace_invalid_characters_level", "replace_invalid_characters_value", "storage_serialization_policy_level", "storage_serialization_policy_value", "suspend_task_after_num_failures_level", "suspend_task_after_num_failures_value", "task_auto_retry_attempts_level", "task_auto_retry_attempts_value", "trace_level_level", "trace_level_value", "user_task_managed_initial_warehouse_size_level", "user_task_managed_initial_warehouse_size_value", "user_task_minimum_trigger_interval_in_seconds_level", "user_task_minimum_trigger_interval_in_seconds_value", "user_task_timeout_ms_level", "user_task_timeout_ms_value"},
				{"comment here", "2024-06-06 00:00:00.000 +0000 UTC", "", "true", "true", "STANDARD", "DB2", "TRANSIENT", "", "ADMIN", "ROLE", "", "10", "DATABASE", "CAT", "DATABASE", "2", "DATABASE", "utf8", "DATABASE", "true", "DATABASE", "S3VOLUME", "DATABASE", "DEBUG", "DATABASE", "20", "DATABASE", "false", "DATABASE", "true", "DATABASE", "AVRO", "DATABASE", "7", "DATABASE", "4", "DATABASE", "PROD", "DATABASE", "SMALL", "DATABASE", "5", "DATABASE", "1000", "DATABASE"},
			},
			expectedOutput: `
resource "snowflake_database" "snowflake_generated_database_DB2" {
  name = "DB2"
  catalog = "CAT"
  comment = "comment here"
  data_retention_time_in_days = 2
  default_ddl_collation = "utf8"
  enable_console_output = true
  external_volume = "S3VOLUME"
  is_transient = true
  log_level = "DEBUG"
  max_data_extension_time_in_days = 20
  quoted_identifiers_ignore_case = false
  replace_invalid_characters = true
  storage_serialization_policy = "AVRO"
  suspend_task_after_num_failures = 7
  task_auto_retry_attempts = 4
  trace_level = "PROD"
  user_task_managed_initial_warehouse_size = "SMALL"
  user_task_minimum_trigger_interval_in_seconds = 5
  user_task_timeout_ms = 1000
}
import {
  to = snowflake_database.snowflake_generated_database_DB2
  id = "\"DB2\""
}
`,
		},
		{
			name: "database with all parameters set on higher level",
			inputRows: [][]string{
				{"comment", "created_on", "dropped_on", "is_current", "is_default", "kind", "name", "options", "origin", "owner", "owner_role_type", "resource_group", "retention_time", "catalog_level", "catalog_value", "data_retention_time_in_days_level", "data_retention_time_in_days_value", "default_ddl_collation_level", "default_ddl_collation_value", "enable_console_output_level", "enable_console_output_value", "external_volume_level", "external_volume_value", "log_level_level", "log_level_value", "max_data_extension_time_in_days_level", "max_data_extension_time_in_days_value", "quoted_identifiers_ignore_case_level", "quoted_identifiers_ignore_case_value", "replace_invalid_characters_level", "replace_invalid_characters_value", "storage_serialization_policy_level", "storage_serialization_policy_value", "suspend_task_after_num_failures_level", "suspend_task_after_num_failures_value", "task_auto_retry_attempts_level", "task_auto_retry_attempts_value", "trace_level_level", "trace_level_value", "user_task_managed_initial_warehouse_size_level", "user_task_managed_initial_warehouse_size_value", "user_task_minimum_trigger_interval_in_seconds_level", "user_task_minimum_trigger_interval_in_seconds_value", "user_task_timeout_ms_level", "user_task_timeout_ms_value"},
				{"comment here", "2024-06-06 00:00:00.000 +0000 UTC", "", "true", "true", "STANDARD", "DB2", "", "", "ADMIN", "ROLE", "", "10", "ACCOUNT", "CAT", "ACCOUNT", "2", "ACCOUNT", "utf8", "ACCOUNT", "true", "ACCOUNT", "S3VOLUME", "ACCOUNT", "DEBUG", "ACCOUNT", "20", "ACCOUNT", "false", "ACCOUNT", "true", "ACCOUNT", "AVRO", "ACCOUNT", "7", "ACCOUNT", "4", "ACCOUNT", "PROD", "ACCOUNT", "SMALL", "ACCOUNT", "5", "ACCOUNT", "1000", "ACCOUNT"},
			},
			expectedOutput: `
resource "snowflake_database" "snowflake_generated_database_DB2" {
  name = "DB2"
  comment = "comment here"
}
import {
  to = snowflake_database.snowflake_generated_database_DB2
  id = "\"DB2\""
}
`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := HandleDatabases(&Config{
				ObjectType: ObjectTypeDatabases,
				ImportFlag: ImportStatementTypeBlock,
			}, tc.inputRows)
			assert.NoError(t, err)
			assert.Equal(t, strings.TrimLeft(tc.expectedOutput, "\n"), strings.TrimLeft(output, "\n"))
		})
	}
}
