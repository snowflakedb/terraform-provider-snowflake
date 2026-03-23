package sdk

import (
	"fmt"
	"slices"
	"strings"
)

type (
	StorageProvider   string
	S3StorageProvider string
	S3EncryptionType  string
	GCSEncryptionType string
	// AzureEncryptionType is currently used only in desc output. Passing encryption for azure storage location is not supported in Snowflake.
	AzureEncryptionType string
	// S3CompatEncryptionType is currently used only in desc output. Passing encryption for s3-compatible storage location is not supported in Snowflake.
	S3CompatEncryptionType string
)

var (
	S3EncryptionTypeSseS3  S3EncryptionType = "AWS_SSE_S3"
	S3EncryptionTypeSseKms S3EncryptionType = "AWS_SSE_KMS"
	S3EncryptionNone       S3EncryptionType = "NONE"

	GCSEncryptionTypeSseKms GCSEncryptionType = "GCS_SSE_KMS"
	GCSEncryptionTypeNone   GCSEncryptionType = "NONE"

	AzureEncryptionTypeNone AzureEncryptionType = "NONE"

	S3CompatEncryptionTypeNone S3CompatEncryptionType = "NONE"

	S3StorageProviderS3    S3StorageProvider = "S3"
	S3StorageProviderS3GOV S3StorageProvider = "S3GOV"

	StorageProviderGCS          StorageProvider = "GCS"
	StorageProviderAzure        StorageProvider = "AZURE"
	StorageProviderS3           StorageProvider = "S3"
	StorageProviderS3GOV        StorageProvider = "S3GOV"
	StorageProviderS3Compatible StorageProvider = "S3COMPAT"
)

var AllStorageProviderValues = []StorageProvider{
	StorageProviderGCS,
	StorageProviderAzure,
	StorageProviderS3,
	StorageProviderS3GOV,
	StorageProviderS3Compatible,
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

func ToAzureEncryptionType(s string) (AzureEncryptionType, error) {
	switch strings.ToUpper(s) {
	case string(AzureEncryptionTypeNone):
		return AzureEncryptionTypeNone, nil
	default:
		return "", fmt.Errorf("invalid azure encryption type: %s", s)
	}
}

func ToS3CompatEncryptionType(s string) (S3CompatEncryptionType, error) {
	switch strings.ToUpper(s) {
	case string(S3CompatEncryptionTypeNone):
		return S3CompatEncryptionTypeNone, nil
	default:
		return "", fmt.Errorf("invalid s3 compat encryption type: %s", s)
	}
}

func ToStorageProvider(s string) (StorageProvider, error) {
	s = strings.ToUpper(s)
	if !slices.Contains(AllStorageProviderValues, StorageProvider(s)) {
		return "", fmt.Errorf("invalid storage provider: %s", s)
	}
	return StorageProvider(s), nil
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
