package handler

import (
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
	"ethos/pkg/jwt"
)

func setupMeRouter(handler *AuthHandler, tokenGen *jwt.TokenGenerator) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.AuthMiddleware(tokenGen))
	router.GET("/api/v1/auth/me", handler.Me)
	return router
}

func TestMe_ValidToken(t *testing.T) {
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)

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
		CreatedAt:     time.Now(),
	}

	mockService.On("GetUserProfile", mock.Anything, userID).Return(expectedProfile, nil)

	router := setupMeRouter(handler, tokenGen)
	req, _ := http.NewRequest("GET", "/api/v1/auth/me", nil)
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

func TestMe_ExpiredToken(t *testing.T) {
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)

	tokenGen := jwt.NewTokenGenerator(
		"test-access-secret",
		"test-refresh-secret",
		15*60*time.Second,
		14*24*60*60*time.Second,
	)

	// Generate an expired token by using a very short expiry
	expiredTokenGen := jwt.NewTokenGenerator(
		"test-access-secret",
		"test-refresh-secret",
		-1*time.Second, // Already expired
		14*24*60*60*time.Second,
	)

	userID := "user-123"
	token, err := expiredTokenGen.GenerateAccessToken(userID)
	assert.NoError(t, err)

	router := setupMeRouter(handler, tokenGen)
	req, _ := http.NewRequest("GET", "/api/v1/auth/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	var errorResponse map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.NoError(t, err)
	assert.Equal(t, "AUTH_TOKEN_EXPIRED", errorResponse["code"])
	mockService.AssertNotCalled(t, "GetUserProfile")
}

func TestMe_InvalidToken(t *testing.T) {
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)

	tokenGen := jwt.NewTokenGenerator(
		"test-access-secret",
		"test-refresh-secret",
		15*60*time.Second,
		14*24*60*60*time.Second,
	)

	router := setupMeRouter(handler, tokenGen)
	req, _ := http.NewRequest("GET", "/api/v1/auth/me", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	var errorResponse map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.NoError(t, err)
	assert.Equal(t, "AUTH_TOKEN_INVALID", errorResponse["code"])
	mockService.AssertNotCalled(t, "GetUserProfile")
}

func TestMe_MissingHeader(t *testing.T) {
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)

	tokenGen := jwt.NewTokenGenerator(
		"test-access-secret",
		"test-refresh-secret",
		15*60*time.Second,
		14*24*60*60*time.Second,
	)

	router := setupMeRouter(handler, tokenGen)
	req, _ := http.NewRequest("GET", "/api/v1/auth/me", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	var errorResponse map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.NoError(t, err)
	assert.Equal(t, "AUTH_TOKEN_INVALID", errorResponse["code"])
	mockService.AssertNotCalled(t, "GetUserProfile")
}

