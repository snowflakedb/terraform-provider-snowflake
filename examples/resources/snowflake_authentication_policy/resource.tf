## Minimal
resource "snowflake_authentication_policy" "basic" {
  database = "database_name"
  schema   = "schema_name"
  name     = "network_policy_name"
}

## Complete (with every optional set)
resource "snowflake_authentication_policy" "complete" {
  database               = "database_name"
  schema                 = "schema_name"
  name                   = "network_policy_name"
  authentication_methods = ["ALL"]
  mfa_enrollment         = "OPTIONAL"
  client_types           = ["ALL"]
  security_integrations  = ["ALL"]

  mfa_policy = {
    allowed_methods                        = ["PASSKEY", "DUO"]
    enforce_mfa_on_external_authentication = "ALL"
  }
  pat_policy = {
    default_expiry_in_days    = 1
    max_expiry_in_days        = 30
    network_policy_evaluation = "NOT_ENFORCED"
  }
  workload_identity_policy = {
    allowed_providers     = ["ALL"]
    allowed_aws_accounts  = ["111122223333"]
    allowed_azure_issuers = ["https://login.microsoftonline.com/tenantid/v2.0"]
    allowed_oidc_issuers  = ["https://example.com"]
  }
  comment = "My authentication policy."
}
