package repository

import (
	"context"

	"ethos/internal/auth/model"
)

// Repository defines the interface for authentication data access
type Repository interface {
	// GetUserByEmail retrieves a user by email address
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)

	// GetUserByID retrieves a user by ID
	GetUserByID(ctx context.Context, userID string) (*model.User, error)

	// CreateUser creates a new user
	CreateUser(ctx context.Context, user *model.User) error

	// UpdateUser updates an existing user
	UpdateUser(ctx context.Context, user *model.User) error

	// SaveRefreshToken saves a refresh token
	SaveRefreshToken(ctx context.Context, userID, tokenHash string, expiresAt int64) error

	// GetRefreshToken retrieves a refresh token by hash
	GetRefreshToken(ctx context.Context, tokenHash string) (string, error)

	// DeleteRefreshToken deletes a refresh token
	DeleteRefreshToken(ctx context.Context, tokenHash string) error
}
