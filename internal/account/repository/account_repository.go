package repository

import (
	"context"

	"ethos/internal/account/model"
)

// Repository defines the interface for account/security data access
type Repository interface {
	// GetSecurityEvents retrieves security events
	GetSecurityEvents(ctx context.Context, userID string, limit, offset int) ([]*model.SecurityEvent, int, error)

	// GetExportStatus retrieves export status
	GetExportStatus(ctx context.Context, userID, exportID string) (*model.DataExport, error)

	// Disable2FA disables two-factor authentication
	Disable2FA(ctx context.Context, userID string) error
}

