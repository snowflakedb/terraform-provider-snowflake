// Code generated by assertions generator; DO NOT EDIT.

package resourceassert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

type OauthIntegrationForPartnerApplicationsResourceAssert struct {
	*assert.ResourceAssert
}

func OauthIntegrationForPartnerApplicationsResource(t *testing.T, name string) *OauthIntegrationForPartnerApplicationsResourceAssert {
	t.Helper()

	return &OauthIntegrationForPartnerApplicationsResourceAssert{
		ResourceAssert: assert.NewResourceAssert(name, "resource"),
	}
}

func ImportedOauthIntegrationForPartnerApplicationsResource(t *testing.T, id string) *OauthIntegrationForPartnerApplicationsResourceAssert {
	t.Helper()

	return &OauthIntegrationForPartnerApplicationsResourceAssert{
		ResourceAssert: assert.NewImportedResourceAssert(id, "imported resource"),
	}
}

///////////////////////////////////
// Attribute value string checks //
///////////////////////////////////

func (o *OauthIntegrationForPartnerApplicationsResourceAssert) HasBlockedRolesListString(expected string) *OauthIntegrationForPartnerApplicationsResourceAssert {
	o.AddAssertion(assert.ValueSet("blocked_roles_list", expected))
	return o
}

func (o *OauthIntegrationForPartnerApplicationsResourceAssert) HasCommentString(expected string) *OauthIntegrationForPartnerApplicationsResourceAssert {
	o.AddAssertion(assert.ValueSet("comment", expected))
	return o
}

func (o *OauthIntegrationForPartnerApplicationsResourceAssert) HasEnabledString(expected string) *OauthIntegrationForPartnerApplicationsResourceAssert {
	o.AddAssertion(assert.ValueSet("enabled", expected))
	return o
}

func (o *OauthIntegrationForPartnerApplicationsResourceAssert) HasFullyQualifiedNameString(expected string) *OauthIntegrationForPartnerApplicationsResourceAssert {
	o.AddAssertion(assert.ValueSet("fully_qualified_name", expected))
	return o
}

func (o *OauthIntegrationForPartnerApplicationsResourceAssert) HasNameString(expected string) *OauthIntegrationForPartnerApplicationsResourceAssert {
	o.AddAssertion(assert.ValueSet("name", expected))
	return o
}

func (o *OauthIntegrationForPartnerApplicationsResourceAssert) HasOauthClientString(expected string) *OauthIntegrationForPartnerApplicationsResourceAssert {
	o.AddAssertion(assert.ValueSet("oauth_client", expected))
	return o
}

func (o *OauthIntegrationForPartnerApplicationsResourceAssert) HasOauthIssueRefreshTokensString(expected string) *OauthIntegrationForPartnerApplicationsResourceAssert {
	o.AddAssertion(assert.ValueSet("oauth_issue_refresh_tokens", expected))
	return o
}

func (o *OauthIntegrationForPartnerApplicationsResourceAssert) HasOauthRedirectUriString(expected string) *OauthIntegrationForPartnerApplicationsResourceAssert {
	o.AddAssertion(assert.ValueSet("oauth_redirect_uri", expected))
	return o
}

func (o *OauthIntegrationForPartnerApplicationsResourceAssert) HasOauthRefreshTokenValidityString(expected string) *OauthIntegrationForPartnerApplicationsResourceAssert {
	o.AddAssertion(assert.ValueSet("oauth_refresh_token_validity", expected))
	return o
}

func (o *OauthIntegrationForPartnerApplicationsResourceAssert) HasOauthUseSecondaryRolesString(expected string) *OauthIntegrationForPartnerApplicationsResourceAssert {
	o.AddAssertion(assert.ValueSet("oauth_use_secondary_roles", expected))
	return o
}

func (o *OauthIntegrationForPartnerApplicationsResourceAssert) HasRelatedParametersString(expected string) *OauthIntegrationForPartnerApplicationsResourceAssert {
	o.AddAssertion(assert.ValueSet("related_parameters", expected))
	return o
}

///////////////////////////////
// Attribute no value checks //
///////////////////////////////

func (o *OauthIntegrationForPartnerApplicationsResourceAssert) HasNoBlockedRolesList() *OauthIntegrationForPartnerApplicationsResourceAssert {
	o.AddAssertion(assert.ValueSet("blocked_roles_list.#", "0"))
	return o
}

func (o *OauthIntegrationForPartnerApplicationsResourceAssert) HasNoComment() *OauthIntegrationForPartnerApplicationsResourceAssert {
	o.AddAssertion(assert.ValueNotSet("comment"))
	return o
}

func (o *OauthIntegrationForPartnerApplicationsResourceAssert) HasNoEnabled() *OauthIntegrationForPartnerApplicationsResourceAssert {
	o.AddAssertion(assert.ValueNotSet("enabled"))
	return o
}

func (o *OauthIntegrationForPartnerApplicationsResourceAssert) HasNoFullyQualifiedName() *OauthIntegrationForPartnerApplicationsResourceAssert {
	o.AddAssertion(assert.ValueNotSet("fully_qualified_name"))
	return o
}

func (o *OauthIntegrationForPartnerApplicationsResourceAssert) HasNoName() *OauthIntegrationForPartnerApplicationsResourceAssert {
	o.AddAssertion(assert.ValueNotSet("name"))
	return o
}

func (o *OauthIntegrationForPartnerApplicationsResourceAssert) HasNoOauthClient() *OauthIntegrationForPartnerApplicationsResourceAssert {
	o.AddAssertion(assert.ValueNotSet("oauth_client"))
	return o
}

func (o *OauthIntegrationForPartnerApplicationsResourceAssert) HasNoOauthIssueRefreshTokens() *OauthIntegrationForPartnerApplicationsResourceAssert {
	o.AddAssertion(assert.ValueNotSet("oauth_issue_refresh_tokens"))
	return o
}

func (o *OauthIntegrationForPartnerApplicationsResourceAssert) HasNoOauthRedirectUri() *OauthIntegrationForPartnerApplicationsResourceAssert {
	o.AddAssertion(assert.ValueNotSet("oauth_redirect_uri"))
	return o
}

func (o *OauthIntegrationForPartnerApplicationsResourceAssert) HasNoOauthRefreshTokenValidity() *OauthIntegrationForPartnerApplicationsResourceAssert {
	o.AddAssertion(assert.ValueNotSet("oauth_refresh_token_validity"))
	return o
}

func (o *OauthIntegrationForPartnerApplicationsResourceAssert) HasNoOauthUseSecondaryRoles() *OauthIntegrationForPartnerApplicationsResourceAssert {
	o.AddAssertion(assert.ValueNotSet("oauth_use_secondary_roles"))
	return o
}

func (o *OauthIntegrationForPartnerApplicationsResourceAssert) HasNoRelatedParameters() *OauthIntegrationForPartnerApplicationsResourceAssert {
	o.AddAssertion(assert.ValueSet("related_parameters.#", "0"))
	return o
}

////////////////////////////
// Attribute empty checks //
////////////////////////////

func (o *OauthIntegrationForPartnerApplicationsResourceAssert) HasCommentEmpty() *OauthIntegrationForPartnerApplicationsResourceAssert {
	o.AddAssertion(assert.ValueSet("comment", ""))
	return o
}

func (o *OauthIntegrationForPartnerApplicationsResourceAssert) HasEnabledEmpty() *OauthIntegrationForPartnerApplicationsResourceAssert {
	o.AddAssertion(assert.ValueSet("enabled", ""))
	return o
}

func (o *OauthIntegrationForPartnerApplicationsResourceAssert) HasFullyQualifiedNameEmpty() *OauthIntegrationForPartnerApplicationsResourceAssert {
	o.AddAssertion(assert.ValueSet("fully_qualified_name", ""))
	return o
}

func (o *OauthIntegrationForPartnerApplicationsResourceAssert) HasOauthIssueRefreshTokensEmpty() *OauthIntegrationForPartnerApplicationsResourceAssert {
	o.AddAssertion(assert.ValueSet("oauth_issue_refresh_tokens", ""))
	return o
}

func (o *OauthIntegrationForPartnerApplicationsResourceAssert) HasOauthRedirectUriEmpty() *OauthIntegrationForPartnerApplicationsResourceAssert {
	o.AddAssertion(assert.ValueSet("oauth_redirect_uri", ""))
	return o
}

func (o *OauthIntegrationForPartnerApplicationsResourceAssert) HasOauthUseSecondaryRolesEmpty() *OauthIntegrationForPartnerApplicationsResourceAssert {
	o.AddAssertion(assert.ValueSet("oauth_use_secondary_roles", ""))
	return o
}

///////////////////////////////
// Attribute presence checks //
///////////////////////////////

func (o *OauthIntegrationForPartnerApplicationsResourceAssert) HasBlockedRolesListNotEmpty() *OauthIntegrationForPartnerApplicationsResourceAssert {
	o.AddAssertion(assert.ValuePresent("blocked_roles_list"))
	return o
}

func (o *OauthIntegrationForPartnerApplicationsResourceAssert) HasCommentNotEmpty() *OauthIntegrationForPartnerApplicationsResourceAssert {
	o.AddAssertion(assert.ValuePresent("comment"))
	return o
}

func (o *OauthIntegrationForPartnerApplicationsResourceAssert) HasEnabledNotEmpty() *OauthIntegrationForPartnerApplicationsResourceAssert {
	o.AddAssertion(assert.ValuePresent("enabled"))
	return o
}

func (o *OauthIntegrationForPartnerApplicationsResourceAssert) HasFullyQualifiedNameNotEmpty() *OauthIntegrationForPartnerApplicationsResourceAssert {
	o.AddAssertion(assert.ValuePresent("fully_qualified_name"))
	return o
}

func (o *OauthIntegrationForPartnerApplicationsResourceAssert) HasNameNotEmpty() *OauthIntegrationForPartnerApplicationsResourceAssert {
	o.AddAssertion(assert.ValuePresent("name"))
	return o
}

func (o *OauthIntegrationForPartnerApplicationsResourceAssert) HasOauthClientNotEmpty() *OauthIntegrationForPartnerApplicationsResourceAssert {
	o.AddAssertion(assert.ValuePresent("oauth_client"))
	return o
}

func (o *OauthIntegrationForPartnerApplicationsResourceAssert) HasOauthIssueRefreshTokensNotEmpty() *OauthIntegrationForPartnerApplicationsResourceAssert {
	o.AddAssertion(assert.ValuePresent("oauth_issue_refresh_tokens"))
	return o
}

func (o *OauthIntegrationForPartnerApplicationsResourceAssert) HasOauthRedirectUriNotEmpty() *OauthIntegrationForPartnerApplicationsResourceAssert {
	o.AddAssertion(assert.ValuePresent("oauth_redirect_uri"))
	return o
}

func (o *OauthIntegrationForPartnerApplicationsResourceAssert) HasOauthRefreshTokenValidityNotEmpty() *OauthIntegrationForPartnerApplicationsResourceAssert {
	o.AddAssertion(assert.ValuePresent("oauth_refresh_token_validity"))
	return o
}

func (o *OauthIntegrationForPartnerApplicationsResourceAssert) HasOauthUseSecondaryRolesNotEmpty() *OauthIntegrationForPartnerApplicationsResourceAssert {
	o.AddAssertion(assert.ValuePresent("oauth_use_secondary_roles"))
	return o
}

func (o *OauthIntegrationForPartnerApplicationsResourceAssert) HasRelatedParametersNotEmpty() *OauthIntegrationForPartnerApplicationsResourceAssert {
	o.AddAssertion(assert.ValuePresent("related_parameters"))
	return o
}
