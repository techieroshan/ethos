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
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=8"`
	FirstName   string `json:"first_name" binding:"required,min=1"`
	LastName    string `json:"last_name" binding:"required,min=1"`
	AcceptTerms bool   `json:"accept_terms" binding:"required"`
}

// RefreshRequest represents a token refresh request
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// ChangePasswordRequest represents a password change request
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required,min=8"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
}

// RequestPasswordResetRequest represents a password reset request
type RequestPasswordResetRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// Setup2FARequest represents a 2FA setup request
type Setup2FARequest struct {
	Password string `json:"password" binding:"required,min=8"`
}

// Setup2FAResponse represents the response after 2FA setup
type Setup2FAResponse struct {
	Secret  string `json:"secret"`
	QRCode  string `json:"qr_code"`
	Message string `json:"message"`
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

	// VerifyEmail marks a user's email as verified
	VerifyEmail(ctx context.Context, token string) error

	// RequestPasswordReset initiates a password reset process
	RequestPasswordReset(ctx context.Context, req *RequestPasswordResetRequest) error

	// ChangePassword changes user's password
	ChangePassword(ctx context.Context, userID string, req *ChangePasswordRequest) error

	// Setup2FA initializes 2FA for a user
	Setup2FA(ctx context.Context, userID string, req *Setup2FARequest) (*Setup2FAResponse, error)

	// MULTI-TENANT METHODS

	// GetUserByID gets a user by ID with tenant memberships loaded
	GetUserByID(ctx context.Context, userID string) (*model.User, error)

	// SwitchUserTenant switches a user's current tenant context
	SwitchUserTenant(ctx context.Context, userID, tenantID string) error
}
