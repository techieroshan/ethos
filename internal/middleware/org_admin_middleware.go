package middleware

import (
	"net/http"
	"strings"

	"ethos/internal/organization/repository"
	"ethos/pkg/errors"
	"ethos/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// OrgAdminMiddleware validates that the authenticated user has org admin privileges for the specified organization
func OrgAdminMiddleware(tokenGen *jwt.TokenGenerator) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from context (set by AuthMiddleware)
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authentication required",
				"code":  "AUTH_REQUIRED",
			})
			c.Abort()
			return
		}

		// Get organization ID from URL
		orgID := c.Param("org_id")
		if orgID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Organization ID is required",
				"code":  "VALIDATION_FAILED",
			})
			c.Abort()
			return
		}

		// Check if user has org admin privileges for this organization
		orgRepo := repository.NewPostgresRepository(nil) // We'll need to inject this properly

		// Check if user is a member of the organization with admin role
		member, err := orgRepo.GetOrganizationMember(c.Request.Context(), orgID, userID.(string))
		if err != nil {
			if err == errors.ErrNotFound {
				c.JSON(http.StatusForbidden, gin.H{
					"error": "Access denied - not a member of this organization",
					"code":  "ACCESS_DENIED",
				})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Failed to verify organization access",
					"code":  "SERVER_ERROR",
				})
			}
			c.Abort()
			return
		}

		// Check if user has admin role in this organization
		hasAdminRole := strings.Contains(member.Role, "admin") || strings.Contains(member.Role, "owner")
		if !hasAdminRole {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Organization admin privileges required",
				"code":  "INSUFFICIENT_ORG_PERMISSIONS",
			})
			c.Abort()
			return
		}

		// Set organization and member info in context for handlers to use
		c.Set("organization_id", orgID)
		c.Set("organization_member", member)
		c.Next()
	}
}
