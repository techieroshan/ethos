package service

import (
	"context"
	"time"

	"ethos/internal/auth/model"
	prefModel "ethos/internal/profile/model"
)

// UpdateProfileRequest represents a profile update request
type UpdateProfileRequest struct {
	Name      string `json:"name"`
	PublicBio string `json:"public_bio"`
}

// UpdatePreferencesRequest represents a preferences update request
type UpdatePreferencesRequest struct {
	NotifyOnLogin *bool   `json:"notify_on_login"`
	Locale        *string `json:"locale"`
}

// OptOutRequest represents an opt-out request
type OptOutRequest struct {
	From   string `json:"from" binding:"required,oneof=public_search analytics_use"`
	Reason string `json:"reason,omitempty"`
}

// AnonymizeResponse represents the response from anonymization
type AnonymizeResponse struct {
	Status             string    `json:"status"`
	ExpectedCompletion time.Time `json:"expected_completion"`
}

// DeleteRequest represents a deletion request
type DeleteRequest struct {
	Confirm bool   `json:"confirm" binding:"required"`
	Reason  string `json:"reason,omitempty"`
}

// DeleteResponse represents the response from deletion request
type DeleteResponse struct {
	Status             string    `json:"status"`
	ExpectedCompletion time.Time `json:"expected_completion"`
}

// Service defines the interface for profile business logic
type Service interface {
	// GetProfile retrieves a user profile by ID (for authenticated user)
	GetProfile(ctx context.Context, userID string) (*model.UserProfile, error)

	// GetUserProfile retrieves a user profile by ID (for any user)
	GetUserProfile(ctx context.Context, userID string) (*model.UserProfile, error)

	// UpdateProfile updates a user profile
	UpdateProfile(ctx context.Context, userID string, req *UpdateProfileRequest) (*model.UserProfile, error)

	// UpdatePreferences updates user preferences
	UpdatePreferences(ctx context.Context, userID string, req *UpdatePreferencesRequest) (*prefModel.UserPreferences, error)

	// DeleteProfile schedules account deletion
	DeleteProfile(ctx context.Context, userID string) error

	// SearchProfiles searches for user profiles
	SearchProfiles(ctx context.Context, query string, limit, offset int) ([]*model.UserProfile, int, error)

	// OptOut handles opt-out requests from certain features
	OptOut(ctx context.Context, userID string, req *OptOutRequest) error

	// Anonymize anonymizes user personal data
	Anonymize(ctx context.Context, userID string) (*AnonymizeResponse, error)

	// RequestDeletion requests account deletion
	RequestDeletion(ctx context.Context, userID string, req *DeleteRequest) (*DeleteResponse, error)
}

