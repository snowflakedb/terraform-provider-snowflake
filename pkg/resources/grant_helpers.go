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

// showGrantsCachedFor issues opts via client.Grants.Show, caching the result in
// providerCtx.GrantShowCache keyed by rendered SQL (sdk.StructToSQL) when experiment is enabled.
func showGrantsCachedFor(ctx context.Context, providerCtx *provider.Context, experiment experimentalfeatures.ExperimentalFeature, opts *sdk.ShowGrantOptions) ([]sdk.Grant, error) {
	if !experimentalfeatures.IsExperimentEnabled(experiment, providerCtx.EnabledExperiments) {
		return providerCtx.Client.Grants.Show(ctx, opts)
	}
	key, err := sdk.StructToSQL(opts)
	if err != nil {
		log.Printf("[WARN] failed to render SHOW GRANTS cache key, falling back to uncached SHOW: %s", err)
		return providerCtx.Client.Grants.Show(ctx, opts)
	}
	return providerCtx.GrantShowCache.GetOrLoad(ctx, key, func(loadCtx context.Context) ([]sdk.Grant, error) {
		return providerCtx.Client.Grants.Show(loadCtx, opts)
	})
}

// showGrantsCached is showGrantsCachedFor gated on GRANTS_SHOW_CACHING, used by
// snowflake_grant_privileges_to_account_role and snowflake_grant_ownership.
func showGrantsCached(ctx context.Context, providerCtx *provider.Context, opts *sdk.ShowGrantOptions) ([]sdk.Grant, error) {
	return showGrantsCachedFor(ctx, providerCtx, experimentalfeatures.GrantsShowCaching, opts)
}

// invalidateGrantsShowCacheFor invalidates the cached entry for opts under experiment, if enabled.
// A no-op if opts is nil, rendering fails, or nothing was cached for it.
func invalidateGrantsShowCacheFor(providerCtx *provider.Context, experiment experimentalfeatures.ExperimentalFeature, opts *sdk.ShowGrantOptions) {
	if opts == nil || !experimentalfeatures.IsExperimentEnabled(experiment, providerCtx.EnabledExperiments) {
		return
	}
	key, err := sdk.StructToSQL(opts)
	if err != nil {
		log.Printf("[WARN] failed to render SHOW GRANTS cache key for invalidation: %s", err)
		return
	}
	providerCtx.GrantShowCache.Invalidate(key)
}

// invalidateGrantsShowCache is invalidateGrantsShowCacheFor gated on GRANTS_SHOW_CACHING.
func invalidateGrantsShowCache(providerCtx *provider.Context, opts *sdk.ShowGrantOptions) {
	invalidateGrantsShowCacheFor(providerCtx, experimentalfeatures.GrantsShowCaching, opts)
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
