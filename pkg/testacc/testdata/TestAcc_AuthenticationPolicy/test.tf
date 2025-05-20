resource "snowflake_authentication_policy" "authentication_policy" {
  name                       = var.name
  database                   = var.database
  schema                     = var.schema
  authentication_methods     = var.authentication_methods
  mfa_authentication_methods = var.mfa_authentication_methods
  mfa_enrollment             = var.mfa_enrollment
  client_types               = var.client_types
  security_integrations      = var.security_integrations
  comment                    = var.comment
}