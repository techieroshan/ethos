package service

import (
	"context"

	authModel "ethos/internal/auth/model"
	"ethos/internal/moderation/model"
	"ethos/internal/moderation/repository"
)

// ModerationService implements the Service interface
type ModerationService struct {
	repo repository.Repository
}

// NewModerationService creates a new moderation service
func NewModerationService(repo repository.Repository) Service {
	return &ModerationService{
		repo: repo,
	}
}

// ListAppeals retrieves all appeals for an organization
func (s *ModerationService) ListAppeals(ctx context.Context, orgID string, limit, offset int) ([]*model.ModerationAppeal, error) {
	return s.repo.ListAppeals(ctx, orgID, limit, offset)
}

// SubmitAppeal submits a moderation appeal
func (s *ModerationService) SubmitAppeal(ctx context.Context, userID, orgID string, req *SubmitAppealRequest) (*model.ModerationAppeal, error) {
	appeal := &model.ModerationAppeal{
		AppealID:        "",
		ModeratedItemID: req.ModeratedItemID,
		ItemType:        req.ItemType,
		Reason:          req.Reason,
		Details:         req.Details,
		Status:          model.AppealStatusPending,
		SubmittedBy: &authModel.UserSummary{
			ID: userID,
		},
	}
	return appeal, nil
}

// GetAppealContext retrieves context for an appeal
func (s *ModerationService) GetAppealContext(ctx context.Context, appealID string) (*model.ModerationContext, error) {
	return s.repo.GetAppealContext(ctx, appealID)
}

// ListModerationActions retrieves moderation actions for an organization
func (s *ModerationService) ListModerationActions(ctx context.Context, orgID string, limit, offset int) ([]*model.ModerationActionResponse, error) {
	actions, err := s.repo.ListModerationActions(ctx, orgID, limit, offset)
	if err != nil {
		return nil, err
	}

	responses := make([]*model.ModerationActionResponse, len(actions))
	for i, action := range actions {
		responses[i] = &model.ModerationActionResponse{
			ID:             action.ID,
			OrganizationID: action.OrganizationID,
			TargetID:       action.TargetID,
			ActionType:     action.ActionType,
			Reason:         action.Reason,
			Details:        action.Details,
			Duration:       action.Duration,
			IssuedBy:       action.IssuedBy,
			ModeratorName:  "Moderator",
			AppealsAllowed: action.AppealsAllowed,
			AppealsUsed:    action.AppealsUsed,
			CreatedAt:      action.CreatedAt,
			ExpiresAt:      action.ExpiresAt,
		}
	}

	return responses, nil
}

// GetModerationHistory retrieves moderation history for a user
func (s *ModerationService) GetModerationHistory(ctx context.Context, orgID, userID string, limit, offset int) ([]*model.ModerationHistoryResponse, error) {
	history, err := s.repo.ListModerationHistory(ctx, orgID, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	responses := make([]*model.ModerationHistoryResponse, len(history))
	for i, h := range history {
		responses[i] = &model.ModerationHistoryResponse{
			ID:             h.ID,
			OrganizationID: h.OrganizationID,
			UserID:         h.UserID,
			UserName:       "User",
			ActionType:     h.ActionType,
			Description:    h.Description,
			Reason:         h.Reason,
			PerformedBy:    h.PerformedBy,
			PerformerName:  "Performer",
			CreatedAt:      h.CreatedAt,
		}
	}

	return responses, nil
}

// GetModerationContext retrieves moderation context for an item
func (s *ModerationService) GetModerationContext(ctx context.Context, itemID, itemType string) (*model.ModerationContext, error) {
	return s.repo.GetModerationContext(ctx, itemID, itemType)
}
