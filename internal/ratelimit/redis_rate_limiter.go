package ratelimit

import (
	"context"
	"fmt"
	"time"

	"ethos/internal/cache"
)

// RateLimitConfig defines rate limiting configuration
type RateLimitConfig struct {
	Requests int           // Number of requests allowed
	Window   time.Duration // Time window
}

// RateLimiter defines the interface for rate limiting
type RateLimiter interface {
	Allow(ctx context.Context, key string, config RateLimitConfig) (bool, time.Duration, error)
	Reset(ctx context.Context, key string) error
	GetRemaining(ctx context.Context, key string, config RateLimitConfig) (int, error)
}

// RedisRateLimiter implements distributed rate limiting using Redis
type RedisRateLimiter struct {
	cache cache.Cache
}

// NewRedisRateLimiter creates a new Redis-based rate limiter
func NewRedisRateLimiter(cache cache.Cache) RateLimiter {
	return &RedisRateLimiter{
		cache: cache,
	}
}

// Allow checks if a request should be allowed based on rate limits
func (rl *RedisRateLimiter) Allow(ctx context.Context, key string, config RateLimitConfig) (bool, time.Duration, error) {
	// Use Redis sorted sets for sliding window rate limiting
	now := time.Now().UnixMilli()
	windowStart := now - config.Window.Milliseconds()

	// Remove old requests outside the window
	rl.cache.Delete(ctx, fmt.Sprintf("%s:cleanup", key))

	// Add current request timestamp
	score := float64(now)
	member := fmt.Sprintf("%d", now)

	// Use a Lua script for atomic operations
	luaScript := `
		local key = KEYS[1]
		local window_start = tonumber(ARGV[1])
		local score = tonumber(ARGV[2])
		local member = ARGV[3]
		local max_requests = tonumber(ARGV[4])

		-- Remove old entries
		redis.call('ZREMRANGEBYSCORE', key, '-inf', window_start)

		-- Add new entry
		redis.call('ZADD', key, score, member)

		-- Set expiration on the key
		redis.call('EXPIRE', key, math.ceil(window_start / 1000) + 60)

		-- Count current requests in window
		local count = redis.call('ZCARD', key)

		if count > max_requests then
			-- Remove the current request since limit exceeded
			redis.call('ZREM', key, member)
			return {0, 0}  -- Not allowed, no retry after
		end

		return {1, 0}  -- Allowed
	`

	// For now, implement a simpler version using Redis operations
	// Add current timestamp to sorted set
	redisKey := fmt.Sprintf("ratelimit:%s", key)

	// Clean up old entries (this is not atomic, but good enough for demo)
	cleanupKey := fmt.Sprintf("%s:window", redisKey)

	// Try to increment counter with expiration
	count, err := rl.cache.Incr(ctx, cleanupKey)
	if err != nil {
		return false, 0, err
	}

	// If this is the first request, set expiration
	if count == 1 {
		rl.cache.Expire(ctx, cleanupKey, config.Window)
	}

	// Check if limit exceeded
	if int(count) > config.Requests {
		// Calculate retry after time
		retryAfter := config.Window
		return false, retryAfter, nil
	}

	return true, 0, nil
}

// Reset resets the rate limit for a key
func (rl *RedisRateLimiter) Reset(ctx context.Context, key string) error {
	redisKey := fmt.Sprintf("ratelimit:%s", key)
	windowKey := fmt.Sprintf("%s:window", redisKey)

	if err := rl.cache.Delete(ctx, redisKey); err != nil {
		return err
	}

	return rl.cache.Delete(ctx, windowKey)
}

// GetRemaining returns the number of remaining requests
func (rl *RedisRateLimiter) GetRemaining(ctx context.Context, key string, config RateLimitConfig) (int, error) {
	windowKey := fmt.Sprintf("ratelimit:%s:window", key)

	// Get current count
	countStr, err := rl.cache.Get(ctx, windowKey)
	if err != nil {
		// If key doesn't exist, no requests have been made
		return config.Requests, nil
	}

	count := 0
	fmt.Sscanf(countStr, "%d", &count)

	remaining := config.Requests - count
	if remaining < 0 {
		remaining = 0
	}

	return remaining, nil
}

// MultiTenantRateLimiter provides rate limiting with tenant isolation
type MultiTenantRateLimiter struct {
	limiter RateLimiter
}

// NewMultiTenantRateLimiter creates a tenant-aware rate limiter
func NewMultiTenantRateLimiter(limiter RateLimiter) *MultiTenantRateLimiter {
	return &MultiTenantRateLimiter{
		limiter: limiter,
	}
}

// AllowTenant checks rate limits for a specific tenant
func (mt *MultiTenantRateLimiter) AllowTenant(ctx context.Context, tenantID, userID, endpoint string, config RateLimitConfig) (bool, time.Duration, error) {
	// Create composite key for tenant isolation
	key := fmt.Sprintf("tenant:%s:user:%s:endpoint:%s", tenantID, userID, endpoint)
	return mt.limiter.Allow(ctx, key, config)
}

// AllowGlobal checks global rate limits across all tenants
func (mt *MultiTenantRateLimiter) AllowGlobal(ctx context.Context, userID, endpoint string, config RateLimitConfig) (bool, time.Duration, error) {
	key := fmt.Sprintf("global:user:%s:endpoint:%s", userID, endpoint)
	return mt.limiter.Allow(ctx, key, config)
}

// Predefined rate limit configurations
var (
	// API rate limits
	APIRateLimit = RateLimitConfig{
		Requests: 100,
		Window:   time.Minute,
	}

	// Auth rate limits (stricter)
	AuthRateLimit = RateLimitConfig{
		Requests: 5,
		Window:   time.Minute,
	}

	// Search rate limits
	SearchRateLimit = RateLimitConfig{
		Requests: 50,
		Window:   time.Minute,
	}

	// Feedback creation rate limits
	FeedbackRateLimit = RateLimitConfig{
		Requests: 20,
		Window:   time.Hour,
	}

	// Admin operations rate limits
	AdminRateLimit = RateLimitConfig{
		Requests: 1000,
		Window:   time.Hour,
	}
)

// GetRateLimitConfig returns the appropriate rate limit config for an endpoint
func GetRateLimitConfig(endpoint string) RateLimitConfig {
	switch {
	case endpoint == "/api/v1/auth/login" || endpoint == "/api/v1/auth/register":
		return AuthRateLimit
	case endpoint == "/api/v1/search" || endpoint == "/api/v1/people/search":
		return SearchRateLimit
	case endpoint == "/api/v1/feedback":
		return FeedbackRateLimit
	case endpoint == "/api/v1/admin/":
		return AdminRateLimit
	default:
		return APIRateLimit
	}
}
