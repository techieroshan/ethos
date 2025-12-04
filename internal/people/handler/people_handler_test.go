package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	authModel "ethos/internal/auth/model"
	"ethos/internal/people"
	"ethos/pkg/jwt"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockPeopleService is a mock implementation of the people service
type MockPeopleService struct {
	mock.Mock
}

func (m *MockPeopleService) SearchPeople(ctx context.Context, query string, limit, offset int) ([]*authModel.UserProfile, int, error) {
	args := m.Called(ctx, query, limit, offset)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*authModel.UserProfile), args.Int(1), args.Error(2)
}

func (m *MockPeopleService) GetRecommendations(ctx context.Context, userID string) ([]*authModel.UserProfile, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*authModel.UserProfile), args.Error(1)
}

func (m *MockPeopleService) SearchPeopleWithFilters(ctx context.Context, query string, limit, offset int, filters *people.PeopleSearchFilters) ([]*authModel.UserProfile, int, error) {
	args := m.Called(ctx, query, limit, offset, filters)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*authModel.UserProfile), args.Int(1), args.Error(2)
}

func setupPeopleRouter(handler *PeopleHandler, _ *jwt.TokenGenerator) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	v1 := router.Group("/api/v1")
	people := v1.Group("/people")
	people.Use(func(c *gin.Context) {
		c.Set("user_id", "test-user-id")
		c.Next()
	})
	people.GET("/search", handler.SearchPeople)
	people.GET("/recommendations", handler.GetRecommendations)
	return router
}

func TestSearchPeople_Success(t *testing.T) {
	mockService := new(MockPeopleService)
	handler := NewPeopleHandler(mockService)
	tokenGen := jwt.NewTokenGenerator("test-secret", "test-refresh-secret", 15, 336)

	profiles := []*authModel.UserProfile{
		{
			ID:            "user-1",
			Name:          "John Doe",
			Email:         "john@example.com",
			EmailVerified: true,
		},
		{
			ID:            "user-2",
			Name:          "Jane Smith",
			Email:         "jane@example.com",
			EmailVerified: true,
		},
	}

	mockService.On("SearchPeople", mock.Anything, "john", 25, 0).Return(profiles, 2, nil)

	router := setupPeopleRouter(handler, tokenGen)
	req, _ := http.NewRequest("GET", "/api/v1/people/search?q=john", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(2), response["count"])
	mockService.AssertExpectations(t)
}

func TestSearchPeople_WithPagination(t *testing.T) {
	mockService := new(MockPeopleService)
	handler := NewPeopleHandler(mockService)
	tokenGen := jwt.NewTokenGenerator("test-secret", "test-refresh-secret", 15, 336)

	profiles := []*authModel.UserProfile{}

	mockService.On("SearchPeople", mock.Anything, "test", 10, 5).Return(profiles, 0, nil)

	router := setupPeopleRouter(handler, tokenGen)
	req, _ := http.NewRequest("GET", "/api/v1/people/search?q=test&limit=10&offset=5", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestSearchPeople_NoResults(t *testing.T) {
	mockService := new(MockPeopleService)
	handler := NewPeopleHandler(mockService)
	tokenGen := jwt.NewTokenGenerator("test-secret", "test-refresh-secret", 15, 336)

	mockService.On("SearchPeople", mock.Anything, "nonexistent", 25, 0).Return([]*authModel.UserProfile{}, 0, nil)

	router := setupPeopleRouter(handler, tokenGen)
	req, _ := http.NewRequest("GET", "/api/v1/people/search?q=nonexistent", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["count"])
	mockService.AssertExpectations(t)
}

func TestGetRecommendations_Success(t *testing.T) {
	mockService := new(MockPeopleService)
	handler := NewPeopleHandler(mockService)
	tokenGen := jwt.NewTokenGenerator("test-secret", "test-refresh-secret", 15, 336)

	recommendations := []*authModel.UserProfile{
		{
			ID:            "user-3",
			Name:          "Alice Johnson",
			Email:         "alice@example.com",
			EmailVerified: true,
		},
	}

	mockService.On("GetRecommendations", mock.Anything, "test-user-id").Return(recommendations, nil)

	router := setupPeopleRouter(handler, tokenGen)
	req, _ := http.NewRequest("GET", "/api/v1/people/recommendations", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotNil(t, response["recommendations"])
	mockService.AssertExpectations(t)
}

func TestGetRecommendations_EmptyList(t *testing.T) {
	mockService := new(MockPeopleService)
	handler := NewPeopleHandler(mockService)
	tokenGen := jwt.NewTokenGenerator("test-secret", "test-refresh-secret", 15, 336)

	mockService.On("GetRecommendations", mock.Anything, "test-user-id").Return([]*authModel.UserProfile{}, nil)

	router := setupPeopleRouter(handler, tokenGen)
	req, _ := http.NewRequest("GET", "/api/v1/people/recommendations", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	recommendations := response["recommendations"].([]interface{})
	assert.Equal(t, 0, len(recommendations))
	mockService.AssertExpectations(t)
}
