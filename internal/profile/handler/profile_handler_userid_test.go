package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	authModel "ethos/internal/auth/model"
	"ethos/pkg/errors"
	"ethos/pkg/jwt"
)

func setupProfileRouterForUserID(handler *ProfileHandler, tokenGen *jwt.TokenGenerator) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	v1 := router.Group("/api/v1")
	profile := v1.Group("/profile")
	profile.GET("/:user_id", handler.GetUserProfileByID)
	return router
}

func TestGetUserProfileByID_Success(t *testing.T) {
	mockService := new(MockProfileService)
	handler := NewProfileHandler(mockService)
	tokenGen := jwt.NewTokenGenerator("test-secret", "test-refresh-secret", 15, 336)

	profile := &authModel.UserProfile{
		ID:            "target-user-id",
		Email:         "target@example.com",
		Name:          "Target User",
		EmailVerified: true,
		PublicBio:     "Test bio",
	}

	mockService.On("GetUserProfile", mock.Anything, "target-user-id").Return(profile, nil)

	router := setupProfileRouterForUserID(handler, tokenGen)
	req, _ := http.NewRequest("GET", "/api/v1/profile/target-user-id", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response authModel.UserProfile
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "target-user-id", response.ID)
	assert.Equal(t, "Target User", response.Name)
	mockService.AssertExpectations(t)
}

func TestGetUserProfileByID_NotFound(t *testing.T) {
	mockService := new(MockProfileService)
	handler := NewProfileHandler(mockService)
	tokenGen := jwt.NewTokenGenerator("test-secret", "test-refresh-secret", 15, 336)

	mockService.On("GetUserProfile", mock.Anything, "nonexistent-id").Return(nil, errors.ErrProfileNotFound)

	router := setupProfileRouterForUserID(handler, tokenGen)
	req, _ := http.NewRequest("GET", "/api/v1/profile/nonexistent-id", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "PROFILE_NOT_FOUND", response["code"])
	mockService.AssertExpectations(t)
}

func TestGetUserProfileByID_MissingUserID(t *testing.T) {
	mockService := new(MockProfileService)
	handler := NewProfileHandler(mockService)
	tokenGen := jwt.NewTokenGenerator("test-secret", "test-refresh-secret", 15, 336)

	router := setupProfileRouterForUserID(handler, tokenGen)
	req, _ := http.NewRequest("GET", "/api/v1/profile/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Should return 404 for missing route or 400 for empty user_id
	assert.True(t, w.Code == http.StatusNotFound || w.Code == http.StatusBadRequest)
}

