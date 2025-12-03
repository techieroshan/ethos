package model

import "time"

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

