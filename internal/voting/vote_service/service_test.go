package vote_service_test

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/larikhide/barradio/internal/voting/vote_service"
)

func setUp(countingInterval time.Duration) (storage *vote_service.MockVoteStorage, service *vote_service.VoteService, err error) {
	storage, err = vote_service.NewMockVoteStorage("")
	if err != nil {
		return nil, nil, errors.New("cannot create storage")
	}
	service, err = vote_service.NewVoteService(storage, countingInterval)
	if err != nil {
		return nil, nil, errors.New("cannot create service")
	}
	return storage, service, nil
}

func TestGetCategories(t *testing.T) {

	countingInterval := 1 * time.Minute
	storage, service, err := setUp(countingInterval)
	if err != nil {
		t.Fatal("setup test failed:", err)
	}

	expected := storage.Categories
	got, err := service.GetCategories()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("invalid result: got %v, expected %v", got, expected)
	}

}
func TestChooseCategory(t *testing.T) {

	countingInterval := 1 * time.Minute
	storage, service, err := setUp(countingInterval)
	if err != nil {
		t.Fatal("setup test failed:", err)
	}
	storage.Votes = []*vote_service.Vote{
		{CreatedAt: time.Now(), Category: "somecategory"},
		{CreatedAt: time.Now(), Category: "somecategory"},
		{CreatedAt: time.Now(), Category: "somecategory"},
	}

	// make "votes"
	categories, _ := service.GetCategories()
	for _, c := range categories {
		err = service.ChooseCategory(c)
		if err != nil {
			t.Fatal("cannot send vote:", err)
		}
	}

	// check if votes have stored
	votesMap := make(map[string]int)
	for _, vote := range storage.Votes {
		votesMap[vote.Category] += 1
	}
	expectedScore := 1
	for _, category := range categories {
		score, ok := votesMap[category]
		if !ok {
			t.Errorf("category %s has no votes, but expected %d", category, expectedScore)
			continue
		}
		if score != expectedScore {
			t.Errorf("category %s invalid votes count: got %d, expected %d", category, score, expectedScore)
		}
	}

}

func nextIntervalStart(interval time.Duration) time.Time {
	return time.Now().Truncate(interval).Add(interval)
}

func TestActualVotingResultMethods(t *testing.T) {

	testCategory := "cheerful"
	countingInterval := 5 * time.Second

	storage, service, err := setUp(countingInterval)
	if err != nil {
		t.Fatal("setup test failed:", err)
	}
	storage.Votes = []*vote_service.Vote{
		{CreatedAt: time.Now(), Category: testCategory},
		{CreatedAt: time.Now(), Category: testCategory},
		{CreatedAt: time.Now(), Category: testCategory},
	}

	t.Run("CurrentCounters", func(t *testing.T) {
		testVotesCount := 3
		// wait until next voting intarvav will started
		// then make some votes
		nextStart := nextIntervalStart(countingInterval)
		time.Sleep(time.Until(nextStart))
		for i := 0; i < testVotesCount; i++ {
			_ = service.ChooseCategory(testCategory)
		}
		// check immediatly
		votes, err := service.CurrentVoting()
		if err != nil {
			t.Error("cannot get vote counters:", err)
		}
		// we should not check datetime, it can differ in millisec
		// but it does not matter
		if votes.Total != testVotesCount {
			t.Errorf("invalid result total: got %d, expected %d", votes.Total, testVotesCount)
		}
		score, ok := votes.Score[testCategory]
		if !ok || score != testVotesCount {
			t.Errorf("invalid result score: got %d, expected %d", score, testVotesCount)
		}
	})

	t.Run("LastCounters", func(t *testing.T) {
		testVotesCount := 3
		// wait until next voting intarvav will started
		// then make some votes
		nextStart := nextIntervalStart(countingInterval)
		time.Sleep(time.Until(nextStart))
		for i := 0; i < testVotesCount; i++ {
			_ = service.ChooseCategory(testCategory)
		}
		// wait next interval before check
		nextStart = nextIntervalStart(countingInterval)
		time.Sleep(time.Until(nextStart))
		votes, err := service.LastVoting()
		if err != nil {
			t.Error("cannot get vote counters:", err)
		}
		if votes.Datetime != nextStart {
			t.Errorf("invalid result datetime: got %s, expected %s", votes.Datetime, nextStart)
		}
		if votes.Total != testVotesCount {
			t.Errorf("invalid result total: got %d, expected %d", votes.Total, testVotesCount)
		}
		score, ok := votes.Score[testCategory]
		if !ok || score != testVotesCount {
			t.Errorf("invalid result score: got %d, expected %d", score, testVotesCount)
		}
	})

}

func fromKitchenTime(value string) time.Time {
	t, _ := time.Parse(time.Kitchen, value)
	return t
}

func TestHistoryVotingResultMethods(t *testing.T) {

	testCategory := "cheerful"
	countingInterval := 1 * time.Hour

	storage, service, err := setUp(countingInterval)
	if err != nil {
		t.Fatal("setup test failed:", err)
	}
	storage.Votes = []*vote_service.Vote{
		{CreatedAt: fromKitchenTime("03:01PM"), Category: testCategory},
		{CreatedAt: fromKitchenTime("04:01PM"), Category: testCategory},
		{CreatedAt: fromKitchenTime("04:02PM"), Category: testCategory},
		{CreatedAt: fromKitchenTime("05:03PM"), Category: testCategory},
		{CreatedAt: fromKitchenTime("05:03PM"), Category: testCategory},
		{CreatedAt: fromKitchenTime("05:03PM"), Category: testCategory},
	}

	t.Run("SelectFullInterval", func(t *testing.T) {
		start := fromKitchenTime("04:00PM")
		end := fromKitchenTime("05:00PM")
		results, err := service.HistoryOfVoting(start, end)
		if err != nil {
			t.Error("fetch history fails:", err)
		}
		expectedResults := 1
		if len(results) != expectedResults {
			t.Errorf("len of results mismatch: got %d, expected %d", len(results), expectedResults)
			return
		}
		expectedTime := fromKitchenTime("05:00PM")
		if results[0].Datetime != expectedTime {
			t.Errorf("invalid interval: got %s, expected %s", results[0].Datetime, expectedTime)
		}
		expectedScore := 2
		score, ok := results[0].Score[testCategory]
		if !ok || score != expectedScore {
			t.Errorf("invalid score: got %d, expected %d", score, expectedScore)
		}
	})

	t.Run("SelectFewIntervals", func(t *testing.T) {
		start := fromKitchenTime("03:00PM")
		end := fromKitchenTime("06:00PM")
		results, err := service.HistoryOfVoting(start, end)
		if err != nil {
			t.Error("fetch history fails:", err)
		}
		expectedResults := 3
		if len(results) != expectedResults {
			t.Errorf("len of results mismatch: got %d, expected %d", len(results), expectedResults)
		}
		expectedTimes := []time.Time{fromKitchenTime("04:00PM"), fromKitchenTime("05:00PM"), fromKitchenTime("06:00PM")}
		expectedScores := []int{1, 2, 3}
		for i, result := range results {
			if result.Datetime != expectedTimes[i] {
				t.Errorf("invalid interval: got %s, expected %s", result.Datetime, expectedTimes[i])
			}
			score, ok := result.Score[testCategory]
			if !ok || score != expectedScores[i] {
				t.Errorf("invalid score: got %d, expected %d", score, expectedScores[i])
			}
		}
	})

	t.Run("SelectEdgeOfIntervals", func(t *testing.T) {
		start := fromKitchenTime("04:30PM")
		end := fromKitchenTime("05:30PM")
		results, err := service.HistoryOfVoting(start, end)
		if err != nil {
			t.Error("fetch history fails:", err)
		}
		expectedResults := 0
		if len(results) != expectedResults {
			t.Errorf("len of results mismatch: got %d, expected %d", len(results), expectedResults)
		}
	})

	t.Run("SelectOutsideIntervals", func(t *testing.T) {
		start := fromKitchenTime("01:00PM")
		end := fromKitchenTime("02:00PM")
		results, err := service.HistoryOfVoting(start, end)
		if err != nil {
			t.Error("fetch history fails:", err)
		}
		expectedResults := 1
		if len(results) != expectedResults {
			t.Errorf("len of results mismatch: got %d, expected %d", len(results), expectedResults)
		}
	})

}
