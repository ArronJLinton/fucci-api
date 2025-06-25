package cache

import (
	"context"
	"time"
)

// CacheInterface defines the interface for cache operations
type CacheInterface interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string, dest interface{}) error
	Exists(ctx context.Context, key string) (bool, error)
	Delete(ctx context.Context, key string) error
	DeletePattern(ctx context.Context, pattern string) error
	FlushAll(ctx context.Context) error
	HealthCheck(ctx context.Context) error
	GetStats(ctx context.Context) (map[string]interface{}, error)
}
