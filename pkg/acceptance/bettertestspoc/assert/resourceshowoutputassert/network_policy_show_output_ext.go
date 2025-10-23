package resourceshowoutputassert

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func (n *NetworkPolicyShowOutputAssert) HasCreatedOnNotEmpty() *NetworkPolicyShowOutputAssert {
	n.AddAssertion(assert.ResourceShowOutputValuePresent("created_on"))
	return n
}
