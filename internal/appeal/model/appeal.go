package model

import "time"

// AppealType represents types of appeals
type AppealType string

const (
	AppealTypeContentRemoval   AppealType = "content_removal"
	AppealTypeAccountSuspension AppealType = "account_suspension"
	AppealTypeFeedbackRemoval  AppealType = "feedback_removal"
	AppealTypeRatingDispute    AppealType = "rating_dispute"
	AppealTypeOther            AppealType = "other"
)

// AppealStatus represents the status of an appeal
type AppealStatus string

const (
	AppealStatusPending   AppealStatus = "pending"
	AppealStatusUnderReview AppealStatus = "under_review"
	AppealStatusApproved   AppealStatus = "approved"
	AppealStatusRejected   AppealStatus = "rejected"
	AppealStatusClosed     AppealStatus = "closed"
)

// Appeal represents an appeal submitted by a user
type Appeal struct {
	AppealID     string       `json:"appeal_id"`
	UserID       string       `json:"user_id"`
	Type         AppealType   `json:"type"`
	ReferenceID  *string      `json:"reference_id,omitempty"`
	Description  string       `json:"description"`
	Status       AppealStatus `json:"status"`
	AdminNotes   *string      `json:"admin_notes,omitempty"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
	ResolvedAt   *time.Time   `json:"resolved_at,omitempty"`
}
