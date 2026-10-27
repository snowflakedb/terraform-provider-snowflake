package defs

import (
	ssotdefs "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot/defs"
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var (
	CatalogIntegrationCatalogSourceTypeEnumDef = g.NewEnum(
		"CatalogIntegrationCatalogSourceType", "CatalogIntegrationCatalogSourceTypes",
		"GLUE", "OBJECT_STORE", "POLARIS", "ICEBERG_REST", "SAP_BDC",
	)
	CatalogIntegrationTableFormatEnumDef = g.NewEnum(
		"CatalogIntegrationTableFormat", "CatalogIntegrationTableFormats",
		"ICEBERG", "DELTA",
	)
	CatalogIntegrationRestAuthenticationTypeEnumDef = g.NewEnum(
		"CatalogIntegrationRestAuthenticationType", "CatalogIntegrationRestAuthenticationTypes",
		"OAUTH", "BEARER", "SIGV4",
	)
	CatalogIntegrationAccessDelegationModeEnumDef = g.NewEnum(
		"CatalogIntegrationAccessDelegationMode", "CatalogIntegrationAccessDelegationModes",
		"VENDED_CREDENTIALS", "EXTERNAL_VOLUME_CREDENTIALS",
	)
	CatalogIntegrationCatalogApiTypeEnumDef = g.NewEnum(
		"CatalogIntegrationCatalogApiType", "CatalogIntegrationCatalogApiTypes",
		"PUBLIC", "PRIVATE", "AWS_API_GATEWAY", "AWS_PRIVATE_API_GATEWAY", "AWS_GLUE", "AWS_PRIVATE_GLUE",
	)
)

func catalogIntegrationDetailsHead(ps *g.Plain) *g.Plain {
	return ps.
		AccountObjectIdentifier().
		Enum("CatalogSource", CatalogIntegrationCatalogSourceTypeEnumDef).
		Enum("TableFormat", CatalogIntegrationTableFormatEnumDef).
		Bool("Enabled").
		Number("RefreshIntervalSeconds").
		Text("Comment")
}

var catalogIntegrationsDef = g.NewInterface(
	"CatalogIntegrations",
	"CatalogIntegration",
	g.KindOfT[sdkcommons.AccountObjectIdentifier](),
).
	CreateVariantOperationsFromObject(ssotdefs.CatalogIntegration).
	AlterOperationFromObject(ssotdefs.CatalogIntegration).
	DropOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/drop-catalog-integration",
		g.NewQueryStruct("DropCatalogIntegration").
			Drop().
			SQL("CATALOG INTEGRATION").
			IfExists().
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	).
	ShowOperationWithPairedStructs(
		"https://docs.snowflake.com/en/sql-reference/sql/show-catalog-integrations",
		g.StructPair("showCatalogIntegrationsDbRow", "CatalogIntegration").
			Text("name").
			Bool("enabled").
			Text("type").
			Text("category").
			OptionalText("comment", g.WithRequiredInPlain()).
			Time("created_on"),
		g.NewQueryStruct("ShowCatalogIntegration").
			Show().
			SQL("CATALOG INTEGRATIONS").
			OptionalLike(),
	).
	DescribeOperationWithPairedStructs(
		g.DescriptionMappingKindSlice,
		"https://docs.snowflake.com/en/sql-reference/sql/desc-catalog-integration",
		g.StructPair("descCatalogIntegrationsDbRow", "CatalogIntegrationProperty").
			Text("property", g.WithPlainFieldName("Name")).
			Text("property_type", g.WithPlainFieldName("Type")).
			Text("property_value", g.WithPlainFieldName("Value")).
			Text("property_default", g.WithPlainFieldName("Default")),
		g.NewQueryStruct("DescribeCatalogIntegration").
			Describe().
			SQL("CATALOG INTEGRATION").
			Name().
			WithValidation(g.ValidIdentifier, "name"),
		catalogIntegrationDetailsHead(g.PlainStruct("CatalogIntegrationAwsGlueDetails")).
			Text("GlueAwsRoleArn").
			Text("GlueCatalogId").
			Text("GlueRegion").
			Text("CatalogNamespace").
			Text("GlueAwsIamUserArn").
			Text("GlueAwsExternalId"),
		catalogIntegrationDetailsHead(g.PlainStruct("CatalogIntegrationObjectStorageDetails")),
		catalogIntegrationDetailsHead(g.PlainStruct("CatalogIntegrationOpenCatalogDetails")).
			Text("CatalogNamespace").
			Field("RestConfig", "OpenCatalogRestConfigDetails").
			Field("RestAuthentication", "OAuthRestAuthenticationDetails"),
		catalogIntegrationDetailsHead(g.PlainStruct("CatalogIntegrationIcebergRestDetails")).
			Text("CatalogNamespace").
			Field("RestConfig", "IcebergRestRestConfigDetails").
			OptionalField("OAuthRestAuthentication", "OAuthRestAuthenticationDetails").
			OptionalField("BearerRestAuthentication", "BearerRestAuthenticationDetails").
			OptionalField("SigV4RestAuthentication", "SigV4RestAuthenticationDetails"),
		catalogIntegrationDetailsHead(g.PlainStruct("CatalogIntegrationAllDetails")).
			Text("GlueAwsRoleArn").
			Text("GlueCatalogId").
			Text("GlueRegion").
			Text("CatalogNamespace").
			Text("GlueAwsIamUserArn").
			Text("GlueAwsExternalId").
			// IcebergRestRestConfigDetails contains all properties of OpenCatalogRestConfigDetails
			OptionalField("RestConfig", "IcebergRestRestConfigDetails").
			OptionalField("OAuthRestAuthentication", "OAuthRestAuthenticationDetails").
			OptionalField("BearerRestAuthentication", "BearerRestAuthenticationDetails").
			OptionalField("SigV4RestAuthentication", "SigV4RestAuthenticationDetails"),
		g.PlainStruct("OpenCatalogRestConfigDetails").
			Text("CatalogUri").
			Enum("CatalogApiType", CatalogIntegrationCatalogApiTypeEnumDef).
			Text("CatalogName").
			Enum("AccessDelegationMode", CatalogIntegrationAccessDelegationModeEnumDef),
		g.PlainStruct("IcebergRestRestConfigDetails").
			Text("CatalogUri").
			Text("Prefix").
			Text("CatalogName").
			Enum("CatalogApiType", CatalogIntegrationCatalogApiTypeEnumDef).
			Enum("AccessDelegationMode", CatalogIntegrationAccessDelegationModeEnumDef),
		g.PlainStruct("OAuthRestAuthenticationDetails").
			Text("OauthTokenUri").
			Text("OauthClientId").
			Text("OauthClientSecret").
			StringList("OauthAllowedScopes"),
		g.PlainStruct("BearerRestAuthenticationDetails").
			Text("BearerToken"),
		g.PlainStruct("SigV4RestAuthenticationDetails").
			Text("Sigv4IamRole").
			Text("Sigv4SigningRegion").
			Text("Sigv4ExternalId"),
	).
	WithCustomInterfaceMethod(
		"DescribeAwsGlueDetails",
		"",
		[]*g.MethodParameter{g.NewMethodParameter("id", g.KindOfT[sdkcommons.AccountObjectIdentifier]())},
		"*CatalogIntegrationAwsGlueDetails", "error",
	).
	WithCustomInterfaceMethod(
		"DescribeObjectStorageDetails",
		"",
		[]*g.MethodParameter{g.NewMethodParameter("id", g.KindOfT[sdkcommons.AccountObjectIdentifier]())},
		"*CatalogIntegrationObjectStorageDetails", "error",
	).
	WithCustomInterfaceMethod(
		"DescribeOpenCatalogDetails",
		"",
		[]*g.MethodParameter{g.NewMethodParameter("id", g.KindOfT[sdkcommons.AccountObjectIdentifier]())},
		"*CatalogIntegrationOpenCatalogDetails", "error",
	).
	WithCustomInterfaceMethod(
		"DescribeIcebergRestDetails",
		"",
		[]*g.MethodParameter{g.NewMethodParameter("id", g.KindOfT[sdkcommons.AccountObjectIdentifier]())},
		"*CatalogIntegrationIcebergRestDetails", "error",
	).
	WithCustomInterfaceMethod(
		"DescribeDetails",
		"DescribeDetails returns combined describe output for all types of catalog integrations.",
		[]*g.MethodParameter{g.NewMethodParameter("id", g.KindOfT[sdkcommons.AccountObjectIdentifier]())},
		"*CatalogIntegrationAllDetails", "error",
	).
	WithEnums(
		CatalogIntegrationCatalogSourceTypeEnumDef,
		CatalogIntegrationTableFormatEnumDef,
		CatalogIntegrationRestAuthenticationTypeEnumDef,
		CatalogIntegrationAccessDelegationModeEnumDef,
		CatalogIntegrationCatalogApiTypeEnumDef,
	)
