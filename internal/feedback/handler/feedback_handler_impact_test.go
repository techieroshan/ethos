package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	fbModel "ethos/internal/feedback/model"
	"ethos/internal/middleware"
	"ethos/pkg/jwt"
)

// MockFeedbackServiceForImpact is a mock implementation for impact analytics
type MockFeedbackServiceForImpact struct {
	mock.Mock
}

func (m *MockFeedbackServiceForImpact) GetFeed(ctx context.Context, limit, offset int) ([]*fbModel.FeedbackItem, int, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*fbModel.FeedbackItem), args.Get(1).(int), args.Error(2)
}

func (m *MockFeedbackServiceForImpact) GetFeedbackByID(ctx context.Context, feedbackID string) (*fbModel.FeedbackItem, error) {
	args := m.Called(ctx, feedbackID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*fbModel.FeedbackItem), args.Error(1)
}

func (m *MockFeedbackServiceForImpact) GetComments(ctx context.Context, feedbackID string, limit, offset int) ([]*fbModel.FeedbackComment, int, error) {
	args := m.Called(ctx, feedbackID, limit, offset)
	if args.Get(0) == nil {
		return nil, 0, args.Error(3)
	}
	return args.Get(0).([]*fbModel.FeedbackComment), args.Get(1).(int), args.Error(3)
}

func (m *MockFeedbackServiceForImpact) CreateFeedback(ctx context.Context, userID string, req interface{}) (*fbModel.FeedbackItem, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*fbModel.FeedbackItem), args.Error(1)
}

func (m *MockFeedbackServiceForImpact) CreateComment(ctx context.Context, userID, feedbackID string, req interface{}) (*fbModel.FeedbackComment, error) {
	args := m.Called(ctx, userID, feedbackID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*fbModel.FeedbackComment), args.Error(1)
}

func (m *MockFeedbackServiceForImpact) AddReaction(ctx context.Context, userID, feedbackID string, reactionType string) error {
	args := m.Called(ctx, userID, feedbackID, reactionType)
	return args.Error(0)
}

func (m *MockFeedbackServiceForImpact) RemoveReaction(ctx context.Context, userID, feedbackID string, reactionType string) error {
	args := m.Called(ctx, userID, feedbackID, reactionType)
	return args.Error(0)
}

func (m *MockFeedbackServiceForImpact) GetTemplates(ctx context.Context, contextFilter, tagsFilter string) ([]*fbModel.FeedbackTemplate, error) {
	args := m.Called(ctx, contextFilter, tagsFilter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*fbModel.FeedbackTemplate), args.Error(1)
}

func (m *MockFeedbackServiceForImpact) SubmitTemplateSuggestion(ctx context.Context, req interface{}) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *MockFeedbackServiceForImpact) GetImpact(ctx context.Context, userID *string, from, to *time.Time) (*fbModel.FeedbackImpact, error) {
	args := m.Called(ctx, userID, from, to)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*fbModel.FeedbackImpact), args.Error(1)
}

func setupFeedbackRouterForImpact(handler *FeedbackHandler, tokenGen *jwt.TokenGenerator) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.AuthMiddleware(tokenGen))
	v1 := router.Group("/api/v1")
	feedback := v1.Group("/feedback")
	feedback.GET("/impact", handler.GetImpact)
	return router
}

func TestGetImpact_Success(t *testing.T) {
	mockService := new(MockFeedbackServiceForImpact)
	handler := NewFeedbackHandler(mockService)
	tokenGen := jwt.NewTokenGenerator("test-secret", "test-refresh-secret", 15, 336)

	impact := &fbModel.FeedbackImpact{
		FeedbackCount:    31,
		AverageHelpfulness: 0.87,
		ReactionTotals: map[string]int{
			"like":      120,
			"helpful":   53,
			"insightful": 12,
		},
		FollowUpCount: 7,
		Trends: []fbModel.FeedbackTrend{
			{
				Date:         time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
				Helpfulness:  0.91,
				FeedbackSubmitted: 4,
			},
		},
	}

	mockService.On("GetImpact", mock.Anything, (*string)(nil), (*time.Time)(nil), (*time.Time)(nil)).Return(impact, nil)

	router := setupFeedbackRouterForImpact(handler, tokenGen)
	req, _ := http.NewRequest("GET", "/api/v1/feedback/impact", nil)
	req.Header.Set("Authorization", "Bearer "+tokenGen.GenerateAccessToken("user-123", []string{}))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response fbModel.FeedbackImpact
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 31, response.FeedbackCount)
	assert.Equal(t, 0.87, response.AverageHelpfulness)
	mockService.AssertExpectations(t)
}

func TestGetImpact_WithUserIDFilter(t *testing.T) {
	mockService := new(MockFeedbackServiceForImpact)
	handler := NewFeedbackHandler(mockService)
	tokenGen := jwt.NewTokenGenerator("test-secret", "test-refresh-secret", 15, 336)

	userID := "user-8822"
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC)

	mockService.On("GetImpact", mock.Anything, &userID, &from, &to).Return(&fbModel.FeedbackImpact{}, nil)

	router := setupFeedbackRouterForImpact(handler, tokenGen)
	req, _ := http.NewRequest("GET", "/api/v1/feedback/impact?user_id=user-8822&from=2024-01-01&to=2024-07-01", nil)
	req.Header.Set("Authorization", "Bearer "+tokenGen.GenerateAccessToken("user-123", []string{}))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestGetImpact_ServiceError(t *testing.T) {
	mockService := new(MockFeedbackServiceForImpact)
	handler := NewFeedbackHandler(mockService)
	tokenGen := jwt.NewTokenGenerator("test-secret", "test-refresh-secret", 15, 336)

	mockService.On("GetImpact", mock.Anything, (*string)(nil), (*time.Time)(nil), (*time.Time)(nil)).Return(nil, assert.AnError)

	router := setupFeedbackRouterForImpact(handler, tokenGen)
	req, _ := http.NewRequest("GET", "/api/v1/feedback/impact", nil)
	req.Header.Set("Authorization", "Bearer "+tokenGen.GenerateAccessToken("user-123", []string{}))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

func TestGetImpact_Unauthorized(t *testing.T) {
	mockService := new(MockFeedbackServiceForImpact)
	handler := NewFeedbackHandler(mockService)
	tokenGen := jwt.NewTokenGenerator("test-secret", "test-refresh-secret", 15, 336)

	router := setupFeedbackRouterForImpact(handler, tokenGen)
	req, _ := http.NewRequest("GET", "/api/v1/feedback/impact", nil)
	// No Authorization header
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
