package resourceshowoutputassert

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (c *ConnectionShowOutputAssert) HasCreatedOnNotEmpty() *ConnectionShowOutputAssert {
	c.AddAssertion(assert.ResourceShowOutputValuePresent("created_on"))
	return c
}

func (c *ConnectionShowOutputAssert) HasPrimaryIdentifier(expected sdk.ExternalObjectIdentifier) *ConnectionShowOutputAssert {
	c.AddAssertion(assert.ResourceShowOutputValueSet("primary", expected.FullyQualifiedName()))
	return c
}

func (c *ConnectionShowOutputAssert) HasFailoverAllowedToAccounts(expected ...sdk.AccountIdentifier) *ConnectionShowOutputAssert {
	c.AddAssertion(assert.ResourceShowOutputValueSet("failover_allowed_to_accounts.#", fmt.Sprintf("%d", len(expected))))
	for _, v := range expected {
		c.AddAssertion(assert.ResourceShowOutputSetElem("failover_allowed_to_accounts", v.Name()))
	}
	return c
}
