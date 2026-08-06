package model

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
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

func (c *CatalogIntegrationIcebergRestModel) WithOauthRestAuthentication(oauth sdk.OAuthRestAuthenticationRequest, oauthClientSecretVariableName string) *CatalogIntegrationIcebergRestModel {
	scopeVars := make([]tfconfig.Variable, len(oauth.OauthAllowedScopes))
	for i, s := range oauth.OauthAllowedScopes {
		scopeVars[i] = tfconfig.StringVariable(s.Value)
	}
	m := map[string]tfconfig.Variable{
		"oauth_client_id":      tfconfig.StringVariable(oauth.OauthClientId),
		"oauth_client_secret":  config.VariableReference(oauthClientSecretVariableName),
		"oauth_allowed_scopes": tfconfig.ListVariable(scopeVars...),
	}
	if oauth.OauthTokenUri != nil {
		m["oauth_token_uri"] = tfconfig.StringVariable(*oauth.OauthTokenUri)
	}
	c.OauthRestAuthentication = tfconfig.ListVariable(tfconfig.ObjectVariable(m))
	return c
}

func (c *CatalogIntegrationIcebergRestModel) WithBearerRestAuthentication(bearerTokenVariableName string) *CatalogIntegrationIcebergRestModel {
	c.BearerRestAuthentication = tfconfig.ListVariable(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
		"bearer_token": config.VariableReference(bearerTokenVariableName),
	}))
	return c
}

func (c *CatalogIntegrationIcebergRestModel) WithSigV4RestAuthentication(sigv4 sdk.SigV4RestAuthenticationRequest) *CatalogIntegrationIcebergRestModel {
	m := map[string]tfconfig.Variable{
		"sigv4_iam_role": tfconfig.StringVariable(sigv4.Sigv4IamRole),
	}
	if sigv4.Sigv4SigningRegion != nil {
		m["sigv4_signing_region"] = tfconfig.StringVariable(*sigv4.Sigv4SigningRegion)
	}
	if sigv4.Sigv4ExternalId != nil {
		m["sigv4_external_id"] = tfconfig.StringVariable(*sigv4.Sigv4ExternalId)
	}
	c.Sigv4RestAuthentication = tfconfig.ListVariable(tfconfig.ObjectVariable(m))
	return c
}

func CatalogIntegrationIcebergRestOAuth(
	resourceName string,
	name string,
	enabled bool,
	restConfig sdk.IcebergRestRestConfigRequest,
	oAuthRestAuthentication sdk.OAuthRestAuthenticationRequest,
	oauthClientSecretVariableName string,
) *CatalogIntegrationIcebergRestModel {
	return CatalogIntegrationIcebergRest(resourceName, name, enabled, []sdk.IcebergRestRestConfigRequest{restConfig}).
		WithOauthRestAuthentication(oAuthRestAuthentication, oauthClientSecretVariableName)
}

func CatalogIntegrationIcebergRestBearer(
	resourceName string,
	name string,
	enabled bool,
	restConfig sdk.IcebergRestRestConfigRequest,
	bearerTokenVariableName string,
) *CatalogIntegrationIcebergRestModel {
	return CatalogIntegrationIcebergRest(resourceName, name, enabled, []sdk.IcebergRestRestConfigRequest{restConfig}).
		WithBearerRestAuthentication(bearerTokenVariableName)
}

func CatalogIntegrationIcebergRestSigV4(
	resourceName string,
	name string,
	enabled bool,
	restConfig sdk.IcebergRestRestConfigRequest,
	sigV4RestAuthentication sdk.SigV4RestAuthenticationRequest,
) *CatalogIntegrationIcebergRestModel {
	return CatalogIntegrationIcebergRest(resourceName, name, enabled, []sdk.IcebergRestRestConfigRequest{restConfig}).
		WithSigV4RestAuthentication(sigV4RestAuthentication)
}
