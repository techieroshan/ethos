package service

import (
	"context"

	"ethos/internal/appeal/model"
	"ethos/internal/appeal/repository"
)

// AppealService implements the Service interface
type AppealService struct {
	repo repository.Repository
}

// NewAppealService creates a new appeal service
func NewAppealService(repo repository.Repository) Service {
	return &AppealService{repo: repo}
}

// SubmitAppeal allows a user to submit an appeal
func (s *AppealService) SubmitAppeal(ctx context.Context, userID string, req *SubmitAppealRequest) (*model.Appeal, error) {
	return s.repo.SubmitAppeal(ctx, userID, req.Type, req.ReferenceID, req.Description)
}

// GetUserAppeals retrieves appeals for a specific user
func (s *AppealService) GetUserAppeals(ctx context.Context, userID string, limit, offset int) ([]*model.Appeal, int, error) {
	return s.repo.GetUserAppeals(ctx, userID, limit, offset)
}

// GetAppealByID retrieves a specific appeal by ID (only if user owns it)
func (s *AppealService) GetAppealByID(ctx context.Context, userID, appealID string) (*model.Appeal, error) {
	return s.repo.GetAppealByID(ctx, userID, appealID)
}
