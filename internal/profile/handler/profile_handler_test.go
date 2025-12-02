package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"ethos/internal/auth/model"
	"ethos/internal/middleware"
	prefModel "ethos/internal/profile/model"
	"ethos/internal/profile/service"
	"ethos/pkg/errors"
	"ethos/pkg/jwt"
)

// MockProfileService is a mock implementation of the profile service
type MockProfileService struct {
	mock.Mock
}

func (m *MockProfileService) GetProfile(ctx context.Context, userID string) (*model.UserProfile, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.UserProfile), args.Error(1)
}

func (m *MockProfileService) UpdateProfile(ctx context.Context, userID string, req *service.UpdateProfileRequest) (*model.UserProfile, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.UserProfile), args.Error(1)
}

func (m *MockProfileService) UpdatePreferences(ctx context.Context, userID string, req *service.UpdatePreferencesRequest) (*prefModel.UserPreferences, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*prefModel.UserPreferences), args.Error(1)
}

func (m *MockProfileService) DeleteProfile(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockProfileService) SearchProfiles(ctx context.Context, query string, limit, offset int) ([]*model.UserProfile, int, error) {
	args := m.Called(ctx, query, limit, offset)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*model.UserProfile), args.Get(1).(int), args.Error(2)
}

func (m *MockProfileService) GetUserProfile(ctx context.Context, userID string) (*model.UserProfile, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.UserProfile), args.Error(1)
}

func (m *MockProfileService) OptOut(ctx context.Context, userID string, req *service.OptOutRequest) error {
	args := m.Called(ctx, userID, req)
	return args.Error(0)
}

func (m *MockProfileService) Anonymize(ctx context.Context, userID string) (*service.AnonymizeResponse, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.AnonymizeResponse), args.Error(1)
}

func (m *MockProfileService) RequestDeletion(ctx context.Context, userID string, req *service.DeleteRequest) (*service.DeleteResponse, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.DeleteResponse), args.Error(1)
}

func setupProfileRouter(handler *ProfileHandler, tokenGen *jwt.TokenGenerator) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.AuthMiddleware(tokenGen))
	router.GET("/api/v1/profile/me", handler.GetProfile)
	return router
}

func TestGetProfile_ValidToken(t *testing.T) {
	mockService := new(MockProfileService)
	handler := NewProfileHandler(mockService)

	tokenGen := jwt.NewTokenGenerator(
		"test-access-secret",
		"test-refresh-secret",
		15*60*time.Second,
		14*24*60*60*time.Second,
	)

	userID := "user-123"
	token, err := tokenGen.GenerateAccessToken(userID)
	assert.NoError(t, err)

	expectedProfile := &model.UserProfile{
		ID:            userID,
		Email:         "user@example.com",
		Name:          "Test User",
		EmailVerified: true,
		PublicBio:     "Test bio",
		CreatedAt:     time.Now(),
	}

	mockService.On("GetProfile", mock.Anything, userID).Return(expectedProfile, nil)

	router := setupProfileRouter(handler, tokenGen)
	req, _ := http.NewRequest("GET", "/api/v1/profile/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response model.UserProfile
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedProfile.ID, response.ID)
	assert.Equal(t, expectedProfile.Email, response.Email)
	mockService.AssertExpectations(t)
}

func TestGetProfile_ExpiredToken(t *testing.T) {
	mockService := new(MockProfileService)
	handler := NewProfileHandler(mockService)

	tokenGen := jwt.NewTokenGenerator(
		"test-access-secret",
		"test-refresh-secret",
		15*60*time.Second,
		14*24*60*60*time.Second,
	)

	expiredTokenGen := jwt.NewTokenGenerator(
		"test-access-secret",
		"test-refresh-secret",
		-1*time.Second,
		14*24*60*60*time.Second,
	)

	userID := "user-123"
	token, err := expiredTokenGen.GenerateAccessToken(userID)
	assert.NoError(t, err)

	router := setupProfileRouter(handler, tokenGen)
	req, _ := http.NewRequest("GET", "/api/v1/profile/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	mockService.AssertNotCalled(t, "GetProfile")
}

func TestGetProfile_ProfileNotFound(t *testing.T) {
	mockService := new(MockProfileService)
	handler := NewProfileHandler(mockService)

	tokenGen := jwt.NewTokenGenerator(
		"test-access-secret",
		"test-refresh-secret",
		15*60*time.Second,
		14*24*60*60*time.Second,
	)

	userID := "user-123"
	token, err := tokenGen.GenerateAccessToken(userID)
	assert.NoError(t, err)

	mockService.On("GetProfile", mock.Anything, userID).Return(nil, errors.ErrProfileNotFound)

	router := setupProfileRouter(handler, tokenGen)
	req, _ := http.NewRequest("GET", "/api/v1/profile/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	var errorResponse map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.NoError(t, err)
	assert.Equal(t, "PROFILE_NOT_FOUND", errorResponse["code"])
	mockService.AssertExpectations(t)
}

