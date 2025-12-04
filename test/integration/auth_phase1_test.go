package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"ethos/internal/auth/handler"
	"ethos/internal/auth/model"
	"ethos/internal/auth/repository"
	"ethos/internal/auth/service"
	"ethos/internal/database"
	"ethos/pkg/jwt"
	"ethos/test/testutil"
)

// TestVerifyEmailSuccess verifies that a valid email verification token marks email as verified
func TestVerifyEmailSuccess(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	user := createTestUser(t, db, "verify@example.com", false)
	verificationToken := fmt.Sprintf("verify_%s", user.ID)

	authRepo := repository.NewPostgresRepository(db)
	tokenGen := jwt.NewTokenGenerator("secret", "secret", 15*time.Minute, 14*24*time.Hour)
	authService := service.NewAuthService(authRepo, tokenGen, nil, nil)
	authHandler := handler.NewAuthHandler(authService)

	ginCtx, w := createGinContext(http.MethodGet, "/api/v1/auth/verify-email/"+verificationToken, nil)
	ginCtx.Params = append(ginCtx.Params, gin.Param{Key: "token", Value: verificationToken})
	authHandler.VerifyEmail(ginCtx)

	assert.Equal(t, http.StatusOK, w.Code)
}

// TestVerifyEmailInvalidToken verifies that an invalid token returns 401
func TestVerifyEmailInvalidToken(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	authRepo := repository.NewPostgresRepository(db)
	tokenGen := jwt.NewTokenGenerator("secret", "secret", 15*time.Minute, 14*24*time.Hour)
	authService := service.NewAuthService(authRepo, tokenGen, nil, nil)
	authHandler := handler.NewAuthHandler(authService)

	ginCtx, w := createGinContext(http.MethodGet, "/api/v1/auth/verify-email/invalid_token", nil)
	ginCtx.Params = append(ginCtx.Params, gin.Param{Key: "token", Value: "invalid_token"})
	authHandler.VerifyEmail(ginCtx)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	var errResp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &errResp)
	require.NoError(t, err)
	assert.Equal(t, "AUTH_TOKEN_INVALID", errResp["code"])
}

// TestVerifyEmailAlreadyVerified verifies that verifying an already verified email is idempotent
func TestVerifyEmailAlreadyVerified(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	user := createTestUser(t, db, "verified@example.com", true)
	verificationToken := fmt.Sprintf("verify_%s", user.ID)

	authRepo := repository.NewPostgresRepository(db)
	tokenGen := jwt.NewTokenGenerator("secret", "secret", 15*time.Minute, 14*24*time.Hour)
	authService := service.NewAuthService(authRepo, tokenGen, nil, nil)
	authHandler := handler.NewAuthHandler(authService)

	ginCtx, w := createGinContext(http.MethodGet, "/api/v1/auth/verify-email/"+verificationToken, nil)
	ginCtx.Params = append(ginCtx.Params, gin.Param{Key: "token", Value: verificationToken})
	authHandler.VerifyEmail(ginCtx)

	assert.Equal(t, http.StatusOK, w.Code)
}

// TestChangePasswordSuccess verifies password can be changed with correct current password
func TestChangePasswordSuccess(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	user := createTestUser(t, db, "change@example.com", true)
	accessToken := generateTestJWT(t, user.ID)

	authRepo := repository.NewPostgresRepository(db)
	tokenGen := jwt.NewTokenGenerator("secret", "secret", 15*time.Minute, 14*24*time.Hour)
	authService := service.NewAuthService(authRepo, tokenGen, nil, nil)
	authHandler := handler.NewAuthHandler(authService)

	changeReq := map[string]string{
		"current_password": "oldpassword",
		"new_password":     "newpassword123",
	}
	body, _ := json.Marshal(changeReq)
	ginCtx, w := createGinContext(http.MethodPost, "/api/v1/auth/change-password", bytes.NewReader(body))
	ginCtx.Request.Header.Set("Authorization", "Bearer "+accessToken)
	ginCtx.Set("user_id", user.ID)
	authHandler.ChangePassword(ginCtx)

	assert.Equal(t, http.StatusOK, w.Code)
}

// TestChangePasswordWrongCurrent verifies that wrong current password fails
func TestChangePasswordWrongCurrent(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	user := createTestUser(t, db, "wrong@example.com", true)
	accessToken := generateTestJWT(t, user.ID)

	authRepo := repository.NewPostgresRepository(db)
	tokenGen := jwt.NewTokenGenerator("secret", "secret", 15*time.Minute, 14*24*time.Hour)
	authService := service.NewAuthService(authRepo, tokenGen, nil, nil)
	authHandler := handler.NewAuthHandler(authService)

	changeReq := map[string]string{
		"current_password": "wrongpassword",
		"new_password":     "newpassword123",
	}
	body, _ := json.Marshal(changeReq)
	ginCtx, w := createGinContext(http.MethodPost, "/api/v1/auth/change-password", bytes.NewReader(body))
	ginCtx.Request.Header.Set("Authorization", "Bearer "+accessToken)
	ginCtx.Set("user_id", user.ID)
	authHandler.ChangePassword(ginCtx)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestChangePasswordUnauthorized verifies that unauthenticated user cannot change password
func TestChangePasswordUnauthorized(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	authRepo := repository.NewPostgresRepository(db)
	tokenGen := jwt.NewTokenGenerator("secret", "secret", 15*time.Minute, 14*24*time.Hour)
	authService := service.NewAuthService(authRepo, tokenGen, nil, nil)
	authHandler := handler.NewAuthHandler(authService)

	changeReq := map[string]string{
		"current_password": "oldpassword",
		"new_password":     "newpassword123",
	}
	body, _ := json.Marshal(changeReq)
	ginCtx, w := createGinContext(http.MethodPost, "/api/v1/auth/change-password", bytes.NewReader(body))
	authHandler.ChangePassword(ginCtx)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestSetup2FASuccess verifies 2FA can be set up
func TestSetup2FASuccess(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	user := createTestUser(t, db, "2fa@example.com", true)
	accessToken := generateTestJWT(t, user.ID)

	authRepo := repository.NewPostgresRepository(db)
	tokenGen := jwt.NewTokenGenerator("secret", "secret", 15*time.Minute, 14*24*time.Hour)
	authService := service.NewAuthService(authRepo, tokenGen, nil, nil)
	authHandler := handler.NewAuthHandler(authService)

	setup2FAReq := map[string]string{
		"password": "oldpassword",
	}
	body, _ := json.Marshal(setup2FAReq)
	ginCtx, w := createGinContext(http.MethodPost, "/api/v1/auth/setup-2fa", bytes.NewReader(body))
	ginCtx.Request.Header.Set("Authorization", "Bearer "+accessToken)
	ginCtx.Set("user_id", user.ID)
	authHandler.Setup2FA(ginCtx)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.NotEmpty(t, resp["secret"])
	assert.NotEmpty(t, resp["qr_code"])
}

// TestSetup2FAUnauthorized verifies that unauthenticated user cannot setup 2FA
func TestSetup2FAUnauthorized(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	authRepo := repository.NewPostgresRepository(db)
	tokenGen := jwt.NewTokenGenerator("secret", "secret", 15*time.Minute, 14*24*time.Hour)
	authService := service.NewAuthService(authRepo, tokenGen, nil, nil)
	authHandler := handler.NewAuthHandler(authService)

	setup2FAReq := map[string]string{
		"password": "oldpassword",
	}
	body, _ := json.Marshal(setup2FAReq)
	ginCtx, w := createGinContext(http.MethodPost, "/api/v1/auth/setup-2fa", bytes.NewReader(body))
	authHandler.Setup2FA(ginCtx)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// Helper functions

func setupTestDB(t *testing.T) *database.DB {
	t.Helper()
	dsn := "postgresql://postgres:postgres@localhost:5432/ethos_test"
	pool := testutil.SetupTestDB(t, dsn)
	return &database.DB{Pool: pool}
}

func createTestUser(t *testing.T, db *database.DB, email string, verified bool) *model.User {
	t.Helper()
	ctx := context.Background()

	user := &model.User{
		ID:            generateID(),
		Email:         email,
		PasswordHash:  hashPassword("oldpassword"),
		EmailVerified: verified,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	repo := repository.NewPostgresRepository(db)
	err := repo.CreateUser(ctx, user)
	require.NoError(t, err)

	return user
}

func createGinContext(method, path string, body interface{}) (*gin.Context, *httptest.ResponseRecorder) {
	var req *http.Request
	var err error

	if body != nil {
		bodyBytes, _ := json.Marshal(body)
		req, err = http.NewRequest(method, path, bytes.NewReader(bodyBytes))
	} else {
		req, err = http.NewRequest(method, path, nil)
	}

	if err != nil {
		panic(err)
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	return c, w
}

func generateTestJWT(t *testing.T, userID string) string {
	t.Helper()
	tokenGen := jwt.NewTokenGenerator("secret", "secret", 15*time.Minute, 14*24*time.Hour)
	token, err := tokenGen.GenerateAccessToken(userID)
	require.NoError(t, err)
	return token
}

func hashPassword(password string) string {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hashed)
}

func generateID() string {
	return "user_" + fmt.Sprintf("%d", time.Now().UnixNano())
}
