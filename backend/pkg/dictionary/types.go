package dictionary

import (
	"context"
	"errors"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
)

var (
	// ErrWordNotFound is returned when a word is not found in the dictionary
	ErrWordNotFound = errors.New("word not found")

	// ErrAPIUnavailable is returned when the API is unavailable
	ErrAPIUnavailable = errors.New("dictionary API unavailable")

	// ErrRateLimitExceeded is returned when rate limit is exceeded
	ErrRateLimitExceeded = errors.New("rate limit exceeded")
)

// Client is the interface for dictionary API clients
type Client interface {
	// LookupWord looks up a word and returns its dictionary entry
	LookupWord(ctx context.Context, word string, language string) (*models.WordEntry, error)

	// GetName returns the name of the dictionary API
	GetName() string
}

// CacheClient is the interface for dictionary cache
type CacheClient interface {
	// Get retrieves a cached word entry
	Get(ctx context.Context, key string) (*models.WordEntry, error)

	// Set stores a word entry in cache
	Set(ctx context.Context, key string, entry *models.WordEntry) error

	// GenerateKey generates a cache key for a word and language
	GenerateKey(word string, language string) string
}
