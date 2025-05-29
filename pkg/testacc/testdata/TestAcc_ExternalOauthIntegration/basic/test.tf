resource "snowflake_external_oauth_integration" "test" {
  name                                            = var.name
  external_oauth_type                             = var.external_oauth_type
  enabled                                         = var.enabled
  external_oauth_issuer                           = var.external_oauth_issuer
  external_oauth_token_user_mapping_claim         = var.external_oauth_token_user_mapping_claim
  external_oauth_snowflake_user_mapping_attribute = var.external_oauth_snowflake_user_mapping_attribute
  external_oauth_jws_keys_url                     = var.external_oauth_jws_keys_url
}
