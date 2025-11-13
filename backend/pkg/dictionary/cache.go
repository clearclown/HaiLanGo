package dictionary

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/go-redis/redis/v8"
)

// RedisCacheClient is a Redis-based cache client for dictionary entries
type RedisCacheClient struct {
	client *redis.Client
	ttl    time.Duration
}

// NewRedisCacheClient creates a new Redis cache client
func NewRedisCacheClient(addr string, password string, db int) (*RedisCacheClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// Test connection
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisCacheClient{
		client: client,
		ttl:    30 * 24 * time.Hour, // 30 days as per requirements
	}, nil
}

// Get retrieves a cached word entry
func (c *RedisCacheClient) Get(ctx context.Context, key string) (*models.WordEntry, error) {
	data, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("cache miss")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get from cache: %w", err)
	}

	var entry models.WordEntry
	if err := json.Unmarshal([]byte(data), &entry); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cached entry: %w", err)
	}

	return &entry, nil
}

// Set stores a word entry in cache
func (c *RedisCacheClient) Set(ctx context.Context, key string, entry *models.WordEntry) error {
	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal entry: %w", err)
	}

	if err := c.client.Set(ctx, key, data, c.ttl).Err(); err != nil {
		return fmt.Errorf("failed to set cache: %w", err)
	}

	return nil
}

// GenerateKey generates a cache key for a word and language
func (c *RedisCacheClient) GenerateKey(word string, language string) string {
	return fmt.Sprintf("dictionary:%s:%s", language, word)
}
