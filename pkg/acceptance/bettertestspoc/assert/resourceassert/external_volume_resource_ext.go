package resourceassert

import (
	"fmt"
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func (e *ExternalVolumeResourceAssert) HasStorageLocationLength(len int) *ExternalVolumeResourceAssert {
	e.AddAssertion(assert.ValueSet("storage_location.#", strconv.FormatInt(int64(len), 10)))
	return e
}

func (e *ExternalVolumeResourceAssert) HasS3StorageLocationAtIndex(
	index int,
	expectedStorageProvider string,
	expectedName string,
	expectedStorageBaseUrl string,
	expectedStorageAwsRoleArn string,
	expectedEncryptionType string,
	expectedEncryptionKmsKeyId string,
	expectedStorageAwsAccessPointArn string,
	expectedUsePrivatelinkEndpoint string,
) *ExternalVolumeResourceAssert {
	prefix := fmt.Sprintf("storage_location.%s", strconv.Itoa(index))
	e.AddAssertion(assert.ValueSet(prefix+".storage_location_name", expectedName))
	e.AddAssertion(assert.ValueSet(prefix+".storage_provider", expectedStorageProvider))
	e.AddAssertion(assert.ValueSet(prefix+".storage_base_url", expectedStorageBaseUrl))
	e.AddAssertion(assert.ValueSet(prefix+".storage_aws_role_arn", expectedStorageAwsRoleArn))
	e.AddAssertion(assert.ValueSet(prefix+".encryption_type", expectedEncryptionType))
	e.AddAssertion(assert.ValueSet(prefix+".encryption_kms_key_id", expectedEncryptionKmsKeyId))
	e.AddAssertion(assert.ValueSet(prefix+".storage_aws_access_point_arn", expectedStorageAwsAccessPointArn))
	e.AddAssertion(assert.ValueSet(prefix+".use_privatelink_endpoint", expectedUsePrivatelinkEndpoint))
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
	e.AddAssertion(assert.ValueSet(prefix+".storage_location_name", expectedName))
	e.AddAssertion(assert.ValueSet(prefix+".storage_provider", "GCS"))
	e.AddAssertion(assert.ValueSet(prefix+".storage_base_url", expectedStorageBaseUrl))
	e.AddAssertion(assert.ValueSet(prefix+".encryption_type", expectedEncryptionType))
	e.AddAssertion(assert.ValueSet(prefix+".encryption_kms_key_id", expectedEncryptionKmsKeyId))
	return e
}

func (e *ExternalVolumeResourceAssert) HasAzureStorageLocationAtIndex(
	index int,
	expectedName string,
	expectedStorageBaseUrl string,
	expectedEncryptionType string,
	expectedAzureTenantId string,
) *ExternalVolumeResourceAssert {
	prefix := fmt.Sprintf("storage_location.%s", strconv.Itoa(index))
	e.AddAssertion(assert.ValueSet(prefix+".storage_location_name", expectedName))
	e.AddAssertion(assert.ValueSet(prefix+".storage_provider", "AZURE"))
	e.AddAssertion(assert.ValueSet(prefix+".storage_base_url", expectedStorageBaseUrl))
	e.AddAssertion(assert.ValueSet(prefix+".encryption_type", expectedEncryptionType))
	e.AddAssertion(assert.ValueSet(prefix+".azure_tenant_id", expectedAzureTenantId))
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
	e.AddAssertion(assert.ValueSet(prefix+".storage_location_name", expectedName))
	e.AddAssertion(assert.ValueSet(prefix+".storage_provider", "S3COMPAT"))
	e.AddAssertion(assert.ValueSet(prefix+".storage_base_url", expectedStorageBaseUrl))
	e.AddAssertion(assert.ValueSet(prefix+".storage_endpoint", expectedStorageEndpoint))
	e.AddAssertion(assert.ValueSet(prefix+".storage_aws_key_id", expectedStorageAwsKeyId))
	return e
}
