package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

// OptimizedDB provides optimized database operations with connection pooling
type OptimizedDB struct {
	pool *pgxpool.Pool
	sqlDB *sql.DB
}

// NewOptimizedDB creates a new optimized database connection
func NewOptimizedDB(databaseURL string) (*OptimizedDB, error) {
	// Parse the connection config
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database URL: %w", err)
	}

	// Optimize connection pool settings
	config.MaxConns = 25                    // Maximum connections
	config.MinConns = 5                     // Minimum connections
	config.MaxConnLifetime = 1 * time.Hour  // Maximum connection lifetime
	config.MaxConnIdleTime = 30 * time.Minute // Maximum idle time
	config.HealthCheckPeriod = 1 * time.Minute // Health check frequency

	// Create the connection pool
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Create sql.DB wrapper for compatibility
	sqlDB := stdlib.OpenDBFromPool(pool)

	// Configure sql.DB settings
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(1 * time.Hour)
	sqlDB.SetConnMaxIdleTime(30 * time.Minute)

	return &OptimizedDB{
		pool:  pool,
		sqlDB: sqlDB,
	}, nil
}

// Pool returns the underlying pgx pool
func (db *OptimizedDB) Pool() *pgxpool.Pool {
	return db.pool
}

// SQLDB returns the sql.DB wrapper
func (db *OptimizedDB) SQLDB() *sql.DB {
	return db.sqlDB
}

// Close closes the database connections
func (db *OptimizedDB) Close() {
	db.pool.Close()
	db.sqlDB.Close()
}

// Ping tests the database connection
func (db *OptimizedDB) Ping(ctx context.Context) error {
	return db.pool.Ping(ctx)
}

// Stats returns connection pool statistics
func (db *OptimizedDB) Stats() *pgxpool.Stat {
	return db.pool.Stat()
}

// HealthCheck performs a comprehensive health check
func (db *OptimizedDB) HealthCheck(ctx context.Context) error {
	// Test basic connectivity
	if err := db.Ping(ctx); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	// Test a simple query
	var result int
	err := db.pool.QueryRow(ctx, "SELECT 1").Scan(&result)
	if err != nil {
		return fmt.Errorf("database query test failed: %w", err)
	}

	if result != 1 {
		return fmt.Errorf("database query test returned unexpected result: %d", result)
	}

	return nil
}

// ExecuteInTransaction executes a function within a database transaction
func (db *OptimizedDB) ExecuteInTransaction(ctx context.Context, fn func(tx pgx.Tx) error) error {
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := fn(tx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// QueryWithTimeout executes a query with a timeout
func (db *OptimizedDB) QueryWithTimeout(ctx context.Context, timeout time.Duration, query string, args ...interface{}) (pgx.Rows, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return db.pool.Query(ctx, query, args...)
}

// QueryRowWithTimeout executes a row query with a timeout
func (db *OptimizedDB) QueryRowWithTimeout(ctx context.Context, timeout time.Duration, query string, args ...interface{}) pgx.Row {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return db.pool.QueryRow(ctx, query, args...)
}

// ExecuteWithTimeout executes a command with a timeout
func (db *OptimizedDB) ExecuteWithTimeout(ctx context.Context, timeout time.Duration, query string, args ...interface{}) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	_, err := db.pool.Exec(ctx, query, args...)
	return err
}

// PreparedStatement holds a prepared statement
type PreparedStatement struct {
	Name string
	SQL  string
}

// PrepareStatements prepares commonly used statements for better performance
func (db *OptimizedDB) PrepareStatements(ctx context.Context, statements []PreparedStatement) error {
	conn, err := db.pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("failed to acquire connection: %w", err)
	}
	defer conn.Release()

	for _, stmt := range statements {
		if _, err := conn.Conn().Prepare(ctx, stmt.Name, stmt.SQL); err != nil {
			return fmt.Errorf("failed to prepare statement %s: %w", stmt.Name, err)
		}
	}

	return nil
}

// Common prepared statements for the application
var CommonStatements = []PreparedStatement{
	{
		Name: "get_user_by_id",
		SQL:  "SELECT id, email, password_hash, name, email_verified, public_bio, created_at, updated_at FROM users WHERE id = $1",
	},
	{
		Name: "get_user_with_roles",
		SQL: `SELECT u.id, u.email, u.password_hash, u.name, u.email_verified, u.public_bio, u.created_at, u.updated_at,
			r.id, r.name, r.description, ura.assigned_at, ura.is_active
			FROM users u
			LEFT JOIN user_role_assignments ura ON u.id = ura.user_id AND ura.is_active = TRUE
			LEFT JOIN roles r ON ura.role_id = r.id
			WHERE u.id = $1`,
	},
	{
		Name: "get_feedback_by_id",
		SQL:  "SELECT id, author_id, recipient_id, content, rating, visibility, created_at, updated_at FROM feedback WHERE id = $1",
	},
	{
		Name: "get_tenant_memberships",
		SQL:  "SELECT tenant_id, role, joined_at, is_active FROM user_tenant_memberships WHERE user_id = $1 AND is_active = TRUE",
	},
}

// InitializePreparedStatements prepares common statements on startup
func (db *OptimizedDB) InitializePreparedStatements(ctx context.Context) error {
	return db.PrepareStatements(ctx, CommonStatements)
}

// BatchInsert performs batch insert operations for better performance
func (db *OptimizedDB) BatchInsert(ctx context.Context, table string, columns []string, values [][]interface{}) error {
	if len(values) == 0 {
		return nil
	}

	// Build the INSERT statement
	query := fmt.Sprintf("INSERT INTO %s (", table)
	for i, col := range columns {
		if i > 0 {
			query += ", "
		}
		query += col
	}
	query += ") VALUES "

	// Build value placeholders
	args := make([]interface{}, 0)
	for i, row := range values {
		if i > 0 {
			query += ", "
		}
		query += "("
		for j := range row {
			if j > 0 {
				query += ", "
			}
			query += fmt.Sprintf("$%d", len(args)+1)
			args = append(args, row[j])
		}
		query += ")"
	}

	return db.ExecuteWithTimeout(ctx, 30*time.Second, query, args...)
}

// GetQueryMetrics collects query performance metrics
type QueryMetrics struct {
	Query     string
	Duration  time.Duration
	RowCount  int
	Error     error
	Timestamp time.Time
}

// QueryMetricsCollector collects and stores query metrics
type QueryMetricsCollector struct {
	metrics []QueryMetrics
}

// NewQueryMetricsCollector creates a new metrics collector
func NewQueryMetricsCollector() *QueryMetricsCollector {
	return &QueryMetricsCollector{
		metrics: make([]QueryMetrics, 0),
	}
}

// RecordMetric records a query metric
func (mc *QueryMetricsCollector) RecordMetric(query string, duration time.Duration, rowCount int, err error) {
	metric := QueryMetrics{
		Query:     query,
		Duration:  duration,
		RowCount:  rowCount,
		Error:     err,
		Timestamp: time.Now(),
	}
	mc.metrics = append(mc.metrics, metric)
}

// GetMetrics returns collected metrics
func (mc *QueryMetricsCollector) GetMetrics() []QueryMetrics {
	return mc.metrics
}

// GetSlowQueries returns queries that took longer than the specified duration
func (mc *QueryMetricsCollector) GetSlowQueries(threshold time.Duration) []QueryMetrics {
	slowQueries := make([]QueryMetrics, 0)
	for _, metric := range mc.metrics {
		if metric.Duration > threshold {
			slowQueries = append(slowQueries, metric)
		}
	}
	return slowQueries
}

// ClearMetrics clears all collected metrics
func (mc *QueryMetricsCollector) ClearMetrics() {
	mc.metrics = make([]QueryMetrics, 0)
}
