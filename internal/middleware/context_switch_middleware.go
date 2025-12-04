package middleware

import (
	"net/http"

	"ethos/internal/organization/service"
	"ethos/pkg/errors"

	"github.com/gin-gonic/gin"
)

// ContextSwitchMiddleware validates context switching operations and enforces organization membership
func ContextSwitchMiddleware(contextService service.UserContextService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user from context (set by AuthMiddleware)
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authentication required",
				"code":  "AUTH_REQUIRED",
			})
			c.Abort()
			return
		}

		userIDStr := userID.(string)

		// Get current organization context
		currentContext, err := contextService.GetCurrentContext(c.Request.Context(), userIDStr)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to load user context",
				"code":  "CONTEXT_LOAD_FAILED",
			})
			c.Abort()
			return
		}

		// Store current context in request
		if currentContext != nil {
			c.Set("current_organization_id", currentContext.OrganizationID)
			c.Set("current_organization_name", currentContext.OrganizationName)
			c.Set("user_role", currentContext.Role)
			c.Set("user_permissions", currentContext.Permissions)
		}

		// Add organization context headers for frontend awareness
		if currentContext != nil {
			c.Header("X-Current-Organization-ID", currentContext.OrganizationID)
			c.Header("X-Current-Organization-Name", currentContext.OrganizationName)
			c.Header("X-User-Role", currentContext.Role)
		}

		c.Next()
	}
}

// ValidateOrganizationMembership ensures user is a member of the requested organization
func ValidateOrganizationMembership(contextService service.UserContextService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user from context (set by AuthMiddleware)
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authentication required",
				"code":  "AUTH_REQUIRED",
			})
			c.Abort()
			return
		}

		// Get organization ID from URL parameter
		organizationID := c.Param("org_id")
		if organizationID == "" {
			c.Next()
			return
		}

		userIDStr := userID.(string)

		// Check if user is a member of this organization
		isMember, err := contextService.ValidateUserInOrganization(c.Request.Context(), userIDStr, organizationID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to validate membership",
				"code":  "MEMBERSHIP_CHECK_FAILED",
			})
			c.Abort()
			return
		}

		if !isMember {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "User is not a member of this organization",
				"code":  "NOT_ORGANIZATION_MEMBER",
			})
			c.Abort()
			return
		}

		// Get user's role in this organization
		role, err := contextService.GetUserRoleInOrganization(c.Request.Context(), userIDStr, organizationID)
		if err != nil && err != errors.ErrNotFound {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to determine user role",
				"code":  "ROLE_CHECK_FAILED",
			})
			c.Abort()
			return
		}

		// Store organization context
		c.Set("target_organization_id", organizationID)
		c.Set("user_role_in_org", role)

		// Add headers for frontend awareness
		c.Header("X-Target-Organization-ID", organizationID)
		if role != "" {
			c.Header("X-User-Role-In-Org", role)
		}

		c.Next()
	}
}

// EnforceOrganizationContext ensures requests have proper organization context
func EnforceOrganizationContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		// For protected organization endpoints, require organization context
		if c.GetString("current_organization_id") == "" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Organization context required",
				"code":  "MISSING_ORG_CONTEXT",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// LogContextOperation logs context switching and organization-level operations for audit
func LogContextOperation(operationType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")
		organizationID := c.GetString("current_organization_id")

		if userID != "" && organizationID != "" {
			// TODO: Implement audit logging
			// This should log:
			// - Operation type (context_switch, org_access, etc.)
			// - User ID
			// - Organization ID
			// - Request path and method
			// - IP address
			// - User agent
			// - Timestamp
			// - Response status code (available in c.Writer.Status() after Next())

			c.Set("audit_operation_type", operationType)
			c.Set("audit_user_id", userID)
			c.Set("audit_organization_id", organizationID)
		}

		c.Next()
	}
}
