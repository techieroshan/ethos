package service

import (
	"context"

	"ethos/internal/notifications/model"
	"ethos/internal/notifications/repository"
)

// NotificationService implements the Service interface
type NotificationService struct {
	repo repository.Repository
}

// NewNotificationService creates a new notification service
func NewNotificationService(repo repository.Repository) Service {
	return &NotificationService{
		repo: repo,
	}
}

// GetNotifications retrieves notifications for a user
func (s *NotificationService) GetNotifications(ctx context.Context, userID string, limit, offset int) ([]*model.Notification, int, int, error) {
	notifications, count, unreadCount, err := s.repo.GetNotifications(ctx, userID, limit, offset)
	if err != nil {
		return nil, 0, 0, err
	}

	return notifications, count, unreadCount, nil
}

// GetPreferences retrieves notification preferences for a user
func (s *NotificationService) GetPreferences(ctx context.Context, userID string) (*model.NotificationPreferences, error) {
	prefs, err := s.repo.GetPreferences(ctx, userID)
	if err != nil {
		return nil, err
	}

	return prefs, nil
}

// UpdatePreferences updates notification preferences
func (s *NotificationService) UpdatePreferences(ctx context.Context, userID string, req *UpdatePreferencesRequest) (*model.NotificationPreferences, error) {
	prefs, err := s.repo.UpdatePreferences(ctx, userID, req.Email, req.Push, req.InApp)
	if err != nil {
		return nil, err
	}

	return prefs, nil
}

