# Simple usage
data "snowflake_database_roles" "simple" {
  in_database = "database-name"
}

output "simple_output" {
  value = data.snowflake_database_roles.simple.database_roles
}

# Filtering (like)
data "snowflake_database_roles" "like" {
  in_database = "database-name"
  like        = "database_role-name"
}

output "like_output" {
  value = data.snowflake_database_roles.like.database_roles
}

# Filtering (limit)
data "snowflake_database_roles" "limit" {
  in_database = "database-name"
  limit {
    rows = 10
    from = "prefix-"
  }
}

output "limit_output" {
  value = data.snowflake_database_roles.limit.database_roles
}

# Ensure the number of database roles is equal to at least one element (with the use of postcondition)
data "snowflake_database_roles" "assert_with_postcondition" {
  in_database = "database-name"
  like        = "database_role-name-%"
  lifecycle {
    postcondition {
      condition     = length(self.database_roles) > 0
      error_message = "there should be at least one database role"
    }
  }
}

# Ensure the number of database roles is equal to exactly one element (with the use of check block)
check "database_role_check" {
  data "snowflake_database_roles" "assert_with_check_block" {
    in_database = "database-name"
    like        = "database_role-name"
  }

  assert {
    condition     = length(data.snowflake_database_roles.assert_with_check_block.database_roles) == 1
    error_message = "Database roles filtered by '${data.snowflake_database_roles.assert_with_check_block.like}' returned ${length(data.snowflake_database_roles.assert_with_check_block.database_roles)} database roles where one was expected"
  }
}
