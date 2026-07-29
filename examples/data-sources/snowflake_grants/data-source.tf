##################################
### SHOW GRANTS ON ...
##################################

# account
data "snowflake_grants" "example_on_account" {
  grants_on {
    account = true
  }
}

# account object (e.g. database)
data "snowflake_grants" "example_on_account_object" {
  grants_on {
    object_name = "some_database"
    object_type = "DATABASE"
  }
}

# database object (e.g. schema)
data "snowflake_grants" "example_on_database_object" {
  grants_on {
    object_name = "\"some_database\".\"some_schema\""
    object_type = "SCHEMA"
  }
}

# schema object (e.g. table)
data "snowflake_grants" "example_on_schema_object" {
  grants_on {
    object_name = "\"some_database\".\"some_schema\".\"some_table\""
    object_type = "TABLE"
  }
}

##################################
### SHOW GRANTS TO ...
##################################

# application
data "snowflake_grants" "example_to_application" {
  grants_to {
    application = "some_application"
  }
}

# application role
data "snowflake_grants" "example_to_application_role" {
  grants_to {
    application_role = "\"some_application\".\"some_application_role\""
  }
}

# account role
data "snowflake_grants" "example_to_role" {
  grants_to {
    account_role = "some_role"
  }
}

# database role
data "snowflake_grants" "example_to_database_role" {
  grants_to {
    database_role = "\"some_database\".\"some_database_role\""
  }
}

# share
data "snowflake_grants" "example_to_share" {
  grants_to {
    share {
      share_name = "some_share"
    }
  }
}

# user
data "snowflake_grants" "example_to_user" {
  grants_to {
    user = "some_user"
  }
}

##################################
### SHOW GRANTS OF ...
##################################

# application role
data "snowflake_grants" "example_of_application_role" {
  grants_of {
    application_role = "\"some_application\".\"some_application_role\""
  }
}

# database role
data "snowflake_grants" "example_of_database_role" {
  grants_of {
    database_role = "\"some_database\".\"some_database_role\""
  }
}

# account role
data "snowflake_grants" "example_of_role" {
  grants_of {
    account_role = "some_role"
  }
}

# share
data "snowflake_grants" "example_of_share" {
  grants_of {
    share = "some_share"
  }
}

##################################
### SHOW FUTURE GRANTS IN ...
##################################

# database
data "snowflake_grants" "example_future_in_database" {
  future_grants_in {
    database = "some_database"
  }
}

# schema
data "snowflake_grants" "example_future_in_schema" {
  future_grants_in {
    schema = "\"some_database\".\"some_schema\""
  }
}

##################################
### SHOW FUTURE GRANTS TO ...
##################################

# account role
data "snowflake_grants" "example_future_to_role" {
  future_grants_to {
    account_role = "some_role"
  }
}

# database role
data "snowflake_grants" "example_future_to_database_role" {
  future_grants_to {
    database_role = "\"some_database\".\"some_database_role\""
  }
}
