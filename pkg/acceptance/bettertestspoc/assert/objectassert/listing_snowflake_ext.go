package objectassert

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (l *ListingAssert) HasGlobalNameNotEmpty() *ListingAssert {
	l.AddAssertion(func(t *testing.T, o *sdk.Listing) error {
		t.Helper()
		if o.GlobalName == "" {
			return fmt.Errorf("expected global_name to be not empty")
		}
		return nil
	})
	return l
}

func (l *ListingAssert) HasCreatedOnNotEmpty() *ListingAssert {
	l.AddAssertion(func(t *testing.T, o *sdk.Listing) error {
		t.Helper()
		if o.CreatedOn == "" {
			return fmt.Errorf("expected created_on to be not empty")
		}
		return nil
	})
	return l
}

func (l *ListingAssert) HasUpdatedOnNotEmpty() *ListingAssert {
	l.AddAssertion(func(t *testing.T, o *sdk.Listing) error {
		t.Helper()
		if o.UpdatedOn == "" {
			return fmt.Errorf("expected updated_on to be not empty")
		}
		return nil
	})
	return l
}

func (l *ListingAssert) HasDetailedTargetAccountsNotEmpty() *ListingAssert {
	l.AddAssertion(func(t *testing.T, o *sdk.Listing) error {
		t.Helper()
		if *o.DetailedTargetAccounts == "" {
			return fmt.Errorf("expected detailed_target_accounts to be not empty")
		}
		return nil
	})
	return l
}
