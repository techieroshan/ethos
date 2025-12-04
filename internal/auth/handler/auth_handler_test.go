package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"ethos/internal/auth/model"
	"ethos/internal/auth/service"
	"ethos/pkg/errors"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAuthService is a mock implementation of the auth service
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Login(ctx context.Context, req *service.LoginRequest) (*service.LoginResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.LoginResponse), args.Error(1)
}

func (m *MockAuthService) Register(ctx context.Context, req *service.RegisterRequest) (*model.UserProfile, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.UserProfile), args.Error(1)
}

func (m *MockAuthService) RefreshToken(ctx context.Context, req *service.RefreshRequest) (*service.LoginResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.LoginResponse), args.Error(1)
}

func (m *MockAuthService) GetUserProfile(ctx context.Context, userID string) (*model.UserProfile, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.UserProfile), args.Error(1)
}

func (m *MockAuthService) ChangePassword(ctx context.Context, userID string, req *service.ChangePasswordRequest) error {
	args := m.Called(ctx, userID, req)
	return args.Error(0)
}

func (m *MockAuthService) VerifyEmail(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockAuthService) RequestPasswordReset(ctx context.Context, req *service.RequestPasswordResetRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *MockAuthService) Setup2FA(ctx context.Context, userID string, req *service.Setup2FARequest) (*service.Setup2FAResponse, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.Setup2FAResponse), args.Error(1)
}

func (m *MockAuthService) GetUserByID(ctx context.Context, userID string) (*model.User, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockAuthService) SwitchUserTenant(ctx context.Context, userID, tenantID string) error {
	args := m.Called(ctx, userID, tenantID)
	return args.Error(0)
}

func setupRouter(handler *AuthHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/v1/auth/login", handler.Login)
	return router
}

func TestLogin_ValidCredentials(t *testing.T) {
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)

	loginReq := service.LoginRequest{
		Email:    "user@example.com",
		Password: "ValidPassword123!",
	}

	expectedResponse := &service.LoginResponse{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
	}

	mockService.On("Login", mock.Anything, &loginReq).Return(expectedResponse, nil)

	router := setupRouter(handler)
	body, _ := json.Marshal(loginReq)
	req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response service.LoginResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse.AccessToken, response.AccessToken)
	assert.Equal(t, expectedResponse.RefreshToken, response.RefreshToken)
	mockService.AssertExpectations(t)
}

func TestLogin_InvalidEmail(t *testing.T) {
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)

	loginReq := service.LoginRequest{
		Email:    "nonexistent@example.com",
		Password: "Password123!",
	}

	mockService.On("Login", mock.Anything, &loginReq).Return(nil, errors.ErrInvalidCredentials)

	router := setupRouter(handler)
	body, _ := json.Marshal(loginReq)
	req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	var errorResponse map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.NoError(t, err)
	assert.Equal(t, "AUTH_INVALID_CREDENTIALS", errorResponse["code"])
	mockService.AssertExpectations(t)
}

func TestLogin_InvalidPassword(t *testing.T) {
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)

	loginReq := service.LoginRequest{
		Email:    "user@example.com",
		Password: "WrongPassword123!",
	}

	mockService.On("Login", mock.Anything, &loginReq).Return(nil, errors.ErrInvalidCredentials)

	router := setupRouter(handler)
	body, _ := json.Marshal(loginReq)
	req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	var errorResponse map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.NoError(t, err)
	assert.Equal(t, "AUTH_INVALID_CREDENTIALS", errorResponse["code"])
	mockService.AssertExpectations(t)
}

func TestLogin_UnverifiedEmail(t *testing.T) {
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)

	loginReq := service.LoginRequest{
		Email:    "user@example.com",
		Password: "ValidPassword123!",
	}

	mockService.On("Login", mock.Anything, &loginReq).Return(nil, errors.ErrEmailUnverified)

	router := setupRouter(handler)
	body, _ := json.Marshal(loginReq)
	req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	var errorResponse map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.NoError(t, err)
	assert.Equal(t, "AUTH_EMAIL_UNVERIFIED", errorResponse["code"])
	mockService.AssertExpectations(t)
}

func TestLogin_MissingBody(t *testing.T) {
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)

	router := setupRouter(handler)
	req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertNotCalled(t, "Login")
}
