package service

import (
	"context"

	authModel "ethos/internal/auth/model"
	"ethos/internal/people"
	"ethos/internal/people/repository"
	peoplepb "ethos/api/proto/people"
	"ethos/pkg/grpc/converter"
)

// PeopleClient defines the interface for people data access (REST or gRPC)
type PeopleClient interface {
	// SearchPeople searches for people
	SearchPeople(ctx context.Context, query string, limit, offset int) ([]*authModel.UserProfile, int, error)

	// SearchPeopleWithFilters searches for people with enhanced filtering
	SearchPeopleWithFilters(ctx context.Context, query string, limit, offset int, filters *people.PeopleSearchFilters) ([]*authModel.UserProfile, int, error)

	// GetRecommendations gets people recommendations
	GetRecommendations(ctx context.Context, userID string) ([]*authModel.UserProfile, error)
}

// RESTPeopleClient implements PeopleClient using REST (current repository)
type RESTPeopleClient struct {
	repo repository.Repository
}

// NewRESTPeopleClient creates a new REST people client
func NewRESTPeopleClient(repo repository.Repository) PeopleClient {
	return &RESTPeopleClient{repo: repo}
}

// SearchPeople implements PeopleClient interface using REST
func (c *RESTPeopleClient) SearchPeople(ctx context.Context, query string, limit, offset int) ([]*authModel.UserProfile, int, error) {
	return c.repo.SearchPeople(ctx, query, limit, offset)
}

// SearchPeopleWithFilters implements PeopleClient interface using REST
func (c *RESTPeopleClient) SearchPeopleWithFilters(ctx context.Context, query string, limit, offset int, filters *people.PeopleSearchFilters) ([]*authModel.UserProfile, int, error) {
	return c.repo.SearchPeopleWithFilters(ctx, query, limit, offset, filters)
}

// GetRecommendations implements PeopleClient interface using REST
func (c *RESTPeopleClient) GetRecommendations(ctx context.Context, userID string) ([]*authModel.UserProfile, error) {
	return c.repo.GetRecommendations(ctx, userID)
}

// GRPCPeopleClient implements PeopleClient using gRPC
type GRPCPeopleClient struct {
	client peoplepb.PeopleServiceClient
}

// NewGRPCPeopleClient creates a new gRPC people client
func NewGRPCPeopleClient(client peoplepb.PeopleServiceClient) PeopleClient {
	return &GRPCPeopleClient{client: client}
}

// SearchPeople implements PeopleClient interface using gRPC
func (c *GRPCPeopleClient) SearchPeople(ctx context.Context, query string, limit, offset int) ([]*authModel.UserProfile, int, error) {
	req := &peoplepb.SearchPeopleRequest{
		Query:  query,
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	resp, err := c.client.SearchPeople(ctx, req)
	if err != nil {
		return nil, 0, err
	}

	profiles := make([]*authModel.UserProfile, 0, len(resp.Results))
	for _, pbProfile := range resp.Results {
		profile := converter.ProtoToUserProfile(pbProfile)
		if profile != nil {
			profiles = append(profiles, profile)
		}
	}

	return profiles, int(resp.Count), nil
}

// SearchPeopleWithFilters implements PeopleClient interface using gRPC
func (c *GRPCPeopleClient) SearchPeopleWithFilters(ctx context.Context, query string, limit, offset int, filters *people.PeopleSearchFilters) ([]*authModel.UserProfile, int, error) {
	// For now, implement as basic search - gRPC proto would need to be updated for full filtering
	return c.SearchPeople(ctx, query, limit, offset)
}

// GetRecommendations implements PeopleClient interface using gRPC
func (c *GRPCPeopleClient) GetRecommendations(ctx context.Context, userID string) ([]*authModel.UserProfile, error) {
	req := &peoplepb.GetRecommendationsRequest{
		UserId: userID,
	}

	resp, err := c.client.GetRecommendations(ctx, req)
	if err != nil {
		return nil, err
	}

	profiles := make([]*authModel.UserProfile, 0, len(resp.Recommendations))
	for _, pbProfile := range resp.Recommendations {
		profile := converter.ProtoToUserProfile(pbProfile)
		if profile != nil {
			profiles = append(profiles, profile)
		}
	}

	return profiles, nil
}

