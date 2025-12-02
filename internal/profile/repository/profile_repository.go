package repository

import (
	"context"

	"ethos/internal/auth/model"
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
}

