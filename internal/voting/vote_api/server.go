package vote_api

import (
	"net/http"
	"time"

	// third party
	"github.com/gin-gonic/gin"

	// custom
	"github.com/larikhide/barradio/internal/voting/vote_service"
)

var ServerTimeout = 60 * time.Second

func NewVoteAPIServer(addr string, service vote_service.VoteService) (srv *http.Server) {

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	apihandler := VoteAPIHandler{service: service}

	// endpoints
	api := router.Group("/api")
	api.GET("/hello", apihandler.Hello)
	api.GET("/voting/category", apihandler.RetrieveVoteCategories)
	api.POST("/voting/choice", apihandler.MakeCategoryChoice)
	api.GET("/voting/result/last", apihandler.RetrieveLastVoteResult)
	api.GET("/voting/result/history", apihandler.RetrieveVoteResultHistory)

	server := http.Server{
		Addr:    addr,
		Handler: router,
	}

	return &server
}

// VoteAPIHandler is container which stores handler depencies
type VoteAPIHandler struct {
	service vote_service.VoteService
}
