package service

import (
	"context"

	dashModel "ethos/internal/dashboard/model"
	"ethos/internal/dashboard/repository"
	dashboardpb "ethos/api/proto/dashboard"
	"ethos/pkg/grpc/converter"
)

// DashboardClient defines the interface for dashboard data access (REST or gRPC)
type DashboardClient interface {
	// GetDashboard retrieves a dashboard snapshot for a user
	GetDashboard(ctx context.Context, userID string) (*dashModel.DashboardSnapshot, error)
}

// RESTDashboardClient implements DashboardClient using REST (current repository)
type RESTDashboardClient struct {
	repo repository.Repository
}

// NewRESTDashboardClient creates a new REST dashboard client
func NewRESTDashboardClient(repo repository.Repository) DashboardClient {
	return &RESTDashboardClient{repo: repo}
}

// GetDashboard implements DashboardClient interface using REST
func (c *RESTDashboardClient) GetDashboard(ctx context.Context, userID string) (*dashModel.DashboardSnapshot, error) {
	return c.repo.GetDashboard(ctx, userID)
}

// GRPCDashboardClient implements DashboardClient using gRPC
type GRPCDashboardClient struct {
	client dashboardpb.DashboardServiceClient
}

// NewGRPCDashboardClient creates a new gRPC dashboard client
func NewGRPCDashboardClient(client dashboardpb.DashboardServiceClient) DashboardClient {
	return &GRPCDashboardClient{client: client}
}

// GetDashboard implements DashboardClient interface using gRPC
func (c *GRPCDashboardClient) GetDashboard(ctx context.Context, userID string) (*dashModel.DashboardSnapshot, error) {
	req := &dashboardpb.GetDashboardRequest{
		UserId: userID,
	}

	resp, err := c.client.GetDashboard(ctx, req)
	if err != nil {
		return nil, err
	}

	return converter.ProtoToDashboardSnapshot(resp.Snapshot), nil
}

