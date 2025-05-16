package resourceassert

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
)

func (a *AccountResourceAssert) HasAdminUserType(expected sdk.UserType) *AccountResourceAssert {
	a.AddAssertion(assert.ValueSet("admin_user_type", string(expected)))
	return a
}
