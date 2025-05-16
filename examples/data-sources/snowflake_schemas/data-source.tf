# Simple usage
data "snowflake_schemas" "simple" {
}

output "simple_output" {
  value = data.snowflake_schemas.simple.schemas
}

# Filtering (like)
data "snowflake_schemas" "like" {
  like = "schema-name"
}

output "like_output" {
  value = data.snowflake_schemas.like.schemas
}

# Filtering by prefix (like)
data "snowflake_schemas" "like_prefix" {
  like = "prefix%"
}

output "like_prefix_output" {
  value = data.snowflake_schemas.like_prefix.schemas
}

# Filtering (limit)
data "snowflake_schemas" "limit" {
  limit {
    rows = 10
    from = "prefix-"
  }
}

output "limit_output" {
  value = data.snowflake_schemas.limit.schemas
}

# Filtering (in)
data "snowflake_schemas" "in" {
  in {
    database = "database"
  }
}

output "in_output" {
  value = data.snowflake_schemas.in.schemas
}

# Without additional data (to limit the number of calls make for every found schema)
data "snowflake_schemas" "only_show" {
  # with_describe is turned on by default and it calls DESCRIBE SCHEMA for every schema found and attaches its output to schemas.*.describe_output field
  with_describe = false
  # with_parameters is turned on by default and it calls SHOW PARAMETERS FOR SCHEMA for every schema found and attaches its output to schemas.*.parameters field
  with_parameters = false
}

output "only_show_output" {
  value = data.snowflake_schemas.only_show.schemas
}

# Ensure the number of schemas is equal to at least one element (with the use of postcondition)
data "snowflake_schemas" "assert_with_postcondition" {
  like = "schema-name%"
  lifecycle {
    postcondition {
      condition     = length(self.schemas) > 0
      error_message = "there should be at least one schema"
    }
  }
}

# Ensure the number of schemas is equal to exactly one element (with the use of check block)
check "schema_check" {
  data "snowflake_schemas" "assert_with_check_block" {
    like = "schema-name"
  }

  assert {
    condition     = length(data.snowflake_schemas.assert_with_check_block.schemas) == 1
    error_message = "schemas filtered by '${data.snowflake_schemas.assert_with_check_block.like}' returned ${length(data.snowflake_schemas.assert_with_check_block.schemas)} schemas where one was expected"
  }
}
