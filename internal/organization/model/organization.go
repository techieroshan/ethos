package model

import "time"

// Organization represents an organization/tenant in the system
type Organization struct {
	ID          string
	Name        string
	Domain      string
	OwnerID     string
	Description string
	Status      string // "active", "suspended", "trial"
	Plan        string // "trial", "pro", "enterprise"
	MaxUsers    int
	CreatedAt   time.Time
	UpdatedAt   time.Time
	SuspendedAt *time.Time
}

// OrganizationResponse represents an organization for API responses
type OrganizationResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Domain      string    `json:"domain"`
	OwnerID     string    `json:"owner_id"`
	Description string    `json:"description,omitempty"`
	Status      string    `json:"status"`
	Plan        string    `json:"plan"`
	MaxUsers    int       `json:"max_users"`
	UserCount   int       `json:"user_count"`
	AdminCount  int       `json:"admin_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateOrganizationRequest represents a request to create an organization
type CreateOrganizationRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=255"`
	Domain      string `json:"domain" binding:"required,min=1,max=255"`
	Description string `json:"description" binding:"max=500"`
	Plan        string `json:"plan" binding:"required,oneof=trial pro enterprise"`
}

// UpdateOrganizationRequest represents a request to update an organization
type UpdateOrganizationRequest struct {
	Name        string `json:"name" binding:"omitempty,min=1,max=255"`
	Description string `json:"description" binding:"omitempty,max=500"`
	Status      string `json:"status" binding:"omitempty,oneof=active suspended"`
}

// OrganizationMember represents a member of an organization
type OrganizationMember struct {
	ID             string
	OrganizationID string
	UserID         string
	Role           string // "admin", "moderator", "user"
	JoinedAt       time.Time
	LastActiveAt   *time.Time
}

// OrganizationMemberResponse represents a member for API responses
type OrganizationMemberResponse struct {
	ID           string     `json:"id"`
	UserID       string     `json:"user_id"`
	UserName     string     `json:"user_name"`
	UserEmail    string     `json:"user_email"`
	Role         string     `json:"role"`
	Status       string     `json:"status"` // "active", "suspended", "pending"
	JoinedAt     time.Time  `json:"joined_at"`
	LastActiveAt *time.Time `json:"last_active_at,omitempty"`
}

// AddMemberRequest represents a request to add a member to an organization
type AddMemberRequest struct {
	Email string `json:"email" binding:"required,email"`
	Role  string `json:"role" binding:"required,oneof=admin moderator user"`
}

// UpdateMemberRequest represents a request to update a member's role
type UpdateMemberRequest struct {
	Role string `json:"role" binding:"required,oneof=admin moderator user"`
}

// OrganizationSettings represents organization-wide settings
type OrganizationSettings struct {
	ID                       string
	OrganizationID           string
	RequireEmailVerification bool
	AllowPublicProfiles      bool
	EnableModeration         bool
	RequireApproval          bool
	DataRetentionDays        int
	UpdatedAt                time.Time
}

// UpdateSettingsRequest represents a request to update organization settings
type UpdateSettingsRequest struct {
	RequireEmailVerification *bool `json:"require_email_verification"`
	AllowPublicProfiles      *bool `json:"allow_public_profiles"`
	EnableModeration         *bool `json:"enable_moderation"`
	RequireApproval          *bool `json:"require_approval"`
	DataRetentionDays        *int  `json:"data_retention_days"`
}
