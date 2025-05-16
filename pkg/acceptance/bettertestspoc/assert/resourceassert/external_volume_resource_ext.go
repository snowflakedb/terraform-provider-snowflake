package resourceassert

import (
	"fmt"
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/assert"
)

func (e *ExternalVolumeResourceAssert) HasStorageLocationLength(len int) *ExternalVolumeResourceAssert {
	e.AddAssertion(assert.ValueSet("storage_location.#", strconv.FormatInt(int64(len), 10)))
	return e
}

func (e *ExternalVolumeResourceAssert) HasStorageLocationAtIndex(
	index int,
	expectedName string,
	expectedStorageProvider string,
	expectedStorageBaseUrl string,
	expectedStorageAwsRoleArn string,
	expectedEncryptionType string,
	expectedEncryptionKmsKeyId string,
	expectedAzureTenantId string,
) *ExternalVolumeResourceAssert {
	e.AddAssertion(assert.ValueSet(fmt.Sprintf("storage_location.%s.storage_location_name", strconv.Itoa(index)), expectedName))
	e.AddAssertion(assert.ValueSet(fmt.Sprintf("storage_location.%s.storage_provider", strconv.Itoa(index)), expectedStorageProvider))
	e.AddAssertion(assert.ValueSet(fmt.Sprintf("storage_location.%s.storage_base_url", strconv.Itoa(index)), expectedStorageBaseUrl))
	e.AddAssertion(assert.ValueSet(fmt.Sprintf("storage_location.%s.storage_aws_role_arn", strconv.Itoa(index)), expectedStorageAwsRoleArn))
	e.AddAssertion(assert.ValueSet(fmt.Sprintf("storage_location.%s.encryption_type", strconv.Itoa(index)), expectedEncryptionType))
	e.AddAssertion(assert.ValueSet(fmt.Sprintf("storage_location.%s.encryption_kms_key_id", strconv.Itoa(index)), expectedEncryptionKmsKeyId))
	e.AddAssertion(assert.ValueSet(fmt.Sprintf("storage_location.%s.azure_tenant_id", strconv.Itoa(index)), expectedAzureTenantId))
	return e
}
