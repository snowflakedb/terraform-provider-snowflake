package resourceshowoutputassert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

type ExternalOauthSecurityIntegrationDescOutputAssert struct {
	*assert.ResourceAssert
}

func ExternalOauthSecurityIntegrationDescOutput(t *testing.T, name string) *ExternalOauthSecurityIntegrationDescOutputAssert {
	t.Helper()

	e := ExternalOauthSecurityIntegrationDescOutputAssert{
		ResourceAssert: assert.NewResourceAssert(name, "describe_output"),
	}
	e.AddAssertion(assert.ValueSet("describe_output.#", "1"))
	return &e
}

func (e *ExternalOauthSecurityIntegrationDescOutputAssert) HasEnabled(expected string) *ExternalOauthSecurityIntegrationDescOutputAssert {
	e.AddAssertion(assert.ResourceDescribeOutputValueSet("enabled.0.value", expected))
	return e
}

func (e *ExternalOauthSecurityIntegrationDescOutputAssert) HasExternalOauthIssuer(expected string) *ExternalOauthSecurityIntegrationDescOutputAssert {
	e.AddAssertion(assert.ResourceDescribeOutputValueSet("external_oauth_issuer.0.value", expected))
	return e
}

func (e *ExternalOauthSecurityIntegrationDescOutputAssert) HasExternalOauthJwsKeysUrl(expected string) *ExternalOauthSecurityIntegrationDescOutputAssert {
	e.AddAssertion(assert.ResourceDescribeOutputValueSet("external_oauth_jws_keys_url.0.value", expected))
	return e
}

func (e *ExternalOauthSecurityIntegrationDescOutputAssert) HasExternalOauthAnyRoleMode(expected string) *ExternalOauthSecurityIntegrationDescOutputAssert {
	e.AddAssertion(assert.ResourceDescribeOutputValueSet("external_oauth_any_role_mode.0.value", expected))
	return e
}

func (e *ExternalOauthSecurityIntegrationDescOutputAssert) HasExternalOauthRsaPublicKey(expected string) *ExternalOauthSecurityIntegrationDescOutputAssert {
	e.AddAssertion(assert.ResourceDescribeOutputValueSet("external_oauth_rsa_public_key.0.value", expected))
	return e
}

func (e *ExternalOauthSecurityIntegrationDescOutputAssert) HasExternalOauthRsaPublicKey2(expected string) *ExternalOauthSecurityIntegrationDescOutputAssert {
	e.AddAssertion(assert.ResourceDescribeOutputValueSet("external_oauth_rsa_public_key_2.0.value", expected))
	return e
}

func (e *ExternalOauthSecurityIntegrationDescOutputAssert) HasExternalOauthBlockedRolesList(expected string) *ExternalOauthSecurityIntegrationDescOutputAssert {
	e.AddAssertion(assert.ResourceDescribeOutputValueSet("external_oauth_blocked_roles_list.0.value", expected))
	return e
}

func (e *ExternalOauthSecurityIntegrationDescOutputAssert) HasExternalOauthAllowedRolesList(expected string) *ExternalOauthSecurityIntegrationDescOutputAssert {
	e.AddAssertion(assert.ResourceDescribeOutputValueSet("external_oauth_allowed_roles_list.0.value", expected))
	return e
}

func (e *ExternalOauthSecurityIntegrationDescOutputAssert) HasExternalOauthAudienceList(expected string) *ExternalOauthSecurityIntegrationDescOutputAssert {
	e.AddAssertion(assert.ResourceDescribeOutputValueSet("external_oauth_audience_list.0.value", expected))
	return e
}

func (e *ExternalOauthSecurityIntegrationDescOutputAssert) HasExternalOauthTokenUserMappingClaim(expected string) *ExternalOauthSecurityIntegrationDescOutputAssert {
	e.AddAssertion(assert.ResourceDescribeOutputValueSet("external_oauth_token_user_mapping_claim.0.value", expected))
	return e
}

func (e *ExternalOauthSecurityIntegrationDescOutputAssert) HasExternalOauthSnowflakeUserMappingAttribute(expected string) *ExternalOauthSecurityIntegrationDescOutputAssert {
	e.AddAssertion(assert.ResourceDescribeOutputValueSet("external_oauth_snowflake_user_mapping_attribute.0.value", expected))
	return e
}

func (e *ExternalOauthSecurityIntegrationDescOutputAssert) HasExternalOauthScopeDelimiter(expected string) *ExternalOauthSecurityIntegrationDescOutputAssert {
	e.AddAssertion(assert.ResourceDescribeOutputValueSet("external_oauth_scope_delimiter.0.value", expected))
	return e
}

func (e *ExternalOauthSecurityIntegrationDescOutputAssert) HasComment(expected string) *ExternalOauthSecurityIntegrationDescOutputAssert {
	e.AddAssertion(assert.ResourceDescribeOutputValueSet("comment.0.value", expected))
	return e
}
