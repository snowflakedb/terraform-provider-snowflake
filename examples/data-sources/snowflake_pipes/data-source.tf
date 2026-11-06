# Simple usage
data "snowflake_pipes" "simple" {
}

output "simple_output" {
  value = data.snowflake_pipes.simple.pipes
}

# Filtering (like)
data "snowflake_pipes" "like" {
  like = "pipe-name"
}

output "like_output" {
  value = data.snowflake_pipes.like.pipes
}

# Filtering by prefix (like)
data "snowflake_pipes" "like_prefix" {
  like = "prefix%"
}

output "like_prefix_output" {
  value = data.snowflake_pipes.like_prefix.pipes
}

# Filtering (in)
data "snowflake_pipes" "in_account" {
  in {
    account = true
  }
}

data "snowflake_pipes" "in_database" {
  in {
    database = "<database_name>"
  }
}

data "snowflake_pipes" "in_schema" {
  in {
    schema = "\"<database_name>\".\"<schema_name>\""
  }
}

data "snowflake_pipes" "in_application" {
  in {
    application = "<application_name>"
  }
}

data "snowflake_pipes" "in_application_package" {
  in {
    application_package = "<application_package_name>"
  }
}

output "in_filtered" {
  value = {
    "account" : data.snowflake_pipes.in_account.pipes,
    "database" : data.snowflake_pipes.in_database.pipes,
    "schema" : data.snowflake_pipes.in_schema.pipes,
    "application" : data.snowflake_pipes.in_application.pipes,
    "application_package" : data.snowflake_pipes.in_application_package.pipes,
  }
}

# Without additional data (to limit the number of calls made for every found pipe)
data "snowflake_pipes" "only_show" {
  # with_describe is turned on by default and it calls DESCRIBE PIPE for every pipe found and attaches its output to pipes.*.describe_output field
  with_describe = false
  # with_parameters is turned on by default and it calls SHOW PARAMETERS FOR PIPE for every pipe found and attaches its output to pipes.*.parameters field
  with_parameters = false
}

output "only_show_output" {
  value = data.snowflake_pipes.only_show.pipes
}

# Ensure the number of pipes is equal to at least one element (with the use of postcondition)
data "snowflake_pipes" "assert_with_postcondition" {
  like = "pipe-name%"
  lifecycle {
    postcondition {
      condition     = length(self.pipes) > 0
      error_message = "there should be at least one pipe"
    }
  }
}

# Ensure the number of pipes is equal to exactly one element (with the use of check block)
check "pipe_check" {
  data "snowflake_pipes" "assert_with_check_block" {
    like = "pipe-name"
  }

  assert {
    condition     = length(data.snowflake_pipes.assert_with_check_block.pipes) == 1
    error_message = "pipes filtered by '${data.snowflake_pipes.assert_with_check_block.like}' returned ${length(data.snowflake_pipes.assert_with_check_block.pipes)} pipes where one was expected"
  }
}
