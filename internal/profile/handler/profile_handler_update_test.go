package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"ethos/internal/auth/model"
	"ethos/internal/middleware"
	"ethos/internal/profile"
	"ethos/pkg/jwt"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupUpdateProfileRouter(handler *ProfileHandler, tokenGen *jwt.TokenGenerator) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.AuthMiddleware(tokenGen))
	router.PUT("/api/v1/profile/me", handler.UpdateProfile)
	return router
}

func TestUpdateProfile_ValidUpdate(t *testing.T) {
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

	updateReq := profile.UpdateProfileRequest{
		Name:      "Updated Name",
		PublicBio: "Updated bio",
	}

	expectedProfile := &model.UserProfile{
		ID:            userID,
		Email:         "user@example.com",
		Name:          "Updated Name",
		EmailVerified: true,
		PublicBio:     "Updated bio",
		CreatedAt:     time.Now(),
	}

	mockService.On("UpdateProfile", mock.Anything, userID, &updateReq).Return(expectedProfile, nil)

	router := setupUpdateProfileRouter(handler, tokenGen)
	body, _ := json.Marshal(updateReq)
	req, _ := http.NewRequest("PUT", "/api/v1/profile/me", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response model.UserProfile
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedProfile.Name, response.Name)
	assert.Equal(t, expectedProfile.PublicBio, response.PublicBio)
	mockService.AssertExpectations(t)
}

func TestUpdateProfile_InvalidFields(t *testing.T) {
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

	updateReq := map[string]interface{}{}

	router := setupUpdateProfileRouter(handler, tokenGen)
	body, _ := json.Marshal(updateReq)
	req, _ := http.NewRequest("PUT", "/api/v1/profile/me", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertNotCalled(t, "UpdateProfile")
}

func TestUpdateProfile_ExpiredToken(t *testing.T) {
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

	updateReq := profile.UpdateProfileRequest{
		Name: "Updated Name",
	}

	router := setupUpdateProfileRouter(handler, tokenGen)
	body, _ := json.Marshal(updateReq)
	req, _ := http.NewRequest("PUT", "/api/v1/profile/me", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	mockService.AssertNotCalled(t, "UpdateProfile")
}
