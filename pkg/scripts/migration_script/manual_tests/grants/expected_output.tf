# Beware of non-deterministic ordering which might cause diffs.

resource "snowflake_grant_privileges_to_account_role" "snowflake_generated_grant_on_DATABASE_MIGRATION_TEST_GRANT_DB_to_MIGRATION_TEST_GRANT_PRIV_ROLE_without_grant_option" {
  account_role_name = "MIGRATION_TEST_GRANT_PRIV_ROLE"
  on_account_object {
    object_name = "MIGRATION_TEST_GRANT_DB"
    object_type = "DATABASE"
  }
  privileges = ["CREATE SCHEMA", "USAGE"]
  with_grant_option = false
}

resource "snowflake_grant_privileges_to_database_role" "snowflake_generated_grant_on_TABLE_MIGRATION_TEST_GRANT_DB_MIGRATION_TEST_GRANT_SCHEMA_MIGRATION_TEST_GRANT_TABLE_to_MIGRATION_TEST_GRANT_DB_MIGRATION_TEST_GRANT_PRIV_DBROLE_with_grant_option" {
  database_role_name = "\"MIGRATION_TEST_GRANT_DB\".\"MIGRATION_TEST_GRANT_PRIV_DBROLE\""
  on_schema_object {
    object_name = "\"MIGRATION_TEST_GRANT_DB\".\"MIGRATION_TEST_GRANT_SCHEMA\".\"MIGRATION_TEST_GRANT_TABLE\""
    object_type = "TABLE"
  }
  privileges = ["INSERT"]
  with_grant_option = true
}

resource "snowflake_grant_privileges_to_account_role" "snowflake_generated_grant_on_WAREHOUSE_MIGRATION_TEST_GRANT_WH_to_MIGRATION_TEST_GRANT_PRIV_ROLE_without_grant_option" {
  account_role_name = "MIGRATION_TEST_GRANT_PRIV_ROLE"
  on_account_object {
    object_name = "MIGRATION_TEST_GRANT_WH"
    object_type = "WAREHOUSE"
  }
  privileges = ["OPERATE", "USAGE"]
  with_grant_option = false
}

resource "snowflake_grant_database_role" "snowflake_generated_grant_MIGRATION_TEST_GRANT_DB_MIGRATION_TEST_GRANT_CHILD_DBROLE_to_database_role_MIGRATION_TEST_GRANT_DB_MIGRATION_TEST_GRANT_PARENT_DBROLE" {
  database_role_name = "\"MIGRATION_TEST_GRANT_DB\".\"MIGRATION_TEST_GRANT_CHILD_DBROLE\""
  parent_database_role_name = "\"MIGRATION_TEST_GRANT_DB\".\"MIGRATION_TEST_GRANT_PARENT_DBROLE\""
}

resource "snowflake_grant_privileges_to_database_role" "snowflake_generated_grant_on_database_MIGRATION_TEST_GRANT_DB_to_MIGRATION_TEST_GRANT_DB_MIGRATION_TEST_GRANT_PRIV_DBROLE_without_grant_option" {
  database_role_name = "\"MIGRATION_TEST_GRANT_DB\".\"MIGRATION_TEST_GRANT_PRIV_DBROLE\""
  on_database = "MIGRATION_TEST_GRANT_DB"
  privileges = ["USAGE"]
  with_grant_option = false
}

resource "snowflake_grant_database_role" "snowflake_generated_grant_MIGRATION_TEST_GRANT_DB_MIGRATION_TEST_GRANT_PRIV_DBROLE_to_role_MIGRATION_TEST_GRANT_PRIV_ROLE" {
  database_role_name = "\"MIGRATION_TEST_GRANT_DB\".\"MIGRATION_TEST_GRANT_PRIV_DBROLE\""
  parent_role_name = "MIGRATION_TEST_GRANT_PRIV_ROLE"
}

resource "snowflake_grant_privileges_to_account_role" "snowflake_generated_grant_on_account_to_MIGRATION_TEST_GRANT_PRIV_ROLE_without_grant_option" {
  account_role_name = "MIGRATION_TEST_GRANT_PRIV_ROLE"
  on_account = true
  privileges = ["CREATE DATABASE", "CREATE WAREHOUSE"]
  with_grant_option = false
}

resource "snowflake_grant_privileges_to_account_role" "snowflake_generated_grant_on_TABLE_MIGRATION_TEST_GRANT_DB_MIGRATION_TEST_GRANT_SCHEMA_MIGRATION_TEST_GRANT_TABLE_to_MIGRATION_TEST_GRANT_PRIV_ROLE_without_grant_option" {
  account_role_name = "MIGRATION_TEST_GRANT_PRIV_ROLE"
  on_schema_object {
    object_name = "\"MIGRATION_TEST_GRANT_DB\".\"MIGRATION_TEST_GRANT_SCHEMA\".\"MIGRATION_TEST_GRANT_TABLE\""
    object_type = "TABLE"
  }
  privileges = ["INSERT", "SELECT", "UPDATE"]
  with_grant_option = false
}

resource "snowflake_grant_privileges_to_database_role" "snowflake_generated_grant_on_TABLE_MIGRATION_TEST_GRANT_DB_MIGRATION_TEST_GRANT_SCHEMA_MIGRATION_TEST_GRANT_TABLE_to_MIGRATION_TEST_GRANT_DB_MIGRATION_TEST_GRANT_PRIV_DBROLE_without_grant_option" {
  database_role_name = "\"MIGRATION_TEST_GRANT_DB\".\"MIGRATION_TEST_GRANT_PRIV_DBROLE\""
  on_schema_object {
    object_name = "\"MIGRATION_TEST_GRANT_DB\".\"MIGRATION_TEST_GRANT_SCHEMA\".\"MIGRATION_TEST_GRANT_TABLE\""
    object_type = "TABLE"
  }
  privileges = ["SELECT"]
  with_grant_option = false
}

resource "snowflake_grant_privileges_to_database_role" "snowflake_generated_grant_on_database_MIGRATION_TEST_GRANT_DB_to_MIGRATION_TEST_GRANT_DB_MIGRATION_TEST_GRANT_PARENT_DBROLE_without_grant_option" {
  database_role_name = "\"MIGRATION_TEST_GRANT_DB\".\"MIGRATION_TEST_GRANT_PARENT_DBROLE\""
  on_database = "MIGRATION_TEST_GRANT_DB"
  privileges = ["USAGE"]
  with_grant_option = false
}

resource "snowflake_grant_privileges_to_database_role" "snowflake_generated_grant_on_schema_MIGRATION_TEST_GRANT_DB_MIGRATION_TEST_GRANT_SCHEMA_to_MIGRATION_TEST_GRANT_DB_MIGRATION_TEST_GRANT_PRIV_DBROLE_without_grant_option" {
  database_role_name = "\"MIGRATION_TEST_GRANT_DB\".\"MIGRATION_TEST_GRANT_PRIV_DBROLE\""
  on_schema {
    schema_name = "\"MIGRATION_TEST_GRANT_DB\".\"MIGRATION_TEST_GRANT_SCHEMA\""
  }
  privileges = ["CREATE TABLE", "USAGE"]
  with_grant_option = false
}

resource "snowflake_grant_privileges_to_account_role" "snowflake_generated_grant_on_schema_MIGRATION_TEST_GRANT_DB_MIGRATION_TEST_GRANT_SCHEMA_to_MIGRATION_TEST_GRANT_PRIV_ROLE_without_grant_option" {
  account_role_name = "MIGRATION_TEST_GRANT_PRIV_ROLE"
  on_schema {
    schema_name = "\"MIGRATION_TEST_GRANT_DB\".\"MIGRATION_TEST_GRANT_SCHEMA\""
  }
  privileges = ["CREATE TABLE", "CREATE VIEW", "USAGE"]
  with_grant_option = false
}

resource "snowflake_grant_privileges_to_account_role" "snowflake_generated_grant_on_TABLE_MIGRATION_TEST_GRANT_DB_MIGRATION_TEST_GRANT_SCHEMA_MIGRATION_TEST_GRANT_TABLE_to_MIGRATION_TEST_GRANT_PRIV_ROLE_with_grant_option" {
  account_role_name = "MIGRATION_TEST_GRANT_PRIV_ROLE"
  on_schema_object {
    object_name = "\"MIGRATION_TEST_GRANT_DB\".\"MIGRATION_TEST_GRANT_SCHEMA\".\"MIGRATION_TEST_GRANT_TABLE\""
    object_type = "TABLE"
  }
  privileges = ["DELETE"]
  with_grant_option = true
}
import {
  to = snowflake_grant_privileges_to_account_role.snowflake_generated_grant_on_DATABASE_MIGRATION_TEST_GRANT_DB_to_MIGRATION_TEST_GRANT_PRIV_ROLE_without_grant_option
  id = "\"MIGRATION_TEST_GRANT_PRIV_ROLE\"|false|false|CREATE SCHEMA,USAGE|OnAccountObject|DATABASE|\"MIGRATION_TEST_GRANT_DB\""
}
import {
  to = snowflake_grant_privileges_to_database_role.snowflake_generated_grant_on_TABLE_MIGRATION_TEST_GRANT_DB_MIGRATION_TEST_GRANT_SCHEMA_MIGRATION_TEST_GRANT_TABLE_to_MIGRATION_TEST_GRANT_DB_MIGRATION_TEST_GRANT_PRIV_DBROLE_with_grant_option
  id = "\"MIGRATION_TEST_GRANT_DB\".\"MIGRATION_TEST_GRANT_PRIV_DBROLE\"|true|false|INSERT|OnSchemaObject|OnObject|TABLE|\"MIGRATION_TEST_GRANT_DB\".\"MIGRATION_TEST_GRANT_SCHEMA\".\"MIGRATION_TEST_GRANT_TABLE\""
}
import {
  to = snowflake_grant_privileges_to_account_role.snowflake_generated_grant_on_WAREHOUSE_MIGRATION_TEST_GRANT_WH_to_MIGRATION_TEST_GRANT_PRIV_ROLE_without_grant_option
  id = "\"MIGRATION_TEST_GRANT_PRIV_ROLE\"|false|false|OPERATE,USAGE|OnAccountObject|WAREHOUSE|\"MIGRATION_TEST_GRANT_WH\""
}
import {
  to = snowflake_grant_database_role.snowflake_generated_grant_MIGRATION_TEST_GRANT_DB_MIGRATION_TEST_GRANT_CHILD_DBROLE_to_database_role_MIGRATION_TEST_GRANT_DB_MIGRATION_TEST_GRANT_PARENT_DBROLE
  id = "\"MIGRATION_TEST_GRANT_DB\".\"MIGRATION_TEST_GRANT_CHILD_DBROLE\"|DATABASE ROLE|\"MIGRATION_TEST_GRANT_DB\".\"MIGRATION_TEST_GRANT_PARENT_DBROLE\""
}
import {
  to = snowflake_grant_privileges_to_database_role.snowflake_generated_grant_on_database_MIGRATION_TEST_GRANT_DB_to_MIGRATION_TEST_GRANT_DB_MIGRATION_TEST_GRANT_PRIV_DBROLE_without_grant_option
  id = "\"MIGRATION_TEST_GRANT_DB\".\"MIGRATION_TEST_GRANT_PRIV_DBROLE\"|false|false|USAGE|OnDatabase|\"MIGRATION_TEST_GRANT_DB\""
}
import {
  to = snowflake_grant_database_role.snowflake_generated_grant_MIGRATION_TEST_GRANT_DB_MIGRATION_TEST_GRANT_PRIV_DBROLE_to_role_MIGRATION_TEST_GRANT_PRIV_ROLE
  id = "\"MIGRATION_TEST_GRANT_DB\".\"MIGRATION_TEST_GRANT_PRIV_DBROLE\"|ROLE|\"MIGRATION_TEST_GRANT_PRIV_ROLE\""
}
import {
  to = snowflake_grant_privileges_to_account_role.snowflake_generated_grant_on_account_to_MIGRATION_TEST_GRANT_PRIV_ROLE_without_grant_option
  id = "\"MIGRATION_TEST_GRANT_PRIV_ROLE\"|false|false|CREATE DATABASE,CREATE WAREHOUSE|OnAccount"
}
import {
  to = snowflake_grant_privileges_to_account_role.snowflake_generated_grant_on_TABLE_MIGRATION_TEST_GRANT_DB_MIGRATION_TEST_GRANT_SCHEMA_MIGRATION_TEST_GRANT_TABLE_to_MIGRATION_TEST_GRANT_PRIV_ROLE_without_grant_option
  id = "\"MIGRATION_TEST_GRANT_PRIV_ROLE\"|false|false|INSERT,SELECT,UPDATE|OnSchemaObject|OnObject|TABLE|\"MIGRATION_TEST_GRANT_DB\".\"MIGRATION_TEST_GRANT_SCHEMA\".\"MIGRATION_TEST_GRANT_TABLE\""
}
import {
  to = snowflake_grant_privileges_to_database_role.snowflake_generated_grant_on_TABLE_MIGRATION_TEST_GRANT_DB_MIGRATION_TEST_GRANT_SCHEMA_MIGRATION_TEST_GRANT_TABLE_to_MIGRATION_TEST_GRANT_DB_MIGRATION_TEST_GRANT_PRIV_DBROLE_without_grant_option
  id = "\"MIGRATION_TEST_GRANT_DB\".\"MIGRATION_TEST_GRANT_PRIV_DBROLE\"|false|false|SELECT|OnSchemaObject|OnObject|TABLE|\"MIGRATION_TEST_GRANT_DB\".\"MIGRATION_TEST_GRANT_SCHEMA\".\"MIGRATION_TEST_GRANT_TABLE\""
}
import {
  to = snowflake_grant_privileges_to_database_role.snowflake_generated_grant_on_database_MIGRATION_TEST_GRANT_DB_to_MIGRATION_TEST_GRANT_DB_MIGRATION_TEST_GRANT_PARENT_DBROLE_without_grant_option
  id = "\"MIGRATION_TEST_GRANT_DB\".\"MIGRATION_TEST_GRANT_PARENT_DBROLE\"|false|false|USAGE|OnDatabase|\"MIGRATION_TEST_GRANT_DB\""
}
import {
  to = snowflake_grant_privileges_to_database_role.snowflake_generated_grant_on_schema_MIGRATION_TEST_GRANT_DB_MIGRATION_TEST_GRANT_SCHEMA_to_MIGRATION_TEST_GRANT_DB_MIGRATION_TEST_GRANT_PRIV_DBROLE_without_grant_option
  id = "\"MIGRATION_TEST_GRANT_DB\".\"MIGRATION_TEST_GRANT_PRIV_DBROLE\"|false|false|CREATE TABLE,USAGE|OnSchema|OnSchema|\"MIGRATION_TEST_GRANT_DB\".\"MIGRATION_TEST_GRANT_SCHEMA\""
}
import {
  to = snowflake_grant_privileges_to_account_role.snowflake_generated_grant_on_schema_MIGRATION_TEST_GRANT_DB_MIGRATION_TEST_GRANT_SCHEMA_to_MIGRATION_TEST_GRANT_PRIV_ROLE_without_grant_option
  id = "\"MIGRATION_TEST_GRANT_PRIV_ROLE\"|false|false|CREATE TABLE,CREATE VIEW,USAGE|OnSchema|OnSchema|\"MIGRATION_TEST_GRANT_DB\".\"MIGRATION_TEST_GRANT_SCHEMA\""
}
import {
  to = snowflake_grant_privileges_to_account_role.snowflake_generated_grant_on_TABLE_MIGRATION_TEST_GRANT_DB_MIGRATION_TEST_GRANT_SCHEMA_MIGRATION_TEST_GRANT_TABLE_to_MIGRATION_TEST_GRANT_PRIV_ROLE_with_grant_option
  id = "\"MIGRATION_TEST_GRANT_PRIV_ROLE\"|true|false|DELETE|OnSchemaObject|OnObject|TABLE|\"MIGRATION_TEST_GRANT_DB\".\"MIGRATION_TEST_GRANT_SCHEMA\".\"MIGRATION_TEST_GRANT_TABLE\""
}
