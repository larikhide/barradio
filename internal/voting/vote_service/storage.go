package vote_service

import "time"

// VoteStorage defines set of methods which should be impemented
// by storage to handle VoteService data
type VoteStorage interface {
	Close() error
	GetVoteCategories() ([]string, error)
	SaveVoteForCategory(Vote) error
	GetVotesCountForInterval(start, end time.Time) (map[string]int, error)
}

type Vote struct {
	Category  string
	CreatedAt time.Time
}

type VotingResult struct {
	Datetime time.Time
	Total    int
	Score    map[string]int
}
