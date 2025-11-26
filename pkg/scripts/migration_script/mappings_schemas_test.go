package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleSchemaMappings(t *testing.T) {
	testCases := []struct {
		name           string
		inputRows      [][]string
		expectedOutput string
	}{
		{
			name: "schema without parameters",
			inputRows: [][]string{
				{"comment", "created_on", "database_name", "dropped_on", "is_current", "is_default", "name", "options", "owner", "owner_role_type", "retention_time"},
				{"", "2024-06-06 00:00:00.000 +0000 UTC", "DB", "", "false", "false", "SCHEMA1", "", "ADMIN", "ROLE", "1"},
			},
			expectedOutput: `
resource "snowflake_schema" "snowflake_generated_schema_DB_SCHEMA1" {
  database = "DB"
  name = "SCHEMA1"
}
import {
  to = snowflake_schema.snowflake_generated_schema_DB_SCHEMA1
  id = "\"DB\".\"SCHEMA1\""
}
`,
		},
		{
			name: "minimal schema",
			inputRows: [][]string{
				{"comment", "created_on", "database_name", "dropped_on", "is_current", "is_default", "name", "options", "owner", "owner_role_type", "retention_time", "catalog_value", "catalog_level", "data_retention_time_in_days_value", "data_retention_time_in_days_level", "default_ddl_collation_value", "default_ddl_collation_level", "enable_console_output_value", "enable_console_output_level", "external_volume_value", "external_volume_level", "log_level_value", "log_level_level", "max_data_extension_time_in_days_value", "max_data_extension_time_in_days_level", "pipe_execution_paused_value", "pipe_execution_paused_level", "quoted_identifiers_ignore_case_value", "quoted_identifiers_ignore_case_level", "replace_invalid_characters_value", "replace_invalid_characters_level", "storage_serialization_policy_value", "storage_serialization_policy_level", "suspend_task_after_num_failures_value", "suspend_task_after_num_failures_level", "task_auto_retry_attempts_value", "task_auto_retry_attempts_level", "trace_level_value", "trace_level_level", "user_task_managed_initial_warehouse_size_value", "user_task_managed_initial_warehouse_size_level", "user_task_minimum_trigger_interval_in_seconds_value", "user_task_minimum_trigger_interval_in_seconds_level", "user_task_timeout_ms_value", "user_task_timeout_ms_level"},
				{"", "2024-06-06 00:00:00.000 +0000 UTC", "DB", "", "false", "false", "SCHEMA1", "", "ADMIN", "ROLE", "1", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
			},
			expectedOutput: `
resource "snowflake_schema" "snowflake_generated_schema_DB_SCHEMA1" {
  database = "DB"
  name = "SCHEMA1"
}
import {
  to = snowflake_schema.snowflake_generated_schema_DB_SCHEMA1
  id = "\"DB\".\"SCHEMA1\""
}
`,
		},
		{
			name: "schema with all fields",
			inputRows: [][]string{
				{"comment", "created_on", "database_name", "dropped_on", "is_current", "is_default", "name", "options", "owner", "owner_role_type", "retention_time", "catalog_value", "catalog_level", "data_retention_time_in_days_value", "data_retention_time_in_days_level", "default_ddl_collation_value", "default_ddl_collation_level", "enable_console_output_value", "enable_console_output_level", "external_volume_value", "external_volume_level", "log_level_value", "log_level_level", "max_data_extension_time_in_days_value", "max_data_extension_time_in_days_level", "pipe_execution_paused_value", "pipe_execution_paused_level", "quoted_identifiers_ignore_case_value", "quoted_identifiers_ignore_case_level", "replace_invalid_characters_value", "replace_invalid_characters_level", "storage_serialization_policy_value", "storage_serialization_policy_level", "suspend_task_after_num_failures_value", "suspend_task_after_num_failures_level", "task_auto_retry_attempts_value", "task_auto_retry_attempts_level", "trace_level_value", "trace_level_level", "user_task_managed_initial_warehouse_size_value", "user_task_managed_initial_warehouse_size_level", "user_task_minimum_trigger_interval_in_seconds_value", "user_task_minimum_trigger_interval_in_seconds_level", "user_task_timeout_ms_value", "user_task_timeout_ms_level"},
				{"comment here", "2024-06-06 00:00:00.000 +0000 UTC", "DB", "", "true", "true", "SCHEMA2", "TRANSIENT,MANAGED ACCESS", "ADMIN", "ROLE", "10", "CAT", "SCHEMA", "2", "SCHEMA", "utf8", "SCHEMA", "true", "SCHEMA", "S3VOLUME", "SCHEMA", "DEBUG", "SCHEMA", "20", "SCHEMA", "true", "SCHEMA", "false", "SCHEMA", "true", "SCHEMA", "AVRO", "SCHEMA", "7", "SCHEMA", "4", "SCHEMA", "PROD", "SCHEMA", "SMALL", "SCHEMA", "5", "SCHEMA", "1000", "SCHEMA"},
			},
			expectedOutput: `
resource "snowflake_schema" "snowflake_generated_schema_DB_SCHEMA2" {
  database = "DB"
  name = "SCHEMA2"
  catalog = "CAT"
  comment = "comment here"
  data_retention_time_in_days = 2
  default_ddl_collation = "utf8"
  enable_console_output = true
  external_volume = "S3VOLUME"
  is_transient = "true"
  log_level = "DEBUG"
  max_data_extension_time_in_days = 20
  pipe_execution_paused = true
  quoted_identifiers_ignore_case = false
  replace_invalid_characters = true
  storage_serialization_policy = "AVRO"
  suspend_task_after_num_failures = 7
  task_auto_retry_attempts = 4
  trace_level = "PROD"
  user_task_managed_initial_warehouse_size = "SMALL"
  user_task_minimum_trigger_interval_in_seconds = 5
  user_task_timeout_ms = 1000
  with_managed_access = "true"
}
import {
  to = snowflake_schema.snowflake_generated_schema_DB_SCHEMA2
  id = "\"DB\".\"SCHEMA2\""
}
`,
		},
		{
			name: "schema with all parameters set on higher level",
			inputRows: [][]string{
				{"comment", "created_on", "database_name", "dropped_on", "is_current", "is_default", "name", "options", "owner", "owner_role_type", "retention_time", "catalog_value", "catalog_level", "data_retention_time_in_days_value", "data_retention_time_in_days_level", "default_ddl_collation_value", "default_ddl_collation_level", "enable_console_output_value", "enable_console_output_level", "external_volume_value", "external_volume_level", "log_level_value", "log_level_level", "max_data_extension_time_in_days_value", "max_data_extension_time_in_days_level", "pipe_execution_paused_value", "pipe_execution_paused_level", "quoted_identifiers_ignore_case_value", "quoted_identifiers_ignore_case_level", "replace_invalid_characters_value", "replace_invalid_characters_level", "storage_serialization_policy_value", "storage_serialization_policy_level", "suspend_task_after_num_failures_value", "suspend_task_after_num_failures_level", "task_auto_retry_attempts_value", "task_auto_retry_attempts_level", "trace_level_value", "trace_level_level", "user_task_managed_initial_warehouse_size_value", "user_task_managed_initial_warehouse_size_level", "user_task_minimum_trigger_interval_in_seconds_value", "user_task_minimum_trigger_interval_in_seconds_level", "user_task_timeout_ms_value", "user_task_timeout_ms_level"},
				{"comment here", "2024-06-06 00:00:00.000 +0000 UTC", "DB", "", "true", "true", "SCHEMA2", "", "ADMIN", "ROLE", "10", "CAT", "ACCOUNT", "2", "ACCOUNT", "utf8", "ACCOUNT", "true", "ACCOUNT", "S3VOLUME", "ACCOUNT", "DEBUG", "ACCOUNT", "20", "ACCOUNT", "true", "ACCOUNT", "false", "ACCOUNT", "true", "ACCOUNT", "AVRO", "ACCOUNT", "7", "ACCOUNT", "4", "ACCOUNT", "PROD", "ACCOUNT", "SMALL", "ACCOUNT", "5", "ACCOUNT", "1000", "ACCOUNT"},
			},
			expectedOutput: `
resource "snowflake_schema" "snowflake_generated_schema_DB_SCHEMA2" {
  database = "DB"
  name = "SCHEMA2"
  comment = "comment here"
}
import {
  to = snowflake_schema.snowflake_generated_schema_DB_SCHEMA2
  id = "\"DB\".\"SCHEMA2\""
}
`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := HandleSchemas(&Config{
				ObjectType: ObjectTypeSchemas,
				ImportFlag: ImportStatementTypeBlock,
			}, tc.inputRows)
			assert.NoError(t, err)
			assert.Equal(t, strings.TrimLeft(tc.expectedOutput, "\n"), strings.TrimLeft(output, "\n"))
		})
	}
}
