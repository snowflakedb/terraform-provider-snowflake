package objectassert

import (
	"fmt"
	"slices"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

// HasAllowedValuesSet checks that the allowed_values field contains the expected values (order independent)
func (t *TagAssert) HasAllowedValuesSet(expected ...string) *TagAssert {
	t.AddAssertion(func(t *testing.T, o *sdk.Tag) error {
		t.Helper()
		if len(o.AllowedValues) != len(expected) {
			return fmt.Errorf("expected allowed values length: %d; got: %d", len(expected), len(o.AllowedValues))
		}

		// Sort both slices for comparison
		actualSorted := make([]string, len(o.AllowedValues))
		copy(actualSorted, o.AllowedValues)
		slices.Sort(actualSorted)

		expectedSorted := make([]string, len(expected))
		copy(expectedSorted, expected)
		slices.Sort(expectedSorted)

		if !slices.Equal(actualSorted, expectedSorted) {
			return fmt.Errorf("expected allowed values: %v; got: %v", expected, o.AllowedValues)
		}
		return nil
	})
	return t
}
