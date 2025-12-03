package repository

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"

	"ethos/internal/auth/model"
	"ethos/internal/database"
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

// GetUserByEmail retrieves a user by email address
func (r *PostgresRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "repository.GetUserByEmail")
	defer span.End()

	var user model.User
	query := `
		SELECT id, email, password_hash, name, email_verified, public_bio, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	err := r.db.Pool.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Name,
		&user.EmailVerified,
		&user.PublicBio,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		if err == pgx.ErrNoRows {
			return nil, errors.ErrUserNotFound
		}
		return nil, errors.WrapError(err, "failed to get user by email")
	}

	span.SetStatus(codes.Ok, "")
	return &user, nil
}

// GetUserByID retrieves a user by ID
func (r *PostgresRepository) GetUserByID(ctx context.Context, userID string) (*model.User, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "repository.GetUserByID")
	defer span.End()

	var user model.User
	query := `
		SELECT id, email, password_hash, name, email_verified, public_bio, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	err := r.db.Pool.QueryRow(ctx, query, userID).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Name,
		&user.EmailVerified,
		&user.PublicBio,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		if err == pgx.ErrNoRows {
			return nil, errors.ErrUserNotFound
		}
		return nil, errors.WrapError(err, "failed to get user by ID")
	}

	span.SetStatus(codes.Ok, "")
	return &user, nil
}

// CreateUser creates a new user
func (r *PostgresRepository) CreateUser(ctx context.Context, user *model.User) error {
	ctx, span := otel.Tracer("repository").Start(ctx, "repository.CreateUser")
	defer span.End()

	if user.ID == "" {
		user.ID = "user-" + uuid.New().String()
	}

	query := `
		INSERT INTO users (id, email, password_hash, name, email_verified, public_bio, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	now := time.Now()
	_, err := r.db.Pool.Exec(ctx, query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.Name,
		user.EmailVerified,
		user.PublicBio,
		now,
		now,
	)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		// Check for unique constraint violation (PostgreSQL error code 23505)
		if strings.Contains(err.Error(), "duplicate key value") || strings.Contains(err.Error(), "23505") {
			return errors.ErrEmailAlreadyExists
		}
		return errors.WrapError(err, "failed to create user")
	}

	user.CreatedAt = now
	user.UpdatedAt = now
	span.SetStatus(codes.Ok, "")
	return nil
}

// UpdateUser updates an existing user
func (r *PostgresRepository) UpdateUser(ctx context.Context, user *model.User) error {
	ctx, span := otel.Tracer("repository").Start(ctx, "repository.UpdateUser")
	defer span.End()

	query := `
		UPDATE users
		SET email = $1, password_hash = $2, name = $3, email_verified = $4, public_bio = $5, updated_at = $6
		WHERE id = $7
	`

	now := time.Now()
	result, err := r.db.Pool.Exec(ctx, query,
		user.Email,
		user.PasswordHash,
		user.Name,
		user.EmailVerified,
		user.PublicBio,
		now,
		user.ID,
	)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return errors.WrapError(err, "failed to update user")
	}

	if result.RowsAffected() == 0 {
		span.RecordError(errors.ErrUserNotFound)
		span.SetStatus(codes.Error, "user not found")
		return errors.ErrUserNotFound
	}

	user.UpdatedAt = now
	span.SetStatus(codes.Ok, "")
	return nil
}

// SaveRefreshToken saves a refresh token
func (r *PostgresRepository) SaveRefreshToken(ctx context.Context, userID, tokenHash string, expiresAt int64) error {
	ctx, span := otel.Tracer("repository").Start(ctx, "repository.SaveRefreshToken")
	defer span.End()

	tokenID := "token-" + uuid.New().String()
	expiresAtTime := time.Unix(expiresAt, 0)

	query := `
		INSERT INTO refresh_tokens (token_id, user_id, token_hash, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		tokenID,
		userID,
		tokenHash,
		expiresAtTime,
		time.Now(),
	)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return errors.WrapError(err, "failed to save refresh token")
	}

	span.SetStatus(codes.Ok, "")
	return nil
}

// GetRefreshToken retrieves a refresh token by hash
func (r *PostgresRepository) GetRefreshToken(ctx context.Context, tokenHash string) (string, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "repository.GetRefreshToken")
	defer span.End()

	var userID string
	query := `
		SELECT user_id
		FROM refresh_tokens
		WHERE token_hash = $1 AND expires_at > NOW()
	`

	err := r.db.Pool.QueryRow(ctx, query, tokenHash).Scan(&userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		if err == pgx.ErrNoRows {
			return "", errors.ErrTokenInvalid
		}
		return "", errors.WrapError(err, "failed to get refresh token")
	}

	span.SetStatus(codes.Ok, "")
	return userID, nil
}

// DeleteRefreshToken deletes a refresh token
func (r *PostgresRepository) DeleteRefreshToken(ctx context.Context, tokenHash string) error {
	ctx, span := otel.Tracer("repository").Start(ctx, "repository.DeleteRefreshToken")
	defer span.End()

	query := `DELETE FROM refresh_tokens WHERE token_hash = $1`

	_, err := r.db.Pool.Exec(ctx, query, tokenHash)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return errors.WrapError(err, "failed to delete refresh token")
	}

	span.SetStatus(codes.Ok, "")
	return nil
}

// hashToken creates a SHA256 hash of a token for storage
func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
