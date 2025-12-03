package service

import (
	"context"

	"ethos/internal/notifications/model"
	"ethos/internal/notifications/repository"
)

// NotificationService implements the Service interface
type NotificationService struct {
	client NotificationClient // Can be REST or gRPC client
	repo   repository.Repository // Kept for write operations (GetPreferences, UpdatePreferences)
}

// NewNotificationService creates a new notification service with REST client
func NewNotificationService(repo repository.Repository) Service {
	return &NotificationService{
		client: NewRESTNotificationClient(repo),
		repo:   repo,
	}
}

// NewNotificationServiceWithClient creates a notification service with a custom client (REST or gRPC)
func NewNotificationServiceWithClient(client NotificationClient, repo repository.Repository) Service {
	return &NotificationService{
		client: client,
		repo:   repo,
	}
}

// GetNotifications retrieves notifications for a user
func (s *NotificationService) GetNotifications(ctx context.Context, userID string, limit, offset int) ([]*model.Notification, int, int, error) {
	return s.client.GetNotifications(ctx, userID, limit, offset)
}

// GetPreferences retrieves notification preferences for a user
func (s *NotificationService) GetPreferences(ctx context.Context, userID string) (*model.NotificationPreferences, error) {
	prefs, err := s.repo.GetPreferences(ctx, userID)
	if err != nil {
		return nil, err
	}

	return prefs, nil
}

// MarkAsRead marks a notification as read or unread
func (s *NotificationService) MarkAsRead(ctx context.Context, userID, notificationID string, read bool) error {
	return s.repo.MarkAsRead(ctx, userID, notificationID, read)
}

// MarkAllAsRead marks all notifications as read for a user
func (s *NotificationService) MarkAllAsRead(ctx context.Context, userID string) error {
	return s.repo.MarkAllAsRead(ctx, userID)
}

// UpdatePreferences updates notification preferences
func (s *NotificationService) UpdatePreferences(ctx context.Context, userID string, req *UpdatePreferencesRequest) (*model.NotificationPreferences, error) {
	prefs, err := s.repo.UpdatePreferences(ctx, userID, req.Email, req.Push, req.InApp)
	if err != nil {
		return nil, err
	}

	return prefs, nil
}

