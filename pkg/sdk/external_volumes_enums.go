package sdk

import (
	"fmt"
	"strings"
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
