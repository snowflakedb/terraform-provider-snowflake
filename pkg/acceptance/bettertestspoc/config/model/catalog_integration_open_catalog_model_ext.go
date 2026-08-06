package model

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
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
	return c.withRestAuthentication(restAuthentication, "")
}

func (c *CatalogIntegrationOpenCatalogModel) withRestAuthentication(restAuthentication []sdk.OAuthRestAuthenticationRequest, oauthClientSecretVariableName string) *CatalogIntegrationOpenCatalogModel {
	if len(restAuthentication) == 0 {
		return c
	}
	ra := restAuthentication[0]
	scopeVars := make([]tfconfig.Variable, len(ra.OauthAllowedScopes))
	for i, s := range ra.OauthAllowedScopes {
		scopeVars[i] = tfconfig.StringVariable(s.Value)
	}
	var oauthClientSecret tfconfig.Variable
	if oauthClientSecretVariableName != "" {
		oauthClientSecret = config.VariableReference(oauthClientSecretVariableName)
	} else {
		oauthClientSecret = tfconfig.StringVariable(ra.OauthClientSecret)
	}
	m := map[string]tfconfig.Variable{
		"oauth_client_id":      tfconfig.StringVariable(ra.OauthClientId),
		"oauth_client_secret":  oauthClientSecret,
		"oauth_allowed_scopes": tfconfig.ListVariable(scopeVars...),
	}
	if ra.OauthTokenUri != nil {
		m["oauth_token_uri"] = tfconfig.StringVariable(*ra.OauthTokenUri)
	}
	c.RestAuthentication = tfconfig.ListVariable(tfconfig.ObjectVariable(m))
	return c
}

func CatalogIntegrationOpenCatalogVar(
	resourceName string,
	name string,
	enabled bool,
	restAuthentication []sdk.OAuthRestAuthenticationRequest,
	restConfig []sdk.OpenCatalogRestConfigRequest,
	oauthClientSecretVariableName string,
) *CatalogIntegrationOpenCatalogModel {
	c := &CatalogIntegrationOpenCatalogModel{ResourceModelMeta: config.Meta(resourceName, resources.CatalogIntegrationOpenCatalog)}
	c.WithName(name)
	c.WithEnabled(enabled)
	c.withRestAuthentication(restAuthentication, oauthClientSecretVariableName)
	c.WithRestConfig(restConfig)
	return c
}
