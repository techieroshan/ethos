package service

import (
	"context"

	"ethos/internal/organization/model"
)

// Service defines the interface for organization business logic
type Service interface {
	// GetOrganization retrieves an organization by ID
	GetOrganization(ctx context.Context, orgID string) (*model.OrganizationResponse, error)

	// GetOrganizationByDomain retrieves an organization by domain
	GetOrganizationByDomain(ctx context.Context, domain string) (*model.OrganizationResponse, error)

	// ListOrganizations retrieves all organizations
	ListOrganizations(ctx context.Context, limit, offset int) ([]*model.OrganizationResponse, error)

	// CreateOrganization creates a new organization
	CreateOrganization(ctx context.Context, ownerID string, req *model.CreateOrganizationRequest) (*model.OrganizationResponse, error)

	// UpdateOrganization updates an organization
	UpdateOrganization(ctx context.Context, orgID string, req *model.UpdateOrganizationRequest) (*model.OrganizationResponse, error)

	// DeleteOrganization deletes an organization
	DeleteOrganization(ctx context.Context, orgID string) error

	// ListOrganizationMembers retrieves members of an organization
	ListOrganizationMembers(ctx context.Context, orgID string, limit, offset int) ([]*model.OrganizationMemberResponse, error)

	// AddOrganizationMember adds a user to an organization
	AddOrganizationMember(ctx context.Context, orgID string, req *model.AddMemberRequest) (*model.OrganizationMemberResponse, error)

	// UpdateOrganizationMemberRole updates a member's role
	UpdateOrganizationMemberRole(ctx context.Context, orgID, userID string, req *model.UpdateMemberRequest) (*model.OrganizationMemberResponse, error)

	// RemoveOrganizationMember removes a user from an organization
	RemoveOrganizationMember(ctx context.Context, orgID, userID string) error

	// GetOrganizationSettings retrieves organization settings
	GetOrganizationSettings(ctx context.Context, orgID string) (*model.OrganizationSettings, error)

	// UpdateOrganizationSettings updates organization settings
	UpdateOrganizationSettings(ctx context.Context, orgID string, req *model.UpdateSettingsRequest) (*model.OrganizationSettings, error)

	// ADMIN METHODS - Platform-wide operations

	// ListAllUsers lists all users across all organizations (admin only)
	ListAllUsers(ctx context.Context, limit, offset int, search, status string) ([]*model.UserAdminResponse, int, error)

	// GetUserDetails gets detailed user information (admin only)
	GetUserDetails(ctx context.Context, userID string) (*model.UserAdminResponse, error)

	// SuspendUser suspends a user account (admin only)
	SuspendUser(ctx context.Context, userID, reason string, duration *int, adminID string) error

	// BanUser permanently bans a user (admin only)
	BanUser(ctx context.Context, userID, reason, adminID string) error

	// UnbanUser removes a ban from a user (admin only)
	UnbanUser(ctx context.Context, userID, adminID string) error

	// DeleteUser permanently deletes a user account (admin only)
	DeleteUser(ctx context.Context, userID, adminID string) error

	// GetSystemAnalytics gets system-wide analytics (admin only)
	GetSystemAnalytics(ctx context.Context) (*model.SystemAnalytics, error)

	// GetUserAnalytics gets user-related analytics (admin only)
	GetUserAnalytics(ctx context.Context) (*model.UserAnalytics, error)

	// GetContentAnalytics gets content-related analytics (admin only)
	GetContentAnalytics(ctx context.Context) (*model.ContentAnalytics, error)

	// GetAuditLogs gets audit logs (admin only)
	GetAuditLogs(ctx context.Context, limit, offset int, userID, action, startDate, endDate string) ([]*model.AuditLogEntry, int, error)

	// GetAuditEntry gets a specific audit log entry (admin only)
	GetAuditEntry(ctx context.Context, entryID string) (*model.AuditLogEntry, error)

	// GetSystemSettings gets system-wide settings (admin only)
	GetSystemSettings(ctx context.Context) (*model.SystemSettings, error)

	// UpdateSystemSettings updates system-wide settings (admin only)
	UpdateSystemSettings(ctx context.Context, settings map[string]interface{}, adminID string) (*model.SystemSettings, error)

	// BulkSuspendUsers suspends multiple users at once (admin only)
	BulkSuspendUsers(ctx context.Context, userIDs []string, reason string, duration *int, adminID string) (*model.BulkOperationResult, error)
}
