package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"ethos/internal/feedback/service"
	"ethos/pkg/errors"
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

	items, count, err := h.service.GetFeed(c.Request.Context(), limitInt, offsetInt)
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

