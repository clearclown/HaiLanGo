package cache

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMockCache_SetAndGet(t *testing.T) {
	ctx := context.Background()
	cache := NewMockCache()

	key := "test_key"
	value := []byte("test value")

	// Set
	err := cache.Set(ctx, key, value, 0)
	require.NoError(t, err)

	// Get
	got, err := cache.Get(ctx, key)
	require.NoError(t, err)
	assert.Equal(t, value, got)
}

func TestMockCache_GetMiss(t *testing.T) {
	ctx := context.Background()
	cache := NewMockCache()

	// 存在しないキー
	_, err := cache.Get(ctx, "nonexistent")
	require.Error(t, err)
	assert.True(t, IsCacheMiss(err))
}

func TestMockCache_Delete(t *testing.T) {
	ctx := context.Background()
	cache := NewMockCache()

	key := "test_key"
	value := []byte("test value")

	// Set
	err := cache.Set(ctx, key, value, 0)
	require.NoError(t, err)

	// Delete
	err = cache.Delete(ctx, key)
	require.NoError(t, err)

	// Get after delete
	_, err = cache.Get(ctx, key)
	require.Error(t, err)
	assert.True(t, IsCacheMiss(err))
}

func TestMockCache_Exists(t *testing.T) {
	ctx := context.Background()
	cache := NewMockCache()

	key := "test_key"
	value := []byte("test value")

	// 存在しない場合
	exists, err := cache.Exists(ctx, key)
	require.NoError(t, err)
	assert.False(t, exists)

	// Set
	err = cache.Set(ctx, key, value, 0)
	require.NoError(t, err)

	// 存在する場合
	exists, err = cache.Exists(ctx, key)
	require.NoError(t, err)
	assert.True(t, exists)
}

func TestMockCache_TTL(t *testing.T) {
	ctx := context.Background()
	cache := NewMockCache()

	key := "test_key"
	value := []byte("test value")
	ttl := 50 * time.Millisecond

	// Set with TTL
	err := cache.Set(ctx, key, value, ttl)
	require.NoError(t, err)

	// すぐに取得できる
	got, err := cache.Get(ctx, key)
	require.NoError(t, err)
	assert.Equal(t, value, got)

	// TTL後は取得できない
	time.Sleep(100 * time.Millisecond)
	_, err = cache.Get(ctx, key)
	require.Error(t, err)
	assert.True(t, IsCacheMiss(err))
}

func TestMockCache_TTL_Exists(t *testing.T) {
	ctx := context.Background()
	cache := NewMockCache()

	key := "test_key"
	value := []byte("test value")
	ttl := 50 * time.Millisecond

	// Set with TTL
	err := cache.Set(ctx, key, value, ttl)
	require.NoError(t, err)

	// すぐに存在する
	exists, err := cache.Exists(ctx, key)
	require.NoError(t, err)
	assert.True(t, exists)

	// TTL後は存在しない
	time.Sleep(100 * time.Millisecond)
	exists, err = cache.Exists(ctx, key)
	require.NoError(t, err)
	assert.False(t, exists)
}

func TestMockCache_Clear(t *testing.T) {
	ctx := context.Background()
	cache := NewMockCache()

	// 複数のキーを設定
	cache.Set(ctx, "key1", []byte("value1"), 0)
	cache.Set(ctx, "key2", []byte("value2"), 0)
	cache.Set(ctx, "key3", []byte("value3"), 0)

	// Clear
	cache.Clear()

	// すべてのキーが削除される
	_, err := cache.Get(ctx, "key1")
	assert.True(t, IsCacheMiss(err))

	_, err = cache.Get(ctx, "key2")
	assert.True(t, IsCacheMiss(err))

	_, err = cache.Get(ctx, "key3")
	assert.True(t, IsCacheMiss(err))
}

func TestMockCache_Concurrent(t *testing.T) {
	ctx := context.Background()
	cache := NewMockCache()

	// 並行アクセスのテスト
	done := make(chan bool)

	// 書き込みゴルーチン
	go func() {
		for i := 0; i < 100; i++ {
			cache.Set(ctx, "key", []byte("value"), 0)
		}
		done <- true
	}()

	// 読み込みゴルーチン
	go func() {
		for i := 0; i < 100; i++ {
			cache.Get(ctx, "key")
		}
		done <- true
	}()

	// 完了を待つ
	<-done
	<-done
}

func TestErrCacheMiss(t *testing.T) {
	err := &ErrCacheMiss{Key: "test_key"}
	assert.Equal(t, "cache miss for key: test_key", err.Error())

	// IsCacheMiss
	assert.True(t, IsCacheMiss(err))
	assert.False(t, IsCacheMiss(assert.AnError))
}
