package provider

import (
	"context"
	"fmt"
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
// the rendered SQL of the SHOW statement (see sdk.StructToSQL), so identical keys are
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
// Invalidate never cancels a load already in flight for its key: that load may be shared (via
// singleflight) with callers unrelated to whatever mutation triggered the invalidation, and turning
// their read into an error would be worse than letting it complete normally. Instead, each key has a
// generation counter that Invalidate bumps. The generation is folded into the singleflight key, so a
// concurrent Invalidate simply causes the in-flight load's result to go uncached (a subsequent
// GetOrLoad reloads under the new generation) rather than forcing it to fail.
//
// Cache is generic over the cached value type so it can be reused for other lookups in the future.
type Cache[T any] struct {
	mu         sync.RWMutex
	data       map[string]T
	generation map[string]uint64
	group      singleflight.Group
}

// NewCache returns an initialized, empty cache.
func NewCache[T any]() *Cache[T] {
	return &Cache[T]{data: make(map[string]T), generation: make(map[string]uint64)}
}

// GetOrLoad returns the cached value for key. On a cache miss it calls loadFn, stores the
// result, and returns it. Concurrent misses on the same key and generation (see Invalidate) are
// collapsed into a single loadFn call whose result is shared by all callers. If loadFn returns an
// error the result is not cached and the error is propagated to the caller.
//
// ctx is passed to loadFn unchanged: canceling it (e.g. a resource Timeout or plan cancellation)
// cancels the load for every caller currently sharing it, exactly as it would for a direct call
// to loadFn.
func (c *Cache[T]) GetOrLoad(ctx context.Context, key string, loadFn func(ctx context.Context) (T, error)) (T, error) {
	// Fast path: warm cache hit, read lock only. Skips singleflight entirely.
	c.mu.RLock()
	if v, ok := c.data[key]; ok {
		c.mu.RUnlock()
		return v, nil
	}
	gen := c.generation[key]
	c.mu.RUnlock()

	// Folding the generation into the singleflight key means an Invalidate landing between
	// here and the Do call below starts a fresh singleflight group for key, rather than
	// joining (or having to cancel) whatever's already in flight for it.
	sfKey := fmt.Sprintf("%d:%s", gen, key)

	v, err, _ := c.group.Do(sfKey, func() (any, error) {
		// Re-check under the read lock: another caller may have populated the entry
		// after our fast-path miss but before this call began executing.
		c.mu.RLock()
		if v, ok := c.data[key]; ok {
			c.mu.RUnlock()
			return v, nil
		}
		c.mu.RUnlock()

		loaded, loadErr := loadFn(ctx)
		if loadErr != nil {
			return nil, loadErr
		}

		c.mu.Lock()
		// Only cache if no Invalidate landed on key since gen was read above; otherwise this
		// result may reflect pre-mutation state, so leave it uncached and let the next
		// GetOrLoad start a fresh load under the new generation.
		if c.generation[key] == gen {
			c.data[key] = loaded
		}
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

// Invalidate removes the cached result for key and bumps its generation, forcing the next
// GetOrLoad call to start a fresh load rather than reusing one already in flight.
func (c *Cache[T]) Invalidate(key string) {
	c.mu.Lock()
	delete(c.data, key)
	c.generation[key]++
	c.mu.Unlock()
}
