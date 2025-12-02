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
}
