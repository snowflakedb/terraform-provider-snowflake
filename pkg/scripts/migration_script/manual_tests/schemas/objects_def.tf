# =============================================================================
# Schemas Test Objects
# =============================================================================
# These resources create test schemas on Snowflake for migration testing.
# Run with: terraform apply
# Cleanup:  terraform destroy
#
# Edge cases covered:
# - Basic schema (minimal config)
# - Schema with comment
# - Schema with long comment
# - Schema with special characters in comment
# - Transient schema
# - Schema with managed access
# - Schema with data retention time
# - Schema with log level and trace level
# - Schema with multiple parameters
# - Schema with pipe_execution_paused
# =============================================================================

terraform {
  required_providers {
    snowflake = {
      source = "snowflakedb/snowflake"
    }
  }
}

provider "snowflake" {}

# ------------------------------------------------------------------------------
# Test Database (required for schemas)
# ------------------------------------------------------------------------------
resource "snowflake_database" "test_db" {
  name    = "MIGRATION_TEST_DB_FOR_SCHEMAS"
  comment = "Database for migration script testing - schemas"
}

# ------------------------------------------------------------------------------
# Basic Schema (minimal configuration)
# ------------------------------------------------------------------------------
resource "snowflake_schema" "basic" {
  database = snowflake_database.test_db.name
  name     = "MIGRATION_TEST_SCHEMA_BASIC"
}

# ------------------------------------------------------------------------------
# Schema with Comment
# ------------------------------------------------------------------------------
resource "snowflake_schema" "with_comment" {
  database = snowflake_database.test_db.name
  name     = "MIGRATION_TEST_SCHEMA_COMMENT"
  comment  = "This schema is used for migration testing purposes"
}

# ------------------------------------------------------------------------------
# Schema with Long Comment
# ------------------------------------------------------------------------------
resource "snowflake_schema" "long_comment" {
  database = snowflake_database.test_db.name
  name     = "MIGRATION_TEST_SCHEMA_LONG_COMMENT"
  comment  = "This is a very long comment that tests how the migration script handles comments with many characters. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."
}

# ------------------------------------------------------------------------------
# Schema with Special Characters in Comment
# ------------------------------------------------------------------------------
resource "snowflake_schema" "special_chars" {
  database = snowflake_database.test_db.name
  name     = "MIGRATION_TEST_SCHEMA_SPECIAL"
  comment  = "Comment with special chars: <>&\""
}

# ------------------------------------------------------------------------------
# Transient Schema
# ------------------------------------------------------------------------------
resource "snowflake_schema" "transient" {
  database     = snowflake_database.test_db.name
  name         = "MIGRATION_TEST_SCHEMA_TRANSIENT"
  is_transient = true
  comment      = "Transient schema for testing"
}

# ------------------------------------------------------------------------------
# Schema with Managed Access
# ------------------------------------------------------------------------------
resource "snowflake_schema" "managed_access" {
  database            = snowflake_database.test_db.name
  name                = "MIGRATION_TEST_SCHEMA_MANAGED"
  with_managed_access = true
  comment             = "Schema with managed access enabled"
}

# ------------------------------------------------------------------------------
# Schema with Custom Data Retention Time
# ------------------------------------------------------------------------------
resource "snowflake_schema" "data_retention" {
  database                    = snowflake_database.test_db.name
  name                        = "MIGRATION_TEST_SCHEMA_RETENTION"
  comment                     = "Schema with custom data retention"
  data_retention_time_in_days = 7
}

# ------------------------------------------------------------------------------
# Schema with Log Level and Trace Level
# ------------------------------------------------------------------------------
resource "snowflake_schema" "log_trace" {
  database    = snowflake_database.test_db.name
  name        = "MIGRATION_TEST_SCHEMA_LOG_TRACE"
  comment     = "Schema with log and trace level"
  log_level   = "INFO"
  trace_level = "ON_EVENT"
}

# ------------------------------------------------------------------------------
# Schema with Multiple Parameters
# ------------------------------------------------------------------------------
resource "snowflake_schema" "multi_params" {
  database                        = snowflake_database.test_db.name
  name                            = "MIGRATION_TEST_SCHEMA_MULTI_PARAMS"
  comment                         = "Schema with multiple parameters"
  data_retention_time_in_days     = 14
  max_data_extension_time_in_days = 28
  log_level                       = "DEBUG"
  trace_level                     = "ALWAYS"
  suspend_task_after_num_failures = 5
  task_auto_retry_attempts        = 3
  quoted_identifiers_ignore_case  = true
  replace_invalid_characters      = true
}

# ------------------------------------------------------------------------------
# Schema with Pipe Execution Paused
# ------------------------------------------------------------------------------
resource "snowflake_schema" "pipe_paused" {
  database              = snowflake_database.test_db.name
  name                  = "MIGRATION_TEST_SCHEMA_PIPE_PAUSED"
  comment               = "Schema with pipe execution paused"
  pipe_execution_paused = true
}

# ------------------------------------------------------------------------------
# Schema with Underscore in Name
# ------------------------------------------------------------------------------
resource "snowflake_schema" "underscore_name" {
  database = snowflake_database.test_db.name
  name     = "MIGRATION_TEST_SCHEMA_WITH_UNDERSCORE"
  comment  = "Schema with underscores in name"
}

# ------------------------------------------------------------------------------
# Schema with Numbers in Name
# ------------------------------------------------------------------------------
resource "snowflake_schema" "with_numbers" {
  database = snowflake_database.test_db.name
  name     = "MIGRATION_TEST_SCHEMA_123"
  comment  = "Schema with numbers in name"
}

# ------------------------------------------------------------------------------
# Transient Schema with Managed Access (edge case: both options)
# ------------------------------------------------------------------------------
resource "snowflake_schema" "transient_managed" {
  database            = snowflake_database.test_db.name
  name                = "MIGRATION_TEST_SCHEMA_TRANS_MANAGED"
  is_transient        = true
  with_managed_access = true
  comment             = "Transient schema with managed access"
}

# ------------------------------------------------------------------------------
# Schema with All Common Parameters
# ------------------------------------------------------------------------------
resource "snowflake_schema" "complete" {
  database                        = snowflake_database.test_db.name
  name                            = "MIGRATION_TEST_SCHEMA_COMPLETE"
  comment                         = "Complete schema with all common parameters"
  data_retention_time_in_days     = 30
  max_data_extension_time_in_days = 90
  log_level                       = "WARN"
  trace_level                     = "ON_EVENT"
  suspend_task_after_num_failures = 10
  task_auto_retry_attempts        = 2
  pipe_execution_paused           = false
  quoted_identifiers_ignore_case  = false
  replace_invalid_characters      = false
  enable_console_output           = true
}

