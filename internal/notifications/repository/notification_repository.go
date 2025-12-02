package repository

import (
	"context"

	"ethos/internal/notifications/model"
)

// Repository defines the interface for notification data access
type Repository interface {
	// GetNotifications retrieves notifications for a user
	GetNotifications(ctx context.Context, userID string, limit, offset int) ([]*model.Notification, int, int, error)

	// GetPreferences retrieves notification preferences
	GetPreferences(ctx context.Context, userID string) (*model.NotificationPreferences, error)

	// UpdatePreferences updates notification preferences
	UpdatePreferences(ctx context.Context, userID string, email, push, inApp *bool) (*model.NotificationPreferences, error)
}

