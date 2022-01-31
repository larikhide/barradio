package playback_api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Composition struct {
	Singer         string `json:"singer"`
	Name           string `json:"name"`
	ImageURL       string `json:"image_url"`
	CompositionURL string `json:"composition_url,omitempty"`
}

func RetrieveCompositions(c *gin.Context) {
	category := c.Param("name")

	_ = category // TODO fetck matching playlist content

	c.JSON(http.StatusOK, []Composition{
		{
			Singer:         "Polina Gagarina",
			Name:           "Kukushka",
			ImageURL:       "https://singer/singer.jpg",
			CompositionURL: "https://singer/singer.mp3",
		},
		{
			Singer:         "Polina Gagarina",
			Name:           "Kukushka",
			ImageURL:       "https://singer/singer.jpg",
			CompositionURL: "https://singer/singer.mp3",
		},
		{
			Singer:         "Polina Gagarina",
			Name:           "Kukushka",
			ImageURL:       "https://singer/singer.jpg",
			CompositionURL: "https://singer/singer.mp3",
		},
	})
}
