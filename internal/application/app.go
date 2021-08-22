package application

import (
	"context"
	"net/http"

	"github.com/bkielbasa/go-web-app/internal/dependency"
)

type App struct {
	httpServer *http.Server
}

func New(ctx context.Context) *App {
	mux := http.NewServeMux()
	healthy := dependency.New()
	mux.HandleFunc("/healthyz", healthy.Healthy)
	mux.HandleFunc("/readyz", healthy.Ready)

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	return &App{
		httpServer: httpServer,
	}
}

func (app *App) Run() error {
	err := app.httpServer.ListenAndServe()
	if err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (app *App) Close(ctx context.Context) error {
	return app.httpServer.Shutdown(ctx)
}
