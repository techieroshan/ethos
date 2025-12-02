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
	"ethos/internal/feedback/service"
	"ethos/internal/middleware"
	"ethos/pkg/jwt"
)

// MockFeedbackServiceForFilters is a mock implementation for filtering tests
type MockFeedbackServiceForFilters struct {
	mock.Mock
}

func (m *MockFeedbackServiceForFilters) GetFeed(ctx context.Context, limit, offset int) ([]*fbModel.FeedbackItem, int, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*fbModel.FeedbackItem), args.Get(1).(int), args.Error(2)
}

func (m *MockFeedbackServiceForFilters) GetFeedbackByID(ctx context.Context, feedbackID string) (*fbModel.FeedbackItem, error) {
	args := m.Called(ctx, feedbackID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*fbModel.FeedbackItem), args.Error(1)
}

func (m *MockFeedbackServiceForFilters) GetComments(ctx context.Context, feedbackID string, limit, offset int) ([]*fbModel.FeedbackComment, int, error) {
	args := m.Called(ctx, feedbackID, limit, offset)
	if args.Get(0) == nil {
		return nil, 0, args.Error(3)
	}
	return args.Get(0).([]*fbModel.FeedbackComment), args.Get(1).(int), args.Error(3)
}

func (m *MockFeedbackServiceForFilters) CreateFeedback(ctx context.Context, userID string, req interface{}) (*fbModel.FeedbackItem, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*fbModel.FeedbackItem), args.Error(1)
}

func (m *MockFeedbackServiceForFilters) CreateComment(ctx context.Context, userID, feedbackID string, req interface{}) (*fbModel.FeedbackComment, error) {
	args := m.Called(ctx, userID, feedbackID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*fbModel.FeedbackComment), args.Error(1)
}

func (m *MockFeedbackServiceForFilters) AddReaction(ctx context.Context, userID, feedbackID string, reactionType string) error {
	args := m.Called(ctx, userID, feedbackID, reactionType)
	return args.Error(0)
}

func (m *MockFeedbackServiceForFilters) RemoveReaction(ctx context.Context, userID, feedbackID string, reactionType string) error {
	args := m.Called(ctx, userID, feedbackID, reactionType)
	return args.Error(0)
}

func (m *MockFeedbackServiceForFilters) GetTemplates(ctx context.Context, contextFilter, tagsFilter string) ([]*fbModel.FeedbackTemplate, error) {
	args := m.Called(ctx, contextFilter, tagsFilter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*fbModel.FeedbackTemplate), args.Error(1)
}

func (m *MockFeedbackServiceForFilters) SubmitTemplateSuggestion(ctx context.Context, req interface{}) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *MockFeedbackServiceForFilters) GetImpact(ctx context.Context, userID *string, from, to *time.Time) (*fbModel.FeedbackImpact, error) {
	args := m.Called(ctx, userID, from, to)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*fbModel.FeedbackImpact), args.Error(1)
}

func (m *MockFeedbackServiceForFilters) CreateBatchFeedback(ctx context.Context, userID string, req *service.BatchFeedbackRequest) (*service.BatchFeedbackResponse, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.BatchFeedbackResponse), args.Error(1)
}

func (m *MockFeedbackServiceForFilters) GetFeedWithFilters(ctx context.Context, limit, offset int, filters *service.FeedFilters) ([]*fbModel.FeedbackItem, int, error) {
	args := m.Called(ctx, limit, offset, filters)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*fbModel.FeedbackItem), args.Get(1).(int), args.Error(2)
}

func setupFeedbackRouterForFilters(handler *FeedbackHandler, tokenGen *jwt.TokenGenerator) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.AuthMiddleware(tokenGen))
	v1 := router.Group("/api/v1")
	feedback := v1.Group("/feedback")
	feedback.GET("/feed", handler.GetFeed)
	return router
}

func TestGetFeed_WithReviewerTypeFilter(t *testing.T) {
	mockService := new(MockFeedbackServiceForFilters)
	handler := NewFeedbackHandler(mockService)
	tokenGen := jwt.NewTokenGenerator("test-secret", "test-refresh-secret", 15, 336)

	reviewerType := "org"
	filters := &service.FeedFilters{ReviewerType: &reviewerType}

	mockService.On("GetFeedWithFilters", mock.Anything, 20, 0, filters).Return([]*fbModel.FeedbackItem{}, 0, nil)

	router := setupFeedbackRouterForFilters(handler, tokenGen)
	req, _ := http.NewRequest("GET", "/api/v1/feedback/feed?reviewer_type=org", nil)
	req.Header.Set("Authorization", "Bearer "+tokenGen.GenerateAccessToken("user-123", []string{}))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestGetFeed_WithContextFilter(t *testing.T) {
	mockService := new(MockFeedbackServiceForFilters)
	handler := NewFeedbackHandler(mockService)
	tokenGen := jwt.NewTokenGenerator("test-secret", "test-refresh-secret", 15, 336)

	contextFilter := "project"
	filters := &service.FeedFilters{Context: &contextFilter}

	mockService.On("GetFeedWithFilters", mock.Anything, 20, 0, filters).Return([]*fbModel.FeedbackItem{}, 0, nil)

	router := setupFeedbackRouterForFilters(handler, tokenGen)
	req, _ := http.NewRequest("GET", "/api/v1/feedback/feed?context=project", nil)
	req.Header.Set("Authorization", "Bearer "+tokenGen.GenerateAccessToken("user-123", []string{}))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestGetFeed_WithVerificationFilter(t *testing.T) {
	mockService := new(MockFeedbackServiceForFilters)
	handler := NewFeedbackHandler(mockService)
	tokenGen := jwt.NewTokenGenerator("test-secret", "test-refresh-secret", 15, 336)

	verification := "verified"
	filters := &service.FeedFilters{Verification: &verification}

	mockService.On("GetFeedWithFilters", mock.Anything, 20, 0, filters).Return([]*fbModel.FeedbackItem{}, 0, nil)

	router := setupFeedbackRouterForFilters(handler, tokenGen)
	req, _ := http.NewRequest("GET", "/api/v1/feedback/feed?verification=verified", nil)
	req.Header.Set("Authorization", "Bearer "+tokenGen.GenerateAccessToken("user-123", []string{}))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestGetFeed_WithTagsFilter(t *testing.T) {
	mockService := new(MockFeedbackServiceForFilters)
	handler := NewFeedbackHandler(mockService)
	tokenGen := jwt.NewTokenGenerator("test-secret", "test-refresh-secret", 15, 336)

	filters := &service.FeedFilters{Tags: []string{"leadership", "initiative"}}

	mockService.On("GetFeedWithFilters", mock.Anything, 20, 0, filters).Return([]*fbModel.FeedbackItem{}, 0, nil)

	router := setupFeedbackRouterForFilters(handler, tokenGen)
	req, _ := http.NewRequest("GET", "/api/v1/feedback/feed?tags=leadership,initiative", nil)
	req.Header.Set("Authorization", "Bearer "+tokenGen.GenerateAccessToken("user-123", []string{}))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestGetFeed_MultipleFilters(t *testing.T) {
	mockService := new(MockFeedbackServiceForFilters)
	handler := NewFeedbackHandler(mockService)
	tokenGen := jwt.NewTokenGenerator("test-secret", "test-refresh-secret", 15, 336)

	reviewerType := "org"
	contextFilter := "team"
	filters := &service.FeedFilters{
		ReviewerType: &reviewerType,
		Context:      &contextFilter,
		Tags:         []string{"leadership"},
	}

	mockService.On("GetFeedWithFilters", mock.Anything, 20, 0, filters).Return([]*fbModel.FeedbackItem{}, 0, nil)

	router := setupFeedbackRouterForFilters(handler, tokenGen)
	req, _ := http.NewRequest("GET", "/api/v1/feedback/feed?reviewer_type=org&context=team&tags=leadership", nil)
	req.Header.Set("Authorization", "Bearer "+tokenGen.GenerateAccessToken("user-123", []string{}))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestGetFeed_NoFilters(t *testing.T) {
	mockService := new(MockFeedbackServiceForFilters)
	handler := NewFeedbackHandler(mockService)
	tokenGen := jwt.NewTokenGenerator("test-secret", "test-refresh-secret", 15, 336)

	mockService.On("GetFeed", mock.Anything, 20, 0).Return([]*fbModel.FeedbackItem{}, 0, nil)

	router := setupFeedbackRouterForFilters(handler, tokenGen)
	req, _ := http.NewRequest("GET", "/api/v1/feedback/feed", nil)
	req.Header.Set("Authorization", "Bearer "+tokenGen.GenerateAccessToken("user-123", []string{}))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}
