package service

import (
	"context"

	"ethos/internal/notifications/model"
)

// UpdatePreferencesRequest represents a notification preferences update request
type UpdatePreferencesRequest struct {
	Email *bool `json:"email"`
	Push  *bool `json:"push"`
	InApp *bool `json:"in_app"`
}

// Service defines the interface for notification business logic
type Service interface {
	// GetNotifications retrieves notifications for a user
	GetNotifications(ctx context.Context, userID string, limit, offset int) ([]*model.Notification, int, int, error)

	// MarkAsRead marks a notification as read or unread
	MarkAsRead(ctx context.Context, userID, notificationID string, read bool) error

	// MarkAllAsRead marks all notifications as read for a user
	MarkAllAsRead(ctx context.Context, userID string) error

	// GetPreferences retrieves notification preferences for a user
	GetPreferences(ctx context.Context, userID string) (*model.NotificationPreferences, error)

	// UpdatePreferences updates notification preferences
	UpdatePreferences(ctx context.Context, userID string, req *UpdatePreferencesRequest) (*model.NotificationPreferences, error)
}

