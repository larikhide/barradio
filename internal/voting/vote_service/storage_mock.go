package vote_service

import "time"

var _ VoteStorage = &MockVoteStorage{}

// MockVoteStorage is in-memory storage implementation for tests
type MockVoteStorage struct {
	Categories []string
	Votes      []*Vote
}

// NewMockVoteStorage returns new ready to use instance of storage
func NewMockVoteStorage(dsn string) (*MockVoteStorage, error) {
	return &MockVoteStorage{
		Categories: []string{"cheerful", "relaxed", "lyrical"},
		Votes:      make([]*Vote, 0),
	}, nil
}

// Close disconnects from DB
func (s *MockVoteStorage) Close() error {
	return nil
}

func (s *MockVoteStorage) GetVoteCategories() ([]string, error) {
	return s.Categories, nil
}

func (s *MockVoteStorage) SaveVoteForCategory(vote Vote) error {
	s.Votes = append(s.Votes, &Vote{
		CreatedAt: vote.CreatedAt,
		Category:  vote.Category,
	})
	return nil
}

func (s *MockVoteStorage) GetVotesCountForInterval(start, end time.Time) (map[string]int, error) {
	result := make(map[string]int)

	for _, vote := range s.Votes {
		if vote.CreatedAt.After(start) && !vote.CreatedAt.After(end) {
			result[vote.Category] += 1
		}
	}

	return result, nil
}
