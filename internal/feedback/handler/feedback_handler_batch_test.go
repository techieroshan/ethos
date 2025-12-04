package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"ethos/internal/feedback"
	fbModel "ethos/internal/feedback/model"
	"ethos/internal/feedback/service"
	"ethos/internal/middleware"
	"ethos/pkg/jwt"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
) // MockFeedbackServiceForBatch is a mock implementation for batch tests
type MockFeedbackServiceForBatch struct {
	mock.Mock
}

func (m *MockFeedbackServiceForBatch) GetFeed(ctx context.Context, limit, offset int) ([]*fbModel.FeedbackItem, int, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*fbModel.FeedbackItem), args.Get(1).(int), args.Error(2)
}

func (m *MockFeedbackServiceForBatch) GetFeedbackByID(ctx context.Context, feedbackID string) (*fbModel.FeedbackItem, error) {
	args := m.Called(ctx, feedbackID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*fbModel.FeedbackItem), args.Error(1)
}

func (m *MockFeedbackServiceForBatch) GetComments(ctx context.Context, feedbackID string, limit, offset int) ([]*fbModel.FeedbackComment, int, error) {
	args := m.Called(ctx, feedbackID, limit, offset)
	if args.Get(0) == nil {
		return nil, 0, args.Error(3)
	}
	return args.Get(0).([]*fbModel.FeedbackComment), args.Get(1).(int), args.Error(3)
}

func (m *MockFeedbackServiceForBatch) CreateFeedback(ctx context.Context, userID string, req *service.CreateFeedbackRequest) (*fbModel.FeedbackItem, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*fbModel.FeedbackItem), args.Error(1)
}

func (m *MockFeedbackServiceForBatch) CreateComment(ctx context.Context, userID, feedbackID string, req *service.CreateCommentRequest) (*fbModel.FeedbackComment, error) {
	args := m.Called(ctx, userID, feedbackID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*fbModel.FeedbackComment), args.Error(1)
}

func (m *MockFeedbackServiceForBatch) AddReaction(ctx context.Context, userID, feedbackID string, reactionType string) error {
	args := m.Called(ctx, userID, feedbackID, reactionType)
	return args.Error(0)
}

func (m *MockFeedbackServiceForBatch) RemoveReaction(ctx context.Context, userID, feedbackID string, reactionType string) error {
	args := m.Called(ctx, userID, feedbackID, reactionType)
	return args.Error(0)
}

func (m *MockFeedbackServiceForBatch) GetTemplates(ctx context.Context, contextFilter, tagsFilter string) ([]*fbModel.FeedbackTemplate, error) {
	args := m.Called(ctx, contextFilter, tagsFilter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*fbModel.FeedbackTemplate), args.Error(1)
}

func (m *MockFeedbackServiceForBatch) SubmitTemplateSuggestion(ctx context.Context, req *feedback.TemplateSuggestionRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *MockFeedbackServiceForBatch) GetImpact(ctx context.Context, userID *string, from, to *time.Time) (*fbModel.FeedbackImpact, error) {
	args := m.Called(ctx, userID, from, to)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*fbModel.FeedbackImpact), args.Error(1)
}

func (m *MockFeedbackServiceForBatch) CreateBatchFeedback(ctx context.Context, userID string, req *feedback.BatchFeedbackRequest) (*feedback.BatchFeedbackResponse, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*feedback.BatchFeedbackResponse), args.Error(1)
}

func (m *MockFeedbackServiceForBatch) GetFeedWithFilters(ctx context.Context, limit, offset int, filters *feedback.FeedFilters) ([]*fbModel.FeedbackItem, int, error) {
	args := m.Called(ctx, limit, offset, filters)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*fbModel.FeedbackItem), args.Get(1).(int), args.Error(2)
}

func (m *MockFeedbackServiceForBatch) GetBookmarks(ctx context.Context, userID string, limit, offset int) ([]*fbModel.FeedbackItem, int, error) {
	args := m.Called(ctx, userID, limit, offset)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*fbModel.FeedbackItem), args.Get(1).(int), args.Error(2)
}

func (m *MockFeedbackServiceForBatch) AddBookmark(ctx context.Context, userID, feedbackID string) error {
	args := m.Called(ctx, userID, feedbackID)
	return args.Error(0)
}

func (m *MockFeedbackServiceForBatch) RemoveBookmark(ctx context.Context, userID, feedbackID string) error {
	args := m.Called(ctx, userID, feedbackID)
	return args.Error(0)
}

func (m *MockFeedbackServiceForBatch) ExportFeedback(ctx context.Context, filters *feedback.FeedFilters, format string) (*feedback.ExportResponse, error) {
	args := m.Called(ctx, filters, format)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*feedback.ExportResponse), args.Error(1)
}

func (m *MockFeedbackServiceForBatch) UpdateFeedback(ctx context.Context, userID, feedbackID string, req *service.UpdateFeedbackRequest) (*fbModel.FeedbackItem, error) {
	args := m.Called(ctx, userID, feedbackID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*fbModel.FeedbackItem), args.Error(1)
}

func (m *MockFeedbackServiceForBatch) DeleteFeedback(ctx context.Context, userID, feedbackID string) error {
	args := m.Called(ctx, userID, feedbackID)
	return args.Error(0)
}

func (m *MockFeedbackServiceForBatch) UpdateComment(ctx context.Context, userID, feedbackID, commentID string, req *service.UpdateCommentRequest) (*fbModel.FeedbackComment, error) {
	args := m.Called(ctx, userID, feedbackID, commentID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*fbModel.FeedbackComment), args.Error(1)
}

func (m *MockFeedbackServiceForBatch) DeleteComment(ctx context.Context, userID, feedbackID, commentID string) error {
	args := m.Called(ctx, userID, feedbackID, commentID)
	return args.Error(0)
}

func (m *MockFeedbackServiceForBatch) GetFeedbackAnalytics(ctx context.Context, userID *string, from, to *time.Time) (*fbModel.FeedbackAnalytics, error) {
	args := m.Called(ctx, userID, from, to)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*fbModel.FeedbackAnalytics), args.Error(1)
}

func setupFeedbackRouterForBatch(handler *FeedbackHandler, tokenGen *jwt.TokenGenerator) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.AuthMiddleware(tokenGen))
	v1 := router.Group("/api/v1")
	feedback := v1.Group("/feedback")
	feedback.POST("/batch", handler.CreateBatchFeedback)
	return router
}

func TestCreateBatchFeedback_Success(t *testing.T) {
	mockService := new(MockFeedbackServiceForBatch)
	handler := NewFeedbackHandler(mockService)
	tokenGen := jwt.NewTokenGenerator("test-secret", "test-refresh-secret", 15*time.Second, 336*time.Hour)

	fbTypeAppreciation := fbModel.FeedbackTypeAppreciation

	batchResponse := &feedback.BatchFeedbackResponse{
		Submitted: []feedback.BatchFeedbackResult{},
	}

	mockService.On("CreateBatchFeedback", mock.Anything, "user-123", mock.Anything).Return(batchResponse, nil)

	router := setupFeedbackRouterForBatch(handler, tokenGen)

	typeStr := string(fbTypeAppreciation)
	batchReq := feedback.BatchFeedbackRequest{
		Items: []feedback.BatchFeedbackItem{
			{
				Content: "Great feedback",
				Type:    &typeStr,
			},
		},
	}

	body, _ := json.Marshal(batchReq)
	req := httptest.NewRequest("POST", "/api/v1/feedback/batch", strings.NewReader(string(body)))
	token, err := tokenGen.GenerateAccessToken("user-123")
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}
