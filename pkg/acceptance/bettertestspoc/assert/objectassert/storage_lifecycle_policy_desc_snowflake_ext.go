package objectassert

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
)

// HasSignature compares each element by name and uses datatypes.AreTheSame for the data type,
// because TableColumnSignature.Type is a datatypes.DataType interface that cannot be compared with ==.
func (s *StorageLifecyclePolicyDetailsAssert) HasSignature(expected ...sdk.TableColumnSignature) *StorageLifecyclePolicyDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.StorageLifecyclePolicyDetails) error {
		t.Helper()
		if len(o.Signature) != len(expected) {
			return fmt.Errorf("expected signature: %v; got: %v", expected, o.Signature)
		}
		for i := range expected {
			if o.Signature[i].Name != expected[i].Name || !datatypes.AreTheSame(o.Signature[i].Type, expected[i].Type) {
				return fmt.Errorf("expected signature: %v; got: %v", expected, o.Signature)
			}
		}
		return nil
	})
	return s
}
