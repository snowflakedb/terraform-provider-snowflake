package resourceshowoutputassert

import (
	"fmt"
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (e *ExternalVolumeDescribeOutputAssert) HasActiveEmpty() *ExternalVolumeDescribeOutputAssert {
	e.StringValueSet("active", "")
	return e
}

func (e *ExternalVolumeDescribeOutputAssert) HasCommentEmpty() *ExternalVolumeDescribeOutputAssert {
	e.StringValueSet("comment", "")
	return e
}

func (e *ExternalVolumeDescribeOutputAssert) HasStorageLocations(expected []sdk.ExternalVolumeStorageLocationDetails) *ExternalVolumeDescribeOutputAssert {
	e.StringValueSet("storage_locations.#", strconv.Itoa(len(expected)))
	for i, loc := range expected {
		prefix := fmt.Sprintf("storage_locations.%d", i)
		e.StringValueSet(prefix+".name", loc.Name)
		e.StringValueSet(prefix+".storage_provider", loc.StorageProvider)
		e.StringValueSet(prefix+".storage_base_url", loc.StorageBaseUrl)
		e.StringValueSet(prefix+".encryption_type", loc.EncryptionType)

		e.StringValueSet(prefix+".storage_allowed_locations.#", "1")
		e.StringValueSet(prefix+".storage_allowed_locations.0", fmt.Sprintf("%s*", loc.StorageBaseUrl))

		if loc.StorageProvider == string(sdk.StorageProviderS3) || loc.StorageProvider == string(sdk.StorageProviderS3gov) {
			s3Prefix := prefix + ".s3_storage_location.0"
			e.StringValueSet(prefix+".s3_storage_location.#", "1")
			e.StringValueSet(s3Prefix+".storage_aws_role_arn", loc.S3StorageLocation.StorageAwsRoleArn)
			e.ValuePresent(s3Prefix + ".storage_aws_iam_user_arn")
			e.ValuePresent(s3Prefix + ".storage_aws_external_id")
			e.StringValueSet(s3Prefix+".storage_aws_access_point_arn", loc.S3StorageLocation.StorageAwsAccessPointArn)
			e.StringValueSet(s3Prefix+".encryption_kms_key_id", loc.S3StorageLocation.EncryptionKmsKeyId)
			if loc.S3StorageLocation.UsePrivatelinkEndpoint != nil {
				e.StringValueSet(s3Prefix+".use_privatelink_endpoint", fmt.Sprintf("%t", *loc.S3StorageLocation.UsePrivatelinkEndpoint))
			} else {
				e.StringValueSet(s3Prefix+".use_privatelink_endpoint", "")
			}
		} else {
			e.StringValueSet(prefix+".s3_storage_location.#", "0")
		}

		if loc.StorageProvider == string(sdk.StorageProviderGcs) {
			gcsPrefix := prefix + ".gcs_storage_location.0"
			e.StringValueSet(prefix+".gcs_storage_location.#", "1")
			e.ValuePresent(gcsPrefix + ".storage_gcp_service_account")
			var encryptionKmsKeyId string
			if loc.GCSStorageLocation != nil && loc.GCSStorageLocation.EncryptionKmsKeyId != "" {
				encryptionKmsKeyId = loc.GCSStorageLocation.EncryptionKmsKeyId
			}
			e.StringValueSet(gcsPrefix+".encryption_kms_key_id", encryptionKmsKeyId)
		} else {
			e.StringValueSet(prefix+".gcs_storage_location.#", "0")
		}

		if loc.StorageProvider == string(sdk.StorageProviderAzure) {
			azurePrefix := prefix + ".azure_storage_location.0"
			e.StringValueSet(prefix+".azure_storage_location.#", "1")
			e.StringValueSet(azurePrefix+".azure_tenant_id", loc.AzureStorageLocation.AzureTenantId)
			e.ValuePresent(azurePrefix + ".azure_multi_tenant_app_name")
			e.ValuePresent(azurePrefix + ".azure_consent_url")
			if loc.AzureStorageLocation.UsePrivatelinkEndpoint != nil {
				e.StringValueSet(azurePrefix+".use_privatelink_endpoint", fmt.Sprintf("%t", *loc.AzureStorageLocation.UsePrivatelinkEndpoint))
			} else {
				e.StringValueSet(azurePrefix+".use_privatelink_endpoint", "")
			}
		} else {
			e.StringValueSet(prefix+".azure_storage_location.#", "0")
		}

		if loc.StorageProvider == string(sdk.StorageProviderS3compat) {
			s3cPrefix := prefix + ".s3_compat_storage_location.0"
			e.StringValueSet(prefix+".s3_compat_storage_location.#", "1")
			e.StringValueSet(s3cPrefix+".endpoint", loc.S3CompatStorageLocation.Endpoint)
			e.StringValueSet(s3cPrefix+".aws_access_key_id", loc.S3CompatStorageLocation.AwsAccessKeyId)
			e.StringValueSet(s3cPrefix+".encryption_kms_key_id", loc.S3CompatStorageLocation.EncryptionKmsKeyId)
		} else {
			e.StringValueSet(prefix+".s3_compat_storage_location.#", "0")
		}
	}
	return e
}
