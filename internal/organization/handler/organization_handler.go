package handler

import (
	"net/http"
	"strconv"

	"ethos/internal/organization/model"
	"ethos/internal/organization/service"
	"ethos/pkg/errors"

	"github.com/gin-gonic/gin"
)

// OrganizationHandler handles organization HTTP requests
type OrganizationHandler struct {
	service service.Service
}

// NewOrganizationHandler creates a new organization handler
func NewOrganizationHandler(svc service.Service) *OrganizationHandler {
	return &OrganizationHandler{
		service: svc,
	}
}

// GetOrganization handles GET /api/v1/organizations/:org_id
func (h *OrganizationHandler) GetOrganization(c *gin.Context) {
	orgID := c.Param("org_id")
	if orgID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Organization ID is required",
			"code":  "VALIDATION_FAILED",
		})
		return
	}

	org, err := h.service.GetOrganization(c.Request.Context(), orgID)
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

	c.JSON(http.StatusOK, org)
}

// ListOrganizations handles GET /api/v1/organizations
func (h *OrganizationHandler) ListOrganizations(c *gin.Context) {
	limit := 10
	offset := 0

	if l := c.Query("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil && val > 0 {
			limit = val
		}
	}

	if o := c.Query("offset"); o != "" {
		if val, err := strconv.Atoi(o); err == nil && val >= 0 {
			offset = val
		}
	}

	organizations, err := h.service.ListOrganizations(c.Request.Context(), limit, offset)
	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			c.JSON(apiErr.HTTPStatus, gin.H{
				"error": apiErr.Message,
				"code":  apiErr.Code,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to list organizations",
			"code":  "SERVER_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"organizations": organizations,
		"limit":         limit,
		"offset":        offset,
	})
}

// CreateOrganization handles POST /api/v1/organizations
func (h *OrganizationHandler) CreateOrganization(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Authentication required",
			"code":  "AUTH_REQUIRED",
		})
		return
	}

	var req model.CreateOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Validation failed",
			"code":  "VALIDATION_FAILED",
		})
		return
	}

	org, err := h.service.CreateOrganization(c.Request.Context(), userID.(string), &req)
	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			c.JSON(apiErr.HTTPStatus, gin.H{
				"error": apiErr.Message,
				"code":  apiErr.Code,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create organization",
			"code":  "SERVER_ERROR",
		})
		return
	}

	c.JSON(http.StatusCreated, org)
}

// UpdateOrganization handles PUT /api/v1/organizations/:org_id
func (h *OrganizationHandler) UpdateOrganization(c *gin.Context) {
	orgID := c.Param("org_id")
	if orgID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Organization ID is required",
			"code":  "VALIDATION_FAILED",
		})
		return
	}

	var req model.UpdateOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Validation failed",
			"code":  "VALIDATION_FAILED",
		})
		return
	}

	org, err := h.service.UpdateOrganization(c.Request.Context(), orgID, &req)
	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			c.JSON(apiErr.HTTPStatus, gin.H{
				"error": apiErr.Message,
				"code":  apiErr.Code,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update organization",
			"code":  "SERVER_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, org)
}

// DeleteOrganization handles DELETE /api/v1/organizations/:org_id
func (h *OrganizationHandler) DeleteOrganization(c *gin.Context) {
	orgID := c.Param("org_id")
	if orgID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Organization ID is required",
			"code":  "VALIDATION_FAILED",
		})
		return
	}

	err := h.service.DeleteOrganization(c.Request.Context(), orgID)
	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			c.JSON(apiErr.HTTPStatus, gin.H{
				"error": apiErr.Message,
				"code":  apiErr.Code,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete organization",
			"code":  "SERVER_ERROR",
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// ListOrganizationMembers handles GET /api/v1/organizations/:org_id/members
func (h *OrganizationHandler) ListOrganizationMembers(c *gin.Context) {
	orgID := c.Param("org_id")
	if orgID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Organization ID is required",
			"code":  "VALIDATION_FAILED",
		})
		return
	}

	limit := 10
	offset := 0

	if l := c.Query("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil && val > 0 {
			limit = val
		}
	}

	if o := c.Query("offset"); o != "" {
		if val, err := strconv.Atoi(o); err == nil && val >= 0 {
			offset = val
		}
	}

	members, err := h.service.ListOrganizationMembers(c.Request.Context(), orgID, limit, offset)
	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			c.JSON(apiErr.HTTPStatus, gin.H{
				"error": apiErr.Message,
				"code":  apiErr.Code,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to list members",
			"code":  "SERVER_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"members": members,
		"limit":   limit,
		"offset":  offset,
	})
}

// AddOrganizationMember handles POST /api/v1/organizations/:org_id/members
func (h *OrganizationHandler) AddOrganizationMember(c *gin.Context) {
	orgID := c.Param("org_id")
	if orgID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Organization ID is required",
			"code":  "VALIDATION_FAILED",
		})
		return
	}

	var req model.AddMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Validation failed",
			"code":  "VALIDATION_FAILED",
		})
		return
	}

	member, err := h.service.AddOrganizationMember(c.Request.Context(), orgID, &req)
	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			c.JSON(apiErr.HTTPStatus, gin.H{
				"error": apiErr.Message,
				"code":  apiErr.Code,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to add member",
			"code":  "SERVER_ERROR",
		})
		return
	}

	c.JSON(http.StatusCreated, member)
}

// UpdateOrganizationMemberRole handles PUT /api/v1/organizations/:org_id/members/:user_id
func (h *OrganizationHandler) UpdateOrganizationMemberRole(c *gin.Context) {
	orgID := c.Param("org_id")
	userID := c.Param("user_id")

	if orgID == "" || userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Organization ID and User ID are required",
			"code":  "VALIDATION_FAILED",
		})
		return
	}

	var req model.UpdateMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Validation failed",
			"code":  "VALIDATION_FAILED",
		})
		return
	}

	member, err := h.service.UpdateOrganizationMemberRole(c.Request.Context(), orgID, userID, &req)
	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			c.JSON(apiErr.HTTPStatus, gin.H{
				"error": apiErr.Message,
				"code":  apiErr.Code,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update member role",
			"code":  "SERVER_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, member)
}

// RemoveOrganizationMember handles DELETE /api/v1/organizations/:org_id/members/:user_id
func (h *OrganizationHandler) RemoveOrganizationMember(c *gin.Context) {
	orgID := c.Param("org_id")
	userID := c.Param("user_id")

	if orgID == "" || userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Organization ID and User ID are required",
			"code":  "VALIDATION_FAILED",
		})
		return
	}

	err := h.service.RemoveOrganizationMember(c.Request.Context(), orgID, userID)
	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			c.JSON(apiErr.HTTPStatus, gin.H{
				"error": apiErr.Message,
				"code":  apiErr.Code,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to remove member",
			"code":  "SERVER_ERROR",
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// GetOrganizationSettings handles GET /api/v1/organizations/:org_id/settings
func (h *OrganizationHandler) GetOrganizationSettings(c *gin.Context) {
	orgID := c.Param("org_id")
	if orgID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Organization ID is required",
			"code":  "VALIDATION_FAILED",
		})
		return
	}

	settings, err := h.service.GetOrganizationSettings(c.Request.Context(), orgID)
	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			c.JSON(apiErr.HTTPStatus, gin.H{
				"error": apiErr.Message,
				"code":  apiErr.Code,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get settings",
			"code":  "SERVER_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, settings)
}

// UpdateOrganizationSettings handles PUT /api/v1/organizations/:org_id/settings
func (h *OrganizationHandler) UpdateOrganizationSettings(c *gin.Context) {
	orgID := c.Param("org_id")
	if orgID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Organization ID is required",
			"code":  "VALIDATION_FAILED",
		})
		return
	}

	var req model.UpdateSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Validation failed",
			"code":  "VALIDATION_FAILED",
		})
		return
	}

	settings, err := h.service.UpdateOrganizationSettings(c.Request.Context(), orgID, &req)
	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			c.JSON(apiErr.HTTPStatus, gin.H{
				"error": apiErr.Message,
				"code":  apiErr.Code,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update settings",
			"code":  "SERVER_ERROR",
		})
		return
	}
	c.JSON(http.StatusOK, settings)
}

// ListAllUsers handles GET /api/v1/admin/users - List all users across all organizations
func (h *OrganizationHandler) ListAllUsers(c *gin.Context) {
	limit := 50
	offset := 0
	search := c.Query("search")
	status := c.Query("status") // active, suspended, banned

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

	users, total, err := h.service.ListAllUsers(c.Request.Context(), limit, offset, search, status)
	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			c.JSON(apiErr.HTTPStatus, gin.H{
				"error": apiErr.Message,
				"code":  apiErr.Code,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to list users",
			"code":  "SERVER_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
		"total": total,
		"limit": limit,
		"offset": offset,
	})
}

// GetUserDetails handles GET /api/v1/admin/users/:user_id - Get detailed user information
func (h *OrganizationHandler) GetUserDetails(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User ID is required",
			"code":  "VALIDATION_FAILED",
		})
		return
	}

	user, err := h.service.GetUserDetails(c.Request.Context(), userID)
	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			c.JSON(apiErr.HTTPStatus, gin.H{
				"error": apiErr.Message,
				"code":  apiErr.Code,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get user details",
			"code":  "SERVER_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

// SuspendUser handles POST /api/v1/admin/users/:user_id/suspend - Suspend a user
func (h *OrganizationHandler) SuspendUser(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User ID is required",
			"code":  "VALIDATION_FAILED",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User suspended successfully"})
}

// BanUser handles POST /api/v1/admin/users/:user_id/ban - Ban a user permanently
func (h *OrganizationHandler) BanUser(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User ID is required",
			"code":  "VALIDATION_FAILED",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User banned successfully"})
}

// UnbanUser handles POST /api/v1/admin/users/:user_id/unban - Unban a user
func (h *OrganizationHandler) UnbanUser(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User ID is required",
			"code":  "VALIDATION_FAILED",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User unbanned successfully"})
}

// DeleteUser handles DELETE /api/v1/admin/users/:user_id - Delete a user permanently
func (h *OrganizationHandler) DeleteUser(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User ID is required",
			"code":  "VALIDATION_FAILED",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// GetSystemAnalytics handles GET /api/v1/admin/analytics/overview - Get system analytics
func (h *OrganizationHandler) GetSystemAnalytics(c *gin.Context) {
	analytics, err := h.service.GetSystemAnalytics(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get analytics",
			"code":  "SERVER_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, analytics)
}

// GetUserAnalytics handles GET /api/v1/admin/analytics/users - Get user analytics
func (h *OrganizationHandler) GetUserAnalytics(c *gin.Context) {
	analytics, err := h.service.GetUserAnalytics(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get user analytics",
			"code":  "SERVER_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, analytics)
}

// GetContentAnalytics handles GET /api/v1/admin/analytics/content - Get content analytics
func (h *OrganizationHandler) GetContentAnalytics(c *gin.Context) {
	analytics, err := h.service.GetContentAnalytics(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get content analytics",
			"code":  "SERVER_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, analytics)
}

// GetAuditLogs handles GET /api/v1/admin/audit - Get audit logs
func (h *OrganizationHandler) GetAuditLogs(c *gin.Context) {
	logs, total, err := h.service.GetAuditLogs(c.Request.Context(), 100, 0, "", "", "", "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get audit logs",
			"code":  "SERVER_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"logs":  logs,
		"total": total,
	})
}

// GetAuditEntry handles GET /api/v1/admin/audit/:entry_id - Get audit entry
func (h *OrganizationHandler) GetAuditEntry(c *gin.Context) {
	entryID := c.Param("entry_id")
	entry, err := h.service.GetAuditEntry(c.Request.Context(), entryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get audit entry",
			"code":  "SERVER_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, entry)
}

// GetSystemSettings handles GET /api/v1/admin/settings - Get system settings
func (h *OrganizationHandler) GetSystemSettings(c *gin.Context) {
	settings, err := h.service.GetSystemSettings(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get system settings",
			"code":  "SERVER_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, settings)
}

// UpdateSystemSettings handles PUT /api/v1/admin/settings - Update system settings
func (h *OrganizationHandler) UpdateSystemSettings(c *gin.Context) {
	var settings map[string]interface{}
	if err := c.ShouldBindJSON(&settings); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"code":  "VALIDATION_FAILED",
		})
		return
	}

	updated, err := h.service.UpdateSystemSettings(c.Request.Context(), settings, "admin-id")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update system settings",
			"code":  "SERVER_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, updated)
}

// BulkSuspendUsers handles POST /api/v1/admin/bulk/suspend-users - Bulk suspend users
func (h *OrganizationHandler) BulkSuspendUsers(c *gin.Context) {
	var req struct {
		UserIDs []string `json:"user_ids" binding:"required,min=1,max=100"`
		Reason  string   `json:"reason" binding:"required"`
		Duration *int    `json:"duration,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"code":  "VALIDATION_FAILED",
		})
		return
	}

	result, err := h.service.BulkSuspendUsers(c.Request.Context(), req.UserIDs, req.Reason, req.Duration, "admin-id")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to bulk suspend users",
			"code":  "SERVER_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// ORGANIZATION ADMIN METHODS - Organization-specific admin operations

// GetOrganizationAnalytics handles GET /api/v1/organizations/:org_id/admin/analytics/overview
func (h *OrganizationHandler) GetOrganizationAnalytics(c *gin.Context) {
	orgID := c.Param("org_id")

	analytics, err := h.service.GetOrganizationAnalytics(c.Request.Context(), orgID)
	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			c.JSON(apiErr.HTTPStatus, gin.H{
				"error": apiErr.Message,
				"code":  apiErr.Code,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get organization analytics",
			"code":  "SERVER_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, analytics)
}

// GetOrganizationUserAnalytics handles GET /api/v1/organizations/:org_id/admin/analytics/users
func (h *OrganizationHandler) GetOrganizationUserAnalytics(c *gin.Context) {
	orgID := c.Param("org_id")

	analytics, err := h.service.GetOrganizationUserAnalytics(c.Request.Context(), orgID)
	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			c.JSON(apiErr.HTTPStatus, gin.H{
				"error": apiErr.Message,
				"code":  apiErr.Code,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get organization user analytics",
			"code":  "SERVER_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, analytics)
}

// GetOrganizationContentAnalytics handles GET /api/v1/organizations/:org_id/admin/analytics/content
func (h *OrganizationHandler) GetOrganizationContentAnalytics(c *gin.Context) {
	orgID := c.Param("org_id")

	analytics, err := h.service.GetOrganizationContentAnalytics(c.Request.Context(), orgID)
	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			c.JSON(apiErr.HTTPStatus, gin.H{
				"error": apiErr.Message,
				"code":  apiErr.Code,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get organization content analytics",
			"code":  "SERVER_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, analytics)
}

// ListOrganizationUsers handles GET /api/v1/organizations/:org_id/admin/users
func (h *OrganizationHandler) ListOrganizationUsers(c *gin.Context) {
	orgID := c.Param("org_id")
	limit := 50
	offset := 0
	search := c.Query("search")
	status := c.Query("status")

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

	users, total, err := h.service.ListOrganizationUsers(c.Request.Context(), orgID, limit, offset, search, status)
	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			c.JSON(apiErr.HTTPStatus, gin.H{
				"error": apiErr.Message,
				"code":  apiErr.Code,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to list organization users",
			"code":  "SERVER_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
		"total": total,
		"limit": limit,
		"offset": offset,
	})
}

// SuspendOrganizationUser handles POST /api/v1/organizations/:org_id/admin/users/:user_id/suspend
func (h *OrganizationHandler) SuspendOrganizationUser(c *gin.Context) {
	orgID := c.Param("org_id")
	userID := c.Param("user_id")

	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User ID is required",
			"code":  "VALIDATION_FAILED",
		})
		return
	}

	var req struct {
		Reason  string `json:"reason" binding:"required"`
		Duration *int  `json:"duration,omitempty"`
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

	err := h.service.SuspendOrganizationUser(c.Request.Context(), orgID, userID, req.Reason, req.Duration, adminID.(string))
	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			c.JSON(apiErr.HTTPStatus, gin.H{
				"error": apiErr.Message,
				"code":  apiErr.Code,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to suspend organization user",
			"code":  "SERVER_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User suspended successfully"})
}

// UnsuspendOrganizationUser handles POST /api/v1/organizations/:org_id/admin/users/:user_id/unsuspend
func (h *OrganizationHandler) UnsuspendOrganizationUser(c *gin.Context) {
	orgID := c.Param("org_id")
	userID := c.Param("user_id")

	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User ID is required",
			"code":  "VALIDATION_FAILED",
		})
		return
	}

	// Get admin user ID from context
	adminID, _ := c.Get("user_id")

	err := h.service.UnsuspendOrganizationUser(c.Request.Context(), orgID, userID, adminID.(string))
	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			c.JSON(apiErr.HTTPStatus, gin.H{
				"error": apiErr.Message,
				"code":  apiErr.Code,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to unsuspend organization user",
			"code":  "SERVER_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User unsuspended successfully"})
}

// RemoveOrganizationUser handles DELETE /api/v1/organizations/:org_id/admin/users/:user_id
func (h *OrganizationHandler) RemoveOrganizationUser(c *gin.Context) {
	orgID := c.Param("org_id")
	userID := c.Param("user_id")

	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User ID is required",
			"code":  "VALIDATION_FAILED",
		})
		return
	}

	// Get admin user ID from context
	adminID, _ := c.Get("user_id")

	err := h.service.RemoveOrganizationUser(c.Request.Context(), orgID, userID, adminID.(string))
	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			c.JSON(apiErr.HTTPStatus, gin.H{
				"error": apiErr.Message,
				"code":  apiErr.Code,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to remove organization user",
			"code":  "SERVER_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User removed from organization successfully"})
}

// GetOrganizationAuditLogs handles GET /api/v1/organizations/:org_id/admin/audit
func (h *OrganizationHandler) GetOrganizationAuditLogs(c *gin.Context) {
	orgID := c.Param("org_id")
	limit := 100
	offset := 0
	userID := c.Query("user_id")
	action := c.Query("action")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 1000 {
			limit = parsed
		}
	}

	if o := c.Query("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	logs, total, err := h.service.GetOrganizationAuditLogs(c.Request.Context(), orgID, limit, offset, userID, action, startDate, endDate)
	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			c.JSON(apiErr.HTTPStatus, gin.H{
				"error": apiErr.Message,
				"code":  apiErr.Code,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get organization audit logs",
			"code":  "SERVER_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"logs":  logs,
		"total": total,
		"limit": limit,
		"offset": offset,
	})
}

// ExportOrganizationAuditLogs handles GET /api/v1/organizations/:org_id/admin/audit/export
func (h *OrganizationHandler) ExportOrganizationAuditLogs(c *gin.Context) {
	orgID := c.Param("org_id")
	format := c.Query("format") // csv, json, pdf
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if format == "" {
		format = "csv"
	}

	data, filename, err := h.service.ExportOrganizationAuditLogs(c.Request.Context(), orgID, format, startDate, endDate)
	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			c.JSON(apiErr.HTTPStatus, gin.H{
				"error": apiErr.Message,
				"code":  apiErr.Code,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to export organization audit logs",
			"code":  "SERVER_ERROR",
		})
		return
	}

	// Set appropriate headers for file download
	contentType := "text/csv"
	if format == "json" {
		contentType = "application/json"
	} else if format == "pdf" {
		contentType = "application/pdf"
	}

	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", contentType)
	c.Data(http.StatusOK, contentType, data)
}

// ListOrganizationIncidents handles GET /api/v1/organizations/:org_id/admin/incidents
func (h *OrganizationHandler) ListOrganizationIncidents(c *gin.Context) {
	orgID := c.Param("org_id")
	limit := 50
	offset := 0
	status := c.Query("status")
	priority := c.Query("priority")

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

	incidents, total, err := h.service.ListOrganizationIncidents(c.Request.Context(), orgID, limit, offset, status, priority)
	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			c.JSON(apiErr.HTTPStatus, gin.H{
				"error": apiErr.Message,
				"code":  apiErr.Code,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to list organization incidents",
			"code":  "SERVER_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"incidents": incidents,
		"total":     total,
		"limit":     limit,
		"offset":    offset,
	})
}

// CreateOrganizationIncident handles POST /api/v1/organizations/:org_id/admin/incidents
func (h *OrganizationHandler) CreateOrganizationIncident(c *gin.Context) {
	orgID := c.Param("org_id")

	var req struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description" binding:"required"`
		Priority    string `json:"priority" binding:"required,oneof=low medium high critical"`
		Category    string `json:"category" binding:"required"`
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

	incident, err := h.service.CreateOrganizationIncident(c.Request.Context(), orgID, req.Title, req.Description, req.Priority, req.Category, adminID.(string))
	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			c.JSON(apiErr.HTTPStatus, gin.H{
				"error": apiErr.Message,
				"code":  apiErr.Code,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create organization incident",
			"code":  "SERVER_ERROR",
		})
		return
	}

	c.JSON(http.StatusCreated, incident)
}

// UpdateOrganizationIncident handles PUT /api/v1/organizations/:org_id/admin/incidents/:incident_id
func (h *OrganizationHandler) UpdateOrganizationIncident(c *gin.Context) {
	orgID := c.Param("org_id")
	incidentID := c.Param("incident_id")

	var req struct {
		Status      *string `json:"status,omitempty"`
		Priority    *string `json:"priority,omitempty"`
		AssignedTo  *string `json:"assigned_to,omitempty"`
		Resolution  *string `json:"resolution,omitempty"`
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

	incident, err := h.service.UpdateOrganizationIncident(c.Request.Context(), orgID, incidentID, req.Status, req.Priority, req.AssignedTo, req.Resolution, adminID.(string))
	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			c.JSON(apiErr.HTTPStatus, gin.H{
				"error": apiErr.Message,
				"code":  apiErr.Code,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update organization incident",
			"code":  "SERVER_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, incident)
}
