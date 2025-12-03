package people

import (
	"context"

	"ethos/internal/auth/model"
)

// PeopleSearchFilters represents filtering options for people search
type PeopleSearchFilters struct {
	OrganizationID *string  `json:"organization_id,omitempty"`
	Role           *string  `json:"role,omitempty"`
	Skills         []string `json:"skills,omitempty"`
	Location       *string  `json:"location,omitempty"`
	Department     *string  `json:"department,omitempty"`
	ReviewerType   *string  `json:"reviewer_type,omitempty"` // "public" or "org"
	Context        *string  `json:"context,omitempty"`       // e.g., "project", "team", "initiative"
	Verification   *string  `json:"verification,omitempty"`  // "verified", "unverified"
	Tags           []string `json:"tags,omitempty"`          // Comma-separated tags
}

// Repository defines the interface for people data access
type Repository interface {
	// SearchPeople searches for people with optional filtering
	SearchPeople(ctx context.Context, query string, limit, offset int) ([]*model.UserProfile, int, error)

	// SearchPeopleWithFilters searches for people with enhanced filtering
	SearchPeopleWithFilters(ctx context.Context, query string, limit, offset int, filters *PeopleSearchFilters) ([]*model.UserProfile, int, error)

	// GetRecommendations gets people recommendations for a user
	GetRecommendations(ctx context.Context, userID string, limit, offset int) ([]*model.UserProfile, int, error)
}
