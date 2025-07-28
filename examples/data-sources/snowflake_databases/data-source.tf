# Simple usage
data "snowflake_databases" "simple" {
}

output "simple_output" {
  value = data.snowflake_databases.simple.databases
}

# Filtering (like)
data "snowflake_databases" "like" {
  like = "database-name"
}

output "like_output" {
  value = data.snowflake_databases.like.databases
}

# Filtering (starts_with)
data "snowflake_databases" "starts_with" {
  starts_with = "database-"
}

output "starts_with_output" {
  value = data.snowflake_databases.starts_with.databases
}

# Filtering (limit)
data "snowflake_databases" "limit" {
  limit {
    rows = 10
    from = "database-"
  }
}

output "limit_output" {
  value = data.snowflake_databases.limit.databases
}

# Without additional data (to limit the number of calls make for every found database)
data "snowflake_databases" "only_show" {
  # with_describe is turned on by default and it calls DESCRIBE DATABASE for every database found and attaches its output to databases.*.describe_output field
  with_describe = false

  # with_parameters is turned on by default and it calls SHOW PARAMETERS FOR DATABASE for every database found and attaches its output to databases.*.parameters field
  with_parameters = false
}

output "only_show_output" {
  value = data.snowflake_databases.only_show.databases
}

# Ensure the number of databases is equal to at least one element (with the use of postcondition)
data "snowflake_databases" "assert_with_postcondition" {
  starts_with = "database-name"
  lifecycle {
    postcondition {
      condition     = length(self.databases) > 0
      error_message = "there should be at least one database"
    }
  }
}

# Ensure the number of databases is equal to exactly one element (with the use of check block)
check "database_check" {
  data "snowflake_databases" "assert_with_check_block" {
    like = "database-name"
  }

  assert {
    condition     = length(data.snowflake_databases.assert_with_check_block.databases) == 1
    error_message = "Databases filtered by '${data.snowflake_databases.assert_with_check_block.like}' returned ${length(data.snowflake_databases.assert_with_check_block.databases)} databases where one was expected"
  }
}
