package repository

import (
	"context"

	"ethos/internal/moderation/model"
)

// Repository defines the interface for moderation data access
type Repository interface {
	// Appeal-related methods
	// GetAppeal retrieves an appeal by ID
	GetAppeal(ctx context.Context, appealID string) (*model.ModerationAppeal, error)

	// ListAppeals retrieves all appeals for an organization
	ListAppeals(ctx context.Context, orgID string, limit, offset int) ([]*model.ModerationAppeal, error)

	// CreateAppeal creates a new appeal
	CreateAppeal(ctx context.Context, appeal *model.ModerationAppeal) error

	// UpdateAppeal updates an appeal
	UpdateAppeal(ctx context.Context, appeal *model.ModerationAppeal) error

	// GetAppealContext retrieves context for an appeal
	GetAppealContext(ctx context.Context, appealID string) (*model.ModerationContext, error)

	// Action-related methods
	// GetModerationAction retrieves a moderation action by ID
	GetModerationAction(ctx context.Context, actionID string) (*model.ModerationAction, error)

	// ListModerationActions retrieves moderation actions for an organization
	ListModerationActions(ctx context.Context, orgID string, limit, offset int) ([]*model.ModerationAction, error)

	// CreateModerationAction creates a new moderation action
	CreateModerationAction(ctx context.Context, action *model.ModerationAction) error

	// History-related methods
	// ListModerationHistory retrieves moderation history for a user
	ListModerationHistory(ctx context.Context, orgID, userID string, limit, offset int) ([]*model.ModerationHistory, error)

	// CreateModerationHistory records a moderation action in history
	CreateModerationHistory(ctx context.Context, history *model.ModerationHistory) error

	// Context-related methods
	// GetModerationContext retrieves moderation context for an item
	GetModerationContext(ctx context.Context, itemID, itemType string) (*model.ModerationContext, error)
}
