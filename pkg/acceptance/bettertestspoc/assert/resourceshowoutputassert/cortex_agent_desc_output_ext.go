package resourceshowoutputassert

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func (c *CortexAgentDescribeOutputAssert) HasCreatedOnNotEmpty() *CortexAgentDescribeOutputAssert {
	c.AddAssertion(assert.ResourceDescribeOutputValuePresent("created_on"))
	return c
}
