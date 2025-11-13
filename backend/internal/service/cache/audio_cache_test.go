package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAudioCacheSet はキャッシュ保存のテスト
func TestAudioCacheSet(t *testing.T) {
	cache := NewAudioCache()
	key := "test-key"
	audioURL := "https://example.com/audio.mp3"
	ttl := 7 * 24 * time.Hour // 7日間

	err := cache.Set(key, audioURL, ttl)
	require.NoError(t, err)
}

// TestAudioCacheGet はキャッシュ取得のテスト
func TestAudioCacheGet(t *testing.T) {
	cache := NewAudioCache()
	key := "test-key"
	audioURL := "https://example.com/audio.mp3"
	ttl := 7 * 24 * time.Hour

	// キャッシュに保存
	err := cache.Set(key, audioURL, ttl)
	require.NoError(t, err)

	// キャッシュから取得
	cachedURL, found := cache.Get(key)
	assert.True(t, found)
	assert.Equal(t, audioURL, cachedURL)
}

// TestAudioCacheNotFound はキャッシュが見つからない場合のテスト
func TestAudioCacheNotFound(t *testing.T) {
	cache := NewAudioCache()
	key := "non-existent-key"

	cachedURL, found := cache.Get(key)
	assert.False(t, found)
	assert.Empty(t, cachedURL)
}

// TestAudioCacheExpiration はキャッシュの有効期限のテスト
func TestAudioCacheExpiration(t *testing.T) {
	cache := NewAudioCache()
	key := "expiring-key"
	audioURL := "https://example.com/audio.mp3"
	ttl := 1 * time.Second // 1秒で期限切れ

	// キャッシュに保存
	err := cache.Set(key, audioURL, ttl)
	require.NoError(t, err)

	// すぐに取得（成功）
	cachedURL, found := cache.Get(key)
	assert.True(t, found)
	assert.Equal(t, audioURL, cachedURL)

	// 2秒待機（期限切れ）
	time.Sleep(2 * time.Second)

	// 取得（失敗）
	cachedURL, found = cache.Get(key)
	assert.False(t, found)
	assert.Empty(t, cachedURL)
}

// TestAudioCacheDelete はキャッシュ削除のテスト
func TestAudioCacheDelete(t *testing.T) {
	cache := NewAudioCache()
	key := "delete-key"
	audioURL := "https://example.com/audio.mp3"
	ttl := 7 * 24 * time.Hour

	// キャッシュに保存
	err := cache.Set(key, audioURL, ttl)
	require.NoError(t, err)

	// 削除
	err = cache.Delete(key)
	require.NoError(t, err)

	// 取得（失敗）
	cachedURL, found := cache.Get(key)
	assert.False(t, found)
	assert.Empty(t, cachedURL)
}

// TestAudioCacheGenerateKey はキャッシュキー生成のテスト
func TestAudioCacheGenerateKey(t *testing.T) {
	cache := NewAudioCache()

	text := "Hello, world!"
	lang := "en"
	quality := "standard"
	speed := 1.0

	key := cache.GenerateKey(text, lang, quality, speed)
	assert.NotEmpty(t, key)

	// 同じパラメータで同じキーが生成される
	key2 := cache.GenerateKey(text, lang, quality, speed)
	assert.Equal(t, key, key2)

	// 異なるパラメータで異なるキーが生成される
	key3 := cache.GenerateKey(text, "ja", quality, speed)
	assert.NotEqual(t, key, key3)
}

// TestAudioCacheHitRate はキャッシュヒット率のテスト
func TestAudioCacheHitRate(t *testing.T) {
	cache := NewAudioCache()

	// 10個のアイテムをキャッシュに保存
	for i := 0; i < 10; i++ {
		key := cache.GenerateKey("text", "en", "standard", float64(i)/10.0)
		err := cache.Set(key, "https://example.com/audio.mp3", 7*24*time.Hour)
		require.NoError(t, err)
	}

	hits := 0
	misses := 0

	// 20回取得を試行（10回はヒット、10回はミス）
	for i := 0; i < 20; i++ {
		key := cache.GenerateKey("text", "en", "standard", float64(i)/10.0)
		_, found := cache.Get(key)
		if found {
			hits++
		} else {
			misses++
		}
	}

	assert.Equal(t, 10, hits)
	assert.Equal(t, 10, misses)

	hitRate := float64(hits) / float64(hits+misses)
	assert.Equal(t, 0.5, hitRate)
}

// TestAudioCacheClear はキャッシュクリアのテスト
func TestAudioCacheClear(t *testing.T) {
	cache := NewAudioCache()

	// 複数のアイテムをキャッシュに保存
	for i := 0; i < 5; i++ {
		key := cache.GenerateKey("text", "en", "standard", float64(i))
		err := cache.Set(key, "https://example.com/audio.mp3", 7*24*time.Hour)
		require.NoError(t, err)
	}

	// クリア
	err := cache.Clear()
	require.NoError(t, err)

	// すべてのアイテムが削除されている
	for i := 0; i < 5; i++ {
		key := cache.GenerateKey("text", "en", "standard", float64(i))
		_, found := cache.Get(key)
		assert.False(t, found)
	}
}
