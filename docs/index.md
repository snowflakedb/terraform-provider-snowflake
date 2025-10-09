---
page_title: "Provider: Snowflake"
description: Manage SnowflakeDB with Terraform.
---

## Support

For official support and urgent, production-impacting issues, please [contact Snowflake Support](https://community.snowflake.com/s/article/How-To-Submit-a-Support-Case-in-Snowflake-Lodge).

~> **Keep in mind** that the official support starts with the [v2.0.0](https://registry.terraform.io/providers/snowflakedb/snowflake/2.0.0) version for stable resources only. All previous versions and preview resources are not officially supported. Also, consult [supported architectures](#supported-architectures).

Please follow [creating issues guidelines](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/CREATING_ISSUES.md), [FAQ](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/FAQ.md), and [known issues](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/KNOWN_ISSUES.md) before submitting an issue on GitHub or directly to Snowflake Support.

# Snowflake Provider

~> **Disclaimer** The project is in GA version, but some features are in preview. Such resources and data sources are considered preview features in the provider, regardless of their state in Snowflake. We do not guarantee their stability. They will be reworked and marked as a stable feature in future releases. Breaking changes in these features are expected, even without bumping the major version. They are disabled by default. To use them, add the relevant feature name to `preview_features_enabled` field in the [provider configuration](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs#schema). The list of preview features is available below. Please always refer to the [Getting Help](https://github.com/snowflakedb/terraform-provider-snowflake?tab=readme-ov-file#getting-help) section in our Github repo to best determine how to get help for your questions.

~> **Note** Please check the [migration guide](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/MIGRATION_GUIDE.md) when changing the version of the provider.

!> **Sensitive values** Important: Do not include credentials, personal identifiers, or other regulated or sensitive information (e.g., GDPR, HIPAA, PCI-DSS data) in non-sensitive fields. Snowflake marks specific fields as sensitive—such as passwords, private keys, and tokens, meaning these fields will not appear in logs. Each sensitive field is properly marked in the documentation. All other fields are treated as non-sensitive by default. Some of them, like [task's](./resources/task) configuration, may contain sensitive information but are not marked as sensitive - you are responsible for safeguarding these fields according to your organization's security standards and regulatory requirements. Snowflake will not be liable for any exposure of data placed in non-sensitive fields. Read more in the [Sensitive values limitations](#sensitive-values-limitations) section.

-> **Note** The current roadmap is available in our GitHub repository: [ROADMAP.md](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/ROADMAP.md).

This is a terraform provider plugin for managing [Snowflake](https://www.snowflake.com/) accounts.
Coverage is focused on part of Snowflake related to access control.

## Supported architectures

We have compiled a list to clarify which binaries are officially supported and which are provided additionally but not officially supported.
The lists are based on what the underlying [gosnowflake driver](https://github.com/snowflakedb/gosnowflake) supports and what [HashiCorp recommends for Terraform providers](https://developer.hashicorp.com/terraform/registry/providers/os-arch).

The provider officially supports the binaries built for the following OSes and architectures:
- Windows: amd64
- Linux: amd64 and arm64
- Darwin: amd64 and arm64

Currently, we also provide the binaries for the following OSes and architectures, but they are not officially supported, and we do not prioritize fixes for them:
- Windows: arm64 and 386
- Linux: 386
- Darwin: 386
- Freebsd: any architecture

## Example Provider Configuration

This is an example configuration of the provider in `main.tf` in a configuration directory. More examples are provided [below](#order-precedence).

```terraform
terraform {
  required_providers {
    snowflake = {
      source = "snowflakedb/snowflake"
    }
  }
}

# A simple configuration of the provider with a default authentication.
# A default value for `authenticator` is `snowflake`, enabling authentication with `user` and `password`.
provider "snowflake" {
  organization_name = "..." # required if not using profile. Can also be set via SNOWFLAKE_ORGANIZATION_NAME env var
  account_name      = "..." # required if not using profile. Can also be set via SNOWFLAKE_ACCOUNT_NAME env var
  user              = "..." # required if not using profile or token. Can also be set via SNOWFLAKE_USER env var
  password          = "..."

  // optional
  role      = "..."
  host      = "..."
  warehouse = "..."
  params = {
    query_tag = "..."
  }
}

# A simple configuration of the provider with private key authentication.
provider "snowflake" {
  organization_name      = "..." # required if not using profile. Can also be set via SNOWFLAKE_ORGANIZATION_NAME env var
  account_name           = "..." # required if not using profile. Can also be set via SNOWFLAKE_ACCOUNT_NAME env var
  user                   = "..." # required if not using profile or token. Can also be set via SNOWFLAKE_USER env var
  authenticator          = "SNOWFLAKE_JWT"
  private_key            = file("~/.ssh/snowflake_key.p8")
  private_key_passphrase = var.private_key_passphrase
}

# Remember to provide the passphrase securely.
variable "private_key_passphrase" {
  type      = string
  sensitive = true
}

# By using the `profile` field, missing fields will be populated from ~/.snowflake/config TOML file
provider "snowflake" {
  profile = "securityadmin"
}
```

## Configuration Schema

**Warning: these values are passed directly to the gosnowflake library, which may not work exactly the way you expect. See the [gosnowflake docs](https://godoc.org/github.com/snowflakedb/gosnowflake#hdr-Connection_Parameters) for more.**

-> **Note**: In Go Snowflake driver 1.12.1 ([release notes](https://docs.snowflake.com/en/release-notes/clients-drivers/golang-2024#version-1-12-1-december-05-2024)), configuration field `InsecureMode` has been deprecated in favor of `DisableOCSPChecks`. This field is not available in the provider yet. Please use `InsecureMode` instead, which has the same behavior. We are planning to support this new field and deprecate the old one.

-> **Note** If a field has a default value, it is shown next to the type in the schema. Most of the values in provider schema can be sourced from environment value (check field descriptions), but If a specified environment variable is not found, then the driver's default value is used instead.

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `account_name` (String) Specifies your Snowflake account name assigned by Snowflake. For information about account identifiers, see the [Snowflake documentation](https://docs.snowflake.com/en/user-guide/admin-account-identifier#account-name). Required unless using `profile`. Can also be sourced from the `SNOWFLAKE_ACCOUNT_NAME` environment variable.
- `authenticator` (String) Specifies the [authentication type](https://pkg.go.dev/github.com/snowflakedb/gosnowflake#AuthType) to use when connecting to Snowflake. Valid options are: `SNOWFLAKE` | `OAUTH` | `EXTERNALBROWSER` | `OKTA` | `SNOWFLAKE_JWT` | `TOKENACCESSOR` | `USERNAMEPASSWORDMFA` | `PROGRAMMATIC_ACCESS_TOKEN` | `OAUTH_CLIENT_CREDENTIALS` | `OAUTH_AUTHORIZATION_CODE`. Can also be sourced from the `SNOWFLAKE_AUTHENTICATOR` environment variable.
- `client_ip` (String) IP address for network checks. Can also be sourced from the `SNOWFLAKE_CLIENT_IP` environment variable.
- `client_request_mfa_token` (String) When true the MFA token is cached in the credential manager. True by default in Windows/OSX. False for Linux. Can also be sourced from the `SNOWFLAKE_CLIENT_REQUEST_MFA_TOKEN` environment variable.
- `client_store_temporary_credential` (String) When true the ID token is cached in the credential manager. True by default in Windows/OSX. False for Linux. Can also be sourced from the `SNOWFLAKE_CLIENT_STORE_TEMPORARY_CREDENTIAL` environment variable.
- `client_timeout` (Number) The timeout in seconds for the client to complete the authentication. Can also be sourced from the `SNOWFLAKE_CLIENT_TIMEOUT` environment variable.
- `disable_console_login` (String) Indicates whether console login should be disabled in the driver. Can also be sourced from the `SNOWFLAKE_DISABLE_CONSOLE_LOGIN` environment variable.
- `disable_query_context_cache` (Boolean) Disables HTAP query context cache in the driver. Can also be sourced from the `SNOWFLAKE_DISABLE_QUERY_CONTEXT_CACHE` environment variable.
- `disable_telemetry` (Boolean) Disables telemetry in the driver. Can also be sourced from the `DISABLE_TELEMETRY` environment variable.
- `driver_tracing` (String) Specifies the logging level to be used by the driver. Valid options are: `trace` | `debug` | `info` | `print` | `warning` | `error` | `fatal` | `panic`. Can also be sourced from the `SNOWFLAKE_DRIVER_TRACING` environment variable.
- `enable_single_use_refresh_tokens` (Boolean) Enables single use refresh tokens for Snowflake IdP. Can also be sourced from the `SNOWFLAKE_ENABLE_SINGLE_USE_REFRESH_TOKENS` environment variable.
- `external_browser_timeout` (Number) The timeout in seconds for the external browser to complete the authentication. Can also be sourced from the `SNOWFLAKE_EXTERNAL_BROWSER_TIMEOUT` environment variable.
- `host` (String) Specifies a custom host value used by the driver for privatelink connections. Can also be sourced from the `SNOWFLAKE_HOST` environment variable.
- `include_retry_reason` (String) Should retried request contain retry reason. Can also be sourced from the `SNOWFLAKE_INCLUDE_RETRY_REASON` environment variable.
- `insecure_mode` (Boolean) If true, bypass the Online Certificate Status Protocol (OCSP) certificate revocation check. IMPORTANT: Change the default value for testing or emergency situations only. Can also be sourced from the `SNOWFLAKE_INSECURE_MODE` environment variable.
- `jwt_client_timeout` (Number) The timeout in seconds for the JWT client to complete the authentication. Can also be sourced from the `SNOWFLAKE_JWT_CLIENT_TIMEOUT` environment variable.
- `jwt_expire_timeout` (Number) JWT expire after timeout in seconds. Can also be sourced from the `SNOWFLAKE_JWT_EXPIRE_TIMEOUT` environment variable.
- `keep_session_alive` (Boolean) Enables the session to persist even after the connection is closed. Can also be sourced from the `SNOWFLAKE_KEEP_SESSION_ALIVE` environment variable.
- `login_timeout` (Number) Login retry timeout in seconds EXCLUDING network roundtrip and read out http response. Can also be sourced from the `SNOWFLAKE_LOGIN_TIMEOUT` environment variable.
- `max_retry_count` (Number) Specifies how many times non-periodic HTTP request can be retried by the driver. Can also be sourced from the `SNOWFLAKE_MAX_RETRY_COUNT` environment variable.
- `oauth_authorization_url` (String, Sensitive) Authorization URL of OAuth2 external IdP. See [Snowflake OAuth documentation](https://docs.snowflake.com/en/user-guide/oauth). Can also be sourced from the `SNOWFLAKE_OAUTH_AUTHORIZATION_URL` environment variable.
- `oauth_client_id` (String, Sensitive) Client id for OAuth2 external IdP. See [Snowflake OAuth documentation](https://docs.snowflake.com/en/user-guide/oauth). Can also be sourced from the `SNOWFLAKE_OAUTH_CLIENT_ID` environment variable.
- `oauth_client_secret` (String, Sensitive) Client secret for OAuth2 external IdP. See [Snowflake OAuth documentation](https://docs.snowflake.com/en/user-guide/oauth). Can also be sourced from the `SNOWFLAKE_OAUTH_CLIENT_SECRET` environment variable.
- `oauth_redirect_uri` (String, Sensitive) Redirect URI registered in IdP. See [Snowflake OAuth documentation](https://docs.snowflake.com/en/user-guide/oauth). Can also be sourced from the `SNOWFLAKE_OAUTH_REDIRECT_URI` environment variable.
- `oauth_scope` (String) Comma separated list of scopes. If empty it is derived from role. See [Snowflake OAuth documentation](https://docs.snowflake.com/en/user-guide/oauth). Can also be sourced from the `SNOWFLAKE_OAUTH_SCOPE` environment variable.
- `oauth_token_request_url` (String, Sensitive) Token request URL of OAuth2 external IdP. See [Snowflake OAuth documentation](https://docs.snowflake.com/en/user-guide/oauth). Can also be sourced from the `SNOWFLAKE_OAUTH_TOKEN_REQUEST_URL` environment variable.
- `ocsp_fail_open` (String) True represents OCSP fail open mode. False represents OCSP fail closed mode. Fail open true by default. Can also be sourced from the `SNOWFLAKE_OCSP_FAIL_OPEN` environment variable.
- `okta_url` (String) The URL of the Okta server. e.g. https://example.okta.com. Okta URL host needs to to have a suffix `okta.com`. Read more in Snowflake [docs](https://docs.snowflake.com/en/user-guide/oauth-okta). Can also be sourced from the `SNOWFLAKE_OKTA_URL` environment variable.
- `organization_name` (String) Specifies your Snowflake organization name assigned by Snowflake. For information about account identifiers, see the [Snowflake documentation](https://docs.snowflake.com/en/user-guide/admin-account-identifier#organization-name). Required unless using `profile`. Can also be sourced from the `SNOWFLAKE_ORGANIZATION_NAME` environment variable.
- `params` (Map of String) Sets other connection (i.e. session) parameters. [Parameters](https://docs.snowflake.com/en/sql-reference/parameters). This field can not be set with environmental variables.
- `passcode` (String, Sensitive) Specifies the passcode provided by Duo when using multi-factor authentication (MFA) for login. Can also be sourced from the `SNOWFLAKE_PASSCODE` environment variable.
- `passcode_in_password` (Boolean) False by default. Set to true if the MFA passcode is embedded to the configured password. Can also be sourced from the `SNOWFLAKE_PASSCODE_IN_PASSWORD` environment variable.
- `password` (String, Sensitive) Password for user + password or [token](https://docs.snowflake.com/en/user-guide/programmatic-access-tokens#generating-a-programmatic-access-token) for [PAT auth](https://docs.snowflake.com/en/user-guide/programmatic-access-tokens). Cannot be used with `private_key` and `private_key_passphrase`. Can also be sourced from the `SNOWFLAKE_PASSWORD` environment variable.
- `port` (Number) Specifies a custom port value used by the driver for privatelink connections. Can also be sourced from the `SNOWFLAKE_PORT` environment variable.
- `preview_features_enabled` (Set of String) A list of preview features that are handled by the provider. See [preview features list](https://github.com/Snowflake-Labs/terraform-provider-snowflake/blob/main/v1-preparations/LIST_OF_PREVIEW_FEATURES_FOR_V1.md). Preview features may have breaking changes in future releases, even without raising the major version. This field can not be set with environmental variables. Valid options are: `snowflake_account_authentication_policy_attachment_resource` | `snowflake_account_password_policy_attachment_resource` | `snowflake_alert_resource` | `snowflake_alerts_datasource` | `snowflake_api_integration_resource` | `snowflake_authentication_policy_resource` | `snowflake_compute_pool_resource` | `snowflake_compute_pools_datasource` | `snowflake_cortex_search_service_resource` | `snowflake_cortex_search_services_datasource` | `snowflake_current_account_resource` | `snowflake_current_account_datasource` | `snowflake_current_organization_account_resource` | `snowflake_database_datasource` | `snowflake_database_role_datasource` | `snowflake_dynamic_table_resource` | `snowflake_dynamic_tables_datasource` | `snowflake_external_function_resource` | `snowflake_external_functions_datasource` | `snowflake_external_table_resource` | `snowflake_external_tables_datasource` | `snowflake_external_volume_resource` | `snowflake_failover_group_resource` | `snowflake_failover_groups_datasource` | `snowflake_file_format_resource` | `snowflake_file_formats_datasource` | `snowflake_function_java_resource` | `snowflake_function_javascript_resource` | `snowflake_function_python_resource` | `snowflake_function_scala_resource` | `snowflake_function_sql_resource` | `snowflake_functions_datasource` | `snowflake_git_repository_resource` | `snowflake_git_repositories_datasource` | `snowflake_image_repository_resource` | `snowflake_image_repositories_datasource` | `snowflake_job_service_resource` | `snowflake_listing_resource` | `snowflake_managed_account_resource` | `snowflake_materialized_view_resource` | `snowflake_materialized_views_datasource` | `snowflake_network_policy_attachment_resource` | `snowflake_network_rule_resource` | `snowflake_email_notification_integration_resource` | `snowflake_notification_integration_resource` | `snowflake_object_parameter_resource` | `snowflake_password_policy_resource` | `snowflake_pipe_resource` | `snowflake_pipes_datasource` | `snowflake_current_role_datasource` | `snowflake_service_resource` | `snowflake_services_datasource` | `snowflake_sequence_resource` | `snowflake_sequences_datasource` | `snowflake_share_resource` | `snowflake_shares_datasource` | `snowflake_parameters_datasource` | `snowflake_procedure_java_resource` | `snowflake_procedure_javascript_resource` | `snowflake_procedure_python_resource` | `snowflake_procedure_scala_resource` | `snowflake_procedure_sql_resource` | `snowflake_procedures_datasource` | `snowflake_stage_resource` | `snowflake_stages_datasource` | `snowflake_storage_integration_resource` | `snowflake_storage_integrations_datasource` | `snowflake_system_generate_scim_access_token_datasource` | `snowflake_system_get_aws_sns_iam_policy_datasource` | `snowflake_system_get_privatelink_config_datasource` | `snowflake_system_get_snowflake_platform_info_datasource` | `snowflake_table_column_masking_policy_application_resource` | `snowflake_table_constraint_resource` | `snowflake_table_resource` | `snowflake_tables_datasource` | `snowflake_user_authentication_policy_attachment_resource` | `snowflake_user_public_keys_resource` | `snowflake_user_password_policy_attachment_resource` | `snowflake_user_programmatic_access_token_resource` | `snowflake_user_programmatic_access_tokens_datasource`.
- `private_key` (String, Sensitive) Private Key for username+private-key auth. Cannot be used with `password`. Can also be sourced from the `SNOWFLAKE_PRIVATE_KEY` environment variable.
- `private_key_passphrase` (String, Sensitive) Supports the encryption ciphers aes-128-cbc, aes-128-gcm, aes-192-cbc, aes-192-gcm, aes-256-cbc, aes-256-gcm, and des-ede3-cbc. Can also be sourced from the `SNOWFLAKE_PRIVATE_KEY_PASSPHRASE` environment variable.
- `profile` (String) Sets the profile to read from ~/.snowflake/config file. Can also be sourced from the `SNOWFLAKE_PROFILE` environment variable.
- `protocol` (String) A protocol used in the connection. Valid options are: `http` | `https`. Can also be sourced from the `SNOWFLAKE_PROTOCOL` environment variable.
- `request_timeout` (Number) request retry timeout in seconds EXCLUDING network roundtrip and read out http response. Can also be sourced from the `SNOWFLAKE_REQUEST_TIMEOUT` environment variable.
- `role` (String) Specifies the role to use by default for accessing Snowflake objects in the client session. Can also be sourced from the `SNOWFLAKE_ROLE` environment variable.
- `skip_toml_file_permission_verification` (Boolean) False by default. Skips TOML configuration file permission verification. This flag has no effect on Windows systems, as the permissions are not checked on this platform. Instead of skipping the permissions verification, we recommend setting the proper privileges - see [the section below](#toml-file-limitations). Can also be sourced from the `SNOWFLAKE_SKIP_TOML_FILE_PERMISSION_VERIFICATION` environment variable.
- `tmp_directory_path` (String) Sets temporary directory used by the driver for operations like encrypting, compressing etc. Can also be sourced from the `SNOWFLAKE_TMP_DIRECTORY_PATH` environment variable.
- `token` (String, Sensitive) Token to use for OAuth and other forms of token based auth. When this field is set here, or in the TOML file, the provider sets the `authenticator` to `OAUTH`. Optionally, set the `authenticator` field to the authenticator you want to use. Can also be sourced from the `SNOWFLAKE_TOKEN` environment variable.
- `token_accessor` (Block List, Max: 1) If you are using the OAuth authentication flows, use the dedicated `authenticator` and `oauth...` fields instead. See our [authentication methods guide](./guides/authentication_methods) for more information. (see [below for nested schema](#nestedblock--token_accessor))
- `use_legacy_toml_file` (Boolean) False by default. When this is set to true, the provider expects the legacy TOML format. Otherwise, it expects the new format. See more in [the section below](#examples) Can also be sourced from the `SNOWFLAKE_USE_LEGACY_TOML_FILE` environment variable.
- `user` (String) Username. Required unless using `profile`. Can also be sourced from the `SNOWFLAKE_USER` environment variable.
- `validate_default_parameters` (String) True by default. If false, disables the validation checks for Database, Schema, Warehouse and Role at the time a connection is established. Can also be sourced from the `SNOWFLAKE_VALIDATE_DEFAULT_PARAMETERS` environment variable.
- `warehouse` (String) Specifies the virtual warehouse to use by default for queries, loading, etc. in the client session. Can also be sourced from the `SNOWFLAKE_WAREHOUSE` environment variable.

<a id="nestedblock--token_accessor"></a>
### Nested Schema for `token_accessor`

Required:

- `client_id` (String, Sensitive) The client ID for the OAuth provider when using a refresh token to renew access token. Can also be sourced from the `SNOWFLAKE_TOKEN_ACCESSOR_CLIENT_ID` environment variable.
- `client_secret` (String, Sensitive) The client secret for the OAuth provider when using a refresh token to renew access token. Can also be sourced from the `SNOWFLAKE_TOKEN_ACCESSOR_CLIENT_SECRET` environment variable.
- `redirect_uri` (String, Sensitive) The redirect URI for the OAuth provider when using a refresh token to renew access token. Can also be sourced from the `SNOWFLAKE_TOKEN_ACCESSOR_REDIRECT_URI` environment variable.
- `refresh_token` (String, Sensitive) The refresh token for the OAuth provider when using a refresh token to renew access token. Can also be sourced from the `SNOWFLAKE_TOKEN_ACCESSOR_REFRESH_TOKEN` environment variable.
- `token_endpoint` (String, Sensitive) The token endpoint for the OAuth provider e.g. https://{yourDomain}/oauth/token when using a refresh token to renew access token. Can also be sourced from the `SNOWFLAKE_TOKEN_ACCESSOR_TOKEN_ENDPOINT` environment variable.

## Authentication

The Snowflake provider supports multiple ways to authenticate:

* Password
* PAT (Personal Access Token)
* OAuth Access Token
* OAuth Refresh Token
* Browser Auth
* Private Key
* Config File
* Oauth with Client Credentials
* Oauth with Authorization Code

In all cases `organization_name`, and `account_name` are required. In all cases except for Oauth with Client Credentials, `user` is required.

-> **Note** Storing the credentials and other secret values safely is on the users' side. Read more in [Authentication Methods guide](./guides/authentication_methods).

### Keypair Authentication Environment Variables

You should generate the public and private keys and set up environment variables.

```shell

cd ~/.ssh
openssl genrsa -out snowflake_key 4096
openssl rsa -in snowflake_key -pubout -out snowflake_key.pub
```

To export the variables into your provider:

```shell
export SNOWFLAKE_USER="..."
export SNOWFLAKE_PRIVATE_KEY=$(cat ~/.ssh/snowflake_key.p8)
```

### Keypair Authentication Passphrase

If your private key requires a passphrase, then this can be supplied via the
environment variable `SNOWFLAKE_PRIVATE_KEY_PASSPHRASE`.

Only the ciphers aes-128-cbc, aes-128-gcm, aes-192-cbc, aes-192-gcm,
aes-256-cbc, aes-256-gcm, and des-ede3-cbc are supported on the private key

```shell
cd ~/.ssh
openssl genrsa -out snowflake_key 4096
openssl rsa -in snowflake_key -pubout -out snowflake_key.pub
openssl pkcs8 -topk8 -inform pem -in snowflake_key -outform PEM -v2 aes-256-cbc -out snowflake_key.p8
```

To export the variables into your provider:

```shell
export SNOWFLAKE_USER="..."
export SNOWFLAKE_PRIVATE_KEY=$(cat ~/.ssh/snowflake_key.p8)
export SNOWFLAKE_PRIVATE_KEY_PASSPHRASE="..."
```

### OAuth Access Token

If you have an OAuth access token, export these credentials as environment variables:

```shell
export SNOWFLAKE_USER='...'
export SNOWFLAKE_TOKEN='...'
```

Note that once this access token expires, you'll need to request a new one through an external application.

### OAuth Refresh Token

If you have an OAuth Refresh token, export these credentials as environment variables:

```shell
export SNOWFLAKE_TOKEN_ACCESSOR_REFRESH_TOKEN='...'
export SNOWFLAKE_TOKEN_ACCESSOR_CLIENT_ID='...'
export SNOWFLAKE_TOKEN_ACCESSOR_CLIENT_SECRET='...'
export SNOWFLAKE_TOKEN_ACCESSOR_TOKEN_ENDPOINT='...'
export SNOWFLAKE_TOKEN_ACCESSOR_REDIRECT_URI='https://localhost.com'
```

Note because access token have a short life; typically 10 minutes, by passing refresh token new access token will be generated.

### Username and Password Environment Variables

If you choose to use Username and Password Authentication, export these credentials:

```shell
export SNOWFLAKE_USER='...'
export SNOWFLAKE_PASSWORD='...'
```

### Oauth with Client Credentials Environment Variables

If you choose to use Oauth with Client Credentials Authentication, export these credentials:

```shell
export SNOWFLAKE_OAUTH_CLIENT_ID='...'
export SNOWFLAKE_OAUTH_CLIENT_SECRET='...'
export SNOWFLAKE_OAUTH_TOKEN_REQUEST_URL='...'
```

### Oauth with Authorization Code Environment Variables

If you choose to use Oauth with Authorization Code Authentication, export these credentials:

```shell
export SNOWFLAKE_OAUTH_CLIENT_ID='...'
export SNOWFLAKE_OAUTH_CLIENT_SECRET='...'
export SNOWFLAKE_OAUTH_AUTHORIZATION_URL='...'
export SNOWFLAKE_OAUTH_TOKEN_REQUEST_URL='...'
export SNOWFLAKE_OAUTH_REDIRECT_URI='...'
export SNOWFLAKE_OAUTH_SCOPE='...'
```
## Order Precedence

Currently, the provider can be configured in three ways:
1. In a Terraform file located in the Terraform module with other resources.
2. In environmental variables (envs). This is mainly used to provide sensitive values.
3. In a TOML file (default in `~/.snowflake/config`).

### Terraform file located in the Terraform module with other resources
One of the methods of configuring the provider is in the Terraform module. Read more in the [Terraform docs](https://developer.hashicorp.com/terraform/language/providers/configuration).

Example content of the Terraform file configuration:

```terraform
provider "snowflake" {
    organization_name = "..."
    account_name = "..."
    username = "..."
    password = "..."
}
```

### Environmental variables
The second method is to use environmental variables. This is mainly used to provide sensitive values.

```bash
export SNOWFLAKE_USER="..."
export SNOWFLAKE_PRIVATE_KEY=$(cat ~/.ssh/snowflake_key.p8)
```

### TOML file
The third method is to use a TOML configuration file (default location in ~/.snowflake/config). Notice the use of different profiles. The profile name needs to be specified in the Terraform configuration file in `profile` field. When this is not specified, `default` profile is loaded.
When a `default` profile is not present in the TOML file, it is treated as "empty", without failing.

Read [TOML](https://toml.io/en/) specification for more details on the syntax.

-> **Note** This configuration file is distinct from the ones used to configure [Snowflake CLI](https://docs.snowflake.com/en/developer-guide/snowflake-cli/connecting/configure-cli) or [SnowSQL](https://docs.snowflake.com/en/user-guide/snowsql-config).

Example content of the Terraform file configuration:

```terraform
provider "snowflake" {
    profile = "default"
}
```

Example content of the TOML file configuration is listed below. Note that this example follows a new TOML format, for the legacy format see [examples](#examples) section.

```toml
[default]
organization_name='organization_name'
account_name='account_name'
user='user'
password='password'
role='ACCOUNTADMIN'

[secondary_test_account]
organization_name='organization_name'
account_name='account2_name'
user='user'
password='password'
role='ACCOUNTADMIN'
```

#### TOML file limitations
To ensure a better security of the provider, the following limitations are introduced:

-> **Note** TOML file size is limited to 10MB.

-> **Note** Only TOML file with restricted privileges can be read. Any privileges for group or others cannot be set (the maximum valid privilege is `700`). You can set the expected privileges like `chmod 0600 ~/.snowflake/config`. This is checked only on non-Windows platforms. If you are using the provider on Windows, please make sure that your configuration file has not too permissive privileges.

### Source Hierarchy
Not all fields must be configured in one source; users can choose which fields are configured in which source.
Provider uses an established hierarchy of sources. The current behavior is that for each field:
1. Check if it is present in the provider configuration. If yes, use this value. If not, go to step 2.
1. Check if it is present in the environment variables. If yes, use this value. If not, go to step 3.
1. Check if it is present in the TOML config file (specifically, use the profile name configured in one of the steps above). If yes, use this value. If not, the value is considered empty.

-> **Note** Currently `private_key` and `private_key_passphrase` are coupled and must be set in one source (both on Terraform side or both in TOML config, see https://github.com/snowflakedb/terraform-provider-snowflake/issues/3332). This will be fixed in the future.

-> **Note** Currently both legacy and new formats are supported. The new format can be enabled with setting `use_legacy_toml_file = false` in the provider configuration. We encourage using the new format for now, as it will be a default one in v2 version of the provider. The differences between these formats are:
- The keys in the provider contain an underscore (`_`) as a separator, but the TOML schema has fields without any separator.
- The field `driver_tracing` in the provider is related to `tracing` in the TOML schema.

### Examples

An example new TOML file contents:

```toml
[example]
account_name = 'account_name'
organization_name = 'organization_name'
user = 'user'
password = 'password'
warehouse = 'SNOWFLAKE'
role = 'ACCOUNTADMIN'
client_ip = '1.2.3.4'
protocol = 'https'
port = 443
okta_url = 'https://example.com'
client_timeout = 10
jwt_client_timeout = 20
login_timeout = 30
request_timeout = 40
jwt_expire_timeout = 50
external_browser_timeout = 60
max_retry_count = 1
authenticator = 'snowflake'
insecure_mode = true
ocsp_fail_open = true
keep_session_alive = true
disable_telemetry = true
validate_default_parameters = true
client_request_mfa_token = true
client_store_temporary_credential = true
driver_tracing = 'info'
tmp_dir_path = '/tmp/terraform-provider/'
disable_query_context_cache = true
include_retry_reason = true
disable_console_login = true
oauth_client_id = 'oauth_client_id'
oauth_client_secret = 'oauth_client_secret'
oauth_token_request_url = 'oauth_token_request_url'
oauth_authorization_url = 'oauth_authorization_url'
oauth_redirect_uri = 'oauth_redirect_uri'
oauth_scope = 'oauth_scope'
enable_single_use_refresh_tokens = true

[example.params]
param_key = 'param_value'
```

An example legacy TOML file contents:

```toml
[example]
accountname = 'account_name'
organizationname = 'organization_name'
user = 'user'
password = 'password'
warehouse = 'SNOWFLAKE'
role = 'ACCOUNTADMIN'
clientip = '1.2.3.4'
protocol = 'https'
port = 443
oktaurl = 'https://example.com'
clienttimeout = 10
jwtclienttimeout = 20
logintimeout = 30
requesttimeout = 40
jwtexpiretimeout = 50
externalbrowsertimeout = 60
maxretrycount = 1
authenticator = 'snowflake'
insecuremode = true
ocspfailopen = true
keepsessionalive = true
disabletelemetry = true
validatedefaultparameters = true
clientrequestmfatoken = true
clientstoretemporarycredential = true
tracing = 'info'
tmpdirpath = '/tmp/terraform-provider/'
disablequerycontextcache = true
includeretryreason = true
disableconsolelogin = true
oauthclientid = 'oauth_client_id'
oauthclientsecret = 'oauth_client_secret'
oauthtokenrequesturl = 'oauth_token_request_url'
oauthauthorizationurl = 'oauth_authorization_url'
oauthredirecturi = 'oauth_redirect_uri'
oauthscope = 'oauth_scope'
enablesingleuserefreshtokens = true

[example.params]
param_key = 'param_value'
```

An example terraform configuration file equivalent:

```terraform
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
```

<!-- Section of deprecated resources -->

<!-- Section of deprecated data sources -->

## Sensitive values limitations

The provider marks fields containing access credentials and other such information as sensitive. This means that the values of these fields will not be logged.

There are some limitations to this mechanism:
- Sensitive values are stored as plaintext in the state file. This is a limitation of Terraform itself ([reference](https://developer.hashicorp.com/terraform/tutorials/configuration-language/sensitive-variables#sensitive-values-in-state)). You should take care to secure access to the state file.
- In [Plugin SDK](https://developer.hashicorp.com/terraform/plugin/sdkv2) there is no possibility to mark sensitive values conditionally ([reference](https://github.com/hashicorp/terraform-plugin-sdk/issues/736)). This means it is not possible to mark sensitive values based on other fields, like marking `body` based on the value of `secure` field in views, functions, and procedures. As a result, this field is not marked as sensitive. For such cases, we add disclaimers in the resource documentation.
- In Plugin SDK, there is no possibility to mark sensitive values in nested fields ([reference](https://github.com/hashicorp/terraform-plugin-sdk/issues/201)). This means the nested fields, like these in `show_output` and `describe_output` cannot be sensitive.
As a result, such nested fields are not marked as sensitive. For such cases, we add disclaimers in the resource documentation. Additionally, some fields are missing from `show_output` and `describe_output`. However, these fields are present in the resource's root, so they can still be referenced.
The alternative solution we considered was setting the whole `show_output` and `describe_output` as sensitive. However, this solution could reduce the provider functionality and would require changes in user's configurations.

As a general rule, please ensure that no personal data, sensitive data, export-controlled data, or other regulated data is entered as metadata when using the provider. If you use one of these fields, they may be present in logs, so ensure that the provider logs are properly restricted. For more information, see [Sensitive values limitations](../#sensitive-values-limitations) and [Metadata fields in Snowflake](https://docs.snowflake.com/en/sql-reference/metadata).

Read more about sensitive values in the [Terraform documentation](https://developer.hashicorp.com/terraform/tutorials/configuration-language/sensitive-variables).

We are planning to research migration to Plugin Framework and we will investigate if the limitations coming from Plugin SDK can be addressed.

## Features

### Operation Timeouts
By default, Terraform sets resource operation timeouts to 20 minutes ([reference](https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/retries-and-customizable-timeouts#default-timeouts-and-deadline-exceeded-errors)). Now, the provider enables configuration of these values by users in `timeouts` block in each resource.
The default timeouts are in general aligned with the Terraform defaults. If a resource has different timeouts, it is specified in the resource documentation.

Data sources will be supported in the future.
Read more in following [official documentation](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts)).

You can specify the timeouts like the following:
```terraform

resource "snowflake_execute" "test" {
  execute = "CREATE DATABASE ABC"
  revert  = "DROP DATABASE ABC"
  query   = "SHOW DATABASES LIKE '%ABC%'"

  timeouts {
    create = "10m"
    read   = "10m"
    update = "10m"
    delete = "10m"
  }
}
```

-> Note: Timeouts can be also set at driver's level (see [driver documentation](https://pkg.go.dev/github.com/snowflakedb/gosnowflake)). These timeouts are independent. We recommend tweaking the timeouts on Terraform level first.

## A list of preview and stable resources and data sources

The provider supports a number of Snowflake features. Within the provider, some features are stable, while others are in preview
(stability of the feature in the provider is not connected to the stability of the feature in Snowflake).

Preview features are **experimental** and may introduce **breaking changes**, even between non-major versions of the provider.
Eventually, every preview resource will be promoted to stable, but the timeline for each feature is not defined (you can find more details on the current/future plans in [our roadmap](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/ROADMAP.md)).
New resources will be introduced as preview ones and promoted over time to stable as we gain more confidence in their stability.

Preview features are disabled by default and should be used with caution.
To use them, add the relevant feature name to the `preview_features_enabled` field in the [provider configuration](#preview_features_enabled-1).

<!-- Section of stable resources -->
## Currently stable resources 

- [snowflake_account](./docs/resources/account)
- [snowflake_account_parameter](./docs/resources/account_parameter)
- [snowflake_account_role](./docs/resources/account_role)
- [snowflake_api_authentication_integration_with_authorization_code_grant](./docs/resources/api_authentication_integration_with_authorization_code_grant)
- [snowflake_api_authentication_integration_with_client_credentials](./docs/resources/api_authentication_integration_with_client_credentials)
- [snowflake_api_authentication_integration_with_jwt_bearer](./docs/resources/api_authentication_integration_with_jwt_bearer)
- [snowflake_database](./docs/resources/database)
- [snowflake_database_role](./docs/resources/database_role)
- [snowflake_execute](./docs/resources/execute)
- [snowflake_external_oauth_integration](./docs/resources/external_oauth_integration)
- [snowflake_grant_account_role](./docs/resources/grant_account_role)
- [snowflake_grant_application_role](./docs/resources/grant_application_role)
- [snowflake_grant_database_role](./docs/resources/grant_database_role)
- [snowflake_grant_ownership](./docs/resources/grant_ownership)
- [snowflake_grant_privileges_to_account_role](./docs/resources/grant_privileges_to_account_role)
- [snowflake_grant_privileges_to_database_role](./docs/resources/grant_privileges_to_database_role)
- [snowflake_grant_privileges_to_share](./docs/resources/grant_privileges_to_share)
- [snowflake_legacy_service_user](./docs/resources/legacy_service_user)
- [snowflake_masking_policy](./docs/resources/masking_policy)
- [snowflake_network_policy](./docs/resources/network_policy)
- [snowflake_oauth_integration_for_custom_clients](./docs/resources/oauth_integration_for_custom_clients)
- [snowflake_oauth_integration_for_partner_applications](./docs/resources/oauth_integration_for_partner_applications)
- [snowflake_primary_connection](./docs/resources/primary_connection)
- [snowflake_resource_monitor](./docs/resources/resource_monitor)
- [snowflake_row_access_policy](./docs/resources/row_access_policy)
- [snowflake_saml2_integration](./docs/resources/saml2_integration)
- [snowflake_schema](./docs/resources/schema)
- [snowflake_scim_integration](./docs/resources/scim_integration)
- [snowflake_secondary_connection](./docs/resources/secondary_connection)
- [snowflake_secondary_database](./docs/resources/secondary_database)
- [snowflake_secret_with_authorization_code_grant](./docs/resources/secret_with_authorization_code_grant)
- [snowflake_secret_with_basic_authentication](./docs/resources/secret_with_basic_authentication)
- [snowflake_secret_with_client_credentials](./docs/resources/secret_with_client_credentials)
- [snowflake_secret_with_generic_string](./docs/resources/secret_with_generic_string)
- [snowflake_service_user](./docs/resources/service_user)
- [snowflake_shared_database](./docs/resources/shared_database)
- [snowflake_stream_on_directory_table](./docs/resources/stream_on_directory_table)
- [snowflake_stream_on_external_table](./docs/resources/stream_on_external_table)
- [snowflake_stream_on_table](./docs/resources/stream_on_table)
- [snowflake_stream_on_view](./docs/resources/stream_on_view)
- [snowflake_streamlit](./docs/resources/streamlit)
- [snowflake_tag](./docs/resources/tag)
- [snowflake_tag_association](./docs/resources/tag_association)
- [snowflake_task](./docs/resources/task)
- [snowflake_user](./docs/resources/user)
- [snowflake_view](./docs/resources/view)
- [snowflake_warehouse](./docs/resources/warehouse)

<!-- Section of stable data sources -->
## Currently stable data sources 

- [snowflake_account_roles](./docs/data-sources/account_roles)
- [snowflake_accounts](./docs/data-sources/accounts)
- [snowflake_connections](./docs/data-sources/connections)
- [snowflake_database_roles](./docs/data-sources/database_roles)
- [snowflake_databases](./docs/data-sources/databases)
- [snowflake_grants](./docs/data-sources/grants)
- [snowflake_masking_policies](./docs/data-sources/masking_policies)
- [snowflake_network_policies](./docs/data-sources/network_policies)
- [snowflake_resource_monitors](./docs/data-sources/resource_monitors)
- [snowflake_row_access_policies](./docs/data-sources/row_access_policies)
- [snowflake_schemas](./docs/data-sources/schemas)
- [snowflake_secrets](./docs/data-sources/secrets)
- [snowflake_security_integrations](./docs/data-sources/security_integrations)
- [snowflake_streamlits](./docs/data-sources/streamlits)
- [snowflake_streams](./docs/data-sources/streams)
- [snowflake_tags](./docs/data-sources/tags)
- [snowflake_tasks](./docs/data-sources/tasks)
- [snowflake_users](./docs/data-sources/users)
- [snowflake_views](./docs/data-sources/views)
- [snowflake_warehouses](./docs/data-sources/warehouses)

<!-- Section of preview resources -->
## Currently preview resources 

- [snowflake_account_authentication_policy_attachment](./docs/resources/account_authentication_policy_attachment)
- [snowflake_account_password_policy_attachment](./docs/resources/account_password_policy_attachment)
- [snowflake_alert](./docs/resources/alert)
- [snowflake_api_integration](./docs/resources/api_integration)
- [snowflake_authentication_policy](./docs/resources/authentication_policy)
- [snowflake_compute_pool](./docs/resources/compute_pool)
- [snowflake_cortex_search_service](./docs/resources/cortex_search_service)
- [snowflake_current_account](./docs/resources/current_account)
- [snowflake_current_organization_account](./docs/resources/current_organization_account)
- [snowflake_dynamic_table](./docs/resources/dynamic_table)
- [snowflake_email_notification_integration](./docs/resources/email_notification_integration)
- [snowflake_external_function](./docs/resources/external_function)
- [snowflake_external_table](./docs/resources/external_table)
- [snowflake_external_volume](./docs/resources/external_volume)
- [snowflake_failover_group](./docs/resources/failover_group)
- [snowflake_file_format](./docs/resources/file_format)
- [snowflake_function_java](./docs/resources/function_java)
- [snowflake_function_javascript](./docs/resources/function_javascript)
- [snowflake_function_python](./docs/resources/function_python)
- [snowflake_function_scala](./docs/resources/function_scala)
- [snowflake_function_sql](./docs/resources/function_sql)
- [snowflake_git_repository](./docs/resources/git_repository)
- [snowflake_image_repository](./docs/resources/image_repository)
- [snowflake_job_service](./docs/resources/job_service)
- [snowflake_listing](./docs/resources/listing)
- [snowflake_managed_account](./docs/resources/managed_account)
- [snowflake_materialized_view](./docs/resources/materialized_view)
- [snowflake_network_policy_attachment](./docs/resources/network_policy_attachment)
- [snowflake_network_rule](./docs/resources/network_rule)
- [snowflake_notification_integration](./docs/resources/notification_integration)
- [snowflake_object_parameter](./docs/resources/object_parameter)
- [snowflake_password_policy](./docs/resources/password_policy)
- [snowflake_pipe](./docs/resources/pipe)
- [snowflake_procedure_java](./docs/resources/procedure_java)
- [snowflake_procedure_javascript](./docs/resources/procedure_javascript)
- [snowflake_procedure_python](./docs/resources/procedure_python)
- [snowflake_procedure_scala](./docs/resources/procedure_scala)
- [snowflake_procedure_sql](./docs/resources/procedure_sql)
- [snowflake_sequence](./docs/resources/sequence)
- [snowflake_service](./docs/resources/service)
- [snowflake_share](./docs/resources/share)
- [snowflake_stage](./docs/resources/stage)
- [snowflake_storage_integration](./docs/resources/storage_integration)
- [snowflake_table](./docs/resources/table)
- [snowflake_table_column_masking_policy_application](./docs/resources/table_column_masking_policy_application)
- [snowflake_table_constraint](./docs/resources/table_constraint)
- [snowflake_user_authentication_policy_attachment](./docs/resources/user_authentication_policy_attachment)
- [snowflake_user_password_policy_attachment](./docs/resources/user_password_policy_attachment)
- [snowflake_user_programmatic_access_token](./docs/resources/user_programmatic_access_token)
- [snowflake_user_public_keys](./docs/resources/user_public_keys)

<!-- Section of preview data sources -->
## Currently preview data sources 

- [snowflake_alerts](./docs/data-sources/alerts)
- [snowflake_compute_pools](./docs/data-sources/compute_pools)
- [snowflake_cortex_search_services](./docs/data-sources/cortex_search_services)
- [snowflake_current_account](./docs/data-sources/current_account)
- [snowflake_current_role](./docs/data-sources/current_role)
- [snowflake_database](./docs/data-sources/database)
- [snowflake_database_role](./docs/data-sources/database_role)
- [snowflake_dynamic_tables](./docs/data-sources/dynamic_tables)
- [snowflake_external_functions](./docs/data-sources/external_functions)
- [snowflake_external_tables](./docs/data-sources/external_tables)
- [snowflake_failover_groups](./docs/data-sources/failover_groups)
- [snowflake_file_formats](./docs/data-sources/file_formats)
- [snowflake_functions](./docs/data-sources/functions)
- [snowflake_git_repositories](./docs/data-sources/git_repositories)
- [snowflake_image_repositories](./docs/data-sources/image_repositories)
- [snowflake_materialized_views](./docs/data-sources/materialized_views)
- [snowflake_parameters](./docs/data-sources/parameters)
- [snowflake_pipes](./docs/data-sources/pipes)
- [snowflake_procedures](./docs/data-sources/procedures)
- [snowflake_sequences](./docs/data-sources/sequences)
- [snowflake_services](./docs/data-sources/services)
- [snowflake_shares](./docs/data-sources/shares)
- [snowflake_stages](./docs/data-sources/stages)
- [snowflake_storage_integrations](./docs/data-sources/storage_integrations)
- [snowflake_system_generate_scim_access_token](./docs/data-sources/system_generate_scim_access_token)
- [snowflake_system_get_aws_sns_iam_policy](./docs/data-sources/system_get_aws_sns_iam_policy)
- [snowflake_system_get_privatelink_config](./docs/data-sources/system_get_privatelink_config)
- [snowflake_system_get_snowflake_platform_info](./docs/data-sources/system_get_snowflake_platform_info)
- [snowflake_tables](./docs/data-sources/tables)
- [snowflake_user_programmatic_access_tokens](./docs/data-sources/user_programmatic_access_tokens)
