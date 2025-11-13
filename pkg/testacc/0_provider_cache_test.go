package testacc

import "sync"

// TODO [SNOW-2661409]: add interface instead of any
// providerInitializationCache is a simple cache used throughout the acceptance tests to save time by reusing the initialized providers.
// It's parametrized with the provider context (Terraform Plugin Framework and REST API PoCs have different contexts).
type providerInitializationCache[V any] struct {
	mu   sync.Mutex
	data map[string]V
}

// newProviderInitializationCache creates a new cache with empty map.
func newProviderInitializationCache[V any]() *providerInitializationCache[V] {
	return &providerInitializationCache[V]{
		data: make(map[string]V),
	}
}

// getOrInit provides the already existing initialization for the given key or creates a new one given the initFn.
// The first check is done without locking, as during the tests we will only have a few entries, and we won't clear/re-initialize them.
// If the entry is missing, we lock, check if it was not created in the meantime, and create.
func (c *providerInitializationCache[V]) getOrInit(key string, initFn func() V) V {
	// Return existing if present without locking
	if v, ok := c.data[key]; ok {
		accTestLog.Printf("[DEBUG] Returning cached provider configuration result for key %s", key)
		return v
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	// Return existing if present (it could be added in the meantime)
	if v, ok := c.data[key]; ok {
		accTestLog.Printf("[DEBUG] Returning cached provider configuration result for key %s", key)
		return v
	}
	accTestLog.Printf("[DEBUG] No cached provider configuration found for key %s or caching is not enabled; configuring a new provider", key)

	// Initialize, store, and return
	v := initFn()
	c.data[key] = v
	return v
}
