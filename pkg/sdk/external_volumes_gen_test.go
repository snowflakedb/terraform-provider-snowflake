package sdk

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExternalVolumes_Create(t *testing.T) {
	id := randomAccountObjectIdentifier()
	// Minimal valid CreateExternalVolumeOptions
	defaultOpts := func() *CreateExternalVolumeOptions {
		return &CreateExternalVolumeOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateExternalVolumeOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: conflicting fields for [opts.OrReplace opts.IfNotExists]", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.IfNotExists = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateExternalVolumeOptions", "OrReplace", "IfNotExists"))
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptyAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: exactly one field from [opts.StorageLocations[i].S3StorageLocationParams opts.StorageLocations[i].GCSStorageLocationParams opts.StorageLocations[i].AzureStorageLocationParams opts.StorageLocations[i].S3CompatStorageLocationParams] should be present", func(t *testing.T) {
		opts := defaultOpts()
		opts.StorageLocations = []ExternalVolumeStorageLocationItem{
			{ExternalVolumeStorageLocation: ExternalVolumeStorageLocation{Name: "s3_basic", S3StorageLocationParams: &S3StorageLocationParams{StorageProvider: S3StorageProviderS3, StorageAwsRoleArn: "arn:aws:iam::123456789012:role/myrole", StorageBaseUrl: "s3://my-bucket/path"}}},
			{},
			{ExternalVolumeStorageLocation: ExternalVolumeStorageLocation{
				Name:                       "multi",
				S3StorageLocationParams:    &S3StorageLocationParams{StorageProvider: S3StorageProviderS3, StorageAwsRoleArn: "arn:aws:iam::123456789012:role/myrole", StorageBaseUrl: "s3://my-bucket/path"},
				GCSStorageLocationParams:   &GCSStorageLocationParams{StorageBaseUrl: "gcs://my-bucket/path"},
				AzureStorageLocationParams: &AzureStorageLocationParams{AzureTenantId: "a123b4cd-1abc-12ab-12ab-1a2b34c5d678", StorageBaseUrl: "azure://myaccount.blob.core.windows.net/mycontainer/path"},
			}},
		}
		assertOptsInvalidJoinedErrors(
			t,
			opts,
			errExactlyOneOf("CreateExternalVolumeOptions.StorageLocation[1]", "S3StorageLocationParams", "GCSStorageLocationParams", "AzureStorageLocationParams", "S3CompatStorageLocationParams"),
			errExactlyOneOf("CreateExternalVolumeOptions.StorageLocation[2]", "S3StorageLocationParams", "GCSStorageLocationParams", "AzureStorageLocationParams", "S3CompatStorageLocationParams"),
		)
	})

	t.Run("validation: length of opts.StorageLocations is > 0", func(t *testing.T) {
		opts := defaultOpts()
		opts.StorageLocations = []ExternalVolumeStorageLocationItem{}
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("CreateExternalVolumeOptions", "StorageLocations"))
	})

	t.Run("basic - s3", func(t *testing.T) {
		opts := defaultOpts()
		opts.StorageLocations = []ExternalVolumeStorageLocationItem{{ExternalVolumeStorageLocation: ExternalVolumeStorageLocation{Name: "s3_basic", S3StorageLocationParams: &S3StorageLocationParams{StorageProvider: S3StorageProviderS3, StorageAwsRoleArn: "arn:aws:iam::123456789012:role/myrole", StorageBaseUrl: "s3://my-bucket/path"}}}}
		assertOptsValidAndSQLEquals(t, opts, `CREATE EXTERNAL VOLUME %s STORAGE_LOCATIONS = ((NAME = 's3_basic' STORAGE_PROVIDER = 'S3' STORAGE_AWS_ROLE_ARN = 'arn:aws:iam::123456789012:role/myrole' STORAGE_BASE_URL = 's3://my-bucket/path'))`, id.FullyQualifiedName())
	})

	t.Run("basic - gcs", func(t *testing.T) {
		opts := defaultOpts()
		opts.StorageLocations = []ExternalVolumeStorageLocationItem{{ExternalVolumeStorageLocation: ExternalVolumeStorageLocation{Name: "gcs_basic", GCSStorageLocationParams: &GCSStorageLocationParams{StorageBaseUrl: "gcs://my-bucket/path"}}}}
		assertOptsValidAndSQLEquals(t, opts, `CREATE EXTERNAL VOLUME %s STORAGE_LOCATIONS = ((NAME = 'gcs_basic' STORAGE_PROVIDER = 'GCS' STORAGE_BASE_URL = 'gcs://my-bucket/path'))`, id.FullyQualifiedName())
	})

	t.Run("basic - azure", func(t *testing.T) {
		opts := defaultOpts()
		opts.StorageLocations = []ExternalVolumeStorageLocationItem{{ExternalVolumeStorageLocation: ExternalVolumeStorageLocation{Name: "azure_basic", AzureStorageLocationParams: &AzureStorageLocationParams{AzureTenantId: "a123b4cd-1abc-12ab-12ab-1a2b34c5d678", StorageBaseUrl: "azure://myaccount.blob.core.windows.net/mycontainer/path"}}}}
		assertOptsValidAndSQLEquals(t, opts, `CREATE EXTERNAL VOLUME %s STORAGE_LOCATIONS = ((NAME = 'azure_basic' STORAGE_PROVIDER = 'AZURE' AZURE_TENANT_ID = 'a123b4cd-1abc-12ab-12ab-1a2b34c5d678' STORAGE_BASE_URL = 'azure://myaccount.blob.core.windows.net/mycontainer/path'))`, id.FullyQualifiedName())
	})

	t.Run("basic - s3compat", func(t *testing.T) {
		opts := defaultOpts()
		opts.StorageLocations = []ExternalVolumeStorageLocationItem{{ExternalVolumeStorageLocation: ExternalVolumeStorageLocation{Name: "s3compat_basic", S3CompatStorageLocationParams: &S3CompatStorageLocationParams{StorageBaseUrl: "s3compat://my-bucket/path", StorageEndpoint: "https://s3-compatible.example.com"}}}}
		assertOptsValidAndSQLEquals(t, opts, `CREATE EXTERNAL VOLUME %s STORAGE_LOCATIONS = ((NAME = 's3compat_basic' STORAGE_PROVIDER = 'S3COMPAT' STORAGE_BASE_URL = 's3compat://my-bucket/path' STORAGE_ENDPOINT = 'https://s3-compatible.example.com'))`, id.FullyQualifiedName())
	})

	t.Run("all options - s3", func(t *testing.T) {
		opts := defaultOpts()
		opts.StorageLocations = []ExternalVolumeStorageLocationItem{{ExternalVolumeStorageLocation: ExternalVolumeStorageLocation{Name: "s3_complete", S3StorageLocationParams: &S3StorageLocationParams{StorageProvider: S3StorageProviderS3, StorageAwsRoleArn: "arn:aws:iam::123456789012:role/myrole", StorageBaseUrl: "s3://my-bucket/path", StorageAwsExternalId: String("external_id_123"), StorageAwsAccessPointArn: String("arn:aws:s3:us-west-2:123456789012:accesspoint/my-access-point"), UsePrivatelinkEndpoint: Bool(true), Encryption: &ExternalVolumeS3Encryption{EncryptionType: S3EncryptionTypeSseKms, KmsKeyId: String("1234abcd-12ab-34cd-56ef-1234567890ab")}}}}}
		assertOptsValidAndSQLEquals(t, opts, `CREATE EXTERNAL VOLUME %s STORAGE_LOCATIONS = ((NAME = 's3_complete' STORAGE_PROVIDER = 'S3' STORAGE_AWS_ROLE_ARN = 'arn:aws:iam::123456789012:role/myrole' STORAGE_BASE_URL = 's3://my-bucket/path' STORAGE_AWS_EXTERNAL_ID = 'external_id_123' STORAGE_AWS_ACCESS_POINT_ARN = 'arn:aws:s3:us-west-2:123456789012:accesspoint/my-access-point' USE_PRIVATELINK_ENDPOINT = true ENCRYPTION = (TYPE = 'AWS_SSE_KMS' KMS_KEY_ID = '1234abcd-12ab-34cd-56ef-1234567890ab')))`, id.FullyQualifiedName())
	})

	t.Run("all options - gcs", func(t *testing.T) {
		opts := defaultOpts()
		opts.StorageLocations = []ExternalVolumeStorageLocationItem{{ExternalVolumeStorageLocation: ExternalVolumeStorageLocation{Name: "gcs_complete", GCSStorageLocationParams: &GCSStorageLocationParams{StorageBaseUrl: "gcs://my-bucket/path", Encryption: &ExternalVolumeGCSEncryption{EncryptionType: GCSEncryptionTypeSseKms, KmsKeyId: String("1234abcd-12ab-34cd-56ef-1234567890ab")}}}}}
		assertOptsValidAndSQLEquals(t, opts, `CREATE EXTERNAL VOLUME %s STORAGE_LOCATIONS = ((NAME = 'gcs_complete' STORAGE_PROVIDER = 'GCS' STORAGE_BASE_URL = 'gcs://my-bucket/path' ENCRYPTION = (TYPE = 'GCS_SSE_KMS' KMS_KEY_ID = '1234abcd-12ab-34cd-56ef-1234567890ab')))`, id.FullyQualifiedName())
	})

	t.Run("all options - azure", func(t *testing.T) {
		opts := defaultOpts()
		opts.StorageLocations = []ExternalVolumeStorageLocationItem{{ExternalVolumeStorageLocation: ExternalVolumeStorageLocation{Name: "azure_complete", AzureStorageLocationParams: &AzureStorageLocationParams{AzureTenantId: "a123b4cd-1abc-12ab-12ab-1a2b34c5d678", StorageBaseUrl: "azure://myaccount.blob.core.windows.net/mycontainer/path", UsePrivatelinkEndpoint: Bool(true)}}}}
		assertOptsValidAndSQLEquals(t, opts, `CREATE EXTERNAL VOLUME %s STORAGE_LOCATIONS = ((NAME = 'azure_complete' STORAGE_PROVIDER = 'AZURE' AZURE_TENANT_ID = 'a123b4cd-1abc-12ab-12ab-1a2b34c5d678' STORAGE_BASE_URL = 'azure://myaccount.blob.core.windows.net/mycontainer/path' USE_PRIVATELINK_ENDPOINT = true))`, id.FullyQualifiedName())
	})

	t.Run("all options - s3compat", func(t *testing.T) {
		opts := defaultOpts()
		opts.StorageLocations = []ExternalVolumeStorageLocationItem{{ExternalVolumeStorageLocation: ExternalVolumeStorageLocation{Name: "s3compat_complete", S3CompatStorageLocationParams: &S3CompatStorageLocationParams{StorageBaseUrl: "s3compat://my-bucket/path", StorageEndpoint: "https://s3-compatible.example.com", Credentials: &ExternalVolumeS3CompatCredentials{AwsKeyId: "AKIAIOSFODNN7EXAMPLE", AwsSecretKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"}}}}}
		assertOptsValidAndSQLEquals(t, opts, `CREATE EXTERNAL VOLUME %s STORAGE_LOCATIONS = ((NAME = 's3compat_complete' STORAGE_PROVIDER = 'S3COMPAT' STORAGE_BASE_URL = 's3compat://my-bucket/path' STORAGE_ENDPOINT = 'https://s3-compatible.example.com' CREDENTIALS = (AWS_KEY_ID = 'AKIAIOSFODNN7EXAMPLE' AWS_SECRET_KEY = 'wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY')))`, id.FullyQualifiedName())
	})

	t.Run("all storage location types with volume options", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.StorageLocations = []ExternalVolumeStorageLocationItem{
			{ExternalVolumeStorageLocation: ExternalVolumeStorageLocation{Name: "s3_complete", S3StorageLocationParams: &S3StorageLocationParams{StorageProvider: S3StorageProviderS3, StorageAwsRoleArn: "arn:aws:iam::123456789012:role/myrole", StorageBaseUrl: "s3://my-bucket/path", StorageAwsExternalId: String("external_id_123"), StorageAwsAccessPointArn: String("arn:aws:s3:us-west-2:123456789012:accesspoint/my-access-point"), UsePrivatelinkEndpoint: Bool(true), Encryption: &ExternalVolumeS3Encryption{EncryptionType: S3EncryptionTypeSseKms, KmsKeyId: String("1234abcd-12ab-34cd-56ef-1234567890ab")}}}},
			{ExternalVolumeStorageLocation: ExternalVolumeStorageLocation{Name: "gcs_complete", GCSStorageLocationParams: &GCSStorageLocationParams{StorageBaseUrl: "gcs://my-bucket/path", Encryption: &ExternalVolumeGCSEncryption{EncryptionType: GCSEncryptionTypeSseKms, KmsKeyId: String("1234abcd-12ab-34cd-56ef-1234567890ab")}}}},
			{ExternalVolumeStorageLocation: ExternalVolumeStorageLocation{Name: "azure_complete", AzureStorageLocationParams: &AzureStorageLocationParams{AzureTenantId: "a123b4cd-1abc-12ab-12ab-1a2b34c5d678", StorageBaseUrl: "azure://myaccount.blob.core.windows.net/mycontainer/path", UsePrivatelinkEndpoint: Bool(true)}}},
			{ExternalVolumeStorageLocation: ExternalVolumeStorageLocation{Name: "s3compat_complete", S3CompatStorageLocationParams: &S3CompatStorageLocationParams{StorageBaseUrl: "s3compat://my-bucket/path", StorageEndpoint: "https://s3-compatible.example.com", Credentials: &ExternalVolumeS3CompatCredentials{AwsKeyId: "AKIAIOSFODNN7EXAMPLE", AwsSecretKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"}}}},
		}
		opts.AllowWrites = Bool(true)
		opts.Comment = String("some comment")
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE EXTERNAL VOLUME %s STORAGE_LOCATIONS = ((NAME = 's3_complete' STORAGE_PROVIDER = 'S3' STORAGE_AWS_ROLE_ARN = 'arn:aws:iam::123456789012:role/myrole' STORAGE_BASE_URL = 's3://my-bucket/path' STORAGE_AWS_EXTERNAL_ID = 'external_id_123' STORAGE_AWS_ACCESS_POINT_ARN = 'arn:aws:s3:us-west-2:123456789012:accesspoint/my-access-point' USE_PRIVATELINK_ENDPOINT = true ENCRYPTION = (TYPE = 'AWS_SSE_KMS' KMS_KEY_ID = '1234abcd-12ab-34cd-56ef-1234567890ab')), (NAME = 'gcs_complete' STORAGE_PROVIDER = 'GCS' STORAGE_BASE_URL = 'gcs://my-bucket/path' ENCRYPTION = (TYPE = 'GCS_SSE_KMS' KMS_KEY_ID = '1234abcd-12ab-34cd-56ef-1234567890ab')), (NAME = 'azure_complete' STORAGE_PROVIDER = 'AZURE' AZURE_TENANT_ID = 'a123b4cd-1abc-12ab-12ab-1a2b34c5d678' STORAGE_BASE_URL = 'azure://myaccount.blob.core.windows.net/mycontainer/path' USE_PRIVATELINK_ENDPOINT = true), (NAME = 's3compat_complete' STORAGE_PROVIDER = 'S3COMPAT' STORAGE_BASE_URL = 's3compat://my-bucket/path' STORAGE_ENDPOINT = 'https://s3-compatible.example.com' CREDENTIALS = (AWS_KEY_ID = 'AKIAIOSFODNN7EXAMPLE' AWS_SECRET_KEY = 'wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY'))) ALLOW_WRITES = true COMMENT = 'some comment'`, id.FullyQualifiedName())
	})
}

func TestExternalVolumes_Alter(t *testing.T) {
	id := randomAccountObjectIdentifier()
	// Minimal valid AlterExternalVolumeOptions
	defaultOpts := func() *AlterExternalVolumeOptions {
		return &AlterExternalVolumeOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *AlterExternalVolumeOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: exactly one field from [opts.RemoveStorageLocation opts.Set opts.AddStorageLocation opts.UpdateStorageLocation] should be present - zero set", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterExternalVolumeOptions", "RemoveStorageLocation", "Set", "AddStorageLocation", "UpdateStorageLocation"))
	})

	t.Run("validation: exactly one field from [opts.RemoveStorageLocation opts.Set opts.AddStorageLocation opts.UpdateStorageLocation] should be present - two set", func(t *testing.T) {
		removeAndSetOpts := defaultOpts()
		removeAndAddOpts := defaultOpts()
		setAndAddOpts := defaultOpts()

		removeAndSetOpts.RemoveStorageLocation = String("some storage location")
		removeAndSetOpts.Set = &AlterExternalVolumeSet{AllowWrites: Bool(true)}

		removeAndAddOpts.RemoveStorageLocation = String("some storage location")
		removeAndAddOpts.AddStorageLocation = &ExternalVolumeStorageLocationItem{ExternalVolumeStorageLocation: ExternalVolumeStorageLocation{Name: "s3_complete", S3StorageLocationParams: &S3StorageLocationParams{StorageProvider: S3StorageProviderS3, StorageAwsRoleArn: "arn:aws:iam::123456789012:role/myrole", StorageBaseUrl: "s3://my-bucket/path", StorageAwsExternalId: String("external_id_123"), StorageAwsAccessPointArn: String("arn:aws:s3:us-west-2:123456789012:accesspoint/my-access-point"), UsePrivatelinkEndpoint: Bool(true), Encryption: &ExternalVolumeS3Encryption{EncryptionType: S3EncryptionTypeSseKms, KmsKeyId: String("1234abcd-12ab-34cd-56ef-1234567890ab")}}}}

		setAndAddOpts.Set = &AlterExternalVolumeSet{AllowWrites: Bool(true)}
		setAndAddOpts.AddStorageLocation = &ExternalVolumeStorageLocationItem{ExternalVolumeStorageLocation: ExternalVolumeStorageLocation{Name: "s3_complete", S3StorageLocationParams: &S3StorageLocationParams{StorageProvider: S3StorageProviderS3, StorageAwsRoleArn: "arn:aws:iam::123456789012:role/myrole", StorageBaseUrl: "s3://my-bucket/path", StorageAwsExternalId: String("external_id_123"), StorageAwsAccessPointArn: String("arn:aws:s3:us-west-2:123456789012:accesspoint/my-access-point"), UsePrivatelinkEndpoint: Bool(true), Encryption: &ExternalVolumeS3Encryption{EncryptionType: S3EncryptionTypeSseKms, KmsKeyId: String("1234abcd-12ab-34cd-56ef-1234567890ab")}}}}

		assertOptsInvalidJoinedErrors(t, removeAndSetOpts, errExactlyOneOf("AlterExternalVolumeOptions", "RemoveStorageLocation", "Set", "AddStorageLocation", "UpdateStorageLocation"))
		assertOptsInvalidJoinedErrors(t, removeAndAddOpts, errExactlyOneOf("AlterExternalVolumeOptions", "RemoveStorageLocation", "Set", "AddStorageLocation", "UpdateStorageLocation"))
		assertOptsInvalidJoinedErrors(t, setAndAddOpts, errExactlyOneOf("AlterExternalVolumeOptions", "RemoveStorageLocation", "Set", "AddStorageLocation", "UpdateStorageLocation"))
	})

	t.Run("validation: exactly one field from [opts.RemoveStorageLocation opts.Set opts.AddStorageLocation opts.UpdateStorageLocation] should be present - three set", func(t *testing.T) {
		opts := defaultOpts()
		opts.RemoveStorageLocation = String("some storage location")
		opts.Set = &AlterExternalVolumeSet{AllowWrites: Bool(true)}
		opts.AddStorageLocation = &ExternalVolumeStorageLocationItem{ExternalVolumeStorageLocation: ExternalVolumeStorageLocation{Name: "s3_complete", S3StorageLocationParams: &S3StorageLocationParams{StorageProvider: S3StorageProviderS3, StorageAwsRoleArn: "arn:aws:iam::123456789012:role/myrole", StorageBaseUrl: "s3://my-bucket/path", StorageAwsExternalId: String("external_id_123"), StorageAwsAccessPointArn: String("arn:aws:s3:us-west-2:123456789012:accesspoint/my-access-point"), UsePrivatelinkEndpoint: Bool(true), Encryption: &ExternalVolumeS3Encryption{EncryptionType: S3EncryptionTypeSseKms, KmsKeyId: String("1234abcd-12ab-34cd-56ef-1234567890ab")}}}}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterExternalVolumeOptions", "RemoveStorageLocation", "Set", "AddStorageLocation", "UpdateStorageLocation"))
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptyAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: exactly one field from [opts.AddStorageLocation.ExternalVolumeStorageLocation.S3StorageLocationParams opts.AddStorageLocation.ExternalVolumeStorageLocation.GCSStorageLocationParams opts.AddStorageLocation.ExternalVolumeStorageLocation.AzureStorageLocationParams opts.AddStorageLocation.ExternalVolumeStorageLocation.S3CompatStorageLocationParams] should be present - none set", func(t *testing.T) {
		opts := defaultOpts()
		opts.AddStorageLocation = &ExternalVolumeStorageLocationItem{ExternalVolumeStorageLocation: ExternalVolumeStorageLocation{}}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterExternalVolumeOptions.AddStorageLocation.ExternalVolumeStorageLocation", "S3StorageLocationParams", "GCSStorageLocationParams", "AzureStorageLocationParams", "S3CompatStorageLocationParams"))
	})

	t.Run("validation: exactly one field from [opts.AddStorageLocation.ExternalVolumeStorageLocation.S3StorageLocationParams opts.AddStorageLocation.ExternalVolumeStorageLocation.GCSStorageLocationParams opts.AddStorageLocation.ExternalVolumeStorageLocation.AzureStorageLocationParams opts.AddStorageLocation.ExternalVolumeStorageLocation.S3CompatStorageLocationParams] should be present - two set", func(t *testing.T) {
		s3AndGcsOpts := defaultOpts()
		s3AndAzureOpts := defaultOpts()
		gcsAndAzureOpts := defaultOpts()
		s3AndGcsOpts.AddStorageLocation = &ExternalVolumeStorageLocationItem{ExternalVolumeStorageLocation: ExternalVolumeStorageLocation{
			Name:                     "s3_and_gcs",
			S3StorageLocationParams:  &S3StorageLocationParams{StorageProvider: S3StorageProviderS3, StorageAwsRoleArn: "arn:aws:iam::123456789012:role/myrole", StorageBaseUrl: "s3://my-bucket/path"},
			GCSStorageLocationParams: &GCSStorageLocationParams{StorageBaseUrl: "gcs://my-bucket/path"},
		}}
		s3AndAzureOpts.AddStorageLocation = &ExternalVolumeStorageLocationItem{ExternalVolumeStorageLocation: ExternalVolumeStorageLocation{
			Name:                       "s3_and_azure",
			S3StorageLocationParams:    &S3StorageLocationParams{StorageProvider: S3StorageProviderS3, StorageAwsRoleArn: "arn:aws:iam::123456789012:role/myrole", StorageBaseUrl: "s3://my-bucket/path"},
			AzureStorageLocationParams: &AzureStorageLocationParams{AzureTenantId: "a123b4cd-1abc-12ab-12ab-1a2b34c5d678", StorageBaseUrl: "azure://myaccount.blob.core.windows.net/mycontainer/path"},
		}}
		gcsAndAzureOpts.AddStorageLocation = &ExternalVolumeStorageLocationItem{ExternalVolumeStorageLocation: ExternalVolumeStorageLocation{
			Name:                       "gcs_and_azure",
			GCSStorageLocationParams:   &GCSStorageLocationParams{StorageBaseUrl: "gcs://my-bucket/path"},
			AzureStorageLocationParams: &AzureStorageLocationParams{AzureTenantId: "a123b4cd-1abc-12ab-12ab-1a2b34c5d678", StorageBaseUrl: "azure://myaccount.blob.core.windows.net/mycontainer/path"},
		}}
		assertOptsInvalidJoinedErrors(t, s3AndGcsOpts, errExactlyOneOf("AlterExternalVolumeOptions.AddStorageLocation.ExternalVolumeStorageLocation", "S3StorageLocationParams", "GCSStorageLocationParams", "AzureStorageLocationParams", "S3CompatStorageLocationParams"))
		assertOptsInvalidJoinedErrors(t, s3AndAzureOpts, errExactlyOneOf("AlterExternalVolumeOptions.AddStorageLocation.ExternalVolumeStorageLocation", "S3StorageLocationParams", "GCSStorageLocationParams", "AzureStorageLocationParams", "S3CompatStorageLocationParams"))
		assertOptsInvalidJoinedErrors(t, gcsAndAzureOpts, errExactlyOneOf("AlterExternalVolumeOptions.AddStorageLocation.ExternalVolumeStorageLocation", "S3StorageLocationParams", "GCSStorageLocationParams", "AzureStorageLocationParams", "S3CompatStorageLocationParams"))
	})

	t.Run("validation: exactly one field from [opts.AddStorageLocation.ExternalVolumeStorageLocation.S3StorageLocationParams opts.AddStorageLocation.ExternalVolumeStorageLocation.GCSStorageLocationParams opts.AddStorageLocation.ExternalVolumeStorageLocation.AzureStorageLocationParams opts.AddStorageLocation.ExternalVolumeStorageLocation.S3CompatStorageLocationParams] should be present - three set", func(t *testing.T) {
		opts := defaultOpts()
		opts.AddStorageLocation = &ExternalVolumeStorageLocationItem{ExternalVolumeStorageLocation: ExternalVolumeStorageLocation{
			Name:                       "multi",
			S3StorageLocationParams:    &S3StorageLocationParams{StorageProvider: S3StorageProviderS3, StorageAwsRoleArn: "arn:aws:iam::123456789012:role/myrole", StorageBaseUrl: "s3://my-bucket/path"},
			GCSStorageLocationParams:   &GCSStorageLocationParams{StorageBaseUrl: "gcs://my-bucket/path"},
			AzureStorageLocationParams: &AzureStorageLocationParams{AzureTenantId: "a123b4cd-1abc-12ab-12ab-1a2b34c5d678", StorageBaseUrl: "azure://myaccount.blob.core.windows.net/mycontainer/path"},
		}}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterExternalVolumeOptions.AddStorageLocation.ExternalVolumeStorageLocation", "S3StorageLocationParams", "GCSStorageLocationParams", "AzureStorageLocationParams", "S3CompatStorageLocationParams"))
	})

	t.Run("remove storage location", func(t *testing.T) {
		opts := defaultOpts()
		opts.RemoveStorageLocation = String("some storage location")
		assertOptsValidAndSQLEquals(t, opts, `ALTER EXTERNAL VOLUME %s REMOVE STORAGE_LOCATION 'some storage location'`, id.FullyQualifiedName())
	})

	t.Run("set - all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &AlterExternalVolumeSet{AllowWrites: Bool(true), Comment: String("some comment")}
		assertOptsValidAndSQLEquals(t, opts, `ALTER EXTERNAL VOLUME %s SET ALLOW_WRITES = true COMMENT = 'some comment'`, id.FullyQualifiedName())
	})

	t.Run("add storage location - basic - s3", func(t *testing.T) {
		opts := defaultOpts()
		opts.AddStorageLocation = &ExternalVolumeStorageLocationItem{ExternalVolumeStorageLocation: ExternalVolumeStorageLocation{Name: "s3_basic", S3StorageLocationParams: &S3StorageLocationParams{StorageProvider: S3StorageProviderS3, StorageAwsRoleArn: "arn:aws:iam::123456789012:role/myrole", StorageBaseUrl: "s3://my-bucket/path"}}}
		assertOptsValidAndSQLEquals(t, opts, `ALTER EXTERNAL VOLUME %s ADD STORAGE_LOCATION = (NAME = 's3_basic' STORAGE_PROVIDER = 'S3' STORAGE_AWS_ROLE_ARN = 'arn:aws:iam::123456789012:role/myrole' STORAGE_BASE_URL = 's3://my-bucket/path')`, id.FullyQualifiedName())
	})

	t.Run("add storage location - basic - gcs", func(t *testing.T) {
		opts := defaultOpts()
		opts.AddStorageLocation = &ExternalVolumeStorageLocationItem{ExternalVolumeStorageLocation: ExternalVolumeStorageLocation{Name: "gcs_basic", GCSStorageLocationParams: &GCSStorageLocationParams{StorageBaseUrl: "gcs://my-bucket/path"}}}
		assertOptsValidAndSQLEquals(t, opts, `ALTER EXTERNAL VOLUME %s ADD STORAGE_LOCATION = (NAME = 'gcs_basic' STORAGE_PROVIDER = 'GCS' STORAGE_BASE_URL = 'gcs://my-bucket/path')`, id.FullyQualifiedName())
	})

	t.Run("add storage location - basic - azure", func(t *testing.T) {
		opts := defaultOpts()
		opts.AddStorageLocation = &ExternalVolumeStorageLocationItem{ExternalVolumeStorageLocation: ExternalVolumeStorageLocation{Name: "azure_basic", AzureStorageLocationParams: &AzureStorageLocationParams{AzureTenantId: "a123b4cd-1abc-12ab-12ab-1a2b34c5d678", StorageBaseUrl: "azure://myaccount.blob.core.windows.net/mycontainer/path"}}}
		assertOptsValidAndSQLEquals(t, opts, `ALTER EXTERNAL VOLUME %s ADD STORAGE_LOCATION = (NAME = 'azure_basic' STORAGE_PROVIDER = 'AZURE' AZURE_TENANT_ID = 'a123b4cd-1abc-12ab-12ab-1a2b34c5d678' STORAGE_BASE_URL = 'azure://myaccount.blob.core.windows.net/mycontainer/path')`, id.FullyQualifiedName())
	})

	t.Run("add storage location - basic - s3compat", func(t *testing.T) {
		opts := defaultOpts()
		opts.AddStorageLocation = &ExternalVolumeStorageLocationItem{ExternalVolumeStorageLocation: ExternalVolumeStorageLocation{Name: "s3compat_basic", S3CompatStorageLocationParams: &S3CompatStorageLocationParams{StorageBaseUrl: "s3compat://my-bucket/path", StorageEndpoint: "https://s3-compatible.example.com"}}}
		assertOptsValidAndSQLEquals(t, opts, `ALTER EXTERNAL VOLUME %s ADD STORAGE_LOCATION = (NAME = 's3compat_basic' STORAGE_PROVIDER = 'S3COMPAT' STORAGE_BASE_URL = 's3compat://my-bucket/path' STORAGE_ENDPOINT = 'https://s3-compatible.example.com')`, id.FullyQualifiedName())
	})

	t.Run("add storage location - all options - s3", func(t *testing.T) {
		opts := defaultOpts()
		opts.AddStorageLocation = &ExternalVolumeStorageLocationItem{ExternalVolumeStorageLocation: ExternalVolumeStorageLocation{Name: "s3_complete", S3StorageLocationParams: &S3StorageLocationParams{StorageProvider: S3StorageProviderS3, StorageAwsRoleArn: "arn:aws:iam::123456789012:role/myrole", StorageBaseUrl: "s3://my-bucket/path", StorageAwsExternalId: String("external_id_123"), StorageAwsAccessPointArn: String("arn:aws:s3:us-west-2:123456789012:accesspoint/my-access-point"), UsePrivatelinkEndpoint: Bool(true), Encryption: &ExternalVolumeS3Encryption{EncryptionType: S3EncryptionTypeSseKms, KmsKeyId: String("1234abcd-12ab-34cd-56ef-1234567890ab")}}}}
		assertOptsValidAndSQLEquals(t, opts, `ALTER EXTERNAL VOLUME %s ADD STORAGE_LOCATION = (NAME = 's3_complete' STORAGE_PROVIDER = 'S3' STORAGE_AWS_ROLE_ARN = 'arn:aws:iam::123456789012:role/myrole' STORAGE_BASE_URL = 's3://my-bucket/path' STORAGE_AWS_EXTERNAL_ID = 'external_id_123' STORAGE_AWS_ACCESS_POINT_ARN = 'arn:aws:s3:us-west-2:123456789012:accesspoint/my-access-point' USE_PRIVATELINK_ENDPOINT = true ENCRYPTION = (TYPE = 'AWS_SSE_KMS' KMS_KEY_ID = '1234abcd-12ab-34cd-56ef-1234567890ab'))`, id.FullyQualifiedName())
	})

	t.Run("add storage location - all options - gcs", func(t *testing.T) {
		opts := defaultOpts()
		opts.AddStorageLocation = &ExternalVolumeStorageLocationItem{ExternalVolumeStorageLocation: ExternalVolumeStorageLocation{Name: "gcs_complete", GCSStorageLocationParams: &GCSStorageLocationParams{StorageBaseUrl: "gcs://my-bucket/path", Encryption: &ExternalVolumeGCSEncryption{EncryptionType: GCSEncryptionTypeSseKms, KmsKeyId: String("1234abcd-12ab-34cd-56ef-1234567890ab")}}}}
		assertOptsValidAndSQLEquals(t, opts, `ALTER EXTERNAL VOLUME %s ADD STORAGE_LOCATION = (NAME = 'gcs_complete' STORAGE_PROVIDER = 'GCS' STORAGE_BASE_URL = 'gcs://my-bucket/path' ENCRYPTION = (TYPE = 'GCS_SSE_KMS' KMS_KEY_ID = '1234abcd-12ab-34cd-56ef-1234567890ab'))`, id.FullyQualifiedName())
	})

	t.Run("add storage location - all options - azure", func(t *testing.T) {
		opts := defaultOpts()
		opts.AddStorageLocation = &ExternalVolumeStorageLocationItem{ExternalVolumeStorageLocation: ExternalVolumeStorageLocation{Name: "azure_complete", AzureStorageLocationParams: &AzureStorageLocationParams{AzureTenantId: "a123b4cd-1abc-12ab-12ab-1a2b34c5d678", StorageBaseUrl: "azure://myaccount.blob.core.windows.net/mycontainer/path", UsePrivatelinkEndpoint: Bool(true)}}}
		assertOptsValidAndSQLEquals(t, opts, `ALTER EXTERNAL VOLUME %s ADD STORAGE_LOCATION = (NAME = 'azure_complete' STORAGE_PROVIDER = 'AZURE' AZURE_TENANT_ID = 'a123b4cd-1abc-12ab-12ab-1a2b34c5d678' STORAGE_BASE_URL = 'azure://myaccount.blob.core.windows.net/mycontainer/path' USE_PRIVATELINK_ENDPOINT = true)`, id.FullyQualifiedName())
	})

	t.Run("add storage location - all options - s3compat", func(t *testing.T) {
		opts := defaultOpts()
		opts.AddStorageLocation = &ExternalVolumeStorageLocationItem{ExternalVolumeStorageLocation: ExternalVolumeStorageLocation{Name: "s3compat_complete", S3CompatStorageLocationParams: &S3CompatStorageLocationParams{StorageBaseUrl: "s3compat://my-bucket/path", StorageEndpoint: "https://s3-compatible.example.com", Credentials: &ExternalVolumeS3CompatCredentials{AwsKeyId: "AKIAIOSFODNN7EXAMPLE", AwsSecretKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"}}}}
		assertOptsValidAndSQLEquals(t, opts, `ALTER EXTERNAL VOLUME %s ADD STORAGE_LOCATION = (NAME = 's3compat_complete' STORAGE_PROVIDER = 'S3COMPAT' STORAGE_BASE_URL = 's3compat://my-bucket/path' STORAGE_ENDPOINT = 'https://s3-compatible.example.com' CREDENTIALS = (AWS_KEY_ID = 'AKIAIOSFODNN7EXAMPLE' AWS_SECRET_KEY = 'wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY'))`, id.FullyQualifiedName())
	})

	t.Run("update storage location", func(t *testing.T) {
		opts := defaultOpts()
		opts.UpdateStorageLocation = &AlterExternalVolumeUpdateStorageLocation{
			StorageLocation: "some_location",
			Credentials: ExternalVolumeUpdateCredentials{
				AwsKeyId:     "AKIAIOSFODNN7EXAMPLE",
				AwsSecretKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER EXTERNAL VOLUME %s UPDATE STORAGE_LOCATION 'some_location' CREDENTIALS = (AWS_KEY_ID = 'AKIAIOSFODNN7EXAMPLE' AWS_SECRET_KEY = 'wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY')`, id.FullyQualifiedName())
	})
}

func TestExternalVolumes_Drop(t *testing.T) {
	id := randomAccountObjectIdentifier()
	// Minimal valid DropExternalVolumeOptions
	defaultOpts := func() *DropExternalVolumeOptions {
		return &DropExternalVolumeOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DropExternalVolumeOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptyAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `DROP EXTERNAL VOLUME %s`, id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, `DROP EXTERNAL VOLUME IF EXISTS %s`, id.FullyQualifiedName())
	})
}

func TestExternalVolumes_Describe(t *testing.T) {
	id := randomAccountObjectIdentifier()
	// Minimal valid DescribeExternalVolumeOptions
	defaultOpts := func() *DescribeExternalVolumeOptions {
		return &DescribeExternalVolumeOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DescribeExternalVolumeOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})
	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptyAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `DESCRIBE EXTERNAL VOLUME %s`, id.FullyQualifiedName())
	})
}

func TestExternalVolumes_Show(t *testing.T) {
	id := randomAccountObjectIdentifier()
	// Minimal valid ShowExternalVolumeOptions
	defaultOpts := func() *ShowExternalVolumeOptions {
		return &ShowExternalVolumeOptions{}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *ShowExternalVolumeOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "SHOW EXTERNAL VOLUMES")
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{
			Pattern: String(id.Name()),
		}
		assertOptsValidAndSQLEquals(t, opts, "SHOW EXTERNAL VOLUMES LIKE '%s'", id.Name())
	})
}

func Test_ExternalVolumes_ToS3EncryptionType(t *testing.T) {
	type test struct {
		input string
		want  S3EncryptionType
	}

	valid := []test{
		{input: "aws_sse_s3", want: S3EncryptionTypeSseS3},
		{input: "AWS_SSE_S3", want: S3EncryptionTypeSseS3},
		{input: "AWS_SSE_KMS", want: S3EncryptionTypeSseKms},
		{input: "NONE", want: S3EncryptionNone},
	}

	invalid := []test{
		{input: ""},
		{input: "foo"},
	}

	for _, tc := range valid {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ToS3EncryptionType(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}

	for _, tc := range invalid {
		t.Run(tc.input, func(t *testing.T) {
			_, err := ToS3EncryptionType(tc.input)
			require.Error(t, err)
		})
	}
}

func Test_ExternalVolumes_ToStorageProvider(t *testing.T) {
	type test struct {
		input string
		want  StorageProvider
	}

	valid := []test{
		{input: "s3", want: StorageProviderS3},
		{input: "S3", want: StorageProviderS3},
		{input: "s3gov", want: StorageProviderS3GOV},
		{input: "S3GOV", want: StorageProviderS3GOV},
		{input: "gcs", want: StorageProviderGCS},
		{input: "GCS", want: StorageProviderGCS},
		{input: "azure", want: StorageProviderAzure},
		{input: "AZURE", want: StorageProviderAzure},
	}

	invalid := []test{
		{input: ""},
		{input: "foo"},
		// TODO(next PRs): make this supported in the resource
		{input: "s3compat", want: StorageProviderS3Compatible},
	}

	for _, tc := range valid {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ToStorageProvider(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}

	for _, tc := range invalid {
		t.Run(tc.input, func(t *testing.T) {
			_, err := ToStorageProvider(tc.input)
			require.Error(t, err)
		})
	}
}

func Test_ExternalVolumes_ToS3StorageProvider(t *testing.T) {
	type test struct {
		input string
		want  S3StorageProvider
	}

	valid := []test{
		{input: "s3", want: S3StorageProviderS3},
		{input: "S3", want: S3StorageProviderS3},
		{input: "s3gov", want: S3StorageProviderS3GOV},
		{input: "S3GOV", want: S3StorageProviderS3GOV},
	}

	invalid := []test{
		{input: ""},
		{input: "foo"},
	}

	for _, tc := range valid {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ToS3StorageProvider(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}

	for _, tc := range invalid {
		t.Run(tc.input, func(t *testing.T) {
			_, err := ToS3StorageProvider(tc.input)
			require.Error(t, err)
		})
	}
}

func Test_ExternalVolumes_ToGCSEncryptionType(t *testing.T) {
	type test struct {
		input string
		want  GCSEncryptionType
	}

	valid := []test{
		{input: "gcs_sse_kms", want: GCSEncryptionTypeSseKms},
		{input: "GCS_SSE_KMS", want: GCSEncryptionTypeSseKms},
		{input: "NONE", want: GCSEncryptionTypeNone},
		{input: "none", want: GCSEncryptionTypeNone},
	}

	invalid := []test{
		{input: ""},
		{input: "foo"},
	}

	for _, tc := range valid {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ToGCSEncryptionType(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}

	for _, tc := range invalid {
		t.Run(tc.input, func(t *testing.T) {
			_, err := ToGCSEncryptionType(tc.input)
			require.Error(t, err)
		})
	}
}

// External volume helper tests

var s3StorageAwsExternalId = "1234567890"

func Test_CopySentinelStorageLocation(t *testing.T) {
	tempStorageLocationName := "terraform_provider_sentinel_storage_location"
	s3StorageLocationName := "s3Test"
	s3StorageBaseUrl := "s3://my_example_bucket"
	s3StorageAwsRoleArn := "arn:aws:iam::123456789012:role/myrole"
	s3EncryptionKmsKeyId := "123456789"
	s3StorageAwsAccessPointArn := "arn:aws:s3:us-west-2:123456789012:accesspoint/my-access-point"

	gcsStorageLocationName := "gcsTest"
	gcsStorageBaseUrl := "gcs://my_example_bucket"
	gcsEncryptionKmsKeyId := "123456789"

	azureStorageLocationName := "azureTest"
	azureStorageBaseUrl := "azure://123456789.blob.core.windows.net/my_example_container"
	azureTenantId := "123456789"

	s3StorageLocationA := S3StorageLocationParams{
		StorageProvider:          S3StorageProviderS3,
		StorageBaseUrl:           s3StorageBaseUrl,
		StorageAwsRoleArn:        s3StorageAwsRoleArn,
		StorageAwsExternalId:     &s3StorageAwsExternalId,
		StorageAwsAccessPointArn: &s3StorageAwsAccessPointArn,
		UsePrivatelinkEndpoint:   Bool(true),
		Encryption: &ExternalVolumeS3Encryption{
			EncryptionType: S3EncryptionTypeSseKms,
			KmsKeyId:       &s3EncryptionKmsKeyId,
		},
	}

	azureStorageLocationA := AzureStorageLocationParams{
		StorageBaseUrl:         azureStorageBaseUrl,
		AzureTenantId:          azureTenantId,
		UsePrivatelinkEndpoint: Bool(true),
	}

	gcsStorageLocationA := GCSStorageLocationParams{
		StorageBaseUrl: gcsStorageBaseUrl,
		Encryption: &ExternalVolumeGCSEncryption{
			EncryptionType: GCSEncryptionTypeSseKms,
			KmsKeyId:       &gcsEncryptionKmsKeyId,
		},
	}

	t.Run("S3 storage location", func(t *testing.T) {
		storageLocationInput := ExternalVolumeStorageLocation{Name: s3StorageLocationName, S3StorageLocationParams: &s3StorageLocationA}
		copiedStorageLocationItem, err := CopySentinelStorageLocation(ExternalVolumeStorageLocationItem{ExternalVolumeStorageLocation: storageLocationInput})
		require.NoError(t, err)
		copiedStorageLocation := copiedStorageLocationItem.ExternalVolumeStorageLocation
		assert.Equal(t, copiedStorageLocation.Name, tempStorageLocationName)
		assert.Equal(t, copiedStorageLocation.S3StorageLocationParams.StorageProvider, s3StorageLocationA.StorageProvider)
		assert.Equal(t, copiedStorageLocation.S3StorageLocationParams.StorageBaseUrl, s3StorageLocationA.StorageBaseUrl)
		assert.Equal(t, copiedStorageLocation.S3StorageLocationParams.StorageAwsRoleArn, s3StorageLocationA.StorageAwsRoleArn)
		assert.Equal(t, copiedStorageLocation.S3StorageLocationParams.StorageAwsExternalId, s3StorageLocationA.StorageAwsExternalId)
		assert.Equal(t, copiedStorageLocation.S3StorageLocationParams.StorageAwsAccessPointArn, s3StorageLocationA.StorageAwsAccessPointArn)
		assert.Equal(t, copiedStorageLocation.S3StorageLocationParams.UsePrivatelinkEndpoint, s3StorageLocationA.UsePrivatelinkEndpoint)
		assert.Equal(t, copiedStorageLocation.S3StorageLocationParams.Encryption.EncryptionType, s3StorageLocationA.Encryption.EncryptionType)
		assert.Equal(t, *copiedStorageLocation.S3StorageLocationParams.Encryption.KmsKeyId, *s3StorageLocationA.Encryption.KmsKeyId)
	})

	t.Run("GCS storage location", func(t *testing.T) {
		storageLocationInput := ExternalVolumeStorageLocation{Name: gcsStorageLocationName, GCSStorageLocationParams: &gcsStorageLocationA}
		copiedStorageLocationItem, err := CopySentinelStorageLocation(ExternalVolumeStorageLocationItem{ExternalVolumeStorageLocation: storageLocationInput})
		copiedStorageLocation := copiedStorageLocationItem.ExternalVolumeStorageLocation
		require.NoError(t, err)
		assert.Equal(t, copiedStorageLocation.Name, tempStorageLocationName)
		assert.Equal(t, copiedStorageLocation.GCSStorageLocationParams.StorageBaseUrl, gcsStorageLocationA.StorageBaseUrl)
		assert.Equal(t, copiedStorageLocation.GCSStorageLocationParams.Encryption.EncryptionType, gcsStorageLocationA.Encryption.EncryptionType)
		assert.Equal(t, *copiedStorageLocation.GCSStorageLocationParams.Encryption.KmsKeyId, *gcsStorageLocationA.Encryption.KmsKeyId)
	})

	t.Run("Azure storage location", func(t *testing.T) {
		storageLocationInput := ExternalVolumeStorageLocation{Name: azureStorageLocationName, AzureStorageLocationParams: &azureStorageLocationA}
		copiedStorageLocationItem, err := CopySentinelStorageLocation(ExternalVolumeStorageLocationItem{ExternalVolumeStorageLocation: storageLocationInput})
		copiedStorageLocation := copiedStorageLocationItem.ExternalVolumeStorageLocation
		require.NoError(t, err)
		assert.Equal(t, copiedStorageLocation.Name, tempStorageLocationName)
		assert.Equal(t, copiedStorageLocation.AzureStorageLocationParams.StorageBaseUrl, azureStorageLocationA.StorageBaseUrl)
		assert.Equal(t, copiedStorageLocation.AzureStorageLocationParams.AzureTenantId, azureStorageLocationA.AzureTenantId)
		assert.Equal(t, copiedStorageLocation.AzureStorageLocationParams.UsePrivatelinkEndpoint, azureStorageLocationA.UsePrivatelinkEndpoint)
	})

	s3CompatStorageLocationA := S3CompatStorageLocationParams{
		StorageBaseUrl:  "s3compat://my-bucket/my-path",
		StorageEndpoint: "https://s3-compatible.example.com",
		Credentials: &ExternalVolumeS3CompatCredentials{
			AwsKeyId:     "some_key_id",
			AwsSecretKey: "some_secret_key",
		},
	}

	t.Run("S3Compatible storage location", func(t *testing.T) {
		storageLocationInput := ExternalVolumeStorageLocation{Name: "s3compatTest", S3CompatStorageLocationParams: &s3CompatStorageLocationA}
		copiedStorageLocationItem, err := CopySentinelStorageLocation(ExternalVolumeStorageLocationItem{ExternalVolumeStorageLocation: storageLocationInput})
		copiedStorageLocation := copiedStorageLocationItem.ExternalVolumeStorageLocation
		require.NoError(t, err)
		assert.Equal(t, copiedStorageLocation.Name, tempStorageLocationName)
		assert.Equal(t, copiedStorageLocation.S3CompatStorageLocationParams.StorageBaseUrl, s3CompatStorageLocationA.StorageBaseUrl)
		assert.Equal(t, copiedStorageLocation.S3CompatStorageLocationParams.StorageEndpoint, s3CompatStorageLocationA.StorageEndpoint)
		assert.Equal(t, copiedStorageLocation.S3CompatStorageLocationParams.Credentials.AwsKeyId, s3CompatStorageLocationA.Credentials.AwsKeyId)
		assert.Equal(t, copiedStorageLocation.S3CompatStorageLocationParams.Credentials.AwsSecretKey, s3CompatStorageLocationA.Credentials.AwsSecretKey)
	})

	invalidTestCases := []struct {
		Name            string
		StorageLocation ExternalVolumeStorageLocation
	}{
		{
			Name:            "Empty S3 storage location",
			StorageLocation: ExternalVolumeStorageLocation{S3StorageLocationParams: &S3StorageLocationParams{}},
		},
		{
			Name:            "Empty GCS storage location",
			StorageLocation: ExternalVolumeStorageLocation{GCSStorageLocationParams: &GCSStorageLocationParams{}},
		},
		{
			Name:            "Empty Azure storage location",
			StorageLocation: ExternalVolumeStorageLocation{AzureStorageLocationParams: &AzureStorageLocationParams{}},
		},
		{
			Name:            "Empty S3Compatible storage location",
			StorageLocation: ExternalVolumeStorageLocation{S3CompatStorageLocationParams: &S3CompatStorageLocationParams{}},
		},
		{
			Name:            "Empty storage location",
			StorageLocation: ExternalVolumeStorageLocation{},
		},
	}

	for _, tc := range invalidTestCases {
		t.Run(tc.Name, func(t *testing.T) {
			_, err := CopySentinelStorageLocation(ExternalVolumeStorageLocationItem{ExternalVolumeStorageLocation: tc.StorageLocation})
			require.Error(t, err)
		})
	}
}

func Test_ParseExternalVolumeDescribed(t *testing.T) {
	azureStorageLocationName := "azureTest"
	azureStorageProvider := "AZURE"
	azureStorageBaseUrl := "azure://123456789.blob.core.windows.net/my_example_container"
	azureTenantId := "123456789"
	azureEncryptionTypeNone := "NONE"
	azureStorageLocationStandard := fmt.Sprintf(
		`{"NAME":"%s","STORAGE_PROVIDER":"%s","STORAGE_BASE_URL":"%s","STORAGE_ALLOWED_LOCATIONS":["azure://123456789.blob.core.windows.net/my_example_container"],"AZURE_TENANT_ID":"%s","AZURE_MULTI_TENANT_APP_NAME":"test12","AZURE_CONSENT_URL":"https://login.microsoftonline.com/123456789/oauth2/authorize?client_id=test&response_type=test","ENCRYPTION_TYPE":"NONE","ENCRYPTION_KMS_KEY_ID":""}`,
		azureStorageLocationName,
		azureStorageProvider,
		azureStorageBaseUrl,
		azureTenantId,
	)

	azureStorageLocationWithExtraFields := fmt.Sprintf(
		`{"NAME":"%s","STORAGE_PROVIDER":"%s","STORAGE_BASE_URL":"%s","STORAGE_ALLOWED_LOCATIONS":["azure://123456789.blob.core.windows.net/my_example_container"],"AZURE_TENANT_ID":"%s","AZURE_MULTI_TENANT_APP_NAME":"test12","AZURE_CONSENT_URL":"https://login.microsoftonline.com/123456789/oauth2/authorize?client_id=test&response_type=test","ENCRYPTION_TYPE":"NONE","ENCRYPTION_KMS_KEY_ID":"","EXTRA_FIELD_ONE":"testing","EXTRA_FIELD_TWO":"123456"}`,
		azureStorageLocationName,
		azureStorageProvider,
		azureStorageBaseUrl,
		azureTenantId,
	)

	azureStorageLocationMissingTenantId := fmt.Sprintf(
		`{"NAME":"%s","STORAGE_PROVIDER":"%s","STORAGE_BASE_URL":"%s","STORAGE_ALLOWED_LOCATIONS":["azure://123456789.blob.core.windows.net/my_example_container"],"AZURE_MULTI_TENANT_APP_NAME":"test12","AZURE_CONSENT_URL":"https://login.microsoftonline.com/123456789/oauth2/authorize?client_id=test&response_type=test","ENCRYPTION_TYPE":"NONE","ENCRYPTION_KMS_KEY_ID":""}`,
		azureStorageLocationName,
		azureStorageProvider,
		azureStorageBaseUrl,
	)

	gcsStorageLocationName := "gcsTest"
	gcsStorageProvider := "GCS"
	gcsStorageBaseUrl := "gcs://my_example_bucket"
	gcsEncryptionTypeNone := "NONE"
	gcsEncryptionTypeSseKms := "GCS_SSE_KMS"
	gcsEncryptionKmsKeyId := "123456789"
	gcsStorageLocationStandard := fmt.Sprintf(
		`{"NAME":"%s","STORAGE_PROVIDER":"%s","STORAGE_BASE_URL":"%s","STORAGE_ALLOWED_LOCATIONS":["gcs://my_example_bucket/*"],"STORAGE_GCP_SERVICE_ACCOUNT":"test@test.iam.test.com","ENCRYPTION_TYPE":"%s","ENCRYPTION_KMS_KEY_ID":""}`,
		gcsStorageLocationName,
		gcsStorageProvider,
		gcsStorageBaseUrl,
		gcsEncryptionTypeNone,
	)

	gcsStorageLocationWithExtraFields := fmt.Sprintf(
		`{"NAME":"%s","STORAGE_PROVIDER":"%s","STORAGE_BASE_URL":"%s","STORAGE_ALLOWED_LOCATIONS":["gcs://my_example_bucket/*"],"STORAGE_GCP_SERVICE_ACCOUNT":"test@test.iam.test.com","ENCRYPTION_TYPE":"%s","ENCRYPTION_KMS_KEY_ID":"","EXTRA_FIELD_ONE":"testing","EXTRA_FIELD_TWO":"123456"}`,
		gcsStorageLocationName,
		gcsStorageProvider,
		gcsStorageBaseUrl,
		gcsEncryptionTypeNone,
	)

	gcsStorageLocationKmsEncryption := fmt.Sprintf(
		`{"NAME":"%s","STORAGE_PROVIDER":"%s","STORAGE_BASE_URL":"%s","STORAGE_ALLOWED_LOCATIONS":["gcs://my_example_bucket/*"],"STORAGE_GCP_SERVICE_ACCOUNT":"test@test.iam.test.com","ENCRYPTION_TYPE":"%s","ENCRYPTION_KMS_KEY_ID":"%s"}`,
		gcsStorageLocationName,
		gcsStorageProvider,
		gcsStorageBaseUrl,
		gcsEncryptionTypeSseKms,
		gcsEncryptionKmsKeyId,
	)

	gcsStorageLocationMissingBaseUrl := fmt.Sprintf(
		`{"NAME":"%s","STORAGE_PROVIDER":"%s","STORAGE_ALLOWED_LOCATIONS":["gcs://my_example_bucket/*"],"STORAGE_GCP_SERVICE_ACCOUNT":"test@test.iam.test.com","ENCRYPTION_TYPE":"%s","ENCRYPTION_KMS_KEY_ID":""}`,
		gcsStorageLocationName,
		gcsStorageProvider,
		gcsEncryptionTypeNone,
	)

	s3StorageLocationName := "s3Test"
	s3StorageProvider := "S3"
	s3StorageBaseUrl := "s3://my_example_bucket"
	s3StorageAwsRoleArn := "arn:aws:iam::123456789012:role/myrole"
	s3StorageAwsExternalId := "123456789"
	s3EncryptionTypeNone := "NONE"
	s3EncryptionTypeSseS3 := "AWS_SSE_S3"
	s3EncryptionTypeSseKms := "AWS_SSE_KMS"
	s3EncryptionKmsKeyId := "123456789"

	s3StorageLocationStandard := fmt.Sprintf(
		`{"NAME":"%s","STORAGE_PROVIDER":"%s","STORAGE_BASE_URL":"%s","STORAGE_ALLOWED_LOCATIONS":["s3://my_example_bucket/*"],"STORAGE_AWS_ROLE_ARN":"%s","STORAGE_AWS_IAM_USER_ARN":"arn:aws:iam::123456789:user/a11b0000-s","STORAGE_AWS_EXTERNAL_ID":"%s","ENCRYPTION_TYPE":"%s"}`,
		s3StorageLocationName,
		s3StorageProvider,
		s3StorageBaseUrl,
		s3StorageAwsRoleArn,
		s3StorageAwsExternalId,
		s3EncryptionTypeNone,
	)

	s3StorageLocationWithExtraFields := fmt.Sprintf(
		`{"NAME":"%s","STORAGE_PROVIDER":"%s","STORAGE_BASE_URL":"%s","STORAGE_ALLOWED_LOCATIONS":["s3://my_example_bucket/*"],"STORAGE_AWS_ROLE_ARN":"%s","STORAGE_AWS_IAM_USER_ARN":"arn:aws:iam::123456789:user/a11b0000-s","STORAGE_AWS_EXTERNAL_ID":"%s","ENCRYPTION_TYPE":"%s","ENCRYPTION_KMS_KEY_ID":"%s","EXTRA_FIELD_ONE":"testing","EXTRA_FIELD_TWO":"123456"}`,
		s3StorageLocationName,
		s3StorageProvider,
		s3StorageBaseUrl,
		s3StorageAwsRoleArn,
		s3StorageAwsExternalId,
		s3EncryptionTypeSseKms,
		s3EncryptionKmsKeyId,
	)

	s3StorageLocationSseS3Encryption := fmt.Sprintf(
		`{"NAME":"%s","STORAGE_PROVIDER":"%s","STORAGE_BASE_URL":"%s","STORAGE_ALLOWED_LOCATIONS":["s3://my_example_bucket/*"],"STORAGE_AWS_ROLE_ARN":"%s","STORAGE_AWS_IAM_USER_ARN":"arn:aws:iam::123456789:user/a11b0000-s","STORAGE_AWS_EXTERNAL_ID":"%s","ENCRYPTION_TYPE":"%s"}`,
		s3StorageLocationName,
		s3StorageProvider,
		s3StorageBaseUrl,
		s3StorageAwsRoleArn,
		s3StorageAwsExternalId,
		s3EncryptionTypeSseS3,
	)

	s3StorageLocationSseKmsEncryption := fmt.Sprintf(
		`{"NAME":"%s","STORAGE_PROVIDER":"%s","STORAGE_BASE_URL":"%s","STORAGE_ALLOWED_LOCATIONS":["s3://my_example_bucket/*"],"STORAGE_AWS_ROLE_ARN":"%s","STORAGE_AWS_IAM_USER_ARN":"arn:aws:iam::123456789:user/a11b0000-s","STORAGE_AWS_EXTERNAL_ID":"%s","ENCRYPTION_TYPE":"%s", "ENCRYPTION_KMS_KEY_ID":"%s"}`,
		s3StorageLocationName,
		s3StorageProvider,
		s3StorageBaseUrl,
		s3StorageAwsRoleArn,
		s3StorageAwsExternalId,
		s3EncryptionTypeSseKms,
		s3EncryptionKmsKeyId,
	)

	s3StorageLocationMissingRoleArn := fmt.Sprintf(
		`{"NAME":"%s","STORAGE_PROVIDER":"%s","STORAGE_BASE_URL":"%s","STORAGE_ALLOWED_LOCATIONS":["s3://my_example_bucket/*"],"STORAGE_AWS_IAM_USER_ARN":"arn:aws:iam::123456789:user/a11b0000-s","STORAGE_AWS_EXTERNAL_ID":"%s","ENCRYPTION_TYPE":"%s"}`,
		s3StorageLocationName,
		s3StorageProvider,
		s3StorageBaseUrl,
		s3StorageAwsExternalId,
		s3EncryptionTypeNone,
	)

	s3CompatStorageLocationName := "s3compatTest"
	s3CompatStorageProvider := "S3COMPAT"
	s3CompatStorageBaseUrl := "s3compat://my_example_bucket"
	s3CompatStorageEndpoint := "https://s3-compatible.example.com"
	s3CompatStorageLocationStandard := fmt.Sprintf(
		`{"NAME":"%s","STORAGE_PROVIDER":"%s","STORAGE_BASE_URL":"%s","STORAGE_ENDPOINT":"%s","ENCRYPTION_TYPE":"NONE","ENCRYPTION_KMS_KEY_ID":""}`,
		s3CompatStorageLocationName,
		s3CompatStorageProvider,
		s3CompatStorageBaseUrl,
		s3CompatStorageEndpoint,
	)

	s3CompatStorageLocationMissingEndpoint := fmt.Sprintf(
		`{"NAME":"%s","STORAGE_PROVIDER":"%s","STORAGE_BASE_URL":"%s","ENCRYPTION_TYPE":"NONE","ENCRYPTION_KMS_KEY_ID":""}`,
		s3CompatStorageLocationName,
		s3CompatStorageProvider,
		s3CompatStorageBaseUrl,
	)

	allowWritesTrue := "true"
	allowWritesFalse := "false"
	comment := "some comment"
	validCases := []struct {
		Name                 string
		DescribeOutput       []ExternalVolumeProperty
		ParsedDescribeOutput ExternalVolumeDetails
	}{
		{
			Name:           "Volume with azure storage location",
			DescribeOutput: GenerateParseExternalVolumeDescribedInput(comment, allowWritesFalse, []string{azureStorageLocationStandard}, ""),
			ParsedDescribeOutput: ExternalVolumeDetails{
				StorageLocations: []ExternalVolumeStorageLocationDetails{
					{
						Name:                    azureStorageLocationName,
						StorageProvider:         azureStorageProvider,
						StorageBaseUrl:          azureStorageBaseUrl,
						StorageAllowedLocations: []string{"azure://123456789.blob.core.windows.net/my_example_container"},
						EncryptionType:          azureEncryptionTypeNone,
						AzureStorageLocation: &StorageLocationAzureDetails{
							AzureTenantId:           azureTenantId,
							AzureMultiTenantAppName: "test12",
							AzureConsentUrl:         "https://login.microsoftonline.com/123456789/oauth2/authorize?client_id=test&response_type=test",
						},
					},
				},
				Active:      "",
				Comment:     comment,
				AllowWrites: allowWritesFalse,
			},
		},
		{
			Name:           "Volume with azure storage location, with extra fields",
			DescribeOutput: GenerateParseExternalVolumeDescribedInput(comment, allowWritesFalse, []string{azureStorageLocationWithExtraFields}, ""),
			ParsedDescribeOutput: ExternalVolumeDetails{
				StorageLocations: []ExternalVolumeStorageLocationDetails{
					{
						Name:                    azureStorageLocationName,
						StorageProvider:         azureStorageProvider,
						StorageBaseUrl:          azureStorageBaseUrl,
						StorageAllowedLocations: []string{"azure://123456789.blob.core.windows.net/my_example_container"},
						EncryptionType:          azureEncryptionTypeNone,
						AzureStorageLocation: &StorageLocationAzureDetails{
							AzureTenantId:           azureTenantId,
							AzureMultiTenantAppName: "test12",
							AzureConsentUrl:         "https://login.microsoftonline.com/123456789/oauth2/authorize?client_id=test&response_type=test",
						},
					},
				},
				Active:      "",
				Comment:     comment,
				AllowWrites: allowWritesFalse,
			},
		},
		{
			Name:           "Volume with gcs storage location",
			DescribeOutput: GenerateParseExternalVolumeDescribedInput(comment, allowWritesTrue, []string{gcsStorageLocationStandard}, ""),
			ParsedDescribeOutput: ExternalVolumeDetails{
				StorageLocations: []ExternalVolumeStorageLocationDetails{
					{
						Name:                    gcsStorageLocationName,
						StorageProvider:         gcsStorageProvider,
						StorageBaseUrl:          gcsStorageBaseUrl,
						StorageAllowedLocations: []string{"gcs://my_example_bucket/*"},
						EncryptionType:          gcsEncryptionTypeNone,
						GCSStorageLocation: &StorageLocationGcsDetails{
							StorageGcpServiceAccount: "test@test.iam.test.com",
						},
					},
				},
				Active:      "",
				Comment:     comment,
				AllowWrites: allowWritesTrue,
			},
		},
		{
			Name:           "Volume with gcs storage location, with extra fields",
			DescribeOutput: GenerateParseExternalVolumeDescribedInput(comment, allowWritesTrue, []string{gcsStorageLocationWithExtraFields}, ""),
			ParsedDescribeOutput: ExternalVolumeDetails{
				StorageLocations: []ExternalVolumeStorageLocationDetails{
					{
						Name:                    gcsStorageLocationName,
						StorageProvider:         gcsStorageProvider,
						StorageBaseUrl:          gcsStorageBaseUrl,
						StorageAllowedLocations: []string{"gcs://my_example_bucket/*"},
						EncryptionType:          gcsEncryptionTypeNone,
						GCSStorageLocation: &StorageLocationGcsDetails{
							StorageGcpServiceAccount: "test@test.iam.test.com",
						},
					},
				},
				Active:      "",
				Comment:     comment,
				AllowWrites: allowWritesTrue,
			},
		},
		{
			Name:           "Volume with gcs storage location, sse kms encryption",
			DescribeOutput: GenerateParseExternalVolumeDescribedInput(comment, allowWritesTrue, []string{gcsStorageLocationKmsEncryption}, ""),
			ParsedDescribeOutput: ExternalVolumeDetails{
				StorageLocations: []ExternalVolumeStorageLocationDetails{
					{
						Name:                    gcsStorageLocationName,
						StorageProvider:         gcsStorageProvider,
						StorageBaseUrl:          gcsStorageBaseUrl,
						StorageAllowedLocations: []string{"gcs://my_example_bucket/*"},
						EncryptionType:          gcsEncryptionTypeSseKms,
						GCSStorageLocation: &StorageLocationGcsDetails{
							StorageGcpServiceAccount: "test@test.iam.test.com",
							EncryptionKmsKeyId:       gcsEncryptionKmsKeyId,
						},
					},
				},
				Active:      "",
				Comment:     comment,
				AllowWrites: allowWritesTrue,
			},
		},
		{
			Name:           "Volume with s3 storage location",
			DescribeOutput: GenerateParseExternalVolumeDescribedInput(comment, allowWritesTrue, []string{s3StorageLocationStandard}, ""),
			ParsedDescribeOutput: ExternalVolumeDetails{
				StorageLocations: []ExternalVolumeStorageLocationDetails{
					{
						Name:                    s3StorageLocationName,
						StorageProvider:         s3StorageProvider,
						StorageBaseUrl:          s3StorageBaseUrl,
						StorageAllowedLocations: []string{"s3://my_example_bucket/*"},
						EncryptionType:          s3EncryptionTypeNone,
						S3StorageLocation: &StorageLocationS3Details{
							StorageAwsRoleArn:    s3StorageAwsRoleArn,
							StorageAwsIamUserArn: "arn:aws:iam::123456789:user/a11b0000-s",
							StorageAwsExternalId: s3StorageAwsExternalId,
						},
					},
				},
				Active:      "",
				Comment:     comment,
				AllowWrites: allowWritesTrue,
			},
		},
		{
			Name:           "Volume with s3 storage location, with extra fields",
			DescribeOutput: GenerateParseExternalVolumeDescribedInput(comment, allowWritesTrue, []string{s3StorageLocationWithExtraFields}, ""),
			ParsedDescribeOutput: ExternalVolumeDetails{
				StorageLocations: []ExternalVolumeStorageLocationDetails{
					{
						Name:                    s3StorageLocationName,
						StorageProvider:         s3StorageProvider,
						StorageBaseUrl:          s3StorageBaseUrl,
						StorageAllowedLocations: []string{"s3://my_example_bucket/*"},
						EncryptionType:          s3EncryptionTypeSseKms,
						S3StorageLocation: &StorageLocationS3Details{
							StorageAwsRoleArn:    s3StorageAwsRoleArn,
							StorageAwsIamUserArn: "arn:aws:iam::123456789:user/a11b0000-s",
							StorageAwsExternalId: s3StorageAwsExternalId,
							EncryptionKmsKeyId:   s3EncryptionKmsKeyId,
						},
					},
				},
				Active:      "",
				Comment:     comment,
				AllowWrites: allowWritesTrue,
			},
		},
		{
			Name:           "Volume with s3 storage location, sse s3 encryption",
			DescribeOutput: GenerateParseExternalVolumeDescribedInput(comment, allowWritesTrue, []string{s3StorageLocationSseS3Encryption}, ""),
			ParsedDescribeOutput: ExternalVolumeDetails{
				StorageLocations: []ExternalVolumeStorageLocationDetails{
					{
						Name:                    s3StorageLocationName,
						StorageProvider:         s3StorageProvider,
						StorageBaseUrl:          s3StorageBaseUrl,
						StorageAllowedLocations: []string{"s3://my_example_bucket/*"},
						EncryptionType:          s3EncryptionTypeSseS3,
						S3StorageLocation: &StorageLocationS3Details{
							StorageAwsRoleArn:    s3StorageAwsRoleArn,
							StorageAwsIamUserArn: "arn:aws:iam::123456789:user/a11b0000-s",
							StorageAwsExternalId: s3StorageAwsExternalId,
						},
					},
				},
				Active:      "",
				Comment:     comment,
				AllowWrites: allowWritesTrue,
			},
		},
		{
			Name:           "Volume with s3 storage location, sse kms encryption",
			DescribeOutput: GenerateParseExternalVolumeDescribedInput(comment, allowWritesTrue, []string{s3StorageLocationSseKmsEncryption}, ""),
			ParsedDescribeOutput: ExternalVolumeDetails{
				StorageLocations: []ExternalVolumeStorageLocationDetails{
					{
						Name:                    s3StorageLocationName,
						StorageProvider:         s3StorageProvider,
						StorageBaseUrl:          s3StorageBaseUrl,
						StorageAllowedLocations: []string{"s3://my_example_bucket/*"},
						EncryptionType:          s3EncryptionTypeSseKms,
						S3StorageLocation: &StorageLocationS3Details{
							StorageAwsRoleArn:    s3StorageAwsRoleArn,
							StorageAwsIamUserArn: "arn:aws:iam::123456789:user/a11b0000-s",
							StorageAwsExternalId: s3StorageAwsExternalId,
							EncryptionKmsKeyId:   s3EncryptionKmsKeyId,
						},
					},
				},
				Active:      "",
				Comment:     comment,
				AllowWrites: allowWritesTrue,
			},
		},
		{
			Name: "Volume with multiple storage locations and active set",
			DescribeOutput: GenerateParseExternalVolumeDescribedInput(
				comment,
				allowWritesTrue,
				[]string{s3StorageLocationStandard, gcsStorageLocationStandard, azureStorageLocationStandard},
				s3StorageLocationName,
			),
			ParsedDescribeOutput: ExternalVolumeDetails{
				StorageLocations: []ExternalVolumeStorageLocationDetails{
					{
						Name:                    s3StorageLocationName,
						StorageProvider:         s3StorageProvider,
						StorageBaseUrl:          s3StorageBaseUrl,
						StorageAllowedLocations: []string{"s3://my_example_bucket/*"},
						EncryptionType:          s3EncryptionTypeNone,
						S3StorageLocation: &StorageLocationS3Details{
							StorageAwsRoleArn:    s3StorageAwsRoleArn,
							StorageAwsIamUserArn: "arn:aws:iam::123456789:user/a11b0000-s",
							StorageAwsExternalId: s3StorageAwsExternalId,
						},
					},
					{
						Name:                    gcsStorageLocationName,
						StorageProvider:         gcsStorageProvider,
						StorageBaseUrl:          gcsStorageBaseUrl,
						StorageAllowedLocations: []string{"gcs://my_example_bucket/*"},
						EncryptionType:          gcsEncryptionTypeNone,
						GCSStorageLocation: &StorageLocationGcsDetails{
							StorageGcpServiceAccount: "test@test.iam.test.com",
						},
					},
					{
						Name:                    azureStorageLocationName,
						StorageProvider:         azureStorageProvider,
						StorageBaseUrl:          azureStorageBaseUrl,
						StorageAllowedLocations: []string{"azure://123456789.blob.core.windows.net/my_example_container"},
						EncryptionType:          azureEncryptionTypeNone,
						AzureStorageLocation: &StorageLocationAzureDetails{
							AzureTenantId:           azureTenantId,
							AzureMultiTenantAppName: "test12",
							AzureConsentUrl:         "https://login.microsoftonline.com/123456789/oauth2/authorize?client_id=test&response_type=test",
						},
					},
				},
				Active:      s3StorageLocationName,
				Comment:     comment,
				AllowWrites: allowWritesTrue,
			},
		},
		{
			Name: "Volume with s3 storage location that has no comment set (in this case describe doesn't contain a comment property)",
			DescribeOutput: []ExternalVolumeProperty{
				{
					Parent:  "",
					Name:    "ALLOW_WRITES",
					Type:    "Boolean",
					Value:   allowWritesTrue,
					Default: "true",
				},
				{
					Parent:  "STORAGE_LOCATIONS",
					Name:    "STORAGE_LOCATION_1",
					Type:    "String",
					Value:   s3StorageLocationSseKmsEncryption,
					Default: "",
				},
				{
					Parent:  "STORAGE_LOCATIONS",
					Name:    "ACTIVE",
					Type:    "String",
					Value:   s3StorageLocationName,
					Default: "",
				},
			},
			ParsedDescribeOutput: ExternalVolumeDetails{
				StorageLocations: []ExternalVolumeStorageLocationDetails{
					{
						Name:                    s3StorageLocationName,
						StorageProvider:         s3StorageProvider,
						StorageBaseUrl:          s3StorageBaseUrl,
						StorageAllowedLocations: []string{"s3://my_example_bucket/*"},
						EncryptionType:          s3EncryptionTypeSseKms,
						S3StorageLocation: &StorageLocationS3Details{
							StorageAwsRoleArn:    s3StorageAwsRoleArn,
							StorageAwsIamUserArn: "arn:aws:iam::123456789:user/a11b0000-s",
							StorageAwsExternalId: s3StorageAwsExternalId,
							EncryptionKmsKeyId:   s3EncryptionKmsKeyId,
						},
					},
				},
				Active:      s3StorageLocationName,
				Comment:     "",
				AllowWrites: allowWritesTrue,
			},
		},
		{
			Name:           "Volume with s3compat storage location",
			DescribeOutput: GenerateParseExternalVolumeDescribedInput(comment, allowWritesTrue, []string{s3CompatStorageLocationStandard}, ""),
			ParsedDescribeOutput: ExternalVolumeDetails{
				StorageLocations: []ExternalVolumeStorageLocationDetails{
					{
						Name:            s3CompatStorageLocationName,
						StorageProvider: s3CompatStorageProvider,
						StorageBaseUrl:  s3CompatStorageBaseUrl,
						EncryptionType:  "NONE",
						S3CompatStorageLocation: &StorageLocationS3CompatDetails{
							StorageEndpoint: s3CompatStorageEndpoint,
						},
					},
				},
				Active:      "",
				Comment:     comment,
				AllowWrites: allowWritesTrue,
			},
		},
	}

	invalidCases := []struct {
		Name           string
		DescribeOutput []ExternalVolumeProperty
	}{
		{
			Name:           "Volume with s3 storage location, missing STORAGE_AWS_ROLE_ARN",
			DescribeOutput: GenerateParseExternalVolumeDescribedInput(comment, allowWritesTrue, []string{s3StorageLocationMissingRoleArn}, ""),
		},
		{
			Name:           "Volume with azure storage location, missing AZURE_TENANT_ID",
			DescribeOutput: GenerateParseExternalVolumeDescribedInput(comment, allowWritesTrue, []string{azureStorageLocationMissingTenantId}, ""),
		},
		{
			Name:           "Volume with gcs storage location, missing STORAGE_BASE_URL",
			DescribeOutput: GenerateParseExternalVolumeDescribedInput(comment, allowWritesTrue, []string{gcsStorageLocationMissingBaseUrl}, ""),
		},
		{
			Name:           "Volume with s3compat storage location, missing STORAGE_ENDPOINT",
			DescribeOutput: GenerateParseExternalVolumeDescribedInput(comment, allowWritesTrue, []string{s3CompatStorageLocationMissingEndpoint}, ""),
		},
		{
			Name:           "Volume with no storage locations",
			DescribeOutput: GenerateParseExternalVolumeDescribedInput(comment, allowWritesTrue, []string{}, ""),
		},
		{
			Name: "Volume with no allow writes",
			DescribeOutput: []ExternalVolumeProperty{
				{
					Parent:  "STORAGE_LOCATIONS",
					Name:    "STORAGE_LOCATION_1",
					Type:    "String",
					Value:   s3StorageLocationSseKmsEncryption,
					Default: "",
				},
				{
					Parent:  "STORAGE_LOCATIONS",
					Name:    "ACTIVE",
					Type:    "String",
					Value:   s3StorageLocationName,
					Default: "",
				},
			},
		},
	}

	for _, tc := range validCases {
		t.Run(tc.Name, func(t *testing.T) {
			parsed, err := ParseExternalVolumeDescribed(tc.DescribeOutput)
			require.NoError(t, err)
			assert.True(t, reflect.DeepEqual(tc.ParsedDescribeOutput, parsed))
		})
	}

	for _, tc := range invalidCases {
		t.Run(tc.Name, func(t *testing.T) {
			_, err := ParseExternalVolumeDescribed(tc.DescribeOutput)
			require.Error(t, err)
		})
	}
}

// Generate input to the ParseExternalVolumeDescribedInput, useful for testing purposes
func GenerateParseExternalVolumeDescribedInput(comment string, allowWrites string, storageLocations []string, active string) []ExternalVolumeProperty {
	storageLocationProperties := make([]ExternalVolumeProperty, len(storageLocations))
	allowWritesProperty := ExternalVolumeProperty{
		Parent:  "",
		Name:    "ALLOW_WRITES",
		Type:    "Boolean",
		Value:   allowWrites,
		Default: "true",
	}

	commentProperty := ExternalVolumeProperty{
		Parent:  "",
		Name:    "COMMENT",
		Type:    "String",
		Value:   comment,
		Default: "",
	}

	activeProperty := ExternalVolumeProperty{
		Parent:  "STORAGE_LOCATIONS",
		Name:    "ACTIVE",
		Type:    "String",
		Value:   active,
		Default: "",
	}

	for i, property := range storageLocations {
		storageLocationProperties[i] = ExternalVolumeProperty{
			Parent:  "STORAGE_LOCATIONS",
			Name:    fmt.Sprintf("STORAGE_LOCATION_%s", strconv.Itoa(i+1)),
			Type:    "String",
			Value:   property,
			Default: "",
		}
	}

	return append(append([]ExternalVolumeProperty{allowWritesProperty, commentProperty}, storageLocationProperties...), activeProperty)
}

func Test_GenerateParseExternalVolumeDescribedInput(t *testing.T) {
	azureStorageLocationName := "azureTest"
	azureStorageProvider := "AZURE"
	azureStorageBaseUrl := "azure://123456789.blob.core.windows.net/my_example_container"
	azureTenantId := "123456789"
	azureStorageLocationStandard := fmt.Sprintf(
		`{"NAME":"%s","STORAGE_PROVIDER":"%s","STORAGE_BASE_URL":"%s","STORAGE_ALLOWED_LOCATIONS":["azure://123456789.blob.core.windows.net/my_example_container"],"AZURE_TENANT_ID":"%s","AZURE_MULTI_TENANT_APP_NAME":"test12","AZURE_CONSENT_URL":"https://login.microsoftonline.com/123456789/oauth2/authorize?client_id=test&response_type=test","ENCRYPTION_TYPE":"NONE","ENCRYPTION_KMS_KEY_ID":""}`,
		azureStorageLocationName,
		azureStorageProvider,
		azureStorageBaseUrl,
		azureTenantId,
	)

	gcsStorageLocationName := "gcsTest"
	gcsStorageProvider := "GCS"
	gcsStorageBaseUrl := "gcs://my_example_bucket"
	gcsEncryptionTypeNone := "NONE"
	gcsStorageLocationStandard := fmt.Sprintf(
		`{"NAME":"%s","STORAGE_PROVIDER":"%s","STORAGE_BASE_URL":"%s","STORAGE_ALLOWED_LOCATIONS":["gcs://my_example_bucket/*"],"STORAGE_GCP_SERVICE_ACCOUNT":"test@test.iam.test.com","ENCRYPTION_TYPE":"%s","ENCRYPTION_KMS_KEY_ID":""}`,
		gcsStorageLocationName,
		gcsStorageProvider,
		gcsStorageBaseUrl,
		gcsEncryptionTypeNone,
	)

	s3StorageLocationName := "s3Test"
	s3StorageProvider := "S3"
	s3StorageBaseUrl := "s3://my_example_bucket"
	s3StorageAwsRoleArn := "arn:aws:iam::123456789012:role/myrole"
	s3StorageAwsExternalId := "123456789"
	s3EncryptionTypeNone := "NONE"

	s3StorageLocationStandard := fmt.Sprintf(
		`{"NAME":"%s","STORAGE_PROVIDER":"%s","STORAGE_BASE_URL":"%s","STORAGE_ALLOWED_LOCATIONS":["s3://my_example_bucket/*"],"STORAGE_AWS_ROLE_ARN":"%s","STORAGE_AWS_IAM_USER_ARN":"arn:aws:iam::123456789:user/a11b0000-s","STORAGE_AWS_EXTERNAL_ID":"%s","ENCRYPTION_TYPE":"%s"}`,
		s3StorageLocationName,
		s3StorageProvider,
		s3StorageBaseUrl,
		s3StorageAwsRoleArn,
		s3StorageAwsExternalId,
		s3EncryptionTypeNone,
	)

	allowWritesTrue := "true"
	comment := "some comment"
	cases := []struct {
		TestName         string
		Comment          string
		AllowWrites      string
		StorageLocations []string
		Active           string
		ExpectedOutput   []ExternalVolumeProperty
	}{
		{
			TestName:         "Generate input",
			Comment:          comment,
			AllowWrites:      allowWritesTrue,
			StorageLocations: []string{s3StorageLocationStandard},
			Active:           "",
			ExpectedOutput: []ExternalVolumeProperty{
				{
					Parent:  "",
					Name:    "ALLOW_WRITES",
					Type:    "Boolean",
					Value:   allowWritesTrue,
					Default: "true",
				},
				{
					Parent:  "",
					Name:    "COMMENT",
					Type:    "String",
					Value:   comment,
					Default: "",
				},
				{
					Parent:  "STORAGE_LOCATIONS",
					Name:    "STORAGE_LOCATION_1",
					Type:    "String",
					Value:   s3StorageLocationStandard,
					Default: "",
				},
				{
					Parent:  "STORAGE_LOCATIONS",
					Name:    "ACTIVE",
					Type:    "String",
					Value:   "",
					Default: "",
				},
			},
		},
		{
			TestName:         "Generate input - multiple locations and active set",
			Comment:          comment,
			AllowWrites:      allowWritesTrue,
			StorageLocations: []string{s3StorageLocationStandard, azureStorageLocationStandard, gcsStorageLocationStandard},
			Active:           s3StorageLocationName,
			ExpectedOutput: []ExternalVolumeProperty{
				{
					Parent:  "",
					Name:    "ALLOW_WRITES",
					Type:    "Boolean",
					Value:   allowWritesTrue,
					Default: "true",
				},
				{
					Parent:  "",
					Name:    "COMMENT",
					Type:    "String",
					Value:   comment,
					Default: "",
				},
				{
					Parent:  "STORAGE_LOCATIONS",
					Name:    "STORAGE_LOCATION_1",
					Type:    "String",
					Value:   s3StorageLocationStandard,
					Default: "",
				},
				{
					Parent:  "STORAGE_LOCATIONS",
					Name:    "STORAGE_LOCATION_2",
					Type:    "String",
					Value:   azureStorageLocationStandard,
					Default: "",
				},
				{
					Parent:  "STORAGE_LOCATIONS",
					Name:    "STORAGE_LOCATION_3",
					Type:    "String",
					Value:   gcsStorageLocationStandard,
					Default: "",
				},
				{
					Parent:  "STORAGE_LOCATIONS",
					Name:    "ACTIVE",
					Type:    "String",
					Value:   s3StorageLocationName,
					Default: "",
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.TestName, func(t *testing.T) {
			generatedInput := GenerateParseExternalVolumeDescribedInput(
				tc.Comment,
				tc.AllowWrites,
				tc.StorageLocations,
				tc.Active,
			)

			assert.Len(t, generatedInput, len(tc.ExpectedOutput))
			for i := range generatedInput {
				assert.Equal(t, tc.ExpectedOutput[i].Parent, generatedInput[i].Parent)
				assert.Equal(t, tc.ExpectedOutput[i].Name, generatedInput[i].Name)
				assert.Equal(t, tc.ExpectedOutput[i].Type, generatedInput[i].Type)
				assert.Equal(t, tc.ExpectedOutput[i].Value, generatedInput[i].Value)
				assert.Equal(t, tc.ExpectedOutput[i].Default, generatedInput[i].Default)
			}
		})
	}
}
