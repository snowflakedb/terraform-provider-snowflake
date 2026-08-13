# Step 1: create the needed objects (e.g. locally, with a profile that can create users).
provider "snowflake" {
  profile = "default"
}

terraform {
  required_providers {
    snowflake = {
      source  = "snowflakedb/snowflake"
      version = "= 2.20.0"
    }
  }
}

# The claims come from the TFC/TFE-issued workload identity token. Decode the value of
# TFC_WORKLOAD_IDENTITY_TOKEN_SNOWFLAKE from a run to read them.
# Note that TFC/TFE issues a different `sub` claim for the plan and the apply phase, which is what
# allows backing the two phases with different Snowflake identities.
resource "snowflake_service_user" "auth_test" {
  name = "AUTH_TEST"
  default_workload_identity {
    oidc {
      issuer             = var.workload_identity_oidc.issuer
      subject            = var.workload_identity_oidc.subject
      oidc_audience_list = [var.workload_identity_oidc.oidc_audience]
    }
  }
}

variable "workload_identity_oidc" {
  type = object({
    issuer        = string
    subject       = string
    oidc_audience = string
  })
  sensitive = true
}

# Step 2: check the authentication. This step only works inside a TFC/TFE run, because the token is
# provided by TFC/TFE through the TFC_WORKLOAD_IDENTITY_TOKEN_SNOWFLAKE environment variable.
# Configure a manually generated workload identity token with the SNOWFLAKE tag for the workspace, see
# https://developer.hashicorp.com/terraform/enterprise/workspaces/dynamic-provider-credentials/manual-generation.
# No token is passed here - neither through the `token` field nor through SNOWFLAKE_TOKEN.
provider "snowflake" {
  alias                           = "tfc_wif_auth"
  organization_name               = "ORGANIZATION_NAME"
  account_name                    = "ACCOUNT_NAME"
  user                            = snowflake_service_user.auth_test.name
  authenticator                   = "WORKLOAD_IDENTITY"
  workload_identity_provider      = "OIDC"
  tfc_workload_identity_token_tag = "SNOWFLAKE"
  role                            = "ROLE_NAME"
}

resource "snowflake_execute" "test" {
  provider = snowflake.tfc_wif_auth
  execute  = "SELECT CURRENT_USER()"
  revert   = "SELECT CURRENT_USER()"
}
