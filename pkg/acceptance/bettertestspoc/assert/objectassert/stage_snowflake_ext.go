package objectassert

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (s *StageAssert) HasNoStorageIntegration() *StageAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.Stage) error {
		t.Helper()
		if o.StorageIntegration != nil {
			return fmt.Errorf("expected storage integration to be nil; got: %s", *o.StorageIntegration)
		}
		return nil
	})
	return s
}
