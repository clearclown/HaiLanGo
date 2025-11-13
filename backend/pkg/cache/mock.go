package cache

import (
	"context"
	"sync"
	"time"
)

// MockCache はモックキャッシュ実装
type MockCache struct {
	mu    sync.RWMutex
	data  map[string]cacheEntry
}

type cacheEntry struct {
	value      []byte
	expiration time.Time
}

// NewMockCache は新しいモックキャッシュを作成する
func NewMockCache() *MockCache {
	return &MockCache{
		data: make(map[string]cacheEntry),
	}
}

// Get はキーに対応する値を取得する
func (m *MockCache) Get(ctx context.Context, key string) ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	entry, ok := m.data[key]
	if !ok {
		return nil, &ErrCacheMiss{Key: key}
	}

	// 有効期限チェック
	if !entry.expiration.IsZero() && time.Now().After(entry.expiration) {
		return nil, &ErrCacheMiss{Key: key}
	}

	return entry.value, nil
}

// Set はキーと値を保存する
func (m *MockCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	entry := cacheEntry{
		value: value,
	}

	if ttl > 0 {
		entry.expiration = time.Now().Add(ttl)
	}

	m.data[key] = entry
	return nil
}

// Delete はキーを削除する
func (m *MockCache) Delete(ctx context.Context, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.data, key)
	return nil
}

// Exists はキーが存在するかチェックする
func (m *MockCache) Exists(ctx context.Context, key string) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	entry, ok := m.data[key]
	if !ok {
		return false, nil
	}

	// 有効期限チェック
	if !entry.expiration.IsZero() && time.Now().After(entry.expiration) {
		return false, nil
	}

	return true, nil
}

// Clear はすべてのキャッシュをクリアする（テスト用）
func (m *MockCache) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.data = make(map[string]cacheEntry)
}
