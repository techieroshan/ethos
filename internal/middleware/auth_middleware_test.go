package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"ethos/pkg/jwt"
)

func setupTestRouter(middleware gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware)
	router.GET("/protected", func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user_id not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"user_id": userID})
	})
	return router
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	// Create a token generator with test secrets
	tokenGen := jwt.NewTokenGenerator(
		"test-access-secret",
		"test-refresh-secret",
		15*60*time.Second, // 15 minutes
		14*24*60*60*time.Second, // 14 days
	)

	// Generate a valid token
	userID := "user-123"
	token, err := tokenGen.GenerateAccessToken(userID)
	assert.NoError(t, err)

	// Create middleware
	middleware := AuthMiddleware(tokenGen)

	// Setup router
	router := setupTestRouter(middleware)

	// Make request with valid token
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	tokenGen := jwt.NewTokenGenerator(
		"test-access-secret",
		"test-refresh-secret",
		15*60*time.Second,
		14*24*60*60*time.Second,
	)

	middleware := AuthMiddleware(tokenGen)
	router := setupTestRouter(middleware)

	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "AUTH_TOKEN_INVALID", response["code"])
}

func TestAuthMiddleware_MissingToken(t *testing.T) {
	tokenGen := jwt.NewTokenGenerator(
		"test-access-secret",
		"test-refresh-secret",
		15*60*time.Second,
		14*24*60*60*time.Second,
	)

	middleware := AuthMiddleware(tokenGen)
	router := setupTestRouter(middleware)

	req, _ := http.NewRequest("GET", "/protected", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "AUTH_TOKEN_INVALID", response["code"])
}

func TestAuthMiddleware_MalformedHeader(t *testing.T) {
	tokenGen := jwt.NewTokenGenerator(
		"test-access-secret",
		"test-refresh-secret",
		15*60*time.Second,
		14*24*60*60*time.Second,
	)

	middleware := AuthMiddleware(tokenGen)
	router := setupTestRouter(middleware)

	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "InvalidFormat token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

