package main

import (
	"context"
	"database/sql"
	"log"
	_ "net/http/pprof"
	"time"

	"github.com/bkielbasa/go-web-app/internal/application"
	"github.com/bkielbasa/go-web-app/internal/cmd"
	"github.com/bkielbasa/go-web-app/internal/dependency"
	"github.com/bkielbasa/go-web-app/links"
	"github.com/bkielbasa/go-web-app/links/infra"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
)

const tearDownTimeout = 5 * time.Second

type Config struct {
	HostName           string        `default:"http://localhost:8080"`
	ServerPort         int           `default:"8080"`
	TearDownTimeout    time.Duration `default:"5s"`
	PostgresConnString string
}

func main() {
	conf := Config{}
	err := envconfig.Process("", &conf)
	if err != nil {
		log.Fatal(err.Error())
	}

	// in this context, we have handling of SIGINT and SIGTERM signals
	ctx := cmd.Context()
	app := application.New(ctx, conf.ServerPort)

	db, err := sql.Open("postgres", conf.PostgresConnString)
	if err != nil {
		log.Fatal(err.Error())
	}
	// the DB connection is closed in the dependency
	app.AddDependency(dependency.NewSQL(db))

	storage := infra.NewPostgresStorage(db)

	app.AddModule(links.New(conf.HostName, storage))

	go func() {
		_ = app.Run()
	}()

	log.Printf("server started on port %d", conf.ServerPort)

	// we are waiting for the cancellation signal
	<-ctx.Done()
	log.Printf("stopping the server")

	// we give some time to close all opened connection and tidy up everything
	ctx, cancel := context.WithTimeout(context.Background(), tearDownTimeout)
	defer cancel()

	err = app.Close(ctx)
	if err != nil {
		log.Printf("cannot tear down clearly: %s", err)
	}
	log.Printf("server stopped")
}
