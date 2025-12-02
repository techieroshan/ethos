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
	"ethos/internal/auth/model"
	fbModel "ethos/internal/feedback/model"
	"ethos/internal/feedback/service"
	"ethos/internal/middleware"
	"ethos/pkg/jwt"
)

// MockFeedbackService is a mock implementation of the feedback service
type MockFeedbackService struct {
	mock.Mock
}

func (m *MockFeedbackService) GetFeed(ctx context.Context, limit, offset int) ([]*fbModel.FeedbackItem, int, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*fbModel.FeedbackItem), args.Get(1).(int), args.Error(2)
}

func (m *MockFeedbackService) GetFeedbackByID(ctx context.Context, feedbackID string) (*fbModel.FeedbackItem, error) {
	args := m.Called(ctx, feedbackID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*fbModel.FeedbackItem), args.Error(1)
}

func (m *MockFeedbackService) GetComments(ctx context.Context, feedbackID string, limit, offset int) ([]*fbModel.FeedbackComment, int, error) {
	args := m.Called(ctx, feedbackID, limit, offset)
	if args.Get(0) == nil {
		return nil, 0, args.Error(3)
	}
	return args.Get(0).([]*fbModel.FeedbackComment), args.Get(1).(int), args.Error(3)
}

func (m *MockFeedbackService) CreateFeedback(ctx context.Context, userID string, req *service.CreateFeedbackRequest) (*fbModel.FeedbackItem, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*fbModel.FeedbackItem), args.Error(1)
}

func (m *MockFeedbackService) CreateComment(ctx context.Context, userID, feedbackID string, req *service.CreateCommentRequest) (*fbModel.FeedbackComment, error) {
	args := m.Called(ctx, userID, feedbackID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*fbModel.FeedbackComment), args.Error(1)
}

func (m *MockFeedbackService) AddReaction(ctx context.Context, userID, feedbackID string, reactionType string) error {
	args := m.Called(ctx, userID, feedbackID, reactionType)
	return args.Error(0)
}

func (m *MockFeedbackService) RemoveReaction(ctx context.Context, userID, feedbackID string, reactionType string) error {
	args := m.Called(ctx, userID, feedbackID, reactionType)
	return args.Error(0)
}

func (m *MockFeedbackService) GetTemplates(ctx context.Context, contextFilter, tagsFilter string) ([]*fbModel.FeedbackTemplate, error) {
	args := m.Called(ctx, contextFilter, tagsFilter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*fbModel.FeedbackTemplate), args.Error(1)
}

func (m *MockFeedbackService) SubmitTemplateSuggestion(ctx context.Context, req *service.TemplateSuggestionRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *MockFeedbackService) GetImpact(ctx context.Context, userID *string, from, to *time.Time) (*fbModel.FeedbackImpact, error) {
	args := m.Called(ctx, userID, from, to)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*fbModel.FeedbackImpact), args.Error(1)
}

func (m *MockFeedbackService) CreateBatchFeedback(ctx context.Context, userID string, req *service.BatchFeedbackRequest) (*service.BatchFeedbackResponse, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.BatchFeedbackResponse), args.Error(1)
}

func (m *MockFeedbackService) GetFeedWithFilters(ctx context.Context, limit, offset int, filters *service.FeedFilters) ([]*fbModel.FeedbackItem, int, error) {
	args := m.Called(ctx, limit, offset, filters)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*fbModel.FeedbackItem), args.Get(1).(int), args.Error(2)
}

func (m *MockFeedbackService) GetBookmarks(ctx context.Context, userID string, limit, offset int) ([]*fbModel.FeedbackItem, int, error) {
	args := m.Called(ctx, userID, limit, offset)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*fbModel.FeedbackItem), args.Get(1).(int), args.Error(2)
}

func (m *MockFeedbackService) AddBookmark(ctx context.Context, userID, feedbackID string) error {
	args := m.Called(ctx, userID, feedbackID)
	return args.Error(0)
}

func (m *MockFeedbackService) RemoveBookmark(ctx context.Context, userID, feedbackID string) error {
	args := m.Called(ctx, userID, feedbackID)
	return args.Error(0)
}

func (m *MockFeedbackService) ExportFeedback(ctx context.Context, filters *service.FeedFilters, format string) (*service.ExportResponse, error) {
	args := m.Called(ctx, filters, format)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.ExportResponse), args.Error(1)
}

func setupFeedbackRouter(handler *FeedbackHandler, tokenGen *jwt.TokenGenerator) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.AuthMiddleware(tokenGen))
	router.GET("/api/v1/feedback/feed", handler.GetFeed)
	router.GET("/api/v1/feedback/:feedback_id", handler.GetFeedbackByID)
	router.GET("/api/v1/feedback/:feedback_id/comments", handler.GetComments)
	router.POST("/api/v1/feedback", handler.CreateFeedback)
	router.POST("/api/v1/feedback/:feedback_id/comments", handler.CreateComment)
	router.POST("/api/v1/feedback/:feedback_id/react", handler.AddReaction)
	router.DELETE("/api/v1/feedback/:feedback_id/react", handler.RemoveReaction)
	return router
}

func TestGetFeed_ValidRequest(t *testing.T) {
	mockService := new(MockFeedbackService)
	handler := NewFeedbackHandler(mockService)

	tokenGen := jwt.NewTokenGenerator(
		"test-access-secret",
		"test-refresh-secret",
		15*60*time.Second,
		14*24*60*60*time.Second,
	)

	userID := "user-123"
	token, err := tokenGen.GenerateAccessToken(userID)
	assert.NoError(t, err)

	expectedItems := []*fbModel.FeedbackItem{
		{
			FeedbackID: "f-001",
			Author: &model.UserSummary{ID: "user-234", Name: "Lisa K."},
			Content: "Really enjoying the new feature!",
			Reactions: map[string]int{"like": 5, "helpful": 2},
			CommentsCount: 3,
			CreatedAt: time.Now(),
		},
	}

	mockService.On("GetFeed", mock.Anything, 20, 0).Return(expectedItems, 1, nil)

	router := setupFeedbackRouter(handler, tokenGen)
	req, _ := http.NewRequest("GET", "/api/v1/feedback/feed?limit=20", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotNil(t, response["results"])
	mockService.AssertExpectations(t)
}

