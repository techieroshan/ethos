package handler

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	feedbackPkg "ethos/internal/feedback"
	"ethos/internal/feedback/model"
	"ethos/internal/feedback/service"
	"ethos/pkg/errors"

	"github.com/gin-gonic/gin"
)

// FeedbackHandler handles feedback HTTP requests
type FeedbackHandler struct {
	service service.Service
}

// NewFeedbackHandler creates a new feedback handler
func NewFeedbackHandler(svc service.Service) *FeedbackHandler {
	return &FeedbackHandler{
		service: svc,
	}
}

// GetFeed handles GET /api/v1/feedback/feed
func (h *FeedbackHandler) GetFeed(c *gin.Context) {
	limit := c.DefaultQuery("limit", "20")
	offset := c.DefaultQuery("offset", "0")

	limitInt := 20
	offsetInt := 0
	if l, err := strconv.Atoi(limit); err == nil && l > 0 {
		limitInt = l
	}
	if o, err := strconv.Atoi(offset); err == nil && o >= 0 {
		offsetInt = o
	}

	// Parse filtering parameters
	filters := &feedbackPkg.FeedFilters{}

	if reviewerType := c.Query("reviewer_type"); reviewerType != "" {
		filters.ReviewerType = &reviewerType
	}

	if contextFilter := c.Query("context"); contextFilter != "" {
		filters.Context = &contextFilter
	}

	if verification := c.Query("verification"); verification != "" {
		filters.Verification = &verification
	}

	if tagsStr := c.Query("tags"); tagsStr != "" {
		filters.Tags = strings.Split(tagsStr, ",")
		// Trim spaces from tags
		for i, tag := range filters.Tags {
			filters.Tags[i] = strings.TrimSpace(tag)
		}
	}

	var items []*model.FeedbackItem
	var count int
	var err error

	// Use filtered feed if any filters are provided
	if filters.ReviewerType != nil || filters.Context != nil || filters.Verification != nil || len(filters.Tags) > 0 {
		items, count, err = h.service.GetFeedWithFilters(c.Request.Context(), limitInt, offsetInt, filters)
	} else {
		// Fallback to original GetFeed for backward compatibility
		items, count, err = h.service.GetFeed(c.Request.Context(), limitInt, offsetInt)
	}

	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			c.JSON(apiErr.HTTPStatus, gin.H{
				"error": apiErr.Message,
				"code":  apiErr.Code,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
			"code":  "SERVER_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"results": items,
		"count":   count,
	})
}

// GetFeedbackByID handles GET /api/v1/feedback/:feedback_id
func (h *FeedbackHandler) GetFeedbackByID(c *gin.Context) {
	feedbackID := c.Param("feedback_id")

	item, err := h.service.GetFeedbackByID(c.Request.Context(), feedbackID)
	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			c.JSON(apiErr.HTTPStatus, gin.H{
				"error": apiErr.Message,
				"code":  apiErr.Code,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
			"code":  "SERVER_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, item)
}

// GetComments handles GET /api/v1/feedback/:feedback_id/comments
func (h *FeedbackHandler) GetComments(c *gin.Context) {
	feedbackID := c.Param("feedback_id")
	limit := c.DefaultQuery("limit", "20")
	offset := c.DefaultQuery("offset", "0")

	limitInt := 20
	offsetInt := 0
	if l, err := strconv.Atoi(limit); err == nil && l > 0 {
		limitInt = l
	}
	if o, err := strconv.Atoi(offset); err == nil && o >= 0 {
		offsetInt = o
	}

	comments, count, err := h.service.GetComments(c.Request.Context(), feedbackID, limitInt, offsetInt)
	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			c.JSON(apiErr.HTTPStatus, gin.H{
				"error": apiErr.Message,
				"code":  apiErr.Code,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
			"code":  "SERVER_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"comments": comments,
		"count":    count,
	})
}

// CreateFeedback handles POST /api/v1/feedback
func (h *FeedbackHandler) CreateFeedback(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid token",
			"code":  "AUTH_TOKEN_INVALID",
		})
		return
	}

	var req service.CreateFeedbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Validation failed",
			"code":  "VALIDATION_FAILED",
		})
		return
	}

	item, err := h.service.CreateFeedback(c.Request.Context(), userID.(string), &req)
	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			c.JSON(apiErr.HTTPStatus, gin.H{
				"error": apiErr.Message,
				"code":  apiErr.Code,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
			"code":  "SERVER_ERROR",
		})
		return
	}

	c.JSON(http.StatusCreated, item)
}

// CreateComment handles POST /api/v1/feedback/:feedback_id/comments
func (h *FeedbackHandler) CreateComment(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid token",
			"code":  "AUTH_TOKEN_INVALID",
		})
		return
	}

	feedbackID := c.Param("feedback_id")

	var req service.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Validation failed",
			"code":  "VALIDATION_FAILED",
		})
		return
	}

	comment, err := h.service.CreateComment(c.Request.Context(), userID.(string), feedbackID, &req)
	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			c.JSON(apiErr.HTTPStatus, gin.H{
				"error": apiErr.Message,
				"code":  apiErr.Code,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
			"code":  "SERVER_ERROR",
		})
		return
	}

	c.JSON(http.StatusCreated, comment)
}

// AddReaction handles POST /api/v1/feedback/:feedback_id/react
func (h *FeedbackHandler) AddReaction(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid token",
			"code":  "AUTH_TOKEN_INVALID",
		})
		return
	}

	feedbackID := c.Param("feedback_id")

	var req service.AddReactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Validation failed",
			"code":  "VALIDATION_FAILED",
		})
		return
	}

	err := h.service.AddReaction(c.Request.Context(), userID.(string), feedbackID, req.ReactionType)
	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			c.JSON(apiErr.HTTPStatus, gin.H{
				"error": apiErr.Message,
				"code":  apiErr.Code,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
			"code":  "SERVER_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"feedback_id": feedbackID,
		"message":     "Reaction added",
	})
}

// RemoveReaction handles DELETE /api/v1/feedback/:feedback_id/react
func (h *FeedbackHandler) RemoveReaction(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid token",
			"code":  "AUTH_TOKEN_INVALID",
		})
		return
	}

	feedbackID := c.Param("feedback_id")
	reactionType := c.Query("reaction_type")
	if reactionType == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "reaction_type query parameter is required",
			"code":  "VALIDATION_FAILED",
		})
		return
	}

	err := h.service.RemoveReaction(c.Request.Context(), userID.(string), feedbackID, reactionType)
	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			c.JSON(apiErr.HTTPStatus, gin.H{
				"error": apiErr.Message,
				"code":  apiErr.Code,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
			"code":  "SERVER_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"feedback_id": feedbackID,
		"message":     "Reaction removed",
	})
}

// GetTemplates handles GET /api/feedback/templates
func (h *FeedbackHandler) GetTemplates(c *gin.Context) {
	contextFilter := c.Query("context")
	tagsFilter := c.Query("tags")

	templates, err := h.service.GetTemplates(c.Request.Context(), contextFilter, tagsFilter)
	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			c.JSON(apiErr.HTTPStatus, gin.H{
				"error": apiErr.Message,
				"code":  apiErr.Code,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
			"code":  "SERVER_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"results": templates,
	})
}

// PostTemplateSuggestions handles POST /api/feedback/template_suggestions
func (h *FeedbackHandler) PostTemplateSuggestions(c *gin.Context) {
	var req feedbackPkg.TemplateSuggestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
			"code":  "VALIDATION_FAILED",
		})
		return
	}

	err := h.service.SubmitTemplateSuggestion(c.Request.Context(), &req)
	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			c.JSON(apiErr.HTTPStatus, gin.H{
				"error": apiErr.Message,
				"code":  apiErr.Code,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
			"code":  "SERVER_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "suggestion_received",
	})
}

// GetImpact handles GET /api/feedback/impact
func (h *FeedbackHandler) GetImpact(c *gin.Context) {
	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid token",
			"code":  "AUTH_TOKEN_INVALID",
		})
		return
	}

	// Parse optional query parameters
	var targetUserID *string
	if uid := c.Query("user_id"); uid != "" {
		targetUserID = &uid
	}

	var from, to *time.Time
	if fromStr := c.Query("from"); fromStr != "" {
		if parsed, err := time.Parse("2006-01-02", fromStr); err == nil {
			from = &parsed
		}
	}
	if toStr := c.Query("to"); toStr != "" {
		if parsed, err := time.Parse("2006-01-02", toStr); err == nil {
			to = &parsed
		}
	}

	impact, err := h.service.GetImpact(c.Request.Context(), targetUserID, from, to)
	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			c.JSON(apiErr.HTTPStatus, gin.H{
				"error": apiErr.Message,
				"code":  apiErr.Code,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
			"code":  "SERVER_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, impact)
}

// GetBookmarks handles GET /api/feedback/bookmarks
func (h *FeedbackHandler) GetBookmarks(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid token",
			"code":  "AUTH_TOKEN_INVALID",
		})
		return
	}

	limit := c.DefaultQuery("limit", "20")
	offset := c.DefaultQuery("offset", "0")

	limitInt := 20
	offsetInt := 0
	if l, err := strconv.Atoi(limit); err == nil && l > 0 {
		limitInt = l
	}
	if o, err := strconv.Atoi(offset); err == nil && o >= 0 {
		offsetInt = o
	}

	bookmarks, count, err := h.service.GetBookmarks(c.Request.Context(), userID.(string), limitInt, offsetInt)
	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			c.JSON(apiErr.HTTPStatus, gin.H{
				"error": apiErr.Message,
				"code":  apiErr.Code,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
			"code":  "SERVER_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"results": bookmarks,
		"count":   count,
	})
}

// AddBookmark handles POST /api/feedback/bookmarks/:feedback_id
func (h *FeedbackHandler) AddBookmark(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid token",
			"code":  "AUTH_TOKEN_INVALID",
		})
		return
	}

	feedbackID := c.Param("feedback_id")
	if feedbackID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Feedback ID is required",
			"code":  "VALIDATION_FAILED",
		})
		return
	}

	err := h.service.AddBookmark(c.Request.Context(), userID.(string), feedbackID)
	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			c.JSON(apiErr.HTTPStatus, gin.H{
				"error": apiErr.Message,
				"code":  apiErr.Code,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
			"code":  "SERVER_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "bookmark_added",
	})
}

// RemoveBookmark handles DELETE /api/feedback/bookmarks/:feedback_id
func (h *FeedbackHandler) RemoveBookmark(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid token",
			"code":  "AUTH_TOKEN_INVALID",
		})
		return
	}

	feedbackID := c.Param("feedback_id")
	if feedbackID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Feedback ID is required",
			"code":  "VALIDATION_FAILED",
		})
		return
	}

	err := h.service.RemoveBookmark(c.Request.Context(), userID.(string), feedbackID)
	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			c.JSON(apiErr.HTTPStatus, gin.H{
				"error": apiErr.Message,
				"code":  apiErr.Code,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
			"code":  "SERVER_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "bookmark_removed",
	})
}

// ExportFeedback handles GET /api/feedback/export
func (h *FeedbackHandler) ExportFeedback(c *gin.Context) {
	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid token",
			"code":  "AUTH_TOKEN_INVALID",
		})
		return
	}

	// Parse format parameter
	format := c.DefaultQuery("format", "json")
	if format != "json" && format != "csv" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid format. Supported formats: json, csv",
			"code":  "VALIDATION_FAILED",
		})
		return
	}

	// Parse filtering parameters (same as feed filtering)
	filters := &feedbackPkg.FeedFilters{}

	if reviewerType := c.Query("reviewer_type"); reviewerType != "" {
		filters.ReviewerType = &reviewerType
	}

	if contextFilter := c.Query("context"); contextFilter != "" {
		filters.Context = &contextFilter
	}

	if verification := c.Query("verification"); verification != "" {
		filters.Verification = &verification
	}

	if tagsStr := c.Query("tags"); tagsStr != "" {
		filters.Tags = strings.Split(tagsStr, ",")
		// Trim spaces from tags
		for i, tag := range filters.Tags {
			filters.Tags[i] = strings.TrimSpace(tag)
		}
	}

	exportResponse, err := h.service.ExportFeedback(c.Request.Context(), filters, format)
	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			c.JSON(apiErr.HTTPStatus, gin.H{
				"error": apiErr.Message,
				"code":  apiErr.Code,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
			"code":  "SERVER_ERROR",
		})
		return
	}

	// Set appropriate headers for file download
	c.Header("Content-Type", exportResponse.ContentType)
	c.Header("Content-Disposition", "attachment; filename="+exportResponse.Filename)

	// Return the export data directly
	c.String(http.StatusOK, exportResponse.Data)
}

// CreateBatchFeedback handles POST /api/feedback/batch
func (h *FeedbackHandler) CreateBatchFeedback(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid token",
			"code":  "AUTH_TOKEN_INVALID",
		})
		return
	}

	var req feedbackPkg.BatchFeedbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
			"code":  "VALIDATION_FAILED",
		})
		return
	}

	// Validate that at least one item is provided
	if len(req.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "At least one feedback item is required",
			"code":  "VALIDATION_FAILED",
		})
		return
	}

	response, err := h.service.CreateBatchFeedback(c.Request.Context(), userID.(string), &req)
	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			c.JSON(apiErr.HTTPStatus, gin.H{
				"error": apiErr.Message,
				"code":  apiErr.Code,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
			"code":  "SERVER_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// UpdateFeedback handles PUT /api/v1/feedback/:feedback_id
func (h *FeedbackHandler) UpdateFeedback(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid token",
			"code":  "AUTH_TOKEN_INVALID",
		})
		return
	}

	feedbackID := c.Param("feedback_id")
	var req service.UpdateFeedbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
			"code":  "VALIDATION_FAILED",
		})
		return
	}

	item, err := h.service.UpdateFeedback(c.Request.Context(), userID.(string), feedbackID, &req)
	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			c.JSON(apiErr.HTTPStatus, gin.H{
				"error": apiErr.Message,
				"code":  apiErr.Code,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
			"code":  "SERVER_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, item)
}

// DeleteFeedback handles DELETE /api/v1/feedback/:feedback_id
func (h *FeedbackHandler) DeleteFeedback(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid token",
			"code":  "AUTH_TOKEN_INVALID",
		})
		return
	}

	feedbackID := c.Param("feedback_id")
	err := h.service.DeleteFeedback(c.Request.Context(), userID.(string), feedbackID)
	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			c.JSON(apiErr.HTTPStatus, gin.H{
				"error": apiErr.Message,
				"code":  apiErr.Code,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
			"code":  "SERVER_ERROR",
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// UpdateComment handles PUT /api/v1/feedback/:feedback_id/comments/:comment_id
func (h *FeedbackHandler) UpdateComment(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid token",
			"code":  "AUTH_TOKEN_INVALID",
		})
		return
	}

	feedbackID := c.Param("feedback_id")
	commentID := c.Param("comment_id")

	var req service.UpdateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
			"code":  "VALIDATION_FAILED",
		})
		return
	}

	comment, err := h.service.UpdateComment(c.Request.Context(), userID.(string), feedbackID, commentID, &req)
	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			c.JSON(apiErr.HTTPStatus, gin.H{
				"error": apiErr.Message,
				"code":  apiErr.Code,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
			"code":  "SERVER_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, comment)
}

// DeleteComment handles DELETE /api/v1/feedback/:feedback_id/comments/:comment_id
func (h *FeedbackHandler) DeleteComment(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid token",
			"code":  "AUTH_TOKEN_INVALID",
		})
		return
	}

	feedbackID := c.Param("feedback_id")
	commentID := c.Param("comment_id")

	err := h.service.DeleteComment(c.Request.Context(), userID.(string), feedbackID, commentID)
	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			c.JSON(apiErr.HTTPStatus, gin.H{
				"error": apiErr.Message,
				"code":  apiErr.Code,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
			"code":  "SERVER_ERROR",
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// GetFeedbackAnalytics handles GET /api/v1/feedback/analytics
func (h *FeedbackHandler) GetFeedbackAnalytics(c *gin.Context) {
	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid token",
			"code":  "AUTH_TOKEN_INVALID",
		})
		return
	}

	// Parse optional date range
	var from, to *time.Time
	if fromStr := c.Query("from"); fromStr != "" {
		if parsedFrom, err := time.Parse(time.RFC3339, fromStr); err == nil {
			from = &parsedFrom
		}
	}
	if toStr := c.Query("to"); toStr != "" {
		if parsedTo, err := time.Parse(time.RFC3339, toStr); err == nil {
			to = &parsedTo
		}
	}

	analytics, err := h.service.GetFeedbackAnalytics(c.Request.Context(), nil, from, to)
	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			c.JSON(apiErr.HTTPStatus, gin.H{
				"error": apiErr.Message,
				"code":  apiErr.Code,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
			"code":  "SERVER_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, analytics)
}
