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
