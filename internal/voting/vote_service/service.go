package vote_service

// VoteService implements business logic of voting
type VoteService struct {
	store VoteStorage
}

// NewVoteService returns ready to use instance of service
func NewVoteService(storage VoteStorage) (*VoteService, error) {
	return &VoteService{store: storage}, nil
}
