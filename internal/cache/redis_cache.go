package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// Cache defines the interface for caching operations
type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) bool
	Expire(ctx context.Context, key string, expiration time.Duration) error
	Incr(ctx context.Context, key string) (int64, error)
	FlushAll(ctx context.Context) error
	Close() error
}

// RedisCache implements the Cache interface using Redis
type RedisCache struct {
	client *redis.Client
}

// NewRedisCache creates a new Redis cache instance
func NewRedisCache(addr, password string, db int) Cache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &RedisCache{
		client: rdb,
	}
}

// Get retrieves a value from cache
func (c *RedisCache) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

// Set stores a value in cache with expiration
func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	var val string

	switch v := value.(type) {
	case string:
		val = v
	case int, int64, float64, bool:
		val = fmt.Sprintf("%v", v)
	default:
		// Try to marshal to JSON
		jsonBytes, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal value to JSON: %w", err)
		}
		val = string(jsonBytes)
	}

	return c.client.Set(ctx, key, val, expiration).Err()
}

// Delete removes a key from cache
func (c *RedisCache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

// Exists checks if a key exists in cache
func (c *RedisCache) Exists(ctx context.Context, key string) bool {
	count, err := c.client.Exists(ctx, key).Result()
	return err == nil && count > 0
}

// Expire sets expiration on a key
func (c *RedisCache) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return c.client.Expire(ctx, key, expiration).Err()
}

// Incr increments a numeric value
func (c *RedisCache) Incr(ctx context.Context, key string) (int64, error) {
	return c.client.Incr(ctx, key).Result()
}

// FlushAll clears all cache data
func (c *RedisCache) FlushAll(ctx context.Context) error {
	return c.client.FlushAll(ctx).Err()
}

// Close closes the cache connection
func (c *RedisCache) Close() error {
	return c.client.Close()
}

// GetJSON retrieves and unmarshals a JSON value from cache
func (c *RedisCache) GetJSON(ctx context.Context, key string, dest interface{}) error {
	val, err := c.Get(ctx, key)
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(val), dest)
}

// SetJSON marshals and stores a JSON value in cache
func (c *RedisCache) SetJSON(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.Set(ctx, key, value, expiration)
}

// GetOrSet implements cache-aside pattern
func (c *RedisCache) GetOrSet(ctx context.Context, key string, getter func() (interface{}, error), expiration time.Duration) (interface{}, error) {
	// Try to get from cache first
	if val, err := c.Get(ctx, key); err == nil {
		return val, nil
	}

	// Get from source
	value, err := getter()
	if err != nil {
		return nil, err
	}

	// Cache the result
	if err := c.Set(ctx, key, value, expiration); err != nil {
		// Log error but don't fail the operation
		fmt.Printf("Failed to cache value for key %s: %v\n", key, err)
	}

	return value, nil
}
