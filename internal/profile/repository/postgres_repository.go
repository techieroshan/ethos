package repository

import (
	"context"
	"net/http"
	"time"

	"ethos/internal/auth/model"
	"ethos/internal/database"
	"ethos/internal/profile"
	prefModel "ethos/internal/profile/model"
	"ethos/pkg/errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

// PostgresRepository implements the Repository interface using PostgreSQL
type PostgresRepository struct {
	db *database.DB
}

// NewPostgresRepository creates a new PostgreSQL repository
func NewPostgresRepository(db *database.DB) Repository {
	return &PostgresRepository{db: db}
}

// GetUserProfile retrieves a user profile by ID
func (r *PostgresRepository) GetUserProfile(ctx context.Context, userID string) (*model.UserProfile, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "repository.GetUserProfile")
	defer span.End()

	var user model.User
	query := `
		SELECT id, email, password_hash, first_name, last_name, email_verified, public_bio, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	err := r.db.Pool.QueryRow(ctx, query, userID).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.EmailVerified,
		&user.PublicBio,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		if err == pgx.ErrNoRows {
			return nil, errors.ErrProfileNotFound
		}
		return nil, errors.WrapError(err, "failed to get user profile")
	}

	span.SetStatus(codes.Ok, "")
	return user.ToProfile(), nil
}

// UpdateUserProfile updates a user profile
func (r *PostgresRepository) UpdateUserProfile(ctx context.Context, userID string, name, publicBio string) (*model.UserProfile, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "repository.UpdateUserProfile")
	defer span.End()

	now := time.Now()
	// Only update fields that are provided (not empty)
	var user model.User
	var err error

	if name != "" && publicBio != "" {
		query := `
			UPDATE users
			SET first_name = $1, public_bio = $2, updated_at = $3
			WHERE id = $4
			RETURNING id, email, password_hash, first_name, last_name, email_verified, public_bio, created_at, updated_at
		`
		err = r.db.Pool.QueryRow(ctx, query, name, publicBio, now, userID).Scan(
			&user.ID,
			&user.Email,
			&user.PasswordHash,
			&user.FirstName,
			&user.LastName,
			&user.EmailVerified,
			&user.PublicBio,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
	} else if name != "" {
		query := `
			UPDATE users
			SET first_name = $1, updated_at = $2
			WHERE id = $3
			RETURNING id, email, password_hash, first_name, last_name, email_verified, public_bio, created_at, updated_at
		`
		err = r.db.Pool.QueryRow(ctx, query, name, now, userID).Scan(
			&user.ID,
			&user.Email,
			&user.PasswordHash,
			&user.FirstName,
			&user.LastName,
			&user.EmailVerified,
			&user.PublicBio,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
	} else if publicBio != "" {
		query := `
			UPDATE users
			SET public_bio = $1, updated_at = $2
			WHERE id = $3
			RETURNING id, email, password_hash, first_name, last_name, email_verified, public_bio, created_at, updated_at
		`
		err = r.db.Pool.QueryRow(ctx, query, publicBio, now, userID).Scan(
			&user.ID,
			&user.Email,
			&user.PasswordHash,
			&user.FirstName,
			&user.LastName,
			&user.EmailVerified,
			&user.PublicBio,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
	} else {
		// If both are empty, just return current profile
		return r.GetUserProfile(ctx, userID)
	}

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		if err == pgx.ErrNoRows {
			return nil, errors.ErrProfileNotFound
		}
		return nil, errors.WrapError(err, "failed to update user profile")
	}

	span.SetStatus(codes.Ok, "")
	return user.ToProfile(), nil
}

// UpdateUserPreferences updates user preferences
func (r *PostgresRepository) UpdateUserPreferences(ctx context.Context, userID string, notifyOnLogin *bool, locale *string) (*prefModel.UserPreferences, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "repository.UpdateUserPreferences")
	defer span.End()

	now := time.Now()

	// Check if preferences exist
	var exists bool
	err := r.db.Pool.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM user_preferences WHERE user_id = $1)", userID).Scan(&exists)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, errors.WrapError(err, "failed to check preferences existence")
	}

	var prefs prefModel.UserPreferences
	if exists {
		// Update existing preferences
		query := `
			UPDATE user_preferences
			SET notify_on_login = COALESCE($1, notify_on_login),
			    locale = COALESCE($2, locale),
			    updated_at = $3
			WHERE user_id = $4
			RETURNING user_id, notify_on_login, locale
		`
		err = r.db.Pool.QueryRow(ctx, query, notifyOnLogin, locale, now, userID).Scan(
			&prefs.UserID,
			&prefs.NotifyOnLogin,
			&prefs.Locale,
		)
	} else {
		// Insert new preferences
		notifyVal := true
		localeVal := "en-US"
		if notifyOnLogin != nil {
			notifyVal = *notifyOnLogin
		}
		if locale != nil {
			localeVal = *locale
		}
		query := `
			INSERT INTO user_preferences (user_id, notify_on_login, locale, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING user_id, notify_on_login, locale
		`
		err = r.db.Pool.QueryRow(ctx, query, userID, notifyVal, localeVal, now, now).Scan(
			&prefs.UserID,
			&prefs.NotifyOnLogin,
			&prefs.Locale,
		)
	}

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, errors.WrapError(err, "failed to update user preferences")
	}

	span.SetStatus(codes.Ok, "")
	return &prefs, nil
}

// ScheduleAccountDeletion schedules an account for deletion
func (r *PostgresRepository) ScheduleAccountDeletion(ctx context.Context, userID string) error {
	ctx, span := otel.Tracer("repository").Start(ctx, "repository.ScheduleAccountDeletion")
	defer span.End()

	deletionID := "del-" + uuid.New().String()
	scheduledAt := time.Now().Add(30 * 24 * time.Hour) // 30 days from now

	query := `
		INSERT INTO account_deletions (deletion_id, user_id, scheduled_at, status, created_at)
		VALUES ($1, $2, $3, 'pending', $4)
	`

	_, err := r.db.Pool.Exec(ctx, query, deletionID, userID, scheduledAt, time.Now())
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return errors.WrapError(err, "failed to schedule account deletion")
	}

	span.SetStatus(codes.Ok, "")
	return nil
}

// SearchUserProfiles searches for user profiles
func (r *PostgresRepository) SearchUserProfiles(ctx context.Context, query string, limit, offset int) ([]*model.UserProfile, int, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "repository.SearchUserProfiles")
	defer span.End()

	searchPattern := "%" + query + "%"

	// Get total count
	var totalCount int
	countQuery := `
		SELECT COUNT(*) FROM users
		WHERE (CONCAT(first_name, ' ', last_name) ILIKE $1 OR email ILIKE $1 OR public_bio ILIKE $1)
	`
	err := r.db.Pool.QueryRow(ctx, countQuery, searchPattern).Scan(&totalCount)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, 0, errors.WrapError(err, "failed to count profiles")
	}

	// Get profiles
	profilesQuery := `
		SELECT id, email, password_hash, first_name, last_name, email_verified, public_bio, created_at, updated_at
		FROM users
		WHERE (CONCAT(first_name, ' ', last_name) ILIKE $1 OR email ILIKE $1 OR public_bio ILIKE $1)
		ORDER BY CONCAT(first_name, ' ', last_name) ASC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Pool.Query(ctx, profilesQuery, searchPattern, limit, offset)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, 0, errors.WrapError(err, "failed to search profiles")
	}
	defer rows.Close()

	var profiles []*model.UserProfile
	for rows.Next() {
		var user model.User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.PasswordHash,
			&user.FirstName,
			&user.LastName,
			&user.EmailVerified,
			&user.PublicBio,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			span.RecordError(err)
			continue
		}
		profiles = append(profiles, user.ToProfile())
	}

	if err = rows.Err(); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, 0, errors.WrapError(err, "failed to iterate profiles")
	}

	span.SetStatus(codes.Ok, "")
	return profiles, totalCount, nil
}

// OptOut handles opt-out requests from certain features
func (r *PostgresRepository) OptOut(ctx context.Context, userID string, req *profile.OptOutRequest) error {
	ctx, span := otel.Tracer("repository").Start(ctx, "repository.OptOut")
	defer span.End()

	// Update the opt_outs array in the users table
	_, err := r.db.Pool.Exec(ctx, `
		UPDATE users
		SET opt_outs = array_append(COALESCE(opt_outs, ARRAY[]::TEXT[]), $1)
		WHERE id = $2 AND NOT ($1 = ANY(COALESCE(opt_outs, ARRAY[]::TEXT[])))
	`, req.From, userID)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return errors.WrapError(err, "failed to opt out")
	}

	span.SetStatus(codes.Ok, "")
	return nil
}

// Anonymize anonymizes user personal data
func (r *PostgresRepository) Anonymize(ctx context.Context, userID string) (*profile.AnonymizeResponse, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "repository.Anonymize")
	defer span.End()

	// In a real implementation, this would:
	// 1. Replace personal data with anonymized versions
	// 2. Queue a background job for complete anonymization
	// 3. Update all references to this user across the system

	// For now, mark as anonymized and return expected completion
	anonymizedAt := time.Now()
	expectedCompletion := anonymizedAt.Add(24 * time.Hour) // 24 hours for full anonymization

	_, err := r.db.Pool.Exec(ctx, `
		UPDATE users
		SET anonymized_at = $1,
		    first_name = 'Anonymous',
		    last_name = 'User',
		    public_bio = NULL,
		    email = CONCAT('anonymous_', id, '@ethos.local')
		WHERE id = $2
	`, anonymizedAt, userID)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, errors.WrapError(err, "failed to anonymize user")
	}

	response := &profile.AnonymizeResponse{
		Status:             "in_progress",
		ExpectedCompletion: expectedCompletion,
	}

	span.SetStatus(codes.Ok, "")
	return response, nil
}

// RequestDeletion requests account deletion
func (r *PostgresRepository) RequestDeletion(ctx context.Context, userID string, req *profile.DeleteRequest) (*profile.DeleteResponse, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "repository.RequestDeletion")
	defer span.End()

	if !req.Confirm {
		return nil, &errors.APIError{
			Message:    "Confirmation required for account deletion",
			Code:       "VALIDATION_FAILED",
			HTTPStatus: http.StatusBadRequest,
		}
	}

	// Mark deletion requested
	deletionRequestedAt := time.Now()
	expectedCompletion := deletionRequestedAt.Add(30 * 24 * time.Hour) // 30 days for deletion

	_, err := r.db.Pool.Exec(ctx, `
		UPDATE users
		SET delete_requested_at = $1
		WHERE id = $2
	`, deletionRequestedAt, userID)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, errors.WrapError(err, "failed to request deletion")
	}

	// In a real implementation, this would also create a record in account_deletions table
	// and queue background jobs for cleanup

	response := &profile.DeleteResponse{
		Status:             "delete_requested",
		ExpectedCompletion: expectedCompletion,
	}

	span.SetStatus(codes.Ok, "")
	return response, nil
}
