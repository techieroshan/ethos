package repository

import (
	"context"

	"ethos/internal/dashboard/model"
)

// Repository defines the interface for dashboard data access
type Repository interface {
	// GetDashboard retrieves a dashboard snapshot
	GetDashboard(ctx context.Context, userID string) (*model.DashboardSnapshot, error)
}

