package resourceshowoutputassert

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func (c *SemanticViewShowOutputAssert) HasCreatedOnNotEmpty() *SemanticViewShowOutputAssert {
	c.AddAssertion(assert.ResourceShowOutputValuePresent("created_on"))
	return c
}
