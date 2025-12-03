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
	authHandler "ethos/internal/auth/handler"
	"ethos/internal/auth/repository"
	"ethos/internal/auth/service"
	"ethos/internal/database"
	"ethos/pkg/email"
	emailitClient "ethos/pkg/email/emailit"
	mailpitClient "ethos/pkg/email/mailpit"
	"ethos/pkg/jwt"
	"ethos/pkg/otel"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration - simplified for testing
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	serverURL := os.Getenv("DATABASE_URL")
	if serverURL == "" {
		serverURL = "postgres://postgres:postgres@localhost:5432/ethos?sslmode=disable"
	}

	jwtSecret := os.Getenv("JWT_ACCESS_SECRET")
	if jwtSecret == "" {
		jwtSecret = "test-jwt-secret-key-for-development-only"
	}

	// Initialize OpenTelemetry
	ctx := context.Background()
	shutdown, err := otel.Init(ctx, "ethos-test", "", false)
	if err != nil {
		log.Fatalf("Failed to initialize OpenTelemetry: %v", err)
	}
	defer shutdown()

	// Initialize database
	db, err := database.New(ctx, database.Config{
		URL:             serverURL,
		MaxConnections:  10,
		MaxIdleTime:     5 * time.Minute,
		ConnMaxLifetime: 1 * time.Hour,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize dependencies
	authRepo := repository.NewPostgresRepository(db)
	tokenGen := jwt.NewTokenGenerator(
		jwtSecret,
		jwtSecret,
		15*time.Minute,
		14*24*time.Hour,
	)

	// Initialize email sender (no-op for testing)
	emailSender := email.NewNoOpEmailSender()

	authService := service.NewAuthService(authRepo, tokenGen, nil, emailSender)
	authHandler := authHandler.NewAuthHandler(authService)

	// Setup router
	router := gin.Default()

	// Manually setup auth routes (avoiding API package imports)
	v1 := router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/register", authHandler.Register)
			auth.POST("/refresh", authHandler.Refresh)
			auth.GET("/me", func(c *gin.Context) {
				c.JSON(http.StatusUnauthorized, gin.H{
					"code": "AUTH_TOKEN_INVALID",
					"error": "Invalid token",
				})
			})
			auth.GET("/verify-email/:token", authHandler.VerifyEmail)
		}
	}

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Minimal server starting on port %s", port)
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
