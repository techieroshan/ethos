package integration

import (
	"context"
	"testing"
	"time"

	"ethos/internal/database"
	"ethos/internal/organization/repository"
	"ethos/internal/organization/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestServiceLayerFoundation tests basic service layer functionality
// Used as foundation for E2E tests
func TestServiceLayerFoundation(t *testing.T) {
	ctx := context.Background()
	db, err := database.New(ctx, database.Config{
		URL:             "postgres://ethos:ethos@localhost:5432/ethos_test",
		MaxConnections:  10,
		MaxIdleTime:     30 * time.Second,
		ConnMaxLifetime: 1 * time.Hour,
	})
	if err != nil {
		t.Skipf("Failed to connect to test database: %v", err)
	}
	defer db.Close()

	// Initialize repository and service
	contextRepo := repository.NewPostgresContextRepository(db)
	contextService := service.NewUserContextService(contextRepo)

	// Verify that service implements the interface
	require.NotNil(t, contextService)

	// Test user ID for foundation tests
	testUserID := "foundation_user_" + time.Now().Format("20060102150405")
	testOrgID := "foundation_org_" + time.Now().Format("20060102150405")

	t.Run("Service methods are callable and return expected types", func(t *testing.T) {
		// STEP 1: GetAvailableContexts returns slice
		contexts, err := contextService.GetAvailableContexts(ctx, testUserID)
		assert.NotNil(t, contexts, "GetAvailableContexts should return a slice")
		_ = err // May error if user doesn't exist

		// STEP 2: GetCurrentContext returns context or nil
		currentCtx, err := contextService.GetCurrentContext(ctx, testUserID)
		_ = err
		_ = currentCtx // May be nil for non-existent user

		// STEP 3: ValidateUserInOrganization returns boolean
		isMember, err := contextService.ValidateUserInOrganization(ctx, testUserID, testOrgID)
		require.NoError(t, err)
		assert.False(t, isMember, "Non-existent user should not be member")

		// STEP 4: GetUserRoleInOrganization returns role string
		role, err := contextService.GetUserRoleInOrganization(ctx, testUserID, testOrgID)
		_ = err
		assert.Equal(t, "", role, "Non-member should have empty role")

		// STEP 5: GetContextSwitchHistory returns records
		records, err := contextService.GetContextSwitchHistory(ctx, testUserID, 50, 0)
		require.NoError(t, err)
		assert.NotNil(t, records, "Should return history records")
	})

	t.Run("Multi-tenant isolation at service layer", func(t *testing.T) {
		// STEP 1: Non-member users cannot access organizations
		restrictedOrgID := "org_not_for_user_" + time.Now().Format("20060102150405")
		isMember, err := contextService.ValidateUserInOrganization(ctx, testUserID, restrictedOrgID)
		require.NoError(t, err)
		assert.False(t, isMember, "Non-member should be rejected")

		// STEP 2: Role checks fail for non-members
		role, err := contextService.GetUserRoleInOrganization(ctx, testUserID, restrictedOrgID)
		_ = err
		assert.Equal(t, "", role, "Non-member role should be empty")
	})
}
