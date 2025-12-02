package converter

import (
	"time"

	notificationpb "ethos/api/proto/notifications"
	notifModel "ethos/internal/notifications/model"
)

// ProtoToNotificationType converts proto NotificationType to domain model
func ProtoToNotificationType(pb notificationpb.NotificationType) notifModel.NotificationType {
	switch pb {
	case notificationpb.NotificationType_NOTIFICATION_TYPE_FEEDBACK_REPLY:
		return notifModel.NotificationTypeFeedbackReply
	case notificationpb.NotificationType_NOTIFICATION_TYPE_FEEDBACK_RECEIVED:
		return notifModel.NotificationTypeFeedbackReceived
	case notificationpb.NotificationType_NOTIFICATION_TYPE_NEW_COMMENT:
		return notifModel.NotificationTypeNewComment
	case notificationpb.NotificationType_NOTIFICATION_TYPE_SYSTEM_ALERT:
		return notifModel.NotificationTypeSystemAlert
	case notificationpb.NotificationType_NOTIFICATION_TYPE_REMINDER:
		return notifModel.NotificationTypeReminder
	case notificationpb.NotificationType_NOTIFICATION_TYPE_OTHER:
		return notifModel.NotificationTypeOther
	default:
		return notifModel.NotificationTypeOther
	}
}

// ProtoToNotification converts proto Notification to domain model
func ProtoToNotification(pb *notificationpb.Notification) *notifModel.Notification {
	if pb == nil {
		return nil
	}

	notif := &notifModel.Notification{
		NotificationID: pb.NotificationId,
		Type:           ProtoToNotificationType(pb.Type),
		Message:        pb.Message,
		Read:           pb.Read,
	}

	// Convert timestamp
	if pb.CreatedAt != nil {
		notif.CreatedAt = pb.CreatedAt.AsTime()
	} else {
		notif.CreatedAt = time.Now()
	}

	return notif
}

// ProtoToNotificationPreferences converts proto NotificationPreferences to domain model
func ProtoToNotificationPreferences(pb *notificationpb.NotificationPreferences) *notifModel.NotificationPreferences {
	if pb == nil {
		return nil
	}

	return &notifModel.NotificationPreferences{
		Email: pb.Email,
		Push:  pb.Push,
		InApp: pb.InApp,
	}
}

