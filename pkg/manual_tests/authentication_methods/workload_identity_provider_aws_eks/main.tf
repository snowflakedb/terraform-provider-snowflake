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
  default_workload_identity {
    oidc {
      issuer  = var.workload_identity_oidc.oidc_issuer_url
      subject = "system:serviceaccount:${var.workload_identity_oidc.namespace}:${var.workload_identity_oidc.service_account}"
    }
  }
}

variable "workload_identity_oidc" {
  type = object({
    oidc_issuer_url = string
    namespace       = string
    service_account = string
  })
  sensitive = true
}

# Step 2: check the authentication
# usually this needs to run in an environment which has access to the infrastructure of the cloud provider, e.g. in a gitlab runner running in an EKS cluster
# for the `AWS EKS Workload Identity Federation with OIDC` case, the token from the token_file has to be provided as an envvar or through the token field.
provider "snowflake" {
  alias                      = "wif_auth"
  organization_name          = "auxmoney"
  account_name               = "terraformtest"
  user                       = snowflake_service_user.auth_test.name
  authenticator              = "WORKLOAD_IDENTITY"
  workload_identity_provider = "OIDC"
  role                       = "ACCOUNTADMIN"
  token                      = file("<token_file_path>")
}

resource "snowflake_execute" "test" {
  provider = snowflake.wif_auth
  execute  = "SELECT CURRENT_USER()"
  revert   = "SELECT CURRENT_USER()"
}
