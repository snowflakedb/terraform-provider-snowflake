package objectassert

import (
	"fmt"
	"slices"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

// HasFailoverAllowedToAccountsInAnyOrder is the order-insensitive variant of HasFailoverAllowedToAccounts;
// Snowflake does not guarantee the order of the failover accounts in SHOW CONNECTIONS.
func (c *ConnectionAssert) HasFailoverAllowedToAccountsInAnyOrder(expected ...sdk.AccountIdentifier) *ConnectionAssert {
	c.AddAssertion(func(t *testing.T, o *sdk.Connection) error {
		t.Helper()
		toSorted := func(accounts []sdk.AccountIdentifier) []string {
			names := collections.Map(accounts, func(account sdk.AccountIdentifier) string { return account.FullyQualifiedName() })
			slices.Sort(names)
			return names
		}
		if !slices.Equal(toSorted(o.FailoverAllowedToAccounts), toSorted(expected)) {
			return fmt.Errorf("expected failover allowed to accounts (in any order): %v; got: %v", expected, o.FailoverAllowedToAccounts)
		}
		return nil
	})
	return c
}
