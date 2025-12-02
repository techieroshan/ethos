package api

import (
	"github.com/gin-gonic/gin"
	accountHandler "ethos/internal/account/handler"
	"ethos/internal/auth/handler"
	communityHandler "ethos/internal/community/handler"
	dashboardHandler "ethos/internal/dashboard/handler"
	feedbackHandler "ethos/internal/feedback/handler"
	notificationHandler "ethos/internal/notifications/handler"
	peopleHandler "ethos/internal/people/handler"
	"ethos/internal/middleware"
	profileHandler "ethos/internal/profile/handler"
	"ethos/pkg/jwt"
)

// SetupRoutes configures all API routes
func SetupRoutes(router *gin.Engine, authHandler *handler.AuthHandler, profileHandler *profileHandler.ProfileHandler, feedbackHandler *feedbackHandler.FeedbackHandler, notificationHandler *notificationHandler.NotificationHandler, dashboardHandler *dashboardHandler.DashboardHandler, peopleHandler *peopleHandler.PeopleHandler, communityHandler *communityHandler.CommunityHandler, accountHandler *accountHandler.AccountHandler, tokenGen *jwt.TokenGenerator) {
	v1 := router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/register", authHandler.Register)
			auth.POST("/refresh", authHandler.Refresh)
			auth.GET("/me", middleware.AuthMiddleware(tokenGen), authHandler.Me)
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
			}
			profile.GET("/user-profile", profileHandler.SearchProfiles)
		}

		feedback := v1.Group("/feedback")
		feedback.Use(middleware.AuthMiddleware(tokenGen))
		{
			feedback.GET("/feed", feedbackHandler.GetFeed)
			feedback.GET("/:feedback_id", feedbackHandler.GetFeedbackByID)
			feedback.GET("/:feedback_id/comments", feedbackHandler.GetComments)
			feedback.POST("", feedbackHandler.CreateFeedback)
			feedback.POST("/:feedback_id/comments", feedbackHandler.CreateComment)
			feedback.POST("/:feedback_id/react", feedbackHandler.AddReaction)
			feedback.DELETE("/:feedback_id/react", feedbackHandler.RemoveReaction)
		}

		notifications := v1.Group("/notifications")
		notifications.Use(middleware.AuthMiddleware(tokenGen))
		{
			notifications.GET("", notificationHandler.GetNotifications)
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

// SetupMiddleware configures global middleware
func SetupMiddleware(router *gin.Engine) {
	router.Use(middleware.TracingMiddleware())
	router.Use(gin.Recovery())
}
