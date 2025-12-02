package repository

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	fbModel "ethos/internal/feedback/model"
	"ethos/internal/database"
	"ethos/internal/dashboard/model"
)

// PostgresRepository implements the Repository interface using PostgreSQL
type PostgresRepository struct {
	db *database.DB
}

// NewPostgresRepository creates a new PostgreSQL repository
func NewPostgresRepository(db *database.DB) Repository {
	return &PostgresRepository{db: db}
}

// GetDashboard retrieves a dashboard snapshot
func (r *PostgresRepository) GetDashboard(ctx context.Context, userID string) (*model.DashboardSnapshot, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "repository.GetDashboard")
	defer span.End()

	snapshot := &model.DashboardSnapshot{
		Stats:            make(map[string]int),
		SuggestedActions: []string{},
		RecentFeedback:   []*fbModel.FeedbackItem{},
	}

	// Get recent feedback (limit to 5)
	feedbackQuery := `
		SELECT f.feedback_id, f.author_id, f.content, f.type, f.visibility, f.created_at,
		       u.id, u.name
		FROM feedback_items f
		JOIN users u ON f.author_id = u.id
		WHERE f.author_id = $1 OR f.visibility = 'public'
		ORDER BY f.created_at DESC
		LIMIT 5
	`

	rows, err := r.db.Pool.Query(ctx, feedbackQuery, userID)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			// Simplified - would need full parsing like in feedback repository
			// For now, just count
		}
	}

	// Get stats
	var feedbackGiven, comments int
	statsQuery := `
		SELECT 
			(SELECT COUNT(*) FROM feedback_items WHERE author_id = $1) as feedback_given,
			(SELECT COUNT(*) FROM feedback_comments WHERE author_id = $1) as comments
	`
	err = r.db.Pool.QueryRow(ctx, statsQuery, userID).Scan(&feedbackGiven, &comments)
	if err == nil {
		snapshot.Stats["feedback_given"] = feedbackGiven
		snapshot.Stats["comments"] = comments
	}

	// Suggested actions
	if feedbackGiven == 0 {
		snapshot.SuggestedActions = append(snapshot.SuggestedActions, "Give feedback")
	}
	if comments == 0 {
		snapshot.SuggestedActions = append(snapshot.SuggestedActions, "Comment on feedback")
	}

	span.SetStatus(codes.Ok, "")
	return snapshot, nil
}

