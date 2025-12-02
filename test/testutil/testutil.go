package testutil

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

// SetupTestDB creates a test database connection
func SetupTestDB(t *testing.T, dsn string) *pgxpool.Pool {
	t.Helper()

	pool, err := pgxpool.New(context.Background(), dsn)
	require.NoError(t, err)

	// Test connection
	err = pool.Ping(context.Background())
	require.NoError(t, err)

	t.Cleanup(func() {
		pool.Close()
	})

	return pool
}

// CleanupTestDB cleans up test data
func CleanupTestDB(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()

	ctx := context.Background()
	_, err := pool.Exec(ctx, "TRUNCATE TABLE refresh_tokens, users CASCADE")
	if err != nil {
		t.Logf("Failed to cleanup test database: %v", err)
	}
}

