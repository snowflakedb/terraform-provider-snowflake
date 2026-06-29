package resourceassert

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (c *CatalogIntegrationOpenCatalogResourceAssert) HasRestConfig(restConfig *sdk.OpenCatalogRestConfigDetails) *CatalogIntegrationOpenCatalogResourceAssert {
	c.ValueSet("rest_config.0.catalog_uri", restConfig.CatalogUri)
	c.ValueSet("rest_config.0.catalog_name", restConfig.CatalogName)
	c.ValueSet("rest_config.0.catalog_api_type", string(restConfig.CatalogApiType))
	c.ValueSet("rest_config.0.access_delegation_mode", string(restConfig.AccessDelegationMode))
	return c
}

func (c *CatalogIntegrationOpenCatalogResourceAssert) HasRestAuthentication(restAuth *sdk.OAuthRestAuthenticationDetails) *CatalogIntegrationOpenCatalogResourceAssert {
	c.ValueSet("rest_authentication.0.oauth_token_uri", restAuth.OauthTokenUri)
	c.ValueSet("rest_authentication.0.oauth_client_id", restAuth.OauthClientId)
	c.ValueSet("rest_authentication.0.oauth_client_secret", restAuth.OauthClientSecret)
	c.ListContainsExactlyStringValuesInOrder("rest_authentication.0.oauth_allowed_scopes", restAuth.OauthAllowedScopes...)
	return c
}
