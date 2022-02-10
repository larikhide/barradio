package vote_service

import (
	"time"
)

// VoteService implements business logic of voting
type VoteService struct {
	store VoteStorage
	// countingInterval is duration for vote accumulation
	// before counting
	countingInterval time.Duration
}

// NewVoteService returns ready to use instance of service
func NewVoteService(storage VoteStorage, countingInterval time.Duration) (*VoteService, error) {
	service := &VoteService{
		store:            storage,
		countingInterval: countingInterval,
	}
	return service, nil
}

// GetCategories returns list of music categories which
// user can choose
func (v *VoteService) GetCategories() ([]string, error) {
	categories, err := v.store.GetVoteCategories()
	return categories, err
}

// ChooseCategory takes category name to store as user vote
func (v *VoteService) ChooseCategory(name string) error {
	vote := Vote{
		Category:  name,
		CreatedAt: time.Now(),
	}
	err := v.store.SaveVoteForCategory(vote)
	return err
}

// LastVoting returns voting results for the current
// uncomlited interval
func (v *VoteService) CurrentVoting() (*VotingResult, error) {
	intervalEnd := time.Now()
	intervalStart := time.Now().Truncate(v.countingInterval)
	return v.countVotesForInterval(intervalStart, intervalEnd)
}

// LastVoting returns voting results for the last
// completed interval
func (v *VoteService) LastVoting() (*VotingResult, error) {
	intervalEnd := time.Now().Truncate(v.countingInterval)
	intervalStart := intervalEnd.Add(-v.countingInterval)
	return v.countVotesForInterval(intervalStart, intervalEnd)
}

// HistoryOfVoting returns all results from <start> till <end>
func (v *VoteService) HistoryOfVoting(start, end time.Time) ([]*VotingResult, error) {
	intervalStart := start.Truncate(v.countingInterval)
	if intervalStart.Before(start) {
		intervalStart = intervalStart.Add(v.countingInterval)
	}
	intervalEnd := intervalStart.Add(v.countingInterval)

	result := make([]*VotingResult, 0)
	for intervalEnd.Before(end) {
		score, err := v.countVotesForInterval(intervalStart, intervalEnd)
		if err != nil {
			return []*VotingResult{}, err
		}
		result = append(result, score)
	}

	return result, nil
}

func (v *VoteService) countVotesForInterval(start, end time.Time) (*VotingResult, error) {
	votes, err := v.store.GetVotesCountForInterval(start, end)
	if err != nil {
		return &VotingResult{}, err
	}

	// fill map with existing categories
	categories, err := v.store.GetVoteCategories()
	if err != nil {
		return nil, err
	}
	score := make(map[string]int)
	total := 0
	for _, category := range categories {
		// get actual score from storage
		// if category votes did not exists - default 0 will be used
		score[category] = votes[category]
		total += votes[category]
	}

	return &VotingResult{
		Datetime: end,
		Total:    total,
		Score:    score,
	}, nil
}
