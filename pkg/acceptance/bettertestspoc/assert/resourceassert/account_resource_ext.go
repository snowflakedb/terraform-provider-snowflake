package resourceassert

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (a *AccountResourceAssert) HasAdminUserTypeEnum(expected sdk.UserType) *AccountResourceAssert {
	a.ValueSet("admin_user_type", string(expected))
	return a
}
