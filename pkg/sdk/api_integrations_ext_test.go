package sdk

import (
	"strings"
)

func init() {
	id := apiIntegrationsTestIdAccountObjectIdentifier

	const (
		awsAllowedPrefix    = "https://123456.execute-api.us-west-2.amazonaws.com/prod/"
		azureAllowedPrefix  = "https://apim-hello-world.azure-api.net/"
		googleAllowedPrefix = "https://gateway-id-123456.uc.gateway.dev/"
		gitAllowedPrefix    = "https://github.com/my-org/"
		mcpAllowedPrefix    = "https://mcp.example.com/"

		apiAwsRoleArn              = "arn:aws:iam::000000000001:/role/test"
		azureTenantId              = "00000000-0000-0000-0000-000000000000"
		azureAdApplicationId       = "11111111-1111-1111-1111-111111111111"
		googleAudience             = "api-gateway-id-123456.apigateway.gcp-project.cloud.goog"
		oauthAuthorizationEndpoint = "https://auth.example.com/authorize"
		oauthTokenEndpoint         = "https://auth.example.com/token" //nolint:gosec // test credentials
		oauthClientId              = "oauth-client-id-123"
		oauthClientSecret          = "oauth-client-secret-456" //nolint:gosec // test credentials
		oauthResourceUrl           = "https://resource.example.com"
	)

	allowedSecret := NewSchemaObjectIdentifier("db", "schema", "secret")
	cert := NewSchemaObjectIdentifier("db", "schema", "cert")

	apiIntegrationsTests.Create.
		withModifyAndExpectedSqlf(
			case_ApiIntegrations_sql_Create_basic,
			func(opts *CreateApiIntegrationOptions) {
				opts.AwsApiProviderParams = &AwsApiParams{
					ApiProvider:   ApiIntegrationAwsApiProviderTypeAwsApiGateway,
					ApiAwsRoleArn: apiAwsRoleArn,
				}
				opts.ApiAllowedPrefixes = []ApiIntegrationEndpointPrefix{{Path: awsAllowedPrefix}}
				opts.Enabled = true
			},
			strings.Join([]string{
				"CREATE API INTEGRATION %s",
				"API_PROVIDER = aws_api_gateway",
				"API_AWS_ROLE_ARN = '%s'",
				"API_ALLOWED_PREFIXES = ('%s')",
				"ENABLED = true",
			}, " "), id.FullyQualifiedName(), apiAwsRoleArn, awsAllowedPrefix,
		).
		withModifyAndExpectedSqlf(
			case_ApiIntegrations_sql_Create_all,
			func(opts *CreateApiIntegrationOptions) {
				opts.IfNotExists = new(true)
				opts.AwsApiProviderParams = &AwsApiParams{
					ApiProvider:   ApiIntegrationAwsApiProviderTypeAwsPrivateApiGateway,
					ApiAwsRoleArn: apiAwsRoleArn,
					ApiKey:        new("key"),
				}
				opts.ApiAllowedPrefixes = []ApiIntegrationEndpointPrefix{{Path: awsAllowedPrefix}}
				opts.ApiBlockedPrefixes = []ApiIntegrationEndpointPrefix{{Path: googleAllowedPrefix}, {Path: azureAllowedPrefix}}
				opts.Enabled = false
				opts.Comment = new("some comment")
			},
			strings.Join([]string{
				"CREATE API INTEGRATION IF NOT EXISTS %s",
				"API_PROVIDER = aws_private_api_gateway",
				"API_AWS_ROLE_ARN = '%s'",
				"API_KEY = 'key'",
				"API_ALLOWED_PREFIXES = ('%s')",
				"API_BLOCKED_PREFIXES = ('%s', '%s')",
				"ENABLED = false",
				"COMMENT = 'some comment'",
			}, " "), id.FullyQualifiedName(), apiAwsRoleArn, awsAllowedPrefix, googleAllowedPrefix, azureAllowedPrefix,
		)

	apiIntegrationsTests.Create.
		withAdditionalSqlCasef(
			"sql_Create_all_azure",
			func(opts *CreateApiIntegrationOptions) {
				opts.IfNotExists = new(true)
				opts.AzureApiProviderParams = &AzureApiParams{
					AzureTenantId:        azureTenantId,
					AzureAdApplicationId: azureAdApplicationId,
					ApiKey:               new("key"),
				}
				opts.ApiAllowedPrefixes = []ApiIntegrationEndpointPrefix{{Path: azureAllowedPrefix}}
				opts.ApiBlockedPrefixes = []ApiIntegrationEndpointPrefix{{Path: awsAllowedPrefix}, {Path: googleAllowedPrefix}}
				opts.Enabled = false
				opts.Comment = new("some comment")
			},
			strings.Join([]string{
				"CREATE API INTEGRATION IF NOT EXISTS %s",
				"API_PROVIDER = azure_api_management",
				"AZURE_TENANT_ID = '%s'",
				"AZURE_AD_APPLICATION_ID = '%s'",
				"API_KEY = 'key'",
				"API_ALLOWED_PREFIXES = ('%s')",
				"API_BLOCKED_PREFIXES = ('%s', '%s')",
				"ENABLED = false",
				"COMMENT = 'some comment'",
			}, " "), id.FullyQualifiedName(), azureTenantId, azureAdApplicationId, azureAllowedPrefix, awsAllowedPrefix, googleAllowedPrefix,
		).
		withAdditionalSqlCasef(
			"sql_Create_all_google",
			func(opts *CreateApiIntegrationOptions) {
				opts.IfNotExists = new(true)
				opts.GoogleApiProviderParams = &GoogleApiParams{
					GoogleAudience: googleAudience,
				}
				opts.ApiAllowedPrefixes = []ApiIntegrationEndpointPrefix{{Path: googleAllowedPrefix}}
				opts.ApiBlockedPrefixes = []ApiIntegrationEndpointPrefix{{Path: awsAllowedPrefix}, {Path: azureAllowedPrefix}}
				opts.Enabled = false
				opts.Comment = new("some comment")
			},
			strings.Join([]string{
				"CREATE API INTEGRATION IF NOT EXISTS %s",
				"API_PROVIDER = google_api_gateway",
				"GOOGLE_AUDIENCE = '%s'",
				"API_ALLOWED_PREFIXES = ('%s')",
				"API_BLOCKED_PREFIXES = ('%s', '%s')",
				"ENABLED = false",
				"COMMENT = 'some comment'",
			}, " "), id.FullyQualifiedName(), googleAudience, googleAllowedPrefix, awsAllowedPrefix, azureAllowedPrefix,
		).
		withAdditionalSqlCasef(
			"sql_Create_git_token",
			func(opts *CreateApiIntegrationOptions) {
				opts.IfNotExists = new(true)
				opts.GitHttpsApiTokenBasedProviderParams = &GitHttpsApiTokenBasedParams{
					AllowedAuthenticationSecrets: &ApiIntegrationAllowedAuthenticationSecrets{
						AllSecrets: new(true),
					},
				}
				opts.ApiAllowedPrefixes = []ApiIntegrationEndpointPrefix{{Path: gitAllowedPrefix}}
				opts.ApiBlockedPrefixes = []ApiIntegrationEndpointPrefix{{Path: awsAllowedPrefix}}
				opts.Enabled = false
				opts.Comment = new("some comment")
			},
			strings.Join([]string{
				"CREATE API INTEGRATION IF NOT EXISTS %s",
				"API_PROVIDER = git_https_api",
				"ALLOWED_AUTHENTICATION_SECRETS = ALL",
				"API_ALLOWED_PREFIXES = ('%s')",
				"API_BLOCKED_PREFIXES = ('%s')",
				"ENABLED = false",
				"COMMENT = 'some comment'",
			}, " "), id.FullyQualifiedName(), gitAllowedPrefix, awsAllowedPrefix,
		).
		withAdditionalSqlCasef(
			"sql_Create_git_github_app",
			func(opts *CreateApiIntegrationOptions) {
				opts.IfNotExists = new(true)
				opts.GitHttpsApiGithubAppProviderParams = &GitHttpsApiGithubAppParams{}
				opts.ApiAllowedPrefixes = []ApiIntegrationEndpointPrefix{{Path: gitAllowedPrefix}}
				opts.ApiBlockedPrefixes = []ApiIntegrationEndpointPrefix{{Path: awsAllowedPrefix}}
				opts.Enabled = false
				opts.Comment = new("some comment")
			},
			strings.Join([]string{
				"CREATE API INTEGRATION IF NOT EXISTS %s",
				"API_PROVIDER = git_https_api",
				"API_USER_AUTHENTICATION = (TYPE = SNOWFLAKE_GITHUB_APP)",
				"API_ALLOWED_PREFIXES = ('%s')",
				"API_BLOCKED_PREFIXES = ('%s')",
				"ENABLED = false",
				"COMMENT = 'some comment'",
			}, " "), id.FullyQualifiedName(), gitAllowedPrefix, awsAllowedPrefix,
		).
		withAdditionalSqlCasef(
			"sql_Create_git_oauth2",
			func(opts *CreateApiIntegrationOptions) {
				opts.IfNotExists = new(true)
				opts.GitHttpsApiOAuth2ProviderParams = &GitHttpsApiOAuth2Params{
					ApiUserAuthentication: OAuth2GitUserAuthentication{
						OauthAuthorizationEndpoint: oauthAuthorizationEndpoint,
						OauthTokenEndpoint:         oauthTokenEndpoint,
						OauthClientId:              oauthClientId,
						OauthClientSecret:          oauthClientSecret,
						OauthAccessTokenValidity:   new(3600),
						OauthRefreshTokenValidity:  new(86400),
						OauthAllowedScopes: []ApiIntegrationOauthAllowedScopeItem{
							{Scope: ApiIntegrationOauthAllowedScopeReadApi},
							{Scope: ApiIntegrationOauthAllowedScopeReadRepository},
						},
						OauthUsername: new("user@example.com"),
					},
				}
				opts.ApiAllowedPrefixes = []ApiIntegrationEndpointPrefix{{Path: gitAllowedPrefix}}
				opts.ApiBlockedPrefixes = []ApiIntegrationEndpointPrefix{{Path: awsAllowedPrefix}}
				opts.Enabled = false
				opts.Comment = new("some comment")
			},
			strings.Join([]string{
				"CREATE API INTEGRATION IF NOT EXISTS %s",
				"API_PROVIDER = git_https_api",
				"API_USER_AUTHENTICATION = (" + strings.Join([]string{
					"TYPE = OAUTH2",
					"OAUTH_AUTHORIZATION_ENDPOINT = '%s'",
					"OAUTH_TOKEN_ENDPOINT = '%s'",
					"OAUTH_CLIENT_ID = '%s'",
					"OAUTH_CLIENT_SECRET = '%s'",
					"OAUTH_ACCESS_TOKEN_VALIDITY = 3600",
					"OAUTH_REFRESH_TOKEN_VALIDITY = 86400",
					"OAUTH_ALLOWED_SCOPES = ('read_api', 'read_repository')",
					"OAUTH_USERNAME = 'user@example.com'",
				}, " ") + ")",
				"API_ALLOWED_PREFIXES = ('%s')",
				"API_BLOCKED_PREFIXES = ('%s')",
				"ENABLED = false",
				"COMMENT = 'some comment'",
			}, " "), id.FullyQualifiedName(), oauthAuthorizationEndpoint, oauthTokenEndpoint, oauthClientId, oauthClientSecret, gitAllowedPrefix, awsAllowedPrefix,
		).
		withAdditionalSqlCasef(
			"sql_Create_git_private_link",
			func(opts *CreateApiIntegrationOptions) {
				opts.IfNotExists = new(true)
				opts.GitHttpsApiPrivateLinkProviderParams = &GitHttpsApiPrivateLinkParams{
					AllowedAuthenticationSecrets: &ApiIntegrationAllowedAuthenticationSecrets{
						AllowedList: []SchemaObjectIdentifier{allowedSecret},
					},
					UsePrivatelinkEndpoint: true,
					TlsTrustedCertificates: []SchemaObjectIdentifier{cert},
				}
				opts.ApiAllowedPrefixes = []ApiIntegrationEndpointPrefix{{Path: gitAllowedPrefix}}
				opts.ApiBlockedPrefixes = []ApiIntegrationEndpointPrefix{{Path: awsAllowedPrefix}}
				opts.Enabled = false
				opts.Comment = new("some comment")
			},
			strings.Join([]string{
				"CREATE API INTEGRATION IF NOT EXISTS %s",
				"API_PROVIDER = git_https_api",
				"ALLOWED_AUTHENTICATION_SECRETS = (%s)",
				"USE_PRIVATELINK_ENDPOINT = true",
				"TLS_TRUSTED_CERTIFICATES = (%s)",
				"API_ALLOWED_PREFIXES = ('%s')",
				"API_BLOCKED_PREFIXES = ('%s')",
				"ENABLED = false",
				"COMMENT = 'some comment'",
			}, " "), id.FullyQualifiedName(), allowedSecret.FullyQualifiedName(), cert.FullyQualifiedName(), gitAllowedPrefix, awsAllowedPrefix,
		).
		withAdditionalSqlCasef(
			"sql_Create_mcp_oauth2",
			func(opts *CreateApiIntegrationOptions) {
				opts.IfNotExists = new(true)
				opts.ExternalMcpOAuth2ProviderParams = &ExternalMcpOAuth2Params{
					ApiUserAuthentication: OAuth2McpUserAuthentication{
						OauthClientId:              oauthClientId,
						OauthClientSecret:          oauthClientSecret,
						OauthTokenEndpoint:         oauthTokenEndpoint,
						OauthAuthorizationEndpoint: oauthAuthorizationEndpoint,
						OauthClientAuthMethod:      new(ApiIntegrationOauthClientAuthMethodClientSecretPost),
						OauthDiscoveryUrl:          new("https://auth.example.com/.well-known/openid-configuration"),
						OauthRefreshTokenValidity:  new(86400),
					},
				}
				opts.ApiAllowedPrefixes = []ApiIntegrationEndpointPrefix{{Path: mcpAllowedPrefix}}
				opts.ApiBlockedPrefixes = []ApiIntegrationEndpointPrefix{{Path: awsAllowedPrefix}}
				opts.Enabled = false
				opts.Comment = new("some comment")
			},
			strings.Join([]string{
				"CREATE API INTEGRATION IF NOT EXISTS %s",
				"API_PROVIDER = external_mcp",
				"API_USER_AUTHENTICATION = (" + strings.Join([]string{
					"TYPE = OAUTH2",
					"OAUTH_CLIENT_ID = '%s'",
					"OAUTH_CLIENT_SECRET = '%s'",
					"OAUTH_TOKEN_ENDPOINT = '%s'",
					"OAUTH_AUTHORIZATION_ENDPOINT = '%s'",
					"OAUTH_CLIENT_AUTH_METHOD = CLIENT_SECRET_POST",
					"OAUTH_DISCOVERY_URL = 'https://auth.example.com/.well-known/openid-configuration'",
					"OAUTH_REFRESH_TOKEN_VALIDITY = 86400",
				}, " ") + ")",
				"API_ALLOWED_PREFIXES = ('%s')",
				"API_BLOCKED_PREFIXES = ('%s')",
				"ENABLED = false",
				"COMMENT = 'some comment'",
			}, " "), id.FullyQualifiedName(), oauthClientId, oauthClientSecret, oauthTokenEndpoint, oauthAuthorizationEndpoint, mcpAllowedPrefix, awsAllowedPrefix,
		).
		withAdditionalSqlCasef(
			"sql_Create_mcp_dynamic_client",
			func(opts *CreateApiIntegrationOptions) {
				opts.IfNotExists = new(true)
				opts.ExternalMcpDynamicClientProviderParams = &ExternalMcpDynamicClientParams{
					ApiUserAuthentication: DynamicClientMcpUserAuthentication{
						OauthResourceUrl: oauthResourceUrl,
					},
				}
				opts.ApiAllowedPrefixes = []ApiIntegrationEndpointPrefix{{Path: mcpAllowedPrefix}}
				opts.ApiBlockedPrefixes = []ApiIntegrationEndpointPrefix{{Path: awsAllowedPrefix}}
				opts.Enabled = false
				opts.Comment = new("some comment")
			},
			strings.Join([]string{
				"CREATE API INTEGRATION IF NOT EXISTS %s",
				"API_PROVIDER = external_mcp",
				"API_USER_AUTHENTICATION = (TYPE = OAUTH_DYNAMIC_CLIENT OAUTH_RESOURCE_URL = '%s')",
				"API_ALLOWED_PREFIXES = ('%s')",
				"API_BLOCKED_PREFIXES = ('%s')",
				"ENABLED = false",
				"COMMENT = 'some comment'",
			}, " "), id.FullyQualifiedName(), oauthResourceUrl, mcpAllowedPrefix, awsAllowedPrefix,
		).
		withAdditionalSqlCasef(
			"sql_Create_orReplace",
			func(opts *CreateApiIntegrationOptions) {
				opts.OrReplace = new(true)
				opts.AwsApiProviderParams = &AwsApiParams{
					ApiProvider:   ApiIntegrationAwsApiProviderTypeAwsApiGateway,
					ApiAwsRoleArn: apiAwsRoleArn,
				}
				opts.ApiAllowedPrefixes = []ApiIntegrationEndpointPrefix{{Path: awsAllowedPrefix}}
				opts.Enabled = true
			},
			strings.Join([]string{
				"CREATE OR REPLACE API INTEGRATION %s",
				"API_PROVIDER = aws_api_gateway",
				"API_AWS_ROLE_ARN = '%s'",
				"API_ALLOWED_PREFIXES = ('%s')",
				"ENABLED = true",
			}, " "), id.FullyQualifiedName(), apiAwsRoleArn, awsAllowedPrefix,
		)

	apiIntegrationsTests.Alter.
		withModify(
			case_ApiIntegrations_validation_Alter_opts_ConflictingFields_IfExists_SetTags,
			func(opts *AlterApiIntegrationOptions) {
				opts.IfExists = new(true)
				opts.SetTags = []TagAssociation{{Name: NewAccountObjectIdentifier("name"), Value: "value"}}
			},
		).
		withModify(
			case_ApiIntegrations_validation_Alter_opts_ConflictingFields_IfExists_UnsetTags,
			func(opts *AlterApiIntegrationOptions) {
				opts.IfExists = new(true)
				opts.UnsetTags = []ObjectIdentifier{NewAccountObjectIdentifier("one")}
			},
		).
		withModifyAndExpectedSqlf(
			case_ApiIntegrations_sql_Alter_Set,
			func(opts *AlterApiIntegrationOptions) {
				opts.Set = &ApiIntegrationSet{
					AwsParams: &SetAwsApiParams{
						ApiAwsRoleArn: new("new-aws-role-arn"),
						ApiKey:        new("key"),
					},
					Enabled:            new(true),
					ApiAllowedPrefixes: []ApiIntegrationEndpointPrefix{{Path: awsAllowedPrefix}},
					ApiBlockedPrefixes: []ApiIntegrationEndpointPrefix{{Path: azureAllowedPrefix}, {Path: googleAllowedPrefix}},
					Comment:            new("comment"),
				}
			},
			strings.Join([]string{
				"ALTER API INTEGRATION %s SET",
				"API_AWS_ROLE_ARN = 'new-aws-role-arn'",
				"API_KEY = 'key'",
				"ENABLED = true",
				"API_ALLOWED_PREFIXES = ('%s')",
				"API_BLOCKED_PREFIXES = ('%s', '%s')",
				"COMMENT = 'comment'",
			}, " "), id.FullyQualifiedName(), awsAllowedPrefix, azureAllowedPrefix, googleAllowedPrefix,
		).
		withModifyAndExpectedSqlf(
			case_ApiIntegrations_sql_Alter_Unset,
			func(opts *AlterApiIntegrationOptions) {
				opts.Unset = &ApiIntegrationUnset{
					Enabled:            new(true),
					ApiBlockedPrefixes: new(true),
					Comment:            new(true),
				}
			},
			"ALTER API INTEGRATION %s UNSET ENABLED, API_BLOCKED_PREFIXES, COMMENT", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_ApiIntegrations_sql_Alter_SetTags,
			func(opts *AlterApiIntegrationOptions) {
				opts.SetTags = []TagAssociation{
					{Name: NewAccountObjectIdentifier("name"), Value: "value"},
					{Name: NewAccountObjectIdentifier("second-name"), Value: "second-value"},
				}
			},
			`ALTER API INTEGRATION %s SET TAG "name" = 'value', "second-name" = 'second-value'`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_ApiIntegrations_sql_Alter_UnsetTags,
			func(opts *AlterApiIntegrationOptions) {
				opts.UnsetTags = []ObjectIdentifier{
					NewAccountObjectIdentifier("name"),
					NewAccountObjectIdentifier("second-name"),
				}
			},
			`ALTER API INTEGRATION %s UNSET TAG "name", "second-name"`, id.FullyQualifiedName(),
		)

	apiIntegrationsTests.Alter.
		withAdditionalSqlCasef(
			"sql_Alter_set_azure",
			func(opts *AlterApiIntegrationOptions) {
				opts.Set = &ApiIntegrationSet{
					AzureParams: &SetAzureApiParams{
						AzureAdApplicationId: new("new-azure-ad-application-id"),
						ApiKey:               new("key"),
					},
					Enabled:            new(true),
					ApiAllowedPrefixes: []ApiIntegrationEndpointPrefix{{Path: azureAllowedPrefix}},
					ApiBlockedPrefixes: []ApiIntegrationEndpointPrefix{{Path: awsAllowedPrefix}, {Path: googleAllowedPrefix}},
					Comment:            new("comment"),
				}
			},
			strings.Join([]string{
				"ALTER API INTEGRATION %s SET",
				"AZURE_AD_APPLICATION_ID = 'new-azure-ad-application-id'",
				"API_KEY = 'key'",
				"ENABLED = true",
				"API_ALLOWED_PREFIXES = ('%s')",
				"API_BLOCKED_PREFIXES = ('%s', '%s')",
				"COMMENT = 'comment'",
			}, " "), id.FullyQualifiedName(), azureAllowedPrefix, awsAllowedPrefix, googleAllowedPrefix,
		).
		withAdditionalSqlCasef(
			"sql_Alter_set_google",
			func(opts *AlterApiIntegrationOptions) {
				opts.Set = &ApiIntegrationSet{
					Enabled:            new(true),
					ApiAllowedPrefixes: []ApiIntegrationEndpointPrefix{{Path: googleAllowedPrefix}},
					ApiBlockedPrefixes: []ApiIntegrationEndpointPrefix{{Path: awsAllowedPrefix}, {Path: azureAllowedPrefix}},
					Comment:            new("comment"),
				}
			},
			strings.Join([]string{
				"ALTER API INTEGRATION %s SET",
				"ENABLED = true",
				"API_ALLOWED_PREFIXES = ('%s')",
				"API_BLOCKED_PREFIXES = ('%s', '%s')",
				"COMMENT = 'comment'",
			}, " "), id.FullyQualifiedName(), googleAllowedPrefix, awsAllowedPrefix, azureAllowedPrefix,
		).
		withAdditionalSqlCasef(
			"sql_Alter_set_git_token",
			func(opts *AlterApiIntegrationOptions) {
				opts.Set = &ApiIntegrationSet{
					GitHttpsApiTokenBasedParams: &SetGitHttpsApiTokenBasedParams{
						AllowedAuthenticationSecrets: &ApiIntegrationAllowedAuthenticationSecrets{
							NoSecrets: new(true),
						},
					},
				}
			},
			"ALTER API INTEGRATION %s SET ALLOWED_AUTHENTICATION_SECRETS = NONE", id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_set_git_private_link",
			func(opts *AlterApiIntegrationOptions) {
				opts.Set = &ApiIntegrationSet{
					GitHttpsApiPrivateLinkParams: &SetGitHttpsApiPrivateLinkParams{
						AllowedAuthenticationSecrets: &ApiIntegrationAllowedAuthenticationSecrets{
							AllSecrets: new(true),
						},
						UsePrivatelinkEndpoint: new(true),
						TlsTrustedCertificates: []SchemaObjectIdentifier{cert},
					},
				}
			},
			strings.Join([]string{
				"ALTER API INTEGRATION %s SET",
				"ALLOWED_AUTHENTICATION_SECRETS = ALL",
				"USE_PRIVATELINK_ENDPOINT = true",
				"TLS_TRUSTED_CERTIFICATES = (%s)",
			}, " "), id.FullyQualifiedName(), cert.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_set_mcp_oauth2",
			func(opts *AlterApiIntegrationOptions) {
				opts.Set = &ApiIntegrationSet{
					ExternalMcpOAuth2Params: &SetExternalMcpOAuth2Params{
						ApiUserAuthentication: OAuth2McpUserAuthentication{
							OauthClientId:              oauthClientId,
							OauthClientSecret:          oauthClientSecret,
							OauthTokenEndpoint:         oauthTokenEndpoint,
							OauthAuthorizationEndpoint: oauthAuthorizationEndpoint,
						},
					},
				}
			},
			strings.Join([]string{
				"ALTER API INTEGRATION %s SET",
				"API_USER_AUTHENTICATION = (" + strings.Join([]string{
					"TYPE = OAUTH2",
					"OAUTH_CLIENT_ID = '%s'",
					"OAUTH_CLIENT_SECRET = '%s'",
					"OAUTH_TOKEN_ENDPOINT = '%s'",
					"OAUTH_AUTHORIZATION_ENDPOINT = '%s'",
				}, " ") + ")",
			}, " "), id.FullyQualifiedName(), oauthClientId, oauthClientSecret, oauthTokenEndpoint, oauthAuthorizationEndpoint,
		).
		withAdditionalSqlCasef(
			"sql_Alter_unset_single_common",
			func(opts *AlterApiIntegrationOptions) {
				opts.Unset = &ApiIntegrationUnset{Comment: new(true)}
			},
			"ALTER API INTEGRATION %s UNSET COMMENT", id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_unset_aws",
			func(opts *AlterApiIntegrationOptions) {
				opts.Unset = &ApiIntegrationUnset{AwsParams: &UnsetAwsApiParams{ApiKey: new(true)}}
			},
			"ALTER API INTEGRATION %s UNSET API_KEY", id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_unset_azure",
			func(opts *AlterApiIntegrationOptions) {
				opts.Unset = &ApiIntegrationUnset{AzureParams: &UnsetAzureApiParams{ApiKey: new(true)}}
			},
			"ALTER API INTEGRATION %s UNSET API_KEY", id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_unset_git_token",
			func(opts *AlterApiIntegrationOptions) {
				opts.Unset = &ApiIntegrationUnset{GitHttpsApiTokenBasedParams: &UnsetGitHttpsApiTokenBasedParams{AllowedAuthenticationSecrets: new(true)}}
			},
			"ALTER API INTEGRATION %s UNSET ALLOWED_AUTHENTICATION_SECRETS", id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_unset_git_private_link",
			func(opts *AlterApiIntegrationOptions) {
				opts.Unset = &ApiIntegrationUnset{
					GitHttpsApiPrivateLinkParams: &UnsetGitHttpsApiPrivateLinkParams{
						AllowedAuthenticationSecrets: new(true),
						TlsTrustedCertificates:       new(true),
						UsePrivatelinkEndpoint:       new(true),
					},
				}
			},
			"ALTER API INTEGRATION %s UNSET ALLOWED_AUTHENTICATION_SECRETS, TLS_TRUSTED_CERTIFICATES, USE_PRIVATELINK_ENDPOINT", id.FullyQualifiedName(),
		)

	apiIntegrationsTests.Drop.
		withExpectedSqlf(case_ApiIntegrations_sql_Drop_basic,
			"DROP API INTEGRATION %s", id.FullyQualifiedName()).
		withModifyAndExpectedSqlf(
			case_ApiIntegrations_sql_Drop_all,
			func(opts *DropApiIntegrationOptions) { opts.IfExists = new(true) },
			"DROP API INTEGRATION IF EXISTS %s", id.FullyQualifiedName(),
		)

	apiIntegrationsTests.Show.
		withExpectedSql(case_ApiIntegrations_sql_Show_basic, "SHOW API INTEGRATIONS").
		withModifyAndExpectedSqlf(
			case_ApiIntegrations_sql_Show_all,
			func(opts *ShowApiIntegrationOptions) { opts.Like = &Like{Pattern: new(id.Name())} },
			"SHOW API INTEGRATIONS LIKE '%s'", id.Name(),
		).
		withModifyAndExpectedSqlf(
			case_ApiIntegrations_sql_Show_Like,
			func(opts *ShowApiIntegrationOptions) { opts.Like = &Like{Pattern: new(id.Name())} },
			"SHOW API INTEGRATIONS LIKE '%s'", id.Name(),
		)

	apiIntegrationsTests.Describe.
		withExpectedSqlf(case_ApiIntegrations_sql_Describe_basic,
			"DESCRIBE API INTEGRATION %s", id.FullyQualifiedName())
}
