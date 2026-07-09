# basic resource
resource "snowflake_api_integration_external_mcp_oauth2" "basic" {
  name                         = "external_mcp_oauth2_integration"
  oauth_client_id              = "my-client-id"
  oauth_client_secret          = "my-client-secret"
  oauth_token_endpoint         = "https://mcp-server.example.com/oauth/token"
  oauth_authorization_endpoint = "https://mcp-server.example.com/oauth/authorize"
  api_allowed_prefixes         = ["https://mcp-server.example.com/"]
  enabled                      = true
}

# complete resource
resource "snowflake_api_integration_external_mcp_oauth2" "complete" {
  name                         = "external_mcp_oauth2_integration_complete"
  oauth_client_id              = "my-client-id"
  oauth_client_secret          = "my-client-secret"
  oauth_token_endpoint         = "https://mcp-server.example.com/oauth/token"
  oauth_authorization_endpoint = "https://mcp-server.example.com/oauth/authorize"
  oauth_client_auth_method     = "CLIENT_SECRET_POST"
  oauth_refresh_token_validity = 86400
  api_allowed_prefixes         = ["https://mcp-server.example.com/"]
  api_blocked_prefixes         = ["https://mcp-server.example.com/blocked/"]
  enabled                      = true
  comment                      = "Example External MCP OAuth2 integration"
}
