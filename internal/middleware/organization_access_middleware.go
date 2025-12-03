package middleware

import (
	"net/http"

	authRepo "ethos/internal/auth/repository"
	orgRepo "ethos/internal/organization/repository"
	"ethos/pkg/errors"
	"ethos/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// OrganizationAccessMiddleware validates that the authenticated user has access to the specified organization
func OrganizationAccessMiddleware(tokenGen *jwt.TokenGenerator) gin.HandlerFunc {
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

		// Check if user has platform admin privileges (they can access any organization)
		userRepo := authRepo.NewPostgresRepository(nil) // We'll need to inject this properly
		user, err := userRepo.GetUserByID(c.Request.Context(), userID.(string))
		if err == nil && user.IsAdmin() {
			// Platform admin has access to all organizations
			c.Set("organization_id", orgID)
			c.Set("user", user)
			c.Next()
			return
		}

		// Check if user is a member of the organization
		orgRepository := orgRepo.NewPostgresRepository(nil) // We'll need to inject this properly
		member, err := orgRepository.GetOrganizationMember(c.Request.Context(), orgID, userID.(string))
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

		// User has access to the organization
		c.Set("organization_id", orgID)
		c.Set("organization_member", member)
		c.Set("user", user)
		c.Next()
	}
}
