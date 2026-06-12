package provider

import (
	"sync"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

// ShowGrantsOfRoleCache is a per-plan in-memory cache for SHOW GRANTS OF ROLE results.
//
// Without caching, every snowflake_grant_account_role instance issues an independent
// SHOW GRANTS OF ROLE <name> call during Read. Because one SHOW returns all grants for
// a role, N resources sharing the same role name trigger N identical round-trips that
// each return the same full result set — only 1 is needed per plan.
//
// The cache is scoped to a single provider instance (= one Terraform plan/apply cycle),
// so there is no risk of carrying stale data across separate runs. Within a single apply,
// Create and Delete invalidate the relevant cache entry so subsequent Reads in the same
// cycle observe the mutations.
type ShowGrantsOfRoleCache struct {
	mu   sync.RWMutex
	data map[string][]sdk.Grant
}

// NewShowGrantsOfRoleCache returns an initialized, empty cache.
func NewShowGrantsOfRoleCache() *ShowGrantsOfRoleCache {
	return &ShowGrantsOfRoleCache{data: make(map[string][]sdk.Grant)}
}

// GetOrLoad returns the cached grants for key. On a cache miss it calls loadFn,
// stores the result, and returns it. If loadFn returns an error the result is not
// cached and the error is propagated to the caller.
func (c *ShowGrantsOfRoleCache) GetOrLoad(key string, loadFn func() ([]sdk.Grant, error)) ([]sdk.Grant, error) {
	// Fast path: read lock only.
	c.mu.RLock()
	if v, ok := c.data[key]; ok {
		c.mu.RUnlock()
		return v, nil
	}
	c.mu.RUnlock()

	// Slow path: upgrade to write lock, re-check, then load.
	c.mu.Lock()
	defer c.mu.Unlock()
	if v, ok := c.data[key]; ok {
		return v, nil
	}
	v, err := loadFn()
	if err != nil {
		return nil, err
	}
	c.data[key] = v
	return v, nil
}

// Invalidate removes the cached result for key, forcing the next GetOrLoad call
// to re-fetch from Snowflake.
func (c *ShowGrantsOfRoleCache) Invalidate(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
}
