# =============================================================================
# Databases Test Objects
# =============================================================================
# These resources create test databases on Snowflake for migration testing.
# Run with: terraform apply
# Cleanup:  terraform destroy
#
# Edge cases covered:
# - Basic database (minimal config)
# - Transient database
# - Database with comment
# - Database with long comment
# - Database with special characters
# - Database with custom data retention
# - Database with log/trace level
# - Database with multiple parameters
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
# Basic Database (minimal configuration)
# ------------------------------------------------------------------------------
resource "snowflake_database" "basic" {
  name = "MIGRATION_TEST_DB_BASIC"
}

# ------------------------------------------------------------------------------
# Database with Comment
# ------------------------------------------------------------------------------
resource "snowflake_database" "with_comment" {
  name    = "MIGRATION_TEST_DB_COMMENT"
  comment = "This database is used for migration testing purposes"
}

# ------------------------------------------------------------------------------
# Transient Database
# ------------------------------------------------------------------------------
resource "snowflake_database" "transient" {
  name         = "MIGRATION_TEST_DB_TRANSIENT"
  is_transient = true
  comment      = "Transient database for testing"
}

# ------------------------------------------------------------------------------
# Database with Long Comment
# ------------------------------------------------------------------------------
resource "snowflake_database" "long_comment" {
  name    = "MIGRATION_TEST_DB_LONG_COMMENT"
  comment = "This is a very long comment that tests how the migration script handles comments with many characters. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."
}

# ------------------------------------------------------------------------------
# Database with Special Characters in Comment
# ------------------------------------------------------------------------------
resource "snowflake_database" "special_chars" {
  name    = "MIGRATION_TEST_DB_SPECIAL"
  comment = "Comment with special chars: <>&\""
}

# ------------------------------------------------------------------------------
# Database with Custom Data Retention Time
# ------------------------------------------------------------------------------
resource "snowflake_database" "data_retention" {
  name                        = "MIGRATION_TEST_DB_RETENTION"
  comment                     = "Database with custom data retention"
  data_retention_time_in_days = 7
}

# ------------------------------------------------------------------------------
# Database with Log Level and Trace Level
# ------------------------------------------------------------------------------
resource "snowflake_database" "log_trace" {
  name        = "MIGRATION_TEST_DB_LOG_TRACE"
  comment     = "Database with log and trace level"
  log_level   = "INFO"
  trace_level = "ON_EVENT"
}

# ------------------------------------------------------------------------------
# Database with Multiple Parameters
# ------------------------------------------------------------------------------
resource "snowflake_database" "multi_params" {
  name                            = "MIGRATION_TEST_DB_MULTI_PARAMS"
  comment                         = "Database with multiple parameters"
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
# Database with Underscore in Name
# ------------------------------------------------------------------------------
resource "snowflake_database" "underscore_name" {
  name    = "MIGRATION_TEST_DB_WITH_UNDERSCORE"
  comment = "Database with underscores in name"
}

# ------------------------------------------------------------------------------
# Database with Numbers in Name
# ------------------------------------------------------------------------------
resource "snowflake_database" "with_numbers" {
  name    = "MIGRATION_TEST_DB_123"
  comment = "Database with numbers in name"
}

# ------------------------------------------------------------------------------
# Transient Database with Parameters (edge case: transient + params)
# ------------------------------------------------------------------------------
resource "snowflake_database" "transient_with_params" {
  name                        = "MIGRATION_TEST_DB_TRANSIENT_PARAMS"
  is_transient                = true
  comment                     = "Transient database with parameters"
  data_retention_time_in_days = 0  # Transient databases often have 0 retention
  log_level                   = "WARN"
}

