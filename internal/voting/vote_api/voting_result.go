package vote_api

import (
	"net/http"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/larikhide/barradio/internal/voting/vote_service"
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

func mapVotingResult(votes *vote_service.VotingResult) *VotingResult {

	results := make([]CategoryScore, 0, len(votes.Score))
	for name, score := range votes.Score {
		results = append(results, CategoryScore{
			Name:  name,
			Score: score,
		})
	}
	sort.Slice(results, func(i, j int) bool {
		if results[i].Score != results[j].Score {
			// firstly, order by score desc
			return results[i].Score > results[j].Score
		}
		// secondly, by Name asc
		return results[i].Name < results[j].Name
	})

	return &VotingResult{
		Datetime: votes.Datetime,
		Total:    votes.Total,
		Results:  results,
	}

}

// RetrieveLastVoteResult returns aggregated result of current
// unfinished voting
func (h *VoteAPIHandler) RetrieveCurrentVoteResult(c *gin.Context) {

	votes, err := h.service.CurrentVoting()
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIBaseError{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, mapVotingResult(votes))
}

// RetrieveLastVoteResult returns aggregated result of last voting
func (h *VoteAPIHandler) RetrieveLastVoteResult(c *gin.Context) {

	votes, err := h.service.LastVoting()
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIBaseError{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, mapVotingResult(votes))
}

// RetrieveVoteResultHistory returns aggregated results of voting
// for defined period
func (h *VoteAPIHandler) RetrieveVoteResultHistory(c *gin.Context) {

	searchInterval := struct {
		start time.Time
		end   time.Time
	}{
		start: time.Now().Add(-h.service.DefaultHistoryDepth),
		end:   time.Now(),
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

	results := make([]*VotingResult, 0, len(votes))
	for _, vote := range votes {
		results = append(results, mapVotingResult(vote))
	}

	c.JSON(http.StatusOK, results)
}
