package repository

import (
	"context"

	"ethos/internal/feedback/model"
)

// Repository defines the interface for feedback data access
type Repository interface {
	// GetFeed retrieves a paginated feed of feedback items
	GetFeed(ctx context.Context, limit, offset int) ([]*model.FeedbackItem, int, error)

	// GetFeedbackByID retrieves a feedback item by ID
	GetFeedbackByID(ctx context.Context, feedbackID string) (*model.FeedbackItem, error)

	// GetComments retrieves comments for a feedback item
	GetComments(ctx context.Context, feedbackID string, limit, offset int) ([]*model.FeedbackComment, int, error)

	// CreateFeedback creates a new feedback item
	CreateFeedback(ctx context.Context, userID string, content string, feedbackType *model.FeedbackType, visibility *model.FeedbackVisibility) (*model.FeedbackItem, error)

	// CreateComment creates a new comment
	CreateComment(ctx context.Context, userID, feedbackID string, content string, parentCommentID *string) (*model.FeedbackComment, error)

	// AddReaction adds a reaction to a feedback item
	AddReaction(ctx context.Context, userID, feedbackID string, reactionType string) error

	// RemoveReaction removes a reaction from a feedback item
	RemoveReaction(ctx context.Context, userID, feedbackID string, reactionType string) error

	// GetReactionsCount gets reaction counts for a feedback item
	GetReactionsCount(ctx context.Context, feedbackID string) (map[string]int, error)

	// GetCommentsCount gets comment count for a feedback item
	GetCommentsCount(ctx context.Context, feedbackID string) (int, error)
}

