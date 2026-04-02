package model

import (
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (c *CatalogIntegrationOpenCatalogModel) WithRestConfig(restConfig []sdk.OpenCatalogRestConfigRequest) *CatalogIntegrationOpenCatalogModel {
	if len(restConfig) == 0 {
		return c
	}
	rc := restConfig[0]
	m := map[string]tfconfig.Variable{
		"catalog_uri":  tfconfig.StringVariable(rc.CatalogUri),
		"catalog_name": tfconfig.StringVariable(rc.CatalogName),
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

func (c *CatalogIntegrationOpenCatalogModel) WithRestAuthentication(restAuthentication []sdk.OAuthRestAuthenticationRequest) *CatalogIntegrationOpenCatalogModel {
	if len(restAuthentication) == 0 {
		return c
	}
	ra := restAuthentication[0]
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
	c.RestAuthentication = tfconfig.ListVariable(tfconfig.ObjectVariable(m))
	return c
}
