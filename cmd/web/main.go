package main

import (
	"context"
	"log"
	"time"

	"github.com/bkielbasa/go-web-app/internal/application"
	"github.com/bkielbasa/go-web-app/internal/cmd"
	"github.com/bkielbasa/go-web-app/links"
)

const tearDownTimeout = 5 * time.Second

func main() {
	ctx := cmd.Context()
	app := application.New(ctx)
	app.AddModule(links.New("https://linkio.io"))

	go func() {
		_ = app.Run()
	}()

	log.Printf("server started")
	<-ctx.Done()
	log.Printf("stopping the server")

	// we give some time to close all opened connection and tidy up everything
	ctx, cancel := context.WithTimeout(context.Background(), tearDownTimeout)
	defer cancel()

	err := app.Close(ctx)
	if err != nil {
		log.Printf("cannot tear down clearly: %s", err)
	}
	log.Printf("server stopped")
}
