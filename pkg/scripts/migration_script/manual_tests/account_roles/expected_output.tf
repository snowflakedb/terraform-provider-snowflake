# =============================================================================
# Expected Migration Script Output - Account Roles
# =============================================================================
# This file contains the EXPECTED output from running:
#   go run .. -import=block account_roles < objects.csv
#
# Use this to compare with actual output:
#   go run .. -import=block account_roles < objects.csv > actual_output.tf
#   diff expected_output.tf actual_output.tf
#
# NOTE:
# - Resources are sorted alphabetically by name
# - Import blocks are grouped at the end
# - Special characters are Unicode-escaped in HCL output
# =============================================================================

resource "snowflake_account_role" "snowflake_generated_account_role_MIGRATION_TEST_ROLE_BASIC" {
  name = "MIGRATION_TEST_ROLE_BASIC"
}

resource "snowflake_account_role" "snowflake_generated_account_role_MIGRATION_TEST_ROLE_CHILD" {
  name = "MIGRATION_TEST_ROLE_CHILD"
  comment = "Child role for hierarchy testing"
}

resource "snowflake_account_role" "snowflake_generated_account_role_MIGRATION_TEST_ROLE_COMMENT" {
  name = "MIGRATION_TEST_ROLE_COMMENT"
  comment = "This role is used for migration testing purposes"
}

resource "snowflake_account_role" "snowflake_generated_account_role_MIGRATION_TEST_ROLE_LONG_COMMENT" {
  name = "MIGRATION_TEST_ROLE_LONG_COMMENT"
  comment = "This is a very long comment that tests how the migration script handles comments with many characters. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."
}

resource "snowflake_account_role" "snowflake_generated_account_role_MIGRATION_TEST_ROLE_PARENT" {
  name = "MIGRATION_TEST_ROLE_PARENT"
  comment = "Parent role for hierarchy testing"
}

resource "snowflake_account_role" "snowflake_generated_account_role_MIGRATION_TEST_ROLE_SPECIAL" {
  name = "MIGRATION_TEST_ROLE_SPECIAL"
  comment = "Comment with special chars: \u003c\u003e\u0026\""
}

import {
  to = snowflake_account_role.snowflake_generated_account_role_MIGRATION_TEST_ROLE_BASIC
  id = "\"MIGRATION_TEST_ROLE_BASIC\""
}
import {
  to = snowflake_account_role.snowflake_generated_account_role_MIGRATION_TEST_ROLE_CHILD
  id = "\"MIGRATION_TEST_ROLE_CHILD\""
}
import {
  to = snowflake_account_role.snowflake_generated_account_role_MIGRATION_TEST_ROLE_COMMENT
  id = "\"MIGRATION_TEST_ROLE_COMMENT\""
}
import {
  to = snowflake_account_role.snowflake_generated_account_role_MIGRATION_TEST_ROLE_LONG_COMMENT
  id = "\"MIGRATION_TEST_ROLE_LONG_COMMENT\""
}
import {
  to = snowflake_account_role.snowflake_generated_account_role_MIGRATION_TEST_ROLE_PARENT
  id = "\"MIGRATION_TEST_ROLE_PARENT\""
}
import {
  to = snowflake_account_role.snowflake_generated_account_role_MIGRATION_TEST_ROLE_SPECIAL
  id = "\"MIGRATION_TEST_ROLE_SPECIAL\""
}

