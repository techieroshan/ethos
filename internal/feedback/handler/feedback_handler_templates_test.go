package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	fbModel "ethos/internal/feedback/model"
	"ethos/pkg/jwt"
)

// MockFeedbackServiceForTemplates is a mock implementation of the feedback service for templates
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

func (m *MockFeedbackServiceForTemplates) CreateFeedback(ctx context.Context, userID string, req interface{}) (*fbModel.FeedbackItem, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*fbModel.FeedbackItem), args.Error(1)
}

func (m *MockFeedbackServiceForTemplates) CreateComment(ctx context.Context, userID, feedbackID string, req interface{}) (*fbModel.FeedbackComment, error) {
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

func (m *MockFeedbackServiceForTemplates) SubmitTemplateSuggestion(ctx context.Context, req interface{}) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func setupFeedbackRouterForTemplates(handler *FeedbackHandler, tokenGen *jwt.TokenGenerator) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	v1 := router.Group("/api/v1")
	feedback := v1.Group("/feedback")
	feedback.GET("/templates", handler.GetTemplates)
	feedback.POST("/template_suggestions", handler.PostTemplateSuggestions)
	return router
}

func TestGetTemplates_Success(t *testing.T) {
	mockService := new(MockFeedbackServiceForTemplates)
	handler := NewFeedbackHandler(mockService)
	tokenGen := jwt.NewTokenGenerator("test-secret", "test-refresh-secret", 15, 336)

	templates := []*fbModel.FeedbackTemplate{
		{
			TemplateID:  "t-001",
			Name:        "Appreciation Template",
			Description: "A short message to acknowledge great work.",
			ContextTags: []string{"appreciation", "general"},
			TemplateFields: map[string]interface{}{
				"fields": []interface{}{
					map[string]interface{}{
						"name":  "what_went_well",
						"type":  "text",
						"label": "What went well?",
					},
				},
			},
		},
	}

	mockService.On("GetTemplates", mock.Anything, "", "").Return(templates, nil)

	router := setupFeedbackRouterForTemplates(handler, tokenGen)
	req, _ := http.NewRequest("GET", "/api/v1/feedback/templates", nil)
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

func TestGetTemplates_WithContextFilter(t *testing.T) {
	mockService := new(MockFeedbackServiceForTemplates)
	handler := NewFeedbackHandler(mockService)
	tokenGen := jwt.NewTokenGenerator("test-secret", "test-refresh-secret", 15, 336)

	mockService.On("GetTemplates", mock.Anything, "performance_review", "").Return([]*model.FeedbackTemplate{}, nil)

	router := setupFeedbackRouterForTemplates(handler, tokenGen)
	req, _ := http.NewRequest("GET", "/api/v1/feedback/templates?context=performance_review", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestGetTemplates_WithTagsFilter(t *testing.T) {
	mockService := new(MockFeedbackServiceForTemplates)
	handler := NewFeedbackHandler(mockService)
	tokenGen := jwt.NewTokenGenerator("test-secret", "test-refresh-secret", 15, 336)

	mockService.On("GetTemplates", mock.Anything, "", "leadership,initiative").Return([]*model.FeedbackTemplate{}, nil)

	router := setupFeedbackRouterForTemplates(handler, tokenGen)
	req, _ := http.NewRequest("GET", "/api/v1/feedback/templates?tags=leadership,initiative", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestGetTemplates_ServiceError(t *testing.T) {
	mockService := new(MockFeedbackServiceForTemplates)
	handler := NewFeedbackHandler(mockService)
	tokenGen := jwt.NewTokenGenerator("test-secret", "test-refresh-secret", 15, 336)

	mockService.On("GetTemplates", mock.Anything, "", "").Return(nil, assert.AnError)

	router := setupFeedbackRouterForTemplates(handler, tokenGen)
	req, _ := http.NewRequest("GET", "/api/v1/feedback/templates", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

func TestPostTemplateSuggestions_Success(t *testing.T) {
	mockService := new(MockFeedbackServiceForTemplates)
	handler := NewFeedbackHandler(mockService)
	tokenGen := jwt.NewTokenGenerator("test-secret", "test-refresh-secret", 15, 336)

	mockService.On("SubmitTemplateSuggestion", mock.Anything, mock.Anything).Return(nil)

	router := setupFeedbackRouterForTemplates(handler, tokenGen)
	suggestion := `{
		"usage_context": "peer review",
		"details": "Need a short, positive feedback template for quick peer recognition.",
		"desired_fields": [{"name": "positive_note", "type": "text"}]
	}`
	req, _ := http.NewRequest("POST", "/api/v1/feedback/template_suggestions", strings.NewReader(suggestion))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "suggestion_received", response["status"])
	mockService.AssertExpectations(t)
}

func TestPostTemplateSuggestions_InvalidJSON(t *testing.T) {
	mockService := new(MockFeedbackServiceForTemplates)
	handler := NewFeedbackHandler(mockService)
	tokenGen := jwt.NewTokenGenerator("test-secret", "test-refresh-secret", 15, 336)

	router := setupFeedbackRouterForTemplates(handler, tokenGen)
	req, _ := http.NewRequest("POST", "/api/v1/feedback/template_suggestions", strings.NewReader("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
