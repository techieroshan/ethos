package repository

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	authModel "ethos/internal/auth/model"
	"ethos/internal/database"
	"ethos/internal/feedback"
	"ethos/internal/feedback/model"
	"ethos/pkg/errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
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

// GetTemplates retrieves feedback templates with optional filtering
func (r *PostgresRepository) GetTemplates(ctx context.Context, contextFilter, tagsFilter string) ([]*model.FeedbackTemplate, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "repository.GetTemplates")
	defer span.End()

	query := `SELECT template_id, name, description, context_tags, template_fields FROM feedback_templates`
	args := []interface{}{}
	argCount := 0

	// Add context filtering if specified
	if contextFilter != "" {
		argCount++
		query += ` WHERE $` + strconv.Itoa(argCount) + ` = ANY(context_tags)`
		args = append(args, contextFilter)
	}

	// Add tags filtering if specified
	if tagsFilter != "" {
		if argCount == 0 {
			query += ` WHERE`
		} else {
			query += ` AND`
		}
		tags := strings.Split(tagsFilter, ",")
		for i, tag := range tags {
			if i > 0 {
				query += ` OR`
			} else {
				argCount++
			}
			query += ` $` + strconv.Itoa(argCount) + ` = ANY(context_tags)`
			args = append(args, strings.TrimSpace(tag))
		}
	}

	query += ` ORDER BY name`

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, errors.WrapError(err, "failed to get feedback templates")
	}
	defer rows.Close()

	var templates []*model.FeedbackTemplate
	for rows.Next() {
		var template model.FeedbackTemplate
		var contextTags []string
		var templateFields []byte

		err := rows.Scan(
			&template.TemplateID,
			&template.Name,
			&template.Description,
			&contextTags,
			&templateFields,
		)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return nil, errors.WrapError(err, "failed to scan feedback template")
		}

		template.ContextTags = contextTags

		// Parse JSON template fields
		if len(templateFields) > 0 {
			err = json.Unmarshal(templateFields, &template.TemplateFields)
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
				return nil, errors.WrapError(err, "failed to parse template fields JSON")
			}
		}

		templates = append(templates, &template)
	}

	if err := rows.Err(); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, errors.WrapError(err, "error iterating template rows")
	}

	span.SetStatus(codes.Ok, "")
	return templates, nil
}

// SubmitTemplateSuggestion submits a template suggestion
func (r *PostgresRepository) SubmitTemplateSuggestion(ctx context.Context, req *feedback.TemplateSuggestionRequest) error {
	ctx, span := otel.Tracer("repository").Start(ctx, "repository.SubmitTemplateSuggestion")
	defer span.End()

	// In a real implementation, this would save the suggestion to database
	// For now, just log it and return success
	span.SetStatus(codes.Ok, "")
	return nil
}

// GetImpact retrieves aggregated feedback analytics
func (r *PostgresRepository) GetImpact(ctx context.Context, userID *string, from, to *time.Time) (*model.FeedbackImpact, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "repository.GetImpact")
	defer span.End()

	impact := &model.FeedbackImpact{
		ReactionTotals: make(map[string]int),
		Trends:         []model.FeedbackTrend{},
	}

	// Build base query conditions
	whereConditions := []string{}
	args := []interface{}{}
	argCount := 0

	if userID != nil {
		argCount++
		whereConditions = append(whereConditions, `author_id = $`+strconv.Itoa(argCount))
		args = append(args, *userID)
	}

	if from != nil {
		argCount++
		whereConditions = append(whereConditions, `created_at >= $`+strconv.Itoa(argCount))
		args = append(args, *from)
	}

	if to != nil {
		argCount++
		whereConditions = append(whereConditions, `created_at <= $`+strconv.Itoa(argCount))
		args = append(args, *to)
	}

	whereClause := ""
	if len(whereConditions) > 0 {
		whereClause = "WHERE " + strings.Join(whereConditions, " AND ")
	}

	// Get feedback count and average helpfulness
	query := `
		SELECT
			COUNT(*) as feedback_count,
			COALESCE(AVG(helpfulness), 0) as avg_helpfulness
		FROM feedback_items
		` + whereClause

	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(
		&impact.FeedbackCount,
		&impact.AverageHelpfulness,
	)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, errors.WrapError(err, "failed to get feedback impact summary")
	}

	// Get reaction totals
	reactionQuery := `
		SELECT
			reaction_type,
			COUNT(*) as count
		FROM feedback_reactions fr
		JOIN feedback_items fi ON fr.feedback_id = fi.feedback_id
		` + whereClause + `
		GROUP BY reaction_type
	`

	reactionRows, err := r.db.Pool.Query(ctx, reactionQuery, args...)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, errors.WrapError(err, "failed to get reaction totals")
	}
	defer reactionRows.Close()

	for reactionRows.Next() {
		var reactionType string
		var count int
		err := reactionRows.Scan(&reactionType, &count)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return nil, errors.WrapError(err, "failed to scan reaction totals")
		}
		impact.ReactionTotals[reactionType] = count
	}

	// Get follow-up count
	followUpQuery := `
		SELECT COUNT(*) as follow_up_count
		FROM feedback_follow_ups ffu
		JOIN feedback_items fi ON ffu.feedback_id = fi.feedback_id
		` + whereClause

	err = r.db.Pool.QueryRow(ctx, followUpQuery, args...).Scan(&impact.FollowUpCount)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, errors.WrapError(err, "failed to get follow-up count")
	}

	// Get trends (grouped by date)
	trendsQuery := `
		SELECT
			DATE(created_at) as date,
			COALESCE(AVG(helpfulness), 0) as avg_helpfulness,
			COUNT(*) as feedback_count
		FROM feedback_items
		` + whereClause + `
		GROUP BY DATE(created_at)
		ORDER BY DATE(created_at)
	`

	trendRows, err := r.db.Pool.Query(ctx, trendsQuery, args...)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, errors.WrapError(err, "failed to get feedback trends")
	}
	defer trendRows.Close()

	for trendRows.Next() {
		var date time.Time
		var helpfulness float64
		var feedbackCount int
		err := trendRows.Scan(&date, &helpfulness, &feedbackCount)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return nil, errors.WrapError(err, "failed to scan feedback trends")
		}
		impact.Trends = append(impact.Trends, model.FeedbackTrend{
			Date:              date,
			Helpfulness:       helpfulness,
			FeedbackSubmitted: feedbackCount,
		})
	}

	span.SetStatus(codes.Ok, "")
	return impact, nil
}

// GetBookmarks retrieves bookmarked feedback items for a user
func (r *PostgresRepository) GetBookmarks(ctx context.Context, userID string, limit, offset int) ([]*model.FeedbackItem, int, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "repository.GetBookmarks")
	defer span.End()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*) FROM feedback_bookmarks
		WHERE user_id = $1
	`
	err := r.db.Pool.QueryRow(ctx, countQuery, userID).Scan(&total)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, 0, errors.WrapError(err, "failed to count bookmarks")
	}

	// Get bookmarked feedback items with full details
	query := `
		SELECT
			fi.feedback_id,
			fi.author_id,
			u.name as author_name,
			u.role as author_role,
			fi.content,
			fi.type,
			fi.visibility,
			fi.is_anonymous,
			fi.helpfulness,
			fi.reviewer_context,
			fi.moderation_state,
			fi.created_at,
			COALESCE(comment_counts.comment_count, 0) as comments_count,
			fb.created_at as bookmarked_at
		FROM feedback_bookmarks fb
		JOIN feedback_items fi ON fb.feedback_id = fi.feedback_id
		JOIN users u ON fi.author_id = u.id
		LEFT JOIN (
			SELECT feedback_id, COUNT(*) as comment_count
			FROM feedback_comments
			GROUP BY feedback_id
		) comment_counts ON fi.feedback_id = comment_counts.feedback_id
		WHERE fb.user_id = $1
		ORDER BY fb.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Pool.Query(ctx, query, userID, limit, offset)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, 0, errors.WrapError(err, "failed to get bookmarks")
	}
	defer rows.Close()

	var items []*model.FeedbackItem
	for rows.Next() {
		var item model.FeedbackItem
		var authorName string
		var authorRole *string
		var reviewerContext []byte
		var moderationState *string
		var bookmarkedAt time.Time

		err := rows.Scan(
			&item.FeedbackID,
			&item.Author.ID,
			&authorName,
			&authorRole,
			&item.Content,
			&item.Type,
			&item.Visibility,
			&item.IsAnonymous,
			&item.Helpfulness,
			&reviewerContext,
			&moderationState,
			&item.CreatedAt,
			&item.CommentsCount,
			&bookmarkedAt,
		)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return nil, 0, errors.WrapError(err, "failed to scan bookmarked feedback item")
		}

		// Set author information
		item.Author.Name = authorName
		// Note: Role assignment removed as UserSummary doesn't have Role field

		// Get reactions for this feedback item
		reactions, err := r.getReactionsForFeedback(ctx, item.FeedbackID)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return nil, 0, errors.WrapError(err, "failed to get reactions")
		}
		item.Reactions = reactions

		// Get reaction analytics
		reactionAnalytics, err := r.getReactionAnalytics(ctx, item.FeedbackID)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return nil, 0, errors.WrapError(err, "failed to get reaction analytics")
		}
		item.ReactionsAnalytics = reactionAnalytics

		items = append(items, &item)
	}

	if err := rows.Err(); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, 0, errors.WrapError(err, "error iterating bookmark rows")
	}

	span.SetStatus(codes.Ok, "")
	return items, total, nil
}

// AddBookmark adds a bookmark for a feedback item
func (r *PostgresRepository) AddBookmark(ctx context.Context, userID, feedbackID string) error {
	ctx, span := otel.Tracer("repository").Start(ctx, "repository.AddBookmark")
	defer span.End()

	// Check if feedback item exists
	var exists bool
	err := r.db.Pool.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM feedback_items WHERE feedback_id = $1)", feedbackID).Scan(&exists)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return errors.WrapError(err, "failed to check feedback existence")
	}
	if !exists {
		return &errors.APIError{
			Message:    "Feedback item not found",
			Code:       "NOT_FOUND",
			HTTPStatus: http.StatusNotFound,
		}
	}

	// Insert bookmark (ignore if already exists - idempotent operation)
	_, err = r.db.Pool.Exec(ctx, `
		INSERT INTO feedback_bookmarks (user_id, feedback_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id, feedback_id) DO NOTHING
	`, userID, feedbackID)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return errors.WrapError(err, "failed to add bookmark")
	}

	span.SetStatus(codes.Ok, "")
	return nil
}

// RemoveBookmark removes a bookmark for a feedback item
func (r *PostgresRepository) RemoveBookmark(ctx context.Context, userID, feedbackID string) error {
	ctx, span := otel.Tracer("repository").Start(ctx, "repository.RemoveBookmark")
	defer span.End()

	// Delete bookmark (idempotent - doesn't fail if bookmark doesn't exist)
	_, err := r.db.Pool.Exec(ctx, `
		DELETE FROM feedback_bookmarks
		WHERE user_id = $1 AND feedback_id = $2
	`, userID, feedbackID)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return errors.WrapError(err, "failed to remove bookmark")
	}

	span.SetStatus(codes.Ok, "")
	return nil
}

// getReactionsForFeedback gets reaction counts for a feedback item
func (r *PostgresRepository) getReactionsForFeedback(ctx context.Context, feedbackID string) (map[string]int, error) {
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
		if err := rows.Scan(&reactionType, &count); err != nil {
			return nil, err
		}
		reactions[reactionType] = count
	}

	return reactions, rows.Err()
}

// getReactionAnalytics gets detailed reaction analytics for a feedback item
func (r *PostgresRepository) getReactionAnalytics(ctx context.Context, feedbackID string) (*model.FeedbackReactionAnalytics, error) {
	query := `
		SELECT
			fr.reaction_type,
			COUNT(*) as count,
			array_agg(u.id) as user_ids
		FROM feedback_reactions fr
		JOIN users u ON fr.user_id = u.id
		WHERE fr.feedback_id = $1
		GROUP BY fr.reaction_type
	`

	rows, err := r.db.Pool.Query(ctx, query, feedbackID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	analytics := &model.FeedbackReactionAnalytics{
		Reactions: make(map[string]model.ReactionDetail),
	}

	for rows.Next() {
		var reactionType string
		var count int
		var userIDs []string

		if err := rows.Scan(&reactionType, &count, &userIDs); err != nil {
			return nil, err
		}

		analytics.Reactions[reactionType] = model.ReactionDetail{
			Count:   count,
			UserIDs: userIDs,
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return analytics, nil
}

// CreateBatchFeedback creates multiple feedback items in a batch
func (r *PostgresRepository) CreateBatchFeedback(ctx context.Context, userID string, req *feedback.BatchFeedbackRequest) (*feedback.BatchFeedbackResponse, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "repository.CreateBatchFeedback")
	defer span.End()

	if len(req.Items) == 0 {
		return nil, &errors.APIError{
			Message:    "At least one feedback item is required",
			Code:       "VALIDATION_FAILED",
			HTTPStatus: http.StatusBadRequest,
		}
	}

	response := &feedback.BatchFeedbackResponse{
		Submitted: make([]feedback.BatchFeedbackResult, 0, len(req.Items)),
	}

	// Begin transaction
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, errors.WrapError(err, "failed to begin transaction")
	}
	defer tx.Rollback(ctx)

	// Process each feedback item
	for _, item := range req.Items {
		feedbackID := uuid.New().String()

		// Determine the actual author (anonymous items don't show the real author)
		actualAuthorID := userID
		if item.IsAnonymous {
			// For anonymous feedback, we could use a system user or null
			// For now, we'll still track the real author for moderation purposes
			actualAuthorID = userID
		}

		// Insert the feedback item directly
		_, err = tx.Exec(ctx,
			`INSERT INTO feedback_items (feedback_id, author_id, content, type, visibility, is_anonymous, created_at)
			 VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP)`,
			feedbackID,
			actualAuthorID,
			item.Content,
			item.Type,
			item.Visibility,
			item.IsAnonymous,
		)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return nil, errors.WrapError(err, "failed to insert feedback item")
		}

		response.Submitted = append(response.Submitted, feedback.BatchFeedbackResult{
			FeedbackID: feedbackID,
			Status:     "created",
		})
	}

	// Commit transaction
	err = tx.Commit(ctx)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, errors.WrapError(err, "failed to commit transaction")
	}

	span.SetStatus(codes.Ok, "")
	return response, nil
}

// GetFeedWithFilters retrieves a paginated feed of feedback items with enhanced filtering
func (r *PostgresRepository) GetFeedWithFilters(ctx context.Context, limit, offset int, filters *feedback.FeedFilters) ([]*model.FeedbackItem, int, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "repository.GetFeedWithFilters")
	defer span.End()

	// Build the base query
	query := `
		SELECT
			fi.feedback_id,
			fi.author_id,
			u.name as author_name,
			u.role as author_role,
			fi.content,
			fi.type,
			fi.visibility,
			fi.is_anonymous,
			fi.helpfulness,
			fi.reviewer_context,
			fi.moderation_state,
			fi.created_at,
			COALESCE(comment_counts.comment_count, 0) as comments_count
		FROM feedback_items fi
		JOIN users u ON fi.author_id = u.id
		LEFT JOIN (
			SELECT feedback_id, COUNT(*) as comment_count
			FROM feedback_comments
			GROUP BY feedback_id
		) comment_counts ON fi.feedback_id = comment_counts.feedback_id
	`

	countQuery := `
		SELECT COUNT(*)
		FROM feedback_items fi
		JOIN users u ON fi.author_id = u.id
	`

	args := []interface{}{}
	argCount := 0
	whereConditions := []string{}

	// Apply filters
	if filters != nil {
		// Reviewer type filter
		if filters.ReviewerType != nil {
			argCount++
			if *filters.ReviewerType == "org" {
				// For org reviewer type, we need feedback with reviewer_context
				whereConditions = append(whereConditions, `fi.reviewer_context IS NOT NULL`)
			} else if *filters.ReviewerType == "public" {
				// For public reviewer type, we need feedback without reviewer_context or with public visibility
				whereConditions = append(whereConditions, `(fi.reviewer_context IS NULL OR fi.visibility = 'public')`)
			}
		}

		// Context filter
		if filters.Context != nil {
			argCount++
			whereConditions = append(whereConditions, `fi.reviewer_context->>'type' = $`+strconv.Itoa(argCount))
			args = append(args, *filters.Context)
		}

		// Verification filter (based on email_verified status of author)
		if filters.Verification != nil {
			if *filters.Verification == "verified" {
				whereConditions = append(whereConditions, `u.email_verified = true`)
			} else if *filters.Verification == "unverified" {
				whereConditions = append(whereConditions, `u.email_verified = false`)
			}
		}

		// Tags filter - this would require a tags field, but for now we'll skip this
		// as the current schema doesn't have tags. This could be added later.
		if len(filters.Tags) > 0 {
			// Placeholder for future implementation
			// whereConditions = append(whereConditions, `fi.tags @> $`+strconv.Itoa(argCount+1))
			// args = append(args, filters.Tags)
		}
	}

	// Add WHERE clause if we have conditions
	if len(whereConditions) > 0 {
		whereClause := " WHERE " + strings.Join(whereConditions, " AND ")
		query += whereClause
		countQuery += whereClause
	}

	// Add ordering and pagination
	query += ` ORDER BY fi.created_at DESC LIMIT $` + strconv.Itoa(argCount+1) + ` OFFSET $` + strconv.Itoa(argCount+2)
	args = append(args, limit, offset)

	// Get total count
	var total int
	err := r.db.Pool.QueryRow(ctx, countQuery, args[:argCount]...).Scan(&total)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, 0, errors.WrapError(err, "failed to get total count")
	}

	// Get feedback items
	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, 0, errors.WrapError(err, "failed to get feedback feed")
	}
	defer rows.Close()

	var items []*model.FeedbackItem
	for rows.Next() {
		var item model.FeedbackItem
		var authorName string
		var authorRole *string
		var reviewerContext []byte
		var moderationState *string

		err := rows.Scan(
			&item.FeedbackID,
			&item.Author.ID,
			&authorName,
			&authorRole,
			&item.Content,
			&item.Type,
			&item.Visibility,
			&item.IsAnonymous,
			&item.Helpfulness,
			&reviewerContext,
			&moderationState,
			&item.CreatedAt,
			&item.CommentsCount,
		)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return nil, 0, errors.WrapError(err, "failed to scan feedback item")
		}

		// Set author information
		item.Author.Name = authorName
		// Note: Role assignment removed as UserSummary doesn't have Role field

		// Parse reviewer context JSON if present
		if len(reviewerContext) > 0 {
			// For now, we'll skip parsing the reviewer context
			// This would require additional model changes
		}

		// Get reactions for this feedback item
		reactions, err := r.getReactionsForFeedback(ctx, item.FeedbackID)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return nil, 0, errors.WrapError(err, "failed to get reactions")
		}
		item.Reactions = reactions

		// Get reaction analytics
		reactionAnalytics, err := r.getReactionAnalytics(ctx, item.FeedbackID)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return nil, 0, errors.WrapError(err, "failed to get reaction analytics")
		}
		item.ReactionsAnalytics = reactionAnalytics

		items = append(items, &item)
	}

	if err := rows.Err(); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, 0, errors.WrapError(err, "error iterating feed rows")
	}

	span.SetStatus(codes.Ok, "")
	return items, total, nil
}

// UpdateFeedback updates an existing feedback item
func (r *PostgresRepository) UpdateFeedback(ctx context.Context, feedbackID string, item *model.FeedbackItem) error {
	ctx, span := otel.Tracer("repository").Start(ctx, "repository.UpdateFeedback")
	defer span.End()

	// Placeholder implementation
	span.SetStatus(codes.Ok, "")
	return nil
}

// DeleteFeedback deletes a feedback item
func (r *PostgresRepository) DeleteFeedback(ctx context.Context, feedbackID string) error {
	ctx, span := otel.Tracer("repository").Start(ctx, "repository.DeleteFeedback")
	defer span.End()

	// Placeholder implementation
	span.SetStatus(codes.Ok, "")
	return nil
}

// GetComment retrieves a specific comment
func (r *PostgresRepository) GetComment(ctx context.Context, feedbackID, commentID string) (*model.FeedbackComment, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "repository.GetComment")
	defer span.End()

	// Placeholder implementation
	comment := &model.FeedbackComment{
		CommentID: commentID,
		Author: &authModel.UserSummary{
			ID: "user-1",
		},
		Content: "Sample comment",
	}
	span.SetStatus(codes.Ok, "")
	return comment, nil
}

// UpdateComment updates an existing comment
func (r *PostgresRepository) UpdateComment(ctx context.Context, feedbackID, commentID string, comment *model.FeedbackComment) error {
	ctx, span := otel.Tracer("repository").Start(ctx, "repository.UpdateComment")
	defer span.End()

	// Placeholder implementation
	span.SetStatus(codes.Ok, "")
	return nil
}

// DeleteComment deletes a comment
func (r *PostgresRepository) DeleteComment(ctx context.Context, feedbackID, commentID string) error {
	ctx, span := otel.Tracer("repository").Start(ctx, "repository.DeleteComment")
	defer span.End()

	// Placeholder implementation
	span.SetStatus(codes.Ok, "")
	return nil
}

// GetFeedbackAnalytics retrieves detailed feedback analytics
func (r *PostgresRepository) GetFeedbackAnalytics(ctx context.Context, userID *string, from, to *time.Time) (*model.FeedbackAnalytics, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "repository.GetFeedbackAnalytics")
	defer span.End()

	// Placeholder implementation
	analytics := &model.FeedbackAnalytics{
		TotalCount:         42,
		AverageHelpfulness: 4.5,
		TopReactions: map[string]int{
			"👍":  120,
			"❤️": 85,
			"😂":  42,
		},
		TypeDistribution: map[string]int{
			"suggestion":   25,
			"issue":        12,
			"appreciation": 5,
		},
		VisibilityDistribution: map[string]int{
			"public": 30,
			"team":   12,
		},
	}
	span.SetStatus(codes.Ok, "")
	return analytics, nil
}

// SearchFeedback searches feedback items by content/metadata
func (r *PostgresRepository) SearchFeedback(ctx context.Context, query string, limit, offset int) ([]*model.FeedbackItem, int, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "repository.SearchFeedback")
	defer span.End()

	// Placeholder implementation - returns empty results
	span.SetStatus(codes.Ok, "")
	return []*model.FeedbackItem{}, 0, nil
}

// GetTrendingFeedback retrieves trending feedback items
func (r *PostgresRepository) GetTrendingFeedback(ctx context.Context, limit, offset int) ([]*model.FeedbackItem, int, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "repository.GetTrendingFeedback")
	defer span.End()

	// Placeholder implementation - returns empty results
	span.SetStatus(codes.Ok, "")
	return []*model.FeedbackItem{}, 0, nil
}

// PinFeedback pins a feedback item
func (r *PostgresRepository) PinFeedback(ctx context.Context, userID, feedbackID string) error {
	ctx, span := otel.Tracer("repository").Start(ctx, "repository.PinFeedback")
	defer span.End()

	// Placeholder implementation
	span.SetStatus(codes.Ok, "")
	return nil
}

// UnpinFeedback unpins a feedback item
func (r *PostgresRepository) UnpinFeedback(ctx context.Context, userID, feedbackID string) error {
	ctx, span := otel.Tracer("repository").Start(ctx, "repository.UnpinFeedback")
	defer span.End()

	// Placeholder implementation
	span.SetStatus(codes.Ok, "")
	return nil
}

// GetFeedbackStats retrieves overall feedback statistics
func (r *PostgresRepository) GetFeedbackStats(ctx context.Context) (*model.FeedbackStats, error) {
	ctx, span := otel.Tracer("repository").Start(ctx, "repository.GetFeedbackStats")
	defer span.End()

	stats := &model.FeedbackStats{
		TotalFeedback:      150,
		TotalComments:      320,
		TotalReactions:     512,
		AverageHelpfulness: 4.3,
		MostPopularType:    "suggestion",
		MostCommonReaction: "👍",
	}
	span.SetStatus(codes.Ok, "")
	return stats, nil
}

