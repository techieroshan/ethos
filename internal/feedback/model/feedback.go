package model

import (
	"time"

	authModel "ethos/internal/auth/model"
)

// FeedbackVisibility represents visibility levels for feedback
type FeedbackVisibility string

const (
	FeedbackVisibilityPublic  FeedbackVisibility = "public"
	FeedbackVisibilityPrivate FeedbackVisibility = "private"
	FeedbackVisibilityTeam    FeedbackVisibility = "team"
)

// FeedbackType represents types of feedback
type FeedbackType string

const (
	FeedbackTypeAppreciation FeedbackType = "appreciation"
	FeedbackTypeSuggestion   FeedbackType = "suggestion"
	FeedbackTypeIssue        FeedbackType = "issue"
	FeedbackTypeOther        FeedbackType = "other"
)

// FeedbackItem represents a feedback post
type FeedbackItem struct {
	FeedbackID         string                     `json:"feedback_id"`
	Author             *authModel.UserSummary     `json:"author"`
	Content            string                     `json:"content"`
	Type               *FeedbackType              `json:"type,omitempty"`
	Visibility         *FeedbackVisibility        `json:"visibility,omitempty"`
	Reactions          map[string]int             `json:"reactions"`
	ReactionsAnalytics *FeedbackReactionAnalytics `json:"reactions_analytics,omitempty"`
	FollowUps          []FeedbackFollowUp         `json:"follow_ups,omitempty"`
	IsAnonymous        bool                       `json:"is_anonymous,omitempty"`
	Helpfulness        float64                    `json:"helpfulness,omitempty"`
	Dimensions         []FeedbackDimensionScore   `json:"dimensions,omitempty"`
	CommentsCount      int                        `json:"comments_count"`
	CreatedAt          time.Time                  `json:"created_at"`
}

// FeedbackDimensionScore represents dimension-level scoring
type FeedbackDimensionScore struct {
	Dimension string `json:"dimension"`
	Score     int    `json:"score"`
}

// FeedbackComment represents a comment on feedback
type FeedbackComment struct {
	CommentID       string                 `json:"comment_id"`
	Author          *authModel.UserSummary `json:"author"`
	Content         string                 `json:"content"`
	CreatedAt       time.Time              `json:"created_at"`
	ParentCommentID *string                `json:"parent_comment_id,omitempty"`
}

// FeedbackReactionAnalytics represents detailed reaction analytics
type FeedbackReactionAnalytics struct {
	Reactions map[string]ReactionDetail `json:"reactions,omitempty"`
}

// ReactionDetail represents detailed information about a specific reaction
type ReactionDetail struct {
	Count   int      `json:"count"`
	UserIDs []string `json:"user_ids"`
}

// FeedbackFollowUp represents a follow-up discussion on feedback
type FeedbackFollowUp struct {
	FollowUpID string                 `json:"follow_up_id"`
	Content    string                 `json:"content"`
	Author     *authModel.UserSummary `json:"author"`
	CreatedAt  time.Time              `json:"created_at"`
}

// FeedbackTemplate represents a feedback template
type FeedbackTemplate struct {
	TemplateID     string                 `json:"template_id"`
	Name           string                 `json:"name"`
	Description    string                 `json:"description"`
	ContextTags    []string               `json:"context_tags"`
	TemplateFields map[string]interface{} `json:"template_fields"`
}

// FeedbackImpact represents aggregated feedback analytics
type FeedbackImpact struct {
	FeedbackCount      int             `json:"feedback_count"`
	AverageHelpfulness float64         `json:"average_helpfulness"`
	ReactionTotals     map[string]int  `json:"reaction_totals"`
	FollowUpCount      int             `json:"follow_up_count"`
	Trends             []FeedbackTrend `json:"trends"`
}

// FeedbackTrend represents feedback analytics over time
type FeedbackTrend struct {
	Date              time.Time `json:"date"`
	Helpfulness       float64   `json:"helpfulness"`
	FeedbackSubmitted int       `json:"feedback_submitted"`
}

// FeedbackAnalytics represents detailed feedback analytics
type FeedbackAnalytics struct {
	TotalCount             int                 `json:"total_count"`
	AverageHelpfulness     float64             `json:"average_helpfulness"`
	TopReactions           map[string]int      `json:"top_reactions"`
	TrendByDay             []FeedbackTrendData `json:"trend_by_day"`
	TypeDistribution       map[string]int      `json:"type_distribution"`
	VisibilityDistribution map[string]int      `json:"visibility_distribution"`
}

// FeedbackTrendData represents trend data for a specific day
type FeedbackTrendData struct {
	Date               string  `json:"date"`
	Count              int     `json:"count"`
	AverageHelpfulness float64 `json:"average_helpfulness"`
}

// PinnedFeedback represents a pinned feedback item
type PinnedFeedback struct {
	FeedbackID string    `json:"feedback_id"`
	UserID     string    `json:"user_id"`
	PinnedAt   time.Time `json:"pinned_at"`
}

// FeedbackStats represents statistics for feedback
type FeedbackStats struct {
	TotalFeedback     int     `json:"total_feedback"`
	TotalComments     int     `json:"total_comments"`
	TotalReactions    int     `json:"total_reactions"`
	AverageHelpfulness float64 `json:"average_helpfulness"`
	MostPopularType    string  `json:"most_popular_type"`
	MostCommonReaction string  `json:"most_common_reaction"`
}

// SearchFeedbackRequest represents a request to search feedback
type SearchFeedbackRequest struct {
	Query string `json:"query" binding:"required"`
}

// TrendingFeedbackResponse represents the response for trending feedback
type TrendingFeedbackResponse struct {
	Items []*FeedbackItem `json:"items"`
}
