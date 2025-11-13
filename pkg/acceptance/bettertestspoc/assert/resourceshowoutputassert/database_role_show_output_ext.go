package resourceshowoutputassert

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func (d *DatabaseRoleShowOutputAssert) HasCreatedOnNotEmpty() *DatabaseRoleShowOutputAssert {
	d.AddAssertion(assert.ResourceShowOutputValuePresent("created_on"))
	return d
}

func (d *DatabaseRoleShowOutputAssert) HasOwnerNotEmpty() *DatabaseRoleShowOutputAssert {
	d.AddAssertion(assert.ResourceShowOutputValuePresent("owner"))
	return d
}

func (d *DatabaseRoleShowOutputAssert) HasOwnerRoleTypeNotEmpty() *DatabaseRoleShowOutputAssert {
	d.AddAssertion(assert.ResourceShowOutputValuePresent("owner_role_type"))
	return d
}
