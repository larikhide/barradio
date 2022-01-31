package apiserver

import (
	"net/http"
	"time"

	// third party
	"github.com/gin-gonic/gin"

	// custom
	"github.com/larikhide/barradio/internal/playback/playback_api"
	"github.com/larikhide/barradio/internal/voting/vote_api"
	"github.com/larikhide/barradio/internal/voting/vote_service"
)

var ServerTimeout = 60 * time.Second

func NewAPIServer(addr string, service vote_service.VoteService) (srv *http.Server) {

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	apihandler := vote_api.NewVoteAPIHandler(service)

	// endpoints
	api := router.Group("/api")
	api.GET("/hello", apihandler.Hello)

	voting := api.Group("/voting")
	voting.GET("/category", apihandler.RetrieveVoteCategories)
	voting.POST("/choice", apihandler.MakeCategoryChoice)
	voting.GET("/result/last", apihandler.RetrieveLastVoteResult)
	voting.GET("/result/history", apihandler.RetrieveVoteResultHistory)

	playback := api.Group("/playback")
	playback.GET("/category/:name", playback_api.RetrieveCompositions)

	server := http.Server{
		Addr:    addr,
		Handler: router,
	}

	return &server
}
