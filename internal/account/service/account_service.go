package service

import (
	"context"

	"ethos/internal/account/model"
)

// Service defines the interface for account/security business logic
type Service interface {
	// GetSecurityEvents retrieves security events for a user
	GetSecurityEvents(ctx context.Context, userID string, limit, offset int) ([]*model.SecurityEvent, int, error)

	// GetExportStatus retrieves the status of a data export
	GetExportStatus(ctx context.Context, userID, exportID string) (*model.DataExport, error)

	// Disable2FA disables two-factor authentication
	Disable2FA(ctx context.Context, userID string) error
}

