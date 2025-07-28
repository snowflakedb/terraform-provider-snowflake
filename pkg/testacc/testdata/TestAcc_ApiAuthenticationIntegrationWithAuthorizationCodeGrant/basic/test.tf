resource "snowflake_api_authentication_integration_with_authorization_code_grant" "test" {
  enabled             = var.enabled
  name                = var.name
  oauth_client_id     = var.oauth_client_id
  oauth_client_secret = var.oauth_client_secret
}
