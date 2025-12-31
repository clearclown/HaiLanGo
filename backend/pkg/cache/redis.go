package cache

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisCache はRedisを使用したキャッシュ実装
type RedisCache struct {
	client *redis.Client
}

// NewRedisCache は新しいRedisキャッシュを作成する
func NewRedisCache(addr, password string, db int) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// 接続テスト
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisCache{
		client: client,
	}, nil
}

// NewRedisCacheFromEnv は環境変数からRedisキャッシュを作成する
func NewRedisCacheFromEnv() (*RedisCache, error) {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}
	password := os.Getenv("REDIS_PASSWORD")

	return NewRedisCache(addr, password, 0)
}

// Get はキーに対応する値を取得する
func (c *RedisCache) Get(ctx context.Context, key string) ([]byte, error) {
	data, err := c.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, &ErrCacheMiss{Key: key}
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get from cache: %w", err)
	}
	return data, nil
}

// Set はキーと値を保存する
func (c *RedisCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	if err := c.client.Set(ctx, key, value, ttl).Err(); err != nil {
		return fmt.Errorf("failed to set cache: %w", err)
	}
	return nil
}

// Delete はキーを削除する
func (c *RedisCache) Delete(ctx context.Context, key string) error {
	if err := c.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete from cache: %w", err)
	}
	return nil
}

// Exists はキーが存在するかチェックする
func (c *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	result, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check existence: %w", err)
	}
	return result > 0, nil
}

// Close はRedis接続を閉じる
func (c *RedisCache) Close() error {
	return c.client.Close()
}

// Ensure RedisCache implements Cache interface
var _ Cache = (*RedisCache)(nil)

// InMemoryCache はインメモリキャッシュ実装（開発・テスト用）
type InMemoryCache struct {
	data map[string]*memoryCacheItem
}

type memoryCacheItem struct {
	value      []byte
	expiration time.Time
}

// NewInMemoryCache は新しいインメモリキャッシュを作成する
func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{
		data: make(map[string]*memoryCacheItem),
	}
}

// Get はキーに対応する値を取得する
func (c *InMemoryCache) Get(ctx context.Context, key string) ([]byte, error) {
	item, found := c.data[key]
	if !found {
		return nil, &ErrCacheMiss{Key: key}
	}

	// 有効期限チェック
	if !item.expiration.IsZero() && time.Now().After(item.expiration) {
		delete(c.data, key)
		return nil, &ErrCacheMiss{Key: key}
	}

	return item.value, nil
}

// Set はキーと値を保存する
func (c *InMemoryCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	var expiration time.Time
	if ttl > 0 {
		expiration = time.Now().Add(ttl)
	}

	c.data[key] = &memoryCacheItem{
		value:      value,
		expiration: expiration,
	}
	return nil
}

// Delete はキーを削除する
func (c *InMemoryCache) Delete(ctx context.Context, key string) error {
	delete(c.data, key)
	return nil
}

// Exists はキーが存在するかチェックする
func (c *InMemoryCache) Exists(ctx context.Context, key string) (bool, error) {
	item, found := c.data[key]
	if !found {
		return false, nil
	}

	// 有効期限チェック
	if !item.expiration.IsZero() && time.Now().After(item.expiration) {
		delete(c.data, key)
		return false, nil
	}

	return true, nil
}

// Ensure InMemoryCache implements Cache interface
var _ Cache = (*InMemoryCache)(nil)

// NewCache は環境変数に基づいて適切なキャッシュを返す
func NewCache() (Cache, error) {
	useMock := os.Getenv("USE_MOCK_APIS") == "true" ||
		os.Getenv("TEST_USE_MOCKS") == "true"

	if useMock {
		return NewInMemoryCache(), nil
	}

	cache, err := NewRedisCacheFromEnv()
	if err != nil {
		// Redisに接続できない場合はインメモリにフォールバック
		return NewInMemoryCache(), nil
	}

	return cache, nil
}
