resource "snowflake_database" "db" {
  name = "database"
}

resource "snowflake_schema" "my_schema" {
  database = snowflake_database.db.name
  name     = "my_schema"
}

resource "snowflake_account_role" "db_role" {
  name = "role_name"
}

##################################
### on account privileges
##################################

# list of privileges
resource "snowflake_grant_privileges_to_account_role" "example" {
  privileges        = ["CREATE DATABASE", "CREATE USER"]
  account_role_name = snowflake_account_role.db_role.name
  on_account        = true
}

## ID: "\"role_name\"|false|false|CREATE DATABASE,CREATE USER|OnAccount"

# all privileges + grant option
resource "snowflake_grant_privileges_to_account_role" "example" {
  account_role_name = snowflake_account_role.db_role.name
  on_account        = true
  all_privileges    = true
  with_grant_option = true
}

## ID: "\"role_name\"|true|false|ALL|OnAccount"

# all privileges + grant option + always apply
resource "snowflake_grant_privileges_to_account_role" "example" {
  account_role_name = snowflake_account_role.db_role.name
  on_account        = true
  always_apply      = true
  all_privileges    = true
  with_grant_option = true
}

## ID: "\"role_name\"|true|true|ALL|OnAccount"

##################################
### on account object privileges
##################################

# list of privileges
resource "snowflake_grant_privileges_to_account_role" "example" {
  privileges        = ["CREATE SCHEMA", "CREATE DATABASE ROLE"]
  account_role_name = snowflake_account_role.db_role.name
  on_account_object {
    object_type = "DATABASE"
    object_name = snowflake_database.db.name
  }
}

## ID: "\"role_name\"|false|false|CREATE SCHEMA,CREATE DATABASE ROLE|OnAccountObject|DATABASE|\"database\""

# all privileges + grant option
resource "snowflake_grant_privileges_to_account_role" "example" {
  account_role_name = snowflake_account_role.db_role.name
  on_account_object {
    object_type = "DATABASE"
    object_name = snowflake_database.db.name
  }
  all_privileges    = true
  with_grant_option = true
}

## ID: "\"role_name\"|true|false|ALL|OnAccountObject|DATABASE|\"database\""

# grant IMPORTED PRIVILEGES on SNOWFLAKE application
resource "snowflake_grant_privileges_to_account_role" "example" {
  account_role_name = snowflake_account_role.db_role.name
  privileges        = ["IMPORTED PRIVILEGES"]
  on_account_object {
    object_type = "DATABASE" # All applications should be using DATABASE object_type
    object_name = "SNOWFLAKE"
  }
}

## ID: "\"role_name\"|false|false|IMPORTED PRIVILEGES|OnAccountObject|DATABASE|\"SNOWFLAKE\""

# all privileges + grant option + always apply
resource "snowflake_grant_privileges_to_account_role" "example" {
  account_role_name = snowflake_account_role.db_role.name
  on_account_object {
    object_type = "DATABASE"
    object_name = snowflake_database.db.name
  }
  always_apply      = true
  all_privileges    = true
  with_grant_option = true
}

## ID: "\"role_name\"|true|true|ALL|OnAccountObject|DATABASE|\"database\""

##################################
### schema privileges
##################################

# list of privileges
resource "snowflake_grant_privileges_to_account_role" "example" {
  privileges        = ["MODIFY", "CREATE TABLE"]
  account_role_name = snowflake_account_role.db_role.name
  on_schema {
    schema_name = snowflake_schema.my_schema.fully_qualified_name # note this is a fully qualified name!
  }
}

## ID: "\"role_name\"|false|false|MODIFY,CREATE TABLE|OnSchema|OnSchema|\"database\".\"my_schema\""

# all privileges + grant option
resource "snowflake_grant_privileges_to_account_role" "example" {
  account_role_name = snowflake_account_role.db_role.name
  on_schema {
    schema_name = snowflake_schema.my_schema.fully_qualified_name # note this is a fully qualified name!
  }
  all_privileges    = true
  with_grant_option = true
}

## ID: "\"role_name\"|true|false|MODIFY,CREATE TABLE|OnSchema|OnSchema|\"database\".\"my_schema\""

# all schemas in database
resource "snowflake_grant_privileges_to_account_role" "example" {
  privileges        = ["MODIFY", "CREATE TABLE"]
  account_role_name = snowflake_account_role.db_role.name
  on_schema {
    all_schemas_in_database = snowflake_database.db.name
  }
}

## ID: "\"role_name\"|false|false|MODIFY,CREATE TABLE|OnSchema|OnAllSchemasInDatabase|\"database\""

# future schemas in database
resource "snowflake_grant_privileges_to_account_role" "example" {
  privileges        = ["MODIFY", "CREATE TABLE"]
  account_role_name = snowflake_account_role.db_role.name
  on_schema {
    future_schemas_in_database = snowflake_database.db.name
  }
}

## ID: "\"role_name\"|false|false|MODIFY,CREATE TABLE|OnSchema|OnFutureSchemasInDatabase|\"database\""

##################################
### schema object privileges
##################################

# list of privileges
resource "snowflake_grant_privileges_to_account_role" "example" {
  privileges        = ["SELECT", "REFERENCES"]
  account_role_name = snowflake_account_role.db_role.name
  on_schema_object {
    object_type = "VIEW"
    object_name = snowflake_view.my_view.fully_qualified_name # note this is a fully qualified name!
  }
}

## ID: "\"role_name\"|false|false|SELECT,REFERENCES|OnSchemaObject|VIEW|\"database\".\"my_schema\".\"my_view\""

# all privileges + grant option
resource "snowflake_grant_privileges_to_account_role" "example" {
  account_role_name = snowflake_account_role.db_role.name
  on_schema_object {
    object_type = "VIEW"
    object_name = snowflake_view.my_view.fully_qualified_name # note this is a fully qualified name!
  }
  all_privileges    = true
  with_grant_option = true
}

## ID: "\"role_name\"|true|false|ALL|OnSchemaObject|OnObject|VIEW|\"database\".\"my_schema\".\"my_view\""

# all in database
resource "snowflake_grant_privileges_to_account_role" "example" {
  privileges        = ["SELECT", "INSERT"]
  account_role_name = snowflake_account_role.db_role.name
  on_schema_object {
    all {
      object_type_plural = "TABLES"
      in_database        = snowflake_database.db.name
    }
  }
}

## ID: "\"role_name\"|false|false|SELECT,INSERT|OnSchemaObject|OnAll|TABLES|InDatabase|\"database\""

# all in schema
resource "snowflake_grant_privileges_to_account_role" "example" {
  privileges        = ["SELECT", "INSERT"]
  account_role_name = snowflake_account_role.db_role.name
  on_schema_object {
    all {
      object_type_plural = "TABLES"
      in_schema          = snowflake_schema.my_schema.fully_qualified_name # note this is a fully qualified name!
    }
  }
}

## ID: "\"role_name\"|false|false|SELECT,INSERT|OnSchemaObject|OnAll|TABLES|InSchema|\"database\".\"my_schema\""

# future in database
resource "snowflake_grant_privileges_to_account_role" "example" {
  privileges        = ["SELECT", "INSERT"]
  account_role_name = snowflake_account_role.db_role.name
  on_schema_object {
    future {
      object_type_plural = "TABLES"
      in_database        = snowflake_database.db.name
    }
  }
}

## ID: "\"role_name\"|false|false|SELECT,INSERT|OnSchemaObject|OnFuture|TABLES|InDatabase|\"database\""

# future in schema
resource "snowflake_grant_privileges_to_account_role" "example" {
  privileges        = ["SELECT", "INSERT"]
  account_role_name = snowflake_account_role.db_role.name
  on_schema_object {
    future {
      object_type_plural = "TABLES"
      in_schema          = snowflake_schema.my_schema.fully_qualified_name # note this is a fully qualified name!
    }
  }
}

## ID: "\"role_name\"|false|false|SELECT,INSERT|OnSchemaObject|OnFuture|TABLES|InSchema|\"database\".\"my_schema\""
