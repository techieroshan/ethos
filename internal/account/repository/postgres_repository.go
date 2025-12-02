package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"ethos/internal/account/model"
	"ethos/internal/database"
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

// GetSecurityEvents retrieves security events
func (r *PostgresRepository) GetSecurityEvents(ctx context.Context, userID string, limit, offset int) ([]*model.SecurityEvent, int, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "repository.GetSecurityEvents")
	defer span.End()

	// Get total count
	var totalCount int
	countQuery := `SELECT COUNT(*) FROM security_events WHERE user_id = $1`
	err := r.db.Pool.QueryRow(ctx, countQuery, userID).Scan(&totalCount)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, 0, errors.WrapError(err, "failed to count security events")
	}

	// Get events
	query := `
		SELECT event_id, type, timestamp, ip, location
		FROM security_events
		WHERE user_id = $1
		ORDER BY timestamp DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Pool.Query(ctx, query, userID, limit, offset)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, 0, errors.WrapError(err, "failed to get security events")
	}
	defer rows.Close()

	var events []*model.SecurityEvent
	for rows.Next() {
		event := &model.SecurityEvent{}
		err := rows.Scan(
			&event.EventID,
			&event.Type,
			&event.Timestamp,
			&event.IP,
			&event.Location,
		)
		if err != nil {
			span.RecordError(err)
			continue
		}
		events = append(events, event)
	}

	if err = rows.Err(); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, 0, errors.WrapError(err, "failed to iterate security events")
	}

	span.SetStatus(codes.Ok, "")
	return events, totalCount, nil
}

// GetExportStatus retrieves export status
func (r *PostgresRepository) GetExportStatus(ctx context.Context, userID, exportID string) (*model.DataExport, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "repository.GetExportStatus")
	defer span.End()

	query := `
		SELECT export_id, status, download_url, expires_at
		FROM data_exports
		WHERE export_id = $1 AND user_id = $2
	`

	export := &model.DataExport{}
	var expiresAt *time.Time

	err := r.db.Pool.QueryRow(ctx, query, exportID, userID).Scan(
		&export.ExportID,
		&export.Status,
		&export.DownloadURL,
		&expiresAt,
	)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		if err == pgx.ErrNoRows {
			return nil, errors.NewValidationError("export not found")
		}
		return nil, errors.WrapError(err, "failed to get export status")
	}

	export.ExpiresAt = expiresAt
	span.SetStatus(codes.Ok, "")
	return export, nil
}

// Disable2FA disables two-factor authentication
func (r *PostgresRepository) Disable2FA(ctx context.Context, userID string) error {
	ctx, span := otel.Tracer("repository").Start(ctx, "repository.Disable2FA")
	defer span.End()

	query := `
		UPDATE users
		SET two_factor_enabled = FALSE, two_factor_secret = NULL
		WHERE id = $1
	`

	result, err := r.db.Pool.Exec(ctx, query, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return errors.WrapError(err, "failed to disable 2FA")
	}

	if result.RowsAffected() == 0 {
		return errors.ErrUserNotFound
	}

	span.SetStatus(codes.Ok, "")
	return nil
}

