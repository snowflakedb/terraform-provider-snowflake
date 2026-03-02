package resourceshowoutputassert

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type ExternalVolumeDescribeOutputAssert struct {
	*assert.ResourceAssert
}

func ExternalVolumeDescribeOutput(t *testing.T, name string) *ExternalVolumeDescribeOutputAssert {
	t.Helper()

	externalVolumeAssert := ExternalVolumeDescribeOutputAssert{
		ResourceAssert: assert.NewResourceAssert(name, "describe_output"),
	}
	externalVolumeAssert.AddAssertion(assert.ValueSet("describe_output.#", "1"))
	return &externalVolumeAssert
}

func ImportedExternalVolumeDescribeOutput(t *testing.T, id string) *ExternalVolumeDescribeOutputAssert {
	t.Helper()

	externalVolumeAssert := ExternalVolumeDescribeOutputAssert{
		ResourceAssert: assert.NewImportedResourceAssert(id, "describe_output"),
	}
	externalVolumeAssert.AddAssertion(assert.ValueSet("describe_output.#", "1"))
	return &externalVolumeAssert
}

////////////////////////////
// Attribute value checks //
////////////////////////////

func (e *ExternalVolumeDescribeOutputAssert) HasActive(expected string) *ExternalVolumeDescribeOutputAssert {
	e.AddAssertion(assert.ResourceDescribeOutputValueSet("active", expected))
	return e
}

func (e *ExternalVolumeDescribeOutputAssert) HasComment(expected string) *ExternalVolumeDescribeOutputAssert {
	e.AddAssertion(assert.ResourceDescribeOutputValueSet("comment", expected))
	return e
}

func (e *ExternalVolumeDescribeOutputAssert) HasAllowWrites(expected string) *ExternalVolumeDescribeOutputAssert {
	e.AddAssertion(assert.ResourceDescribeOutputValueSet("allow_writes", expected))
	return e
}

func (e *ExternalVolumeDescribeOutputAssert) HasStorageLocations(expected []sdk.ExternalVolumeStorageLocationDetails) *ExternalVolumeDescribeOutputAssert {
	e.AddAssertion(assert.ResourceDescribeOutputValueSet("storage_locations.#", strconv.Itoa(len(expected))))
	for i, loc := range expected {
		prefix := fmt.Sprintf("storage_locations.%d", i)
		e.AddAssertion(assert.ResourceDescribeOutputValueSet(prefix+".name", loc.Name))
		e.AddAssertion(assert.ResourceDescribeOutputValueSet(prefix+".storage_provider", loc.StorageProvider))
		e.AddAssertion(assert.ResourceDescribeOutputValueSet(prefix+".storage_base_url", loc.StorageBaseUrl))
		e.AddAssertion(assert.ResourceDescribeOutputValueSet(prefix+".encryption_type", loc.EncryptionType))

		e.AddAssertion(assert.ResourceDescribeOutputValueSet(prefix+".storage_allowed_locations.#", "1"))
		e.AddAssertion(assert.ResourceDescribeOutputValueSet(prefix+".storage_allowed_locations.0", fmt.Sprintf("%s/*", loc.StorageBaseUrl)))

		if loc.StorageProvider == string(sdk.StorageProviderS3) || loc.StorageProvider == string(sdk.StorageProviderS3GOV) {
			s3Prefix := prefix + ".s3_storage_location.0"
			e.AddAssertion(assert.ResourceDescribeOutputValueSet(prefix+".s3_storage_location.#", "1"))
			e.AddAssertion(assert.ResourceDescribeOutputValueSet(s3Prefix+".storage_aws_role_arn", loc.S3StorageLocation.StorageAwsRoleArn))
			// StorageAwsIamUserArn and StorageAwsExternalId are Snowflake-generated - assert presence only
			e.AddAssertion(assert.ResourceDescribeOutputValuePresent(s3Prefix + ".storage_aws_iam_user_arn"))
			e.AddAssertion(assert.ResourceDescribeOutputValuePresent(s3Prefix + ".storage_aws_external_id"))
			e.AddAssertion(assert.ResourceDescribeOutputValueSet(s3Prefix+".storage_aws_access_point_arn", loc.S3StorageLocation.StorageAwsAccessPointArn))
			e.AddAssertion(assert.ResourceDescribeOutputValueSet(s3Prefix+".use_privatelink_endpoint", loc.S3StorageLocation.UsePrivatelinkEndpoint))
			e.AddAssertion(assert.ResourceDescribeOutputValueSet(s3Prefix+".encryption_kms_key_id", loc.S3StorageLocation.EncryptionKmsKeyId))
		} else {
			e.AddAssertion(assert.ResourceDescribeOutputValueSet(prefix+".s3_storage_location.#", "0"))
		}

		if loc.StorageProvider == string(sdk.StorageProviderGCS) {
			gcsPrefix := prefix + ".gcs_storage_location.0"
			e.AddAssertion(assert.ResourceDescribeOutputValueSet(prefix+".gcs_storage_location.#", "1"))
			// StorageGcpServiceAccount is Snowflake-generated - assert presence only
			e.AddAssertion(assert.ResourceDescribeOutputValuePresent(gcsPrefix + ".storage_gcp_service_account"))
			var encryptionKmsKeyId string
			if loc.GCSStorageLocation != nil && loc.GCSStorageLocation.EncryptionKmsKeyId != "" {
				encryptionKmsKeyId = loc.GCSStorageLocation.EncryptionKmsKeyId
			}
			e.AddAssertion(assert.ResourceDescribeOutputValueSet(gcsPrefix+".encryption_kms_key_id", encryptionKmsKeyId))
		} else {
			e.AddAssertion(assert.ResourceDescribeOutputValueSet(prefix+".gcs_storage_location.#", "0"))
		}

		if loc.StorageProvider == string(sdk.StorageProviderAzure) {
			azurePrefix := prefix + ".azure_storage_location.0"
			e.AddAssertion(assert.ResourceDescribeOutputValueSet(prefix+".azure_storage_location.#", "1"))
			e.AddAssertion(assert.ResourceDescribeOutputValueSet(azurePrefix+".azure_tenant_id", loc.AzureStorageLocation.AzureTenantId))
			// AzureMultiTenantAppName and AzureConsentUrl are Snowflake-generated - assert presence only
			e.AddAssertion(assert.ResourceDescribeOutputValuePresent(azurePrefix + ".azure_multi_tenant_app_name"))
			e.AddAssertion(assert.ResourceDescribeOutputValuePresent(azurePrefix + ".azure_consent_url"))
		} else {
			e.AddAssertion(assert.ResourceDescribeOutputValueSet(prefix+".azure_storage_location.#", "0"))
		}

		if loc.StorageProvider == string(sdk.StorageProviderS3Compatible) {
			s3cPrefix := prefix + ".s3_compat_storage_location.0"
			e.AddAssertion(assert.ResourceDescribeOutputValueSet(prefix+".s3_compat_storage_location.#", "1"))
			e.AddAssertion(assert.ResourceDescribeOutputValueSet(s3cPrefix+".endpoint", loc.S3CompatStorageLocation.Endpoint))
			e.AddAssertion(assert.ResourceDescribeOutputValueSet(s3cPrefix+".aws_access_key_id", loc.S3CompatStorageLocation.AwsAccessKeyId))
			e.AddAssertion(assert.ResourceDescribeOutputValueSet(s3cPrefix+".encryption_kms_key_id", loc.S3CompatStorageLocation.EncryptionKmsKeyId))
		} else {
			e.AddAssertion(assert.ResourceDescribeOutputValueSet(prefix+".s3_compat_storage_location.#", "0"))
		}
	}
	return e
}

///////////////////////////////
// Attribute no value checks //
///////////////////////////////

func (e *ExternalVolumeDescribeOutputAssert) HasNoActive() *ExternalVolumeDescribeOutputAssert {
	e.AddAssertion(assert.ResourceDescribeOutputValueNotSet("active"))
	return e
}

func (e *ExternalVolumeDescribeOutputAssert) HasNoComment() *ExternalVolumeDescribeOutputAssert {
	e.AddAssertion(assert.ResourceDescribeOutputValueNotSet("comment"))
	return e
}

func (e *ExternalVolumeDescribeOutputAssert) HasCommentEmpty() *ExternalVolumeDescribeOutputAssert {
	e.AddAssertion(assert.ResourceDescribeOutputValueSet("comment", ""))
	return e
}

func (e *ExternalVolumeDescribeOutputAssert) HasNoAllowWrites() *ExternalVolumeDescribeOutputAssert {
	e.AddAssertion(assert.ResourceDescribeOutputValueNotSet("allow_writes"))
	return e
}

func (e *ExternalVolumeDescribeOutputAssert) HasNoStorageLocations() *ExternalVolumeDescribeOutputAssert {
	e.AddAssertion(assert.ResourceDescribeOutputValueSet("storage_locations.#", "0"))
	return e
}
