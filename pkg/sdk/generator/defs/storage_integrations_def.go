package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var StorageLocationDef = g.NewQueryStruct("StorageLocation").Text("Path", g.KeywordOptions().SingleQuotes().Required())

var storageIntegrationsDef = g.NewInterface(
	"StorageIntegrations",
	"StorageIntegration",
	g.KindOfT[sdkcommons.AccountObjectIdentifier](),
).
	CreateOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/create-storage-integration",
		g.NewQueryStruct("CreateStorageIntegration").
			Create().
			OrReplace().
			SQL("STORAGE INTEGRATION").
			IfNotExists().
			Name().
			SQLWithCustomFieldName("externalStageType", "TYPE = EXTERNAL_STAGE").
			OptionalQueryStructField(
				"S3StorageProviderParams",
				g.NewQueryStruct("S3StorageParams").
					WithField(g.EnumLegacy[sdkcommons.S3Protocol]("Protocol", g.ParameterOptions().SQL("STORAGE_PROVIDER").SingleQuotes().Required())).
					TextAssignment("STORAGE_AWS_ROLE_ARN", g.ParameterOptions().SingleQuotes().Required()).
					OptionalTextAssignment("STORAGE_AWS_EXTERNAL_ID", g.ParameterOptions().SingleQuotes()).
					OptionalTextAssignment("STORAGE_AWS_OBJECT_ACL", g.ParameterOptions().SingleQuotes()).
					OptionalBooleanAssignment("USE_PRIVATELINK_ENDPOINT", g.ParameterOptions()),
				g.KeywordOptions(),
			).
			OptionalQueryStructField(
				"GCSStorageProviderParams",
				g.NewQueryStruct("GCSStorageParams").
					SQLWithCustomFieldName("storageProvider", "STORAGE_PROVIDER = 'GCS'"),
				g.KeywordOptions(),
			).
			OptionalQueryStructField(
				"AzureStorageProviderParams",
				g.NewQueryStruct("AzureStorageParams").
					SQLWithCustomFieldName("storageProvider", "STORAGE_PROVIDER = 'AZURE'").
					TextAssignment("AZURE_TENANT_ID", g.ParameterOptions().SingleQuotes().Required()).
					OptionalBooleanAssignment("USE_PRIVATELINK_ENDPOINT", g.ParameterOptions()),
				g.KeywordOptions(),
			).
			// Enabled is required even though it can be UNSET. Not using it in create results in:
			// 002029 (42601): SQL compilation error: Missing option(s): ENABLED
			BooleanAssignment("ENABLED", g.ParameterOptions().Required()).
			ListAssignment("STORAGE_ALLOWED_LOCATIONS", "StorageLocation", g.ParameterOptions().Parentheses().Required()).
			ListAssignment("STORAGE_BLOCKED_LOCATIONS", "StorageLocation", g.ParameterOptions().Parentheses()).
			OptionalComment().
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ConflictingFields, "IfNotExists", "OrReplace").
			WithValidation(g.ExactlyOneValueSet, "S3StorageProviderParams", "GCSStorageProviderParams", "AzureStorageProviderParams"),
		StorageLocationDef,
	).
	AlterOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/alter-storage-integration",
		g.NewQueryStruct("AlterStorageIntegration").
			Alter().
			SQL("STORAGE INTEGRATION").
			IfExists().
			Name().
			OptionalQueryStructField(
				"Set",
				g.NewQueryStruct("StorageIntegrationSet").
					OptionalQueryStructField(
						"S3Params",
						g.NewQueryStruct("SetS3StorageParams").
							OptionalTextAssignment("STORAGE_AWS_ROLE_ARN", g.ParameterOptions().SingleQuotes()).
							OptionalTextAssignment("STORAGE_AWS_EXTERNAL_ID", g.ParameterOptions().SingleQuotes()).
							OptionalTextAssignment("STORAGE_AWS_OBJECT_ACL", g.ParameterOptions().SingleQuotes()).
							OptionalBooleanAssignment("USE_PRIVATELINK_ENDPOINT", g.ParameterOptions()).
							WithValidation(g.AtLeastOneValueSet, "StorageAwsRoleArn", "StorageAwsExternalId", "StorageAwsObjectAcl", "UsePrivatelinkEndpoint"),
						g.KeywordOptions(),
					).
					OptionalQueryStructField(
						"AzureParams",
						g.NewQueryStruct("SetAzureStorageParams").
							OptionalTextAssignment("AZURE_TENANT_ID", g.ParameterOptions().SingleQuotes()).
							OptionalBooleanAssignment("USE_PRIVATELINK_ENDPOINT", g.ParameterOptions()).
							WithValidation(g.AtLeastOneValueSet, "AzureTenantId", "UsePrivatelinkEndpoint"),
						g.KeywordOptions(),
					).
					OptionalBooleanAssignment("ENABLED", g.ParameterOptions()).
					ListAssignment("STORAGE_ALLOWED_LOCATIONS", "StorageLocation", g.ParameterOptions().Parentheses()).
					ListAssignment("STORAGE_BLOCKED_LOCATIONS", "StorageLocation", g.ParameterOptions().Parentheses()).
					OptionalComment().
					WithValidation(g.ConflictingFields, "S3Params", "AzureParams").
					WithValidation(g.AtLeastOneValueSet, "S3Params", "AzureParams", "Enabled", "StorageAllowedLocations", "StorageBlockedLocations", "Comment"),
				g.KeywordOptions().SQL("SET"),
			).
			OptionalQueryStructField(
				"Unset",
				g.NewQueryStruct("StorageIntegrationUnset").
					OptionalQueryStructField(
						"S3Params",
						g.NewQueryStruct("UnsetS3StorageParams").
							OptionalSQL("STORAGE_AWS_EXTERNAL_ID").
							OptionalSQL("STORAGE_AWS_OBJECT_ACL").
							OptionalSQL("USE_PRIVATELINK_ENDPOINT").
							WithValidation(g.AtLeastOneValueSet, "StorageAwsExternalId", "StorageAwsObjectAcl", "UsePrivatelinkEndpoint"),
						g.ListOptions(),
					).
					OptionalQueryStructField(
						"AzureParams",
						g.NewQueryStruct("UnsetAzureStorageParams").
							OptionalSQL("USE_PRIVATELINK_ENDPOINT").
							WithValidation(g.AtLeastOneValueSet, "UsePrivatelinkEndpoint"),
						g.ListOptions(),
					).
					OptionalSQL("ENABLED").
					OptionalSQL("STORAGE_BLOCKED_LOCATIONS").
					OptionalSQL("COMMENT").
					WithValidation(g.ConflictingFields, "S3Params", "AzureParams").
					WithValidation(g.AtLeastOneValueSet, "S3Params", "AzureParams", "Enabled", "StorageBlockedLocations", "Comment"),
				g.ListOptions().SQL("UNSET"),
			).
			OptionalSetTags().
			OptionalUnsetTags().
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ConflictingFields, "IfExists", "UnsetTags").
			WithValidation(g.ExactlyOneValueSet, "Set", "Unset", "SetTags", "UnsetTags"),
	).
	DropOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/drop-integration",
		g.NewQueryStruct("DropStorageIntegration").
			Drop().
			SQL("STORAGE INTEGRATION").
			IfExists().
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	).
	ShowOperationWithPairedStructs(
		"https://docs.snowflake.com/en/sql-reference/sql/show-integrations",
		g.StructPair("showStorageIntegrationsDbRow", "StorageIntegration").
			Text("name").
			Text("type", g.WithPlainFieldName("StorageType")).
			Text("category").
			Bool("enabled").
			OptionalText("comment", g.WithRequiredInPlain()).
			Time("created_on").
			WithConvertGeneration(),
		g.NewQueryStruct("ShowStorageIntegrations").
			Show().
			SQL("STORAGE INTEGRATIONS").
			OptionalLike(),
	).
	DescribeOperationWithPairedStructs(
		g.DescriptionMappingKindSlice,
		"https://docs.snowflake.com/en/sql-reference/sql/desc-integration",
		g.StructPair("descStorageIntegrationsDbRow", "StorageIntegrationProperty").
			Text("property", g.WithPlainFieldName("Name")).
			Text("property_type", g.WithPlainFieldName("Type")).
			Text("property_value", g.WithPlainFieldName("Value")).
			Text("property_default", g.WithPlainFieldName("Default")).
			WithConvertGeneration(),
		g.NewQueryStruct("DescribeStorageIntegration").
			Describe().
			SQL("STORAGE INTEGRATION").
			Name().
			WithValidation(g.ValidIdentifier, "name"),
		g.PlainStruct("StorageIntegrationAwsDetails").
			AccountObjectIdentifier().
			Bool("Enabled").
			// TODO [next PRs]: enum?
			Text("Provider").
			StringList("AllowedLocations").
			StringList("BlockedLocations").
			Text("Comment").
			Bool("UsePrivatelinkEndpoint").
			Text("IamUserArn").
			Text("RoleArn").
			Text("ObjectAcl").
			Text("ExternalId"),
		g.PlainStruct("StorageIntegrationAzureDetails").
			AccountObjectIdentifier().
			Bool("Enabled").
			Text("Provider").
			StringList("AllowedLocations").
			StringList("BlockedLocations").
			Text("Comment").
			Bool("UsePrivatelinkEndpoint").
			Text("TenantId").
			Text("ConsentUrl").
			Text("MultiTenantAppName"),
		g.PlainStruct("StorageIntegrationGcsDetails").
			AccountObjectIdentifier().
			Bool("Enabled").
			Text("Provider").
			StringList("AllowedLocations").
			StringList("BlockedLocations").
			Text("Comment").
			Bool("UsePrivatelinkEndpoint").
			Text("ServiceAccount"),
		g.PlainStruct("StorageIntegrationAllDetails").
			AccountObjectIdentifier().
			Bool("Enabled").
			Text("Provider").
			StringList("AllowedLocations").
			StringList("BlockedLocations").
			Text("Comment").
			Bool("UsePrivatelinkEndpoint").
			Text("IamUserArn").
			Text("RoleArn").
			Text("ObjectAcl").
			Text("ExternalId").
			Text("TenantId").
			Text("ConsentUrl").
			Text("MultiTenantAppName").
			Text("ServiceAccount"),
	).
	WithCustomInterfaceMethod(
		"DescribeAwsDetails",
		"DescribeAwsDetails returns converted describe output for AWS storage integrations.",
		[]*g.MethodParameter{g.NewMethodParameter("id", g.KindOfT[sdkcommons.AccountObjectIdentifier]())},
		"*StorageIntegrationAwsDetails", "error",
	).
	WithCustomInterfaceMethod(
		"DescribeAzureDetails",
		"DescribeAzureDetails returns converted describe output for Azure storage integrations.",
		[]*g.MethodParameter{g.NewMethodParameter("id", g.KindOfT[sdkcommons.AccountObjectIdentifier]())},
		"*StorageIntegrationAzureDetails", "error",
	).
	WithCustomInterfaceMethod(
		"DescribeGcsDetails",
		"DescribeGcsDetails returns converted describe output for GCS storage integrations.",
		[]*g.MethodParameter{g.NewMethodParameter("id", g.KindOfT[sdkcommons.AccountObjectIdentifier]())},
		"*StorageIntegrationGcsDetails", "error",
	).
	WithCustomInterfaceMethod(
		"DescribeDetails",
		"DescribeDetails returns combined describe output for all types of storage integrations.",
		[]*g.MethodParameter{g.NewMethodParameter("id", g.KindOfT[sdkcommons.AccountObjectIdentifier]())},
		"*StorageIntegrationAllDetails", "error",
	)
