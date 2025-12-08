# =============================================================================
# Expected Migration Script Output - Schemas
# =============================================================================
# This file contains the EXPECTED output from running:
#   go run .. -import=block schemas < objects.csv
#
# Use this to compare with actual output:
#   go run .. -import=block schemas < objects.csv > actual_output.tf
#   diff expected_output.tf actual_output.tf
#
# NOTE:
# - Resources are sorted alphabetically by name
# - Import blocks are grouped at the end
# - Special characters are Unicode-escaped in HCL output
# - Boolean values like is_transient and with_managed_access are strings ("true")
# - Only SCHEMA-level parameters are included (not DATABASE/ACCOUNT-level defaults)
# - Import IDs use fully qualified format: "DATABASE"."SCHEMA"
# =============================================================================

resource "snowflake_schema" "snowflake_generated_schema_MIGRATION_TEST_DB_FOR_SCHEMAS_MIGRATION_TEST_SCHEMA_123" {
  database = "MIGRATION_TEST_DB_FOR_SCHEMAS"
  name = "MIGRATION_TEST_SCHEMA_123"
  comment = "Schema with numbers in name"
}

resource "snowflake_schema" "snowflake_generated_schema_MIGRATION_TEST_DB_FOR_SCHEMAS_MIGRATION_TEST_SCHEMA_BASIC" {
  database = "MIGRATION_TEST_DB_FOR_SCHEMAS"
  name = "MIGRATION_TEST_SCHEMA_BASIC"
}

resource "snowflake_schema" "snowflake_generated_schema_MIGRATION_TEST_DB_FOR_SCHEMAS_MIGRATION_TEST_SCHEMA_COMMENT" {
  database = "MIGRATION_TEST_DB_FOR_SCHEMAS"
  name = "MIGRATION_TEST_SCHEMA_COMMENT"
  comment = "This schema is used for migration testing purposes"
}

resource "snowflake_schema" "snowflake_generated_schema_MIGRATION_TEST_DB_FOR_SCHEMAS_MIGRATION_TEST_SCHEMA_COMPLETE" {
  database = "MIGRATION_TEST_DB_FOR_SCHEMAS"
  name = "MIGRATION_TEST_SCHEMA_COMPLETE"
  comment = "Complete schema with all common parameters"
  data_retention_time_in_days = 30
  enable_console_output = true
  log_level = "WARN"
  max_data_extension_time_in_days = 90
  pipe_execution_paused = false
  quoted_identifiers_ignore_case = false
  replace_invalid_characters = false
  suspend_task_after_num_failures = 10
  task_auto_retry_attempts = 2
  trace_level = "ON_EVENT"
}

resource "snowflake_schema" "snowflake_generated_schema_MIGRATION_TEST_DB_FOR_SCHEMAS_MIGRATION_TEST_SCHEMA_LOG_TRACE" {
  database = "MIGRATION_TEST_DB_FOR_SCHEMAS"
  name = "MIGRATION_TEST_SCHEMA_LOG_TRACE"
  comment = "Schema with log and trace level"
  log_level = "INFO"
  trace_level = "ON_EVENT"
}

resource "snowflake_schema" "snowflake_generated_schema_MIGRATION_TEST_DB_FOR_SCHEMAS_MIGRATION_TEST_SCHEMA_LONG_COMMENT" {
  database = "MIGRATION_TEST_DB_FOR_SCHEMAS"
  name = "MIGRATION_TEST_SCHEMA_LONG_COMMENT"
  comment = "This is a very long comment that tests how the migration script handles comments with many characters. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."
}

resource "snowflake_schema" "snowflake_generated_schema_MIGRATION_TEST_DB_FOR_SCHEMAS_MIGRATION_TEST_SCHEMA_MANAGED" {
  database = "MIGRATION_TEST_DB_FOR_SCHEMAS"
  name = "MIGRATION_TEST_SCHEMA_MANAGED"
  comment = "Schema with managed access enabled"
  with_managed_access = "true"
}

resource "snowflake_schema" "snowflake_generated_schema_MIGRATION_TEST_DB_FOR_SCHEMAS_MIGRATION_TEST_SCHEMA_MULTI_PARAMS" {
  database = "MIGRATION_TEST_DB_FOR_SCHEMAS"
  name = "MIGRATION_TEST_SCHEMA_MULTI_PARAMS"
  comment = "Schema with multiple parameters"
  data_retention_time_in_days = 14
  log_level = "DEBUG"
  max_data_extension_time_in_days = 28
  quoted_identifiers_ignore_case = true
  replace_invalid_characters = true
  suspend_task_after_num_failures = 5
  task_auto_retry_attempts = 3
  trace_level = "ALWAYS"
}

resource "snowflake_schema" "snowflake_generated_schema_MIGRATION_TEST_DB_FOR_SCHEMAS_MIGRATION_TEST_SCHEMA_PIPE_PAUSED" {
  database = "MIGRATION_TEST_DB_FOR_SCHEMAS"
  name = "MIGRATION_TEST_SCHEMA_PIPE_PAUSED"
  comment = "Schema with pipe execution paused"
  pipe_execution_paused = true
}

resource "snowflake_schema" "snowflake_generated_schema_MIGRATION_TEST_DB_FOR_SCHEMAS_MIGRATION_TEST_SCHEMA_RETENTION" {
  database = "MIGRATION_TEST_DB_FOR_SCHEMAS"
  name = "MIGRATION_TEST_SCHEMA_RETENTION"
  comment = "Schema with custom data retention"
  data_retention_time_in_days = 7
}

resource "snowflake_schema" "snowflake_generated_schema_MIGRATION_TEST_DB_FOR_SCHEMAS_MIGRATION_TEST_SCHEMA_SPECIAL" {
  database = "MIGRATION_TEST_DB_FOR_SCHEMAS"
  name = "MIGRATION_TEST_SCHEMA_SPECIAL"
  comment = "Comment with special chars: \u003c\u003e\u0026\""
}

resource "snowflake_schema" "snowflake_generated_schema_MIGRATION_TEST_DB_FOR_SCHEMAS_MIGRATION_TEST_SCHEMA_TRANSIENT" {
  database = "MIGRATION_TEST_DB_FOR_SCHEMAS"
  name = "MIGRATION_TEST_SCHEMA_TRANSIENT"
  comment = "Transient schema for testing"
  is_transient = "true"
}

resource "snowflake_schema" "snowflake_generated_schema_MIGRATION_TEST_DB_FOR_SCHEMAS_MIGRATION_TEST_SCHEMA_TRANS_MANAGED" {
  database = "MIGRATION_TEST_DB_FOR_SCHEMAS"
  name = "MIGRATION_TEST_SCHEMA_TRANS_MANAGED"
  comment = "Transient schema with managed access"
  is_transient = "true"
  with_managed_access = "true"
}

resource "snowflake_schema" "snowflake_generated_schema_MIGRATION_TEST_DB_FOR_SCHEMAS_MIGRATION_TEST_SCHEMA_WITH_UNDERSCORE" {
  database = "MIGRATION_TEST_DB_FOR_SCHEMAS"
  name = "MIGRATION_TEST_SCHEMA_WITH_UNDERSCORE"
  comment = "Schema with underscores in name"
}

import {
  to = snowflake_schema.snowflake_generated_schema_MIGRATION_TEST_DB_FOR_SCHEMAS_MIGRATION_TEST_SCHEMA_123
  id = "\"MIGRATION_TEST_DB_FOR_SCHEMAS\".\"MIGRATION_TEST_SCHEMA_123\""
}
import {
  to = snowflake_schema.snowflake_generated_schema_MIGRATION_TEST_DB_FOR_SCHEMAS_MIGRATION_TEST_SCHEMA_BASIC
  id = "\"MIGRATION_TEST_DB_FOR_SCHEMAS\".\"MIGRATION_TEST_SCHEMA_BASIC\""
}
import {
  to = snowflake_schema.snowflake_generated_schema_MIGRATION_TEST_DB_FOR_SCHEMAS_MIGRATION_TEST_SCHEMA_COMMENT
  id = "\"MIGRATION_TEST_DB_FOR_SCHEMAS\".\"MIGRATION_TEST_SCHEMA_COMMENT\""
}
import {
  to = snowflake_schema.snowflake_generated_schema_MIGRATION_TEST_DB_FOR_SCHEMAS_MIGRATION_TEST_SCHEMA_COMPLETE
  id = "\"MIGRATION_TEST_DB_FOR_SCHEMAS\".\"MIGRATION_TEST_SCHEMA_COMPLETE\""
}
import {
  to = snowflake_schema.snowflake_generated_schema_MIGRATION_TEST_DB_FOR_SCHEMAS_MIGRATION_TEST_SCHEMA_LOG_TRACE
  id = "\"MIGRATION_TEST_DB_FOR_SCHEMAS\".\"MIGRATION_TEST_SCHEMA_LOG_TRACE\""
}
import {
  to = snowflake_schema.snowflake_generated_schema_MIGRATION_TEST_DB_FOR_SCHEMAS_MIGRATION_TEST_SCHEMA_LONG_COMMENT
  id = "\"MIGRATION_TEST_DB_FOR_SCHEMAS\".\"MIGRATION_TEST_SCHEMA_LONG_COMMENT\""
}
import {
  to = snowflake_schema.snowflake_generated_schema_MIGRATION_TEST_DB_FOR_SCHEMAS_MIGRATION_TEST_SCHEMA_MANAGED
  id = "\"MIGRATION_TEST_DB_FOR_SCHEMAS\".\"MIGRATION_TEST_SCHEMA_MANAGED\""
}
import {
  to = snowflake_schema.snowflake_generated_schema_MIGRATION_TEST_DB_FOR_SCHEMAS_MIGRATION_TEST_SCHEMA_MULTI_PARAMS
  id = "\"MIGRATION_TEST_DB_FOR_SCHEMAS\".\"MIGRATION_TEST_SCHEMA_MULTI_PARAMS\""
}
import {
  to = snowflake_schema.snowflake_generated_schema_MIGRATION_TEST_DB_FOR_SCHEMAS_MIGRATION_TEST_SCHEMA_PIPE_PAUSED
  id = "\"MIGRATION_TEST_DB_FOR_SCHEMAS\".\"MIGRATION_TEST_SCHEMA_PIPE_PAUSED\""
}
import {
  to = snowflake_schema.snowflake_generated_schema_MIGRATION_TEST_DB_FOR_SCHEMAS_MIGRATION_TEST_SCHEMA_RETENTION
  id = "\"MIGRATION_TEST_DB_FOR_SCHEMAS\".\"MIGRATION_TEST_SCHEMA_RETENTION\""
}
import {
  to = snowflake_schema.snowflake_generated_schema_MIGRATION_TEST_DB_FOR_SCHEMAS_MIGRATION_TEST_SCHEMA_SPECIAL
  id = "\"MIGRATION_TEST_DB_FOR_SCHEMAS\".\"MIGRATION_TEST_SCHEMA_SPECIAL\""
}
import {
  to = snowflake_schema.snowflake_generated_schema_MIGRATION_TEST_DB_FOR_SCHEMAS_MIGRATION_TEST_SCHEMA_TRANSIENT
  id = "\"MIGRATION_TEST_DB_FOR_SCHEMAS\".\"MIGRATION_TEST_SCHEMA_TRANSIENT\""
}
import {
  to = snowflake_schema.snowflake_generated_schema_MIGRATION_TEST_DB_FOR_SCHEMAS_MIGRATION_TEST_SCHEMA_TRANS_MANAGED
  id = "\"MIGRATION_TEST_DB_FOR_SCHEMAS\".\"MIGRATION_TEST_SCHEMA_TRANS_MANAGED\""
}
import {
  to = snowflake_schema.snowflake_generated_schema_MIGRATION_TEST_DB_FOR_SCHEMAS_MIGRATION_TEST_SCHEMA_WITH_UNDERSCORE
  id = "\"MIGRATION_TEST_DB_FOR_SCHEMAS\".\"MIGRATION_TEST_SCHEMA_WITH_UNDERSCORE\""
}
