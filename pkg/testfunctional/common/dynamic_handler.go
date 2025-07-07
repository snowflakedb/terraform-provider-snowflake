package common

import (
	"encoding/json"
	"net/http"
)

// TODO [mux-PRs]: make it possible to reuse simultaneously from multiple tests (e.g. map per test)
type DynamicHandler[T any] struct {
	currentValue    T
	replaceWithFunc func(T, T) T
}

func (h *DynamicHandler[T]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(h.currentValue)
	case http.MethodPost:
		w.WriteHeader(http.StatusCreated)
		var newValue T
		_ = json.NewDecoder(r.Body).Decode(&newValue)
		h.currentValue = h.replaceWithFunc(h.currentValue, newValue)
	}
}

func (h *DynamicHandler[T]) SetCurrentValue(valueProvider T) {
	h.currentValue = valueProvider
}

func NewDynamicHandler[T any]() *DynamicHandler[T] {
	return &DynamicHandler[T]{}
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
