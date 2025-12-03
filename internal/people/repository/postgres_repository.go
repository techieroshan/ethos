package repository

import (
	"context"
	"strconv"

	"ethos/internal/auth/model"
	"ethos/internal/database"
	"ethos/internal/people"
	"ethos/pkg/errors"

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

// SearchPeople searches for people
func (r *PostgresRepository) SearchPeople(ctx context.Context, query string, limit, offset int) ([]*model.UserProfile, int, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "repository.SearchPeople")
	defer span.End()

	searchPattern := "%" + query + "%"

	// Get total count
	var totalCount int
	countQuery := `
		SELECT COUNT(*) FROM users
		WHERE (CONCAT(first_name, ' ', last_name) ILIKE $1 OR email ILIKE $1)
	`
	err := r.db.Pool.QueryRow(ctx, countQuery, searchPattern).Scan(&totalCount)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, 0, errors.WrapError(err, "failed to count people")
	}

	// Get profiles
	profilesQuery := `
		SELECT id, email, password_hash, first_name, last_name, email_verified, public_bio, created_at, updated_at
		FROM users
		WHERE (CONCAT(first_name, ' ', last_name) ILIKE $1 OR email ILIKE $1)
		ORDER BY CONCAT(first_name, ' ', last_name) ASC
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
		SELECT DISTINCT u.id, u.email, u.password_hash, u.first_name, u.last_name, u.email_verified, u.public_bio, u.created_at, u.updated_at
		FROM users u
		JOIN feedback_items f ON u.id = f.author_id
		WHERE u.id != $1
		ORDER BY CONCAT(u.first_name, ' ', u.last_name) ASC
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

// SearchPeopleWithFilters searches for people with enhanced filtering
func (r *PostgresRepository) SearchPeopleWithFilters(ctx context.Context, query string, limit, offset int, filters *people.PeopleSearchFilters) ([]*model.UserProfile, int, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "repository.SearchPeopleWithFilters")
	defer span.End()

	searchPattern := "%" + query + "%"

	// Build base query
	querySQL := `
		SELECT id, email, password_hash, first_name, last_name, email_verified, public_bio, created_at, updated_at
		FROM users
		WHERE (CONCAT(first_name, ' ', last_name) ILIKE $1 OR email ILIKE $1)
	`

	countQuery := `
		SELECT COUNT(*)
		FROM users
		WHERE (CONCAT(first_name, ' ', last_name) ILIKE $1 OR email ILIKE $1)
	`

	args := []interface{}{searchPattern}
	argCount := 1

	// Apply filters
	if filters != nil {
		// Verification filter
		if filters.Verification != nil {
			if *filters.Verification == "verified" {
				querySQL += ` AND email_verified = true`
				countQuery += ` AND email_verified = true`
			} else if *filters.Verification == "unverified" {
				querySQL += ` AND email_verified = false`
				countQuery += ` AND email_verified = false`
			}
		}

		// For reviewer_type and context filters, we would need to join with feedback data
		// This is a simplified implementation - in a real system, you'd have user roles/attributes
		if filters.ReviewerType != nil {
			if *filters.ReviewerType == "org" {
				// Users who have given feedback (simplified org reviewer logic)
				querySQL += ` AND EXISTS (SELECT 1 FROM feedback_items WHERE author_id = users.id)`
				countQuery += ` AND EXISTS (SELECT 1 FROM feedback_items WHERE author_id = users.id)`
			}
			// For "public" reviewer type, no additional filter needed (default)
		}

		// Context filter would require additional user attributes or role-based filtering
		// This is simplified - in practice you'd have departments, teams, etc.
		if filters.Context != nil {
			// Placeholder - would filter by user context/attributes
			// span.SetAttributes(attribute.String("filter.context", *filters.Context))
		}

		// Tags filter would require user tagging system
		if len(filters.Tags) > 0 {
			// Placeholder - would filter by user tags
			// span.SetAttributes(attribute.String("filter.tags", strings.Join(filters.Tags, ",")))
		}
	}

	// Add ordering and pagination
	querySQL += ` ORDER BY CONCAT(first_name, ' ', last_name) ASC LIMIT $` + strconv.Itoa(argCount+1) + ` OFFSET $` + strconv.Itoa(argCount+2)
	args = append(args, limit, offset)

	// Get total count
	var totalCount int
	err := r.db.Pool.QueryRow(ctx, countQuery, args[:argCount]...).Scan(&totalCount)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, 0, errors.WrapError(err, "failed to count filtered people")
	}

	// Get profiles
	rows, err := r.db.Pool.Query(ctx, querySQL, args...)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, 0, errors.WrapError(err, "failed to search filtered people")
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
			span.SetStatus(codes.Error, err.Error())
			return nil, 0, errors.WrapError(err, "failed to scan user")
		}
		profiles = append(profiles, user.ToProfile())
	}

	if err = rows.Err(); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, 0, errors.WrapError(err, "failed to iterate filtered people")
	}

	span.SetStatus(codes.Ok, "")
	return profiles, totalCount, nil
}
