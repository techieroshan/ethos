package model

import "time"

// NotificationType represents types of notifications
type NotificationType string

const (
	NotificationTypeFeedbackReply    NotificationType = "feedback_reply"
	NotificationTypeFeedbackReceived NotificationType = "feedback_received"
	NotificationTypeNewComment       NotificationType = "new_comment"
	NotificationTypeSystemAlert      NotificationType = "system_alert"
	NotificationTypeReminder        NotificationType = "reminder"
	NotificationTypeOther            NotificationType = "other"
)

// Notification represents a notification
type Notification struct {
	NotificationID string           `json:"notification_id"`
	Type           NotificationType `json:"type"`
	Message        string           `json:"message"`
	Read           bool             `json:"read"`
	CreatedAt      time.Time        `json:"created_at"`
}

// NotificationPreferences represents notification delivery preferences
type NotificationPreferences struct {
	Email  bool `json:"email"`
	Push   bool `json:"push"`
	InApp  bool `json:"in_app"`
}

