package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims represents JWT claims
type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// TokenGenerator handles JWT token generation and validation
type TokenGenerator struct {
	accessSecret  []byte
	refreshSecret []byte
	accessExpiry  time.Duration
	refreshExpiry time.Duration
}

// NewTokenGenerator creates a new token generator
func NewTokenGenerator(accessSecret, refreshSecret string, accessExpiry, refreshExpiry time.Duration) *TokenGenerator {
	return &TokenGenerator{
		accessSecret:  []byte(accessSecret),
		refreshSecret: []byte(refreshSecret),
		accessExpiry:  accessExpiry,
		refreshExpiry: refreshExpiry,
	}
}

// GenerateAccessToken generates a new access token
func (tg *TokenGenerator) GenerateAccessToken(userID string) (string, error) {
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tg.accessExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(tg.accessSecret)
}

// GenerateRefreshToken generates a new refresh token
func (tg *TokenGenerator) GenerateRefreshToken(userID string) (string, error) {
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tg.refreshExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(tg.refreshSecret)
}

// ValidateAccessToken validates an access token and returns the user ID
func (tg *TokenGenerator) ValidateAccessToken(tokenString string) (string, error) {
	return tg.validateToken(tokenString, tg.accessSecret)
}

// ValidateRefreshToken validates a refresh token and returns the user ID
func (tg *TokenGenerator) ValidateRefreshToken(tokenString string) (string, error) {
	return tg.validateToken(tokenString, tg.refreshSecret)
}

func (tg *TokenGenerator) validateToken(tokenString string, secret []byte) (string, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return "", errors.New("token expired")
		}
		return "", fmt.Errorf("invalid token: %w", err)
	}

	if !token.Valid {
		return "", errors.New("invalid token")
	}

	return claims.UserID, nil
}

