package provider

import (
	"errors"
	"sync"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShowGrantsOfRoleCache_MissCallsLoadFn(t *testing.T) {
	cache := NewShowGrantsOfRoleCache()
	calls := 0
	grants := []sdk.Grant{{GrantedTo: sdk.ObjectTypeRole}}

	result, err := cache.GetOrLoad("ROLE_A", func() ([]sdk.Grant, error) {
		calls++
		return grants, nil
	})

	require.NoError(t, err)
	assert.Equal(t, grants, result)
	assert.Equal(t, 1, calls)
}

func TestShowGrantsOfRoleCache_HitSkipsLoadFn(t *testing.T) {
	cache := NewShowGrantsOfRoleCache()
	calls := 0
	grants := []sdk.Grant{{GrantedTo: sdk.ObjectTypeRole}}

	for range 3 {
		result, err := cache.GetOrLoad("ROLE_A", func() ([]sdk.Grant, error) {
			calls++
			return grants, nil
		})
		require.NoError(t, err)
		assert.Equal(t, grants, result)
	}
	assert.Equal(t, 1, calls, "loadFn should be called exactly once for repeated reads of the same key")
}

func TestShowGrantsOfRoleCache_InvalidateForcesMiss(t *testing.T) {
	cache := NewShowGrantsOfRoleCache()
	calls := 0

	_, _ = cache.GetOrLoad("ROLE_A", func() ([]sdk.Grant, error) { calls++; return nil, nil })
	cache.Invalidate("ROLE_A")
	_, _ = cache.GetOrLoad("ROLE_A", func() ([]sdk.Grant, error) { calls++; return nil, nil })

	assert.Equal(t, 2, calls, "loadFn should be called again after invalidation")
}

func TestShowGrantsOfRoleCache_InvalidateUnknownKeyIsNoop(t *testing.T) {
	cache := NewShowGrantsOfRoleCache()
	assert.NotPanics(t, func() { cache.Invalidate("NONEXISTENT") })
}

func TestShowGrantsOfRoleCache_ErrorNotCached(t *testing.T) {
	cache := NewShowGrantsOfRoleCache()
	calls := 0
	boom := errors.New("snowflake unavailable")

	for range 2 {
		_, err := cache.GetOrLoad("ROLE_A", func() ([]sdk.Grant, error) {
			calls++
			return nil, boom
		})
		assert.ErrorIs(t, err, boom)
	}
	assert.Equal(t, 2, calls, "errors must not be cached; loadFn must be called on every miss")
}

func TestShowGrantsOfRoleCache_KeysAreIndependent(t *testing.T) {
	cache := NewShowGrantsOfRoleCache()
	grantsA := []sdk.Grant{{GrantedTo: sdk.ObjectTypeRole}}
	grantsB := []sdk.Grant{{GrantedTo: sdk.ObjectTypeUser}}

	resultA, _ := cache.GetOrLoad("ROLE_A", func() ([]sdk.Grant, error) { return grantsA, nil })
	resultB, _ := cache.GetOrLoad("ROLE_B", func() ([]sdk.Grant, error) { return grantsB, nil })

	assert.Equal(t, grantsA, resultA)
	assert.Equal(t, grantsB, resultB)

	// Second load for A must still return A's grants, not B's.
	resultA2, _ := cache.GetOrLoad("ROLE_A", func() ([]sdk.Grant, error) { return grantsB, nil })
	assert.Equal(t, grantsA, resultA2)
}

func TestShowGrantsOfRoleCache_ConcurrentReadsAreSafe(t *testing.T) {
	cache := NewShowGrantsOfRoleCache()
	grants := []sdk.Grant{{GrantedTo: sdk.ObjectTypeRole}}
	// Prime the cache.
	_, _ = cache.GetOrLoad("ROLE_A", func() ([]sdk.Grant, error) { return grants, nil })

	var wg sync.WaitGroup
	for range 100 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			result, err := cache.GetOrLoad("ROLE_A", func() ([]sdk.Grant, error) { return grants, nil })
			assert.NoError(t, err)
			assert.Equal(t, grants, result)
		}()
	}
	wg.Wait()
}

func TestShowGrantsOfRoleCache_ConcurrentWritesAreSafe(t *testing.T) {
	cache := NewShowGrantsOfRoleCache()
	grants := []sdk.Grant{{GrantedTo: sdk.ObjectTypeRole}}

	var wg sync.WaitGroup
	for i := range 50 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_, _ = cache.GetOrLoad("ROLE_A", func() ([]sdk.Grant, error) { return grants, nil })
			if i%5 == 0 {
				cache.Invalidate("ROLE_A")
			}
		}(i)
	}
	wg.Wait()
}
