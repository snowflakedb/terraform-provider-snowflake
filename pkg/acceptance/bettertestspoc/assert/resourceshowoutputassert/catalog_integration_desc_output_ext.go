package resourceshowoutputassert

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (c *CatalogIntegrationDescribeOutputAssert) HasRestConfigCatalogUri(expected string) *CatalogIntegrationDescribeOutputAssert {
	c.StringValueSet("rest_config.0.catalog_uri", expected)
	return c
}

func (c *CatalogIntegrationDescribeOutputAssert) HasRestConfigPrefix(expected string) *CatalogIntegrationDescribeOutputAssert {
	c.StringValueSet("rest_config.0.prefix", expected)
	return c
}

func (c *CatalogIntegrationDescribeOutputAssert) HasRestConfigCatalogName(expected string) *CatalogIntegrationDescribeOutputAssert {
	c.StringValueSet("rest_config.0.catalog_name", expected)
	return c
}

func (c *CatalogIntegrationDescribeOutputAssert) HasRestConfigCatalogApiType(expected sdk.CatalogIntegrationCatalogApiType) *CatalogIntegrationDescribeOutputAssert {
	c.StringValueSet("rest_config.0.catalog_api_type", string(expected))
	return c
}

func (c *CatalogIntegrationDescribeOutputAssert) HasRestConfigAccessDelegationMode(expected sdk.CatalogIntegrationAccessDelegationMode) *CatalogIntegrationDescribeOutputAssert {
	c.StringValueSet("rest_config.0.access_delegation_mode", string(expected))
	return c
}

func (c *CatalogIntegrationDescribeOutputAssert) HasSigv4RestAuthenticationSigv4IamRole(expected string) *CatalogIntegrationDescribeOutputAssert {
	c.StringValueSet("sigv4_rest_authentication.0.sigv4_iam_role", expected)
	return c
}

func (c *CatalogIntegrationDescribeOutputAssert) HasSigv4RestAuthenticationSigv4SigningRegion(expected string) *CatalogIntegrationDescribeOutputAssert {
	c.StringValueSet("sigv4_rest_authentication.0.sigv4_signing_region", expected)
	return c
}

func (c *CatalogIntegrationDescribeOutputAssert) HasSigv4RestAuthenticationSigv4ExternalId(expected string) *CatalogIntegrationDescribeOutputAssert {
	c.StringValueSet("sigv4_rest_authentication.0.sigv4_external_id", expected)
	return c
}

func (c *CatalogIntegrationDescribeOutputAssert) HasOAuthRestAuthenticationOauthTokenUri(expected string) *CatalogIntegrationDescribeOutputAssert {
	c.StringValueSet("oauth_rest_authentication.0.oauth_token_uri", expected)
	return c
}

func (c *CatalogIntegrationDescribeOutputAssert) HasOAuthRestAuthenticationOauthClientId(expected string) *CatalogIntegrationDescribeOutputAssert {
	c.StringValueSet("oauth_rest_authentication.0.oauth_client_id", expected)
	return c
}

func (c *CatalogIntegrationDescribeOutputAssert) HasOAuthRestAuthenticationOauthClientSecret(expected string) *CatalogIntegrationDescribeOutputAssert {
	c.StringValueSet("oauth_rest_authentication.0.oauth_client_secret", expected)
	return c
}

func (c *CatalogIntegrationDescribeOutputAssert) HasOAuthRestAuthenticationOauthAllowedScopes(expected ...string) *CatalogIntegrationDescribeOutputAssert {
	c.StringValueSet("oauth_rest_authentication.0.oauth_allowed_scopes.#", fmt.Sprintf("%d", len(expected)))
	for i, v := range expected {
		c.StringValueSet(fmt.Sprintf("oauth_rest_authentication.0.oauth_allowed_scopes.%d", i), v)
	}
	return c
}

func (c *CatalogIntegrationDescribeOutputAssert) HasNoOAuthRestAuthenticationOauthTokenUri() *CatalogIntegrationDescribeOutputAssert {
	c.ValueNotSet("oauth_rest_authentication.0.oauth_token_uri")
	return c
}

func (c *CatalogIntegrationDescribeOutputAssert) HasNoOAuthRestAuthenticationOauthClientId() *CatalogIntegrationDescribeOutputAssert {
	c.ValueNotSet("oauth_rest_authentication.0.oauth_client_id")
	return c
}

func (c *CatalogIntegrationDescribeOutputAssert) HasNoOAuthRestAuthenticationOauthClientSecret() *CatalogIntegrationDescribeOutputAssert {
	c.ValueNotSet("oauth_rest_authentication.0.oauth_client_secret")
	return c
}

func (c *CatalogIntegrationDescribeOutputAssert) HasNoOAuthRestAuthenticationOauthAllowedScopes() *CatalogIntegrationDescribeOutputAssert {
	c.StringValueSet("oauth_rest_authentication.0.oauth_allowed_scopes.#", "0")
	return c
}

func (c *CatalogIntegrationDescribeOutputAssert) HasGlueAwsIamUserArnNotEmpty() *CatalogIntegrationDescribeOutputAssert {
	c.ValuePresent("glue_aws_iam_user_arn")
	return c
}

func (c *CatalogIntegrationDescribeOutputAssert) HasGlueAwsExternalIdNotEmpty() *CatalogIntegrationDescribeOutputAssert {
	c.ValuePresent("glue_aws_external_id")
	return c
}
