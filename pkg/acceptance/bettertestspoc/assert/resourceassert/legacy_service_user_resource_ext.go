package resourceassert

import (
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (u *LegacyServiceUserResourceAssert) HasDisabledBool(expected bool) *LegacyServiceUserResourceAssert {
	u.ValueSet("disabled", strconv.FormatBool(expected))
	return u
}

func (u *LegacyServiceUserResourceAssert) HasMustChangePasswordBool(expected bool) *LegacyServiceUserResourceAssert {
	u.ValueSet("must_change_password", strconv.FormatBool(expected))
	return u
}

func (u *LegacyServiceUserResourceAssert) HasDefaultSecondaryRolesOptionEnum(expected sdk.SecondaryRolesOption) *LegacyServiceUserResourceAssert {
	return u.HasDefaultSecondaryRolesOptionString(string(expected))
}

func (u *LegacyServiceUserResourceAssert) HasDefaultWorkloadIdentityOidc(issuer, subject string, audienceList ...string) *LegacyServiceUserResourceAssert {
	userApplyDefaultWorkloadIdentityOidcChecks(u.ResourceAssert, issuer, subject, audienceList...)
	return u
}

func (u *LegacyServiceUserResourceAssert) HasDefaultWorkloadIdentityAws(arn string, issuer ...string) *LegacyServiceUserResourceAssert {
	userApplyDefaultWorkloadIdentityAwsChecks(u.ResourceAssert, arn, issuer...)
	return u
}

func (u *LegacyServiceUserResourceAssert) HasDefaultWorkloadIdentityAzure(issuer, subject string) *LegacyServiceUserResourceAssert {
	userApplyDefaultWorkloadIdentityAzureChecks(u.ResourceAssert, issuer, subject)
	return u
}

func (u *LegacyServiceUserResourceAssert) HasDefaultWorkloadIdentityGcp(subject string) *LegacyServiceUserResourceAssert {
	userApplyDefaultWorkloadIdentityGcpChecks(u.ResourceAssert, subject)
	return u
}
