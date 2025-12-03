package service

import (
	"context"

	"ethos/internal/auth/model"
	"ethos/internal/profile"
	"ethos/internal/profile/repository"
	prefModel "ethos/internal/profile/model"
	"ethos/pkg/errors"
)

// ProfileService implements the Service interface
type ProfileService struct {
	repo repository.Repository
}

// NewProfileService creates a new profile service
func NewProfileService(repo repository.Repository) Service {
	return &ProfileService{
		repo: repo,
	}
}

// GetProfile retrieves a user profile by ID (for authenticated user)
func (s *ProfileService) GetProfile(ctx context.Context, userID string) (*model.UserProfile, error) {
	profile, err := s.repo.GetUserProfile(ctx, userID)
	if err != nil {
		return nil, err
	}

	return profile, nil
}

// GetUserProfile retrieves a user profile by ID (for any user)
func (s *ProfileService) GetUserProfile(ctx context.Context, userID string) (*model.UserProfile, error) {
	profile, err := s.repo.GetUserProfile(ctx, userID)
	if err != nil {
		return nil, err
	}

	return profile, nil
}

// UpdateProfile updates a user profile
func (s *ProfileService) UpdateProfile(ctx context.Context, userID string, req *profile.UpdateProfileRequest) (*model.UserProfile, error) {
	// Validate request
	if req.Name == "" && req.PublicBio == "" {
		return nil, errors.NewValidationError("at least one field must be provided")
	}

	profile, err := s.repo.UpdateUserProfile(ctx, userID, req.Name, req.PublicBio)
	if err != nil {
		return nil, err
	}

	return profile, nil
}

// UpdatePreferences updates user preferences
func (s *ProfileService) UpdatePreferences(ctx context.Context, userID string, req *profile.UpdatePreferencesRequest) (*prefModel.UserPreferences, error) {
	prefs, err := s.repo.UpdateUserPreferences(ctx, userID, req.NotifyOnLogin, req.Locale)
	if err != nil {
		return nil, err
	}

	return prefs, nil
}

// DeleteProfile schedules account deletion
func (s *ProfileService) DeleteProfile(ctx context.Context, userID string) error {
	err := s.repo.ScheduleAccountDeletion(ctx, userID)
	if err != nil {
		return err
	}

	return nil
}

// OptOut handles opt-out requests from certain features
func (s *ProfileService) OptOut(ctx context.Context, userID string, req *profile.OptOutRequest) error {
	return s.repo.OptOut(ctx, userID, req)
}

// Anonymize anonymizes user personal data
func (s *ProfileService) Anonymize(ctx context.Context, userID string) (*profile.AnonymizeResponse, error) {
	return s.repo.Anonymize(ctx, userID)
}

// RequestDeletion requests account deletion
func (s *ProfileService) RequestDeletion(ctx context.Context, userID string, req *profile.DeleteRequest) (*profile.DeleteResponse, error) {
	return s.repo.RequestDeletion(ctx, userID, req)
}

// SearchProfiles searches for user profiles
func (s *ProfileService) SearchProfiles(ctx context.Context, query string, limit, offset int) ([]*model.UserProfile, int, error) {
	if query == "" {
		return nil, 0, errors.NewValidationError("search query is required")
	}

	profiles, count, err := s.repo.SearchUserProfiles(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return profiles, count, nil
}

