package handler

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
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

// MockFeedbackServiceForExport is a mock implementation for export tests
type MockFeedbackServiceForExport struct {
	mock.Mock
}

func (m *MockFeedbackServiceForExport) GetFeed(ctx context.Context, limit, offset int) ([]*fbModel.FeedbackItem, int, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*fbModel.FeedbackItem), args.Get(1).(int), args.Error(2)
}

func (m *MockFeedbackServiceForExport) GetFeedbackByID(ctx context.Context, feedbackID string) (*fbModel.FeedbackItem, error) {
	args := m.Called(ctx, feedbackID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*fbModel.FeedbackItem), args.Error(1)
}

func (m *MockFeedbackServiceForExport) GetComments(ctx context.Context, feedbackID string, limit, offset int) ([]*fbModel.FeedbackComment, int, error) {
	args := m.Called(ctx, feedbackID, limit, offset)
	if args.Get(0) == nil {
		return nil, 0, args.Error(3)
	}
	return args.Get(0).([]*fbModel.FeedbackComment), args.Get(1).(int), args.Error(3)
}

func (m *MockFeedbackServiceForExport) CreateFeedback(ctx context.Context, userID string, req interface{}) (*fbModel.FeedbackItem, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*fbModel.FeedbackItem), args.Error(1)
}

func (m *MockFeedbackServiceForExport) CreateComment(ctx context.Context, userID, feedbackID string, req interface{}) (*fbModel.FeedbackComment, error) {
	args := m.Called(ctx, userID, feedbackID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*fbModel.FeedbackComment), args.Error(1)
}

func (m *MockFeedbackServiceForExport) AddReaction(ctx context.Context, userID, feedbackID string, reactionType string) error {
	args := m.Called(ctx, userID, feedbackID, reactionType)
	return args.Error(0)
}

func (m *MockFeedbackServiceForExport) RemoveReaction(ctx context.Context, userID, feedbackID string, reactionType string) error {
	args := m.Called(ctx, userID, feedbackID, reactionType)
	return args.Error(0)
}

func (m *MockFeedbackServiceForExport) GetTemplates(ctx context.Context, contextFilter, tagsFilter string) ([]*fbModel.FeedbackTemplate, error) {
	args := m.Called(ctx, contextFilter, tagsFilter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*fbModel.FeedbackTemplate), args.Error(1)
}

func (m *MockFeedbackServiceForExport) SubmitTemplateSuggestion(ctx context.Context, req interface{}) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *MockFeedbackServiceForExport) GetImpact(ctx context.Context, userID *string, from, to *time.Time) (*fbModel.FeedbackImpact, error) {
	args := m.Called(ctx, userID, from, to)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*fbModel.FeedbackImpact), args.Error(1)
}

func (m *MockFeedbackServiceForExport) CreateBatchFeedback(ctx context.Context, userID string, req *service.BatchFeedbackRequest) (*service.BatchFeedbackResponse, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.BatchFeedbackResponse), args.Error(1)
}

func (m *MockFeedbackServiceForExport) GetFeedWithFilters(ctx context.Context, limit, offset int, filters *service.FeedFilters) ([]*fbModel.FeedbackItem, int, error) {
	args := m.Called(ctx, limit, offset, filters)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*fbModel.FeedbackItem), args.Get(1).(int), args.Error(2)
}

func (m *MockFeedbackServiceForExport) GetBookmarks(ctx context.Context, userID string, limit, offset int) ([]*fbModel.FeedbackItem, int, error) {
	args := m.Called(ctx, userID, limit, offset)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*fbModel.FeedbackItem), args.Get(1).(int), args.Error(2)
}

func (m *MockFeedbackServiceForExport) AddBookmark(ctx context.Context, userID, feedbackID string) error {
	args := m.Called(ctx, userID, feedbackID)
	return args.Error(0)
}

func (m *MockFeedbackServiceForExport) RemoveBookmark(ctx context.Context, userID, feedbackID string) error {
	args := m.Called(ctx, userID, feedbackID)
	return args.Error(0)
}

func (m *MockFeedbackServiceForExport) ExportFeedback(ctx context.Context, filters *service.FeedFilters, format string) (*service.ExportResponse, error) {
	args := m.Called(ctx, filters, format)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.ExportResponse), args.Error(1)
}

func setupFeedbackRouterForExport(handler *FeedbackHandler, tokenGen *jwt.TokenGenerator) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.AuthMiddleware(tokenGen))
	v1 := router.Group("/api/v1")
	feedback := v1.Group("/feedback")
	feedback.GET("/export", handler.ExportFeedback)
	return router
}

func TestExportFeedback_JSONFormat_Success(t *testing.T) {
	mockService := new(MockFeedbackServiceForExport)
	handler := NewFeedbackHandler(mockService)
	tokenGen := jwt.NewTokenGenerator("test-secret", "test-refresh-secret", 15, 336)

	filters := &service.FeedFilters{
		ReviewerType: stringPtr("org"),
	}
	exportResponse := &service.ExportResponse{
		Format:      "json",
		ContentType: "application/json",
		Data:        `[{"feedback_id":"f-001","content":"Test feedback"}]`,
		Count:       1,
	}

	mockService.On("ExportFeedback", mock.Anything, filters, "json").Return(exportResponse, nil)

	router := setupFeedbackRouterForExport(handler, tokenGen)
	req, _ := http.NewRequest("GET", "/api/v1/feedback/export?reviewer_type=org&format=json", nil)
	req.Header.Set("Authorization", "Bearer "+tokenGen.GenerateAccessToken("user-123", []string{}))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Equal(t, "attachment; filename=feedback_export.json", w.Header().Get("Content-Disposition"))

	var response service.ExportResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "json", response.Format)
	assert.Equal(t, 1, response.Count)
	mockService.AssertExpectations(t)
}

func TestExportFeedback_CSVFormat_Success(t *testing.T) {
	mockService := new(MockFeedbackServiceForExport)
	handler := NewFeedbackHandler(mockService)
	tokenGen := jwt.NewTokenGenerator("test-secret", "test-refresh-secret", 15, 336)

	filters := &service.FeedFilters{}
	exportResponse := &service.ExportResponse{
		Format:      "csv",
		ContentType: "text/csv",
		Data:        "feedback_id,content,author_name\nf-001,Test feedback,John Doe\n",
		Count:       1,
	}

	mockService.On("ExportFeedback", mock.Anything, filters, "csv").Return(exportResponse, nil)

	router := setupFeedbackRouterForExport(handler, tokenGen)
	req, _ := http.NewRequest("GET", "/api/v1/feedback/export?format=csv", nil)
	req.Header.Set("Authorization", "Bearer "+tokenGen.GenerateAccessToken("user-123", []string{}))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "text/csv", w.Header().Get("Content-Type"))
	assert.Equal(t, "attachment; filename=feedback_export.csv", w.Header().Get("Content-Disposition"))
	mockService.AssertExpectations(t)
}

func TestExportFeedback_WithAllFilters(t *testing.T) {
	mockService := new(MockFeedbackServiceForExport)
	handler := NewFeedbackHandler(mockService)
	tokenGen := jwt.NewTokenGenerator("test-secret", "test-refresh-secret", 15, 336)

	filters := &service.FeedFilters{
		ReviewerType: stringPtr("org"),
		Context:      stringPtr("project"),
		Verification: stringPtr("verified"),
		Tags:         []string{"leadership", "teamwork"},
	}

	exportResponse := &service.ExportResponse{
		Format:      "json",
		ContentType: "application/json",
		Data:        `[{"feedback_id":"f-001","content":"Filtered feedback"}]`,
		Count:       1,
	}

	mockService.On("ExportFeedback", mock.Anything, filters, "json").Return(exportResponse, nil)

	router := setupFeedbackRouterForExport(handler, tokenGen)
	req, _ := http.NewRequest("GET", "/api/v1/feedback/export?reviewer_type=org&context=project&verification=verified&tags=leadership,teamwork&format=json", nil)
	req.Header.Set("Authorization", "Bearer "+tokenGen.GenerateAccessToken("user-123", []string{}))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestExportFeedback_DefaultFormat(t *testing.T) {
	mockService := new(MockFeedbackServiceForExport)
	handler := NewFeedbackHandler(mockService)
	tokenGen := jwt.NewTokenGenerator("test-secret", "test-refresh-secret", 15, 336)

	filters := &service.FeedFilters{}
	exportResponse := &service.ExportResponse{
		Format:      "json",
		ContentType: "application/json",
		Data:        `[{"feedback_id":"f-001","content":"Default format"}]`,
		Count:       1,
	}

	mockService.On("ExportFeedback", mock.Anything, filters, "json").Return(exportResponse, nil)

	router := setupFeedbackRouterForExport(handler, tokenGen)
	req, _ := http.NewRequest("GET", "/api/v1/feedback/export", nil)
	req.Header.Set("Authorization", "Bearer "+tokenGen.GenerateAccessToken("user-123", []string{}))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	mockService.AssertExpectations(t)
}

func TestExportFeedback_InvalidFormat(t *testing.T) {
	mockService := new(MockFeedbackServiceForExport)
	handler := NewFeedbackHandler(mockService)
	tokenGen := jwt.NewTokenGenerator("test-secret", "test-refresh-secret", 15, 336)

	router := setupFeedbackRouterForExport(handler, tokenGen)
	req, _ := http.NewRequest("GET", "/api/v1/feedback/export?format=invalid", nil)
	req.Header.Set("Authorization", "Bearer "+tokenGen.GenerateAccessToken("user-123", []string{}))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "VALIDATION_FAILED", response["code"])
}

func TestExportFeedback_ServiceError(t *testing.T) {
	mockService := new(MockFeedbackServiceForExport)
	handler := NewFeedbackHandler(mockService)
	tokenGen := jwt.NewTokenGenerator("test-secret", "test-refresh-secret", 15, 336)

	mockService.On("ExportFeedback", mock.Anything, mock.Anything, "json").Return(nil, assert.AnError)

	router := setupFeedbackRouterForExport(handler, tokenGen)
	req, _ := http.NewRequest("GET", "/api/v1/feedback/export", nil)
	req.Header.Set("Authorization", "Bearer "+tokenGen.GenerateAccessToken("user-123", []string{}))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

func TestExportFeedback_Unauthorized(t *testing.T) {
	mockService := new(MockFeedbackServiceForExport)
	handler := NewFeedbackHandler(mockService)
	tokenGen := jwt.NewTokenGenerator("test-secret", "test-refresh-secret", 15, 336)

	router := setupFeedbackRouterForExport(handler, tokenGen)
	req, _ := http.NewRequest("GET", "/api/v1/feedback/export", nil)
	// No Authorization header
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestExportFeedback_EmptyResult(t *testing.T) {
	mockService := new(MockFeedbackServiceForExport)
	handler := NewFeedbackHandler(mockService)
	tokenGen := jwt.NewTokenGenerator("test-secret", "test-refresh-secret", 15, 336)

	filters := &service.FeedFilters{}
	exportResponse := &service.ExportResponse{
		Format:      "json",
		ContentType: "application/json",
		Data:        `[]`,
		Count:       0,
	}

	mockService.On("ExportFeedback", mock.Anything, filters, "json").Return(exportResponse, nil)

	router := setupFeedbackRouterForExport(handler, tokenGen)
	req, _ := http.NewRequest("GET", "/api/v1/feedback/export", nil)
	req.Header.Set("Authorization", "Bearer "+tokenGen.GenerateAccessToken("user-123", []string{}))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response service.ExportResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 0, response.Count)
	mockService.AssertExpectations(t)
}

// Helper function
func stringPtr(s string) *string {
	return &s
}
