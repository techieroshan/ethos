package service

import (
	"context"
	"time"

	"ethos/internal/organization/repository"
)

// UserContextService defines methods for managing user context and multi-tenant switching
type UserContextService interface {
	// GetAvailableContexts retrieves all organizations a user belongs to
	GetAvailableContexts(ctx context.Context, userID string) ([]*repository.UserContext, error)

	// GetCurrentContext retrieves the user's current organization context
	GetCurrentContext(ctx context.Context, userID string) (*repository.UserContext, error)

	// SwitchContext switches the user's current organization context
	SwitchContext(ctx context.Context, userID, organizationID string, ipAddress, userAgent string) (*repository.UserContext, error)

	// CreateUserSession creates a new user session for multi-tenant context
	CreateUserSession(ctx context.Context, userID, organizationID, tokenHash, refreshTokenHash, ipAddress, userAgent, deviceName string, expiresAt time.Time) (*repository.UserSession, error)

	// GetUserSession retrieves a user session by token
	GetUserSession(ctx context.Context, tokenHash string) (*repository.UserSession, error)

	// RevokeUserSession revokes a user session
	RevokeUserSession(ctx context.Context, sessionID string) error

	// GetContextSwitchHistory retrieves the user's context switch history
	GetContextSwitchHistory(ctx context.Context, userID string, limit, offset int) ([]*repository.ContextSwitchRecord, error)

	// ValidateUserInOrganization checks if user is a member of an organization
	ValidateUserInOrganization(ctx context.Context, userID, organizationID string) (bool, error)

	// GetUserRoleInOrganization gets the user's role in an organization
	GetUserRoleInOrganization(ctx context.Context, userID, organizationID string) (string, error)
}

// SwitchContextRequest represents a request to switch user context
type SwitchContextRequest struct {
	OrganizationID string `json:"organization_id" validate:"required,uuid"`
}

// SwitchContextResponse represents the response when switching context
type SwitchContextResponse struct {
	Context   *repository.UserContext `json:"context"`
	Message   string                  `json:"message"`
	Timestamp time.Time               `json:"timestamp"`
}

// ContextSwitchHistoryResponse represents the response for context switch history
type ContextSwitchHistoryResponse struct {
	Records []*repository.ContextSwitchRecord `json:"records"`
	Total   int64                             `json:"total"`
}

// AvailableContextsResponse represents all available contexts for a user
type AvailableContextsResponse struct {
	Contexts []*repository.UserContext `json:"contexts"`
	Current  *repository.UserContext   `json:"current"`
	Total    int                       `json:"total"`
}

// RevokeSessionRequest represents a request to revoke a session
type RevokeSessionRequest struct {
	SessionID string `json:"session_id" validate:"required,uuid"`
}

// RevokeSessionResponse represents the response when revoking a session
type RevokeSessionResponse struct {
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}
