package service

import (
	"context"

	notifModel "ethos/internal/notifications/model"
	"ethos/internal/notifications/repository"
	notificationpb "ethos/api/proto/notifications"
	"ethos/pkg/grpc/converter"
)

// NotificationClient defines the interface for notification data access (REST or gRPC)
type NotificationClient interface {
	// GetNotifications retrieves notifications for a user
	GetNotifications(ctx context.Context, userID string, limit, offset int) ([]*notifModel.Notification, int, int, error)
}

// RESTNotificationClient implements NotificationClient using REST (current repository)
type RESTNotificationClient struct {
	repo repository.Repository
}

// NewRESTNotificationClient creates a new REST notification client
func NewRESTNotificationClient(repo repository.Repository) NotificationClient {
	return &RESTNotificationClient{repo: repo}
}

// GetNotifications implements NotificationClient interface using REST
func (c *RESTNotificationClient) GetNotifications(ctx context.Context, userID string, limit, offset int) ([]*notifModel.Notification, int, int, error) {
	return c.repo.GetNotifications(ctx, userID, limit, offset)
}

// GRPCNotificationClient implements NotificationClient using gRPC
type GRPCNotificationClient struct {
	client notificationpb.NotificationServiceClient
}

// NewGRPCNotificationClient creates a new gRPC notification client
func NewGRPCNotificationClient(client notificationpb.NotificationServiceClient) NotificationClient {
	return &GRPCNotificationClient{client: client}
}

// GetNotifications implements NotificationClient interface using gRPC
func (c *GRPCNotificationClient) GetNotifications(ctx context.Context, userID string, limit, offset int) ([]*notifModel.Notification, int, int, error) {
	req := &notificationpb.GetNotificationsRequest{
		UserId: userID,
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	resp, err := c.client.GetNotifications(ctx, req)
	if err != nil {
		return nil, 0, 0, err
	}

	notifications := make([]*notifModel.Notification, 0, len(resp.Notifications))
	unreadCount := 0
	for _, pbNotif := range resp.Notifications {
		notif := converter.ProtoToNotification(pbNotif)
		if notif != nil {
			notifications = append(notifications, notif)
			if !notif.Read {
				unreadCount++
			}
		}
	}

	return notifications, int(resp.Count), unreadCount, nil
}

