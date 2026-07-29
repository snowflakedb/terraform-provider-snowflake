resource "snowflake_database" "db" {
  name = "database"
}

resource "snowflake_schema" "my_schema" {
  database = snowflake_database.db.name
  name     = "my_schema"
}

resource "snowflake_database_role" "db_role" {
  database = snowflake_database.db.name
  name     = "db_role_name"
}

##################################
### on database privileges
##################################

# list of privileges
resource "snowflake_grant_privileges_to_database_role" "example" {
  privileges         = ["CREATE", "MONITOR"]
  database_role_name = snowflake_database_role.db_role.fully_qualified_name
  on_database        = snowflake_database_role.db_role.database
}

# all privileges + grant option
resource "snowflake_grant_privileges_to_database_role" "example" {
  database_role_name = snowflake_database_role.db_role.fully_qualified_name
  on_database        = snowflake_database_role.db_role.database
  all_privileges     = true
  with_grant_option  = true
}

# all privileges + grant option + always apply
resource "snowflake_grant_privileges_to_database_role" "example" {
  database_role_name = snowflake_database_role.db_role.fully_qualified_name
  on_database        = snowflake_database_role.db_role.database
  always_apply       = true
  all_privileges     = true
  with_grant_option  = true
}

##################################
### schema privileges
##################################

# list of privileges
resource "snowflake_grant_privileges_to_database_role" "example" {
  privileges         = ["MODIFY", "CREATE TABLE"]
  database_role_name = snowflake_database_role.db_role.fully_qualified_name
  on_schema {
    schema_name = snowflake_schema.my_schema.fully_qualified_name # note this is a fully qualified name!
  }
}

# all privileges + grant option
resource "snowflake_grant_privileges_to_database_role" "example" {
  database_role_name = snowflake_database_role.db_role.fully_qualified_name
  on_schema {
    schema_name = snowflake_schema.my_schema.fully_qualified_name # note this is a fully qualified name!
  }
  all_privileges    = true
  with_grant_option = true
}

# all schemas in database
resource "snowflake_grant_privileges_to_database_role" "example" {
  privileges         = ["MODIFY", "CREATE TABLE"]
  database_role_name = snowflake_database_role.db_role.fully_qualified_name
  on_schema {
    all_schemas_in_database = snowflake_database_role.db_role.database
  }
}

# future schemas in database
resource "snowflake_grant_privileges_to_database_role" "example" {
  privileges         = ["MODIFY", "CREATE TABLE"]
  database_role_name = snowflake_database_role.db_role.fully_qualified_name
  on_schema {
    future_schemas_in_database = snowflake_database_role.db_role.database
  }
}

##################################
### schema object privileges
##################################

# list of privileges
resource "snowflake_grant_privileges_to_database_role" "example" {
  privileges         = ["SELECT", "REFERENCES"]
  database_role_name = snowflake_database_role.db_role.fully_qualified_name
  on_schema_object {
    object_type = "VIEW"
    object_name = snowflake_view.my_view.fully_qualified_name # note this is a fully qualified name!
  }
}

# all privileges + grant option
resource "snowflake_grant_privileges_to_database_role" "example" {
  database_role_name = snowflake_database_role.db_role.fully_qualified_name
  on_schema_object {
    object_type = "VIEW"
    object_name = snowflake_view.my_view.fully_qualified_name # note this is a fully qualified name!
  }
  all_privileges    = true
  with_grant_option = true
}

# all in database
resource "snowflake_grant_privileges_to_database_role" "example" {
  privileges         = ["SELECT", "INSERT"]
  database_role_name = snowflake_database_role.db_role.fully_qualified_name
  on_schema_object {
    all {
      object_type_plural = "TABLES"
      in_database        = snowflake_database_role.db_role.database
    }
  }
}

# all in schema
resource "snowflake_grant_privileges_to_database_role" "example" {
  privileges         = ["SELECT", "INSERT"]
  database_role_name = snowflake_database_role.db_role.fully_qualified_name
  on_schema_object {
    all {
      object_type_plural = "TABLES"
      in_schema          = snowflake_schema.my_schema.fully_qualified_name # note this is a fully qualified name!
    }
  }
}

# future in database
resource "snowflake_grant_privileges_to_database_role" "example" {
  privileges         = ["SELECT", "INSERT"]
  database_role_name = snowflake_database_role.db_role.fully_qualified_name
  on_schema_object {
    future {
      object_type_plural = "TABLES"
      in_database        = snowflake_database_role.db_role.database
    }
  }
}

# future in schema
resource "snowflake_grant_privileges_to_database_role" "example" {
  privileges         = ["SELECT", "INSERT"]
  database_role_name = snowflake_database_role.db_role.fully_qualified_name
  on_schema_object {
    future {
      object_type_plural = "TABLES"
      in_schema          = snowflake_schema.my_schema.fully_qualified_name # note this is a fully qualified name!
    }
  }
}
