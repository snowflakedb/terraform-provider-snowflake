package resourceassert

import (
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
)

func (u *ServiceUserResourceAssert) HasDisabled(expected bool) *ServiceUserResourceAssert {
	u.AddAssertion(assert.ValueSet("disabled", strconv.FormatBool(expected)))
	return u
}

func (u *ServiceUserResourceAssert) HasDefaultSecondaryRolesOption(expected sdk.SecondaryRolesOption) *ServiceUserResourceAssert {
	return u.HasDefaultSecondaryRolesOptionString(string(expected))
}
