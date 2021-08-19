package run

import (
	"context"
	"log"
	"net/http"

	"github.com/bkielbasa/go-web-app/pkg/infra"
)

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
			err := srv.ListenAndServe()
			if err != nil {
				log.Print(err)
			}
		},
		func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		}, nil
}
