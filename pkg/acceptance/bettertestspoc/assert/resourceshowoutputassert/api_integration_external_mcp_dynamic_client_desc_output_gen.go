// file edited manually; all ShowOutput changed to DescribeOutput; only DynamicClient-relevant fields included

package resourceshowoutputassert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type ApiIntegrationExternalMcpDynamicClientDescribeOutputAssert struct {
	*assert.ResourceAssert
}

func ApiIntegrationExternalMcpDynamicClientDescribeOutput(t *testing.T, name string) *ApiIntegrationExternalMcpDynamicClientDescribeOutputAssert {
	t.Helper()

	a := ApiIntegrationExternalMcpDynamicClientDescribeOutputAssert{
		ResourceAssert: assert.NewResourceAssert(name, "describe_output"),
	}
	a.AddAssertion(assert.ValueSet("describe_output.#", "1"))
	return &a
}

func ImportedApiIntegrationExternalMcpDynamicClientDescribeOutput(t *testing.T, id string) *ApiIntegrationExternalMcpDynamicClientDescribeOutputAssert {
	t.Helper()

	a := ApiIntegrationExternalMcpDynamicClientDescribeOutputAssert{
		ResourceAssert: assert.NewImportedResourceAssert(id, "describe_output"),
	}
	a.AddAssertion(assert.ValueSet("describe_output.#", "1"))
	return &a
}

////////////////////////////
// Attribute value checks //
////////////////////////////

func (a *ApiIntegrationExternalMcpDynamicClientDescribeOutputAssert) HasId(expected sdk.AccountObjectIdentifier) *ApiIntegrationExternalMcpDynamicClientDescribeOutputAssert {
	a.AddAssertion(assert.ResourceDescribeOutputStringUnderlyingValueSet("id", expected.Name()))
	return a
}

func (a *ApiIntegrationExternalMcpDynamicClientDescribeOutputAssert) HasEnabled(expected bool) *ApiIntegrationExternalMcpDynamicClientDescribeOutputAssert {
	a.AddAssertion(assert.ResourceDescribeOutputBoolValueSet("enabled", expected))
	return a
}

func (a *ApiIntegrationExternalMcpDynamicClientDescribeOutputAssert) HasApiProvider(expected string) *ApiIntegrationExternalMcpDynamicClientDescribeOutputAssert {
	a.AddAssertion(assert.ResourceDescribeOutputValueSet("api_provider", expected))
	return a
}

func (a *ApiIntegrationExternalMcpDynamicClientDescribeOutputAssert) HasUserAuthType(expected string) *ApiIntegrationExternalMcpDynamicClientDescribeOutputAssert {
	a.AddAssertion(assert.ResourceDescribeOutputValueSet("user_auth_type", expected))
	return a
}

func (a *ApiIntegrationExternalMcpDynamicClientDescribeOutputAssert) HasOauthResourceUrl(expected string) *ApiIntegrationExternalMcpDynamicClientDescribeOutputAssert {
	a.AddAssertion(assert.ResourceDescribeOutputValueSet("oauth_resource_url", expected))
	return a
}

func (a *ApiIntegrationExternalMcpDynamicClientDescribeOutputAssert) HasComment(expected string) *ApiIntegrationExternalMcpDynamicClientDescribeOutputAssert {
	a.AddAssertion(assert.ResourceDescribeOutputValueSet("comment", expected))
	return a
}

///////////////////////////////
// Attribute no value checks //
///////////////////////////////

func (a *ApiIntegrationExternalMcpDynamicClientDescribeOutputAssert) HasNoBlockedPrefixes() *ApiIntegrationExternalMcpDynamicClientDescribeOutputAssert {
	a.AddAssertion(assert.ResourceDescribeOutputValueSet("blocked_prefixes.#", "0"))
	return a
}

func (a *ApiIntegrationExternalMcpDynamicClientDescribeOutputAssert) HasNoComment() *ApiIntegrationExternalMcpDynamicClientDescribeOutputAssert {
	a.AddAssertion(assert.ResourceDescribeOutputValueNotSet("comment"))
	return a
}
