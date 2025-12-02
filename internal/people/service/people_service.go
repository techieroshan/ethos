package service

import (
	"context"

	"ethos/internal/auth/model"
)

// PeopleSearchFilters represents filtering options for people search
type PeopleSearchFilters struct {
	ReviewerType *string  `json:"reviewer_type,omitempty"` // "public" or "org"
	Context      *string  `json:"context,omitempty"`       // e.g., "project", "team", "initiative"
	Verification *string  `json:"verification,omitempty"`  // "verified", "unverified"
	Tags         []string `json:"tags,omitempty"`          // Comma-separated tags
}

// Service defines the interface for people search business logic
type Service interface {
	// SearchPeople searches for people
	SearchPeople(ctx context.Context, query string, limit, offset int) ([]*model.UserProfile, int, error)

	// SearchPeopleWithFilters searches for people with enhanced filtering
	SearchPeopleWithFilters(ctx context.Context, query string, limit, offset int, filters *PeopleSearchFilters) ([]*model.UserProfile, int, error)

	// GetRecommendations gets people recommendations
	GetRecommendations(ctx context.Context, userID string) ([]*model.UserProfile, error)
}

