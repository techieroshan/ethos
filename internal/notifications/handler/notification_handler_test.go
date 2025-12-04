package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	notifModel "ethos/internal/notifications/model"
	notifService "ethos/internal/notifications/service"
	"ethos/pkg/jwt"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockNotificationService is a mock implementation of the notification service
type MockNotificationService struct {
	mock.Mock
}

func (m *MockNotificationService) GetNotifications(ctx context.Context, userID string, limit, offset int) ([]*notifModel.Notification, int, int, error) {
	args := m.Called(ctx, userID, limit, offset)
	if args.Get(0) == nil {
		return nil, 0, 0, args.Error(3)
	}
	return args.Get(0).([]*notifModel.Notification), args.Int(1), args.Int(2), args.Error(3)
}

func (m *MockNotificationService) GetPreferences(ctx context.Context, userID string) (*notifModel.NotificationPreferences, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*notifModel.NotificationPreferences), args.Error(1)
}

func (m *MockNotificationService) UpdatePreferences(ctx context.Context, userID string, req *notifService.UpdatePreferencesRequest) (*notifModel.NotificationPreferences, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*notifModel.NotificationPreferences), args.Error(1)
}

func (m *MockNotificationService) MarkAllAsRead(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockNotificationService) MarkAsRead(ctx context.Context, userID, notificationID string, read bool) error {
	args := m.Called(ctx, userID, notificationID, read)
	return args.Error(0)
}

func setupNotificationRouter(handler *NotificationHandler, _ *jwt.TokenGenerator) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	v1 := router.Group("/api/v1")
	notifications := v1.Group("/notifications")
	notifications.Use(func(c *gin.Context) {
		c.Set("user_id", "test-user-id")
		c.Next()
	})
	notifications.GET("", handler.GetNotifications)
	notifications.GET("/preferences", handler.GetPreferences)
	notifications.PUT("/preferences", handler.UpdatePreferences)
	return router
}

func TestGetNotifications_Success(t *testing.T) {
	mockService := new(MockNotificationService)
	handler := NewNotificationHandler(mockService)
	tokenGen := jwt.NewTokenGenerator("test-secret", "test-refresh-secret", 15, 336)

	notifications := []*notifModel.Notification{
		{
			NotificationID: "notif-1",
			Type:           notifModel.NotificationTypeFeedbackReply,
			Message:        "You received a reply",
			Read:           false,
		},
		{
			NotificationID: "notif-2",
			Type:           notifModel.NotificationTypeNewComment,
			Message:        "New comment on your feedback",
			Read:           true,
		},
	}

	mockService.On("GetNotifications", mock.Anything, "test-user-id", 20, 0).Return(notifications, 2, 1, nil)

	router := setupNotificationRouter(handler, tokenGen)
	req, _ := http.NewRequest("GET", "/api/v1/notifications", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(2), response["count"])
	assert.Equal(t, float64(1), response["unread_count"])
	mockService.AssertExpectations(t)
}

func TestGetNotifications_WithPagination(t *testing.T) {
	mockService := new(MockNotificationService)
	handler := NewNotificationHandler(mockService)
	tokenGen := jwt.NewTokenGenerator("test-secret", "test-refresh-secret", 15, 336)

	notifications := []*notifModel.Notification{}

	mockService.On("GetNotifications", mock.Anything, "test-user-id", 10, 5).Return(notifications, 0, 0, nil)

	router := setupNotificationRouter(handler, tokenGen)
	req, _ := http.NewRequest("GET", "/api/v1/notifications?limit=10&offset=5", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestGetNotifications_EmptyList(t *testing.T) {
	mockService := new(MockNotificationService)
	handler := NewNotificationHandler(mockService)
	tokenGen := jwt.NewTokenGenerator("test-secret", "test-refresh-secret", 15, 336)

	mockService.On("GetNotifications", mock.Anything, "test-user-id", 20, 0).Return([]*notifModel.Notification{}, 0, 0, nil)

	router := setupNotificationRouter(handler, tokenGen)
	req, _ := http.NewRequest("GET", "/api/v1/notifications", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["count"])
	assert.Equal(t, float64(0), response["unread_count"])
	mockService.AssertExpectations(t)
}

func TestGetPreferences_Success(t *testing.T) {
	mockService := new(MockNotificationService)
	handler := NewNotificationHandler(mockService)
	tokenGen := jwt.NewTokenGenerator("test-secret", "test-refresh-secret", 15, 336)

	prefs := &notifModel.NotificationPreferences{
		Email: true,
		Push:  false,
		InApp: true,
	}

	mockService.On("GetPreferences", mock.Anything, "test-user-id").Return(prefs, nil)

	router := setupNotificationRouter(handler, tokenGen)
	req, _ := http.NewRequest("GET", "/api/v1/notifications/preferences", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotNil(t, response["preferences"])
	mockService.AssertExpectations(t)
}

func TestUpdatePreferences_Success(t *testing.T) {
	mockService := new(MockNotificationService)
	handler := NewNotificationHandler(mockService)
	tokenGen := jwt.NewTokenGenerator("test-secret", "test-refresh-secret", 15, 336)

	prefs := &notifModel.NotificationPreferences{
		Email: true,
		Push:  true,
		InApp: false,
	}

	mockService.On("UpdatePreferences", mock.Anything, "test-user-id", mock.Anything).Return(prefs, nil)

	router := setupNotificationRouter(handler, tokenGen)
	reqBody := `{"email": true, "push": true, "in_app": false}`
	req, _ := http.NewRequest("PUT", "/api/v1/notifications/preferences", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotNil(t, response["preferences"])
	mockService.AssertExpectations(t)
}

func TestUpdatePreferences_InvalidJSON(t *testing.T) {
	mockService := new(MockNotificationService)
	handler := NewNotificationHandler(mockService)
	tokenGen := jwt.NewTokenGenerator("test-secret", "test-refresh-secret", 15, 336)

	router := setupNotificationRouter(handler, tokenGen)
	req, _ := http.NewRequest("PUT", "/api/v1/notifications/preferences", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "VALIDATION_FAILED", response["code"])
}
