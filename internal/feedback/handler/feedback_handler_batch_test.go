package handler

import (
	"bytes"
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

// MockFeedbackServiceForBatch is a mock implementation for batch feedback
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

func (m *MockFeedbackServiceForBatch) CreateFeedback(ctx context.Context, userID string, req interface{}) (*fbModel.FeedbackItem, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*fbModel.FeedbackItem), args.Error(1)
}

func (m *MockFeedbackServiceForBatch) CreateComment(ctx context.Context, userID, feedbackID string, req interface{}) (*fbModel.FeedbackComment, error) {
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

func (m *MockFeedbackServiceForBatch) SubmitTemplateSuggestion(ctx context.Context, req interface{}) error {
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

func (m *MockFeedbackServiceForBatch) CreateBatchFeedback(ctx context.Context, userID string, req *service.BatchFeedbackRequest) (*service.BatchFeedbackResponse, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.BatchFeedbackResponse), args.Error(1)
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
	tokenGen := jwt.NewTokenGenerator("test-secret", "test-refresh-secret", 15, 336)

	request := &service.BatchFeedbackRequest{
		Items: []service.BatchFeedbackItem{
			{
				Content:    "Great presentation in the meeting!",
				Type:       "appreciation",
				Visibility: "public",
				IsAnonymous: false,
			},
			{
				Content:     "Consider slowing down during the Q&A.",
				Type:        "suggestion",
				Visibility:  "org",
				IsAnonymous: true,
			},
		},
	}

	response := &service.BatchFeedbackResponse{
		Submitted: []service.BatchFeedbackResult{
			{FeedbackID: "f-741", Status: "created"},
			{FeedbackID: "f-742", Status: "created"},
		},
	}

	mockService.On("CreateBatchFeedback", mock.Anything, "user-123", request).Return(response, nil)

	router := setupFeedbackRouterForBatch(handler, tokenGen)

	requestBody, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", "/api/v1/feedback/batch", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tokenGen.GenerateAccessToken("user-123", []string{}))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp service.BatchFeedbackResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(resp.Submitted))
	assert.Equal(t, "f-741", resp.Submitted[0].FeedbackID)
	mockService.AssertExpectations(t)
}

func TestCreateBatchFeedback_EmptyItems(t *testing.T) {
	mockService := new(MockFeedbackServiceForBatch)
	handler := NewFeedbackHandler(mockService)
	tokenGen := jwt.NewTokenGenerator("test-secret", "test-refresh-secret", 15, 336)

	request := &service.BatchFeedbackRequest{
		Items: []service.BatchFeedbackItem{},
	}

	router := setupFeedbackRouterForBatch(handler, tokenGen)

	requestBody, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", "/api/v1/feedback/batch", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tokenGen.GenerateAccessToken("user-123", []string{}))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateBatchFeedback_InvalidJSON(t *testing.T) {
	mockService := new(MockFeedbackServiceForBatch)
	handler := NewFeedbackHandler(mockService)
	tokenGen := jwt.NewTokenGenerator("test-secret", "test-refresh-secret", 15, 336)

	router := setupFeedbackRouterForBatch(handler, tokenGen)
	req, _ := http.NewRequest("POST", "/api/v1/feedback/batch", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tokenGen.GenerateAccessToken("user-123", []string{}))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateBatchFeedback_ServiceError(t *testing.T) {
	mockService := new(MockFeedbackServiceForBatch)
	handler := NewFeedbackHandler(mockService)
	tokenGen := jwt.NewTokenGenerator("test-secret", "test-refresh-secret", 15, 336)

	request := &service.BatchFeedbackRequest{
		Items: []service.BatchFeedbackItem{
			{Content: "Test feedback"},
		},
	}

	mockService.On("CreateBatchFeedback", mock.Anything, "user-123", mock.Anything).Return(nil, assert.AnError)

	router := setupFeedbackRouterForBatch(handler, tokenGen)

	requestBody, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", "/api/v1/feedback/batch", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tokenGen.GenerateAccessToken("user-123", []string{}))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

func TestCreateBatchFeedback_Unauthorized(t *testing.T) {
	mockService := new(MockFeedbackServiceForBatch)
	handler := NewFeedbackHandler(mockService)
	tokenGen := jwt.NewTokenGenerator("test-secret", "test-refresh-secret", 15, 336)

	router := setupFeedbackRouterForBatch(handler, tokenGen)
	req, _ := http.NewRequest("POST", "/api/v1/feedback/batch", bytes.NewBufferString("{}"))
	req.Header.Set("Content-Type", "application/json")
	// No Authorization header
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
