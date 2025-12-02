package model

import (
	"time"

	authModel "ethos/internal/auth/model"
)

// FeedbackVisibility represents visibility levels for feedback
type FeedbackVisibility string

const (
	FeedbackVisibilityPublic FeedbackVisibility = "public"
	FeedbackVisibilityPrivate FeedbackVisibility = "private"
	FeedbackVisibilityTeam   FeedbackVisibility = "team"
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
	FeedbackID   string                    `json:"feedback_id"`
	Author       *authModel.UserSummary        `json:"author"`
	Content      string                    `json:"content"`
	Type         *FeedbackType             `json:"type,omitempty"`
	Visibility   *FeedbackVisibility       `json:"visibility,omitempty"`
	Reactions    map[string]int            `json:"reactions"`
	Dimensions   []FeedbackDimensionScore  `json:"dimensions,omitempty"`
	CommentsCount int                      `json:"comments_count"`
	CreatedAt    time.Time                 `json:"created_at"`
}

// FeedbackDimensionScore represents dimension-level scoring
type FeedbackDimensionScore struct {
	Dimension string `json:"dimension"`
	Score     int    `json:"score"`
}

// FeedbackComment represents a comment on feedback
type FeedbackComment struct {
	CommentID       string             `json:"comment_id"`
	Author          *authModel.UserSummary `json:"author"`
	Content         string             `json:"content"`
	CreatedAt       time.Time          `json:"created_at"`
	ParentCommentID *string            `json:"parent_comment_id,omitempty"`
}

