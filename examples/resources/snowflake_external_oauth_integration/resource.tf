# basic resource
resource "snowflake_external_oauth_integration" "test" {
  enabled                                         = true
  external_oauth_issuer                           = "issuer"
  external_oauth_snowflake_user_mapping_attribute = "LOGIN_NAME"
  external_oauth_token_user_mapping_claim         = ["upn"]
  name                                            = "test"
  external_oauth_type                             = "CUSTOM"
}
# resource with all fields set (jws keys url and allowed roles)
resource "snowflake_external_oauth_integration" "test" {
  comment                                         = "comment"
  enabled                                         = true
  external_oauth_allowed_roles_list               = [snowflake_role.one.fully_qualified_name]
  external_oauth_any_role_mode                    = "ENABLE"
  external_oauth_audience_list                    = ["https://example.com"]
  external_oauth_issuer                           = "issuer"
  external_oauth_jws_keys_url                     = ["https://example.com"]
  external_oauth_scope_delimiter                  = ","
  external_oauth_scope_mapping_attribute          = "scope"
  external_oauth_snowflake_user_mapping_attribute = "LOGIN_NAME"
  is_agentic                                      = "true"
  external_oauth_token_user_mapping_claim         = ["upn"]
  name                                            = "test"
  external_oauth_type                             = "CUSTOM"
}
# resource with all fields set (rsa public keys and blocked roles)
resource "snowflake_external_oauth_integration" "test" {
  comment                                         = "comment"
  enabled                                         = true
  external_oauth_any_role_mode                    = "ENABLE"
  external_oauth_audience_list                    = ["https://example.com"]
  external_oauth_blocked_roles_list               = [snowflake_role.one.fully_qualified_name]
  external_oauth_issuer                           = "issuer"
  external_oauth_rsa_public_key                   = file("key.pem")
  external_oauth_rsa_public_key_2                 = file("key2.pem")
  external_oauth_scope_delimiter                  = ","
  external_oauth_scope_mapping_attribute          = "scope"
  external_oauth_snowflake_user_mapping_attribute = "LOGIN_NAME"
  is_agentic                                      = "true"
  external_oauth_token_user_mapping_claim         = ["upn"]
  name                                            = "test"
  external_oauth_type                             = "CUSTOM"
}
