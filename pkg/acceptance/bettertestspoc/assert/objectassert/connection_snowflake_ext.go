package objectassert

import (
	"fmt"
	"slices"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
)

func (c *ConnectionAssert) HasFailoverAllowedToAccounts(expected ...sdk.AccountIdentifier) *ConnectionAssert {
	c.AddAssertion(func(t *testing.T, o *sdk.Connection) error {
		t.Helper()
		if !slices.Equal(expected, o.FailoverAllowedToAccounts) {
			return fmt.Errorf("expected failover_allowed_to_accounts to be: %v; got: %v", expected, o.FailoverAllowedToAccounts)
		}
		return nil
	})
	return c
}

func (c *ConnectionAssert) HasNoComment() *ConnectionAssert {
	c.AddAssertion(func(t *testing.T, o *sdk.Connection) error {
		t.Helper()
		if o.Comment != nil {
			return fmt.Errorf("expected comment to have nil; got: %s", *o.Comment)
		}
		return nil
	})
	return c
}

func (c *ConnectionAssert) HasConnectionUrlNotEmpty() *ConnectionAssert {
	c.AddAssertion(func(t *testing.T, o *sdk.Connection) error {
		t.Helper()
		if o.ConnectionUrl == "" {
			return fmt.Errorf("expected connection url not empty, got: %s", o.ConnectionUrl)
		}
		return nil
	})

	return c
}

func (c *ConnectionAssert) HasPrimaryIdentifier(expected sdk.ExternalObjectIdentifier) *ConnectionAssert {
	c.AddAssertion(func(t *testing.T, o *sdk.Connection) error {
		t.Helper()
		if o.Primary != expected {
			return fmt.Errorf("expected primary identifier: %v; got: %v", expected.FullyQualifiedName(), o.Primary)
		}
		return nil
	})
	return c
}
