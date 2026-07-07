package resourceshowoutputassert

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (c *ConnectionShowOutputAssert) HasCreatedOnNotEmpty() *ConnectionShowOutputAssert {
	c.ValuePresent("created_on")
	return c
}

func (c *ConnectionShowOutputAssert) HasPrimaryIdentifier(expected sdk.ExternalObjectIdentifier) *ConnectionShowOutputAssert {
	c.StringValueSet("primary", expected.FullyQualifiedName())
	return c
}

func (c *ConnectionShowOutputAssert) HasFailoverAllowedToAccounts(expected ...sdk.AccountIdentifier) *ConnectionShowOutputAssert {
	c.StringValueSet("failover_allowed_to_accounts.#", fmt.Sprintf("%d", len(expected)))
	for _, v := range expected {
		c.SetContainsElem("failover_allowed_to_accounts", v.Name())
	}
	return c
}
