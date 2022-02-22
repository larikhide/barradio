package playback_api

import (
	"github.com/larikhide/barradio/internal/playback/playback_service"
	"github.com/larikhide/barradio/internal/voting/vote_service"
)

// PlaybackAPIHandler is container which stores handler depencies
type PlaybackAPIHandler struct {
	stream *playback_service.PlaybackService
	voting *vote_service.VoteService
}

func NewPlaybackHandler(stream *playback_service.PlaybackService, voting *vote_service.VoteService) PlaybackAPIHandler {
	return PlaybackAPIHandler{
		stream: stream,
		voting: voting,
	}
}
