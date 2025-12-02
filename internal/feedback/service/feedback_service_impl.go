package service

import (
	"context"

	"ethos/internal/feedback/model"
	"ethos/internal/feedback/repository"
)

// FeedbackService implements the Service interface
type FeedbackService struct {
	client FeedbackClient // Can be REST or gRPC client
	repo   repository.Repository // Kept for write operations (CreateFeedback, CreateComment, AddReaction, RemoveReaction)
}

// NewFeedbackService creates a new feedback service with REST client
func NewFeedbackService(repo repository.Repository) Service {
	return &FeedbackService{
		client: NewRESTFeedbackClient(repo),
		repo:   repo,
	}
}

// NewFeedbackServiceWithClient creates a feedback service with a custom client (REST or gRPC)
func NewFeedbackServiceWithClient(client FeedbackClient, repo repository.Repository) Service {
	return &FeedbackService{
		client: client,
		repo:   repo,
	}
}

// GetFeed retrieves a paginated feed of feedback items
func (s *FeedbackService) GetFeed(ctx context.Context, limit, offset int) ([]*model.FeedbackItem, int, error) {
	return s.client.GetFeed(ctx, limit, offset)
}

// GetFeedbackByID retrieves a feedback item by ID
func (s *FeedbackService) GetFeedbackByID(ctx context.Context, feedbackID string) (*model.FeedbackItem, error) {
	return s.client.GetFeedbackByID(ctx, feedbackID)
}

// GetComments retrieves comments for a feedback item
func (s *FeedbackService) GetComments(ctx context.Context, feedbackID string, limit, offset int) ([]*model.FeedbackComment, int, error) {
	return s.client.GetComments(ctx, feedbackID, limit, offset)
}

// CreateFeedback creates a new feedback item
func (s *FeedbackService) CreateFeedback(ctx context.Context, userID string, req *CreateFeedbackRequest) (*model.FeedbackItem, error) {
	item, err := s.repo.CreateFeedback(ctx, userID, req.Content, req.Type, req.Visibility)
	if err != nil {
		return nil, err
	}

	return item, nil
}

// CreateComment creates a new comment on a feedback item
func (s *FeedbackService) CreateComment(ctx context.Context, userID, feedbackID string, req *CreateCommentRequest) (*model.FeedbackComment, error) {
	comment, err := s.repo.CreateComment(ctx, userID, feedbackID, req.Content, req.ParentCommentID)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

// AddReaction adds a reaction to a feedback item
func (s *FeedbackService) AddReaction(ctx context.Context, userID, feedbackID string, reactionType string) error {
	err := s.repo.AddReaction(ctx, userID, feedbackID, reactionType)
	if err != nil {
		return err
	}

	return nil
}

// RemoveReaction removes a reaction from a feedback item
func (s *FeedbackService) RemoveReaction(ctx context.Context, userID, feedbackID string, reactionType string) error {
	err := s.repo.RemoveReaction(ctx, userID, feedbackID, reactionType)
	if err != nil {
		return err
	}

	return nil
}

