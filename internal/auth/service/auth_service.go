package service

import (
	"context"

	"ethos/internal/auth/model"
)

// LoginRequest represents a login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// RegisterRequest represents a registration request
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required"`
}

// RefreshRequest represents a token refresh request
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// Service defines the interface for authentication business logic
type Service interface {
	// Login authenticates a user and returns tokens
	Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error)

	// Register creates a new user account
	Register(ctx context.Context, req *RegisterRequest) (*model.UserProfile, error)

	// RefreshToken generates a new access token from a refresh token
	RefreshToken(ctx context.Context, req *RefreshRequest) (*LoginResponse, error)

	// GetUserProfile retrieves a user profile by ID
	GetUserProfile(ctx context.Context, userID string) (*model.UserProfile, error)
}

