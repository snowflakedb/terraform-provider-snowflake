## Basic example with a network rule
resource "snowflake_external_access_integration" "basic" {
  name                  = "my_external_access_integration"
  allowed_network_rules = ["mydb.myschema.my_network_rule"]
  enabled               = true
}

## Example with an authentication secret and comment
resource "snowflake_external_access_integration" "with_secret" {
  name                           = "my_external_access_integration_with_secret"
  allowed_network_rules          = ["mydb.myschema.my_network_rule"]
  allowed_authentication_secrets = ["mydb.myschema.my_secret"]
  enabled                        = true
  comment                        = "Integration for accessing external API with authentication"
}
