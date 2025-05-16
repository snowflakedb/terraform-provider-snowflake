package resourceassert

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
)

func (c *PrimaryConnectionResourceAssert) HasExactlyFailoverToAccountsInOrder(expected ...sdk.AccountIdentifier) *PrimaryConnectionResourceAssert {
	c.AddAssertion(assert.ValueSet("enable_failover_to_accounts.#", fmt.Sprintf("%d", len(expected))))
	for i, v := range expected {
		c.AddAssertion(assert.ValueSet(fmt.Sprintf("enable_failover_to_accounts.%d", i), v.Name()))
	}
	return c
}

func (c *PrimaryConnectionResourceAssert) HasNoEnableFailoverToAccounts() *PrimaryConnectionResourceAssert {
	c.AddAssertion(assert.ValueSet("enable_failover_to_accounts.#", "0"))
	return c
}
