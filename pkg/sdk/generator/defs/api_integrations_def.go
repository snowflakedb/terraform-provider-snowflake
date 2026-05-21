package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var (
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
)

var ApiIntegrationEndpointPrefixDef = g.NewQueryStruct("ApiIntegrationEndpointPrefix").Text("Path", g.KeywordOptions().SingleQuotes().Required())

// ALLOWED_AUTHENTICATION_SECRETS = ALL | NONE | ( secrets... )
// Pattern follows sessionPolicySecondaryRoles in session_policies_def.go
var apiIntegrationAllowedAuthSecretsDef = g.NewQueryStruct("ApiIntegrationAllowedAuthenticationSecrets").
	OptionalSQLWithCustomFieldName("AllSecrets", "ALL").
	OptionalSQLWithCustomFieldName("NoSecrets", "NONE").
	List("AllowedList", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.ListOptions().Parentheses()).
	WithValidation(g.ExactlyOneValueSet, "AllSecrets", "NoSecrets", "AllowedList")

// API_USER_AUTHENTICATION = ( TYPE = OAUTH2 ... ) for git_https_api
// OAUTH_ALLOWED_SCOPES reuses ApiIntegrationScope from secrets_gen.go
var oauth2GitAuthDef = g.NewQueryStruct("OAuth2GitUserAuthentication").
	SQLWithCustomFieldName("authType", "TYPE = OAUTH2").
	TextAssignment("OAUTH_AUTHORIZATION_ENDPOINT", g.ParameterOptions().SingleQuotes().Required()).
	TextAssignment("OAUTH_TOKEN_ENDPOINT", g.ParameterOptions().SingleQuotes().Required()).
	TextAssignment("OAUTH_CLIENT_ID", g.ParameterOptions().SingleQuotes().Required()).
	TextAssignment("OAUTH_CLIENT_SECRET", g.ParameterOptions().SingleQuotes().Required()).
	OptionalNumberAssignment("OAUTH_ACCESS_TOKEN_VALIDITY", g.ParameterOptions().NoQuotes()).
	OptionalNumberAssignment("OAUTH_REFRESH_TOKEN_VALIDITY", g.ParameterOptions().NoQuotes()).
	ListAssignment("OAUTH_ALLOWED_SCOPES", "ApiIntegrationScope", g.ParameterOptions().Parentheses()).
	OptionalTextAssignment("OAUTH_USERNAME", g.ParameterOptions().SingleQuotes())

// API_USER_AUTHENTICATION = ( TYPE = OAUTH2 ... ) for external_mcp
var oauth2McpAuthDef = g.NewQueryStruct("OAuth2McpUserAuthentication").
	SQLWithCustomFieldName("authType", "TYPE = OAUTH2").
	TextAssignment("OAUTH_CLIENT_ID", g.ParameterOptions().SingleQuotes().Required()).
	TextAssignment("OAUTH_CLIENT_SECRET", g.ParameterOptions().SingleQuotes().Required()).
	TextAssignment("OAUTH_TOKEN_ENDPOINT", g.ParameterOptions().SingleQuotes().Required()).
	TextAssignment("OAUTH_AUTHORIZATION_ENDPOINT", g.ParameterOptions().SingleQuotes().Required()).
	OptionalEnumAssignment("OAUTH_CLIENT_AUTH_METHOD", ApiIntegrationOauthClientAuthMethodEnum, g.ParameterOptions().NoQuotes()).
	OptionalTextAssignment("OAUTH_DISCOVERY_URL", g.ParameterOptions().SingleQuotes()).
	OptionalNumberAssignment("OAUTH_REFRESH_TOKEN_VALIDITY", g.ParameterOptions().NoQuotes())

// API_USER_AUTHENTICATION = ( TYPE = OAUTH_DYNAMIC_CLIENT ... ) for external_mcp
var dynamicClientMcpAuthDef = g.NewQueryStruct("DynamicClientMcpUserAuthentication").
	SQLWithCustomFieldName("authType", "TYPE = OAUTH_DYNAMIC_CLIENT").
	TextAssignment("OAUTH_RESOURCE_URL", g.ParameterOptions().SingleQuotes().Required())

// TODO [SNOW-1016561]: all integrations reuse almost the same show, drop, and describe. For now we are copying it. Consider reusing in linked issue.
var apiIntegrationsDef = g.NewInterface(
	"ApiIntegrations",
	"ApiIntegration",
	g.KindOfT[sdkcommons.AccountObjectIdentifier](),
).
	WithEnums(ApiIntegrationAwsApiProviderTypeEnum, ApiIntegrationOauthClientAuthMethodEnum).
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
					SQLWithCustomFieldName("apiUserAuthentication", "API_USER_AUTHENTICATION = (TYPE = SNOWFLAKE_GITHUB_APP)"),
				g.KeywordOptions(),
			).
			OptionalQueryStructField(
				"GitHttpsApiOAuth2ProviderParams",
				g.NewQueryStruct("GitHttpsApiOAuth2Params").
					// TODO: Turn into generated enum (? and in other cases?)
					SQLWithCustomFieldName("apiProvider", "API_PROVIDER = git_https_api").
					QueryStructField(
						"ApiUserAuthentication",
						oauth2GitAuthDef,
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
						oauth2McpAuthDef,
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
						dynamicClientMcpAuthDef,
						g.ListOptions().SQL("API_USER_AUTHENTICATION =").Parentheses().NoComma(),
					),
				g.KeywordOptions(),
			).
			ListAssignment("API_ALLOWED_PREFIXES", "ApiIntegrationEndpointPrefix", g.ParameterOptions().Parentheses().Required()).
			// API_BLOCKED_PREFIXES is not supported for Amazon API Gateway; enforced in validation and github and maybe other types (maybe should be a part of particular type params)
			ListAssignment("API_BLOCKED_PREFIXES", "ApiIntegrationEndpointPrefix", g.ParameterOptions().Parentheses()).
			BooleanAssignment("ENABLED", g.ParameterOptions().Required()).
			OptionalComment().
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ConflictingFields, "IfNotExists", "OrReplace").
			WithValidation(g.ExactlyOneValueSet, "AwsApiProviderParams", "AzureApiProviderParams", "GoogleApiProviderParams", "GitHttpsApiTokenBasedProviderParams", "GitHttpsApiGithubAppProviderParams", "GitHttpsApiOAuth2ProviderParams", "GitHttpsApiPrivateLinkProviderParams", "ExternalMcpOAuth2ProviderParams", "ExternalMcpDynamicClientProviderParams"),
		ApiIntegrationEndpointPrefixDef,
		apiIntegrationAllowedAuthSecretsDef,
		oauth2GitAuthDef,
		oauth2McpAuthDef,
		dynamicClientMcpAuthDef,
	).
	AlterOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/alter-api-integration",
		g.NewQueryStruct("AlterApiIntegration").
			Alter().
			SQL("API INTEGRATION").
			IfExists().
			Name().
			// TODO: Validate with create options, the docs seem to be "simplistic"
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
						"GitHttpsApiGithubAppParams",
						g.NewQueryStruct("SetGitHttpsApiGithubAppParams").
							SQLWithCustomFieldName("ApiUserAuthentication", "API_USER_AUTHENTICATION = (TYPE = SNOWFLAKE_GITHUB_APP)"),
						g.KeywordOptions(),
					).
					OptionalQueryStructField(
						"GitHttpsApiOAuth2Params",
						g.NewQueryStruct("SetGitHttpsApiOAuth2Params").
							QueryStructField(
								"ApiUserAuthentication",
								oauth2GitAuthDef,
								g.ListOptions().SQL("API_USER_AUTHENTICATION =").Parentheses().NoComma(),
							),
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
								oauth2McpAuthDef,
								g.ListOptions().SQL("API_USER_AUTHENTICATION =").Parentheses().NoComma(),
							),
						g.KeywordOptions(),
					).
					OptionalQueryStructField(
						"ExternalMcpDynamicClientParams",
						g.NewQueryStruct("SetExternalMcpDynamicClientParams").
							QueryStructField(
								"ApiUserAuthentication",
								dynamicClientMcpAuthDef,
								g.ListOptions().SQL("API_USER_AUTHENTICATION =").Parentheses().NoComma(),
							),
						g.KeywordOptions(),
					).
					OptionalBooleanAssignment("ENABLED", g.ParameterOptions()).
					ListAssignment("API_ALLOWED_PREFIXES", "ApiIntegrationEndpointPrefix", g.ParameterOptions().Parentheses()).
					ListAssignment("API_BLOCKED_PREFIXES", "ApiIntegrationEndpointPrefix", g.ParameterOptions().Parentheses()).
					OptionalComment().
					// TODO [SNOW-2324252]: ConflictingFields generates everyValueSet; change to moreThanOneValueSet manually after regeneration
					WithValidation(g.ConflictingFields, "AwsParams", "AzureParams", "GoogleParams", "GitHttpsApiTokenBasedParams", "GitHttpsApiGithubAppParams", "GitHttpsApiOAuth2Params", "GitHttpsApiPrivateLinkParams", "ExternalMcpOAuth2Params", "ExternalMcpDynamicClientParams").
					WithValidation(g.AtLeastOneValueSet, "AwsParams", "AzureParams", "GoogleParams", "GitHttpsApiTokenBasedParams", "GitHttpsApiGithubAppParams", "GitHttpsApiOAuth2Params", "GitHttpsApiPrivateLinkParams", "ExternalMcpOAuth2Params", "ExternalMcpDynamicClientParams", "Enabled", "ApiAllowedPrefixes", "ApiBlockedPrefixes", "Comment"),
				g.KeywordOptions().SQL("SET"),
			).
			OptionalQueryStructField(
				"Unset",
				g.NewQueryStruct("ApiIntegrationUnset").
					OptionalSQL("API_KEY").
					OptionalSQL("ENABLED").
					OptionalSQL("API_BLOCKED_PREFIXES").
					OptionalSQL("ALLOWED_AUTHENTICATION_SECRETS").
					OptionalSQL("USE_PRIVATELINK_ENDPOINT").
					OptionalSQL("COMMENT").
					WithValidation(g.AtLeastOneValueSet, "ApiKey", "Enabled", "ApiBlockedPrefixes", "AllowedAuthenticationSecrets", "UsePrivatelinkEndpoint", "Comment"),
				g.ListOptions().NoParentheses().SQL("UNSET"),
			).
			OptionalSetTags().
			OptionalUnsetTags().
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ConflictingFields, "IfExists", "SetTags").
			WithValidation(g.ConflictingFields, "IfExists", "UnsetTags").
			WithValidation(g.ExactlyOneValueSet, "Set", "Unset", "SetTags", "UnsetTags"),
	).
	// TODO: Pull out common drop operation and reuse it
	DropOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/drop-integration",
		g.NewQueryStruct("DropApiIntegration").
			Drop().
			SQL("API INTEGRATION").
			IfExists().
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	).
	// TODO: Pull out common show operation and reuse it
	ShowOperationWithPairedStructs(
		"https://docs.snowflake.com/en/sql-reference/sql/show-integrations",
		// TODO: Generate convert function
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
	// TODO: Pull out common describe operation and reuse it
	DescribeOperationWithPairedStructs(
		g.DescriptionMappingKindSlice,
		"https://docs.snowflake.com/en/sql-reference/sql/desc-integration",
		// TODO: Generate convert function
		// TODO: Create common struct for properties as in other property-based describes (e.g. storage_integrations_ext.go)
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
	WithShowObjectType("Integration")
