# Simple usage
data "snowflake_streams" "simple" {
}

output "simple_output" {
  value = data.snowflake_streams.simple.streams
}

# Filtering (like)
data "snowflake_streams" "like" {
  like = "stream-name"
}

output "like_output" {
  value = data.snowflake_streams.like.streams
}

# Filtering by prefix (like)
data "snowflake_streams" "like_prefix" {
  like = "prefix%"
}

output "like_prefix_output" {
  value = data.snowflake_streams.like_prefix.streams
}

# Filtering (limit)
data "snowflake_streams" "limit" {
  limit {
    rows = 10
    from = "prefix-"
  }
}

output "limit_output" {
  value = data.snowflake_streams.limit.streams
}

# Filtering (in)
data "snowflake_streams" "in_account" {
  in {
    account = true
  }
}

data "snowflake_streams" "in_database" {
  in {
    database = "<database_name>"
  }
}

data "snowflake_streams" "in_schema" {
  in {
    schema = "<database_name>.<schema_name>"
  }
}

data "snowflake_streams" "in_application" {
  in {
    application = "<application_name>"
  }
}

data "snowflake_streams" "in_application_package" {
  in {
    application_package = "<application_package_name>"
  }
}

output "in_output" {
  value = {
    "account" : data.snowflake_streams.in_account.streams,
    "database" : data.snowflake_streams.in_database.streams,
    "schema" : data.snowflake_streams.in_schema.streams,
    "application" : data.snowflake_streams.in_application.streams,
    "application_package" : data.snowflake_streams.in_application_package.streams,
  }
}

output "in_output" {
  value = data.snowflake_streams.in.streams
}

# Without additional data (to limit the number of calls make for every found stream)
data "snowflake_streams" "only_show" {
  # with_describe is turned on by default and it calls DESCRIBE STREAM for every stream found and attaches its output to streams.*.describe_output field
  with_describe = false
}

output "only_show_output" {
  value = data.snowflake_streams.only_show.streams
}

# Ensure the number of streams is equal to at least one element (with the use of postcondition)
data "snowflake_streams" "assert_with_postcondition" {
  like = "stream-name%"
  lifecycle {
    postcondition {
      condition     = length(self.streams) > 0
      error_message = "there should be at least one stream"
    }
  }
}

# Ensure the number of streams is equal to exactly one element (with the use of check block)
check "stream_check" {
  data "snowflake_streams" "assert_with_check_block" {
    like = "stream-name"
  }

  assert {
    condition     = length(data.snowflake_streams.assert_with_check_block.streams) == 1
    error_message = "streams filtered by '${data.snowflake_streams.assert_with_check_block.like}' returned ${length(data.snowflake_streams.assert_with_check_block.streams)} streams where one was expected"
  }
}
