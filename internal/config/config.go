package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Cache    CacheConfig
	JWT      JWTConfig
	OTEL     OTELConfig
	Checker  CheckerConfig
	Emailit  EmailitConfig
	Mailpit  MailpitConfig
	GRPC     GRPCConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// DatabaseConfig holds database connection configuration
type DatabaseConfig struct {
	URL             string
	MaxConnections  int
	MaxIdleTime     time.Duration
	ConnMaxLifetime time.Duration
}

// CacheConfig holds Redis cache configuration
type CacheConfig struct {
	URL      string
	Password string
	DB       int
	Enabled  bool
}

// JWTConfig holds JWT token configuration
type JWTConfig struct {
	AccessTokenSecret  string
	RefreshTokenSecret string
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
}

// OTELConfig holds OpenTelemetry configuration
type OTELConfig struct {
	ServiceName string
	JaegerURL   string
	Enabled     bool
}

// CheckerConfig holds Checker API configuration for email validation
type CheckerConfig struct {
	APIKey  string
	BaseURL string
	Timeout time.Duration
	Retries int
}

// EmailitConfig holds Emailit API configuration for email sending
type EmailitConfig struct {
	APIKey  string
	BaseURL string
	Timeout time.Duration
	Retries int
}

// MailpitConfig holds Mailpit configuration for local email testing
type MailpitConfig struct {
	SMTPHost  string
	SMTPPort  int
	FromEmail string
	Enabled   bool
}

// GRPCConfig holds gRPC client configuration
type GRPCConfig struct {
	FeedbackEndpoint      string
	DashboardEndpoint     string
	NotificationsEndpoint string
	PeopleEndpoint        string
	Timeout               time.Duration
	Retries               int
	Enabled               bool
	// Protocol selection per service (defaults to "rest" if not set)
	FeedbackProtocol      string // "rest" or "grpc"
	DashboardProtocol     string // "rest" or "grpc"
	NotificationsProtocol string // "rest" or "grpc"
	PeopleProtocol         string // "rest" or "grpc"
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Try to load .env file, but don't fail if it doesn't exist
	_ = godotenv.Load()

	cfg := &Config{
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", "8000"),
			ReadTimeout:  getDurationEnv("SERVER_READ_TIMEOUT", 15*time.Second),
			WriteTimeout: getDurationEnv("SERVER_WRITE_TIMEOUT", 15*time.Second),
		},
		Database: DatabaseConfig{
			URL:             getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/ethos?sslmode=disable"),
			MaxConnections:  getIntEnv("DB_MAX_CONNECTIONS", 25),
			MaxIdleTime:     getDurationEnv("DB_MAX_IDLE_TIME", 5*time.Minute),
			ConnMaxLifetime: getDurationEnv("DB_CONN_MAX_LIFETIME", 1*time.Hour),
		},
		Cache: CacheConfig{
			URL:      getEnv("REDIS_URL", "redis://localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getIntEnv("REDIS_DB", 0),
			Enabled:  getBoolEnv("REDIS_ENABLED", true),
		},
		JWT: JWTConfig{
			AccessTokenSecret:  getEnv("JWT_ACCESS_SECRET", "your-access-secret-key-change-in-production"),
			RefreshTokenSecret: getEnv("JWT_REFRESH_SECRET", "your-refresh-secret-key-change-in-production"),
			AccessTokenExpiry:  getDurationEnv("JWT_ACCESS_EXPIRY", 15*time.Minute),
			RefreshTokenExpiry: getDurationEnv("JWT_REFRESH_EXPIRY", 14*24*time.Hour),
		},
		OTEL: OTELConfig{
			ServiceName: getEnv("OTEL_SERVICE_NAME", "ethos-api"),
			JaegerURL:   getEnv("JAEGER_URL", "http://localhost:14268/api/traces"),
			Enabled:     getBoolEnv("OTEL_ENABLED", false),
		},
		Checker: CheckerConfig{
			APIKey:  getEnv("CHECKER_API_KEY", ""),
			BaseURL: getEnv("CHECKER_BASE_URL", "https://api.checker.com"),
			Timeout: getDurationEnv("CHECKER_TIMEOUT", 5*time.Second),
			Retries: getIntEnv("CHECKER_RETRIES", 2),
		},
		Emailit: EmailitConfig{
			APIKey:  getEnv("EMAILIT_API_KEY", ""),
			BaseURL: getEnv("EMAILIT_BASE_URL", "https://api.emailit.com"),
			Timeout: getDurationEnv("EMAILIT_TIMEOUT", 10*time.Second),
			Retries: getIntEnv("EMAILIT_RETRIES", 2),
		},
		Mailpit: MailpitConfig{
			SMTPHost:  getEnv("MAILPIT_SMTP_HOST", "localhost"),
			SMTPPort:  getIntEnv("MAILPIT_SMTP_PORT", 1025),
			FromEmail: getEnv("MAILPIT_FROM_EMAIL", "noreply@ethos.test"),
			Enabled:   getBoolEnv("MAILPIT_ENABLED", false),
		},
		GRPC: GRPCConfig{
			FeedbackEndpoint:      getEnv("GRPC_FEEDBACK_ENDPOINT", "localhost:50051"),
			DashboardEndpoint:     getEnv("GRPC_DASHBOARD_ENDPOINT", "localhost:50052"),
			NotificationsEndpoint: getEnv("GRPC_NOTIFICATIONS_ENDPOINT", "localhost:50053"),
			PeopleEndpoint:        getEnv("GRPC_PEOPLE_ENDPOINT", "localhost:50054"),
			Timeout:               getDurationEnv("GRPC_TIMEOUT", 5*time.Second),
			Retries:               getIntEnv("GRPC_RETRIES", 3),
			Enabled:               getBoolEnv("GRPC_ENABLED", false),
			FeedbackProtocol:      getEnv("GRPC_FEEDBACK_PROTOCOL", "rest"),
			DashboardProtocol:     getEnv("GRPC_DASHBOARD_PROTOCOL", "rest"),
			NotificationsProtocol: getEnv("GRPC_NOTIFICATIONS_PROTOCOL", "rest"),
			PeopleProtocol:         getEnv("GRPC_PEOPLE_PROTOCOL", "rest"),
		},
	}

	// Validate required fields
	if cfg.JWT.AccessTokenSecret == "" || cfg.JWT.RefreshTokenSecret == "" {
		return nil, fmt.Errorf("JWT secrets must be set")
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

