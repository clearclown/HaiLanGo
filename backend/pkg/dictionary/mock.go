package dictionary

import (
	"context"
	"errors"
	"fmt"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
)

// MockCacheClient is a mock implementation of CacheClient for testing
type MockCacheClient struct {
	data map[string]*models.WordEntry
}

// NewMockCacheClient creates a new mock cache client
func NewMockCacheClient() *MockCacheClient {
	return &MockCacheClient{
		data: make(map[string]*models.WordEntry),
	}
}

// Get retrieves a cached word entry
func (m *MockCacheClient) Get(ctx context.Context, key string) (*models.WordEntry, error) {
	if entry, ok := m.data[key]; ok {
		return entry, nil
	}
	return nil, errors.New("not found in cache")
}

// Set stores a word entry in cache
func (m *MockCacheClient) Set(ctx context.Context, key string, entry *models.WordEntry) error {
	m.data[key] = entry
	return nil
}

// GenerateKey generates a cache key for a word and language
func (m *MockCacheClient) GenerateKey(word string, language string) string {
	return fmt.Sprintf("dictionary:%s:%s", language, word)
}
