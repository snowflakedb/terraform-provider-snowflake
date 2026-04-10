package model

import (
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (e *ExternalVolumeModel) WithStorageLocation(storageLocation []sdk.ExternalVolumeStorageLocationRequest) *ExternalVolumeModel {
	maps := make([]tfconfig.Variable, len(storageLocation))
	for i, v := range storageLocation {
		switch {
		case v.S3StorageLocationParams != nil:
			m := map[string]tfconfig.Variable{
				"storage_location_name": tfconfig.StringVariable(v.Name),
				"storage_provider":      tfconfig.StringVariable(string(sdk.StorageProviderS3)),
				"storage_aws_role_arn":  tfconfig.StringVariable(v.S3StorageLocationParams.StorageAwsRoleArn),
				"storage_base_url":      tfconfig.StringVariable(v.S3StorageLocationParams.StorageBaseUrl),
			}
			if v.S3StorageLocationParams.StorageAwsExternalId != nil {
				m["storage_aws_external_id"] = tfconfig.StringVariable(*v.S3StorageLocationParams.StorageAwsExternalId)
			}
			if v.S3StorageLocationParams.Encryption != nil {
				m["encryption_type"] = tfconfig.StringVariable(string(v.S3StorageLocationParams.Encryption.EncryptionType))
				if v.S3StorageLocationParams.Encryption.KmsKeyId != nil {
					m["encryption_kms_key_id"] = tfconfig.StringVariable(*v.S3StorageLocationParams.Encryption.KmsKeyId)
				}
			}
			maps[i] = tfconfig.MapVariable(m)
		case v.GCSStorageLocationParams != nil:
			m := map[string]tfconfig.Variable{
				"storage_location_name": tfconfig.StringVariable(v.Name),
				"storage_provider":      tfconfig.StringVariable(string(sdk.StorageProviderGCS)),
				"storage_base_url":      tfconfig.StringVariable(v.GCSStorageLocationParams.StorageBaseUrl),
			}
			if v.GCSStorageLocationParams.Encryption != nil {
				m["encryption_type"] = tfconfig.StringVariable(string(v.GCSStorageLocationParams.Encryption.EncryptionType))
				if v.GCSStorageLocationParams.Encryption.KmsKeyId != nil {
					m["encryption_kms_key_id"] = tfconfig.StringVariable(*v.GCSStorageLocationParams.Encryption.KmsKeyId)
				}
			}
			maps[i] = tfconfig.MapVariable(m)
		case v.AzureStorageLocationParams != nil:
			m := map[string]tfconfig.Variable{
				"storage_location_name": tfconfig.StringVariable(v.Name),
				"storage_provider":      tfconfig.StringVariable(string(sdk.StorageProviderAzure)),
				"azure_tenant_id":       tfconfig.StringVariable(v.AzureStorageLocationParams.AzureTenantId),
				"storage_base_url":      tfconfig.StringVariable(v.AzureStorageLocationParams.StorageBaseUrl),
			}
			maps[i] = tfconfig.MapVariable(m)
		case v.S3CompatStorageLocationParams != nil:
			m := map[string]tfconfig.Variable{
				"storage_location_name": tfconfig.StringVariable(v.Name),
				"storage_provider":      tfconfig.StringVariable(string(sdk.StorageProviderS3Compatible)),
				"storage_base_url":      tfconfig.StringVariable(v.S3CompatStorageLocationParams.StorageBaseUrl),
				"storage_endpoint":      tfconfig.StringVariable(v.S3CompatStorageLocationParams.StorageEndpoint),
			}
			maps[i] = tfconfig.MapVariable(m)
		}
	}
	e.StorageLocation = tfconfig.ListVariable(maps...)
	return e
}
