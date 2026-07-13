package provider

import (
	"sync"

	"golang.org/x/sync/singleflight"
)

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
// Concurrent cache misses are deduplicated per key via a singleflight.Group: concurrent misses on
// the SAME key share a single in-flight loadFn call and its result, while misses on DIFFERENT keys
// proceed fully in parallel. The mutex guards only the short critical sections that read or write
// the data map — it is never held across loadFn, so a slow lookup for one key cannot serialize
// lookups for other keys (which is exactly what Terraform's parallel resource reads rely on).
//
// Cache is generic over the cached value type so it can be reused for other lookups in the future.
type Cache[T any] struct {
	mu    sync.RWMutex
	data  map[string]T
	group singleflight.Group
}

// NewCache returns an initialized, empty cache.
func NewCache[T any]() *Cache[T] {
	return &Cache[T]{data: make(map[string]T)}
}

// GetOrLoad returns the cached value for key. On a cache miss it calls loadFn,
// stores the result, and returns it. Concurrent misses on the same key are collapsed
// into a single loadFn call whose result is shared by all callers. If loadFn returns an
// error the result is not cached and the error is propagated to the caller.
func (c *Cache[T]) GetOrLoad(key string, loadFn func() (T, error)) (T, error) {
	// Fast path: warm cache hit, read lock only. Skips singleflight entirely.
	c.mu.RLock()
	if v, ok := c.data[key]; ok {
		c.mu.RUnlock()
		return v, nil
	}
	c.mu.RUnlock()

	// Slow path: deduplicate concurrent misses per key. singleflight.Group.Do serializes
	// only callers sharing the same key; different keys run concurrently. We do NOT hold
	// c.mu across loadFn.
	v, err, _ := c.group.Do(key, func() (any, error) {
		// Re-check under the read lock: another caller may have populated the entry
		// after our fast-path miss but before this call began executing.
		c.mu.RLock()
		if v, ok := c.data[key]; ok {
			c.mu.RUnlock()
			return v, nil
		}
		c.mu.RUnlock()

		loaded, loadErr := loadFn()
		if loadErr != nil {
			return nil, loadErr
		}

		c.mu.Lock()
		c.data[key] = loaded
		c.mu.Unlock()
		return loaded, nil
	})
	if err != nil {
		var zero T
		return zero, err
	}
	// Do returns any (it predates generics). Every non-error return path above returns a
	// concrete T, so this assertion is total.
	return v.(T), nil
}

// Invalidate removes the cached result for key, forcing the next GetOrLoad call
// to re-fetch from the source.
func (c *Cache[T]) Invalidate(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
}
