package repository

import (
	"context"

	"ethos/internal/auth/model"
	"ethos/internal/people/service"
)

// Repository defines the interface for people search data access
type Repository interface {
	// SearchPeople searches for people
	SearchPeople(ctx context.Context, query string, limit, offset int) ([]*model.UserProfile, int, error)

	// SearchPeopleWithFilters searches for people with enhanced filtering
	SearchPeopleWithFilters(ctx context.Context, query string, limit, offset int, filters *service.PeopleSearchFilters) ([]*model.UserProfile, int, error)

	// GetRecommendations gets people recommendations
	GetRecommendations(ctx context.Context, userID string) ([]*model.UserProfile, error)
}

