package resourceshowoutputassert

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"

func (c *ComputePoolShowOutputAssert) HasCreatedOnNotEmpty() *ComputePoolShowOutputAssert {
	c.AddAssertion(assert.ResourceShowOutputValuePresent("created_on"))
	return c
}

func (c *ComputePoolShowOutputAssert) HasResumedOnNotEmpty() *ComputePoolShowOutputAssert {
	c.AddAssertion(assert.ResourceShowOutputValuePresent("resumed_on"))
	return c
}

func (c *ComputePoolShowOutputAssert) HasUpdatedOnNotEmpty() *ComputePoolShowOutputAssert {
	c.AddAssertion(assert.ResourceShowOutputValuePresent("updated_on"))
	return c
}

func (c *ComputePoolShowOutputAssert) HasNoApplication() *ComputePoolShowOutputAssert {
	c.AddAssertion(assert.ResourceShowOutputValueSet("application", ""))
	return c
}
