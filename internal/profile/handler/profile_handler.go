package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"ethos/internal/profile/service"
	"ethos/pkg/errors"
)

// ProfileHandler handles profile HTTP requests
type ProfileHandler struct {
	service service.Service
}

// NewProfileHandler creates a new profile handler
func NewProfileHandler(svc service.Service) *ProfileHandler {
	return &ProfileHandler{
		service: svc,
	}
}

// GetProfile handles GET /api/v1/profile/me
func (h *ProfileHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid token",
			"code":  "AUTH_TOKEN_INVALID",
		})
		return
	}

	profile, err := h.service.GetProfile(c.Request.Context(), userID.(string))
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

	c.JSON(http.StatusOK, profile)
}

// UpdateProfile handles PUT /api/v1/profile/me
func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid token",
			"code":  "AUTH_TOKEN_INVALID",
		})
		return
	}

	var req service.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Validation failed",
			"code":  "VALIDATION_FAILED",
		})
		return
	}

	// Validate that at least one field is provided
	if req.Name == "" && req.PublicBio == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "At least one field must be provided",
			"code":  "VALIDATION_FAILED",
		})
		return
	}

	profile, err := h.service.UpdateProfile(c.Request.Context(), userID.(string), &req)
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

	c.JSON(http.StatusOK, profile)
}

// UpdatePreferences handles PATCH /api/v1/profile/me/preferences
func (h *ProfileHandler) UpdatePreferences(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid token",
			"code":  "AUTH_TOKEN_INVALID",
		})
		return
	}

	var req service.UpdatePreferencesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Validation failed",
			"code":  "VALIDATION_FAILED",
		})
		return
	}

	prefs, err := h.service.UpdatePreferences(c.Request.Context(), userID.(string), &req)
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

	c.JSON(http.StatusOK, gin.H{"preferences": prefs})
}

// DeleteProfile handles DELETE /api/v1/profile/me
func (h *ProfileHandler) DeleteProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid token",
			"code":  "AUTH_TOKEN_INVALID",
		})
		return
	}

	err := h.service.DeleteProfile(c.Request.Context(), userID.(string))
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
		"message": "Account scheduled for deletion. You will receive a confirmation email.",
	})
}

// SearchProfiles handles GET /api/v1/profile/user-profile
func (h *ProfileHandler) SearchProfiles(c *gin.Context) {
	query := c.Query("q")
	limit := c.DefaultQuery("limit", "25")
	offset := c.DefaultQuery("offset", "0")

	limitInt := 25
	offsetInt := 0
	if l, err := strconv.Atoi(limit); err == nil && l > 0 {
		limitInt = l
	}
	if o, err := strconv.Atoi(offset); err == nil && o >= 0 {
		offsetInt = o
	}

	profiles, count, err := h.service.SearchProfiles(c.Request.Context(), query, limitInt, offsetInt)
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
		"results": profiles,
		"count":   count,
	})
}
