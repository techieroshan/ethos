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
	"ethos/internal/auth/model"
	"ethos/internal/auth/service"
	"ethos/pkg/errors"
)

func setupRegisterRouter(handler *AuthHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/v1/auth/register", handler.Register)
	return router
}

func TestRegister_ValidRegistration(t *testing.T) {
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)

	registerReq := service.RegisterRequest{
		Email:    "newuser@example.com",
		Password: "ValidPassword123!",
		Name:     "New User",
	}

	expectedProfile := &model.UserProfile{
		ID:            "user-123",
		Email:         "newuser@example.com",
		Name:          "New User",
		EmailVerified: false,
	}

	mockService.On("Register", mock.Anything, &registerReq).Return(expectedProfile, nil)

	router := setupRegisterRouter(handler)
	body, _ := json.Marshal(registerReq)
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var response model.UserProfile
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedProfile.Email, response.Email)
	assert.Equal(t, expectedProfile.Name, response.Name)
	mockService.AssertExpectations(t)
}

func TestRegister_DuplicateEmail(t *testing.T) {
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)

	registerReq := service.RegisterRequest{
		Email:    "existing@example.com",
		Password: "ValidPassword123!",
		Name:     "Existing User",
	}

	mockService.On("Register", mock.Anything, &registerReq).Return(nil, errors.ErrEmailAlreadyExists)

	router := setupRegisterRouter(handler)
	body, _ := json.Marshal(registerReq)
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
	var errorResponse map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.NoError(t, err)
	assert.Equal(t, "EMAIL_ALREADY_EXISTS", errorResponse["code"])
	mockService.AssertExpectations(t)
}

func TestRegister_InvalidEmailFormat(t *testing.T) {
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)

	registerReq := service.RegisterRequest{
		Email:    "invalid-email",
		Password: "ValidPassword123!",
		Name:     "User",
	}

	router := setupRegisterRouter(handler)
	body, _ := json.Marshal(registerReq)
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var errorResponse map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.NoError(t, err)
	assert.Equal(t, "VALIDATION_FAILED", errorResponse["code"])
	mockService.AssertNotCalled(t, "Register")
}

func TestRegister_WeakPassword(t *testing.T) {
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)

	registerReq := service.RegisterRequest{
		Email:    "user@example.com",
		Password: "short",
		Name:     "User",
	}

	router := setupRegisterRouter(handler)
	body, _ := json.Marshal(registerReq)
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var errorResponse map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.NoError(t, err)
	assert.Equal(t, "VALIDATION_FAILED", errorResponse["code"])
	mockService.AssertNotCalled(t, "Register")
}

func TestRegister_MissingBody(t *testing.T) {
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)

	router := setupRegisterRouter(handler)
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertNotCalled(t, "Register")
}

