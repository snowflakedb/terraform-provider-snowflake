package resources

import (
	"context"
	"log"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/experimentalfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

// toPrivileges validates and normalizes a list of privilege names using sdk.ToPrivilege.
// It is used when parsing untrusted input (e.g. resource identifiers during import).
func toPrivileges(privileges []string) ([]string, error) {
	return collections.MapErr(privileges, sdk.ToPrivilege)
}

// showGrantsCached issues opts via client.Grants.Show, transparently caching the result in
// providerCtx.GrantShowCache when the GRANTS_SHOW_CACHING experiment is enabled. It is shared by
// snowflake_grant_privileges_to_account_role and snowflake_grant_ownership; it is NOT used by
// snowflake_grant_account_role, which is gated by the separate, longer-standing
// GRANT_ACCOUNT_ROLE_SHOW_CACHING experiment (see showGrantsOfRoleCacheKey in grant_account_role.go).
//
// The cache key is the rendered SQL of opts (sdk.ShowGrantOptionsToSQL), so identical SHOW
// statements issued from different resource instances/types share one cache entry. Callers that
// mutate what a given SHOW would return (grant/revoke/ownership-transfer) must invalidate the same
// key — see grantShowCacheKey — strictly after the mutating SQL executes, not before, so a
// concurrent Read can't repopulate the cache with pre-mutation data that nothing will invalidate
// again.
func showGrantsCached(ctx context.Context, providerCtx *provider.Context, opts *sdk.ShowGrantOptions) ([]sdk.Grant, error) {
	if !experimentalfeatures.IsExperimentEnabled(experimentalfeatures.GrantsShowCaching, providerCtx.EnabledExperiments) {
		return providerCtx.Client.Grants.Show(ctx, opts)
	}
	key, err := sdk.ShowGrantOptionsToSQL(opts)
	if err != nil {
		// Cache-key rendering should never fail for a well-formed opts value built by this
		// provider; fall back to an uncached SHOW rather than fail the caller over it.
		log.Printf("[WARN] failed to render SHOW GRANTS cache key, falling back to uncached SHOW: %s", err)
		return providerCtx.Client.Grants.Show(ctx, opts)
	}
	return providerCtx.GrantShowCache.GetOrLoad(ctx, key, func(loadCtx context.Context) ([]sdk.Grant, error) {
		return providerCtx.Client.Grants.Show(loadCtx, opts)
	})
}

// invalidateGrantsShowCache invalidates the GrantShowCache entry for opts, if the GRANTS_SHOW_CACHING
// experiment is enabled. A no-op (not an error) if opts fails to render or the cache was never
// populated for it — Invalidate on an absent key is already a documented no-op.
func invalidateGrantsShowCache(providerCtx *provider.Context, opts *sdk.ShowGrantOptions) {
	if !experimentalfeatures.IsExperimentEnabled(experimentalfeatures.GrantsShowCaching, providerCtx.EnabledExperiments) {
		return
	}
	key, err := sdk.ShowGrantOptionsToSQL(opts)
	if err != nil {
		log.Printf("[WARN] failed to render SHOW GRANTS cache key for invalidation: %s", err)
		return
	}
	providerCtx.GrantShowCache.Invalidate(key)
}

func isNotOwnershipGrant() func(value any, path cty.Path) diag.Diagnostics {
	return func(value any, path cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics
		if privilege, ok := value.(string); ok && strings.ToUpper(privilege) == "OWNERSHIP" {
			diags = append(diags, diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       "Unsupported privilege 'OWNERSHIP'",
				Detail:        "Granting ownership is only allowed in snowflake_grant_ownership resource.",
				AttributePath: nil,
			})
		}
		return diags
	}
}
