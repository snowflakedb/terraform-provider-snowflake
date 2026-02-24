package objectassert

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (s *StorageLocationGcsDetailsAssert) HasStorageGcpServiceAccountNotEmpty() *StorageLocationGcsDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.StorageLocationGcsDetails) error {
		t.Helper()
		if o.StorageGcpServiceAccount == "" {
			return fmt.Errorf("expected storage gcp service account not empty; got empty")
		}
		return nil
	})
	return s
}
