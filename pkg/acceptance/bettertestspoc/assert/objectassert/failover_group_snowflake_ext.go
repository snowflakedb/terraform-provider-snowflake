package objectassert

import (
	"fmt"
	"slices"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

// The SDK object keeps these values as slices in the order of SHOW, which differs from the order of the request.
// The InAnyOrder variants compare them as sets.

func (f *FailoverGroupAssert) HasObjectTypesInAnyOrder(expected ...sdk.FailoverGroupObjectType) *FailoverGroupAssert {
	f.AddAssertion(func(t *testing.T, o *sdk.FailoverGroup) error {
		t.Helper()
		if !slices.Equal(sortedStrings(o.ObjectTypes), sortedStrings(expected)) {
			return fmt.Errorf("expected object types (in any order): %v; got: %v", expected, o.ObjectTypes)
		}
		return nil
	})
	return f
}

func (f *FailoverGroupAssert) HasAllowedIntegrationTypesInAnyOrder(expected ...sdk.FailoverGroupIntegrationType) *FailoverGroupAssert {
	f.AddAssertion(func(t *testing.T, o *sdk.FailoverGroup) error {
		t.Helper()
		if !slices.Equal(sortedStrings(o.AllowedIntegrationTypes), sortedStrings(expected)) {
			return fmt.Errorf("expected allowed integration types (in any order): %v; got: %v", expected, o.AllowedIntegrationTypes)
		}
		return nil
	})
	return f
}

func (f *FailoverGroupAssert) HasAllowedAccountsInAnyOrder(expected ...sdk.AccountIdentifier) *FailoverGroupAssert {
	f.AddAssertion(func(t *testing.T, o *sdk.FailoverGroup) error {
		t.Helper()
		toSorted := func(accounts []sdk.AccountIdentifier) []string {
			return sortedStrings(collections.Map(accounts, func(account sdk.AccountIdentifier) string { return account.FullyQualifiedName() }))
		}
		if !slices.Equal(toSorted(o.AllowedAccounts), toSorted(expected)) {
			return fmt.Errorf("expected allowed accounts (in any order): %v; got: %v", expected, o.AllowedAccounts)
		}
		return nil
	})
	return f
}

func sortedStrings[T ~string](values []T) []string {
	result := collections.Map(values, func(value T) string { return string(value) })
	slices.Sort(result)
	return result
}
