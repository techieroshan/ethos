package repository

import (
	"context"

	"ethos/internal/appeal/model"
)

// Repository defines the interface for appeal data access
type Repository interface {
	// SubmitAppeal creates a new appeal
	SubmitAppeal(ctx context.Context, userID string, appealType model.AppealType, referenceID *string, description string) (*model.Appeal, error)

	// GetUserAppeals retrieves appeals for a specific user
	GetUserAppeals(ctx context.Context, userID string, limit, offset int) ([]*model.Appeal, int, error)

	// GetAppealByID retrieves a specific appeal by ID (only if user owns it)
	GetAppealByID(ctx context.Context, userID, appealID string) (*model.Appeal, error)

	// UpdateAppealStatus updates the status of an appeal (for admins/moderators)
	UpdateAppealStatus(ctx context.Context, appealID string, status model.AppealStatus, adminNotes *string) error
}
