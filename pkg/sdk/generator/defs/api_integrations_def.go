package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var ApiIntegrationEndpointPrefixDef = g.NewQueryStruct("ApiIntegrationEndpointPrefix").Text("Path", g.KeywordOptions().SingleQuotes().Required())

// TODO [SNOW-1016561]: all integrations reuse almost the same show, drop, and describe. For now we are copying it. Consider reusing in linked issue.
var apiIntegrationsDef = g.NewInterface(
	"ApiIntegrations",
	"ApiIntegration",
	g.KindOfT[sdkcommons.AccountObjectIdentifier](),
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
					// TODO: Turn into generated enum
					Assignment("API_PROVIDER", g.KindOfT[sdkcommons.ApiIntegrationAwsApiProviderType](), g.ParameterOptions().NoQuotes().Required()).
					TextAssignment("API_AWS_ROLE_ARN", g.ParameterOptions().SingleQuotes().Required()).
					OptionalTextAssignment("API_KEY", g.ParameterOptions().SingleQuotes()),
				g.KeywordOptions(),
			).
			OptionalQueryStructField(
				"AzureApiProviderParams",
				g.NewQueryStruct("AzureApiParams").
					// TODO: Turn into generated enum (?)
					SQLWithCustomFieldName("apiProvider", "API_PROVIDER = azure_api_management").
					TextAssignment("AZURE_TENANT_ID", g.ParameterOptions().SingleQuotes().Required()).
					TextAssignment("AZURE_AD_APPLICATION_ID", g.ParameterOptions().SingleQuotes().Required()).
					OptionalTextAssignment("API_KEY", g.ParameterOptions().SingleQuotes()),
				g.KeywordOptions(),
			).
			OptionalQueryStructField(
				"GoogleApiProviderParams",
				g.NewQueryStruct("GoogleApiParams").
					// TODO: Turn into generated enum (?)
					SQLWithCustomFieldName("apiProvider", "API_PROVIDER = google_api_gateway").
					TextAssignment("GOOGLE_AUDIENCE", g.ParameterOptions().SingleQuotes().Required()),
				g.KeywordOptions(),
			).
			// TODO: Add git repository variant (with internal 4 variants: Token, GitHub OAuth, OAuth2 parameters, Private Link)
			// TODO: Add external MCP server variant (with internal 2 variants: OAuth2 and Dynamic Client Registration)
			ListAssignment("API_ALLOWED_PREFIXES", "ApiIntegrationEndpointPrefix", g.ParameterOptions().Parentheses().Required()).
			// TODO: API_BLOCKED_PREFIXES for amazon api gateway is not supported (?)
			ListAssignment("API_BLOCKED_PREFIXES", "ApiIntegrationEndpointPrefix", g.ParameterOptions().Parentheses()).
			BooleanAssignment("ENABLED", g.ParameterOptions().Required()).
			OptionalComment().
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ConflictingFields, "IfNotExists", "OrReplace").
			WithValidation(g.ExactlyOneValueSet, "AwsApiProviderParams", "AzureApiProviderParams", "GoogleApiProviderParams"),
		ApiIntegrationEndpointPrefixDef,
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
					OptionalBooleanAssignment("ENABLED", g.ParameterOptions()).
					ListAssignment("API_ALLOWED_PREFIXES", "ApiIntegrationEndpointPrefix", g.ParameterOptions().Parentheses()).
					ListAssignment("API_BLOCKED_PREFIXES", "ApiIntegrationEndpointPrefix", g.ParameterOptions().Parentheses()).
					OptionalComment().
					// resulting validation changed to moreThanOneValueSet (not yet supported in the generator)
					WithValidation(g.ConflictingFields, "AwsParams", "AzureParams", "GoogleParams").
					WithValidation(g.AtLeastOneValueSet, "AwsParams", "AzureParams", "GoogleParams", "Enabled", "ApiAllowedPrefixes", "ApiBlockedPrefixes", "Comment"),
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
