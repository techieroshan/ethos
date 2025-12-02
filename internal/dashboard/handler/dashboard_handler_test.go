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
	dashModel "ethos/internal/dashboard/model"
	fbModel "ethos/internal/feedback/model"
	"ethos/pkg/jwt"
)

// MockDashboardService is a mock implementation of the dashboard service
type MockDashboardService struct {
	mock.Mock
}

func (m *MockDashboardService) GetDashboard(ctx context.Context, userID string) (*dashModel.DashboardSnapshot, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dashModel.DashboardSnapshot), args.Error(1)
}

func setupDashboardRouter(handler *DashboardHandler, tokenGen *jwt.TokenGenerator) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	v1 := router.Group("/api/v1")
	dashboard := v1.Group("/dashboard")
	dashboard.Use(func(c *gin.Context) {
		c.Set("user_id", "test-user-id")
		c.Next()
	})
	dashboard.GET("", handler.GetDashboard)
	return router
}

func TestGetDashboard_Success(t *testing.T) {
	mockService := new(MockDashboardService)
	handler := NewDashboardHandler(mockService)
	tokenGen := jwt.NewTokenGenerator("test-secret", "test-refresh-secret", 15, 336)

	snapshot := &dashModel.DashboardSnapshot{
		RecentFeedback: []*fbModel.FeedbackItem{},
		Stats: map[string]int{
			"feedback_given": 5,
			"comments":       10,
		},
		SuggestedActions: []string{"Give feedback", "Connect with team"},
	}

	mockService.On("GetDashboard", mock.Anything, "test-user-id").Return(snapshot, nil)

	router := setupDashboardRouter(handler, tokenGen)
	req, _ := http.NewRequest("GET", "/api/v1/dashboard", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response dashModel.DashboardSnapshot
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 5, response.Stats["feedback_given"])
	assert.Equal(t, 10, response.Stats["comments"])
	assert.Equal(t, 2, len(response.SuggestedActions))
	mockService.AssertExpectations(t)
}

func TestGetDashboard_ServiceError(t *testing.T) {
	mockService := new(MockDashboardService)
	handler := NewDashboardHandler(mockService)
	tokenGen := jwt.NewTokenGenerator("test-secret", "test-refresh-secret", 15, 336)

	mockService.On("GetDashboard", mock.Anything, "test-user-id").Return(nil, assert.AnError)

	router := setupDashboardRouter(handler, tokenGen)
	req, _ := http.NewRequest("GET", "/api/v1/dashboard", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "SERVER_ERROR", response["code"])
	mockService.AssertExpectations(t)
}

func TestGetDashboard_EmptyStats(t *testing.T) {
	mockService := new(MockDashboardService)
	handler := NewDashboardHandler(mockService)
	tokenGen := jwt.NewTokenGenerator("test-secret", "test-refresh-secret", 15, 336)

	snapshot := &dashModel.DashboardSnapshot{
		RecentFeedback:   []*fbModel.FeedbackItem{},
		Stats:            map[string]int{},
		SuggestedActions: []string{},
	}

	mockService.On("GetDashboard", mock.Anything, "test-user-id").Return(snapshot, nil)

	router := setupDashboardRouter(handler, tokenGen)
	req, _ := http.NewRequest("GET", "/api/v1/dashboard", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response dashModel.DashboardSnapshot
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(response.Stats))
	assert.Equal(t, 0, len(response.SuggestedActions))
	mockService.AssertExpectations(t)
}

