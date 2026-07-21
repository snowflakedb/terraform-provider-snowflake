package objectassert

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

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
