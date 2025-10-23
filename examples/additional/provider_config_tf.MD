provider "snowflake" {
	organization_name = "organization_name"
	account_name = "account_name"
	user = "user"
	password = "password"
	warehouse = "SNOWFLAKE"
	protocol = "https"
	port = "443"
	role = "ACCOUNTADMIN"
	validate_default_parameters = true
	client_ip = "1.2.3.4"
	authenticator = "snowflake"
	okta_url = "https://example.com"
	login_timeout = 10
	request_timeout = 20
	jwt_expire_timeout = 30
	client_timeout = 40
	jwt_client_timeout = 50
	external_browser_timeout = 60
	insecure_mode = true
	ocsp_fail_open = true
	keep_session_alive = true
	disable_telemetry = true
	client_request_mfa_token = true
	client_store_temporary_credential = true
	disable_query_context_cache = true
	include_retry_reason = true
	max_retry_count = 3
	driver_tracing = "info"
	tmp_directory_path = "/tmp/terraform-provider/"
	disable_console_login = true
	oauth_client_id         = var.oauth_client_id
	oauth_client_secret     = var.oauth_client_secret
	oauth_token_request_url = var.oauth_token_request_url
	oauth_authorization_url = var.oauth_authorization_url
	oauth_redirect_uri = var.oauth_redirect_uri
	oauth_scope = "session:role:PUBLIC"
	workload_identity_provider = "azure"
	workload_identity_entra_resource = "workload_identity_entra_resource"
	enable_single_use_refresh_tokens = true

	params = {
		param_key = "param_value"
	}
}

# Client ID from the Okta application.
variable "oauth_client_id" {
  type      = string
  sensitive = true
}

# Client Secret from the Okta application.
variable "oauth_client_secret" {
  type      = string
  sensitive = true
}

# Client Token Request URL from the Okta API Authorization Server.
variable "oauth_token_request_url" {
  type      = string
  sensitive = true
}

# Authorization URL for the Oauth flow.
variable "oauth_authorization_url" {
  type      = string
  sensitive = true
}

# Redirect URI for the Oauth flow.
variable "oauth_redirect_uri" {
  type      = string
  sensitive = true
}
