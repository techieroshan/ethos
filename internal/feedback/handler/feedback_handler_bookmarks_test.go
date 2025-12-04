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
authModel "ethos/internal/auth/model"
"ethos/internal/feedback"
fbModel "ethos/internal/feedback/model"
"ethos/internal/feedback/service"
"ethos/internal/middleware"
"ethos/pkg/jwt"
)

// MockFeedbackServiceForBookmarks is a mock implementation for bookmark tests
type MockFeedbackServiceForBookmarks struct {
	mock.Mock
}

func (m *MockFeedbackServiceForBookmarks) GetFeed(ctx context.Context, limit, offset int) ([]*fbModel.FeedbackItem, int, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*fbModel.FeedbackItem), args.Get(1).(int), args.Error(2)
}

func (m *MockFeedbackServiceForBookmarks) GetFeedbackByID(ctx context.Context, feedbackID string) (*fbModel.FeedbackItem, error) {
	args := m.Called(ctx, feedbackID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*fbModel.FeedbackItem), args.Error(1)
}

func (m *MockFeedbackServiceForBookmarks) GetComments(ctx context.Context, feedbackID string, limit, offset int) ([]*fbModel.FeedbackComment, int, error) {
	args := m.Called(ctx, feedbackID, limit, offset)
	if args.Get(0) == nil {
		return nil, 0, args.Error(3)
	}
	return args.Get(0).([]*fbModel.FeedbackComment), args.Get(1).(int), args.Error(3)
}

func (m *MockFeedbackServiceForBookmarks) CreateFeedback(ctx context.Context, userID string, req *service.CreateFeedbackRequest) (*fbModel.FeedbackItem, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*fbModel.FeedbackItem), args.Error(1)
}

func (m *MockFeedbackServiceForBookmarks) CreateComment(ctx context.Context, userID, feedbackID string, req *service.CreateCommentRequest) (*fbModel.FeedbackComment, error) {
	args := m.Called(ctx, userID, feedbackID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*fbModel.FeedbackComment), args.Error(1)
}

func (m *MockFeedbackServiceForBookmarks) AddReaction(ctx context.Context, userID, feedbackID string, reactionType string) error {
	args := m.Called(ctx, userID, feedbackID, reactionType)
	return args.Error(0)
}

func (m *MockFeedbackServiceForBookmarks) RemoveReaction(ctx context.Context, userID, feedbackID string, reactionType string) error {
	args := m.Called(ctx, userID, feedbackID, reactionType)
	return args.Error(0)
}

func (m *MockFeedbackServiceForBookmarks) GetTemplates(ctx context.Context, contextFilter, tagsFilter string) ([]*fbModel.FeedbackTemplate, error) {
	args := m.Called(ctx, contextFilter, tagsFilter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*fbModel.FeedbackTemplate), args.Error(1)
}

func (m *MockFeedbackServiceForBookmarks) SubmitTemplateSuggestion(ctx context.Context, req *feedback.TemplateSuggestionRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *MockFeedbackServiceForBookmarks) GetImpact(ctx context.Context, userID *string, from, to *time.Time) (*fbModel.FeedbackImpact, error) {
	args := m.Called(ctx, userID, from, to)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*fbModel.FeedbackImpact), args.Error(1)
}

func (m *MockFeedbackServiceForBookmarks) CreateBatchFeedback(ctx context.Context, userID string, req *feedback.BatchFeedbackRequest) (*feedback.BatchFeedbackResponse, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*feedback.BatchFeedbackResponse), args.Error(1)
}

func (m *MockFeedbackServiceForBookmarks) GetFeedWithFilters(ctx context.Context, limit, offset int, filters *feedback.FeedFilters) ([]*fbModel.FeedbackItem, int, error) {
	args := m.Called(ctx, limit, offset, filters)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*fbModel.FeedbackItem), args.Get(1).(int), args.Error(2)
}

func (m *MockFeedbackServiceForBookmarks) GetBookmarks(ctx context.Context, userID string, limit, offset int) ([]*fbModel.FeedbackItem, int, error) {
	args := m.Called(ctx, userID, limit, offset)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*fbModel.FeedbackItem), args.Get(1).(int), args.Error(2)
}

func (m *MockFeedbackServiceForBookmarks) AddBookmark(ctx context.Context, userID, feedbackID string) error {
	args := m.Called(ctx, userID, feedbackID)
	return args.Error(0)
}

func (m *MockFeedbackServiceForBookmarks) RemoveBookmark(ctx context.Context, userID, feedbackID string) error {
	args := m.Called(ctx, userID, feedbackID)
	return args.Error(0)
}

func (m *MockFeedbackServiceForBookmarks) ExportFeedback(ctx context.Context, filters *feedback.FeedFilters, format string) (*feedback.ExportResponse, error) {
	args := m.Called(ctx, filters, format)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*feedback.ExportResponse), args.Error(1)
}

func (m *MockFeedbackServiceForBookmarks) UpdateFeedback(ctx context.Context, userID, feedbackID string, req *service.UpdateFeedbackRequest) (*fbModel.FeedbackItem, error) {
	args := m.Called(ctx, userID, feedbackID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*fbModel.FeedbackItem), args.Error(1)
}

func (m *MockFeedbackServiceForBookmarks) DeleteFeedback(ctx context.Context, userID, feedbackID string) error {
	args := m.Called(ctx, userID, feedbackID)
	return args.Error(0)
}

func (m *MockFeedbackServiceForBookmarks) UpdateComment(ctx context.Context, userID, feedbackID, commentID string, req *service.UpdateCommentRequest) (*fbModel.FeedbackComment, error) {
	args := m.Called(ctx, userID, feedbackID, commentID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*fbModel.FeedbackComment), args.Error(1)
}

func (m *MockFeedbackServiceForBookmarks) DeleteComment(ctx context.Context, userID, feedbackID, commentID string) error {
	args := m.Called(ctx, userID, feedbackID, commentID)
	return args.Error(0)
}

func (m *MockFeedbackServiceForBookmarks) GetFeedbackAnalytics(ctx context.Context, userID *string, from, to *time.Time) (*fbModel.FeedbackAnalytics, error) {
	args := m.Called(ctx, userID, from, to)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*fbModel.FeedbackAnalytics), args.Error(1)
}

func setupFeedbackRouterForBookmarks(handler *FeedbackHandler, tokenGen *jwt.TokenGenerator) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.AuthMiddleware(tokenGen))
	v1 := router.Group("/api/v1")
	feedback := v1.Group("/feedback")
	feedback.GET("/bookmarks", handler.GetBookmarks)
	feedback.POST("/bookmarks/:feedback_id", handler.AddBookmark)
	feedback.DELETE("/bookmarks/:feedback_id", handler.RemoveBookmark)
	return router
}

func TestGetBookmarks_Success(t *testing.T) {
	mockService := new(MockFeedbackServiceForBookmarks)
	handler := NewFeedbackHandler(mockService)
	tokenGen := jwt.NewTokenGenerator("test-secret", "test-refresh-secret", 15*time.Second, 336*time.Hour)

	fbTypeAppreciation := fbModel.FeedbackTypeAppreciation
	fbVisPublic := fbModel.FeedbackVisibilityPublic

	bookmarks := []*fbModel.FeedbackItem{
		{
			FeedbackID:    "f-001",
			Author:        &authModel.UserSummary{ID: "user-1", Name: "John Doe"},
			Content:       "Great work on the project!",
			Type:          &fbTypeAppreciation,
			Visibility:    &fbVisPublic,
			Reactions:     map[string]int{"like": 5},
			CommentsCount: 2,
			CreatedAt:     time.Now(),
		},
	}

	mockService.On("GetBookmarks", mock.Anything, "user-123", 20, 0).Return(bookmarks, 1, nil)

	router := setupFeedbackRouterForBookmarks(handler, tokenGen)
	req, _ := http.NewRequest("GET", "/api/v1/feedback/bookmarks", nil)
	token, err := tokenGen.GenerateAccessToken("user-123")
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	err2 := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err2)
	results := response["results"].([]interface{})
	assert.Equal(t, 1, len(results))
	mockService.AssertExpectations(t)
}
