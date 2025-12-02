package service

import (
	"context"

	"ethos/internal/moderation/model"
)

// SubmitAppealRequest represents a request to submit a moderation appeal
type SubmitAppealRequest struct {
	ModeratedItemID string `json:"moderated_item_id" binding:"required"`
	ItemType        string `json:"item_type" binding:"required,oneof=feedback comment profile"`
	Reason          string `json:"reason" binding:"required"`
	Details         string `json:"details,omitempty"`
}

// Service defines the interface for moderation business logic
type Service interface {
	// Appeal-related methods
	// ListAppeals retrieves all appeals for an organization
	ListAppeals(ctx context.Context, orgID string, limit, offset int) ([]*model.ModerationAppeal, error)

	// SubmitAppeal submits a moderation appeal
	SubmitAppeal(ctx context.Context, userID, orgID string, req *SubmitAppealRequest) (*model.ModerationAppeal, error)

	// GetAppealContext retrieves context for an appeal
	GetAppealContext(ctx context.Context, appealID string) (*model.ModerationContext, error)

	// Action-related methods
	// ListModerationActions retrieves moderation actions for an organization
	ListModerationActions(ctx context.Context, orgID string, limit, offset int) ([]*model.ModerationActionResponse, error)

	// History-related methods
	// GetModerationHistory retrieves moderation history for a user
	GetModerationHistory(ctx context.Context, orgID, userID string, limit, offset int) ([]*model.ModerationHistoryResponse, error)

	// Context-related methods
	// GetModerationContext retrieves moderation context for an item
	GetModerationContext(ctx context.Context, itemID, itemType string) (*model.ModerationContext, error)
}
