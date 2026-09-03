package util

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPollUntil(t *testing.T) {
	t.Run("returns immediately when done", func(t *testing.T) {
		got, err := PollUntil(context.Background(), time.Millisecond, func() (int, error) {
			return 7, nil
		}, func(v int) bool { return v == 7 }, "timeout")
		require.NoError(t, err)
		require.Equal(t, 7, got)
	})

	t.Run("returns fetch error", func(t *testing.T) {
		fetchErr := errors.New("fetch failed")
		_, err := PollUntil(context.Background(), time.Millisecond, func() (int, error) {
			return 0, fetchErr
		}, func(int) bool { return false }, "timeout")
		require.ErrorIs(t, err, fetchErr)
	})

	t.Run("returns context timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()
		time.Sleep(50 * time.Millisecond)
		_, err := PollUntil(ctx, time.Millisecond, func() (int, error) {
			return 0, nil
		}, func(int) bool { return false }, "still waiting")
		require.ErrorIs(t, err, context.DeadlineExceeded)
		require.ErrorContains(t, err, "still waiting")
	})
}
