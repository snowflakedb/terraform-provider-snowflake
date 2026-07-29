package provider

import "sync"

// Cache is a simple per-plan in-memory, concurrency-safe cache keyed by string.
//
// It is intended for caching expensive, read-only Snowflake lookups whose result is shared
// by many resource instances within a single Terraform plan/apply cycle. The canonical use is
// SHOW GRANTS OF ROLE: without caching, every snowflake_grant_account_role instance issues an
// independent SHOW GRANTS OF ROLE <name> call during Read. Because one SHOW returns all grants
// for a role, N resources sharing the same role name trigger N identical round-trips that each
// return the same full result set — only 1 is needed per plan.
//
// The cache is scoped to a single provider instance (= one Terraform plan/apply cycle), so there
// is no risk of carrying stale data across separate runs. Within a single apply, callers that
// mutate the underlying object (e.g. Create/Delete) must Invalidate the relevant key so subsequent
// Reads in the same cycle observe the mutation.
//
// Cache is generic over the cached value type so it can be reused for other lookups in the future.
type Cache[T any] struct {
	mu   sync.RWMutex
	data map[string]T
}

// NewCache returns an initialized, empty cache.
func NewCache[T any]() *Cache[T] {
	return &Cache[T]{data: make(map[string]T)}
}

// GetOrLoad returns the cached value for key. On a cache miss it calls loadFn,
// stores the result, and returns it. If loadFn returns an error the result is not
// cached and the error is propagated to the caller.
func (c *Cache[T]) GetOrLoad(key string, loadFn func() (T, error)) (T, error) {
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
		var zero T
		return zero, err
	}
	c.data[key] = v
	return v, nil
}

// Invalidate removes the cached result for key, forcing the next GetOrLoad call
// to re-fetch from the source.
func (c *Cache[T]) Invalidate(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
}
