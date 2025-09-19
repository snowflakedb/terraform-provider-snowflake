package common

import (
	"encoding/json"
	"net/http"
)

// TODO [SNOW-2296350]: make it possible to reuse simultaneously from multiple tests (e.g. map per test) (https://go.dev/blog/routing-enhancements)
type DynamicHandler[T any] struct {
	currentValue    T
	defaultValue    T
	replaceWithFunc func(T, T, T) T
}

// TODO [SNOW-2296350]: Log nicer values (use interface)
// TODO [SNOW-2296350]: Handle set/unset instead just opts; update the helper functions and tests while doing it.
func (h *DynamicHandler[T]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		logger.Printf("[DEBUG] Received get request. Current value %v", h.currentValue)
		w.WriteHeader(http.StatusOK)
		// Not handling the error on purpose - it's a test helper, there is no test context here, and we will know on the assertion level.
		_ = json.NewEncoder(w).Encode(h.currentValue)
	case http.MethodPost:
		w.WriteHeader(http.StatusCreated)
		var newValue T
		// Not handling the error on purpose - it's a test helper, there is no test context here, and we will know on the assertion level.
		_ = json.NewDecoder(r.Body).Decode(&newValue)
		logger.Printf("[DEBUG] Received post request. New value %v", newValue)
		h.currentValue = h.replaceWithFunc(h.currentValue, h.defaultValue, newValue)
	}
}

func (h *DynamicHandler[T]) SetCurrentValue(valueProvider T) {
	h.currentValue = valueProvider
}

func NewDynamicHandler[T any]() *DynamicHandler[T] {
	return &DynamicHandler[T]{
		replaceWithFunc: AlwaysReplace[T],
	}
}

func NewDynamicHandlerWithInitialValueAndReplaceWithFunc[T any](initialValue T, replaceWithFunc func(T, T, T) T) *DynamicHandler[T] {
	return &DynamicHandler[T]{
		currentValue:    initialValue,
		replaceWithFunc: replaceWithFunc,
	}
}

func NewDynamicHandlerWithDefaultValueAndReplaceWithFunc[T any](defaultValue T, replaceWithFunc func(T, T, T) T) *DynamicHandler[T] {
	return &DynamicHandler[T]{
		defaultValue:    defaultValue,
		replaceWithFunc: replaceWithFunc,
	}
}

func AlwaysReplace[T any](_ T, _T, replaceWith T) T {
	return replaceWith
}
