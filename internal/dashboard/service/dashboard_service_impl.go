package service

import (
	"context"

	"ethos/internal/dashboard/model"
	"ethos/internal/dashboard/repository"
)

// DashboardService implements the Service interface
type DashboardService struct {
	client DashboardClient // Can be REST or gRPC client
}

// NewDashboardService creates a new dashboard service with REST client
func NewDashboardService(repo repository.Repository) Service {
	return &DashboardService{
		client: NewRESTDashboardClient(repo),
	}
}

// NewDashboardServiceWithClient creates a dashboard service with a custom client (REST or gRPC)
func NewDashboardServiceWithClient(client DashboardClient) Service {
	return &DashboardService{
		client: client,
	}
}

// GetDashboard retrieves a dashboard snapshot for a user
func (s *DashboardService) GetDashboard(ctx context.Context, userID string) (*model.DashboardSnapshot, error) {
	return s.client.GetDashboard(ctx, userID)
}

