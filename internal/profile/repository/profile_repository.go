package repository

import (
	"context"

	"ethos/internal/auth/model"
	"ethos/internal/profile"
	prefModel "ethos/internal/profile/model"
)

// Repository defines the interface for profile data access
type Repository interface {
	// GetUserProfile retrieves a user profile by ID
	GetUserProfile(ctx context.Context, userID string) (*model.UserProfile, error)

	// UpdateUserProfile updates a user profile
	UpdateUserProfile(ctx context.Context, userID string, name, publicBio string) (*model.UserProfile, error)

	// UpdateUserPreferences updates user preferences
	UpdateUserPreferences(ctx context.Context, userID string, notifyOnLogin *bool, locale *string) (*prefModel.UserPreferences, error)

	// ScheduleAccountDeletion schedules an account for deletion
	ScheduleAccountDeletion(ctx context.Context, userID string) error

	// SearchUserProfiles searches for user profiles
	SearchUserProfiles(ctx context.Context, query string, limit, offset int) ([]*model.UserProfile, int, error)

	// OptOut handles opt-out requests
	OptOut(ctx context.Context, userID string, req *profile.OptOutRequest) error

	// Anonymize anonymizes user personal data
	Anonymize(ctx context.Context, userID string) (*profile.AnonymizeResponse, error)

	// RequestDeletion requests account deletion
	RequestDeletion(ctx context.Context, userID string, req *profile.DeleteRequest) (*profile.DeleteResponse, error)
}

