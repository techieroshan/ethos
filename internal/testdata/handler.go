package testdata

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// Handler handles test data endpoints
type Handler struct {
	service *TestDataService
}

// NewHandler creates a new test data handler
func NewHandler() *Handler {
	return &Handler{
		service: NewTestDataService(),
	}
}

// RegisterRoutes registers test data routes
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	// Only register routes in dev/staging environments
	if !isDevEnvironment() {
		return
	}

	mux.HandleFunc("/api/v1/test-data/seed", h.SeedTestDataHTTP)
	mux.HandleFunc("/api/v1/test-data/cleanup", h.CleanupTestDataHTTP)
}

// SeedTestData creates all test data
// POST /api/v1/test-data/seed
func (h *Handler) SeedTestData(c *gin.Context) {
	// Security check: only dev/staging
	if !isDevEnvironment() {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Test data endpoints only available in development",
		})
		return
	}

	// Create test data
	result, err := h.service.CreateTestData(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Return success response
	c.JSON(http.StatusCreated, result)
}

// CleanupTestData removes all test data
// DELETE /api/v1/test-data/cleanup
func (h *Handler) CleanupTestData(c *gin.Context) {
	// Security check: only dev/staging
	if !isDevEnvironment() {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Test data endpoints only available in development",
		})
		return
	}

	// Cleanup test data
	result, err := h.service.CleanupTestData(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, result)
}

// SeedTestDataHTTP creates all test data (HTTP handler version)
// POST /api/v1/test-data/seed
func (h *Handler) SeedTestDataHTTP(w http.ResponseWriter, r *http.Request) {
	// Only allow POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Security check: only dev/staging
	if !isDevEnvironment() {
		http.Error(w, "Test data endpoints only available in development", http.StatusForbidden)
		return
	}

	// Create test data
	result, err := h.service.CreateTestData(r.Context())
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result)
}

// CleanupTestDataHTTP removes all test data (HTTP handler version)
// DELETE /api/v1/test-data/cleanup
func (h *Handler) CleanupTestDataHTTP(w http.ResponseWriter, r *http.Request) {
	// Only allow DELETE
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Security check: only dev/staging
	if !isDevEnvironment() {
		http.Error(w, "Test data endpoints only available in development", http.StatusForbidden)
		return
	}

	// Cleanup test data
	result, err := h.service.CleanupTestData(r.Context())
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

// isDevEnvironment checks if we're in a development environment
func isDevEnvironment() bool {
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = os.Getenv("ENV")
	}
	if env == "" {
		env = os.Getenv("NODE_ENV")
	}

	return env == "development" || env == "dev" || env == "staging" || env == "test"
}

