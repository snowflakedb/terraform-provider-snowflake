package provider

import (
	"context"
	"sync"

	"golang.org/x/sync/singleflight"
)

// Cache is a simple per-plan in-memory, concurrency-safe cache keyed by string.
//
// It is intended for caching expensive, read-only Snowflake lookups whose result is shared
// by many resource instances within a single Terraform plan/apply cycle. The canonical use is
// SHOW GRANTS: without caching, every grant resource instance (snowflake_grant_account_role,
// snowflake_grant_privileges_to_account_role, snowflake_grant_ownership, ...) issues an independent
// SHOW GRANTS call during Read. Because one SHOW returns the full result set for whatever it's
// scoped to (a role, an object, a container), N resource instances that resolve to the same SHOW
// statement trigger N identical round-trips — only 1 is needed per plan. Callers key the cache by
// the rendered SQL of the SHOW statement (see sdk.ShowGrantOptionsToSQL), so identical keys are
// guaranteed to represent identical queries.
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
	mu      sync.RWMutex
	data    map[string]T
	group   singleflight.Group
	cancels map[string]context.CancelFunc
}

// NewCache returns an initialized, empty cache.
func NewCache[T any]() *Cache[T] {
	return &Cache[T]{data: make(map[string]T), cancels: make(map[string]context.CancelFunc)}
}

// GetOrLoad returns the cached value for key. On a cache miss it calls loadFn,
// stores the result, and returns it. Concurrent misses on the same key are collapsed
// into a single loadFn call whose result is shared by all callers. If loadFn returns an
// error the result is not cached and the error is propagated to the caller.
//
// ctx is the caller's Terraform-provided context (so resource Timeouts and plan cancellation
// propagate into the Snowflake call on a cache miss). Because concurrent misses on the same key
// are collapsed into a single loadFn call via singleflight, only the ctx of whichever caller
// triggers the load actually governs it; a later caller joining an in-flight load is not able to
// cancel it via its own ctx, only via Invalidate.
func (c *Cache[T]) GetOrLoad(ctx context.Context, key string, loadFn func(ctx context.Context) (T, error)) (T, error) {
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

		loadCtx, cancel := context.WithCancel(ctx)
		c.mu.Lock()
		c.cancels[key] = cancel
		c.mu.Unlock()
		defer func() {
			c.mu.Lock()
			delete(c.cancels, key)
			c.mu.Unlock()
			cancel()
		}()

		loaded, loadErr := loadFn(loadCtx)
		if loadErr != nil {
			return nil, loadErr
		}
		if loadCtx.Err() != nil {
			// Invalidate ran while loadFn was in flight; the result may reflect
			// pre-mutation state, so don't let it overwrite the invalidation.
			return nil, loadCtx.Err()
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
	delete(c.data, key)
	cancel, ok := c.cancels[key]
	c.mu.Unlock()
	if ok {
		cancel()
	}
	c.group.Forget(key)
}
