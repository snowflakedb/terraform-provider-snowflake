# basic resource
# each pem file contains a base64 encoded IdP signing certificate on a single line without the leading -----BEGIN CERTIFICATE----- and ending -----END CERTIFICATE----- markers.
resource "snowflake_saml2_integration" "saml_integration" {
  name            = "saml_integration"
  saml2_provider  = "CUSTOM"
  saml2_issuer    = "test_issuer"
  saml2_sso_url   = "https://example.com"
  saml2_x509_cert = file("cert.pem")
}
# resource with all fields set
resource "snowflake_saml2_integration" "test" {
  allowed_email_patterns              = ["^(.+dev)@example.com$"]
  allowed_user_domains                = ["example.com"]
  comment                             = "foo"
  enabled                             = true
  name                                = "saml_integration"
  saml2_enable_sp_initiated           = true
  saml2_force_authn                   = true
  saml2_issuer                        = "foo"
  saml2_post_logout_redirect_url      = "https://example.com"
  saml2_provider                      = "CUSTOM"
  saml2_requested_nameid_format       = "urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified"
  saml2_sign_request                  = true
  saml2_snowflake_acs_url             = "example.snowflakecomputing.com/fed/login"
  saml2_snowflake_issuer_url          = "example.snowflakecomputing.com/fed/login"
  saml2_snowflake_x509_cert           = file("snowflake_cert.pem")
  saml2_sp_initiated_login_page_label = "foo"
  saml2_sso_url                       = "https://example.com"
  saml2_x509_cert                     = file("cert.pem")
}
