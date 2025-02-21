// Code generated by assertions generator; DO NOT EDIT.

package resourceassert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

type OauthIntegrationForCustomClientsResourceAssert struct {
	*assert.ResourceAssert
}

func OauthIntegrationForCustomClientsResource(t *testing.T, name string) *OauthIntegrationForCustomClientsResourceAssert {
	t.Helper()

	return &OauthIntegrationForCustomClientsResourceAssert{
		ResourceAssert: assert.NewResourceAssert(name, "resource"),
	}
}

func ImportedOauthIntegrationForCustomClientsResource(t *testing.T, id string) *OauthIntegrationForCustomClientsResourceAssert {
	t.Helper()

	return &OauthIntegrationForCustomClientsResourceAssert{
		ResourceAssert: assert.NewImportedResourceAssert(id, "imported resource"),
	}
}

///////////////////////////////////
// Attribute value string checks //
///////////////////////////////////

func (o *OauthIntegrationForCustomClientsResourceAssert) HasBlockedRolesListString(expected string) *OauthIntegrationForCustomClientsResourceAssert {
	o.AddAssertion(assert.ValueSet("blocked_roles_list", expected))
	return o
}

func (o *OauthIntegrationForCustomClientsResourceAssert) HasCommentString(expected string) *OauthIntegrationForCustomClientsResourceAssert {
	o.AddAssertion(assert.ValueSet("comment", expected))
	return o
}

func (o *OauthIntegrationForCustomClientsResourceAssert) HasEnabledString(expected string) *OauthIntegrationForCustomClientsResourceAssert {
	o.AddAssertion(assert.ValueSet("enabled", expected))
	return o
}

func (o *OauthIntegrationForCustomClientsResourceAssert) HasFullyQualifiedNameString(expected string) *OauthIntegrationForCustomClientsResourceAssert {
	o.AddAssertion(assert.ValueSet("fully_qualified_name", expected))
	return o
}

func (o *OauthIntegrationForCustomClientsResourceAssert) HasNameString(expected string) *OauthIntegrationForCustomClientsResourceAssert {
	o.AddAssertion(assert.ValueSet("name", expected))
	return o
}

func (o *OauthIntegrationForCustomClientsResourceAssert) HasNetworkPolicyString(expected string) *OauthIntegrationForCustomClientsResourceAssert {
	o.AddAssertion(assert.ValueSet("network_policy", expected))
	return o
}

func (o *OauthIntegrationForCustomClientsResourceAssert) HasOauthAllowNonTlsRedirectUriString(expected string) *OauthIntegrationForCustomClientsResourceAssert {
	o.AddAssertion(assert.ValueSet("oauth_allow_non_tls_redirect_uri", expected))
	return o
}

func (o *OauthIntegrationForCustomClientsResourceAssert) HasOauthClientRsaPublicKeyString(expected string) *OauthIntegrationForCustomClientsResourceAssert {
	o.AddAssertion(assert.ValueSet("oauth_client_rsa_public_key", expected))
	return o
}

func (o *OauthIntegrationForCustomClientsResourceAssert) HasOauthClientRsaPublicKey2String(expected string) *OauthIntegrationForCustomClientsResourceAssert {
	o.AddAssertion(assert.ValueSet("oauth_client_rsa_public_key_2", expected))
	return o
}

func (o *OauthIntegrationForCustomClientsResourceAssert) HasOauthClientTypeString(expected string) *OauthIntegrationForCustomClientsResourceAssert {
	o.AddAssertion(assert.ValueSet("oauth_client_type", expected))
	return o
}

func (o *OauthIntegrationForCustomClientsResourceAssert) HasOauthEnforcePkceString(expected string) *OauthIntegrationForCustomClientsResourceAssert {
	o.AddAssertion(assert.ValueSet("oauth_enforce_pkce", expected))
	return o
}

func (o *OauthIntegrationForCustomClientsResourceAssert) HasOauthIssueRefreshTokensString(expected string) *OauthIntegrationForCustomClientsResourceAssert {
	o.AddAssertion(assert.ValueSet("oauth_issue_refresh_tokens", expected))
	return o
}

func (o *OauthIntegrationForCustomClientsResourceAssert) HasOauthRedirectUriString(expected string) *OauthIntegrationForCustomClientsResourceAssert {
	o.AddAssertion(assert.ValueSet("oauth_redirect_uri", expected))
	return o
}

func (o *OauthIntegrationForCustomClientsResourceAssert) HasOauthRefreshTokenValidityString(expected string) *OauthIntegrationForCustomClientsResourceAssert {
	o.AddAssertion(assert.ValueSet("oauth_refresh_token_validity", expected))
	return o
}

func (o *OauthIntegrationForCustomClientsResourceAssert) HasOauthUseSecondaryRolesString(expected string) *OauthIntegrationForCustomClientsResourceAssert {
	o.AddAssertion(assert.ValueSet("oauth_use_secondary_roles", expected))
	return o
}

func (o *OauthIntegrationForCustomClientsResourceAssert) HasPreAuthorizedRolesListString(expected string) *OauthIntegrationForCustomClientsResourceAssert {
	o.AddAssertion(assert.ValueSet("pre_authorized_roles_list", expected))
	return o
}

func (o *OauthIntegrationForCustomClientsResourceAssert) HasRelatedParametersString(expected string) *OauthIntegrationForCustomClientsResourceAssert {
	o.AddAssertion(assert.ValueSet("related_parameters", expected))
	return o
}

///////////////////////////////
// Attribute no value checks //
///////////////////////////////

func (o *OauthIntegrationForCustomClientsResourceAssert) HasNoBlockedRolesList() *OauthIntegrationForCustomClientsResourceAssert {
	o.AddAssertion(assert.ValueSet("blocked_roles_list.#", "0"))
	return o
}

func (o *OauthIntegrationForCustomClientsResourceAssert) HasNoComment() *OauthIntegrationForCustomClientsResourceAssert {
	o.AddAssertion(assert.ValueNotSet("comment"))
	return o
}

func (o *OauthIntegrationForCustomClientsResourceAssert) HasNoEnabled() *OauthIntegrationForCustomClientsResourceAssert {
	o.AddAssertion(assert.ValueNotSet("enabled"))
	return o
}

func (o *OauthIntegrationForCustomClientsResourceAssert) HasNoFullyQualifiedName() *OauthIntegrationForCustomClientsResourceAssert {
	o.AddAssertion(assert.ValueNotSet("fully_qualified_name"))
	return o
}

func (o *OauthIntegrationForCustomClientsResourceAssert) HasNoName() *OauthIntegrationForCustomClientsResourceAssert {
	o.AddAssertion(assert.ValueNotSet("name"))
	return o
}

func (o *OauthIntegrationForCustomClientsResourceAssert) HasNoNetworkPolicy() *OauthIntegrationForCustomClientsResourceAssert {
	o.AddAssertion(assert.ValueNotSet("network_policy"))
	return o
}

func (o *OauthIntegrationForCustomClientsResourceAssert) HasNoOauthAllowNonTlsRedirectUri() *OauthIntegrationForCustomClientsResourceAssert {
	o.AddAssertion(assert.ValueNotSet("oauth_allow_non_tls_redirect_uri"))
	return o
}

func (o *OauthIntegrationForCustomClientsResourceAssert) HasNoOauthClientRsaPublicKey() *OauthIntegrationForCustomClientsResourceAssert {
	o.AddAssertion(assert.ValueNotSet("oauth_client_rsa_public_key"))
	return o
}

func (o *OauthIntegrationForCustomClientsResourceAssert) HasNoOauthClientRsaPublicKey2() *OauthIntegrationForCustomClientsResourceAssert {
	o.AddAssertion(assert.ValueNotSet("oauth_client_rsa_public_key_2"))
	return o
}

func (o *OauthIntegrationForCustomClientsResourceAssert) HasNoOauthClientType() *OauthIntegrationForCustomClientsResourceAssert {
	o.AddAssertion(assert.ValueNotSet("oauth_client_type"))
	return o
}

func (o *OauthIntegrationForCustomClientsResourceAssert) HasNoOauthEnforcePkce() *OauthIntegrationForCustomClientsResourceAssert {
	o.AddAssertion(assert.ValueNotSet("oauth_enforce_pkce"))
	return o
}

func (o *OauthIntegrationForCustomClientsResourceAssert) HasNoOauthIssueRefreshTokens() *OauthIntegrationForCustomClientsResourceAssert {
	o.AddAssertion(assert.ValueNotSet("oauth_issue_refresh_tokens"))
	return o
}

func (o *OauthIntegrationForCustomClientsResourceAssert) HasNoOauthRedirectUri() *OauthIntegrationForCustomClientsResourceAssert {
	o.AddAssertion(assert.ValueNotSet("oauth_redirect_uri"))
	return o
}

func (o *OauthIntegrationForCustomClientsResourceAssert) HasNoOauthRefreshTokenValidity() *OauthIntegrationForCustomClientsResourceAssert {
	o.AddAssertion(assert.ValueNotSet("oauth_refresh_token_validity"))
	return o
}

func (o *OauthIntegrationForCustomClientsResourceAssert) HasNoOauthUseSecondaryRoles() *OauthIntegrationForCustomClientsResourceAssert {
	o.AddAssertion(assert.ValueNotSet("oauth_use_secondary_roles"))
	return o
}

func (o *OauthIntegrationForCustomClientsResourceAssert) HasNoPreAuthorizedRolesList() *OauthIntegrationForCustomClientsResourceAssert {
	o.AddAssertion(assert.ValueSet("pre_authorized_roles_list.#", "0"))
	return o
}

func (o *OauthIntegrationForCustomClientsResourceAssert) HasNoRelatedParameters() *OauthIntegrationForCustomClientsResourceAssert {
	o.AddAssertion(assert.ValueSet("related_parameters.#", "0"))
	return o
}

////////////////////////////
// Attribute empty checks //
////////////////////////////

func (o *OauthIntegrationForCustomClientsResourceAssert) HasCommentEmpty() *OauthIntegrationForCustomClientsResourceAssert {
	o.AddAssertion(assert.ValueSet("comment", ""))
	return o
}

func (o *OauthIntegrationForCustomClientsResourceAssert) HasEnabledEmpty() *OauthIntegrationForCustomClientsResourceAssert {
	o.AddAssertion(assert.ValueSet("enabled", ""))
	return o
}

func (o *OauthIntegrationForCustomClientsResourceAssert) HasFullyQualifiedNameEmpty() *OauthIntegrationForCustomClientsResourceAssert {
	o.AddAssertion(assert.ValueSet("fully_qualified_name", ""))
	return o
}

func (o *OauthIntegrationForCustomClientsResourceAssert) HasNameEmpty() *OauthIntegrationForCustomClientsResourceAssert {
	o.AddAssertion(assert.ValueSet("name", ""))
	return o
}

func (o *OauthIntegrationForCustomClientsResourceAssert) HasNetworkPolicyEmpty() *OauthIntegrationForCustomClientsResourceAssert {
	o.AddAssertion(assert.ValueSet("network_policy", ""))
	return o
}

func (o *OauthIntegrationForCustomClientsResourceAssert) HasOauthAllowNonTlsRedirectUriEmpty() *OauthIntegrationForCustomClientsResourceAssert {
	o.AddAssertion(assert.ValueSet("oauth_allow_non_tls_redirect_uri", ""))
	return o
}

func (o *OauthIntegrationForCustomClientsResourceAssert) HasOauthClientRsaPublicKeyEmpty() *OauthIntegrationForCustomClientsResourceAssert {
	o.AddAssertion(assert.ValueSet("oauth_client_rsa_public_key", ""))
	return o
}

func (o *OauthIntegrationForCustomClientsResourceAssert) HasOauthClientRsaPublicKey2Empty() *OauthIntegrationForCustomClientsResourceAssert {
	o.AddAssertion(assert.ValueSet("oauth_client_rsa_public_key_2", ""))
	return o
}

func (o *OauthIntegrationForCustomClientsResourceAssert) HasOauthClientTypeEmpty() *OauthIntegrationForCustomClientsResourceAssert {
	o.AddAssertion(assert.ValueSet("oauth_client_type", ""))
	return o
}

func (o *OauthIntegrationForCustomClientsResourceAssert) HasOauthEnforcePkceEmpty() *OauthIntegrationForCustomClientsResourceAssert {
	o.AddAssertion(assert.ValueSet("oauth_enforce_pkce", ""))
	return o
}

func (o *OauthIntegrationForCustomClientsResourceAssert) HasOauthIssueRefreshTokensEmpty() *OauthIntegrationForCustomClientsResourceAssert {
	o.AddAssertion(assert.ValueSet("oauth_issue_refresh_tokens", ""))
	return o
}

func (o *OauthIntegrationForCustomClientsResourceAssert) HasOauthRedirectUriEmpty() *OauthIntegrationForCustomClientsResourceAssert {
	o.AddAssertion(assert.ValueSet("oauth_redirect_uri", ""))
	return o
}

func (o *OauthIntegrationForCustomClientsResourceAssert) HasOauthUseSecondaryRolesEmpty() *OauthIntegrationForCustomClientsResourceAssert {
	o.AddAssertion(assert.ValueSet("oauth_use_secondary_roles", ""))
	return o
}

///////////////////////////////
// Attribute presence checks //
///////////////////////////////

func (o *OauthIntegrationForCustomClientsResourceAssert) HasBlockedRolesListNotEmpty() *OauthIntegrationForCustomClientsResourceAssert {
	o.AddAssertion(assert.ValuePresent("blocked_roles_list"))
	return o
}

func (o *OauthIntegrationForCustomClientsResourceAssert) HasCommentNotEmpty() *OauthIntegrationForCustomClientsResourceAssert {
	o.AddAssertion(assert.ValuePresent("comment"))
	return o
}

func (o *OauthIntegrationForCustomClientsResourceAssert) HasEnabledNotEmpty() *OauthIntegrationForCustomClientsResourceAssert {
	o.AddAssertion(assert.ValuePresent("enabled"))
	return o
}

func (o *OauthIntegrationForCustomClientsResourceAssert) HasFullyQualifiedNameNotEmpty() *OauthIntegrationForCustomClientsResourceAssert {
	o.AddAssertion(assert.ValuePresent("fully_qualified_name"))
	return o
}

func (o *OauthIntegrationForCustomClientsResourceAssert) HasNameNotEmpty() *OauthIntegrationForCustomClientsResourceAssert {
	o.AddAssertion(assert.ValuePresent("name"))
	return o
}

func (o *OauthIntegrationForCustomClientsResourceAssert) HasNetworkPolicyNotEmpty() *OauthIntegrationForCustomClientsResourceAssert {
	o.AddAssertion(assert.ValuePresent("network_policy"))
	return o
}

func (o *OauthIntegrationForCustomClientsResourceAssert) HasOauthAllowNonTlsRedirectUriNotEmpty() *OauthIntegrationForCustomClientsResourceAssert {
	o.AddAssertion(assert.ValuePresent("oauth_allow_non_tls_redirect_uri"))
	return o
}

func (o *OauthIntegrationForCustomClientsResourceAssert) HasOauthClientRsaPublicKeyNotEmpty() *OauthIntegrationForCustomClientsResourceAssert {
	o.AddAssertion(assert.ValuePresent("oauth_client_rsa_public_key"))
	return o
}

func (o *OauthIntegrationForCustomClientsResourceAssert) HasOauthClientRsaPublicKey2NotEmpty() *OauthIntegrationForCustomClientsResourceAssert {
	o.AddAssertion(assert.ValuePresent("oauth_client_rsa_public_key_2"))
	return o
}

func (o *OauthIntegrationForCustomClientsResourceAssert) HasOauthClientTypeNotEmpty() *OauthIntegrationForCustomClientsResourceAssert {
	o.AddAssertion(assert.ValuePresent("oauth_client_type"))
	return o
}

func (o *OauthIntegrationForCustomClientsResourceAssert) HasOauthEnforcePkceNotEmpty() *OauthIntegrationForCustomClientsResourceAssert {
	o.AddAssertion(assert.ValuePresent("oauth_enforce_pkce"))
	return o
}

func (o *OauthIntegrationForCustomClientsResourceAssert) HasOauthIssueRefreshTokensNotEmpty() *OauthIntegrationForCustomClientsResourceAssert {
	o.AddAssertion(assert.ValuePresent("oauth_issue_refresh_tokens"))
	return o
}

func (o *OauthIntegrationForCustomClientsResourceAssert) HasOauthRedirectUriNotEmpty() *OauthIntegrationForCustomClientsResourceAssert {
	o.AddAssertion(assert.ValuePresent("oauth_redirect_uri"))
	return o
}

func (o *OauthIntegrationForCustomClientsResourceAssert) HasOauthRefreshTokenValidityNotEmpty() *OauthIntegrationForCustomClientsResourceAssert {
	o.AddAssertion(assert.ValuePresent("oauth_refresh_token_validity"))
	return o
}

func (o *OauthIntegrationForCustomClientsResourceAssert) HasOauthUseSecondaryRolesNotEmpty() *OauthIntegrationForCustomClientsResourceAssert {
	o.AddAssertion(assert.ValuePresent("oauth_use_secondary_roles"))
	return o
}

func (o *OauthIntegrationForCustomClientsResourceAssert) HasPreAuthorizedRolesListNotEmpty() *OauthIntegrationForCustomClientsResourceAssert {
	o.AddAssertion(assert.ValuePresent("pre_authorized_roles_list"))
	return o
}

func (o *OauthIntegrationForCustomClientsResourceAssert) HasRelatedParametersNotEmpty() *OauthIntegrationForCustomClientsResourceAssert {
	o.AddAssertion(assert.ValuePresent("related_parameters"))
	return o
}
