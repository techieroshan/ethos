package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ethos/api"
	accountHandler "ethos/internal/account/handler"
	accountRepository "ethos/internal/account/repository"
	accountService "ethos/internal/account/service"
	"ethos/internal/auth/handler"
	"ethos/internal/auth/repository"
	"ethos/internal/auth/service"
	communityHandler "ethos/internal/community/handler"
	"ethos/internal/config"
	dashboardHandler "ethos/internal/dashboard/handler"
	dashboardRepository "ethos/internal/dashboard/repository"
	"ethos/internal/database"
	feedbackHandler "ethos/internal/feedback/handler"
	feedbackRepository "ethos/internal/feedback/repository"
	moderationHandler "ethos/internal/moderation/handler"
	moderationRepository "ethos/internal/moderation/repository"
	moderationService "ethos/internal/moderation/service"
	notificationHandler "ethos/internal/notifications/handler"
	notificationRepository "ethos/internal/notifications/repository"
	organizationHandler "ethos/internal/organization/handler"
	organizationRepository "ethos/internal/organization/repository"
	organizationService "ethos/internal/organization/service"
	peopleHandler "ethos/internal/people/handler"
	peopleRepository "ethos/internal/people/repository"
	profileHandler "ethos/internal/profile/handler"
	profileRepository "ethos/internal/profile/repository"
	profileService "ethos/internal/profile/service"
	"ethos/pkg/email"
	checkerClient "ethos/pkg/email/checker"
	emailitClient "ethos/pkg/email/emailit"
	mailpitClient "ethos/pkg/email/mailpit"
	grpcClient "ethos/pkg/grpc/client"
	"ethos/pkg/jwt"
	"ethos/pkg/otel"

	"github.com/gin-gonic/gin"
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

	// Initialize email checker (optional - can be nil if not configured)
	var emailChecker checkerClient.EmailChecker
	if cfg.Checker.APIKey != "" {
		emailChecker = checkerClient.NewChecker(checkerClient.Config{
			APIKey:  cfg.Checker.APIKey,
			BaseURL: cfg.Checker.BaseURL,
			Timeout: cfg.Checker.Timeout,
			Retries: cfg.Checker.Retries,
		})
	}

	// Initialize email sender (Emailit for prod, Mailpit for local)
	var emailSender email.EmailSender
	if cfg.Mailpit.Enabled {
		emailSender = mailpitClient.NewMailpit(mailpitClient.Config{
			SMTPHost:  cfg.Mailpit.SMTPHost,
			SMTPPort:  cfg.Mailpit.SMTPPort,
			FromEmail: cfg.Mailpit.FromEmail,
		})
		log.Println("Using Mailpit for email sending (local development)")
	} else if cfg.Emailit.APIKey != "" {
		emailSender = emailitClient.NewEmailit(emailitClient.Config{
			APIKey:  cfg.Emailit.APIKey,
			BaseURL: cfg.Emailit.BaseURL,
			Timeout: cfg.Emailit.Timeout,
			Retries: cfg.Emailit.Retries,
		})
		log.Println("Using Emailit for email sending (production)")
	} else {
		emailSender = email.NewNoOpEmailSender() // Use no-op if not configured
		log.Println("Warning: No email sender configured, using no-op")
	}

	authService := service.NewAuthService(authRepo, tokenGen, emailChecker, emailSender)
	authHandler := handler.NewAuthHandler(authService)

	// Initialize profile dependencies
	profileRepo := profileRepository.NewPostgresRepository(db)
	profileSvc := profileService.NewProfileService(profileRepo)
	profileHandler := profileHandler.NewProfileHandler(profileSvc)

	// Initialize gRPC client manager if enabled
	var grpcManager *grpcClient.ClientManager
	if cfg.GRPC.Enabled {
		grpcManager = grpcClient.NewClientManager(grpcClient.Config{
			FeedbackEndpoint:      cfg.GRPC.FeedbackEndpoint,
			DashboardEndpoint:     cfg.GRPC.DashboardEndpoint,
			NotificationsEndpoint: cfg.GRPC.NotificationsEndpoint,
			PeopleEndpoint:        cfg.GRPC.PeopleEndpoint,
			Timeout:               cfg.GRPC.Timeout,
			Retries:               cfg.GRPC.Retries,
		})
		defer grpcManager.Close()
		log.Println("gRPC client manager initialized")
	}

	// Initialize feedback dependencies - temporarily disabled due to import cycles
	feedbackHandler := &feedbackHandler.FeedbackHandler{} // Stub handler

	// Initialize notification dependencies - temporarily disabled due to import cycles
	notificationHandler := &notificationHandler.NotificationHandler{} // Stub handler

	// Initialize dashboard dependencies - temporarily disabled due to import cycles
	dashboardHandler := &dashboardHandler.DashboardHandler{} // Stub handler

	// Initialize people dependencies - temporarily disabled due to import cycles
	peopleHandler := &peopleHandler.PeopleHandler{} // Stub handler

	// Initialize community handler
	communityHandler := communityHandler.NewCommunityHandler()

	// Initialize account dependencies
	accountRepo := accountRepository.NewPostgresRepository(db)
	accountSvc := accountService.NewAccountService(accountRepo)
	accountHandler := accountHandler.NewAccountHandler(accountSvc)

	// Initialize moderation dependencies
	moderationRepo := moderationRepository.NewPostgresRepository(db)
	moderationSvc := moderationService.NewModerationService(moderationRepo)
	moderationHandler := moderationHandler.NewModerationHandler(moderationSvc)

	// Initialize organization dependencies
	orgRepo := organizationRepository.NewPostgresRepository(db)
	orgSvc := organizationService.NewOrganizationService(orgRepo)
	orgHandler := organizationHandler.NewOrganizationHandler(orgSvc)

	// Setup router
	router := gin.New()
	api.SetupMiddleware(router)
	api.SetupRoutes(router, authHandler, profileHandler, feedbackHandler, notificationHandler, dashboardHandler, orgHandler, peopleHandler, communityHandler, accountHandler, moderationHandler, tokenGen)

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
