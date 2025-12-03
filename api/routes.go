package api

import (
	accountHandler "ethos/internal/account/handler"
	"ethos/internal/auth/handler"
	communityHandler "ethos/internal/community/handler"
	dashboardHandler "ethos/internal/dashboard/handler"
	feedbackHandler "ethos/internal/feedback/handler"
	"ethos/internal/middleware"
	moderationHandler "ethos/internal/moderation/handler"
	notificationHandler "ethos/internal/notifications/handler"
	organizationHandler "ethos/internal/organization/handler"
	peopleHandler "ethos/internal/people/handler"
	profileHandler "ethos/internal/profile/handler"
	"ethos/pkg/jwt"

github.com/gin-gonic/gin"
)

// SetupRoutes configures all API routes
func SetupRoutes(router *gin.Engine, authHandler *handler.AuthHandler, profileHandler *profileHandler.ProfileHandler, feedbackHandler *feedbackHandler.FeedbackHandler, notificationHandler *notificationHandler.NotificationHandler, dashboardHandler *dashboardHandler.DashboardHandler, organizationHandler *organizationHandler.OrganizationHandler, peopleHandler *peopleHandler.PeopleHandler, communityHandler *communityHandler.CommunityHandler, accountHandler *accountHandler.AccountHandler, moderationHandler *moderationHandler.ModerationHandler, tokenGen *jwt.TokenGenerator) {
	// Global OPTIONS handler for all API routes
	router.OPTIONS("/api/*path", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Length, Content-Type, Authorization, Accept, X-Requested-With")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "43200")
		c.Status(200)
	})

	v1 := router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/register", authHandler.Register)
			auth.POST("/refresh", authHandler.Refresh)
			auth.GET("/me", middleware.AuthMiddleware(tokenGen), authHandler.Me)
			auth.GET("/verify-email/:token", authHandler.VerifyEmail)
			auth.POST("/change-password", middleware.AuthMiddleware(tokenGen), authHandler.ChangePassword)
			auth.POST("/setup-2fa", middleware.AuthMiddleware(tokenGen), authHandler.Setup2FA)
			auth.DELETE("/setup-2fa", middleware.AuthMiddleware(tokenGen), accountHandler.Disable2FA)
		}

		profile := v1.Group("/profile")
		{
			profileProtected := profile.Group("")
			profileProtected.Use(middleware.AuthMiddleware(tokenGen))
			{
				profileProtected.GET("/me", profileHandler.GetProfile)
				profileProtected.PUT("/me", profileHandler.UpdateProfile)
				profileProtected.PATCH("/me/preferences", profileHandler.UpdatePreferences)
				profileProtected.DELETE("/me", profileHandler.DeleteProfile)
				profileProtected.POST("/opt-out", profileHandler.OptOut)
				profileProtected.POST("/anonymize", profileHandler.Anonymize)
				profileProtected.POST("/delete_request", profileHandler.RequestDeletion)
			}
			profile.GET("/user-profile", profileHandler.SearchProfiles)
			profile.GET("/:user_id", profileHandler.GetUserProfileByID)
		}

		organizations := v1.Group("/organizations")
		organizations.Use(middleware.AuthMiddleware(tokenGen))
		{
			organizations.GET("", organizationHandler.ListOrganizations)
			organizations.POST("", organizationHandler.CreateOrganization)
			organizations.GET("/:org_id", organizationHandler.GetOrganization)
			organizations.PUT("/:org_id", organizationHandler.UpdateOrganization)
			organizations.DELETE("/:org_id", organizationHandler.DeleteOrganization)
			organizations.GET("/:org_id/members", organizationHandler.ListOrganizationMembers)
			organizations.POST("/:org_id/members", organizationHandler.AddOrganizationMember)
			organizations.PUT("/:org_id/members/:user_id", organizationHandler.UpdateOrganizationMemberRole)
			organizations.DELETE("/:org_id/members/:user_id", organizationHandler.RemoveOrganizationMember)
			organizations.GET("/:org_id/settings", organizationHandler.GetOrganizationSettings)
			organizations.PUT("/:org_id/settings", organizationHandler.UpdateOrganizationSettings)

			// Moderation routes nested under organizations
			moderation := organizations.Group("/:org_id/moderation")
			{
				moderation.GET("/appeals", moderationHandler.ListAppeals)
				moderation.POST("/appeals", moderationHandler.SubmitAppeal)
				moderation.GET("/appeals/:appeal_id/context", moderationHandler.GetAppealContext)
				moderation.GET("/actions", moderationHandler.ListModerationActions)
				moderation.GET("/history/:user_id", moderationHandler.GetModerationHistory)
			}
		}

		feedback := v1.Group("/feedback")
		{
			feedback.GET("/feed", middleware.AuthMiddleware(tokenGen), feedbackHandler.GetFeed)
			feedback.GET("/:feedback_id", middleware.AuthMiddleware(tokenGen), feedbackHandler.GetFeedbackByID)
			feedback.GET("/:feedback_id/comments", middleware.AuthMiddleware(tokenGen), feedbackHandler.GetComments)
			feedback.POST("", middleware.AuthMiddleware(tokenGen), feedbackHandler.CreateFeedback)
			feedback.POST("/:feedback_id/comments", middleware.AuthMiddleware(tokenGen), feedbackHandler.CreateComment)
			feedback.POST("/:feedback_id/react", middleware.AuthMiddleware(tokenGen), feedbackHandler.AddReaction)
			feedback.DELETE("/:feedback_id/react", middleware.AuthMiddleware(tokenGen), feedbackHandler.RemoveReaction)
			feedback.GET("/templates", feedbackHandler.GetTemplates)
			feedback.POST("/template_suggestions", feedbackHandler.PostTemplateSuggestions)
			feedback.GET("/impact", feedbackHandler.GetImpact)
			feedback.POST("/batch", feedbackHandler.CreateBatchFeedback)
			feedback.GET("/bookmarks", feedbackHandler.GetBookmarks)
			feedback.POST("/bookmarks/:feedback_id", feedbackHandler.AddBookmark)
			feedback.DELETE("/bookmarks/:feedback_id", feedbackHandler.RemoveBookmark)
			feedback.GET("/export", feedbackHandler.ExportFeedback)
			feedback.GET("/analytics", middleware.AuthMiddleware(tokenGen), feedbackHandler.GetFeedbackAnalytics)
			feedback.PUT("/:feedback_id", middleware.AuthMiddleware(tokenGen), feedbackHandler.UpdateFeedback)
			feedback.DELETE("/:feedback_id", middleware.AuthMiddleware(tokenGen), feedbackHandler.DeleteFeedback)
			feedback.PUT("/:feedback_id/comments/:comment_id", middleware.AuthMiddleware(tokenGen), feedbackHandler.UpdateComment)
			feedback.DELETE("/:feedback_id/comments/:comment_id", middleware.AuthMiddleware(tokenGen), feedbackHandler.DeleteComment)
			feedback.GET("/search", middleware.AuthMiddleware(tokenGen), feedbackHandler.SearchFeedback)
			feedback.GET("/trending", middleware.AuthMiddleware(tokenGen), feedbackHandler.GetTrendingFeedback)
			feedback.POST("/:feedback_id/pin", middleware.AuthMiddleware(tokenGen), feedbackHandler.PinFeedback)
			feedback.DELETE("/:feedback_id/pin", middleware.AuthMiddleware(tokenGen), feedbackHandler.UnpinFeedback)
			feedback.GET("/stats", middleware.AuthMiddleware(tokenGen), feedbackHandler.GetFeedbackStats)
		}

		notifications := v1.Group("/notifications")
		notifications.Use(middleware.AuthMiddleware(tokenGen))
		{
			notifications.GET("", notificationHandler.GetNotifications)
			notifications.PUT("/:id/read", notificationHandler.MarkAsRead)
			notifications.PUT("/mark-all-read", notificationHandler.MarkAllAsRead)
			notifications.GET("/preferences", notificationHandler.GetPreferences)
			notifications.PUT("/preferences", notificationHandler.UpdatePreferences)
		}

		dashboard := v1.Group("/dashboard")
		dashboard.Use(middleware.AuthMiddleware(tokenGen))
		{
			dashboard.GET("", dashboardHandler.GetDashboard)
		}

		people := v1.Group("/people")
		people.Use(middleware.AuthMiddleware(tokenGen))
		{
			people.GET("/search", peopleHandler.SearchPeople)
			people.GET("/recommendations", peopleHandler.GetRecommendations)
		}

		community := v1.Group("/community")
		{
			community.GET("/rules", communityHandler.GetRules)
		}

		account := v1.Group("/account")
		account.Use(middleware.AuthMiddleware(tokenGen))
		{
			account.GET("/security-events", accountHandler.GetSecurityEvents)
			account.GET("/export-data/:export_id/status", accountHandler.GetExportStatus)
		}
	}
}


// corsHandler handles CORS for all requests
func corsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Add CORS headers to all responses
		c.Header("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Length, Content-Type, Authorization, Accept, X-Requested-With")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "43200")

		// Handle preflight OPTIONS requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		c.Next()
	}
}

// SetupMiddleware configures global middleware
func SetupMiddleware(router *gin.Engine) {
	router.Use(corsHandler())
	router.Use(gin.Logger())
	router.Use(middleware.TracingMiddleware())
	router.Use(gin.Recovery())
}
