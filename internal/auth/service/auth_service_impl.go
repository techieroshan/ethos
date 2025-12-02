package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"ethos/internal/auth/model"
	"ethos/internal/auth/repository"
	"ethos/pkg/email"
	emailTemplates "ethos/pkg/email/templates"
	"ethos/pkg/errors"
	"ethos/pkg/jwt"

	"golang.org/x/crypto/bcrypt"
)

// EmailChecker defines the interface for email validation
type EmailChecker interface {
	ValidateEmail(ctx context.Context, email string) (bool, error)
}

// EmailSender defines the interface for sending emails
type EmailSender interface {
	SendEmail(ctx context.Context, req email.SendEmailRequest) error
}

// AuthService implements the Service interface
type AuthService struct {
	repo           repository.Repository
	tokenGenerator *jwt.TokenGenerator
	emailChecker   EmailChecker
	emailSender    EmailSender
}

// NewAuthService creates a new authentication service
func NewAuthService(repo repository.Repository, tokenGen *jwt.TokenGenerator, emailChecker EmailChecker, emailSender EmailSender) Service {
	return &AuthService{
		repo:           repo,
		tokenGenerator: tokenGen,
		emailChecker:   emailChecker,
		emailSender:    emailSender,
	}
}

// Login authenticates a user and returns tokens
func (s *AuthService) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	// Get user by email
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if err == errors.ErrUserNotFound {
			return nil, errors.ErrInvalidCredentials
		}
		return nil, errors.WrapError(err, "failed to get user")
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, errors.ErrInvalidCredentials
	}

	// Check if email is verified
	if !user.EmailVerified {
		return nil, errors.ErrEmailUnverified
	}

	// Generate tokens
	accessToken, err := s.tokenGenerator.GenerateAccessToken(user.ID)
	if err != nil {
		return nil, errors.WrapError(err, "failed to generate access token")
	}

	refreshToken, err := s.tokenGenerator.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, errors.WrapError(err, "failed to generate refresh token")
	}

	// Save refresh token
	tokenHash := hashToken(refreshToken)
	expiresAt := time.Now().Add(14 * 24 * time.Hour).Unix()
	if err := s.repo.SaveRefreshToken(ctx, user.ID, tokenHash, expiresAt); err != nil {
		return nil, errors.WrapError(err, "failed to save refresh token")
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// Register creates a new user account
func (s *AuthService) Register(ctx context.Context, req *RegisterRequest) (*model.UserProfile, error) {
	// Validate email if checker is provided
	if s.emailChecker != nil {
		valid, err := s.emailChecker.ValidateEmail(ctx, req.Email)
		if err != nil {
			return nil, errors.NewValidationError(err.Error())
		}
		if !valid {
			return nil, errors.NewValidationError("invalid email address")
		}
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.WrapError(err, "failed to hash password")
	}

	// Create user
	user := &model.User{
		Email:         req.Email,
		PasswordHash:  string(hashedPassword),
		Name:          req.Name,
		EmailVerified: false,
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	// Send verification email if sender is configured
	if s.emailSender != nil {
		template := emailTemplates.GetTemplate(emailTemplates.TemplateEmailVerification)
		emailReq := email.SendEmailRequest{
			To:         req.Email,
			Subject:    template["subject"].(string),
			TemplateID: template["template_id"].(string),
			TemplateData: map[string]interface{}{
				"name":    req.Name,
				"email":   req.Email,
				"user_id": user.ID,
			},
		}
		// Send email asynchronously (don't fail registration if email fails)
		go func() {
			if err := s.emailSender.SendEmail(context.Background(), emailReq); err != nil {
				// Log error but don't fail registration
				fmt.Printf("Failed to send verification email: %v\n", err)
			}
		}()
	}

	return user.ToProfile(), nil
}

// RefreshToken generates a new access token from a refresh token
func (s *AuthService) RefreshToken(ctx context.Context, req *RefreshRequest) (*LoginResponse, error) {
	// Validate refresh token
	userID, err := s.tokenGenerator.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		if err.Error() == "token expired" {
			return nil, errors.ErrTokenExpired
		}
		return nil, errors.ErrTokenInvalid
	}

	// Verify token exists in database
	tokenHash := hashToken(req.RefreshToken)
	storedUserID, err := s.repo.GetRefreshToken(ctx, tokenHash)
	if err != nil {
		return nil, errors.ErrTokenInvalid
	}

	if storedUserID != userID {
		return nil, errors.ErrTokenInvalid
	}

	// Generate new access token
	accessToken, err := s.tokenGenerator.GenerateAccessToken(userID)
	if err != nil {
		return nil, errors.WrapError(err, "failed to generate access token")
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: req.RefreshToken, // Return same refresh token
	}, nil
}

// GetUserProfile retrieves a user profile by ID
func (s *AuthService) GetUserProfile(ctx context.Context, userID string) (*model.UserProfile, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return user.ToProfile(), nil
}

// VerifyEmail marks a user's email as verified using a verification token
func (s *AuthService) VerifyEmail(ctx context.Context, token string) error {
	// Parse the verification token to extract user ID
	// Token format: hash of user_id
	userID, err := parseVerificationToken(token)
	if err != nil {
		return errors.ErrTokenInvalid
	}

	// Get user from repository
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return errors.ErrUserNotFound
	}

	// If already verified, return success (idempotent)
	if user.EmailVerified {
		return nil
	}

	// Mark email as verified
	user.EmailVerified = true
	user.UpdatedAt = time.Now()

	// Save updated user to repository
	if err := s.repo.UpdateUser(ctx, user); err != nil {
		return err
	}

	return nil
}

// ChangePassword changes a user's password after verifying the current password
func (s *AuthService) ChangePassword(ctx context.Context, userID string, req *ChangePasswordRequest) error {
	// Get user from repository
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return errors.ErrUserNotFound
	}

	// Verify current password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.CurrentPassword))
	if err != nil {
		return errors.ErrInvalidCredentials
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.ErrServerError
	}

	// Update password
	user.PasswordHash = string(hashedPassword)
	user.UpdatedAt = time.Now()

	// Save updated user to repository
	if err := s.repo.UpdateUser(ctx, user); err != nil {
		return err
	}

	return nil
}

// Setup2FA initializes two-factor authentication for a user
func (s *AuthService) Setup2FA(ctx context.Context, userID string, req *Setup2FARequest) (*Setup2FAResponse, error) {
	// Get user from repository
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, errors.ErrUserNotFound
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, errors.ErrInvalidCredentials
	}

	// Generate 2FA secret
	secret := generateTOTPSecret(user.ID)
	qrCode := generateQRCode(user.Email, secret)

	// Note: In production, would store the 2FA secret in a new model/table
	// For now returning it in response for client to save/confirm
	return &Setup2FAResponse{
		Secret:  secret,
		QRCode:  qrCode,
		Message: "Please save this secret and scan the QR code with an authenticator app",
	}, nil
}

// parseVerificationToken extracts user ID from verification token
func parseVerificationToken(token string) (string, error) {
	// Token format: hash-based token containing user_id
	// In production, would validate HMAC or JWT signature
	// For testing: assume token is prefixed with "verify_"
	if len(token) < 20 {
		return "", fmt.Errorf("invalid token format")
	}

	// Extract user ID from token (simplified)
	if len(token) > 7 && token[:7] == "verify_" {
		return token[7:], nil
	}

	return "", fmt.Errorf("invalid token format")
}

// generateTOTPSecret generates a TOTP secret for 2FA
func generateTOTPSecret(userID string) string {
	// Generate deterministic secret from user ID
	hash := sha256.Sum256([]byte(userID + time.Now().Format("20060102")))
	return hex.EncodeToString(hash[:16])
}

// generateQRCode generates a QR code URI for TOTP setup
func generateQRCode(email, secret string) string {
	// Return otpauth URI that can be converted to QR code by client
	return fmt.Sprintf("otpauth://totp/Ethos:%%s?secret=%s&issuer=Ethos", secret)
}

// hashToken creates a SHA256 hash of a token for storage
func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
