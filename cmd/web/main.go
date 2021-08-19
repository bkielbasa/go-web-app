package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/bkielbasa/go-web-app/cmd/web/run"
)

const tearDownTimeout = 5 * time.Second

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		oscall := <-c
		log.Printf("system call:%+v", oscall)
		cancel()
	}()

	run, teardown, err := run.App(ctx)
	if err != nil {
		log.Panic(err)
	}

	go run()

	log.Printf("server started")
	<-ctx.Done()
	log.Printf("stopping the server")

	// we give some time to close all opened connection and tidy up everything
	ctx, cancel = context.WithTimeout(context.Background(), tearDownTimeout)
	defer cancel()

	err = teardown(ctx)
	if err != nil {
		log.Printf("cannot tear down clearly: %s", err)
	}
	log.Printf("server stopped")
}
