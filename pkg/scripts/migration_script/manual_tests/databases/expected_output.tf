# =============================================================================
# Expected Migration Script Output - Databases
# =============================================================================
# This file contains the EXPECTED output from running:
#   go run .. -import=block databases < objects.csv
#
# Use this to compare with actual output:
#   go run .. -import=block databases < objects.csv > actual_output.tf
#   diff expected_output.tf actual_output.tf
#
# NOTE:
# - Resources are sorted alphabetically by name
# - Import blocks are grouped at the end
# - Special characters are Unicode-escaped in HCL output
# - Only DATABASE-level parameters are included (not ACCOUNT-level defaults)
# =============================================================================

resource "snowflake_database" "snowflake_generated_database_MIGRATION_TEST_DB_123" {
  name = "MIGRATION_TEST_DB_123"
  comment = "Database with numbers in name"
}

resource "snowflake_database" "snowflake_generated_database_MIGRATION_TEST_DB_BASIC" {
  name = "MIGRATION_TEST_DB_BASIC"
}

resource "snowflake_database" "snowflake_generated_database_MIGRATION_TEST_DB_COMMENT" {
  name = "MIGRATION_TEST_DB_COMMENT"
  comment = "This database is used for migration testing purposes"
}

resource "snowflake_database" "snowflake_generated_database_MIGRATION_TEST_DB_LOG_TRACE" {
  name = "MIGRATION_TEST_DB_LOG_TRACE"
  comment = "Database with log and trace level"
  log_level = "INFO"
  trace_level = "ON_EVENT"
}

resource "snowflake_database" "snowflake_generated_database_MIGRATION_TEST_DB_LONG_COMMENT" {
  name = "MIGRATION_TEST_DB_LONG_COMMENT"
  comment = "This is a very long comment that tests how the migration script handles comments with many characters. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."
}

resource "snowflake_database" "snowflake_generated_database_MIGRATION_TEST_DB_MULTI_PARAMS" {
  name = "MIGRATION_TEST_DB_MULTI_PARAMS"
  comment = "Database with multiple parameters"
  data_retention_time_in_days = 14
  log_level = "DEBUG"
  max_data_extension_time_in_days = 28
  quoted_identifiers_ignore_case = true
  replace_invalid_characters = true
  suspend_task_after_num_failures = 5
  task_auto_retry_attempts = 3
  trace_level = "ALWAYS"
}

resource "snowflake_database" "snowflake_generated_database_MIGRATION_TEST_DB_RETENTION" {
  name = "MIGRATION_TEST_DB_RETENTION"
  comment = "Database with custom data retention"
  data_retention_time_in_days = 7
}

resource "snowflake_database" "snowflake_generated_database_MIGRATION_TEST_DB_SPECIAL" {
  name = "MIGRATION_TEST_DB_SPECIAL"
  comment = "Comment with special chars: \u003c\u003e\u0026\""
}

resource "snowflake_database" "snowflake_generated_database_MIGRATION_TEST_DB_TRANSIENT" {
  name = "MIGRATION_TEST_DB_TRANSIENT"
  comment = "Transient database for testing"
  is_transient = true
}

resource "snowflake_database" "snowflake_generated_database_MIGRATION_TEST_DB_TRANSIENT_PARAMS" {
  name = "MIGRATION_TEST_DB_TRANSIENT_PARAMS"
  comment = "Transient database with parameters"
  data_retention_time_in_days = 0
  is_transient = true
  log_level = "WARN"
}

resource "snowflake_database" "snowflake_generated_database_MIGRATION_TEST_DB_WITH_UNDERSCORE" {
  name = "MIGRATION_TEST_DB_WITH_UNDERSCORE"
  comment = "Database with underscores in name"
}

import {
  to = snowflake_database.snowflake_generated_database_MIGRATION_TEST_DB_123
  id = "\"MIGRATION_TEST_DB_123\""
}
import {
  to = snowflake_database.snowflake_generated_database_MIGRATION_TEST_DB_BASIC
  id = "\"MIGRATION_TEST_DB_BASIC\""
}
import {
  to = snowflake_database.snowflake_generated_database_MIGRATION_TEST_DB_COMMENT
  id = "\"MIGRATION_TEST_DB_COMMENT\""
}
import {
  to = snowflake_database.snowflake_generated_database_MIGRATION_TEST_DB_LOG_TRACE
  id = "\"MIGRATION_TEST_DB_LOG_TRACE\""
}
import {
  to = snowflake_database.snowflake_generated_database_MIGRATION_TEST_DB_LONG_COMMENT
  id = "\"MIGRATION_TEST_DB_LONG_COMMENT\""
}
import {
  to = snowflake_database.snowflake_generated_database_MIGRATION_TEST_DB_MULTI_PARAMS
  id = "\"MIGRATION_TEST_DB_MULTI_PARAMS\""
}
import {
  to = snowflake_database.snowflake_generated_database_MIGRATION_TEST_DB_RETENTION
  id = "\"MIGRATION_TEST_DB_RETENTION\""
}
import {
  to = snowflake_database.snowflake_generated_database_MIGRATION_TEST_DB_SPECIAL
  id = "\"MIGRATION_TEST_DB_SPECIAL\""
}
import {
  to = snowflake_database.snowflake_generated_database_MIGRATION_TEST_DB_TRANSIENT
  id = "\"MIGRATION_TEST_DB_TRANSIENT\""
}
import {
  to = snowflake_database.snowflake_generated_database_MIGRATION_TEST_DB_TRANSIENT_PARAMS
  id = "\"MIGRATION_TEST_DB_TRANSIENT_PARAMS\""
}
import {
  to = snowflake_database.snowflake_generated_database_MIGRATION_TEST_DB_WITH_UNDERSCORE
  id = "\"MIGRATION_TEST_DB_WITH_UNDERSCORE\""
}

