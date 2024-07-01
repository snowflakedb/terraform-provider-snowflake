# basic resource
resource "snowflake_api_authentication_integration_with_jwt_bearer" "test" {
  enabled                = true
  name                   = "foo"
  oauth_client_id        = "foo"
  oauth_client_secret    = "foo"
  oauth_assertion_issuer = "foo"
}
# resource with all fields set
resource "snowflake_api_authentication_integration_with_jwt_bearer" "test" {
  comment                      = "foo"
  enabled                      = true
  name                         = "foo"
  oauth_access_token_validity  = 42
  oauth_authorization_endpoint = "https://example.com"
  oauth_client_auth_method     = "CLIENT_SECRET_POST"
  oauth_client_id              = "foo"
  oauth_client_secret          = "foo"
  oauth_refresh_token_validity = 42
  oauth_token_endpoint         = "https://example.com"
  oauth_assertion_issuer       = "foo"
}
