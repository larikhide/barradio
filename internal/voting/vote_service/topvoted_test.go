package vote_service_test

import (
	"reflect"
	"testing"

	"github.com/larikhide/barradio/internal/voting/vote_service"
)

func TestSortCategoriesByScore(t *testing.T) {

	testTable := []struct {
		input    *vote_service.VotingResult
		expected []*vote_service.CategoryScore
	}{
		{
			input: &vote_service.VotingResult{
				Score: map[string]int{
					"cheerful": 3,
					"relaxed":  5,
				},
			},
			expected: []*vote_service.CategoryScore{
				{Name: "relaxed", Score: 5},
				{Name: "cheerful", Score: 3},
			},
		},
		{
			input: &vote_service.VotingResult{
				Score: map[string]int{
					"cheerful": 5,
					"relaxed":  5,
				},
			},
			expected: []*vote_service.CategoryScore{
				{Name: "cheerful", Score: 5},
				{Name: "relaxed", Score: 5},
			},
		},
		{
			input: &vote_service.VotingResult{
				Score: map[string]int{},
			},
			expected: []*vote_service.CategoryScore{},
		},
	}

	for i, test := range testTable {
		got := vote_service.SortCategoriesByScore(test.input)
		if !reflect.DeepEqual(got, test.expected) {
			t.Errorf("invalid sord order in test %d: got %v, expected %v", i, got, test.expected)
		}
	}

}
