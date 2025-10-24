package resourceshowoutputassert

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func (r *RoleShowOutputAssert) HasCreatedOnNotEmpty() *RoleShowOutputAssert {
	r.AddAssertion(assert.ResourceShowOutputValuePresent("created_on"))
	return r
}
