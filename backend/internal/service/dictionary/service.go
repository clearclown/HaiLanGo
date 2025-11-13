package dictionary

import (
	"context"
	"fmt"
	"os"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/clearclown/HaiLanGo/backend/pkg/dictionary"
)

// Service is the internal service for dictionary operations
type Service struct {
	dictService *dictionary.Service
}

// NewService creates a new dictionary internal service
func NewService() (*Service, error) {
	// Create dictionary clients
	oxfordAPIKey := os.Getenv("OXFORD_DICTIONARY_API_KEY")
	oxfordAppID := os.Getenv("OXFORD_DICTIONARY_APP_ID")
	oxfordClient := dictionary.NewOxfordClient(oxfordAPIKey, oxfordAppID)

	freeDictClient := dictionary.NewFreeDictionaryClient()
	wiktionaryClient := dictionary.NewWiktionaryClient()

	// Create cache client
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	redisPassword := os.Getenv("REDIS_PASSWORD")

	// Check if we should use mock mode
	useMock := os.Getenv("USE_MOCK_APIS") == "true" || os.Getenv("TEST_USE_MOCKS") == "true"

	var cache dictionary.CacheClient
	if useMock {
		cache = dictionary.NewMockCacheClient()
	} else {
		var err error
		cache, err = dictionary.NewRedisCacheClient(redisAddr, redisPassword, 0)
		if err != nil {
			// If Redis is not available, use mock cache
			cache = dictionary.NewMockCacheClient()
		}
	}

	// Create dictionary service with fallback clients
	dictService := dictionary.NewService(
		oxfordClient,
		[]dictionary.Client{freeDictClient, wiktionaryClient},
		cache,
	)

	return &Service{
		dictService: dictService,
	}, nil
}

// LookupWord looks up a word in the dictionary
func (s *Service) LookupWord(ctx context.Context, word string, language string) (*models.WordEntry, error) {
	if word == "" {
		return nil, fmt.Errorf("word cannot be empty")
	}

	if language == "" {
		language = "en" // Default to English
	}

	return s.dictService.LookupWord(ctx, word, language)
}

// LookupWordDetails provides detailed information about a word
func (s *Service) LookupWordDetails(ctx context.Context, word string, language string) (*models.WordEntry, error) {
	if word == "" {
		return nil, fmt.Errorf("word cannot be empty")
	}

	if language == "" {
		language = "en" // Default to English
	}

	return s.dictService.LookupWordDetails(ctx, word, language)
}
