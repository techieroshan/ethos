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
