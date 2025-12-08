# =============================================================================
# Expected Migration Script Output - Database Roles
# =============================================================================
# This file contains the EXPECTED output from running:
#   go run .. -import=block database_roles < objects.csv
#
# Use this to compare with actual output:
#   go run .. -import=block database_roles < objects.csv > actual_output.tf
#   diff expected_output.tf actual_output.tf
#
# NOTE:
# - Resources are sorted alphabetically by name
# - Import blocks are grouped at the end
# - Special characters are Unicode-escaped in HCL output
# - Database roles use fully qualified import IDs: "DATABASE"."ROLE_NAME"
# =============================================================================

resource "snowflake_database_role" "snowflake_generated_database_role_MIGRATION_TEST_DB_FOR_ROLES_MIGRATION_TEST_DBROLE_123" {
  database = "MIGRATION_TEST_DB_FOR_ROLES"
  name = "MIGRATION_TEST_DBROLE_123"
  comment = "Role with numbers in name"
}

resource "snowflake_database_role" "snowflake_generated_database_role_MIGRATION_TEST_DB_FOR_ROLES_MIGRATION_TEST_DBROLE_BASIC" {
  database = "MIGRATION_TEST_DB_FOR_ROLES"
  name = "MIGRATION_TEST_DBROLE_BASIC"
}

resource "snowflake_database_role" "snowflake_generated_database_role_MIGRATION_TEST_DB_FOR_ROLES_MIGRATION_TEST_DBROLE_CHILD" {
  database = "MIGRATION_TEST_DB_FOR_ROLES"
  name = "MIGRATION_TEST_DBROLE_CHILD"
  comment = "Child database role for hierarchy testing"
}

resource "snowflake_database_role" "snowflake_generated_database_role_MIGRATION_TEST_DB_FOR_ROLES_MIGRATION_TEST_DBROLE_COMMENT" {
  database = "MIGRATION_TEST_DB_FOR_ROLES"
  name = "MIGRATION_TEST_DBROLE_COMMENT"
  comment = "This database role is used for migration testing purposes"
}

resource "snowflake_database_role" "snowflake_generated_database_role_MIGRATION_TEST_DB_FOR_ROLES_MIGRATION_TEST_DBROLE_LONG_COMMENT" {
  database = "MIGRATION_TEST_DB_FOR_ROLES"
  name = "MIGRATION_TEST_DBROLE_LONG_COMMENT"
  comment = "This is a very long comment that tests how the migration script handles comments with many characters. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."
}

resource "snowflake_database_role" "snowflake_generated_database_role_MIGRATION_TEST_DB_FOR_ROLES_MIGRATION_TEST_DBROLE_PARENT" {
  database = "MIGRATION_TEST_DB_FOR_ROLES"
  name = "MIGRATION_TEST_DBROLE_PARENT"
  comment = "Parent database role for hierarchy testing"
}

resource "snowflake_database_role" "snowflake_generated_database_role_MIGRATION_TEST_DB_FOR_ROLES_MIGRATION_TEST_DBROLE_SPECIAL" {
  database = "MIGRATION_TEST_DB_FOR_ROLES"
  name = "MIGRATION_TEST_DBROLE_SPECIAL"
  comment = "Comment with special chars: \u003c\u003e\u0026\""
}

resource "snowflake_database_role" "snowflake_generated_database_role_MIGRATION_TEST_DB_FOR_ROLES_MIGRATION_TEST_DBROLE_WITH_UNDERSCORE" {
  database = "MIGRATION_TEST_DB_FOR_ROLES"
  name = "MIGRATION_TEST_DBROLE_WITH_UNDERSCORE"
  comment = "Role with underscores in name"
}

import {
  to = snowflake_database_role.snowflake_generated_database_role_MIGRATION_TEST_DB_FOR_ROLES_MIGRATION_TEST_DBROLE_123
  id = "\"MIGRATION_TEST_DB_FOR_ROLES\".\"MIGRATION_TEST_DBROLE_123\""
}
import {
  to = snowflake_database_role.snowflake_generated_database_role_MIGRATION_TEST_DB_FOR_ROLES_MIGRATION_TEST_DBROLE_BASIC
  id = "\"MIGRATION_TEST_DB_FOR_ROLES\".\"MIGRATION_TEST_DBROLE_BASIC\""
}
import {
  to = snowflake_database_role.snowflake_generated_database_role_MIGRATION_TEST_DB_FOR_ROLES_MIGRATION_TEST_DBROLE_CHILD
  id = "\"MIGRATION_TEST_DB_FOR_ROLES\".\"MIGRATION_TEST_DBROLE_CHILD\""
}
import {
  to = snowflake_database_role.snowflake_generated_database_role_MIGRATION_TEST_DB_FOR_ROLES_MIGRATION_TEST_DBROLE_COMMENT
  id = "\"MIGRATION_TEST_DB_FOR_ROLES\".\"MIGRATION_TEST_DBROLE_COMMENT\""
}
import {
  to = snowflake_database_role.snowflake_generated_database_role_MIGRATION_TEST_DB_FOR_ROLES_MIGRATION_TEST_DBROLE_LONG_COMMENT
  id = "\"MIGRATION_TEST_DB_FOR_ROLES\".\"MIGRATION_TEST_DBROLE_LONG_COMMENT\""
}
import {
  to = snowflake_database_role.snowflake_generated_database_role_MIGRATION_TEST_DB_FOR_ROLES_MIGRATION_TEST_DBROLE_PARENT
  id = "\"MIGRATION_TEST_DB_FOR_ROLES\".\"MIGRATION_TEST_DBROLE_PARENT\""
}
import {
  to = snowflake_database_role.snowflake_generated_database_role_MIGRATION_TEST_DB_FOR_ROLES_MIGRATION_TEST_DBROLE_SPECIAL
  id = "\"MIGRATION_TEST_DB_FOR_ROLES\".\"MIGRATION_TEST_DBROLE_SPECIAL\""
}
import {
  to = snowflake_database_role.snowflake_generated_database_role_MIGRATION_TEST_DB_FOR_ROLES_MIGRATION_TEST_DBROLE_WITH_UNDERSCORE
  id = "\"MIGRATION_TEST_DB_FOR_ROLES\".\"MIGRATION_TEST_DBROLE_WITH_UNDERSCORE\""
}

