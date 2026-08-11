## Minimal
# The referenced network rule has to be created with mode = "EGRESS".
resource "snowflake_external_access_integration" "basic" {
  name                  = "external_access_integration_name"
  enabled               = true
  allowed_network_rules = [snowflake_network_rule.egress.fully_qualified_name]
}

## Complete (with every optional set)
resource "snowflake_external_access_integration" "complete" {
  name                  = "external_access_integration_name"
  enabled               = true
  allowed_network_rules = [snowflake_network_rule.egress.fully_qualified_name, snowflake_network_rule.other_egress.fully_qualified_name]
  allowed_api_authentication_integrations {
    integrations = [snowflake_api_authentication_integration_with_client_credentials.example.name]
  }
  allowed_authentication_secrets {
    secrets = [snowflake_secret_with_basic_authentication.example.fully_qualified_name]
  }
  comment = "my external access integration"
}

## Allowing every secret in the account
resource "snowflake_external_access_integration" "all_secrets" {
  name                  = "external_access_integration_name"
  enabled               = true
  allowed_network_rules = [snowflake_network_rule.egress.fully_qualified_name]
  allowed_authentication_secrets {
    all = true
  }
}

## Explicitly allowing no secrets and no API authentication integrations
resource "snowflake_external_access_integration" "none_allowed" {
  name                  = "external_access_integration_name"
  enabled               = true
  allowed_network_rules = [snowflake_network_rule.egress.fully_qualified_name]
  allowed_authentication_secrets {
    none = true
  }
  allowed_api_authentication_integrations {
    none = true
  }
}
