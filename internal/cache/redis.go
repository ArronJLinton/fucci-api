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
		fmt.Printf("Error marshaling data for cache key %s: %v\n", key, err)
		return err
	}

	err = c.client.Set(ctx, key, data, expiration).Err()
	if err != nil {
		fmt.Printf("Error setting cache key %s: %v\n", key, err)
		return err
	}
	fmt.Printf("Successfully set cache key %s with expiration %v\n", key, expiration)
	return nil
}

// Get retrieves a value from the cache by key and unmarshals it into the provided interface
func (c *Cache) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			fmt.Printf("Cache miss for key %s\n", key)
			return nil // Key not found
		}
		fmt.Printf("Error getting cache key %s: %v\n", key, err)
		return err
	}

	err = json.Unmarshal(data, dest)
	if err != nil {
		fmt.Printf("Error unmarshaling data for cache key %s: %v\n", key, err)
		return err
	}
	fmt.Printf("Cache hit for key %s\n", key)
	return nil
}

// Delete removes a key from the cache
func (c *Cache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

// Exists checks if a key exists in the cache
func (c *Cache) Exists(ctx context.Context, key string) (bool, error) {
	result, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		fmt.Printf("Error checking existence of cache key %s: %v\n", key, err)
		return false, err
	}
	exists := result > 0
	fmt.Printf("Cache key %s exists: %v\n", key, exists)
	return exists, nil
}
