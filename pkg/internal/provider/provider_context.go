package provider

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

type Context struct {
	Client             *sdk.Client
	EnabledFeatures    []string
	EnabledExperiments []string
	// GrantShowCache caches SHOW GRANTS results, keyed by the rendered SQL of the ShowGrantOptions
	// that produced them (see sdk.ShowGrantOptionsToSQL). Shared by snowflake_grant_account_role
	// (gated by GRANT_ACCOUNT_ROLE_SHOW_CACHING), and snowflake_grant_privileges_to_account_role /
	// snowflake_grant_ownership (gated by GRANTS_SHOW_CACHING).
	GrantShowCache *Cache[[]sdk.Grant]
}
