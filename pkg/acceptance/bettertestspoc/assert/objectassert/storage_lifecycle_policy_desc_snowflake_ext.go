package objectassert

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (s *StorageLifecyclePolicyDetailsAssert) HasNoArchiveForDays() *StorageLifecyclePolicyDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.StorageLifecyclePolicyDetails) error {
		t.Helper()
		if o.ArchiveForDays != nil {
			return fmt.Errorf("expected archive for days to be nil; got: %d", *o.ArchiveForDays)
		}
		return nil
	})
	return s
}
