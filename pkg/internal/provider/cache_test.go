package provider

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// waitTimeout waits for wg with an upper bound, failing the test instead of hanging CI
// forever if the goroutines never complete (e.g. because concurrent misses got serialized
// behind a shared lock and deadlocked on a rendezvous).
func waitTimeout(t *testing.T, wg *sync.WaitGroup, timeout time.Duration) {
	t.Helper()
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(timeout):
		t.Fatalf("timed out after %s waiting for GetOrLoad calls to complete", timeout)
	}
}

func TestCache_MissCallsLoadFn(t *testing.T) {
	cache := NewCache[[]sdk.Grant]()
	calls := 0
	grants := []sdk.Grant{{GrantedTo: sdk.ObjectTypeRole}}

	result, err := cache.GetOrLoad("ROLE_A", func(context.Context) ([]sdk.Grant, error) {
		calls++
		return grants, nil
	})

	require.NoError(t, err)
	assert.Equal(t, grants, result)
	assert.Equal(t, 1, calls)
}

func TestCache_HitSkipsLoadFn(t *testing.T) {
	cache := NewCache[[]sdk.Grant]()
	calls := 0
	grants := []sdk.Grant{{GrantedTo: sdk.ObjectTypeRole}}

	for range 3 {
		result, err := cache.GetOrLoad("ROLE_A", func(context.Context) ([]sdk.Grant, error) {
			calls++
			return grants, nil
		})
		require.NoError(t, err)
		assert.Equal(t, grants, result)
	}
	assert.Equal(t, 1, calls, "loadFn should be called exactly once for repeated reads of the same key")
}

func TestCache_InvalidateForcesMiss(t *testing.T) {
	cache := NewCache[[]sdk.Grant]()
	calls := 0

	_, _ = cache.GetOrLoad("ROLE_A", func(context.Context) ([]sdk.Grant, error) { calls++; return nil, nil })
	cache.Invalidate("ROLE_A")
	_, _ = cache.GetOrLoad("ROLE_A", func(context.Context) ([]sdk.Grant, error) { calls++; return nil, nil })

	assert.Equal(t, 2, calls, "loadFn should be called again after invalidation")
}

func TestCache_InvalidateUnknownKeyIsNoop(t *testing.T) {
	cache := NewCache[[]sdk.Grant]()
	assert.NotPanics(t, func() { cache.Invalidate("NONEXISTENT") })
}

func TestCache_ErrorNotCached(t *testing.T) {
	cache := NewCache[[]sdk.Grant]()
	calls := 0
	boom := errors.New("snowflake unavailable")

	for range 2 {
		_, err := cache.GetOrLoad("ROLE_A", func(context.Context) ([]sdk.Grant, error) {
			calls++
			return nil, boom
		})
		assert.ErrorIs(t, err, boom)
	}
	assert.Equal(t, 2, calls, "errors must not be cached; loadFn must be called on every miss")
}

func TestCache_KeysAreIndependent(t *testing.T) {
	cache := NewCache[[]sdk.Grant]()
	grantsA := []sdk.Grant{{GrantedTo: sdk.ObjectTypeRole}}
	grantsB := []sdk.Grant{{GrantedTo: sdk.ObjectTypeUser}}

	resultA, _ := cache.GetOrLoad("ROLE_A", func(context.Context) ([]sdk.Grant, error) { return grantsA, nil })
	resultB, _ := cache.GetOrLoad("ROLE_B", func(context.Context) ([]sdk.Grant, error) { return grantsB, nil })

	assert.Equal(t, grantsA, resultA)
	assert.Equal(t, grantsB, resultB)

	// Second load for A must still return A's grants, not B's.
	resultA2, _ := cache.GetOrLoad("ROLE_A", func(context.Context) ([]sdk.Grant, error) { return grantsB, nil })
	assert.Equal(t, grantsA, resultA2)
}

func TestCache_ConcurrentReadsAreSafe(t *testing.T) {
	cache := NewCache[[]sdk.Grant]()
	grants := []sdk.Grant{{GrantedTo: sdk.ObjectTypeRole}}
	// Prime the cache.
	_, _ = cache.GetOrLoad("ROLE_A", func(context.Context) ([]sdk.Grant, error) { return grants, nil })

	var wg sync.WaitGroup
	for range 100 {
		wg.Go(func() {
			result, err := cache.GetOrLoad("ROLE_A", func(context.Context) ([]sdk.Grant, error) { return grants, nil })
			assert.NoError(t, err)
			assert.Equal(t, grants, result)
		})
	}
	wg.Wait()
}

func TestCache_ConcurrentWritesAreSafe(t *testing.T) {
	cache := NewCache[[]sdk.Grant]()
	grants := []sdk.Grant{{GrantedTo: sdk.ObjectTypeRole}}

	var wg sync.WaitGroup
	for i := range 50 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_, _ = cache.GetOrLoad("ROLE_A", func(context.Context) ([]sdk.Grant, error) { return grants, nil })
			if i%5 == 0 {
				cache.Invalidate("ROLE_A")
			}
		}(i)
	}
	wg.Wait()
}

// TestCache_GenericValueType verifies the cache works for value types other than []sdk.Grant,
// confirming the parametrization is genuinely generic and reusable for future lookups.
func TestCache_GenericValueType(t *testing.T) {
	cache := NewCache[int]()
	calls := 0

	for range 3 {
		result, err := cache.GetOrLoad("answer", func(context.Context) (int, error) {
			calls++
			return 42, nil
		})
		require.NoError(t, err)
		assert.Equal(t, 42, result)
	}
	assert.Equal(t, 1, calls, "loadFn should be called exactly once for repeated reads of the same key")

	// Error path returns the zero value of the type parameter.
	boom := errors.New("boom")
	result, err := cache.GetOrLoad("missing", func(context.Context) (int, error) { return 7, boom })
	assert.ErrorIs(t, err, boom)
	assert.Equal(t, 0, result, "error path must return the zero value of T")
}

// TestCache_ConcurrentMissesOnDifferentKeysRunInParallel proves that cache misses on
// distinct keys are NOT serialized behind a shared lock. Each loadFn signals that it has
// started, then blocks on a barrier until every loadFn has started. This can only complete
// if all loadFns run concurrently. With a single global write lock held across loadFn, only
// one goroutine reaches loadFn while the rest block on the lock, the barrier never fills, and
// the test fails via the bounded timeout instead of hanging CI forever.
func TestCache_ConcurrentMissesOnDifferentKeysRunInParallel(t *testing.T) {
	cache := NewCache[int]()
	const n = 8

	started := make(chan struct{}, n)
	release := make(chan struct{})

	var wg sync.WaitGroup
	for i := range n {
		wg.Go(func() {
			key := fmt.Sprintf("ROLE_%d", i)
			_, err := cache.GetOrLoad(key, func(context.Context) (int, error) {
				started <- struct{}{} // I have started running.
				<-release             // Block until everyone has started.
				return i, nil
			})
			assert.NoError(t, err)
		})
	}

	// Barrier: once all n loadFns are confirmed running, release them together.
	go func() {
		for range n {
			<-started
		}
		close(release)
	}()

	waitTimeout(t, &wg, 30*time.Second)
}

// TestCache_ConcurrentMissesOnSameKeyCollapseToOneLoad proves that concurrent misses on the
// SAME key share a single loadFn invocation, under real concurrency (not just the sequential
// case already covered by TestCache_HitSkipsLoadFn). The in-flight call is held open until at
// least one loadFn is confirmed running, widening the window for a duplicate load if
// deduplication were broken.
func TestCache_ConcurrentMissesOnSameKeyCollapseToOneLoad(t *testing.T) {
	cache := NewCache[int]()
	const n = 50

	var calls atomic.Int64
	started := make(chan struct{}, n) // buffered so a broken (multi-call) impl never blocks on send
	proceed := make(chan struct{})

	var wg sync.WaitGroup
	for range n {
		wg.Go(func() {
			result, err := cache.GetOrLoad("ROLE_A", func(context.Context) (int, error) {
				calls.Add(1)
				started <- struct{}{}
				<-proceed // keep the in-flight call open to widen the dedup window
				return 42, nil
			})
			assert.NoError(t, err)
			assert.Equal(t, 42, result)
		})
	}

	// Wait until at least one loadFn is in flight, then release it (and everyone sharing it).
	<-started
	close(proceed)

	waitTimeout(t, &wg, 30*time.Second)
	assert.Equal(t, int64(1), calls.Load(),
		"concurrent misses on the same key must collapse to exactly one loadFn call")
}
