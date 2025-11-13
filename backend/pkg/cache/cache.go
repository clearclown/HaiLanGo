package cache

import (
	"context"
	"time"
)

// Cache はキャッシュのインターフェース
type Cache interface {
	// Get はキーに対応する値を取得する
	Get(ctx context.Context, key string) ([]byte, error)

	// Set はキーと値を保存する
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error

	// Delete はキーを削除する
	Delete(ctx context.Context, key string) error

	// Exists はキーが存在するかチェックする
	Exists(ctx context.Context, key string) (bool, error)
}

// ErrCacheMiss はキャッシュミスを示すエラー
type ErrCacheMiss struct {
	Key string
}

func (e *ErrCacheMiss) Error() string {
	return "cache miss for key: " + e.Key
}

// IsCacheMiss はエラーがキャッシュミスかどうかをチェックする
func IsCacheMiss(err error) bool {
	_, ok := err.(*ErrCacheMiss)
	return ok
}
