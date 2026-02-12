# Grant types covered:
# 1. Grant Account Role to Role
# 2. Grant Account Role to User
# 3. Grant Database Role to Database Role
# 4. Grant Database Role to Account Role
# 5. Grant Privileges to Account Role:
#    - On Account
#    - On Account Object (DATABASE, WAREHOUSE)
#    - On Schema
#    - On Schema Object (TABLE)
# 6. Grant Privileges to Database Role:
#    - On Database
#    - On Schema
#    - On Schema Object (TABLE)
# 7. Grant with grant_option = true

terraform {
  required_providers {
    snowflake = {
      source  = "snowflakedb/snowflake"
      version = "= 2.12.0"
    }
    random = {
      source = "hashicorp/random"
    }
  }
}

provider "snowflake" {
  # Uses default configuration from ~/.snowflake/config or environment variables
  preview_features_enabled = ["snowflake_table_resource"]
}

# ------------------------------------------------------------------------------
# UUID for unique naming
# ------------------------------------------------------------------------------
resource "random_uuid" "test_run" {}

locals {
  # Use first 8 chars of UUID for shorter names
  test_id = upper(substr(random_uuid.test_run.result, 0, 8))
  prefix  = "MIGRATION_TEST_${local.test_id}"
}

# Output the prefix for use in datasource
output "test_prefix" {
  description = "The unique prefix used for all test objects"
  value       = local.prefix
}

# ==============================================================================
# PREREQUISITE OBJECTS
# ==============================================================================

# ------------------------------------------------------------------------------
# Account Roles
# ------------------------------------------------------------------------------
resource "snowflake_account_role" "parent_role" {
  name    = "${local.prefix}_PARENT_ROLE"
  comment = "Parent role for grant testing"
}

resource "snowflake_account_role" "child_role" {
  name    = "${local.prefix}_CHILD_ROLE"
  comment = "Child role for grant testing"
}

resource "snowflake_account_role" "priv_role" {
  name    = "${local.prefix}_PRIV_ROLE"
  comment = "Role for privilege grant testing"
}

# ------------------------------------------------------------------------------
# Database
# ------------------------------------------------------------------------------
resource "snowflake_database" "test_db" {
  name    = "${local.prefix}_DB"
  comment = "Database for grant testing\nLine 2 of comment\nLine 3 with special chars: \"quotes\" and \\backslash"
}

# ------------------------------------------------------------------------------
# Database Roles
# ------------------------------------------------------------------------------
resource "snowflake_database_role" "parent_db_role" {
  database = snowflake_database.test_db.name
  name     = "${local.prefix}_PARENT_DBROLE"
  comment  = "Parent database role for grant testing"
}

resource "snowflake_database_role" "child_db_role" {
  database = snowflake_database.test_db.name
  name     = "${local.prefix}_CHILD_DBROLE"
  comment  = "Child database role for grant testing"
}

resource "snowflake_database_role" "priv_db_role" {
  database = snowflake_database.test_db.name
  name     = "${local.prefix}_PRIV_DBROLE"
  comment  = "Database role for privilege grant testing"
}

# ------------------------------------------------------------------------------
# Schema
# ------------------------------------------------------------------------------
resource "snowflake_schema" "test_schema" {
  database = snowflake_database.test_db.name
  name     = "${local.prefix}_SCHEMA"
  comment  = "Schema for grant testing"
}

# ------------------------------------------------------------------------------
# Table (for schema object grants)
# ------------------------------------------------------------------------------
resource "snowflake_table" "test_table" {
  database = snowflake_database.test_db.name
  schema   = snowflake_schema.test_schema.name
  name     = "${local.prefix}_TABLE"
  comment  = "Table for grant testing"

  column {
    name = "ID"
    type = "NUMBER(38,0)"
  }

  column {
    name = "NAME"
    type = "VARCHAR(100)"
  }
}

# ------------------------------------------------------------------------------
# Warehouse (for account object grants)
# ------------------------------------------------------------------------------
resource "snowflake_warehouse" "test_wh" {
  name           = "${local.prefix}_WH"
  warehouse_size = "XSMALL"
  comment        = "Warehouse for grant testing"
}

# ==============================================================================
# GRANT ACCOUNT ROLE TESTS
# ==============================================================================

# ------------------------------------------------------------------------------
# Grant Account Role to Role
# ------------------------------------------------------------------------------
resource "snowflake_grant_account_role" "role_to_role" {
  role_name        = snowflake_account_role.child_role.name
  parent_role_name = snowflake_account_role.parent_role.name
}

# ==============================================================================
# GRANT DATABASE ROLE TESTS
# ==============================================================================

# ------------------------------------------------------------------------------
# Grant Database Role to Database Role
# ------------------------------------------------------------------------------
resource "snowflake_grant_database_role" "dbrole_to_dbrole" {
  database_role_name        = snowflake_database_role.child_db_role.fully_qualified_name
  parent_database_role_name = snowflake_database_role.parent_db_role.fully_qualified_name
}

# ------------------------------------------------------------------------------
# Grant Database Role to Account Role
# ------------------------------------------------------------------------------
resource "snowflake_grant_database_role" "dbrole_to_role" {
  database_role_name = snowflake_database_role.priv_db_role.fully_qualified_name
  parent_role_name   = snowflake_account_role.priv_role.name
}

# ==============================================================================
# GRANT PRIVILEGES TO ACCOUNT ROLE TESTS
# ==============================================================================

# ------------------------------------------------------------------------------
# Grant Privileges on Account
# ------------------------------------------------------------------------------
resource "snowflake_grant_privileges_to_account_role" "on_account" {
  account_role_name = snowflake_account_role.priv_role.name
  privileges        = ["CREATE DATABASE"]
  on_account        = true
}

# ------------------------------------------------------------------------------
# Grant Privileges on Account Object (Database)
# ------------------------------------------------------------------------------
resource "snowflake_grant_privileges_to_account_role" "on_database" {
  account_role_name = snowflake_account_role.priv_role.name
  privileges        = ["USAGE", "CREATE SCHEMA"]

  on_account_object {
    object_type = "DATABASE"
    object_name = snowflake_database.test_db.name
  }
}

# ------------------------------------------------------------------------------
# Grant Privileges on Account Object (Warehouse)
# ------------------------------------------------------------------------------
resource "snowflake_grant_privileges_to_account_role" "on_warehouse" {
  account_role_name = snowflake_account_role.priv_role.name
  privileges        = ["USAGE", "OPERATE"]

  on_account_object {
    object_type = "WAREHOUSE"
    object_name = snowflake_warehouse.test_wh.name
  }
}

# ------------------------------------------------------------------------------
# Grant Privileges on Schema
# ------------------------------------------------------------------------------
resource "snowflake_grant_privileges_to_account_role" "on_schema" {
  account_role_name = snowflake_account_role.priv_role.name
  privileges        = ["USAGE", "CREATE TABLE", "CREATE VIEW"]

  on_schema {
    schema_name = snowflake_schema.test_schema.fully_qualified_name
  }
}

# ------------------------------------------------------------------------------
# Grant Privileges on Schema Object (Table)
# ------------------------------------------------------------------------------
resource "snowflake_grant_privileges_to_account_role" "on_table" {
  account_role_name = snowflake_account_role.priv_role.name
  privileges        = ["SELECT", "INSERT", "UPDATE"]

  on_schema_object {
    object_type = "TABLE"
    object_name = snowflake_table.test_table.fully_qualified_name
  }
}

# ------------------------------------------------------------------------------
# Grant Privileges with grant_option = true
# ------------------------------------------------------------------------------
resource "snowflake_grant_privileges_to_account_role" "with_grant_option" {
  account_role_name = snowflake_account_role.priv_role.name
  privileges        = ["DELETE"]
  with_grant_option = true

  on_schema_object {
    object_type = "TABLE"
    object_name = snowflake_table.test_table.fully_qualified_name
  }
}

# ==============================================================================
# GRANT PRIVILEGES TO DATABASE ROLE TESTS
# ==============================================================================

# ------------------------------------------------------------------------------
# Grant Privileges on Database (to database role)
# ------------------------------------------------------------------------------
resource "snowflake_grant_privileges_to_database_role" "on_database" {
  database_role_name = snowflake_database_role.priv_db_role.fully_qualified_name
  privileges         = ["USAGE"]
  on_database        = snowflake_database.test_db.name
}

# ------------------------------------------------------------------------------
# Grant Privileges on Schema (to database role)
# ------------------------------------------------------------------------------
resource "snowflake_grant_privileges_to_database_role" "on_schema" {
  database_role_name = snowflake_database_role.priv_db_role.fully_qualified_name
  privileges         = ["USAGE", "CREATE TABLE"]

  on_schema {
    schema_name = snowflake_schema.test_schema.fully_qualified_name
  }
}

# ------------------------------------------------------------------------------
# Grant Privileges on Schema Object (Table, to database role)
# ------------------------------------------------------------------------------
resource "snowflake_grant_privileges_to_database_role" "on_table" {
  database_role_name = snowflake_database_role.priv_db_role.fully_qualified_name
  privileges         = ["SELECT"]

  on_schema_object {
    object_type = "TABLE"
    object_name = snowflake_table.test_table.fully_qualified_name
  }
}

# ------------------------------------------------------------------------------
# Grant Privileges to Database Role with grant_option
# ------------------------------------------------------------------------------
resource "snowflake_grant_privileges_to_database_role" "with_grant_option" {
  database_role_name = snowflake_database_role.priv_db_role.fully_qualified_name
  privileges         = ["INSERT"]
  with_grant_option  = true

  on_schema_object {
    object_type = "TABLE"
    object_name = snowflake_table.test_table.fully_qualified_name
  }
}
