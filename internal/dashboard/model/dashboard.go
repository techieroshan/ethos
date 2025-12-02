package model

import (
	fbModel "ethos/internal/feedback/model"
)

// DashboardSnapshot represents a dashboard snapshot
type DashboardSnapshot struct {
	RecentFeedback   []*fbModel.FeedbackItem `json:"recent_feedback"`
	Stats            map[string]int           `json:"stats"`
	SuggestedActions []string                 `json:"suggested_actions"`
}

