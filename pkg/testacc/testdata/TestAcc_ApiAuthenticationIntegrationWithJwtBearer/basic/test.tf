resource "snowflake_api_authentication_integration_with_jwt_bearer" "test" {
  name                   = var.name
  enabled                = var.enabled
  oauth_client_id        = var.oauth_client_id
  oauth_client_secret    = var.oauth_client_secret
  oauth_assertion_issuer = var.oauth_assertion_issuer
}
