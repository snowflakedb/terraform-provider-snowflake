package resourceassert

import (
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
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

func (u *LegacyServiceUserResourceAssert) HasDefaultWorkloadIdentityOidc(issuer, subject string, audienceList ...string) *LegacyServiceUserResourceAssert {
	for _, assertion := range UserHasDefaultWorkloadIdentityOidcAssertions(issuer, subject, audienceList...) {
		u.AddAssertion(assertion)
	}
	return u
}

func (u *LegacyServiceUserResourceAssert) HasDefaultWorkloadIdentityAws(arn string) *LegacyServiceUserResourceAssert {
	for _, assertion := range UserHasDefaultWorkloadIdentityAwsAssertions(arn) {
		u.AddAssertion(assertion)
	}
	return u
}

func (u *LegacyServiceUserResourceAssert) HasDefaultWorkloadIdentityAzure(issuer, subject string) *LegacyServiceUserResourceAssert {
	for _, assertion := range UserHasDefaultWorkloadIdentityAzureAssertions(issuer, subject) {
		u.AddAssertion(assertion)
	}
	return u
}

func (u *LegacyServiceUserResourceAssert) HasDefaultWorkloadIdentityGcp(subject string) *LegacyServiceUserResourceAssert {
	for _, assertion := range UserHasDefaultWorkloadIdentityGcpAssertions(subject) {
		u.AddAssertion(assertion)
	}
	return u
}
