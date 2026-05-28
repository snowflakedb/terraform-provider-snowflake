package resourceshowoutputassert

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (c *CortexAgentShowOutputAssert) HasCreatedOnNotEmpty() *CortexAgentShowOutputAssert {
	c.AddAssertion(assert.ResourceShowOutputValuePresent("created_on"))
	return c
}

func (c *CortexAgentShowOutputAssert) HasProfile(expected sdk.CortexAgentProfile) *CortexAgentShowOutputAssert {
	c.AddAssertion(assert.ResourceShowOutputValueSet("profile.#", "1"))
	if expected.DisplayName != nil {
		c.AddAssertion(assert.ResourceShowOutputValueSet("profile.0.display_name", *expected.DisplayName))
	} else {
		c.AddAssertion(assert.ResourceShowOutputValueSet("profile.0.display_name", ""))
	}
	if expected.Avatar != nil {
		c.AddAssertion(assert.ResourceShowOutputValueSet("profile.0.avatar", *expected.Avatar))
	} else {
		c.AddAssertion(assert.ResourceShowOutputValueSet("profile.0.avatar", ""))
	}
	if expected.Color != nil {
		c.AddAssertion(assert.ResourceShowOutputValueSet("profile.0.color", *expected.Color))
	} else {
		c.AddAssertion(assert.ResourceShowOutputValueSet("profile.0.color", ""))
	}
	return c
}
