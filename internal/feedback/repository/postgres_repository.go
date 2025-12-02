package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	authModel "ethos/internal/auth/model"
	"ethos/internal/database"
	"ethos/internal/feedback/model"
	"ethos/pkg/errors"
)

// PostgresRepository implements the Repository interface using PostgreSQL
type PostgresRepository struct {
	db *database.DB
}

// NewPostgresRepository creates a new PostgreSQL repository
func NewPostgresRepository(db *database.DB) Repository {
	return &PostgresRepository{db: db}
}

// GetFeed retrieves a paginated feed of feedback items
func (r *PostgresRepository) GetFeed(ctx context.Context, limit, offset int) ([]*model.FeedbackItem, int, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "repository.GetFeed")
	defer span.End()

	// Get total count
	var totalCount int
	countQuery := `SELECT COUNT(*) FROM feedback_items WHERE visibility = 'public'`
	err := r.db.Pool.QueryRow(ctx, countQuery).Scan(&totalCount)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, 0, errors.WrapError(err, "failed to count feedback items")
	}

	// Get feedback items
	query := `
		SELECT f.feedback_id, f.author_id, f.content, f.type, f.visibility, f.created_at,
		       u.id, u.name
		FROM feedback_items f
		JOIN users u ON f.author_id = u.id
		WHERE f.visibility = 'public'
		ORDER BY f.created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Pool.Query(ctx, query, limit, offset)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, 0, errors.WrapError(err, "failed to get feedback feed")
	}
	defer rows.Close()

	var items []*model.FeedbackItem
	for rows.Next() {
		item := &model.FeedbackItem{
			Reactions: make(map[string]int),
		}
		var authorID, authorName string
		var feedbackType, visibility *string

		err := rows.Scan(
			&item.FeedbackID,
			&authorID,
			&item.Content,
			&feedbackType,
			&visibility,
			&item.CreatedAt,
			&authorID,
			&authorName,
		)
		if err != nil {
			span.RecordError(err)
			continue
		}

		item.Author = &authModel.UserSummary{ID: authorID, Name: authorName}
		if feedbackType != nil {
			ft := model.FeedbackType(*feedbackType)
			item.Type = &ft
		}
		if visibility != nil {
			v := model.FeedbackVisibility(*visibility)
			item.Visibility = &v
		}

		// Get reactions count
		reactions, _ := r.GetReactionsCount(ctx, item.FeedbackID)
		item.Reactions = reactions

		// Get comments count
		commentsCount, _ := r.GetCommentsCount(ctx, item.FeedbackID)
		item.CommentsCount = commentsCount

		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, 0, errors.WrapError(err, "failed to iterate feedback items")
	}

	span.SetStatus(codes.Ok, "")
	return items, totalCount, nil
}

// GetFeedbackByID retrieves a feedback item by ID
func (r *PostgresRepository) GetFeedbackByID(ctx context.Context, feedbackID string) (*model.FeedbackItem, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "repository.GetFeedbackByID")
	defer span.End()

	query := `
		SELECT f.feedback_id, f.author_id, f.content, f.type, f.visibility, f.created_at,
		       u.id, u.name
		FROM feedback_items f
		JOIN users u ON f.author_id = u.id
		WHERE f.feedback_id = $1
	`

	item := &model.FeedbackItem{
		Reactions: make(map[string]int),
	}
	var authorID, authorName string
	var feedbackType, visibility *string
	var scannedAuthorID string

	err := r.db.Pool.QueryRow(ctx, query, feedbackID).Scan(
		&item.FeedbackID,
		&scannedAuthorID,
		&item.Content,
		&feedbackType,
		&visibility,
		&item.CreatedAt,
		&authorID,
		&authorName,
	)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		if err == pgx.ErrNoRows {
			return nil, errors.NewValidationError("feedback not found")
		}
		return nil, errors.WrapError(err, "failed to get feedback")
	}

	item.Author = &authModel.UserSummary{ID: authorID, Name: authorName}
	if feedbackType != nil {
		ft := model.FeedbackType(*feedbackType)
		item.Type = &ft
	}
	if visibility != nil {
		v := model.FeedbackVisibility(*visibility)
		item.Visibility = &v
	}

	// Get reactions count
	reactions, _ := r.GetReactionsCount(ctx, feedbackID)
	item.Reactions = reactions

	// Get comments count
	commentsCount, _ := r.GetCommentsCount(ctx, feedbackID)
	item.CommentsCount = commentsCount

	span.SetStatus(codes.Ok, "")
	return item, nil
}

// GetComments retrieves comments for a feedback item
func (r *PostgresRepository) GetComments(ctx context.Context, feedbackID string, limit, offset int) ([]*model.FeedbackComment, int, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "repository.GetComments")
	defer span.End()

	// Get total count
	var totalCount int
	countQuery := `SELECT COUNT(*) FROM feedback_comments WHERE feedback_id = $1`
	err := r.db.Pool.QueryRow(ctx, countQuery, feedbackID).Scan(&totalCount)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, 0, errors.WrapError(err, "failed to count comments")
	}

	// Get comments
	query := `
		SELECT c.comment_id, c.author_id, c.content, c.created_at, c.parent_comment_id,
		       u.id, u.name
		FROM feedback_comments c
		JOIN users u ON c.author_id = u.id
		WHERE c.feedback_id = $1
		ORDER BY c.created_at ASC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Pool.Query(ctx, query, feedbackID, limit, offset)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, 0, errors.WrapError(err, "failed to get comments")
	}
	defer rows.Close()

	var comments []*model.FeedbackComment
	for rows.Next() {
		comment := &model.FeedbackComment{}
		var authorID, authorName string
		var parentCommentID *string
		var scannedAuthorID string

		err := rows.Scan(
			&comment.CommentID,
			&scannedAuthorID,
			&comment.Content,
			&comment.CreatedAt,
			&parentCommentID,
			&authorID,
			&authorName,
		)
		if err != nil {
			span.RecordError(err)
			continue
		}

		comment.Author = &authModel.UserSummary{ID: authorID, Name: authorName}
		comment.ParentCommentID = parentCommentID
		comments = append(comments, comment)
	}

	if err = rows.Err(); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, 0, errors.WrapError(err, "failed to iterate comments")
	}

	span.SetStatus(codes.Ok, "")
	return comments, totalCount, nil
}

// CreateFeedback creates a new feedback item
func (r *PostgresRepository) CreateFeedback(ctx context.Context, userID string, content string, feedbackType *model.FeedbackType, visibility *model.FeedbackVisibility) (*model.FeedbackItem, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "repository.CreateFeedback")
	defer span.End()

	feedbackID := "f-" + uuid.New().String()
	now := time.Now()

	var typeStr, visibilityStr *string
	if feedbackType != nil {
		s := string(*feedbackType)
		typeStr = &s
	}
	if visibility != nil {
		s := string(*visibility)
		visibilityStr = &s
	} else {
		defaultVis := "public"
		visibilityStr = &defaultVis
	}

	query := `
		INSERT INTO feedback_items (feedback_id, author_id, content, type, visibility, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING feedback_id, author_id, content, type, visibility, created_at
	`

	item := &model.FeedbackItem{
		Reactions: make(map[string]int),
	}
	var authorID string

	err := r.db.Pool.QueryRow(ctx, query, feedbackID, userID, content, typeStr, visibilityStr, now, now).Scan(
		&item.FeedbackID,
		&authorID,
		&item.Content,
		&typeStr,
		&visibilityStr,
		&item.CreatedAt,
	)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, errors.WrapError(err, "failed to create feedback")
	}

	// Get author info
	var authorName string
	err = r.db.Pool.QueryRow(ctx, "SELECT name FROM users WHERE id = $1", authorID).Scan(&authorName)
	if err == nil {
		item.Author = &authModel.UserSummary{ID: authorID, Name: authorName}
	}

	if typeStr != nil {
		ft := model.FeedbackType(*typeStr)
		item.Type = &ft
	}
	if visibilityStr != nil {
		v := model.FeedbackVisibility(*visibilityStr)
		item.Visibility = &v
	}

	span.SetStatus(codes.Ok, "")
	return item, nil
}

// CreateComment creates a new comment
func (r *PostgresRepository) CreateComment(ctx context.Context, userID, feedbackID string, content string, parentCommentID *string) (*model.FeedbackComment, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "repository.CreateComment")
	defer span.End()

	commentID := "c-" + uuid.New().String()
	now := time.Now()

	query := `
		INSERT INTO feedback_comments (comment_id, feedback_id, author_id, content, parent_comment_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING comment_id, author_id, content, created_at, parent_comment_id
	`

	comment := &model.FeedbackComment{}
	var authorID string

	err := r.db.Pool.QueryRow(ctx, query, commentID, feedbackID, userID, content, parentCommentID, now, now).Scan(
		&comment.CommentID,
		&authorID,
		&comment.Content,
		&comment.CreatedAt,
		&comment.ParentCommentID,
	)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, errors.WrapError(err, "failed to create comment")
	}

	// Get author info
	var authorName string
	err = r.db.Pool.QueryRow(ctx, "SELECT name FROM users WHERE id = $1", authorID).Scan(&authorName)
	if err == nil {
		comment.Author = &authModel.UserSummary{ID: authorID, Name: authorName}
	}

	span.SetStatus(codes.Ok, "")
	return comment, nil
}

// AddReaction adds a reaction to a feedback item
func (r *PostgresRepository) AddReaction(ctx context.Context, userID, feedbackID string, reactionType string) error {
	ctx, span := otel.Tracer("repository").Start(ctx, "repository.AddReaction")
	defer span.End()

	reactionID := "react-" + uuid.New().String()

	query := `
		INSERT INTO feedback_reactions (reaction_id, feedback_id, user_id, reaction_type, created_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (feedback_id, user_id, reaction_type) DO NOTHING
	`

	_, err := r.db.Pool.Exec(ctx, query, reactionID, feedbackID, userID, reactionType, time.Now())
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return errors.WrapError(err, "failed to add reaction")
	}

	span.SetStatus(codes.Ok, "")
	return nil
}

// RemoveReaction removes a reaction from a feedback item
func (r *PostgresRepository) RemoveReaction(ctx context.Context, userID, feedbackID string, reactionType string) error {
	ctx, span := otel.Tracer("repository").Start(ctx, "repository.RemoveReaction")
	defer span.End()

	query := `
		DELETE FROM feedback_reactions
		WHERE feedback_id = $1 AND user_id = $2 AND reaction_type = $3
	`

	result, err := r.db.Pool.Exec(ctx, query, feedbackID, userID, reactionType)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return errors.WrapError(err, "failed to remove reaction")
	}

	if result.RowsAffected() == 0 {
		return errors.NewValidationError("reaction not found")
	}

	span.SetStatus(codes.Ok, "")
	return nil
}

// GetReactionsCount gets reaction counts for a feedback item
func (r *PostgresRepository) GetReactionsCount(ctx context.Context, feedbackID string) (map[string]int, error) {
	query := `
		SELECT reaction_type, COUNT(*) as count
		FROM feedback_reactions
		WHERE feedback_id = $1
		GROUP BY reaction_type
	`

	rows, err := r.db.Pool.Query(ctx, query, feedbackID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reactions := make(map[string]int)
	for rows.Next() {
		var reactionType string
		var count int
		if err := rows.Scan(&reactionType, &count); err == nil {
			reactions[reactionType] = count
		}
	}

	return reactions, nil
}

// GetCommentsCount gets comment count for a feedback item
func (r *PostgresRepository) GetCommentsCount(ctx context.Context, feedbackID string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM feedback_comments WHERE feedback_id = $1`
	err := r.db.Pool.QueryRow(ctx, query, feedbackID).Scan(&count)
	return count, err
}

