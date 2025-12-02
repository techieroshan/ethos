package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupCommunityRouter(handler *CommunityHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	v1 := router.Group("/api/v1")
	community := v1.Group("/community")
	community.GET("/rules", handler.GetRules)
	return router
}

func TestGetRules_Success(t *testing.T) {
	handler := NewCommunityHandler()
	router := setupCommunityRouter(handler)
	req, _ := http.NewRequest("GET", "/api/v1/community/rules", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Community Rules", response["title"])
	assert.NotEmpty(t, response["content"])
}

func TestGetRules_NoAuthRequired(t *testing.T) {
	handler := NewCommunityHandler()
	router := setupCommunityRouter(handler)
	// No auth token required
	req, _ := http.NewRequest("GET", "/api/v1/community/rules", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

