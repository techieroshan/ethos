package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"ethos/internal/feedback"
	fbModel "ethos/internal/feedback/model"
	"ethos/internal/feedback/service"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
) // MockFeedbackServiceForTemplates is a mock implementation for templates tests
type MockFeedbackServiceForTemplates struct {
	mock.Mock
}

func (m *MockFeedbackServiceForTemplates) GetFeed(ctx context.Context, limit, offset int) ([]*fbModel.FeedbackItem, int, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*fbModel.FeedbackItem), args.Get(1).(int), args.Error(2)
}

func (m *MockFeedbackServiceForTemplates) GetFeedbackByID(ctx context.Context, feedbackID string) (*fbModel.FeedbackItem, error) {
	args := m.Called(ctx, feedbackID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*fbModel.FeedbackItem), args.Error(1)
}

func (m *MockFeedbackServiceForTemplates) GetComments(ctx context.Context, feedbackID string, limit, offset int) ([]*fbModel.FeedbackComment, int, error) {
	args := m.Called(ctx, feedbackID, limit, offset)
	if args.Get(0) == nil {
		return nil, 0, args.Error(3)
	}
	return args.Get(0).([]*fbModel.FeedbackComment), args.Get(1).(int), args.Error(3)
}

func (m *MockFeedbackServiceForTemplates) CreateFeedback(ctx context.Context, userID string, req *service.CreateFeedbackRequest) (*fbModel.FeedbackItem, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*fbModel.FeedbackItem), args.Error(1)
}

func (m *MockFeedbackServiceForTemplates) CreateComment(ctx context.Context, userID, feedbackID string, req *service.CreateCommentRequest) (*fbModel.FeedbackComment, error) {
	args := m.Called(ctx, userID, feedbackID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*fbModel.FeedbackComment), args.Error(1)
}

func (m *MockFeedbackServiceForTemplates) AddReaction(ctx context.Context, userID, feedbackID string, reactionType string) error {
	args := m.Called(ctx, userID, feedbackID, reactionType)
	return args.Error(0)
}

func (m *MockFeedbackServiceForTemplates) RemoveReaction(ctx context.Context, userID, feedbackID string, reactionType string) error {
	args := m.Called(ctx, userID, feedbackID, reactionType)
	return args.Error(0)
}

func (m *MockFeedbackServiceForTemplates) GetTemplates(ctx context.Context, contextFilter, tagsFilter string) ([]*fbModel.FeedbackTemplate, error) {
	args := m.Called(ctx, contextFilter, tagsFilter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*fbModel.FeedbackTemplate), args.Error(1)
}

func (m *MockFeedbackServiceForTemplates) SubmitTemplateSuggestion(ctx context.Context, req *feedback.TemplateSuggestionRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *MockFeedbackServiceForTemplates) GetImpact(ctx context.Context, userID *string, from, to *time.Time) (*fbModel.FeedbackImpact, error) {
	args := m.Called(ctx, userID, from, to)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*fbModel.FeedbackImpact), args.Error(1)
}

func (m *MockFeedbackServiceForTemplates) CreateBatchFeedback(ctx context.Context, userID string, req *feedback.BatchFeedbackRequest) (*feedback.BatchFeedbackResponse, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*feedback.BatchFeedbackResponse), args.Error(1)
}

func (m *MockFeedbackServiceForTemplates) GetFeedWithFilters(ctx context.Context, limit, offset int, filters *feedback.FeedFilters) ([]*fbModel.FeedbackItem, int, error) {
	args := m.Called(ctx, limit, offset, filters)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*fbModel.FeedbackItem), args.Get(1).(int), args.Error(2)
}

func (m *MockFeedbackServiceForTemplates) GetBookmarks(ctx context.Context, userID string, limit, offset int) ([]*fbModel.FeedbackItem, int, error) {
	args := m.Called(ctx, userID, limit, offset)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*fbModel.FeedbackItem), args.Get(1).(int), args.Error(2)
}

func (m *MockFeedbackServiceForTemplates) AddBookmark(ctx context.Context, userID, feedbackID string) error {
	args := m.Called(ctx, userID, feedbackID)
	return args.Error(0)
}

func (m *MockFeedbackServiceForTemplates) RemoveBookmark(ctx context.Context, userID, feedbackID string) error {
	args := m.Called(ctx, userID, feedbackID)
	return args.Error(0)
}

func (m *MockFeedbackServiceForTemplates) ExportFeedback(ctx context.Context, filters *feedback.FeedFilters, format string) (*feedback.ExportResponse, error) {
	args := m.Called(ctx, filters, format)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*feedback.ExportResponse), args.Error(1)
}

func (m *MockFeedbackServiceForTemplates) UpdateFeedback(ctx context.Context, userID, feedbackID string, req *service.UpdateFeedbackRequest) (*fbModel.FeedbackItem, error) {
	args := m.Called(ctx, userID, feedbackID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*fbModel.FeedbackItem), args.Error(1)
}

func (m *MockFeedbackServiceForTemplates) DeleteFeedback(ctx context.Context, userID, feedbackID string) error {
	args := m.Called(ctx, userID, feedbackID)
	return args.Error(0)
}

func (m *MockFeedbackServiceForTemplates) UpdateComment(ctx context.Context, userID, feedbackID, commentID string, req *service.UpdateCommentRequest) (*fbModel.FeedbackComment, error) {
	args := m.Called(ctx, userID, feedbackID, commentID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*fbModel.FeedbackComment), args.Error(1)
}

func (m *MockFeedbackServiceForTemplates) DeleteComment(ctx context.Context, userID, feedbackID, commentID string) error {
	args := m.Called(ctx, userID, feedbackID, commentID)
	return args.Error(0)
}

func (m *MockFeedbackServiceForTemplates) GetFeedbackAnalytics(ctx context.Context, userID *string, from, to *time.Time) (*fbModel.FeedbackAnalytics, error) {
	args := m.Called(ctx, userID, from, to)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*fbModel.FeedbackAnalytics), args.Error(1)
}

func TestGetTemplates_Success(t *testing.T) {
	mockService := new(MockFeedbackServiceForTemplates)
	handler := NewFeedbackHandler(mockService)

	templates := []*fbModel.FeedbackTemplate{
		{
			TemplateID:     "t-001",
			Name:           "Performance Review",
			Description:    "For performance reviews",
			ContextTags:    []string{"performance_review"},
			TemplateFields: map[string]interface{}{"question": "What are your strengths?"},
		},
	}

	mockService.On("GetTemplates", mock.Anything, "performance_review", "").Return(templates, nil)

	gin.SetMode(gin.TestMode)
	req := httptest.NewRequest("GET", "/api/v1/feedback/templates?context=performance_review", nil)
	w := httptest.NewRecorder()

	router := gin.New()
	router.GET("/api/v1/feedback/templates", handler.GetTemplates)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestGetTemplates_WithTags(t *testing.T) {
	mockService := new(MockFeedbackServiceForTemplates)
	handler := NewFeedbackHandler(mockService)

	templates := []*fbModel.FeedbackTemplate{
		{
			TemplateID:     "t-002",
			Name:           "Initiative Review",
			Description:    "For initiative feedback",
			ContextTags:    []string{"leadership", "initiative"},
			TemplateFields: map[string]interface{}{},
		},
	}

	mockService.On("GetTemplates", mock.Anything, "", "leadership,initiative").Return(templates, nil)

	gin.SetMode(gin.TestMode)
	req := httptest.NewRequest("GET", "/api/v1/feedback/templates?tags=leadership,initiative", nil)
	w := httptest.NewRecorder()

	router := gin.New()
	router.GET("/api/v1/feedback/templates", handler.GetTemplates)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}
