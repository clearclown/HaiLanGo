package tts

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	// テスト実行時は自動的にモックを使用
	if err := os.Setenv("TEST_USE_MOCKS", "true"); err != nil {
		panic(err)
	}
	if err := os.Setenv("USE_MOCK_APIS", "true"); err != nil {
		panic(err)
	}
	code := m.Run()
	os.Exit(code)
}

// TestGenerateAudio は音声生成のテスト
func TestGenerateAudio(t *testing.T) {
	ctx := context.Background()
	service := NewTTSService()

	text := "Hello, world!"
	lang := "en"
	quality := "standard"
	speed := 1.0

	audioURL, err := service.GenerateAudio(ctx, text, lang, quality, speed)

	require.NoError(t, err)
	assert.NotEmpty(t, audioURL)
	assert.Contains(t, audioURL, "http")
}

// TestCacheHit はキャッシュヒットのテスト
func TestCacheHit(t *testing.T) {
	ctx := context.Background()
	service := NewTTSService()

	text := "Cache test"
	lang := "en"
	quality := "standard"
	speed := 1.0

	// 1回目: キャッシュミス
	startTime1 := time.Now()
	audioURL1, err := service.GenerateAudio(ctx, text, lang, quality, speed)
	duration1 := time.Since(startTime1)
	require.NoError(t, err)
	assert.NotEmpty(t, audioURL1)

	// 2回目: キャッシュヒット（高速）
	startTime2 := time.Now()
	audioURL2, err := service.GenerateAudio(ctx, text, lang, quality, speed)
	duration2 := time.Since(startTime2)
	require.NoError(t, err)
	assert.NotEmpty(t, audioURL2)

	// 同じURLが返される
	assert.Equal(t, audioURL1, audioURL2)

	// 2回目の方が速い（キャッシュヒット）
	assert.Less(t, duration2, duration1)
}

// TestBatchGenerate はバッチ生成のテスト
func TestBatchGenerate(t *testing.T) {
	ctx := context.Background()
	service := NewTTSService()

	texts := []string{
		"Hello",
		"Goodbye",
		"Thank you",
	}
	lang := "en"
	quality := "standard"
	speed := 1.0

	audioURLs, err := service.BatchGenerate(ctx, texts, lang, quality, speed)

	require.NoError(t, err)
	assert.Len(t, audioURLs, len(texts))

	for _, url := range audioURLs {
		assert.NotEmpty(t, url)
		assert.Contains(t, url, "http")
	}
}

// TestRateLimit はレート制限のテスト
func TestRateLimit(t *testing.T) {
	ctx := context.Background()
	service := NewTTSService()

	// 連続して大量のリクエストを送信
	for i := 0; i < 10; i++ {
		_, err := service.GenerateAudio(ctx, "Test", "en", "standard", 1.0)
		if err != nil {
			// レート制限エラーが発生する可能性がある
			assert.Contains(t, err.Error(), "rate limit")
			return
		}
	}
}

// TestDifferentQualityLevels は異なる品質レベルのテスト
func TestDifferentQualityLevels(t *testing.T) {
	ctx := context.Background()
	service := NewTTSService()

	text := "Quality test"
	lang := "en"
	speed := 1.0

	qualities := []string{"standard", "premium"}

	for _, quality := range qualities {
		audioURL, err := service.GenerateAudio(ctx, text, lang, quality, speed)
		require.NoError(t, err, "Quality %s should work", quality)
		assert.NotEmpty(t, audioURL)
	}
}

// TestTTSLatency はTTSレイテンシのテスト（1秒以内）
func TestTTSLatency(t *testing.T) {
	ctx := context.Background()
	service := NewTTSService()

	text := "Latency test"
	lang := "en"
	quality := "standard"
	speed := 1.0

	startTime := time.Now()
	_, err := service.GenerateAudio(ctx, text, lang, quality, speed)
	duration := time.Since(startTime)

	require.NoError(t, err)
	// モック環境では1秒以内に完了するはず
	assert.Less(t, duration, 1*time.Second)
}

// TestConcurrentGeneration は並列生成のテスト
func TestConcurrentGeneration(t *testing.T) {
	ctx := context.Background()
	service := NewTTSService()

	texts := []string{"Test 1", "Test 2", "Test 3", "Test 4", "Test 5"}
	lang := "en"
	quality := "standard"
	speed := 1.0

	results := make(chan string, len(texts))
	errors := make(chan error, len(texts))

	for _, text := range texts {
		go func(txt string) {
			audioURL, err := service.GenerateAudio(ctx, txt, lang, quality, speed)
			if err != nil {
				errors <- err
			} else {
				results <- audioURL
			}
		}(text)
	}

	// 結果を収集
	for i := 0; i < len(texts); i++ {
		select {
		case url := <-results:
			assert.NotEmpty(t, url)
		case err := <-errors:
			t.Errorf("Concurrent generation failed: %v", err)
		case <-time.After(5 * time.Second):
			t.Fatal("Timeout waiting for concurrent generation")
		}
	}
}

// TestErrorHandling はエラーハンドリングのテスト
func TestErrorHandling(t *testing.T) {
	ctx := context.Background()
	service := NewTTSService()

	testCases := []struct {
		name    string
		text    string
		lang    string
		quality string
		speed   float64
		wantErr bool
	}{
		{
			name:    "Empty text",
			text:    "",
			lang:    "en",
			quality: "standard",
			speed:   1.0,
			wantErr: true,
		},
		{
			name:    "Invalid speed (too low)",
			text:    "Test",
			lang:    "en",
			quality: "standard",
			speed:   0.1,
			wantErr: true,
		},
		{
			name:    "Invalid speed (too high)",
			text:    "Test",
			lang:    "en",
			quality: "standard",
			speed:   5.0,
			wantErr: true,
		},
		{
			name:    "Valid input",
			text:    "Test",
			lang:    "en",
			quality: "standard",
			speed:   1.0,
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := service.GenerateAudio(ctx, tc.text, tc.lang, tc.quality, tc.speed)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
