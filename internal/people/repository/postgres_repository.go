package repository

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"ethos/internal/auth/model"
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

// SearchPeople searches for people
func (r *PostgresRepository) SearchPeople(ctx context.Context, query string, limit, offset int) ([]*model.UserProfile, int, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "repository.SearchPeople")
	defer span.End()

	searchPattern := "%" + query + "%"
	
	// Get total count
	var totalCount int
	countQuery := `
		SELECT COUNT(*) FROM users
		WHERE (name ILIKE $1 OR email ILIKE $1)
	`
	err := r.db.Pool.QueryRow(ctx, countQuery, searchPattern).Scan(&totalCount)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, 0, errors.WrapError(err, "failed to count people")
	}

	// Get profiles
	profilesQuery := `
		SELECT id, email, password_hash, name, email_verified, public_bio, created_at, updated_at
		FROM users
		WHERE (name ILIKE $1 OR email ILIKE $1)
		ORDER BY name ASC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Pool.Query(ctx, profilesQuery, searchPattern, limit, offset)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, 0, errors.WrapError(err, "failed to search people")
	}
	defer rows.Close()

	var profiles []*model.UserProfile
	for rows.Next() {
		var user model.User
		err := rows.Scan(
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
			continue
		}
		profiles = append(profiles, user.ToProfile())
	}

	if err = rows.Err(); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, 0, errors.WrapError(err, "failed to iterate people")
	}

	span.SetStatus(codes.Ok, "")
	return profiles, totalCount, nil
}

// GetRecommendations gets people recommendations
func (r *PostgresRepository) GetRecommendations(ctx context.Context, userID string) ([]*model.UserProfile, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "repository.GetRecommendations")
	defer span.End()

	// Simple recommendation: get users who have given feedback (excluding current user)
	query := `
		SELECT DISTINCT u.id, u.email, u.password_hash, u.name, u.email_verified, u.public_bio, u.created_at, u.updated_at
		FROM users u
		JOIN feedback_items f ON u.id = f.author_id
		WHERE u.id != $1
		ORDER BY u.name ASC
		LIMIT 10
	`

	rows, err := r.db.Pool.Query(ctx, query, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, errors.WrapError(err, "failed to get recommendations")
	}
	defer rows.Close()

	var recommendations []*model.UserProfile
	for rows.Next() {
		var user model.User
		err := rows.Scan(
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
			continue
		}
		recommendations = append(recommendations, user.ToProfile())
	}

	if err = rows.Err(); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, errors.WrapError(err, "failed to iterate recommendations")
	}

	span.SetStatus(codes.Ok, "")
	return recommendations, nil
}

