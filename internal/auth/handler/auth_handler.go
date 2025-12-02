package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"ethos/internal/auth/service"
	"ethos/pkg/errors"
)

// AuthHandler handles authentication HTTP requests
type AuthHandler struct {
	service service.Service
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(svc service.Service) *AuthHandler {
	return &AuthHandler{
		service: svc,
	}
}

// Login handles POST /api/v1/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Validation failed",
			"code":  "VALIDATION_FAILED",
		})
		return
	}

		resp, err := h.service.Login(c.Request.Context(), &req)
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

	c.JSON(http.StatusOK, resp)
}

// Register handles POST /api/v1/auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req service.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Validation failed",
			"code":  "VALIDATION_FAILED",
		})
		return
	}

	profile, err := h.service.Register(c.Request.Context(), &req)
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

	c.JSON(http.StatusCreated, profile)
}

// Refresh handles POST /api/v1/auth/refresh
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req service.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Validation failed",
			"code":  "VALIDATION_FAILED",
		})
		return
	}

	resp, err := h.service.RefreshToken(c.Request.Context(), &req)
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

	c.JSON(http.StatusOK, resp)
}

// Me handles GET /api/v1/auth/me
func (h *AuthHandler) Me(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid token",
			"code":  "AUTH_TOKEN_INVALID",
		})
		return
	}

	profile, err := h.service.GetUserProfile(c.Request.Context(), userID.(string))
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

