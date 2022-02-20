package apiserver

import (
	"net/http"
	"time"

	// third party
	"github.com/gin-gonic/gin"

	// custom
	"github.com/larikhide/barradio/internal/playback/playback_api"
	"github.com/larikhide/barradio/internal/playback/playback_service"
	"github.com/larikhide/barradio/internal/voting/vote_api"
	"github.com/larikhide/barradio/internal/voting/vote_service"
)

var ServerTimeout = 60 * time.Second

func NewAPIServer(addr string, voteSvc *vote_service.VoteService, playbackSvc *playback_service.PlaybackService) (srv *http.Server) {

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	voteHandler := vote_api.NewVoteAPIHandler(voteSvc)
	playbackHandler := playback_api.NewPlaybackHandler(playbackSvc, voteSvc)

	// endpoints
	api := router.Group("/api")
	api.GET("/hello", voteHandler.Hello)

	voting := api.Group("/voting")
	voting.GET("/category", voteHandler.RetrieveVoteCategories)
	voting.POST("/choice", voteHandler.MakeCategoryChoice)
	voting.GET("/result/current", voteHandler.RetrieveCurrentVoteResult)
	voting.GET("/result/last", voteHandler.RetrieveLastVoteResult)
	voting.GET("/result/history", voteHandler.RetrieveVoteResultHistory)

	playback := api.Group("/playback")
	playback.GET("/compositions/:name", playbackHandler.RetrieveCompositions)
	playback.GET("/topvoted/playlist/compositions", playbackHandler.TopvotedPlaylistCompositions)
	playback.GET("/topvoted/playlist", playbackHandler.TopvotedPlaylist)

	server := http.Server{
		Addr:    addr,
		Handler: router,
	}

	return &server
}
