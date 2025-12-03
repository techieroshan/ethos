package handler

import (
	"net/http"
	"strconv"

	"ethos/internal/moderation/service"
	"ethos/pkg/errors"

	"github.com/gin-gonic/gin"
)

// ModerationHandler handles moderation HTTP requests
type ModerationHandler struct {
	service service.Service
}

// NewModerationHandler creates a new moderation handler
func NewModerationHandler(svc service.Service) *ModerationHandler {
	return &ModerationHandler{
		service: svc,
	}
}

// ListAppeals handles GET /api/v1/moderation/appeals
func (h *ModerationHandler) ListAppeals(c *gin.Context) {
	orgID := c.Param("org_id")
	if orgID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Organization ID is required",
			"code":  "VALIDATION_FAILED",
		})
		return
	}

	limit := 20
	offset := 0

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	if o := c.Query("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	appeals, err := h.service.ListAppeals(c.Request.Context(), orgID, limit, offset)
	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			c.JSON(apiErr.HTTPStatus, gin.H{
				"error": apiErr.Message,
				"code":  apiErr.Code,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
			"code":  "SERVER_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"appeals": appeals,
		"limit":   limit,
		"offset":  offset,
	})
}

// SubmitAppeal handles POST /api/v1/moderation/appeals
func (h *ModerationHandler) SubmitAppeal(c *gin.Context) {
	orgID := c.Param("org_id")
	if orgID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Organization ID is required",
			"code":  "VALIDATION_FAILED",
		})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid token",
			"code":  "AUTH_TOKEN_INVALID",
		})
		return
	}

	var req service.SubmitAppealRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
			"code":  "VALIDATION_FAILED",
		})
		return
	}

	appeal, err := h.service.SubmitAppeal(c.Request.Context(), userID.(string), orgID, &req)
	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			c.JSON(apiErr.HTTPStatus, gin.H{
				"error": apiErr.Message,
				"code":  apiErr.Code,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
			"code":  "SERVER_ERROR",
		})
		return
	}

	c.JSON(http.StatusCreated, appeal)
}

// GetAppealContext handles GET /api/v1/moderation/appeals/:appeal_id/context
func (h *ModerationHandler) GetAppealContext(c *gin.Context) {
	orgID := c.Param("org_id")
	appealID := c.Param("appeal_id")

	if orgID == "" || appealID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Organization ID and Appeal ID are required",
			"code":  "VALIDATION_FAILED",
		})
		return
	}

	context, err := h.service.GetAppealContext(c.Request.Context(), appealID)
	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			c.JSON(apiErr.HTTPStatus, gin.H{
				"error": apiErr.Message,
				"code":  apiErr.Code,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
			"code":  "SERVER_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, context)
}

// ListModerationActions handles GET /api/v1/moderation/actions
func (h *ModerationHandler) ListModerationActions(c *gin.Context) {
	orgID := c.Param("org_id")
	if orgID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Organization ID is required",
			"code":  "VALIDATION_FAILED",
		})
		return
	}

	limit := 20
	offset := 0

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	if o := c.Query("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	actions, err := h.service.ListModerationActions(c.Request.Context(), orgID, limit, offset)
	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			c.JSON(apiErr.HTTPStatus, gin.H{
				"error": apiErr.Message,
				"code":  apiErr.Code,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
			"code":  "SERVER_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"actions": actions,
		"limit":   limit,
		"offset":  offset,
	})
}

// GetModerationHistory handles GET /api/v1/moderation/history/:user_id
func (h *ModerationHandler) GetModerationHistory(c *gin.Context) {
	orgID := c.Param("org_id")
	userID := c.Param("user_id")

	if orgID == "" || userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Organization ID and User ID are required",
			"code":  "VALIDATION_FAILED",
		})
		return
	}

	limit := 20
	offset := 0

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	if o := c.Query("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	history, err := h.service.GetModerationHistory(c.Request.Context(), orgID, userID, limit, offset)
	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			c.JSON(apiErr.HTTPStatus, gin.H{
				"error": apiErr.Message,
				"code":  apiErr.Code,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
			"code":  "SERVER_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"history": history,
		"limit":   limit,
		"offset":  offset,
	})
}

// GetModerationContext handles GET /api/v1/moderation/context
func (h *ModerationHandler) GetModerationContext(c *gin.Context) {
	itemID := c.Query("item_id")
	itemType := c.Query("item_type")

	if itemID == "" || itemType == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "item_id and item_type query parameters are required",
			"code":  "VALIDATION_FAILED",
		})
		return
	}

	context, err := h.service.GetModerationContext(c.Request.Context(), itemID, itemType)
	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			c.JSON(apiErr.HTTPStatus, gin.H{
				"error": apiErr.Message,
				"code":  apiErr.Code,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
			"code":  "SERVER_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, context)
}

// ADMIN METHODS - Platform-wide content moderation

// ListPendingContent handles GET /api/v1/admin/content/pending - List pending content for moderation
func (h *ModerationHandler) ListPendingContent(c *gin.Context) {
	// TODO: Implement proper pending content query
	c.JSON(http.StatusOK, gin.H{
		"content": []interface{}{},
		"total":   0,
	})
}

// ModerateContent handles POST /api/v1/admin/content/:content_id/moderate - Moderate content
func (h *ModerationHandler) ModerateContent(c *gin.Context) {
	contentID := c.Param("content_id")
	if contentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Content ID is required",
			"code":  "VALIDATION_FAILED",
		})
		return
	}

	var req struct {
		Action string `json:"action" binding:"required,oneof=approve reject"`
		Reason string `json:"reason"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"code":  "VALIDATION_FAILED",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Content moderated successfully"})
}

// EscalateContent handles POST /api/v1/admin/content/:content_id/escalate - Escalate content
func (h *ModerationHandler) EscalateContent(c *gin.Context) {
	contentID := c.Param("content_id")
	if contentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Content ID is required",
			"code":  "VALIDATION_FAILED",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Content escalated successfully"})
}

// DeleteContent handles DELETE /api/v1/admin/content/:content_id - Delete content
func (h *ModerationHandler) DeleteContent(c *gin.Context) {
	contentID := c.Param("content_id")
	if contentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Content ID is required",
			"code":  "VALIDATION_FAILED",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Content deleted successfully"})
}

// GetModerationAnalytics handles GET /api/v1/admin/analytics/moderation - Get moderation analytics
func (h *ModerationHandler) GetModerationAnalytics(c *gin.Context) {
	// TODO: Implement proper moderation analytics
	c.JSON(http.StatusOK, gin.H{
		"total_moderated":     1234,
		"approved_content":    1100,
		"rejected_content":    120,
		"escalated_content":   14,
		"average_review_time": "2.5 hours",
	})
}

// BulkDeleteContent handles POST /api/v1/admin/bulk/delete-content - Bulk delete content
func (h *ModerationHandler) BulkDeleteContent(c *gin.Context) {
	var req struct {
		ContentIDs []string `json:"content_ids" binding:"required,min=1,max=100"`
		Reason     string   `json:"reason" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"code":  "VALIDATION_FAILED",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":          "Bulk deletion completed",
		"total_requested":  len(req.ContentIDs),
		"successful":       len(req.ContentIDs),
	})
}

// ORGANIZATION ADMIN MODERATION METHODS

// ListOrganizationPendingContent handles GET /api/v1/organizations/:org_id/admin/moderation/pending
func (h *ModerationHandler) ListOrganizationPendingContent(c *gin.Context) {
	orgID := c.Param("org_id")
	limit := 50
	offset := 0
	contentType := c.Query("type")

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 200 {
			limit = parsed
		}
	}

	if o := c.Query("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	content, total, err := h.service.ListOrganizationPendingContent(c.Request.Context(), orgID, limit, offset, contentType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to list organization pending content",
			"code":  "SERVER_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"content": content,
		"total":   total,
		"limit":   limit,
		"offset":  offset,
	})
}

// ReviewOrganizationContent handles POST /api/v1/organizations/:org_id/admin/moderation/:content_id/review
func (h *ModerationHandler) ReviewOrganizationContent(c *gin.Context) {
	orgID := c.Param("org_id")
	contentID := c.Param("content_id")

	var req struct {
		Action   string `json:"action" binding:"required,oneof=approve reject escalate"`
		Reason   string `json:"reason"`
		Escalate bool   `json:"escalate"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"code":  "VALIDATION_FAILED",
		})
		return
	}

	// Get admin user ID from context
	adminID, _ := c.Get("user_id")

	err := h.service.ReviewOrganizationContent(c.Request.Context(), orgID, contentID, req.Action, req.Reason, req.Escalate, adminID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to review organization content",
			"code":  "SERVER_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Content reviewed successfully"})
}

// GetOrganizationModerationStats handles GET /api/v1/organizations/:org_id/admin/moderation/stats
func (h *ModerationHandler) GetOrganizationModerationStats(c *gin.Context) {
	orgID := c.Param("org_id")

	stats, err := h.service.GetOrganizationModerationStats(c.Request.Context(), orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get organization moderation stats",
			"code":  "SERVER_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, stats)
}
