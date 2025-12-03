package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"

	"ethos/internal/appeal/model"
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

// SubmitAppeal creates a new appeal
func (r *PostgresRepository) SubmitAppeal(ctx context.Context, userID string, appealType model.AppealType, referenceID *string, description string) (*model.Appeal, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "repository.SubmitAppeal")
	defer span.End()

	now := time.Now()
	appeal := &model.Appeal{
		UserID:      userID,
		Type:        appealType,
		ReferenceID: referenceID,
		Description: description,
		Status:      model.AppealStatusPending,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	query := `
		INSERT INTO user_appeals (user_id, type, reference_id, description, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING appeal_id
	`

	err := r.db.Pool.QueryRow(ctx, query,
		appeal.UserID,
		string(appeal.Type),
		appeal.ReferenceID,
		appeal.Description,
		string(appeal.Status),
		appeal.CreatedAt,
		appeal.UpdatedAt,
	).Scan(&appeal.AppealID)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, errors.WrapError(err, "failed to submit appeal")
	}

	span.SetStatus(codes.Ok, "")
	return appeal, nil
}

// GetUserAppeals retrieves appeals for a specific user
func (r *PostgresRepository) GetUserAppeals(ctx context.Context, userID string, limit, offset int) ([]*model.Appeal, int, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "repository.GetUserAppeals")
	defer span.End()

	// Get total count
	var totalCount int
	countQuery := `SELECT COUNT(*) FROM user_appeals WHERE user_id = $1`
	err := r.db.Pool.QueryRow(ctx, countQuery, userID).Scan(&totalCount)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, 0, errors.WrapError(err, "failed to count appeals")
	}

	// Get appeals
	query := `
		SELECT appeal_id, user_id, type, reference_id, description, status, admin_notes, created_at, updated_at, resolved_at
		FROM user_appeals
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Pool.Query(ctx, query, userID, limit, offset)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, 0, errors.WrapError(err, "failed to get appeals")
	}
	defer rows.Close()

	var appeals []*model.Appeal
	for rows.Next() {
		appeal := &model.Appeal{}
		var appealType, status string

		err := rows.Scan(
			&appeal.AppealID,
			&appeal.UserID,
			&appealType,
			&appeal.ReferenceID,
			&appeal.Description,
			&status,
			&appeal.AdminNotes,
			&appeal.CreatedAt,
			&appeal.UpdatedAt,
			&appeal.ResolvedAt,
		)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return nil, 0, errors.WrapError(err, "failed to scan appeal")
		}

		appeal.Type = model.AppealType(appealType)
		appeal.Status = model.AppealStatus(status)

		appeals = append(appeals, appeal)
	}

	if err = rows.Err(); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, 0, errors.WrapError(err, "error iterating appeals")
	}

	span.SetStatus(codes.Ok, "")
	return appeals, totalCount, nil
}

// GetAppealByID retrieves a specific appeal by ID (only if user owns it)
func (r *PostgresRepository) GetAppealByID(ctx context.Context, userID, appealID string) (*model.Appeal, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "repository.GetAppealByID")
	defer span.End()

	query := `
		SELECT appeal_id, user_id, type, reference_id, description, status, admin_notes, created_at, updated_at, resolved_at
		FROM user_appeals
		WHERE appeal_id = $1 AND user_id = $2
	`

	appeal := &model.Appeal{}
	var appealType, status string

	err := r.db.Pool.QueryRow(ctx, query, appealID, userID).Scan(
		&appeal.AppealID,
		&appeal.UserID,
		&appealType,
		&appeal.ReferenceID,
		&appeal.Description,
		&status,
		&appeal.AdminNotes,
		&appeal.CreatedAt,
		&appeal.UpdatedAt,
		&appeal.ResolvedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("appeal not found")
		}
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, errors.WrapError(err, "failed to get appeal")
	}

	appeal.Type = model.AppealType(appealType)
	appeal.Status = model.AppealStatus(status)

	span.SetStatus(codes.Ok, "")
	return appeal, nil
}

// UpdateAppealStatus updates the status of an appeal (for admins/moderators)
func (r *PostgresRepository) UpdateAppealStatus(ctx context.Context, appealID string, status model.AppealStatus, adminNotes *string) error {
	ctx, span := otel.Tracer("repository").Start(ctx, "repository.UpdateAppealStatus")
	defer span.End()

	now := time.Now()
	var resolvedAt *time.Time
	if status == model.AppealStatusApproved || status == model.AppealStatusRejected || status == model.AppealStatusClosed {
		resolvedAt = &now
	}

	query := `
		UPDATE user_appeals
		SET status = $1, admin_notes = $2, updated_at = $3, resolved_at = $4
		WHERE appeal_id = $5
	`

	_, err := r.db.Pool.Exec(ctx, query, string(status), adminNotes, now, resolvedAt, appealID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return errors.WrapError(err, "failed to update appeal status")
	}

	span.SetStatus(codes.Ok, "")
	return nil
}
