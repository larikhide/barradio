package vote_storage

import (
	"time"

	"github.com/larikhide/barradio/internal/voting/vote_service"
)

var _ vote_service.VoteStorage = &PostgresVoteStorage{}

// PostgresVoteStorage encapsulates all DB work
type PostgresVoteStorage struct {

	// TODO
	// now data stores in memory - just for tests
	// remove after implementing DB calls
	votes []struct {
		CreatedAt time.Time
		Category  string
	}
}

// NewPostgresVoteStorage returns new rady to use instance of DB connector
func NewPostgresVoteStorage(dsn string) (*PostgresVoteStorage, error) {
	// dummy to just start app server
	// TODO implement due https://github.com/larikhide/barradio/issues/8
	return &PostgresVoteStorage{}, nil
}

// Close disconnects from DB
func (s *PostgresVoteStorage) Close() error {
	// dummy to clean shutdown app server
	// TODO implement due https://github.com/larikhide/barradio/issues/8
	return nil
}

func (s *PostgresVoteStorage) GetVoteCategories() ([]string, error) {

	// TODO make actual DB request

	categories := []string{"cheerful", "relaxed", "lyrical"}
	return categories, nil
}

func (s *PostgresVoteStorage) SaveVoteForCategory(vote vote_service.Vote) error {

	// TODO store in DB instead of memory

	s.votes = append(s.votes, struct {
		CreatedAt time.Time
		Category  string
	}{
		CreatedAt: vote.CreatedAt,
		Category:  vote.Category,
	})
	return nil
}

func (s *PostgresVoteStorage) GetVotesForInterval(start, end time.Time) ([]vote_service.Vote, error) {

	// TODO fetch from DB

	result := make([]vote_service.Vote, len(s.votes))

	for _, vote := range s.votes {
		if vote.CreatedAt.After(start) && vote.CreatedAt.Before(end) {
			result = append(result, vote_service.Vote{
				Category:  vote.Category,
				CreatedAt: vote.CreatedAt,
			})
		}
	}

	return result, nil
}
