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

		e.AddAssertion(assert.ResourceDescribeOutputValueSet(fmt.Sprintf("%s.storage_allowed_locations.#", prefix), strconv.Itoa(len(loc.StorageAllowedLocations))))
		for j, allowedLoc := range loc.StorageAllowedLocations {
			e.AddAssertion(assert.ResourceDescribeOutputValueSet(fmt.Sprintf("%s.storage_allowed_locations.%d", prefix, j), allowedLoc))
		}

		if loc.S3StorageLocation != nil {
			s3Prefix := prefix + ".s3_storage_location.0"
			e.AddAssertion(assert.ResourceDescribeOutputValueSet(prefix+".s3_storage_location.#", "1"))
			e.AddAssertion(assert.ResourceDescribeOutputValueSet(s3Prefix+".storage_aws_role_arn", loc.S3StorageLocation.StorageAwsRoleArn))
			e.AddAssertion(assert.ResourceDescribeOutputValueSet(s3Prefix+".storage_aws_iam_user_arn", loc.S3StorageLocation.StorageAwsIamUserArn))
			e.AddAssertion(assert.ResourceDescribeOutputValueSet(s3Prefix+".storage_aws_external_id", loc.S3StorageLocation.StorageAwsExternalId))
			e.AddAssertion(assert.ResourceDescribeOutputValueSet(s3Prefix+".storage_aws_access_point_arn", loc.S3StorageLocation.StorageAwsAccessPointArn))
			e.AddAssertion(assert.ResourceDescribeOutputValueSet(s3Prefix+".use_privatelink_endpoint", loc.S3StorageLocation.UsePrivatelinkEndpoint))
			e.AddAssertion(assert.ResourceDescribeOutputValueSet(s3Prefix+".encryption_kms_key_id", loc.S3StorageLocation.EncryptionKmsKeyId))
		} else {
			e.AddAssertion(assert.ResourceDescribeOutputValueSet(prefix+".s3_storage_location.#", "0"))
		}

		if loc.GCSStorageLocation != nil {
			gcsPrefix := prefix + ".gcs_storage_location.0"
			e.AddAssertion(assert.ResourceDescribeOutputValueSet(prefix+".gcs_storage_location.#", "1"))
			e.AddAssertion(assert.ResourceDescribeOutputValueSet(gcsPrefix+".storage_gcp_service_account", loc.GCSStorageLocation.StorageGcpServiceAccount))
			e.AddAssertion(assert.ResourceDescribeOutputValueSet(gcsPrefix+".encryption_kms_key_id", loc.GCSStorageLocation.EncryptionKmsKeyId))
		} else {
			e.AddAssertion(assert.ResourceDescribeOutputValueSet(prefix+".gcs_storage_location.#", "0"))
		}

		if loc.AzureStorageLocation != nil {
			azurePrefix := prefix + ".azure_storage_location.0"
			e.AddAssertion(assert.ResourceDescribeOutputValueSet(prefix+".azure_storage_location.#", "1"))
			e.AddAssertion(assert.ResourceDescribeOutputValueSet(azurePrefix+".azure_tenant_id", loc.AzureStorageLocation.AzureTenantId))
			e.AddAssertion(assert.ResourceDescribeOutputValueSet(azurePrefix+".azure_multi_tenant_app_name", loc.AzureStorageLocation.AzureMultiTenantAppName))
			e.AddAssertion(assert.ResourceDescribeOutputValueSet(azurePrefix+".azure_consent_url", loc.AzureStorageLocation.AzureConsentUrl))
		} else {
			e.AddAssertion(assert.ResourceDescribeOutputValueSet(prefix+".azure_storage_location.#", "0"))
		}

		if loc.S3CompatStorageLocation != nil {
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
