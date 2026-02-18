package sdk

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Storage location structs for testing

// Basic variants (required fields only)
var s3StorageLocationParamsBasic = &S3StorageLocationParams{
	Name:              "s3_basic",
	StorageProvider:   S3StorageProviderS3,
	StorageAwsRoleArn: "arn:aws:iam::123456789012:role/myrole",
	StorageBaseUrl:    "s3://my-bucket/path",
}

var gcsStorageLocationParamsBasic = &GCSStorageLocationParams{
	Name:           "gcs_basic",
	StorageBaseUrl: "gcs://my-bucket/path",
}

var azureStorageLocationParamsBasic = &AzureStorageLocationParams{
	Name:           "azure_basic",
	AzureTenantId:  "a123b4cd-1abc-12ab-12ab-1a2b34c5d678",
	StorageBaseUrl: "azure://myaccount.blob.core.windows.net/mycontainer/path",
}

var s3CompatStorageLocationParamsBasic = &S3CompatStorageLocationParams{
	Name:            "s3compat_basic",
	StorageBaseUrl:  "s3compat://my-bucket/path",
	StorageEndpoint: "https://s3-compatible.example.com",
}

// Complete variants (all optional fields)
var s3StorageLocationParamsComplete = &S3StorageLocationParams{
	Name:                     "s3_complete",
	StorageProvider:          S3StorageProviderS3,
	StorageAwsRoleArn:        "arn:aws:iam::123456789012:role/myrole",
	StorageBaseUrl:           "s3://my-bucket/path",
	StorageAwsExternalId:     String("external_id_123"),
	StorageAwsAccessPointArn: String("arn:aws:s3:us-west-2:123456789012:accesspoint/my-access-point"),
	UsePrivatelinkEndpoint:   Bool(true),
	Encryption: &ExternalVolumeS3Encryption{
		EncryptionType: S3EncryptionTypeSseKms,
		KmsKeyId:       String("1234abcd-12ab-34cd-56ef-1234567890ab"),
	},
}

var gcsStorageLocationParamsComplete = &GCSStorageLocationParams{
	Name:           "gcs_complete",
	StorageBaseUrl: "gcs://my-bucket/path",
	Encryption: &ExternalVolumeGCSEncryption{
		EncryptionType: GCSEncryptionTypeSseKms,
		KmsKeyId:       String("1234abcd-12ab-34cd-56ef-1234567890ab"),
	},
}

var azureStorageLocationParamsComplete = &AzureStorageLocationParams{
	Name:                   "azure_complete",
	AzureTenantId:          "a123b4cd-1abc-12ab-12ab-1a2b34c5d678",
	StorageBaseUrl:         "azure://myaccount.blob.core.windows.net/mycontainer/path",
	UsePrivatelinkEndpoint: Bool(true),
}

var s3CompatStorageLocationParamsComplete = &S3CompatStorageLocationParams{
	Name:            "s3compat_complete",
	StorageBaseUrl:  "s3compat://my-bucket/path",
	StorageEndpoint: "https://s3-compatible.example.com",
	Credentials: &ExternalVolumeS3CompatCredentials{
		AwsKeyId:     "AKIAIOSFODNN7EXAMPLE",
		AwsSecretKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
	},
}

// Helper variables for Alter ADD tests

// Basic variants for Alter (required fields only)
var s3StorageLocationParamsAlterBasic = &S3StorageLocationParams{
	Name:              "s3_alter_basic",
	StorageProvider:   S3StorageProviderS3,
	StorageAwsRoleArn: "arn:aws:iam::123456789012:role/myrole",
	StorageBaseUrl:    "s3://my-bucket/alter-path",
}

var gcsStorageLocationParamsAlterBasic = &GCSStorageLocationParams{
	Name:           "gcs_alter_basic",
	StorageBaseUrl: "gcs://my-bucket/alter-path",
}

var azureStorageLocationParamsAlterBasic = &AzureStorageLocationParams{
	Name:           "azure_alter_basic",
	AzureTenantId:  "a123b4cd-1abc-12ab-12ab-1a2b34c5d678",
	StorageBaseUrl: "azure://myaccount.blob.core.windows.net/mycontainer/alter-path",
}

var s3CompatStorageLocationParamsAlterBasic = &S3CompatStorageLocationParams{
	Name:            "s3compat_alter_basic",
	StorageBaseUrl:  "s3compat://my-bucket/alter-path",
	StorageEndpoint: "https://s3-compatible.example.com",
}

// Complete variants for Alter (all optional fields)
var s3StorageLocationParamsAlterComplete = &S3StorageLocationParams{
	Name:                     "s3_alter_complete",
	StorageProvider:          S3StorageProviderS3,
	StorageAwsRoleArn:        "arn:aws:iam::123456789012:role/myrole",
	StorageBaseUrl:           "s3://my-bucket/alter-path",
	StorageAwsExternalId:     String("external_id_alter"),
	StorageAwsAccessPointArn: String("arn:aws:s3:us-west-2:123456789012:accesspoint/my-access-point"),
	UsePrivatelinkEndpoint:   Bool(true),
	Encryption: &ExternalVolumeS3Encryption{
		EncryptionType: S3EncryptionTypeSseKms,
		KmsKeyId:       String("1234abcd-alter-kms-key"),
	},
}

var gcsStorageLocationParamsAlterComplete = &GCSStorageLocationParams{
	Name:           "gcs_alter_complete",
	StorageBaseUrl: "gcs://my-bucket/alter-path",
	Encryption: &ExternalVolumeGCSEncryption{
		EncryptionType: GCSEncryptionTypeSseKms,
		KmsKeyId:       String("1234abcd-alter-kms-key"),
	},
}

var azureStorageLocationParamsAlterComplete = &AzureStorageLocationParams{
	Name:                   "azure_alter_complete",
	AzureTenantId:          "a123b4cd-1abc-12ab-12ab-1a2b34c5d678",
	StorageBaseUrl:         "azure://myaccount.blob.core.windows.net/mycontainer/alter-path",
	UsePrivatelinkEndpoint: Bool(true),
}

var s3CompatStorageLocationParamsAlterComplete = &S3CompatStorageLocationParams{
	Name:            "s3compat_alter_complete",
	StorageBaseUrl:  "s3compat://my-bucket/alter-path",
	StorageEndpoint: "https://s3-compatible.example.com",
	Credentials: &ExternalVolumeS3CompatCredentials{
		AwsKeyId:     "AKIAIOSFODNN7EXAMPLE",
		AwsSecretKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
	},
}

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
		opts.StorageLocations = []ExternalVolumeStorageLocation{
			{S3StorageLocationParams: s3StorageLocationParamsBasic},
			{},
			{
				S3StorageLocationParams:    s3StorageLocationParamsBasic,
				GCSStorageLocationParams:   gcsStorageLocationParamsBasic,
				AzureStorageLocationParams: azureStorageLocationParamsBasic,
			},
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
		opts.StorageLocations = []ExternalVolumeStorageLocation{}
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("CreateExternalVolumeOptions", "StorageLocations"))
	})

	t.Run("basic - s3", func(t *testing.T) {
		opts := defaultOpts()
		opts.StorageLocations = []ExternalVolumeStorageLocation{{S3StorageLocationParams: s3StorageLocationParamsBasic}}
		assertOptsValidAndSQLEquals(t, opts, `CREATE EXTERNAL VOLUME %s STORAGE_LOCATIONS = ((NAME = 's3_basic' STORAGE_PROVIDER = 'S3' STORAGE_AWS_ROLE_ARN = 'arn:aws:iam::123456789012:role/myrole' STORAGE_BASE_URL = 's3://my-bucket/path'))`, id.FullyQualifiedName())
	})

	t.Run("basic - gcs", func(t *testing.T) {
		opts := defaultOpts()
		opts.StorageLocations = []ExternalVolumeStorageLocation{{GCSStorageLocationParams: gcsStorageLocationParamsBasic}}
		assertOptsValidAndSQLEquals(t, opts, `CREATE EXTERNAL VOLUME %s STORAGE_LOCATIONS = ((NAME = 'gcs_basic' STORAGE_PROVIDER = 'GCS' STORAGE_BASE_URL = 'gcs://my-bucket/path'))`, id.FullyQualifiedName())
	})

	t.Run("basic - azure", func(t *testing.T) {
		opts := defaultOpts()
		opts.StorageLocations = []ExternalVolumeStorageLocation{{AzureStorageLocationParams: azureStorageLocationParamsBasic}}
		assertOptsValidAndSQLEquals(t, opts, `CREATE EXTERNAL VOLUME %s STORAGE_LOCATIONS = ((NAME = 'azure_basic' STORAGE_PROVIDER = 'AZURE' AZURE_TENANT_ID = 'a123b4cd-1abc-12ab-12ab-1a2b34c5d678' STORAGE_BASE_URL = 'azure://myaccount.blob.core.windows.net/mycontainer/path'))`, id.FullyQualifiedName())
	})

	t.Run("basic - s3compat", func(t *testing.T) {
		opts := defaultOpts()
		opts.StorageLocations = []ExternalVolumeStorageLocation{{S3CompatStorageLocationParams: s3CompatStorageLocationParamsBasic}}
		assertOptsValidAndSQLEquals(t, opts, `CREATE EXTERNAL VOLUME %s STORAGE_LOCATIONS = ((NAME = 's3compat_basic' STORAGE_PROVIDER = 'S3COMPAT' STORAGE_BASE_URL = 's3compat://my-bucket/path' STORAGE_ENDPOINT = 'https://s3-compatible.example.com'))`, id.FullyQualifiedName())
	})

	t.Run("all options - s3", func(t *testing.T) {
		opts := defaultOpts()
		opts.StorageLocations = []ExternalVolumeStorageLocation{{S3StorageLocationParams: s3StorageLocationParamsComplete}}
		assertOptsValidAndSQLEquals(t, opts, `CREATE EXTERNAL VOLUME %s STORAGE_LOCATIONS = ((NAME = 's3_complete' STORAGE_PROVIDER = 'S3' STORAGE_AWS_ROLE_ARN = 'arn:aws:iam::123456789012:role/myrole' STORAGE_BASE_URL = 's3://my-bucket/path' STORAGE_AWS_EXTERNAL_ID = 'external_id_123' STORAGE_AWS_ACCESS_POINT_ARN = 'arn:aws:s3:us-west-2:123456789012:accesspoint/my-access-point' USE_PRIVATELINK_ENDPOINT = true ENCRYPTION = (TYPE = 'AWS_SSE_KMS' KMS_KEY_ID = '1234abcd-12ab-34cd-56ef-1234567890ab')))`, id.FullyQualifiedName())
	})

	t.Run("all options - gcs", func(t *testing.T) {
		opts := defaultOpts()
		opts.StorageLocations = []ExternalVolumeStorageLocation{{GCSStorageLocationParams: gcsStorageLocationParamsComplete}}
		assertOptsValidAndSQLEquals(t, opts, `CREATE EXTERNAL VOLUME %s STORAGE_LOCATIONS = ((NAME = 'gcs_complete' STORAGE_PROVIDER = 'GCS' STORAGE_BASE_URL = 'gcs://my-bucket/path' ENCRYPTION = (TYPE = 'GCS_SSE_KMS' KMS_KEY_ID = '1234abcd-12ab-34cd-56ef-1234567890ab')))`, id.FullyQualifiedName())
	})

	t.Run("all options - azure", func(t *testing.T) {
		opts := defaultOpts()
		opts.StorageLocations = []ExternalVolumeStorageLocation{{AzureStorageLocationParams: azureStorageLocationParamsComplete}}
		assertOptsValidAndSQLEquals(t, opts, `CREATE EXTERNAL VOLUME %s STORAGE_LOCATIONS = ((NAME = 'azure_complete' STORAGE_PROVIDER = 'AZURE' AZURE_TENANT_ID = 'a123b4cd-1abc-12ab-12ab-1a2b34c5d678' STORAGE_BASE_URL = 'azure://myaccount.blob.core.windows.net/mycontainer/path' USE_PRIVATELINK_ENDPOINT = true))`, id.FullyQualifiedName())
	})

	t.Run("all options - s3compat", func(t *testing.T) {
		opts := defaultOpts()
		opts.StorageLocations = []ExternalVolumeStorageLocation{{S3CompatStorageLocationParams: s3CompatStorageLocationParamsComplete}}
		assertOptsValidAndSQLEquals(t, opts, `CREATE EXTERNAL VOLUME %s STORAGE_LOCATIONS = ((NAME = 's3compat_complete' STORAGE_PROVIDER = 'S3COMPAT' STORAGE_BASE_URL = 's3compat://my-bucket/path' STORAGE_ENDPOINT = 'https://s3-compatible.example.com' CREDENTIALS = (AWS_KEY_ID = 'AKIAIOSFODNN7EXAMPLE' AWS_SECRET_KEY = 'wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY')))`, id.FullyQualifiedName())
	})

	t.Run("all storage location types with volume options", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.StorageLocations = []ExternalVolumeStorageLocation{
			{S3StorageLocationParams: s3StorageLocationParamsComplete},
			{GCSStorageLocationParams: gcsStorageLocationParamsComplete},
			{AzureStorageLocationParams: azureStorageLocationParamsComplete},
			{S3CompatStorageLocationParams: s3CompatStorageLocationParamsComplete},
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

	t.Run("validation: exactly one field from [opts.RemoveStorageLocation opts.Set opts.AddStorageLocation] should be present - zero set", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterExternalVolumeOptions", "RemoveStorageLocation", "Set", "AddStorageLocation"))
	})

	t.Run("validation: exactly one field from [opts.RemoveStorageLocation opts.Set opts.AddStorageLocation] should be present - two set", func(t *testing.T) {
		removeAndSetOpts := defaultOpts()
		removeAndAddOpts := defaultOpts()
		setAndAddOpts := defaultOpts()

		removeAndSetOpts.RemoveStorageLocation = String("some storage location")
		removeAndSetOpts.Set = &AlterExternalVolumeSet{AllowWrites: Bool(true)}

		removeAndAddOpts.RemoveStorageLocation = String("some storage location")
		removeAndAddOpts.AddStorageLocation = &ExternalVolumeStorageLocation{S3StorageLocationParams: s3StorageLocationParamsComplete}

		setAndAddOpts.Set = &AlterExternalVolumeSet{AllowWrites: Bool(true)}
		setAndAddOpts.AddStorageLocation = &ExternalVolumeStorageLocation{S3StorageLocationParams: s3StorageLocationParamsComplete}

		assertOptsInvalidJoinedErrors(t, removeAndSetOpts, errExactlyOneOf("AlterExternalVolumeOptions", "RemoveStorageLocation", "Set", "AddStorageLocation"))
		assertOptsInvalidJoinedErrors(t, removeAndAddOpts, errExactlyOneOf("AlterExternalVolumeOptions", "RemoveStorageLocation", "Set", "AddStorageLocation"))
		assertOptsInvalidJoinedErrors(t, setAndAddOpts, errExactlyOneOf("AlterExternalVolumeOptions", "RemoveStorageLocation", "Set", "AddStorageLocation"))
	})

	t.Run("validation: exactly one field from [opts.RemoveStorageLocation opts.Set opts.AddStorageLocation] should be present - three set", func(t *testing.T) {
		opts := defaultOpts()
		opts.RemoveStorageLocation = String("some storage location")
		opts.Set = &AlterExternalVolumeSet{AllowWrites: Bool(true)}
		opts.AddStorageLocation = &ExternalVolumeStorageLocation{S3StorageLocationParams: s3StorageLocationParamsComplete}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterExternalVolumeOptions", "RemoveStorageLocation", "Set", "AddStorageLocation"))
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptyAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: exactly one field from [opts.AddStorageLocation.S3StorageLocationParams opts.AddStorageLocation.GCSStorageLocationParams opts.AddStorageLocation.AzureStorageLocationParams opts.AddStorageLocation.S3CompatStorageLocationParams] should be present - none set", func(t *testing.T) {
		opts := defaultOpts()
		opts.AddStorageLocation = &ExternalVolumeStorageLocation{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterExternalVolumeOptions.AddStorageLocation", "S3StorageLocationParams", "GCSStorageLocationParams", "AzureStorageLocationParams", "S3CompatStorageLocationParams"))
	})

	t.Run("validation: exactly one field from [opts.AddStorageLocation.S3StorageLocationParams opts.AddStorageLocation.GCSStorageLocationParams opts.AddStorageLocation.AzureStorageLocationParams opts.AddStorageLocation.S3CompatStorageLocationParams] should be present - two set", func(t *testing.T) {
		s3AndGcsOpts := defaultOpts()
		s3AndAzureOpts := defaultOpts()
		gcsAndAzureOpts := defaultOpts()
		s3AndGcsOpts.AddStorageLocation = &ExternalVolumeStorageLocation{
			S3StorageLocationParams:  s3StorageLocationParamsBasic,
			GCSStorageLocationParams: gcsStorageLocationParamsBasic,
		}
		s3AndAzureOpts.AddStorageLocation = &ExternalVolumeStorageLocation{
			S3StorageLocationParams:    s3StorageLocationParamsBasic,
			AzureStorageLocationParams: azureStorageLocationParamsBasic,
		}
		gcsAndAzureOpts.AddStorageLocation = &ExternalVolumeStorageLocation{
			GCSStorageLocationParams:   gcsStorageLocationParamsBasic,
			AzureStorageLocationParams: azureStorageLocationParamsBasic,
		}
		assertOptsInvalidJoinedErrors(t, s3AndGcsOpts, errExactlyOneOf("AlterExternalVolumeOptions.AddStorageLocation", "S3StorageLocationParams", "GCSStorageLocationParams", "AzureStorageLocationParams", "S3CompatStorageLocationParams"))
		assertOptsInvalidJoinedErrors(t, s3AndAzureOpts, errExactlyOneOf("AlterExternalVolumeOptions.AddStorageLocation", "S3StorageLocationParams", "GCSStorageLocationParams", "AzureStorageLocationParams", "S3CompatStorageLocationParams"))
		assertOptsInvalidJoinedErrors(t, gcsAndAzureOpts, errExactlyOneOf("AlterExternalVolumeOptions.AddStorageLocation", "S3StorageLocationParams", "GCSStorageLocationParams", "AzureStorageLocationParams", "S3CompatStorageLocationParams"))
	})

	t.Run("validation: exactly one field from [opts.AddStorageLocation.S3StorageLocationParams opts.AddStorageLocation.GCSStorageLocationParams opts.AddStorageLocation.AzureStorageLocationParams opts.AddStorageLocation.S3CompatStorageLocationParams] should be present - three set", func(t *testing.T) {
		opts := defaultOpts()
		opts.AddStorageLocation = &ExternalVolumeStorageLocation{
			S3StorageLocationParams:    s3StorageLocationParamsBasic,
			GCSStorageLocationParams:   gcsStorageLocationParamsBasic,
			AzureStorageLocationParams: azureStorageLocationParamsBasic,
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterExternalVolumeOptions.AddStorageLocation", "S3StorageLocationParams", "GCSStorageLocationParams", "AzureStorageLocationParams", "S3CompatStorageLocationParams"))
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
		opts.AddStorageLocation = &ExternalVolumeStorageLocation{S3StorageLocationParams: s3StorageLocationParamsAlterBasic}
		assertOptsValidAndSQLEquals(t, opts, `ALTER EXTERNAL VOLUME %s ADD STORAGE_LOCATION = (NAME = 's3_alter_basic' STORAGE_PROVIDER = 'S3' STORAGE_AWS_ROLE_ARN = 'arn:aws:iam::123456789012:role/myrole' STORAGE_BASE_URL = 's3://my-bucket/alter-path')`, id.FullyQualifiedName())
	})

	t.Run("add storage location - basic - gcs", func(t *testing.T) {
		opts := defaultOpts()
		opts.AddStorageLocation = &ExternalVolumeStorageLocation{GCSStorageLocationParams: gcsStorageLocationParamsAlterBasic}
		assertOptsValidAndSQLEquals(t, opts, `ALTER EXTERNAL VOLUME %s ADD STORAGE_LOCATION = (NAME = 'gcs_alter_basic' STORAGE_PROVIDER = 'GCS' STORAGE_BASE_URL = 'gcs://my-bucket/alter-path')`, id.FullyQualifiedName())
	})

	t.Run("add storage location - basic - azure", func(t *testing.T) {
		opts := defaultOpts()
		opts.AddStorageLocation = &ExternalVolumeStorageLocation{AzureStorageLocationParams: azureStorageLocationParamsAlterBasic}
		assertOptsValidAndSQLEquals(t, opts, `ALTER EXTERNAL VOLUME %s ADD STORAGE_LOCATION = (NAME = 'azure_alter_basic' STORAGE_PROVIDER = 'AZURE' AZURE_TENANT_ID = 'a123b4cd-1abc-12ab-12ab-1a2b34c5d678' STORAGE_BASE_URL = 'azure://myaccount.blob.core.windows.net/mycontainer/alter-path')`, id.FullyQualifiedName())
	})

	t.Run("add storage location - basic - s3compat", func(t *testing.T) {
		opts := defaultOpts()
		opts.AddStorageLocation = &ExternalVolumeStorageLocation{S3CompatStorageLocationParams: s3CompatStorageLocationParamsAlterBasic}
		assertOptsValidAndSQLEquals(t, opts, `ALTER EXTERNAL VOLUME %s ADD STORAGE_LOCATION = (NAME = 's3compat_alter_basic' STORAGE_PROVIDER = 'S3COMPAT' STORAGE_BASE_URL = 's3compat://my-bucket/alter-path' STORAGE_ENDPOINT = 'https://s3-compatible.example.com')`, id.FullyQualifiedName())
	})

	t.Run("add storage location - all options - s3", func(t *testing.T) {
		opts := defaultOpts()
		opts.AddStorageLocation = &ExternalVolumeStorageLocation{S3StorageLocationParams: s3StorageLocationParamsAlterComplete}
		assertOptsValidAndSQLEquals(t, opts, `ALTER EXTERNAL VOLUME %s ADD STORAGE_LOCATION = (NAME = 's3_alter_complete' STORAGE_PROVIDER = 'S3' STORAGE_AWS_ROLE_ARN = 'arn:aws:iam::123456789012:role/myrole' STORAGE_BASE_URL = 's3://my-bucket/alter-path' STORAGE_AWS_EXTERNAL_ID = 'external_id_alter' STORAGE_AWS_ACCESS_POINT_ARN = 'arn:aws:s3:us-west-2:123456789012:accesspoint/my-access-point' USE_PRIVATELINK_ENDPOINT = true ENCRYPTION = (TYPE = 'AWS_SSE_KMS' KMS_KEY_ID = '1234abcd-alter-kms-key'))`, id.FullyQualifiedName())
	})

	t.Run("add storage location - all options - gcs", func(t *testing.T) {
		opts := defaultOpts()
		opts.AddStorageLocation = &ExternalVolumeStorageLocation{GCSStorageLocationParams: gcsStorageLocationParamsAlterComplete}
		assertOptsValidAndSQLEquals(t, opts, `ALTER EXTERNAL VOLUME %s ADD STORAGE_LOCATION = (NAME = 'gcs_alter_complete' STORAGE_PROVIDER = 'GCS' STORAGE_BASE_URL = 'gcs://my-bucket/alter-path' ENCRYPTION = (TYPE = 'GCS_SSE_KMS' KMS_KEY_ID = '1234abcd-alter-kms-key'))`, id.FullyQualifiedName())
	})

	t.Run("add storage location - all options - azure", func(t *testing.T) {
		opts := defaultOpts()
		opts.AddStorageLocation = &ExternalVolumeStorageLocation{AzureStorageLocationParams: azureStorageLocationParamsAlterComplete}
		assertOptsValidAndSQLEquals(t, opts, `ALTER EXTERNAL VOLUME %s ADD STORAGE_LOCATION = (NAME = 'azure_alter_complete' STORAGE_PROVIDER = 'AZURE' AZURE_TENANT_ID = 'a123b4cd-1abc-12ab-12ab-1a2b34c5d678' STORAGE_BASE_URL = 'azure://myaccount.blob.core.windows.net/mycontainer/alter-path' USE_PRIVATELINK_ENDPOINT = true)`, id.FullyQualifiedName())
	})

	t.Run("add storage location - all options - s3compat", func(t *testing.T) {
		opts := defaultOpts()
		opts.AddStorageLocation = &ExternalVolumeStorageLocation{S3CompatStorageLocationParams: s3CompatStorageLocationParamsAlterComplete}
		assertOptsValidAndSQLEquals(t, opts, `ALTER EXTERNAL VOLUME %s ADD STORAGE_LOCATION = (NAME = 's3compat_alter_complete' STORAGE_PROVIDER = 'S3COMPAT' STORAGE_BASE_URL = 's3compat://my-bucket/alter-path' STORAGE_ENDPOINT = 'https://s3-compatible.example.com' CREDENTIALS = (AWS_KEY_ID = 'AKIAIOSFODNN7EXAMPLE' AWS_SECRET_KEY = 'wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY'))`, id.FullyQualifiedName())
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
		{input: "s3compat", want: StorageProviderS3COMPAT},
		{input: "S3COMPAT", want: StorageProviderS3COMPAT},
	}

	invalid := []test{
		{input: ""},
		{input: "foo"},
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

func Test_GetStorageLocationName(t *testing.T) {
	s3StorageLocationName := "s3Test"
	s3StorageBaseUrl := "s3://my_example_bucket"
	s3StorageAwsRoleArn := "arn:aws:iam::123456789012:role/myrole"
	s3EncryptionKmsKeyId := "123456789"

	gcsStorageLocationName := "gcsTest"
	gcsStorageBaseUrl := "gcs://my_example_bucket"
	gcsEncryptionKmsKeyId := "123456789"

	azureStorageLocationName := "azureTest"
	azureStorageBaseUrl := "azure://123456789.blob.core.windows.net/my_example_container"
	azureTenantId := "123456789"

	s3StorageLocationA := S3StorageLocationParams{
		Name:                 s3StorageLocationName,
		StorageProvider:      S3StorageProviderS3,
		StorageBaseUrl:       s3StorageBaseUrl,
		StorageAwsRoleArn:    s3StorageAwsRoleArn,
		StorageAwsExternalId: &s3StorageAwsExternalId,
		Encryption: &ExternalVolumeS3Encryption{
			EncryptionType: S3EncryptionTypeSseKms,
			KmsKeyId:       &s3EncryptionKmsKeyId,
		},
	}

	azureStorageLocationA := AzureStorageLocationParams{
		Name:           azureStorageLocationName,
		StorageBaseUrl: azureStorageBaseUrl,
		AzureTenantId:  azureTenantId,
	}

	gcsStorageLocationA := GCSStorageLocationParams{
		Name:           gcsStorageLocationName,
		StorageBaseUrl: gcsStorageBaseUrl,
		Encryption: &ExternalVolumeGCSEncryption{
			EncryptionType: GCSEncryptionTypeSseKms,
			KmsKeyId:       &gcsEncryptionKmsKeyId,
		},
	}

	s3GovStorageLocationA := S3StorageLocationParams{
		Name:              s3StorageLocationName,
		StorageProvider:   S3StorageProviderS3GOV,
		StorageBaseUrl:    s3StorageBaseUrl,
		StorageAwsRoleArn: s3StorageAwsRoleArn,
	}

	s3CompatStorageLocationA := S3CompatStorageLocationParams{
		Name:            "s3compatTest",
		StorageBaseUrl:  "s3compat://my-bucket/my-path",
		StorageEndpoint: "https://s3-compatible.example.com",
	}

	testCases := []struct {
		Name            string
		StorageLocation ExternalVolumeStorageLocation
		ExpectedName    string
	}{
		{
			Name:            "S3 storage location name successfully read",
			StorageLocation: ExternalVolumeStorageLocation{S3StorageLocationParams: &s3StorageLocationA},
			ExpectedName:    s3StorageLocationA.Name,
		},
		{
			Name:            "S3GOV storage location name successfully read",
			StorageLocation: ExternalVolumeStorageLocation{S3StorageLocationParams: &s3GovStorageLocationA},
			ExpectedName:    s3GovStorageLocationA.Name,
		},
		{
			Name:            "GCS storage location name successfully read",
			StorageLocation: ExternalVolumeStorageLocation{GCSStorageLocationParams: &gcsStorageLocationA},
			ExpectedName:    gcsStorageLocationA.Name,
		},
		{
			Name:            "Azure storage location name successfully read",
			StorageLocation: ExternalVolumeStorageLocation{AzureStorageLocationParams: &azureStorageLocationA},
			ExpectedName:    azureStorageLocationA.Name,
		},
		{
			Name:            "S3COMPAT storage location name successfully read",
			StorageLocation: ExternalVolumeStorageLocation{S3CompatStorageLocationParams: &s3CompatStorageLocationA},
			ExpectedName:    s3CompatStorageLocationA.Name,
		},
	}

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
			Name:            "Empty S3COMPAT storage location",
			StorageLocation: ExternalVolumeStorageLocation{S3CompatStorageLocationParams: &S3CompatStorageLocationParams{}},
		},
		{
			Name:            "Empty storage location",
			StorageLocation: ExternalVolumeStorageLocation{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			name, err := GetStorageLocationName(tc.StorageLocation)
			require.NoError(t, err)
			assert.Equal(t, tc.ExpectedName, name)
		})
	}
	for _, tc := range invalidTestCases {
		t.Run(tc.Name, func(t *testing.T) {
			_, err := GetStorageLocationName(tc.StorageLocation)
			require.Error(t, err)
		})
	}
}

func Test_GetStorageLocationStorageProvider(t *testing.T) {
	s3StorageLocationName := "s3Test"
	s3StorageBaseUrl := "s3://my_example_bucket"
	s3StorageAwsRoleArn := "arn:aws:iam::123456789012:role/myrole"
	s3EncryptionKmsKeyId := "123456789"

	gcsStorageLocationName := "gcsTest"
	gcsStorageBaseUrl := "gcs://my_example_bucket"
	gcsEncryptionKmsKeyId := "123456789"

	azureStorageLocationName := "azureTest"
	azureStorageBaseUrl := "azure://123456789.blob.core.windows.net/my_example_container"
	azureTenantId := "123456789"

	s3StorageLocationA := S3StorageLocationParams{
		Name:                 s3StorageLocationName,
		StorageProvider:      S3StorageProviderS3,
		StorageBaseUrl:       s3StorageBaseUrl,
		StorageAwsRoleArn:    s3StorageAwsRoleArn,
		StorageAwsExternalId: &s3StorageAwsExternalId,
		Encryption: &ExternalVolumeS3Encryption{
			EncryptionType: S3EncryptionTypeSseKms,
			KmsKeyId:       &s3EncryptionKmsKeyId,
		},
	}

	azureStorageLocationA := AzureStorageLocationParams{
		Name:           azureStorageLocationName,
		StorageBaseUrl: azureStorageBaseUrl,
		AzureTenantId:  azureTenantId,
	}

	gcsStorageLocationA := GCSStorageLocationParams{
		Name:           gcsStorageLocationName,
		StorageBaseUrl: gcsStorageBaseUrl,
		Encryption: &ExternalVolumeGCSEncryption{
			EncryptionType: GCSEncryptionTypeSseKms,
			KmsKeyId:       &gcsEncryptionKmsKeyId,
		},
	}

	s3GovStorageLocationA := S3StorageLocationParams{
		Name:              s3StorageLocationName,
		StorageProvider:   S3StorageProviderS3GOV,
		StorageBaseUrl:    s3StorageBaseUrl,
		StorageAwsRoleArn: s3StorageAwsRoleArn,
	}

	s3CompatStorageLocationA := S3CompatStorageLocationParams{
		Name:            "s3compatTest",
		StorageBaseUrl:  "s3compat://my-bucket/my-path",
		StorageEndpoint: "https://s3-compatible.example.com",
	}

	testCases := []struct {
		Name                    string
		StorageLocation         ExternalVolumeStorageLocation
		ExpectedStorageProvider StorageProvider
	}{
		{
			Name:                    "S3 storage provider",
			StorageLocation:         ExternalVolumeStorageLocation{S3StorageLocationParams: &s3StorageLocationA},
			ExpectedStorageProvider: StorageProviderS3,
		},
		{
			Name:                    "S3GOV storage provider",
			StorageLocation:         ExternalVolumeStorageLocation{S3StorageLocationParams: &s3GovStorageLocationA},
			ExpectedStorageProvider: StorageProviderS3GOV,
		},
		{
			Name:                    "GCS storage provider",
			StorageLocation:         ExternalVolumeStorageLocation{GCSStorageLocationParams: &gcsStorageLocationA},
			ExpectedStorageProvider: StorageProviderGCS,
		},
		{
			Name:                    "Azure storage provider",
			StorageLocation:         ExternalVolumeStorageLocation{AzureStorageLocationParams: &azureStorageLocationA},
			ExpectedStorageProvider: StorageProviderAzure,
		},
		{
			Name:                    "S3COMPAT storage provider",
			StorageLocation:         ExternalVolumeStorageLocation{S3CompatStorageLocationParams: &s3CompatStorageLocationA},
			ExpectedStorageProvider: StorageProviderS3COMPAT,
		},
	}

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
			Name:            "Empty S3COMPAT storage location",
			StorageLocation: ExternalVolumeStorageLocation{S3CompatStorageLocationParams: &S3CompatStorageLocationParams{}},
		},
		{
			Name:            "Empty storage location",
			StorageLocation: ExternalVolumeStorageLocation{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			storageProvider, err := GetStorageLocationStorageProvider(tc.StorageLocation)
			require.NoError(t, err)
			assert.Equal(t, tc.ExpectedStorageProvider, storageProvider)
		})
	}
	for _, tc := range invalidTestCases {
		t.Run(tc.Name, func(t *testing.T) {
			_, err := GetStorageLocationName(tc.StorageLocation)
			require.Error(t, err)
		})
	}
}

var s3StorageAwsExternalId = "1234567890"

func Test_CopySentinelStorageLocation(t *testing.T) {
	tempStorageLocationName := "terraform_provider_sentinel_storage_location"
	s3StorageLocationName := "s3Test"
	s3StorageBaseUrl := "s3://my_example_bucket"
	s3StorageAwsRoleArn := "arn:aws:iam::123456789012:role/myrole"
	s3EncryptionKmsKeyId := "123456789"

	gcsStorageLocationName := "gcsTest"
	gcsStorageBaseUrl := "gcs://my_example_bucket"
	gcsEncryptionKmsKeyId := "123456789"

	azureStorageLocationName := "azureTest"
	azureStorageBaseUrl := "azure://123456789.blob.core.windows.net/my_example_container"
	azureTenantId := "123456789"

	s3StorageLocationA := S3StorageLocationParams{
		Name:                 s3StorageLocationName,
		StorageProvider:      S3StorageProviderS3,
		StorageBaseUrl:       s3StorageBaseUrl,
		StorageAwsRoleArn:    s3StorageAwsRoleArn,
		StorageAwsExternalId: &s3StorageAwsExternalId,
		Encryption: &ExternalVolumeS3Encryption{
			EncryptionType: S3EncryptionTypeSseKms,
			KmsKeyId:       &s3EncryptionKmsKeyId,
		},
	}

	azureStorageLocationA := AzureStorageLocationParams{
		Name:           azureStorageLocationName,
		StorageBaseUrl: azureStorageBaseUrl,
		AzureTenantId:  azureTenantId,
	}

	gcsStorageLocationA := GCSStorageLocationParams{
		Name:           gcsStorageLocationName,
		StorageBaseUrl: gcsStorageBaseUrl,
		Encryption: &ExternalVolumeGCSEncryption{
			EncryptionType: GCSEncryptionTypeSseKms,
			KmsKeyId:       &gcsEncryptionKmsKeyId,
		},
	}

	t.Run("S3 storage location", func(t *testing.T) {
		storageLocationInput := ExternalVolumeStorageLocation{S3StorageLocationParams: &s3StorageLocationA}
		copiedStorageLocation, err := CopySentinelStorageLocation(storageLocationInput)
		require.NoError(t, err)
		assert.Equal(t, copiedStorageLocation.S3StorageLocationParams.Name, tempStorageLocationName)
		assert.Equal(t, copiedStorageLocation.S3StorageLocationParams.StorageProvider, s3StorageLocationA.StorageProvider)
		assert.Equal(t, copiedStorageLocation.S3StorageLocationParams.StorageBaseUrl, s3StorageLocationA.StorageBaseUrl)
		assert.Equal(t, copiedStorageLocation.S3StorageLocationParams.StorageAwsRoleArn, s3StorageLocationA.StorageAwsRoleArn)
		assert.Equal(t, copiedStorageLocation.S3StorageLocationParams.StorageAwsExternalId, s3StorageLocationA.StorageAwsExternalId)
		assert.Equal(t, copiedStorageLocation.S3StorageLocationParams.Encryption.EncryptionType, s3StorageLocationA.Encryption.EncryptionType)
		assert.Equal(t, *copiedStorageLocation.S3StorageLocationParams.Encryption.KmsKeyId, *s3StorageLocationA.Encryption.KmsKeyId)
	})

	t.Run("GCS storage location", func(t *testing.T) {
		storageLocationInput := ExternalVolumeStorageLocation{GCSStorageLocationParams: &gcsStorageLocationA}
		copiedStorageLocation, err := CopySentinelStorageLocation(storageLocationInput)
		require.NoError(t, err)
		assert.Equal(t, copiedStorageLocation.GCSStorageLocationParams.Name, tempStorageLocationName)
		assert.Equal(t, copiedStorageLocation.GCSStorageLocationParams.StorageBaseUrl, gcsStorageLocationA.StorageBaseUrl)
		assert.Equal(t, copiedStorageLocation.GCSStorageLocationParams.Encryption.EncryptionType, gcsStorageLocationA.Encryption.EncryptionType)
		assert.Equal(t, *copiedStorageLocation.GCSStorageLocationParams.Encryption.KmsKeyId, *gcsStorageLocationA.Encryption.KmsKeyId)
	})

	t.Run("Azure storage location", func(t *testing.T) {
		storageLocationInput := ExternalVolumeStorageLocation{AzureStorageLocationParams: &azureStorageLocationA}
		copiedStorageLocation, err := CopySentinelStorageLocation(storageLocationInput)
		require.NoError(t, err)
		assert.Equal(t, copiedStorageLocation.AzureStorageLocationParams.Name, tempStorageLocationName)
		assert.Equal(t, copiedStorageLocation.AzureStorageLocationParams.StorageBaseUrl, azureStorageLocationA.StorageBaseUrl)
		assert.Equal(t, copiedStorageLocation.AzureStorageLocationParams.AzureTenantId, azureStorageLocationA.AzureTenantId)
	})

	s3CompatStorageLocationA := S3CompatStorageLocationParams{
		Name:            "s3compatTest",
		StorageBaseUrl:  "s3compat://my-bucket/my-path",
		StorageEndpoint: "https://s3-compatible.example.com",
		Credentials: &ExternalVolumeS3CompatCredentials{
			AwsKeyId:     "some_key_id",
			AwsSecretKey: "some_secret_key",
		},
	}

	t.Run("S3COMPAT storage location", func(t *testing.T) {
		storageLocationInput := ExternalVolumeStorageLocation{S3CompatStorageLocationParams: &s3CompatStorageLocationA}
		copiedStorageLocation, err := CopySentinelStorageLocation(storageLocationInput)
		require.NoError(t, err)
		assert.Equal(t, copiedStorageLocation.S3CompatStorageLocationParams.Name, tempStorageLocationName)
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
			Name:            "Empty S3COMPAT storage location",
			StorageLocation: ExternalVolumeStorageLocation{S3CompatStorageLocationParams: &S3CompatStorageLocationParams{}},
		},
		{
			Name:            "Empty storage location",
			StorageLocation: ExternalVolumeStorageLocation{},
		},
	}

	for _, tc := range invalidTestCases {
		t.Run(tc.Name, func(t *testing.T) {
			_, err := CopySentinelStorageLocation(tc.StorageLocation)
			require.Error(t, err)
		})
	}
}

func Test_CommonPrefixLastIndex(t *testing.T) {
	s3StorageLocationName := "s3Test"
	s3StorageLocationName2 := "s3Test2"
	s3StorageBaseUrl := "s3://my_example_bucket"
	s3StorageAwsRoleArn := "arn:aws:iam::123456789012:role/myrole"
	s3EncryptionKmsKeyId := "123456789"

	gcsStorageLocationName := "gcsTest"
	gcsStorageLocationName2 := "gcsTest2"
	gcsStorageBaseUrl := "gcs://my_example_bucket"
	gcsEncryptionKmsKeyId := "123456789"

	azureStorageLocationName := "azureTest"
	azureStorageLocationName2 := "azureTest2"
	azureStorageBaseUrl := "azure://123456789.blob.core.windows.net/my_example_container"
	azureTenantId := "123456789"

	s3StorageLocationA := S3StorageLocationParams{
		Name:                 s3StorageLocationName,
		StorageProvider:      S3StorageProviderS3,
		StorageBaseUrl:       s3StorageBaseUrl,
		StorageAwsRoleArn:    s3StorageAwsRoleArn,
		StorageAwsExternalId: &s3StorageAwsExternalId,
		Encryption: &ExternalVolumeS3Encryption{
			EncryptionType: S3EncryptionTypeSseKms,
			KmsKeyId:       &s3EncryptionKmsKeyId,
		},
	}

	s3StorageLocationB := S3StorageLocationParams{
		Name:                 s3StorageLocationName2,
		StorageProvider:      S3StorageProviderS3,
		StorageBaseUrl:       s3StorageBaseUrl,
		StorageAwsRoleArn:    s3StorageAwsRoleArn,
		StorageAwsExternalId: &s3StorageAwsExternalId,
		Encryption: &ExternalVolumeS3Encryption{
			EncryptionType: S3EncryptionTypeSseKms,
			KmsKeyId:       &s3EncryptionKmsKeyId,
		},
	}

	azureStorageLocationA := AzureStorageLocationParams{
		Name:           azureStorageLocationName,
		StorageBaseUrl: azureStorageBaseUrl,
		AzureTenantId:  azureTenantId,
	}

	azureStorageLocationB := AzureStorageLocationParams{
		Name:           azureStorageLocationName2,
		StorageBaseUrl: azureStorageBaseUrl,
		AzureTenantId:  azureTenantId,
	}

	gcsStorageLocationA := GCSStorageLocationParams{
		Name:           gcsStorageLocationName,
		StorageBaseUrl: gcsStorageBaseUrl,
		Encryption: &ExternalVolumeGCSEncryption{
			EncryptionType: GCSEncryptionTypeSseKms,
			KmsKeyId:       &gcsEncryptionKmsKeyId,
		},
	}

	gcsStorageLocationB := GCSStorageLocationParams{
		Name:           gcsStorageLocationName2,
		StorageBaseUrl: gcsStorageBaseUrl,
		Encryption: &ExternalVolumeGCSEncryption{
			EncryptionType: GCSEncryptionTypeSseKms,
			KmsKeyId:       &gcsEncryptionKmsKeyId,
		},
	}

	gcsStorageLocationC := GCSStorageLocationParams{
		Name:           "test",
		StorageBaseUrl: gcsStorageBaseUrl,
		Encryption: &ExternalVolumeGCSEncryption{
			EncryptionType: GCSEncryptionTypeSseKms,
			KmsKeyId:       &gcsEncryptionKmsKeyId,
		},
	}

	s3GovStorageLocationA := S3StorageLocationParams{
		Name:              s3StorageLocationName,
		StorageProvider:   S3StorageProviderS3GOV,
		StorageBaseUrl:    s3StorageBaseUrl,
		StorageAwsRoleArn: s3StorageAwsRoleArn,
	}

	testCases := []struct {
		Name           string
		ListA          []ExternalVolumeStorageLocation
		ListB          []ExternalVolumeStorageLocation
		ExpectedOutput int
	}{
		{
			Name:           "Two empty lists",
			ListA:          []ExternalVolumeStorageLocation{},
			ListB:          []ExternalVolumeStorageLocation{},
			ExpectedOutput: -1,
		},
		{
			Name:           "First list empty",
			ListA:          []ExternalVolumeStorageLocation{},
			ListB:          []ExternalVolumeStorageLocation{{S3StorageLocationParams: &s3StorageLocationA}},
			ExpectedOutput: -1,
		},
		{
			Name:           "Second list empty",
			ListA:          []ExternalVolumeStorageLocation{{S3StorageLocationParams: &s3StorageLocationA}},
			ListB:          []ExternalVolumeStorageLocation{},
			ExpectedOutput: -1,
		},
		{
			Name:           "Lists with no common prefix - length 1",
			ListA:          []ExternalVolumeStorageLocation{{S3StorageLocationParams: &s3StorageLocationA}},
			ListB:          []ExternalVolumeStorageLocation{{S3StorageLocationParams: &s3StorageLocationB}},
			ExpectedOutput: -1,
		},
		{
			Name:           "Lists with no common prefix - length 2",
			ListA:          []ExternalVolumeStorageLocation{{S3StorageLocationParams: &s3StorageLocationA}, {AzureStorageLocationParams: &azureStorageLocationA}},
			ListB:          []ExternalVolumeStorageLocation{{S3StorageLocationParams: &s3StorageLocationB}, {AzureStorageLocationParams: &azureStorageLocationB}},
			ExpectedOutput: -1,
		},
		{
			Name:           "Identical lists - length 1",
			ListA:          []ExternalVolumeStorageLocation{{S3StorageLocationParams: &s3StorageLocationA}},
			ListB:          []ExternalVolumeStorageLocation{{S3StorageLocationParams: &s3StorageLocationA}},
			ExpectedOutput: 0,
		},
		{
			Name:           "Identical lists - length 2",
			ListA:          []ExternalVolumeStorageLocation{{S3StorageLocationParams: &s3StorageLocationA}, {AzureStorageLocationParams: &azureStorageLocationA}},
			ListB:          []ExternalVolumeStorageLocation{{S3StorageLocationParams: &s3StorageLocationA}, {AzureStorageLocationParams: &azureStorageLocationA}},
			ExpectedOutput: 1,
		},
		{
			Name: "Identical lists - length 3",
			ListA: []ExternalVolumeStorageLocation{
				{S3StorageLocationParams: &s3StorageLocationA},
				{AzureStorageLocationParams: &azureStorageLocationA},
				{S3StorageLocationParams: &s3GovStorageLocationA},
			},
			ListB: []ExternalVolumeStorageLocation{
				{S3StorageLocationParams: &s3StorageLocationA},
				{AzureStorageLocationParams: &azureStorageLocationA},
				{S3StorageLocationParams: &s3GovStorageLocationA},
			},
			ExpectedOutput: 2,
		},
		{
			Name: "Lists with a common prefix - length 3, matching up to and including index 1",
			ListA: []ExternalVolumeStorageLocation{
				{S3StorageLocationParams: &s3StorageLocationA},
				{AzureStorageLocationParams: &azureStorageLocationA},
				{GCSStorageLocationParams: &gcsStorageLocationA},
			},
			ListB: []ExternalVolumeStorageLocation{
				{S3StorageLocationParams: &s3StorageLocationA},
				{AzureStorageLocationParams: &azureStorageLocationA},
				{GCSStorageLocationParams: &gcsStorageLocationB},
			},
			ExpectedOutput: 1,
		},
		{
			Name: "Lists with a common prefix - length 4, matching up to and including index 2",
			ListA: []ExternalVolumeStorageLocation{
				{S3StorageLocationParams: &s3StorageLocationA},
				{AzureStorageLocationParams: &azureStorageLocationA},
				{GCSStorageLocationParams: &gcsStorageLocationA},
				{GCSStorageLocationParams: &gcsStorageLocationB},
			},
			ListB: []ExternalVolumeStorageLocation{
				{S3StorageLocationParams: &s3StorageLocationA},
				{AzureStorageLocationParams: &azureStorageLocationA},
				{GCSStorageLocationParams: &gcsStorageLocationA},
				{GCSStorageLocationParams: &gcsStorageLocationC},
			},
			ExpectedOutput: 2,
		},
		{
			Name: "Lists with a common prefix - length 4, matching up to and including index 1",
			ListA: []ExternalVolumeStorageLocation{
				{S3StorageLocationParams: &s3StorageLocationA},
				{AzureStorageLocationParams: &azureStorageLocationA},
				{GCSStorageLocationParams: &gcsStorageLocationA},
				{GCSStorageLocationParams: &gcsStorageLocationC},
			},
			ListB: []ExternalVolumeStorageLocation{
				{S3StorageLocationParams: &s3StorageLocationA},
				{AzureStorageLocationParams: &azureStorageLocationA},
				{GCSStorageLocationParams: &gcsStorageLocationB},
				{GCSStorageLocationParams: &gcsStorageLocationC},
			},
			ExpectedOutput: 1,
		},
		{
			Name: "Lists with a common prefix - different lengths, matching up to and including index 1 (last index of shorter list)",
			ListA: []ExternalVolumeStorageLocation{
				{S3StorageLocationParams: &s3StorageLocationA},
				{AzureStorageLocationParams: &azureStorageLocationA},
				{GCSStorageLocationParams: &gcsStorageLocationA},
			},
			ListB: []ExternalVolumeStorageLocation{
				{S3StorageLocationParams: &s3StorageLocationA},
				{AzureStorageLocationParams: &azureStorageLocationA},
			},
			ExpectedOutput: 1,
		},
		{
			Name: "Lists with a common prefix - different lengths, matching up to and including index 2",
			ListA: []ExternalVolumeStorageLocation{
				{S3StorageLocationParams: &s3StorageLocationA},
				{AzureStorageLocationParams: &azureStorageLocationA},
				{S3StorageLocationParams: &s3StorageLocationB},
				{GCSStorageLocationParams: &gcsStorageLocationA},
				{GCSStorageLocationParams: &gcsStorageLocationB},
				{AzureStorageLocationParams: &azureStorageLocationB},
			},
			ListB: []ExternalVolumeStorageLocation{
				{S3StorageLocationParams: &s3StorageLocationA},
				{AzureStorageLocationParams: &azureStorageLocationA},
				{S3StorageLocationParams: &s3StorageLocationB},
				{GCSStorageLocationParams: &gcsStorageLocationB},
				{AzureStorageLocationParams: &azureStorageLocationB},
			},
			ExpectedOutput: 2,
		},
		{
			Name:           "Empty S3 storage location",
			ListA:          []ExternalVolumeStorageLocation{{S3StorageLocationParams: &s3StorageLocationA}},
			ListB:          []ExternalVolumeStorageLocation{{S3StorageLocationParams: &S3StorageLocationParams{}}},
			ExpectedOutput: -1,
		},
		{
			Name:           "Empty GCS storage location",
			ListA:          []ExternalVolumeStorageLocation{{S3StorageLocationParams: &s3StorageLocationA}},
			ListB:          []ExternalVolumeStorageLocation{{GCSStorageLocationParams: &GCSStorageLocationParams{}}},
			ExpectedOutput: -1,
		},
		{
			Name:           "Empty Azure storage location",
			ListA:          []ExternalVolumeStorageLocation{{S3StorageLocationParams: &s3StorageLocationA}},
			ListB:          []ExternalVolumeStorageLocation{{AzureStorageLocationParams: &AzureStorageLocationParams{}}},
			ExpectedOutput: -1,
		},
		{
			Name:           "Empty storage location",
			ListA:          []ExternalVolumeStorageLocation{{S3StorageLocationParams: &s3StorageLocationA}},
			ListB:          []ExternalVolumeStorageLocation{{}},
			ExpectedOutput: -1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			commonPrefixLastIndex, err := CommonPrefixLastIndex(tc.ListA, tc.ListB)
			require.NoError(t, err)
			assert.Equal(t, tc.ExpectedOutput, commonPrefixLastIndex)
		})
	}
}
