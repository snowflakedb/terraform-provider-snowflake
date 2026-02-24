package objectassert

import (
	"errors"
	"fmt"
	"slices"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type ExternalVolumeDetailsAssert struct {
	*assert.SnowflakeObjectAssert[sdk.ExternalVolumeDetails, sdk.AccountObjectIdentifier]
}

func ExternalVolumeDetails(t *testing.T, id sdk.AccountObjectIdentifier) *ExternalVolumeDetailsAssert {
	t.Helper()
	return &ExternalVolumeDetailsAssert{
		assert.NewSnowflakeObjectAssertWithTestClientObjectProvider(
			sdk.ObjectType("EXTERNAL_VOLUME_DETAILS"),
			id,
			func(testClient *helpers.TestClient) assert.ObjectProvider[sdk.ExternalVolumeDetails, sdk.AccountObjectIdentifier] {
				return testClient.ExternalVolume.Describe
			}),
	}
}

func (e *ExternalVolumeDetailsAssert) HasActive(expected string) *ExternalVolumeDetailsAssert {
	e.AddAssertion(func(t *testing.T, o *sdk.ExternalVolumeDetails) error {
		t.Helper()
		if o.Active != expected {
			return fmt.Errorf("expected active: %v; got: %v", expected, o.Active)
		}
		return nil
	})
	return e
}

func (e *ExternalVolumeDetailsAssert) HasComment(expected string) *ExternalVolumeDetailsAssert {
	e.AddAssertion(func(t *testing.T, o *sdk.ExternalVolumeDetails) error {
		t.Helper()
		if o.Comment != expected {
			return fmt.Errorf("expected comment: %v; got: %v", expected, o.Comment)
		}
		return nil
	})
	return e
}

func (e *ExternalVolumeDetailsAssert) HasAllowWrites(expected string) *ExternalVolumeDetailsAssert {
	e.AddAssertion(func(t *testing.T, o *sdk.ExternalVolumeDetails) error {
		t.Helper()
		if o.AllowWrites != expected {
			return fmt.Errorf("expected allow writes: %v; got: %v", expected, o.AllowWrites)
		}
		return nil
	})
	return e
}

// HasStorageLocations compares storage locations by user-controllable fields only,
// ignoring Snowflake-populated fields.
func (e *ExternalVolumeDetailsAssert) HasStorageLocations(expected ...sdk.ExternalVolumeStorageLocationDetails) *ExternalVolumeDetailsAssert {
	e.AddAssertion(func(t *testing.T, o *sdk.ExternalVolumeDetails) error {
		t.Helper()
		if len(o.StorageLocations) != len(expected) {
			return fmt.Errorf("expected %d storage locations; got: %d\nexpected: %v\ngot: %v", len(expected), len(o.StorageLocations), expected, o.StorageLocations)
		}
		var errs []error
		for i := range expected {
			var locationErrors []error
			actual := o.StorageLocations[i]
			exp := expected[i]

			// Common user-controllable fields
			if actual.Name != exp.Name {
				locationErrors = append(locationErrors, fmt.Errorf("Name: expected %q, got %q", exp.Name, actual.Name))
			}
			if actual.StorageProvider != exp.StorageProvider {
				locationErrors = append(locationErrors, fmt.Errorf("StorageProvider: expected %q, got %q", exp.StorageProvider, actual.StorageProvider))
			}
			if actual.StorageBaseUrl != exp.StorageBaseUrl {
				locationErrors = append(locationErrors, fmt.Errorf("StorageBaseUrl: expected %q, got %q", exp.StorageBaseUrl, actual.StorageBaseUrl))
			}
			if !slices.Equal(actual.StorageAllowedLocations, exp.StorageAllowedLocations) {
				locationErrors = append(locationErrors, fmt.Errorf("StorageAllowedLocations: expected %v, got %v", exp.StorageAllowedLocations, actual.StorageAllowedLocations))
			}
			if actual.EncryptionType != exp.EncryptionType {
				locationErrors = append(locationErrors, fmt.Errorf("EncryptionType: expected %q, got %q", exp.EncryptionType, actual.EncryptionType))
			}

			// Provider-specific sub-struct presence
			if err := compareS3StorageLocationDetails(exp.S3StorageLocation, actual.S3StorageLocation); err != nil {
				locationErrors = append(locationErrors, err)
			}
			if err := compareGCSStorageLocationDetails(exp.GCSStorageLocation, actual.GCSStorageLocation); err != nil {
				locationErrors = append(locationErrors, err)
			}
			if err := compareAzureStorageLocationDetails(exp.AzureStorageLocation, actual.AzureStorageLocation); err != nil {
				locationErrors = append(locationErrors, err)
			}
			if err := compareS3CompatStorageLocationDetails(exp.S3CompatStorageLocation, actual.S3CompatStorageLocation); err != nil {
				locationErrors = append(locationErrors, err)
			}

			if len(locationErrors) > 0 {
				errs = append(errs, fmt.Errorf("storage location at index %d differs:\n%s", i, errors.Join(locationErrors...)))
			}
		}
		return errors.Join(errs...)
	})
	return e
}

func compareS3StorageLocationDetails(expected *sdk.StorageLocationS3Details, got *sdk.StorageLocationS3Details) error {
	errs := []error{}
	if expected == nil && got == nil {
		return nil
	}
	if expected == nil || got == nil {
		return fmt.Errorf("expected s3 storage location to have value; got: nil")
	}
	if expected.StorageAwsRoleArn != got.StorageAwsRoleArn {
		errs = append(errs, fmt.Errorf("StorageAwsRoleArn: expected %q, got %q", expected.StorageAwsRoleArn, got.StorageAwsRoleArn))
	}
	if expected.StorageAwsExternalId != got.StorageAwsExternalId {
		errs = append(errs, fmt.Errorf("StorageAwsExternalId: expected %q, got %q", expected.StorageAwsExternalId, got.StorageAwsExternalId))
	}
	if expected.StorageAwsAccessPointArn != got.StorageAwsAccessPointArn {
		errs = append(errs, fmt.Errorf("StorageAwsAccessPointArn: expected %q, got %q", expected.StorageAwsAccessPointArn, got.StorageAwsAccessPointArn))
	}
	if expected.UsePrivatelinkEndpoint != got.UsePrivatelinkEndpoint {
		errs = append(errs, fmt.Errorf("UsePrivatelinkEndpoint: expected %q, got %q", expected.UsePrivatelinkEndpoint, got.UsePrivatelinkEndpoint))
	}
	if expected.EncryptionKmsKeyId != got.EncryptionKmsKeyId {
		errs = append(errs, fmt.Errorf("EncryptionKmsKeyId: expected %q, got %q", expected.EncryptionKmsKeyId, got.EncryptionKmsKeyId))
	}
	// StorageAwsIamUserArn is set by Snowflake
	if got.StorageAwsIamUserArn == "" {
		errs = append(errs, fmt.Errorf("StorageAwsIamUserArn: expected non-empty; got empty"))
	}
	return errors.Join(errs...)
}

func compareGCSStorageLocationDetails(expected *sdk.StorageLocationGcsDetails, got *sdk.StorageLocationGcsDetails) error {
	errs := []error{}
	if expected == nil && got == nil {
		return nil
	}
	if expected == nil || got == nil {
		return fmt.Errorf("expected gcs storage location to have value; got: nil")
	}
	if expected.EncryptionKmsKeyId != got.EncryptionKmsKeyId {
		errs = append(errs, fmt.Errorf("EncryptionKmsKeyId: expected %q, got %q", expected.EncryptionKmsKeyId, got.EncryptionKmsKeyId))
	}
	// StorageGcpServiceAccount is set by Snowflake
	if got.StorageGcpServiceAccount == "" {
		errs = append(errs, fmt.Errorf("StorageGcpServiceAccount: expected non-empty; got empty"))
	}
	return errors.Join(errs...)
}

func compareAzureStorageLocationDetails(expected *sdk.StorageLocationAzureDetails, got *sdk.StorageLocationAzureDetails) error {
	errs := []error{}
	if expected == nil && got == nil {
		return nil
	}
	if expected == nil || got == nil {
		return fmt.Errorf("expected azure storage location to have value; got: nil")
	}
	if expected.AzureTenantId != got.AzureTenantId {
		errs = append(errs, fmt.Errorf("AzureTenantId: expected %q, got %q", expected.AzureTenantId, got.AzureTenantId))
	}
	// AzureMultiTenantAppName is set by Snowflake
	if got.AzureMultiTenantAppName == "" {
		errs = append(errs, fmt.Errorf("AzureMultiTenantAppName: expected non-empty; got empty"))
	}
	// AzureConsentUrl is set by Snowflake
	if got.AzureConsentUrl == "" {
		errs = append(errs, fmt.Errorf("AzureConsentUrl: expected non-empty; got empty"))
	}
	return errors.Join(errs...)
}

func compareS3CompatStorageLocationDetails(expected *sdk.StorageLocationS3CompatDetails, got *sdk.StorageLocationS3CompatDetails) error {
	errs := []error{}
	if expected == nil && got == nil {
		return nil
	}
	if expected == nil || got == nil {
		return fmt.Errorf("expected s3 compat storage location to have value; got: nil")
	}
	if expected.StorageEndpoint != got.StorageEndpoint {
		errs = append(errs, fmt.Errorf("StorageEndpoint: expected %q, got %q", expected.StorageEndpoint, got.StorageEndpoint))
	}
	if expected.AwsAccessKeyId != got.AwsAccessKeyId {
		errs = append(errs, fmt.Errorf("AwsAccessKeyId: expected %q, got %q", expected.AwsAccessKeyId, got.AwsAccessKeyId))
	}
	if expected.EncryptionKmsKeyId != got.EncryptionKmsKeyId {
		errs = append(errs, fmt.Errorf("EncryptionKmsKeyId: expected %q, got %q", expected.EncryptionKmsKeyId, got.EncryptionKmsKeyId))
	}
	return errors.Join(errs...)
}
