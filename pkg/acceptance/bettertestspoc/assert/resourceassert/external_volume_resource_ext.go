package resourceassert

import (
	"fmt"
	"strconv"
)

func (e *ExternalVolumeResourceAssert) HasStorageLocationLength(len int) *ExternalVolumeResourceAssert {
	e.ValueSet("storage_location.#", strconv.FormatInt(int64(len), 10))
	return e
}

func (e *ExternalVolumeResourceAssert) HasS3StorageLocationAtIndex(
	index int,
	expectedStorageProvider string,
	expectedName string,
	expectedStorageBaseUrl string,
	expectedStorageAwsRoleArn string,
	expectedStorageAwsExternalId string,
	expectedEncryptionType string,
	expectedEncryptionKmsKeyId string,
	expectedStorageAwsAccessPointArn string,
	expectedUsePrivatelinkEndpoint string,
) *ExternalVolumeResourceAssert {
	prefix := fmt.Sprintf("storage_location.%s", strconv.Itoa(index))
	e.ValueSet(prefix+".storage_location_name", expectedName)
	e.ValueSet(prefix+".storage_provider", expectedStorageProvider)
	e.ValueSet(prefix+".storage_base_url", expectedStorageBaseUrl)
	e.ValueSet(prefix+".storage_aws_role_arn", expectedStorageAwsRoleArn)
	e.ValueSet(prefix+".storage_aws_external_id", expectedStorageAwsExternalId)
	e.ValueSet(prefix+".encryption_type", expectedEncryptionType)
	e.ValueSet(prefix+".encryption_kms_key_id", expectedEncryptionKmsKeyId)
	e.ValueSet(prefix+".storage_aws_access_point_arn", expectedStorageAwsAccessPointArn)
	e.ValueSet(prefix+".use_privatelink_endpoint", expectedUsePrivatelinkEndpoint)
	return e
}

func (e *ExternalVolumeResourceAssert) HasGCSStorageLocationAtIndex(
	index int,
	expectedName string,
	expectedStorageBaseUrl string,
	expectedEncryptionType string,
	expectedEncryptionKmsKeyId string,
) *ExternalVolumeResourceAssert {
	prefix := fmt.Sprintf("storage_location.%s", strconv.Itoa(index))
	e.ValueSet(prefix+".storage_location_name", expectedName)
	e.ValueSet(prefix+".storage_provider", "GCS")
	e.ValueSet(prefix+".storage_base_url", expectedStorageBaseUrl)
	e.ValueSet(prefix+".encryption_type", expectedEncryptionType)
	e.ValueSet(prefix+".encryption_kms_key_id", expectedEncryptionKmsKeyId)
	return e
}

func (e *ExternalVolumeResourceAssert) HasAzureStorageLocationAtIndex(
	index int,
	expectedName string,
	expectedStorageBaseUrl string,
	expectedEncryptionType string,
	expectedAzureTenantId string,
	expectedUsePrivatelinkEndpoint string,
) *ExternalVolumeResourceAssert {
	prefix := fmt.Sprintf("storage_location.%s", strconv.Itoa(index))
	e.ValueSet(prefix+".storage_location_name", expectedName)
	e.ValueSet(prefix+".storage_provider", "AZURE")
	e.ValueSet(prefix+".storage_base_url", expectedStorageBaseUrl)
	e.ValueSet(prefix+".encryption_type", expectedEncryptionType)
	e.ValueSet(prefix+".azure_tenant_id", expectedAzureTenantId)
	e.ValueSet(prefix+".use_privatelink_endpoint", expectedUsePrivatelinkEndpoint)
	return e
}

func (e *ExternalVolumeResourceAssert) HasS3CompatStorageLocationAtIndex(
	index int,
	expectedName string,
	expectedStorageBaseUrl string,
	expectedStorageEndpoint string,
	expectedStorageAwsKeyId string,
) *ExternalVolumeResourceAssert {
	prefix := fmt.Sprintf("storage_location.%s", strconv.Itoa(index))
	e.ValueSet(prefix+".storage_location_name", expectedName)
	e.ValueSet(prefix+".storage_provider", "S3COMPAT")
	e.ValueSet(prefix+".storage_base_url", expectedStorageBaseUrl)
	e.ValueSet(prefix+".storage_endpoint", expectedStorageEndpoint)
	e.ValueSet(prefix+".storage_aws_key_id", expectedStorageAwsKeyId)
	return e
}
