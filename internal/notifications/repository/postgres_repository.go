package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"ethos/internal/database"
	"ethos/internal/notifications/model"
	"ethos/pkg/errors"
)

// PostgresRepository implements the Repository interface using PostgreSQL
type PostgresRepository struct {
	db *database.DB
}

// NewPostgresRepository creates a new PostgreSQL repository
func NewPostgresRepository(db *database.DB) Repository {
	return &PostgresRepository{db: db}
}

// GetNotifications retrieves notifications for a user
func (r *PostgresRepository) GetNotifications(ctx context.Context, userID string, limit, offset int) ([]*model.Notification, int, int, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "repository.GetNotifications")
	defer span.End()

	// Get total count
	var totalCount int
	countQuery := `SELECT COUNT(*) FROM notifications WHERE user_id = $1`
	err := r.db.Pool.QueryRow(ctx, countQuery, userID).Scan(&totalCount)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, 0, 0, errors.WrapError(err, "failed to count notifications")
	}

	// Get unread count
	var unreadCount int
	unreadQuery := `SELECT COUNT(*) FROM notifications WHERE user_id = $1 AND read = FALSE`
	err = r.db.Pool.QueryRow(ctx, unreadQuery, userID).Scan(&unreadCount)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, 0, 0, errors.WrapError(err, "failed to count unread notifications")
	}

	// Get notifications
	query := `
		SELECT notification_id, type, message, read, created_at
		FROM notifications
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Pool.Query(ctx, query, userID, limit, offset)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, 0, 0, errors.WrapError(err, "failed to get notifications")
	}
	defer rows.Close()

	var notifications []*model.Notification
	for rows.Next() {
		notification := &model.Notification{}
		var notificationType string

		err := rows.Scan(
			&notification.NotificationID,
			&notificationType,
			&notification.Message,
			&notification.Read,
			&notification.CreatedAt,
		)
		if err != nil {
			span.RecordError(err)
			continue
		}

		notification.Type = model.NotificationType(notificationType)
		notifications = append(notifications, notification)
	}

	if err = rows.Err(); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, 0, 0, errors.WrapError(err, "failed to iterate notifications")
	}

	span.SetStatus(codes.Ok, "")
	return notifications, totalCount, unreadCount, nil
}

// GetPreferences retrieves notification preferences
func (r *PostgresRepository) GetPreferences(ctx context.Context, userID string) (*model.NotificationPreferences, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "repository.GetPreferences")
	defer span.End()

	query := `
		SELECT email, push, in_app
		FROM notification_preferences
		WHERE user_id = $1
	`

	prefs := &model.NotificationPreferences{}
	err := r.db.Pool.QueryRow(ctx, query, userID).Scan(
		&prefs.Email,
		&prefs.Push,
		&prefs.InApp,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			// Return defaults if preferences don't exist
			return &model.NotificationPreferences{
				Email: true,
				Push:  true,
				InApp: true,
			}, nil
		}
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, errors.WrapError(err, "failed to get preferences")
	}

	span.SetStatus(codes.Ok, "")
	return prefs, nil
}

// UpdatePreferences updates notification preferences
func (r *PostgresRepository) UpdatePreferences(ctx context.Context, userID string, email, push, inApp *bool) (*model.NotificationPreferences, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "repository.UpdatePreferences")
	defer span.End()

	now := time.Now()

	// Check if preferences exist
	var exists bool
	err := r.db.Pool.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM notification_preferences WHERE user_id = $1)", userID).Scan(&exists)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, errors.WrapError(err, "failed to check preferences existence")
	}

	var prefs model.NotificationPreferences
	if exists {
		// Update existing preferences
		query := `
			UPDATE notification_preferences
			SET email = COALESCE($1, email),
			    push = COALESCE($2, push),
			    in_app = COALESCE($3, in_app),
			    updated_at = $4
			WHERE user_id = $5
			RETURNING email, push, in_app
		`
		err = r.db.Pool.QueryRow(ctx, query, email, push, inApp, now, userID).Scan(
			&prefs.Email,
			&prefs.Push,
			&prefs.InApp,
		)
	} else {
		// Insert new preferences
		emailVal := true
		pushVal := true
		inAppVal := true
		if email != nil {
			emailVal = *email
		}
		if push != nil {
			pushVal = *push
		}
		if inApp != nil {
			inAppVal = *inApp
		}
		query := `
			INSERT INTO notification_preferences (user_id, email, push, in_app, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING email, push, in_app
		`
		err = r.db.Pool.QueryRow(ctx, query, userID, emailVal, pushVal, inAppVal, now, now).Scan(
			&prefs.Email,
			&prefs.Push,
			&prefs.InApp,
		)
	}

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, errors.WrapError(err, "failed to update preferences")
	}

	span.SetStatus(codes.Ok, "")
	return &prefs, nil
}

// MarkAsRead marks a notification as read or unread
func (r *PostgresRepository) MarkAsRead(ctx context.Context, userID, notificationID string, read bool) error {
	ctx, span := otel.Tracer("repository").Start(ctx, "repository.MarkAsRead")
	defer span.End()

	query := `UPDATE notifications SET read = $1 WHERE notification_id = $2 AND user_id = $3`
	_, err := r.db.Pool.Exec(ctx, query, read, notificationID, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return errors.WrapError(err, "failed to mark notification as read")
	}

	span.SetStatus(codes.Ok, "")
	return nil
}

// MarkAllAsRead marks all notifications as read for a user
func (r *PostgresRepository) MarkAllAsRead(ctx context.Context, userID string) error {
	ctx, span := otel.Tracer("repository").Start(ctx, "repository.MarkAllAsRead")
	defer span.End()

	query := `UPDATE notifications SET read = true WHERE user_id = $1 AND read = false`
	_, err := r.db.Pool.Exec(ctx, query, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return errors.WrapError(err, "failed to mark all notifications as read")
	}

	span.SetStatus(codes.Ok, "")
	return nil
}

