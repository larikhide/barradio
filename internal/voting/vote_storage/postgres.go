package vote_storage

import "github.com/larikhide/barradio/internal/voting/vote_service"

var _ vote_service.VoteStorage = &PostgresVoteStorage{}

type PostgresVoteStorage struct{}

func NewPostgresVoteStorage(dsn string) (*PostgresVoteStorage, error) {
	// dummy to just start app server
	// TODO implement due https://github.com/larikhide/barradio/issues/8
	return &PostgresVoteStorage{}, nil
}

func (s *PostgresVoteStorage) Close() error {
	// dummy to clean shutdown app server
	// TODO implement due https://github.com/larikhide/barradio/issues/8
	return nil
}
