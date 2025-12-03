package middleware

import (
	"net/http"
	"strings"

	authRepo "ethos/internal/auth/repository"
	"ethos/internal/auth/model"
	"ethos/pkg/errors"
	"ethos/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// TenantContextMiddleware manages tenant context and boundary enforcement
func TenantContextMiddleware(tokenGen *jwt.TokenGenerator) gin.HandlerFunc {
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

		// Get user with tenant memberships
		userRepo := authRepo.NewPostgresRepository(nil) // We'll need to inject this properly
		user, err := userRepo.GetUserByID(c.Request.Context(), userID.(string))
		if err != nil {
			if err == errors.ErrUserNotFound {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "User not found",
					"code":  "USER_NOT_FOUND",
				})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Failed to verify user",
					"code":  "SERVER_ERROR",
				})
			}
			c.Abort()
			return
		}

		// Check if user has any tenant memberships
		if len(user.TenantMemberships) == 0 {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "User is not a member of any tenant",
				"code":  "NO_TENANT_ACCESS",
			})
			c.Abort()
			return
		}

		// Determine tenant context from various sources:
		// 1. Explicit tenant header
		// 2. URL parameter (for org-specific routes)
		// 3. User's current tenant
		// 4. User's default/first active tenant

		var tenantID string
		var tenantSource string

		// Check for explicit tenant header
		if tenantHeader := c.GetHeader("X-Tenant-ID"); tenantHeader != "" {
			tenantID = tenantHeader
			tenantSource = "header"
		}

		// Check for tenant in URL (for org-specific routes like /organizations/:org_id/...)
		if tenantParam := c.Param("org_id"); tenantParam != "" {
			tenantID = tenantParam
			tenantSource = "url_param"
		}

		// If no explicit tenant, use user's current tenant
		if tenantID == "" && user.CurrentTenantID != nil {
			tenantID = *user.CurrentTenantID
			tenantSource = "user_current"
		}

		// If still no tenant, use first active tenant as default
		if tenantID == "" {
			for _, membership := range user.TenantMemberships {
				if membership.IsActive {
					tenantID = membership.TenantID
					tenantSource = "user_default"
					break
				}
			}
		}

		// Validate tenant access
		if tenantID == "" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "No valid tenant context available",
				"code":  "NO_TENANT_CONTEXT",
			})
			c.Abort()
			return
		}

		if !user.HasTenantAccess(tenantID) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Access denied - user does not have access to this tenant",
				"code":  "TENANT_ACCESS_DENIED",
			})
			c.Abort()
			return
		}

		// Get tenant membership for context
		var currentMembership *model.TenantMembership
		for _, membership := range user.TenantMemberships {
			if membership.TenantID == tenantID && membership.IsActive {
				currentMembership = &membership
				break
			}
		}

		if currentMembership == nil {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Invalid tenant membership",
				"code":  "INVALID_TENANT_MEMBERSHIP",
			})
			c.Abort()
			return
		}

		// Set tenant context in request
		c.Set("tenant_id", tenantID)
		c.Set("tenant_membership", currentMembership)
		c.Set("tenant_source", tenantSource)
		c.Set("user", user) // Update user with full context

		// Add tenant context to response headers for frontend awareness
		c.Header("X-Current-Tenant-ID", tenantID)
		c.Header("X-Current-Tenant-Role", currentMembership.Role)

		// Log tenant context switch for audit
		if tenantSource == "header" || tenantSource == "url_param" {
			// TODO: Log tenant context switches for audit purposes
		}

		c.Next()
	}
}

// TenantBoundaryMiddleware enforces tenant data isolation
func TenantBoundaryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get tenant context
		tenantID, exists := c.Get("tenant_id")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Tenant context not set",
				"code":  "MISSING_TENANT_CONTEXT",
			})
			c.Abort()
			return
		}

		// Add tenant filter to query parameters for database operations
		// This ensures all database queries are automatically scoped to the tenant
		query := c.Request.URL.Query()
		query.Set("_tenant_id", tenantID.(string))
		c.Request.URL.RawQuery = query.Encode()

		// Set tenant context in response for frontend
		c.Header("X-Tenant-ID", tenantID.(string))

		c.Next()
	}
}

// CrossTenantAuditMiddleware logs cross-tenant operations
func CrossTenantAuditMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get tenant context
		_, tenantExists := c.Get("tenant_id")
		_, userExists := c.Get("user_id")

		if tenantExists && userExists {
			// Check if this is a cross-tenant operation
			// (This would be more sophisticated in production)
			isCrossTenantOp := strings.Contains(c.Request.URL.Path, "/admin") ||
				strings.Contains(c.Request.URL.Path, "/organizations")

			if isCrossTenantOp {
				// TODO: Log cross-tenant operation for audit
				// This should include: user_id, tenant_id, operation, timestamp, IP, etc.
			}
		}

		c.Next()
	}
}
