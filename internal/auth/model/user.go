package model

import (
	"time"

	"ethos/pkg/errors"
)

// User represents a user in the system
type User struct {
	ID            string
	Email         string
	PasswordHash  string
	Name          string
	EmailVerified bool
	PublicBio     string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Roles         []UserRole `json:"-"` // Loaded separately
	// Multi-tenant context
	CurrentTenantID *string    `json:"current_tenant_id,omitempty"`
	TenantMemberships []TenantMembership `json:"-"` // Loaded separately
}

// UserRole represents a role assigned to a user
type UserRole struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Permissions map[string]interface{} `json:"permissions"`
	AssignedAt  time.Time `json:"assigned_at"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	IsActive    bool      `json:"is_active"`
}

// HasRole checks if user has a specific role
func (u *User) HasRole(roleName string) bool {
	for _, role := range u.Roles {
		if role.Name == roleName && role.IsActive {
			return true
		}
	}
	return false
}

// IsAdmin checks if user has admin privileges (platform_admin or org_admin)
func (u *User) IsAdmin() bool {
	return u.HasRole("platform_admin") || u.HasRole("org_admin")
}

// IsModerator checks if user has moderation privileges
func (u *User) IsModerator() bool {
	return u.HasRole("moderator") || u.IsAdmin()
}

// TenantMembership represents a user's membership in a tenant/organization
type TenantMembership struct {
	TenantID   string    `json:"tenant_id"`
	TenantName string    `json:"tenant_name"`
	Role       string    `json:"role"`
	JoinedAt   time.Time `json:"joined_at"`
	IsActive   bool      `json:"is_active"`
}

// GetCurrentTenantMembership returns the user's current tenant membership
func (u *User) GetCurrentTenantMembership() *TenantMembership {
	if u.CurrentTenantID == nil {
		return nil
	}

	for _, membership := range u.TenantMemberships {
		if membership.TenantID == *u.CurrentTenantID && membership.IsActive {
			return &membership
		}
	}
	return nil
}

// HasTenantAccess checks if user has access to a specific tenant
func (u *User) HasTenantAccess(tenantID string) bool {
	for _, membership := range u.TenantMemberships {
		if membership.TenantID == tenantID && membership.IsActive {
			return true
		}
	}
	return false
}

// HasTenantRole checks if user has a specific role in a tenant
func (u *User) HasTenantRole(tenantID, role string) bool {
	for _, membership := range u.TenantMemberships {
		if membership.TenantID == tenantID && membership.Role == role && membership.IsActive {
			return true
		}
	}
	return false
}

// IsTenantAdmin checks if user is an admin in the current tenant
func (u *User) IsTenantAdmin() bool {
	if u.CurrentTenantID == nil {
		return false
	}
	return u.HasTenantRole(*u.CurrentTenantID, "admin") || u.HasTenantRole(*u.CurrentTenantID, "owner") || u.IsAdmin()
}

// SwitchTenant switches the user's current tenant context
func (u *User) SwitchTenant(tenantID string) error {
	if !u.HasTenantAccess(tenantID) {
		return errors.NewValidationError("user does not have access to this tenant")
	}

	u.CurrentTenantID = &tenantID
	return nil
}

// UserSummary represents a minimal user reference
type UserSummary struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// UserProfile represents a user profile for API responses
type UserProfile struct {
	ID            string    `json:"id"`
	Email         string    `json:"email"`
	Name          string    `json:"name"`
	EmailVerified bool      `json:"email_verified"`
	PublicBio     string    `json:"public_bio,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty"`
}

// ToSummary converts a User to UserSummary
func (u *User) ToSummary() *UserSummary {
	return &UserSummary{
		ID:   u.ID,
		Name: u.Name,
	}
}

// ToProfile converts a User to UserProfile
func (u *User) ToProfile() *UserProfile {
	var updatedAt *time.Time
	if !u.UpdatedAt.IsZero() {
		updatedAt = &u.UpdatedAt
	}

	return &UserProfile{
		ID:            u.ID,
		Email:         u.Email,
		Name:          u.Name,
		EmailVerified: u.EmailVerified,
		PublicBio:     u.PublicBio,
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     updatedAt,
	}
}

