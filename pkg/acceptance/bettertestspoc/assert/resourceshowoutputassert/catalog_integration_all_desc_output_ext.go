package resourceshowoutputassert

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

// Composite methods

func (c *CatalogIntegrationAllDescribeOutputAssert) HasOAuthRestAuthentication(tokenUri, clientId string, scopes ...string) *CatalogIntegrationAllDescribeOutputAssert {
	catalogIntegrationApplyOAuthChecks(c.ResourceAssert, "oauth_rest_authentication", tokenUri, clientId, scopes...)
	return c
}

func (c *CatalogIntegrationAllDescribeOutputAssert) HasSigV4RestAuthentication(iamRole, signingRegion, externalId string) *CatalogIntegrationAllDescribeOutputAssert {
	catalogIntegrationApplySigV4Checks(c.ResourceAssert, iamRole, signingRegion, externalId)
	return c
}

// Individual RestConfig methods

func (c *CatalogIntegrationAllDescribeOutputAssert) HasRestConfigCatalogUri(expected string) *CatalogIntegrationAllDescribeOutputAssert {
	c.StringValueSet("rest_config.0.catalog_uri", expected)
	return c
}

func (c *CatalogIntegrationAllDescribeOutputAssert) HasRestConfigPrefix(expected string) *CatalogIntegrationAllDescribeOutputAssert {
	c.StringValueSet("rest_config.0.prefix", expected)
	return c
}

func (c *CatalogIntegrationAllDescribeOutputAssert) HasRestConfigCatalogName(expected string) *CatalogIntegrationAllDescribeOutputAssert {
	c.StringValueSet("rest_config.0.catalog_name", expected)
	return c
}

func (c *CatalogIntegrationAllDescribeOutputAssert) HasRestConfigCatalogApiType(expected sdk.CatalogIntegrationCatalogApiType) *CatalogIntegrationAllDescribeOutputAssert {
	c.StringValueSet("rest_config.0.catalog_api_type", string(expected))
	return c
}

func (c *CatalogIntegrationAllDescribeOutputAssert) HasRestConfigAccessDelegationMode(expected sdk.CatalogIntegrationAccessDelegationMode) *CatalogIntegrationAllDescribeOutputAssert {
	c.StringValueSet("rest_config.0.access_delegation_mode", string(expected))
	return c
}

// Individual SigV4 methods

func (c *CatalogIntegrationAllDescribeOutputAssert) HasSigv4RestAuthenticationSigv4IamRole(expected string) *CatalogIntegrationAllDescribeOutputAssert {
	c.StringValueSet("sigv4_rest_authentication.0.sigv4_iam_role", expected)
	return c
}

func (c *CatalogIntegrationAllDescribeOutputAssert) HasSigv4RestAuthenticationSigv4SigningRegion(expected string) *CatalogIntegrationAllDescribeOutputAssert {
	c.StringValueSet("sigv4_rest_authentication.0.sigv4_signing_region", expected)
	return c
}

func (c *CatalogIntegrationAllDescribeOutputAssert) HasSigv4RestAuthenticationSigv4ExternalId(expected string) *CatalogIntegrationAllDescribeOutputAssert {
	c.StringValueSet("sigv4_rest_authentication.0.sigv4_external_id", expected)
	return c
}

// Individual OAuth methods

func (c *CatalogIntegrationAllDescribeOutputAssert) HasOAuthRestAuthenticationOauthTokenUri(expected string) *CatalogIntegrationAllDescribeOutputAssert {
	c.StringValueSet("oauth_rest_authentication.0.oauth_token_uri", expected)
	return c
}

func (c *CatalogIntegrationAllDescribeOutputAssert) HasOAuthRestAuthenticationOauthClientId(expected string) *CatalogIntegrationAllDescribeOutputAssert {
	c.StringValueSet("oauth_rest_authentication.0.oauth_client_id", expected)
	return c
}

func (c *CatalogIntegrationAllDescribeOutputAssert) HasOAuthRestAuthenticationOauthClientSecret(expected string) *CatalogIntegrationAllDescribeOutputAssert {
	c.StringValueSet("oauth_rest_authentication.0.oauth_client_secret", expected)
	return c
}

func (c *CatalogIntegrationAllDescribeOutputAssert) HasOAuthRestAuthenticationOauthAllowedScopes(expected ...string) *CatalogIntegrationAllDescribeOutputAssert {
	catalogIntegrationApplyOAuthScopesCheck(c.ResourceAssert, "oauth_rest_authentication", expected...)
	return c
}

// No-value OAuth methods

func (c *CatalogIntegrationAllDescribeOutputAssert) HasNoOAuthRestAuthenticationOauthTokenUri() *CatalogIntegrationAllDescribeOutputAssert {
	c.ValueNotSet("oauth_rest_authentication.0.oauth_token_uri")
	return c
}

func (c *CatalogIntegrationAllDescribeOutputAssert) HasNoOAuthRestAuthenticationOauthClientId() *CatalogIntegrationAllDescribeOutputAssert {
	c.ValueNotSet("oauth_rest_authentication.0.oauth_client_id")
	return c
}

func (c *CatalogIntegrationAllDescribeOutputAssert) HasNoOAuthRestAuthenticationOauthClientSecret() *CatalogIntegrationAllDescribeOutputAssert {
	c.ValueNotSet("oauth_rest_authentication.0.oauth_client_secret")
	return c
}

func (c *CatalogIntegrationAllDescribeOutputAssert) HasNoOAuthRestAuthenticationOauthAllowedScopes() *CatalogIntegrationAllDescribeOutputAssert {
	c.StringValueSet("oauth_rest_authentication.0.oauth_allowed_scopes.#", "0")
	return c
}

// Individual Glue methods

func (c *CatalogIntegrationAllDescribeOutputAssert) HasGlueAwsIamUserArnNotEmpty() *CatalogIntegrationAllDescribeOutputAssert {
	c.ValuePresent("glue_aws_iam_user_arn")
	return c
}

func (c *CatalogIntegrationAllDescribeOutputAssert) HasGlueAwsExternalIdNotEmpty() *CatalogIntegrationAllDescribeOutputAssert {
	c.ValuePresent("glue_aws_external_id")
	return c
}
