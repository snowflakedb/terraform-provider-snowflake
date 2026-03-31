package resourceshowoutputassert

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

type OAuthRestAuthenticationDescribeOutputAssert struct {
	*assert.ResourceAssert
	containingField string
}

func OAuthRestAuthenticationDescribeOutput(t *testing.T, name string, containingField string) *OAuthRestAuthenticationDescribeOutputAssert {
	t.Helper()

	oAuthRestAuthenticationAssert := OAuthRestAuthenticationDescribeOutputAssert{
		ResourceAssert:  assert.NewResourceAssert(name, "describe_output.0."+containingField),
		containingField: containingField,
	}
	oAuthRestAuthenticationAssert.AddAssertion(assert.ValueSet(fmt.Sprintf("describe_output.0.%s.#", containingField), "1"))
	return &oAuthRestAuthenticationAssert
}

func ImportedOAuthRestAuthenticationDescribeOutput(t *testing.T, id string, containingField string) *OAuthRestAuthenticationDescribeOutputAssert {
	t.Helper()

	oAuthRestAuthenticationAssert := OAuthRestAuthenticationDescribeOutputAssert{
		ResourceAssert:  assert.NewImportedResourceAssert(id, "describe_output.0."+containingField),
		containingField: containingField,
	}
	oAuthRestAuthenticationAssert.AddAssertion(assert.ValueSet(fmt.Sprintf("describe_output.0.%s.#", containingField), "1"))
	return &oAuthRestAuthenticationAssert
}

////////////////////////////
// Attribute value checks //
////////////////////////////

func (o *OAuthRestAuthenticationDescribeOutputAssert) HasOauthTokenUri(expected string) *OAuthRestAuthenticationDescribeOutputAssert {
	o.AddAssertion(assert.ResourceDescribeOutputValueSet(o.containingField+".0.oauth_token_uri", expected))
	return o
}

func (o *OAuthRestAuthenticationDescribeOutputAssert) HasOauthClientId(expected string) *OAuthRestAuthenticationDescribeOutputAssert {
	o.AddAssertion(assert.ResourceDescribeOutputValueSet(o.containingField+".0.oauth_client_id", expected))
	return o
}

func (o *OAuthRestAuthenticationDescribeOutputAssert) HasOauthClientSecret(expected string) *OAuthRestAuthenticationDescribeOutputAssert {
	o.AddAssertion(assert.ResourceDescribeOutputValueSet(o.containingField+".0.oauth_client_secret", expected))
	return o
}

func (o *OAuthRestAuthenticationDescribeOutputAssert) HasOauthAllowedScopes(expected ...string) *OAuthRestAuthenticationDescribeOutputAssert {
	o.AddAssertion(assert.ResourceDescribeOutputValueSet(o.containingField+".0.oauth_allowed_scopes.#", fmt.Sprintf("%d", len(expected))))
	for i, v := range expected {
		o.AddAssertion(assert.ResourceDescribeOutputValueSet(fmt.Sprintf("%s.0.oauth_allowed_scopes.%d", o.containingField, i), v))
	}
	return o
}

///////////////////////////////
// Attribute no value checks //
///////////////////////////////

func (o *OAuthRestAuthenticationDescribeOutputAssert) HasNoOauthTokenUri() *OAuthRestAuthenticationDescribeOutputAssert {
	o.AddAssertion(assert.ResourceDescribeOutputValueNotSet(o.containingField + ".0.oauth_token_uri"))
	return o
}

func (o *OAuthRestAuthenticationDescribeOutputAssert) HasNoOauthClientId() *OAuthRestAuthenticationDescribeOutputAssert {
	o.AddAssertion(assert.ResourceDescribeOutputValueNotSet(o.containingField + ".0.oauth_client_id"))
	return o
}

func (o *OAuthRestAuthenticationDescribeOutputAssert) HasNoOauthClientSecret() *OAuthRestAuthenticationDescribeOutputAssert {
	o.AddAssertion(assert.ResourceDescribeOutputValueNotSet(o.containingField + ".0.oauth_client_secret"))
	return o
}

func (o *OAuthRestAuthenticationDescribeOutputAssert) HasNoOauthAllowedScopes() *OAuthRestAuthenticationDescribeOutputAssert {
	o.AddAssertion(assert.ResourceDescribeOutputValueSet(o.containingField+".0.oauth_allowed_scopes.#", "0"))
	return o
}
