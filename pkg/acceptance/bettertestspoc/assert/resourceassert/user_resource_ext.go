package resourceassert

import (
	"strconv"

	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (u *UserResourceAssert) HasDisabledBool(expected bool) *UserResourceAssert {
	u.ValueSet("disabled", strconv.FormatBool(expected))
	return u
}

func (u *UserResourceAssert) HasEmptyPassword() *UserResourceAssert {
	u.ValueSet("password", "")
	return u
}

func (u *UserResourceAssert) HasNotEmptyPassword() *UserResourceAssert {
	u.ValuePresent("password")
	return u
}

func (u *UserResourceAssert) HasMustChangePasswordBool(expected bool) *UserResourceAssert {
	u.ValueSet("must_change_password", strconv.FormatBool(expected))
	return u
}

func (u *UserResourceAssert) HasDefaultSecondaryRolesOptionEnum(expected sdk.SecondaryRolesOption) *UserResourceAssert {
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
		HasDefaultSecondaryRolesOptionEnum(expectedDefaultSecondaryRoles).
		HasMinsToBypassMfaString(r.IntDefaultString).
		HasRsaPublicKeyEmpty().
		HasRsaPublicKey2Empty().
		HasCommentEmpty().
		HasDisableMfaString(r.BooleanDefault).
		HasFullyQualifiedNameString(userId.FullyQualifiedName())
}
