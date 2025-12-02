package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"ethos/internal/auth/service"
	"ethos/pkg/errors"
)

func setupRefreshRouter(handler *AuthHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/v1/auth/refresh", handler.Refresh)
	return router
}

func TestRefresh_ValidToken(t *testing.T) {
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)

	refreshReq := service.RefreshRequest{
		RefreshToken: "valid-refresh-token",
	}

	expectedResponse := &service.LoginResponse{
		AccessToken:  "new-access-token",
		RefreshToken: "valid-refresh-token",
	}

	mockService.On("RefreshToken", mock.Anything, &refreshReq).Return(expectedResponse, nil)

	router := setupRefreshRouter(handler)
	body, _ := json.Marshal(refreshReq)
	req, _ := http.NewRequest("POST", "/api/v1/auth/refresh", bytes.NewBuffer(body))
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

func TestRefresh_ExpiredToken(t *testing.T) {
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)

	refreshReq := service.RefreshRequest{
		RefreshToken: "expired-refresh-token",
	}

	mockService.On("RefreshToken", mock.Anything, &refreshReq).Return(nil, errors.ErrTokenExpired)

	router := setupRefreshRouter(handler)
	body, _ := json.Marshal(refreshReq)
	req, _ := http.NewRequest("POST", "/api/v1/auth/refresh", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	var errorResponse map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.NoError(t, err)
	assert.Equal(t, "AUTH_TOKEN_EXPIRED", errorResponse["code"])
	mockService.AssertExpectations(t)
}

func TestRefresh_InvalidToken(t *testing.T) {
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)

	refreshReq := service.RefreshRequest{
		RefreshToken: "invalid-refresh-token",
	}

	mockService.On("RefreshToken", mock.Anything, &refreshReq).Return(nil, errors.ErrTokenInvalid)

	router := setupRefreshRouter(handler)
	body, _ := json.Marshal(refreshReq)
	req, _ := http.NewRequest("POST", "/api/v1/auth/refresh", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	var errorResponse map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.NoError(t, err)
	assert.Equal(t, "AUTH_TOKEN_INVALID", errorResponse["code"])
	mockService.AssertExpectations(t)
}

func TestRefresh_MissingToken(t *testing.T) {
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)

	router := setupRefreshRouter(handler)
	req, _ := http.NewRequest("POST", "/api/v1/auth/refresh", bytes.NewBuffer([]byte("{}")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var errorResponse map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.NoError(t, err)
	assert.Equal(t, "VALIDATION_FAILED", errorResponse["code"])
	mockService.AssertNotCalled(t, "RefreshToken")
}

func TestRefresh_InvalidBody(t *testing.T) {
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)

	router := setupRefreshRouter(handler)
	req, _ := http.NewRequest("POST", "/api/v1/auth/refresh", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertNotCalled(t, "RefreshToken")
}

