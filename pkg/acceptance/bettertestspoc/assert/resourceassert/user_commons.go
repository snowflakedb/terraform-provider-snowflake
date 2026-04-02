package resourceassert

import (
	"fmt"
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func UserHasDefaultWorkloadIdentityOidcAssertions(issuer, subject string, audienceList ...string) []assert.ResourceAssertion {
	assertions := []assert.ResourceAssertion{
		assert.ValueSet("default_workload_identity.0.azure.#", "0"),
		assert.ValueSet("default_workload_identity.0.gcp.#", "0"),
		assert.ValueSet("default_workload_identity.0.aws.#", "0"),
		assert.ValueSet("default_workload_identity.0.oidc.#", "1"),
		assert.ValueSet("default_workload_identity.0.oidc.0.issuer", issuer),
		assert.ValueSet("default_workload_identity.0.oidc.0.subject", subject),
		assert.ValueSet("default_workload_identity.0.oidc.0.oidc_audience_list.#", strconv.Itoa(len(audienceList))),
	}
	for i, audience := range audienceList {
		assertions = append(assertions, assert.ValueSet(fmt.Sprintf("default_workload_identity.0.oidc.0.oidc_audience_list.%d", i), audience))
	}
	return assertions
}

func UserHasDefaultWorkloadIdentityAwsAssertions(arn string) []assert.ResourceAssertion {
	return []assert.ResourceAssertion{
		assert.ValueSet("default_workload_identity.0.oidc.#", "0"),
		assert.ValueSet("default_workload_identity.0.azure.#", "0"),
		assert.ValueSet("default_workload_identity.0.gcp.#", "0"),
		assert.ValueSet("default_workload_identity.0.aws.#", "1"),
		assert.ValueSet("default_workload_identity.0.aws.0.arn", arn),
	}
}

func UserHasDefaultWorkloadIdentityAzureAssertions(issuer, subject string) []assert.ResourceAssertion {
	return []assert.ResourceAssertion{
		assert.ValueSet("default_workload_identity.0.oidc.#", "0"),
		assert.ValueSet("default_workload_identity.0.gcp.#", "0"),
		assert.ValueSet("default_workload_identity.0.aws.#", "0"),
		assert.ValueSet("default_workload_identity.0.azure.#", "1"),
		assert.ValueSet("default_workload_identity.0.azure.0.issuer", issuer),
		assert.ValueSet("default_workload_identity.0.azure.0.subject", subject),
	}
}

func UserHasDefaultWorkloadIdentityGcpAssertions(subject string) []assert.ResourceAssertion {
	return []assert.ResourceAssertion{
		assert.ValueSet("default_workload_identity.0.oidc.#", "0"),
		assert.ValueSet("default_workload_identity.0.azure.#", "0"),
		assert.ValueSet("default_workload_identity.0.aws.#", "0"),
		assert.ValueSet("default_workload_identity.0.gcp.#", "1"),
		assert.ValueSet("default_workload_identity.0.gcp.0.subject", subject),
	}
}
