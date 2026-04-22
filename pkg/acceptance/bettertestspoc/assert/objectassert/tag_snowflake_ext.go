package objectassert

import (
	"errors"
	"fmt"
	"slices"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (t *TagAssert) HasAllowedValuesUnordered(expected ...string) *TagAssert {
	t.AddAssertion(func(t *testing.T, o *sdk.Tag) error {
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
	return t
}

func (t *TagAssert) HasAllowedValuesNil() *TagAssert {
	t.AddAssertion(func(t *testing.T, o *sdk.Tag) error {
		t.Helper()
		if o.AllowedValues != nil {
			return fmt.Errorf("expected allowed values to be nil; got: %v", o.AllowedValues)
		}
		return nil
	})
	return t
}

func (t *TagAssert) HasAllowedValuesEmpty() *TagAssert {
	t.AddAssertion(func(t *testing.T, o *sdk.Tag) error {
		t.Helper()
		if o.AllowedValues == nil {
			return fmt.Errorf("expected allowed values to be empty (non-nil); got: nil")
		}
		if len(o.AllowedValues) != 0 {
			return fmt.Errorf("expected allowed values to be empty; got: %v", o.AllowedValues)
		}
		return nil
	})
	return t
}

func (t *TagAssert) HasPropagateEnum(expected sdk.TagPropagation) *TagAssert {
	return t.HasPropagate(string(expected))
}

// HasOnConflict asserts that the tag's OnConflict value matches the expected string.
// This is only populated when BCR-2291 is enabled.
func (t *TagAssert) HasOnConflict(expected string) *TagAssert {
	t.AddAssertion(func(t *testing.T, o *sdk.Tag) error {
		t.Helper()
		if o.OnConflict == nil {
			return fmt.Errorf("expected on_conflict: %q; got nil (BCR-2291 may not be enabled)", expected)
		}
		if *o.OnConflict != expected {
			return fmt.Errorf("expected on_conflict: %q; got: %q", expected, *o.OnConflict)
		}
		return nil
	})
	return t
}

// HasOnConflictNil asserts that the tag's OnConflict field is nil (BCR-2291 not enabled or on_conflict not set).
func (t *TagAssert) HasOnConflictNil() *TagAssert {
	t.AddAssertion(func(t *testing.T, o *sdk.Tag) error {
		t.Helper()
		if o.OnConflict != nil {
			return fmt.Errorf("expected on_conflict to be nil; got: %q", *o.OnConflict)
		}
		return nil
	})
	return t
}
