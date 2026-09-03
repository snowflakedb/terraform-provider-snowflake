package util

import (
	"context"
	"fmt"
	"time"
)

// PollUntil calls fetch until done returns true or ctx is canceled.
func PollUntil[T any](ctx context.Context, interval time.Duration, fetch func() (T, error), done func(T) bool, timeoutMsg string) (T, error) {
	for {
		val, err := fetch()
		if err != nil {
			var zero T
			return zero, err
		}
		if done(val) {
			return val, nil
		}
		select {
		case <-ctx.Done():
			var zero T
			return zero, fmt.Errorf("%s: %w", timeoutMsg, ctx.Err())
		case <-time.After(interval):
		}
	}
}
