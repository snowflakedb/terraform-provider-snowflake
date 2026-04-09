package objectassert

import (
	"errors"
	"fmt"
	"slices"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (s *SessionPolicyDetailsAssert) HasAllowedSecondaryRolesUnordered(expected ...string) *SessionPolicyDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.SessionPolicyDetails) error {
		t.Helper()
		if len(o.AllowedSecondaryRoles) != len(expected) {
			return fmt.Errorf("expected allowed values length: %v; got: %v", len(expected), len(o.AllowedSecondaryRoles))
		}
		var errs []error
		for _, wantElem := range expected {
			if !slices.ContainsFunc(o.AllowedSecondaryRoles, func(gotElem string) bool {
				return wantElem == gotElem
			}) {
				errs = append(errs, fmt.Errorf("expected value: %s, to be in the value list: %v", wantElem, o.AllowedSecondaryRoles))
			}
		}
		return errors.Join(errs...)
	})
	return s
}
