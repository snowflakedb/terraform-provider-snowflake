package resourceassert

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (c *CatalogIntegrationIcebergRestResourceAssert) HasRestConfig(restConfig *sdk.IcebergRestRestConfigDetails) *CatalogIntegrationIcebergRestResourceAssert {
	c.ValueSet("rest_config.0.catalog_uri", restConfig.CatalogUri)
	c.ValueSet("rest_config.0.prefix", restConfig.Prefix)
	c.ValueSet("rest_config.0.catalog_name", restConfig.CatalogName)
	c.ValueSet("rest_config.0.catalog_api_type", string(restConfig.CatalogApiType))
	c.ValueSet("rest_config.0.access_delegation_mode", string(restConfig.AccessDelegationMode))
	return c
}

func (c *CatalogIntegrationIcebergRestResourceAssert) HasOauthRestAuthentication(restAuth *sdk.OAuthRestAuthenticationDetails) *CatalogIntegrationIcebergRestResourceAssert {
	c.ValueSet("oauth_rest_authentication.0.oauth_token_uri", restAuth.OauthTokenUri)
	c.ValueSet("oauth_rest_authentication.0.oauth_client_id", restAuth.OauthClientId)
	c.ValueSet("oauth_rest_authentication.0.oauth_client_secret", restAuth.OauthClientSecret)
	c.ListContainsExactlyStringValuesInOrder("oauth_rest_authentication.0.oauth_allowed_scopes", restAuth.OauthAllowedScopes...)
	return c
}

func (c *CatalogIntegrationIcebergRestResourceAssert) HasBearerRestAuthentication(restAuth *sdk.BearerRestAuthenticationDetails) *CatalogIntegrationIcebergRestResourceAssert {
	c.ValueSet("bearer_rest_authentication.0.bearer_token", restAuth.BearerToken)
	return c
}

func (c *CatalogIntegrationIcebergRestResourceAssert) HasSigV4RestAuthentication(restAuth *sdk.SigV4RestAuthenticationDetails) *CatalogIntegrationIcebergRestResourceAssert {
	c.ValueSet("sigv4_rest_authentication.0.sigv4_iam_role", restAuth.Sigv4IamRole)
	c.ValueSet("sigv4_rest_authentication.0.sigv4_signing_region", restAuth.Sigv4SigningRegion)
	c.ValueSet("sigv4_rest_authentication.0.sigv4_external_id", restAuth.Sigv4ExternalId)
	return c
}
