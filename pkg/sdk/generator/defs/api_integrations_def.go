package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var (
	ApiIntegrationAllowedAuthenticationSecretsValueEnum = g.NewEnum(
		"ApiIntegrationAllowedAuthenticationSecretsValue", "ApiIntegrationAllowedAuthenticationSecretsValues",
		"ALL",
		"NONE",
	)
	ApiIntegrationAwsApiProviderTypeEnum = g.NewEnum(
		"ApiIntegrationAwsApiProviderType", "ApiIntegrationAwsApiProviderTypes",
		"aws_api_gateway",
		"aws_private_api_gateway",
		"aws_gov_api_gateway",
		"aws_gov_private_api_gateway",
	)
	ApiIntegrationOauthClientAuthMethodEnum = g.NewEnum(
		"ApiIntegrationOauthClientAuthMethod", "ApiIntegrationOauthClientAuthMethods",
		"CLIENT_SECRET_BASIC",
		"CLIENT_SECRET_POST",
	)
	ApiIntegrationOauthAllowedScopeEnum = g.NewEnum(
		"ApiIntegrationOauthAllowedScope", "ApiIntegrationOauthAllowedScopes",
		"read_api",
		"read_repository",
		"write_repository",
	)
	ApiIntegrationAzureApiProviderTypeEnum = g.NewEnum(
		"ApiIntegrationAzureApiProviderType", "ApiIntegrationAzureApiProviderTypes",
		"azure_api_management",
	)
	ApiIntegrationGoogleApiProviderTypeEnum = g.NewEnum(
		"ApiIntegrationGoogleApiProviderType", "ApiIntegrationGoogleApiProviderTypes",
		"google_api_gateway",
	)
	ApiIntegrationGitApiProviderTypeEnum = g.NewEnum(
		"ApiIntegrationGitApiProviderType", "ApiIntegrationGitApiProviderTypes",
		"git_https_api",
	)
	ApiIntegrationMcpApiProviderTypeEnum = g.NewEnum(
		"ApiIntegrationMcpApiProviderType", "ApiIntegrationMcpApiProviderTypes",
		"external_mcp",
	)
	ApiIntegrationUserAuthTypeEnum = g.NewEnum(
		"ApiIntegrationUserAuthType", "ApiIntegrationUserAuthTypes",
		"OAUTH2",
		"OAUTH_DYNAMIC_CLIENT",
		"SNOWFLAKE_GITHUB_APP",
	)
)

var apiIntegrationEndpointPrefixDef = g.NewQueryStruct("ApiIntegrationEndpointPrefix").Text("Path", g.KeywordOptions().SingleQuotes().Required())

var apiIntegrationAllowedAuthSecretsDef = g.NewQueryStruct("ApiIntegrationAllowedAuthenticationSecrets").
	OptionalSQLWithCustomFieldName("AllSecrets", "ALL").
	OptionalSQLWithCustomFieldName("NoSecrets", "NONE").
	List("AllowedList", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.ListOptions().Parentheses()).
	WithValidation(g.ExactlyOneValueSet, "AllSecrets", "NoSecrets", "AllowedList")

var apiIntegrationOauthAllowedScopeItemDef = g.NewQueryStruct("ApiIntegrationOauthAllowedScopeItem").
	Enum("Scope", ApiIntegrationOauthAllowedScopeEnum, g.KeywordOptions().SingleQuotes().Required())

var apiIntegrationOauth2GitAuthDef = g.NewQueryStruct("OAuth2GitUserAuthentication").
	SQLWithCustomFieldName("authType", "TYPE = OAUTH2").
	TextAssignment("OAUTH_AUTHORIZATION_ENDPOINT", g.ParameterOptions().SingleQuotes().Required()).
	TextAssignment("OAUTH_TOKEN_ENDPOINT", g.ParameterOptions().SingleQuotes().Required()).
	TextAssignment("OAUTH_CLIENT_ID", g.ParameterOptions().SingleQuotes().Required()).
	TextAssignment("OAUTH_CLIENT_SECRET", g.ParameterOptions().SingleQuotes().Required()).
	OptionalNumberAssignment("OAUTH_ACCESS_TOKEN_VALIDITY", g.ParameterOptions().NoQuotes()).
	OptionalNumberAssignment("OAUTH_REFRESH_TOKEN_VALIDITY", g.ParameterOptions().NoQuotes()).
	ListAssignment("OAUTH_ALLOWED_SCOPES", "ApiIntegrationOauthAllowedScopeItem", g.ParameterOptions().Parentheses()).
	OptionalTextAssignment("OAUTH_USERNAME", g.ParameterOptions().SingleQuotes())

var apiIntegrationOauth2McpAuthDef = g.NewQueryStruct("OAuth2McpUserAuthentication").
	SQLWithCustomFieldName("authType", "TYPE = OAUTH2").
	TextAssignment("OAUTH_CLIENT_ID", g.ParameterOptions().SingleQuotes().Required()).
	TextAssignment("OAUTH_CLIENT_SECRET", g.ParameterOptions().SingleQuotes().Required()).
	TextAssignment("OAUTH_TOKEN_ENDPOINT", g.ParameterOptions().SingleQuotes().Required()).
	TextAssignment("OAUTH_AUTHORIZATION_ENDPOINT", g.ParameterOptions().SingleQuotes().Required()).
	OptionalEnumAssignment("OAUTH_CLIENT_AUTH_METHOD", ApiIntegrationOauthClientAuthMethodEnum, g.ParameterOptions().NoQuotes()).
	OptionalTextAssignment("OAUTH_DISCOVERY_URL", g.ParameterOptions().SingleQuotes()).
	OptionalNumberAssignment("OAUTH_REFRESH_TOKEN_VALIDITY", g.ParameterOptions().NoQuotes())

var apiIntegrationDynamicClientMcpAuthDef = g.NewQueryStruct("DynamicClientMcpUserAuthentication").
	SQLWithCustomFieldName("authType", "TYPE = OAUTH_DYNAMIC_CLIENT").
	TextAssignment("OAUTH_RESOURCE_URL", g.ParameterOptions().SingleQuotes().Required())

var apiIntegrationGithubAppAuthDef = g.NewQueryStruct("GithubAppUserAuthentication").
	SQLWithCustomFieldName("authType", "TYPE = SNOWFLAKE_GITHUB_APP")

// TODO(SNOW-1016561): all integrations reuse almost the same show, drop, and describe. For now we are copying it. Consider reusing in linked issue.
var apiIntegrationsDef = g.NewInterface(
	"ApiIntegrations",
	"ApiIntegration",
	g.KindOfT[sdkcommons.AccountObjectIdentifier](),
).
	WithEnums(
		ApiIntegrationAllowedAuthenticationSecretsValueEnum,
		ApiIntegrationAwsApiProviderTypeEnum,
		ApiIntegrationOauthClientAuthMethodEnum,
		ApiIntegrationOauthAllowedScopeEnum,
		ApiIntegrationAzureApiProviderTypeEnum,
		ApiIntegrationGoogleApiProviderTypeEnum,
		ApiIntegrationGitApiProviderTypeEnum,
		ApiIntegrationMcpApiProviderTypeEnum,
		ApiIntegrationUserAuthTypeEnum,
	).
	CreateOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/create-api-integration",
		g.NewQueryStruct("CreateApiIntegration").
			Create().
			OrReplace().
			SQL("API INTEGRATION").
			IfNotExists().
			Name().
			OptionalQueryStructField(
				"AwsApiProviderParams",
				g.NewQueryStruct("AwsApiParams").
					EnumAssignment("API_PROVIDER", ApiIntegrationAwsApiProviderTypeEnum, g.ParameterOptions().NoQuotes().Required()).
					TextAssignment("API_AWS_ROLE_ARN", g.ParameterOptions().SingleQuotes().Required()).
					OptionalTextAssignment("API_KEY", g.ParameterOptions().SingleQuotes()),
				g.KeywordOptions(),
			).
			OptionalQueryStructField(
				"AzureApiProviderParams",
				g.NewQueryStruct("AzureApiParams").
					SQLWithCustomFieldName("apiProvider", "API_PROVIDER = azure_api_management").
					TextAssignment("AZURE_TENANT_ID", g.ParameterOptions().SingleQuotes().Required()).
					TextAssignment("AZURE_AD_APPLICATION_ID", g.ParameterOptions().SingleQuotes().Required()).
					OptionalTextAssignment("API_KEY", g.ParameterOptions().SingleQuotes()),
				g.KeywordOptions(),
			).
			OptionalQueryStructField(
				"GoogleApiProviderParams",
				g.NewQueryStruct("GoogleApiParams").
					SQLWithCustomFieldName("apiProvider", "API_PROVIDER = google_api_gateway").
					TextAssignment("GOOGLE_AUDIENCE", g.ParameterOptions().SingleQuotes().Required()),
				g.KeywordOptions(),
			).
			OptionalQueryStructField(
				"GitHttpsApiTokenBasedProviderParams",
				g.NewQueryStruct("GitHttpsApiTokenBasedParams").
					SQLWithCustomFieldName("apiProvider", "API_PROVIDER = git_https_api").
					OptionalQueryStructField(
						"AllowedAuthenticationSecrets",
						apiIntegrationAllowedAuthSecretsDef,
						g.KeywordOptions().SQL("ALLOWED_AUTHENTICATION_SECRETS ="),
					),
				g.KeywordOptions(),
			).
			OptionalQueryStructField(
				"GitHttpsApiGithubAppProviderParams",
				g.NewQueryStruct("GitHttpsApiGithubAppParams").
					SQLWithCustomFieldName("apiProvider", "API_PROVIDER = git_https_api").
					QueryStructField(
						"ApiUserAuthentication",
						apiIntegrationGithubAppAuthDef,
						g.ListOptions().SQL("API_USER_AUTHENTICATION =").Parentheses().NoComma(),
					),
				g.KeywordOptions(),
			).
			OptionalQueryStructField(
				"GitHttpsApiOAuth2ProviderParams",
				g.NewQueryStruct("GitHttpsApiOAuth2Params").
					SQLWithCustomFieldName("apiProvider", "API_PROVIDER = git_https_api").
					QueryStructField(
						"ApiUserAuthentication",
						apiIntegrationOauth2GitAuthDef,
						g.ListOptions().SQL("API_USER_AUTHENTICATION =").Parentheses().NoComma(),
					),
				g.KeywordOptions(),
			).
			OptionalQueryStructField(
				"GitHttpsApiPrivateLinkProviderParams",
				g.NewQueryStruct("GitHttpsApiPrivateLinkParams").
					SQLWithCustomFieldName("apiProvider", "API_PROVIDER = git_https_api").
					OptionalQueryStructField(
						"AllowedAuthenticationSecrets",
						apiIntegrationAllowedAuthSecretsDef,
						g.KeywordOptions().SQL("ALLOWED_AUTHENTICATION_SECRETS ="),
					).
					BooleanAssignment("USE_PRIVATELINK_ENDPOINT", g.ParameterOptions().Required()).
					ListAssignment("TLS_TRUSTED_CERTIFICATES", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.ParameterOptions().Parentheses()),
				g.KeywordOptions(),
			).
			OptionalQueryStructField(
				"ExternalMcpOAuth2ProviderParams",
				g.NewQueryStruct("ExternalMcpOAuth2Params").
					SQLWithCustomFieldName("apiProvider", "API_PROVIDER = external_mcp").
					QueryStructField(
						"ApiUserAuthentication",
						apiIntegrationOauth2McpAuthDef,
						g.ListOptions().SQL("API_USER_AUTHENTICATION =").Parentheses().NoComma(),
					),
				g.KeywordOptions(),
			).
			OptionalQueryStructField(
				"ExternalMcpDynamicClientProviderParams",
				g.NewQueryStruct("ExternalMcpDynamicClientParams").
					SQLWithCustomFieldName("apiProvider", "API_PROVIDER = external_mcp").
					QueryStructField(
						"ApiUserAuthentication",
						apiIntegrationDynamicClientMcpAuthDef,
						g.ListOptions().SQL("API_USER_AUTHENTICATION =").Parentheses().NoComma(),
					),
				g.KeywordOptions(),
			).
			ListAssignment("API_ALLOWED_PREFIXES", "ApiIntegrationEndpointPrefix", g.ParameterOptions().Parentheses().Required()).
			// TODO(next pr): API_BLOCKED_PREFIXES is not supported for Amazon API Gateway, github and maybe other types; in next PR test the compatibility of common params for each api_provider and sub-type
			ListAssignment("API_BLOCKED_PREFIXES", "ApiIntegrationEndpointPrefix", g.ParameterOptions().Parentheses()).
			BooleanAssignment("ENABLED", g.ParameterOptions().Required()).
			OptionalComment().
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ConflictingFields, "IfNotExists", "OrReplace").
			WithValidation(g.ConflictingFields, "OrReplace", "ExternalMcpOAuth2ProviderParams").
			WithValidation(g.ConflictingFields, "OrReplace", "ExternalMcpDynamicClientProviderParams").
			WithValidation(
				g.ExactlyOneValueSet,
				"AwsApiProviderParams",
				"AzureApiProviderParams",
				"GoogleApiProviderParams",
				"GitHttpsApiTokenBasedProviderParams",
				"GitHttpsApiGithubAppProviderParams",
				"GitHttpsApiOAuth2ProviderParams",
				"GitHttpsApiPrivateLinkProviderParams",
				"ExternalMcpOAuth2ProviderParams",
				"ExternalMcpDynamicClientProviderParams",
			),
		apiIntegrationEndpointPrefixDef,
		apiIntegrationAllowedAuthSecretsDef,
		apiIntegrationOauthAllowedScopeItemDef,
		apiIntegrationOauth2GitAuthDef,
		apiIntegrationOauth2McpAuthDef,
		apiIntegrationDynamicClientMcpAuthDef,
		apiIntegrationGithubAppAuthDef,
	).
	AlterOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/alter-api-integration",
		g.NewQueryStruct("AlterApiIntegration").
			Alter().
			SQL("API INTEGRATION").
			IfExists().
			Name().
			OptionalQueryStructField(
				"Set",
				g.NewQueryStruct("ApiIntegrationSet").
					OptionalQueryStructField(
						"AwsParams",
						g.NewQueryStruct("SetAwsApiParams").
							OptionalTextAssignment("API_AWS_ROLE_ARN", g.ParameterOptions().SingleQuotes()).
							OptionalTextAssignment("API_KEY", g.ParameterOptions().SingleQuotes()).
							WithValidation(g.AtLeastOneValueSet, "ApiAwsRoleArn", "ApiKey"),
						g.KeywordOptions(),
					).
					OptionalQueryStructField(
						"AzureParams",
						g.NewQueryStruct("SetAzureApiParams").
							OptionalTextAssignment("AZURE_TENANT_ID", g.ParameterOptions().SingleQuotes()).
							OptionalTextAssignment("AZURE_AD_APPLICATION_ID", g.ParameterOptions().SingleQuotes()).
							OptionalTextAssignment("API_KEY", g.ParameterOptions().SingleQuotes()).
							WithValidation(g.AtLeastOneValueSet, "AzureTenantId", "AzureAdApplicationId", "ApiKey"),
						g.KeywordOptions(),
					).
					OptionalQueryStructField(
						"GoogleParams",
						g.NewQueryStruct("SetGoogleApiParams").
							TextAssignment("GOOGLE_AUDIENCE", g.ParameterOptions().SingleQuotes().Required()),
						g.KeywordOptions(),
					).
					OptionalQueryStructField(
						"GitHttpsApiTokenBasedParams",
						g.NewQueryStruct("SetGitHttpsApiTokenBasedParams").
							OptionalQueryStructField(
								"AllowedAuthenticationSecrets",
								apiIntegrationAllowedAuthSecretsDef,
								g.KeywordOptions().SQL("ALLOWED_AUTHENTICATION_SECRETS ="),
							).
							WithValidation(g.AtLeastOneValueSet, "AllowedAuthenticationSecrets"),
						g.KeywordOptions(),
					).
					OptionalQueryStructField(
						"GitHttpsApiPrivateLinkParams",
						g.NewQueryStruct("SetGitHttpsApiPrivateLinkParams").
							OptionalQueryStructField(
								"AllowedAuthenticationSecrets",
								apiIntegrationAllowedAuthSecretsDef,
								g.KeywordOptions().SQL("ALLOWED_AUTHENTICATION_SECRETS ="),
							).
							OptionalBooleanAssignment("USE_PRIVATELINK_ENDPOINT", g.ParameterOptions()).
							ListAssignment("TLS_TRUSTED_CERTIFICATES", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.ParameterOptions().Parentheses()).
							WithValidation(g.AtLeastOneValueSet, "AllowedAuthenticationSecrets", "UsePrivatelinkEndpoint", "TlsTrustedCertificates"),
						g.KeywordOptions(),
					).
					OptionalQueryStructField(
						"ExternalMcpOAuth2Params",
						g.NewQueryStruct("SetExternalMcpOAuth2Params").
							QueryStructField(
								"ApiUserAuthentication",
								apiIntegrationOauth2McpAuthDef,
								g.ListOptions().SQL("API_USER_AUTHENTICATION =").Parentheses().NoComma(),
							),
						g.KeywordOptions(),
					).
					OptionalBooleanAssignment("ENABLED", g.ParameterOptions()).
					ListAssignment("API_ALLOWED_PREFIXES", "ApiIntegrationEndpointPrefix", g.ParameterOptions().Parentheses()).
					ListAssignment("API_BLOCKED_PREFIXES", "ApiIntegrationEndpointPrefix", g.ParameterOptions().Parentheses()).
					OptionalComment().
					WithValidation(
						g.MoreThanOneValueSet,
						"AwsParams",
						"AzureParams",
						"GoogleParams",
						"GitHttpsApiTokenBasedParams",
						"GitHttpsApiPrivateLinkParams",
						"ExternalMcpOAuth2Params",
					).
					WithValidation(
						g.AtLeastOneValueSet,
						"AwsParams",
						"AzureParams",
						"GoogleParams",
						"GitHttpsApiTokenBasedParams",
						"GitHttpsApiPrivateLinkParams",
						"ExternalMcpOAuth2Params",
						"Enabled",
						"ApiAllowedPrefixes",
						"ApiBlockedPrefixes",
						"Comment",
					),
				g.KeywordOptions().SQL("SET"),
			).
			OptionalQueryStructField(
				"Unset",
				g.NewQueryStruct("ApiIntegrationUnset").
					OptionalQueryStructField(
						"AwsParams",
						g.NewQueryStruct("UnsetAwsApiParams").
							OptionalSQL("API_KEY").
							WithValidation(g.AtLeastOneValueSet, "ApiKey"),
						g.ListOptions().NoParentheses(),
					).
					OptionalQueryStructField(
						"AzureParams",
						g.NewQueryStruct("UnsetAzureApiParams").
							OptionalSQL("API_KEY").
							WithValidation(g.AtLeastOneValueSet, "ApiKey"),
						g.ListOptions().NoParentheses(),
					).
					OptionalQueryStructField(
						"GitHttpsApiTokenBasedParams",
						g.NewQueryStruct("UnsetGitHttpsApiTokenBasedParams").
							OptionalSQL("ALLOWED_AUTHENTICATION_SECRETS").
							WithValidation(g.AtLeastOneValueSet, "AllowedAuthenticationSecrets"),
						g.ListOptions().NoParentheses(),
					).
					OptionalQueryStructField(
						"GitHttpsApiPrivateLinkParams",
						g.NewQueryStruct("UnsetGitHttpsApiPrivateLinkParams").
							OptionalSQL("ALLOWED_AUTHENTICATION_SECRETS").
							OptionalSQL("TLS_TRUSTED_CERTIFICATES").
							OptionalSQL("USE_PRIVATELINK_ENDPOINT").
							WithValidation(g.AtLeastOneValueSet, "AllowedAuthenticationSecrets", "TlsTrustedCertificates", "UsePrivatelinkEndpoint"),
						g.ListOptions().NoParentheses(),
					).
					OptionalSQL("ENABLED").
					OptionalSQL("API_BLOCKED_PREFIXES").
					OptionalSQL("COMMENT").
					WithValidation(g.MoreThanOneValueSet, "AwsParams", "AzureParams", "GitHttpsApiTokenBasedParams", "GitHttpsApiPrivateLinkParams").
					WithValidation(g.AtLeastOneValueSet, "AwsParams", "AzureParams", "GitHttpsApiTokenBasedParams", "GitHttpsApiPrivateLinkParams", "Enabled", "ApiBlockedPrefixes", "Comment"),
				g.ListOptions().NoParentheses().SQL("UNSET"),
			).
			OptionalSetTags().
			OptionalUnsetTags().
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ConflictingFields, "IfExists", "SetTags").
			WithValidation(g.ConflictingFields, "IfExists", "UnsetTags").
			WithValidation(g.ExactlyOneValueSet, "Set", "Unset", "SetTags", "UnsetTags"),
	).
	// TODO(SNOW-1016561): Pull out common drop operation and reuse it
	DropOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/drop-integration",
		g.NewQueryStruct("DropApiIntegration").
			Drop().
			SQL("API INTEGRATION").
			IfExists().
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	).
	// TODO(SNOW-1016561): Pull out common show operation and reuse it
	ShowOperationWithPairedStructs(
		"https://docs.snowflake.com/en/sql-reference/sql/show-integrations",
		g.StructPair("showApiIntegrationsDbRow", "ApiIntegration").
			Text("name").
			Text("type", g.WithPlainFieldName("ApiType")).
			Text("category").
			Bool("enabled").
			OptionalText("comment", g.WithRequiredInPlain()).
			Time("created_on"),
		g.NewQueryStruct("ShowApiIntegrations").
			Show().
			SQL("API INTEGRATIONS").
			OptionalLike(),
	).
	// TODO(SNOW-1016561): Pull out common describe operation and reuse it
	DescribeOperationWithPairedStructs(
		g.DescriptionMappingKindSlice,
		"https://docs.snowflake.com/en/sql-reference/sql/desc-integration",
		g.StructPair("descApiIntegrationsDbRow", "ApiIntegrationProperty").
			Text("property", g.WithPlainFieldName("Name")).
			Text("property_type", g.WithPlainFieldName("Type")).
			Text("property_value", g.WithPlainFieldName("Value")).
			Text("property_default", g.WithPlainFieldName("Default")),
		g.NewQueryStruct("DescribeApiIntegration").
			Describe().
			SQL("API INTEGRATION").
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	).
	WithShowObjectType("Integration").
	WithCustomInterfaceMethod(
		"DescribeAwsDetails",
		"DescribeAwsDetails returns converted describe output for AWS API integrations.",
		[]*g.MethodParameter{g.NewMethodParameter("id", g.KindOfT[sdkcommons.AccountObjectIdentifier]())},
		"*ApiIntegrationAwsDetails", "error",
	).
	WithCustomInterfaceMethod(
		"DescribeAzureDetails",
		"DescribeAzureDetails returns converted describe output for Azure API integrations.",
		[]*g.MethodParameter{g.NewMethodParameter("id", g.KindOfT[sdkcommons.AccountObjectIdentifier]())},
		"*ApiIntegrationAzureDetails", "error",
	).
	WithCustomInterfaceMethod(
		"DescribeGoogleDetails",
		"DescribeGoogleDetails returns converted describe output for Google API integrations.",
		[]*g.MethodParameter{g.NewMethodParameter("id", g.KindOfT[sdkcommons.AccountObjectIdentifier]())},
		"*ApiIntegrationGoogleDetails", "error",
	).
	WithCustomInterfaceMethod(
		"DescribeGitHttpsApiDetails",
		"DescribeGitHttpsApiDetails returns converted describe output for git HTTPS API integrations.",
		[]*g.MethodParameter{g.NewMethodParameter("id", g.KindOfT[sdkcommons.AccountObjectIdentifier]())},
		"*ApiIntegrationGitHttpsApiDetails", "error",
	).
	WithCustomInterfaceMethod(
		"DescribeExternalMcpDetails",
		"DescribeExternalMcpDetails returns converted describe output for external MCP API integrations.",
		[]*g.MethodParameter{g.NewMethodParameter("id", g.KindOfT[sdkcommons.AccountObjectIdentifier]())},
		"*ApiIntegrationExternalMcpDetails", "error",
	).
	WithCustomInterfaceMethod(
		"DescribeAllDetails",
		"DescribeAllDetails returns parsed describe output for any API integration type.",
		[]*g.MethodParameter{g.NewMethodParameter("id", g.KindOfT[sdkcommons.AccountObjectIdentifier]())},
		"*ApiIntegrationAllDetails", "error",
	)
