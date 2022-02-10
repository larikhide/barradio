package vote_api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/larikhide/barradio/internal/voting/vote_service"
)

// VoteAPIHandler is container which stores handler depencies
type VoteAPIHandler struct {
	service             vote_service.VoteService
	defaultHistoryDepth time.Duration
}

func NewVoteAPIHandler(service vote_service.VoteService, defaultHistoryDepth time.Duration) VoteAPIHandler {
	return VoteAPIHandler{service: service, defaultHistoryDepth: defaultHistoryDepth}
}

// Hello just dummy handler to check if app is working
func (h *VoteAPIHandler) Hello(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]string{"message": "Hello"})
}
