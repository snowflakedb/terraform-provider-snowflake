package common

import (
	"encoding/json"
	"net/http"
)

type DynamicHandler[T any] struct {
	currentValue T
}

func (h *DynamicHandler[T]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.WriteHeader(http.StatusOK)
		resp, _ := json.Marshal(h.currentValue)
		_, _ = w.Write(resp)
	case http.MethodPost:
		w.WriteHeader(http.StatusCreated)
		var newValue T
		_ = json.NewDecoder(r.Body).Decode(&newValue)
		h.currentValue = newValue
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
