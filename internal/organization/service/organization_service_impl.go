package service

import (
	"context"
	"fmt"
	"time"

	"ethos/internal/organization/model"
	"ethos/internal/organization/repository"
	"ethos/pkg/errors"

	"github.com/google/uuid"
)

// OrganizationService implements the Service interface
type OrganizationService struct {
	repo repository.Repository
}

// NewOrganizationService creates a new organization service
func NewOrganizationService(repo repository.Repository) Service {
	return &OrganizationService{
		repo: repo,
	}
}

// GetOrganization retrieves an organization by ID
func (s *OrganizationService) GetOrganization(ctx context.Context, orgID string) (*model.OrganizationResponse, error) {
	org, err := s.repo.GetOrganization(ctx, orgID)
	if err != nil {
		return nil, err
	}

	userCount, _ := s.repo.GetOrganizationUserCount(ctx, orgID)
	adminCount, _ := s.repo.GetOrganizationAdminCount(ctx, orgID)

	return &model.OrganizationResponse{
		ID:          org.ID,
		Name:        org.Name,
		Domain:      org.Domain,
		OwnerID:     org.OwnerID,
		Description: org.Description,
		Status:      org.Status,
		Plan:        org.Plan,
		MaxUsers:    org.MaxUsers,
		UserCount:   userCount,
		AdminCount:  adminCount,
		CreatedAt:   org.CreatedAt,
		UpdatedAt:   org.UpdatedAt,
	}, nil
}

// GetOrganizationByDomain retrieves an organization by domain
func (s *OrganizationService) GetOrganizationByDomain(ctx context.Context, domain string) (*model.OrganizationResponse, error) {
	org, err := s.repo.GetOrganizationByDomain(ctx, domain)
	if err != nil {
		return nil, err
	}

	userCount, _ := s.repo.GetOrganizationUserCount(ctx, org.ID)
	adminCount, _ := s.repo.GetOrganizationAdminCount(ctx, org.ID)

	return &model.OrganizationResponse{
		ID:          org.ID,
		Name:        org.Name,
		Domain:      org.Domain,
		OwnerID:     org.OwnerID,
		Description: org.Description,
		Status:      org.Status,
		Plan:        org.Plan,
		MaxUsers:    org.MaxUsers,
		UserCount:   userCount,
		AdminCount:  adminCount,
		CreatedAt:   org.CreatedAt,
		UpdatedAt:   org.UpdatedAt,
	}, nil
}

// ListOrganizations retrieves all organizations
func (s *OrganizationService) ListOrganizations(ctx context.Context, limit, offset int) ([]*model.OrganizationResponse, error) {
	orgs, err := s.repo.ListOrganizations(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	var responses []*model.OrganizationResponse
	for _, org := range orgs {
		userCount, _ := s.repo.GetOrganizationUserCount(ctx, org.ID)
		adminCount, _ := s.repo.GetOrganizationAdminCount(ctx, org.ID)

		responses = append(responses, &model.OrganizationResponse{
			ID:          org.ID,
			Name:        org.Name,
			Domain:      org.Domain,
			OwnerID:     org.OwnerID,
			Description: org.Description,
			Status:      org.Status,
			Plan:        org.Plan,
			MaxUsers:    org.MaxUsers,
			UserCount:   userCount,
			AdminCount:  adminCount,
			CreatedAt:   org.CreatedAt,
			UpdatedAt:   org.UpdatedAt,
		})
	}

	return responses, nil
}

// CreateOrganization creates a new organization
func (s *OrganizationService) CreateOrganization(ctx context.Context, ownerID string, req *model.CreateOrganizationRequest) (*model.OrganizationResponse, error) {
	// Check if domain is already in use
	existing, _ := s.repo.GetOrganizationByDomain(ctx, req.Domain)
	if existing != nil {
		return nil, errors.ErrEmailAlreadyExists // Using as generic "already exists"
	}

	org := &model.Organization{
		ID:          "org-" + uuid.New().String(),
		Name:        req.Name,
		Domain:      req.Domain,
		OwnerID:     ownerID,
		Description: req.Description,
		Status:      "active",
		Plan:        req.Plan,
		MaxUsers:    100, // Default max users
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.repo.CreateOrganization(ctx, org); err != nil {
		return nil, err
	}

	// Add owner as admin
	member := &model.OrganizationMember{
		ID:             "member-" + uuid.New().String(),
		OrganizationID: org.ID,
		UserID:         ownerID,
		Role:           "admin",
		JoinedAt:       time.Now(),
	}
	s.repo.AddOrganizationMember(ctx, member)

	return &model.OrganizationResponse{
		ID:          org.ID,
		Name:        org.Name,
		Domain:      org.Domain,
		OwnerID:     org.OwnerID,
		Description: org.Description,
		Status:      org.Status,
		Plan:        org.Plan,
		MaxUsers:    org.MaxUsers,
		UserCount:   1,
		AdminCount:  1,
		CreatedAt:   org.CreatedAt,
		UpdatedAt:   org.UpdatedAt,
	}, nil
}

// UpdateOrganization updates an organization
func (s *OrganizationService) UpdateOrganization(ctx context.Context, orgID string, req *model.UpdateOrganizationRequest) (*model.OrganizationResponse, error) {
	org, err := s.repo.GetOrganization(ctx, orgID)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		org.Name = req.Name
	}
	if req.Description != "" {
		org.Description = req.Description
	}
	if req.Status != "" {
		org.Status = req.Status
		if req.Status == "suspended" {
			now := time.Now()
			org.SuspendedAt = &now
		}
	}

	org.UpdatedAt = time.Now()
	if err := s.repo.UpdateOrganization(ctx, org); err != nil {
		return nil, err
	}

	userCount, _ := s.repo.GetOrganizationUserCount(ctx, org.ID)
	adminCount, _ := s.repo.GetOrganizationAdminCount(ctx, org.ID)

	return &model.OrganizationResponse{
		ID:          org.ID,
		Name:        org.Name,
		Domain:      org.Domain,
		OwnerID:     org.OwnerID,
		Description: org.Description,
		Status:      org.Status,
		Plan:        org.Plan,
		MaxUsers:    org.MaxUsers,
		UserCount:   userCount,
		AdminCount:  adminCount,
		CreatedAt:   org.CreatedAt,
		UpdatedAt:   org.UpdatedAt,
	}, nil
}

// DeleteOrganization deletes an organization
func (s *OrganizationService) DeleteOrganization(ctx context.Context, orgID string) error {
	return s.repo.DeleteOrganization(ctx, orgID)
}

// ListOrganizationMembers retrieves members of an organization
func (s *OrganizationService) ListOrganizationMembers(ctx context.Context, orgID string, limit, offset int) ([]*model.OrganizationMemberResponse, error) {
	members, err := s.repo.GetOrganizationMembers(ctx, orgID, limit, offset)
	if err != nil {
		return nil, err
	}

	var responses []*model.OrganizationMemberResponse
	for _, member := range members {
		responses = append(responses, &model.OrganizationMemberResponse{
			ID:           member.ID,
			UserID:       member.UserID,
			UserName:     fmt.Sprintf("User %s", member.UserID[:8]),
			UserEmail:    fmt.Sprintf("user+%s@example.com", member.UserID[:8]),
			Role:         member.Role,
			Status:       "active",
			JoinedAt:     member.JoinedAt,
			LastActiveAt: member.LastActiveAt,
		})
	}

	return responses, nil
}

// AddOrganizationMember adds a user to an organization
func (s *OrganizationService) AddOrganizationMember(ctx context.Context, orgID string, req *model.AddMemberRequest) (*model.OrganizationMemberResponse, error) {
	member := &model.OrganizationMember{
		ID:             "member-" + uuid.New().String(),
		OrganizationID: orgID,
		UserID:         req.Email, // Simplified - would lookup user by email
		Role:           req.Role,
		JoinedAt:       time.Now(),
	}

	if err := s.repo.AddOrganizationMember(ctx, member); err != nil {
		return nil, err
	}

	return &model.OrganizationMemberResponse{
		ID:           member.ID,
		UserID:       member.UserID,
		UserName:     fmt.Sprintf("User %s", member.UserID[:min(8, len(member.UserID))]),
		UserEmail:    req.Email,
		Role:         member.Role,
		Status:       "active",
		JoinedAt:     member.JoinedAt,
		LastActiveAt: member.LastActiveAt,
	}, nil
}

// UpdateOrganizationMemberRole updates a member's role
func (s *OrganizationService) UpdateOrganizationMemberRole(ctx context.Context, orgID, userID string, req *model.UpdateMemberRequest) (*model.OrganizationMemberResponse, error) {
	member, err := s.repo.GetOrganizationMember(ctx, orgID, userID)
	if err != nil {
		return nil, err
	}

	member.Role = req.Role
	if err := s.repo.UpdateOrganizationMember(ctx, member); err != nil {
		return nil, err
	}

	return &model.OrganizationMemberResponse{
		ID:           member.ID,
		UserID:       member.UserID,
		UserName:     fmt.Sprintf("User %s", member.UserID[:min(8, len(member.UserID))]),
		UserEmail:    fmt.Sprintf("user+%s@example.com", member.UserID[:8]),
		Role:         member.Role,
		Status:       "active",
		JoinedAt:     member.JoinedAt,
		LastActiveAt: member.LastActiveAt,
	}, nil
}

// RemoveOrganizationMember removes a user from an organization
func (s *OrganizationService) RemoveOrganizationMember(ctx context.Context, orgID, userID string) error {
	return s.repo.RemoveOrganizationMember(ctx, orgID, userID)
}

// GetOrganizationSettings retrieves organization settings
func (s *OrganizationService) GetOrganizationSettings(ctx context.Context, orgID string) (*model.OrganizationSettings, error) {
	return s.repo.GetOrganizationSettings(ctx, orgID)
}

// UpdateOrganizationSettings updates organization settings
func (s *OrganizationService) UpdateOrganizationSettings(ctx context.Context, orgID string, req *model.UpdateSettingsRequest) (*model.OrganizationSettings, error) {
	settings, err := s.repo.GetOrganizationSettings(ctx, orgID)
	if err != nil {
		return nil, err
	}

	if req.RequireEmailVerification != nil {
		settings.RequireEmailVerification = *req.RequireEmailVerification
	}
	if req.AllowPublicProfiles != nil {
		settings.AllowPublicProfiles = *req.AllowPublicProfiles
	}
	if req.EnableModeration != nil {
		settings.EnableModeration = *req.EnableModeration
	}
	if req.RequireApproval != nil {
		settings.RequireApproval = *req.RequireApproval
	}
	if req.DataRetentionDays != nil {
		settings.DataRetentionDays = *req.DataRetentionDays
	}

	settings.UpdatedAt = time.Now()
	if err := s.repo.UpdateOrganizationSettings(ctx, settings); err != nil {
		return nil, err
	}

	return settings, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
