package repository

import (
	"context"
	"time"
)

// UserContext represents a user's context within an organization
type UserContext struct {
	UserID           string
	OrganizationID   string
	OrganizationName string
	Role             string
	Permissions      []string
	JoinedAt         time.Time
	LastSwitchedAt   time.Time
}

// UserSession represents an authenticated user session
type UserSession struct {
	ID               string
	UserID           string
	OrganizationID   string
	TokenHash        string
	RefreshTokenHash string
	IPAddress        string
	UserAgent        string
	DeviceName       string
	LastActivityAt   time.Time
	ExpiresAt        time.Time
	RevokedAt        *time.Time
	CreatedAt        time.Time
}

// ContextSwitchRecord represents a context switch event
type ContextSwitchRecord struct {
	ID                 string
	UserID             string
	FromOrganizationID *string
	ToOrganizationID   string
	SessionID          string
	IPAddress          string
	Timestamp          time.Time
}

// ContextRepository defines database operations for user context and sessions
type ContextRepository interface {
	// GetUserOrganizations retrieves all organizations a user belongs to
	GetUserOrganizations(ctx context.Context, userID string) ([]*UserContext, error)

	// GetUserCurrentOrganization retrieves the user's current organization
	GetUserCurrentOrganization(ctx context.Context, userID string) (*UserContext, error)

	// UpdateUserCurrentOrganization updates the user's current organization
	UpdateUserCurrentOrganization(ctx context.Context, userID, organizationID string) error

	// CreateUserSession creates a new user session
	CreateUserSession(ctx context.Context, userID, organizationID, tokenHash, refreshTokenHash, ipAddress, userAgent, deviceName string, expiresAt time.Time) (*UserSession, error)

	// GetUserSessionByToken retrieves a user session by token hash
	GetUserSessionByToken(ctx context.Context, tokenHash string) (*UserSession, error)

	// GetUserSessionByID retrieves a user session by ID
	GetUserSessionByID(ctx context.Context, sessionID string) (*UserSession, error)

	// RevokeUserSession revokes (marks as deleted) a user session
	RevokeUserSession(ctx context.Context, sessionID string) error

	// RevokeAllUserSessions revokes all sessions for a user
	RevokeAllUserSessions(ctx context.Context, userID string) error

	// CleanupExpiredSessions deletes expired sessions older than the given time
	CleanupExpiredSessions(ctx context.Context, beforeTime time.Time) (int, error)

	// RecordContextSwitch records a user switching organizations
	RecordContextSwitch(ctx context.Context, userID, fromOrgID, toOrgID, sessionID, ipAddress string) (*ContextSwitchRecord, error)

	// GetContextSwitchHistory retrieves context switch history for a user
	GetContextSwitchHistory(ctx context.Context, userID string, limit, offset int) ([]*ContextSwitchRecord, int64, error)

	// GetUserSessionsByOrganization retrieves all active sessions for a user in an organization
	GetUserSessionsByOrganization(ctx context.Context, userID, organizationID string) ([]*UserSession, error)

	// IsUserInOrganization checks if a user is a member of an organization
	IsUserInOrganization(ctx context.Context, userID, organizationID string) (bool, error)

	// GetUserRoleInOrganization gets the user's role in an organization
	GetUserRoleInOrganization(ctx context.Context, userID, organizationID string) (string, error)

	// LogOrganizationActivity logs an action in the organization activity log
	LogOrganizationActivity(ctx context.Context, organizationID, userID, action, resourceType, resourceID, ipAddress, userAgent string, changes map[string]interface{}) error

	// GetOrganizationActivity retrieves the activity log for an organization
	GetOrganizationActivity(ctx context.Context, organizationID string, limit, offset int) (interface{}, int64, error)
}
