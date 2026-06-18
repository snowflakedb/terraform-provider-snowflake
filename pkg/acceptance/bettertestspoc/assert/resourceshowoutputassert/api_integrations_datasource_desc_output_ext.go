package resourceshowoutputassert

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

type ApiIntegrationAllDetailsDescribeOutputAssert struct {
	*assert.ResourceAssert
}

func ApiIntegrationsDatasourceDescribeOutput(t *testing.T, datasourceReference string) *ApiIntegrationAllDetailsDescribeOutputAssert {
	t.Helper()
	s := ApiIntegrationAllDetailsDescribeOutputAssert{
		ResourceAssert: assert.NewDatasourceAssert(datasourceReference, "describe_output", "api_integrations.0."),
	}
	s.AddAssertion(assert.ValueSet("describe_output.#", "1"))
	return &s
}

func (a *ApiIntegrationAllDetailsDescribeOutputAssert) HasEnabled(expected bool) *ApiIntegrationAllDetailsDescribeOutputAssert {
	a.AddAssertion(assert.ResourceDescribeOutputBoolValueSet("enabled", expected))
	return a
}

func (a *ApiIntegrationAllDetailsDescribeOutputAssert) HasApiKey(expected string) *ApiIntegrationAllDetailsDescribeOutputAssert {
	a.AddAssertion(assert.ResourceDescribeOutputValueSet("api_key", expected))
	return a
}

func (a *ApiIntegrationAllDetailsDescribeOutputAssert) HasApiProvider(expected string) *ApiIntegrationAllDetailsDescribeOutputAssert {
	a.AddAssertion(assert.ResourceDescribeOutputValueSet("api_provider", expected))
	return a
}

func (a *ApiIntegrationAllDetailsDescribeOutputAssert) HasComment(expected string) *ApiIntegrationAllDetailsDescribeOutputAssert {
	a.AddAssertion(assert.ResourceDescribeOutputValueSet("comment", expected))
	return a
}

func (a *ApiIntegrationAllDetailsDescribeOutputAssert) HasAllowedPrefixes(expected ...string) *ApiIntegrationAllDetailsDescribeOutputAssert {
	a.AddAssertion(assert.ResourceDescribeOutputValueSet("allowed_prefixes.#", fmt.Sprintf("%d", len(expected))))
	return a
}

func (a *ApiIntegrationAllDetailsDescribeOutputAssert) HasBlockedPrefixes(expected ...string) *ApiIntegrationAllDetailsDescribeOutputAssert {
	a.AddAssertion(assert.ResourceDescribeOutputValueSet("blocked_prefixes.#", fmt.Sprintf("%d", len(expected))))
	return a
}

// AWS-specific

func (a *ApiIntegrationAllDetailsDescribeOutputAssert) HasApiAwsRoleArn(expected string) *ApiIntegrationAllDetailsDescribeOutputAssert {
	a.AddAssertion(assert.ResourceDescribeOutputValueSet("api_aws_role_arn", expected))
	return a
}

func (a *ApiIntegrationAllDetailsDescribeOutputAssert) HasApiAwsIamUserArnNotEmpty() *ApiIntegrationAllDetailsDescribeOutputAssert {
	a.AddAssertion(assert.ResourceDescribeOutputValuePresent("api_aws_iam_user_arn"))
	return a
}

func (a *ApiIntegrationAllDetailsDescribeOutputAssert) HasApiAwsExternalIdNotEmpty() *ApiIntegrationAllDetailsDescribeOutputAssert {
	a.AddAssertion(assert.ResourceDescribeOutputValuePresent("api_aws_external_id"))
	return a
}

// Azure-specific

func (a *ApiIntegrationAllDetailsDescribeOutputAssert) HasAzureTenantId(expected string) *ApiIntegrationAllDetailsDescribeOutputAssert {
	a.AddAssertion(assert.ResourceDescribeOutputValueSet("azure_tenant_id", expected))
	return a
}

func (a *ApiIntegrationAllDetailsDescribeOutputAssert) HasAzureAdApplicationId(expected string) *ApiIntegrationAllDetailsDescribeOutputAssert {
	a.AddAssertion(assert.ResourceDescribeOutputValueSet("azure_ad_application_id", expected))
	return a
}

// Google-specific

func (a *ApiIntegrationAllDetailsDescribeOutputAssert) HasGoogleAudience(expected string) *ApiIntegrationAllDetailsDescribeOutputAssert {
	a.AddAssertion(assert.ResourceDescribeOutputValueSet("google_audience", expected))
	return a
}

func (a *ApiIntegrationAllDetailsDescribeOutputAssert) HasGoogleApiServiceAccountNotEmpty() *ApiIntegrationAllDetailsDescribeOutputAssert {
	a.AddAssertion(assert.ResourceDescribeOutputValuePresent("google_api_service_account"))
	return a
}

// Git HTTPS / External MCP

func (a *ApiIntegrationAllDetailsDescribeOutputAssert) HasUserAuthType(expected string) *ApiIntegrationAllDetailsDescribeOutputAssert {
	a.AddAssertion(assert.ResourceDescribeOutputValueSet("user_auth_type", expected))
	return a
}

func (a *ApiIntegrationAllDetailsDescribeOutputAssert) HasAllowedAuthenticationSecrets(expected string) *ApiIntegrationAllDetailsDescribeOutputAssert {
	a.AddAssertion(assert.ResourceDescribeOutputValueSet("allowed_authentication_secrets", expected))
	return a
}
