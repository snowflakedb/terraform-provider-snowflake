package defs

import (
	ssotdefs "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot/defs"
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

func externalS3StorageLocation() *g.QueryStruct {
	return g.NewQueryStruct("S3StorageLocationParams").
		EnumAssignment("STORAGE_PROVIDER", ssotdefs.S3StorageProvider, g.ParameterOptions().SingleQuotes().Required()).
		TextAssignment("STORAGE_AWS_ROLE_ARN", g.ParameterOptions().SingleQuotes().Required()).
		TextAssignment("STORAGE_BASE_URL", g.ParameterOptions().SingleQuotes().Required()).
		OptionalTextAssignment("STORAGE_AWS_EXTERNAL_ID", g.ParameterOptions().SingleQuotes()).
		OptionalTextAssignment("STORAGE_AWS_ACCESS_POINT_ARN", g.ParameterOptions().SingleQuotes()).
		OptionalBooleanAssignment("USE_PRIVATELINK_ENDPOINT", g.ParameterOptions()).
		OptionalQueryStructField(
			"Encryption",
			g.NewQueryStruct("ExternalVolumeS3Encryption").
				EnumAssignmentWithFieldName("TYPE", ssotdefs.S3EncryptionType, g.ParameterOptions().SingleQuotes().Required(), "EncryptionType").
				OptionalTextAssignment("KMS_KEY_ID", g.ParameterOptions().SingleQuotes()),
			g.ListOptions().Parentheses().NoComma().SQL("ENCRYPTION ="),
		)
}

func externalGCSStorageLocation() *g.QueryStruct {
	return g.NewQueryStruct("GCSStorageLocationParams").
		PredefinedQueryStructField("StorageProviderGcs", "string", g.StaticOptions().SQL("STORAGE_PROVIDER = 'GCS'")).
		TextAssignment("STORAGE_BASE_URL", g.ParameterOptions().SingleQuotes().Required()).
		OptionalQueryStructField(
			"Encryption",
			g.NewQueryStruct("ExternalVolumeGCSEncryption").
				EnumAssignmentWithFieldName("TYPE", ssotdefs.GCSEncryptionType, g.ParameterOptions().SingleQuotes().Required(), "EncryptionType").
				OptionalTextAssignment("KMS_KEY_ID", g.ParameterOptions().SingleQuotes()),
			g.ListOptions().Parentheses().NoComma().SQL("ENCRYPTION ="),
		)
}

func externalAzureStorageLocation() *g.QueryStruct {
	return g.NewQueryStruct("AzureStorageLocationParams").
		PredefinedQueryStructField("StorageProviderAzure", "string", g.StaticOptions().SQL("STORAGE_PROVIDER = 'AZURE'")).
		TextAssignment("AZURE_TENANT_ID", g.ParameterOptions().SingleQuotes().Required()).
		TextAssignment("STORAGE_BASE_URL", g.ParameterOptions().SingleQuotes().Required()).
		OptionalBooleanAssignment("USE_PRIVATELINK_ENDPOINT", g.ParameterOptions())
}

func externalS3CompatStorageLocation() *g.QueryStruct {
	return g.NewQueryStruct("S3CompatStorageLocationParams").
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
}

func storageLocation() *g.QueryStruct {
	return g.NewQueryStruct("ExternalVolumeStorageLocation").
		TextAssignment("NAME", g.ParameterOptions().SingleQuotes().Required()).
		OptionalQueryStructField(
			"S3StorageLocationParams",
			externalS3StorageLocation(),
			g.ListOptions().NoComma(),
		).
		OptionalQueryStructField(
			"GCSStorageLocationParams",
			externalGCSStorageLocation(),
			g.ListOptions().NoComma(),
		).
		OptionalQueryStructField(
			"AzureStorageLocationParams",
			externalAzureStorageLocation(),
			g.ListOptions().NoComma(),
		).
		OptionalQueryStructField(
			"S3CompatStorageLocationParams",
			externalS3CompatStorageLocation(),
			g.ListOptions().NoComma(),
		).
		WithValidation(g.ExactlyOneValueSet, "S3StorageLocationParams", "GCSStorageLocationParams", "AzureStorageLocationParams", "S3CompatStorageLocationParams")
}

func storageLocationItem() *g.QueryStruct {
	return g.NewQueryStruct("ExternalVolumeStorageLocationItem").
		QueryStructField(
			"ExternalVolumeStorageLocation",
			storageLocation(),
			g.ListOptions().Parentheses().NoComma().Required(),
		)
}

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
	CreateOperationFromObject(ssotdefs.ExternalVolumeObject, storageLocationItem()).
	AlterOperationFromObject(ssotdefs.ExternalVolumeObject, storageLocationItem(), updateStorageLocationDef).
	DropOperationFromObject(ssotdefs.ExternalVolumeObject).
	DescribeOperationFromObject(ssotdefs.ExternalVolumeObject).
	ShowOperationFromObject(ssotdefs.ExternalVolumeObject).
	WithEnums(
		ssotdefs.StorageProvider,
		ssotdefs.S3StorageProvider,
		ssotdefs.S3EncryptionType,
		ssotdefs.GCSEncryptionType,
		ssotdefs.AzureEncryptionType,
		ssotdefs.S3CompatEncryptionType,
	)
