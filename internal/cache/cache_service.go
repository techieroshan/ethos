package cache

import (
	"context"
	"crypto/md5"
	"fmt"
	"time"
)

// CacheService provides high-level caching operations with different strategies
type CacheService struct {
	cache Cache
}

// NewCacheService creates a new cache service
func NewCacheService(cache Cache) *CacheService {
	return &CacheService{
		cache: cache,
	}
}

// CacheKey generates a cache key from components
func (s *CacheService) CacheKey(components ...string) string {
	key := ""
	for _, component := range components {
		key += component + ":"
	}
	// Remove trailing colon
	if len(key) > 0 {
		key = key[:len(key)-1]
	}
	return key
}

// HashKey generates an MD5 hash of the key for long keys
func (s *CacheService) HashKey(key string) string {
	hash := md5.Sum([]byte(key))
	return fmt.Sprintf("%x", hash)
}

// UserProfileKey generates cache key for user profile
func (s *CacheService) UserProfileKey(userID string) string {
	return s.CacheKey("user", "profile", userID)
}

// UserTenantsKey generates cache key for user tenants
func (s *CacheService) UserTenantsKey(userID string) string {
	return s.CacheKey("user", "tenants", userID)
}

// FeedbackListKey generates cache key for feedback list with filters
func (s *CacheService) FeedbackListKey(tenantID, userID string, filters map[string]interface{}) string {
	key := s.CacheKey("feedback", "list", tenantID, userID)

	// Include filter parameters in key
	if filters != nil {
		for k, v := range filters {
			key = s.CacheKey(key, fmt.Sprintf("%s=%v", k, v))
		}
	}

	// Hash long keys to avoid Redis key length limits
	if len(key) > 250 {
		key = s.HashKey(key)
	}

	return key
}

// OrganizationStatsKey generates cache key for organization statistics
func (s *CacheService) OrganizationStatsKey(orgID string) string {
	return s.CacheKey("org", "stats", orgID)
}

// SearchResultsKey generates cache key for search results
func (s *CacheService) SearchResultsKey(query, tenantID string, filters map[string]interface{}) string {
	key := s.CacheKey("search", query, tenantID)

	if filters != nil {
		for k, v := range filters {
			key = s.CacheKey(key, fmt.Sprintf("%s=%v", k, v))
		}
	}

	if len(key) > 250 {
		key = s.HashKey(key)
	}

	return key
}

// GetUserProfile retrieves cached user profile
func (s *CacheService) GetUserProfile(ctx context.Context, userID string) (string, error) {
	return s.cache.Get(ctx, s.UserProfileKey(userID))
}

// SetUserProfile caches user profile
func (s *CacheService) SetUserProfile(ctx context.Context, userID string, profile interface{}) error {
	return s.cache.Set(ctx, s.UserProfileKey(userID), profile, 30*time.Minute)
}

// GetUserTenants retrieves cached user tenants
func (s *CacheService) GetUserTenants(ctx context.Context, userID string) (string, error) {
	return s.cache.Get(ctx, s.UserTenantsKey(userID))
}

// SetUserTenants caches user tenants
func (s *CacheService) SetUserTenants(ctx context.Context, userID string, tenants interface{}) error {
	return s.cache.Set(ctx, s.UserTenantsKey(userID), tenants, 1*time.Hour)
}

// GetFeedbackList retrieves cached feedback list
func (s *CacheService) GetFeedbackList(ctx context.Context, tenantID, userID string, filters map[string]interface{}) (string, error) {
	return s.cache.Get(ctx, s.FeedbackListKey(tenantID, userID, filters))
}

// SetFeedbackList caches feedback list
func (s *CacheService) SetFeedbackList(ctx context.Context, tenantID, userID string, filters map[string]interface{}, feedbackList interface{}) error {
	return s.cache.Set(ctx, s.FeedbackListKey(tenantID, userID, filters), feedbackList, 15*time.Minute)
}

// GetOrganizationStats retrieves cached organization statistics
func (s *CacheService) GetOrganizationStats(ctx context.Context, orgID string) (string, error) {
	return s.cache.Get(ctx, s.OrganizationStatsKey(orgID))
}

// SetOrganizationStats caches organization statistics
func (s *CacheService) SetOrganizationStats(ctx context.Context, orgID string, stats interface{}) error {
	return s.cache.Set(ctx, s.OrganizationStatsKey(orgID), stats, 10*time.Minute)
}

// GetSearchResults retrieves cached search results
func (s *CacheService) GetSearchResults(ctx context.Context, query, tenantID string, filters map[string]interface{}) (string, error) {
	return s.cache.Get(ctx, s.SearchResultsKey(query, tenantID, filters))
}

// SetSearchResults caches search results
func (s *CacheService) SetSearchResults(ctx context.Context, query, tenantID string, filters map[string]interface{}, results interface{}) error {
	return s.cache.Set(ctx, s.SearchResultsKey(query, tenantID, filters), results, 5*time.Minute)
}

// InvalidateUserCache invalidates all user-related cache entries
func (s *CacheService) InvalidateUserCache(ctx context.Context, userID string) error {
	keys := []string{
		s.UserProfileKey(userID),
		s.UserTenantsKey(userID),
	}

	for _, key := range keys {
		if err := s.cache.Delete(ctx, key); err != nil {
			// Log error but continue
			fmt.Printf("Failed to delete cache key %s: %v\n", key, err)
		}
	}

	return nil
}

// InvalidateOrganizationCache invalidates organization-related cache entries
func (s *CacheService) InvalidateOrganizationCache(ctx context.Context, orgID string) error {
	// For simplicity, we'll invalidate the stats cache
	// In a production system, you might want to use Redis key patterns or pub/sub
	key := s.OrganizationStatsKey(orgID)
	return s.cache.Delete(ctx, key)
}

// InvalidateFeedbackCache invalidates feedback-related cache entries for an organization
func (s *CacheService) InvalidateFeedbackCache(ctx context.Context, orgID string) error {
	// This is a simplified implementation
	// In production, you might want to use Redis SCAN or maintain a separate index
	return nil
}

// WarmupCache pre-loads frequently accessed data into cache
func (s *CacheService) WarmupCache(ctx context.Context) error {
	// This would be called during application startup
	// Implementation would depend on your specific warmup requirements
	return nil
}

// HealthCheck verifies cache connectivity
func (s *CacheService) HealthCheck(ctx context.Context) error {
	testKey := "health_check"
	testValue := "ok"

	// Try to set and get a test value
	if err := s.cache.Set(ctx, testKey, testValue, 1*time.Minute); err != nil {
		return fmt.Errorf("failed to set test value: %w", err)
	}

	if _, err := s.cache.Get(ctx, testKey); err != nil {
		return fmt.Errorf("failed to get test value: %w", err)
	}

	// Clean up
	s.cache.Delete(ctx, testKey)

	return nil
}
