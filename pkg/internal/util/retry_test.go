package util

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_Retry(t *testing.T) {
	sentinel := errors.New("boom")
	tests := []struct {
		name            string
		attempts        int
		f               func(calls *int) (error, bool)
		wantErr         error
		wantErrContains string
		wantCalls       int
	}{
		{
			name:     "succeeds on first attempt",
			attempts: 3,
			f: func(calls *int) (error, bool) {
				*calls++
				return nil, true
			},
			wantCalls: 1,
		},
		{
			name:     "succeeds after retries",
			attempts: 3,
			f: func(calls *int) (error, bool) {
				*calls++
				return nil, *calls >= 2
			},
			wantCalls: 2,
		},
		{
			name:     "gives up after attempts",
			attempts: 3,
			f: func(calls *int) (error, bool) {
				*calls++
				return nil, false
			},
			wantErrContains: "giving up after 3 attempts",
			wantCalls:       3,
		},
		{
			name:     "returns immediate error",
			attempts: 3,
			f: func(calls *int) (error, bool) {
				*calls++
				return sentinel, false
			},
			wantErr:   sentinel,
			wantCalls: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := 0
			err := Retry(tt.attempts, time.Millisecond, func() (error, bool) {
				return tt.f(&calls)
			})
			switch {
			case tt.wantErr != nil:
				require.ErrorIs(t, err, tt.wantErr)
			case tt.wantErrContains != "":
				require.ErrorContains(t, err, tt.wantErrContains)
			default:
				require.NoError(t, err)
			}
			require.Equal(t, tt.wantCalls, calls)
		})
	}
}

func Test_RetryWithContext(t *testing.T) {
	sentinel := errors.New("boom")
	tests := []struct {
		name            string
		ctx             func(t *testing.T) (context.Context, context.CancelFunc)
		sleep           time.Duration
		f               func(calls *int) (error, bool)
		wantErr         error
		wantErrContains string
		wantCalls       int
	}{
		{
			name: "succeeds on first attempt",
			ctx: func(t *testing.T) (context.Context, context.CancelFunc) {
				t.Helper()
				return context.WithTimeout(t.Context(), time.Second)
			},
			sleep: time.Millisecond,
			f: func(calls *int) (error, bool) {
				*calls++
				return nil, true
			},
			wantCalls: 1,
		},
		{
			name: "succeeds after retries",
			ctx: func(t *testing.T) (context.Context, context.CancelFunc) {
				t.Helper()
				return context.WithTimeout(t.Context(), time.Second)
			},
			sleep: time.Millisecond,
			f: func(calls *int) (error, bool) {
				*calls++
				return nil, *calls >= 3
			},
			wantCalls: 3,
		},
		{
			name: "gives up when context deadline exceeded",
			ctx: func(t *testing.T) (context.Context, context.CancelFunc) {
				t.Helper()
				return context.WithTimeout(t.Context(), 20*time.Millisecond)
			},
			sleep: time.Hour,
			f: func(calls *int) (error, bool) {
				*calls++
				return nil, false
			},
			wantErr:         context.DeadlineExceeded,
			wantErrContains: "giving up",
			wantCalls:       1,
		},
		{
			name: "already canceled",
			ctx: func(t *testing.T) (context.Context, context.CancelFunc) {
				t.Helper()
				ctx, cancel := context.WithCancel(t.Context())
				cancel()
				return ctx, cancel
			},
			sleep: time.Millisecond,
			f: func(calls *int) (error, bool) {
				*calls++
				return nil, true
			},
			wantErr:   context.Canceled,
			wantCalls: 0,
		},
		{
			name: "returns immediate error",
			ctx: func(t *testing.T) (context.Context, context.CancelFunc) {
				t.Helper()
				return context.WithTimeout(t.Context(), time.Second)
			},
			sleep: time.Millisecond,
			f: func(calls *int) (error, bool) {
				*calls++
				return sentinel, false
			},
			wantErr:   sentinel,
			wantCalls: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := tt.ctx(t)
			defer cancel()

			calls := 0
			err := RetryWithContext(ctx, tt.sleep, func() (error, bool) {
				return tt.f(&calls)
			})
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			}
			if tt.wantErrContains != "" {
				require.ErrorContains(t, err, tt.wantErrContains)
			}
			if tt.wantErr == nil && tt.wantErrContains == "" {
				require.NoError(t, err)
			}
			require.Equal(t, tt.wantCalls, calls)
		})
	}
}
