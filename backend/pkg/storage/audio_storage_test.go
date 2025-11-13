package storage

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSaveAudio は音声ファイル保存のテスト
func TestSaveAudio(t *testing.T) {
	storage := NewAudioStorage()

	audioData := []byte("fake audio data")
	filename := "test-audio.mp3"

	url, err := storage.Save(audioData, filename)

	require.NoError(t, err)
	assert.NotEmpty(t, url)
	assert.Contains(t, url, filename)
}

// TestGetAudio は音声ファイル取得のテスト
func TestGetAudio(t *testing.T) {
	storage := NewAudioStorage()

	audioData := []byte("fake audio data")
	filename := "test-audio.mp3"

	// 保存
	url, err := storage.Save(audioData, filename)
	require.NoError(t, err)

	// 取得
	retrievedData, err := storage.Get(url)
	require.NoError(t, err)
	assert.Equal(t, audioData, retrievedData)
}

// TestDeleteAudio は音声ファイル削除のテスト
func TestDeleteAudio(t *testing.T) {
	storage := NewAudioStorage()

	audioData := []byte("fake audio data")
	filename := "test-audio.mp3"

	// 保存
	url, err := storage.Save(audioData, filename)
	require.NoError(t, err)

	// 削除
	err = storage.Delete(url)
	require.NoError(t, err)

	// 取得（失敗）
	_, err = storage.Get(url)
	assert.Error(t, err)
}

// TestGenerateFilename はファイル名生成のテスト
func TestGenerateFilename(t *testing.T) {
	storage := NewAudioStorage()

	text := "Hello, world!"
	lang := "en"
	quality := "standard"
	speed := 1.0

	filename := storage.GenerateFilename(text, lang, quality, speed)
	assert.NotEmpty(t, filename)
	assert.Contains(t, filename, ".mp3")

	// 同じパラメータで同じファイル名が生成される
	filename2 := storage.GenerateFilename(text, lang, quality, speed)
	assert.Equal(t, filename, filename2)
}

// TestLargeAudioFile は大きな音声ファイルのテスト
func TestLargeAudioFile(t *testing.T) {
	storage := NewAudioStorage()

	// 1MBのダミーデータ
	largeData := bytes.Repeat([]byte("a"), 1024*1024)
	filename := "large-audio.mp3"

	url, err := storage.Save(largeData, filename)
	require.NoError(t, err)
	assert.NotEmpty(t, url)

	// 取得
	retrievedData, err := storage.Get(url)
	require.NoError(t, err)
	assert.Equal(t, len(largeData), len(retrievedData))
}

// TestConcurrentSave は並列保存のテスト
func TestConcurrentSave(t *testing.T) {
	storage := NewAudioStorage()

	numGoroutines := 10
	results := make(chan string, numGoroutines)
	errors := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(index int) {
			audioData := []byte("fake audio data")
			filename := storage.GenerateFilename("text", "en", "standard", float64(index))
			url, err := storage.Save(audioData, filename)
			if err != nil {
				errors <- err
			} else {
				results <- url
			}
		}(i)
	}

	// 結果を収集
	for i := 0; i < numGoroutines; i++ {
		select {
		case url := <-results:
			assert.NotEmpty(t, url)
		case err := <-errors:
			t.Errorf("Concurrent save failed: %v", err)
		}
	}
}
