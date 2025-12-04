package service

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	"ethos/internal/organization/repository"
	"ethos/pkg/errors"
)

// UserContextServiceImpl implements the UserContextService interface
type UserContextServiceImpl struct {
	repo repository.ContextRepository
}

// NewUserContextService creates a new user context service
func NewUserContextService(repo repository.ContextRepository) UserContextService {
	return &UserContextServiceImpl{
		repo: repo,
	}
}

// GetAvailableContexts retrieves all organizations a user belongs to
func (s *UserContextServiceImpl) GetAvailableContexts(ctx context.Context, userID string) ([]*repository.UserContext, error) {
	return s.repo.GetUserOrganizations(ctx, userID)
}

// GetCurrentContext retrieves the user's current organization context
func (s *UserContextServiceImpl) GetCurrentContext(ctx context.Context, userID string) (*repository.UserContext, error) {
	return s.repo.GetUserCurrentOrganization(ctx, userID)
}

// SwitchContext switches the user's current organization context
func (s *UserContextServiceImpl) SwitchContext(ctx context.Context, userID, organizationID string, ipAddress, userAgent string) (*repository.UserContext, error) {
	// Verify user is in the organization
	isMember, err := s.repo.IsUserInOrganization(ctx, userID, organizationID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, errors.ErrForbidden
	}

	// Get current organization for audit trail
	currentCtx, err := s.repo.GetUserCurrentOrganization(ctx, userID)
	var currentOrgID string
	if err == nil && currentCtx != nil {
		currentOrgID = currentCtx.OrganizationID
	}

	// Update the user's current organization
	if err := s.repo.UpdateUserCurrentOrganization(ctx, userID, organizationID); err != nil {
		return nil, err
	}

	// Create a new session for this context switch
	tokenHash := hashToken(generateRandomToken())
	refreshTokenHash := hashToken(generateRandomToken())
	expiresAt := time.Now().Add(24 * time.Hour) // 24-hour session

	_, err = s.repo.CreateUserSession(ctx, userID, organizationID, tokenHash, refreshTokenHash, ipAddress, userAgent, "", expiresAt)
	if err != nil {
		return nil, err
	}

	// Record the context switch for audit trail
	if _, err := s.repo.RecordContextSwitch(ctx, userID, currentOrgID, organizationID, "", ipAddress); err != nil {
		// Log but don't fail if audit logging fails
		fmt.Printf("Failed to record context switch: %v\n", err)
	}

	// Log the activity
	if err := s.repo.LogOrganizationActivity(ctx, organizationID, userID, "context_switch", "user", userID, ipAddress, userAgent, map[string]interface{}{
		"from_organization": currentOrgID,
		"to_organization":   organizationID,
	}); err != nil {
		// Log but don't fail if activity logging fails
		fmt.Printf("Failed to log organization activity: %v\n", err)
	}

	// Return the new context
	return s.repo.GetUserCurrentOrganization(ctx, userID)
}

// CreateUserSession creates a new user session
func (s *UserContextServiceImpl) CreateUserSession(ctx context.Context, userID, organizationID, tokenHash, refreshTokenHash, ipAddress, userAgent, deviceName string, expiresAt time.Time) (*repository.UserSession, error) {
	// Verify user is in the organization
	isMember, err := s.repo.IsUserInOrganization(ctx, userID, organizationID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, errors.ErrForbidden
	}

	return s.repo.CreateUserSession(ctx, userID, organizationID, tokenHash, refreshTokenHash, ipAddress, userAgent, deviceName, expiresAt)
}

// GetUserSession retrieves a user session by token
func (s *UserContextServiceImpl) GetUserSession(ctx context.Context, tokenHash string) (*repository.UserSession, error) {
	return s.repo.GetUserSessionByToken(ctx, tokenHash)
}

// RevokeUserSession revokes a user session
func (s *UserContextServiceImpl) RevokeUserSession(ctx context.Context, sessionID string) error {
	return s.repo.RevokeUserSession(ctx, sessionID)
}

// GetContextSwitchHistory retrieves the user's context switch history
func (s *UserContextServiceImpl) GetContextSwitchHistory(ctx context.Context, userID string, limit, offset int) ([]*repository.ContextSwitchRecord, error) {
	records, _, err := s.repo.GetContextSwitchHistory(ctx, userID, limit, offset)
	return records, err
}

// ValidateUserInOrganization checks if user is a member of an organization
func (s *UserContextServiceImpl) ValidateUserInOrganization(ctx context.Context, userID, organizationID string) (bool, error) {
	return s.repo.IsUserInOrganization(ctx, userID, organizationID)
}

// GetUserRoleInOrganization gets the user's role in an organization
func (s *UserContextServiceImpl) GetUserRoleInOrganization(ctx context.Context, userID, organizationID string) (string, error) {
	return s.repo.GetUserRoleInOrganization(ctx, userID, organizationID)
}

// Helper functions

// generateRandomToken generates a random token (in production, use crypto/rand)
func generateRandomToken() string {
	return fmt.Sprintf("%d-%d", time.Now().UnixNano(), time.Now().UnixNano()%1000000)
}

// hashToken hashes a token using SHA256
func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return fmt.Sprintf("%x", hash)
}
