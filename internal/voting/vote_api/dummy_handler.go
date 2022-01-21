package vote_api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Hello just dummy handler to start app, TODO: remove it later
func (h *VoteAPIHandler) Hello(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]string{"msg": "Hello"})
}
