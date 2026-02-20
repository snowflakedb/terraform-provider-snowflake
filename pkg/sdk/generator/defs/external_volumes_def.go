package defs

import (
	"fmt"

	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var externalS3StorageLocationDef = g.NewQueryStruct("S3StorageLocationParams").
	TextAssignment("NAME", g.ParameterOptions().SingleQuotes().Required()).
	Assignment("STORAGE_PROVIDER", g.KindOfT[sdkcommons.S3StorageProvider](), g.ParameterOptions().SingleQuotes().Required()).
	TextAssignment("STORAGE_AWS_ROLE_ARN", g.ParameterOptions().SingleQuotes().Required()).
	TextAssignment("STORAGE_BASE_URL", g.ParameterOptions().SingleQuotes().Required()).
	OptionalTextAssignment("STORAGE_AWS_EXTERNAL_ID", g.ParameterOptions().SingleQuotes()).
	OptionalTextAssignment("STORAGE_AWS_ACCESS_POINT_ARN", g.ParameterOptions().SingleQuotes()).
	OptionalBooleanAssignment("USE_PRIVATELINK_ENDPOINT", g.ParameterOptions()).
	OptionalQueryStructField(
		"Encryption",
		g.NewQueryStruct("ExternalVolumeS3Encryption").
			AssignmentWithFieldName("TYPE", g.KindOfT[sdkcommons.S3EncryptionType](), g.ParameterOptions().SingleQuotes().Required(), "EncryptionType").
			OptionalTextAssignment("KMS_KEY_ID", g.ParameterOptions().SingleQuotes()),
		g.ListOptions().Parentheses().NoComma().SQL("ENCRYPTION ="),
	)

var externalGCSStorageLocationDef = g.NewQueryStruct("GCSStorageLocationParams").
	TextAssignment("NAME", g.ParameterOptions().SingleQuotes().Required()).
	PredefinedQueryStructField("StorageProviderGcs", "string", g.StaticOptions().SQL(fmt.Sprintf("STORAGE_PROVIDER = '%s'", sdkcommons.StorageProviderGCS))).
	TextAssignment("STORAGE_BASE_URL", g.ParameterOptions().SingleQuotes().Required()).
	OptionalQueryStructField(
		"Encryption",
		g.NewQueryStruct("ExternalVolumeGCSEncryption").
			AssignmentWithFieldName("TYPE", g.KindOfT[sdkcommons.GCSEncryptionType](), g.ParameterOptions().SingleQuotes().Required(), "EncryptionType").
			OptionalTextAssignment("KMS_KEY_ID", g.ParameterOptions().SingleQuotes()),
		g.ListOptions().Parentheses().NoComma().SQL("ENCRYPTION ="),
	)

var externalAzureStorageLocationDef = g.NewQueryStruct("AzureStorageLocationParams").
	TextAssignment("NAME", g.ParameterOptions().SingleQuotes().Required()).
	PredefinedQueryStructField("StorageProviderAzure", "string", g.StaticOptions().SQL(fmt.Sprintf("STORAGE_PROVIDER = '%s'", sdkcommons.StorageProviderAzure))).
	TextAssignment("AZURE_TENANT_ID", g.ParameterOptions().SingleQuotes().Required()).
	TextAssignment("STORAGE_BASE_URL", g.ParameterOptions().SingleQuotes().Required()).
	OptionalBooleanAssignment("USE_PRIVATELINK_ENDPOINT", g.ParameterOptions())

var externalS3CompatStorageLocationDef = g.NewQueryStruct("S3CompatStorageLocationParams").
	TextAssignment("NAME", g.ParameterOptions().SingleQuotes().Required()).
	PredefinedQueryStructField("StorageProviderS3Compat", "string", g.StaticOptions().SQL("STORAGE_PROVIDER = 'S3COMPAT'")).
	TextAssignment("STORAGE_BASE_URL", g.ParameterOptions().SingleQuotes().Required()).
	TextAssignment("STORAGE_ENDPOINT", g.ParameterOptions().SingleQuotes().Required()).
	OptionalQueryStructField(
		"Credentials",
		g.NewQueryStruct("ExternalVolumeS3CompatCredentials").
			TextAssignment("AWS_KEY_ID", g.ParameterOptions().SingleQuotes().Required()).
			TextAssignment("AWS_SECRET_KEY", g.ParameterOptions().SingleQuotes().Required()),
		g.ListOptions().Parentheses().NoComma().SQL("CREDENTIALS ="),
	)

// Can't name StorageLocation due to naming clash with type in storage integration
var storageLocationDef = g.NewQueryStruct("ExternalVolumeStorageLocation").
	OptionalQueryStructField(
		"S3StorageLocationParams",
		externalS3StorageLocationDef,
		g.ListOptions().Parentheses().NoComma(),
	).
	OptionalQueryStructField(
		"GCSStorageLocationParams",
		externalGCSStorageLocationDef,
		g.ListOptions().Parentheses().NoComma(),
	).
	OptionalQueryStructField(
		"AzureStorageLocationParams",
		externalAzureStorageLocationDef,
		g.ListOptions().Parentheses().NoComma(),
	).
	OptionalQueryStructField(
		"S3CompatStorageLocationParams",
		externalS3CompatStorageLocationDef,
		g.ListOptions().Parentheses().NoComma(),
	).
	WithValidation(g.ExactlyOneValueSet, "S3StorageLocationParams", "GCSStorageLocationParams", "AzureStorageLocationParams", "S3CompatStorageLocationParams")

var updateStorageLocationDef = g.NewQueryStruct("AlterExternalVolumeUpdateStorageLocation").
	TextAssignment("STORAGE_LOCATION", g.ParameterOptions().SingleQuotes().NoEquals().Required()).
	OptionalQueryStructField(
		"Credentials",
		g.NewQueryStruct("ExternalVolumeUpdateCredentials").
			TextAssignment("AWS_KEY_ID", g.ParameterOptions().SingleQuotes().Required()).
			TextAssignment("AWS_SECRET_KEY", g.ParameterOptions().SingleQuotes().Required()),
		g.ListOptions().Parentheses().NoComma().SQL("CREDENTIALS ="),
	)

var externalVolumesDef = g.NewInterface(
	"ExternalVolumes",
	"ExternalVolume",
	g.KindOfT[sdkcommons.AccountObjectIdentifier](),
).
	CreateOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/create-external-volume",
		g.NewQueryStruct("CreateExternalVolume").
			Create().
			OrReplace().
			SQL("EXTERNAL VOLUME").
			IfNotExists().
			Name().
			ListAssignment("STORAGE_LOCATIONS", "ExternalVolumeStorageLocation", g.ParameterOptions().Parentheses().Required()).
			OptionalBooleanAssignment("ALLOW_WRITES", nil).
			OptionalComment().
			WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists").
			WithValidation(g.ValidIdentifier, "name"),
		storageLocationDef,
	).
	AlterOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/alter-external-volume",
		g.NewQueryStruct("AlterExternalVolume").
			Alter().
			SQL("EXTERNAL VOLUME").
			IfExists().
			Name().
			OptionalTextAssignment("REMOVE STORAGE_LOCATION", g.ParameterOptions().SingleQuotes().NoEquals()).
			OptionalQueryStructField(
				"Set",
				g.NewQueryStruct("AlterExternalVolumeSet").
					OptionalBooleanAssignment("ALLOW_WRITES", g.ParameterOptions()).
					OptionalComment(),
				g.KeywordOptions().SQL("SET"),
			).
			OptionalQueryStructField(
				"AddStorageLocation",
				storageLocationDef,
				g.ParameterOptions().SQL("ADD STORAGE_LOCATION"),
			).
			OptionalQueryStructField(
				"UpdateStorageLocation",
				updateStorageLocationDef,
				g.KeywordOptions().SQL("UPDATE"),
			).
			WithValidation(g.ExactlyOneValueSet, "RemoveStorageLocation", "Set", "AddStorageLocation", "UpdateStorageLocation").
			WithValidation(g.ValidIdentifier, "name"),
	).
	DropOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/drop-external-volume",
		g.NewQueryStruct("DropExternalVolume").
			Drop().
			SQL("EXTERNAL VOLUME").
			IfExists().
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	).
	DescribeOperation(
		g.DescriptionMappingKindSlice,
		"https://docs.snowflake.com/en/sql-reference/sql/desc-external-volume",
		g.DbStruct("externalVolumeDescRow").
			Text("parent_property").
			Text("property").
			Text("property_type").
			Text("property_value").
			Text("property_default"),
		g.PlainStruct("ExternalVolumeProperty").
			Text("Parent").
			Text("Name").
			Text("Type").
			Text("Value").
			Text("Default"),
		g.NewQueryStruct("DescExternalVolume").
			Describe().
			SQL("EXTERNAL VOLUME").
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	).
	ShowOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/show-external-volumes",
		g.DbStruct("externalVolumeShowRow").
			Text("name").
			Bool("allow_writes").
			OptionalText("comment"),
		g.PlainStruct("ExternalVolume").
			Text("Name").
			Bool("AllowWrites").
			Text("Comment"),
		g.NewQueryStruct("ShowExternalVolumes").
			Show().
			SQL("EXTERNAL VOLUMES").
			OptionalLike(),
	).
	ShowByIdOperationWithFiltering(
		g.ShowByIDLikeFiltering,
	)
