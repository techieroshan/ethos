package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode"

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
	// Simple in-memory rate limiting (in production, use Redis)
	loginAttempts  map[string]int
	lockedAccounts map[string]time.Time
}

// NewAuthService creates a new authentication service
func NewAuthService(repo repository.Repository, tokenGen *jwt.TokenGenerator, emailChecker EmailChecker, emailSender EmailSender) Service {
	return &AuthService{
		repo:           repo,
		tokenGenerator: tokenGen,
		emailChecker:   emailChecker,
		emailSender:    emailSender,
		loginAttempts:  make(map[string]int),
		lockedAccounts: make(map[string]time.Time),
	}
}

// Login authenticates a user and returns tokens
func (s *AuthService) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	// Check if account is locked
	if lockTime, locked := s.lockedAccounts[req.Email]; locked {
		if time.Now().Before(lockTime) {
			return nil, errors.NewValidationError("account is temporarily locked due to too many failed login attempts")
		}
		// Lock has expired, remove it
		delete(s.lockedAccounts, req.Email)
		delete(s.loginAttempts, req.Email)
	}

	// Check rate limiting
	if attempts := s.loginAttempts[req.Email]; attempts >= 5 {
		// Lock account for 15 minutes
		s.lockedAccounts[req.Email] = time.Now().Add(15 * time.Minute)
		return nil, errors.NewValidationError("too many failed login attempts - account locked for 15 minutes")
	}

	// Get user by email
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if err == errors.ErrUserNotFound {
			s.loginAttempts[req.Email]++
			return nil, errors.ErrInvalidCredentials
		}
		return nil, errors.WrapError(err, "failed to get user")
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		s.loginAttempts[req.Email]++
		return nil, errors.ErrInvalidCredentials
	}

	// Successful login - reset attempts
	delete(s.loginAttempts, req.Email)

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
	// DEBUG: Log received email at service layer
	fmt.Printf("[DEBUG] Service Register - Received email: '%s' (length: %d, bytes: %v)\n", req.Email, len(req.Email), []byte(req.Email))

	// Combine FirstName and LastName
	fullName := strings.TrimSpace(req.FirstName + " " + req.LastName)

	// Comprehensive input validation
	if fullName == "" {
		return nil, errors.NewValidationError("name is required")
	}
	if len(fullName) < 2 || len(fullName) > 100 {
		return nil, errors.NewValidationError("name must be between 2 and 100 characters")
	}

	// Email validation
	if req.Email == "" {
		return nil, errors.NewValidationError("email is required")
	}
	if !isValidEmailFormat(req.Email) {
		return nil, errors.NewValidationError("invalid email format")
	}

	// Password validation
	if req.Password == "" {
		return nil, errors.NewValidationError("password is required")
	}
	if len(req.Password) < 8 {
		return nil, errors.NewValidationError("password must be at least 8 characters long")
	}
	if !hasPasswordRequirements(req.Password) {
		return nil, errors.NewValidationError("password must contain at least one uppercase letter, one lowercase letter, one number, and one special character")
	}

	// Terms acceptance validation
	if !req.AcceptTerms {
		return nil, errors.NewValidationError("you must accept the terms and conditions")
	}

	// Validate email if checker is provided
	if s.emailChecker != nil {
		valid, err := s.emailChecker.ValidateEmail(ctx, req.Email)
		if err != nil {
			return nil, errors.NewValidationError("email validation service unavailable")
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
		FirstName:     req.FirstName,
		LastName:      req.LastName,
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
				"FirstName": req.FirstName,
				"Name":      user.FirstName + " " + user.LastName,
				"email":     req.Email,
				"user_id":   user.ID,
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

	// Send welcome email after successful verification
	if s.emailSender != nil {
		template := emailTemplates.GetTemplate(emailTemplates.TemplateWelcomeStandardUser)
		emailReq := email.SendEmailRequest{
			To:         user.Email,
			Subject:    template["subject"].(string),
			TemplateID: template["template_id"].(string),
			TemplateData: map[string]interface{}{
				"name":           user.FirstName + " " + user.LastName,
				"email":          user.Email,
				"role":           "Standard User",
				"dashboardURL":   "http://localhost:5173/dashboard",              // TODO: Make configurable
				"unsubscribeURL": "http://localhost:5173/settings/notifications", // TODO: Make configurable
			},
		}

		// Send email asynchronously (don't fail verification if email fails)
		go func() {
			if err := s.emailSender.SendEmail(context.Background(), emailReq); err != nil {
				// Log error but don't fail verification
				fmt.Printf("Failed to send welcome email: %v\n", err)
			}
		}()
	}

	return nil
}

// RequestPasswordReset initiates a password reset process by sending a reset email
func (s *AuthService) RequestPasswordReset(ctx context.Context, req *RequestPasswordResetRequest) error {
	// Get user by email
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if err == errors.ErrUserNotFound {
			// Don't reveal if email exists or not for security
			return nil
		}
		return errors.WrapError(err, "failed to get user")
	}

	// Generate reset token (simplified - in production, use secure token generation)
	resetToken := fmt.Sprintf("reset_%s_%d", user.ID, time.Now().Unix())

	// Send password reset email if sender is configured
	if s.emailSender != nil {
		template := emailTemplates.GetTemplate(emailTemplates.TemplatePasswordReset)
		emailReq := email.SendEmailRequest{
			To:         user.Email,
			Subject:    template["subject"].(string),
			TemplateID: template["template_id"].(string),
			TemplateData: map[string]interface{}{
				"name":           user.FirstName + " " + user.LastName,
				"email":          user.Email,
				"resetURL":       fmt.Sprintf("http://localhost:5173/reset-password?token=%s", resetToken), // TODO: Make configurable
				"unsubscribeURL": "http://localhost:5173/settings/notifications",                           // TODO: Make configurable
			},
		}

		// Send email asynchronously
		go func() {
			if err := s.emailSender.SendEmail(context.Background(), emailReq); err != nil {
				// Log error but don't fail the request
				fmt.Printf("Failed to send password reset email: %v\n", err)
			}
		}()
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

// isValidEmailFormat validates email format using regex
func isValidEmailFormat(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// hasPasswordRequirements checks if password meets security requirements
func hasPasswordRequirements(password string) bool {
	if len(password) < 8 {
		return false
	}

	var hasUpper, hasLower, hasNumber, hasSpecial bool
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasNumber && hasSpecial
}

// MULTI-TENANT METHODS IMPLEMENTATION

// GetUserByID gets a user by ID with tenant memberships loaded
func (s *AuthService) GetUserByID(ctx context.Context, userID string) (*model.User, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Load tenant memberships for the user
	// TODO: Implement tenant membership loading from database
	// For now, return mock data
	now := time.Now()
	user.TenantMemberships = []model.TenantMembership{
		{
			TenantID:   "org-1",
			TenantName: "Acme Corporation",
			Role:       "user",
			JoinedAt:   now,
			IsActive:   true,
		},
		{
			TenantID:   "org-2",
			TenantName: "Tech Startup Inc",
			Role:       "admin",
			JoinedAt:   now,
			IsActive:   true,
		},
	}

	// Set default current tenant if not set
	if user.CurrentTenantID == nil && len(user.TenantMemberships) > 0 {
		user.CurrentTenantID = &user.TenantMemberships[0].TenantID
	}

	return user, nil
}

// SwitchUserTenant switches a user's current tenant context
func (s *AuthService) SwitchUserTenant(ctx context.Context, userID, tenantID string) error {
	user, err := s.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	err = user.SwitchTenant(tenantID)
	if err != nil {
		return err
	}

	// TODO: Persist the tenant switch to database
	// For now, just validate the switch

	return nil
}
