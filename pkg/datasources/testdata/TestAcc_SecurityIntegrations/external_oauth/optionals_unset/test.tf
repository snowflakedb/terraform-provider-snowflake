resource "snowflake_external_oauth_integration" "test" {
  comment                                         = var.comment
  enabled                                         = var.enabled
  external_oauth_allowed_roles_list               = var.external_oauth_allowed_roles_list
  external_oauth_any_role_mode                    = var.external_oauth_any_role_mode
  external_oauth_audience_list                    = var.external_oauth_audience_list
  external_oauth_issuer                           = var.external_oauth_issuer
  external_oauth_jws_keys_url                     = var.external_oauth_jws_keys_url
  external_oauth_scope_delimiter                  = var.external_oauth_scope_delimiter
  external_oauth_scope_mapping_attribute          = var.external_oauth_scope_mapping_attribute
  external_oauth_snowflake_user_mapping_attribute = var.external_oauth_snowflake_user_mapping_attribute
  external_oauth_token_user_mapping_claim         = var.external_oauth_token_user_mapping_claim
  name                                            = var.name
  external_oauth_type                             = var.external_oauth_type
}

data "snowflake_security_integrations" "test" {
  depends_on = [snowflake_external_oauth_integration.test]

  with_describe = false
  like          = var.name
}
