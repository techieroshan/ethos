package repository

import (
	"context"

	"ethos/internal/organization/model"
)

// Repository defines the interface for organization data access
type Repository interface {
	// GetOrganization retrieves an organization by ID
	GetOrganization(ctx context.Context, orgID string) (*model.Organization, error)

	// GetOrganizationByDomain retrieves an organization by domain
	GetOrganizationByDomain(ctx context.Context, domain string) (*model.Organization, error)

	// ListOrganizations retrieves all organizations with optional filtering
	ListOrganizations(ctx context.Context, limit, offset int) ([]*model.Organization, error)

	// CreateOrganization creates a new organization
	CreateOrganization(ctx context.Context, org *model.Organization) error

	// UpdateOrganization updates an existing organization
	UpdateOrganization(ctx context.Context, org *model.Organization) error

	// DeleteOrganization deletes an organization (soft delete)
	DeleteOrganization(ctx context.Context, orgID string) error

	// GetOrganizationMembers retrieves members of an organization
	GetOrganizationMembers(ctx context.Context, orgID string, limit, offset int) ([]*model.OrganizationMember, error)

	// GetOrganizationMember retrieves a specific member
	GetOrganizationMember(ctx context.Context, orgID, userID string) (*model.OrganizationMember, error)

	// AddOrganizationMember adds a user to an organization
	AddOrganizationMember(ctx context.Context, member *model.OrganizationMember) error

	// UpdateOrganizationMember updates a member's role
	UpdateOrganizationMember(ctx context.Context, member *model.OrganizationMember) error

	// RemoveOrganizationMember removes a user from an organization
	RemoveOrganizationMember(ctx context.Context, orgID, userID string) error

	// GetOrganizationSettings retrieves organization settings
	GetOrganizationSettings(ctx context.Context, orgID string) (*model.OrganizationSettings, error)

	// UpdateOrganizationSettings updates organization settings
	UpdateOrganizationSettings(ctx context.Context, settings *model.OrganizationSettings) error

	// GetOrganizationUserCount retrieves the number of users in an organization
	GetOrganizationUserCount(ctx context.Context, orgID string) (int, error)

	// GetOrganizationAdminCount retrieves the number of admins in an organization
	GetOrganizationAdminCount(ctx context.Context, orgID string) (int, error)
}
