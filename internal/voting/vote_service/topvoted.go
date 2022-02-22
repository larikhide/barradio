package vote_service

import (
	"sort"
)

type CategoryScore struct {
	Name  string
	Score int
}

// SortCategoriesByScore takes score map as arg and returns
// ordered by score desc list of categories
func SortCategoriesByScore(votes *VotingResult) []*CategoryScore {

	results := make([]*CategoryScore, 0, len(votes.Score))
	for name, score := range votes.Score {
		results = append(results, &CategoryScore{
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

	return results

}
