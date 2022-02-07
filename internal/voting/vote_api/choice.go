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

	// TODO call service methods here to actually save vote

	c.Status(http.StatusCreated)
}
