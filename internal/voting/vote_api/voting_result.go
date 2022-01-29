package vote_api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type VotingWinner struct {
	Name  string `json:"name"`
	Score int    `json:"score"`
}

type VotingResult struct {
	Datetime time.Time    `json:"datetime"`
	Total    int          `json:"total"`
	Winner   VotingWinner `json:"winner"`
}

// RetrieveLastVoteResult returns aggregated result of last voting
func (h *VoteAPIHandler) RetrieveLastVoteResult(c *gin.Context) {

	// TODO call service methods here to actually fetch data

	c.JSON(http.StatusOK, VotingResult{
		Datetime: time.Now().Round(1 * time.Hour),
		Total:    100,
		Winner: VotingWinner{
			Name:  "Happy",
			Score: 43,
		},
	})
}

// RetrieveVoteResultHistory returns aggregated results of voting
// for defined period
func (h *VoteAPIHandler) RetrieveVoteResultHistory(c *gin.Context) {

	start := c.Query("start")
	end := c.Query("end")

	_, _ = start, end

	// TODO call service methods here to actually fetch data

	lastHour := time.Now().Round(1 * time.Hour)

	c.JSON(http.StatusOK, []VotingResult{
		{
			Datetime: lastHour.Add(-1 * time.Hour),
			Total:    100,
			Winner: VotingWinner{
				Name:  "Happy",
				Score: 43,
			},
		},
		{
			Datetime: lastHour.Add(-2 * time.Hour),
			Total:    78,
			Winner: VotingWinner{
				Name:  "Happy",
				Score: 51,
			},
		},
		{
			Datetime: lastHour.Add(-3 * time.Hour),
			Total:    15,
			Winner: VotingWinner{
				Name:  "Lyrical",
				Score: 11,
			},
		},
	})
}
