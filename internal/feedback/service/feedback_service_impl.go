package service

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	feedbackPkg "ethos/internal/feedback"
	"ethos/internal/feedback/model"
	"ethos/internal/feedback/repository"
	"ethos/pkg/errors"
)

// FeedbackService implements the Service interface
type FeedbackService struct {
	client FeedbackClient        // Can be REST or gRPC client
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

// GetTemplates retrieves feedback templates with optional filtering
func (s *FeedbackService) GetTemplates(ctx context.Context, contextFilter, tagsFilter string) ([]*model.FeedbackTemplate, error) {
	templates, err := s.repo.GetTemplates(ctx, contextFilter, tagsFilter)
	if err != nil {
		return nil, err
	}

	return templates, nil
}

// SubmitTemplateSuggestion submits a template suggestion
func (s *FeedbackService) SubmitTemplateSuggestion(ctx context.Context, req *feedbackPkg.TemplateSuggestionRequest) error {
	// In a real implementation, this might save to database or send to a queue
	// For now, we'll just log it (the API spec shows it just returns success)
	return s.repo.SubmitTemplateSuggestion(ctx, req)
}

// GetImpact retrieves aggregated feedback analytics
func (s *FeedbackService) GetImpact(ctx context.Context, userID *string, from, to *time.Time) (*model.FeedbackImpact, error) {
	return s.repo.GetImpact(ctx, userID, from, to)
}

// CreateBatchFeedback creates multiple feedback items in a batch
func (s *FeedbackService) CreateBatchFeedback(ctx context.Context, userID string, req *feedbackPkg.BatchFeedbackRequest) (*feedbackPkg.BatchFeedbackResponse, error) {
	return s.repo.CreateBatchFeedback(ctx, userID, req)
}

// GetFeedWithFilters retrieves a paginated feed of feedback items with enhanced filtering
func (s *FeedbackService) GetFeedWithFilters(ctx context.Context, limit, offset int, filters *feedbackPkg.FeedFilters) ([]*model.FeedbackItem, int, error) {
	return s.repo.GetFeedWithFilters(ctx, limit, offset, filters)
}

// GetBookmarks retrieves bookmarked feedback items for a user
func (s *FeedbackService) GetBookmarks(ctx context.Context, userID string, limit, offset int) ([]*model.FeedbackItem, int, error) {
	return s.repo.GetBookmarks(ctx, userID, limit, offset)
}

// AddBookmark adds a bookmark for a feedback item
func (s *FeedbackService) AddBookmark(ctx context.Context, userID, feedbackID string) error {
	return s.repo.AddBookmark(ctx, userID, feedbackID)
}

// RemoveBookmark removes a bookmark for a feedback item
func (s *FeedbackService) RemoveBookmark(ctx context.Context, userID, feedbackID string) error {
	return s.repo.RemoveBookmark(ctx, userID, feedbackID)
}

// ExportFeedback exports feedback data with optional filtering
func (s *FeedbackService) ExportFeedback(ctx context.Context, filters *feedbackPkg.FeedFilters, format string) (*feedbackPkg.ExportResponse, error) {
	// Validate format
	if format != "json" && format != "csv" {
		return nil, errors.ErrValidationFailed
	}

	// Get all matching feedback items (no pagination for export)
	items, _, err := s.repo.GetFeedWithFilters(ctx, 10000, 0, filters) // Reasonable limit for export
	if err != nil {
		return nil, err
	}

	var data string
	var contentType string
	var filename string

	switch format {
	case "json":
		data, err = s.formatJSONExport(items)
		contentType = "application/json"
		filename = "feedback_export.json"
	case "csv":
		data, err = s.formatCSVExport(items)
		contentType = "text/csv"
		filename = "feedback_export.csv"
	}

	if err != nil {
		return nil, err
	}

	return &feedbackPkg.ExportResponse{
		Format:      format,
		ContentType: contentType,
		Data:        data,
		Count:       len(items),
		Filename:    filename,
	}, nil
}

// formatJSONExport formats feedback items as JSON
func (s *FeedbackService) formatJSONExport(items []*model.FeedbackItem) (string, error) {
	type exportItem struct {
		FeedbackID    string         `json:"feedback_id"`
		Content       string         `json:"content"`
		AuthorName    string         `json:"author_name"`
		AuthorRole    string         `json:"author_role,omitempty"`
		Type          *string        `json:"type,omitempty"`
		Visibility    *string        `json:"visibility,omitempty"`
		IsAnonymous   bool           `json:"is_anonymous"`
		Helpfulness   float64        `json:"helpfulness"`
		Reactions     map[string]int `json:"reactions"`
		CommentsCount int            `json:"comments_count"`
		CreatedAt     time.Time      `json:"created_at"`
	}

	var exportItems []exportItem
	for _, item := range items {
		exportItem := exportItem{
			FeedbackID:    item.FeedbackID,
			Content:       item.Content,
			IsAnonymous:   item.IsAnonymous,
			Helpfulness:   item.Helpfulness,
			Reactions:     item.Reactions,
			CommentsCount: item.CommentsCount,
			CreatedAt:     item.CreatedAt,
		}

		if item.Author != nil {
			exportItem.AuthorName = item.Author.Name
			// Note: Role field doesn't exist on UserSummary, commented out
			// if item.Author.Role != "" {
			// 	exportItem.AuthorRole = item.Author.Role
			// }
		}

		if item.Type != nil {
			typeStr := string(*item.Type)
			exportItem.Type = &typeStr
		}

		if item.Visibility != nil {
			visStr := string(*item.Visibility)
			exportItem.Visibility = &visStr
		}

		exportItems = append(exportItems, exportItem)
	}

	data, err := json.Marshal(exportItems)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// formatCSVExport formats feedback items as CSV
func (s *FeedbackService) formatCSVExport(items []*model.FeedbackItem) (string, error) {
	var buf strings.Builder
	writer := csv.NewWriter(&buf)

	// Write header
	header := []string{
		"feedback_id",
		"content",
		"author_name",
		"author_role",
		"type",
		"visibility",
		"is_anonymous",
		"helpfulness",
		"reactions_like",
		"reactions_helpful",
		"comments_count",
		"created_at",
	}
	if err := writer.Write(header); err != nil {
		return "", err
	}

	// Write data rows
	for _, item := range items {
		row := []string{
			item.FeedbackID,
			item.Content,
			item.Author.Name,
			"", // AuthorRole - field doesn't exist on UserSummary
			stringPtrToString(item.Type),
			visibilityToString(item.Visibility),
			boolToString(item.IsAnonymous),
			floatToString(item.Helpfulness),
			intToString(item.Reactions["like"]),
			intToString(item.Reactions["helpful"]),
			intToString(item.CommentsCount),
			item.CreatedAt.Format(time.RFC3339),
		}
		if err := writer.Write(row); err != nil {
			return "", err
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// UpdateFeedback updates an existing feedback item
func (s *FeedbackService) UpdateFeedback(ctx context.Context, userID, feedbackID string, req *UpdateFeedbackRequest) (*model.FeedbackItem, error) {
	// Retrieve existing feedback first to ensure it exists and user owns it
	item, err := s.repo.GetFeedbackByID(ctx, feedbackID)
	if err != nil {
		return nil, err
	}

	// Check if user owns the feedback
	if item.Author == nil || item.Author.ID != userID {
		return nil, errors.ErrForbidden
	}

	// Update only provided fields
	if req.Content != nil {
		item.Content = *req.Content
	}
	if req.Type != nil {
		item.Type = req.Type
	}
	if req.Visibility != nil {
		item.Visibility = req.Visibility
	}

	// Persist the update
	err = s.repo.UpdateFeedback(ctx, feedbackID, item)
	if err != nil {
		return nil, err
	}

	return item, nil
}

// DeleteFeedback deletes a feedback item
func (s *FeedbackService) DeleteFeedback(ctx context.Context, userID, feedbackID string) error {
	// Retrieve feedback to check ownership
	item, err := s.repo.GetFeedbackByID(ctx, feedbackID)
	if err != nil {
		return err
	}

	// Check if user owns the feedback
	if item.Author == nil || item.Author.ID != userID {
		return errors.ErrForbidden
	}

	return s.repo.DeleteFeedback(ctx, feedbackID)
}

// UpdateComment updates an existing comment
func (s *FeedbackService) UpdateComment(ctx context.Context, userID, feedbackID, commentID string, req *UpdateCommentRequest) (*model.FeedbackComment, error) {
	// Retrieve comment to check ownership
	comment, err := s.repo.GetComment(ctx, feedbackID, commentID)
	if err != nil {
		return nil, err
	}

	// Check if user owns the comment
	if comment.Author == nil || comment.Author.ID != userID {
		return nil, errors.ErrForbidden
	}

	// Update content
	comment.Content = req.Content

	// Persist the update
	err = s.repo.UpdateComment(ctx, feedbackID, commentID, comment)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

// DeleteComment deletes a comment
func (s *FeedbackService) DeleteComment(ctx context.Context, userID, feedbackID, commentID string) error {
	// Retrieve comment to check ownership
	comment, err := s.repo.GetComment(ctx, feedbackID, commentID)
	if err != nil {
		return err
	}

	// Check if user owns the comment
	if comment.Author == nil || comment.Author.ID != userID {
		return errors.ErrForbidden
	}

	return s.repo.DeleteComment(ctx, feedbackID, commentID)
}

// GetFeedbackAnalytics retrieves detailed feedback analytics
func (s *FeedbackService) GetFeedbackAnalytics(ctx context.Context, userID *string, from, to *time.Time) (*model.FeedbackAnalytics, error) {
	analytics, err := s.repo.GetFeedbackAnalytics(ctx, userID, from, to)
	if err != nil {
		return nil, err
	}

	return analytics, nil
}

// Helper functions for CSV formatting
func stringPtrToString(s *model.FeedbackType) string {
	if s == nil {
		return ""
	}
	return string(*s)
}

func visibilityToString(v *model.FeedbackVisibility) string {
	if v == nil {
		return ""
	}
	return string(*v)
}

func boolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

func floatToString(f float64) string {
	return strconv.FormatFloat(f, 'f', 2, 64)
}

func intToString(i int) string {
	return strconv.Itoa(i)
}
