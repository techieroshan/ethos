package service

import (
	"context"

	fbModel "ethos/internal/feedback/model"
	"ethos/internal/feedback/repository"
	feedbackpb "ethos/api/proto/feedback"
	"ethos/pkg/grpc/converter"
)

// FeedbackClient defines the interface for feedback data access (REST or gRPC)
type FeedbackClient interface {
	// GetFeed retrieves a paginated feed of feedback posts
	GetFeed(ctx context.Context, limit, offset int) ([]*fbModel.FeedbackItem, int, error)
	
	// GetFeedbackByID retrieves a specific feedback item
	GetFeedbackByID(ctx context.Context, feedbackID string) (*fbModel.FeedbackItem, error)
	
	// GetComments retrieves comments for a feedback item
	GetComments(ctx context.Context, feedbackID string, limit, offset int) ([]*fbModel.FeedbackComment, int, error)
}

// RESTFeedbackClient implements FeedbackClient using REST (current repository)
type RESTFeedbackClient struct {
	repo repository.Repository
}

// NewRESTFeedbackClient creates a new REST feedback client
func NewRESTFeedbackClient(repo repository.Repository) FeedbackClient {
	return &RESTFeedbackClient{repo: repo}
}

// GetFeed implements FeedbackClient interface using REST
func (c *RESTFeedbackClient) GetFeed(ctx context.Context, limit, offset int) ([]*fbModel.FeedbackItem, int, error) {
	return c.repo.GetFeed(ctx, limit, offset)
}

// GetFeedbackByID implements FeedbackClient interface using REST
func (c *RESTFeedbackClient) GetFeedbackByID(ctx context.Context, feedbackID string) (*fbModel.FeedbackItem, error) {
	return c.repo.GetFeedbackByID(ctx, feedbackID)
}

// GetComments implements FeedbackClient interface using REST
func (c *RESTFeedbackClient) GetComments(ctx context.Context, feedbackID string, limit, offset int) ([]*fbModel.FeedbackComment, int, error) {
	return c.repo.GetComments(ctx, feedbackID, limit, offset)
}

// GRPCFeedbackClient implements FeedbackClient using gRPC
type GRPCFeedbackClient struct {
	client feedbackpb.FeedbackServiceClient
}

// NewGRPCFeedbackClient creates a new gRPC feedback client
func NewGRPCFeedbackClient(client feedbackpb.FeedbackServiceClient) FeedbackClient {
	return &GRPCFeedbackClient{client: client}
}

// GetFeed implements FeedbackClient interface using gRPC
func (c *GRPCFeedbackClient) GetFeed(ctx context.Context, limit, offset int) ([]*fbModel.FeedbackItem, int, error) {
	req := &feedbackpb.GetFeedRequest{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	resp, err := c.client.GetFeed(ctx, req)
	if err != nil {
		return nil, 0, err
	}

	items := make([]*fbModel.FeedbackItem, 0, len(resp.Results))
	for _, pbItem := range resp.Results {
		item := converter.ProtoToFeedbackItem(pbItem)
		if item != nil {
			items = append(items, item)
		}
	}

	return items, int(resp.Count), nil
}

// GetFeedbackByID implements FeedbackClient interface using gRPC
func (c *GRPCFeedbackClient) GetFeedbackByID(ctx context.Context, feedbackID string) (*fbModel.FeedbackItem, error) {
	req := &feedbackpb.GetFeedbackRequest{
		FeedbackId: feedbackID,
	}

	resp, err := c.client.GetFeedback(ctx, req)
	if err != nil {
		return nil, err
	}

	return converter.ProtoToFeedbackItem(resp.Feedback), nil
}

// GetComments implements FeedbackClient interface using gRPC
func (c *GRPCFeedbackClient) GetComments(ctx context.Context, feedbackID string, limit, offset int) ([]*fbModel.FeedbackComment, int, error) {
	req := &feedbackpb.GetCommentsRequest{
		FeedbackId: feedbackID,
		Limit:      int32(limit),
		Offset:     int32(offset),
	}

	resp, err := c.client.GetComments(ctx, req)
	if err != nil {
		return nil, 0, err
	}

	comments := make([]*fbModel.FeedbackComment, 0, len(resp.Comments))
	for _, pbComment := range resp.Comments {
		comment := converter.ProtoToFeedbackComment(pbComment)
		if comment != nil {
			comments = append(comments, comment)
		}
	}

	return comments, int(resp.Count), nil
}

