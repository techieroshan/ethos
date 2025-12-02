package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"ethos/internal/auth/model"
	"ethos/pkg/errors"
)

// MockChecker is a mock email checker
type MockChecker struct {
	mock.Mock
}

func (m *MockChecker) ValidateEmail(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

func TestAuthService_Register_WithEmailValidation(t *testing.T) {
	tests := []struct {
		name          string
		email         string
		checkerValid  bool
		checkerErr    error
		repoErr       error
		wantErr       bool
		wantErrCode   string
		expectChecker bool
	}{
		{
			name:          "successful registration with valid email",
			email:         "user@example.com",
			checkerValid:  true,
			checkerErr:    nil,
			repoErr:       nil,
			wantErr:       false,
			expectChecker: true,
		},
		{
			name:          "registration fails with temporary email",
			email:         "user@tempmail.com",
			checkerValid:  false,
			checkerErr:    errors.NewValidationError("temporary email addresses are not allowed"),
			repoErr:       nil,
			wantErr:       true,
			wantErrCode:   "VALIDATION_FAILED",
			expectChecker: true,
		},
		{
			name:          "registration fails with invalid email",
			email:         "invalid-email",
			checkerValid:  false,
			checkerErr:    errors.NewValidationError("invalid email format"),
			repoErr:       nil,
			wantErr:       true,
			wantErrCode:   "VALIDATION_FAILED",
			expectChecker: true,
		},
		{
			name:          "registration fails when email already exists",
			email:         "existing@example.com",
			checkerValid:  true,
			checkerErr:    nil,
			repoErr:       errors.ErrEmailAlreadyExists,
			wantErr:       true,
			wantErrCode:   "EMAIL_ALREADY_EXISTS",
			expectChecker: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockRepo := new(MockRepository)
			mockChecker := new(MockChecker)

			// Setup checker expectations
			if tt.expectChecker {
				mockChecker.On("ValidateEmail", mock.Anything, tt.email).
					Return(tt.checkerValid, tt.checkerErr).Once()
			}

			// Setup repository expectations (only if checker passes)
			if tt.checkerValid && tt.checkerErr == nil {
				if tt.repoErr == nil {
					mockRepo.On("CreateUser", mock.Anything, mock.AnythingOfType("*model.User")).
						Return(nil).Once()
				} else {
					mockRepo.On("CreateUser", mock.Anything, mock.AnythingOfType("*model.User")).
						Return(tt.repoErr).Once()
				}
			}

			// Create service with checker
			service := &AuthService{
				repo:          mockRepo,
				tokenGenerator: nil, // Not needed for registration test
				emailChecker:  mockChecker,
			}

			// Test registration
			req := &RegisterRequest{
				Email:    tt.email,
				Password: "SecurePassword123!",
				Name:     "Test User",
			}

			profile, err := service.Register(context.Background(), req)

			// Assertions
			if tt.wantErr {
				assert.Error(t, err)
				if tt.wantErrCode != "" {
					if apiErr, ok := err.(*errors.APIError); ok {
						assert.Equal(t, tt.wantErrCode, apiErr.Code)
					}
				}
				assert.Nil(t, profile)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, profile)
			}

			// Verify mocks
			mockChecker.AssertExpectations(t)
			if tt.checkerValid && tt.checkerErr == nil {
				mockRepo.AssertExpectations(t)
			}
		})
	}
}

// MockRepository is a mock repository for testing
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) CreateUser(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockRepository) GetUserByID(ctx context.Context, userID string) (*model.User, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockRepository) SaveRefreshToken(ctx context.Context, userID, tokenHash string, expiresAt int64) error {
	args := m.Called(ctx, userID, tokenHash, expiresAt)
	return args.Error(0)
}

func (m *MockRepository) GetRefreshToken(ctx context.Context, tokenHash string) (string, error) {
	args := m.Called(ctx, tokenHash)
	return args.String(0), args.Error(1)
}

func (m *MockRepository) DeleteRefreshToken(ctx context.Context, tokenHash string) error {
	args := m.Called(ctx, tokenHash)
	return args.Error(0)
}

// MockTokenGenerator is a mock token generator for testing
type MockTokenGenerator struct {
	mock.Mock
}

func (m *MockTokenGenerator) GenerateAccessToken(userID string) (string, error) {
	args := m.Called(userID)
	return args.String(0), args.Error(1)
}

func (m *MockTokenGenerator) GenerateRefreshToken(userID string) (string, error) {
	args := m.Called(userID)
	return args.String(0), args.Error(1)
}

func (m *MockTokenGenerator) ValidateAccessToken(token string) (string, error) {
	args := m.Called(token)
	return args.String(0), args.Error(1)
}

func (m *MockTokenGenerator) ValidateRefreshToken(token string) (string, error) {
	args := m.Called(token)
	return args.String(0), args.Error(1)
}

