package model

import (
	"time"

	authModel "ethos/internal/auth/model"
)

// ModerationState represents the state of moderated content
type ModerationState string

const (
	ModerationStatePending   ModerationState = "pending"
	ModerationStateWarned    ModerationState = "warned"
	ModerationStateActioned  ModerationState = "actioned"
	ModerationStateEscalated ModerationState = "escalated"
	ModerationStateAppealed  ModerationState = "appealed"
)

// AppealStatus represents the status of an appeal request
type AppealStatus string

const (
	AppealStatusPending   AppealStatus = "pending"
	AppealStatusReviewing AppealStatus = "reviewing"
	AppealStatusApproved  AppealStatus = "approved"
	AppealStatusRejected  AppealStatus = "rejected"
	AppealStatusResolved  AppealStatus = "resolved"
)

// ModerationAppeal represents an appeal for a moderation decision
type ModerationAppeal struct {
	AppealID        string                 `json:"appeal_id"`
	ModeratedItemID string                 `json:"moderated_item_id"`
	ItemType        string                 `json:"item_type"` // "feedback", "comment", "profile"
	Reason          string                 `json:"reason"`
	Details         string                 `json:"details,omitempty"`
	Status          AppealStatus           `json:"status"`
	SubmittedBy     *authModel.UserSummary `json:"submitted_by"`
	SubmittedAt     time.Time              `json:"submitted_at"`
	ReviewedAt      *time.Time             `json:"reviewed_at,omitempty"`
	ReviewerNotes   string                 `json:"reviewer_notes,omitempty"`
}

// ModerationRule represents a moderation rule
type ModerationRule struct {
	RuleID      string `json:"rule_id"`
	Description string `json:"description"`
	Status      string `json:"status"` // "applied", "not_applied"
}

// ModerationContext represents the moderation context for an item
type ModerationContext struct {
	ItemID        string           `json:"item_id"`
	ItemType      string           `json:"item_type"`
	CurrentState  ModerationState  `json:"current_state"`
	RulesApplied  []ModerationRule `json:"rules_applied"`
	ReviewerNotes string           `json:"reviewer_notes,omitempty"`
}

// ModerationAction represents a moderation action taken on a user or content
type ModerationAction struct {
	ID             string
	OrganizationID string
	TargetID       string // User ID or Content ID
	ActionType     string // "warning", "suspension", "ban", "content_removal"
	Reason         string
	Details        string
	Duration       *int   // Duration in days (for suspension/temporary bans)
	IssuedBy       string // Moderator/Admin ID
	AppealsAllowed int    // Number of appeals allowed
	AppealsUsed    int    // Number of appeals used
	CreatedAt      time.Time
	ExpiresAt      *time.Time
}

// ModerationActionResponse represents a moderation action for API responses
type ModerationActionResponse struct {
	ID             string     `json:"id"`
	OrganizationID string     `json:"organization_id"`
	TargetID       string     `json:"target_id"`
	ActionType     string     `json:"action_type"`
	Reason         string     `json:"reason"`
	Details        string     `json:"details"`
	Duration       *int       `json:"duration,omitempty"`
	IssuedBy       string     `json:"issued_by"`
	ModeratorName  string     `json:"moderator_name"`
	AppealsAllowed int        `json:"appeals_allowed"`
	AppealsUsed    int        `json:"appeals_used"`
	CreatedAt      time.Time  `json:"created_at"`
	ExpiresAt      *time.Time `json:"expires_at,omitempty"`
}

// ModerationHistory represents historical moderation action record
type ModerationHistory struct {
	ID             string
	OrganizationID string
	UserID         string
	ActionType     string
	Description    string
	Reason         string
	PerformedBy    string
	CreatedAt      time.Time
}

// ModerationHistoryResponse represents moderation history for API responses
type ModerationHistoryResponse struct {
	ID             string    `json:"id"`
	OrganizationID string    `json:"organization_id"`
	UserID         string    `json:"user_id"`
	UserName       string    `json:"user_name"`
	ActionType     string    `json:"action_type"`
	Description    string    `json:"description"`
	Reason         string    `json:"reason"`
	PerformedBy    string    `json:"performed_by"`
	PerformerName  string    `json:"performer_name"`
	CreatedAt      time.Time `json:"created_at"`
}
// ORGANIZATION ADMIN MODERATION MODELS

// PendingContentItem represents a content item pending moderation
type PendingContentItem struct {
	ID          string `json:"id"`
	Type        string `json:"type"` // feedback, comment
	AuthorID    string `json:"author_id"`
	AuthorName  string `json:"author_name"`
	Content     string `json:"content"`
	SubmittedAt time.Time `json:"submitted_at"`
	Flags       []string `json:"flags"` // spam, harassment, inappropriate, etc.
	Priority    string `json:"priority"` // low, medium, high
}

// OrganizationModerationStats represents moderation statistics for an organization
type OrganizationModerationStats struct {
	TotalModerated     int64 `json:"total_moderated"`
	PendingContent     int64 `json:"pending_content"`
	ApprovedContent    int64 `json:"approved_content"`
	RejectedContent    int64 `json:"rejected_content"`
	EscalatedContent   int64 `json:"escalated_content"`
	AverageReviewTime  float64 `json:"average_review_time_hours"`
	ModeratorWorkload  int64 `json:"moderator_workload"`
}
