package sdk

import (
	"fmt"
	"strings"

	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"
)

type (
	StorageProvider   string
	S3StorageProvider string
	S3EncryptionType  string
	GCSEncryptionType string
)

var (
	S3EncryptionTypeSseS3   S3EncryptionType  = "AWS_SSE_S3"
	S3EncryptionTypeSseKms  S3EncryptionType  = "AWS_SSE_KMS"
	S3EncryptionNone        S3EncryptionType  = "NONE"
	GCSEncryptionTypeSseKms GCSEncryptionType = "GCS_SSE_KMS"
	GCSEncryptionTypeNone   GCSEncryptionType = "NONE"
	S3StorageProviderS3     S3StorageProvider = "S3"
	S3StorageProviderS3GOV  S3StorageProvider = "S3GOV"
	StorageProviderGCS      StorageProvider   = "GCS"
	StorageProviderAzure    StorageProvider   = "AZURE"
	StorageProviderS3       StorageProvider   = "S3"
	StorageProviderS3GOV    StorageProvider   = "S3GOV"
)

var AllStorageProviderValues = []StorageProvider{
	StorageProviderGCS,
	StorageProviderAzure,
	StorageProviderS3,
	StorageProviderS3GOV,
}

func ToS3EncryptionType(s string) (S3EncryptionType, error) {
	switch strings.ToUpper(s) {
	case string(S3EncryptionTypeSseS3):
		return S3EncryptionTypeSseS3, nil
	case string(S3EncryptionTypeSseKms):
		return S3EncryptionTypeSseKms, nil
	case string(S3EncryptionNone):
		return S3EncryptionNone, nil
	default:
		return "", fmt.Errorf("invalid s3 encryption type: %s", s)
	}
}

func ToGCSEncryptionType(s string) (GCSEncryptionType, error) {
	switch strings.ToUpper(s) {
	case string(GCSEncryptionTypeSseKms):
		return GCSEncryptionTypeSseKms, nil
	case string(GCSEncryptionTypeNone):
		return GCSEncryptionTypeNone, nil
	default:
		return "", fmt.Errorf("invalid gcs encryption type: %s", s)
	}
}

func ToStorageProvider(s string) (StorageProvider, error) {
	switch strings.ToUpper(s) {
	case string(StorageProviderGCS):
		return StorageProviderGCS, nil
	case string(StorageProviderAzure):
		return StorageProviderAzure, nil
	case string(StorageProviderS3):
		return StorageProviderS3, nil
	case string(StorageProviderS3GOV):
		return StorageProviderS3GOV, nil
	default:
		return "", fmt.Errorf("invalid storage provider: %s", s)
	}
}

func ToS3StorageProvider(s string) (S3StorageProvider, error) {
	switch strings.ToUpper(s) {
	case string(S3StorageProviderS3):
		return S3StorageProviderS3, nil
	case string(S3StorageProviderS3GOV):
		return S3StorageProviderS3GOV, nil
	default:
		return "", fmt.Errorf("invalid s3 storage provider: %s", s)
	}
}

var externalS3StorageLocationDef = g.NewQueryStruct("S3StorageLocationParams").
	TextAssignment("NAME", g.ParameterOptions().SingleQuotes().Required()).
	Assignment("STORAGE_PROVIDER", g.KindOfT[S3StorageProvider](), g.ParameterOptions().SingleQuotes().Required()).
	TextAssignment("STORAGE_AWS_ROLE_ARN", g.ParameterOptions().SingleQuotes().Required()).
	TextAssignment("STORAGE_BASE_URL", g.ParameterOptions().SingleQuotes().Required()).
	OptionalTextAssignment("STORAGE_AWS_EXTERNAL_ID", g.ParameterOptions().SingleQuotes()).
	OptionalQueryStructField(
		"Encryption",
		g.NewQueryStruct("ExternalVolumeS3Encryption").
			Assignment("TYPE", g.KindOfT[S3EncryptionType](), g.ParameterOptions().SingleQuotes().Required()).
			OptionalTextAssignment("KMS_KEY_ID", g.ParameterOptions().SingleQuotes()),
		g.ListOptions().Parentheses().NoComma().SQL("ENCRYPTION ="),
	)

var externalGCSStorageLocationDef = g.NewQueryStruct("GCSStorageLocationParams").
	TextAssignment("NAME", g.ParameterOptions().SingleQuotes().Required()).
	PredefinedQueryStructField("StorageProviderGcs", "string", g.StaticOptions().SQL(fmt.Sprintf("STORAGE_PROVIDER = '%s'", StorageProviderGCS))).
	TextAssignment("STORAGE_BASE_URL", g.ParameterOptions().SingleQuotes().Required()).
	OptionalQueryStructField(
		"Encryption",
		g.NewQueryStruct("ExternalVolumeGCSEncryption").
			Assignment("TYPE", g.KindOfT[GCSEncryptionType](), g.ParameterOptions().SingleQuotes().Required()).
			OptionalTextAssignment("KMS_KEY_ID", g.ParameterOptions().SingleQuotes()),
		g.ListOptions().Parentheses().NoComma().SQL("ENCRYPTION ="),
	)

var externalAzureStorageLocationDef = g.NewQueryStruct("AzureStorageLocationParams").
	TextAssignment("NAME", g.ParameterOptions().SingleQuotes().Required()).
	PredefinedQueryStructField("StorageProviderAzure", "string", g.StaticOptions().SQL(fmt.Sprintf("STORAGE_PROVIDER = '%s'", StorageProviderAzure))).
	TextAssignment("AZURE_TENANT_ID", g.ParameterOptions().SingleQuotes().Required()).
	TextAssignment("STORAGE_BASE_URL", g.ParameterOptions().SingleQuotes().Required())

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
	WithValidation(g.ExactlyOneValueSet, "S3StorageLocationParams", "GCSStorageLocationParams", "AzureStorageLocationParams")

var ExternalVolumesDef = g.NewInterface(
	"ExternalVolumes",
	"ExternalVolume",
	g.KindOfT[AccountObjectIdentifier](),
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
			WithValidation(g.ExactlyOneValueSet, "RemoveStorageLocation", "Set", "AddStorageLocation").
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
