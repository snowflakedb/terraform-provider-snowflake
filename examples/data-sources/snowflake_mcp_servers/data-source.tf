# Simple usage
data "snowflake_mcp_servers" "simple" {
}

output "simple_output" {
  value = data.snowflake_mcp_servers.simple.mcp_servers
}

# Filtering (like)
data "snowflake_mcp_servers" "like" {
  like = "mcp-server-name"
}

output "like_output" {
  value = data.snowflake_mcp_servers.like.mcp_servers
}

# Filtering by prefix (like)
data "snowflake_mcp_servers" "like_prefix" {
  like = "prefix%"
}

output "like_prefix_output" {
  value = data.snowflake_mcp_servers.like_prefix.mcp_servers
}

# Filtering (in)
data "snowflake_mcp_servers" "in_database" {
  in {
    database = "<database_name>"
  }
}

data "snowflake_mcp_servers" "in_schema" {
  in {
    schema = "\"<database_name>\".\"<schema_name>\""
  }
}

output "in_filtered" {
  value = {
    "database" : data.snowflake_mcp_servers.in_database.mcp_servers,
    "schema" : data.snowflake_mcp_servers.in_schema.mcp_servers,
  }
}

# Without additional data (to limit the number of calls made for every found MCP server)
data "snowflake_mcp_servers" "only_show" {
  # with_describe is turned on by default and it calls DESCRIBE MCP SERVER for every MCP server found and attaches its output to mcp_servers.*.describe_output field
  with_describe = false
}

output "only_show_output" {
  value = data.snowflake_mcp_servers.only_show.mcp_servers
}

# Ensure the number of MCP servers is equal to at least one element (with the use of postcondition)
data "snowflake_mcp_servers" "assert_with_postcondition" {
  like = "mcp-server-name%"
  lifecycle {
    postcondition {
      condition     = length(self.mcp_servers) > 0
      error_message = "there should be at least one MCP server"
    }
  }
}

# Ensure the number of MCP servers is equal to exactly one element (with the use of check block)
check "mcp_server_check" {
  data "snowflake_mcp_servers" "assert_with_check_block" {
    like = "mcp-server-name"
  }

  assert {
    condition     = length(data.snowflake_mcp_servers.assert_with_check_block.mcp_servers) == 1
    error_message = "MCP servers filtered by '${data.snowflake_mcp_servers.assert_with_check_block.like}' returned ${length(data.snowflake_mcp_servers.assert_with_check_block.mcp_servers)} MCP servers where one was expected"
  }
}
