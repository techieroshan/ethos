package service

import (
	"context"

	"ethos/internal/dashboard/model"
)

// Service defines the interface for dashboard business logic
type Service interface {
	// GetDashboard retrieves a dashboard snapshot for a user
	GetDashboard(ctx context.Context, userID string) (*model.DashboardSnapshot, error)
}

