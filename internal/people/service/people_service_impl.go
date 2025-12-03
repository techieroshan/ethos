package service

import (
	"context"

	"ethos/internal/auth/model"
	"ethos/internal/people"
	"ethos/internal/people/repository"
)

// PeopleService implements the Service interface
type PeopleService struct {
	client PeopleClient // Can be REST or gRPC client
}

// NewPeopleService creates a new people service with REST client
func NewPeopleService(repo repository.Repository) Service {
	return &PeopleService{
		client: NewRESTPeopleClient(repo),
	}
}

// NewPeopleServiceWithClient creates a people service with a custom client (REST or gRPC)
func NewPeopleServiceWithClient(client PeopleClient) Service {
	return &PeopleService{
		client: client,
	}
}

// SearchPeople searches for people
func (s *PeopleService) SearchPeople(ctx context.Context, query string, limit, offset int) ([]*model.UserProfile, int, error) {
	return s.client.SearchPeople(ctx, query, limit, offset)
}

// SearchPeopleWithFilters searches for people with enhanced filtering
func (s *PeopleService) SearchPeopleWithFilters(ctx context.Context, query string, limit, offset int, filters *people.PeopleSearchFilters) ([]*model.UserProfile, int, error) {
	return s.client.SearchPeopleWithFilters(ctx, query, limit, offset, filters)
}

// GetRecommendations gets people recommendations
func (s *PeopleService) GetRecommendations(ctx context.Context, userID string) ([]*model.UserProfile, error) {
	return s.client.GetRecommendations(ctx, userID)
}

