package resourceshowoutputassert

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func (p *PasswordPolicyShowOutputAssert) HasCreatedOnNotEmpty() *PasswordPolicyShowOutputAssert {
	p.AddAssertion(assert.ResourceShowOutputValuePresent("created_on"))
	return p
}
