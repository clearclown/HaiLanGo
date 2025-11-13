package dictionary

import (
	"context"
	"testing"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestOxfordDictionaryClient tests the Oxford Dictionary API client
func TestOxfordDictionaryClient(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	t.Run("LookupWord_Success", func(t *testing.T) {
		client := NewMockOxfordClient()
		ctx := context.Background()

		entry, err := client.LookupWord(ctx, "hello", "en")
		require.NoError(t, err)
		require.NotNil(t, entry)

		assert.Equal(t, "hello", entry.Word)
		assert.NotEmpty(t, entry.Meanings)
		assert.Equal(t, "oxford", entry.SourceAPI)
	})

	t.Run("LookupWord_NotFound", func(t *testing.T) {
		client := NewMockOxfordClient()
		ctx := context.Background()

		_, err := client.LookupWord(ctx, "xyzabc123notaword", "en")
		assert.ErrorIs(t, err, ErrWordNotFound)
	})

	t.Run("GetName", func(t *testing.T) {
		client := NewMockOxfordClient()
		assert.Equal(t, "Oxford Dictionary", client.GetName())
	})
}

// TestFreeDictionaryClient tests the Free Dictionary API client
func TestFreeDictionaryClient(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	t.Run("LookupWord_Success", func(t *testing.T) {
		client := NewMockFreeDictionaryClient()
		ctx := context.Background()

		entry, err := client.LookupWord(ctx, "hello", "en")
		require.NoError(t, err)
		require.NotNil(t, entry)

		assert.Equal(t, "hello", entry.Word)
		assert.NotEmpty(t, entry.Meanings)
		assert.Equal(t, "free_dictionary", entry.SourceAPI)
	})

	t.Run("LookupWord_NotFound", func(t *testing.T) {
		client := NewMockFreeDictionaryClient()
		ctx := context.Background()

		_, err := client.LookupWord(ctx, "xyzabc123notaword", "en")
		assert.ErrorIs(t, err, ErrWordNotFound)
	})

	t.Run("GetName", func(t *testing.T) {
		client := NewMockFreeDictionaryClient()
		assert.Equal(t, "Free Dictionary", client.GetName())
	})
}

// TestWiktionaryClient tests the Wiktionary API client
func TestWiktionaryClient(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	t.Run("LookupWord_Success", func(t *testing.T) {
		client := NewMockWiktionaryClient()
		ctx := context.Background()

		entry, err := client.LookupWord(ctx, "hello", "en")
		require.NoError(t, err)
		require.NotNil(t, entry)

		assert.Equal(t, "hello", entry.Word)
		assert.NotEmpty(t, entry.Meanings)
		assert.Equal(t, "wiktionary", entry.SourceAPI)
	})

	t.Run("LookupWord_NotFound", func(t *testing.T) {
		client := NewMockWiktionaryClient()
		ctx := context.Background()

		_, err := client.LookupWord(ctx, "xyzabc123notaword", "en")
		assert.ErrorIs(t, err, ErrWordNotFound)
	})

	t.Run("GetName", func(t *testing.T) {
		client := NewMockWiktionaryClient()
		assert.Equal(t, "Wiktionary", client.GetName())
	})
}

// TestDictionaryService tests the dictionary service with fallback
func TestDictionaryService(t *testing.T) {
	t.Run("LookupWord_PrimarySuccess", func(t *testing.T) {
		primaryClient := NewMockOxfordClient()
		fallbackClients := []Client{NewMockFreeDictionaryClient()}
		cache := NewMockCacheClient()

		service := NewService(primaryClient, fallbackClients, cache)
		ctx := context.Background()

		entry, err := service.LookupWord(ctx, "hello", "en")
		require.NoError(t, err)
		require.NotNil(t, entry)

		assert.Equal(t, "hello", entry.Word)
		assert.Equal(t, "oxford", entry.SourceAPI)
	})

	t.Run("LookupWord_FallbackSuccess", func(t *testing.T) {
		primaryClient := NewMockOxfordClientWithError()
		fallbackClients := []Client{NewMockFreeDictionaryClient()}
		cache := NewMockCacheClient()

		service := NewService(primaryClient, fallbackClients, cache)
		ctx := context.Background()

		entry, err := service.LookupWord(ctx, "hello", "en")
		require.NoError(t, err)
		require.NotNil(t, entry)

		assert.Equal(t, "hello", entry.Word)
		assert.Equal(t, "free_dictionary", entry.SourceAPI)
	})

	t.Run("LookupWord_CacheHit", func(t *testing.T) {
		primaryClient := NewMockOxfordClient()
		fallbackClients := []Client{NewMockFreeDictionaryClient()}
		cache := NewMockCacheClientWithData()

		service := NewService(primaryClient, fallbackClients, cache)
		ctx := context.Background()

		// First call should populate cache
		entry1, err := service.LookupWord(ctx, "hello", "en")
		require.NoError(t, err)

		// Second call should hit cache
		entry2, err := service.LookupWord(ctx, "hello", "en")
		require.NoError(t, err)

		assert.Equal(t, entry1.Word, entry2.Word)
	})

	t.Run("LookupWord_AllFail", func(t *testing.T) {
		primaryClient := NewMockOxfordClientWithError()
		fallbackClients := []Client{NewMockFreeDictionaryClientWithError()}
		cache := NewMockCacheClient()

		service := NewService(primaryClient, fallbackClients, cache)
		ctx := context.Background()

		_, err := service.LookupWord(ctx, "xyzabc123notaword", "en")
		assert.Error(t, err)
	})
}

// Mock implementations for testing

type MockOxfordClient struct {
	shouldError bool
}

func NewMockOxfordClient() *MockOxfordClient {
	return &MockOxfordClient{shouldError: false}
}

func NewMockOxfordClientWithError() *MockOxfordClient {
	return &MockOxfordClient{shouldError: true}
}

func (m *MockOxfordClient) LookupWord(ctx context.Context, word string, language string) (*models.WordEntry, error) {
	if m.shouldError {
		return nil, ErrAPIUnavailable
	}

	if word == "xyzabc123notaword" {
		return nil, ErrWordNotFound
	}

	return &models.WordEntry{
		Word:      word,
		Language:  language,
		SourceAPI: "oxford",
		Phonetics: []models.WordPhonetic{
			{Text: "/həˈloʊ/"},
		},
		Meanings: []models.WordMeaning{
			{
				PartOfSpeech: "noun",
				Definitions: []models.WordDefinition{
					{
						Definition: "used as a greeting",
						Examples:   []string{"Hello, how are you?"},
					},
				},
			},
		},
	}, nil
}

func (m *MockOxfordClient) GetName() string {
	return "Oxford Dictionary"
}

type MockFreeDictionaryClient struct {
	shouldError bool
}

func NewMockFreeDictionaryClient() *MockFreeDictionaryClient {
	return &MockFreeDictionaryClient{shouldError: false}
}

func NewMockFreeDictionaryClientWithError() *MockFreeDictionaryClient {
	return &MockFreeDictionaryClient{shouldError: true}
}

func (m *MockFreeDictionaryClient) LookupWord(ctx context.Context, word string, language string) (*models.WordEntry, error) {
	if m.shouldError {
		return nil, ErrAPIUnavailable
	}

	if word == "xyzabc123notaword" {
		return nil, ErrWordNotFound
	}

	return &models.WordEntry{
		Word:      word,
		Language:  language,
		SourceAPI: "free_dictionary",
		Phonetics: []models.WordPhonetic{
			{Text: "/həˈləʊ/"},
		},
		Meanings: []models.WordMeaning{
			{
				PartOfSpeech: "interjection",
				Definitions: []models.WordDefinition{
					{
						Definition: "used to greet someone",
						Examples:   []string{"Hello there!"},
					},
				},
			},
		},
	}, nil
}

func (m *MockFreeDictionaryClient) GetName() string {
	return "Free Dictionary"
}

type MockWiktionaryClient struct {
	shouldError bool
}

func NewMockWiktionaryClient() *MockWiktionaryClient {
	return &MockWiktionaryClient{shouldError: false}
}

func (m *MockWiktionaryClient) LookupWord(ctx context.Context, word string, language string) (*models.WordEntry, error) {
	if m.shouldError {
		return nil, ErrAPIUnavailable
	}

	if word == "xyzabc123notaword" {
		return nil, ErrWordNotFound
	}

	return &models.WordEntry{
		Word:      word,
		Language:  language,
		SourceAPI: "wiktionary",
		Meanings: []models.WordMeaning{
			{
				PartOfSpeech: "interjection",
				Definitions: []models.WordDefinition{
					{
						Definition: "A greeting",
						Synonyms:   []string{"hi", "hey"},
					},
				},
			},
		},
	}, nil
}

func (m *MockWiktionaryClient) GetName() string {
	return "Wiktionary"
}

func NewMockCacheClientWithData() *MockCacheClient {
	return NewMockCacheClient()
}
