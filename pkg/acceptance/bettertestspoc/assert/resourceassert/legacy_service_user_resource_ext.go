package resourceassert

import (
	"fmt"
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
	u.AddAssertion(assert.ValueSet("default_workload_identity.0.aws.#", "0"))
	u.AddAssertion(assert.ValueSet("default_workload_identity.0.azure.#", "0"))
	u.AddAssertion(assert.ValueSet("default_workload_identity.0.gcp.#", "0"))
	u.AddAssertion(assert.ValueSet("default_workload_identity.0.oidc.0.issuer", issuer))
	u.AddAssertion(assert.ValueSet("default_workload_identity.0.oidc.0.subject", subject))
	u.AddAssertion(assert.ValueSet("default_workload_identity.0.oidc.0.oidc_audience_list.#", strconv.Itoa(len(audienceList))))
	for i, audience := range audienceList {
		u.AddAssertion(assert.ValueSet(fmt.Sprintf("default_workload_identity.0.oidc.0.oidc_audience_list.%d", i), audience))
	}
	return u
}

func (u *LegacyServiceUserResourceAssert) HasDefaultWorkloadIdentityAws(arn string) *LegacyServiceUserResourceAssert {
	u.AddAssertion(assert.ValueSet("default_workload_identity.0.azure.#", "0"))
	u.AddAssertion(assert.ValueSet("default_workload_identity.0.gcp.#", "0"))
	u.AddAssertion(assert.ValueSet("default_workload_identity.0.oidc.#", "0"))
	u.AddAssertion(assert.ValueSet("default_workload_identity.0.aws.0.arn", arn))
	return u
}

func (u *LegacyServiceUserResourceAssert) HasDefaultWorkloadIdentityAzure(issuer, subject string) *LegacyServiceUserResourceAssert {
	u.AddAssertion(assert.ValueSet("default_workload_identity.0.aws.#", "0"))
	u.AddAssertion(assert.ValueSet("default_workload_identity.0.gcp.#", "0"))
	u.AddAssertion(assert.ValueSet("default_workload_identity.0.oidc.#", "0"))
	u.AddAssertion(assert.ValueSet("default_workload_identity.0.azure.0.issuer", issuer))
	u.AddAssertion(assert.ValueSet("default_workload_identity.0.azure.0.subject", subject))
	return u
}

func (u *LegacyServiceUserResourceAssert) HasDefaultWorkloadIdentityGcp(subject string) *LegacyServiceUserResourceAssert {
	u.AddAssertion(assert.ValueSet("default_workload_identity.0.aws.#", "0"))
	u.AddAssertion(assert.ValueSet("default_workload_identity.0.azure.#", "0"))
	u.AddAssertion(assert.ValueSet("default_workload_identity.0.oidc.#", "0"))
	u.AddAssertion(assert.ValueSet("default_workload_identity.0.gcp.0.subject", subject))
	return u
}
