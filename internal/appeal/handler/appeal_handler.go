package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"ethos/internal/appeal/service"
	"ethos/pkg/errors"
)

// AppealHandler handles appeal HTTP requests
type AppealHandler struct {
	service service.Service
}

// NewAppealHandler creates a new appeal handler
func NewAppealHandler(svc service.Service) *AppealHandler {
	return &AppealHandler{
		service: svc,
	}
}

// SubmitAppeal handles POST /api/v1/appeals
func (h *AppealHandler) SubmitAppeal(c *gin.Context) {
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
			"error": "Validation failed",
			"code":  "VALIDATION_FAILED",
		})
		return
	}

	appeal, err := h.service.SubmitAppeal(c.Request.Context(), userID.(string), &req)
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

	c.JSON(http.StatusCreated, gin.H{"appeal": appeal})
}

// GetUserAppeals handles GET /api/v1/appeals
func (h *AppealHandler) GetUserAppeals(c *gin.Context) {
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

	appeals, count, err := h.service.GetUserAppeals(c.Request.Context(), userID.(string), limitInt, offsetInt)
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
		"count":   count,
	})
}

// GetAppealByID handles GET /api/v1/appeals/:appeal_id
func (h *AppealHandler) GetAppealByID(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid token",
			"code":  "AUTH_TOKEN_INVALID",
		})
		return
	}

	appealID := c.Param("appeal_id")
	if appealID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Appeal ID is required",
			"code":  "VALIDATION_FAILED",
		})
		return
	}

	appeal, err := h.service.GetAppealByID(c.Request.Context(), userID.(string), appealID)
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

	c.JSON(http.StatusOK, gin.H{"appeal": appeal})
}
