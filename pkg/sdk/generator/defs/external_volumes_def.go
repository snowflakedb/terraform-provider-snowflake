package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var (
	StorageProviderEnumDef = g.NewEnum(
		"StorageProvider", "StorageProviders",
		"GCS", "AZURE", "S3", "S3GOV", "S3COMPAT",
	)
	S3StorageProviderEnumDef = g.NewEnum(
		"S3StorageProvider", "S3StorageProviders",
		"S3", "S3GOV",
	)
	S3EncryptionTypeEnumDef = g.NewEnum(
		"S3EncryptionType", "S3EncryptionTypes",
		"AWS_SSE_S3", "AWS_SSE_KMS", "NONE",
	)
	GCSEncryptionTypeEnumDef = g.NewEnum(
		"GCSEncryptionType", "GCSEncryptionTypes",
		"GCS_SSE_KMS", "NONE",
	)
	AzureEncryptionTypeEnumDef = g.NewEnum(
		"AzureEncryptionType", "AzureEncryptionTypes",
		"NONE",
	)
	S3CompatEncryptionTypeEnumDef = g.NewEnum(
		"S3CompatEncryptionType", "S3CompatEncryptionTypes",
		"NONE",
	)
)

var externalS3StorageLocationDef = g.NewQueryStruct("S3StorageLocationParams").
	EnumAssignment("STORAGE_PROVIDER", S3StorageProviderEnumDef, g.ParameterOptions().SingleQuotes().Required()).
	TextAssignment("STORAGE_AWS_ROLE_ARN", g.ParameterOptions().SingleQuotes().Required()).
	TextAssignment("STORAGE_BASE_URL", g.ParameterOptions().SingleQuotes().Required()).
	OptionalTextAssignment("STORAGE_AWS_EXTERNAL_ID", g.ParameterOptions().SingleQuotes()).
	OptionalTextAssignment("STORAGE_AWS_ACCESS_POINT_ARN", g.ParameterOptions().SingleQuotes()).
	OptionalBooleanAssignment("USE_PRIVATELINK_ENDPOINT", g.ParameterOptions()).
	OptionalQueryStructField(
		"Encryption",
		g.NewQueryStruct("ExternalVolumeS3Encryption").
			EnumAssignmentWithFieldName("TYPE", S3EncryptionTypeEnumDef, g.ParameterOptions().SingleQuotes().Required(), "EncryptionType").
			OptionalTextAssignment("KMS_KEY_ID", g.ParameterOptions().SingleQuotes()),
		g.ListOptions().Parentheses().NoComma().SQL("ENCRYPTION ="),
	)

var externalGCSStorageLocationDef = g.NewQueryStruct("GCSStorageLocationParams").
	PredefinedQueryStructField("StorageProviderGcs", "string", g.StaticOptions().SQL("STORAGE_PROVIDER = 'GCS'")).
	TextAssignment("STORAGE_BASE_URL", g.ParameterOptions().SingleQuotes().Required()).
	OptionalQueryStructField(
		"Encryption",
		g.NewQueryStruct("ExternalVolumeGCSEncryption").
			EnumAssignmentWithFieldName("TYPE", GCSEncryptionTypeEnumDef, g.ParameterOptions().SingleQuotes().Required(), "EncryptionType").
			OptionalTextAssignment("KMS_KEY_ID", g.ParameterOptions().SingleQuotes()),
		g.ListOptions().Parentheses().NoComma().SQL("ENCRYPTION ="),
	)

var externalAzureStorageLocationDef = g.NewQueryStruct("AzureStorageLocationParams").
	PredefinedQueryStructField("StorageProviderAzure", "string", g.StaticOptions().SQL("STORAGE_PROVIDER = 'AZURE'")).
	TextAssignment("AZURE_TENANT_ID", g.ParameterOptions().SingleQuotes().Required()).
	TextAssignment("STORAGE_BASE_URL", g.ParameterOptions().SingleQuotes().Required()).
	OptionalBooleanAssignment("USE_PRIVATELINK_ENDPOINT", g.ParameterOptions())

var externalS3CompatStorageLocationDef = g.NewQueryStruct("S3CompatStorageLocationParams").
	PredefinedQueryStructField("StorageProviderS3Compat", "string", g.StaticOptions().SQL("STORAGE_PROVIDER = 'S3COMPAT'")).
	TextAssignment("STORAGE_BASE_URL", g.ParameterOptions().SingleQuotes().Required()).
	TextAssignment("STORAGE_ENDPOINT", g.ParameterOptions().SingleQuotes().Required()).
	QueryStructField(
		"Credentials",
		g.NewQueryStruct("ExternalVolumeS3CompatCredentials").
			TextAssignment("AWS_KEY_ID", g.ParameterOptions().SingleQuotes().Required()).
			TextAssignment("AWS_SECRET_KEY", g.ParameterOptions().SingleQuotes().Required()),
		g.ListOptions().Parentheses().NoComma().SQL("CREDENTIALS =").Required(),
	)

var storageLocationDef = g.NewQueryStruct("ExternalVolumeStorageLocation").
	TextAssignment("NAME", g.ParameterOptions().SingleQuotes().Required()).
	OptionalQueryStructField(
		"S3StorageLocationParams",
		externalS3StorageLocationDef,
		g.ListOptions().NoComma(),
	).
	OptionalQueryStructField(
		"GCSStorageLocationParams",
		externalGCSStorageLocationDef,
		g.ListOptions().NoComma(),
	).
	OptionalQueryStructField(
		"AzureStorageLocationParams",
		externalAzureStorageLocationDef,
		g.ListOptions().NoComma(),
	).
	OptionalQueryStructField(
		"S3CompatStorageLocationParams",
		externalS3CompatStorageLocationDef,
		g.ListOptions().NoComma(),
	).
	WithValidation(g.ExactlyOneValueSet, "S3StorageLocationParams", "GCSStorageLocationParams", "AzureStorageLocationParams", "S3CompatStorageLocationParams")

var storageLocationItemDef = g.NewQueryStruct("ExternalVolumeStorageLocationItem").
	QueryStructField(
		"ExternalVolumeStorageLocation",
		storageLocationDef,
		g.ListOptions().Parentheses().NoComma().Required(),
	)

var updateStorageLocationDef = g.NewQueryStruct("AlterExternalVolumeUpdateStorageLocation").
	TextAssignment("STORAGE_LOCATION", g.ParameterOptions().SingleQuotes().NoEquals().Required()).
	QueryStructField(
		"Credentials",
		g.NewQueryStruct("ExternalVolumeUpdateCredentials").
			TextAssignment("AWS_KEY_ID", g.ParameterOptions().SingleQuotes().Required()).
			TextAssignment("AWS_SECRET_KEY", g.ParameterOptions().SingleQuotes().Required()),
		g.ListOptions().Parentheses().NoComma().SQL("CREDENTIALS =").Required(),
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
			ListAssignment("STORAGE_LOCATIONS", "ExternalVolumeStorageLocationItem", g.ParameterOptions().Parentheses().Required()).
			OptionalBooleanAssignment("ALLOW_WRITES", nil).
			OptionalComment().
			WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists").
			WithValidation(g.ValidIdentifier, "name"),
		storageLocationItemDef,
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
				storageLocationItemDef,
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
	DescribeOperationWithPairedStructs(
		g.DescriptionMappingKindSlice,
		"https://docs.snowflake.com/en/sql-reference/sql/desc-external-volume",
		g.StructPair("externalVolumeDescRow", "ExternalVolumeProperty").
			Text("parent_property", g.WithPlainFieldName("Parent")).
			Text("property", g.WithPlainFieldName("Name")).
			Text("property_type", g.WithPlainFieldName("Type")).
			Text("property_value", g.WithPlainFieldName("Value")).
			Text("property_default", g.WithPlainFieldName("Default")),
		g.NewQueryStruct("DescExternalVolume").
			Describe().
			SQL("EXTERNAL VOLUME").
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	).
	ShowOperationWithPairedStructs(
		"https://docs.snowflake.com/en/sql-reference/sql/show-external-volumes",
		g.StructPair("externalVolumeShowRow", "ExternalVolume").
			Text("name").
			Bool("allow_writes").
			OptionalText("comment", g.WithRequiredInPlain()),
		g.NewQueryStruct("ShowExternalVolumes").
			Show().
			SQL("EXTERNAL VOLUMES").
			OptionalLike(),
	).
	WithEnums(
		StorageProviderEnumDef,
		S3StorageProviderEnumDef,
		S3EncryptionTypeEnumDef,
		GCSEncryptionTypeEnumDef,
		AzureEncryptionTypeEnumDef,
		S3CompatEncryptionTypeEnumDef,
	)
