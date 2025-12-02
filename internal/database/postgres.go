package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// DB wraps the database connection pool
type DB struct {
	Pool *pgxpool.Pool
}

// Config holds database configuration
type Config struct {
	URL             string
	MaxConnections  int
	MaxIdleTime     time.Duration
	ConnMaxLifetime time.Duration
}

// New creates a new database connection pool
func New(ctx context.Context, cfg Config) (*DB, error) {
	ctx, span := otel.Tracer("database").Start(ctx, "database.New")
	defer span.End()

	poolConfig, err := pgxpool.ParseConfig(cfg.URL)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to parse database URL: %w", err)
	}

	poolConfig.MaxConns = int32(cfg.MaxConnections)
	poolConfig.MaxConnIdleTime = cfg.MaxIdleTime
	poolConfig.MaxConnLifetime = cfg.ConnMaxLifetime

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test the connection
	if err := pool.Ping(ctx); err != nil {
		span.RecordError(err)
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	span.SetAttributes(
		attribute.String("db.system", "postgresql"),
		attribute.Int("db.max_connections", cfg.MaxConnections),
	)

	return &DB{Pool: pool}, nil
}

// Close closes the database connection pool
func (db *DB) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
}

// HealthCheck performs a health check on the database
func (db *DB) HealthCheck(ctx context.Context) error {
	ctx, span := otel.Tracer("database").Start(ctx, "database.HealthCheck")
	defer span.End()

	if err := db.Pool.Ping(ctx); err != nil {
		span.RecordError(err)
		return fmt.Errorf("database health check failed: %w", err)
	}

	return nil
}

// WithTracing wraps a database operation with tracing
func WithTracing(ctx context.Context, operation string, fn func(context.Context) error) error {
	ctx, span := otel.Tracer("database").Start(ctx, operation)
	defer span.End()

	err := fn(ctx)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	} else {
		span.SetStatus(codes.Ok, "")
	}

	return err
}

