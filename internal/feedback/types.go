package feedback

import (
	"context"
	"time"

	"ethos/internal/feedback/model"
)

// TemplateSuggestionRequest represents a request to suggest a new template
type TemplateSuggestionRequest struct {
	SuggestedBy   *string             `json:"suggested_by,omitempty"`
	UsageContext  string              `json:"usage_context" binding:"required"`
	Details       string              `json:"details" binding:"required"`
	DesiredFields []map[string]string `json:"desired_fields,omitempty"`
}

// BatchFeedbackItem represents a single feedback item in a batch request
type BatchFeedbackItem struct {
	Content     string  `json:"content" binding:"required"`
	Type        *string `json:"type,omitempty"`
	Visibility  *string `json:"visibility,omitempty"`
	IsAnonymous bool    `json:"is_anonymous,omitempty"`
}

// BatchFeedbackRequest represents a request to create multiple feedback items
type BatchFeedbackRequest struct {
	Items []BatchFeedbackItem `json:"items" binding:"required,dive"`
}

// BatchFeedbackResult represents the result of creating a single feedback item
type BatchFeedbackResult struct {
	FeedbackID string `json:"feedback_id"`
	Status     string `json:"status"`
}

// BatchFeedbackResponse represents the response from a batch feedback creation
type BatchFeedbackResponse struct {
	Submitted []BatchFeedbackResult `json:"submitted"`
}

// FeedFilters represents filtering options for the feedback feed
type FeedFilters struct {
	ReviewerType *string  `json:"reviewer_type,omitempty"` // "public" or "org"
	Context      *string  `json:"context,omitempty"`       // e.g., "project", "team", "initiative"
	Verification *string  `json:"verification,omitempty"`  // "verified", "unverified"
	Tags         []string `json:"tags,omitempty"`          // Comma-separated tags
}

// ExportResponse represents the response from a feedback export
type ExportResponse struct {
	Format      string `json:"format"`
	ContentType string `json:"content_type"`
	Data        string `json:"data"`
	Count       int    `json:"count"`
	Filename    string `json:"filename,omitempty"`
}

// Repository defines the interface for feedback data access
type Repository interface {
	// Core CRUD operations
	GetFeed(ctx context.Context, limit, offset int) ([]*model.FeedbackItem, int, error)
	GetFeedbackByID(ctx context.Context, feedbackID string) (*model.FeedbackItem, error)
	GetComments(ctx context.Context, feedbackID string, limit, offset int) ([]*model.FeedbackComment, int, error)
	CreateFeedback(ctx context.Context, userID string, item *model.FeedbackItem) (*model.FeedbackItem, error)
	CreateComment(ctx context.Context, userID, feedbackID string, comment *model.FeedbackComment) (*model.FeedbackComment, error)
	AddReaction(ctx context.Context, userID, feedbackID, reactionType string) error
	RemoveReaction(ctx context.Context, userID, feedbackID, reactionType string) error

	// Template operations
	GetTemplates(ctx context.Context, contextFilter, tagsFilter string) ([]*model.FeedbackTemplate, error)
	SubmitTemplateSuggestion(ctx context.Context, req *TemplateSuggestionRequest) error

	// Analytics and impact
	GetImpact(ctx context.Context, userID *string, from, to *time.Time) (*model.FeedbackImpact, error)
	GetFeedbackAnalytics(ctx context.Context, userID *string, from, to *time.Time) (*model.FeedbackAnalytics, error)

	// Batch operations
	CreateBatchFeedback(ctx context.Context, userID string, req *BatchFeedbackRequest) (*BatchFeedbackResponse, error)

	// Enhanced filtering
	GetFeedWithFilters(ctx context.Context, limit, offset int, filters *FeedFilters) ([]*model.FeedbackItem, int, error)

	// Bookmarks
	GetBookmarks(ctx context.Context, userID string, limit, offset int) ([]*model.FeedbackItem, int, error)
	AddBookmark(ctx context.Context, userID, feedbackID string) error
	RemoveBookmark(ctx context.Context, userID, feedbackID string) error

	// Export
	ExportFeedback(ctx context.Context, filters *FeedFilters, format string) (*ExportResponse, error)

	// Updates and deletes
	UpdateFeedback(ctx context.Context, userID, feedbackID string, content *string, feedbackType *model.FeedbackType, visibility *model.FeedbackVisibility) (*model.FeedbackItem, error)
	DeleteFeedback(ctx context.Context, userID, feedbackID string) error
	UpdateComment(ctx context.Context, userID, feedbackID, commentID string, content string) (*model.FeedbackComment, error)
	DeleteComment(ctx context.Context, userID, feedbackID, commentID string) error

	// Search and trending
	SearchFeedback(ctx context.Context, query string, limit, offset int) ([]*model.FeedbackItem, int, error)
	GetTrendingFeedback(ctx context.Context, limit, offset int) ([]*model.FeedbackItem, int, error)

	// Pinning
	PinFeedback(ctx context.Context, userID, feedbackID string) error
	UnpinFeedback(ctx context.Context, userID, feedbackID string) error

	// Stats
	GetFeedbackStats(ctx context.Context) (*model.FeedbackStats, error)
}
