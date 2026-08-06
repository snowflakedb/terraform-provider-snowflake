package provider

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

type Context struct {
	Client               *sdk.Client
	EnabledFeatures      []string
	EnabledExperiments   []string
	GrantShowOfRoleCache *Cache[[]sdk.Grant]
	// RoleShowCache caches the result of looking up an account role by identifier (SHOW ROLES LIKE
	// '<name>', via Roles.ShowByID/ShowByIDSafely), keyed by the role's fully qualified name. Shared
	// by snowflake_account_role, snowflake_grant_application_role, and
	// snowflake_grant_privileges_to_account_role (gated by ACCOUNT_ROLE_SHOW_CACHING).
	RoleShowCache *Cache[*sdk.Role]
}
