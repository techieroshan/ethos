package service

import (
	"context"

	"ethos/internal/auth/model"
	"ethos/internal/profile"
	prefModel "ethos/internal/profile/model"
)

// Service defines the interface for profile business logic
type Service interface {
	// GetProfile retrieves a user profile by ID (for authenticated user)
	GetProfile(ctx context.Context, userID string) (*model.UserProfile, error)

	// GetUserProfile retrieves a user profile by ID (for any user)
	GetUserProfile(ctx context.Context, userID string) (*model.UserProfile, error)

	// UpdateProfile updates a user profile
	UpdateProfile(ctx context.Context, userID string, req *profile.UpdateProfileRequest) (*model.UserProfile, error)

	// UpdatePreferences updates user preferences
	UpdatePreferences(ctx context.Context, userID string, req *profile.UpdatePreferencesRequest) (*prefModel.UserPreferences, error)

	// DeleteProfile schedules account deletion
	DeleteProfile(ctx context.Context, userID string) error

	// SearchProfiles searches for user profiles
	SearchProfiles(ctx context.Context, query string, limit, offset int) ([]*model.UserProfile, int, error)

	// OptOut handles opt-out requests from certain features
	OptOut(ctx context.Context, userID string, req *profile.OptOutRequest) error

	// Anonymize anonymizes user personal data
	Anonymize(ctx context.Context, userID string) (*profile.AnonymizeResponse, error)

	// RequestDeletion requests account deletion
	RequestDeletion(ctx context.Context, userID string, req *profile.DeleteRequest) (*profile.DeleteResponse, error)
}

