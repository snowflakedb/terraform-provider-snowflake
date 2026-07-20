package objectassert

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (s *StorageIntegrationGcsDetailsAssert) HasServiceAccountNotEmpty() *StorageIntegrationGcsDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.StorageIntegrationGcsDetails) error {
		t.Helper()
		if o.ServiceAccount == "" {
			return fmt.Errorf("expected service account not empty; got empty")
		}
		return nil
	})
	return s
}
