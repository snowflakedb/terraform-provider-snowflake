
# Step 1: create needed objects
provider "snowflake" {
  profile = "default"
}

terraform {
  required_providers {
    snowflake = {
      source  = "snowflakedb/snowflake"
      version = "= 2.8.0"
    }
  }
}

resource "snowflake_oauth_integration_for_custom_clients" "test" {
  name                             = "OAUTH_CODE"
  enabled                          = true
  oauth_client_type                = "CONFIDENTIAL"
  oauth_allow_non_tls_redirect_uri = true
  oauth_redirect_uri               = "http://localhost:8001"
  oauth_issue_refresh_tokens       = true
  oauth_enforce_pkce               = false
  pre_authorized_roles_list        = ["PUBLIC"]
  oauth_refresh_token_validity     = 86400
}

# Step 2: Get the client id and client secret
resource "snowflake_execute" "show_oauth_client_secrets" {
  query   = "SELECT SYSTEM$SHOW_OAUTH_CLIENT_SECRETS('${snowflake_oauth_integration_for_custom_clients.test.name}')"
  execute = "SELECT 1"
  revert  = "SELECT 1"
}

# Step 3: check the authentication

provider "snowflake" {
  alias                   = "oauth"
  organization_name       = var.organization_name
  account_name            = var.account_name
  authenticator           = "OAUTH_AUTHORIZATION_CODE"
  role                    = "PUBLIC"
  user                    = var.login_name
  oauth_client_id         = var.oauth_client_id
  oauth_client_secret     = var.oauth_client_secret
  oauth_authorization_url = "${var.issuer}/oauth/authorize"
  oauth_token_request_url = "${var.issuer}/oauth/token-request"
  oauth_redirect_uri      = "http://localhost:8001"
  oauth_scope             = "session:role:PUBLIC"
}

variable "organization_name" {
  type      = string
  sensitive = true
}

variable "account_name" {
  type      = string
  sensitive = true
}

variable "oauth_client_id" {
  type      = string
  sensitive = true
}

variable "oauth_client_secret" {
  type      = string
  sensitive = true
}

variable "login_name" {
  type      = string
  sensitive = true
}

variable "issuer" {
  type      = string
  sensitive = true
}

resource "snowflake_execute" "test" {
  provider   = snowflake.oauth
  execute    = "SELECT CURRENT_USER()"
  revert     = "SELECT CURRENT_USER()"
  depends_on = [snowflake_execute.show_oauth_client_secrets]
}
