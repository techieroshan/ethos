package service

import (
	"context"

	"ethos/internal/account/model"
	"ethos/internal/account/repository"
)

// AccountService implements the Service interface
type AccountService struct {
	repo repository.Repository
}

// NewAccountService creates a new account service
func NewAccountService(repo repository.Repository) Service {
	return &AccountService{
		repo: repo,
	}
}

// GetSecurityEvents retrieves security events for a user
func (s *AccountService) GetSecurityEvents(ctx context.Context, userID string, limit, offset int) ([]*model.SecurityEvent, int, error) {
	events, count, err := s.repo.GetSecurityEvents(ctx, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return events, count, nil
}

// GetExportStatus retrieves the status of a data export
func (s *AccountService) GetExportStatus(ctx context.Context, userID, exportID string) (*model.DataExport, error) {
	export, err := s.repo.GetExportStatus(ctx, userID, exportID)
	if err != nil {
		return nil, err
	}

	return export, nil
}

// Disable2FA disables two-factor authentication
func (s *AccountService) Disable2FA(ctx context.Context, userID string) error {
	err := s.repo.Disable2FA(ctx, userID)
	if err != nil {
		return err
	}

	return nil
}

