package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"ethos/internal/dashboard/service"
	"ethos/pkg/errors"
)

// DashboardHandler handles dashboard HTTP requests
type DashboardHandler struct {
	service service.Service
}

// NewDashboardHandler creates a new dashboard handler
func NewDashboardHandler(svc service.Service) *DashboardHandler {
	return &DashboardHandler{
		service: svc,
	}
}

// GetDashboard handles GET /api/v1/dashboard
func (h *DashboardHandler) GetDashboard(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid token",
			"code":  "AUTH_TOKEN_INVALID",
		})
		return
	}

	snapshot, err := h.service.GetDashboard(c.Request.Context(), userID.(string))
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

	c.JSON(http.StatusOK, snapshot)
}

