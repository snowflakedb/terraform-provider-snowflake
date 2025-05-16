package resourceassert

import (
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
)

func (u *LegacyServiceUserResourceAssert) HasDisabled(expected bool) *LegacyServiceUserResourceAssert {
	u.AddAssertion(assert.ValueSet("disabled", strconv.FormatBool(expected)))
	return u
}

func (u *LegacyServiceUserResourceAssert) HasMustChangePassword(expected bool) *LegacyServiceUserResourceAssert {
	u.AddAssertion(assert.ValueSet("must_change_password", strconv.FormatBool(expected)))
	return u
}

func (u *LegacyServiceUserResourceAssert) HasDefaultSecondaryRolesOption(expected sdk.SecondaryRolesOption) *LegacyServiceUserResourceAssert {
	return u.HasDefaultSecondaryRolesOptionString(string(expected))
}
