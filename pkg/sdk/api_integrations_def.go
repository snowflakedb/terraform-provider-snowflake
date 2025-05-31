package sdk

import g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"

//go:generate go run ./poc/main.go

type ApiIntegrationAwsApiProviderType string

var (
	ApiIntegrationAwsApiGateway           ApiIntegrationAwsApiProviderType = "aws_api_gateway"
	ApiIntegrationAwsPrivateApiGateway    ApiIntegrationAwsApiProviderType = "aws_private_api_gateway"
	ApiIntegrationAwsGovApiGateway        ApiIntegrationAwsApiProviderType = "aws_gov_api_gateway"
	ApiIntegrationAwsGovPrivateApiGateway ApiIntegrationAwsApiProviderType = "aws_gov_private_api_gateway"
)

type AllowedAuthenticationSecretOption string

var (
	ApiIntegrationEndpointPrefixDef         = g.NewQueryStruct("ApiIntegrationEndpointPrefix").Text("Path", g.KeywordOptions().SingleQuotes().Required())
	AllowedAuthenticationSecretListItemsDef = g.NewQueryStruct("AllowedAuthenticationSecretListItems").Text("Secret", g.KeywordOptions().SingleQuotes().Required())
	AllowedAuthenticationSecretDef          = g.NewQueryStruct("AllowedAuthenticationSecret").
						OptionalListAssignment("AllowedAuthenticationSecretList", "AllowedAuthenticationSecretListItems", g.ParameterOptions().NoKey().NoQuotes()).
						OptionalTextAssignment("AllowedAuthenticationSecretOption", g.ParameterOptions().NoKey().NoQuotes()).
						WithValidation(g.ExactlyOneValueSet, "AllowedAuthenticationSecretList", "AllowedAuthenticationSecretOption")
)

// TODO [SNOW-1016561]: all integrations reuse almost the same show, drop, and describe. For now we are copying it. Consider reusing in linked issue.
var ApiIntegrationsDef = g.NewInterface(
	"ApiIntegrations",
	"ApiIntegration",
	g.KindOfT[AccountObjectIdentifier](),
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
					Assignment("API_PROVIDER", g.KindOfT[ApiIntegrationAwsApiProviderType](), g.ParameterOptions().NoQuotes().Required()).
					TextAssignment("API_AWS_ROLE_ARN", g.ParameterOptions().SingleQuotes().Required()).
					OptionalTextAssignment("API_KEY", g.ParameterOptions().SingleQuotes()),
				g.KeywordOptions(),
			).
			OptionalQueryStructField(
				"AzureApiProviderParams",
				g.NewQueryStruct("AzureApiParams").
					PredefinedQueryStructField("apiProvider", "string", g.StaticOptions().SQL("API_PROVIDER = azure_api_management")).
					TextAssignment("AZURE_TENANT_ID", g.ParameterOptions().SingleQuotes().Required()).
					TextAssignment("AZURE_AD_APPLICATION_ID", g.ParameterOptions().SingleQuotes().Required()).
					OptionalTextAssignment("API_KEY", g.ParameterOptions().SingleQuotes()),
				g.KeywordOptions(),
			).
			OptionalQueryStructField(
				"GoogleApiProviderParams",
				g.NewQueryStruct("GoogleApiParams").
					PredefinedQueryStructField("apiProvider", "string", g.StaticOptions().SQL("API_PROVIDER = google_api_gateway")).
					TextAssignment("GOOGLE_AUDIENCE", g.ParameterOptions().SingleQuotes().Required()),
				g.KeywordOptions(),
			).
			OptionalQueryStructField(
				"GitApiProviderParams",
				g.NewQueryStruct("GitApiParams").
					PredefinedQueryStructField("apiProvider", "string", g.StaticOptions().SQL("API_PROVIDER = git_https_api")).
					OptionalQueryStructField("AllowedAuthenticationSecret", AllowedAuthenticationSecretDef, g.KeywordOptions().SQL("ALLOWED_AUTHENTICATION_SECRETS = ")),
				g.KeywordOptions(),
			).
			ListAssignment("API_ALLOWED_PREFIXES", "ApiIntegrationEndpointPrefix", g.ParameterOptions().Parentheses().Required()).
			ListAssignment("API_BLOCKED_PREFIXES", "ApiIntegrationEndpointPrefix", g.ParameterOptions().Parentheses()).
			BooleanAssignment("ENABLED", g.ParameterOptions().Required()).
			OptionalComment().
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ConflictingFields, "IfNotExists", "OrReplace").
			WithValidation(g.ExactlyOneValueSet, "AwsApiProviderParams", "AzureApiProviderParams", "GoogleApiProviderParams", "GitApiProviderParams"),
		ApiIntegrationEndpointPrefixDef,
		AllowedAuthenticationSecretDef,
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
						"GitParams",
						g.NewQueryStruct("SetGitApiParams").
							QueryStructField("AllowedAuthenticationSecret", AllowedAuthenticationSecretDef, g.KeywordOptions().SQL("ALLOWED_AUTHENTICATION_SECRETS = ")),
						g.KeywordOptions(),
					).
					OptionalBooleanAssignment("ENABLED", g.ParameterOptions()).
					ListAssignment("API_ALLOWED_PREFIXES", "ApiIntegrationEndpointPrefix", g.ParameterOptions().Parentheses()).
					ListAssignment("API_BLOCKED_PREFIXES", "ApiIntegrationEndpointPrefix", g.ParameterOptions().Parentheses()).
					OptionalComment().
					// resulting validation changed to moreThanOneValueSet (not yet supported in the generator)
					WithValidation(g.ConflictingFields, "AwsParams", "AzureParams", "GoogleParams", "GitParams").
					WithValidation(g.AtLeastOneValueSet, "AwsParams", "AzureParams", "GoogleParams", "GitParams", "Enabled", "ApiAllowedPrefixes", "ApiBlockedPrefixes", "Comment"),
				g.KeywordOptions().SQL("SET"),
			).
			OptionalQueryStructField(
				"Unset",
				g.NewQueryStruct("ApiIntegrationUnset").
					OptionalSQL("API_KEY").
					OptionalSQL("ENABLED").
					OptionalSQL("API_BLOCKED_PREFIXES").
					OptionalSQL("COMMENT").
					WithValidation(g.AtLeastOneValueSet, "ApiKey", "Enabled", "ApiBlockedPrefixes", "Comment"),
				g.ListOptions().NoParentheses().SQL("UNSET"),
			).
			OptionalSetTags().
			OptionalUnsetTags().
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ConflictingFields, "IfExists", "SetTags").
			WithValidation(g.ConflictingFields, "IfExists", "UnsetTags").
			WithValidation(g.ExactlyOneValueSet, "Set", "Unset", "SetTags", "UnsetTags"),
	).
	DropOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/drop-integration",
		g.NewQueryStruct("DropApiIntegration").
			Drop().
			SQL("API INTEGRATION").
			IfExists().
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	).
	ShowOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/show-integrations",
		g.DbStruct("showApiIntegrationsDbRow").
			Text("name").
			Text("type").
			Text("category").
			Bool("enabled").
			OptionalText("comment").
			Time("created_on"),
		g.PlainStruct("ApiIntegration").
			Text("Name").
			Text("ApiType").
			Text("Category").
			Bool("Enabled").
			Text("Comment").
			Time("CreatedOn"),
		g.NewQueryStruct("ShowApiIntegrations").
			Show().
			SQL("API INTEGRATIONS").
			OptionalLike(),
	).
	ShowByIdOperationWithFiltering(
		g.ShowByIDLikeFiltering,
	).
	DescribeOperation(
		g.DescriptionMappingKindSlice,
		"https://docs.snowflake.com/en/sql-reference/sql/desc-integration",
		g.DbStruct("descApiIntegrationsDbRow").
			Text("property").
			Text("property_type").
			Text("property_value").
			Text("property_default"),
		g.PlainStruct("ApiIntegrationProperty").
			Text("Name").
			Text("Type").
			Text("Value").
			Text("Default"),
		g.NewQueryStruct("DescribeApiIntegration").
			Describe().
			SQL("API INTEGRATION").
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	)
