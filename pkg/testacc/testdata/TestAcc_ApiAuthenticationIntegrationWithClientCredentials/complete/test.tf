resource "snowflake_api_authentication_integration_with_client_credentials" "test" {
  comment                      = var.comment
  enabled                      = var.enabled
  name                         = var.name
  oauth_access_token_validity  = var.oauth_access_token_validity
  oauth_refresh_token_validity = var.oauth_refresh_token_validity
  oauth_client_auth_method     = var.oauth_client_auth_method
  oauth_client_id              = var.oauth_client_id
  oauth_client_secret          = var.oauth_client_secret
  oauth_token_endpoint         = var.oauth_token_endpoint
  oauth_allowed_scopes         = var.oauth_allowed_scopes
}
