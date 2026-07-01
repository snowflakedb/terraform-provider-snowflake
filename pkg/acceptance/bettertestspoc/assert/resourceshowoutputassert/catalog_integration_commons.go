package resourceshowoutputassert

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func catalogIntegrationApplyRestConfigChecks(e *assert.ResourceAssert, catalogUri, prefix, catalogName string, catalogApiType sdk.CatalogIntegrationCatalogApiType, accessDelegationMode sdk.CatalogIntegrationAccessDelegationMode) {
	e.StringValueSet("rest_config.0.catalog_uri", catalogUri)
	e.StringValueSet("rest_config.0.prefix", prefix)
	e.StringValueSet("rest_config.0.catalog_name", catalogName)
	e.StringValueSet("rest_config.0.catalog_api_type", string(catalogApiType))
	e.StringValueSet("rest_config.0.access_delegation_mode", string(accessDelegationMode))
}

func catalogIntegrationApplyOAuthChecks(e *assert.ResourceAssert, fieldPrefix, tokenUri, clientId string, scopes ...string) {
	e.StringValueSet(fieldPrefix+".0.oauth_token_uri", tokenUri)
	e.StringValueSet(fieldPrefix+".0.oauth_client_id", clientId)
	e.StringValueSet(fieldPrefix+".0.oauth_allowed_scopes.#", fmt.Sprintf("%d", len(scopes)))
	for i, v := range scopes {
		e.StringValueSet(fmt.Sprintf("%s.0.oauth_allowed_scopes.%d", fieldPrefix, i), v)
	}
}

func catalogIntegrationApplySigV4Checks(e *assert.ResourceAssert, iamRole, signingRegion, externalId string) {
	e.StringValueSet("sigv4_rest_authentication.0.sigv4_iam_role", iamRole)
	e.StringValueSet("sigv4_rest_authentication.0.sigv4_signing_region", signingRegion)
	e.StringValueSet("sigv4_rest_authentication.0.sigv4_external_id", externalId)
}

func catalogIntegrationApplyOAuthScopesCheck(e *assert.ResourceAssert, fieldPrefix string, scopes ...string) {
	e.StringValueSet(fieldPrefix+".0.oauth_allowed_scopes.#", fmt.Sprintf("%d", len(scopes)))
	for i, v := range scopes {
		e.StringValueSet(fmt.Sprintf("%s.0.oauth_allowed_scopes.%d", fieldPrefix, i), v)
	}
}
