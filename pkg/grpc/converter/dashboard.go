package converter

import (
	dashboardpb "ethos/api/proto/dashboard"
	dashModel "ethos/internal/dashboard/model"
	fbModel "ethos/internal/feedback/model"
)

// ProtoToDashboardSnapshot converts proto DashboardSnapshot to domain model
func ProtoToDashboardSnapshot(pb *dashboardpb.DashboardSnapshot) *dashModel.DashboardSnapshot {
	if pb == nil {
		return nil
	}

	snapshot := &dashModel.DashboardSnapshot{
		RecentFeedback:   make([]*fbModel.FeedbackItem, 0, len(pb.RecentFeedback)),
		Stats:             make(map[string]int),
		SuggestedActions:  pb.SuggestedActions,
	}

	// Convert recent feedback
	for _, fb := range pb.RecentFeedback {
		// Use feedback converter for FeedbackItem
		item := ProtoToFeedbackItem(fb)
		if item != nil {
			snapshot.RecentFeedback = append(snapshot.RecentFeedback, item)
		}
	}

	// Convert stats
	if pb.Stats != nil {
		snapshot.Stats["feedback_given"] = int(pb.Stats.FeedbackGiven)
		snapshot.Stats["comments"] = int(pb.Stats.Comments)
		// Convert additional stats
		if pb.Stats.AdditionalStats != nil {
			for k, v := range pb.Stats.AdditionalStats {
				snapshot.Stats[k] = int(v)
			}
		}
	}

	return snapshot
}

