package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	authModel "ethos/internal/auth/model"
	"ethos/internal/people"
	"ethos/internal/people/service"
	"ethos/pkg/errors"
)

// PeopleHandler handles people search HTTP requests
type PeopleHandler struct {
	service service.Service
}

// NewPeopleHandler creates a new people handler
func NewPeopleHandler(svc service.Service) *PeopleHandler {
	return &PeopleHandler{
		service: svc,
	}
}

// SearchPeople handles GET /api/v1/people/search
func (h *PeopleHandler) SearchPeople(c *gin.Context) {
	query := c.Query("q")
	limit := c.DefaultQuery("limit", "25")
	offset := c.DefaultQuery("offset", "0")

	limitInt := 25
	offsetInt := 0
	if l, err := strconv.Atoi(limit); err == nil && l > 0 {
		limitInt = l
	}
	if o, err := strconv.Atoi(offset); err == nil && o >= 0 {
		offsetInt = o
	}

	// Parse filtering parameters
	filters := &people.PeopleSearchFilters{}

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

	var results []*authModel.UserProfile
	var count int
	var err error

	// Use filtered search if any filters are provided
	if filters.ReviewerType != nil || filters.Context != nil || filters.Verification != nil || len(filters.Tags) > 0 {
		results, count, err = h.service.SearchPeopleWithFilters(c.Request.Context(), query, limitInt, offsetInt, filters)
	} else {
		// Fallback to original SearchPeople for backward compatibility
		results, count, err = h.service.SearchPeople(c.Request.Context(), query, limitInt, offsetInt)
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
		"results": results,
		"count":   count,
	})
}

// GetRecommendations handles GET /api/v1/people/recommendations
func (h *PeopleHandler) GetRecommendations(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid token",
			"code":  "AUTH_TOKEN_INVALID",
		})
		return
	}

	recommendations, err := h.service.GetRecommendations(c.Request.Context(), userID.(string))
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
		"recommendations": recommendations,
	})
}

