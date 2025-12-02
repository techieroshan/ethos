package service

import (
	"context"

	"ethos/internal/dashboard/model"
	"ethos/internal/dashboard/repository"
)

// DashboardService implements the Service interface
type DashboardService struct {
	repo repository.Repository
}

// NewDashboardService creates a new dashboard service
func NewDashboardService(repo repository.Repository) Service {
	return &DashboardService{
		repo: repo,
	}
}

// GetDashboard retrieves a dashboard snapshot for a user
func (s *DashboardService) GetDashboard(ctx context.Context, userID string) (*model.DashboardSnapshot, error) {
	snapshot, err := s.repo.GetDashboard(ctx, userID)
	if err != nil {
		return nil, err
	}

	return snapshot, nil
}

