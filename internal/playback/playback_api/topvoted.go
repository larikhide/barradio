package playback_api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

var ErrTopvotedNotExists = errors.New("cannot get voting winner")

type Playlist struct {
	Name        string  `json:"name"`
	URL         string  `json:"url,omitempty"`
	TotalTracks int     `json:"total_tracks"`
	Duration    float64 `json:"duration_sec"`
}

func (h *PlaybackAPIHandler) getTopvotedCategory() (string, error) {
	result := ""
	votes, err := h.voting.LastVoting()
	if err != nil {
		return result, err
	}
	if votes.Total == 0 {
		return result, ErrTopvotedNotExists
	}
	maxScore := 0
	for category, score := range votes.Score {
		if score >= maxScore {
			maxScore = score
			result = category
		}
	}
	return result, nil
}

// TopvotedPlaylist return data of playlist related to category-winner
// of the last user voting
func (h *PlaybackAPIHandler) TopvotedPlaylist(c *gin.Context) {
	category, err := h.getTopvotedCategory()
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIBaseError{Message: err.Error()})
		return
	}
	playlist, err := h.stream.GetTracklistByCategory(category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIBaseError{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, Playlist{
		Name:        playlist.Name,
		URL:         playlist.ExternalURL,
		Duration:    playlist.TotalDuration,
		TotalTracks: playlist.TotalTracks,
	})
}

// TopvotedPlaylist return data of compositions related to category-winner
// of the last user voting
func (h *PlaybackAPIHandler) TopvotedPlaylistCompositions(c *gin.Context) {
	category, err := h.getTopvotedCategory()
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIBaseError{Message: err.Error()})
		return
	}
	tracks, err := h.stream.GetTracksByCategory(category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIBaseError{Message: err.Error()})
		return
	}

	result := make([]Composition, 0, len(tracks))
	for _, tr := range tracks {
		result = append(result, Composition{
			Singer:         tr.Artist,
			Name:           tr.Name,
			ImageURL:       tr.ImageURL,
			CompositionURL: tr.ExternalURL,
			Duration:       tr.Duration,
		})
	}

	c.JSON(http.StatusOK, result)
}
