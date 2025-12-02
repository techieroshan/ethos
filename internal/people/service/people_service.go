package service

import (
	"context"

	"ethos/internal/auth/model"
)

// Service defines the interface for people search business logic
type Service interface {
	// SearchPeople searches for people
	SearchPeople(ctx context.Context, query string, limit, offset int) ([]*model.UserProfile, int, error)

	// GetRecommendations gets people recommendations
	GetRecommendations(ctx context.Context, userID string) ([]*model.UserProfile, error)
}

