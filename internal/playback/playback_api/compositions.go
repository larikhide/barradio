package playback_api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/larikhide/barradio/internal/playback/playback_service"
)

type Composition struct {
	Singer         string  `json:"singer"`
	Name           string  `json:"name"`
	ImageURL       string  `json:"image_url"`
	CompositionURL string  `json:"composition_url,omitempty"`
	Duration       float64 `json:"duration_sec"`
}

func (h *PlaybackAPIHandler) RetrieveCompositions(c *gin.Context) {
	category := c.Param("name")

	tracks, err := h.service.GetTreksByCategory(category)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, playback_service.ErrUnknownCategory) {
			status = http.StatusBadRequest
		}
		c.JSON(status, APIBaseError{Message: err.Error()})
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
