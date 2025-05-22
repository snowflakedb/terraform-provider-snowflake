package objectassert

import (
	"errors"
	"fmt"
	"slices"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
)

func (s *TagAssert) HasAllowedValues(expected ...string) *TagAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.Tag) error {
		t.Helper()
		if len(o.AllowedValues) != len(expected) {
			return fmt.Errorf("expected allowed values length: %v; got: %v", len(expected), len(o.AllowedValues))
		}
		var errs []error
		for _, wantElem := range expected {
			if !slices.ContainsFunc(o.AllowedValues, func(gotElem string) bool {
				return wantElem == gotElem
			}) {
				errs = append(errs, fmt.Errorf("expected value: %s, to be in the value list: %v", wantElem, o.AllowedValues))
			}
		}
		return errors.Join(errs...)
	})
	return s
}
