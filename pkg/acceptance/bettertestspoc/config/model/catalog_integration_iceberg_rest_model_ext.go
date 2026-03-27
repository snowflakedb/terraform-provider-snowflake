package model

import (
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (c *CatalogIntegrationIcebergRestModel) WithRestConfig(restConfig []sdk.IcebergRestRestConfigRequest) *CatalogIntegrationIcebergRestModel {
	if len(restConfig) == 0 {
		return c
	}
	rc := restConfig[0]
	m := map[string]tfconfig.Variable{
		"catalog_uri": tfconfig.StringVariable(rc.CatalogUri),
	}
	if rc.Prefix != nil {
		m["prefix"] = tfconfig.StringVariable(*rc.Prefix)
	}
	if rc.CatalogName != nil {
		m["catalog_name"] = tfconfig.StringVariable(*rc.CatalogName)
	}
	if rc.CatalogApiType != nil {
		m["catalog_api_type"] = tfconfig.StringVariable(string(*rc.CatalogApiType))
	}
	if rc.AccessDelegationMode != nil {
		m["access_delegation_mode"] = tfconfig.StringVariable(string(*rc.AccessDelegationMode))
	}
	c.RestConfig = tfconfig.ListVariable(tfconfig.ObjectVariable(m))
	return c
}

func (c *CatalogIntegrationIcebergRestModel) WithOauthRestAuthentication(oauth []sdk.OAuthRestAuthenticationRequest) *CatalogIntegrationIcebergRestModel {
	if len(oauth) == 0 {
		return c
	}
	ra := oauth[0]
	scopeVars := make([]tfconfig.Variable, len(ra.OauthAllowedScopes))
	for i, s := range ra.OauthAllowedScopes {
		scopeVars[i] = tfconfig.StringVariable(s.Value)
	}
	m := map[string]tfconfig.Variable{
		"oauth_client_id":      tfconfig.StringVariable(ra.OauthClientId),
		"oauth_client_secret":  tfconfig.StringVariable(ra.OauthClientSecret),
		"oauth_allowed_scopes": tfconfig.ListVariable(scopeVars...),
	}
	if ra.OauthTokenUri != nil {
		m["oauth_token_uri"] = tfconfig.StringVariable(*ra.OauthTokenUri)
	}
	c.OauthRestAuthentication = tfconfig.ListVariable(tfconfig.ObjectVariable(m))
	return c
}

func (c *CatalogIntegrationIcebergRestModel) WithBearerRestAuthentication(bearer []sdk.BearerRestAuthenticationRequest) *CatalogIntegrationIcebergRestModel {
	if len(bearer) == 0 {
		return c
	}
	b := bearer[0]
	c.BearerRestAuthentication = tfconfig.ListVariable(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
		"bearer_token": tfconfig.StringVariable(b.BearerToken),
	}))
	return c
}

func (c *CatalogIntegrationIcebergRestModel) WithSigv4RestAuthentication(sigv4 []sdk.SigV4RestAuthenticationRequest) *CatalogIntegrationIcebergRestModel {
	if len(sigv4) == 0 {
		return c
	}
	s := sigv4[0]
	m := map[string]tfconfig.Variable{
		"sigv4_iam_role": tfconfig.StringVariable(s.Sigv4IamRole),
	}
	if s.Sigv4SigningRegion != nil {
		m["sigv4_signing_region"] = tfconfig.StringVariable(*s.Sigv4SigningRegion)
	}
	if s.Sigv4ExternalId != nil {
		m["sigv4_external_id"] = tfconfig.StringVariable(*s.Sigv4ExternalId)
	}
	c.Sigv4RestAuthentication = tfconfig.ListVariable(tfconfig.ObjectVariable(m))
	return c
}

// CatalogIntegrationIcebergRestOAuth is a convenience constructor for the common case (OAuth + REST config).
func CatalogIntegrationIcebergRestOAuth(
	resourceName string,
	name string,
	enabled bool,
	oauth []sdk.OAuthRestAuthenticationRequest,
	restConfig []sdk.IcebergRestRestConfigRequest,
) *CatalogIntegrationIcebergRestModel {
	return CatalogIntegrationIcebergRest(resourceName, name, enabled, restConfig).WithOauthRestAuthentication(oauth)
}
