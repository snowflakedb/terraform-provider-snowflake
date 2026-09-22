package provider

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/experimentalfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type Context struct {
	Client               *sdk.Client
	EnabledFeatures      []string
	Experiments          experimentalfeatures.Experiments
	SpanID               string // correlates telemetry events for one provider configure/run
	GrantShowOfRoleCache *Cache[[]sdk.Grant]
	RoleShowCache        *Cache[*sdk.Role]
	// GrantShowCache caches SHOW GRANTS results, keyed by rendered SQL (see sdk.StructToSQL).
	GrantShowCache *Cache[[]sdk.Grant]
}
