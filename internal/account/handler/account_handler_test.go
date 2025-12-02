package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	accountModel "ethos/internal/account/model"
	"ethos/pkg/jwt"
)

// MockAccountService is a mock implementation of the account service
type MockAccountService struct {
	mock.Mock
}

func (m *MockAccountService) GetSecurityEvents(ctx context.Context, userID string, limit, offset int) ([]*accountModel.SecurityEvent, int, error) {
	args := m.Called(ctx, userID, limit, offset)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*accountModel.SecurityEvent), args.Int(1), args.Error(2)
}

func (m *MockAccountService) GetExportStatus(ctx context.Context, userID, exportID string) (*accountModel.DataExport, error) {
	args := m.Called(ctx, userID, exportID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*accountModel.DataExport), args.Error(1)
}

func (m *MockAccountService) Disable2FA(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func setupAccountRouter(handler *AccountHandler, tokenGen *jwt.TokenGenerator) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	v1 := router.Group("/api/v1")
	account := v1.Group("/account")
	account.Use(func(c *gin.Context) {
		c.Set("user_id", "test-user-id")
		c.Next()
	})
	account.GET("/security-events", handler.GetSecurityEvents)
	account.GET("/export-data/:export_id/status", handler.GetExportStatus)
	
	auth := v1.Group("/auth")
	auth.Use(func(c *gin.Context) {
		c.Set("user_id", "test-user-id")
		c.Next()
	})
	auth.DELETE("/setup-2fa", handler.Disable2FA)
	return router
}

func TestGetSecurityEvents_Success(t *testing.T) {
	mockService := new(MockAccountService)
	handler := NewAccountHandler(mockService)
	tokenGen := jwt.NewTokenGenerator("test-secret", "test-refresh-secret", 15, 336)

	events := []*accountModel.SecurityEvent{
		{
			EventID: "event-1",
			Type:    "login",
			IP:      "192.168.1.1",
		},
		{
			EventID: "event-2",
			Type:    "password_change",
			IP:      "192.168.1.1",
		},
	}

	mockService.On("GetSecurityEvents", mock.Anything, "test-user-id", 20, 0).Return(events, 2, nil)

	router := setupAccountRouter(handler, tokenGen)
	req, _ := http.NewRequest("GET", "/api/v1/account/security-events", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(2), response["count"])
	mockService.AssertExpectations(t)
}

func TestGetSecurityEvents_WithPagination(t *testing.T) {
	mockService := new(MockAccountService)
	handler := NewAccountHandler(mockService)
	tokenGen := jwt.NewTokenGenerator("test-secret", "test-refresh-secret", 15, 336)

	mockService.On("GetSecurityEvents", mock.Anything, "test-user-id", 10, 5).Return([]*accountModel.SecurityEvent{}, 0, nil)

	router := setupAccountRouter(handler, tokenGen)
	req, _ := http.NewRequest("GET", "/api/v1/account/security-events?limit=10&offset=5", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestGetSecurityEvents_EmptyList(t *testing.T) {
	mockService := new(MockAccountService)
	handler := NewAccountHandler(mockService)
	tokenGen := jwt.NewTokenGenerator("test-secret", "test-refresh-secret", 15, 336)

	mockService.On("GetSecurityEvents", mock.Anything, "test-user-id", 20, 0).Return([]*accountModel.SecurityEvent{}, 0, nil)

	router := setupAccountRouter(handler, tokenGen)
	req, _ := http.NewRequest("GET", "/api/v1/account/security-events", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["count"])
	mockService.AssertExpectations(t)
}

func TestGetExportStatus_Success(t *testing.T) {
	mockService := new(MockAccountService)
	handler := NewAccountHandler(mockService)
	tokenGen := jwt.NewTokenGenerator("test-secret", "test-refresh-secret", 15, 336)

	export := &accountModel.DataExport{
		ExportID: "export-123",
		Status:   "completed",
	}

	mockService.On("GetExportStatus", mock.Anything, "test-user-id", "export-123").Return(export, nil)

	router := setupAccountRouter(handler, tokenGen)
	req, _ := http.NewRequest("GET", "/api/v1/account/export-data/export-123/status", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response accountModel.DataExport
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "export-123", response.ExportID)
	assert.Equal(t, "completed", response.Status)
	mockService.AssertExpectations(t)
}

func TestGetExportStatus_NotFound(t *testing.T) {
	mockService := new(MockAccountService)
	handler := NewAccountHandler(mockService)
	tokenGen := jwt.NewTokenGenerator("test-secret", "test-refresh-secret", 15, 336)

	mockService.On("GetExportStatus", mock.Anything, "test-user-id", "nonexistent").Return(nil, assert.AnError)

	router := setupAccountRouter(handler, tokenGen)
	req, _ := http.NewRequest("GET", "/api/v1/account/export-data/nonexistent/status", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

func TestDisable2FA_Success(t *testing.T) {
	mockService := new(MockAccountService)
	handler := NewAccountHandler(mockService)
	tokenGen := jwt.NewTokenGenerator("test-secret", "test-refresh-secret", 15, 336)

	mockService.On("Disable2FA", mock.Anything, "test-user-id").Return(nil)

	router := setupAccountRouter(handler, tokenGen)
	req, _ := http.NewRequest("DELETE", "/api/v1/auth/setup-2fa", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["message"], "Two-factor authentication disabled")
	mockService.AssertExpectations(t)
}

func TestDisable2FA_NotEnabled(t *testing.T) {
	mockService := new(MockAccountService)
	handler := NewAccountHandler(mockService)
	tokenGen := jwt.NewTokenGenerator("test-secret", "test-refresh-secret", 15, 336)

	mockService.On("Disable2FA", mock.Anything, "test-user-id").Return(assert.AnError)

	router := setupAccountRouter(handler, tokenGen)
	req, _ := http.NewRequest("DELETE", "/api/v1/auth/setup-2fa", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

