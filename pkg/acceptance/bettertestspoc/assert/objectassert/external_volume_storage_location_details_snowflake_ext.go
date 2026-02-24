package objectassert

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (e *ExternalVolumeStorageLocationDetailsAssert) HasGCSStorageLocationCustomFields(expected sdk.StorageLocationGcsDetails) *ExternalVolumeStorageLocationDetailsAssert {
	e.AddAssertion(func(t *testing.T, o *sdk.ExternalVolumeStorageLocationDetails) error {
		t.Helper()
		if o.GCSStorageLocation == nil {
			return fmt.Errorf("expected gcs storage location to have value; got: nil")
		}
		if o.GCSStorageLocation.EncryptionKmsKeyId != expected.EncryptionKmsKeyId {
			return fmt.Errorf("expected storage encryption kms key id: %v; got: %v", expected.EncryptionKmsKeyId, o.GCSStorageLocation.EncryptionKmsKeyId)
		}

		// read-only fields
		if o.GCSStorageLocation.StorageGcpServiceAccount == "" {
			return fmt.Errorf("expected storage gcp service account to not be empty; got empty")
		}
		return nil
	})
	return e
}
