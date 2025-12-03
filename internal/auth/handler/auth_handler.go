package handler

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"ethos/internal/auth/service"
	"ethos/pkg/errors"

	"github.com/gin-gonic/gin"
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
		// DEBUG: Log binding error
		debugFile, _ := os.OpenFile("/tmp/ethos_debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if debugFile != nil {
			fmt.Fprintf(debugFile, "[%s] Binding Error: %v\n", time.Now().Format(time.RFC3339), err)
			debugFile.Close()
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"code":    "VALIDATION_FAILED",
			"details": err.Error(),
		})
		return
	}

	// DEBUG: Log received email to file
	debugFile, _ := os.OpenFile("/tmp/ethos_debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if debugFile != nil {
		fmt.Fprintf(debugFile, "[%s] Register Handler - Email: '%s' (len:%d)\n", time.Now().Format(time.RFC3339), req.Email, len(req.Email))
		debugFile.Close()
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

// RequestPasswordReset handles POST /api/v1/auth/request-password-reset
func (h *AuthHandler) RequestPasswordReset(c *gin.Context) {
	var req service.RequestPasswordResetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Validation failed",
			"code":  "VALIDATION_FAILED",
		})
		return
	}

	err := h.service.RequestPasswordReset(c.Request.Context(), &req)
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

	// Always return success for security (don't reveal if email exists)
	c.JSON(http.StatusOK, gin.H{
		"message": "If an account with this email exists, a password reset link has been sent.",
	})
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

// VerifyEmail handles GET /api/v1/auth/verify-email/:token
func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Verification token is required",
			"code":  "VALIDATION_FAILED",
		})
		return
	}

	err := h.service.VerifyEmail(c.Request.Context(), token)
	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			c.JSON(apiErr.HTTPStatus, gin.H{
				"error": apiErr.Message,
				"code":  apiErr.Code,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to verify email",
			"code":  "SERVER_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "email_verified",
		"message": "Your email has been successfully verified",
	})
}

// ChangePassword handles POST /api/v1/auth/change-password
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Authentication required",
			"code":  "AUTH_REQUIRED",
		})
		return
	}

	var req service.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Validation failed",
			"code":  "VALIDATION_FAILED",
		})
		return
	}

	err := h.service.ChangePassword(c.Request.Context(), userID.(string), &req)
	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			c.JSON(apiErr.HTTPStatus, gin.H{
				"error": apiErr.Message,
				"code":  apiErr.Code,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to change password",
			"code":  "SERVER_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "password_changed",
		"message": "Password has been successfully changed",
	})
}

// Setup2FA handles POST /api/v1/auth/setup-2fa
func (h *AuthHandler) Setup2FA(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Authentication required",
			"code":  "AUTH_REQUIRED",
		})
		return
	}

	var req service.Setup2FARequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Validation failed",
			"code":  "VALIDATION_FAILED",
		})
		return
	}

	resp, err := h.service.Setup2FA(c.Request.Context(), userID.(string), &req)
	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			c.JSON(apiErr.HTTPStatus, gin.H{
				"error": apiErr.Message,
				"code":  apiErr.Code,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to setup 2FA",
			"code":  "SERVER_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// MULTI-TENANT FUNCTIONALITY HANDLERS

// ListUserTenants handles GET /api/v1/auth/tenants
func (h *AuthHandler) ListUserTenants(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Authentication required",
			"code":  "AUTH_REQUIRED",
		})
		return
	}

	user, err := h.service.GetUserByID(c.Request.Context(), userID.(string))
	if err != nil {
		if err == errors.ErrUserNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User not found",
				"code":  "USER_NOT_FOUND",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to get user",
				"code":  "SERVER_ERROR",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tenants":           user.TenantMemberships,
		"current_tenant_id": user.CurrentTenantID,
	})
}

// SwitchTenant handles POST /api/v1/auth/tenants/:tenant_id/switch
func (h *AuthHandler) SwitchTenant(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Authentication required",
			"code":  "AUTH_REQUIRED",
		})
		return
	}

	tenantID := c.Param("tenant_id")
	if tenantID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Tenant ID is required",
			"code":  "VALIDATION_FAILED",
		})
		return
	}

	err := h.service.SwitchUserTenant(c.Request.Context(), userID.(string), tenantID)
	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			c.JSON(apiErr.HTTPStatus, gin.H{
				"error": apiErr.Message,
				"code":  apiErr.Code,
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to switch tenant",
				"code":  "SERVER_ERROR",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Successfully switched tenant",
		"tenant_id": tenantID,
	})
}

// GetCurrentTenant handles GET /api/v1/auth/tenants/current
func (h *AuthHandler) GetCurrentTenant(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Authentication required",
			"code":  "AUTH_REQUIRED",
		})
		return
	}

	user, err := h.service.GetUserByID(c.Request.Context(), userID.(string))
	if err != nil {
		if err == errors.ErrUserNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User not found",
				"code":  "USER_NOT_FOUND",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to get user",
				"code":  "SERVER_ERROR",
			})
		}
		return
	}

	currentMembership := user.GetCurrentTenantMembership()
	if currentMembership == nil {
		c.JSON(http.StatusOK, gin.H{
			"current_tenant": nil,
			"message":        "No current tenant set",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"current_tenant": currentMembership,
	})
}
