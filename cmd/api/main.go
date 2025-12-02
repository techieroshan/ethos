package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"ethos/api"
	accountHandler "ethos/internal/account/handler"
	accountRepository "ethos/internal/account/repository"
	accountService "ethos/internal/account/service"
	"ethos/internal/auth/handler"
	"ethos/internal/auth/repository"
	"ethos/internal/auth/service"
	communityHandler "ethos/internal/community/handler"
	dashboardHandler "ethos/internal/dashboard/handler"
	dashboardRepository "ethos/internal/dashboard/repository"
	dashboardService "ethos/internal/dashboard/service"
	feedbackHandler "ethos/internal/feedback/handler"
	feedbackRepository "ethos/internal/feedback/repository"
	feedbackService "ethos/internal/feedback/service"
	notificationHandler "ethos/internal/notifications/handler"
	notificationRepository "ethos/internal/notifications/repository"
	notificationService "ethos/internal/notifications/service"
	peopleHandler "ethos/internal/people/handler"
	peopleRepository "ethos/internal/people/repository"
	peopleService "ethos/internal/people/service"
	profileHandler "ethos/internal/profile/handler"
	profileRepository "ethos/internal/profile/repository"
	profileService "ethos/internal/profile/service"
	"ethos/internal/config"
	"ethos/internal/database"
	"ethos/pkg/jwt"
	"ethos/pkg/otel"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize OpenTelemetry
	ctx := context.Background()
	shutdown, err := otel.Init(ctx, cfg.OTEL.ServiceName, cfg.OTEL.JaegerURL, cfg.OTEL.Enabled)
	if err != nil {
		log.Fatalf("Failed to initialize OpenTelemetry: %v", err)
	}
	defer shutdown()

	// Initialize database
	db, err := database.New(ctx, database.Config{
		URL:             cfg.Database.URL,
		MaxConnections:  cfg.Database.MaxConnections,
		MaxIdleTime:     cfg.Database.MaxIdleTime,
		ConnMaxLifetime: cfg.Database.ConnMaxLifetime,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize dependencies
	authRepo := repository.NewPostgresRepository(db)
	tokenGen := jwt.NewTokenGenerator(
		cfg.JWT.AccessTokenSecret,
		cfg.JWT.RefreshTokenSecret,
		cfg.JWT.AccessTokenExpiry,
		cfg.JWT.RefreshTokenExpiry,
	)
	authService := service.NewAuthService(authRepo, tokenGen)
	authHandler := handler.NewAuthHandler(authService)

	// Initialize profile dependencies
	profileRepo := profileRepository.NewPostgresRepository(db)
	profileSvc := profileService.NewProfileService(profileRepo)
	profileHandler := profileHandler.NewProfileHandler(profileSvc)

	// Initialize feedback dependencies
	feedbackRepo := feedbackRepository.NewPostgresRepository(db)
	feedbackSvc := feedbackService.NewFeedbackService(feedbackRepo)
	feedbackHandler := feedbackHandler.NewFeedbackHandler(feedbackSvc)

	// Initialize notification dependencies
	notificationRepo := notificationRepository.NewPostgresRepository(db)
	notificationSvc := notificationService.NewNotificationService(notificationRepo)
	notificationHandler := notificationHandler.NewNotificationHandler(notificationSvc)

	// Initialize dashboard dependencies
	dashboardRepo := dashboardRepository.NewPostgresRepository(db)
	dashboardSvc := dashboardService.NewDashboardService(dashboardRepo)
	dashboardHandler := dashboardHandler.NewDashboardHandler(dashboardSvc)

	// Initialize people dependencies
	peopleRepo := peopleRepository.NewPostgresRepository(db)
	peopleSvc := peopleService.NewPeopleService(peopleRepo)
	peopleHandler := peopleHandler.NewPeopleHandler(peopleSvc)

	// Initialize community handler
	communityHandler := communityHandler.NewCommunityHandler()

	// Initialize account dependencies
	accountRepo := accountRepository.NewPostgresRepository(db)
	accountSvc := accountService.NewAccountService(accountRepo)
	accountHandler := accountHandler.NewAccountHandler(accountSvc)

	// Setup router
	router := gin.Default()
	api.SetupMiddleware(router)
	api.SetupRoutes(router, authHandler, profileHandler, feedbackHandler, notificationHandler, dashboardHandler, peopleHandler, communityHandler, accountHandler, tokenGen)

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
