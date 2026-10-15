package resources

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_parseStorageLocation(t *testing.T) {
	blocks := []map[string]any{
		{
			"storage_location_name":        "s3-loc",
			"storage_provider":             "S3",
			"storage_base_url":             "s3://bucket/",
			"storage_aws_role_arn":         "arn:aws:iam::123:role/test",
			"storage_aws_external_id":      "ext",
			"storage_aws_access_point_arn": "arn:aws:s3:us-west-2:123:accesspoint/ap",
			"encryption_type":              "AWS_SSE_KMS",
			"encryption_kms_key_id":        "key",
			"use_privatelink_endpoint":     "true",
		},
		{
			"storage_location_name": "gcs-loc",
			"storage_provider":      "GCS",
			"storage_base_url":      "gcs://bucket/",
			"encryption_type":       "NONE",
		},
		{
			"storage_location_name":    "azure-loc",
			"storage_provider":         "AZURE",
			"storage_base_url":         "azure://bucket/",
			"azure_tenant_id":          "tenant",
			"encryption_type":          "NONE",
			"use_privatelink_endpoint": "true",
		},
		{
			"storage_location_name":  "compat-loc",
			"storage_provider":       "S3COMPAT",
			"storage_base_url":       "s3://compat/",
			"storage_endpoint":       "s3.us-west-2.amazonaws.com",
			"storage_aws_key_id":     "key-id",
			"storage_aws_secret_key": "",
			"encryption_type":        "",
		},
	}
	raw := make([]any, len(blocks))
	for i, block := range blocks {
		raw[i] = block
	}

	d := schema.TestResourceDataRaw(t, externalVolumeSchema, map[string]any{
		"name":             "vol",
		"storage_location": raw,
	})
	var fromPrefix []sdk.ExternalVolumeStorageLocationItem
	require.NoError(t, listBlockAttributeCreate(d, "storage_location", &fromPrefix, parseStorageLocation))

	require.Len(t, fromPrefix, 4)
	assert.Equal(t, "s3-loc", fromPrefix[0].ExternalVolumeStorageLocation.Name)
	s3 := fromPrefix[0].ExternalVolumeStorageLocation.S3StorageLocationParams
	require.NotNil(t, s3)
	assert.Equal(t, sdk.S3StorageProviderS3, s3.StorageProvider)
	assert.Equal(t, "ext", *s3.StorageAwsExternalId)
	assert.Equal(t, sdk.S3EncryptionTypeAwsSseKms, s3.Encryption.EncryptionType)
	assert.Equal(t, "key", *s3.Encryption.KmsKeyId)
	require.NotNil(t, s3.UsePrivatelinkEndpoint)
	assert.True(t, *s3.UsePrivatelinkEndpoint)

	gcs := fromPrefix[1].ExternalVolumeStorageLocation.GCSStorageLocationParams
	require.NotNil(t, gcs)
	assert.Equal(t, sdk.GCSEncryptionTypeNone, gcs.Encryption.EncryptionType)
	assert.Nil(t, gcs.Encryption.KmsKeyId)

	azure := fromPrefix[2].ExternalVolumeStorageLocation.AzureStorageLocationParams
	require.NotNil(t, azure)
	assert.Equal(t, "tenant", azure.AzureTenantId)
	require.NotNil(t, azure.UsePrivatelinkEndpoint)
	assert.True(t, *azure.UsePrivatelinkEndpoint)

	compat := fromPrefix[3].ExternalVolumeStorageLocation.S3CompatStorageLocationParams
	require.NotNil(t, compat)
	assert.Equal(t, "", compat.Credentials.AwsSecretKey)
}

func Test_parseStorageLocation_unsetOptionalFields(t *testing.T) {
	block := map[string]any{
		"storage_location_name": "s3-loc",
		"storage_provider":      "s3gov",
		"storage_base_url":      "s3://bucket/",
		"storage_aws_role_arn":  "arn:aws:iam::123:role/test",
	}
	got := parseStorageLocationForTest(t, block)
	s3 := got.ExternalVolumeStorageLocation.S3StorageLocationParams
	require.NotNil(t, s3)
	assert.Equal(t, sdk.S3StorageProviderS3gov, s3.StorageProvider)
	assert.Nil(t, s3.StorageAwsExternalId)
	assert.Nil(t, s3.StorageAwsAccessPointArn)
	assert.Nil(t, s3.UsePrivatelinkEndpoint)
	assert.Nil(t, s3.Encryption)
}

func Test_parseStorageLocation_validation(t *testing.T) {
	base := map[string]any{
		"storage_location_name": "loc",
		"storage_base_url":      "s3://bucket/",
	}
	tests := []struct {
		name    string
		block   map[string]any
		wantErr string
	}{
		{
			name: "s3 requires role arn",
			block: collections.MergeMaps(base, map[string]any{
				"storage_provider": "S3",
			}),
			wantErr: "storage_aws_role_arn is required for s3 storage location",
		},
		{
			name: "s3 rejects azure tenant",
			block: collections.MergeMaps(base, map[string]any{
				"storage_provider":     "S3",
				"storage_aws_role_arn": "arn",
				"azure_tenant_id":      "tenant",
			}),
			wantErr: "azure_tenant_id is not supported for s3 storage location",
		},
		{
			name: "gcs rejects role arn",
			block: collections.MergeMaps(base, map[string]any{
				"storage_provider":     "GCS",
				"storage_aws_role_arn": "arn",
			}),
			wantErr: "storage_aws_role_arn is not supported for gcs storage location",
		},
		{
			name: "gcs rejects external id",
			block: collections.MergeMaps(base, map[string]any{
				"storage_provider":        "GCS",
				"storage_aws_external_id": "ext",
			}),
			wantErr: "storage_aws_external_id is not supported for gcs storage location",
		},
		{
			name: "azure requires tenant",
			block: collections.MergeMaps(base, map[string]any{
				"storage_provider": "AZURE",
			}),
			wantErr: "missing azure_tenant_id provider key in an azure storage location",
		},
		{
			name: "azure rejects role arn",
			block: collections.MergeMaps(base, map[string]any{
				"storage_provider":     "AZURE",
				"azure_tenant_id":      "tenant",
				"storage_aws_role_arn": "arn",
			}),
			wantErr: "storage_aws_role_arn is not supported for azure storage location",
		},
		{
			name: "azure rejects kms key",
			block: collections.MergeMaps(base, map[string]any{
				"storage_provider":      "AZURE",
				"azure_tenant_id":       "tenant",
				"encryption_kms_key_id": "key",
			}),
			wantErr: "encryption_kms_key_id is not supported for azure storage location",
		},
		{
			name: "azure rejects encryption type",
			block: collections.MergeMaps(base, map[string]any{
				"storage_provider": "AZURE",
				"azure_tenant_id":  "tenant",
				"encryption_type":  "GCS_SSE_KMS",
			}),
			wantErr: "encryption_type is not supported for azure storage location",
		},
		{
			name: "azure rejects external id",
			block: collections.MergeMaps(base, map[string]any{
				"storage_provider":        "AZURE",
				"azure_tenant_id":         "tenant",
				"storage_aws_external_id": "ext",
			}),
			wantErr: "storage_aws_external_id is not supported for azure storage location",
		},
		{
			name: "s3compat requires endpoint",
			block: collections.MergeMaps(base, map[string]any{
				"storage_provider":       "S3COMPAT",
				"storage_aws_key_id":     "key-id",
				"storage_aws_secret_key": "secret",
			}),
			wantErr: "storage_endpoint is required for s3compat storage location",
		},
		{
			name: "s3compat requires key id",
			block: collections.MergeMaps(base, map[string]any{
				"storage_provider":       "S3COMPAT",
				"storage_endpoint":       "endpoint",
				"storage_aws_secret_key": "secret",
			}),
			wantErr: "storage_aws_key_id is required for s3compat storage location",
		},
		{
			name: "s3compat rejects external id",
			block: collections.MergeMaps(base, map[string]any{
				"storage_provider":        "S3COMPAT",
				"storage_endpoint":        "endpoint",
				"storage_aws_key_id":      "key-id",
				"storage_aws_secret_key":  "secret",
				"storage_aws_external_id": "ext",
			}),
			wantErr: "storage_aws_external_id is not supported for s3compat storage location",
		},
		{
			name: "invalid provider",
			block: collections.MergeMaps(base, map[string]any{
				"storage_provider": "INVALID-STORAGE-PROVIDER",
			}),
			wantErr: "invalid storage provider: INVALID-STORAGE-PROVIDER",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, externalVolumeSchema, map[string]any{
				"name":             "vol",
				"storage_location": []any{tt.block},
			})
			var got []sdk.ExternalVolumeStorageLocationItem
			err := listBlockAttributeCreate(d, "storage_location", &got, parseStorageLocation)
			require.ErrorContains(t, err, tt.wantErr)
		})
	}
}

func parseStorageLocationForTest(t *testing.T, block map[string]any) sdk.ExternalVolumeStorageLocationItem {
	t.Helper()
	d := schema.TestResourceDataRaw(t, externalVolumeSchema, map[string]any{
		"name":             "vol",
		"storage_location": []any{block},
	})
	var got []sdk.ExternalVolumeStorageLocationItem
	require.NoError(t, listBlockAttributeCreate(d, "storage_location", &got, parseStorageLocation))
	require.Len(t, got, 1)
	return got[0]
}
