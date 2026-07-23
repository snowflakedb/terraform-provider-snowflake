package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

func allowedApiAuthenticationIntegrations() *g.QueryStruct {
	return g.NewQueryStruct("ExternalAccessIntegrationAllowedApiAuthenticationIntegrations").
		OptionalSQLWithCustomFieldName("None", "none").
		List("Integrations", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.ListOptions().Parentheses()).
		WithValidation(g.ExactlyOneValueSet, "None", "Integrations")
}

func allowedAuthenticationSecrets() *g.QueryStruct {
	return g.NewQueryStruct("ExternalAccessIntegrationAllowedAuthenticationSecrets").
		OptionalSQLWithCustomFieldName("All", "all").
		OptionalSQLWithCustomFieldName("None", "none").
		List("Secrets", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.ListOptions().Parentheses()).
		WithValidation(g.ExactlyOneValueSet, "All", "None", "Secrets")
}

var externalAccessIntegrationsDef = g.NewInterface(
	"ExternalAccessIntegrations",
	"ExternalAccessIntegration",
	g.KindOfT[sdkcommons.AccountObjectIdentifier](),
).CreateOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/create-external-access-integration",
	g.NewQueryStruct("CreateExternalAccessIntegration").
		Create().
		OrReplace().
		SQL("EXTERNAL ACCESS INTEGRATION").
		Name().
		ListAssignment("ALLOWED_NETWORK_RULES", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.ParameterOptions().Parentheses().Required()).
		OptionalQueryStructField(
			"AllowedApiAuthenticationIntegrations",
			allowedApiAuthenticationIntegrations(),
			g.KeywordOptions().SQL("ALLOWED_API_AUTHENTICATION_INTEGRATIONS ="),
		).
		OptionalQueryStructField(
			"AllowedAuthenticationSecrets",
			allowedAuthenticationSecrets(),
			g.KeywordOptions().SQL("ALLOWED_AUTHENTICATION_SECRETS ="),
		).
		BooleanAssignment("ENABLED", g.ParameterOptions().Required()).
		OptionalComment().
		WithValidation(g.ValidIdentifier, "name"),
).AlterOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/alter-external-access-integration",
	g.NewQueryStruct("AlterExternalAccessIntegration").
		Alter().
		SQL("EXTERNAL ACCESS INTEGRATION").
		IfExists().
		Name().
		OptionalQueryStructField(
			"Set",
			g.NewQueryStruct("ExternalAccessIntegrationSet").
				ListAssignment("ALLOWED_NETWORK_RULES", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.ParameterOptions().Parentheses()).
				OptionalQueryStructField(
					"AllowedApiAuthenticationIntegrations",
					allowedApiAuthenticationIntegrations(),
					g.KeywordOptions().SQL("ALLOWED_API_AUTHENTICATION_INTEGRATIONS ="),
				).
				OptionalQueryStructField(
					"AllowedAuthenticationSecrets",
					allowedAuthenticationSecrets(),
					g.KeywordOptions().SQL("ALLOWED_AUTHENTICATION_SECRETS ="),
				).
				OptionalBooleanAssignment("ENABLED", g.ParameterOptions()).
				OptionalComment().
				PredefinedQueryStructField("SetTags", "[]TagAssociation", g.KeywordOptions().SQL("TAG")).
				WithValidation(g.AtLeastOneValueSet, "AllowedNetworkRules", "AllowedApiAuthenticationIntegrations", "AllowedAuthenticationSecrets", "Enabled", "Comment", "SetTags"),
			g.KeywordOptions().SQL("SET"),
		).
		OptionalQueryStructField(
			"Unset",
			g.NewQueryStruct("ExternalAccessIntegrationUnset").
				OptionalSQL("ALLOWED_NETWORK_RULES").
				OptionalSQL("ALLOWED_API_AUTHENTICATION_INTEGRATIONS").
				OptionalSQL("ALLOWED_AUTHENTICATION_SECRETS").
				OptionalSQL("COMMENT").
				PredefinedQueryStructField("UnsetTags", "[]ObjectIdentifier", g.KeywordOptions().SQL("TAG")).
				WithValidation(g.AtLeastOneValueSet, "AllowedNetworkRules", "AllowedApiAuthenticationIntegrations", "AllowedAuthenticationSecrets", "Comment", "UnsetTags"),
			g.ListOptions().NoParentheses().SQL("UNSET"),
		).
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ExactlyOneValueSet, "Set", "Unset"),
).DropOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/drop-integration",
	g.NewQueryStruct("DropExternalAccessIntegration").
		Drop().
		SQL("EXTERNAL ACCESS INTEGRATION").
		IfExists().
		Name().
		WithValidation(g.ValidIdentifier, "name"),
).ShowOperationWithPairedStructs(
	"https://docs.snowflake.com/en/sql-reference/sql/show-integrations",
	g.StructPair("showExternalAccessIntegrationsDbRow", "ExternalAccessIntegration").
		Time("created_on").
		Text("name").
		Text("type").
		Text("category").
		Bool("enabled").
		OptionalText("comment", g.WithRequiredInPlain()),
	g.NewQueryStruct("ShowExternalAccessIntegrations").
		Show().
		SQL("EXTERNAL ACCESS INTEGRATIONS").
		OptionalLike(),
).DescribeOperationWithPairedStructs(
	g.DescriptionMappingKindSlice,
	"https://docs.snowflake.com/en/sql-reference/sql/desc-integration",
	g.StructPair("descExternalAccessIntegrationsDbRow", "ExternalAccessIntegrationProperty").
		Text("property", g.WithPlainFieldName("Name")).
		Text("property_type", g.WithPlainFieldName("Type")).
		Text("property_value", g.WithPlainFieldName("Value")).
		Text("property_default", g.WithPlainFieldName("Default")),
	g.NewQueryStruct("DescribeExternalAccessIntegration").
		Describe().
		SQL("EXTERNAL ACCESS INTEGRATION").
		Name().
		WithValidation(g.ValidIdentifier, "name"),
).WithShowObjectType("Integration")
