package resourceassert

import (
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (u *ServiceUserResourceAssert) HasDisabledBool(expected bool) *ServiceUserResourceAssert {
	u.ValueSet("disabled", strconv.FormatBool(expected))
	return u
}

func (u *ServiceUserResourceAssert) HasDefaultSecondaryRolesOptionEnum(expected sdk.SecondaryRolesOption) *ServiceUserResourceAssert {
	return u.HasDefaultSecondaryRolesOptionString(string(expected))
}

func (u *ServiceUserResourceAssert) HasDefaultWorkloadIdentityOidc(issuer, subject string, audienceList ...string) *ServiceUserResourceAssert {
	userApplyDefaultWorkloadIdentityOidcChecks(u.ResourceAssert, issuer, subject, audienceList...)
	return u
}

func (u *ServiceUserResourceAssert) HasDefaultWorkloadIdentityAws(arn string) *ServiceUserResourceAssert {
	userApplyDefaultWorkloadIdentityAwsChecks(u.ResourceAssert, arn)
	return u
}

func (u *ServiceUserResourceAssert) HasDefaultWorkloadIdentityAzure(issuer, subject string) *ServiceUserResourceAssert {
	userApplyDefaultWorkloadIdentityAzureChecks(u.ResourceAssert, issuer, subject)
	return u
}

func (u *ServiceUserResourceAssert) HasDefaultWorkloadIdentityGcp(subject string) *ServiceUserResourceAssert {
	userApplyDefaultWorkloadIdentityGcpChecks(u.ResourceAssert, subject)
	return u
}
