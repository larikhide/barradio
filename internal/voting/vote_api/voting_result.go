package vote_api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type CategoryScore struct {
	Name  string `json:"name"`
	Score int    `json:"score"`
}

type VotingResult struct {
	Datetime time.Time       `json:"datetime"`
	Total    int             `json:"total"`
	Results  []CategoryScore `json:"results"`
}

// RetrieveLastVoteResult returns aggregated result of last voting
func (h *VoteAPIHandler) RetrieveLastVoteResult(c *gin.Context) {

	// TODO call service methods here to actually fetch data

	c.JSON(http.StatusOK, VotingResult{
		Datetime: time.Now().Round(1 * time.Hour),
		Total:    82,
		Results: []CategoryScore{
			{
				Name:  "cheerful",
				Score: 43,
			},
			{
				Name:  "lyrical",
				Score: 32,
			},
			{
				Name:  "relaxed",
				Score: 7,
			},
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
			Total:    82,
			Results: []CategoryScore{
				{
					Name:  "cheerful",
					Score: 43,
				},
				{
					Name:  "lyrical",
					Score: 32,
				},
				{
					Name:  "relaxed",
					Score: 7,
				},
			},
		},
		{
			Datetime: lastHour.Add(-2 * time.Hour),
			Total:    78,
			Results: []CategoryScore{
				{
					Name:  "lyrical",
					Score: 51,
				},
				{
					Name:  "cheerful",
					Score: 20,
				},
				{
					Name:  "relaxed",
					Score: 7,
				},
			},
		},
		{
			Datetime: lastHour.Add(-3 * time.Hour),
			Total:    15,
			Results: []CategoryScore{
				{
					Name:  "lyrical",
					Score: 10,
				},
				{
					Name:  "relaxed",
					Score: 5,
				},
				{
					Name:  "cheerful",
					Score: 0,
				},
			},
		},
	})
}
