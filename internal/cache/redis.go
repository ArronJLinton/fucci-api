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

	fmt.Printf("Attempting to connect to Redis with URL: %s\n", redisURL)

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
		fmt.Printf("Failed to parse Redis URL: %v\n", err)
		return nil, fmt.Errorf("failed to parse Redis URL: %v", err)
	}

	fmt.Printf("Redis connection options parsed successfully. Host: %s, Port: %d\n", opt.Addr, opt.DB)

	client := redis.NewClient(opt)

	// Test the connection with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Try to ping Redis
	if err := client.Ping(ctx).Err(); err != nil {
		fmt.Printf("Failed to connect to Redis: %v\n", err)
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	fmt.Printf("Successfully connected to Redis\n")

	// Try a test key write/read
	testCtx, testCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer testCancel()

	testKey := "test_connection"
	testValue := "test_value"

	// Try to write
	err = client.Set(testCtx, testKey, testValue, 1*time.Minute).Err()
	if err != nil {
		fmt.Printf("Failed to write test key to Redis: %v\n", err)
		return nil, fmt.Errorf("failed to write test key to Redis: %v", err)
	}

	// Try to read
	val, err := client.Get(testCtx, testKey).Result()
	if err != nil {
		fmt.Printf("Failed to read test key from Redis: %v\n", err)
		return nil, fmt.Errorf("failed to read test key from Redis: %v", err)
	}

	if val != testValue {
		fmt.Printf("Test key value mismatch. Expected: %s, Got: %s\n", testValue, val)
		return nil, fmt.Errorf("test key value mismatch")
	}

	fmt.Printf("Redis read/write test successful\n")

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

// HealthCheck performs a health check on the Redis connection
func (c *Cache) HealthCheck(ctx context.Context) error {
	// Try to ping Redis
	err := c.client.Ping(ctx).Err()
	if err != nil {
		return fmt.Errorf("Redis ping failed: %v", err)
	}

	// Try a test key write/read
	testKey := "health_check"
	testValue := "test_value"

	// Try to write
	err = c.client.Set(ctx, testKey, testValue, 1*time.Minute).Err()
	if err != nil {
		return fmt.Errorf("Redis write test failed: %v", err)
	}

	// Try to read
	val, err := c.client.Get(ctx, testKey).Result()
	if err != nil {
		return fmt.Errorf("Redis read test failed: %v", err)
	}

	if val != testValue {
		return fmt.Errorf("Redis value mismatch")
	}

	return nil
}
