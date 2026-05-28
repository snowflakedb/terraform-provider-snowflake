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
					Assignment("API_PROVIDER", g.KindOfT[sdkcommons.ApiIntegrationAwsApiProviderType](), g.ParameterOptions().NoQuotes().Required()).
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
			ListAssignment("API_ALLOWED_PREFIXES", "ApiIntegrationEndpointPrefix", g.ParameterOptions().Parentheses().Required()).
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
					WithAdditionalValidations().
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
	DropOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/drop-integration",
		g.NewQueryStruct("DropApiIntegration").
			Drop().
			SQL("API INTEGRATION").
			IfExists().
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	).
	ShowOperationWithPairedStructs(
		"https://docs.snowflake.com/en/sql-reference/sql/show-integrations",
		g.StructPair("showApiIntegrationsDbRow", "ApiIntegration").
			Text("name").
			Text("type", g.WithPlainFieldName("ApiType")).
			Text("category").
			Bool("enabled").
			OptionalText("comment", g.WithRequiredInPlain()).
			Time("created_on").
			WithConvertGeneration(),
		g.NewQueryStruct("ShowApiIntegrations").
			Show().
			SQL("API INTEGRATIONS").
			OptionalLike(),
	).
	DescribeOperationWithPairedStructs(
		g.DescriptionMappingKindSlice,
		"https://docs.snowflake.com/en/sql-reference/sql/desc-integration",
		g.StructPair("descApiIntegrationsDbRow", "ApiIntegrationProperty").
			Text("property", g.WithPlainFieldName("Name")).
			Text("property_type", g.WithPlainFieldName("Type")).
			Text("property_value", g.WithPlainFieldName("Value")).
			Text("property_default", g.WithPlainFieldName("Default")).
			WithConvertGeneration(),
		g.NewQueryStruct("DescribeApiIntegration").
			Describe().
			SQL("API INTEGRATION").
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	).
	WithShowObjectType("Integration")
