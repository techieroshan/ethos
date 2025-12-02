package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CommunityHandler handles community HTTP requests
type CommunityHandler struct{}

// NewCommunityHandler creates a new community handler
func NewCommunityHandler() *CommunityHandler {
	return &CommunityHandler{}
}

// GetRules handles GET /api/v1/community/rules
func (h *CommunityHandler) GetRules(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"title":   "Community Rules",
		"content": "Please respect others and do not post prohibited content.",
	})
}

