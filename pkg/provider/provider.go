package provider

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/url"
	"slices"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/datasources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider/docs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider/validators"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/snowflakedb/gosnowflake"
)

func init() {
	// useful links:
	// - https://github.com/hashicorp/terraform-plugin-docs/issues/10#issuecomment-767682837
	// - https://github.com/hashicorp/terraform-plugin-docs/issues/156#issuecomment-1600427216
	// TODO(SNOW-1901053): Rework handling deprecated objects' descriptions.
	schema.ResourceDescriptionBuilder = func(r *schema.Resource) string {
		desc := r.Description
		if r.DeprecationMessage != "" {
			deprecationMessage := r.DeprecationMessage
			replacement, path, ok := docs.GetDeprecatedResourceReplacement(deprecationMessage)
			if ok {
				deprecationMessage = strings.ReplaceAll(deprecationMessage, replacement, docs.RelativeLink(replacement, path))
			}
			// <deprecation> tag is a hack to split description into two parts (deprecation/real description) nicely. This tag won't be rendered.
			// Check resources.md.tmpl for usage example.
			desc = fmt.Sprintf("~> **Deprecation** %v <deprecation>\n\n%s", deprecationMessage, r.Description)
		}
		return strings.TrimSpace(desc)
	}

	schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
		desc := s.Description
		if s.Default != nil {
			if slices.Contains([]any{
				provider.IntDefault,
				provider.BooleanDefault,
			}, s.Default) {
				desc = fmt.Sprintf("(Default: fallback to Snowflake default - uses special value that cannot be set in the configuration manually (`%v`)) %s", s.Default, s.Description)
			} else {
				desc = fmt.Sprintf("(Default: `%v`) %s", s.Default, s.Description)
			}
		}
		return desc
	}
}

// Provider returns a Terraform Provider using configuration. It is based on https://pkg.go.dev/github.com/snowflakedb/gosnowflake#Config.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"account_name": {
				Type:         schema.TypeString,
				Description:  envNameFieldDescription("Specifies your Snowflake account name assigned by Snowflake. For information about account identifiers, see the [Snowflake documentation](https://docs.snowflake.com/en/user-guide/admin-account-identifier#account-name). Required unless using `profile`.", snowflakeenvs.AccountName),
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc(snowflakeenvs.AccountName, nil),
				RequiredWith: []string{"account_name", "organization_name"},
			},
			"organization_name": {
				Type:         schema.TypeString,
				Description:  envNameFieldDescription("Specifies your Snowflake organization name assigned by Snowflake. For information about account identifiers, see the [Snowflake documentation](https://docs.snowflake.com/en/user-guide/admin-account-identifier#organization-name). Required unless using `profile`.", snowflakeenvs.OrganizationName),
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc(snowflakeenvs.OrganizationName, nil),
				RequiredWith: []string{"account_name", "organization_name"},
			},
			"user": {
				Type:             schema.TypeString,
				Description:      envNameFieldDescription("Username. Required unless using `profile`.", snowflakeenvs.User),
				Optional:         true,
				DefaultFunc:      schema.EnvDefaultFunc(snowflakeenvs.User, nil),
				ValidateDiagFunc: validators.IsValidIdentifier[sdk.AccountObjectIdentifier](),
			},
			"password": {
				Type:          schema.TypeString,
				Description:   envNameFieldDescription("Password for user + password auth. Cannot be used with `private_key` and `private_key_passphrase`.", snowflakeenvs.Password),
				Optional:      true,
				Sensitive:     true,
				DefaultFunc:   schema.EnvDefaultFunc(snowflakeenvs.Password, nil),
				ConflictsWith: []string{"private_key", "private_key_passphrase"},
			},
			"warehouse": {
				Type:             schema.TypeString,
				Description:      envNameFieldDescription("Specifies the virtual warehouse to use by default for queries, loading, etc. in the client session.", snowflakeenvs.Warehouse),
				Optional:         true,
				DefaultFunc:      schema.EnvDefaultFunc(snowflakeenvs.Warehouse, nil),
				ValidateDiagFunc: validators.IsValidIdentifier[sdk.AccountObjectIdentifier](),
			},
			"role": {
				Type:             schema.TypeString,
				Description:      envNameFieldDescription("Specifies the role to use by default for accessing Snowflake objects in the client session.", snowflakeenvs.Role),
				Optional:         true,
				ValidateDiagFunc: validators.IsValidIdentifier[sdk.AccountObjectIdentifier](),
				DefaultFunc:      schema.EnvDefaultFunc(snowflakeenvs.Role, nil),
			},
			"validate_default_parameters": {
				Type:             schema.TypeString,
				Description:      envNameFieldDescription("True by default. If false, disables the validation checks for Database, Schema, Warehouse and Role at the time a connection is established.", snowflakeenvs.ValidateDefaultParameters),
				Optional:         true,
				ValidateDiagFunc: validators.ValidateBooleanStringWithDefault,
				DefaultFunc:      schema.EnvDefaultFunc(snowflakeenvs.ValidateDefaultParameters, provider.BooleanDefault),
			},
			// TODO(SNOW-999056): optionally rename to session_params
			"params": {
				Type:        schema.TypeMap,
				Description: "Sets other connection (i.e. session) parameters. [Parameters](https://docs.snowflake.com/en/sql-reference/parameters). This field can not be set with environmental variables.",
				Optional:    true,
			},
			"client_ip": {
				Type:             schema.TypeString,
				Description:      envNameFieldDescription("IP address for network checks.", snowflakeenvs.ClientIp),
				Optional:         true,
				DefaultFunc:      schema.EnvDefaultFunc(snowflakeenvs.ClientIp, nil),
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsIPAddress),
			},
			"protocol": {
				Type:             schema.TypeString,
				Description:      envNameFieldDescription(fmt.Sprintf("A protocol used in the connection. Valid options are: %v.", docs.PossibleValuesListed(allProtocols)), snowflakeenvs.Protocol),
				Optional:         true,
				DefaultFunc:      schema.EnvDefaultFunc(snowflakeenvs.Protocol, nil),
				ValidateDiagFunc: validators.NormalizeValidation(toProtocol),
			},
			"host": {
				Type:        schema.TypeString,
				Description: envNameFieldDescription("Specifies a custom host value used by the driver for privatelink connections.", snowflakeenvs.Host),
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(snowflakeenvs.Host, nil),
			},
			"port": {
				Type:             schema.TypeInt,
				Description:      envNameFieldDescription("Specifies a custom port value used by the driver for privatelink connections.", snowflakeenvs.Port),
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsPortNumberOrZero),
				DefaultFunc:      schema.EnvDefaultFunc(snowflakeenvs.Port, nil),
			},
			"authenticator": {
				Type:             schema.TypeString,
				Description:      envNameFieldDescription(fmt.Sprintf("Specifies the [authentication type](https://pkg.go.dev/github.com/snowflakedb/gosnowflake#AuthType) to use when connecting to Snowflake. Valid options are: %v.", docs.PossibleValuesListed(sdk.AllAuthenticationTypes)), snowflakeenvs.Authenticator),
				Optional:         true,
				DefaultFunc:      schema.EnvDefaultFunc(snowflakeenvs.Authenticator, string(sdk.AuthenticationTypeEmpty)),
				ValidateDiagFunc: validators.NormalizeValidation(sdk.ToExtendedAuthenticatorType),
			},
			"passcode": {
				Type:          schema.TypeString,
				Description:   envNameFieldDescription("Specifies the passcode provided by Duo when using multi-factor authentication (MFA) for login.", snowflakeenvs.Passcode),
				Optional:      true,
				Sensitive:     true,
				ConflictsWith: []string{"passcode_in_password"},
				DefaultFunc:   schema.EnvDefaultFunc(snowflakeenvs.Passcode, nil),
			},
			"passcode_in_password": {
				Type:          schema.TypeBool,
				Description:   envNameFieldDescription("False by default. Set to true if the MFA passcode is embedded to the configured password.", snowflakeenvs.PasscodeInPassword),
				Optional:      true,
				ConflictsWith: []string{"passcode"},
				DefaultFunc:   schema.EnvDefaultFunc(snowflakeenvs.PasscodeInPassword, nil),
			},
			"okta_url": {
				Type:             schema.TypeString,
				Description:      envNameFieldDescription("The URL of the Okta server. e.g. https://example.okta.com. Okta URL host needs to to have a suffix `okta.com`. Read more in Snowflake [docs](https://docs.snowflake.com/en/user-guide/oauth-okta).", snowflakeenvs.OktaUrl),
				Optional:         true,
				DefaultFunc:      schema.EnvDefaultFunc(snowflakeenvs.OktaUrl, nil),
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsURLWithHTTPorHTTPS),
			},
			"login_timeout": {
				Type:             schema.TypeInt,
				Description:      envNameFieldDescription("Login retry timeout in seconds EXCLUDING network roundtrip and read out http response.", snowflakeenvs.LoginTimeout),
				Optional:         true,
				DefaultFunc:      schema.EnvDefaultFunc(snowflakeenvs.LoginTimeout, nil),
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(0)),
			},
			"request_timeout": {
				Type:             schema.TypeInt,
				Description:      envNameFieldDescription("request retry timeout in seconds EXCLUDING network roundtrip and read out http response.", snowflakeenvs.RequestTimeout),
				Optional:         true,
				DefaultFunc:      schema.EnvDefaultFunc(snowflakeenvs.RequestTimeout, nil),
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(0)),
			},
			"jwt_expire_timeout": {
				Type:             schema.TypeInt,
				Description:      envNameFieldDescription("JWT expire after timeout in seconds.", snowflakeenvs.JwtExpireTimeout),
				Optional:         true,
				DefaultFunc:      schema.EnvDefaultFunc(snowflakeenvs.JwtExpireTimeout, nil),
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(0)),
			},
			"client_timeout": {
				Type:             schema.TypeInt,
				Description:      envNameFieldDescription("The timeout in seconds for the client to complete the authentication.", snowflakeenvs.ClientTimeout),
				Optional:         true,
				DefaultFunc:      schema.EnvDefaultFunc(snowflakeenvs.ClientTimeout, nil),
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(0)),
			},
			"jwt_client_timeout": {
				Type:             schema.TypeInt,
				Description:      envNameFieldDescription("The timeout in seconds for the JWT client to complete the authentication.", snowflakeenvs.JwtClientTimeout),
				Optional:         true,
				DefaultFunc:      schema.EnvDefaultFunc(snowflakeenvs.JwtClientTimeout, nil),
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(0)),
			},
			"external_browser_timeout": {
				Type:             schema.TypeInt,
				Description:      envNameFieldDescription("The timeout in seconds for the external browser to complete the authentication.", snowflakeenvs.ExternalBrowserTimeout),
				Optional:         true,
				DefaultFunc:      schema.EnvDefaultFunc(snowflakeenvs.ExternalBrowserTimeout, nil),
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(0)),
			},
			"insecure_mode": {
				Type:        schema.TypeBool,
				Description: envNameFieldDescription("If true, bypass the Online Certificate Status Protocol (OCSP) certificate revocation check. IMPORTANT: Change the default value for testing or emergency situations only.", snowflakeenvs.InsecureMode),
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(snowflakeenvs.InsecureMode, nil),
			},
			"ocsp_fail_open": {
				Type:             schema.TypeString,
				Description:      envNameFieldDescription("True represents OCSP fail open mode. False represents OCSP fail closed mode. Fail open true by default.", snowflakeenvs.OcspFailOpen),
				Optional:         true,
				DefaultFunc:      schema.EnvDefaultFunc(snowflakeenvs.OcspFailOpen, provider.BooleanDefault),
				ValidateDiagFunc: validators.ValidateBooleanStringWithDefault,
			},
			"token": {
				Type:        schema.TypeString,
				Description: envNameFieldDescription("Token to use for OAuth and other forms of token based auth.", snowflakeenvs.Token),
				Sensitive:   true,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(snowflakeenvs.Token, nil),
			},
			"token_accessor": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"token_endpoint": {
							Type:             schema.TypeString,
							Description:      envNameFieldDescription("The token endpoint for the OAuth provider e.g. https://{yourDomain}/oauth/token when using a refresh token to renew access token.", snowflakeenvs.TokenAccessorTokenEndpoint),
							Required:         true,
							Sensitive:        true,
							DefaultFunc:      schema.EnvDefaultFunc(snowflakeenvs.TokenAccessorTokenEndpoint, nil),
							ValidateDiagFunc: validation.ToDiagFunc(validation.IsURLWithHTTPorHTTPS),
						},
						"refresh_token": {
							Type:        schema.TypeString,
							Description: envNameFieldDescription("The refresh token for the OAuth provider when using a refresh token to renew access token.", snowflakeenvs.TokenAccessorRefreshToken),
							Required:    true,
							Sensitive:   true,
							DefaultFunc: schema.EnvDefaultFunc(snowflakeenvs.TokenAccessorRefreshToken, nil),
						},
						"client_id": {
							Type:        schema.TypeString,
							Description: envNameFieldDescription("The client ID for the OAuth provider when using a refresh token to renew access token.", snowflakeenvs.TokenAccessorClientId),
							Required:    true,
							Sensitive:   true,
							DefaultFunc: schema.EnvDefaultFunc(snowflakeenvs.TokenAccessorClientId, nil),
						},
						"client_secret": {
							Type:        schema.TypeString,
							Description: envNameFieldDescription("The client secret for the OAuth provider when using a refresh token to renew access token.", snowflakeenvs.TokenAccessorClientSecret),
							Required:    true,
							Sensitive:   true,
							DefaultFunc: schema.EnvDefaultFunc(snowflakeenvs.TokenAccessorClientSecret, nil),
						},
						"redirect_uri": {
							Type:        schema.TypeString,
							Description: envNameFieldDescription("The redirect URI for the OAuth provider when using a refresh token to renew access token.", snowflakeenvs.TokenAccessorRedirectUri),
							Required:    true,
							Sensitive:   true,
							DefaultFunc: schema.EnvDefaultFunc(snowflakeenvs.TokenAccessorRedirectUri, nil),
						},
					},
				},
			},
			"keep_session_alive": {
				Type:        schema.TypeBool,
				Description: envNameFieldDescription("Enables the session to persist even after the connection is closed.", snowflakeenvs.KeepSessionAlive),
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(snowflakeenvs.KeepSessionAlive, nil),
			},
			"private_key": {
				Type:          schema.TypeString,
				Description:   envNameFieldDescription("Private Key for username+private-key auth. Cannot be used with `password`.", snowflakeenvs.PrivateKey),
				Optional:      true,
				Sensitive:     true,
				DefaultFunc:   schema.EnvDefaultFunc(snowflakeenvs.PrivateKey, nil),
				ConflictsWith: []string{"password"},
			},
			"private_key_passphrase": {
				Type:          schema.TypeString,
				Description:   envNameFieldDescription("Supports the encryption ciphers aes-128-cbc, aes-128-gcm, aes-192-cbc, aes-192-gcm, aes-256-cbc, aes-256-gcm, and des-ede3-cbc.", snowflakeenvs.PrivateKeyPassphrase),
				Optional:      true,
				Sensitive:     true,
				DefaultFunc:   schema.EnvDefaultFunc(snowflakeenvs.PrivateKeyPassphrase, nil),
				ConflictsWith: []string{"password"},
			},
			"disable_telemetry": {
				Type:        schema.TypeBool,
				Description: envNameFieldDescription("Disables telemetry in the driver.", snowflakeenvs.DisableTelemetry),
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(snowflakeenvs.DisableTelemetry, nil),
			},
			"client_request_mfa_token": {
				Type:             schema.TypeString,
				Description:      envNameFieldDescription("When true the MFA token is cached in the credential manager. True by default in Windows/OSX. False for Linux.", snowflakeenvs.ClientRequestMfaToken),
				Optional:         true,
				DefaultFunc:      schema.EnvDefaultFunc(snowflakeenvs.ClientRequestMfaToken, provider.BooleanDefault),
				ValidateDiagFunc: validators.ValidateBooleanStringWithDefault,
			},
			"client_store_temporary_credential": {
				Type:             schema.TypeString,
				Description:      envNameFieldDescription("When true the ID token is cached in the credential manager. True by default in Windows/OSX. False for Linux.", snowflakeenvs.ClientStoreTemporaryCredential),
				Optional:         true,
				DefaultFunc:      schema.EnvDefaultFunc(snowflakeenvs.ClientStoreTemporaryCredential, provider.BooleanDefault),
				ValidateDiagFunc: validators.ValidateBooleanStringWithDefault,
			},
			"disable_query_context_cache": {
				Type:        schema.TypeBool,
				Description: envNameFieldDescription("Disables HTAP query context cache in the driver.", snowflakeenvs.DisableQueryContextCache),
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(snowflakeenvs.DisableQueryContextCache, nil),
			},
			"include_retry_reason": {
				Type:             schema.TypeString,
				Description:      envNameFieldDescription("Should retried request contain retry reason.", snowflakeenvs.IncludeRetryReason),
				Optional:         true,
				DefaultFunc:      schema.EnvDefaultFunc(snowflakeenvs.IncludeRetryReason, resources.BooleanDefault),
				ValidateDiagFunc: validators.ValidateBooleanStringWithDefault,
			},
			"max_retry_count": {
				Type:             schema.TypeInt,
				Description:      envNameFieldDescription("Specifies how many times non-periodic HTTP request can be retried by the driver.", snowflakeenvs.MaxRetryCount),
				Optional:         true,
				DefaultFunc:      schema.EnvDefaultFunc(snowflakeenvs.MaxRetryCount, nil),
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(0)),
			},
			"driver_tracing": {
				Type:             schema.TypeString,
				Description:      envNameFieldDescription(fmt.Sprintf("Specifies the logging level to be used by the driver. Valid options are: %v.", docs.PossibleValuesListed(sdk.AllDriverLogLevels)), snowflakeenvs.DriverTracing),
				Optional:         true,
				DefaultFunc:      schema.EnvDefaultFunc(snowflakeenvs.DriverTracing, nil),
				ValidateDiagFunc: validators.NormalizeValidation(sdk.ToDriverLogLevel),
			},
			"tmp_directory_path": {
				Type:        schema.TypeString,
				Description: envNameFieldDescription("Sets temporary directory used by the driver for operations like encrypting, compressing etc.", snowflakeenvs.TmpDirectoryPath),
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(snowflakeenvs.TmpDirectoryPath, nil),
			},
			"disable_console_login": {
				Type:             schema.TypeString,
				Description:      envNameFieldDescription("Indicates whether console login should be disabled in the driver.", snowflakeenvs.DisableConsoleLogin),
				Optional:         true,
				DefaultFunc:      schema.EnvDefaultFunc(snowflakeenvs.DisableConsoleLogin, resources.BooleanDefault),
				ValidateDiagFunc: validators.ValidateBooleanStringWithDefault,
			},
			// TODO(SNOW-1761318): Add DisableSamlURLCheck.
			// TODO(SNOW-1917271): Add DisableOCSPChecks.
			"profile": {
				Type: schema.TypeString,
				// TODO(SNOW-1754364): Note that a default file path is already filled on sdk side.
				Description: envNameFieldDescription("Sets the profile to read from ~/.snowflake/config file.", snowflakeenvs.Profile),
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(snowflakeenvs.Profile, "default"),
			},
			"preview_features_enabled": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: validators.StringInSlice(previewfeatures.AllPreviewFeatures, true),
				},
				Description: fmt.Sprintf("A list of preview features that are handled by the provider. See [preview features list](https://github.com/Snowflake-Labs/terraform-provider-snowflake/blob/main/v1-preparations/LIST_OF_PREVIEW_FEATURES_FOR_V1.md). Preview features may have breaking changes in future releases, even without raising the major version. This field can not be set with environmental variables. Valid options are: %v.", docs.PossibleValuesListed(previewfeatures.AllPreviewFeatures)),
			},
			"skip_toml_file_permission_verification": {
				Type:        schema.TypeBool,
				Description: envNameFieldDescription("False by default. Skips TOML configuration file permission verification. This flag has no effect on Windows systems, as the permissions are not checked on this platform. Instead of skipping the permissions verification, we recommend setting the proper privileges - see [the section below](#toml-file-limitations).", snowflakeenvs.SkipTomlFilePermissionVerification),
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(snowflakeenvs.SkipTomlFilePermissionVerification, false),
			},
			"use_legacy_toml_file": {
				Type:        schema.TypeBool,
				Description: envNameFieldDescription("False by default. When this is set to true, the provider expects the legacy TOML format. Otherwise, it expects the new format. See more in [the section below](#examples)", snowflakeenvs.UseLegacyTomlFile),
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(snowflakeenvs.UseLegacyTomlFile, false),
			},
		},
		ResourcesMap:         getResources(),
		DataSourcesMap:       getDataSources(),
		ConfigureContextFunc: ConfigureProvider,
		ProviderMetaSchema:   map[string]*schema.Schema{},
	}
}

func getResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"snowflake_account": resources.Account(),
		"snowflake_account_authentication_policy_attachment":                     resources.AccountAuthenticationPolicyAttachment(),
		"snowflake_account_role":                                                 resources.AccountRole(),
		"snowflake_account_password_policy_attachment":                           resources.AccountPasswordPolicyAttachment(),
		"snowflake_account_parameter":                                            resources.AccountParameter(),
		"snowflake_alert":                                                        resources.Alert(),
		"snowflake_api_authentication_integration_with_authorization_code_grant": resources.ApiAuthenticationIntegrationWithAuthorizationCodeGrant(),
		"snowflake_api_authentication_integration_with_client_credentials":       resources.ApiAuthenticationIntegrationWithClientCredentials(),
		"snowflake_api_authentication_integration_with_jwt_bearer":               resources.ApiAuthenticationIntegrationWithJwtBearer(),
		"snowflake_api_integration":                                              resources.APIIntegration(),
		"snowflake_authentication_policy":                                        resources.AuthenticationPolicy(),
		"snowflake_compute_pool":                                                 resources.ComputePool(),
		"snowflake_cortex_search_service":                                        resources.CortexSearchService(),
		"snowflake_database":                                                     resources.Database(),
		"snowflake_database_role":                                                resources.DatabaseRole(),
		"snowflake_dynamic_table":                                                resources.DynamicTable(),
		"snowflake_email_notification_integration":                               resources.EmailNotificationIntegration(),
		"snowflake_execute":                                                      resources.Execute(),
		"snowflake_external_function":                                            resources.ExternalFunction(),
		"snowflake_external_oauth_integration":                                   resources.ExternalOauthIntegration(),
		"snowflake_external_table":                                               resources.ExternalTable(),
		"snowflake_external_volume":                                              resources.ExternalVolume(),
		"snowflake_failover_group":                                               resources.FailoverGroup(),
		"snowflake_file_format":                                                  resources.FileFormat(),
		"snowflake_function_java":                                                resources.FunctionJava(),
		"snowflake_function_javascript":                                          resources.FunctionJavascript(),
		"snowflake_function_python":                                              resources.FunctionPython(),
		"snowflake_function_scala":                                               resources.FunctionScala(),
		"snowflake_function_sql":                                                 resources.FunctionSql(),
		"snowflake_grant_account_role":                                           resources.GrantAccountRole(),
		"snowflake_grant_application_role":                                       resources.GrantApplicationRole(),
		"snowflake_grant_database_role":                                          resources.GrantDatabaseRole(),
		"snowflake_grant_ownership":                                              resources.GrantOwnership(),
		"snowflake_grant_privileges_to_account_role":                             resources.GrantPrivilegesToAccountRole(),
		"snowflake_grant_privileges_to_database_role":                            resources.GrantPrivilegesToDatabaseRole(),
		"snowflake_grant_privileges_to_share":                                    resources.GrantPrivilegesToShare(),
		"snowflake_git_repository":                                               resources.GitRepository(),
		"snowflake_image_repository":                                             resources.ImageRepository(),
		"snowflake_legacy_service_user":                                          resources.LegacyServiceUser(),
		"snowflake_managed_account":                                              resources.ManagedAccount(),
		"snowflake_masking_policy":                                               resources.MaskingPolicy(),
		"snowflake_materialized_view":                                            resources.MaterializedView(),
		"snowflake_network_policy":                                               resources.NetworkPolicy(),
		"snowflake_network_policy_attachment":                                    resources.NetworkPolicyAttachment(),
		"snowflake_network_rule":                                                 resources.NetworkRule(),
		"snowflake_notification_integration":                                     resources.NotificationIntegration(),
		"snowflake_oauth_integration_for_partner_applications":                   resources.OauthIntegrationForPartnerApplications(),
		"snowflake_oauth_integration_for_custom_clients":                         resources.OauthIntegrationForCustomClients(),
		"snowflake_object_parameter":                                             resources.ObjectParameter(),
		"snowflake_password_policy":                                              resources.PasswordPolicy(),
		"snowflake_pipe":                                                         resources.Pipe(),
		"snowflake_primary_connection":                                           resources.PrimaryConnection(),
		"snowflake_procedure_java":                                               resources.ProcedureJava(),
		"snowflake_procedure_javascript":                                         resources.ProcedureJavascript(),
		"snowflake_procedure_python":                                             resources.ProcedurePython(),
		"snowflake_procedure_scala":                                              resources.ProcedureScala(),
		"snowflake_procedure_sql":                                                resources.ProcedureSql(),
		"snowflake_resource_monitor":                                             resources.ResourceMonitor(),
		"snowflake_row_access_policy":                                            resources.RowAccessPolicy(),
		"snowflake_saml2_integration":                                            resources.SAML2Integration(),
		"snowflake_schema":                                                       resources.Schema(),
		"snowflake_scim_integration":                                             resources.SCIMIntegration(),
		"snowflake_secondary_connection":                                         resources.SecondaryConnection(),
		"snowflake_secondary_database":                                           resources.SecondaryDatabase(),
		"snowflake_secret_with_authorization_code_grant":                         resources.SecretWithAuthorizationCodeGrant(),
		"snowflake_secret_with_basic_authentication":                             resources.SecretWithBasicAuthentication(),
		"snowflake_secret_with_client_credentials":                               resources.SecretWithClientCredentials(),
		"snowflake_secret_with_generic_string":                                   resources.SecretWithGenericString(),
		"snowflake_sequence":                                                     resources.Sequence(),
		"snowflake_service_user":                                                 resources.ServiceUser(),
		"snowflake_share":                                                        resources.Share(),
		"snowflake_shared_database":                                              resources.SharedDatabase(),
		"snowflake_stage":                                                        resources.Stage(),
		"snowflake_storage_integration":                                          resources.StorageIntegration(),
		"snowflake_stream_on_directory_table":                                    resources.StreamOnDirectoryTable(),
		"snowflake_stream_on_external_table":                                     resources.StreamOnExternalTable(),
		"snowflake_stream_on_table":                                              resources.StreamOnTable(),
		"snowflake_stream_on_view":                                               resources.StreamOnView(),
		"snowflake_streamlit":                                                    resources.Streamlit(),
		"snowflake_table":                                                        resources.Table(),
		"snowflake_table_column_masking_policy_application":                      resources.TableColumnMaskingPolicyApplication(),
		"snowflake_table_constraint":                                             resources.TableConstraint(),
		"snowflake_tag":                                                          resources.Tag(),
		"snowflake_tag_association":                                              resources.TagAssociation(),
		"snowflake_task":                                                         resources.Task(),
		"snowflake_user":                                                         resources.User(),
		"snowflake_user_authentication_policy_attachment":                        resources.UserAuthenticationPolicyAttachment(),
		"snowflake_user_password_policy_attachment":                              resources.UserPasswordPolicyAttachment(),
		"snowflake_user_public_keys":                                             resources.UserPublicKeys(),
		"snowflake_view":                                                         resources.View(),
		"snowflake_warehouse":                                                    resources.Warehouse(),
	}
}

func getDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"snowflake_accounts":                           datasources.Accounts(),
		"snowflake_account_roles":                      datasources.AccountRoles(),
		"snowflake_alerts":                             datasources.Alerts(),
		"snowflake_compute_pools":                      datasources.ComputePools(),
		"snowflake_connections":                        datasources.Connections(),
		"snowflake_cortex_search_services":             datasources.CortexSearchServices(),
		"snowflake_current_account":                    datasources.CurrentAccount(),
		"snowflake_current_role":                       datasources.CurrentRole(),
		"snowflake_database":                           datasources.Database(),
		"snowflake_database_role":                      datasources.DatabaseRole(),
		"snowflake_database_roles":                     datasources.DatabaseRoles(),
		"snowflake_databases":                          datasources.Databases(),
		"snowflake_dynamic_tables":                     datasources.DynamicTables(),
		"snowflake_external_functions":                 datasources.ExternalFunctions(),
		"snowflake_external_tables":                    datasources.ExternalTables(),
		"snowflake_failover_groups":                    datasources.FailoverGroups(),
		"snowflake_file_formats":                       datasources.FileFormats(),
		"snowflake_functions":                          datasources.Functions(),
		"snowflake_git_repositories":                   datasources.GitRepositories(),
		"snowflake_grants":                             datasources.Grants(),
		"snowflake_image_repositories":                 datasources.ImageRepositories(),
		"snowflake_masking_policies":                   datasources.MaskingPolicies(),
		"snowflake_materialized_views":                 datasources.MaterializedViews(),
		"snowflake_network_policies":                   datasources.NetworkPolicies(),
		"snowflake_parameters":                         datasources.Parameters(),
		"snowflake_pipes":                              datasources.Pipes(),
		"snowflake_procedures":                         datasources.Procedures(),
		"snowflake_resource_monitors":                  datasources.ResourceMonitors(),
		"snowflake_row_access_policies":                datasources.RowAccessPolicies(),
		"snowflake_schemas":                            datasources.Schemas(),
		"snowflake_secrets":                            datasources.Secrets(),
		"snowflake_security_integrations":              datasources.SecurityIntegrations(),
		"snowflake_sequences":                          datasources.Sequences(),
		"snowflake_shares":                             datasources.Shares(),
		"snowflake_stages":                             datasources.Stages(),
		"snowflake_storage_integrations":               datasources.StorageIntegrations(),
		"snowflake_streams":                            datasources.Streams(),
		"snowflake_streamlits":                         datasources.Streamlits(),
		"snowflake_system_generate_scim_access_token":  datasources.SystemGenerateSCIMAccessToken(),
		"snowflake_system_get_aws_sns_iam_policy":      datasources.SystemGetAWSSNSIAMPolicy(),
		"snowflake_system_get_privatelink_config":      datasources.SystemGetPrivateLinkConfig(),
		"snowflake_system_get_snowflake_platform_info": datasources.SystemGetSnowflakePlatformInfo(),
		"snowflake_tables":                             datasources.Tables(),
		"snowflake_tags":                               datasources.Tags(),
		"snowflake_tasks":                              datasources.Tasks(),
		"snowflake_users":                              datasources.Users(),
		"snowflake_views":                              datasources.Views(),
		"snowflake_warehouses":                         datasources.Warehouses(),
	}
}

func ConfigureProvider(ctx context.Context, s *schema.ResourceData) (any, diag.Diagnostics) {
	config, err := getDriverConfigFromTerraform(s)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	var verifyPermissions bool
	if v := s.Get("skip_toml_file_permission_verification"); v.(bool) {
		verifyPermissions = false
	} else {
		verifyPermissions = true
	}
	var useLegacyTomlFile bool
	if v := s.Get("use_legacy_toml_file"); v.(bool) {
		useLegacyTomlFile = true
	} else {
		useLegacyTomlFile = false
	}

	if v, ok := s.GetOk("profile"); ok && v.(string) != "" {
		tomlConfig, err := getDriverConfigFromTOML(v.(string), verifyPermissions, useLegacyTomlFile)
		if err != nil {
			return nil, diag.FromErr(err)
		}
		config = sdk.MergeConfig(config, tomlConfig)
	}

	client, clientErr := sdk.NewClient(config)

	providerCtx := &provider.Context{Client: client}

	if v, ok := s.GetOk("preview_features_enabled"); ok {
		providerCtx.EnabledFeatures = expandStringList(v.(*schema.Set).List())
	}

	if clientErr != nil {
		return nil, diag.FromErr(clientErr)
	}

	return providerCtx, nil
}

// TODO: reuse with the function from resources package
func expandStringList(configured []interface{}) []string {
	vs := make([]string, 0, len(configured))
	for _, v := range configured {
		val, ok := v.(string)
		if ok && val != "" {
			vs = append(vs, val)
		}
	}
	return vs
}

func getDriverConfigFromTOML(profile string, verifyPermissions, useLegacyTomlFile bool) (*gosnowflake.Config, error) {
	if profile == "default" {
		return sdk.DefaultConfig(
			sdk.WithVerifyPermissions(verifyPermissions),
			sdk.WithUseLegacyTomlFormat(useLegacyTomlFile),
		), nil
	}
	path, err := sdk.GetConfigFileName()
	if err != nil {
		return nil, err
	}

	profileConfig, err := sdk.ProfileConfig(
		profile,
		sdk.WithVerifyPermissions(verifyPermissions),
		sdk.WithUseLegacyTomlFormat(useLegacyTomlFile),
	)
	if err != nil {
		return nil, fmt.Errorf(`could not retrieve "%s" profile config from file %s: %w`, profile, path, err)
	}
	if profileConfig == nil {
		return nil, fmt.Errorf(`profile "%s" not found in file %s`, profile, path)
	}
	return profileConfig, nil
}

func getDriverConfigFromTerraform(s *schema.ResourceData) (*gosnowflake.Config, error) {
	config := &gosnowflake.Config{
		Application: "terraform-provider-snowflake",
	}

	err := errors.Join(
		// account_name and organization_name are handled below
		handleStringField(s, "user", &config.User),
		handleStringField(s, "password", &config.Password),
		handleStringField(s, "warehouse", &config.Warehouse),
		handleStringField(s, "role", &config.Role),
		handleBooleanStringAttribute(s, "validate_default_parameters", &config.ValidateDefaultParameters),
		// params are handled below
		// client ip
		func() error {
			if v, ok := s.GetOk("client_ip"); ok && v.(string) != "" {
				config.ClientIP = net.ParseIP(v.(string))
			}
			return nil
		}(),
		// protocol
		func() error {
			if v, ok := s.GetOk("protocol"); ok && v.(string) != "" {
				protocol, err := toProtocol(v.(string))
				if err != nil {
					return err
				}
				config.Protocol = string(protocol)
			}
			return nil
		}(),
		handleStringField(s, "host", &config.Host),
		handleIntAttribute(s, "port", &config.Port),
		// authenticator
		func() error {
			authType, err := sdk.ToExtendedAuthenticatorType(s.Get("authenticator").(string))
			if err != nil {
				return err
			}
			config.Authenticator = authType
			return nil
		}(),
		handleStringField(s, "passcode", &config.Passcode),
		handleBoolField(s, "passcode_in_password", &config.PasscodeInPassword),
		// okta url
		func() error {
			if v, ok := s.GetOk("okta_url"); ok && v.(string) != "" {
				oktaURL, err := url.Parse(v.(string))
				if err != nil {
					return fmt.Errorf("could not parse okta_url err = %w", err)
				}
				config.OktaURL = oktaURL
			}
			return nil
		}(),
		handleDurationInSecondsAttribute(s, "login_timeout", &config.LoginTimeout),
		handleDurationInSecondsAttribute(s, "request_timeout", &config.RequestTimeout),
		handleDurationInSecondsAttribute(s, "jwt_expire_timeout", &config.JWTExpireTimeout),
		handleDurationInSecondsAttribute(s, "client_timeout", &config.ClientTimeout),
		handleDurationInSecondsAttribute(s, "jwt_client_timeout", &config.JWTClientTimeout),
		handleDurationInSecondsAttribute(s, "external_browser_timeout", &config.ExternalBrowserTimeout),
		handleBoolField(s, "insecure_mode", &config.InsecureMode), //nolint:staticcheck
		// ocsp fail open
		func() error {
			if v := s.Get("ocsp_fail_open").(string); v != provider.BooleanDefault {
				parsed, err := provider.BooleanStringToBool(v)
				if err != nil {
					return err
				}
				if parsed {
					config.OCSPFailOpen = gosnowflake.OCSPFailOpenTrue
				} else {
					config.OCSPFailOpen = gosnowflake.OCSPFailOpenFalse
				}
			}
			return nil
		}(),
		// token
		func() error {
			if v, ok := s.GetOk("token"); ok && v.(string) != "" {
				config.Token = v.(string)
				config.Authenticator = gosnowflake.AuthTypeOAuth
			}
			return nil
		}(),
		// token accessor is handled below
		handleBoolField(s, "keep_session_alive", &config.KeepSessionAlive),
		// private key and private key passphrase are handled below
		handleBoolField(s, "disable_telemetry", &config.DisableTelemetry),
		handleBooleanStringAttribute(s, "client_request_mfa_token", &config.ClientRequestMfaToken),
		handleBooleanStringAttribute(s, "client_store_temporary_credential", &config.ClientStoreTemporaryCredential),
		handleBoolField(s, "disable_query_context_cache", &config.DisableQueryContextCache),
		handleBooleanStringAttribute(s, "include_retry_reason", &config.IncludeRetryReason),
		handleIntAttribute(s, "max_retry_count", &config.MaxRetryCount),
		// driver tracing
		func() error {
			if v, ok := s.GetOk("driver_tracing"); ok {
				driverLogLevel, err := sdk.ToDriverLogLevel(v.(string))
				if err != nil {
					return err
				}
				config.Tracing = string(driverLogLevel)
			}
			return nil
		}(),
		handleStringField(s, "tmp_directory_path", &config.TmpDirPath),
		handleBooleanStringAttribute(s, "disable_console_login", &config.DisableConsoleLogin),
		// profile is handled in the calling function
		// TODO(SNOW-1761318): handle DisableSamlURLCheck after upgrading the driver to at least 1.10.1
	)
	if err != nil {
		return nil, err
	}

	// account_name and organization_name override legacy account field
	accountName := s.Get("account_name").(string)
	organizationName := s.Get("organization_name").(string)
	if accountName != "" && organizationName != "" {
		config.Account = strings.Join([]string{organizationName, accountName}, "-")
	}

	m := make(map[string]interface{})
	if v, ok := s.GetOk("params"); ok {
		m = v.(map[string]interface{})
	}

	params := make(map[string]*string)
	for key, value := range m {
		strValue := value.(string)
		params[key] = &strValue
	}
	config.Params = params

	if v, ok := s.GetOk("token_accessor"); ok {
		if len(v.([]any)) > 0 {
			tokenAccessor := v.([]any)[0].(map[string]any)
			tokenEndpoint := tokenAccessor["token_endpoint"].(string)
			refreshToken := tokenAccessor["refresh_token"].(string)
			clientID := tokenAccessor["client_id"].(string)
			clientSecret := tokenAccessor["client_secret"].(string)
			redirectURI := tokenAccessor["redirect_uri"].(string)
			accessToken, err := GetAccessTokenWithRefreshToken(tokenEndpoint, clientID, clientSecret, refreshToken, redirectURI)
			if err != nil {
				return nil, fmt.Errorf("could not retrieve access token from refresh token, err = %w", err)
			}
			config.Token = accessToken
			config.Authenticator = gosnowflake.AuthTypeOAuth
		}
	}

	privateKey := s.Get("private_key").(string)
	privateKeyPassphrase := s.Get("private_key_passphrase").(string)
	v, err := getPrivateKey(privateKey, privateKeyPassphrase)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve private key: %w", err)
	}
	if v != nil {
		config.PrivateKey = v
	}

	return config, nil
}
