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

