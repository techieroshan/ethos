package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	authModel "ethos/internal/auth/model"
	"ethos/internal/middleware"
	"ethos/internal/profile"
	profileModel "ethos/internal/profile/model"
	"ethos/pkg/jwt"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
) // MockProfileServiceForPrivacy is a mock implementation for privacy tests
type MockProfileServiceForPrivacy struct {
	mock.Mock
}

func (m *MockProfileServiceForPrivacy) GetProfile(ctx context.Context, userID string) (*authModel.UserProfile, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*authModel.UserProfile), args.Error(1)
}

func (m *MockProfileServiceForPrivacy) GetUserProfile(ctx context.Context, userID string) (*authModel.UserProfile, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*authModel.UserProfile), args.Error(1)
}

func (m *MockProfileServiceForPrivacy) UpdateProfile(ctx context.Context, userID string, req *profile.UpdateProfileRequest) (*authModel.UserProfile, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*authModel.UserProfile), args.Error(1)
}

func (m *MockProfileServiceForPrivacy) UpdatePreferences(ctx context.Context, userID string, req *profile.UpdatePreferencesRequest) (*profileModel.UserPreferences, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*profileModel.UserPreferences), args.Error(1)
}

func (m *MockProfileServiceForPrivacy) OptOut(ctx context.Context, userID string, req *profile.OptOutRequest) error {
	args := m.Called(ctx, userID, req)
	return args.Error(0)
}

func (m *MockProfileServiceForPrivacy) Anonymize(ctx context.Context, userID string) (*profile.AnonymizeResponse, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*profile.AnonymizeResponse), args.Error(1)
}

func (m *MockProfileServiceForPrivacy) RequestDeletion(ctx context.Context, userID string, req *profile.DeleteRequest) (*profile.DeleteResponse, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*profile.DeleteResponse), args.Error(1)
}

func (m *MockProfileServiceForPrivacy) DeleteProfile(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockProfileServiceForPrivacy) SearchProfiles(ctx context.Context, query string, limit, offset int) ([]*authModel.UserProfile, int, error) {
	args := m.Called(ctx, query, limit, offset)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*authModel.UserProfile), args.Get(1).(int), args.Error(2)
}

func setupProfileRouterForPrivacy(handler *ProfileHandler, tokenGen *jwt.TokenGenerator) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.AuthMiddleware(tokenGen))
	v1 := router.Group("/api/v1")
	profileGrp := v1.Group("/profile")
	profileGrp.POST("/opt-out", handler.OptOut)
	profileGrp.POST("/anonymize", handler.Anonymize)
	profileGrp.POST("/delete_request", handler.RequestDeletion)
	return router
}

func TestOptOut_Success(t *testing.T) {
	mockService := new(MockProfileServiceForPrivacy)
	handler := NewProfileHandler(mockService)
	tokenGen := jwt.NewTokenGenerator("test-secret", "test-refresh-secret", 15*time.Second, 336*time.Hour)

	request := &profile.OptOutRequest{
		From:   "public_search",
		Reason: "Prefer not to show up in company-wide searches.",
	}

	mockService.On("OptOut", mock.Anything, "user-123", request).Return(nil)

	router := setupProfileRouterForPrivacy(handler, tokenGen)
	req := httptest.NewRequest("POST", "/api/v1/profile/opt-out", nil)
	token, err := tokenGen.GenerateAccessToken("user-123")
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestAnonymize_Success(t *testing.T) {
	mockService := new(MockProfileServiceForPrivacy)
	handler := NewProfileHandler(mockService)
	tokenGen := jwt.NewTokenGenerator("test-secret", "test-refresh-secret", 15*time.Second, 336*time.Hour)

	expectedResponse := &profile.AnonymizeResponse{
		Status:             "processing",
		ExpectedCompletion: time.Now().Add(24 * time.Hour),
	}

	mockService.On("Anonymize", mock.Anything, "user-123").Return(expectedResponse, nil)

	router := setupProfileRouterForPrivacy(handler, tokenGen)
	req := httptest.NewRequest("POST", "/api/v1/profile/anonymize", nil)
	token, err := tokenGen.GenerateAccessToken("user-123")
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestRequestDeletion_Success(t *testing.T) {
	mockService := new(MockProfileServiceForPrivacy)
	handler := NewProfileHandler(mockService)
	tokenGen := jwt.NewTokenGenerator("test-secret", "test-refresh-secret", 15*time.Second, 336*time.Hour)

	request := &profile.DeleteRequest{
		Confirm: true,
		Reason:  "Testing deletion request",
	}

	expectedResponse := &profile.DeleteResponse{
		Status:             "pending",
		ExpectedCompletion: time.Now().Add(30 * 24 * time.Hour),
	}

	mockService.On("RequestDeletion", mock.Anything, "user-123", request).Return(expectedResponse, nil)

	router := setupProfileRouterForPrivacy(handler, tokenGen)
	req := httptest.NewRequest("POST", "/api/v1/profile/delete_request", nil)
	token, err := tokenGen.GenerateAccessToken("user-123")
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}
