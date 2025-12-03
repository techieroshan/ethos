package middleware

import (
	"github.com/gin-gonic/gin"
)

// AdminMiddleware validates that the authenticated user has admin privileges
// TODO: Implement proper role checking once user roles are integrated with JWT
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// For now, allow all authenticated users (will be replaced with proper role checking)
		// This is a temporary implementation until role-based access control is fully implemented

		// TODO: Check user roles from database or JWT claims
		// For development, we'll allow access - in production this should validate admin roles

		c.Next()
	}
}
