package sdk

import (
	"fmt"
	"strings"

	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"
)

//go:generate go run ./poc/main.go

type S3Protocol string

const (
	RegularS3Protocol S3Protocol = "S3"
	GovS3Protocol     S3Protocol = "S3GOV"
	ChinaS3Protocol   S3Protocol = "S3CHINA"
)

var (
	AllS3Protocols      = []S3Protocol{RegularS3Protocol, GovS3Protocol, ChinaS3Protocol}
	AllStorageProviders = append(AsStringList(AllS3Protocols), "GCS", "AZURE")
)

func ToS3Protocol(s string) (S3Protocol, error) {
	switch protocol := S3Protocol(strings.ToUpper(s)); protocol {
	case RegularS3Protocol, GovS3Protocol, ChinaS3Protocol:
		return protocol, nil
	default:
		return "", fmt.Errorf("invalid S3 protocol: %s", s)
	}
}

var StorageLocationDef = g.NewQueryStruct("StorageLocation").Text("Path", g.KeywordOptions().SingleQuotes().Required())

var StorageIntegrationDef = g.NewInterface(
	"StorageIntegrations",
	"StorageIntegration",
	g.KindOfT[AccountObjectIdentifier](),
).
	CreateOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/create-storage-integration",
		g.NewQueryStruct("CreateStorageIntegration").
			Create().
			OrReplace().
			SQL("STORAGE INTEGRATION").
			IfNotExists().
			Name().
			PredefinedQueryStructField("externalStageType", "string", g.StaticOptions().SQL("TYPE = EXTERNAL_STAGE")).
			OptionalQueryStructField(
				"S3StorageProviderParams",
				g.NewQueryStruct("S3StorageParams").
					PredefinedQueryStructField("Protocol", g.KindOfT[S3Protocol](), g.ParameterOptions().SQL("STORAGE_PROVIDER").SingleQuotes().Required()).
					TextAssignment("STORAGE_AWS_ROLE_ARN", g.ParameterOptions().SingleQuotes().Required()).
					OptionalTextAssignment("STORAGE_AWS_EXTERNAL_ID", g.ParameterOptions().SingleQuotes()).
					OptionalTextAssignment("STORAGE_AWS_OBJECT_ACL", g.ParameterOptions().SingleQuotes()).
					OptionalBooleanAssignment("USE_PRIVATELINK_ENDPOINT", g.ParameterOptions()),
				g.KeywordOptions(),
			).
			OptionalQueryStructField(
				"GCSStorageProviderParams",
				g.NewQueryStruct("GCSStorageParams").
					PredefinedQueryStructField("storageProvider", "string", g.StaticOptions().SQL("STORAGE_PROVIDER = 'GCS'")),
				g.KeywordOptions(),
			).
			OptionalQueryStructField(
				"AzureStorageProviderParams",
				g.NewQueryStruct("AzureStorageParams").
					PredefinedQueryStructField("storageProvider", "string", g.StaticOptions().SQL("STORAGE_PROVIDER = 'AZURE'")).
					OptionalTextAssignment("AZURE_TENANT_ID", g.ParameterOptions().SingleQuotes().Required()).
					OptionalBooleanAssignment("USE_PRIVATELINK_ENDPOINT", g.ParameterOptions()),
				g.KeywordOptions(),
			).
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
							TextAssignment("STORAGE_AWS_ROLE_ARN", g.ParameterOptions().SingleQuotes().Required()).
							OptionalTextAssignment("STORAGE_AWS_EXTERNAL_ID", g.ParameterOptions().SingleQuotes()).
							OptionalTextAssignment("STORAGE_AWS_OBJECT_ACL", g.ParameterOptions().SingleQuotes()).
							OptionalBooleanAssignment("USE_PRIVATELINK_ENDPOINT", g.ParameterOptions()),
						g.KeywordOptions(),
					).
					OptionalQueryStructField(
						"AzureParams",
						g.NewQueryStruct("SetAzureStorageParams").
							TextAssignment("AZURE_TENANT_ID", g.ParameterOptions().SingleQuotes().Required()).
							OptionalBooleanAssignment("USE_PRIVATELINK_ENDPOINT", g.ParameterOptions()),
						g.KeywordOptions(),
					).
					OptionalBooleanAssignment("ENABLED", g.ParameterOptions()).
					ListAssignment("STORAGE_ALLOWED_LOCATIONS", "StorageLocation", g.ParameterOptions().Parentheses()).
					ListAssignment("STORAGE_BLOCKED_LOCATIONS", "StorageLocation", g.ParameterOptions().Parentheses()).
					OptionalComment(),
				g.KeywordOptions().SQL("SET"),
			).
			OptionalQueryStructField(
				"Unset",
				g.NewQueryStruct("StorageIntegrationUnset").
					OptionalSQL("STORAGE_AWS_EXTERNAL_ID").
					OptionalSQL("STORAGE_AWS_OBJECT_ACL").
					OptionalSQL("ENABLED").
					OptionalSQL("STORAGE_BLOCKED_LOCATIONS").
					OptionalSQL("COMMENT").
					OptionalSQL("USE_PRIVATELINK_ENDPOINT"),
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
	ShowOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/show-integrations",
		g.DbStruct("showStorageIntegrationsDbRow").
			Text("name").
			Text("type").
			Text("category").
			Bool("enabled").
			OptionalText("comment").
			Time("created_on"),
		g.PlainStruct("StorageIntegration").
			Text("Name").
			Text("StorageType").
			Text("Category").
			Bool("Enabled").
			Text("Comment").
			Time("CreatedOn"),
		g.NewQueryStruct("ShowStorageIntegrations").
			Show().
			SQL("STORAGE INTEGRATIONS").
			OptionalLike(),
	).
	ShowByIdOperationWithFiltering(
		g.ShowByIDLikeFiltering,
	).
	DescribeOperation(
		g.DescriptionMappingKindSlice,
		"https://docs.snowflake.com/en/sql-reference/sql/desc-integration",
		g.DbStruct("descStorageIntegrationsDbRow").
			Text("property").
			Text("property_type").
			Text("property_value").
			Text("property_default"),
		g.PlainStruct("StorageIntegrationProperty").
			Text("Name").
			Text("Type").
			Text("Value").
			Text("Default"),
		g.NewQueryStruct("DescribeStorageIntegration").
			Describe().
			SQL("STORAGE INTEGRATION").
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	)
