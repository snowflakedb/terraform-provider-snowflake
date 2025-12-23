package resourceassert

import (
	"strconv"

	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (u *UserResourceAssert) HasDisabled(expected bool) *UserResourceAssert {
	u.AddAssertion(assert.ValueSet("disabled", strconv.FormatBool(expected)))
	return u
}

func (u *UserResourceAssert) HasEmptyPassword() *UserResourceAssert {
	u.AddAssertion(assert.ValueSet("password", ""))
	return u
}

func (u *UserResourceAssert) HasNotEmptyPassword() *UserResourceAssert {
	u.AddAssertion(assert.ValuePresent("password"))
	return u
}

func (u *UserResourceAssert) HasMustChangePassword(expected bool) *UserResourceAssert {
	u.AddAssertion(assert.ValueSet("must_change_password", strconv.FormatBool(expected)))
	return u
}

func (u *UserResourceAssert) HasDefaultSecondaryRolesOption(expected sdk.SecondaryRolesOption) *UserResourceAssert {
	return u.HasDefaultSecondaryRolesOptionString(string(expected))
}

func (u *UserResourceAssert) HasAllDefaults(userId sdk.AccountObjectIdentifier, expectedDefaultSecondaryRoles sdk.SecondaryRolesOption) *UserResourceAssert {
	return u.
		HasNameString(userId.Name()).
		HasNoPassword().
		HasNoLoginName().
		HasNoDisplayName().
		HasFirstNameEmpty().
		HasMiddleNameEmpty().
		HasLastNameEmpty().
		HasEmailEmpty().
		HasMustChangePasswordString(r.BooleanDefault).
		HasDisabledString(r.BooleanDefault).
		HasNoDaysToExpiry().
		HasMinsToUnlockString(r.IntDefaultString).
		HasDefaultWarehouseEmpty().
		HasNoDefaultNamespace().
		HasDefaultRoleEmpty().
		HasDefaultSecondaryRolesOption(expectedDefaultSecondaryRoles).
		HasMinsToBypassMfaString(r.IntDefaultString).
		HasRsaPublicKeyEmpty().
		HasRsaPublicKey2Empty().
		HasCommentEmpty().
		HasDisableMfaString(r.BooleanDefault).
		HasFullyQualifiedNameString(userId.FullyQualifiedName())
}
