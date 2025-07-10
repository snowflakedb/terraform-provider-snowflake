package common

import (
	"encoding/json"
	"net/http"
)

// TODO [mux-PRs]: make it possible to reuse simultaneously from multiple tests (e.g. map per test)
// TODO [mux-PRs]: https://go.dev/blog/routing-enhancements
type DynamicHandler[T any] struct {
	currentValue    T
	replaceWithFunc func(T, T) T
}

// TODO [mux-PRs] Log nicer values (use interface)
func (h *DynamicHandler[T]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		logger.Printf("[DEBUG] Received get request. Current value %v", h.currentValue)
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(h.currentValue)
	case http.MethodPost:
		w.WriteHeader(http.StatusCreated)
		var newValue T
		_ = json.NewDecoder(r.Body).Decode(&newValue)
		logger.Printf("[DEBUG] Received post request. New value %v", newValue)
		h.currentValue = h.replaceWithFunc(h.currentValue, newValue)
	}
}

func (h *DynamicHandler[T]) SetCurrentValue(valueProvider T) {
	h.currentValue = valueProvider
}

func NewDynamicHandler[T any]() *DynamicHandler[T] {
	return &DynamicHandler[T]{
		replaceWithFunc: func(_ T, t2 T) T {
			return t2
		},
	}
}

func NewDynamicHandlerWithInitialValue[T any](initialValue T) *DynamicHandler[T] {
	return &DynamicHandler[T]{
		currentValue: initialValue,
	}
}

func NewDynamicHandlerWithInitialValueAndReplaceWithFunc[T any](initialValue T, replaceWithFunc func(T, T) T) *DynamicHandler[T] {
	return &DynamicHandler[T]{
		currentValue:    initialValue,
		replaceWithFunc: replaceWithFunc,
	}
}
