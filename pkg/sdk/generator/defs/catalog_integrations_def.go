package defs

import (
	"fmt"

	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var polarisRestConfigDef = g.NewQueryStruct("PolarisRestConfig").
	TextAssignment("CATALOG_URI", g.ParameterOptions().SingleQuotes().Required()).
	OptionalTextAssignment("CATALOG_API_TYPE", g.ParameterOptions().NoQuotes()).
	TextAssignment("CATALOG_NAME", g.ParameterOptions().SingleQuotes().Required()).
	OptionalTextAssignment("ACCESS_DELEGATION_MODE", g.ParameterOptions().NoQuotes())

var icebergRestRestConfigDef = g.NewQueryStruct("IcebergRestRestConfig").
	TextAssignment("CATALOG_URI", g.ParameterOptions().SingleQuotes().Required()).
	OptionalTextAssignment("PREFIX", g.ParameterOptions().SingleQuotes()).
	OptionalTextAssignment("CATALOG_NAME", g.ParameterOptions().SingleQuotes()).
	OptionalTextAssignment("CATALOG_API_TYPE", g.ParameterOptions().NoQuotes()).
	OptionalTextAssignment("ACCESS_DELEGATION_MODE", g.ParameterOptions().NoQuotes())

var sapBdcRestConfigDef = g.NewQueryStruct("SapBdcRestConfig").
	TextAssignment("SAP_BDC_INVITATION_LINK", g.ParameterOptions().SingleQuotes().Required()).
	OptionalTextAssignment("ACCESS_DELEGATION_MODE", g.ParameterOptions().NoQuotes())

var oAuthRestAuthenticationDef = g.NewQueryStruct("OAuthRestAuthentication").
	PredefinedQueryStructField("restAuthType", "string", g.StaticOptions().SQL(fmt.Sprintf("TYPE = %s", sdkcommons.RestAuthenticationTypeOAuth))).
	OptionalTextAssignment("OAUTH_TOKEN_URI", g.ParameterOptions().SingleQuotes()).
	TextAssignment("OAUTH_CLIENT_ID", g.ParameterOptions().SingleQuotes().Required()).
	TextAssignment("OAUTH_CLIENT_SECRET", g.ParameterOptions().SingleQuotes().Required()).
	ListAssignment("OAUTH_ALLOWED_SCOPES", "StringListItemWrapper", g.ParameterOptions().Parentheses().Required())

var bearerRestAuthenticationDef = g.NewQueryStruct("BearerRestAuthentication").
	PredefinedQueryStructField("restAuthType", "string", g.StaticOptions().SQL(fmt.Sprintf("TYPE = %s", sdkcommons.RestAuthenticationTypeBearer))).
	TextAssignment("BEARER_TOKEN", g.ParameterOptions().SingleQuotes().Required())

var sigV4RestAuthenticationDef = g.NewQueryStruct("SigV4RestAuthentication").
	PredefinedQueryStructField("restAuthType", "string", g.StaticOptions().SQL(fmt.Sprintf("TYPE = %s", sdkcommons.RestAuthenticationTypeSigV4))).
	TextAssignment("SIGV4_IAM_ROLE", g.ParameterOptions().SingleQuotes().Required()).
	OptionalTextAssignment("SIGV4_SIGNING_REGION", g.ParameterOptions().SingleQuotes()).
	OptionalTextAssignment("SIGV4_EXTERNAL_ID", g.ParameterOptions().SingleQuotes())

var catalogIntegrationsDef = g.NewInterface(
	"CatalogIntegrations",
	"CatalogIntegration",
	g.KindOfT[sdkcommons.AccountObjectIdentifier](),
).
	CreateOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/create-catalog-integration",
		g.NewQueryStruct("CreateCatalogIntegration").
			Create().
			OrReplace().
			SQL("CATALOG INTEGRATION").
			IfNotExists().
			Name().
			OptionalQueryStructField(
				"AwsGlueCatalogSourceParams",
				g.NewQueryStruct("AwsGlueParams").
					PredefinedQueryStructField("catalogSource", "string", g.StaticOptions().SQL(fmt.Sprintf("CATALOG_SOURCE = %s", sdkcommons.CatalogSourceTypeAWSGlue))).
					TextAssignment("GLUE_AWS_ROLE_ARN", g.ParameterOptions().SingleQuotes().Required()).
					TextAssignment("GLUE_CATALOG_ID", g.ParameterOptions().SingleQuotes().Required()).
					OptionalTextAssignment("GLUE_REGION", g.ParameterOptions().SingleQuotes()).
					OptionalTextAssignment("CATALOG_NAMESPACE", g.ParameterOptions().SingleQuotes()),
				g.KeywordOptions()).
			OptionalQueryStructField(
				"ObjectStorageCatalogSourceParams",
				g.NewQueryStruct("ObjectStorageParams").
					PredefinedQueryStructField("catalogSource", "string", g.StaticOptions().SQL(fmt.Sprintf("CATALOG_SOURCE = %s", sdkcommons.CatalogSourceTypeObjectStorage))),
				g.KeywordOptions()).
			OptionalQueryStructField(
				"PolarisCatalogSourceParams",
				g.NewQueryStruct("PolarisParams").
					PredefinedQueryStructField("catalogSource", "string", g.StaticOptions().SQL(fmt.Sprintf("CATALOG_SOURCE = %s", sdkcommons.CatalogSourceTypePolaris))).
					OptionalTextAssignment("CATALOG_NAMESPACE", g.ParameterOptions().SingleQuotes()).
					QueryStructField(
						"RestConfig",
						polarisRestConfigDef,
						g.ListOptions().SQL("REST_CONFIG =").Parentheses().NoComma()).
					QueryStructField(
						"RestAuthentication",
						oAuthRestAuthenticationDef,
						g.ListOptions().SQL("REST_AUTHENTICATION =").Parentheses().NoComma()),
				g.KeywordOptions()).
			OptionalQueryStructField(
				"IcebergRestCatalogSourceParams",
				g.NewQueryStruct("IcebergRestParams").
					PredefinedQueryStructField("catalogSource", "string", g.StaticOptions().SQL(fmt.Sprintf("CATALOG_SOURCE = %s", sdkcommons.CatalogSourceTypeIcebergREST))).
					OptionalTextAssignment("CATALOG_NAMESPACE", g.ParameterOptions().SingleQuotes()).
					QueryStructField(
						"RestConfig",
						icebergRestRestConfigDef,
						g.ListOptions().SQL("REST_CONFIG =").Parentheses().NoComma()).
					OptionalQueryStructField(
						"OAuthRestAuthentication",
						oAuthRestAuthenticationDef,
						g.ListOptions().SQL("REST_AUTHENTICATION =").Parentheses().NoComma()).
					OptionalQueryStructField(
						"BearerRestAuthentication",
						bearerRestAuthenticationDef,
						g.ListOptions().SQL("REST_AUTHENTICATION =").Parentheses().NoComma()).
					OptionalQueryStructField(
						"SigV4RestAuthentication",
						sigV4RestAuthenticationDef,
						g.ListOptions().SQL("REST_AUTHENTICATION =").Parentheses().NoComma()).
					WithValidation(g.ExactlyOneValueSet, "OAuthRestAuthentication", "BearerRestAuthentication", "SigV4RestAuthentication"),
				g.KeywordOptions()).
			OptionalQueryStructField(
				"SapBdcCatalogSourceParams",
				g.NewQueryStruct("SapBdcParams").
					PredefinedQueryStructField("catalogSource", "string", g.StaticOptions().SQL(fmt.Sprintf("CATALOG_SOURCE = %s", sdkcommons.CatalogSourceTypeSAPBusinessDataCloud))).
					QueryStructField(
						"RestConfig",
						sapBdcRestConfigDef,
						g.ListOptions().SQL("REST_CONFIG =").Parentheses().NoComma()),
				g.KeywordOptions()).
			TextAssignment("TABLE_FORMAT", g.ParameterOptions().NoQuotes().Required()).
			BooleanAssignment("ENABLED", g.ParameterOptions().Required()).
			OptionalNumberAssignment("REFRESH_INTERVAL_SECONDS", g.ParameterOptions().NoQuotes()).
			OptionalComment().
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ConflictingFields, "IfNotExists", "OrReplace").
			WithValidation(g.ExactlyOneValueSet, "AwsGlueCatalogSourceParams", "ObjectStorageCatalogSourceParams", "PolarisCatalogSourceParams", "IcebergRestCatalogSourceParams", "SapBdcCatalogSourceParams"),
	).
	AlterOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/alter-catalog-integration",
		g.NewQueryStruct("AlterCatalogIntegration").
			Alter().
			SQL("CATALOG INTEGRATION").
			IfExists().
			Name().
			OptionalQueryStructField(
				"Set",
				g.NewQueryStruct("CatalogIntegrationSet").
					OptionalQueryStructField(
						"SetOAuthRestAuthentication",
						g.NewQueryStruct("SetOAuthRestAuthentication").
							TextAssignment("OAUTH_CLIENT_SECRET", g.ParameterOptions().SingleQuotes().Required()),
						g.ListOptions().SQL("REST_AUTHENTICATION =").Parentheses().NoComma()).
					OptionalQueryStructField(
						"SetBearerRestAuthentication",
						g.NewQueryStruct("SetBearerRestAuthentication").
							TextAssignment("BEARER_TOKEN", g.ParameterOptions().SingleQuotes().Required()),
						g.ListOptions().SQL("REST_AUTHENTICATION =").Parentheses().NoComma()).
					OptionalBooleanAssignment("ENABLED", g.ParameterOptions()).
					OptionalNumberAssignment("REFRESH_INTERVAL_SECONDS", g.ParameterOptions().NoQuotes()).
					// TODO(SNOW-3121221): use COMMENT in unset and here use OptionalComment
					OptionalAssignment("COMMENT", "StringAllowEmpty", g.ParameterOptions()).
					WithValidation(g.ConflictingFields, "SetOAuthRestAuthentication", "SetBearerRestAuthentication").
					WithValidation(g.AtLeastOneValueSet, "SetOAuthRestAuthentication", "SetBearerRestAuthentication", "Enabled", "RefreshIntervalSeconds", "Comment"),
				g.KeywordOptions().SQL("SET"),
			).
			OptionalSetTags().
			OptionalUnsetTags().
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ExactlyOneValueSet, "Set", "SetTags", "UnsetTags"),
	).
	DropOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/drop-catalog-integration",
		g.NewQueryStruct("DropCatalogIntegration").
			Drop().
			SQL("CATALOG INTEGRATION").
			IfExists().
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	).
	ShowOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/show-catalog-integrations",
		g.DbStruct("showCatalogIntegrationsDbRow").
			Text("name").
			Bool("enabled").
			Text("type").
			Text("category").
			OptionalText("comment").
			Time("created_on"),
		g.PlainStruct("CatalogIntegration").
			Text("Name").
			Bool("Enabled").
			Text("Type").
			Text("Category").
			Text("Comment").
			Time("CreatedOn"),
		g.NewQueryStruct("ShowCatalogIntegration").
			Show().
			SQL("CATALOG INTEGRATIONS").
			OptionalLike(),
	).
	ShowByIdOperationWithFiltering(
		g.ShowByIDLikeFiltering,
	).
	DescribeOperation(
		g.DescriptionMappingKindSlice,
		"https://docs.snowflake.com/en/sql-reference/sql/desc-catalog-integration",
		g.DbStruct("descCatalogIntegrationsDbRow").
			Text("property").
			Text("property_type").
			Text("property_value").
			Text("property_default"),
		g.PlainStruct("CatalogIntegrationProperty").
			Text("Name").
			Text("Type").
			Text("Value").
			Text("Default"),
		g.NewQueryStruct("DescribeCatalogIntegration").
			Describe().
			SQL("CATALOG INTEGRATION").
			Name().
			WithValidation(g.ValidIdentifier, "name"))
