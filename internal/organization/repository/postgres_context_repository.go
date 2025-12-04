package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"ethos/internal/database"
	"ethos/pkg/errors"

	"github.com/google/uuid"
)

// PostgresContextRepository implements the ContextRepository interface using PostgreSQL
type PostgresContextRepository struct {
	db *database.DB
}

// NewPostgresContextRepository creates a new PostgreSQL context repository
func NewPostgresContextRepository(db *database.DB) ContextRepository {
	return &PostgresContextRepository{db: db}
}

// GetUserOrganizations retrieves all organizations a user belongs to
func (r *PostgresContextRepository) GetUserOrganizations(ctx context.Context, userID string) ([]*UserContext, error) {
	query := `
		SELECT 
			om.user_id,
			om.organization_id,
			o.name,
			om.role,
			om.permissions,
			om.joined_at,
			COALESCE(u.last_login_at, om.joined_at) as last_switched_at
		FROM organization_members om
		JOIN organizations o ON om.organization_id = o.id
		JOIN users u ON u.id = om.user_id
		WHERE om.user_id = $1 AND o.deleted_at IS NULL
		ORDER BY om.joined_at DESC
	`

	rows, err := r.db.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contexts []*UserContext
	for rows.Next() {
		var ctx UserContext
		var permsJSON []byte

		if err := rows.Scan(
			&ctx.UserID,
			&ctx.OrganizationID,
			&ctx.OrganizationName,
			&ctx.Role,
			&permsJSON,
			&ctx.JoinedAt,
			&ctx.LastSwitchedAt,
		); err != nil {
			return nil, err
		}

		// Parse permissions
		if permsJSON != nil {
			var perms []string
			if err := json.Unmarshal(permsJSON, &perms); err != nil {
				ctx.Permissions = []string{}
			} else {
				ctx.Permissions = perms
			}
		}

		contexts = append(contexts, &ctx)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return contexts, nil
}

// GetUserCurrentOrganization retrieves the user's current organization
func (r *PostgresContextRepository) GetUserCurrentOrganization(ctx context.Context, userID string) (*UserContext, error) {
	query := `
		SELECT 
			u.id,
			u.current_organization_id,
			o.name,
			om.role,
			om.permissions,
			om.joined_at,
			COALESCE(u.last_login_at, om.joined_at) as last_switched_at
		FROM users u
		LEFT JOIN organizations o ON u.current_organization_id = o.id
		LEFT JOIN organization_members om ON om.user_id = u.id AND om.organization_id = o.id
		WHERE u.id = $1
	`

	var userCtx UserContext
	var permsJSON *[]byte
	var currentOrgID *string

	err := r.db.Pool.QueryRow(ctx, query, userID).Scan(
		&userCtx.UserID,
		&currentOrgID,
		&userCtx.OrganizationName,
		&userCtx.Role,
		&permsJSON,
		&userCtx.JoinedAt,
		&userCtx.LastSwitchedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	if currentOrgID != nil {
		userCtx.OrganizationID = *currentOrgID
	}

	// Parse permissions
	if permsJSON != nil {
		var perms []string
		if err := json.Unmarshal(*permsJSON, &perms); err != nil {
			userCtx.Permissions = []string{}
		} else {
			userCtx.Permissions = perms
		}
	}

	return &userCtx, nil
}

// UpdateUserCurrentOrganization updates the user's current organization
func (r *PostgresContextRepository) UpdateUserCurrentOrganization(ctx context.Context, userID, organizationID string) error {
	query := `
		UPDATE users 
		SET current_organization_id = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`

	result, err := r.db.Pool.Exec(ctx, query, organizationID, userID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.ErrUserNotFound
	}

	return nil
}

// CreateUserSession creates a new user session
func (r *PostgresContextRepository) CreateUserSession(ctx context.Context, userID, organizationID, tokenHash, refreshTokenHash, ipAddress, userAgent, deviceName string, expiresAt time.Time) (*UserSession, error) {
	sessionID := uuid.New().String()

	query := `
		INSERT INTO user_sessions 
		(id, user_id, organization_id, token_hash, refresh_token_hash, ip_address, user_agent, device_name, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, user_id, organization_id, token_hash, refresh_token_hash, ip_address, user_agent, device_name, last_activity_at, expires_at, revoked_at, created_at
	`

	var session UserSession
	err := r.db.Pool.QueryRow(ctx, query,
		sessionID, userID, organizationID, tokenHash, refreshTokenHash, ipAddress, userAgent, deviceName, expiresAt,
	).Scan(
		&session.ID,
		&session.UserID,
		&session.OrganizationID,
		&session.TokenHash,
		&session.RefreshTokenHash,
		&session.IPAddress,
		&session.UserAgent,
		&session.DeviceName,
		&session.LastActivityAt,
		&session.ExpiresAt,
		&session.RevokedAt,
		&session.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &session, nil
}

// GetUserSessionByToken retrieves a user session by token hash
func (r *PostgresContextRepository) GetUserSessionByToken(ctx context.Context, tokenHash string) (*UserSession, error) {
	query := `
		SELECT id, user_id, organization_id, token_hash, refresh_token_hash, ip_address, user_agent, device_name, last_activity_at, expires_at, revoked_at, created_at
		FROM user_sessions
		WHERE token_hash = $1 AND revoked_at IS NULL AND expires_at > CURRENT_TIMESTAMP
	`

	var session UserSession
	err := r.db.Pool.QueryRow(ctx, query, tokenHash).Scan(
		&session.ID,
		&session.UserID,
		&session.OrganizationID,
		&session.TokenHash,
		&session.RefreshTokenHash,
		&session.IPAddress,
		&session.UserAgent,
		&session.DeviceName,
		&session.LastActivityAt,
		&session.ExpiresAt,
		&session.RevokedAt,
		&session.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &session, nil
}

// GetUserSessionByID retrieves a user session by ID
func (r *PostgresContextRepository) GetUserSessionByID(ctx context.Context, sessionID string) (*UserSession, error) {
	query := `
		SELECT id, user_id, organization_id, token_hash, refresh_token_hash, ip_address, user_agent, device_name, last_activity_at, expires_at, revoked_at, created_at
		FROM user_sessions
		WHERE id = $1
	`

	var session UserSession
	err := r.db.Pool.QueryRow(ctx, query, sessionID).Scan(
		&session.ID,
		&session.UserID,
		&session.OrganizationID,
		&session.TokenHash,
		&session.RefreshTokenHash,
		&session.IPAddress,
		&session.UserAgent,
		&session.DeviceName,
		&session.LastActivityAt,
		&session.ExpiresAt,
		&session.RevokedAt,
		&session.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &session, nil
}

// RevokeUserSession revokes a user session
func (r *PostgresContextRepository) RevokeUserSession(ctx context.Context, sessionID string) error {
	query := `
		UPDATE user_sessions
		SET revoked_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	result, err := r.db.Pool.Exec(ctx, query, sessionID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.ErrNotFound
	}

	return nil
}

// RevokeAllUserSessions revokes all sessions for a user
func (r *PostgresContextRepository) RevokeAllUserSessions(ctx context.Context, userID string) error {
	query := `
		UPDATE user_sessions
		SET revoked_at = CURRENT_TIMESTAMP
		WHERE user_id = $1 AND revoked_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query, userID)
	return err
}

// CleanupExpiredSessions deletes expired sessions
func (r *PostgresContextRepository) CleanupExpiredSessions(ctx context.Context, beforeTime time.Time) (int, error) {
	query := `
		DELETE FROM user_sessions
		WHERE expires_at < $1 OR (revoked_at IS NOT NULL AND revoked_at < $1)
	`

	result, err := r.db.Pool.Exec(ctx, query, beforeTime)
	if err != nil {
		return 0, err
	}

	return int(result.RowsAffected()), nil
}

// RecordContextSwitch records a context switch
func (r *PostgresContextRepository) RecordContextSwitch(ctx context.Context, userID, fromOrgID, toOrgID, sessionID, ipAddress string) (*ContextSwitchRecord, error) {
	recordID := uuid.New().String()
	var fromOrgIDPtr *string
	if fromOrgID != "" {
		fromOrgIDPtr = &fromOrgID
	}

	query := `
		INSERT INTO user_context_switches 
		(id, user_id, from_organization_id, to_organization_id, session_id, ip_address, timestamp)
		VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP)
		RETURNING id, user_id, from_organization_id, to_organization_id, session_id, ip_address, timestamp
	`

	var record ContextSwitchRecord
	err := r.db.Pool.QueryRow(ctx, query, recordID, userID, fromOrgIDPtr, toOrgID, sessionID, ipAddress).Scan(
		&record.ID,
		&record.UserID,
		&record.FromOrganizationID,
		&record.ToOrganizationID,
		&record.SessionID,
		&record.IPAddress,
		&record.Timestamp,
	)

	if err != nil {
		return nil, err
	}

	return &record, nil
}

// GetContextSwitchHistory retrieves context switch history
func (r *PostgresContextRepository) GetContextSwitchHistory(ctx context.Context, userID string, limit, offset int) ([]*ContextSwitchRecord, int64, error) {
	// Get total count
	countQuery := `SELECT COUNT(*) FROM user_context_switches WHERE user_id = $1`
	var total int64
	if err := r.db.Pool.QueryRow(ctx, countQuery, userID).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Get paginated records
	query := `
		SELECT id, user_id, from_organization_id, to_organization_id, session_id, ip_address, timestamp
		FROM user_context_switches
		WHERE user_id = $1
		ORDER BY timestamp DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Pool.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var records []*ContextSwitchRecord
	for rows.Next() {
		var record ContextSwitchRecord
		if err := rows.Scan(
			&record.ID,
			&record.UserID,
			&record.FromOrganizationID,
			&record.ToOrganizationID,
			&record.SessionID,
			&record.IPAddress,
			&record.Timestamp,
		); err != nil {
			return nil, 0, err
		}
		records = append(records, &record)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

// GetUserSessionsByOrganization retrieves sessions for a user in an organization
func (r *PostgresContextRepository) GetUserSessionsByOrganization(ctx context.Context, userID, organizationID string) ([]*UserSession, error) {
	query := `
		SELECT id, user_id, organization_id, token_hash, refresh_token_hash, ip_address, user_agent, device_name, last_activity_at, expires_at, revoked_at, created_at
		FROM user_sessions
		WHERE user_id = $1 AND organization_id = $2 AND revoked_at IS NULL
		ORDER BY last_activity_at DESC
	`

	rows, err := r.db.Pool.Query(ctx, query, userID, organizationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []*UserSession
	for rows.Next() {
		var session UserSession
		if err := rows.Scan(
			&session.ID,
			&session.UserID,
			&session.OrganizationID,
			&session.TokenHash,
			&session.RefreshTokenHash,
			&session.IPAddress,
			&session.UserAgent,
			&session.DeviceName,
			&session.LastActivityAt,
			&session.ExpiresAt,
			&session.RevokedAt,
			&session.CreatedAt,
		); err != nil {
			return nil, err
		}
		sessions = append(sessions, &session)
	}

	return sessions, rows.Err()
}

// IsUserInOrganization checks if user is a member of an organization
func (r *PostgresContextRepository) IsUserInOrganization(ctx context.Context, userID, organizationID string) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM organization_members
			WHERE user_id = $1 AND organization_id = $2
		)
	`

	var exists bool
	err := r.db.Pool.QueryRow(ctx, query, userID, organizationID).Scan(&exists)
	return exists, err
}

// GetUserRoleInOrganization gets the user's role in an organization
func (r *PostgresContextRepository) GetUserRoleInOrganization(ctx context.Context, userID, organizationID string) (string, error) {
	query := `
		SELECT role FROM organization_members
		WHERE user_id = $1 AND organization_id = $2
	`

	var role string
	err := r.db.Pool.QueryRow(ctx, query, userID, organizationID).Scan(&role)
	if err == sql.ErrNoRows {
		return "", errors.ErrNotFound
	}
	if err != nil {
		return "", err
	}

	return role, nil
}

// LogOrganizationActivity logs an organization activity
func (r *PostgresContextRepository) LogOrganizationActivity(ctx context.Context, organizationID, userID, action, resourceType, resourceID, ipAddress, userAgent string, changes map[string]interface{}) error {
	changesJSON, err := json.Marshal(changes)
	if err != nil {
		changesJSON = []byte("{}")
	}

	query := `
		INSERT INTO organization_activity_log 
		(id, organization_id, user_id, action, resource_type, resource_id, changes, ip_address, user_agent)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err = r.db.Pool.Exec(ctx, query,
		uuid.New().String(),
		organizationID,
		userID,
		action,
		resourceType,
		resourceID,
		changesJSON,
		ipAddress,
		userAgent,
	)

	return err
}

// GetOrganizationActivity retrieves organization activity log
func (r *PostgresContextRepository) GetOrganizationActivity(ctx context.Context, organizationID string, limit, offset int) (interface{}, int64, error) {
	// Get total count
	countQuery := `SELECT COUNT(*) FROM organization_activity_log WHERE organization_id = $1`
	var total int64
	if err := r.db.Pool.QueryRow(ctx, countQuery, organizationID).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Get paginated activity
	query := `
		SELECT id, organization_id, user_id, action, resource_type, resource_id, changes, ip_address, user_agent, created_at
		FROM organization_activity_log
		WHERE organization_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Pool.Query(ctx, query, organizationID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var activities []interface{}
	for rows.Next() {
		var id, orgID, userID, action, resourceType, resourceID, ipAddress, userAgent string
		var changesJSON []byte
		var createdAt time.Time

		if err := rows.Scan(&id, &orgID, &userID, &action, &resourceType, &resourceID, &changesJSON, &ipAddress, &userAgent, &createdAt); err != nil {
			return nil, 0, err
		}

		activities = append(activities, map[string]interface{}{
			"id":              id,
			"organization_id": orgID,
			"user_id":         userID,
			"action":          action,
			"resource_type":   resourceType,
			"resource_id":     resourceID,
			"changes":         json.RawMessage(changesJSON),
			"ip_address":      ipAddress,
			"user_agent":      userAgent,
			"created_at":      createdAt,
		})
	}

	return activities, total, rows.Err()
}
