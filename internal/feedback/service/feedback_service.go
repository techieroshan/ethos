package service

import (
	"context"

	"ethos/internal/feedback/model"
)

// CreateFeedbackRequest represents a request to create feedback
type CreateFeedbackRequest struct {
	Content    string                  `json:"content" binding:"required"`
	Type       *model.FeedbackType     `json:"type,omitempty"`
	Visibility *model.FeedbackVisibility `json:"visibility,omitempty"`
}

// CreateCommentRequest represents a request to create a comment
type CreateCommentRequest struct {
	Content         string `json:"content" binding:"required"`
	ParentCommentID *string `json:"parent_comment_id,omitempty"`
}

// AddReactionRequest represents a request to add a reaction
type AddReactionRequest struct {
	ReactionType string `json:"reaction_type" binding:"required"`
}

// Service defines the interface for feedback business logic
type Service interface {
	// GetFeed retrieves a paginated feed of feedback items
	GetFeed(ctx context.Context, limit, offset int) ([]*model.FeedbackItem, int, error)

	// GetFeedbackByID retrieves a feedback item by ID
	GetFeedbackByID(ctx context.Context, feedbackID string) (*model.FeedbackItem, error)

	// GetComments retrieves comments for a feedback item
	GetComments(ctx context.Context, feedbackID string, limit, offset int) ([]*model.FeedbackComment, int, error)

	// CreateFeedback creates a new feedback item
	CreateFeedback(ctx context.Context, userID string, req *CreateFeedbackRequest) (*model.FeedbackItem, error)

	// CreateComment creates a new comment on a feedback item
	CreateComment(ctx context.Context, userID, feedbackID string, req *CreateCommentRequest) (*model.FeedbackComment, error)

	// AddReaction adds a reaction to a feedback item
	AddReaction(ctx context.Context, userID, feedbackID string, reactionType string) error

	// RemoveReaction removes a reaction from a feedback item
	RemoveReaction(ctx context.Context, userID, feedbackID string, reactionType string) error
}

