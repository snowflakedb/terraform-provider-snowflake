package resourceshowoutputassert

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

// Composite methods

func (c *CatalogIntegrationIcebergRestDescribeOutputAssert) HasRestConfig(catalogUri, prefix, catalogName string, catalogApiType sdk.CatalogIntegrationCatalogApiType, accessDelegationMode sdk.CatalogIntegrationAccessDelegationMode) *CatalogIntegrationIcebergRestDescribeOutputAssert {
	catalogIntegrationApplyRestConfigChecks(c.ResourceAssert, catalogUri, prefix, catalogName, catalogApiType, accessDelegationMode)
	return c
}

func (c *CatalogIntegrationIcebergRestDescribeOutputAssert) HasOAuthRestAuthentication(tokenUri, clientId string, scopes ...string) *CatalogIntegrationIcebergRestDescribeOutputAssert {
	catalogIntegrationApplyOAuthChecks(c.ResourceAssert, "oauth_rest_authentication", tokenUri, clientId, scopes...)
	return c
}

func (c *CatalogIntegrationIcebergRestDescribeOutputAssert) HasSigV4RestAuthentication(iamRole, signingRegion, externalId string) *CatalogIntegrationIcebergRestDescribeOutputAssert {
	catalogIntegrationApplySigV4Checks(c.ResourceAssert, iamRole, signingRegion, externalId)
	return c
}

// Individual RestConfig methods

func (c *CatalogIntegrationIcebergRestDescribeOutputAssert) HasRestConfigCatalogUri(expected string) *CatalogIntegrationIcebergRestDescribeOutputAssert {
	c.StringValueSet("rest_config.0.catalog_uri", expected)
	return c
}

func (c *CatalogIntegrationIcebergRestDescribeOutputAssert) HasRestConfigPrefix(expected string) *CatalogIntegrationIcebergRestDescribeOutputAssert {
	c.StringValueSet("rest_config.0.prefix", expected)
	return c
}

func (c *CatalogIntegrationIcebergRestDescribeOutputAssert) HasRestConfigCatalogName(expected string) *CatalogIntegrationIcebergRestDescribeOutputAssert {
	c.StringValueSet("rest_config.0.catalog_name", expected)
	return c
}

func (c *CatalogIntegrationIcebergRestDescribeOutputAssert) HasRestConfigCatalogApiType(expected sdk.CatalogIntegrationCatalogApiType) *CatalogIntegrationIcebergRestDescribeOutputAssert {
	c.StringValueSet("rest_config.0.catalog_api_type", string(expected))
	return c
}

func (c *CatalogIntegrationIcebergRestDescribeOutputAssert) HasRestConfigAccessDelegationMode(expected sdk.CatalogIntegrationAccessDelegationMode) *CatalogIntegrationIcebergRestDescribeOutputAssert {
	c.StringValueSet("rest_config.0.access_delegation_mode", string(expected))
	return c
}

// Individual SigV4 methods

func (c *CatalogIntegrationIcebergRestDescribeOutputAssert) HasSigv4RestAuthenticationSigv4IamRole(expected string) *CatalogIntegrationIcebergRestDescribeOutputAssert {
	c.StringValueSet("sigv4_rest_authentication.0.sigv4_iam_role", expected)
	return c
}

func (c *CatalogIntegrationIcebergRestDescribeOutputAssert) HasSigv4RestAuthenticationSigv4SigningRegion(expected string) *CatalogIntegrationIcebergRestDescribeOutputAssert {
	c.StringValueSet("sigv4_rest_authentication.0.sigv4_signing_region", expected)
	return c
}

func (c *CatalogIntegrationIcebergRestDescribeOutputAssert) HasSigv4RestAuthenticationSigv4ExternalId(expected string) *CatalogIntegrationIcebergRestDescribeOutputAssert {
	c.StringValueSet("sigv4_rest_authentication.0.sigv4_external_id", expected)
	return c
}

// Individual OAuth methods

func (c *CatalogIntegrationIcebergRestDescribeOutputAssert) HasOAuthRestAuthenticationOauthTokenUri(expected string) *CatalogIntegrationIcebergRestDescribeOutputAssert {
	c.StringValueSet("oauth_rest_authentication.0.oauth_token_uri", expected)
	return c
}

func (c *CatalogIntegrationIcebergRestDescribeOutputAssert) HasOAuthRestAuthenticationOauthClientId(expected string) *CatalogIntegrationIcebergRestDescribeOutputAssert {
	c.StringValueSet("oauth_rest_authentication.0.oauth_client_id", expected)
	return c
}

func (c *CatalogIntegrationIcebergRestDescribeOutputAssert) HasOAuthRestAuthenticationOauthClientSecret(expected string) *CatalogIntegrationIcebergRestDescribeOutputAssert {
	c.StringValueSet("oauth_rest_authentication.0.oauth_client_secret", expected)
	return c
}

func (c *CatalogIntegrationIcebergRestDescribeOutputAssert) HasOAuthRestAuthenticationOauthAllowedScopes(expected ...string) *CatalogIntegrationIcebergRestDescribeOutputAssert {
	catalogIntegrationApplyOAuthScopesCheck(c.ResourceAssert, "oauth_rest_authentication", expected...)
	return c
}

// No-value OAuth methods

func (c *CatalogIntegrationIcebergRestDescribeOutputAssert) HasNoOAuthRestAuthenticationOauthTokenUri() *CatalogIntegrationIcebergRestDescribeOutputAssert {
	c.ValueNotSet("oauth_rest_authentication.0.oauth_token_uri")
	return c
}

func (c *CatalogIntegrationIcebergRestDescribeOutputAssert) HasNoOAuthRestAuthenticationOauthClientId() *CatalogIntegrationIcebergRestDescribeOutputAssert {
	c.ValueNotSet("oauth_rest_authentication.0.oauth_client_id")
	return c
}

func (c *CatalogIntegrationIcebergRestDescribeOutputAssert) HasNoOAuthRestAuthenticationOauthClientSecret() *CatalogIntegrationIcebergRestDescribeOutputAssert {
	c.ValueNotSet("oauth_rest_authentication.0.oauth_client_secret")
	return c
}

func (c *CatalogIntegrationIcebergRestDescribeOutputAssert) HasNoOAuthRestAuthenticationOauthAllowedScopes() *CatalogIntegrationIcebergRestDescribeOutputAssert {
	c.StringValueSet("oauth_rest_authentication.0.oauth_allowed_scopes.#", "0")
	return c
}
