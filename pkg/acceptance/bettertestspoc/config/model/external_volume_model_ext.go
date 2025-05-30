package model

import (
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (e *ExternalVolumeModel) WithStorageLocation(storageLocation []sdk.ExternalVolumeStorageLocation) *ExternalVolumeModel {
	maps := make([]tfconfig.Variable, len(storageLocation))
	for i, v := range storageLocation {
		switch {
		case v.S3StorageLocationParams != nil:
			m := map[string]tfconfig.Variable{
				"storage_location_name": tfconfig.StringVariable(v.S3StorageLocationParams.Name),
				"storage_provider":      tfconfig.StringVariable(string(v.S3StorageLocationParams.StorageProvider)),
				"storage_aws_role_arn":  tfconfig.StringVariable(v.S3StorageLocationParams.StorageAwsRoleArn),
				"storage_base_url":      tfconfig.StringVariable(v.S3StorageLocationParams.StorageBaseUrl),
			}
			if v.S3StorageLocationParams.StorageAwsExternalId != nil {
				m["storage_aws_external_id"] = tfconfig.StringVariable(*v.S3StorageLocationParams.StorageAwsExternalId)
			}
			if v.S3StorageLocationParams.Encryption != nil {
				m["encryption_type"] = tfconfig.StringVariable(string(v.S3StorageLocationParams.Encryption.Type))
				if v.S3StorageLocationParams.Encryption.KmsKeyId != nil {
					m["encryption_kms_key_id"] = tfconfig.StringVariable(*v.S3StorageLocationParams.Encryption.KmsKeyId)
				}
			}
			maps[i] = tfconfig.MapVariable(m)
		case v.GCSStorageLocationParams != nil:
			m := map[string]tfconfig.Variable{
				"storage_location_name": tfconfig.StringVariable(v.GCSStorageLocationParams.Name),
				"storage_provider":      tfconfig.StringVariable(v.GCSStorageLocationParams.StorageProviderGcs),
				"storage_base_url":      tfconfig.StringVariable(v.GCSStorageLocationParams.StorageBaseUrl),
			}
			if v.GCSStorageLocationParams.Encryption != nil {
				m["encryption_type"] = tfconfig.StringVariable(string(v.GCSStorageLocationParams.Encryption.Type))
				if v.GCSStorageLocationParams.Encryption.KmsKeyId != nil {
					m["encryption_kms_key_id"] = tfconfig.StringVariable(*v.GCSStorageLocationParams.Encryption.KmsKeyId)
				}
			}
			maps[i] = tfconfig.MapVariable(m)
		case v.AzureStorageLocationParams != nil:
			m := map[string]tfconfig.Variable{
				"storage_location_name": tfconfig.StringVariable(v.AzureStorageLocationParams.Name),
				"storage_provider":      tfconfig.StringVariable(v.AzureStorageLocationParams.StorageProviderAzure),
				"azure_tenant_id":       tfconfig.StringVariable(v.AzureStorageLocationParams.AzureTenantId),
				"storage_base_url":      tfconfig.StringVariable(v.GCSStorageLocationParams.StorageBaseUrl),
			}
			maps[i] = tfconfig.MapVariable(m)
		}
	}
	e.StorageLocation = tfconfig.ListVariable(maps...)
	return e
}
