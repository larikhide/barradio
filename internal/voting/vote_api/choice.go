package vote_api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// MakeCategoryChoice saves user category choice (vote)
func (h *VoteAPIHandler) MakeCategoryChoice(c *gin.Context) {
	var choice VoteCategory
	err := c.ShouldBindJSON(&choice)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIBaseError{Message: err.Error()})
		return
	}

	err = h.service.ChooseCategory(choice.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIBaseError{Message: err.Error()})
		return
	}

	votes, err := h.service.CurrentVoting()
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIBaseError{Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, mapVotingResult(votes))
}
