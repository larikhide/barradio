package vote_api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type VoteCategory struct {
	ID   string `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
}

// RetrieveVotecategories returns music categories which user can choose
func (h *VoteAPIHandler) RetrieveVoteCategories(c *gin.Context) {
	// TODO call service method to fetch data from storage

	c.JSON(http.StatusOK, []VoteCategory{
		{ID: "c100", Name: "Happy"},
		{ID: "c101", Name: "Calm"},
		{ID: "c102", Name: "Lyrical"},
	})
}
