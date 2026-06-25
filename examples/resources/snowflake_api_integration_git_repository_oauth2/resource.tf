# basic resource
resource "snowflake_api_integration_git_repository_oauth2" "basic" {
  name                         = "git_repository_oauth2_integration"
  oauth_authorization_endpoint = "https://gitlab.com/oauth/authorize"
  oauth_token_endpoint         = "https://gitlab.com/oauth/token"
  oauth_client_id              = "my-client-id"
  oauth_client_secret          = "my-client-secret"
  api_allowed_prefixes         = ["https://gitlab.com/my-org/"]
  enabled                      = true
}

# complete resource
resource "snowflake_api_integration_git_repository_oauth2" "complete" {
  name                         = "git_repository_oauth2_integration_complete"
  oauth_authorization_endpoint = "https://gitlab.com/oauth/authorize"
  oauth_token_endpoint         = "https://gitlab.com/oauth/token"
  oauth_client_id              = "my-client-id"
  oauth_client_secret          = "my-client-secret"
  oauth_access_token_validity  = 3600
  oauth_refresh_token_validity = 86400
  oauth_allowed_scopes         = ["read_repository"]
  oauth_username               = "my-git-user"
  api_allowed_prefixes         = ["https://gitlab.com/my-org/"]
  api_blocked_prefixes         = ["https://gitlab.com/my-org/private-repo/"]
  enabled                      = true
  comment                      = "Example Git Repository OAuth2 integration"
}
