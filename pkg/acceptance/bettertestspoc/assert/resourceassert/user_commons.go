package resourceassert

import (
	"fmt"
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func userApplyDefaultWorkloadIdentityOidcChecks(e *assert.ResourceAssert, issuer, subject string, audienceList ...string) {
	e.ValueSet("default_workload_identity.0.azure.#", "0")
	e.ValueSet("default_workload_identity.0.gcp.#", "0")
	e.ValueSet("default_workload_identity.0.aws.#", "0")
	e.ValueSet("default_workload_identity.0.oidc.#", "1")
	e.ValueSet("default_workload_identity.0.oidc.0.issuer", issuer)
	e.ValueSet("default_workload_identity.0.oidc.0.subject", subject)
	e.ValueSet("default_workload_identity.0.oidc.0.oidc_audience_list.#", strconv.Itoa(len(audienceList)))
	for i, audience := range audienceList {
		e.ValueSet(fmt.Sprintf("default_workload_identity.0.oidc.0.oidc_audience_list.%d", i), audience)
	}
}

func userApplyDefaultWorkloadIdentityAwsChecks(e *assert.ResourceAssert, arn string, issuer ...string) {
	e.ValueSet("default_workload_identity.0.oidc.#", "0")
	e.ValueSet("default_workload_identity.0.azure.#", "0")
	e.ValueSet("default_workload_identity.0.gcp.#", "0")
	e.ValueSet("default_workload_identity.0.aws.#", "1")
	e.ValueSet("default_workload_identity.0.aws.0.arn", arn)
	if len(issuer) > 0 {
		e.ValueSet("default_workload_identity.0.aws.0.issuer", issuer[0])
	}
}

func userApplyDefaultWorkloadIdentityAzureChecks(e *assert.ResourceAssert, issuer, subject string) {
	e.ValueSet("default_workload_identity.0.oidc.#", "0")
	e.ValueSet("default_workload_identity.0.gcp.#", "0")
	e.ValueSet("default_workload_identity.0.aws.#", "0")
	e.ValueSet("default_workload_identity.0.azure.#", "1")
	e.ValueSet("default_workload_identity.0.azure.0.issuer", issuer)
	e.ValueSet("default_workload_identity.0.azure.0.subject", subject)
}

func userApplyDefaultWorkloadIdentityGcpChecks(e *assert.ResourceAssert, subject string) {
	e.ValueSet("default_workload_identity.0.oidc.#", "0")
	e.ValueSet("default_workload_identity.0.azure.#", "0")
	e.ValueSet("default_workload_identity.0.aws.#", "0")
	e.ValueSet("default_workload_identity.0.gcp.#", "1")
	e.ValueSet("default_workload_identity.0.gcp.0.subject", subject)
}
