package retry

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDo_Success(t *testing.T) {
	ctx := context.Background()
	called := 0

	fn := func(ctx context.Context) error {
		called++
		return nil
	}

	err := Do(ctx, DefaultConfig, fn, nil)
	require.NoError(t, err)
	assert.Equal(t, 1, called, "should be called once")
}

func TestDo_SuccessAfterRetry(t *testing.T) {
	ctx := context.Background()
	called := 0

	fn := func(ctx context.Context) error {
		called++
		if called < 3 {
			return errors.New("temporary error")
		}
		return nil
	}

	config := Config{
		MaxRetries:     3,
		InitialBackoff: 10 * time.Millisecond,
		MaxBackoff:     100 * time.Millisecond,
		Multiplier:     2.0,
	}

	err := Do(ctx, config, fn, nil)
	require.NoError(t, err)
	assert.Equal(t, 3, called, "should be called 3 times")
}

func TestDo_MaxRetriesExceeded(t *testing.T) {
	ctx := context.Background()
	called := 0

	fn := func(ctx context.Context) error {
		called++
		return errors.New("persistent error")
	}

	config := Config{
		MaxRetries:     2,
		InitialBackoff: 10 * time.Millisecond,
		MaxBackoff:     100 * time.Millisecond,
		Multiplier:     2.0,
	}

	err := Do(ctx, config, fn, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "max retries exceeded")
	assert.Equal(t, 3, called, "should be called 3 times (initial + 2 retries)")
}

func TestDo_NonRetryableError(t *testing.T) {
	ctx := context.Background()
	called := 0

	fn := func(ctx context.Context) error {
		called++
		return errors.New("non-retryable error")
	}

	shouldRetry := func(err error) bool {
		return false // このエラーはリトライしない
	}

	config := Config{
		MaxRetries:     3,
		InitialBackoff: 10 * time.Millisecond,
		MaxBackoff:     100 * time.Millisecond,
		Multiplier:     2.0,
	}

	err := Do(ctx, config, fn, shouldRetry)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "non-retryable error")
	assert.Equal(t, 1, called, "should be called only once")
}

func TestDo_ContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	called := 0

	fn := func(ctx context.Context) error {
		called++
		if called == 2 {
			cancel() // 2回目の呼び出し後にキャンセル
		}
		return errors.New("error")
	}

	config := Config{
		MaxRetries:     5,
		InitialBackoff: 10 * time.Millisecond,
		MaxBackoff:     100 * time.Millisecond,
		Multiplier:     2.0,
	}

	err := Do(ctx, config, fn, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "context cancelled")
}

func TestCalculateBackoff(t *testing.T) {
	config := Config{
		InitialBackoff: 1 * time.Second,
		MaxBackoff:     30 * time.Second,
		Multiplier:     2.0,
	}

	tests := []struct {
		attempt  int
		expected time.Duration
	}{
		{0, 1 * time.Second},
		{1, 2 * time.Second},
		{2, 4 * time.Second},
		{3, 8 * time.Second},
		{4, 16 * time.Second},
		{5, 30 * time.Second}, // MaxBackoff
		{6, 30 * time.Second}, // MaxBackoff
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("attempt_%d", tt.attempt), func(t *testing.T) {
			backoff := calculateBackoff(config, tt.attempt)
			assert.Equal(t, tt.expected, backoff)
		})
	}
}
