package service

import (
	"context"

	"ethos/internal/appeal/model"
)

// SubmitAppealRequest represents a request to submit an appeal
type SubmitAppealRequest struct {
	Type        model.AppealType `json:"type" binding:"required"`
	ReferenceID *string         `json:"reference_id,omitempty"`
	Description string          `json:"description" binding:"required,min=10,max=1000"`
}

// Service defines the interface for appeal business logic
type Service interface {
	// SubmitAppeal allows a user to submit an appeal
	SubmitAppeal(ctx context.Context, userID string, req *SubmitAppealRequest) (*model.Appeal, error)

	// GetUserAppeals retrieves appeals for a specific user
	GetUserAppeals(ctx context.Context, userID string, limit, offset int) ([]*model.Appeal, int, error)

	// GetAppealByID retrieves a specific appeal by ID (only if user owns it)
	GetAppealByID(ctx context.Context, userID, appealID string) (*model.Appeal, error)
}
