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

	votes, err := h.service.LastVoting()
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIBaseError{Message: err.Error()})
		return
	}

	results := make([]CategoryScore, len(votes.Score))
	for name, score := range votes.Score {
		results = append(results, CategoryScore{
			Name:  name,
			Score: score,
		})
	}

	c.JSON(http.StatusOK, VotingResult{
		Datetime: votes.Datetime,
		Total:    votes.Total,
		Results:  results,
	})
}

// RetrieveVoteResultHistory returns aggregated results of voting
// for defined period
func (h *VoteAPIHandler) RetrieveVoteResultHistory(c *gin.Context) {

	var searchInterval struct {
		start time.Time
		end   time.Time
	}
	err := c.ShouldBindQuery(&searchInterval)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIBaseError{Message: err.Error()})
		return
	}

	votes, err := h.service.HistoryOfVoting(searchInterval.start, searchInterval.end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIBaseError{Message: err.Error()})
		return
	}

	results := make([]VotingResult, len(votes))
	for _, vote := range votes {
		currResult := make([]CategoryScore, len(vote.Score))
		for name, score := range vote.Score {
			currResult = append(currResult, CategoryScore{
				Name:  name,
				Score: score,
			})
		}
		results = append(results, VotingResult{
			Datetime: vote.Datetime,
			Total:    vote.Total,
			Results:  currResult,
		})
	}

	c.JSON(http.StatusOK, results)
}
