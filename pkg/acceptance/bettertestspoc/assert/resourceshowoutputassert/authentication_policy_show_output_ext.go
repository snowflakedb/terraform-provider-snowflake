package resourceshowoutputassert

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func (a *AuthenticationPolicyShowOutputAssert) HasCreatedOnNotEmpty() *AuthenticationPolicyShowOutputAssert {
	a.AddAssertion(assert.ResourceShowOutputValuePresent("created_on"))
	return a
}
