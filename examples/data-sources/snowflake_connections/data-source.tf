# Simple usage
data "snowflake_connections" "simple" {
}

output "simple_output" {
  value = data.snowflake_connections.simple.connections
}

# Filtering (like)
data "snowflake_connections" "like" {
  like = "connection-name"
}

output "like_output" {
  value = data.snowflake_connections.like.connections
}

# Filtering by prefix (like)
data "snowflake_connections" "like_prefix" {
  like = "prefix%"
}

output "like_prefix_output" {
  value = data.snowflake_connections.like_prefix.connections
}

# Ensure the number of connections is equal to exactly one element (with the use of check block)
check "connection_check" {
  data "snowflake_connections" "assert_with_check_block" {
    like = "connection-name"
  }

  assert {
    condition     = length(data.snowflake_connections.assert_with_check_block.connections) == 1
    error_message = "connections filtered by '${data.snowflake_connections.assert_with_check_block.like}' returned ${length(data.snowflake_connections.assert_with_check_block.connections)} connections where one was expected"
  }
}
