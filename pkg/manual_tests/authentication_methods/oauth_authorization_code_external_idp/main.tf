
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

resource "snowflake_user" "test" {
  name                 = "AUTH_TEST"
  login_name           = var.login_name
  password             = var.password
  must_change_password = false
}

variable "password" {
  type      = string
  sensitive = true
}

variable "login_name" {
  type      = string
  sensitive = true
}

resource "snowflake_external_oauth_integration" "test" {
  name                                            = "EXTERNAL_OAUTH_CODE"
  enabled                                         = true
  external_oauth_type                             = "OKTA"
  external_oauth_issuer                           = var.issuer
  external_oauth_jws_keys_url                     = ["${var.issuer}/v1/keys"]
  external_oauth_audience_list                    = [var.audience]
  external_oauth_snowflake_user_mapping_attribute = "LOGIN_NAME"
  external_oauth_token_user_mapping_claim         = ["sub"]
}

variable "issuer" {
  type      = string
  sensitive = true
}

variable "audience" {
  type      = string
  sensitive = true
}

# Step 2: check the authentication

provider "snowflake" {
  alias                   = "oauth"
  organization_name       = var.organization_name
  account_name            = var.account_name
  authenticator           = "OAUTH_AUTHORIZATION_CODE"
  role                    = "PUBLIC"
  user                    = var.login_name
  oauth_client_id         = var.oauth_client_id
  oauth_client_secret     = var.oauth_client_secret
  oauth_authorization_url = "${var.issuer}/v1/authorize"
  oauth_token_request_url = "${var.issuer}/v1/token"
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

resource "snowflake_execute" "test" {
  provider   = snowflake.oauth
  execute    = "SELECT CURRENT_USER()"
  revert     = "SELECT CURRENT_USER()"
  depends_on = [snowflake_external_oauth_integration.test, snowflake_user.test]
}
