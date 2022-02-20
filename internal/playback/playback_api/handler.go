package playback_api

import "github.com/larikhide/barradio/internal/playback/playback_service"

// PlaybackAPIHandler is container which stores handler depencies
type PlaybackAPIHandler struct {
	service *playback_service.PlaybackService
}

func NewPlaybackHandler(service *playback_service.PlaybackService) PlaybackAPIHandler {
	return PlaybackAPIHandler{service: service}
}
