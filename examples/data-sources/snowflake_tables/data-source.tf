# Simple usage
data "snowflake_tables" "simple" {
}

output "simple_output" {
  value = data.snowflake_tables.simple.tables
}

# Filtering (like)
data "snowflake_tables" "like" {
  like = "table-name"
}

output "like_output" {
  value = data.snowflake_tables.like.tables
}

# Filtering by prefix (like)
data "snowflake_tables" "like_prefix" {
  like = "prefix%"
}

output "like_prefix_output" {
  value = data.snowflake_tables.like_prefix.tables
}

# Filtering (starts_with)
data "snowflake_tables" "starts_with" {
  starts_with = "prefix-"
}

output "starts_with_output" {
  value = data.snowflake_tables.starts_with.tables
}

# Filtering (in)
data "snowflake_tables" "in_account" {
  in {
    account = true
  }
}

data "snowflake_tables" "in_database" {
  in {
    database = "<database_name>"
  }
}

data "snowflake_tables" "in_schema" {
  in {
    schema = "<database_name>.<schema_name>"
  }
}

data "snowflake_tables" "in_application" {
  in {
    application = "<application_name>"
  }
}

data "snowflake_tables" "in_application_package" {
  in {
    application_package = "<application_package_name>"
  }
}

output "in_output" {
  value = {
    "account" : data.snowflake_tables.in_account.tables,
    "database" : data.snowflake_tables.in_database.tables,
    "schema" : data.snowflake_tables.in_schema.tables,
    "application" : data.snowflake_tables.in_application.tables,
    "application_package" : data.snowflake_tables.in_application_package.tables,
  }
}

# Filtering (limit)
data "snowflake_tables" "limit" {
  limit {
    rows = 10
    from = "prefix-"
  }
}

output "limit_output" {
  value = data.snowflake_tables.limit.tables
}

# Without additional data (to limit the number of calls make for every found table)
data "snowflake_tables" "only_show" {
  # with_describe is turned on by default and it calls DESCRIBE TABLE for every table found and attaches its output to tables.*.describe_output field
  with_describe = false
}

output "only_show_output" {
  value = data.snowflake_tables.only_show.tables
}

# Ensure the number of tables is equal to at least one element (with the use of postcondition)
data "snowflake_tables" "assert_with_postcondition" {
  like = "table-name%"
  lifecycle {
    postcondition {
      condition     = length(self.tables) > 0
      error_message = "there should be at least one table"
    }
  }
}

# Ensure the number of tables is equal to exactly one element (with the use of check block)
check "table_check" {
  data "snowflake_tables" "assert_with_check_block" {
    like = "table-name"
  }

  assert {
    condition     = length(data.snowflake_tables.assert_with_check_block.tables) == 1
    error_message = "tables filtered by '${data.snowflake_tables.assert_with_check_block.like}' returned ${length(data.snowflake_tables.assert_with_check_block.tables)} tables where one was expected"
  }
}
