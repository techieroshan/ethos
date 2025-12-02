package converter

import (
	"time"

	feedbackpb "ethos/api/proto/feedback"
	fbModel "ethos/internal/feedback/model"
	authModel "ethos/internal/auth/model"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ProtoToFeedbackVisibility converts proto FeedbackVisibility to domain model
func ProtoToFeedbackVisibility(pb feedbackpb.FeedbackVisibility) *fbModel.FeedbackVisibility {
	switch pb {
	case feedbackpb.FeedbackVisibility_FEEDBACK_VISIBILITY_PUBLIC:
		vis := fbModel.FeedbackVisibilityPublic
		return &vis
	case feedbackpb.FeedbackVisibility_FEEDBACK_VISIBILITY_PRIVATE:
		vis := fbModel.FeedbackVisibilityPrivate
		return &vis
	case feedbackpb.FeedbackVisibility_FEEDBACK_VISIBILITY_TEAM:
		vis := fbModel.FeedbackVisibilityTeam
		return &vis
	default:
		return nil
	}
}

// ProtoToFeedbackType converts proto FeedbackType to domain model
func ProtoToFeedbackType(pb feedbackpb.FeedbackType) *fbModel.FeedbackType {
	switch pb {
	case feedbackpb.FeedbackType_FEEDBACK_TYPE_APPRECIATION:
		typ := fbModel.FeedbackTypeAppreciation
		return &typ
	case feedbackpb.FeedbackType_FEEDBACK_TYPE_SUGGESTION:
		typ := fbModel.FeedbackTypeSuggestion
		return &typ
	case feedbackpb.FeedbackType_FEEDBACK_TYPE_ISSUE:
		typ := fbModel.FeedbackTypeIssue
		return &typ
	case feedbackpb.FeedbackType_FEEDBACK_TYPE_OTHER:
		typ := fbModel.FeedbackTypeOther
		return &typ
	default:
		return nil
	}
}

// ProtoToUserSummary converts proto UserSummary to domain model
func ProtoToUserSummary(pb *feedbackpb.UserSummary) *authModel.UserSummary {
	if pb == nil {
		return nil
	}
	return &authModel.UserSummary{
		ID:   pb.Id,
		Name: pb.Name,
	}
}

// ProtoToFeedbackItem converts proto FeedbackItem to domain model
func ProtoToFeedbackItem(pb *feedbackpb.FeedbackItem) *fbModel.FeedbackItem {
	if pb == nil {
		return nil
	}

	item := &fbModel.FeedbackItem{
		FeedbackID:   pb.FeedbackId,
		Author:       ProtoToUserSummary(pb.Author),
		Content:      pb.Content,
		Reactions:    make(map[string]int),
		CommentsCount: int(pb.CommentsCount),
	}

	// Convert type
	if pb.Type != feedbackpb.FeedbackType_FEEDBACK_TYPE_UNSPECIFIED {
		item.Type = ProtoToFeedbackType(pb.Type)
	}

	// Convert visibility
	if pb.Visibility != feedbackpb.FeedbackVisibility_FEEDBACK_VISIBILITY_UNSPECIFIED {
		item.Visibility = ProtoToFeedbackVisibility(pb.Visibility)
	}

	// Convert reactions
	if pb.Reactions != nil {
		for k, v := range pb.Reactions.Reactions {
			item.Reactions[k] = int(v)
		}
	}

	// Convert dimensions
	if len(pb.Dimensions) > 0 {
		item.Dimensions = make([]fbModel.FeedbackDimensionScore, len(pb.Dimensions))
		for i, dim := range pb.Dimensions {
			item.Dimensions[i] = fbModel.FeedbackDimensionScore{
				Dimension: dim.Dimension,
				Score:     int(dim.Score),
			}
		}
	}

	// Convert timestamp
	if pb.CreatedAt != nil {
		item.CreatedAt = pb.CreatedAt.AsTime()
	} else {
		item.CreatedAt = time.Now()
	}

	return item
}

// ProtoToFeedbackComment converts proto FeedbackComment to domain model
func ProtoToFeedbackComment(pb *feedbackpb.FeedbackComment) *fbModel.FeedbackComment {
	if pb == nil {
		return nil
	}

	comment := &fbModel.FeedbackComment{
		CommentID: pb.CommentId,
		Author:    ProtoToUserSummary(pb.Author),
		Content:   pb.Content,
	}

	// Convert parent comment ID
	if pb.ParentCommentId != "" {
		comment.ParentCommentID = &pb.ParentCommentId
	}

	// Convert timestamp
	if pb.CreatedAt != nil {
		comment.CreatedAt = pb.CreatedAt.AsTime()
	} else {
		comment.CreatedAt = time.Now()
	}

	return comment
}

// FeedbackTypeToProto converts domain FeedbackType to proto
func FeedbackTypeToProto(typ *fbModel.FeedbackType) feedbackpb.FeedbackType {
	if typ == nil {
		return feedbackpb.FeedbackType_FEEDBACK_TYPE_UNSPECIFIED
	}
	switch *typ {
	case fbModel.FeedbackTypeAppreciation:
		return feedbackpb.FeedbackType_FEEDBACK_TYPE_APPRECIATION
	case fbModel.FeedbackTypeSuggestion:
		return feedbackpb.FeedbackType_FEEDBACK_TYPE_SUGGESTION
	case fbModel.FeedbackTypeIssue:
		return feedbackpb.FeedbackType_FEEDBACK_TYPE_ISSUE
	case fbModel.FeedbackTypeOther:
		return feedbackpb.FeedbackType_FEEDBACK_TYPE_OTHER
	default:
		return feedbackpb.FeedbackType_FEEDBACK_TYPE_UNSPECIFIED
	}
}

// FeedbackVisibilityToProto converts domain FeedbackVisibility to proto
func FeedbackVisibilityToProto(vis *fbModel.FeedbackVisibility) feedbackpb.FeedbackVisibility {
	if vis == nil {
		return feedbackpb.FeedbackVisibility_FEEDBACK_VISIBILITY_UNSPECIFIED
	}
	switch *vis {
	case fbModel.FeedbackVisibilityPublic:
		return feedbackpb.FeedbackVisibility_FEEDBACK_VISIBILITY_PUBLIC
	case fbModel.FeedbackVisibilityPrivate:
		return feedbackpb.FeedbackVisibility_FEEDBACK_VISIBILITY_PRIVATE
	case fbModel.FeedbackVisibilityTeam:
		return feedbackpb.FeedbackVisibility_FEEDBACK_VISIBILITY_TEAM
	default:
		return feedbackpb.FeedbackVisibility_FEEDBACK_VISIBILITY_UNSPECIFIED
	}
}

// UserSummaryToProto converts domain UserSummary to proto
func UserSummaryToProto(summary *authModel.UserSummary) *feedbackpb.UserSummary {
	if summary == nil {
		return nil
	}
	return &feedbackpb.UserSummary{
		Id:   summary.ID,
		Name: summary.Name,
	}
}

// FeedbackItemToProto converts domain FeedbackItem to proto
func FeedbackItemToProto(item *fbModel.FeedbackItem) *feedbackpb.FeedbackItem {
	if item == nil {
		return nil
	}

	pb := &feedbackpb.FeedbackItem{
		FeedbackId:   item.FeedbackID,
		Author:       UserSummaryToProto(item.Author),
		Content:      item.Content,
		CommentsCount: int32(item.CommentsCount),
		CreatedAt:    timestamppb.New(item.CreatedAt),
	}

	// Convert type
	if item.Type != nil {
		pb.Type = FeedbackTypeToProto(item.Type)
	}

	// Convert visibility
	if item.Visibility != nil {
		pb.Visibility = FeedbackVisibilityToProto(item.Visibility)
	}

	// Convert reactions
	if len(item.Reactions) > 0 {
		pb.Reactions = &feedbackpb.FeedbackReactionSummary{
			Reactions: make(map[string]int32),
		}
		for k, v := range item.Reactions {
			pb.Reactions.Reactions[k] = int32(v)
		}
	}

	// Convert dimensions
	if len(item.Dimensions) > 0 {
		pb.Dimensions = make([]*feedbackpb.FeedbackDimensionScore, len(item.Dimensions))
		for i, dim := range item.Dimensions {
			pb.Dimensions[i] = &feedbackpb.FeedbackDimensionScore{
				Dimension: dim.Dimension,
				Score:     int32(dim.Score),
			}
		}
	}

	return pb
}

