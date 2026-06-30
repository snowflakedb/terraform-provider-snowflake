# basic resource
resource "snowflake_api_integration_external_mcp_dynamic_client" "basic" {
  name                 = "external_mcp_dynamic_client_integration"
  oauth_resource_url   = "https://mcp-server.example.com"
  api_allowed_prefixes = ["https://mcp-server.example.com/"]
  enabled              = true
}

# complete resource
resource "snowflake_api_integration_external_mcp_dynamic_client" "complete" {
  name                 = "external_mcp_dynamic_client_integration_complete"
  oauth_resource_url   = "https://mcp-server.example.com"
  api_allowed_prefixes = ["https://mcp-server.example.com/"]
  api_blocked_prefixes = ["https://mcp-server.example.com/blocked/"]
  enabled              = true
  comment              = "Example External MCP Dynamic Client Registration integration"
}
