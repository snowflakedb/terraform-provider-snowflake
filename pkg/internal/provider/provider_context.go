package provider

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

type Context struct {
	Client             *sdk.Client
	EnabledFeatures    []string
	EnabledExperiments []string
	// GrantShowCache caches SHOW GRANTS results, keyed by rendered SQL (see sdk.StructToSQL).
	GrantShowCache *Cache[[]sdk.Grant]
}
