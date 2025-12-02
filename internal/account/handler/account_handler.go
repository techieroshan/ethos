package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"ethos/internal/account/service"
	"ethos/pkg/errors"
)

// AccountHandler handles account/security HTTP requests
type AccountHandler struct {
	service service.Service
}

// NewAccountHandler creates a new account handler
func NewAccountHandler(svc service.Service) *AccountHandler {
	return &AccountHandler{
		service: svc,
	}
}

// GetSecurityEvents handles GET /api/v1/account/security-events
func (h *AccountHandler) GetSecurityEvents(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid token",
			"code":  "AUTH_TOKEN_INVALID",
		})
		return
	}

	limit := c.DefaultQuery("limit", "20")
	offset := c.DefaultQuery("offset", "0")

	limitInt := 20
	offsetInt := 0
	if l, err := strconv.Atoi(limit); err == nil && l > 0 {
		limitInt = l
	}
	if o, err := strconv.Atoi(offset); err == nil && o >= 0 {
		offsetInt = o
	}

	events, count, err := h.service.GetSecurityEvents(c.Request.Context(), userID.(string), limitInt, offsetInt)
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
		"events": events,
		"count":  count,
	})
}

// GetExportStatus handles GET /api/v1/account/export-data/:export_id/status
func (h *AccountHandler) GetExportStatus(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid token",
			"code":  "AUTH_TOKEN_INVALID",
		})
		return
	}

	exportID := c.Param("export_id")

	export, err := h.service.GetExportStatus(c.Request.Context(), userID.(string), exportID)
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

	c.JSON(http.StatusOK, export)
}

// Disable2FA handles DELETE /api/v1/auth/setup-2fa
func (h *AccountHandler) Disable2FA(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid token",
			"code":  "AUTH_TOKEN_INVALID",
		})
		return
	}

	err := h.service.Disable2FA(c.Request.Context(), userID.(string))
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
		"message": "Two-factor authentication disabled for your account.",
	})
}

