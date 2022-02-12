package vote_api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type VoteCategory struct {
	ID   string `json:"id"`
	Name string `json:"name" binding:"required"`
}

// RetrieveVotecategories returns music categories which user can choose
func (h *VoteAPIHandler) RetrieveVoteCategories(c *gin.Context) {
	categories, err := h.service.GetCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIBaseError{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, categories)
}
