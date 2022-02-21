// Package main implements CLI to setup and run backend server on Heroku platform.
// Makes some specific actions for Heroku Proxy, DB integrations
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"time"

	"github.com/kelseyhightower/envconfig"

	"github.com/larikhide/barradio/internal/apiserver"
	"github.com/larikhide/barradio/internal/voting/vote_service"
	"github.com/larikhide/barradio/internal/voting/vote_storage"
)

const (
	envPrefix   = "BARRADIO"
	envDbPrefix = "DATABASE"
)

var herokuNginxSignalFile = path.Join(os.TempDir(), "app-initialized")

type DBSettings struct {
	URL           string        `default:"host=localhost port=5432 user=postgres password=postgres dbname=barradio sslmode=disable"`
	QueryTimeout  time.Duration `default:"30s"`
	MigrationsDir string        `default:"migrations"`
}

type ServerSettings struct {
	Listen              string        `default:"127.0.0.1:8123"`
	LogLevel            string        `default:"INFO"`
	CountingInterval    time.Duration `default:"30m"`
	DefaultHistoryDepth time.Duration `default:"24h"`
}

func setUp() (srv *ServerSettings, db *DBSettings) {
	srv = &ServerSettings{}
	db = &DBSettings{}
	flag.Usage = func() {
		fmt.Print("-- App server config --\n\n")
		_ = envconfig.Usage(envPrefix, srv)
		fmt.Print("\n-- Database config --\n\n")
		_ = envconfig.Usage(envDbPrefix, db)
	}
	flag.Parse()

	// always try to read env, maybe use defaults
	if err := envconfig.Process(envPrefix, srv); err != nil {
		log.Fatalln(err)
	}
	if err := envconfig.Process(envDbPrefix, db); err != nil {
		log.Fatalln(err)
	}

	return srv, db
}

func main() {
	// register signal handlers
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	// setup app
	srvSetting, dbSettings := setUp()

	// init services
	voteStore, err := vote_storage.NewPostgresVoteStorage(dbSettings.URL, dbSettings.MigrationsDir)
	if err != nil {
		log.Fatalf("cannot initialize storage: %s", err.Error())
	}
	defer voteStore.Close()

	voteService, err := vote_service.NewVoteService(voteStore, srvSetting.CountingInterval)
	if err != nil {
		log.Fatalf("cannot initialize service: %s", err.Error())
	}
	// start api server

	server := apiserver.NewAPIServer(srvSetting.Listen, *voteService, srvSetting.DefaultHistoryDepth)

	go func() {
		// usually server works behind proxy,
		// so just run plain http server
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("server fails: %s", err.Error())
		}
	}()

	// touch file to signal that nginx can start
	_, err = os.Stat(herokuNginxSignalFile)
	if os.IsNotExist(err) {
		f, err := os.Create(herokuNginxSignalFile)
		if err != nil {
			log.Fatalf("cannot create %s file for nginx integration, err: %s", herokuNginxSignalFile, err.Error())
		}
		f.Close()
	} else {
		err = os.Chtimes(herokuNginxSignalFile, time.Now().Local(), time.Now().Local())
		if err != nil {
			log.Fatalf("cannot touch %s file for nginx integration, err: %s", herokuNginxSignalFile, err.Error())
		}
	}

	// wait signal to gracefully shutdown
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("error during shutdown: %s", err.Error())
	}
}
