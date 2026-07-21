package objectassert

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (s *StorageLocationGcsDetailsAssert) HasNoEncryptionKmsKeyId() *StorageLocationGcsDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.StorageLocationGcsDetails) error {
		t.Helper()
		if o.EncryptionKmsKeyId != "" {
			return fmt.Errorf("expected encryption kms key id to be empty; got: %v", o.EncryptionKmsKeyId)
		}
		return nil
	})
	return s
}
