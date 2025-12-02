package service

import (
	"context"

	"ethos/internal/auth/model"
	"ethos/internal/people/repository"
)

// PeopleService implements the Service interface
type PeopleService struct {
	repo repository.Repository
}

// NewPeopleService creates a new people service
func NewPeopleService(repo repository.Repository) Service {
	return &PeopleService{
		repo: repo,
	}
}

// SearchPeople searches for people
func (s *PeopleService) SearchPeople(ctx context.Context, query string, limit, offset int) ([]*model.UserProfile, int, error) {
	profiles, count, err := s.repo.SearchPeople(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return profiles, count, nil
}

// GetRecommendations gets people recommendations
func (s *PeopleService) GetRecommendations(ctx context.Context, userID string) ([]*model.UserProfile, error) {
	recommendations, err := s.repo.GetRecommendations(ctx, userID)
	if err != nil {
		return nil, err
	}

	return recommendations, nil
}

