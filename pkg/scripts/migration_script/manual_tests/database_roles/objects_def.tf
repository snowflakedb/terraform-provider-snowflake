# =============================================================================
# Database Roles Test Objects
# =============================================================================
# These resources create test database roles on Snowflake for migration testing.
# Run with: terraform apply
# Cleanup:  terraform destroy
#
# Note: Database roles require a database to exist first.
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
# Test Database (required for database roles)
# ------------------------------------------------------------------------------
resource "snowflake_database" "test_db" {
  name    = "MIGRATION_TEST_DB_FOR_ROLES"
  comment = "Database for migration script testing - database roles"
}

# ------------------------------------------------------------------------------
# Basic Database Role (minimal configuration)
# ------------------------------------------------------------------------------
resource "snowflake_database_role" "basic" {
  database = snowflake_database.test_db.name
  name     = "MIGRATION_TEST_DBROLE_BASIC"
}

# ------------------------------------------------------------------------------
# Database Role with Comment
# ------------------------------------------------------------------------------
resource "snowflake_database_role" "with_comment" {
  database = snowflake_database.test_db.name
  name     = "MIGRATION_TEST_DBROLE_COMMENT"
  comment  = "This database role is used for migration testing purposes"
}

# ------------------------------------------------------------------------------
# Database Role with Long Comment
# ------------------------------------------------------------------------------
resource "snowflake_database_role" "long_comment" {
  database = snowflake_database.test_db.name
  name     = "MIGRATION_TEST_DBROLE_LONG_COMMENT"
  comment  = "This is a very long comment that tests how the migration script handles comments with many characters. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."
}

# ------------------------------------------------------------------------------
# Database Role with Special Characters in Comment
# ------------------------------------------------------------------------------
resource "snowflake_database_role" "special_chars" {
  database = snowflake_database.test_db.name
  name     = "MIGRATION_TEST_DBROLE_SPECIAL"
  comment  = "Comment with special chars: <>&\""
}

# ------------------------------------------------------------------------------
# Database Role for Testing Hierarchy (Parent)
# ------------------------------------------------------------------------------
resource "snowflake_database_role" "parent" {
  database = snowflake_database.test_db.name
  name     = "MIGRATION_TEST_DBROLE_PARENT"
  comment  = "Parent database role for hierarchy testing"
}

# ------------------------------------------------------------------------------
# Database Role for Testing Hierarchy (Child)
# ------------------------------------------------------------------------------
resource "snowflake_database_role" "child" {
  database = snowflake_database.test_db.name
  name     = "MIGRATION_TEST_DBROLE_CHILD"
  comment  = "Child database role for hierarchy testing"
}

# ------------------------------------------------------------------------------
# Database Role with Underscore in Name
# ------------------------------------------------------------------------------
resource "snowflake_database_role" "underscore_name" {
  database = snowflake_database.test_db.name
  name     = "MIGRATION_TEST_DBROLE_WITH_UNDERSCORE"
  comment  = "Role with underscores in name"
}

# ------------------------------------------------------------------------------
# Database Role with Numbers in Name
# ------------------------------------------------------------------------------
resource "snowflake_database_role" "with_numbers" {
  database = snowflake_database.test_db.name
  name     = "MIGRATION_TEST_DBROLE_123"
  comment  = "Role with numbers in name"
}

