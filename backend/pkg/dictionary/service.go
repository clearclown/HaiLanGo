package dictionary

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
)

// Service is the main dictionary service with fallback support
type Service struct {
	primaryClient   Client
	fallbackClients []Client
	cache           CacheClient
}

// NewService creates a new dictionary service
func NewService(primaryClient Client, fallbackClients []Client, cache CacheClient) *Service {
	return &Service{
		primaryClient:   primaryClient,
		fallbackClients: fallbackClients,
		cache:           cache,
	}
}

// LookupWord looks up a word with fallback and caching
func (s *Service) LookupWord(ctx context.Context, word string, language string) (*models.WordEntry, error) {
	// Generate cache key
	cacheKey := s.cache.GenerateKey(word, language)

	// Try cache first
	if cachedEntry, err := s.cache.Get(ctx, cacheKey); err == nil {
		log.Printf("Cache hit for word: %s (language: %s)", word, language)
		return cachedEntry, nil
	}

	// Try primary client
	entry, err := s.primaryClient.LookupWord(ctx, word, language)
	if err == nil {
		// Cache the result
		if cacheErr := s.cache.Set(ctx, cacheKey, entry); cacheErr != nil {
			log.Printf("Failed to cache entry: %v", cacheErr)
		}
		return entry, nil
	}

	// If word not found, don't try fallback
	if errors.Is(err, ErrWordNotFound) {
		return nil, err
	}

	log.Printf("Primary client (%s) failed: %v. Trying fallback...", s.primaryClient.GetName(), err)

	// Try fallback clients
	var lastErr error
	for _, fallbackClient := range s.fallbackClients {
		entry, err = fallbackClient.LookupWord(ctx, word, language)
		if err == nil {
			// Cache the result
			if cacheErr := s.cache.Set(ctx, cacheKey, entry); cacheErr != nil {
				log.Printf("Failed to cache entry: %v", cacheErr)
			}
			return entry, nil
		}

		// If word not found, don't try more fallbacks
		if errors.Is(err, ErrWordNotFound) {
			return nil, err
		}

		log.Printf("Fallback client (%s) failed: %v", fallbackClient.GetName(), err)
		lastErr = err
	}

	// All clients failed
	if lastErr != nil {
		return nil, fmt.Errorf("all dictionary APIs failed: %w", lastErr)
	}

	return nil, fmt.Errorf("failed to lookup word")
}

// LookupWordDetails provides detailed information about a word
// This is similar to LookupWord but can be extended for more detailed queries
func (s *Service) LookupWordDetails(ctx context.Context, word string, language string) (*models.WordEntry, error) {
	return s.LookupWord(ctx, word, language)
}
