package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

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
			BooleanAssignment("ENABLED", g.ParameterOptions().Required()).
			OptionalQueryStructField(
				"ObjectStoreParams",
				g.NewQueryStruct("ObjectStoreParams").
					PredefinedQueryStructField("catalogSource", "string", g.StaticOptions().SQL("CATALOG_SOURCE = OBJECT_STORE")).
					TextAssignment("TABLE_FORMAT", g.ParameterOptions().NoQuotes().Required()),
				g.KeywordOptions(),
			).
			OptionalQueryStructField(
				"GlueParams",
				g.NewQueryStruct("GlueParams").
					PredefinedQueryStructField("catalogSource", "string", g.StaticOptions().SQL("CATALOG_SOURCE = GLUE")).
					TextAssignment("TABLE_FORMAT", g.ParameterOptions().NoQuotes().Required()).
					TextAssignment("GLUE_AWS_ROLE_ARN", g.ParameterOptions().SingleQuotes().Required()).
					TextAssignment("GLUE_CATALOG_ID", g.ParameterOptions().SingleQuotes().Required()).
					OptionalTextAssignment("GLUE_REGION", g.ParameterOptions().SingleQuotes()),
				g.KeywordOptions(),
			).
			OptionalQueryStructField(
				"IcebergRestParams",
				g.NewQueryStruct("IcebergRestParams").
					PredefinedQueryStructField("catalogSource", "string", g.StaticOptions().SQL("CATALOG_SOURCE = ICEBERG_REST")).
					TextAssignment("TABLE_FORMAT", g.ParameterOptions().NoQuotes().Required()),
				g.KeywordOptions(),
			).
			OptionalQueryStructField(
				"PolarisParams",
				g.NewQueryStruct("PolarisParams").
					PredefinedQueryStructField("catalogSource", "string", g.StaticOptions().SQL("CATALOG_SOURCE = POLARIS")).
					TextAssignment("TABLE_FORMAT", g.ParameterOptions().NoQuotes().Required()),
				g.KeywordOptions(),
			).
			OptionalQueryStructField(
				"SapBdcParams",
				g.NewQueryStruct("SapBdcParams").
					PredefinedQueryStructField("catalogSource", "string", g.StaticOptions().SQL("CATALOG_SOURCE = SAP_BDC")).
					TextAssignment("TABLE_FORMAT", g.ParameterOptions().NoQuotes().Required()),
				g.KeywordOptions(),
			).
			OptionalTextAssignment("CATALOG_NAMESPACE", g.ParameterOptions().SingleQuotes()).
			OptionalComment().
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ConflictingFields, "IfNotExists", "OrReplace").
			WithValidation(g.ExactlyOneValueSet, "ObjectStoreParams", "GlueParams", "IcebergRestParams", "PolarisParams", "SapBdcParams"),
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
					OptionalBooleanAssignment("ENABLED", g.ParameterOptions()).
					OptionalTextAssignment("GLUE_AWS_ROLE_ARN", g.ParameterOptions().SingleQuotes()).
					OptionalTextAssignment("GLUE_CATALOG_ID", g.ParameterOptions().SingleQuotes()).
					OptionalTextAssignment("GLUE_REGION", g.ParameterOptions().SingleQuotes()).
					OptionalTextAssignment("CATALOG_NAMESPACE", g.ParameterOptions().SingleQuotes()).
					OptionalComment().
					WithValidation(g.AtLeastOneValueSet, "Enabled", "GlueAwsRoleArn", "GlueCatalogId", "GlueRegion", "CatalogNamespace", "Comment"),
				g.KeywordOptions().SQL("SET"),
			).
			OptionalQueryStructField(
				"Unset",
				g.NewQueryStruct("CatalogIntegrationUnset").
					OptionalSQL("ENABLED").
					OptionalSQL("CATALOG_NAMESPACE").
					OptionalSQL("COMMENT").
					WithValidation(g.AtLeastOneValueSet, "Enabled", "CatalogNamespace", "Comment"),
				g.ListOptions().NoParentheses().SQL("UNSET"),
			).
			OptionalSetTags().
			OptionalUnsetTags().
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ExactlyOneValueSet, "Set", "Unset", "SetTags", "UnsetTags"),
	).
	DropOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/drop-integration",
		g.NewQueryStruct("DropCatalogIntegration").
			Drop().
			SQL("CATALOG INTEGRATION").
			IfExists().
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	).
	ShowOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/show-integrations",
		g.DbStruct("showCatalogIntegrationsDbRow").
			Text("name").
			Text("type").
			Text("category").
			Bool("enabled").
			OptionalText("comment").
			Time("created_on"),
		g.PlainStruct("CatalogIntegration").
			Text("Name").
			Text("Type").
			Text("Category").
			Bool("Enabled").
			Text("Comment").
			Time("CreatedOn"),
		g.NewQueryStruct("ShowCatalogIntegrations").
			Show().
			SQL("CATALOG INTEGRATIONS").
			OptionalLike(),
	).
	ShowByIdOperationWithFiltering(
		g.ShowByIDLikeFiltering,
	).
	DescribeOperation(
		g.DescriptionMappingKindSlice,
		"https://docs.snowflake.com/en/sql-reference/sql/desc-integration",
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
			WithValidation(g.ValidIdentifier, "name"),
	)
