package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"sync"
	"time"
)

// AudioCache は音声キャッシュのインターフェース
type AudioCache interface {
	Get(key string) (string, bool)
	Set(key string, audioURL string, ttl time.Duration) error
	Delete(key string) error
	GenerateKey(text string, lang string, quality string, speed float64) string
	Clear() error
}

// InMemoryAudioCache はインメモリ音声キャッシュ（開発・テスト用）
type InMemoryAudioCache struct {
	data      map[string]*cacheItem
	mu        sync.RWMutex
	useMock   bool
	redisMode bool
}

type cacheItem struct {
	value      string
	expiration time.Time
}

// NewAudioCache は新しい音声キャッシュを作成
func NewAudioCache() AudioCache {
	useMock := os.Getenv("USE_MOCK_APIS") == "true" ||
		os.Getenv("TEST_USE_MOCKS") == "true"

	cache := &InMemoryAudioCache{
		data:      make(map[string]*cacheItem),
		useMock:   useMock,
		redisMode: false, // TODO: Redis実装時にtrue
	}

	// バックグラウンドで期限切れアイテムを削除
	go cache.cleanupExpired()

	return cache
}

// Get はキャッシュから値を取得
func (c *InMemoryAudioCache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.data[key]
	if !found {
		return "", false
	}

	// 有効期限チェック
	if time.Now().After(item.expiration) {
		return "", false
	}

	return item.value, true
}

// Set はキャッシュに値を保存
func (c *InMemoryAudioCache) Set(key string, audioURL string, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = &cacheItem{
		value:      audioURL,
		expiration: time.Now().Add(ttl),
	}

	return nil
}

// Delete はキャッシュから値を削除
func (c *InMemoryAudioCache) Delete(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.data, key)
	return nil
}

// GenerateKey はキャッシュキーを生成
func (c *InMemoryAudioCache) GenerateKey(text string, lang string, quality string, speed float64) string {
	// ハッシュを使用してキーを生成
	data := fmt.Sprintf("tts:%s:%s:%s:%.2f", text, lang, quality, speed)
	hash := sha256.Sum256([]byte(data))
	return "audio:" + hex.EncodeToString(hash[:])
}

// Clear はすべてのキャッシュをクリア
func (c *InMemoryAudioCache) Clear() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data = make(map[string]*cacheItem)
	return nil
}

// cleanupExpired は期限切れアイテムを定期的に削除
func (c *InMemoryAudioCache) cleanupExpired() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for key, item := range c.data {
			if now.After(item.expiration) {
				delete(c.data, key)
			}
		}
		c.mu.Unlock()
	}
}

// RedisAudioCache はRedis音声キャッシュ（本番用）
// TODO: Redis実装
type RedisAudioCache struct {
	// redis client
}

// NewRedisAudioCache は新しいRedis音声キャッシュを作成
func NewRedisAudioCache() AudioCache {
	// TODO: Redis実装
	return NewAudioCache() // 現在はインメモリを返す
}
