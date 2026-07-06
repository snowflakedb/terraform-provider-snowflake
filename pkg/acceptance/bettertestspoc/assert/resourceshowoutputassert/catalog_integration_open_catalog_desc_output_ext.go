package resourceshowoutputassert

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

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
	c.StringValueSet("rest_authentication.0.oauth_allowed_scopes.#", fmt.Sprintf("%d", len(expected)))
	for i, v := range expected {
		c.StringValueSet(fmt.Sprintf("rest_authentication.0.oauth_allowed_scopes.%d", i), v)
	}
	return c
}

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
