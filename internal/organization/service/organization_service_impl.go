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

// ADMIN METHODS - Platform-wide operations

// ListAllUsers lists all users across all organizations (admin only)
func (s *OrganizationService) ListAllUsers(ctx context.Context, limit, offset int, search, status string) ([]*model.UserAdminResponse, int, error) {
	// TODO: Implement proper database query with joins for roles
	// For now, return mock data
	users := []*model.UserAdminResponse{
		{
			ID:            "user-1",
			Email:         "admin@ethos.com",
			Name:          "Platform Admin",
			EmailVerified: true,
			Status:        "active",
			Roles: []model.UserRole{
				{
					ID:          "role-1",
					Name:        "platform_admin",
					Description: "Platform Administrator",
					AssignedAt:  s.now(),
					IsActive:    true,
				},
			},
			CreatedAt: s.now(),
		},
		{
			ID:            "user-2",
			Email:         "user@example.com",
			Name:          "Standard User",
			EmailVerified: true,
			Status:        "active",
			Roles: []model.UserRole{
				{
					ID:          "role-2",
					Name:        "user",
					Description: "Standard User",
					AssignedAt:  s.now(),
					IsActive:    true,
				},
			},
			CreatedAt: s.now(),
		},
	}

	return users, len(users), nil
}

// GetUserDetails gets detailed user information (admin only)
func (s *OrganizationService) GetUserDetails(ctx context.Context, userID string) (*model.UserAdminResponse, error) {
	// TODO: Implement proper database query
	// Mock response for now
	return &model.UserAdminResponse{
		ID:            userID,
		Email:         "user@example.com",
		Name:          "Test User",
		EmailVerified: true,
		Status:        "active",
		Roles: []model.UserRole{
			{
				ID:          "role-2",
				Name:        "user",
				Description: "Standard User",
				AssignedAt:  s.now(),
				IsActive:    true,
			},
		},
		CreatedAt: s.now(),
	}, nil
}

// SuspendUser suspends a user account (admin only)
func (s *OrganizationService) SuspendUser(ctx context.Context, userID, reason string, duration *int, adminID string) error {
	// TODO: Implement proper database update with moderation action logging
	return nil
}

// BanUser permanently bans a user (admin only)
func (s *OrganizationService) BanUser(ctx context.Context, userID, reason, adminID string) error {
	// TODO: Implement proper database update with moderation action logging
	return nil
}

// UnbanUser removes a ban from a user (admin only)
func (s *OrganizationService) UnbanUser(ctx context.Context, userID, adminID string) error {
	// TODO: Implement proper database update
	return nil
}

// DeleteUser permanently deletes a user account (admin only)
func (s *OrganizationService) DeleteUser(ctx context.Context, userID, adminID string) error {
	// TODO: Implement proper database deletion with audit logging
	return nil
}

// GetSystemAnalytics gets system-wide analytics (admin only)
func (s *OrganizationService) GetSystemAnalytics(ctx context.Context) (*model.SystemAnalytics, error) {
	// TODO: Implement proper analytics queries
	return &model.SystemAnalytics{
		TotalUsers:            1247,
		ActiveUsers:           892,
		TotalOrganizations:    47,
		ActiveOrganizations:   42,
		TotalFeedback:         3456,
		PendingModeration:     23,
		SystemHealth:          "healthy",
	}, nil
}

// GetUserAnalytics gets user-related analytics (admin only)
func (s *OrganizationService) GetUserAnalytics(ctx context.Context) (*model.UserAnalytics, error) {
	// TODO: Implement proper analytics queries
	return &model.UserAnalytics{
		UserGrowth: []model.TimeSeriesPoint{
			{Date: "2024-01-01", Value: 100},
			{Date: "2024-01-02", Value: 150},
		},
		UserRetention: []model.RetentionPoint{
			{Cohort: "2024-01", Day0: 100, Day7: 75, Day30: 60, Day90: 45},
		},
		UserActivity: []model.ActivityPoint{
			{Date: "2024-01-01", ActiveUsers: 892, NewUsers: 23, ReturningUsers: 869},
		},
		GeographicDistribution: []model.GeoPoint{
			{Country: "US", Users: 450, Percent: 36.1},
			{Country: "UK", Users: 180, Percent: 14.4},
		},
	}, nil
}

// GetContentAnalytics gets content-related analytics (admin only)
func (s *OrganizationService) GetContentAnalytics(ctx context.Context) (*model.ContentAnalytics, error) {
	// TODO: Implement proper analytics queries
	return &model.ContentAnalytics{
		FeedbackGrowth: []model.TimeSeriesPoint{
			{Date: "2024-01-01", Value: 100},
			{Date: "2024-01-02", Value: 200},
		},
		ContentModeration: []model.ModerationStats{
			{Date: "2024-01-01", PendingContent: 5, ApprovedContent: 95, RejectedContent: 3, EscalatedContent: 2},
		},
		PopularCategories: []model.CategoryStats{
			{Category: "Engineering", Count: 1200, Percent: 34.7},
			{Category: "Product", Count: 800, Percent: 23.1},
		},
		EngagementMetrics: model.EngagementStats{
			AverageLikes:     4.2,
			AverageComments:  2.8,
			AverageBookmarks: 1.5,
			EngagementRate:   12.3,
		},
	}, nil
}

// GetAuditLogs gets audit logs (admin only)
func (s *OrganizationService) GetAuditLogs(ctx context.Context, limit, offset int, userID, action, startDate, endDate string) ([]*model.AuditLogEntry, int, error) {
	// TODO: Implement proper audit log queries
	logs := []*model.AuditLogEntry{
		{
			ID:        "audit-1",
			UserID:    "user-1",
			UserName:  "Platform Admin",
			Action:    "user_suspend",
			Resource:  "user",
			ResourceID: "user-2",
			Details:   "Suspended user for policy violation",
			IPAddress: "192.168.1.100",
			UserAgent: "Mozilla/5.0...",
			Timestamp: s.now(),
		},
	}

	return logs, len(logs), nil
}

// GetAuditEntry gets a specific audit log entry (admin only)
func (s *OrganizationService) GetAuditEntry(ctx context.Context, entryID string) (*model.AuditLogEntry, error) {
	// TODO: Implement proper audit log query
	return &model.AuditLogEntry{
		ID:        entryID,
		UserID:    "user-1",
		UserName:  "Platform Admin",
		Action:    "user_suspend",
		Resource:  "user",
		ResourceID: "user-2",
		Details:   "Suspended user for policy violation",
		IPAddress: "192.168.1.100",
		UserAgent: "Mozilla/5.0...",
		Timestamp: s.now(),
	}, nil
}

// GetSystemSettings gets system-wide settings (admin only)
func (s *OrganizationService) GetSystemSettings(ctx context.Context) (*model.SystemSettings, error) {
	// TODO: Implement proper settings storage and retrieval
	return &model.SystemSettings{
		ID:                           "system-settings",
		RequireEmailVerification:    true,
		AllowPublicProfiles:         true,
		EnableGlobalModeration:      true,
		MaxFeedbackPerDay:           50,
		MaxCommentsPerHour:          20,
		DataRetentionDays:           365,
		EnableAnalytics:             true,
		MaintenanceMode:             false,
		CustomSettings:              map[string]interface{}{"feature_flag_new_ui": true},
		UpdatedAt:                   s.now(),
		UpdatedBy:                   "system",
	}, nil
}

// UpdateSystemSettings updates system-wide settings (admin only)
func (s *OrganizationService) UpdateSystemSettings(ctx context.Context, settings map[string]interface{}, adminID string) (*model.SystemSettings, error) {
	// TODO: Implement proper settings update with validation
	return &model.SystemSettings{
		ID:                           "system-settings",
		RequireEmailVerification:    true,
		AllowPublicProfiles:         true,
		EnableGlobalModeration:      true,
		MaxFeedbackPerDay:           50,
		MaxCommentsPerHour:          20,
		DataRetentionDays:           365,
		EnableAnalytics:             true,
		MaintenanceMode:             false,
		CustomSettings:              settings,
		UpdatedAt:                   s.now(),
		UpdatedBy:                   adminID,
	}, nil
}

// BulkSuspendUsers suspends multiple users at once (admin only)
func (s *OrganizationService) BulkSuspendUsers(ctx context.Context, userIDs []string, reason string, duration *int, adminID string) (*model.BulkOperationResult, error) {
	// TODO: Implement proper bulk operation with transaction
	return &model.BulkOperationResult{
		TotalRequested: len(userIDs),
		Successful:     len(userIDs),
		Failed:         0,
		Errors:         []model.BulkOperationError{},
	}, nil
}

// Helper method to get current time
func (s *OrganizationService) now() time.Time {
	return time.Now()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
