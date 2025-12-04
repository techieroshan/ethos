package handler

import (
	"net/http"
	"strconv"

	"ethos/internal/organization/repository"
	"ethos/internal/organization/service"

	"github.com/gin-gonic/gin"
)

// ContextSwitchHandler handles multi-tenant context switching endpoints
type ContextSwitchHandler struct {
	contextService service.UserContextService
}

// NewContextSwitchHandler creates a new context switch handler
func NewContextSwitchHandler(contextService service.UserContextService) *ContextSwitchHandler {
	return &ContextSwitchHandler{
		contextService: contextService,
	}
}

// SwitchContextRequest represents a request to switch context
type SwitchContextRequest struct {
	OrganizationID string `json:"organization_id" binding:"required,uuid"`
}

// SwitchContextResponse represents the response from switching context
type SwitchContextResponse struct {
	Context   *repository.UserContext `json:"context"`
	Message   string                  `json:"message"`
	Timestamp string                  `json:"timestamp"`
}

// GetAvailableContexts handles GET /api/v1/profile/available-contexts
// Returns all organizations a user belongs to
func (h *ContextSwitchHandler) GetAvailableContexts(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	contexts, err := h.contextService.GetAvailableContexts(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch available contexts"})
		return
	}

	currentContext, _ := h.contextService.GetCurrentContext(c.Request.Context(), userID)

	c.JSON(http.StatusOK, gin.H{
		"contexts": contexts,
		"current":  currentContext,
		"total":    len(contexts),
	})
}

// GetCurrentContext handles GET /api/v1/profile/current-context
// Returns the user's current organization context
func (h *ContextSwitchHandler) GetCurrentContext(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	currentContext, err := h.contextService.GetCurrentContext(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch current context"})
		return
	}

	if currentContext == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no current context found"})
		return
	}

	c.JSON(http.StatusOK, currentContext)
}

// SwitchContext handles POST /api/v1/profile/switch-context
// Switches the user to a different organization context
func (h *ContextSwitchHandler) SwitchContext(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req SwitchContextRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// Validate user is in the organization
	isMember, err := h.contextService.ValidateUserInOrganization(c.Request.Context(), userID, req.OrganizationID)
	if err != nil || !isMember {
		c.JSON(http.StatusForbidden, gin.H{"error": "user is not a member of this organization"})
		return
	}

	// Get IP address and user agent
	ipAddress := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	// Switch context
	newContext, err := h.contextService.SwitchContext(c.Request.Context(), userID, req.OrganizationID, ipAddress, userAgent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to switch context"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"context": newContext,
		"message": "context switched successfully",
	})
}

// GetContextSwitchHistory handles GET /api/v1/profile/context-switch-history
// Returns the user's context switch history
func (h *ContextSwitchHandler) GetContextSwitchHistory(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	limit := 50
	offset := 0
	if l := c.Query("limit"); l != "" {
		if n, err := parseIntQuery(l); err == nil && n > 0 {
			limit = n
		}
	}
	if o := c.Query("offset"); o != "" {
		if n, err := parseIntQuery(o); err == nil && n >= 0 {
			offset = n
		}
	}

	records, err := h.contextService.GetContextSwitchHistory(c.Request.Context(), userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch switch history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"records": records,
		"total":   len(records),
	})
}

// parseIntQuery is a helper to parse query parameters as integers
func parseIntQuery(s string) (int, error) {
	return strconv.Atoi(s)
}
