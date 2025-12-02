package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"ethos/internal/middleware"
	prefModel "ethos/internal/profile/model"
	"ethos/internal/profile/service"
	"ethos/pkg/jwt"
)

// UpdatePreferences is already defined in profile_handler_test.go

func setupPreferencesRouter(handler *ProfileHandler, tokenGen *jwt.TokenGenerator) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.AuthMiddleware(tokenGen))
	router.PATCH("/api/v1/profile/me/preferences", handler.UpdatePreferences)
	return router
}

func TestUpdatePreferences_ValidUpdate(t *testing.T) {
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

	notifyOnLogin := false
	locale := "en-GB"
	updateReq := service.UpdatePreferencesRequest{
		NotifyOnLogin: &notifyOnLogin,
		Locale:        &locale,
	}

	expectedPrefs := &prefModel.UserPreferences{
		UserID:        userID,
		NotifyOnLogin: false,
		Locale:        "en-GB",
	}

	mockService.On("UpdatePreferences", mock.Anything, userID, &updateReq).Return(expectedPrefs, nil)

	router := setupPreferencesRouter(handler, tokenGen)
	body, _ := json.Marshal(updateReq)
	req, _ := http.NewRequest("PATCH", "/api/v1/profile/me/preferences", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response struct {
		Preferences prefModel.UserPreferences `json:"preferences"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedPrefs.NotifyOnLogin, response.Preferences.NotifyOnLogin)
	assert.Equal(t, expectedPrefs.Locale, response.Preferences.Locale)
	mockService.AssertExpectations(t)
}

