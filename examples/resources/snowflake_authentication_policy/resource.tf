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
  comment                = "My authentication policy."
}
