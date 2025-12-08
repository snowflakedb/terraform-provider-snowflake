# =============================================================================
# Account Roles Test Objects
# =============================================================================
# These resources create test account roles on Snowflake for migration testing.
# Run with: terraform apply
# Cleanup:  terraform destroy
# =============================================================================

terraform {
  required_providers {
    snowflake = {
      source  = "snowflakedb/snowflake"
    }
  }
}

provider "snowflake" {
  # Configure via environment variables:
  # SNOWFLAKE_ACCOUNT, SNOWFLAKE_USER, SNOWFLAKE_PASSWORD, etc.
}

# ------------------------------------------------------------------------------
# Basic Account Role (minimal configuration)
# ------------------------------------------------------------------------------
resource "snowflake_account_role" "basic" {
  name = "MIGRATION_TEST_ROLE_BASIC"
}

# ------------------------------------------------------------------------------
# Account Role with Comment
# ------------------------------------------------------------------------------
resource "snowflake_account_role" "with_comment" {
  name    = "MIGRATION_TEST_ROLE_COMMENT"
  comment = "This role is used for migration testing purposes"
}

# ------------------------------------------------------------------------------
# Account Role with Long Comment
# ------------------------------------------------------------------------------
resource "snowflake_account_role" "long_comment" {
  name    = "MIGRATION_TEST_ROLE_LONG_COMMENT"
  comment = "This is a very long comment that tests how the migration script handles comments with many characters. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."
}

# ------------------------------------------------------------------------------
# Account Role with Special Characters in Comment
# ------------------------------------------------------------------------------
resource "snowflake_account_role" "special_chars" {
  name    = "MIGRATION_TEST_ROLE_SPECIAL"
  comment = "Comment with special chars: <>&\""
}

# ------------------------------------------------------------------------------
# Account Role for Testing Hierarchy (Parent)
# ------------------------------------------------------------------------------
resource "snowflake_account_role" "parent" {
  name    = "MIGRATION_TEST_ROLE_PARENT"
  comment = "Parent role for hierarchy testing"
}

# ------------------------------------------------------------------------------
# Account Role for Testing Hierarchy (Child)
# ------------------------------------------------------------------------------
resource "snowflake_account_role" "child" {
  name    = "MIGRATION_TEST_ROLE_CHILD"
  comment = "Child role for hierarchy testing"
}

