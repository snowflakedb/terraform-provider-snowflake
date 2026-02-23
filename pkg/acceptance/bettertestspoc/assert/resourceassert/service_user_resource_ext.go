package resourceassert

import (
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (u *ServiceUserResourceAssert) HasDisabledBool(expected bool) *ServiceUserResourceAssert {
	u.AddAssertion(assert.ValueSet("disabled", strconv.FormatBool(expected)))
	return u
}

func (u *ServiceUserResourceAssert) HasDefaultSecondaryRolesOptionEnum(expected sdk.SecondaryRolesOption) *ServiceUserResourceAssert {
	return u.HasDefaultSecondaryRolesOptionString(string(expected))
}

func (u *ServiceUserResourceAssert) HasDefaultWorkloadIdentityOidc(issuer, subject string, audienceList ...string) *ServiceUserResourceAssert {
	for _, assertion := range UserHasDefaultWorkloadIdentityOidcAssertions(issuer, subject, audienceList...) {
		u.AddAssertion(assertion)
	}
	return u
}

func (u *ServiceUserResourceAssert) HasDefaultWorkloadIdentityAws(arn string) *ServiceUserResourceAssert {
	for _, assertion := range UserHasDefaultWorkloadIdentityAwsAssertions(arn) {
		u.AddAssertion(assertion)
	}
	return u
}

func (u *ServiceUserResourceAssert) HasDefaultWorkloadIdentityAzure(issuer, subject string) *ServiceUserResourceAssert {
	for _, assertion := range UserHasDefaultWorkloadIdentityAzureAssertions(issuer, subject) {
		u.AddAssertion(assertion)
	}
	return u
}

func (u *ServiceUserResourceAssert) HasDefaultWorkloadIdentityGcp(subject string) *ServiceUserResourceAssert {
	for _, assertion := range UserHasDefaultWorkloadIdentityGcpAssertions(subject) {
		u.AddAssertion(assertion)
	}
	return u
}
