package retry

import (
	"context"
	"fmt"
	"math"
	"time"
)

// Config はリトライ設定
type Config struct {
	MaxRetries     int           // 最大リトライ回数
	InitialBackoff time.Duration // 初期バックオフ時間
	MaxBackoff     time.Duration // 最大バックオフ時間
	Multiplier     float64       // バックオフ倍率
}

// DefaultConfig はデフォルトのリトライ設定
var DefaultConfig = Config{
	MaxRetries:     3,
	InitialBackoff: 1 * time.Second,
	MaxBackoff:     30 * time.Second,
	Multiplier:     2.0,
}

// RetryFunc はリトライ可能な関数の型
type RetryFunc func(ctx context.Context) error

// ShouldRetryFunc はリトライすべきかを判定する関数の型
type ShouldRetryFunc func(err error) bool

// Do は指数バックオフでリトライを実行する
func Do(ctx context.Context, config Config, fn RetryFunc, shouldRetry ShouldRetryFunc) error {
	var lastErr error

	for attempt := 0; attempt <= config.MaxRetries; attempt++ {
		// コンテキストがキャンセルされていないかチェック
		select {
		case <-ctx.Done():
			return fmt.Errorf("context cancelled: %w", ctx.Err())
		default:
		}

		// 関数を実行
		err := fn(ctx)
		if err == nil {
			return nil // 成功
		}

		lastErr = err

		// リトライすべきかチェック
		if shouldRetry != nil && !shouldRetry(err) {
			return fmt.Errorf("non-retryable error: %w", err)
		}

		// 最後の試行の場合はリトライしない
		if attempt == config.MaxRetries {
			break
		}

		// バックオフ時間を計算
		backoff := calculateBackoff(config, attempt)

		// バックオフ待機
		select {
		case <-ctx.Done():
			return fmt.Errorf("context cancelled during backoff: %w", ctx.Err())
		case <-time.After(backoff):
			// 次の試行へ
		}
	}

	return fmt.Errorf("max retries exceeded (%d attempts): %w", config.MaxRetries+1, lastErr)
}

// calculateBackoff は指数バックオフ時間を計算する
func calculateBackoff(config Config, attempt int) time.Duration {
	backoff := float64(config.InitialBackoff) * math.Pow(config.Multiplier, float64(attempt))

	if backoff > float64(config.MaxBackoff) {
		return config.MaxBackoff
	}

	return time.Duration(backoff)
}
