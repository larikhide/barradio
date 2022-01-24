package vote_service

// VoteStorage defines set of methods which should be impemented
// by storage to handle VoteService data
type VoteStorage interface {
	Close() error
}
