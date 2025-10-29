# Step 1: create needed objects (e.g. locally)
provider "snowflake" {
  profile = "default"
}

terraform {
  required_providers {
    snowflake = {
      source  = "snowflakedb/snowflake"
      version = "= 2.10.0"
    }
  }
}

resource "snowflake_service_user" "auth_test" {
  name = "AUTH_TEST"
}

// TODO(SNOW-2272350): Update the example once WIF management is supported in users resources
resource "snowflake_execute" "workload_identity_federation_oidc" {
  execute = <<-SQL
    ALTER USER ${snowflake_service_user.auth_test.name}
      SET WORKLOAD_IDENTITY = (
        TYPE = OIDC
        ISSUER = '${var.workload_identity_oidc.issuer}'
        SUBJECT = '${var.workload_identity_oidc.subject}'
        OIDC_AUDIENCE_LIST = ('${var.workload_identity_oidc.oidc_audience}')
      )
  SQL
  revert  = <<-SQL
    ALTER USER IF EXISTS ${snowflake_service_user.auth_test.name}
      UNSET WORKLOAD_IDENTITY;
  SQL

  depends_on = [snowflake_service_user.auth_test]
}

variable "workload_identity_oidc" {
  type = object({
    issuer        = string
    subject       = string
    oidc_audience = string
  })
  sensitive = true
}

# Step 2: check the authentication
# You need to have the token from the OIDC provider. This token must match the values in the user object.
# The token from the token_file has to be provided as an envvar or through the token field.
provider "snowflake" {
  alias                      = "wif_auth"
  organization_name          = "ORGANIZATION_NAME"
  account_name               = "ACCOUNT_NAME"
  user                       = snowflake_service_user.auth_test.name
  authenticator              = "WORKLOAD_IDENTITY"
  workload_identity_provider = "OIDC"
  role                       = "ROLE_NAME"
  token                      = file("<token_file_path>")
}

resource "snowflake_execute" "test" {
  provider = snowflake.wif_auth
  execute  = "SELECT CURRENT_USER()"
  revert   = "SELECT CURRENT_USER()"
}
