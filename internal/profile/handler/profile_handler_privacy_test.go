package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"ethos/internal/middleware"
	"ethos/internal/profile/service"
	"ethos/pkg/jwt"
)

// MockProfileServiceForPrivacy is a mock implementation for privacy tests
type MockProfileServiceForPrivacy struct {
	mock.Mock
}

func (m *MockProfileServiceForPrivacy) OptOut(ctx context.Context, userID string, req *profileService.OptOutRequest) error {
	args := m.Called(ctx, userID, req)
	return args.Error(0)
}

func (m *MockProfileServiceForPrivacy) Anonymize(ctx context.Context, userID string) (*profileService.AnonymizeResponse, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*profileService.AnonymizeResponse), args.Error(1)
}

func (m *MockProfileServiceForPrivacy) RequestDeletion(ctx context.Context, userID string, req *profileService.DeleteRequest) (*profileService.DeleteResponse, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*profileService.DeleteResponse), args.Error(1)
}

func setupProfileRouterForPrivacy(handler *ProfileHandler, tokenGen *jwt.TokenGenerator) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.AuthMiddleware(tokenGen))
	v1 := router.Group("/api/v1")
	profile := v1.Group("/profile")
	profile.POST("/opt-out", handler.OptOut)
	profile.POST("/anonymize", handler.Anonymize)
	profile.POST("/delete_request", handler.RequestDeletion)
	return router
}

func TestOptOut_Success(t *testing.T) {
	mockService := new(MockProfileServiceForPrivacy)
	handler := NewProfileHandler(mockService)
	tokenGen := jwt.NewTokenGenerator("test-secret", "test-refresh-secret", 15, 336)

	request := &profileService.OptOutRequest{
		From:   "public_search",
		Reason: "Prefer not to show up in company-wide searches.",
	}

	mockService.On("OptOut", mock.Anything, "user-123", request).Return(nil)

	router := setupProfileRouterForPrivacy(handler, tokenGen)

	requestBody, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", "/api/v1/profile/opt-out", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tokenGen.GenerateAccessToken("user-123", []string{}))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "opted_out", response["status"])
	assert.True(t, response["changed"].(bool))
	mockService.AssertExpectations(t)
}

func TestOptOut_InvalidRequest(t *testing.T) {
	mockService := new(MockProfileServiceForPrivacy)
	handler := NewProfileHandler(mockService)
	tokenGen := jwt.NewTokenGenerator("test-secret", "test-refresh-secret", 15, 336)

	// Missing required 'from' field
	invalidRequest := map[string]string{
		"reason": "Test reason",
	}

	router := setupProfileRouterForPrivacy(handler, tokenGen)

	requestBody, _ := json.Marshal(invalidRequest)
	req, _ := http.NewRequest("POST", "/api/v1/profile/opt-out", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tokenGen.GenerateAccessToken("user-123", []string{}))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "VALIDATION_FAILED", response["code"])
}

func TestAnonymize_Success(t *testing.T) {
	mockService := new(MockProfileServiceForPrivacy)
	handler := NewProfileHandler(mockService)
	tokenGen := jwt.NewTokenGenerator("test-secret", "test-refresh-secret", 15, 336)

	expectedResponse := &profileService.AnonymizeResponse{
		Status:         "in_progress",
		ExpectedCompletion: time.Date(2024, 12, 15, 12, 0, 0, 0, time.UTC),
	}

	mockService.On("Anonymize", mock.Anything, "user-123").Return(expectedResponse, nil)

	router := setupProfileRouterForPrivacy(handler, tokenGen)

	req, _ := http.NewRequest("POST", "/api/v1/profile/anonymize", nil)
	req.Header.Set("Authorization", "Bearer "+tokenGen.GenerateAccessToken("user-123", []string{}))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response profileService.AnonymizeResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "in_progress", response.Status)
	assert.NotNil(t, response.ExpectedCompletion)
	mockService.AssertExpectations(t)
}

func TestRequestDeletion_Success(t *testing.T) {
	mockService := new(MockProfileServiceForPrivacy)
	handler := NewProfileHandler(mockService)
	tokenGen := jwt.NewTokenGenerator("test-secret", "test-refresh-secret", 15, 336)

	request := &profileService.DeleteRequest{
		Confirm: true,
		Reason:  "Leaving the company",
	}

	expectedResponse := &profileService.DeleteResponse{
		Status:             "delete_requested",
		ExpectedCompletion: time.Date(2024, 12, 20, 18, 0, 0, 0, time.UTC),
	}

	mockService.On("RequestDeletion", mock.Anything, "user-123", request).Return(expectedResponse, nil)

	router := setupProfileRouterForPrivacy(handler, tokenGen)

	requestBody, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", "/api/v1/profile/delete_request", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tokenGen.GenerateAccessToken("user-123", []string{}))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response profileService.DeleteResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "delete_requested", response.Status)
	assert.NotNil(t, response.ExpectedCompletion)
	mockService.AssertExpectations(t)
}

func TestRequestDeletion_MissingConfirmation(t *testing.T) {
	mockService := new(MockProfileServiceForPrivacy)
	handler := NewProfileHandler(mockService)
	tokenGen := jwt.NewTokenGenerator("test-secret", "test-refresh-secret", 15, 336)

	request := &profileService.DeleteRequest{
		Confirm: false, // Not confirmed
		Reason:  "Leaving the company",
	}

	router := setupProfileRouterForPrivacy(handler, tokenGen)

	requestBody, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", "/api/v1/profile/delete_request", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tokenGen.GenerateAccessToken("user-123", []string{}))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "VALIDATION_FAILED", response["code"])
}

func TestPrivacyControls_Unauthorized(t *testing.T) {
	mockService := new(MockProfileServiceForPrivacy)
	handler := NewProfileHandler(mockService)
	tokenGen := jwt.NewTokenGenerator("test-secret", "test-refresh-secret", 15, 336)

	router := setupProfileRouterForPrivacy(handler, tokenGen)
	req, _ := http.NewRequest("POST", "/api/v1/profile/opt-out", bytes.NewBufferString("{}"))
	req.Header.Set("Content-Type", "application/json")
	// No Authorization header
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
