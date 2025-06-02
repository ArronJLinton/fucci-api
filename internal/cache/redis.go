package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	client *redis.Client
}

func NewCache(redisURL string) (*Cache, error) {
	// Clean up the URL by removing any newlines and extra spaces
	redisURL = strings.TrimSpace(redisURL)
	redisURL = strings.ReplaceAll(redisURL, "\n", "")

	// If URL is empty, return error
	if redisURL == "" {
		return nil, fmt.Errorf("redis URL is empty")
	}

	// If URL is just "redis://", it's incomplete
	if redisURL == "redis://" {
		return nil, fmt.Errorf("incomplete redis URL")
	}

	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %v", err)
	}

	client := redis.NewClient(opt)

	// Test the connection
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	return &Cache{
		client: client,
	}, nil
}

// Set stores a value in the cache with the given key and expiration time
func (c *Cache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, key, data, expiration).Err()
}

// Get retrieves a value from the cache by key and unmarshals it into the provided interface
func (c *Cache) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil // Key not found
		}
		return err
	}

	return json.Unmarshal(data, dest)
}

// Delete removes a key from the cache
func (c *Cache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

// Exists checks if a key exists in the cache
func (c *Cache) Exists(ctx context.Context, key string) (bool, error) {
	result, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}
