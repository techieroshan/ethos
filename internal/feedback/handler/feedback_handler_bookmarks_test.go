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

func (m *MockFeedbackServiceForBookmarks) CreateFeedback(ctx context.Context, userID string, req interface{}) (*fbModel.FeedbackItem, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*fbModel.FeedbackItem), args.Error(1)
}

func (m *MockFeedbackServiceForBookmarks) CreateComment(ctx context.Context, userID, feedbackID string, req interface{}) (*fbModel.FeedbackComment, error) {
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

func (m *MockFeedbackServiceForBookmarks) SubmitTemplateSuggestion(ctx context.Context, req interface{}) error {
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

func (m *MockFeedbackServiceForBookmarks) CreateBatchFeedback(ctx context.Context, userID string, req *service.BatchFeedbackRequest) (*service.BatchFeedbackResponse, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.BatchFeedbackResponse), args.Error(1)
}

func (m *MockFeedbackServiceForBookmarks) GetFeedWithFilters(ctx context.Context, limit, offset int, filters *service.FeedFilters) ([]*fbModel.FeedbackItem, int, error) {
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
	tokenGen := jwt.NewTokenGenerator("test-secret", "test-refresh-secret", 15, 336)

	bookmarks := []*fbModel.FeedbackItem{
		{
			FeedbackID:   "f-001",
			Author:       &fbModel.UserSummary{ID: "user-1", Name: "John Doe"},
			Content:      "Great work on the project!",
			Type:         &fbModel.FeedbackTypeAppreciation,
			Visibility:   &fbModel.FeedbackVisibilityPublic,
			Reactions:    map[string]int{"like": 5},
			CommentsCount: 2,
			CreatedAt:    time.Now(),
		},
	}

	mockService.On("GetBookmarks", mock.Anything, "user-123", 20, 0).Return(bookmarks, 1, nil)

	router := setupFeedbackRouterForBookmarks(handler, tokenGen)
	req, _ := http.NewRequest("GET", "/api/v1/feedback/bookmarks", nil)
	req.Header.Set("Authorization", "Bearer "+tokenGen.GenerateAccessToken("user-123", []string{}))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	results := response["results"].([]interface{})
	assert.Equal(t, 1, len(results))
	mockService.AssertExpectations(t)
}

func TestGetBookmarks_WithPagination(t *testing.T) {
	mockService := new(MockFeedbackServiceForBookmarks)
	handler := NewFeedbackHandler(mockService)
	tokenGen := jwt.NewTokenGenerator("test-secret", "test-refresh-secret", 15, 336)

	mockService.On("GetBookmarks", mock.Anything, "user-123", 10, 5).Return([]*fbModel.FeedbackItem{}, 0, nil)

	router := setupFeedbackRouterForBookmarks(handler, tokenGen)
	req, _ := http.NewRequest("GET", "/api/v1/feedback/bookmarks?limit=10&offset=5", nil)
	req.Header.Set("Authorization", "Bearer "+tokenGen.GenerateAccessToken("user-123", []string{}))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestGetBookmarks_EmptyList(t *testing.T) {
	mockService := new(MockFeedbackServiceForBookmarks)
	handler := NewFeedbackHandler(mockService)
	tokenGen := jwt.NewTokenGenerator("test-secret", "test-refresh-secret", 15, 336)

	mockService.On("GetBookmarks", mock.Anything, "user-123", 20, 0).Return([]*fbModel.FeedbackItem{}, 0, nil)

	router := setupFeedbackRouterForBookmarks(handler, tokenGen)
	req, _ := http.NewRequest("GET", "/api/v1/feedback/bookmarks", nil)
	req.Header.Set("Authorization", "Bearer "+tokenGen.GenerateAccessToken("user-123", []string{}))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	results := response["results"].([]interface{})
	assert.Equal(t, 0, len(results))
	mockService.AssertExpectations(t)
}

func TestGetBookmarks_ServiceError(t *testing.T) {
	mockService := new(MockFeedbackServiceForBookmarks)
	handler := NewFeedbackHandler(mockService)
	tokenGen := jwt.NewTokenGenerator("test-secret", "test-refresh-secret", 15, 336)

	mockService.On("GetBookmarks", mock.Anything, "user-123", 20, 0).Return(nil, 0, assert.AnError)

	router := setupFeedbackRouterForBookmarks(handler, tokenGen)
	req, _ := http.NewRequest("GET", "/api/v1/feedback/bookmarks", nil)
	req.Header.Set("Authorization", "Bearer "+tokenGen.GenerateAccessToken("user-123", []string{}))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

func TestAddBookmark_Success(t *testing.T) {
	mockService := new(MockFeedbackServiceForBookmarks)
	handler := NewFeedbackHandler(mockService)
	tokenGen := jwt.NewTokenGenerator("test-secret", "test-refresh-secret", 15, 336)

	mockService.On("AddBookmark", mock.Anything, "user-123", "f-001").Return(nil)

	router := setupFeedbackRouterForBookmarks(handler, tokenGen)
	req, _ := http.NewRequest("POST", "/api/v1/feedback/bookmarks/f-001", nil)
	req.Header.Set("Authorization", "Bearer "+tokenGen.GenerateAccessToken("user-123", []string{}))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "bookmark_added", response["status"])
	mockService.AssertExpectations(t)
}

func TestAddBookmark_InvalidFeedbackID(t *testing.T) {
	mockService := new(MockFeedbackServiceForBookmarks)
	handler := NewFeedbackHandler(mockService)
	tokenGen := jwt.NewTokenGenerator("test-secret", "test-refresh-secret", 15, 336)

	router := setupFeedbackRouterForBookmarks(handler, tokenGen)
	req, _ := http.NewRequest("POST", "/api/v1/feedback/bookmarks/", nil)
	req.Header.Set("Authorization", "Bearer "+tokenGen.GenerateAccessToken("user-123", []string{}))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestAddBookmark_AlreadyExists(t *testing.T) {
	mockService := new(MockFeedbackServiceForBookmarks)
	handler := NewFeedbackHandler(mockService)
	tokenGen := jwt.NewTokenGenerator("test-secret", "test-refresh-secret", 15, 336)

	// For now, we'll treat this as a successful operation (idempotent)
	mockService.On("AddBookmark", mock.Anything, "user-123", "f-001").Return(nil)

	router := setupFeedbackRouterForBookmarks(handler, tokenGen)
	req, _ := http.NewRequest("POST", "/api/v1/feedback/bookmarks/f-001", nil)
	req.Header.Set("Authorization", "Bearer "+tokenGen.GenerateAccessToken("user-123", []string{}))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestRemoveBookmark_Success(t *testing.T) {
	mockService := new(MockFeedbackServiceForBookmarks)
	handler := NewFeedbackHandler(mockService)
	tokenGen := jwt.NewTokenGenerator("test-secret", "test-refresh-secret", 15, 336)

	mockService.On("RemoveBookmark", mock.Anything, "user-123", "f-001").Return(nil)

	router := setupFeedbackRouterForBookmarks(handler, tokenGen)
	req, _ := http.NewRequest("DELETE", "/api/v1/feedback/bookmarks/f-001", nil)
	req.Header.Set("Authorization", "Bearer "+tokenGen.GenerateAccessToken("user-123", []string{}))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "bookmark_removed", response["status"])
	mockService.AssertExpectations(t)
}

func TestRemoveBookmark_NotFound(t *testing.T) {
	mockService := new(MockFeedbackServiceForBookmarks)
	handler := NewFeedbackHandler(mockService)
	tokenGen := jwt.NewTokenGenerator("test-secret", "test-refresh-secret", 15, 336)

	// For now, we'll treat removal of non-existent bookmark as successful (idempotent)
	mockService.On("RemoveBookmark", mock.Anything, "user-123", "f-999").Return(nil)

	router := setupFeedbackRouterForBookmarks(handler, tokenGen)
	req, _ := http.NewRequest("DELETE", "/api/v1/feedback/bookmarks/f-999", nil)
	req.Header.Set("Authorization", "Bearer "+tokenGen.GenerateAccessToken("user-123", []string{}))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestBookmark_Unauthorized(t *testing.T) {
	mockService := new(MockFeedbackServiceForBookmarks)
	handler := NewFeedbackHandler(mockService)
	tokenGen := jwt.NewTokenGenerator("test-secret", "test-refresh-secret", 15, 336)

	router := setupFeedbackRouterForBookmarks(handler, tokenGen)
	req, _ := http.NewRequest("GET", "/api/v1/feedback/bookmarks", nil)
	// No Authorization header
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
