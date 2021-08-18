package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/bkielbasa/go-web-app/pkg/infra"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		oscall := <-c
		log.Printf("system call:%+v", oscall)
		cancel()
	}()

	run, teardown, err := App(ctx)
	if err != nil {
		log.Panic(err)
	}

	go run()

	log.Printf("server started")
	<-ctx.Done()
	log.Printf("stopping the server")

	// we give some time to close all opened connection and tidy up everything
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	teardown(ctx)
	log.Printf("server stopped")
}

func App(ctx context.Context) (func(), func(context.Context) error, error) {
	mux := http.NewServeMux()
	healthy := infra.NewHealthy()
	mux.HandleFunc("/healthyz", healthy.Healthy)
	mux.HandleFunc("/readyz", healthy.Ready)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	return func() {
			srv.ListenAndServe()
		},
		func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		}, nil
}
