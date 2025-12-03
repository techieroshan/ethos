package service

import (
	"context"
	"time"

	"ethos/internal/feedback"
	"ethos/internal/feedback/model"
)

// CreateFeedbackRequest represents a request to create feedback
type CreateFeedbackRequest struct {
	Content    string                    `json:"content" binding:"required"`
	Type       *model.FeedbackType       `json:"type,omitempty"`
	Visibility *model.FeedbackVisibility `json:"visibility,omitempty"`
}

// CreateCommentRequest represents a request to create a comment
type CreateCommentRequest struct {
	Content         string  `json:"content" binding:"required"`
	ParentCommentID *string `json:"parent_comment_id,omitempty"`
}

// AddReactionRequest represents a request to add a reaction
type AddReactionRequest struct {
	ReactionType string `json:"reaction_type" binding:"required"`
}

// UpdateFeedbackRequest represents a request to update feedback
type UpdateFeedbackRequest struct {
	Content    *string                   `json:"content,omitempty"`
	Type       *model.FeedbackType       `json:"type,omitempty"`
	Visibility *model.FeedbackVisibility `json:"visibility,omitempty"`
}

// UpdateCommentRequest represents a request to update a comment
type UpdateCommentRequest struct {
	Content string `json:"content" binding:"required"`
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

	// GetTemplates retrieves feedback templates with optional filtering
	GetTemplates(ctx context.Context, contextFilter, tagsFilter string) ([]*model.FeedbackTemplate, error)

	// SubmitTemplateSuggestion submits a template suggestion
	SubmitTemplateSuggestion(ctx context.Context, req *feedback.TemplateSuggestionRequest) error

	// GetImpact retrieves aggregated feedback analytics
	GetImpact(ctx context.Context, userID *string, from, to *time.Time) (*model.FeedbackImpact, error)

	// CreateBatchFeedback creates multiple feedback items in a batch
	CreateBatchFeedback(ctx context.Context, userID string, req *feedback.BatchFeedbackRequest) (*feedback.BatchFeedbackResponse, error)

	// GetFeedWithFilters retrieves a paginated feed of feedback items with enhanced filtering
	GetFeedWithFilters(ctx context.Context, limit, offset int, filters *feedback.FeedFilters) ([]*model.FeedbackItem, int, error)

	// GetBookmarks retrieves bookmarked feedback items for a user
	GetBookmarks(ctx context.Context, userID string, limit, offset int) ([]*model.FeedbackItem, int, error)

	// AddBookmark adds a bookmark for a feedback item
	AddBookmark(ctx context.Context, userID, feedbackID string) error

	// RemoveBookmark removes a bookmark for a feedback item
	RemoveBookmark(ctx context.Context, userID, feedbackID string) error

	// ExportFeedback exports feedback data with optional filtering
	ExportFeedback(ctx context.Context, filters *feedback.FeedFilters, format string) (*feedback.ExportResponse, error)

	// UpdateFeedback updates an existing feedback item
	UpdateFeedback(ctx context.Context, userID, feedbackID string, req *UpdateFeedbackRequest) (*model.FeedbackItem, error)

	// DeleteFeedback deletes a feedback item
	DeleteFeedback(ctx context.Context, userID, feedbackID string) error

	// UpdateComment updates an existing comment
	UpdateComment(ctx context.Context, userID, feedbackID, commentID string, req *UpdateCommentRequest) (*model.FeedbackComment, error)

	// DeleteComment deletes a comment
	DeleteComment(ctx context.Context, userID, feedbackID, commentID string) error

	// GetFeedbackAnalytics retrieves detailed feedback analytics
	GetFeedbackAnalytics(ctx context.Context, userID *string, from, to *time.Time) (*model.FeedbackAnalytics, error)

	// SearchFeedback searches feedback items by content/metadata
	SearchFeedback(ctx context.Context, query string, limit, offset int) ([]*model.FeedbackItem, int, error)

	// GetTrendingFeedback retrieves trending feedback items
	GetTrendingFeedback(ctx context.Context, limit, offset int) ([]*model.FeedbackItem, int, error)

	// PinFeedback pins a feedback item
	PinFeedback(ctx context.Context, userID, feedbackID string) error

	// UnpinFeedback unpins a feedback item
	UnpinFeedback(ctx context.Context, userID, feedbackID string) error

	// GetFeedbackStats retrieves overall feedback statistics
	GetFeedbackStats(ctx context.Context) (*model.FeedbackStats, error)
}
