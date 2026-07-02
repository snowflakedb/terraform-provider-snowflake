package resourceshowoutputassert

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

// Composite methods

func (c *CatalogIntegrationOpenCatalogDescribeOutputAssert) HasRestAuthentication(tokenUri, clientId string, scopes ...string) *CatalogIntegrationOpenCatalogDescribeOutputAssert {
	catalogIntegrationApplyOAuthChecks(c.ResourceAssert, "rest_authentication", tokenUri, clientId, scopes...)
	return c
}

// Individual RestConfig methods

func (c *CatalogIntegrationOpenCatalogDescribeOutputAssert) HasRestConfigCatalogUri(expected string) *CatalogIntegrationOpenCatalogDescribeOutputAssert {
	c.StringValueSet("rest_config.0.catalog_uri", expected)
	return c
}

func (c *CatalogIntegrationOpenCatalogDescribeOutputAssert) HasRestConfigPrefix(expected string) *CatalogIntegrationOpenCatalogDescribeOutputAssert {
	c.StringValueSet("rest_config.0.prefix", expected)
	return c
}

func (c *CatalogIntegrationOpenCatalogDescribeOutputAssert) HasRestConfigCatalogName(expected string) *CatalogIntegrationOpenCatalogDescribeOutputAssert {
	c.StringValueSet("rest_config.0.catalog_name", expected)
	return c
}

func (c *CatalogIntegrationOpenCatalogDescribeOutputAssert) HasRestConfigCatalogApiType(expected sdk.CatalogIntegrationCatalogApiType) *CatalogIntegrationOpenCatalogDescribeOutputAssert {
	c.StringValueSet("rest_config.0.catalog_api_type", string(expected))
	return c
}

func (c *CatalogIntegrationOpenCatalogDescribeOutputAssert) HasRestConfigAccessDelegationMode(expected sdk.CatalogIntegrationAccessDelegationMode) *CatalogIntegrationOpenCatalogDescribeOutputAssert {
	c.StringValueSet("rest_config.0.access_delegation_mode", string(expected))
	return c
}

// Individual OAuth methods

func (c *CatalogIntegrationOpenCatalogDescribeOutputAssert) HasRestAuthenticationOauthTokenUri(expected string) *CatalogIntegrationOpenCatalogDescribeOutputAssert {
	c.StringValueSet("rest_authentication.0.oauth_token_uri", expected)
	return c
}

func (c *CatalogIntegrationOpenCatalogDescribeOutputAssert) HasRestAuthenticationOauthClientId(expected string) *CatalogIntegrationOpenCatalogDescribeOutputAssert {
	c.StringValueSet("rest_authentication.0.oauth_client_id", expected)
	return c
}

func (c *CatalogIntegrationOpenCatalogDescribeOutputAssert) HasRestAuthenticationOauthClientSecret(expected string) *CatalogIntegrationOpenCatalogDescribeOutputAssert {
	c.StringValueSet("rest_authentication.0.oauth_client_secret", expected)
	return c
}

func (c *CatalogIntegrationOpenCatalogDescribeOutputAssert) HasRestAuthenticationOauthAllowedScopes(expected ...string) *CatalogIntegrationOpenCatalogDescribeOutputAssert {
	catalogIntegrationApplyOAuthScopesCheck(c.ResourceAssert, "rest_authentication", expected...)
	return c
}

// No-value OAuth methods

func (c *CatalogIntegrationOpenCatalogDescribeOutputAssert) HasNoRestAuthenticationOauthTokenUri() *CatalogIntegrationOpenCatalogDescribeOutputAssert {
	c.ValueNotSet("rest_authentication.0.oauth_token_uri")
	return c
}

func (c *CatalogIntegrationOpenCatalogDescribeOutputAssert) HasNoRestAuthenticationOauthClientId() *CatalogIntegrationOpenCatalogDescribeOutputAssert {
	c.ValueNotSet("rest_authentication.0.oauth_client_id")
	return c
}

func (c *CatalogIntegrationOpenCatalogDescribeOutputAssert) HasNoRestAuthenticationOauthClientSecret() *CatalogIntegrationOpenCatalogDescribeOutputAssert {
	c.ValueNotSet("rest_authentication.0.oauth_client_secret")
	return c
}

func (c *CatalogIntegrationOpenCatalogDescribeOutputAssert) HasNoRestAuthenticationOauthAllowedScopes() *CatalogIntegrationOpenCatalogDescribeOutputAssert {
	c.StringValueSet("rest_authentication.0.oauth_allowed_scopes.#", "0")
	return c
}
